package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Ряды, из которых складываются курс и инфляция.
//
// Курс тенге — не погода, он не «случается». С одной стороны стоит количество
// тенге, с другой — количество долларов, которыми страна располагает. Обе
// величины Нацбанк публикует сам, помесячно и с 1994 года; всё, что делает эта
// часть, — забирает их и кладёт рядом, чтобы одно можно было увидеть на фоне
// другого.
//
// Машинный доступ у Нацбанка ровно один: кнопка «Экспорт» на странице
// статистики, то есть книга Excel по адресу с суффиксом /excel. Ни ленты, ни
// API у него нет.

const (
	// MacroM3 — широкая денежная масса, млн тенге.
	MacroM3 = "m3"
	// MacroM0 — наличные деньги в обращении, млн тенге.
	MacroM0 = "m0"
	// MacroBase — денежная база, млн тенге.
	MacroBase = "base_money"
	// MacroReserves — валовые международные резервы, млн долларов.
	MacroReserves = "reserves"
	// MacroFund — валютные активы Национального фонда, млн долларов.
	MacroFund = "nat_fund"
	// MacroCPI — годовая инфляция, проценты.
	MacroCPI = "cpi"
)

const (
	nbkMoneyURL = "https://nationalbank.kz/ru/monetarybase/" +
		"denezhnaya-baza-i-agregaty-shirokoy-denezhnoy-massy/excel"
	nbkReservesURL = "https://nationalbank.kz/ru/international-reserve-and-asset/" +
		"mezhdunarodnye-rezervy-i-aktivy-nacionalnogo-fonda-rk/excel"
	// Инфляцию Нацбанк машинно не отдаёт: у него это документы, а не ряд.
	// Всемирный банк отдаёт её годовым рядом с 1994 года — с первого полного
	// года тенге.
	worldBankCPIURL = "https://api.worldbank.org/v2/country/KAZ/indicator/" +
		"FP.CPI.TOTL.ZG?format=json&per_page=200"
)

// MacroPoint — точка ряда.
type MacroPoint struct {
	Period time.Time
	Value  float64
}

// MacroStore хранит ряды и отдаёт их.
type MacroStore struct{ db *pgxpool.Pool }

// NewMacroStore строит хранилище над общим пулом.
func NewMacroStore(db *pgxpool.Pool) *MacroStore { return &MacroStore{db: db} }

// Save кладёт ряд целиком, перезаписывая пересчитанные значения: статистику
// уточняют задним числом, и спорить с источником не наше дело.
func (s *MacroStore) Save(ctx context.Context, code string, pts []MacroPoint) error {
	if len(pts) == 0 {
		return nil
	}
	b := &pgx.Batch{}
	for _, p := range pts {
		b.Queue(`
			INSERT INTO macro_series (code, period, value) VALUES ($1,$2,$3)
			ON CONFLICT (code, period) DO UPDATE SET value = EXCLUDED.value`,
			code, p.Period, p.Value)
	}
	br := s.db.SendBatch(ctx, b)
	defer br.Close()
	for range pts {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("сохранение ряда %s: %w", code, err)
		}
	}
	return nil
}

