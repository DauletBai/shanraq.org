package articles

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"
)

// The section "How the tenge rate and inflation are formed".
//
// The page is not written, it is computed. Once a month the National Bank
// publishes how many tenge exist and how many dollars sit in the reserves; once a
// year the World Bank publishes inflation. All that happens here is laying one
// beside the other: the quantity of money, the quantity of currency it can be
// exchanged for, and the rate that follows from the ratio.
//
// No claims about anyone's intentions: the series are shown and set against each
// other, and the reader draws the conclusion. A figure with a source behind it
// persuades better than any adjective.

// MacroChart is two lines in one coordinate system.
type MacroChart struct {
	A      string
	B      string
	Area   string
	Y      []FxTick
	X      []FxAxis
	Width  int
	Height int
	// Log marks a logarithmic scale — the reader has to be told.
	Log bool
	// MinWidth, when set, is the frame's width in pixels inside a sideways
	// scroll: a series too dense to read at the column's width gets room and
	// the reader moves along it.
	MinWidth int
	// Hover carries the readout data for the frame, as ready JSON.
	Hover string
}

// MacroYear is one row of the annual summary.
type MacroYear struct {
	Year  int
	CPI   string
	Money string
	Rate  string
	Fund  string
	// FundDown marks a year in which the National Fund shrank.
	FundDown bool
}

// MacroErode is what a thousand tenge of a given year turned into.
type MacroErode struct {
	Year int
	// Kept is what that thousand is worth in today's tenge by purchasing
	// power.
	Kept string
	// Then and Now are the dollars a thousand tenge bought then and buys now.
	Then string
	Now  string
}

// MacroBlock is the whole section.
type MacroBlock struct {
	Has bool

	// M3 carries two lines: the base money the National Bank issues itself and
	// the broad money the banking system ends up with. One line answered "how
	// much" and left "why" to the reader; two answer both, because the gap
	// between them is exactly the part banks add by lending.
	M3      MacroChart
	M3Last  string
	M3Month string
	M3Year  string
	M3Up    bool
	// BaseLast is the monetary base — the tenge the Bank itself put into
	// circulation. BaseYear is its growth over the year.
	BaseLast string
	BaseYear string
	BaseUp   bool
	// Mult is M3 divided by the base: how many tenge of broad money stand on
	// each tenge the Bank issued.
	Mult string

	Cmp      MacroChart
	CmpFrom  string
	MoneyMul string
	RateMul  string

	Cover      FxChart
	CoverLast  string
	CoverMonth string
	RateAt     string

	Res      MacroChart
	ResLast  string
	FundLast string
	ResMonth string

	CPILast string
	CPIYear string

	Years []MacroYear
	Erode []MacroErode
	// NowDollars is the dollars a thousand tenge buys today.
	NowDollars string

	// Now is the National Bank's indicator panel: rate, inflation, target,
	// TONIA.
	Now MacroNow
	// Formulas are the rules by which prices and the rate follow from those
	// quantities.
	Formulas []MacroFormula

	// RateCPI puts the Bank's rate and annual inflation on one scale.
	RateCPI     MacroChart
	RateCPIFrom string
	RateSince   string
	// The Bank's record against its own target, counted from the same series
	// the chart draws. Hard-coding these would let the verdict drift away from
	// the picture it stands under, which is the one failure this block cannot
	// afford.
	RateYears          string // сколько лет на графике
	RateMissed         string // из них лет выше цели
	RateAvg10          string // средняя инфляция за последнее десятилетие
	RateTarget         string // объявленная цель
	RateAllYearsMissed bool
	// What the record cost a person holding money: a thousand tenge of the
	// oldest year the erosion table reaches, and what it buys now.
	RateErodeYear string
	RateErodeKept string
	// The gap the whole page is about: how many times the money supply grew
	// against how many times the country's output did, and the ratio between
	// them. Counted from the same year the comparison chart starts.
	RateGdpMul   string
	RateMoneyMul string
	RateGap      string
	// Who ended up holding the money and what it cost to borrow: the share of
	// the money supply commercial banks created themselves, and how many of the
	// chart's years credit cost more than prices rose.
	RateBankShare string
	RateDearYears string
	// RateYearsNum is the bare count, for a sentence that supplies its own word
	// for "years": "22 года из 32" needs the second number naked.
	RateYearsNum string
}

// ready reports whether the section was assembled in full.
func (b MacroBlock) ready() bool {
	return b.Has && b.Cmp.A != "" && b.CoverLast != "" && b.Res.A != "" && b.CPILast != ""
}

