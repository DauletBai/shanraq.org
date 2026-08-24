package articles

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Hover readouts for the line charts.
//
// A line answers "what shape"; it does not answer "how much, and when". Without
// figures under the cursor the reader is left estimating a value against the
// axis ticks — measuring by eye where we already hold the exact number.
//
// The readout data travels in an attribute on the frame rather than being
// computed in the browser: the values are already formatted by the same rules as
// the rest of the page, so a thousands separator in the tooltip matches the one
// in the table below it. Recomputing them client-side would mean a second set of
// rounding rules, and eventually a disagreement between the two.

// FxHoverSeries is one series inside a readout.
type FxHoverSeries struct {
	// N names the series. A lone line needs no name: it is the only one in frame.
	N string `json:"n,omitempty"`
	// Y is the point's position in the frame, percent from the top. A gap marks
	// a period the series has no value for.
	Y []*float64 `json:"y"`
	// V is the value in words, already formatted.
	V []string `json:"v"`
}

// FxHover is everything one frame's readout needs.
type FxHover struct {
	// L holds the time-axis labels, one per point.
	L []string        `json:"l"`
	S []FxHoverSeries `json:"s"`
}

// hoverJSON packs a readout for its attribute. An empty string means "no
// readout", and the frame then behaves exactly as it did before.
func hoverJSON(h FxHover) string {
	if len(h.L) == 0 || len(h.S) == 0 {
		return ""
	}
	b, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(b)
}

// hoverPct rounds a point's position to a tenth of a percent: a frame a thousand
// units wide resolves nothing finer, and every extra digit is an extra byte in
// each of the page's several thousand points.
func hoverPct(v float64) *float64 {
	r := math.Round(v*10) / 10
	return &r
}

// Label granularity for a readout. A point carries as much date as its series
// actually measured: no more, because "01.01.2024" credits an annual figure with
// a day it never had, and no less, because a daily rate needs its day.
const (
	fxLabelDay   = "day"
	fxLabelMonth = "month"
	fxLabelYear  = "year"
)

// fxPointLabels prints point dates the way a readout reads them — in full, not in
// the axis's abbreviation.
func fxPointLabels(pts []FxPoint, mode, lang string) []string {
	out := make([]string, 0, len(pts))
	for _, p := range pts {
		switch mode {
		case fxLabelYear:
			out = append(out, fmt.Sprintf("%d", p.Day.Year()))
		case fxLabelMonth:
			out = append(out, macroMonth(p.Day, lang))
		default:
			out = append(out, fxDate(p.Day))
		}
	}
	return out
}

// macroMonthKey is the month key two series are joined on.
func macroMonthKey(p FxPoint) string { return p.Day.Format("2006-01") }

// macroGrid lays two series on one shared time axis.
//
// They used to be laid out by point number rather than by date. While both
// series started in the same month and ran month for month that came to the same
// thing; for the reserves and the National Fund it does not. Reserves run from
// December 1993, the Fund was created in 2001, and a Fund line stretched across
// the full width showed it seven years before it existed.
//
// Returns the shared grid and each series' value on it. A missing value is
// flagged rather than zeroed: on a chart a zero is a point on the axis, not a
// gap.
func macroGrid(a, b []FxPoint) (grid []FxPoint, av, bv []float64, ah, bh []bool) {
	at := make(map[string]float64, len(a))
	bt := make(map[string]float64, len(b))
	seen := map[string]bool{}
	for _, p := range a {
		at[macroMonthKey(p)] = p.Value
	}
	for _, p := range b {
		bt[macroMonthKey(p)] = p.Value
	}
	// Merge in order rather than sorting afterwards: both series arrive
	// ascending, so one pass down the two of them is enough.
	merged := make([]FxPoint, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		switch {
		case j >= len(b) || (i < len(a) && !a[i].Day.After(b[j].Day)):
			merged = append(merged, a[i])
			i++
		default:
			merged = append(merged, b[j])
			j++
		}
	}
	for _, p := range merged {
		k := macroMonthKey(p)
		if seen[k] {
			continue
		}
		seen[k] = true
		grid = append(grid, p)
	}
	for _, p := range grid {
		k := macroMonthKey(p)
		x, okA := at[k]
		y, okB := bt[k]
		av = append(av, x)
		bv = append(bv, y)
		ah = append(ah, okA)
		bh = append(bh, okB)
	}
	return grid, av, bv, ah, bh
}

// macroThinGrid thins the shared grid together with both series, so the three
// stay index-aligned.
func macroThinGrid(grid []FxPoint, av, bv []float64, ah, bh []bool, max int) ([]FxPoint, []float64, []float64, []bool, []bool) {
	n := len(grid)
	if n <= max {
		return grid, av, bv, ah, bh
	}
	g := make([]FxPoint, 0, max)
	a := make([]float64, 0, max)
	b := make([]float64, 0, max)
	ga := make([]bool, 0, max)
	gb := make([]bool, 0, max)
	keep := func(i int) {
		g = append(g, grid[i])
		a = append(a, av[i])
		b = append(b, bv[i])
		ga = append(ga, ah[i])
		gb = append(gb, bh[i])
	}
	for i := 0; i < max-1; i++ {
		keep(i * (n - 1) / (max - 1))
	}
	keep(n - 1)
	return g, a, b, ga, gb
}

// macroPath draws a series along the shared grid, breaking the line wherever the
// series has no value.
func macroPath(vals []float64, has []bool, at func(float64) float64, w float64) string {
	if len(vals) < 2 {
		return ""
	}
	var s strings.Builder
	open := false
	for i, v := range vals {
		if !has[i] {
			open = false
			continue
		}
		x := w * float64(i) / float64(len(vals)-1)
		if !open {
			if s.Len() > 0 {
				s.WriteByte(' ')
			}
			s.WriteByte('M')
			open = true
		} else {
			s.WriteByte(' ')
			s.WriteByte('L')
		}
		writeCoord(&s, x, at(v))
	}
	return s.String()
}

// writeCoord prints a coordinate pair exactly as the old code did: one decimal
// place, decimal point.
func writeCoord(s *strings.Builder, x, y float64) {
	fmt.Fprintf(s, "%.1f %.1f", x, y)
}
