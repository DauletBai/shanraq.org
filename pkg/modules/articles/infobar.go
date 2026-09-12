package articles

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"shanraq.org/internal/config"
)

// The top info bar shows the date, weather and KZT exchange rates plus social
// links. Weather and rates are fetched SERVER-SIDE and cached — the strict CSP
// blocks the browser from calling external APIs, and server-side fetching
// degrades gracefully (a cell simply hides) if a source is blocked or down.

// wxBarReading is one place's temperature for the strip.
type wxBarReading struct {
	at    time.Time
	icon  string
	temp  string
	press string
	// pending marks a fetch already on its way, so a hundred readers from one
	// city ask the forecast service once between them.
	pending bool
}

// wxBarTTL is how long a reading is reused, and wxBarPlaces how many places are
// kept at once. The strip is on every page, so both numbers are about not
// turning a popular site into a load generator: half an hour is finer than the
// weather changes, and a few hundred places cover a country several times over
// while giving a crawler nothing to fill memory with.
const (
	wxBarTTL    = 30 * time.Minute
	wxBarPlaces = 400
)

// InfoBar fetches and caches the weather and exchange rates in the background.
type InfoBar struct {
	mu          sync.RWMutex
	rates       []Rate
	weatherIc   string
	weatherTmp  string
	weatherPres string
	// byPlace holds a reading per rounded coordinate. The default city's own
	// reading stays in the fields above: it is the answer for everyone whose
	// place we do not know, and it must be there before the first request.
	byPlace map[string]wxBarReading

	httpc    *http.Client
	log      *zap.Logger
	lat, lon float64
	social   []SocialLink
	github   string
}

// socialLinks returns the four social profiles in a fixed order so their icons
// always render. A profile with no configured URL yet gets a "#" placeholder
// (shown but non-navigating) until a real URL is filled in.
func socialLinks(cfg config.SocialConfig) []SocialLink {
	out := make([]SocialLink, 0, 4)
	for _, s := range []SocialLink{
		{Name: "telegram", URL: strings.TrimSpace(cfg.Telegram)},
		{Name: "instagram", URL: strings.TrimSpace(cfg.Instagram)},
		{Name: "youtube", URL: strings.TrimSpace(cfg.YouTube)},
		{Name: "facebook", URL: strings.TrimSpace(cfg.Facebook)},
	} {
		if s.URL == "" {
			s.URL = "#"
		}
		out = append(out, s)
	}
	return out
}

// NewInfoBar builds the provider. Weather defaults to Almaty.
func NewInfoBar(log *zap.Logger, social []SocialLink, github string) *InfoBar {
	return &InfoBar{
		httpc:   &http.Client{Timeout: 10 * time.Second},
		log:     log,
		lat:     wxDefaultLat,
		lon:     wxDefaultLon,
		social:  social,
		github:  strings.TrimSpace(github),
		byPlace: map[string]wxBarReading{},
	}
}

// Snapshot returns the cached bar data for the given day string.
func (b *InfoBar) Snapshot(today, todayISO string) InfoBarData {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return InfoBarData{Today: today, TodayISO: todayISO, WeatherIcon: b.weatherIc, WeatherTemp: b.weatherTmp, WeatherPress: b.weatherPres, Rates: b.rates, Social: b.social, GitHub: b.github}
}

// SnapshotAt is the same bar with the temperature of the reader's own place.
//
// The strip showed one city to everybody, which is right for the exchange rate
// and wrong for the sky: a reader in Oral has no use for what it is doing in
// Almaty. When the place is unknown -- no city database, an address it does not
// cover -- the default city answers, exactly as before.
//
// Nothing is fetched while the reader waits. A place we have not seen for half
// an hour is refreshed in the background and shows the default city until it
// arrives, because a page must not wait on somebody else's server.
func (b *InfoBar) SnapshotAt(today, todayISO, place string, lat, lon float64) InfoBarData {
	data := b.Snapshot(today, todayISO)
	if !wxPointOK(lat, lon) {
		return data
	}
	key := wxPlaceKey(lat, lon)

	b.mu.Lock()
	got, ok := b.byPlace[key]
	fresh := ok && time.Since(got.at) < wxBarTTL
	if !fresh && !got.pending {
		got.pending = true
		b.byPlace[key] = got
		go b.fetchPlace(key, lat, lon)
	}
	b.mu.Unlock()

	if !ok || got.temp == "" {
		return data
	}
	data.WeatherIcon, data.WeatherTemp, data.WeatherPress = got.icon, got.temp, got.press
	data.WeatherPlace = place
	return data
}

