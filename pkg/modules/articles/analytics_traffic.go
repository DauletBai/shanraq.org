package articles

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// The traffic chart behind /analytics and the admin panel.
//
// Four figures over time, under two sets of switches: which audience to count
// and how finely to slice it. Everything is read from analytics_slots, so the
// definitions are the table's own -- a visit is a visitor inside a half-hour,
// a host is an address, a visitor is an address and a browser together.
//
// The chart is drawn on the server and the switches are ordinary links. That is
// not thrift for its own sake: this page is opened by advertisers deciding
// whether the audience is worth buying, and a figure that arrives only after a
// script has run is a figure some of them never see.

// trafficAudience is one of the four audience switches.
type trafficAudience struct {
	Code string // url value
	KZ   bool
	Mob  bool
}

var trafficAudiences = []trafficAudience{
	{Code: "all"},
	{Code: "kz", KZ: true},
	{Code: "mob", Mob: true},
	{Code: "mobkz", KZ: true, Mob: true},
}

// trafficPeriod is one of the four time slices, with the window it looks back
// over. The windows differ because a chart is read at a glance: two days of
// hours and two years of months both come to a readable number of points.
type trafficPeriod struct {
	Code   string
	Trunc  string // date_trunc field
	Window time.Duration
}

var trafficPeriods = []trafficPeriod{
	{Code: "hour", Trunc: "hour", Window: 48 * time.Hour},
	{Code: "day", Trunc: "day", Window: 30 * 24 * time.Hour},
	{Code: "week", Trunc: "week", Window: 26 * 7 * 24 * time.Hour},
	{Code: "month", Trunc: "month", Window: 730 * 24 * time.Hour},
}

func audienceByCode(c string) trafficAudience {
	for _, a := range trafficAudiences {
		if a.Code == c {
			return a
		}
	}
	return trafficAudiences[0]
}

func periodByCode(c string) trafficPeriod {
	for _, p := range trafficPeriods {
		if p.Code == c {
			return p
		}
	}
	return trafficPeriods[1] // days: the slice that reads well before any choice is made
}

// TrafficPoint is one column of the chart.
type TrafficPoint struct {
	Label    string
	Hosts    int64
	Visitors int64
	Visits   int64
	Views    int64
}

// TrafficChart is everything the template needs: the points, the scale they are
// drawn against, and which switches are currently on.
type TrafficChart struct {
	Points   []TrafficPoint
	Max      int64
	Audience string
	Period   string
	Empty    bool
}

// trafficChart reads one audience over one period.
func (m *Module) trafficChart(ctx context.Context, audience, period string) TrafficChart {
	a, p := audienceByCode(audience), periodByCode(period)
	out := TrafficChart{Audience: a.Code, Period: p.Code}

	// The audience switches are AND-ed onto the query rather than pre-summed
	// into four buckets, so "mobile in Kazakhstan" needs no table of its own.
	rows, err := m.rt.DB.Query(ctx, `
		SELECT date_trunc($1, slot) AS b,
		       COUNT(DISTINCT host) AS hosts,
		       COUNT(DISTINCT vid)  AS visitors,
		       COUNT(*)             AS visits,
		       COALESCE(SUM(views), 0) AS views
		  FROM analytics_slots
		 WHERE slot >= $2
		   AND ($3::bool IS NOT TRUE OR is_kz)
		   AND ($4::bool IS NOT TRUE OR is_mobile)
		 GROUP BY b
		 ORDER BY b`,
		p.Trunc, time.Now().UTC().Add(-p.Window), a.KZ, a.Mob)
	if err != nil {
		m.rt.Logger.Warn("traffic chart", zap.Error(err))
		out.Empty = true
		return out
	}
	defer rows.Close()

	layout := map[string]string{"hour": "15:04", "day": "02.01", "week": "02.01", "month": "01.2006"}[p.Code]
	for rows.Next() {
		var at time.Time
		var pt TrafficPoint
		if err := rows.Scan(&at, &pt.Hosts, &pt.Visitors, &pt.Visits, &pt.Views); err != nil {
			m.rt.Logger.Warn("traffic chart scan", zap.Error(err))
			break
		}
		pt.Label = at.Format(layout)
		for _, v := range []int64{pt.Hosts, pt.Visitors, pt.Visits, pt.Views} {
			if v > out.Max {
				out.Max = v
			}
		}
		out.Points = append(out.Points, pt)
	}
	out.Empty = len(out.Points) == 0
	return out
}

// trafficChartFrom reads the switches off the query string. Unknown values fall
// back rather than erroring: a hand-edited address should show a chart.
func (m *Module) trafficChartFrom(r *http.Request) TrafficChart {
	return m.trafficChart(r.Context(), r.URL.Query().Get("a"), r.URL.Query().Get("p"))
}

