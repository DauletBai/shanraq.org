package articles

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// The exchange rate page.
//
// One currency, one period, one line. Everything else is a caption to it: where
// the line was highest and lowest, how far it moved over the period, what that is
// in percent. The point of the page is that someone who knows only the word
// "rate" should see the shape of the curve at a glance, while someone who came to
// do arithmetic finds the numbers beside it and does not have to.

// fxPeriods are the periods on offer. The key travels in the address, so it is
// short and it does not change: people link to these.
var fxPeriods = []string{"month", "year", "five", "all"}

// fxDefaultCode is the currency the page opens on.
const fxDefaultCode = "USD"

// fxMaxPoints is how many points are worth drawing. Five years of daily rates is
// about thirteen hundred marks on a line a thousand units wide: neighbouring
// points land closer together than the pen is thick, and the surplus adds not one
// distinguishable bend while it does add weight to the page.
const fxMaxPoints = 480

// FxTick is one mark on the value axis.
type FxTick struct {
	Label string
	Pos   float64 // проценты сверху вниз
}

// FxAxis is one mark on the time axis. Its position is a percentage rather than a
// point number: months differ in length, and an even layout would separate a
// label from the point it belongs to.
type FxAxis struct {
	Label string
	Pos   float64
	Grid  bool
}

// FxChart is a line ready to draw, in the picture's own coordinates.
type FxChart struct {
	Line   string
	Area   string
	Y      []FxTick
	X      []FxAxis
	Width  int
	Height int
	Points int
	// Hover carries the readout data as ready JSON. Empty when there are too few
	// points for pointing at one to mean anything.
	Hover string
}

// FxYear is one row of the annual table.
type FxYear struct {
	Year  int
	Close string
	Diff  string
	Pct   string
	Up    bool
	Down  bool
}

// FxPage is the page's data.
type FxPage struct {
	Base
	Desc       string
	Currencies []FxCurrency
	Code       string
	Name       string
	Quant      int
	Period     string
	Chart      FxChart
	HasData    bool
	Last       string
	LastDay    string
	Diff       string
	Pct        string
	Up         bool
	Down       bool
	Min        string
	MinDay     string
	Max        string
	MaxDay     string
	Avg        string
	Spread     string
	Since      string
	Years      []FxYear
	Monthly    bool
	DeepSince  string
	DailySince string
	Macro      MacroBlock
}

// handleRates serves the exchange rate page.
func (m *Module) handleRates(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	ctx := r.Context()

	page := FxPage{Period: fxPeriod(r.URL.Query().Get("p"))}
	if m.fx != nil {
		page.Currencies, _ = m.fx.Currencies(ctx)
	}
	page.Code = fxCode(r.URL.Query().Get("c"), page.Currencies)
	page.Quant = 1
	for _, c := range page.Currencies {
		if c.Code == page.Code {
			page.Name, page.Quant = c.Name, c.Quant
			page.Since = fxDate(c.Since)
			break
		}
	}
	if page.Name == "" {
		page.Name = page.Code
	}

	if m.fx != nil {
		m.fillRates(ctx, &page, lang)
	}
	page.Macro = m.macroCached(ctx, lang)

	// Every currency gets its own address, title and description. Folding them
	// into one page would mean that someone searching for the rouble rate does not
	// find our rouble rate page: a search engine only knows the address we
	// declared as ours. The period stays out of the canonical — it is the same
	// rate shown at different depths.
	page.Base = m.base(r, T(lang, "fx.title"), lang)
	page.Desc = T(lang, "fx.desc")
	if page.Code != fxDefaultCode {
		filter := "c=" + page.Code
		page.Base.CanonURL = canonURL("/rates", filter, lang)
		page.Base.LangLinks = langLinks("/rates", filter)
		page.Base.Title = fmt.Sprintf(T(lang, "fx.title_cur"), page.Code)
		page.Desc = fmt.Sprintf(T(lang, "fx.desc_cur"), fxSubject(page.Name, page.Code, lang))
	}
	m.render(w, "rates", page)
}

