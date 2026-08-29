package articles

import (
	"net/http"
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

// tzCookieName carries the reader's own IANA zone, set once by the browser.
//
// The server cannot know it: an HTTP request says nothing about the clock on
// the other end, and the address only narrows it to a country -- which settles
// nothing for a reader in Russia or the United States, and guesses wrong for
// anyone travelling. So the browser, which does know, says so once and the
// answer is kept for a year.
//
// The value is a zone name and nothing else. It is not an identifier: millions
// of people share "Asia/Almaty", and it is used for one thing, deciding which
// hour a figure is drawn under.
const tzCookieName = "tz"

// readerLoc is the clock to show a reader their figures in.
//
// Their own if the browser has told us, the site's otherwise -- which is the
// right default, because most readers are here and because a chart with no
// script running is still a chart that has to be correct for somebody.
func readerLoc(r *http.Request) *time.Location {
	if r == nil {
		return siteLoc()
	}
	c, err := r.Cookie(tzCookieName)
	if err != nil || c.Value == "" || len(c.Value) > 64 {
		return siteLoc()
	}
	// LoadLocation is the validation: a name the zone database does not know is
	// simply not a zone, and nothing else is done with the string.
	loc, err := time.LoadLocation(c.Value)
	if err != nil {
		return siteLoc()
	}
	return loc
}

// readerNow is siteNow through the reader's own clock.
func readerNow(r *http.Request) time.Time { return time.Now().In(readerLoc(r)) }
