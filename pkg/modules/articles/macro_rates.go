package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// The National Bank's policy rate, and the live indicators from its front page.
//
// Money supply and reserves explain what the exchange rate is made of. They do
// not explain who decided anything, or when. Monetary policy has exactly one
// decision and it is announced out loud: every six weeks or so the National Bank
// names the rate it lends to banks at, and publishes the date, the size and the
// corridor. It is the only figure on this whole page with an author, a signature
// and a minute of the meeting — and without it the analysis shows consequences
// while staying silent about the cause.
//
// There is no machine access to the rate: no feed, no export. What there is, is
// two tables sitting in the HTML, and they have held unchanged for years.

const (
	// MacroBaseRate is the base rate, percent per annum. Announced since 2015 and
	// the main instrument of policy ever since.
	MacroBaseRate = "base_rate"
	// MacroRefiRate is the refinancing rate, percent per annum. The series runs
	// from 1992; on 27 October 2020 it was set equal to the base rate.
	MacroRefiRate = "refi_rate"
	// MacroCPINow is annual inflation as the National Bank shows it today. Not
	// the same as MacroCPI: that one is annual, comes from the World Bank and
	// lags by nearly a year; this one is the latest published month.
	MacroCPINow = "cpi_now"
	// MacroCPITarget is the inflation target the National Bank announces.
	MacroCPITarget = "cpi_target"
	// MacroTonia is TONIA, the overnight repo rate on the exchange. It shows
	// whether the announced rate is what actually holds on the market.
	MacroTonia = "tonia"
	// MacroGDP is real GDP growth, percent per year.
	MacroGDP = "gdp_growth"
)

const (
	nbkRefiURL = "https://nationalbank.kz/ru/page/stavka-refinansirovaniya-nbk"
	nbkBaseURL = "https://nationalbank.kz/ru/news/grafik-prinyatiya-resheniy-po-bazovoy-stavke"
	nbkHomeURL = "https://nationalbank.kz/ru"
	// Real GDP growth is needed because the equation of exchange subtracts it
	// from money growth. The National Bank does not publish it — it is not its
	// indicator.
	worldBankGDPURL = "https://api.worldbank.org/v2/country/KAZ/indicator/" +
		"NY.GDP.MKTP.KD.ZG?format=json&per_page=200"
)

// nbkBaseRubricMax caps how many yearly base-rate archives to walk. Decisions
// are published one archive per year since 2015; the cap is there so that a
// change in the markup cannot turn the walk into an endless one.
const nbkBaseRubricMax = 20

// htmlTables pulls every table out of a document as rows of finished cells.
//
// Parsing HTML with expressions would be shorter right up to the first day the
// National Bank moves an attribute or adds a line break inside a cell. Parsing
// the tree properly costs one import and survives such edits silently.
func htmlTables(data []byte) [][][]string {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil
	}
	var tables [][][]string
	var walkTable func(*html.Node) [][]string
	var walk func(*html.Node)

	// text collects a cell's whole content into one string. National Bank cells
	// contain <br>, <p> and links, and those have to be joined without a
	// separator: "16,75" must not become "16, 75".
	var text func(*html.Node, *strings.Builder)
	text = func(n *html.Node, b *strings.Builder) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			text(c, b)
		}
	}

	walkTable = func(t *html.Node) [][]string {
		var rows [][]string
		var rec func(*html.Node)
		rec = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "tr" {
				var cells []string
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.ElementNode || (c.Data != "td" && c.Data != "th") {
						continue
					}
					var b strings.Builder
					text(c, &b)
					cells = append(cells, macroCellText(b.String()))
				}
				if len(cells) > 0 {
					rows = append(rows, cells)
				}
				return
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
		}
		rec(t)
		return rows
	}

	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			if rows := walkTable(n); len(rows) > 0 {
				tables = append(tables, rows)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return tables
}

// macroCellText reduces a cell's content to one line, free of non-breaking
// spaces and stray wrapping.
func macroCellText(s string) string {
	s = strings.ReplaceAll(s, " ", " ")
	s = strings.ReplaceAll(s, " ", " ")
	return strings.Join(strings.Fields(s), " ")
}

// macroDayCell reads a date like "27.10.2020" from a cell, dropping footnote
// marks.
func macroDayCell(s string) (time.Time, bool) {
	s = strings.TrimSpace(strings.TrimRight(macroCellText(s), "*†‡ "))
	d, err := time.Parse("02.01.2006", s)
	if err != nil {
		return time.Time{}, false
	}
	return d, true
}