// Pct places one value on the chart's own scale, as a percentage of its tallest
// point. Zero max means an empty chart, which the template hides.
func (c TrafficChart) Pct(v int64) int {
	if c.Max <= 0 {
		return 0
	}
	p := int(v * 100 / c.Max)
	if p < 1 && v > 0 {
		return 1
	}
	return p
}

// Audiences and Periods drive the switch rows, each entry knowing whether it is
// the one currently on.
func (c TrafficChart) Audiences() []map[string]any {
	out := make([]map[string]any, 0, len(trafficAudiences))
	for _, a := range trafficAudiences {
		out = append(out, map[string]any{"Code": a.Code, "On": a.Code == c.Audience})
	}
	return out
}

func (c TrafficChart) Periods() []map[string]any {
	out := make([]map[string]any, 0, len(trafficPeriods))
	for _, p := range trafficPeriods {
		out = append(out, map[string]any{"Code": p.Code, "On": p.Code == c.Period})
	}
	return out
}

// Href builds the address of one switch, keeping whichever switch is not being
// changed.
func (c TrafficChart) Href(base, kind, code string) string {
	a, p := c.Audience, c.Period
	if kind == "a" {
		a = code
	} else {
		p = code
	}
	return fmt.Sprintf("%s?a=%s&p=%s", base, a, p)
}

// ---- drawing ----
//
// The geometry is computed here because a Go template cannot do arithmetic, and
// the alternative -- shipping the numbers and a script to plot them -- would put
// the figures behind JavaScript on the one page whose whole subject is figures.

// chartW and chartH are the viewBox. The SVG scales to its container, so these
// are proportions rather than pixels; the height leaves room under the plot for
// the labels the last point carries.
const (
	chartW = 1000
	chartH = 300
	padR   = 96 // room for the direct label at the end of each line
)

// TrafficSeries is one line, ready to draw.
type TrafficSeries struct {
	Key    string // i18n suffix: hosts | visitors | visits | views
	Slot   int    // colour slot, fixed: never cycled, never reassigned by rank
	Points string // SVG polyline points
	LastX  float64
	LastY  float64
	Last   int64
}

// Series returns the four lines in their fixed order. The order is the colour
// order too: a series keeps its hue whatever the switches are set to, so
// changing the audience never repaints the survivors.
func (c TrafficChart) Series() []TrafficSeries {
	pick := []struct {
		key string
		get func(TrafficPoint) int64
	}{
		{"hosts", func(p TrafficPoint) int64 { return p.Hosts }},
		{"visitors", func(p TrafficPoint) int64 { return p.Visitors }},
		{"visits", func(p TrafficPoint) int64 { return p.Visits }},
		{"views", func(p TrafficPoint) int64 { return p.Views }},
	}
	n := len(c.Points)
	if n == 0 || c.Max <= 0 {
		return nil
	}
	step := float64(chartW-padR) / float64(maxInt(n-1, 1))
	out := make([]TrafficSeries, 0, len(pick))
	for i, s := range pick {
		ser := TrafficSeries{Key: s.key, Slot: i}
		pts := ""
		for j, p := range c.Points {
			v := s.get(p)
			x := float64(j) * step
			y := float64(chartH) - float64(v)*float64(chartH)/float64(c.Max)
			if j > 0 {
				pts += " "
			}
			pts += fmt.Sprintf("%.1f,%.1f", x, y)
			ser.LastX, ser.LastY, ser.Last = x, y, v
		}
		ser.Points = pts
		out = append(out, ser)
	}
	return out
}

// Ticks are the horizontal grid lines and their values, four across the scale.
func (c TrafficChart) Ticks() []map[string]any {
	out := []map[string]any{}
	if c.Max <= 0 {
		return out
	}
	for i := 1; i <= 4; i++ {
		v := c.Max * int64(i) / 4
		out = append(out, map[string]any{
			"Y": float64(chartH) - float64(v)*float64(chartH)/float64(c.Max),
			"N": v,
		})
	}
	return out
}

// XLabels thins the axis to at most eight labels, so hours and months both fit
// without the text colliding.
func (c TrafficChart) XLabels() []map[string]any {
	n := len(c.Points)
	if n == 0 {
		return nil
	}
	every := maxInt(n/8, 1)
	step := float64(chartW-padR) / float64(maxInt(n-1, 1))
	out := []map[string]any{}
	for i, p := range c.Points {
		if i%every == 0 || i == n-1 {
			// Percent, not viewBox units: the SVG scales to its container, so a
			// label placed at a pixel offset drifts off its own gridline.
			out = append(out, map[string]any{"X": float64(i) * step * 100 / chartW, "L": p.Label})
		}
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