// buildMacro assembles the section. Any missing part is simply not shown: half an
// analysis is more use than an empty page.
func (m *Module) buildMacro(ctx context.Context, lang string) MacroBlock {
	var b MacroBlock
	if m.macro == nil || m.fx == nil {
		return b
	}

	m3, _ := m.macro.Series(ctx, MacroM3)
	baseMoney, _ := m.macro.Series(ctx, MacroBase)
	res, _ := m.macro.Series(ctx, MacroReserves)
	fund, _ := m.macro.Series(ctx, MacroFund)
	cpi, _ := m.macro.Series(ctx, MacroCPI)
	rate, _ := m.fx.SeriesMonthly(ctx, fxDefaultCode, time.Time{})
	if len(m3) < 24 {
		return b
	}
	b.Has = true

	// 1. How many tenge exist, and where they came from.
	//
	// Broad money on its own says how much and stays silent on why. The
	// monetary base beside it says why: it is the tenge the National Bank
	// issued itself, and everything above that line is what banks built on top
	// of it by lending. The two together name the author of the growth without
	// a word of commentary.
	//
	// The frame is scaled to trillions, but a readout on January 1994 must not
	// say "0.0 trn ₸": each point names its own magnitude, so the early years
	// read in billions and millions.
	tenge := func(v float64) string { return macroTenge(v*1e6, lang) }
	b.M3 = macroTwoLinesWith(macroScaled(m3, 1e6), macroScaled(baseMoney, 1e6), lang, true,
		fxChartOpts{Name: T(lang, "fx.col_money"), Unit: "₸", Format: tenge,
			AxisFormat: func(v float64) string { return macroTengeRound(v*1e6, lang) }},
		fxChartOpts{Name: T(lang, "fx.col_base"), Unit: "₸", Format: tenge},
	) // trillions of tenge
	last := m3[len(m3)-1]
	b.M3Last = macroTenge(last.Value, lang)
	b.M3Month = macroMonth(last.Period, lang)
	if prev, ok := macroAt(m3, last.Period.AddDate(-1, 0, 0)); ok && prev > 0 {
		d := (last.Value/prev - 1) * 100
		b.M3Year, b.M3Up = fxPct(d), d > 0
	}
	if bl, ok := macroLastPoint(baseMoney); ok && bl.Value > 0 {
		b.BaseLast = macroTenge(bl.Value, lang)
		if prev, ok := macroAt(baseMoney, bl.Period.AddDate(-1, 0, 0)); ok && prev > 0 {
			d := (bl.Value/prev - 1) * 100
			b.BaseYear, b.BaseUp = fxPct(d), d > 0
		}
		// The multiplier is taken for the same month as M3, never for the
		// latest of each: the two series can be published a month apart, and
		// dividing across that gap would invent a jump that never happened.
		if bv, ok := macroAt(baseMoney, last.Period); ok && bv > 0 {
			// The same ratio the formula prints, derived the same way: the
			// note and the formula must not disagree about the multiplier by a
			// hundredth.
			b.Mult = fxFormat(macroShownRatio(last.Value, bv), 2)
		}
	}

	// 2. Money supply and the dollar rate on one scale. Both lines are indexed to
	//    a hundred in their first shared month: one is measured in tenge and the
	//    other in tenge per dollar, and they can only be compared as multiples.
	if len(rate) > 24 {
		b.Cmp, b.CmpFrom, b.MoneyMul, b.RateMul = macroCompare(m3, rate, lang)
	}

	// 3. How many tenge there are per dollar of reserves. Not a rate and not a
	//    forecast: the ratio of two published figures.
	if len(res) > 24 {
		cover := macroCover(m3, res)
		if len(cover) > 24 {
			b.Cover = fxBuildChartWith(cover, "all", lang, fxChartOpts{
				Unit: T(lang, "fx.u_kzt_usd"),
			})
			c := cover[len(cover)-1]
			// The figure in the note is the one the formula prints, derived the
			// same way: the two must not disagree by half a tenge, because a
			// reader who spots that stops trusting both.
			shown := c.Value
			if rv, ok := macroAt(res, c.Day); ok && rv > 0 {
				if mv, ok := macroAt(m3, c.Day); ok && mv > 0 {
					shown = macroShownRatio(mv, rv)
				}
			}
			b.CoverLast, b.CoverMonth = fxNum(shown), macroMonth(c.Day, lang)
			if v, ok := macroRateAt(rate, c.Day); ok {
				b.RateAt = fxNum(v)
			}
		}
	}

	// 4. Reserves and the National Fund.
	if len(res) > 24 && len(fund) > 24 {
		dollars := func(v float64) string { return macroDollars(v*1e3, lang) }
		b.Res = macroTwoLinesWith(macroScaled(res, 1e3), macroScaled(fund, 1e3), lang, false,
			fxChartOpts{Name: T(lang, "fx.res_key_a"), Unit: "$", Format: dollars},
			fxChartOpts{Name: T(lang, "fx.col_fund"), Unit: "$", Format: dollars},
		) // billions of dollars
		b.ResLast = macroDollars(res[len(res)-1].Value, lang)
		b.FundLast = macroDollars(fund[len(fund)-1].Value, lang)
		b.ResMonth = macroMonth(res[len(res)-1].Period, lang)
	}

	// 5. Inflation by year. The standalone frame is gone — a single line of
	//    annual percentages said nothing the rate chart below does not say
	//    better, and it said it beside the rate rather than against it. The
	//    figures stay: the rate chart and the formulas both need them.
	if len(cpi) > 5 {
		c := cpi[len(cpi)-1]
		b.CPILast, b.CPIYear = fxPct(c.Value), fmt.Sprintf("%d", c.Period.Year())
	}

	b.Years = macroYears(cpi, m3, rate, fund)
	b.Erode, b.NowDollars = macroErosion(cpi, rate)
	// The oldest row is the one that lands: a thousand tenge from thirty years
	// ago says more about the record than any percentage.
	if len(b.Erode) > 0 {
		b.RateErodeYear = fmt.Sprintf("%d", b.Erode[0].Year)
		b.RateErodeKept = b.Erode[0].Kept
	}

	// 6. The decision everything follows from: the National Bank's rate. It is
	//    the one quantity in this analysis with a day, a signature and a minute
	//    of the meeting; the other series only measure what happened after it.
	refi, _ := m.macro.Series(ctx, MacroRefiRate)
	baseRate, _ := m.macro.Series(ctx, MacroBaseRate)
	policy := macroPolicyRate(refi, baseRate)
	if len(policy) > 0 {
		b.RateSince = fmt.Sprintf("%d", policy[0].Day.Year())
	}
	if len(policy) > 4 && len(cpi) > 4 {
		ra, rb := macroAlignYears(macroRateByYear(policy), cpi)
		if len(ra) > 4 {
			// Logarithmic scale: in 1994 inflation ran into thousands of
			// percent, and on a linear scale today's two-digit figures would lie
			// in a single line along the bottom.
			var nYears int64
			b.RateDearYears = macroDearYears(ra, rb)
			b.RateYears, b.RateMissed, b.RateAvg10, b.RateAllYearsMissed =
				macroTargetRecord(rb, cpiTargetValue(m, ctx))
			b.RateYearsNum = b.RateYears
			if v, err := strconv.ParseInt(b.RateYears, 10, 64); err == nil {
				nYears = v
				b.RateYears = b.RateYears + " " + macroYearsWord(nYears, lang)
			}
			b.RateCPI = macroTwoLinesWith(ra, rb, lang, true,
				fxChartOpts{Name: T(lang, "fx.rate_key_a"), Unit: "%", Annual: true,
					// Three decades of annual figures at the column's width
					// give twenty pixels a year, and the year-to-year movement
					// — the whole point of the picture — disappears.
					PxPerPoint: 38,
					AxisFormat: func(v float64) string { return fxFormat(v, 0) + " %" }},
				fxChartOpts{Name: T(lang, "fx.col_cpi"), Unit: "%", Annual: true},
			)
			b.RateCPIFrom = fmt.Sprintf("%d", ra[0].Day.Year())
			b.RateTarget = macroPctTrim(cpiTargetValue(m, ctx))
		}
	}

	// Money against output, from the first year both are known. This is the one
	// figure on the page that needs no theory of what causes what: the country
	// produced this much more, and this much more money was printed for it.
	b.RateGdpMul, b.RateMoneyMul, b.RateGap = macroMoneyVsOutput(m3, gdpSeries(m, ctx), lang)
	if last, ok := macroLastPoint(m3); ok && last.Value > 0 {
		if bv, ok := macroAt(baseMoney, last.Period); ok && bv > 0 {
			b.RateBankShare = fxFormat((last.Value-bv)/last.Value*100, 0) + " %"
		}
	}

	cpiNow, _ := m.macro.Series(ctx, MacroCPINow)
	cpiTarget, _ := m.macro.Series(ctx, MacroCPITarget)
	tonia, _ := m.macro.Series(ctx, MacroTonia)
	gdp, _ := m.macro.Series(ctx, MacroGDP)
	in := macroFormulaInput{
		lang: lang, m3: m3, base: baseMoney, res: res, gdp: gdp,
		rate: policy, cpiNow: cpiNow, cpiTarget: cpiTarget, fxRate: rate,
	}
	b.Now = buildNow(in, tonia)
	b.Formulas = buildFormulas(in)
	return b
}

