package articles

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// The rules by which the National Bank's decisions turn into inflation and the
// exchange rate.
//
// The series above show what happened. They do not show where it comes from, and
// so they read like weather: prices went up, the tenge got cheaper. A formula
// gives the analysis back its missing link — the rule by which one follows from
// the other — and substitutes into it not an example but the latest published
// figures.
//
// Three conditions hold everywhere below, or a formula is worth less than its
// absence. Every symbol is spelled out in words. Every figure is named together
// with its source and the month it came from. Every assumption the equality rests
// on is written beside it rather than implied: the equation of exchange becomes an
// estimate of inflation only if the velocity of money is unchanged, and passing
// over that in silence would present a model as a measurement.

// MacroTerm is one symbol in a formula: what it is, what it equals, where it
// came from.
type MacroTerm struct {
	Sym   string
	Name  string
	Value string
	Src   string
}

// MacroFormula is a formula with its figures filled in.
type MacroFormula struct {
	// ID picks the heading, the explanation and the conclusion out of the
	// translation dictionary.
	ID string
	// Symbols is the rule itself, in letters: "π ≈ ΔM − ΔQ".
	Symbols string
	// Filled is the same rule in figures.
	Filled string
	// Out is what came out.
	Out string
	// Note is the conclusion under the formula, figures already substituted.
	Note  string
	Terms []MacroTerm
}

// MacroNow is the indicator panel as the National Bank shows it on its own front
// page. It belongs here because these are the only figures in the analysis with
// an author: the Bank sets the rate, announces the target, publishes inflation.
type MacroNow struct {
	Day    string
	Base   string
	CPI    string
	Target string
	Tonia  string
	// Real is the base rate net of inflation, in points.
	Real    string
	RealPos bool
	// Gap is how far measured inflation sits from the announced target.
	Gap    string
	GapPos bool
	// Decision is the day the current rate was set.
	Decision string
}

// macroPct prints a quantity in percent without a sign: a rate and an inflation
// figure are not changes, and a plus in front of them would read as growth.
func macroPct(v float64) string { return fxFormat(v, 2) + " %" }

// macroPP prints the difference of two percentages, in points, with its sign.
func macroPP(v float64, lang string) string {
	return fxSign(v, fxFormat(v, 2)) + " " + T(lang, "fx.pp")
}

// macroPPAbs prints the same difference without a sign. Needed wherever the
// direction is already said in words: "the rate is −3.25 pp below inflation" is
// two negatives in a row, one of them redundant.
func macroPPAbs(v float64, lang string) string {
	return fxFormat(math.Abs(v), 2) + " " + T(lang, "fx.pp")
}

// macroLastPoint returns a series' last point.
func macroLastPoint(pts []MacroPoint) (MacroPoint, bool) {
	if len(pts) == 0 {
		return MacroPoint{}, false
	}
	return pts[len(pts)-1], true
}

// macroPolicyRate splices the National Bank's rate into a single series.
//
// Until September 2015 the instrument was the refinancing rate, after it the base
// rate; on 27 October 2020 the Bank set the first equal to the second, and where
// they overlap they match value for value. The splice is made at the date of the
// first base rate decision: before it the refinancing rate, from it on the base
// rate.
func macroPolicyRate(refi, base []MacroPoint) []FxPoint {
	if len(base) == 0 {
		return macroPoints(refi)
	}
	cut := base[0].Period
	out := make([]FxPoint, 0, len(refi)+len(base))
	for _, p := range refi {
		if p.Period.Before(cut) {
			out = append(out, FxPoint{Day: p.Period, Value: p.Value})
		}
	}
	for _, p := range base {
		out = append(out, FxPoint{Day: p.Period, Value: p.Value})
	}
	return out
}

