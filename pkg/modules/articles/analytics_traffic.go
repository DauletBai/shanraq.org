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

// Each slice shows the whole of the next unit up: hours make a day, days make a
// month, months make a year. Weeks sat between the first two and told the same
// story as days with fewer points -- a slice that duplicates its neighbour is
// one more thing to click and nothing more to learn.
var trafficPeriods = []trafficPeriod{
	{Code: "hour", Trunc: "hour", Window: 24 * time.Hour},
	{Code: "day", Trunc: "day", Window: 30 * 24 * time.Hour},
	{Code: "month", Trunc: "month", Window: 365 * 24 * time.Hour},
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
	Label string
	// Known separates an hour with no visitors from an hour nobody counted. The
	// axis spans the whole day either way, but a line drawn through the second
	// kind would claim a zero we never measured -- and would dive to the floor
	// across every hour still to come.
	Known bool
	// Running marks the bucket still filling. The hour now in progress holds
	// five minutes of traffic against a full hour before it, and drawn as an
	// ordinary point it reads as a collapse rather than as a count so far.
	Running bool
	// Counted says hosts, visitors and visits were measured for this bucket.
	// Views reach back to the older counter, which never told readers apart, so
	// a day from before that counting began has a real view total and no
	// measurement at all of the other three.
	Counted  bool
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
	// Since is the day hosts, visitors and visits began being counted. Views
	// run further back, from the older counter, so without this line the early
	// part of the chart looks broken rather than merely older.
	Since string
	// Zone is the clock the figures are drawn under, named on the page: a chart
	// of hours that does not say whose hours is a chart two readers will read
	// two different ways.
	Zone string
	// Partial says the last point is a bucket still filling, so the page can
	// say so instead of leaving a cliff unexplained.
	Partial bool
}