// macroErosion works out what a thousand tenge of various years turned into.
//
// Two measures at once, because they answer different questions. Purchasing power
// is how many of today's tenge are needed to buy the same thing: inflation
// accumulated over the years. Dollars are what that same thousand bought in
// currency then and buys now. The first is about prices inside the country, the
// second about what those prices are worth outside it.
func macroErosion(cpi []MacroPoint, rate []FxPoint) ([]MacroErode, string) {
	if len(cpi) < 5 || len(rate) < 12 {
		return nil, ""
	}
	byYear := map[int]float64{}
	for _, c := range cpi {
		byYear[c.Period.Year()] = c.Value
	}
	last := rate[len(rate)-1]
	nowRate := last.Value
	if nowRate <= 0 {
		return nil, ""
	}
	rateAtYear := map[int]float64{}
	for _, p := range rate {
		if _, ok := rateAtYear[p.Day.Year()]; !ok {
			rateAtYear[p.Day.Year()] = p.Value
		}
	}

	out := []MacroErode{}
	for _, y := range []int{1995, 2000, 2005, 2010, 2015, 2020} {
		if y >= last.Day.Year() {
			continue
		}
		// Price growth accumulated from the start of year y to the last year with
		// data.
		growth := 1.0
		ok := true
		for yy := y; yy < last.Day.Year(); yy++ {
			v, has := byYear[yy]
			if !has {
				ok = false
				break
			}
			growth *= 1 + v/100
		}
		r, hasRate := rateAtYear[y]
		if !ok || !hasRate || r <= 0 || growth <= 0 {
			continue
		}
		out = append(out, MacroErode{
			Year: y,
			Kept: fxFormat(1000/growth, 0),
			Then: fxFormat(1000/r, 2),
			Now:  fxFormat(1000/nowRate, 2),
		})
	}
	return out, fxFormat(1000/nowRate, 2)
}