// fxSubject is how to name a currency in the page description. The name arrives
// from the bank in capitals and in Russian only, so the other two languages keep
// the code: an invented translation is worse than an honest "USD".
func fxSubject(name, code, lang string) string {
	if lang != LangRU {
		return code
	}
	n := fxTidyName(name)
	if n == "" {
		return code
	}
	return n + " (" + code + ")"
}

// fxTidyName turns the bank's name into ordinary writing: "ДОЛЛАР США" → "Доллар
// США". Words longer than three letters go lower case, short ones stay as they
// are — those are abbreviations like США, ОАЭ and СДР, and "сша" must not come out
// of them.
func fxTidyName(name string) string {
	words := strings.Fields(strings.TrimSpace(name))
	for i, w := range words {
		if len([]rune(w)) > 3 {
			words[i] = strings.ToLower(w)
		}
	}
	out := strings.Join(words, " ")
	r := []rune(out)
	if len(r) == 0 {
		return ""
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// fillRates assembles the series and every figure that goes with it.
func (m *Module) fillRates(ctx context.Context, page *FxPage, lang string) {
	var (
		pts []FxPoint
		err error
	)
	now := time.Now().UTC()
	switch page.Period {
	case "all":
		page.Monthly = true
		pts, err = m.fx.SeriesMonthly(ctx, page.Code, time.Time{})
	case "five":
		page.Monthly = true
		pts, err = m.fx.SeriesMonthly(ctx, page.Code, monthStart(now.AddDate(-5, 0, 0)))
	case "year":
		pts, err = m.fx.Series(ctx, page.Code, now.AddDate(-1, 0, 0))
	default:
		pts, err = m.fx.Series(ctx, page.Code, now.AddDate(0, 0, -30))
	}
	if err != nil || len(pts) < 2 {
		return
	}

	// Scale to the familiar multiple: the bank quotes yen in hundreds, and the
	// eye looks for exactly the number it saw in the rate ticker.
	q := float64(page.Quant)
	for i := range pts {
		pts[i].Value *= q
	}

	page.HasData = true
	page.Chart = fxBuildChartWith(pts, page.Period, lang, fxChartOpts{
		Name: page.Code, Unit: "₸",
	})

	first, last := pts[0], pts[len(pts)-1]
	lo, hi, sum := pts[0], pts[0], 0.0
	for _, p := range pts {
		if p.Value < lo.Value {
			lo = p
		}
		if p.Value > hi.Value {
			hi = p
		}
		sum += p.Value
	}
	diff := last.Value - first.Value

	page.Last, page.LastDay = fxNum(last.Value), fxDate(last.Day)
	page.Diff = fxDelta(diff, last.Value)
	if first.Value > 0 {
		page.Pct = fxPct(diff / first.Value * 100)
	}
	page.Up, page.Down = diff > 0, diff < 0
	page.Min, page.MinDay = fxNum(lo.Value), fxDate(lo.Day)
	page.Max, page.MaxDay = fxNum(hi.Value), fxDate(hi.Day)
	page.Avg = fxNum(sum / float64(len(pts)))
	page.Spread = fxNum(hi.Value - lo.Value)
	page.Years = fxYears(pts)
}

// fxChartOpts says how to label a frame's points in the readout. Each frame
// formats its values its own way: an exchange rate in tiyn, money supply in
// trillions, inflation in percent — and one rounding for all of them would turn
// half the readouts into rows of identical zeros.
type fxChartOpts struct {
	// Name is the series' name. A lone line needs none: it is the only one in
	// frame.
	Name string
	// Unit is the unit of measurement appended to the value.
	Unit string
	// Monthly marks a series whose point is a month rather than a day.
	Monthly bool
	// Annual marks a series whose point is a whole year.
	Annual bool
	// Format prints the value; empty means print it as an exchange rate.
	Format func(float64) string
}

// label resolves how much of the date a point should show.
func (o fxChartOpts) label(period string) string {
	switch {
	case o.Annual:
		return fxLabelYear
	case o.Monthly || period == "all" || period == "five":
		return fxLabelMonth
	}
	return fxLabelDay
}

// fxBuildChart turns a series into the picture's coordinates.
func fxBuildChart(pts []FxPoint, period, lang string) FxChart {
	return fxBuildChartWith(pts, period, lang, fxChartOpts{})
}

// fxBuildChartWith is the same with labels for the readout.
func fxBuildChartWith(pts []FxPoint, period, lang string, o fxChartOpts) FxChart {
	const w, h = 1000.0, 360.0
	pts = fxThin(pts, fxMaxPoints)

	lo, hi := pts[0].Value, pts[0].Value
	for _, p := range pts {
		lo = math.Min(lo, p.Value)
		hi = math.Max(hi, p.Value)
	}
	lo, hi, step := fxScale(lo, hi)

	x := func(i int) float64 { return w * float64(i) / float64(len(pts)-1) }
	y := func(v float64) float64 { return h - h*(v-lo)/(hi-lo) }

	var line strings.Builder
	for i, p := range pts {
		if i > 0 {
			line.WriteByte(' ')
		}
		if i == 0 {
			line.WriteByte('M')
		} else {
			line.WriteByte('L')
		}
		fmt.Fprintf(&line, "%.1f %.1f", x(i), y(p.Value))
	}
	area := line.String() + fmt.Sprintf(" L%.1f %.1f L0 %.1f Z", w, h, h)

	ticks := []FxTick{}
	for v := hi; v >= lo-step/2; v -= step {
		ticks = append(ticks, FxTick{Label: fxAxisNum(v, step), Pos: (hi - v) / (hi - lo) * 100})
	}

	format := o.Format
	if format == nil {
		format = fxNum
	}
	hover := FxHover{L: fxPointLabels(pts, o.label(period), lang)}
	series := FxHoverSeries{N: o.Name}
	for _, p := range pts {
		series.Y = append(series.Y, hoverPct(y(p.Value)/h*100))
		series.V = append(series.V, fxUnit(format(p.Value), o.Unit))
	}
	hover.S = append(hover.S, series)

	return FxChart{
		Line: line.String(), Area: area, Y: ticks, X: fxAxis(pts, period, lang),
		Width: int(w), Height: int(h), Points: len(pts),
		Hover: hoverJSON(hover),
	}
}

// fxUnit appends a unit of measurement to a value. Percent sits against the
// number; everything else takes a non-breaking space, so that "55,8 трлн ₸" does
// not wrap in the middle of a sum.
func fxUnit(v, unit string) string {
	switch {
	case unit == "":
		return v
	case unit == "%":
		return v + " %"
	}
	return v + "\u00a0" + unit
}

// fxAxis marks the time axis on the chosen period's natural boundaries: a month
// by days, a year by months, five years and the whole history by years. Even
// marks every N points would give labels like "14.03" and "27.07", by which one
// can neither find a place nor compare two years.
func fxAxis(pts []FxPoint, period, lang string) []FxAxis {
	if len(pts) < 2 {
		return nil
	}
	pos := func(i int) float64 { return 100 * float64(i) / float64(len(pts)-1) }
	out := []FxAxis{}

	switch period {
	case "month":
		// Every day by its number: a month is read by day-of-month, not by date.
		for i, p := range pts {
			out = append(out, FxAxis{Label: strconv.Itoa(p.Day.Day()), Pos: pos(i),
				Grid: p.Day.Day() == 1})
		}
	case "year":
		for i, p := range pts {
			if i > 0 && p.Day.Month() == pts[i-1].Day.Month() {
				continue
			}
			lab := fxMonthShort(p.Day.Month(), lang)
			if p.Day.Month() == time.January {
				lab = strconv.Itoa(p.Day.Year())
			}
			out = append(out, FxAxis{Label: lab, Pos: pos(i), Grid: p.Day.Month() == time.January})
		}
	default:
		// Five years and the whole history: label the years, but not three dozen
		// of them — over thirty-three years the labels would fuse into a band.
		step := 1
		if years := pts[len(pts)-1].Day.Year() - pts[0].Day.Year(); years > 12 {
			step = 5
		} else if years > 6 {
			step = 2
		}
		for i, p := range pts {
			if i > 0 && p.Day.Year() == pts[i-1].Day.Year() {
				continue
			}
			if p.Day.Year()%step != 0 {
				continue
			}
			out = append(out, FxAxis{Label: strconv.Itoa(p.Day.Year()), Pos: pos(i), Grid: true})
		}
	}
	return fxSpaceOut(out)
}

// fxMinGap is the smallest gap between axis labels, as a percentage of width.
const fxMinGap = 2.6

// fxSpaceOut removes labels that run into their neighbour on the right. A series
// almost always begins mid-month, and the first point's "авг" lands against the
// "сен" of 1 September, fusing with it into an unreadable smudge. We walk from the
// end because the last label is today's date, and it is needed more than any
// before it.
func fxSpaceOut(a []FxAxis) []FxAxis {
	if len(a) < 2 {
		return a
	}
	last := 101.0
	for i := len(a) - 1; i >= 0; i-- {
		if last-a[i].Pos < fxMinGap {
			// Only the label goes: the grid mark rules the chart without it, and
			// removing it along with the text would lose the month's line too.
			a[i].Label = ""
			continue
		}
		last = a[i].Pos
	}
	return a
}

// fxMonthShort is a month's short name in the reader's language.
func fxMonthShort(m time.Month, lang string) string {
	names := map[string][12]string{
		LangKZ: {"қаң", "ақп", "нау", "сәу", "мам", "мау", "шіл", "там", "қыр", "қаз", "қар", "жел"},
		LangRU: {"янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"},
		LangEN: {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	}
	row, ok := names[lang]
	if !ok {
		row = names[LangRU]
	}
	return row[int(m)-1]
}

// fxScale picks the bounds and step of the value axis so that the labels are
// round numbers. An axis ruled at "472,91" and "478,79" reads worse than none at
// all: by it one can neither estimate a point's value nor compare two charts.
func fxScale(lo, hi float64) (float64, float64, float64) {
	if hi <= lo {
		hi = lo + math.Max(math.Abs(lo)*0.01, 0.01)
	}
	step := fxNiceStep((hi - lo) / 4)
	lo = math.Floor(lo/step) * step
	hi = math.Ceil(hi/step) * step
	// An exchange rate is never negative, and room for a minus on the axis is
	// chart height thrown away.
	if lo < 0 {
		lo = 0
	}
	return lo, hi, step
}

// fxNiceStep rounds a step up to one, two, two and a half, or five times the
// nearest power of ten.
func fxNiceStep(raw float64) float64 {
	if raw <= 0 {
		return 1
	}
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	for _, m := range []float64{1, 2, 2.5, 5} {
		if raw <= m*mag {
			return m * mag
		}
	}
	return 10 * mag
}

// fxThin thins a series down to max points, keeping the first and the last.
func fxThin(pts []FxPoint, max int) []FxPoint {
	if len(pts) <= max {
		return pts
	}
	out := make([]FxPoint, 0, max)
	for i := 0; i < max-1; i++ {
		out = append(out, pts[i*(len(pts)-1)/(max-1)])
	}
	return append(out, pts[len(pts)-1])
}

// fxYears reduces a series to years: what each year closed at and how far it
// moved the rate. A long line cannot be broken into years by eye, and "which year
// was a bad one" is exactly how the question gets asked.
func fxYears(pts []FxPoint) []FxYear {
	last := map[int]FxPoint{}
	years := []int{}
	for _, p := range pts {
		y := p.Day.Year()
		if _, ok := last[y]; !ok {
			years = append(years, y)
		}
		last[y] = p
	}
	if len(years) < 2 {
		return nil
	}
	out := make([]FxYear, 0, len(years))
	for i := len(years) - 1; i >= 0; i-- {
		y := years[i]
		row := FxYear{Year: y, Close: fxNum(last[y].Value)}
		if i > 0 {
			prev := last[years[i-1]].Value
			d := last[y].Value - prev
			row.Diff, row.Up, row.Down = fxDelta(d, last[y].Value), d > 0, d < 0
			if prev > 0 {
				row.Pct = fxPct(d / prev * 100)
			}
		}
		out = append(out, row)
	}
	return out
}

// fxPeriod maps a period from the address onto one we know.
func fxPeriod(s string) string {
	for _, p := range fxPeriods {
		if s == p {
			return p
		}
	}
	return "month"
}

// fxCode maps a currency from the address onto one we hold.
func fxCode(s string, list []FxCurrency) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	for _, c := range list {
		if c.Code == s {
			return s
		}
	}
	// The list runs from the currencies in demand to the rest, and the first is
	// the dollar if we have it. A typo in the address must not leave someone with
	// nothing.
	if len(list) > 0 {
		return list[0].Code
	}
	return fxDefaultCode
}

// fxNum prints a rate: two decimal places for ordinary numbers and four for
// small ones, where two would turn the whole series into identical zeros.
func fxNum(v float64) string { return fxFormat(v, fxDigits(v)) }

// fxAxisNum prints an axis mark to the precision its step sets. A step of ten
// tenge needs no tiyn: "480,00" instead of "480" is two digits of noise in every
// label.
func fxAxisNum(v, step float64) string {
	digits := 0
	if step < 1 {
		digits = int(math.Ceil(-math.Log10(step)))
	}
	return fxFormat(v, digits)
}

// fxFormat prints a number with a comma decimal and spaces between thousands.
func fxFormat(v float64, digits int) string {
	s := strconv.FormatFloat(v, 'f', digits, 64)
	whole, frac, _ := strings.Cut(s, ".")
	neg := strings.HasPrefix(whole, "-")
	whole = strings.TrimPrefix(whole, "-")
	// Thousands are separated by a narrow non-breaking space: it stops the number
	// wrapping mid-thousand and does not push the digits as far apart as an
	// ordinary space would.
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteRune('\u202f')
		}
		b.WriteRune(r)
	}
	out := b.String()
	if frac != "" {
		out += "," + frac
	}
	if neg {
		out = "−" + out
	}
	return out
}

// fxDelta prints a change with its sign, to the same precision as the rate
// itself: in a column holding "+2,25" beside "+0,5900" the eye trips over the
// differing fraction lengths and stops comparing the numbers. The precision comes
// from ref — the size of the rate, not the size of the change.
//
// The minus here is a real one, not a hyphen: in a column of numbers a hyphen does
// not match the plus in height, and the column shimmers.
func fxDelta(v, ref float64) string {
	return fxSign(v, fxFormat(v, fxDigits(ref)))
}

// fxPct prints percentages — always to two places, whatever the rate.
func fxPct(v float64) string {
	return fxSign(v, fxFormat(v, 2)) + "%"
}

// fxSign prepends a plus: the minus is already in the number.
func fxSign(v float64, s string) string {
	if v > 0 {
		return "+" + s
	}
	return s
}

// fxDigits is how many decimal places suit a rate of this size. For a currency
// cheaper than the tenge, two would turn the whole series into identical zeros.
func fxDigits(v float64) int {
	if math.Abs(v) < 1 {
		return 4
	}
	return 2
}

// fxDate prints a date short. One format for all three languages: a numeric date
// reads the same in each and needs no month names in three cases.
func fxDate(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format("02.01.2006")
}
