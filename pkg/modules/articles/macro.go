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

// The series the exchange rate and inflation are made of.
//
// The tenge rate is not weather; it does not "happen". On one side stands the
// quantity of tenge, on the other the quantity of dollars the country holds. The
// National Bank publishes both itself, monthly and back to 1994; all this part
// does is fetch them and lay them side by side so that one can be seen against
// the other.
//
// Machine access amounts to exactly one thing: the "Export" button on the
// statistics page, which is an Excel workbook at the same address with an /excel
// suffix. There is no feed and no API.

const (
	// MacroM3 is broad money, millions of tenge.
	MacroM3 = "m3"
	// MacroM0 is cash in circulation, millions of tenge.
	MacroM0 = "m0"
	// MacroBase is the monetary base, millions of tenge.
	MacroBase = "base_money"
	// MacroReserves is gross international reserves, millions of dollars.
	MacroReserves = "reserves"
	// MacroFund is the National Fund's currency assets, millions of dollars.
	MacroFund = "nat_fund"
	// MacroCPI is annual inflation, percent.
	MacroCPI = "cpi"
)

const (
	nbkMoneyURL = "https://nationalbank.kz/ru/monetarybase/" +
		"denezhnaya-baza-i-agregaty-shirokoy-denezhnoy-massy/excel"
	nbkReservesURL = "https://nationalbank.kz/ru/international-reserve-and-asset/" +
		"mezhdunarodnye-rezervy-i-aktivy-nacionalnogo-fonda-rk/excel"
	// The National Bank does not serve inflation by machine: for it these are
	// documents, not a series. The World Bank serves it as an annual series from
	// 1994 — the tenge's first full year.
	worldBankCPIURL = "https://api.worldbank.org/v2/country/KAZ/indicator/" +
		"FP.CPI.TOTL.ZG?format=json&per_page=200"
)

// MacroPoint is one point of a series.
type MacroPoint struct {
	Period time.Time
	Value  float64
}

// MacroStore keeps the series and serves them.
type MacroStore struct{ db *pgxpool.Pool }

// NewMacroStore builds the store over the shared pool.
func NewMacroStore(db *pgxpool.Pool) *MacroStore { return &MacroStore{db: db} }

// Save writes a whole series, overwriting revised values: statistics are
// restated after the fact, and arguing with the source is not our business.
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

// Series returns a series ascending by period.
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

// Latest is the most recent period that has a value.
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

// parseNBKMonth reads a National Bank month label: " 07.26" → July 2026.
//
// The year is two digits, so the century has to be guessed. The boundary is drawn
// at the tenge's appearance: anything under 90 is the 2000s, the rest the 1990s.
// The series starts in 1994, so there is nowhere to go wrong.
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
	// A one-digit year is a two-digit year whose trailing zero was lost. The
	// bank writes "10.10" for October 2010, the workbook stores that cell as
	// the number 10.1, and what reaches us is "10.1" — which read literally is
	// October 2001. The same happens to every year ending in a zero: 2000,
	// 2010, 2020. Leading zeros survive (May 2003 stays "5.03"), because a
	// number only drops zeros after its last significant digit, so a single
	// digit can only ever be the tens.
	//
	// Left unhandled this was not a gap but a forgery: the 2010 and 2020 rows
	// overwrote 2001 and 2002 in place, and the reserves chart showed the
	// country holding twelve times its actual currency in the autumn of 2001.
	if len(strings.TrimSpace(yy)) == 1 && y >= 0 && y <= 9 {
		y *= 10
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

// parseNBKMoney parses the monetary aggregates workbook. Months run across the
// columns and indicators down the rows, so the columns are dated first and then
// each row of interest is read along them.
func parseNBKMoney(data []byte) (map[string][]MacroPoint, error) {
	sheet, err := readXLSX(data)
	if err != nil {
		return nil, err
	}
	cols, headerRow := macroMonthColumns(sheet)
	if len(cols) < 12 {
		return nil, fmt.Errorf("в книге денежных агрегатов размечено %d месяцев", len(cols))
	}

	// Rows of interest are recognised by the start of their name: the Bank
	// numbers its indicators, and the number together with the name holds for
	// years, while the order of the rows does not.
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

// macroCol is a workbook column and the month it stands for.
type macroCol struct {
	col   string
	month time.Time
}

// macroMonthColumns finds the header row of months and dates the columns by it.
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

// parseNBKReserves parses the international reserves workbook.
//
// Here months run down the rows and indicators across the columns, and each
// indicator has three columns: the volume and two percentage changes. Columns are
// found by the indicator's name rather than by counting: counting positions would
// record net reserves as National Fund assets the moment the source adds a
// column.
func parseNBKReserves(data []byte) (map[string][]MacroPoint, error) {
	sheet, err := readXLSX(data)
	if err != nil {
		return nil, err
	}

	// The indicator's name sits one row above the word "объем", same column.
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

// parseWorldBankCPI parses annual inflation out of the World Bank's answer.
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

// macroFetch fetches a whole document.
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

// Refreshing the series.
//
// The National Bank publishes monetary aggregates and reserves monthly, the World
// Bank inflation yearly. There is no point going for them more than once a day,
// and no way to go less often: the day of publication is not known in advance.

// macroEvery is how often to check whether new data has appeared.
const macroEvery = 24 * time.Hour

// RunMacro keeps the series fresh until the context is cancelled.
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

// refreshMacro fetches every series. The sources are independent: one going
// silent must not leave the page without the others.
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

	// The National Bank's rate and today's indicator panel.
	failed = append(failed, m.refreshRates(ctx, time.Now().UTC().Truncate(24*time.Hour))...)

	if len(failed) > 0 {
		return fmt.Errorf("%s", strings.Join(failed, "; "))
	}
	m.rt.Logger.Info("макроряды обновлены")
	return nil
}
