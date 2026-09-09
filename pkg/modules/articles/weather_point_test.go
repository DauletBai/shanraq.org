package articles

import (
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

// The strip's temperature is per reader, and the machinery that makes it so has
// two jobs: never ask the forecast service while somebody waits, and never
// answer with a place that is not a place.
func TestBarPlaceKeyAndBounds(t *testing.T) {
	if got := wxPlaceKey(43.2389, 76.8897); got != "43.2,76.9" {
		t.Errorf("ключ места = %q, ожидался 43.2,76.9", got)
	}
	// One key for one city: two addresses eleven kilometres apart must not send
	// two requests for the same temperature.
	if wxPlaceKey(51.1605, 71.4704) != wxPlaceKey(51.1988, 71.4501) {
		t.Error("два адреса одного города дали разные ключи")
	}
	for _, bad := range [][2]float64{{0, 0}, {91, 20}, {40, 181}, {-91, -181}} {
		if wxPointOK(bad[0], bad[1]) {
			t.Errorf("точка %v принята, а она не место на Земле", bad)
		}
	}
	if !wxPointOK(47.09, 51.92) {
		t.Error("Атырау отвергнут как невозможная точка")
	}
}

// A reader whose place is not cached yet gets the default city rather than a
// blank cell and rather than a page that waits for open-meteo.
func TestBarFallsBackToTheDefaultCity(t *testing.T) {
	bar := NewInfoBar(zap.NewNop(), nil, "")
	bar.weatherIc, bar.weatherTmp, bar.weatherPres = "wx_sun", "+21°", "742"
	// Seeded as pending so the miss does not start a real fetch in a test.
	bar.byPlace[wxPlaceKey(47.09, 51.92)] = wxBarReading{pending: true}

	got := bar.SnapshotAt("сегодня", "2026-09-09", "Атырау", 47.09, 51.92)
	if got.WeatherTemp != "+21°" {
		t.Errorf("температура = %q, ожидался запасной вариант +21°", got.WeatherTemp)
	}
	if got.WeatherPlace != "" {
		t.Errorf("подписано место %q, хотя показана погода города по умолчанию", got.WeatherPlace)
	}

	bar.byPlace[wxPlaceKey(47.09, 51.92)] = wxBarReading{at: time.Now(), icon: "wx_snow", temp: "-8°", press: "751"}
	got = bar.SnapshotAt("сегодня", "2026-09-09", "Атырау", 47.09, 51.92)
	if got.WeatherTemp != "-8°" || got.WeatherIcon != "wx_snow" || got.WeatherPress != "751" {
		t.Errorf("своя погода не подставилась: %+v", got)
	}
	if got.WeatherPlace != "Атырау" {
		t.Errorf("место = %q, ожидался Атырау", got.WeatherPlace)
	}
	// An unknown place leaves the strip exactly as it was.
	if plain := bar.SnapshotAt("сегодня", "2026-09-09", "", 0, 0); plain.WeatherTemp != "+21°" || plain.WeatherPlace != "" {
		t.Errorf("нулевые координаты изменили полосу: %+v", plain)
	}
}

// The map answers a press with a fragment of the same page. A point that is not
// a point is refused before anything is fetched.
func TestWeatherPointRefusesNonsense(t *testing.T) {
	app := newTestApp(t)
	for _, q := range []string{"", "?lat=abc&lon=51", "?lat=95&lon=51", "?lat=47", "?lat=0&lon=0"} {
		if rec := app.do("GET", "/weather/point"+q, nil); rec.Code != 400 {
			t.Errorf("точка %q ответила %d, ожидался 400", q, rec.Code)
		}
	}
}

// The fragment is a fragment: a piece of the page, not a page. Injecting a whole
// document into a div is how a layout ends up with two headers.
func TestWeatherPointIsAFragment(t *testing.T) {
	app := newTestApp(t)
	rec := app.do("GET", "/weather/point?lat=47.09&lon=51.92", nil)
	if rec.Code != 200 {
		t.Fatalf("точка ответила %d", rec.Code)
	}
	body := rec.Body.String()
	for _, whole := range []string{"<html", "<body", "site_header", "<footer"} {
		if strings.Contains(body, whole) {
			t.Errorf("в ответе оказался кусок целой страницы: %q", whole)
		}
	}
	// Either the forecast arrived or the fragment says so; silence is the one
	// thing a reader who pressed the map must not get.
	if !strings.Contains(body, "chart__h") && !strings.Contains(body, "chart__note") {
		t.Errorf("фрагмент пуст: %.200q", body)
	}
}