// Series отдаёт ряд по возрастанию периода.
func (s *MacroStore) Series(ctx context.Context, code string) ([]MacroPoint, error) {
	rows, err := s.db.Query(ctx,
		`SELECT period, value FROM macro_series WHERE code = $1 ORDER BY period`, code)
	if err != nil {
		return nil, fmt.Errorf("ряд %s: %w", code, err)
	}
	defer rows.Close()
	out := []MacroPoint{}
	for rows.Next() {
		var p MacroPoint
		if err := rows.Scan(&p.Period, &p.Value); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Latest — последний период, за который есть значение.
func (s *MacroStore) Latest(ctx context.Context, code string) (time.Time, error) {
	var d *time.Time
	if err := s.db.QueryRow(ctx,
		`SELECT max(period) FROM macro_series WHERE code = $1`, code).Scan(&d); err != nil {
		return time.Time{}, err
	}
	if d == nil {
		return time.Time{}, nil
	}
	return *d, nil
}

// parseNBKMonth разбирает подпись месяца Нацбанка: « 07.26» → июль 2026.
//
// Год двузначный, и век приходится угадывать. Граница проведена по появлению
// тенге: всё, что меньше 90, — двухтысячные, остальное — девяностые. Ряд
// начинается в 1994 году, так что ошибиться негде.
func parseNBKMonth(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	mm, yy, ok := strings.Cut(s, ".")
	if !ok {
		return time.Time{}, false
	}
	m, err1 := strconv.Atoi(strings.TrimSpace(mm))
	y, err2 := strconv.Atoi(strings.TrimSpace(yy))
	if err1 != nil || err2 != nil || m < 1 || m > 12 {
		return time.Time{}, false
	}
	if y < 90 {
		y += 2000
	} else if y < 100 {
		y += 1900
	}
	if y < 1993 || y > 2100 {
		return time.Time{}, false
	}
	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC), true
}

// parseNBKMoney разбирает книгу денежных агрегатов. Месяцы идут по колонкам,
// показатели — по строкам, поэтому колонки сперва размечаются датами, а потом
// каждая нужная строка читается вдоль них.
func parseNBKMoney(data []byte) (map[string][]MacroPoint, error) {
	sheet, err := readXLSX(data)
	if err != nil {
		return nil, err
	}
	cols, headerRow := macroMonthColumns(sheet)
	if len(cols) < 12 {
		return nil, fmt.Errorf("в книге денежных агрегатов размечено %d месяцев", len(cols))
	}

	// Нужные строки узнаём по началу названия: Нацбанк нумерует показатели, и
	// номер вместе с названием держится годами, а порядок строк — нет.
	want := map[string]string{
		"1. денежная база": MacroBase,
		"2. m0":            MacroM0,
		"5. m3":            MacroM3,
	}
	out := map[string][]MacroPoint{}
	for i, row := range sheet {
		if i <= headerRow {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(row["A"]))
		for prefix, code := range want {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			pts := make([]MacroPoint, 0, len(cols))
			for _, c := range cols {
				v, err := strconv.ParseFloat(strings.TrimSpace(row[c.col]), 64)
				if err != nil || v <= 0 {
					continue
				}
				pts = append(pts, MacroPoint{Period: c.month, Value: v})
			}
			if len(pts) > len(out[code]) {
				out[code] = pts
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("в книге денежных агрегатов не нашлось ни одного нужного показателя")
	}
	return out, nil
}

// macroCol — колонка книги и месяц, который она обозначает.
type macroCol struct {
	col   string
	month time.Time
}

// macroMonthColumns находит строку-шапку с месяцами и размечает по ней колонки.
func macroMonthColumns(sheet xlsxSheet) ([]macroCol, int) {
	best, bestRow := []macroCol{}, -1
	for i, row := range sheet {
		if i > 8 {
			break
		}
		cols := []macroCol{}
		for col, v := range row {
			if col == "A" {
				continue
			}
			if m, ok := parseNBKMonth(v); ok {
				cols = append(cols, macroCol{col: col, month: m})
			}
		}
		if len(cols) > len(best) {
			best, bestRow = cols, i
		}
	}
	sort.Slice(best, func(i, j int) bool {
		return xlsxColNum(best[i].col) < xlsxColNum(best[j].col)
	})
	return best, bestRow
}

// parseNBKReserves разбирает книгу международных резервов.
//
// Здесь месяцы идут по строкам, а показатели по колонкам, и у каждого
// показателя три колонки: объём и два изменения в процентах. Колонки ищутся по
// названию показателя, а не по счёту: считать позиции значило бы записать
// чистые резервы как активы Нацфонда, стоит источнику добавить один столбец.
func parseNBKReserves(data []byte) (map[string][]MacroPoint, error) {
	sheet, err := readXLSX(data)
	if err != nil {
		return nil, err
	}

	// Название показателя стоит строкой выше слова «объем» и в той же колонке.
	want := map[string]string{
		"валовые международные резервы":       MacroReserves,
		"валютные активы национального фонда": MacroFund,
	}
	col := map[string]string{} // код ряда → колонка объёма
	for i := 1; i < len(sheet) && i < 8; i++ {
		for c, v := range sheet[i] {
			if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "объем") {
				continue
			}
			title := strings.ToLower(strings.TrimSpace(sheet[i-1][c]))
			for prefix, code := range want {
				if strings.HasPrefix(title, prefix) {
					col[code] = c
				}
			}
		}
		if len(col) == len(want) {
			break
		}
	}
	if len(col) != len(want) {
		return nil, fmt.Errorf("в книге резервов найдено %d показателей из %d", len(col), len(want))
	}

	out := map[string][]MacroPoint{}
	for _, row := range sheet {
		month, ok := parseNBKMonth(row["A"])
		if !ok {
			continue
		}
		for code, c := range col {
			v, err := strconv.ParseFloat(strings.TrimSpace(row[c]), 64)
			if err != nil || v <= 0 {
				continue
			}
			out[code] = append(out[code], MacroPoint{Period: month, Value: v})
		}
	}
	if len(out[MacroReserves]) < 12 {
		return nil, fmt.Errorf("в книге резервов прочиталось %d месяцев", len(out[MacroReserves]))
	}
	return out, nil
}

// parseWorldBankCPI разбирает годовую инфляцию из ответа Всемирного банка.
func parseWorldBankCPI(body []byte) ([]MacroPoint, error) {
	var doc []json.RawMessage
	if err := json.Unmarshal(body, &doc); err != nil || len(doc) < 2 {
		return nil, fmt.Errorf("ответ Всемирного банка не разобран")
	}
	var rows []struct {
		Date  string   `json:"date"`
		Value *float64 `json:"value"`
	}
	if err := json.Unmarshal(doc[1], &rows); err != nil {
		return nil, fmt.Errorf("ряд инфляции: %w", err)
	}
	out := make([]MacroPoint, 0, len(rows))
	for _, r := range rows {
		if r.Value == nil {
			continue
		}
		y, err := strconv.Atoi(r.Date)
		if err != nil || y < 1990 {
			continue
		}
		out = append(out, MacroPoint{
			Period: time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC),
			Value:  *r.Value,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Period.Before(out[j].Period) })
	return out, nil
}

// macroFetch забирает документ целиком.
func macroFetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "shanraq.org/1.0")
	resp, err := (&http.Client{Timeout: 3 * time.Minute}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("источник ответил %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 64<<20))
}

// Обновление рядов.
//
// Нацбанк публикует денежные агрегаты и резервы раз в месяц, Всемирный банк
// инфляцию — раз в год. Ходить за ними чаще, чем раз в сутки, незачем, а реже
// нельзя: день выхода заранее неизвестен.

// macroEvery — как часто проверять, не вышли ли новые данные.
const macroEvery = 24 * time.Hour

// RunMacro поддерживает ряды в свежем виде до отмены контекста.
func (m *Module) RunMacro(ctx context.Context) {
	if m.macro == nil {
		return
	}
	for {
		if err := m.refreshMacro(ctx); err != nil {
			m.rt.Logger.Warn("макроряды", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(macroEvery):
		}
	}
}

// refreshMacro забирает все ряды. Источники независимы: молчание одного не
// должно оставлять страницу без остальных.
func (m *Module) refreshMacro(ctx context.Context) error {
	var failed []string

	if data, err := macroFetch(ctx, nbkMoneyURL); err != nil {
		failed = append(failed, "денежные агрегаты: "+err.Error())
	} else if series, err := parseNBKMoney(data); err != nil {
		failed = append(failed, "денежные агрегаты: "+err.Error())
	} else {
		for code, pts := range series {
			if err := m.macro.Save(ctx, code, pts); err != nil {
				failed = append(failed, err.Error())
			}
		}
	}

	if data, err := macroFetch(ctx, nbkReservesURL); err != nil {
		failed = append(failed, "резервы: "+err.Error())
	} else if series, err := parseNBKReserves(data); err != nil {
		failed = append(failed, "резервы: "+err.Error())
	} else {
		for code, pts := range series {
			if err := m.macro.Save(ctx, code, pts); err != nil {
				failed = append(failed, err.Error())
			}
		}
	}

	if data, err := macroFetch(ctx, worldBankCPIURL); err != nil {
		failed = append(failed, "инфляция: "+err.Error())
	} else if pts, err := parseWorldBankCPI(data); err != nil {
		failed = append(failed, "инфляция: "+err.Error())
	} else if err := m.macro.Save(ctx, MacroCPI, pts); err != nil {
		failed = append(failed, err.Error())
	}

	if len(failed) > 0 {
		return fmt.Errorf("%s", strings.Join(failed, "; "))
	}
	m.rt.Logger.Info("макроряды обновлены")
	return nil
}
