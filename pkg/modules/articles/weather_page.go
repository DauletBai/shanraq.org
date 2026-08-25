package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// The forecast page.
//
// The strip at the top of every page has shown one temperature for a while, and
// it was the second thing readers clicked after the exchange rate — except it
// led nowhere. A temperature answers "what is it now" and immediately provokes
// "and tomorrow", which is the question people actually came with.
//
// Every place in the reference with coordinates gets its own address, which is
// the point: somebody searching for the weather in Kachar cannot find a page
// that does not exist. That is the same argument the place pages were built on,
// and it applies here with more force, because far more people look up the
// weather in their town than look up news about it.
//
// The data is fetched server-side and cached. The content security policy stops
// the browser reaching an external API, and a forecast that only some readers
// can see is worse than none.

// wxDefaultLat/Lon is the city the strip's temperature has always come from, so
// the bare /weather address lands where that reading was taken.
const (
	wxDefaultLat = 43.2389
	wxDefaultLon = 76.8897
)

// wxForecastDays is how far ahead the page looks. Beyond a week a forecast is
// not information, it is decoration.
const wxForecastDays = 7

// wxHours is how many hours of the hourly series are drawn: two days, which is
// as far as an hourly forecast is worth reading.
const wxHours = 48

// wxTTL is how long a fetched forecast is reused. The model behind it updates
// hourly, so anything shorter spends requests on the same numbers.
const wxTTL = 30 * time.Minute

// WxNow is the current conditions block.
type WxNow struct {
	Temp     string
	Feels    string
	Icon     string
	Desc     string
	Humidity string
	Wind     string
	Pressure string
	At       string
}

// WxDay is one row of the week's forecast.
type WxDay struct {
	Date    string
	Weekday string
	Icon    string
	Desc    string
	Min     string
	Max     string
	Precip  string
	Wind    string
	Sunrise string
	Sunset  string
	// Today marks the current day so the row can be highlighted.
	Today bool
}

// WeatherPage backs /weather and /weather/{slug}.
type WeatherPage struct {
	Base
	Desc      string
	PlaceName string
	// PlaceLabel is the full address of the place: "Качар, Костанайская область".
	PlaceLabel string
	Slug       string
	HasData    bool

	Now  WxNow
	Days []WxDay

	// Temp is the hourly temperature over the next two days; Rain the chance of
	// precipitation over the same hours.
	Temp    FxChart
	Rain    FxChart
	HasRain bool

	// Nearby are other places a reader may want instead, largest first.
	Nearby []GeoNode
	// Updated is when this forecast was fetched.
	Updated string
}

// wxCacheEntry is one place's cached forecast.
type wxCacheEntry struct {
	at   time.Time
	page WeatherPage
}

// wxCache holds forecasts by place and language.
var wxCache = struct {
	mu sync.Mutex
	m  map[string]wxCacheEntry
}{m: map[string]wxCacheEntry{}}

