package articles

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// book assembles an Excel workbook from sheet rows: a list of rows, each a list
// of cells. The real workbook from the source weighs eighty kilobytes, and there
// is no reason to keep it in the repository just to test the parsing.
func book(rows [][]string) []byte {
	var shared []string
	idx := map[string]int{}
	intern := func(s string) int {
		if i, ok := idx[s]; ok {
			return i
		}
		idx[s] = len(shared)
		shared = append(shared, s)
		return len(shared) - 1
	}

	var sheet strings.Builder
	sheet.WriteString(`<?xml version="1.0"?><worksheet><sheetData>`)
	for r, row := range rows {
		fmt.Fprintf(&sheet, `<row r="%d">`, r+1)
		for c, v := range row {
			ref := fmt.Sprintf("%s%d", string(rune('A'+c)), r+1)
			if v == "" {
				fmt.Fprintf(&sheet, `<c r="%s"/>`, ref)
				continue
			}
			if isNumber(v) {
				fmt.Fprintf(&sheet, `<c r="%s"><v>%s</v></c>`, ref, v)
				continue
			}
			fmt.Fprintf(&sheet, `<c r="%s" t="s"><v>%d</v></c>`, ref, intern(v))
		}
		sheet.WriteString(`</row>`)
	}
	sheet.WriteString(`</sheetData></worksheet>`)

	var ss strings.Builder
	ss.WriteString(`<?xml version="1.0"?><sst>`)
	for _, s := range shared {
		fmt.Fprintf(&ss, `<si><t>%s</t></si>`, s)
	}
	ss.WriteString(`</sst>`)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range map[string]string{
		"xl/worksheets/sheet1.xml": sheet.String(),
		"xl/sharedStrings.xml":     ss.String(),
	} {
		w, _ := zw.Create(name)
		_, _ = w.Write([]byte(body))
	}
	_ = zw.Close()
	return buf.Bytes()
}

func isNumber(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && r != '.' && r != '-' {
			return false
		}
	}
	return s != ""
}

// An empty cell is written as a self-closing tag, and it cannot be skipped: the
// match between columns and indicators rests on it.
func TestEmptyCellsKeepTheirPlace(t *testing.T) {
	sheet, err := readXLSX(book([][]string{{"A", "", "C"}}))
	if err != nil {
		t.Fatal(err)
	}
	if len(sheet) != 1 {
		t.Fatalf("прочитано %d строк", len(sheet))
	}
	if _, ok := sheet[0]["B"]; !ok {
		t.Error("пустая ячейка выпала, и колонка C заняла бы место B")
	}
	if sheet[0]["C"] != "C" {
		t.Errorf("колонка C прочиталась как %q", sheet[0]["C"])
	}
}

// The year in the source's labels is two digits, so the century has to be guessed.
// The boundary is drawn at the tenge's appearance.
func TestTwoDigitYearsLandInTheRightCentury(t *testing.T) {
	for in, want := range map[string]string{
		" 07.26": "2026-07",
		"12.93":  "1993-12",
		" 01.94": "1994-01",
		"12.99":  "1999-12",
	} {
		got, ok := parseNBKMonth(in)
		if !ok {
			t.Errorf("%q не разобрано", in)
			continue
		}
		if got.Format("2006-01") != want {
			t.Errorf("%q → %s, ожидалось %s", in, got.Format("2006-01"), want)
		}
	}
	if _, ok := parseNBKMonth("итого"); ok {
		t.Error("подпись «итого» принята за месяц")
	}
}

// In the monetary aggregates workbook months run across the columns and indicators
// down the rows, and the rows of interest are recognised by number plus name
// rather than by counting: the source has changed the row order, never the numbers.
func TestMoneyAggregatesAreFoundByName(t *testing.T) {
	data := book([][]string{
		{"Денежная база и агрегаты"},
		{"млн. теңге"},
		{"", " 01.94", " 02.94", " 03.94"},
		{"1. Денежная база (резервные деньги)", "5087", "6298", "23814"},
		{"Изменения за месяц в %", "1.2", "3.4", "5.6"},
		{"2. M0 (наличные деньги в обращении)", "2542", "2600", "2700"},
		{"5. M3 (денежная масса)", "8232", "9000", "10000"},
	})
	got, err := parseNBKMoney(data)
	if err == nil {
		t.Fatal("три месяца — слишком короткая книга, но она принята")
	}

	// The same workbook, but long enough.
	head := []string{""}
	base := []string{"1. Денежная база (резервные деньги)"}
	m0 := []string{"2. M0 (наличные деньги в обращении)"}
	m3 := []string{"5. M3 (денежная масса)"}
	pct := []string{"Изменения за месяц в %"}
	for i := 1; i <= 14; i++ {
		head = append(head, fmt.Sprintf(" %02d.94", i%12+1))
		base = append(base, fmt.Sprintf("%d", 5000+i))
		m0 = append(m0, fmt.Sprintf("%d", 2500+i))
		m3 = append(m3, fmt.Sprintf("%d", 8000+i))
		pct = append(pct, "1.1")
	}
	got, err = parseNBKMoney(book([][]string{{"Заголовок"}, {"млн"}, head, base, pct, m0, m3}))
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{MacroBase, MacroM0, MacroM3} {
		if len(got[code]) != 14 {
			t.Errorf("%s: %d точек вместо 14", code, len(got[code]))
		}
	}
	if got[MacroM3][0].Value != 8001 {
		t.Errorf("первое значение M3 — %.0f", got[MacroM3][0].Value)
	}
	// A percentages row must not stand in for the indicator itself.
	if got[MacroBase][0].Value == 1.1 {
		t.Error("вместо денежной базы прочитана строка процентов")
	}
}