// macroRateByYear reduces the policy rate to one figure per year: the rate
// actually in force through that year, weighted by the number of days it held.
//
// The obvious alternative — the rate standing on 31 December — is what this did
// first, and it quietly falsified the chart. Inflation for a year is what
// happened over the whole of it; a rate read on the last day is a snapshot. Put
// side by side, the two describe different things, and the mismatch invented a
// lead that was never in the data.
//
// 2015 is the case that exposed it. The tenge was floated that August, the Bank
// answered with a rate of sixteen percent in December, and the devaluation
// reached the shops over the following months. Read at year end the rate showed
// 16.0 against inflation of 6.7 — the rate apparently leaping a full year ahead
// of prices. Averaged over the days it actually held, 2015 was 8.65 against
// 6.68, and 2016 was 14.47 against 14.36: the two lines almost touching, which
// is what happened.
//
// A chart that makes a reader see a sequence the numbers do not contain is
// worse than no chart. This one now compares a year with the same year.
func macroRateByYear(pts []FxPoint) map[int]float64 {
	out := map[int]float64{}
	if len(pts) == 0 {
		return out
	}
	firstYear := pts[0].Day.Year()
	lastYear := pts[len(pts)-1].Day.Year()
	// A rate holds until the next decision, and the next decision may be a long
	// way off — so the series runs to the present, not to the last meeting.
	// Stopping at the last decision dropped whole years in which the Bank
	// simply left the rate alone, which is itself a decision.
	if now := siteNow().Year(); now > lastYear {
		lastYear = now
	}
	// The rate holds until the next decision, so a year with no decisions in it
	// is carried entirely by the last one taken before it.
	for y := firstYear; y <= lastYear; y++ {
		start := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(1, 0, 0)

		cur := 0.0
		haveCur := false
		for _, p := range pts {
			if !p.Day.After(start) {
				cur, haveCur = p.Value, true
			}
		}
		// Decisions taken inside the year, in order.
		type step struct {
			at time.Time
			v  float64
		}
		steps := []step{}
		if haveCur {
			steps = append(steps, step{start, cur})
		}
		for _, p := range pts {
			if p.Day.After(start) && p.Day.Before(end) {
				steps = append(steps, step{p.Day, p.Value})
			}
		}
		if len(steps) == 0 {
			continue
		}
		total, days := 0.0, 0.0
		for i, st := range steps {
			until := end
			if i+1 < len(steps) {
				until = steps[i+1].at
			}
			n := until.Sub(st.at).Hours() / 24
			if n <= 0 {
				continue
			}
			total += st.v * n
			days += n
		}
		if days > 0 {
			out[y] = total / days
		}
	}
	return out
}

// macroAlignYears reduces the rate and inflation to their common years, order
// preserved. Both must carry the same set of years, or the years drift apart and
// the chart starts lying.
func macroAlignYears(rate map[int]float64, cpi []MacroPoint) ([]FxPoint, []FxPoint) {
	var a, b []FxPoint
	for _, c := range cpi {
		y := c.Period.Year()
		r, ok := rate[y]
		if !ok || c.Value <= 0 || r <= 0 {
			continue
		}
		day := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
		a = append(a, FxPoint{Day: day, Value: r})
		b = append(b, FxPoint{Day: day, Value: c.Value})
	}
	return a, b
}

// macroSource assembles a source line: "National Bank · monetary aggregates ·
// July 2026".
func macroSource(lang, key string, when time.Time) string {
	s := T(lang, key)
	if when.IsZero() {
		return s
	}
	return s + " · " + macroMonth(when, lang)
}

// macroSourceYear is the same for an annual series. Real GDP growth has no
// month, and a "January 2025" label would credit an annual quantity with a
// precision it does not have.
func macroSourceYear(lang, key string, when time.Time) string {
	s := T(lang, key)
	if when.IsZero() {
		return s
	}
	return s + " · " + fmt.Sprintf("%d", when.Year())
}

// macroSourceDay is the same with an exact date: a rate decision has a day.
func macroSourceDay(lang, key string, when time.Time) string {
	s := T(lang, key)
	if when.IsZero() {
		return s
	}
	return s + " · " + fxDate(when)
}

// macroFormulaInput is everything the formulas are assembled from.
type macroFormulaInput struct {
	lang string

	m3        []MacroPoint
	base      []MacroPoint
	res       []MacroPoint
	gdp       []MacroPoint
	rate      []FxPoint
	cpiNow    []MacroPoint
	cpiTarget []MacroPoint
	fxRate    []FxPoint
}

