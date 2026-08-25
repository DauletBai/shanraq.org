package articles

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// The verdict under the rate chart is the picture's own arithmetic, so it has
// to be counted from the same series and never typed in. These guard the count
// and the sentence it lands in.
func TestTargetRecordCountsTheSeries(t *testing.T) {
	series := func(vs ...float64) []FxPoint {
		out := make([]FxPoint, 0, len(vs))
		for _, v := range vs {
			out = append(out, FxPoint{Value: v})
		}
		return out
	}
	years, missed, avg, all := macroTargetRecord(series(4, 6, 7), 5)
	if years != "3" || missed != "2" || all {
		t.Errorf("три года, два выше цели: получили %s/%s all=%v", years, missed, all)
	}
	if avg != "5,7 %" {
		t.Errorf("среднее = %q, ожидалось 5,7 %%", avg)
	}
	// Every year above the target is the case the page actually shows, and the
	// flag that selects the harder sentence.
	_, _, _, all = macroTargetRecord(series(6, 7, 8), 5)
	if !all {
		t.Error("все годы выше цели не распознаны")
	}
	// The average covers the last ten points, not the whole series: the verdict
	// speaks about the decade.
	long := make([]float64, 0, 20)
	for i := 0; i < 10; i++ {
		long = append(long, 100)
	}
	for i := 0; i < 10; i++ {
		long = append(long, 10)
	}
	if _, _, a, _ := macroTargetRecord(series(long...), 5); a != "10,0 %" {
		t.Errorf("среднее за десятилетие = %q, ожидалось 10,0 %%", a)
	}
	// An empty series must not produce a verdict at all.
	if y, _, _, _ := macroTargetRecord(nil, 5); y != "" {
		t.Errorf("пустой ряд дал вердикт %q", y)
	}
}

// "За все 32 лет" is a sentence that stumbles, and a sentence that stumbles is
// one the reader stops trusting.
func TestYearsAgreeWithTheirNumeral(t *testing.T) {
	ru := map[int64]string{
		1: "год", 2: "года", 4: "года", 5: "лет", 11: "лет", 14: "лет",
		21: "год", 22: "года", 25: "лет", 31: "год", 32: "года", 100: "лет",
	}
	for n, want := range ru {
		if got := macroYearsWord(n, LangRU); got != want {
			t.Errorf("%d %s — ожидалось %q", n, got, want)
		}
	}
	if macroYearsWord(1, LangEN) != "year" || macroYearsWord(2, LangEN) != "years" {
		t.Error("английское число не согласовано")
	}
	if macroYearsWord(32, LangKZ) != "жыл" {
		t.Error("казахский вариант потерян")
	}
}

// A target announced as five percent must not be printed as "5,00 %": two
// characters of precision the figure never had.
func TestTargetPrintsWithoutFalsePrecision(t *testing.T) {
	if got := macroPctTrim(5); got != "5 %" {
		t.Errorf("macroPctTrim(5) = %q, ожидалось \"5 %%\"", got)
	}
	if got := macroPctTrim(4.75); got != "4,75 %" {
		t.Errorf("macroPctTrim(4.75) = %q, ожидалось \"4,75 %%\"", got)
	}
}

// The verdict must not claim a direction the chart cannot show. An annual
// series settles co-movement and nothing more, and an earlier version of this
// text asserted a lag that the data does not carry.
func TestVerdictClaimsNoLeadOrLag(t *testing.T) {
	forbidden := map[string][]string{
		LangRU: {"идёт следом", "впереди неё", "это реакция"},
		LangEN: {"follows the inflation line", "rather than leading"},
	}
	for lang, phrases := range forbidden {
		text := T(lang, "fx.rate_verdict_1") + " " + T(lang, "fx.rate_verdict_2")
		for _, p := range phrases {
			if strings.Contains(text, p) {
				t.Errorf("вердикт (%s) снова утверждает направление: %q", lang, p)
			}
		}
	}
	// And it must still name who owns both quantities — that is the whole point.
	for _, lang := range []string{LangKZ, LangRU, LangEN} {
		for _, key := range []string{"fx.rate_verdict_1", "fx.rate_verdict_1p", "fx.rate_verdict_2"} {
			if strings.TrimSpace(T(lang, key)) == "" {
				t.Errorf("нет строки %q для языка %q", key, lang)
			}
		}
	}
}