// In the reserves workbook each indicator has three columns, and the National
// Fund's column has to be found by name. Counting positions makes net reserves the
// Fund — and those exist from 1993, whereas the Fund was created in 2000.
func TestTheNationalFundIsFoundByNameNotByPosition(t *testing.T) {
	rows := [][]string{
		{"Международные резервы"},
		{"", "Валовые международные резервы", "", "", "Активы в СКВ", "", "",
			"Монетарное золото", "", "", "Чистые международные резервы", "", "",
			"Валютные активы Национального фонда"},
		{"", "объем, млн. долларов", "изменение к предыдущему месяцу, %", "изменение к началу года, %",
			"объем, млн. долларов", "изменение к предыдущему месяцу, %", "изменение к началу года, %",
			"объем, млн. долларов", "изменение к предыдущему месяцу, %", "изменение к началу года, %",
			"объем, млн. долларов", "изменение к предыдущему месяцу, %", "изменение к началу года, %",
			"объем, млн. долларов"},
	}
	for i := 1; i <= 14; i++ {
		fund := ""
		if i > 7 {
			fund = fmt.Sprintf("%d", 600+i)
		}
		rows = append(rows, []string{
			fmt.Sprintf("%02d.94", i%12+1),
			fmt.Sprintf("%d", 700+i), "1.0", "1.0",
			fmt.Sprintf("%d", 400+i), "1.0", "1.0",
			fmt.Sprintf("%d", 200+i), "1.0", "1.0",
			fmt.Sprintf("%d", 690+i), "1.0", "1.0",
			fund,
		})
	}
	got, err := parseNBKReserves(book(rows))
	if err != nil {
		t.Fatal(err)
	}
	if n := len(got[MacroReserves]); n != 14 {
		t.Errorf("резервов %d точек вместо 14", n)
	}
	// The Fund is filled only from the second half of the series — as in life.
	if n := len(got[MacroFund]); n != 7 {
		t.Errorf("Нацфонд получил %d точек вместо семи: похоже, взята соседняя колонка", n)
	}
	if got[MacroReserves][0].Value != 701 {
		t.Errorf("первое значение резервов — %.0f", got[MacroReserves][0].Value)
	}
	if got[MacroFund][0].Value != 608 {
		t.Errorf("первое значение Нацфонда — %.0f", got[MacroFund][0].Value)
	}
}

// Years with no observation arrive as null, and they must not be taken for zero: a
// zero would mean prices did not rise at all that year.
func TestYearsWithoutAnObservationAreSkipped(t *testing.T) {
	body := []byte(`[{"page":1},[
		{"date":"2025","value":11.39},
		{"date":"2024","value":null},
		{"date":"2023","value":14.7}]]`)
	pts, err := parseWorldBankCPI(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(pts) != 2 {
		t.Fatalf("получено %d точек вместо двух", len(pts))
	}
	if pts[0].Period.Year() != 2023 || pts[1].Period.Year() != 2025 {
		t.Errorf("ряд не упорядочен по годам: %v", pts)
	}
}

// Cover is the ratio of two published figures, and the months have to match: July's
// money supply against July's reserves, not against whatever comes to hand.
func TestCoverMatchesMonthToMonth(t *testing.T) {
	july := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	june := july.AddDate(0, -1, 0)
	got := macroCover(
		[]MacroPoint{{Period: june, Value: 1000}, {Period: july, Value: 2000}},
		[]MacroPoint{{Period: july, Value: 4}},
	)
	if len(got) != 1 {
		t.Fatalf("получено %d точек: месяц без пары должен выпадать", len(got))
	}
	if got[0].Value != 500 {
		t.Errorf("покрытие посчитано как %.0f вместо 500", got[0].Value)
	}
}

// "в 97,2 раз" is a typo the reader trips over instead of looking at the number.
func TestTheRussianWordForTimesIsDeclined(t *testing.T) {
	for n, want := range map[int64]string{
		1: "раз", 2: "раза", 3: "раза", 4: "раза", 5: "раз",
		11: "раз", 12: "раз", 14: "раз", 22: "раза", 97: "раз", 6776: "раз",
	} {
		if got := ruTimesWord(n); got != want {
			t.Errorf("%d %s, а должно быть %s", n, got, want)
		}
	}
	if got := macroMul(97.2, LangRU); got != "97 раз" {
		t.Errorf("кратность напечатана как %q", got)
	}
	if got := macroMul(2.4, LangRU); got != "2,4 раза" {
		t.Errorf("дробная кратность напечатана как %q", got)
	}
	if got := macroMul(97.2, LangEN); got != "97 times" {
		t.Errorf("по-английски напечатано %q", got)
	}
}

// Erosion is computed from accumulated inflation, not from a single year.
func TestSavingsErosionCompoundsEveryYear(t *testing.T) {
	cpi := []MacroPoint{}
	for y := 2020; y <= 2025; y++ {
		cpi = append(cpi, MacroPoint{Period: time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100})
	}
	rate := []FxPoint{}
	for y := 2020; y <= 2025; y++ {
		rate = append(rate, FxPoint{Day: time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC), Value: 400})
		rate = append(rate, FxPoint{Day: time.Date(y, 12, 1, 0, 0, 0, 0, time.UTC), Value: 400})
	}
	got, now := macroErosion(cpi, rate)
	if len(got) != 1 || got[0].Year != 2020 {
		t.Fatalf("получено %v", got)
	}
	// Five years at a hundred percent is a division by thirty-two, not by six.
	if got[0].Kept != "31" {
		t.Errorf("покупательная способность посчитана как %q вместо 31", got[0].Kept)
	}
	if now != "2,50" {
		t.Errorf("сегодняшний эквивалент — %q", now)
	}
}
