package articles

import (
	"testing"
	"time"
)

// The privacy policy tells readers that yesterday's identifier cannot be tied
// to today's. That is a property of the code, so it is asserted here rather
// than trusted.
func TestIdentUnlinkableAcrossSalts(t *testing.T) {
	const ip, ua = "203.0.113.7", "Mozilla/5.0 (Linux; Android 14)"
	v1, h1 := ident([]byte("monday-key-monday-key-monday-key"), ip, ua)
	v2, h2 := ident([]byte("tuesday-key-tuesday-key-tuesdayk"), ip, ua)
	if v1 == v2 {
		t.Error("same visitor id under two different daily keys: readers could be followed across days")
	}
	if h1 == h2 {
		t.Error("same host id under two different daily keys")
	}
}

// Within one day the identifier has to be stable, or every page view would look
// like a new person and the visitor count would just be the view count.
func TestIdentStableWithinSalt(t *testing.T) {
	salt := []byte("one-day-key-one-day-key-one-dayk")
	const ip, ua = "203.0.113.7", "Mozilla/5.0"
	v1, h1 := ident(salt, ip, ua)
	v2, h2 := ident(salt, ip, ua)
	if v1 != v2 || h1 != h2 {
		t.Fatal("identifier is not stable for one reader inside one day")
	}
}

// Two people behind one address are one host and two visitors. That is what the
// words mean in every other counter, and the chart is read against those.
func TestHostSharedVisitorNot(t *testing.T) {
	salt := []byte("one-day-key-one-day-key-one-dayk")
	v1, h1 := ident(salt, "203.0.113.7", "Firefox/1")
	v2, h2 := ident(salt, "203.0.113.7", "Safari/2")
	if h1 != h2 {
		t.Error("one address produced two hosts")
	}
	if v1 == v2 {
		t.Error("two browsers on one address collapsed into one visitor")
	}
}

// A visit is a half-hour window: five pages in ten minutes is one visit, a
// return after the break is another.
func TestSlotOfHalfHourWindows(t *testing.T) {
	at := func(h, m int) time.Time { return time.Date(2026, 8, 29, h, m, 0, 0, time.UTC) }
	if slotOf(at(10, 3)) != slotOf(at(10, 29)) {
		t.Error("two hits ten minutes apart fell into different visits")
	}
	if slotOf(at(10, 29)) == slotOf(at(10, 31)) {
		t.Error("the half-hour boundary did not start a new visit")
	}
	if got, want := slotOf(at(10, 47)), at(10, 30); got != want {
		t.Errorf("slot = %v, want %v", got, want)
	}
}