// wxPointOK rejects coordinates that cannot be a place on Earth, including the
// 0,0 that a database returns when it means "no idea".
func wxPointOK(lat, lon float64) bool {
	if lat == 0 && lon == 0 {
		return false
	}
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

// wxPlaceKey rounds a point to a tenth of a degree -- some eleven kilometres,
// which is one city and one temperature.
func wxPlaceKey(lat, lon float64) string {
	return fmt.Sprintf("%.1f,%.1f", lat, lon)
}

// fetchPlace reads one place's current weather and stores it under key.
func (b *InfoBar) fetchPlace(key string, lat, lon float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	icon, temp, press, err := b.currentWeather(ctx, lat, lon)

	b.mu.Lock()
	defer b.mu.Unlock()
	entry := b.byPlace[key]
	entry.pending = false
	if err != nil {
		// The failure is kept quiet on purpose: an unreachable forecast service
		// is not this reader's problem, and the strip already has an answer.
		b.byPlace[key] = entry
		return
	}
	entry.at, entry.icon, entry.temp, entry.press = time.Now(), icon, temp, press
	// A cache with no ceiling is a way for a crawler walking every address to
	// fill memory. Past the ceiling the oldest reading goes.
	if len(b.byPlace) >= wxBarPlaces {
		oldest, at := "", time.Now()
		for k, v := range b.byPlace {
			if k != key && !v.pending && v.at.Before(at) {
				oldest, at = k, v.at
			}
		}
		if oldest != "" {
			delete(b.byPlace, oldest)
		}
	}
	b.byPlace[key] = entry
}

// Run refreshes weather every 30 min and rates every ~6 h until ctx is done.
func (b *InfoBar) Run(ctx context.Context) {
	b.refreshRates(ctx)
	b.refreshWeather(ctx)
	tick := time.NewTicker(30 * time.Minute)
	defer tick.Stop()
	n := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			b.refreshWeather(ctx)
			// Rates refresh every ~6 h (n%12), but if a transient failure left
			// them empty, retry on every 30-min tick so they recover in minutes
			// rather than staying blank for the full 6-h cycle.
			if n++; n%12 == 0 || b.ratesEmpty() {
				b.refreshRates(ctx)
			}
		}
	}
}

func (b *InfoBar) get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Shanraq/1.0 (+https://shanraq.org)")
	resp, err := b.httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}

// ratesEmpty reports whether no exchange rates are currently cached (e.g. after
// a failed fetch), so the Run loop can retry sooner instead of waiting 6 h.
func (b *InfoBar) ratesEmpty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.rates) == 0
}