// trafficChart reads one audience over one period.
func (m *Module) trafficChart(ctx context.Context, audience, period string, loc *time.Location) TrafficChart {
	a, p := audienceByCode(audience), periodByCode(period)
	out := TrafficChart{Audience: a.Code, Period: p.Code, Zone: loc.String()}

	// The audience switches are AND-ed onto the query rather than pre-summed
	// into four buckets, so "mobile in Kazakhstan" needs no table of its own.
	rows, err := m.rt.DB.Query(ctx, `
		SELECT date_trunc($1, slot AT TIME ZONE $5) AS b,
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
		p.Trunc, time.Now().In(loc).Add(-p.Window), a.KZ, a.Mob, loc.String())
	if err != nil {
		m.rt.Logger.Warn("traffic chart", zap.Error(err))
		out.Empty = true
		return out
	}
	defer rows.Close()

	layout := trafficLayout(p.Code)
	for rows.Next() {
		var at time.Time
		var pt TrafficPoint
		if err := rows.Scan(&at, &pt.Hosts, &pt.Visitors, &pt.Visits, &pt.Views); err != nil {
			m.rt.Logger.Warn("traffic chart scan", zap.Error(err))
			break
		}
		pt.Label, pt.Known, pt.Counted = at.Format(layout), true, true
		out.Points = append(out.Points, pt)
	}
	// Views reach further back than the other three. analytics_daily has counted
	// them since long before this table existed, and hosts, visitors and visits
	// cannot be recovered for those days because what they are computed from was
	// deliberately never recorded. So the older part of the chart carries the one
	// line that is real for it, and the other three begin where counting began --
	// which is honest in a way that back-filling them with zero would not be.
	m.backfillViews(ctx, a, p, &out, loc)

	var first time.Time
	if err := m.rt.DB.QueryRow(ctx, `SELECT MIN(slot) FROM analytics_slots`).Scan(&first); err == nil && !first.IsZero() {
		out.Since = first.In(loc).Format("02.01.2006 15:04")
		if p.Code == "hour" {
			fillHours(&out, first, loc)
		}
	}

	for _, pt := range out.Points {
		for _, v := range []int64{pt.Hosts, pt.Visitors, pt.Visits, pt.Views} {
			if v > out.Max {
				out.Max = v
			}
		}
	}
	out.Empty = len(out.Points) == 0
	return out
}

// backfillViews folds the older per-day view counts into the chart.
//
// The audience switches map onto counters analytics_daily already keeps: the
// page counter for everyone, the country counter for Kazakhstan, the device
// counter for mobile. It keeps no country-by-device cross, so "mobile in
// Kazakhstan" has no history and starts where the new table does.
func (m *Module) backfillViews(ctx context.Context, a trafficAudience, p trafficPeriod, out *TrafficChart, loc *time.Location) {
	if p.Code == "hour" {
		return // analytics_daily is keyed by date; it cannot answer an hour
	}
	kind, label := metricPage, ""
	switch {
	case a.KZ && a.Mob:
		return
	case a.KZ:
		kind, label = metricCountry, "KZ"
	case a.Mob:
		kind, label = metricDevice, "mobile"
	}
	rows, err := m.rt.DB.Query(ctx, `
		SELECT date_trunc($1, day::timestamp) AS b, COALESCE(SUM(n), 0)
		  FROM analytics_daily
		 WHERE day >= $2 AND kind = $3 AND ($4 = '' OR label = $4)
		 GROUP BY b ORDER BY b`,
		p.Trunc, time.Now().In(loc).Add(-p.Window), kind, label)
	if err != nil {
		m.rt.Logger.Warn("traffic backfill", zap.Error(err))
		return
	}
	defer rows.Close()

	layout := trafficLayout(p.Code)
	byLabel := map[string]int{}
	for i, pt := range out.Points {
		byLabel[pt.Label] = i
	}
	older := []TrafficPoint{}
	for rows.Next() {
		var at time.Time
		var n int64
		if err := rows.Scan(&at, &n); err != nil {
			break
		}
		lbl := at.Format(layout)
		if i, ok := byLabel[lbl]; ok {
			// Both sources have this bucket. The older counter is the one that
			// has been running longest, so its view count wins rather than being
			// added to a partial day from the new table.
			if n > out.Points[i].Views {
				out.Points[i].Views = n
			}
			continue
		}
		older = append(older, TrafficPoint{Label: lbl, Views: n, Known: true, Counted: false})
	}
	if len(older) > 0 {
		out.Points = append(older, out.Points...)
	}
}

// trafficLayout is the label format for a period, shared by both sources so a
// bucket lines up whichever table produced it.
func trafficLayout(code string) string {
	switch code {
	case "hour":
		return "15:00"
	case "month":
		return "01.2006"
	default:
		return "02.01"
	}
}

// trafficChartFrom reads the switches off the query string. Unknown values fall
// back rather than erroring: a hand-edited address should show a chart.
func (m *Module) trafficChartFrom(r *http.Request) TrafficChart {
	return m.trafficChart(r.Context(), r.URL.Query().Get("a"), r.URL.Query().Get("p"), readerLoc(r))
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
	Points string // SVG polyline over the settled buckets
	// Tail is the two-point segment into the bucket still filling, drawn dashed.
	Tail  string
	LastX float64
	LastY float64
	Last  int64
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
		pts, lastAt := "", ""
		for j, p := range c.Points {
			// A bucket holds its place on the axis whatever it knows. Views are
			// drawn wherever there is a total; the other three only where they
			// were actually measured, so the line begins at the counting rather
			// than climbing out of a zero nobody recorded.
			if !p.Known || (s.key != "views" && !p.Counted) {
				continue
			}
			v := s.get(p)
			x := float64(j) * step
			y := float64(chartH) - float64(v)*float64(chartH)/float64(c.Max)
			at := fmt.Sprintf("%.1f,%.1f", x, y)
			if p.Running {
				// The settled line ends at the previous point; this one hangs off
				// it on a dashed segment.
				ser.Tail = lastAt + " " + at
			} else {
				if pts != "" {
					pts += " "
				}
				pts += at
				lastAt = at
			}
			ser.LastX, ser.LastY, ser.Last = x, y, v
		}
		if pts == "" && ser.Tail == "" {
			continue
		}
		ser.Points = pts
		if ser.Tail != "" && lastAt == "" {
			ser.Tail = "" // nothing settled to hang it from
		}
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
	for i := 0; i <= 4; i++ {
		v := c.Max * int64(i) / 4
		y := float64(chartH) - float64(v)*float64(chartH)/float64(c.Max)
		out = append(out, map[string]any{
			"Y": y,
			// Top is the same line as a percentage, for the label beside it: the
			// SVG scales to its container, so a label placed in viewBox units
			// drifts away from the gridline it belongs to.
			"Top": y * 100 / chartH,
			"N":   v,
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
	// The hourly view is read as a day, so it is labelled every other hour --
	// twelve marks that a reader counts round rather than eight that only sample
	// it. The other periods thin to eight, which is as many as their longer
	// labels fit.
	every := maxInt(n/8, 1)
	if c.Period == "hour" {
		every = 2
	}
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

// fillHours lays out the whole day, whether or not it has been counted.
//
// The axis runs midnight to midnight in the reader's own clock, so the shape of
// a day is read against a fixed frame rather than against however many hours
// happen to hold data. Hours that were counted and saw nobody are real zeros
// and get a point on the floor. Hours before counting started, and hours still
// to come, are not zero but unmeasured: they hold the x position and carry no
// point, so the line begins where the counting did and ends at the hour now in
// progress.
func fillHours(out *TrafficChart, since time.Time, loc *time.Location) {
	now := time.Now().In(loc)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	from, upto := since.In(loc).Truncate(time.Hour), now.Truncate(time.Hour)

	have := map[string]TrafficPoint{}
	for _, pt := range out.Points {
		have[pt.Label] = pt
	}
	day := make([]TrafficPoint, 0, 24)
	for h := 0; h < 24; h++ {
		at := midnight.Add(time.Duration(h) * time.Hour)
		lbl := at.Format(trafficLayout("hour"))
		running := at.Equal(upto)
		if pt, ok := have[lbl]; ok {
			pt.Running = running
			day = append(day, pt)
			continue
		}
		known := !at.Before(from) && !at.After(upto)
		day = append(day, TrafficPoint{Label: lbl, Known: known, Counted: known, Running: running})
	}
	out.Points = day
	if n := len(day); n > 0 {
		for _, pt := range day {
			if pt.Running && pt.Known {
				out.Partial = true
			}
		}
	}
}

// TrafficHit is one column's hover target: the invisible band a pointer lands
// in, the figures to show for it, and the sentence a screen reader is given
// instead of the shape.
type TrafficHit struct {
	I     int
	X     float64 // left edge, percent
	W     float64 // width, percent
	Mid   float64 // centre, percent -- where the bubble is anchored
	Label string
	Rows  []TrafficHitRow
	Alt   string
	// Flip anchors the bubble by its right edge. Centred on a column near the
	// end of the plot it would hang off the panel; which columns those are is
	// arithmetic, not something a CSS attribute-substring match can be trusted
	// to guess.
	Flip bool
}

// TrafficHitRow is one line inside the bubble.
type TrafficHitRow struct {
	Slot int
	Key  string
	Val  string
}

// Hits builds the hover targets. They are computed here and rendered into the
// page rather than assembled in the browser, so the figures are in the HTML --
// where a screen reader and a page saved to disk can both still find them now
// that the table underneath is gone.
func (c TrafficChart) Hits(lang string) []TrafficHit {
	n := len(c.Points)
	if n == 0 {
		return nil
	}
	span := 100.0 / float64(n)
	step := float64(chartW-padR) / float64(maxInt(n-1, 1)) * 100 / chartW
	out := make([]TrafficHit, 0, n)
	for i, p := range c.Points {
		h := TrafficHit{I: i, X: float64(i) * span, W: span, Mid: float64(i) * step, Label: p.Label}
		dash := "—"
		vals := []struct {
			key string
			v   int64
			ok  bool
		}{
			{"hosts", p.Hosts, p.Counted},
			{"visitors", p.Visitors, p.Counted},
			{"visits", p.Visits, p.Counted},
			{"views", p.Views, p.Known},
		}
		alt := p.Label
		for slot, v := range vals {
			s := dash
			if v.ok {
				s = fmt.Sprintf("%d", v.v)
			}
			h.Rows = append(h.Rows, TrafficHitRow{Slot: slot, Key: v.key, Val: s})
			alt += ", " + T(lang, "tc.s_"+v.key) + " " + s
		}
		h.Alt = alt
		out = append(out, h)
	}
	return out
}