// macroPercentCell reads a percentage from a cell. The National Bank writes the
// decimal separator sometimes as a comma and sometimes as a point, occasionally
// with a percent sign and with spaces inside the number.
func macroPercentCell(s string) (float64, bool) {
	s = macroCellText(s)
	s = strings.NewReplacer("%", "", " ", "", ",", ".").Replace(s)
	s = strings.TrimRight(s, "*†‡")
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 || v > 1000 {
		return 0, false
	}
	return v, true
}

// parseNBKRefi parses the refinancing rate table, which starts in 1992.
//
// Columns: date, resolution, the single rate, and three term rates — six, three
// and one month. In 1995 there was no single rate, only term rates were
// announced; for those months the six-month rate is taken, being the longest of
// the ones named.
func parseNBKRefi(data []byte) ([]MacroPoint, error) {
	out := []MacroPoint{}
	for _, rows := range htmlTables(data) {
		for _, r := range rows {
			if len(r) < 3 {
				continue
			}
			day, ok := macroDayCell(r[0])
			if !ok {
				continue
			}
			v, ok := macroPercentCell(r[2])
			if !ok && len(r) > 3 {
				v, ok = macroPercentCell(r[3])
			}
			if !ok {
				continue
			}
			out = append(out, MacroPoint{Period: day, Value: v})
		}
	}
	if len(out) < 20 {
		return nil, fmt.Errorf("в таблице ставки рефинансирования прочиталось %d решений", len(out))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Period.Before(out[j].Period) })
	return out, nil
}

// parseNBKBase parses one yearly archive of base rate decisions.
//
// Columns: date, size, corridor, links to the press release. Future meetings sit
// in the same table with an empty size — those are skipped, or the chart would
// grow a zero where no decision has been taken yet.
func parseNBKBase(data []byte) []MacroPoint {
	out := []MacroPoint{}
	for _, rows := range htmlTables(data) {
		for _, r := range rows {
			if len(r) < 2 {
				continue
			}
			day, ok := macroDayCell(r[0])
			if !ok {
				continue
			}
			v, ok := macroPercentCell(r[1])
			if !ok {
				continue
			}
			out = append(out, MacroPoint{Period: day, Value: v})
		}
	}
	return out
}

// nbkRubricLinks collects the links to a section's yearly archives.
func nbkRubricLinks(data []byte, section string) []string {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	out := []string{}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				href := a.Val
				if !strings.Contains(href, section) || !strings.Contains(href, "/rubrics/") {
					continue
				}
				if strings.HasPrefix(href, "/") {
					href = "https://nationalbank.kz" + href
				}
				// The same section in Kazakh and English leads to the same
				// decisions: walking those would download everything three times.
				if !strings.Contains(href, "/ru/") || seen[href] {
					continue
				}
				seen[href] = true
				out = append(out, href)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	sort.Strings(out)
	if len(out) > nbkBaseRubricMax {
		out = out[:nbkBaseRubricMax]
	}
	return out
}

// nbkIndicator is one indicator from the National Bank's front page.
type nbkIndicator struct {
	Label string
	Value float64
}

// parseNBKIndicators reads the indicator panel on the National Bank's front
// page.
//
// Four figures, and they are the ones it presents itself as the current state of
// affairs: the inflation target, annual inflation, the base rate and TONIA. The
// panel's markup is a "number, label" pair, and the label is what identifies
// them: the order of the tiles changes, the labels hold.
func parseNBKIndicators(data []byte) []nbkIndicator {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil
	}
	out := []nbkIndicator{}
	var num string
	var text func(*html.Node) string
	text = func(n *html.Node) string {
		var b strings.Builder
		var rec func(*html.Node)
		rec = func(x *html.Node) {
			if x.Type == html.TextNode {
				b.WriteString(x.Data)
			}
			for c := x.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
		}
		rec(n)
		return macroCellText(b.String())
	}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch macroClass(n) {
			case "stats__number":
				num = text(n)
			case "stats__label":
				if v, ok := macroPercentCell(num); ok {
					out = append(out, nbkIndicator{Label: strings.ToLower(text(n)), Value: v})
				}
				num = ""
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return out
}

// macroClass returns a node's class if it is one of the two we know.
func macroClass(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key != "class" {
			continue
		}
		for _, c := range strings.Fields(a.Val) {
			switch c {
			case "stats__number", "stats__label":
				return c
			}
		}
	}
	return ""
}

// nbkIndicatorCode maps a tile's label onto our series code.
func nbkIndicatorCode(label string) string {
	switch {
	case strings.Contains(label, "цель"):
		return MacroCPITarget
	case strings.Contains(label, "инфляц"):
		return MacroCPINow
	case strings.Contains(label, "базовая ставка"):
		return MacroBaseRate
	case strings.Contains(label, "tonia"):
		return MacroTonia
	}
	return ""
}