// refreshRates pulls USD/EUR/RUB from the National Bank of Kazakhstan feed. The
// fetch is retried a few times with a short backoff so a transient network blip
// (a single slow response at boot) doesn't leave the bar without rates.
func (b *InfoBar) refreshRates(ctx context.Context) {
	var body []byte
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		if body, err = b.get(ctx, "https://nationalbank.kz/rss/rates_all.xml"); err == nil {
			break
		}
	}
	if err != nil {
		b.log.Warn("infobar rates fetch", zap.Error(err))
		return
	}
	var doc struct {
		Items []struct {
			Title string `xml:"title"`
			Desc  string `xml:"description"`
			Index string `xml:"index"`
		} `xml:"channel>item"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		b.log.Warn("infobar rates parse", zap.Error(err))
		return
	}
	want := map[string]bool{"USD": true, "EUR": true, "RUB": true}
	order := []string{"USD", "EUR", "RUB"}
	found := map[string]Rate{}
	for _, it := range doc.Items {
		code := strings.TrimSpace(it.Title)
		if !want[code] {
			continue
		}
		dir := ""
		switch strings.ToUpper(strings.TrimSpace(it.Index)) {
		case "UP":
			dir = "up"
		case "DOWN":
			dir = "down"
		}
		val := strings.TrimSpace(it.Desc)
		main, last := val, ""
		if r := []rune(val); len(r) > 1 {
			main, last = string(r[:len(r)-1]), string(r[len(r)-1])
		}
		found[code] = Rate{Code: code, Main: main, Last: last, Dir: dir}
	}
	var rates []Rate
	for _, c := range order {
		if r, ok := found[c]; ok {
			rates = append(rates, r)
		}
	}
	if len(rates) > 0 {
		b.mu.Lock()
		b.rates = rates
		b.mu.Unlock()
	}
}

// refreshWeather pulls the default city's current conditions for the strip.
func (b *InfoBar) refreshWeather(ctx context.Context) {
	icon, temp, press, err := b.currentWeather(ctx, b.lat, b.lon)
	if err != nil {
		b.log.Warn("infobar weather", zap.Error(err))
		return
	}
	b.mu.Lock()
	b.weatherIc, b.weatherTmp, b.weatherPres = icon, temp, press
	b.mu.Unlock()
}

// currentWeather reads one point's conditions from open-meteo and formats them
// the way the strip shows them.
func (b *InfoBar) currentWeather(ctx context.Context, lat, lon float64) (icon, temp, press string, err error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,weather_code,pressure_msl&timezone=auto", lat, lon)
	body, err := b.get(ctx, url)
	if err != nil {
		return "", "", "", err
	}
	var doc struct {
		Current struct {
			Temp  float64 `json:"temperature_2m"`
			Code  int     `json:"weather_code"`
			Press float64 `json:"pressure_msl"` // hPa
		} `json:"current"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return "", "", "", err
	}
	sign := ""
	t := int(doc.Current.Temp)
	if t > 0 {
		sign = "+"
	}
	// Pressure in mm Hg (hPa × 0.750062) — just the number; the template appends
	// the unit in the active UI language ("мм рт.ст." / "mmHg" / "мм сын.бағ.").
	if doc.Current.Press > 0 {
		press = fmt.Sprintf("%d", int(doc.Current.Press*0.750062+0.5))
	}
	return weatherIconName(doc.Current.Code), fmt.Sprintf("%s%d°", sign, t), press, nil
}

// monthNames holds the info-bar month names per UI language. Russian uses the
// genitive case so "18 июля 2026" reads naturally with the day first.
var monthNames = map[string][]string{
	LangRU: {"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"},
	LangKZ: {"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым", "шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан"},
	LangEN: {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
}

// weekdayNames holds the info-bar weekday names per UI language, indexed by
// time.Weekday (Sunday = 0).
var weekdayNames = map[string][]string{
	LangRU: {"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"},
	LangKZ: {"жексенбі", "дүйсенбі", "сейсенбі", "сәрсенбі", "бейсенбі", "жұма", "сенбі"},
	LangEN: {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
}

// localizedDate renders a date as "воскресенье, 18 июля 2026" — weekday first,
// the way all three languages read a date aloud — for the info-bar calendar
// cell.
func localizedDate(lang string, t time.Time) string {
	months, ok := monthNames[lang]
	if !ok {
		months = monthNames[LangRU]
	}
	days, ok := weekdayNames[lang]
	if !ok {
		days = weekdayNames[LangRU]
	}
	return fmt.Sprintf("%s, %d %s %d", days[int(t.Weekday())], t.Day(), months[int(t.Month())-1], t.Year())
}

// weatherIconName maps a WMO weather code to a Shanraq weather icon key.
func weatherIconName(code int) string {
	switch {
	case code == 0:
		return "wx_sun"
	case code <= 3:
		return "wx_cloud"
	case code == 45 || code == 48:
		return "wx_fog"
	case code >= 51 && code <= 67:
		return "wx_rain"
	case code >= 71 && code <= 77:
		return "wx_snow"
	case code >= 80 && code <= 82:
		return "wx_rain"
	case code >= 95:
		return "wx_storm"
	default:
		return "wx_cloud"
	}
}