// handleWeather renders the forecast for a place.
func (m *Module) handleWeather(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	slug := strings.TrimSpace(chi.URLParam(r, "slug"))

	var (
		lat, lon    float64
		name, label string
		node        *GeoNode
	)
	if slug == "" {
		// A reader who told us where they live gets their own town, by way of
		// its real address rather than by rendering it here: one page per place
		// keeps the link shareable, the cache useful and the crawler happy.
		//
		// Temporary, not permanent: the answer depends on who is asking, and a
		// permanent redirect would teach every cache the first reader's town.
		if own := m.readerWeatherSlug(r); own != "" {
			http.Redirect(w, r, "/weather/"+own+"?lang="+lang, http.StatusFound)
			return
		}
		// Everyone else keeps the city the strip has always shown, so the link
		// in the header lands where its temperature came from.
		lat, lon = wxDefaultLat, wxDefaultLon
		name = T(lang, "wx.default_city")
	} else {
		n, err := m.geo.BySlug(r.Context(), slug, lang)
		if err != nil {
			m.rt.Logger.Error("weather place", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		// A place with no coordinates cannot have a forecast, and inventing one
		// from its region would be a forecast for somewhere else.
		if n == nil || n.Lat == nil || n.Lng == nil {
			http.NotFound(w, r)
			return
		}
		node, lat, lon, name = n, *n.Lat, *n.Lng, n.Name
	}

	page, ok := m.weatherCached(r.Context(), slug, lang, lat, lon)
	page.PlaceName = name
	page.Slug = slug
	if !ok {
		m.rt.Logger.Warn("weather unavailable", zap.String("slug", slug))
	}

	if node != nil {
		if id, err := uuid.Parse(node.ID); err == nil {
			if l, err := m.geo.PlaceLabel(r.Context(), id, lang); err == nil {
				label = l
			}
			if kids, err := m.geo.Children(r.Context(), id, lang); err == nil {
				page.Nearby = wxWithCoords(kids)
			}
		}
	}
	page.PlaceLabel = label

	title := fmt.Sprintf(T(lang, "wx.title_place"), name)
	page.Base = m.base(r, title, lang)
	page.Desc = fmt.Sprintf(T(lang, "wx.desc_place"), name)
	if slug != "" {
		page.Base.CanonURL = canonURL("/weather/"+slug, "", lang)
		page.Base.LangLinks = langLinks("/weather/"+slug, "")
	}
	m.render(w, "weather", page)
}

// wxWithCoords keeps only the places a forecast can be built for.
func wxWithCoords(in []GeoNode) []GeoNode {
	out := make([]GeoNode, 0, len(in))
	for _, n := range in {
		if n.Lat != nil && n.Lng != nil {
			out = append(out, n)
		}
	}
	if len(out) > 12 {
		out = out[:12]
	}
	return out
}

// weatherCached returns a forecast, fetching it at most twice an hour per place.
func (m *Module) weatherCached(ctx context.Context, slug, lang string, lat, lon float64) (WeatherPage, bool) {
	key := slug + "|" + lang
	wxCache.mu.Lock()
	e, ok := wxCache.m[key]
	wxCache.mu.Unlock()
	if ok && time.Since(e.at) < wxTTL {
		return e.page, true
	}

	page, err := m.fetchForecast(ctx, lang, lat, lon)
	if err != nil {
		m.rt.Logger.Warn("forecast fetch", zap.Error(err))
		// A stale forecast beats an empty page: yesterday's week is still
		// roughly this week, and the page says when it was taken.
		if ok {
			return e.page, true
		}
		return WeatherPage{}, false
	}
	wxCache.mu.Lock()
	wxCache.m[key] = wxCacheEntry{at: time.Now(), page: page}
	wxCache.mu.Unlock()
	return page, true
}

// wxResponse is the shape of the open-meteo reply we read.
type wxResponse struct {
	Current struct {
		Time     string  `json:"time"`
		Temp     float64 `json:"temperature_2m"`
		Humidity float64 `json:"relative_humidity_2m"`
		Feels    float64 `json:"apparent_temperature"`
		Code     int     `json:"weather_code"`
		Pressure float64 `json:"pressure_msl"`
		Wind     float64 `json:"wind_speed_10m"`
	} `json:"current"`
	Hourly struct {
		Time []string  `json:"time"`
		Temp []float64 `json:"temperature_2m"`
		Code []int     `json:"weather_code"`
		Rain []float64 `json:"precipitation_probability"`
		Wind []float64 `json:"wind_speed_10m"`
	} `json:"hourly"`
	Daily struct {
		Time    []string  `json:"time"`
		Code    []int     `json:"weather_code"`
		Max     []float64 `json:"temperature_2m_max"`
		Min     []float64 `json:"temperature_2m_min"`
		Sunrise []string  `json:"sunrise"`
		Sunset  []string  `json:"sunset"`
		Precip  []float64 `json:"precipitation_sum"`
		Wind    []float64 `json:"wind_speed_10m_max"`
	} `json:"daily"`
}

// fetchForecast pulls and shapes one place's forecast.
func (m *Module) fetchForecast(ctx context.Context, lang string, lat, lon float64) (WeatherPage, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast"+
		"?latitude=%.4f&longitude=%.4f"+
		"&current=temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,pressure_msl,wind_speed_10m"+
		"&hourly=temperature_2m,weather_code,precipitation_probability,wind_speed_10m"+
		"&daily=weather_code,temperature_2m_max,temperature_2m_min,sunrise,sunset,precipitation_sum,wind_speed_10m_max"+
		"&forecast_days=%d&timezone=auto", lat, lon, wxForecastDays)

	body, err := macroFetch(ctx, url)
	if err != nil {
		return WeatherPage{}, err
	}
	var doc wxResponse
	if err := json.Unmarshal(body, &doc); err != nil {
		return WeatherPage{}, fmt.Errorf("forecast: %w", err)
	}
	if len(doc.Daily.Time) == 0 {
		return WeatherPage{}, fmt.Errorf("forecast came back empty")
	}

	page := WeatherPage{HasData: true, Updated: wxClock(doc.Current.Time)}
	page.Now = WxNow{
		Temp:     wxTemp(doc.Current.Temp),
		Feels:    wxTemp(doc.Current.Feels),
		Icon:     weatherIconName(doc.Current.Code),
		Desc:     wxDescribe(doc.Current.Code, lang),
		Humidity: fmt.Sprintf("%.0f %%", doc.Current.Humidity),
		Wind:     wxWind(doc.Current.Wind, lang),
		Pressure: wxPressure(doc.Current.Pressure, lang),
		At:       wxClock(doc.Current.Time),
	}

	// The hourly series starts at midnight of the current day, so the part
	// already behind us is dropped: a forecast for this morning is not a
	// forecast.
	from := wxHourIndex(doc.Hourly.Time, doc.Current.Time)
	temps, rain := []FxPoint{}, []FxPoint{}
	for i := from; i < len(doc.Hourly.Time) && i < from+wxHours; i++ {
		t, err := time.Parse("2006-01-02T15:04", doc.Hourly.Time[i])
		if err != nil {
			continue
		}
		temps = append(temps, FxPoint{Day: t, Value: doc.Hourly.Temp[i]})
		if i < len(doc.Hourly.Rain) {
			rain = append(rain, FxPoint{Day: t, Value: doc.Hourly.Rain[i]})
		}
	}
	if len(temps) > 2 {
		page.Temp = fxBuildChartWith(temps, "hours", lang, fxChartOpts{
			Hourly: true, Unit: "°C", Format: func(v float64) string { return wxTemp(v) },
			AxisFormat: func(v float64) string { return fxFormat(v, 0) + "°" },
		})
	}
	// A flat zero line teaches nothing: the chance of rain is only worth a
	// frame when there is some.
	if len(rain) > 2 && wxAny(rain) {
		page.HasRain = true
		page.Rain = fxBuildChartWith(rain, "hours", lang, fxChartOpts{
			Hourly: true, Unit: "%", Format: func(v float64) string { return fxFormat(v, 0) + " %" },
			AxisFormat: func(v float64) string { return fxFormat(v, 0) + " %" },
		})
	}

	today := time.Now().Format("2006-01-02")
	for i := range doc.Daily.Time {
		d, err := time.Parse("2006-01-02", doc.Daily.Time[i])
		if err != nil {
			continue
		}
		row := WxDay{
			Date:    fmt.Sprintf("%d %s", d.Day(), fxMonthShort(d.Month(), lang)),
			Weekday: wxWeekday(d.Weekday(), lang),
			Icon:    weatherIconName(doc.Daily.Code[i]),
			Desc:    wxDescribe(doc.Daily.Code[i], lang),
			Min:     wxTemp(doc.Daily.Min[i]),
			Max:     wxTemp(doc.Daily.Max[i]),
			Wind:    wxWind(doc.Daily.Wind[i], lang),
			Sunrise: wxClock(doc.Daily.Sunrise[i]),
			Sunset:  wxClock(doc.Daily.Sunset[i]),
			Today:   doc.Daily.Time[i] == today,
		}
		if p := doc.Daily.Precip[i]; p > 0 {
			row.Precip = fxFormat(p, 1) + " " + T(lang, "wx.mm")
		}
		page.Days = append(page.Days, row)
	}
	return page, nil
}

// wxAny reports whether any point is above zero.
func wxAny(pts []FxPoint) bool {
	for _, p := range pts {
		if p.Value > 0 {
			return true
		}
	}
	return false
}

// wxHourIndex finds the hourly slot the current reading belongs to.
func wxHourIndex(times []string, now string) int {
	if len(now) < 13 {
		return 0
	}
	hour := now[:13] // "2026-08-25T07"
	for i, t := range times {
		if len(t) >= 13 && t[:13] >= hour {
			return i
		}
	}
	return 0
}

// wxTemp prints a temperature with its sign, the way a forecast is read.
func wxTemp(v float64) string {
	s := fmt.Sprintf("%.0f", v)
	if v > 0 {
		s = "+" + s
	}
	// A rounded −0.4 prints "-0", which reads as a typo rather than a
	// temperature.
	if s == "-0" || s == "+0" {
		s = "0"
	}
	return s + "°"
}

// wxWind prints wind speed. The API answers in km/h; metres per second is what
// a forecast is read in here.
func wxWind(kmh float64, lang string) string {
	ms := math.Round(kmh / 3.6)
	// Still air is a state, not a measurement of zero: "0 m/s" reads like a
	// missing figure, and every forecast in the language calls this calm.
	if ms < 1 {
		return T(lang, "wx.calm")
	}
	return fmt.Sprintf("%.0f %s", ms, T(lang, "wx.ms"))
}

// wxPressure converts hectopascals to the millimetres of mercury a barometer in
// this part of the world is marked in.
func wxPressure(hpa float64, lang string) string {
	return fmt.Sprintf("%.0f %s", hpa*0.750062, T(lang, "wx.pressure_unit"))
}

// wxClock keeps the time out of an ISO stamp.
func wxClock(iso string) string {
	if i := strings.IndexByte(iso, 'T'); i >= 0 && len(iso) >= i+6 {
		return iso[i+1 : i+6]
	}
	return ""
}

// wxWeekday names the day in the reader's language.
func wxWeekday(d time.Weekday, lang string) string {
	names := map[string][7]string{
		LangKZ: {"жексенбі", "дүйсенбі", "сейсенбі", "сәрсенбі", "бейсенбі", "жұма", "сенбі"},
		LangRU: {"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"},
		LangEN: {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	}
	row, ok := names[lang]
	if !ok {
		row = names[LangRU]
	}
	return row[int(d)]
}

// wxDescribe turns a WMO weather code into words. The codes are a published
// standard; the groupings match the icons the strip already uses.
func wxDescribe(code int, lang string) string {
	key := "wx.c_cloud"
	switch {
	case code == 0:
		key = "wx.c_clear"
	case code <= 2:
		key = "wx.c_partly"
	case code == 3:
		key = "wx.c_cloud"
	case code == 45 || code == 48:
		key = "wx.c_fog"
	case code >= 51 && code <= 57:
		key = "wx.c_drizzle"
	case code >= 61 && code <= 67:
		key = "wx.c_rain"
	case code >= 71 && code <= 77:
		key = "wx.c_snow"
	case code >= 80 && code <= 82:
		key = "wx.c_showers"
	case code == 85 || code == 86:
		key = "wx.c_snow"
	case code >= 95:
		key = "wx.c_storm"
	}
	return T(lang, key)
}

// WeatherPlaceFor picks the place whose forecast answers a reader who lives in
// `place`.
//
// The place itself when it has coordinates — a person in Kachar wants Kachar.
// Otherwise its largest settlement: a region has no coordinates of its own, and
// somebody who told us only "Kostanay region" is better served by the forecast
// for Kostanay than by one for a city a thousand kilometres away.
//
// Returns an empty slug when nothing inside the place can be forecast, and the
// caller then falls back to the default city.
func (s *GeoStore) WeatherPlaceFor(ctx context.Context, place uuid.UUID) (string, error) {
	var slug string
	err := s.db.QueryRow(ctx, `
		WITH RECURSIVE down AS (
			SELECT id, parent_id, slug, lat, lng, population, 0 AS depth
			FROM geo_nodes WHERE id = $1
			UNION ALL
			SELECT g.id, g.parent_id, g.slug, g.lat, g.lng, g.population, down.depth + 1
			FROM geo_nodes g JOIN down ON g.parent_id = down.id
			WHERE down.depth < 4
		)
		SELECT COALESCE(slug,'') FROM down
		WHERE lat IS NOT NULL AND lng IS NOT NULL AND slug IS NOT NULL
		ORDER BY depth, population DESC NULLS LAST
		LIMIT 1`, place).Scan(&slug)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", nil
		}
		return "", fmt.Errorf("weather place for: %w", err)
	}
	return slug, nil
}

