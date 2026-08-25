package articles

import (
	"strings"
	"testing"
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