// macroOne prints a value to one decimal place. Trillions of tenge and billions
// of dollars carry no meaning past the first: "55.8" is the figure a reader
// repeats, "55.83" is two digits of noise in every readout.
func macroOne(v float64) string { return fxFormat(v, 1) }

// macroPoints turns a macro series into chart points.
func macroPoints(pts []MacroPoint) []FxPoint {
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, FxPoint{Day: p.Period, Value: p.Value})
	}
	return out
}

// macroAt returns a series' value for one specific month.
func macroAt(pts []MacroPoint, when time.Time) (float64, bool) {
	for _, p := range pts {
		if p.Period.Equal(when) {
			return p.Value, true
		}
	}
	return 0, false
}

// macroRateAt returns the rate for a point's month, taking the nearest one no
// later than it.
func macroRateAt(rate []FxPoint, when time.Time) (float64, bool) {
	var v float64
	var ok bool
	for _, p := range rate {
		if p.Day.After(when.AddDate(0, 1, 0)) {
			break
		}
		v, ok = p.Value, true
	}
	return v, ok
}

// macroCover works out the tenge per dollar of reserves: money supply in
// millions of tenge divided by reserves in millions of dollars gives tenge per
// dollar.
func macroCover(m3 []MacroPoint, res []MacroPoint) []FxPoint {
	byMonth := make(map[string]float64, len(res))
	for _, r := range res {
		byMonth[r.Period.Format("2006-01")] = r.Value
	}
	out := make([]FxPoint, 0, len(m3))
	for _, p := range m3 {
		r, ok := byMonth[p.Period.Format("2006-01")]
		if !ok || r <= 0 {
			continue
		}
		out = append(out, FxPoint{Day: p.Period, Value: p.Value / r})
	}
	return out
}

// macroTimes is how many times over a series grew from its first value to its
// last, together with the word "times" in the right form.
func macroTimes[T any](pts []T, lang string) string {
	first, last, ok := macroEnds(pts)
	if !ok || first <= 0 {
		return ""
	}
	return macroMul(last/first, lang)
}

// macroEnds takes the first and last value out of a series of either of our
// kinds.
func macroEnds[T any](pts []T) (float64, float64, bool) {
	if len(pts) < 2 {
		return 0, 0, false
	}
	val := func(x any) (float64, bool) {
		switch v := x.(type) {
		case MacroPoint:
			return v.Value, true
		case FxPoint:
			return v.Value, true
		}
		return 0, false
	}
	f, ok1 := val(any(pts[0]))
	l, ok2 := val(any(pts[len(pts)-1]))
	return f, l, ok1 && ok2
}

// macroMul prints a multiple together with its word: "6 776 times", "2.4 times".
//
// The word has to be declined: in Russian 2, 3 and 4 take "раза", the rest take
// "раз", and a fractional number always takes "раза". Without this the text says
// "в 97,2 раз", and the reader trips over the phrase instead of looking at the
// number.
func macroMul(x float64, lang string) string {
	whole := x >= 10
	n := fxFormat(x, 1)
	if whole {
		n = fxFormat(math.Round(x), 0)
	}
	switch lang {
	case LangKZ:
		return n + " есе"
	case LangEN:
		return n + " times"
	}
	if !whole {
		return n + " раза"
	}
	return n + " " + ruTimesWord(int64(math.Round(x)))
}

// ruTimesWord picks the form of the Russian word "раз" for a whole number.
func ruTimesWord(n int64) string {
	if n < 0 {
		n = -n
	}
	if t := n % 100; t >= 11 && t <= 14 {
		return "раз"
	}
	switch n % 10 {
	case 2, 3, 4:
		return "раза"
	}
	return "раз"
}