// buildNow assembles today's indicator panel.
func buildNow(in macroFormulaInput, tonia []MacroPoint) MacroNow {
	var n MacroNow
	lang := in.lang
	if p, ok := macroLastPoint(in.cpiNow); ok {
		n.CPI, n.Day = macroPct(p.Value), fxDate(p.Period)
	}
	if p, ok := macroLastPoint(in.cpiTarget); ok {
		n.Target = macroPct(p.Value)
		if now, ok := macroLastPoint(in.cpiNow); ok {
			d := now.Value - p.Value
			n.Gap, n.GapPos = macroPP(d, lang), d > 0
		}
	}
	if p, ok := macroLastPoint(tonia); ok {
		n.Tonia = macroPct(p.Value)
	}
	if len(in.rate) > 0 {
		last := in.rate[len(in.rate)-1]
		n.Base, n.Decision = macroPct(last.Value), fxDate(last.Day)
		if p, ok := macroLastPoint(in.cpiNow); ok {
			d := last.Value - p.Value
			n.Real, n.RealPos = macroPP(d, lang), d > 0
		}
	}
	return n
}

// buildFormulas assembles the formulas it managed to fill with figures. A
// formula without a fresh number is not shown at all: on this page an example
// with invented quantities is worse than an empty space.
func buildFormulas(in macroFormulaInput) []MacroFormula {
	lang := in.lang
	out := []MacroFormula{}
	// The base comes first: it is the only step in the chain the National Bank
	// performs directly, and every rule after it operates on its result.
	if f, ok := formulaBase(in); ok {
		out = append(out, f)
	}
	if f, ok := formulaExchange(in); ok {
		out = append(out, f)
	}
	if f, ok := formulaReal(in); ok {
		out = append(out, f)
	}
	if f, ok := formulaCover(in); ok {
		out = append(out, f)
	}
	if f, ok := formulaTarget(in); ok {
		out = append(out, f)
	}
	_ = lang
	return out
}

// formulaBase splits the money supply into the part the National Bank issued
// itself and the part banks built on top of it.
//
// M3 = B × m is an identity, not a model: the multiplier is defined as their
// ratio, so nothing here is assumed. That is exactly why it belongs on this
// page — it names the author of each half of the money supply without a single
// contestable step. The base is an act of the National Bank. The multiplier is
// what banks do with that base, under the reserve rules and the rate the same
// Bank sets.
func formulaBase(in macroFormulaInput) (MacroFormula, bool) {
	lang := in.lang
	m3, ok := macroLastPoint(in.m3)
	if !ok {
		return MacroFormula{}, false
	}
	// Both figures must be for the same month: the two series are published
	// separately and can be a month apart.
	base, ok := macroAt(in.base, m3.Period)
	if !ok || base <= 0 {
		return MacroFormula{}, false
	}
	// The multiplier is derived from the figures as they are PRINTED, not from
	// the precise ones behind them. A formula whose own numbers fail to
	// multiply out is worse than no formula: 55,8 = 16,7 × 3,33 comes to 55,61,
	// and the first reader with a calculator finds the page wrong about itself.
	// Rounding each factor on its own put a hundred and ninety billion tenge
	// into that gap.
	mult := macroShownRatio(m3.Value, base)

	f := MacroFormula{
		ID:      "base",
		Symbols: "M3 = B × m",
		Filled: fmt.Sprintf("%s = %s × %s",
			macroTenge(m3.Value, lang), macroTenge(base, lang), fxFormat(mult, 2)),
		Out: macroTenge(base, lang),
		Terms: []MacroTerm{
			{Sym: "B", Name: T(lang, "fx.t_base"), Value: macroTenge(base, lang) + " ₸",
				Src: macroSource(lang, "fx.s_m3", m3.Period)},
			{Sym: "m", Name: T(lang, "fx.t_mult"), Value: fxFormat(mult, 2)},
			{Sym: "M3", Name: T(lang, "fx.t_m3"), Value: macroTenge(m3.Value, lang) + " ₸",
				Src: macroSource(lang, "fx.s_m3", m3.Period)},
		},
	}
	// The year's growth of each half turns the identity into an answer: it
	// shows which of the two actually moved.
	if pm, ok1 := macroAt(in.m3, m3.Period.AddDate(-1, 0, 0)); ok1 && pm > 0 {
		if pb, ok2 := macroAt(in.base, m3.Period.AddDate(-1, 0, 0)); ok2 && pb > 0 {
			f.Note = fmt.Sprintf(T(lang, "fx.f_base_cmp"),
				macroTenge(base, lang), macroTenge(m3.Value, lang),
				fxPct((base/pb-1)*100), fxPct((m3.Value/pm-1)*100))
		}
	}
	return f, true
}

