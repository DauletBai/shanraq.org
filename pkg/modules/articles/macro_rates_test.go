package articles

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// nbkTable wraps rows in the markup the National Bank's own pages use, so the
// parsers are exercised against the shape they meet in production without a
// 700 KB page in the repository.
func nbkTable(rows [][]string) []byte {
	var b strings.Builder
	b.WriteString(`<html><body><div class="content"><table class="t">`)
	for _, r := range rows {
		b.WriteString("<tr>")
		for _, c := range r {
			// Cells arrive with non-breaking spaces and stray markup inside.
			b.WriteString("<td><p>" + c + "</p></td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString(`</table></div></body></html>`)
	return []byte(b.String())
}

func TestParseNBKRefi(t *testing.T) {
	rows := [][]string{
		{"Дата установления ставки (с)", "№, дата постановления", "Размер ставки, %"},
		{"единая", "сроком на:"},
		{"01.02.1992", "№3 от 14.02.92", "25", "", "", ""},
		{"01.07.1993", "№17 от 06.07.93", "110", "", "", ""},
		// 1995 announced only term rates; the six-month one stands in.
		{"01.01.1995", "№25 от 28.12.94", "", "210", "176", "152"},
		{"09.03.1995", "№5 от 6.03.95", "", "170", "144", "130"},
		{"04.02.2020", "", "9.25", "", "", ""},
		// A footnote mark on the date, and a comma decimal.
		{"27.10.2020**", "", "9,00", "", "", ""},
	}
	// Twenty decisions is the parser's own floor, so pad with plausible rows.
	for i := 0; i < 20; i++ {
		rows = append(rows, []string{"0" + string(rune('1'+i%9)) + ".05.2005", "", "8.00", "", "", ""})
	}
	pts, err := parseNBKRefi(nbkTable(rows))
	if err != nil {
		t.Fatalf("parseNBKRefi: %v", err)
	}
	byDay := map[string]float64{}
	for _, p := range pts {
		byDay[p.Period.Format("02.01.2006")] = p.Value
	}
	for day, want := range map[string]float64{
		"01.02.1992": 25, "01.07.1993": 110,
		// The single-rate column is empty here, so the six-month rate is taken.
		"01.01.1995": 210, "09.03.1995": 170,
		"04.02.2020": 9.25, "27.10.2020": 9,
	} {
		if got := byDay[day]; got != want {
			t.Errorf("%s = %v, want %v", day, got, want)
		}
	}
	// Ascending order is what every consumer assumes.
	for i := 1; i < len(pts); i++ {
		if pts[i].Period.Before(pts[i-1].Period) {
			t.Fatalf("series not ascending at %d", i)
		}
	}
}

func TestParseNBKRefiRefusesRubbish(t *testing.T) {
	// A page whose table lost its date column must fail loudly rather than
	// silently saving a three-point series over a good one.
	if _, err := parseNBKRefi(nbkTable([][]string{
		{"что-то", "9.00"}, {"ещё", "9.25"},
	})); err == nil {
		t.Fatal("expected an error on a table with no decisions")
	}
}

func TestParseNBKBase(t *testing.T) {
	pts := parseNBKBase(nbkTable([][]string{
		{"Дата установления базовой ставки", "Размер базовой ставки, %", "Коридор"},
		{"20.01.2025", "15,25", "14,25 - 16,25"},
		{"27.07.2026", "16,75", "15,75 - 17,75"},
		// A meeting that has not happened yet: no size, and no point on the chart.
		{"04.09.2026*", "", ""},
	}))
	if len(pts) != 2 {
		t.Fatalf("got %d decisions, want 2: %+v", len(pts), pts)
	}
	if pts[1].Value != 16.75 || pts[1].Period.Format("02.01.2006") != "27.07.2026" {
		t.Errorf("last decision = %+v", pts[1])
	}
}

func TestParseNBKIndicators(t *testing.T) {
	page := []byte(`<html><body><div class="main-banner-info__stats stats">
	  <a href="/x" class="stats__item"><p class="stats__number">5%</p><p class="stats__label">Цель по инфляции</p></a>
	  <div class="stats__item"><p class="stats__number">10.2%</p><p class="stats__label">Годовая инфляция</p></div>
	  <a class="stats__item"><p class="stats__number">16.75%</p><p class="stats__label">Базовая ставка</p></a>
	  <a class="stats__item"><p class="stats__number">16.04%</p><p class="stats__label">TONIA</p></a>
	  <a class="stats__item"><p class="stats__number">не число</p><p class="stats__label">Прочее</p></a>
	</div></body></html>`)
	got := map[string]float64{}
	for _, ind := range parseNBKIndicators(page) {
		if code := nbkIndicatorCode(ind.Label); code != "" {
			got[code] = ind.Value
		}
	}
	want := map[string]float64{
		MacroCPITarget: 5, MacroCPINow: 10.2, MacroBaseRate: 16.75, MacroTonia: 16.04,
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %v, want %v", k, got[k], v)
		}
	}
	if len(got) != len(want) {
		t.Errorf("read %d indicators, want %d: %v", len(got), len(want), got)
	}
}

func TestNBKRubricLinks(t *testing.T) {
	page := []byte(`<html><body>
	  <a href="/ru/news/grafik-prinyatiya-resheniy-po-bazovoy-stavke/rubrics/1543">2020</a>
	  <a href="/ru/news/grafik-prinyatiya-resheniy-po-bazovoy-stavke/rubrics/2365">2026</a>
	  <a href="/ru/news/grafik-prinyatiya-resheniy-po-bazovoy-stavke/rubrics/2365">2026 again</a>
	  <a href="https://nationalbank.kz/kz/news/grafik-prinyatiya-resheniy-po-bazovoy-stavke/rubrics/2365">kz</a>
	  <a href="/ru/news/something-else/rubrics/99">other section</a>
	  <a href="/ru/news/grafik-prinyatiya-resheniy-po-bazovoy-stavke">the section itself</a>
	</body></html>`)
	got := nbkRubricLinks(page, "grafik-prinyatiya-resheniy-po-bazovoy-stavke")
	if len(got) != 2 {
		t.Fatalf("got %d links, want 2: %v", len(got), got)
	}
	for _, u := range got {
		if !strings.HasPrefix(u, "https://nationalbank.kz/ru/") {
			t.Errorf("link not absolute or not Russian: %s", u)
		}
	}
}

func TestMacroPolicyRateSplice(t *testing.T) {
	day := func(s string) time.Time {
		d, err := time.Parse("02.01.2006", s)
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	refi := []MacroPoint{
		{Period: day("01.02.1992"), Value: 25},
		{Period: day("01.01.2000"), Value: 18},
		// These two mirror the base rate and must not be duplicated into it.
		{Period: day("04.02.2020"), Value: 9.25},
		{Period: day("27.10.2020"), Value: 9},
	}
	base := []MacroPoint{
		{Period: day("02.09.2015"), Value: 12},
		{Period: day("04.02.2020"), Value: 9.25},
		{Period: day("27.07.2026"), Value: 16.75},
	}
	got := macroPolicyRate(refi, base)
	if len(got) != 5 {
		t.Fatalf("got %d points, want 5: %+v", len(got), got)
	}
	if got[1].Day.Year() != 2000 || got[2].Day.Year() != 2015 {
		t.Errorf("splice landed wrong: %+v", got)
	}
	if last := got[len(got)-1]; last.Value != 16.75 {
		t.Errorf("last = %v, want 16.75", last.Value)
	}
}

func TestMacroRateByYearCarriesForward(t *testing.T) {
	day := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}
	rate := macroRateByYear([]FxPoint{
		{Day: day(2020, time.March, 10), Value: 12},
		{Day: day(2020, time.July, 21), Value: 9},
		// No decision in 2021 at all: the year inherits 2020's last one.
		{Day: day(2022, time.February, 24), Value: 13.5},
	})
	for y, want := range map[int]float64{2020: 9, 2021: 9, 2022: 13.5} {
		if rate[y] != want {
			t.Errorf("%d = %v, want %v", y, rate[y], want)
		}
	}
}

// TestMacroGridAlignsByDate is the regression test for the National Fund line.
// Reserves run from 1993 and the Fund from 2001; laid out by point number the
// Fund was drawn seven years before it was created.
func TestMacroGridAlignsByDate(t *testing.T) {
	month := func(y int, m time.Month) time.Time {
		return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	}
	var res, fund []FxPoint
	for y := 1994; y <= 2004; y++ {
		res = append(res, FxPoint{Day: month(y, time.January), Value: float64(y - 1990)})
		if y >= 2001 {
			fund = append(fund, FxPoint{Day: month(y, time.January), Value: float64(y - 2000)})
		}
	}
	grid, av, bv, ah, bh := macroGrid(res, fund)
	if len(grid) != 11 {
		t.Fatalf("grid has %d points, want 11", len(grid))
	}
	for i, p := range grid {
		wantFund := p.Day.Year() >= 2001
		if bh[i] != wantFund {
			t.Errorf("%d: fund present = %v, want %v", p.Day.Year(), bh[i], wantFund)
		}
		if !ah[i] {
			t.Errorf("%d: reserves missing", p.Day.Year())
		}
		if bh[i] && bv[i] != float64(p.Day.Year()-2000) {
			t.Errorf("%d: fund = %v", p.Day.Year(), bv[i])
		}
		if av[i] != float64(p.Day.Year()-1990) {
			t.Errorf("%d: reserves = %v", p.Day.Year(), av[i])
		}
	}

	// And the drawn path must start where the data starts, not at the left edge.
	at := func(v float64) float64 { return 360 - v }
	path := macroPath(bv, bh, at, 1000)
	if !strings.HasPrefix(path, "M") {
		t.Fatalf("path does not open with a move: %q", path)
	}
	var x float64
	if _, err := fmtSscan(path, &x); err != nil {
		t.Fatalf("cannot read the first x: %v", err)
	}
	// 2001 is the eighth of eleven points, so the line must begin well right of
	// the frame's left edge.
	if x < 600 {
		t.Errorf("fund line starts at x=%.0f, expected it past 600", x)
	}
}

// fmtSscan reads the x coordinate out of the opening "M<x> <y>" of a path.
func fmtSscan(path string, x *float64) (int, error) {
	return fmt.Sscanf(strings.TrimPrefix(path, "M"), "%f", x)
}