// macroCompare builds two lines indexed to a hundred in their first shared
// month.
// macroCompare returns the chart, the month both lines start from, and how many
// times each has grown SINCE THAT MONTH.
//
// The multiples used to be measured from each series' own first point, while the
// caption said "since January 1994" — the month the lines are indexed to. The
// two series do not begin together: the money supply starts in January 1994 and
// the exchange rate in November 1993, when the tenge was two months old and
// stood at 4.70. Measured from its own beginning the dollar had grown 97 times;
// measured from the month the sentence actually names, 43. The page printed the
// first number under the second date.
//
// So the multiples are computed here, from the same trimmed series the chart
// draws, and cannot drift away from the caption again.
func macroCompare(m3 []MacroPoint, rate []FxPoint, lang string) (chart MacroChart, from string, moneyMul, rateMul string) {
	a := macroPoints(m3)
	start := a[0].Day
	if rate[0].Day.After(start) {
		start = rate[0].Day
	}
	a = macroSince(a, start)
	b := macroSince(rate, start)
	if len(a) < 2 || len(b) < 2 {
		return MacroChart{}, "", "", ""
	}
	chart = macroTwoLinesWith(macroIndex(a), macroIndex(b), lang, true,
		fxChartOpts{Name: T(lang, "fx.col_money"), Format: macroOne},
		fxChartOpts{Name: T(lang, "fx.col_rate"), Format: macroOne},
	)
	return chart, macroMonthIn(start, lang), macroTimes(a, lang), macroTimes(b, lang)
}

// macroScaled converts a series into convenient units: millions of tenge into
// trillions, millions of dollars into billions. An axis labelled "60 000 000"
// neither reads nor fits — and it is the same 60 trillion.
func macroScaled(pts []MacroPoint, div float64) []FxPoint {
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, FxPoint{Day: p.Period, Value: p.Value / div})
	}
	return out
}

// macroSince drops everything earlier than the given month.
func macroSince(pts []FxPoint, from time.Time) []FxPoint {
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		if p.Day.Before(from) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// macroIndex indexes a series to a hundred at its first point.
func macroIndex(pts []FxPoint) []FxPoint {
	if len(pts) == 0 || pts[0].Value == 0 {
		return pts
	}
	base := pts[0].Value
	out := make([]FxPoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, FxPoint{Day: p.Day, Value: p.Value / base * 100})
	}
	return out
}

// macroTwoLines draws two series in one coordinate system.
//
// log switches on a logarithmic scale. It is needed where the series have parted
// by orders of magnitude: money supply grew more than six thousand times, the
// rate about a hundred, and on a linear scale the second line lies flat on zero
// as though the rate had never moved. A logarithm shows not the difference in
// size but the difference in pace — which is the whole reason two series are put
// side by side.
//
// Both series are laid on one shared time axis rather than side by side by point
// number. Where a series has no value for a period its line breaks, instead of
// being stretched across the full width: reserves run from 1993 and the National
// Fund from 2001, and stretching drew the Fund seven years before it existed.
func macroTwoLines(a, b []FxPoint, lang string, log bool) MacroChart {
	return macroTwoLinesWith(a, b, lang, log, fxChartOpts{}, fxChartOpts{})
}

