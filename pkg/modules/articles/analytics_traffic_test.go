package articles

import (
	"testing"
	"time"
)

// The hourly chart is keyed by a label that carries no date: an hour prints as
// "15:00". A query window reaching back a full day therefore answered for this
// evening with yesterday's, and the chart drew traffic for hours that had not
// happened yet. These tests hold both halves of the fix: the window now starts
// where the axis starts, and a bucket past the one in progress stays unmeasured
// whatever a query returns under its name.

func TestPeriodStartIsTheBoundaryTheAxisDrawsFrom(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 8, 31, 10, 53, 0, 0, loc)
	for _, c := range []struct {
		code string
		want time.Time
	}{
		{"hour", time.Date(2026, 8, 31, 0, 0, 0, 0, loc)},
		{"day", time.Date(2026, 8, 1, 0, 0, 0, 0, loc)},
		{"month", time.Date(2026, 1, 1, 0, 0, 0, 0, loc)},
	} {
		if got := periodStart(c.code, now, loc); !got.Equal(c.want) {
			t.Errorf("periodStart(%q) = %s, want %s", c.code, got, c.want)
		}
	}
}

func TestHoursStillToComeCarryNoMeasurement(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	// Yesterday's traffic, labelled the way today's would be. Before the fix
	// these filled the rest of the axis.
	var pts []TrafficPoint
	for h := 0; h < 24; h++ {
		lbl := time.Date(now.Year(), now.Month(), now.Day(), h, 0, 0, 0, loc).Format(trafficLayout("hour"))
		pts = append(pts, TrafficPoint{Label: lbl, Known: true, Counted: true, Visits: 42, Views: 99})
	}
	out := TrafficChart{Points: pts}
	fillPeriod(&out, now.AddDate(0, 0, -7), loc, "hour")

	if len(out.Points) != 24 {
		t.Fatalf("axis has %d columns, want 24", len(out.Points))
	}
	cur := now.Hour()
	for h, pt := range out.Points {
		switch {
		case h > cur:
			if pt.Known || pt.Counted || pt.Visits != 0 || pt.Views != 0 {
				t.Errorf("hour %02d:00 is still to come but carries %+v", h, pt)
			}
		case h == cur:
			if !pt.Running {
				t.Errorf("hour %02d:00 is the one in progress and is not marked running", h)
			}
		default:
			if !pt.Known {
				t.Errorf("hour %02d:00 has passed and lost its measurement", h)
			}
		}
	}
}

func TestTheHourInProgressMakesTheChartPartial(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	out := TrafficChart{}
	fillPeriod(&out, now.AddDate(0, 0, -7), loc, "hour")
	if !out.Partial {
		t.Error("an hour is always in progress, so the chart is always partial")
	}
}
