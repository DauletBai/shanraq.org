package articles

import (
	"testing"
	"time"
)

// The chart rolls: it ends at the bucket now in progress and reaches back a
// fixed number of buckets. Before that it framed a calendar unit, which drew
// the part of the day that had not happened yet as empty air — and, because an
// hour prints as "15:00" with no date on it, let yesterday evening answer to
// this evening's name and fill that emptiness with the wrong day's traffic.

func TestTheWindowEndsAtTheBucketInProgress(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 8, 31, 10, 53, 0, 0, loc)
	for _, c := range []struct {
		code string
		want time.Time
	}{
		{"hour", time.Date(2026, 8, 30, 11, 0, 0, 0, loc)}, // 24 hours back to 11:00 yesterday
		{"day", time.Date(2026, 8, 2, 0, 0, 0, 0, loc)},    // 30 days back
		{"month", time.Date(2025, 9, 1, 0, 0, 0, 0, loc)},  // 12 months back
	} {
		got := periodStart(c.code, now, loc)
		if !got.Equal(c.want) {
			t.Errorf("periodStart(%q) = %s, want %s", c.code, got, c.want)
		}
		// The window must land exactly on the bucket in progress, never past it.
		n := periodLen(c.code)
		var last time.Time
		switch c.code {
		case "hour":
			last = got.Add(time.Duration(n-1) * time.Hour)
		case "month":
			last = got.AddDate(0, n-1, 0)
		default:
			last = got.AddDate(0, 0, n-1)
		}
		if last.After(now) {
			t.Errorf("%s window ends at %s, past now %s", c.code, last, now)
		}
	}
}

func TestNoTwoColumnsShareALabel(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	for _, code := range []string{"hour", "day", "month"} {
		out := TrafficChart{}
		fillPeriod(&out, now.AddDate(-2, 0, 0), loc, code)
		if len(out.Points) != periodLen(code) {
			t.Fatalf("%s drew %d columns, want %d", code, len(out.Points), periodLen(code))
		}
		seen := map[string]int{}
		for i, pt := range out.Points {
			if j, dup := seen[pt.Label]; dup {
				t.Errorf("%s: columns %d and %d are both labelled %q", code, j, i, pt.Label)
			}
			seen[pt.Label] = i
		}
	}
}

func TestTheNewestColumnIsLastAndStillFilling(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	out := TrafficChart{}
	fillPeriod(&out, now.AddDate(0, 0, -7), loc, "hour")

	last := out.Points[len(out.Points)-1]
	if want := now.Format(trafficLayout("hour")); last.Label != want {
		t.Errorf("last column is %q, want the hour in progress %q", last.Label, want)
	}
	if !last.Running {
		t.Error("the last column is the hour in progress and is not marked running")
	}
	if !out.Partial {
		t.Error("an hour is always in progress, so the chart is always partial")
	}
	for i, pt := range out.Points[:len(out.Points)-1] {
		if pt.Running {
			t.Errorf("column %d is not the last and is marked running", i)
		}
	}
}

func TestTheHourInProgressCarriesALabelOfItsOwn(t *testing.T) {
	c := TrafficChart{Period: "hour"}
	for i := 0; i < 24; i++ {
		c.Points = append(c.Points, TrafficPoint{Label: time.Date(2026, 8, 31, i, 0, 0, 0, time.UTC).Format("15:00")})
	}
	labels := c.XLabels()
	if len(labels) == 0 {
		t.Fatal("no labels")
	}
	if last := labels[len(labels)-1]; last["Last"] != true {
		t.Error("the newest column carries no label of its own")
	}
	// Labels are thinned to every other hour because "22:00" needs two columns
	// to print in. Any pair closer than that overlaps, which is what the forced
	// mark at the end used to do to its neighbour.
	step := float64(chartW) / 23
	want := 2 * step
	for i := 1; i < len(labels); i++ {
		gap := labels[i]["VX"].(float64) - labels[i-1]["VX"].(float64)
		if gap < want-0.01 {
			t.Errorf("labels %q and %q sit %.1f units apart, need %.1f — they overlap",
				labels[i-1]["L"], labels[i]["L"], gap, want)
		}
	}
}