// macroTwoLinesWith is the same with labels for the readout.
func macroTwoLinesWith(a, b []FxPoint, lang string, log bool, oa, ob fxChartOpts) MacroChart {
	const w, h = 1000.0, 360.0
	if len(a) < 2 || len(b) < 2 {
		return MacroChart{}
	}
	grid, av, bv, ah, bh := macroGrid(a, b)
	grid, av, bv, ah, bh = macroThinGrid(grid, av, bv, ah, bh, fxMaxPoints)
	if len(grid) < 2 {
		return MacroChart{}
	}

	lo, hi := math.Inf(1), math.Inf(-1)
	for i := range grid {
		for _, set := range []struct {
			v  float64
			ok bool
		}{{av[i], ah[i]}, {bv[i], bh[i]}} {
			if !set.ok || (log && set.v <= 0) {
				continue
			}
			lo, hi = math.Min(lo, set.v), math.Max(hi, set.v)
		}
	}
	if math.IsInf(lo, 1) || math.IsInf(hi, -1) {
		return MacroChart{}
	}

	var ticks []FxTick
	var at func(float64) float64
	if log {
		loE, hiE := math.Floor(math.Log10(lo)), math.Ceil(math.Log10(hi))
		at = func(v float64) float64 {
			if v <= 0 {
				return h
			}
			return h - h*(math.Log10(v)-loE)/(hiE-loE)
		}
		// Below one, a whole-number label prints "0" — and a logarithmic axis
		// reading 100 / 10 / 1 / 0 / 0 / 0 is worse than no axis at all. Each
		// decade gets exactly the digits it needs, and a series that knows how
		// to name its own magnitude ("40,2 млрд") labels the axis that way
		// instead, so the ticks and the readouts speak the same units.
		label := func(v, e float64) string {
			if oa.AxisFormat != nil {
				return oa.AxisFormat(v)
			}
			if oa.Format != nil {
				return oa.Format(v)
			}
			digits := 0
			if e < 0 {
				digits = int(-e)
			}
			return fxFormat(v, digits)
		}
		for e := hiE; e >= loE-0.5; e-- {
			v := math.Pow(10, e)
			ticks = append(ticks, FxTick{Label: label(v, e), Pos: (hiE - e) / (hiE - loE) * 100})
		}
	} else {
		var step float64
		lo, hi, step = fxScale(lo, hi)
		at = func(v float64) float64 { return h - h*(v-lo)/(hi-lo) }
		for v := hi; v >= lo-step/2; v -= step {
			ticks = append(ticks, FxTick{Label: fxAxisNum(v, step), Pos: (hi - v) / (hi - lo) * 100})
		}
	}

	pa := macroPath(av, ah, at, w)
	// No fill under the first line on a logarithmic scale: the area under a
	// logarithm means nothing, while the eye reads it as a volume.
	area := ""
	if !log && pa != "" {
		area = pa + fmt.Sprintf(" L%.1f %.1f L0 %.1f Z", w, h, h)
	}

	mode := oa.label("all")
	if ob.Annual {
		mode = fxLabelYear
	}
	hover := FxHover{L: fxPointLabels(grid, mode, lang)}
	for _, set := range []struct {
		vals []float64
		has  []bool
		o    fxChartOpts
	}{{av, ah, oa}, {bv, bh, ob}} {
		format := set.o.Format
		if format == nil {
			format = fxNum
		}
		ser := FxHoverSeries{N: set.o.Name}
		for i, v := range set.vals {
			if !set.has[i] {
				ser.Y = append(ser.Y, nil)
				ser.V = append(ser.V, "")
				continue
			}
			ser.Y = append(ser.Y, hoverPct(at(v)/h*100))
			ser.V = append(ser.V, fxUnit(format(v), set.o.Unit))
		}
		hover.S = append(hover.S, ser)
	}

	minWidth := 0
	if oa.PxPerPoint > 0 {
		minWidth = len(grid) * oa.PxPerPoint
	}
	return MacroChart{
		A: pa, B: macroPath(bv, bh, at, w),
		Area:  area,
		Y:     ticks,
		X:     fxAxis(grid, "all", lang),
		Width: int(w), Height: int(h),
		Log:      log,
		Hover:    hoverJSON(hover),
		MinWidth: minWidth,
	}
}

// macroYears brings inflation, money growth and the change in the rate into one
// table by year. Three figures in a row answer the question people come here for:
// what was happening to money, and what to prices.
func macroYears(cpi, m3 []MacroPoint, rate []FxPoint, fund []MacroPoint) []MacroYear {
	byYear := map[int]*MacroYear{}
	get := func(y int) *MacroYear {
		if _, ok := byYear[y]; !ok {
			byYear[y] = &MacroYear{Year: y}
		}
		return byYear[y]
	}
	for _, c := range cpi {
		get(c.Period.Year()).CPI = fxPct(c.Value)
	}
	for y, d := range macroYearChange(macroPoints(m3)) {
		get(y).Money = fxPct(d)
	}
	for y, d := range macroYearChange(rate) {
		get(y).Rate = fxPct(d)
	}
	// The National Fund is the one row where a fall matters more than a rise: a
	// decrease means more was taken out of it than was put in.
	for y, d := range macroYearChange(macroPoints(fund)) {
		r := get(y)
		r.Fund, r.FundDown = fxPct(d), d < 0
	}

	out := make([]MacroYear, 0, len(byYear))
	for _, v := range byYear {
		if v.CPI == "" && v.Money == "" && v.Rate == "" && v.Fund == "" {
			continue
		}
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Year > out[j].Year })
	return out
}

// macroYearChange computes a series' December-to-December change, in percent.
func macroYearChange(pts []FxPoint) map[int]float64 {
	last := map[int]float64{}
	years := []int{}
	for _, p := range pts {
		y := p.Day.Year()
		if _, ok := last[y]; !ok {
			years = append(years, y)
		}
		last[y] = p.Value
	}
	out := map[int]float64{}
	for i := 1; i < len(years); i++ {
		prev := last[years[i-1]]
		if prev <= 0 {
			continue
		}
		out[years[i]] = (last[years[i]]/prev - 1) * 100
	}
	return out
}