// Which way the erosion figure runs.
//
// This is the regression test for a sentence that was published backwards. The
// column holds 1000 ÷ (accumulated price growth): what a thousand tenge of that
// year buys today, counted in that same year's money. With prices up forty-two
// times since 1995, a thousand tenge of 1995 buys today what twenty-four tenge
// bought then — it is NOT that a thousand tenge of 1995 equals twenty-four
// tenge of today, which would mean old money was worth less than new.
func TestErosionRunsFromTheOldYearForward(t *testing.T) {
	year := func(y int, v float64) MacroPoint {
		return MacroPoint{Period: mustYear(y), Value: v}
	}
	// Prices double every year for four years: growth ×8 from 1995 to 1998.
	cpi := []MacroPoint{year(1995, 100), year(1996, 100), year(1997, 100), year(1998, 100),
		year(1999, 0), year(2000, 0)}
	// The erosion table needs a monthly rate series of its own: it also answers
	// "how many dollars did that thousand fetch", and one point a year is not
	// enough for it to run at all.
	rate := []FxPoint{}
	for y := 1995; y <= 2000; y++ {
		for m := time.January; m <= time.December; m++ {
			rate = append(rate, FxPoint{
				Day:   time.Date(y, m, 1, 0, 0, 0, 0, time.UTC),
				Value: 100,
			})
		}
	}
	rows, _ := macroErosion(cpi, rate)
	if len(rows) == 0 {
		t.Fatal("таблица обесценения пуста")
	}
	first := rows[0]
	if first.Year != 1995 {
		t.Fatalf("первая строка = %d, ожидался 1995", first.Year)
	}
	// Prices double in each of 1995–1998 and hold still in 1999, so growth from
	// 1995 to the last year with a rate (2000) is 2⁴ = 16 and the thousand comes
	// to 62. The figure must be far below a thousand either way: money loses
	// value with time, it does not gain it.
	got := strings.ReplaceAll(strings.ReplaceAll(first.Kept, " ", ""), " ", "")
	n, err := strconv.Atoi(got)
	if err != nil {
		t.Fatalf("не число: %q", first.Kept)
	}
	if n >= 1000 {
		t.Errorf("осталось %d из тысячи — направление перевёрнуто: деньги не дорожают со временем", n)
	}
	if n != 62 {
		t.Errorf("осталось %d, ожидалось 62 (тысяча делится на рост цен в 16 раз)", n)
	}
}

func mustYear(y int) time.Time { return time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC) }

// A formula whose own numbers fail to multiply out is worse than no formula:
// the first reader with a calculator finds the page wrong about itself.
func TestPrintedIdentityMultipliesOut(t *testing.T) {
	// The real July 2026 figures, which is where this was caught: rounding each
	// factor on its own gave 55,8 = 16,7 × 3,33, and 16,7 × 3,33 is 55,61.
	const m3, base = 55_783_694.0, 16_749_995.0
	mult := macroShownRatio(m3, base)

	// Both sums as printed, in trillions.
	shownM3 := macroShownMillions(m3) / 1e6
	shownBase := macroShownMillions(base) / 1e6
	rounded := float64(int(mult*100+0.5)) / 100 // as printed, two decimals

	product := shownBase * rounded
	if diff := product - shownM3; diff > 0.05 || diff < -0.05 {
		t.Errorf("%.1f × %.2f = %.2f, а напечатано %.1f — формула не сходится",
			shownBase, rounded, product, shownM3)
	}

	// Across magnitudes too: the reserves line prints trillions over billions,
	// and 55,8 ÷ 63,7 has to come to what the page says it does.
	const res = 63_715.0
	k := macroShownRatio(m3, res)
	if want := 55.8e6 / 63.7e3; k < want-0.01 || k > want+0.01 {
		t.Errorf("K = %.2f, а 55,8 трлн ÷ 63,7 млрд = %.2f", k, want)
	}
	// A zero denominator must not produce an infinity on the page.
	if got := macroShownRatio(100, 0); got != 0 {
		t.Errorf("деление на ноль дало %v", got)
	}
}
