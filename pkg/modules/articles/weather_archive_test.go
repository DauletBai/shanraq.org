package articles

import (
	"strings"
	"testing"
	"time"
)

// The server runs in UTC and the readers do not. From seven in the evening
// until midnight the two disagree about the date, and the strip was printing
// the server's answer — every evening, yesterday.
func TestSiteClockIsTheReadersNotTheServers(t *testing.T) {
	// 20:30 in Almaty on the 25th is 15:30 UTC on the 25th; an hour before
	// midnight there it is still the 25th UTC, and the day must not slip.
	utcEvening := time.Date(2026, 8, 25, 19, 30, 0, 0, time.UTC) // 00:30 on the 26th here
	if got := siteDay(utcEvening); got.Day() != 26 {
		t.Errorf("19:30 UTC — это уже 26-е в Алматы, получили %s", got.Format("2006-01-02"))
	}
	utcMorning := time.Date(2026, 8, 25, 3, 0, 0, 0, time.UTC) // 08:00 here
	if got := siteDay(utcMorning); got.Day() != 25 {
		t.Errorf("03:00 UTC — это 25-е в Алматы, получили %s", got.Format("2006-01-02"))
	}
	// The zone must resolve even in an image with no timezone database.
	if off := siteNow().Format("-07:00"); off != "+05:00" {
		t.Errorf("смещение сайта %s, ожидалось +05:00", off)
	}
}

// Still air is a state, not a measurement of zero.
func TestCalmIsNotZero(t *testing.T) {
	if got := wxWind(1.4, LangRU); got != "Штиль" {
		t.Errorf("1,4 км/ч — штиль, получили %q", got)
	}
	if got := wxWind(18, LangRU); !strings.HasPrefix(got, "5 ") {
		t.Errorf("18 км/ч — это 5 м/с, получили %q", got)
	}
}

// A rounded −0.4 must not print as "-0", which reads like a typo.
func TestTemperatureNeverPrintsNegativeZero(t *testing.T) {
	cases := map[float64]string{
		-0.4: "0°", 0.4: "0°", 0: "0°",
		17.6: "+18°", -12.2: "-12°", 30: "+30°",
	}
	for in, want := range cases {
		if got := wxTemp(in); got != want {
			t.Errorf("wxTemp(%.1f) = %q, ожидалось %q", in, got, want)
		}
	}
}

// The hourly series starts at midnight, so everything already behind us is
// dropped: a forecast for this morning is not a forecast.
func TestHourlySeriesStartsFromNow(t *testing.T) {
	times := []string{
		"2026-08-25T00:00", "2026-08-25T01:00", "2026-08-25T07:00",
		"2026-08-25T08:00", "2026-08-25T09:00",
	}
	if got := wxHourIndex(times, "2026-08-25T07:15"); got != 2 {
		t.Errorf("в 07:15 ряд должен начинаться с 07:00 (индекс 2), получили %d", got)
	}
	if got := wxHourIndex(times, "2026-08-25T00:00"); got != 0 {
		t.Errorf("в полночь ряд начинается с начала, получили %d", got)
	}
	// A reading past the end of the series must not send the slice negative.
	if got := wxHourIndex(times, "2026-08-26T23:00"); got != 0 {
		t.Errorf("за концом ряда ожидался 0, получили %d", got)
	}
}

// WMO codes are a published standard; every group the icons use must also have
// words, in all three languages.
func TestEveryWeatherCodeHasWords(t *testing.T) {
	for _, code := range []int{0, 1, 2, 3, 45, 48, 51, 61, 71, 80, 85, 95, 99} {
		for _, lang := range []string{LangKZ, LangRU, LangEN} {
			if s := wxDescribe(code, lang); strings.TrimSpace(s) == "" {
				t.Errorf("код %d на языке %q без описания", code, lang)
			}
		}
		if weatherIconName(code) == "" {
			t.Errorf("код %d без иконки", code)
		}
	}
}

// The archive answers only for days that have already happened here.
func TestArchiveRefusesTheFuture(t *testing.T) {
	app := newTestApp(t)
	today := siteDay(siteNow())
	cases := []struct {
		path string
		want int
	}{
		{"/archive/" + today.Format("2006-01-02"), 200},
		{"/archive/" + today.AddDate(0, 0, -400).Format("2006-01-02"), 200},
		{"/archive/" + today.AddDate(0, 0, 1).Format("2006-01-02"), 404},
		{"/archive/2026-13-45", 404},
		{"/archive/tomorrow", 404},
	}
	for _, c := range cases {
		rec := app.do("GET", c.path, nil)
		if rec.Code != c.want {
			t.Errorf("%s → %d, ожидалось %d", c.path, rec.Code, c.want)
		}
	}
}

// A day with nothing on it is a legitimate address — the strip links to it on a
// quiet day — but it must not join the index claiming to be about that day.
func TestEmptyArchiveDayIsNotIndexed(t *testing.T) {
	app := newTestApp(t)
	rec := app.do("GET", "/archive/2019-03-07", nil) // до первой публикации
	if rec.Code != 200 {
		t.Fatalf("пустой день ответил %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "noindex") {
		t.Error("пустой день архива не помечен noindex")
	}
}

// A forecast exists only where there are coordinates; anything else is a
// forecast for somewhere else.
func TestWeatherNeedsCoordinates(t *testing.T) {
	app := newTestApp(t)
	// The region has no coordinates of its own in the reference.
	if rec := app.do("GET", "/weather/kostanaiskaya-oblast", nil); rec.Code != 404 {
		t.Errorf("область без координат ответила %d, ожидался 404", rec.Code)
	}
	if rec := app.do("GET", "/weather/nosuchplace", nil); rec.Code != 404 {
		t.Errorf("несуществующее место ответило %d, ожидался 404", rec.Code)
	}
}