// macroTenge prints a tenge sum that arrives in millions, choosing the magnitude
// word that keeps the figure short.
//
// The word is translated. It used to be hard-coded Russian, which put "55,8 трлн"
// into the middle of an English sentence on the English page.
// macroTengeRound is macroTenge without the fraction, for axis ticks: the tick
// stands on an exact power of ten, and "100,0 трлн" spends three characters
// saying nothing while the label runs out of gutter and gets clipped.
func macroTengeRound(millions float64, lang string) string {
	switch {
	case millions >= 1e6:
		return fxFormat(millions/1e6, 0) + " " + T(lang, "fx.mag_trn")
	case millions >= 1e3:
		return fxFormat(millions/1e3, 0) + " " + T(lang, "fx.mag_bln")
	default:
		return fxFormat(millions, 0) + " " + T(lang, "fx.mag_mln")
	}
}

func macroTenge(millions float64, lang string) string {
	switch {
	case millions >= 1e6:
		return fxFormat(millions/1e6, 1) + " " + T(lang, "fx.mag_trn")
	case millions >= 1e3:
		return fxFormat(millions/1e3, 1) + " " + T(lang, "fx.mag_bln")
	default:
		return fxFormat(millions, 0) + " " + T(lang, "fx.mag_mln")
	}
}

// macroDollars prints a dollar sum that arrives in millions.
func macroDollars(millions float64, lang string) string {
	if millions >= 1e3 {
		return fxFormat(millions/1e3, 1) + " " + T(lang, "fx.mag_bln")
	}
	return fxFormat(millions, 0) + " " + T(lang, "fx.mag_mln")
}

// macroMonth prints a month by name with its year. An abbreviation suits a label
// under a chart but not a sentence: "as of Jul 2026" reads like a typo.
func macroMonth(t time.Time, lang string) string {
	if t.IsZero() {
		return ""
	}
	return fxMonthFull(t.Month(), lang) + " " + fmt.Sprintf("%d", t.Year())
}

// macroMonthIn prints a month in the form needed after the preposition "in". In
// Russian that is the prepositional case: "в январе", not "в январь". The other
// two languages need no declension here.
func macroMonthIn(t time.Time, lang string) string {
	if t.IsZero() {
		return ""
	}
	if lang != LangRU {
		return macroMonth(t, lang)
	}
	prep := [12]string{"январе", "феврале", "марте", "апреле", "мае", "июне",
		"июле", "августе", "сентябре", "октябре", "ноябре", "декабре"}
	return prep[int(t.Month())-1] + " " + fmt.Sprintf("%d", t.Year())
}

// fxMonthFull is a month's full name in the reader's language.
func fxMonthFull(m time.Month, lang string) string {
	names := map[string][12]string{
		LangKZ: {"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым",
			"шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан"},
		LangRU: {"январь", "февраль", "март", "апрель", "май", "июнь",
			"июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"},
		LangEN: {"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December"},
	}
	row, ok := names[lang]
	if !ok {
		row = names[LangRU]
	}
	return row[int(m)-1]
}

// Caching the section.
//
// The section is identical for all forty-eight currencies and every period, and
// its data changes once a month. Building it on every visit would mean six queries
// and five charts for a picture that is the same either way.

// macroTTL is how long an assembled section lives.
const macroTTL = time.Hour

// macroCache holds the assembled section per language.
var macroCache = struct {
	mu   sync.Mutex
	at   time.Time
	byLg map[string]MacroBlock
}{byLg: map[string]MacroBlock{}}

// macroCached serves the section, rebuilding it at most once an hour.
func (m *Module) macroCached(ctx context.Context, lang string) MacroBlock {
	macroCache.mu.Lock()
	fresh := time.Since(macroCache.at) < macroTTL
	got, ok := macroCache.byLg[lang]
	macroCache.mu.Unlock()
	if fresh && ok {
		return got
	}

	built := m.buildMacro(ctx, lang)
	if !built.ready() {
		// An incomplete section is not cached. Right after a deployment the
		// series arrive one after another — money supply, then reserves, then
		// inflation — and caching the state between them means showing half an
		// analysis for an hour while the other half is already in the database.
		return built
	}
	macroCache.mu.Lock()
	if !fresh {
		// A stale set is dropped whole: keeping one language from the old hour
		// beside another from the new one means showing two different months on
		// the same page.
		macroCache.byLg = map[string]MacroBlock{}
		macroCache.at = time.Now()
	}
	macroCache.byLg[lang] = built
	macroCache.mu.Unlock()
	return built
}

// macroTargetRecord counts the Bank's record against its own inflation target,
// straight off the series the chart draws.
//
// Counted rather than written down, because a verdict printed under a picture
// has to be the picture's own arithmetic. A number typed into a translation
// string is right on the day it is typed and wrong from the next revision of
// the data onwards, and nobody notices until a reader does.
func macroTargetRecord(cpi []FxPoint, target float64) (years, missed, avg10 string, allMissed bool) {
	if len(cpi) == 0 || target <= 0 {
		return "", "", "", false
	}
	n, over := 0, 0
	for _, p := range cpi {
		n++
		if p.Value > target {
			over++
		}
	}
	// The last ten points, or all of them when the series is shorter.
	tail := cpi
	if len(tail) > 10 {
		tail = tail[len(tail)-10:]
	}
	sum := 0.0
	for _, p := range tail {
		sum += p.Value
	}
	return fmt.Sprintf("%d", n), fmt.Sprintf("%d", over),
		fxFormat(sum/float64(len(tail)), 1) + " %", over == n
}

