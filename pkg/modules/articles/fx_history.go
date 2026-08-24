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

// The deep rate archive.
//
// The National Bank serves daily rates about five years back, and that window
// moves forward: what it showed last year is already out of reach today. So the
// whole history of the tenge before 2021 cannot be asked of it at all.
//
// The Bank for International Settlements serves it: a monthly series of rates
// against the dollar, for the tenge from November 1993 — the month the tenge
// appeared. There are no direct rates against the tenge there, and none are
// needed: both currencies are measured in the same dollar, and tenge per unit of
// currency X is (tenge per dollar) / (X per dollar).
//
// A month instead of a day is a gain here, not a loss. Thirty-two years of daily
// points is twelve thousand marks on one line — a mush in which nothing can be
// made out. Monthly there are four hundred, and they show exactly what people go
// this deep for: devaluations, the move to a floating rate, the long waves.

// bisSeriesURL builds the address of the "end of period" monthly series for our
// currencies. End of period, not the average: it joins onto the National Bank's
// daily series, being one specific day's rate rather than an invented average that
// held on no day of the month.
func bisSeriesURL(codes []string) string {
	return "https://stats.bis.org/api/v2/data/dataflow/BIS/WS_XRU/1.0/M.." +
		strings.Join(codes, "+") + ".E?format=csv"
}

// bisBase is the currency BIS measures all the others in.
const bisBase = "USD"

// bisAnchor is the tenge's code in the same series. The cross rate is computed
// through it.
const bisAnchor = "KZT"

// FxMonth is one monthly point: tenge per unit of currency.
type FxMonth struct {
	Month time.Time
	Code  string
	Value float64
}

// parseBISMonthly parses the BIS export into cross rates against the tenge.
//
// Months with no rate for the tenge itself are skipped whole: without it there is
// nothing to compute a cross rate from. The dollar is the special case: it is the
// series' base, and its rate against the tenge is taken directly.
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

	// period -> currency -> rate against the dollar
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
		// "NaN" in the export means "no rate for this month", but ParseFloat
		// reads it without error and returns NaN, which is neither less than nor
		// greater than zero — a plain positivity check lets it through.
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
		// The dollar is the series' own base, so its rate against the tenge is
		// the tenge row. Taking it from a separate "dollar against dollar" row
		// would add a link that may not be there.
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

// SaveMonthly writes the monthly archive.
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

// MonthlyFresh reports whether the monthly archive is current. BIS lags by about
// two months, so an archive that has reached the month before last counts as
// fresh.
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

// monthStart moves a date to the first of its month.
func monthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// fetchBISMonthly fetches and parses the export.
func (m *Module) fetchBISMonthly(ctx context.Context) ([]FxMonth, error) {
	codes, err := m.fx.KnownCodes(ctx)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		// The daily archive is still empty, so which currencies we need is not
		// known yet.
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