// formulaExchange is the equation of exchange, out of which inflation comes.
//
// M·V = P·Q ties the quantity of money to the speed at which it circulates, the
// price level and the volume produced. If that speed does not change, the growth
// rates on both sides are equal, and price growth stays equal to money growth
// less output growth. That is the part of inflation with someone in charge of it:
// how much money the country holds is decided by whoever issues it.
func formulaExchange(in macroFormulaInput) (MacroFormula, bool) {
	lang := in.lang
	last, ok := macroLastPoint(in.m3)
	if !ok {
		return MacroFormula{}, false
	}
	prev, ok := macroAt(in.m3, last.Period.AddDate(-1, 0, 0))
	if !ok || prev <= 0 {
		return MacroFormula{}, false
	}
	dm := (last.Value/prev - 1) * 100

	g, ok := macroLastPoint(in.gdp)
	if !ok {
		return MacroFormula{}, false
	}
	dq := g.Value
	exp := dm - dq

	f := MacroFormula{
		ID:      "exch",
		Symbols: "M × V = P × Q   ⇒   π ≈ ΔM − ΔQ",
		Filled:  fmt.Sprintf("π ≈ %s − %s = %s", macroPct(dm), macroPct(dq), macroPct(exp)),
		Out:     macroPct(exp),
		Terms: []MacroTerm{
			{Sym: "ΔM", Name: T(lang, "fx.t_dm"), Value: macroPct(dm),
				Src: macroSource(lang, "fx.s_m3", last.Period)},
			{Sym: "ΔQ", Name: T(lang, "fx.t_dq"), Value: macroPct(dq),
				Src: macroSourceYear(lang, "fx.s_gdp", g.Period)},
			{Sym: "V", Name: T(lang, "fx.t_v"), Value: T(lang, "fx.t_v_val")},
			{Sym: "π", Name: T(lang, "fx.t_pi"), Value: macroPct(exp)},
		},
	}
	// The estimate always stands next to measured inflation: the discrepancy
	// between them is precisely the part the equation does not explain.
	if now, ok := macroLastPoint(in.cpiNow); ok {
		f.Note = fmt.Sprintf(T(lang, "fx.f_exch_cmp"), macroPct(exp), macroPct(now.Value),
			macroPPAbs(now.Value-exp, lang))
	}
	return f, true
}

// formulaReal is the real rate: what money costs once inflation is taken out.
func formulaReal(in macroFormulaInput) (MacroFormula, bool) {
	lang := in.lang
	if len(in.rate) == 0 {
		return MacroFormula{}, false
	}
	now, ok := macroLastPoint(in.cpiNow)
	if !ok {
		return MacroFormula{}, false
	}
	last := in.rate[len(in.rate)-1]
	r := last.Value - now.Value

	f := MacroFormula{
		ID:      "real",
		Symbols: "r = i − π",
		Filled:  fmt.Sprintf("r = %s − %s = %s", macroPct(last.Value), macroPct(now.Value), macroPP(r, lang)),
		Out:     macroPP(r, lang),
		Terms: []MacroTerm{
			{Sym: "i", Name: T(lang, "fx.t_i"), Value: macroPct(last.Value),
				Src: macroSourceDay(lang, "fx.s_base", last.Day)},
			{Sym: "π", Name: T(lang, "fx.t_cpi"), Value: macroPct(now.Value),
				Src: macroSourceDay(lang, "fx.s_panel", now.Period)},
			{Sym: "r", Name: T(lang, "fx.t_r"), Value: macroPP(r, lang)},
		},
	}
	key := "fx.f_real_neg"
	if r > 0 {
		key = "fx.f_real_pos"
	}
	f.Note = fmt.Sprintf(T(lang, key), macroPPAbs(r, lang))
	return f, true
}

