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

func TestTheNewestColumnIsMarkedAndTheGridStaysThin(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	c := TrafficChart{Period: "hour"}
	fillPeriod(&c, now.AddDate(0, 0, -7), loc, "hour")
	labels := c.XLabels()
	if len(labels) == 0 {
		t.Fatal("no labels")
	}
	last := labels[len(labels)-1]
	if last["Last"] != true || last["L"] == "" {
		t.Error("the newest column carries no label of its own")
	}
	if last["Grid"] != true {
		t.Error("the newest column carries no gridline; the grid is counted from the left again")
	}
	// Gridlines two columns apart at the narrowest. Counted from the left, a
	// forced line at the end landed one column from its neighbour.
	step := float64(chartW) / 23
	prev := -1.0
	for _, l := range labels {
		if l["Grid"] != true {
			continue
		}
		vx := l["VX"].(float64)
		if prev >= 0 && vx-prev < 2*step-0.01 {
			t.Errorf("gridlines %.1f units apart, need %.1f", vx-prev, 2*step)
		}
		prev = vx
	}
}

// The scale beside the chart held a gutter wide enough for six digits, which
// at every reading ever taken was mostly empty. Thousands are shortened
// instead, so the gutter can be narrow and stay narrow.

func TestKiloShortensOnlyWhereItSaysSomething(t *testing.T) {
	for _, c := range []struct {
		lang string
		n    int64
		want string
	}{
		{"ru", 0, "0"},
		{"ru", 999, "999"},
		{"ru", 1000, "1т"},
		{"ru", 1500, "1,5т"},
		{"ru", 9900, "9,9т"},
		{"ru", 12345, "12т"},
		{"ru", 150000, "150т"},
		{"kz", 1500, "1,5м"},
		{"en", 1500, "1.5k"},
		{"en", 999, "999"},
	} {
		if got := kilo(c.lang, c.n); got != c.want {
			t.Errorf("kilo(%q, %d) = %q, want %q", c.lang, c.n, got, c.want)
		}
	}
}

// The note saying views reach further back than the other three was printed on
// every chart, including the hourly one, which never asks the older counter at
// all. It is now tied to the counter having actually contributed.
func TestTheHourlyChartClaimsNoOlderViews(t *testing.T) {
	c := TrafficChart{Period: "hour"}
	if c.Backfilled {
		t.Error("an hourly chart never reads the older counter and must not say it did")
	}
}

// The axis labelled every second hour and every fourth day because the full
// label is wide. The ones it left out were the columns a reader was trying to
// read. Every column is labelled now, with the part that changes.

func TestEveryColumnCarriesItsOwnLabel(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	for _, code := range []string{"hour", "day", "month"} {
		c := TrafficChart{Period: code}
		fillPeriod(&c, now.AddDate(-2, 0, 0), loc, code)
		labels := c.XLabels()
		if len(labels) != periodLen(code) {
			t.Errorf("%s: %d labels for %d columns", code, len(labels), periodLen(code))
		}
		grids := 0
		for _, l := range labels {
			if l["L"] == "" {
				t.Errorf("%s: a column carries no label", code)
			}
			if l["Grid"] == true {
				grids++
			}
		}
		// Gridlines stay thinned: one behind every column is a fence.
		if grids >= len(labels) {
			t.Errorf("%s: %d gridlines for %d columns, they were not thinned", code, grids, len(labels))
		}
	}
}

func TestTheShortLabelKeepsTheUnitThatTurnedOver(t *testing.T) {
	for _, c := range []struct {
		code, want string
		at         time.Time
	}{
		{"hour", "13", time.Date(2026, 8, 31, 13, 0, 0, 0, time.UTC)},
		{"hour", "09", time.Date(2026, 8, 31, 9, 0, 0, 0, time.UTC)},
		{"day", "04", time.Date(2026, 8, 4, 0, 0, 0, 0, time.UTC)},
		{"day", "01.09", time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
		{"month", "08", time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)},
		{"month", "01.26", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	} {
		if got := shortLabel(c.code, c.at); got != c.want {
			t.Errorf("shortLabel(%q, %s) = %q, want %q", c.code, c.at.Format("2006-01-02 15:04"), got, c.want)
		}
	}
}

// The full label is what the tooltip and the screen reader get: they have room,
// and a reader hovering a column is asking exactly which one it is.
func TestTheTooltipKeepsTheWholeLabel(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	c := TrafficChart{Period: "hour"}
	fillPeriod(&c, now.AddDate(0, 0, -7), loc, "hour")
	for _, h := range c.Hits("ru") {
		if len(h.Label) != len("13:00") {
			t.Fatalf("tooltip label %q is not the whole label", h.Label)
		}
	}
}
