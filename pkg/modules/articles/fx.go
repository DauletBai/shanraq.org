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

// Курсы валют: хранилище, загрузка из Нацбанка и ряды для страницы разбора.
//
// Источник — тот же, что и у строки курсов в шапке, но за произвольную дату:
// nationalbank.kz/rss/get_rates.cfm?fdate=ДД.ММ.ГГГГ. Ответ за 2020 год и
// раньше пустой, поэтому глубина чужого архива — около пяти лет, и она едет
// вперёд. Наша задача — успевать забирать.

// fxRatesURL возвращает адрес официальных курсов на конкретный день.
func fxRatesURL(d time.Time) string {
	return "https://nationalbank.kz/rss/get_rates.cfm?fdate=" + d.Format("02.01.2006")
}

// FxRate — один курс за один день.
type FxRate struct {
	Day   time.Time
	Code  string
	Value float64
	Quant int
	Name  string
}

// parseFxRates разбирает ответ банка. Пустой список — это не ошибка: за
// выходные и за даты вне окна банк отвечает документом без единой валюты.
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

// FxStore хранит курсы и отдаёт ряды.
type FxStore struct{ db *pgxpool.Pool }

// NewFxStore строит хранилище над общим пулом.
func NewFxStore(db *pgxpool.Pool) *FxStore { return &FxStore{db: db} }

// Save кладёт день целиком. Повторный вызов за ту же дату обновляет значения:
// банк изредка уточняет курс задним числом, и спорить с источником не наше
// дело.
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

// Has сообщает, есть ли у нас хоть один курс за этот день. Догрузка спрашивает
// это перед тем, как идти в сеть.
func (s *FxStore) Has(ctx context.Context, day time.Time) (bool, error) {
	var n int
	err := s.db.QueryRow(ctx, `SELECT count(*) FROM fx_rates WHERE day = $1`, day).Scan(&n)
	return n > 0, err
}

// FxPoint — точка ряда.
type FxPoint struct {
	Day   time.Time
	Value float64
}

// Series отдаёт дневной ряд по валюте от from до сегодня, по возрастанию даты.
//
// Значение приводится к одной единице валюты. Банк публикует иены сотнями, а
// вьетнамские донги тысячами, и эта кратность за годы менялась; хранится она
// как есть, а сравнивать между собой можно только приведённое.
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

// SeriesMonthly отдаёт месячный ряд на всю доступную глубину: архив BIS с
// ноября 1993 года, а поверх него — последний день каждого месяца из нашего
// дневного архива за те месяцы, до которых BIS ещё не дошёл. Стык честный:
// обе половины это курс конкретного дня, а не среднее по месяцу.
//
// Точка датируется тем днём, курс которого показывает: последним днём месяца,
// а для текущего, ещё не закрытого месяца — сегодняшним. Датировать её первым
// числом было бы неправдой на подписи оси.
func (s *FxStore) SeriesMonthly(ctx context.Context, code string, from time.Time) ([]FxPoint, error) {
	// Псевдоним намеренно не day: в ORDER BY имя выходной колонки перекрывает
	// одноимённую колонку таблицы, и «последний день месяца» превращался в
	// сортировку по началу месяца, то есть в случайную строку из тридцати.
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

// FxCurrency — валюта в списке выбора.
type FxCurrency struct {
	Code  string
	Name  string
	Quant int
	Days  int
	Since time.Time
}

// Label — название для списка. Банк пишет их прописными, а строка из капса
// посреди обычного текста читается как крик.
func (c FxCurrency) Label() string { return fxTidyName(c.Name) }

// Currencies перечисляет валюты для выбора: код, название, привычную
// кратность и глубину истории. Порядок — от самых спрашиваемых к остальным,
// дальше по алфавиту: доллар с евро листать никто не должен.
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

// fxPopular — валюты, которые спрашивают чаще всего, в порядке спроса.
var fxPopular = []string{"USD", "EUR", "RUB", "CNY", "GBP", "TRY", "KGS", "UZS", "AED"}

// fxRank ставит ходовые валюты в начало списка.
func fxRank(code string) int {
	for i, c := range fxPopular {
		if c == code {
			return i
		}
	}
	return len(fxPopular)
}

// Earliest — первый день, за который у нас есть хоть что-то.
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

// MarkProbed запоминает, что за этот день мы банк уже спрашивали.
func (s *FxStore) MarkProbed(ctx context.Context, day time.Time, found int) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO fx_probed (day, found, probed_at) VALUES ($1,$2,now())
		ON CONFLICT (day) DO UPDATE SET found = EXCLUDED.found, probed_at = now()`, day, found)
	return err
}

// NextToProbe возвращает ближайший неопрошенный день, двигаясь от сегодня в
// прошлое, но не глубже floor. Пустой результат означает, что догружать нечего.
func (s *FxStore) NextToProbe(ctx context.Context, floor time.Time) (time.Time, bool, error) {
	// max, а не ORDER BY … LIMIT 1: когда опрашивать нечего, выборка со
	// строками вернула бы «нет строк», то есть ошибку, и догрузка каждый час
	// жаловалась бы в журнал на то, что она закончила работу. Агрегат в этом
	// случае честно отдаёт одну строку с NULL.
	var d *time.Time
	err := s.db.QueryRow(ctx, `
		SELECT max(g)::date FROM generate_series($1::date, CURRENT_DATE, interval '1 day') g
		WHERE g::date NOT IN (SELECT day FROM fx_probed)`, floor).Scan(&d)
	if err != nil || d == nil {
		return time.Time{}, false, err
	}
	return *d, true, nil
}

// EmptyRunBelow считает, сколько подряд идущих опрошенных дней ниже day не дали
// ни одного курса. Длинная цепочка означает, что мы дошли до края чужого окна:
// дальше банк молчит, и стучаться туда бессмысленно.
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

// KnownCodes перечисляет валюты, которые встречались в дневном архиве.
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
