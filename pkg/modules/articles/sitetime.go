package articles

import (
	"sync"
	"time"
)

// The site's own clock.
//
// The server runs in UTC and the readers do not. Kazakhstan is five hours
// ahead, so from seven in the evening until midnight the two disagree about
// what day it is — and the strip at the top of every page was showing the
// server's answer. Every evening it printed yesterday's date.
//
// That was cosmetic while the date was only text. It stops being cosmetic the
// moment the date becomes a link to what was published that day, and it decides
// which day an article published at half past one in the morning belongs to.

// siteTimeZone is the zone the site keeps its calendar in.
const siteTimeZone = "Asia/Almaty"

// siteTZOffset is the fallback offset in hours, for an image with no timezone
// database. Kazakhstan moved to a single zone at UTC+5 in March 2024 and has no
// daylight saving, so a fixed offset is not an approximation here — it is the
// rule.
const siteTZOffset = 5

var siteLoc = sync.OnceValue(func() *time.Location {
	if loc, err := time.LoadLocation(siteTimeZone); err == nil {
		return loc
	}
	return time.FixedZone(siteTimeZone, siteTZOffset*3600)
})

// siteNow is the current time as a reader in the country experiences it.
func siteNow() time.Time { return time.Now().In(siteLoc()) }

// siteDay truncates a moment to the calendar day it falls on here.
func siteDay(t time.Time) time.Time {
	t = t.In(siteLoc())
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, siteLoc())
}