// parseWorldBankGDP parses annual real GDP growth.
func parseWorldBankGDP(body []byte) ([]MacroPoint, error) {
	var doc []json.RawMessage
	if err := json.Unmarshal(body, &doc); err != nil || len(doc) < 2 {
		return nil, fmt.Errorf("ответ Всемирного банка не разобран")
	}
	var rows []struct {
		Date  string   `json:"date"`
		Value *float64 `json:"value"`
	}
	if err := json.Unmarshal(doc[1], &rows); err != nil {
		return nil, fmt.Errorf("ряд роста ВВП: %w", err)
	}
	out := make([]MacroPoint, 0, len(rows))
	for _, r := range rows {
		if r.Value == nil {
			continue
		}
		y, err := strconv.Atoi(r.Date)
		if err != nil || y < 1990 {
			continue
		}
		out = append(out, MacroPoint{
			Period: time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC),
			Value:  *r.Value,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Period.Before(out[j].Period) })
	if len(out) < 5 {
		return nil, fmt.Errorf("в ряду роста ВВП %d значений", len(out))
	}
	return out, nil
}

// refreshRates fetches the policy rate and the live indicators. It returns what
// failed: one silent source must not bring the others down.
func (m *Module) refreshRates(ctx context.Context, today time.Time) []string {
	var failed []string

	// The refinancing rate: one document, the whole history since 1992.
	if data, err := macroFetch(ctx, nbkRefiURL); err != nil {
		failed = append(failed, "ставка рефинансирования: "+err.Error())
	} else if pts, err := parseNBKRefi(data); err != nil {
		failed = append(failed, "ставка рефинансирования: "+err.Error())
	} else if err := m.macro.Save(ctx, MacroRefiRate, pts); err != nil {
		failed = append(failed, err.Error())
	}

	// The base rate: one archive per year, starting from 2015.
	if data, err := macroFetch(ctx, nbkBaseURL); err != nil {
		failed = append(failed, "базовая ставка: "+err.Error())
	} else {
		pts := parseNBKBase(data)
		for _, u := range nbkRubricLinks(data, "grafik-prinyatiya-resheniy-po-bazovoy-stavke") {
			year, err := macroFetch(ctx, u)
			if err != nil {
				continue
			}
			pts = append(pts, parseNBKBase(year)...)
		}
		if len(pts) < 10 {
			failed = append(failed, fmt.Sprintf("базовая ставка: прочиталось %d решений", len(pts)))
		} else if err := m.macro.Save(ctx, MacroBaseRate, macroDedupe(pts)); err != nil {
			failed = append(failed, err.Error())
		}
	}

	// The front-page indicator panel. The date is the day we read it: the
	// National Bank keeps the latest published value there and names no month
	// alongside it.
	if data, err := macroFetch(ctx, nbkHomeURL); err != nil {
		failed = append(failed, "показатели Нацбанка: "+err.Error())
	} else {
		got := 0
		for _, ind := range parseNBKIndicators(data) {
			code := nbkIndicatorCode(ind.Label)
			// The base rate arrives here too, but dated by the crawl instead of
			// by the decision — such a point does not belong in its series.
			if code == "" || code == MacroBaseRate {
				continue
			}
			got++
			if err := m.macro.Save(ctx, code, []MacroPoint{{Period: today, Value: ind.Value}}); err != nil {
				failed = append(failed, err.Error())
			}
		}
		if got == 0 {
			failed = append(failed, "показатели Нацбанка: панель не разобрана")
		}
	}

	if data, err := macroFetch(ctx, worldBankGDPURL); err != nil {
		failed = append(failed, "рост ВВП: "+err.Error())
	} else if pts, err := parseWorldBankGDP(data); err != nil {
		failed = append(failed, "рост ВВП: "+err.Error())
	} else if err := m.macro.Save(ctx, MacroGDP, pts); err != nil {
		failed = append(failed, err.Error())
	}

	return failed
}

// macroDedupe collapses repeats by date, keeping the last value. The yearly
// archives overlap: a January decision appears both in its own year and in the
// summary table on the section's front page.
func macroDedupe(pts []MacroPoint) []MacroPoint {
	seen := map[string]int{}
	out := []MacroPoint{}
	for _, p := range pts {
		k := p.Period.Format("2006-01-02")
		if i, ok := seen[k]; ok {
			out[i] = p
			continue
		}
		seen[k] = len(out)
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Period.Before(out[j].Period) })
	return out
}
