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

// Rate is one currency's KZT rate for the info bar. Main+Last split the value
// so the last digit can be dropped on narrow phones (Main shown, Last hidden).
type Rate struct {
	Code string // USD / EUR / RUB
	Main string // value without its last character, e.g. "469.8"
	Last string // the last character, e.g. "3"
	Dir  string // "up" | "down" | ""
}

// SocialLink is one configured social profile shown in the bar.
type SocialLink struct {
	Name string // telegram | instagram | youtube | facebook (icon key)
	URL  string
}

// InfoBarData is the per-request snapshot handed to templates.
type InfoBarData struct {
	Today        string
	WeatherIcon  string // icon key, e.g. "wx_sun" ("" when unavailable)
	WeatherTemp  string // e.g. "+25°"
	WeatherPress string // atmospheric pressure, e.g. "742 мм" ("" when unavailable)
	Rates        []Rate // empty when unavailable
	Social       []SocialLink
}

// InfoBar fetches and caches the weather and exchange rates in the background.
type InfoBar struct {
	mu          sync.RWMutex
	rates       []Rate
	weatherIc   string
	weatherTmp  string
	weatherPres string

	httpc    *http.Client
	log      *zap.Logger
	lat, lon float64
	social   []SocialLink
}

// socialLinks returns the four social profiles in a fixed order so their icons
// always render. A profile with no configured URL yet gets a "#" placeholder
// (shown but non-navigating) until a real URL is filled in.
func socialLinks(cfg config.SocialConfig) []SocialLink {
	out := make([]SocialLink, 0, 4)
	for _, s := range []SocialLink{
		{"telegram", strings.TrimSpace(cfg.Telegram)},
		{"instagram", strings.TrimSpace(cfg.Instagram)},
		{"youtube", strings.TrimSpace(cfg.YouTube)},
		{"facebook", strings.TrimSpace(cfg.Facebook)},
	} {
		if s.URL == "" {
			s.URL = "#"
		}
		out = append(out, s)
	}
	return out
}

// NewInfoBar builds the provider. Weather defaults to Almaty.
func NewInfoBar(log *zap.Logger, social []SocialLink) *InfoBar {
	return &InfoBar{
		httpc:  &http.Client{Timeout: 6 * time.Second},
		log:    log,
		lat:    43.2389,
		lon:    76.8897,
		social: social,
	}
}

// Snapshot returns the cached bar data for the given day string.
func (b *InfoBar) Snapshot(today string) InfoBarData {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return InfoBarData{Today: today, WeatherIcon: b.weatherIc, WeatherTemp: b.weatherTmp, WeatherPress: b.weatherPres, Rates: b.rates, Social: b.social}
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
			if n++; n%12 == 0 {
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

// refreshRates pulls USD/EUR/RUB from the National Bank of Kazakhstan feed.
func (b *InfoBar) refreshRates(ctx context.Context) {
	body, err := b.get(ctx, "https://nationalbank.kz/rss/rates_all.xml")
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

// refreshWeather pulls the current temperature and condition from open-meteo.
func (b *InfoBar) refreshWeather(ctx context.Context) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,weather_code,pressure_msl&timezone=auto", b.lat, b.lon)
	body, err := b.get(ctx, url)
	if err != nil {
		b.log.Warn("infobar weather fetch", zap.Error(err))
		return
	}
	var doc struct {
		Current struct {
			Temp  float64 `json:"temperature_2m"`
			Code  int     `json:"weather_code"`
			Press float64 `json:"pressure_msl"` // hPa
		} `json:"current"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		b.log.Warn("infobar weather parse", zap.Error(err))
		return
	}
	sign := ""
	t := int(doc.Current.Temp)
	if t > 0 {
		sign = "+"
	}
	// Pressure in mm Hg (hPa × 0.750062) — just the number; the template appends
	// the unit in the active UI language ("мм рт.ст." / "mmHg" / "мм сын.бағ.").
	press := ""
	if doc.Current.Press > 0 {
		press = fmt.Sprintf("%d", int(doc.Current.Press*0.750062+0.5))
	}
	b.mu.Lock()
	b.weatherIc = weatherIconName(doc.Current.Code)
	b.weatherTmp = fmt.Sprintf("%s%d°", sign, t)
	b.weatherPres = press
	b.mu.Unlock()
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

// curSymbol maps a currency code to its symbol for the compact bar.
func curSymbol(code string) string {
	switch code {
	case "USD":
		return "$"
	case "EUR":
		return "€"
	case "RUB":
		return "₽"
	default:
		return code
	}
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