// formulaCover is the rate implied by the quantity of tenge and the stock of
// currency.
//
// The ratio of money supply to reserves is neither a rate nor a forecast: it is
// the price at which all the country's tenge would clear against all the currency
// it holds. Beside the published rate it shows how far one is backed by the
// other.
func formulaCover(in macroFormulaInput) (MacroFormula, bool) {
	lang := in.lang
	m3, ok := macroLastPoint(in.m3)
	if !ok {
		return MacroFormula{}, false
	}
	res, ok2 := macroAt(in.res, m3.Period)
	if !ok2 || res <= 0 {
		r, ok3 := macroLastPoint(in.res)
		if !ok3 || r.Value <= 0 {
			return MacroFormula{}, false
		}
		res = r.Value
	}
	implied := m3.Value / res

	f := MacroFormula{
		ID:      "cover",
		Symbols: "K = M3 ÷ R",
		Filled: fmt.Sprintf("K = %s ÷ %s = %s ₸/$",
			macroTenge(m3.Value, lang), macroDollars(res, lang), fxNum(implied)),
		Out: fxNum(implied) + " ₸/$",
		Terms: []MacroTerm{
			{Sym: "M3", Name: T(lang, "fx.t_m3"), Value: macroTenge(m3.Value, lang) + " ₸",
				Src: macroSource(lang, "fx.s_m3", m3.Period)},
			{Sym: "R", Name: T(lang, "fx.t_res"), Value: macroDollars(res, lang) + " $",
				Src: macroSource(lang, "fx.s_res", m3.Period)},
			{Sym: "K", Name: T(lang, "fx.t_k"), Value: fxNum(implied) + " ₸/$"},
		},
	}
	if v, ok := macroRateAt(in.fxRate, m3.Period); ok && v > 0 {
		f.Note = fmt.Sprintf(T(lang, "fx.f_cover_cmp"), fxNum(implied), fxNum(v),
			macroMul(implied/v, lang))
		f.Terms = append(f.Terms, MacroTerm{
			Sym: "K₀", Name: T(lang, "fx.t_k0"), Value: fxNum(v) + " ₸/$",
			Src: macroSource(lang, "fx.s_fx", m3.Period),
		})
	}
	return f, true
}

// formulaTarget is the distance from the announced target to measured
// inflation.
func formulaTarget(in macroFormulaInput) (MacroFormula, bool) {
	lang := in.lang
	now, ok := macroLastPoint(in.cpiNow)
	if !ok {
		return MacroFormula{}, false
	}
	tgt, ok := macroLastPoint(in.cpiTarget)
	if !ok || tgt.Value <= 0 {
		return MacroFormula{}, false
	}
	gap := now.Value - tgt.Value

	f := MacroFormula{
		ID:      "target",
		Symbols: "Δ = π − π*",
		Filled: fmt.Sprintf("Δ = %s − %s = %s",
			macroPct(now.Value), macroPct(tgt.Value), macroPP(gap, lang)),
		Out: macroPP(gap, lang),
		Terms: []MacroTerm{
			{Sym: "π", Name: T(lang, "fx.t_cpi"), Value: macroPct(now.Value),
				Src: macroSourceDay(lang, "fx.s_panel", now.Period)},
			{Sym: "π*", Name: T(lang, "fx.t_target"), Value: macroPct(tgt.Value),
				Src: macroSourceDay(lang, "fx.s_panel", tgt.Period)},
			{Sym: "Δ", Name: T(lang, "fx.t_gap"), Value: macroPP(gap, lang)},
		},
	}
	key, times := "fx.f_target_over", macroMul(now.Value/tgt.Value, lang)
	if gap < 0 {
		key, times = "fx.f_target_under", macroMul(tgt.Value/now.Value, lang)
	}
	f.Note = fmt.Sprintf(T(lang, key), times)
	return f, true
}

// macroFormulaWidth is a housekeeping measure: a formula too long for its line
// breaks the card. Nothing is formatted here, but the length is what a test
// checks.
func macroFormulaWidth(f MacroFormula) int {
	return int(math.Max(float64(len([]rune(f.Symbols))), float64(len([]rune(strings.TrimSpace(f.Filled))))))
}

// macroShownRatio divides two sums the way the page shows them, so a printed
// identity holds when the reader checks it.
//
// Both figures reach the page through macroTenge, which rounds to one decimal
// inside its magnitude. When they land in the same magnitude the ratio is taken
// between those rounded values; when they do not — a base in billions under a
// mass in trillions — the exact ratio is the honest one, and no printed line
// invites the multiplication anyway.
func macroShownRatio(numMillions, denMillions float64) float64 {
	if denMillions <= 0 {
		return 0
	}
	nv, nu := macroShownValue(numMillions)
	dv, du := macroShownValue(denMillions)
	if nu == du && dv != 0 {
		return nv / dv
	}
	return numMillions / denMillions
}

// macroShownValue is the number macroTenge prints, with the magnitude it printed
// it in.
func macroShownValue(millions float64) (float64, string) {
	switch {
	case millions >= 1e6:
		return math.Round(millions/1e6*10) / 10, "trn"
	case millions >= 1e3:
		return math.Round(millions/1e3*10) / 10, "bln"
	default:
		return math.Round(millions), "mln"
	}
}
