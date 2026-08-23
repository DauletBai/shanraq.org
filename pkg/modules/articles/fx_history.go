package articles

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Глубокий архив курсов.
//
// Нацбанк отдаёт дневной курс примерно на пять лет назад, и это окно едет
// вперёд: то, что он показывал в прошлом году, сегодня уже недоступно. Значит
// вся история тенге до 2021 года у него не спрашивается вовсе.
//
// Её отдаёт Банк международных расчётов: месячный ряд курсов к доллару, для
// тенге — с ноября 1993 года, с месяца, когда тенге появился. Прямых курсов к
// тенге там нет, но они и не нужны: обе валюты меряются одним долларом, и
// тенге за единицу валюты X это (тенге за доллар) / (X за доллар).
//
// Месяц вместо дня здесь не потеря, а выигрыш. Тридцать два года дневных точек
// это двенадцать тысяч засечек на одной линии — каша, в которой ничего не
// разобрать. Месячных — четыре сотни, и на них видно ровно то, ради чего в
// такую глубину лезут: девальвации, переход к плавающему курсу, длинные волны.

// bisSeriesURL собирает адрес месячного ряда «конец периода» для наших валют.
// Конец периода, а не среднее: он стыкуется с дневным рядом Нацбанка — это
// курс конкретного дня, а не выдуманное среднее, которого не было ни в один
// день месяца.
func bisSeriesURL(codes []string) string {
	return "https://stats.bis.org/api/v2/data/dataflow/BIS/WS_XRU/1.0/M.." +
		strings.Join(codes, "+") + ".E?format=csv"
}

// bisBase — валюта, в которой BIS меряет все остальные.
const bisBase = "USD"

// bisAnchor — код тенге в том же ряду. Через него считается кросс-курс.
const bisAnchor = "KZT"

// FxMonth — одна месячная точка: тенге за единицу валюты.
type FxMonth struct {
	Month time.Time
	Code  string
	Value float64
}

// parseBISMonthly разбирает выгрузку BIS в кросс-курсы к тенге.
//
// Месяцы, где нет курса самого тенге, пропускаются целиком: без него кросс-курс
// не из чего считать. Доллар — особый случай: он и есть база ряда, его курс к
// тенге берётся напрямую.
func parseBISMonthly(r io.Reader) ([]FxMonth, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	head, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("заголовок выгрузки BIS: %w", err)
	}
	col := map[string]int{}
	for i, h := range head {
		col[strings.TrimSpace(h)] = i
	}
	iCur, iPer, iVal := col["CURRENCY"], col["TIME_PERIOD"], col["OBS_VALUE"]
	if iCur == 0 && iPer == 0 && iVal == 0 {
		return nil, fmt.Errorf("выгрузка BIS без нужных колонок")
	}

	// период -> валюта -> курс к доллару
	byMonth := map[string]map[string]float64{}
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("строка выгрузки BIS: %w", err)
		}
		if len(rec) <= iVal || len(rec) <= iCur || len(rec) <= iPer {
			continue
		}
		// «NaN» в выгрузке означает «за этот месяц курса нет», но ParseFloat
		// разбирает его без ошибки и возвращает NaN, а тот не меньше и не
		// больше нуля — простая проверка на положительность его пропускает.
		v, err := strconv.ParseFloat(strings.TrimSpace(rec[iVal]), 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
			continue
		}
		p := strings.TrimSpace(rec[iPer])
		if len(p) != 7 {
			continue
		}
		m := byMonth[p]
		if m == nil {
			m = map[string]float64{}
			byMonth[p] = m
		}
		m[strings.ToUpper(strings.TrimSpace(rec[iCur]))] = v
	}

	out := make([]FxMonth, 0, len(byMonth)*40)
	for p, m := range byMonth {
		kzt, ok := m[bisAnchor]
		if !ok || kzt <= 0 {
			continue
		}
		month, err := time.Parse("2006-01", p)
		if err != nil {
			continue
		}
		// Доллар — сама база ряда, его курс к тенге и есть строка тенге.
		// Брать его из отдельной строки «доллар к доллару» было бы лишним
		// звеном, которого может не оказаться.
		out = append(out, FxMonth{Month: month, Code: bisBase, Value: kzt})
		for code, per := range m {
			if code == bisAnchor || code == bisBase || per <= 0 {
				continue
			}
			out = append(out, FxMonth{Month: month, Code: code, Value: kzt / per})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].Month.Equal(out[j].Month) {
			return out[i].Month.Before(out[j].Month)
		}
		return out[i].Code < out[j].Code
	})
	return out, nil
}

// SaveMonthly кладёт месячный архив.
func (s *FxStore) SaveMonthly(ctx context.Context, rows []FxMonth) error {
	if len(rows) == 0 {
		return nil
	}
	b := &pgx.Batch{}
	for _, r := range rows {
		b.Queue(`
			INSERT INTO fx_monthly (month, code, value) VALUES ($1,$2,$3)
			ON CONFLICT (month, code) DO UPDATE SET value = EXCLUDED.value`,
			r.Month, r.Code, r.Value)
	}
	br := s.db.SendBatch(ctx, b)
	defer br.Close()
	for range rows {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

// MonthlyFresh сообщает, свежий ли месячный архив. BIS отстаёт примерно на два
// месяца, поэтому свежим считаем архив, добравшийся до позапрошлого месяца.
func (s *FxStore) MonthlyFresh(ctx context.Context) (bool, error) {
	var last *time.Time
	if err := s.db.QueryRow(ctx, `SELECT max(month) FROM fx_monthly`).Scan(&last); err != nil {
		return false, err
	}
	if last == nil {
		return false, nil
	}
	return !last.Before(monthStart(time.Now().UTC().AddDate(0, -3, 0))), nil
}

// monthStart сдвигает дату на первое число её месяца.
func monthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// fetchBISMonthly забирает и разбирает выгрузку.
func (m *Module) fetchBISMonthly(ctx context.Context) ([]FxMonth, error) {
	codes, err := m.fx.KnownCodes(ctx)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		// Дневной архив ещё пуст — какие валюты нам нужны, пока неизвестно.
		return nil, nil
	}
	if !contains(codes, bisAnchor) {
		codes = append(codes, bisAnchor)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bisSeriesURL(codes), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "shanraq.org/1.0")
	resp, err := (&http.Client{Timeout: 5 * time.Minute}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BIS ответил %d", resp.StatusCode)
	}
	return parseBISMonthly(io.LimitReader(resp.Body, 64<<20))
}