// macroYearsWord puts the Russian word for "year" into the form the number in
// front of it requires: 32 года, 35 лет, 31 год. Without it the verdict reads
// "за все 32 лет", and a sentence that stumbles is a sentence the reader stops
// trusting. Kazakh and English need no agreement here.
func macroYearsWord(n int64, lang string) string {
	switch lang {
	case LangKZ:
		return "жыл"
	case LangEN:
		if n == 1 {
			return "year"
		}
		return "years"
	}
	if n < 0 {
		n = -n
	}
	if t := n % 100; t >= 11 && t <= 14 {
		return "лет"
	}
	switch n % 10 {
	case 1:
		return "год"
	case 2, 3, 4:
		return "года"
	}
	return "лет"
}

// macroPctTrim prints a percentage without a fraction when it has none: an
// inflation target of "5,00 %" is two characters of false precision on a figure
// that was announced as five.
func macroPctTrim(v float64) string {
	if v == math.Trunc(v) {
		return fxFormat(v, 0) + " %"
	}
	return fxFormat(v, 2) + " %"
}

// cpiTargetValue is the target the Bank publishes, or five percent when the
// panel has not been read yet — the figure it has announced for years.
func cpiTargetValue(m *Module, ctx context.Context) float64 {
	if m.macro == nil {
		return 5
	}
	pts, err := m.macro.Series(ctx, MacroCPITarget)
	if err != nil || len(pts) == 0 {
		return 5
	}
	return pts[len(pts)-1].Value
}

// gdpSeries loads real output growth, or nothing when it has not been fetched.
func gdpSeries(m *Module, ctx context.Context) []MacroPoint {
	if m.macro == nil {
		return nil
	}
	pts, err := m.macro.Series(ctx, MacroGDP)
	if err != nil {
		return nil
	}
	return pts
}

// macroMoneyVsOutput compares how much the money supply grew with how much the
// country's production grew, over the years both are known for.
//
// This is the page's plainest statement and the one that survives every
// argument about which line on a chart moves first. Prices are what money buys;
// when money multiplies far faster than the things it can buy, the difference
// comes out in prices. Both figures are published by the institutions
// themselves, and the arithmetic is a division.
func macroMoneyVsOutput(m3 []MacroPoint, gdp []MacroPoint, lang string) (gdpMul, moneyMul, gap string) {
	if len(m3) < 2 || len(gdp) < 5 {
		return "", "", ""
	}
	// Both spans must start in the same year, or the comparison is between
	// different stretches of history. Output is published from 1991 and the
	// money supply only from 1994, and the three missing years are the deepest
	// contraction the country has had — counting them against money that was
	// not yet measured turned a 3.6-fold rise in output into 2.7.
	//
	// This is the second time the page has been caught making exactly this
	// mistake; the first was the growth multiples under the comparison chart.
	firstYear := gdp[0].Period.Year()
	if m3First := m3[0].Period.Year(); m3First > firstYear {
		firstYear = m3First
	}
	out := 1.0
	years := 0
	for _, p := range gdp {
		if p.Period.Year() < firstYear {
			continue
		}
		out *= 1 + p.Value/100
		years++
	}
	if out <= 0 || years < 5 {
		return "", "", ""
	}
	// Money: from the first month of the same year the output series starts, so
	// the two spans match. Comparing a money multiple since 1994 with output
	// since 1998 would be the same mismatch this page has already been caught
	// making once.
	var start, end float64
	for _, p := range m3 {
		if p.Period.Year() < firstYear {
			continue
		}
		if start == 0 {
			start = p.Value
		}
		end = p.Value
	}
	if start <= 0 || end <= 0 {
		return "", "", ""
	}
	money := end / start
	if money <= 0 {
		return "", "", ""
	}
	return macroMul(out, lang), macroMul(money, lang), macroMul(money/out, lang)
}

// macroDearYears counts the years borrowing cost more than prices rose — the
// rate above inflation.
//
// This is the plainest measure of what credit cost the people who might have
// built something with it. It says nothing about intent and needs no theory:
// either the money was dearer than the inflation it was meant to outrun, or it
// was not.
func macroDearYears(rate, cpi []FxPoint) string {
	if len(rate) == 0 || len(rate) != len(cpi) {
		return ""
	}
	n := 0
	for i := range rate {
		if rate[i].Value > cpi[i].Value {
			n++
		}
	}
	return fmt.Sprintf("%d", n)
}
