package articles

import (
	"context"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Exchange rates: the store, the fetch from the National Bank, and the series the
// rates page draws.
//
// The source is the same one the header ticker uses, but for an arbitrary date:
// nationalbank.kz/rss/get_rates.cfm?fdate=DD.MM.YYYY. The answer for 2020 and
// earlier is empty, so somebody else's archive is about five years deep and that
// window moves forward. Our job is to keep up with it.

// fxRatesURL returns the address of the official rates for one day.
func fxRatesURL(d time.Time) string {
	return "https://nationalbank.kz/rss/get_rates.cfm?fdate=" + d.Format("02.01.2006")
}

// FxRate is one rate for one day.
type FxRate struct {
	Day   time.Time
	Code  string
	Value float64
	Quant int
	Name  string
}

// parseFxRates parses the bank's answer. An empty list is not an error: for
// weekends and for dates outside the window the bank answers with a document
// holding no currencies at all.
func parseFxRates(body []byte, day time.Time) ([]FxRate, error) {
	var doc struct {
		Items []struct {
			Fullname string `xml:"fullname"`
			Title    string `xml:"title"`
			Desc     string `xml:"description"`
			Quant    string `xml:"quant"`
		} `xml:"item"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("разбор курсов: %w", err)
	}
	out := make([]FxRate, 0, len(doc.Items))
	for _, it := range doc.Items {
		code := strings.ToUpper(strings.TrimSpace(it.Title))
		if len(code) != 3 {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(it.Desc), 64)
		if err != nil || v <= 0 {
			continue
		}
		q, err := strconv.Atoi(strings.TrimSpace(it.Quant))
		if err != nil || q <= 0 {
			q = 1
		}
		out = append(out, FxRate{
			Day: day, Code: code, Value: v, Quant: q,
			Name: strings.TrimSpace(it.Fullname),
		})
	}
	return out, nil
}

// FxStore keeps the rates and serves the series.
type FxStore struct{ db *pgxpool.Pool }

// NewFxStore builds the store over the shared pool.
func NewFxStore(db *pgxpool.Pool) *FxStore { return &FxStore{db: db} }

// Save writes a whole day. Calling it again for the same date updates the values:
// the bank occasionally restates a rate after the fact, and arguing with the
// source is not our business.
func (s *FxStore) Save(ctx context.Context, rates []FxRate) error {
	if len(rates) == 0 {
		return nil
	}
	b := make([][]any, 0, len(rates))
	for _, r := range rates {
		b = append(b, []any{r.Day, r.Code, r.Value, r.Quant, r.Name})
	}
	for _, row := range b {
		if _, err := s.db.Exec(ctx, `
			INSERT INTO fx_rates (day, code, value, quant, name)
			VALUES ($1,$2,$3,$4,$5)
			ON CONFLICT (day, code) DO UPDATE
			  SET value = EXCLUDED.value, quant = EXCLUDED.quant,
			      name = CASE WHEN EXCLUDED.name <> '' THEN EXCLUDED.name ELSE fx_rates.name END`,
			row...); err != nil {
			return fmt.Errorf("сохранение курса: %w", err)
		}
	}
	return nil
}

// Has reports whether we hold any rate at all for this day. The backfill asks
// before it goes to the network.
func (s *FxStore) Has(ctx context.Context, day time.Time) (bool, error) {
	var n int
	err := s.db.QueryRow(ctx, `SELECT count(*) FROM fx_rates WHERE day = $1`, day).Scan(&n)
	return n > 0, err
}

// FxPoint is one point of a series.
type FxPoint struct {
	Day   time.Time
	Value float64
}

// Series serves a currency's daily series from `from` to today, ascending by
// date.
//
// The value is normalised to one unit of the currency. The bank publishes yen in
// hundreds and Vietnamese dong in thousands, and that multiple has changed over
// the years; it is stored as it came, and only the normalised values can be
// compared with one another.
func (s *FxStore) Series(ctx context.Context, code string, from time.Time) ([]FxPoint, error) {
	rows, err := s.db.Query(ctx, `
		SELECT day, value / quant FROM fx_rates
		WHERE code = $1 AND day >= $2 AND quant > 0
		ORDER BY day`, strings.ToUpper(code), from)
	if err != nil {
		return nil, fmt.Errorf("ряд курса: %w", err)
	}
	defer rows.Close()
	out := []FxPoint{}
	for rows.Next() {
		var p FxPoint
		if err := rows.Scan(&p.Day, &p.Value); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// SeriesMonthly serves a monthly series to the full available depth: the BIS
// archive from November 1993, and over it the last day of each month from our own
// daily archive for the months BIS has not reached yet. The join is honest: both
// halves are one specific day's rate, not a monthly average.
//
// A point is dated by the day whose rate it shows: the last day of the month, and
// for the current, unclosed month, today. Dating it the first of the month would
// be a falsehood on the axis label.
func (s *FxStore) SeriesMonthly(ctx context.Context, code string, from time.Time) ([]FxPoint, error) {
	// The alias is deliberately not `day`: in ORDER BY an output column's name
	// shadows the table column of the same name, and "the last day of the month"
	// turned into a sort by the month's start — that is, a random row out of
	// thirty.
	rows, err := s.db.Query(ctx, `
		WITH deep AS (
			SELECT (month + INTERVAL '1 month - 1 day')::date AS d, value AS v
			  FROM fx_monthly WHERE code = $1
		), edge AS (
			SELECT COALESCE(max(month), DATE '1900-01-01') AS m FROM fx_monthly WHERE code = $1
		), recent AS (
			SELECT DISTINCT ON (date_trunc('month', r.day))
			       r.day AS d, r.value / r.quant AS v
			  FROM fx_rates r
			 WHERE r.code = $1 AND r.quant > 0
			   AND date_trunc('month', r.day) > (SELECT m FROM edge)
			 ORDER BY date_trunc('month', r.day), r.day DESC
		)
		SELECT d, v FROM deep
		UNION ALL
		SELECT d, v FROM recent
		ORDER BY 1`, strings.ToUpper(code))
	if err != nil {
		return nil, fmt.Errorf("месячный ряд курса: %w", err)
	}
	defer rows.Close()
	out := []FxPoint{}
	for rows.Next() {
		var p FxPoint
		if err := rows.Scan(&p.Day, &p.Value); err != nil {
			return nil, err
		}
		if !from.IsZero() && p.Day.Before(from) {
			continue
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// FxCurrency is one currency in the picker.
type FxCurrency struct {
	Code  string
	Name  string
	Quant int
	Days  int
	Since time.Time
}

// Label is the name for the picker. The bank writes them in capitals, and a line
// of caps in the middle of ordinary text reads as shouting.
func (c FxCurrency) Label() string { return fxTidyName(c.Name) }

// Currencies lists the currencies to choose from: code, name, the familiar
// multiple and the depth of history. The order runs from the most asked-for to the
// rest and then alphabetically: nobody should have to scroll past the dollar and
// the euro.
func (s *FxStore) Currencies(ctx context.Context) ([]FxCurrency, error) {
	rows, err := s.db.Query(ctx, `
		SELECT r.code,
		       (SELECT name FROM fx_rates n
		         WHERE n.code = r.code AND n.name <> '' ORDER BY n.day DESC LIMIT 1),
		       (SELECT quant FROM fx_rates q WHERE q.code = r.code ORDER BY q.day DESC LIMIT 1),
		       count(*),
		       LEAST(min(r.day), COALESCE((SELECT min(month) FROM fx_monthly m WHERE m.code = r.code), min(r.day)))
		  FROM fx_rates r
		 GROUP BY r.code
		 ORDER BY r.code`)
	if err != nil {
		return nil, fmt.Errorf("список валют: %w", err)
	}
	defer rows.Close()
	out := []FxCurrency{}
	for rows.Next() {
		var c FxCurrency
		var name *string
		if err := rows.Scan(&c.Code, &name, &c.Quant, &c.Days, &c.Since); err != nil {
			return nil, err
		}
		if name != nil {
			c.Name = *name
		}
		if c.Quant <= 0 {
			c.Quant = 1
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(out, func(i, j int) bool {
		pi, pj := fxRank(out[i].Code), fxRank(out[j].Code)
		if pi != pj {
			return pi < pj
		}
		return out[i].Code < out[j].Code
	})
	return out, nil
}

// fxPopular are the currencies asked for most often, in order of demand.
var fxPopular = []string{"USD", "EUR", "RUB", "CNY", "GBP", "TRY", "KGS", "UZS", "AED"}

// fxRank puts the currencies in demand at the head of the list.
func fxRank(code string) int {
	for i, c := range fxPopular {
		if c == code {
			return i
		}
	}
	return len(fxPopular)
}

// Earliest is the first day we hold anything at all for.
func (s *FxStore) Earliest(ctx context.Context) (time.Time, error) {
	var d *time.Time
	if err := s.db.QueryRow(ctx, `SELECT min(day) FROM fx_rates`).Scan(&d); err != nil {
		return time.Time{}, err
	}
	if d == nil {
		return time.Time{}, nil
	}
	return *d, nil
}

// MarkProbed records that we have already asked the bank about this day.
func (s *FxStore) MarkProbed(ctx context.Context, day time.Time, found int) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO fx_probed (day, found, probed_at) VALUES ($1,$2,now())
		ON CONFLICT (day) DO UPDATE SET found = EXCLUDED.found, probed_at = now()`, day, found)
	return err
}

// NextToProbe returns the nearest unprobed day, walking back from today but no
// deeper than floor. An empty result means there is nothing left to backfill.
func (s *FxStore) NextToProbe(ctx context.Context, floor time.Time) (time.Time, bool, error) {
	// max, not ORDER BY … LIMIT 1: with nothing left to probe, a row query would
	// return "no rows" — that is, an error — and the backfill would complain to
	// the log every hour about having finished its work. An aggregate honestly
	// returns one row holding NULL.
	var d *time.Time
	err := s.db.QueryRow(ctx, `
		SELECT max(g)::date FROM generate_series($1::date, CURRENT_DATE, interval '1 day') g
		WHERE g::date NOT IN (SELECT day FROM fx_probed)`, floor).Scan(&d)
	if err != nil || d == nil {
		return time.Time{}, false, err
	}
	return *d, true, nil
}

// EmptyRunBelow counts how many consecutive probed days below `day` yielded no
// rate at all. A long run means we have reached the edge of somebody else's
// window: past it the bank is silent, and knocking there is pointless.
func (s *FxStore) EmptyRunBelow(ctx context.Context, day time.Time) (int, error) {
	var n int
	err := s.db.QueryRow(ctx, `
		WITH ordered AS (
			SELECT day, found FROM fx_probed WHERE day <= $1 ORDER BY day DESC
		), marked AS (
			SELECT found, sum(CASE WHEN found > 0 THEN 1 ELSE 0 END)
			       OVER (ORDER BY day DESC ROWS UNBOUNDED PRECEDING) AS grp
			FROM ordered
		)
		SELECT count(*) FROM marked WHERE grp = 0`, day).Scan(&n)
	return n, err
}

// KnownCodes lists the currencies that have appeared in the daily archive.
func (s *FxStore) KnownCodes(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT DISTINCT code FROM fx_rates ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