// readerWeatherSlug is the place this reader's own forecast lives at, or "" for
// a guest and for anyone who never said where they live.
//
// It comes from the place they chose in their profile — never from their
// address on the network. The guide promises we do not work out anyone's town,
// and a convenience is not a reason to start.
func (m *Module) readerWeatherSlug(r *http.Request) string {
	uid, ok := m.authorID(r)
	if !ok || m.geo == nil {
		return ""
	}
	place, err := m.geo.UserPlace(r.Context(), uid)
	if err != nil || place == nil {
		return ""
	}
	slug, err := m.geo.WeatherPlaceFor(r.Context(), *place)
	if err != nil {
		m.rt.Logger.Warn("reader weather place", zap.Error(err))
		return ""
	}
	return slug
}

// WeatherPlaces lists every place a forecast page exists for: those with both
// coordinates and an address of their own.
//
// Russian cities are included along with Kazakhstani ones. The reference holds
// them because listings need them, and readers here have reasons to look up the
// weather across the border that a publication has no business second-guessing.
func (s *GeoStore) WeatherPlaces(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `
		SELECT slug FROM geo_nodes
		WHERE slug IS NOT NULL AND lat IS NOT NULL AND lng IS NOT NULL
		ORDER BY population DESC NULLS LAST, slug`)
	if err != nil {
		return nil, fmt.Errorf("weather places: %w", err)
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		out = append(out, slug)
	}
	return out, rows.Err()
}
