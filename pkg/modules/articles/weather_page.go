package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
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

	// Nearby are other places a reader may want instead, largest first: a
	// region's own settlements, never the whole reference.
	Nearby []GeoNode
	// Population is the place's own, formatted, empty when unknown; PopYear the
	// census it was counted in. They are here because a page that says only what
	// the sky is doing says the same thing as every other forecast page, and a
	// search engine reading four hundred of those sees one page repeated.
	Population string
	PopYear    int
	// Lat and Lng centre the map; Radar is the precipitation tile template,
	// empty when the service could not be reached.
	Lat, Lng float64
	Radar    string
	// Updated is when this forecast was fetched.
	Updated string
	// FromPoint explains a place the reader never named: they pressed the map,
	// and this is the nearest settlement to where they pressed. Empty on a page
	// somebody opened by its own address.
	FromPoint string
}

// WxSwap is the answer to a press on the map: the whole of the page, for the
// place that was pressed.
//
// The map used to answer beneath itself, in a block of its own, and that was the
// wrong answer to the question people were asking. A reader in Kostanay who
// presses Almaty wants Almaty's page — its heading, its degrees now, its two
// charts, its week, its neighbours — not Kostanay's page with a second table
// sewn onto the bottom. So the server renders the same parts the page is built
// from, the browser puts them where the old ones stood, and the address bar
// changes to the place that answered: what is on the screen and what the link
// says stay the same thing.
type WxSwap struct {
	Page WeatherPage
	// URL is where this answer lives — the place's own page, or the coordinates
	// when the point is far from anything named. It is pressed into the reader's
	// history, so the link stays shareable and Back still leads back.
	URL string
	// Title is the browser tab's, spelled the way site_head spells it.
	Title string
}

// wxSiteTitle is the tail of every page title. site_head writes it into the
// markup; here it is written by hand, because the tab has to keep saying the
// same thing after the page changes underneath it.
const wxSiteTitle = " \u00b7 Shanraq.org"

// wxPlace is what a weather page answers about: a point on Earth, a name for it,
// and the reference entry behind it when there is one.
type wxPlace struct {
	slug string
	name string
	node *GeoNode
	lat  float64
	lon  float64
	// note explains a place the reader never named: they pressed the map, and
	// this is how far from that press the answer was taken.
	note string
}

// key is what this place's forecast is cached under. A place with an address of
// its own shares the cache with its page; anything else is kept by coordinates,
// so two readers standing in different steppes are not served each other's sky.
func (p wxPlace) key() string {
	if p.slug != "" {
		return p.slug
	}
	return fmt.Sprintf("@%.2f,%.2f", p.lat, p.lon)
}

// url is the address this answer lives at.
func (p wxPlace) url() string {
	if p.slug != "" {
		return "/weather/" + p.slug
	}
	return fmt.Sprintf("/weather?at=%.2f,%.2f", p.lat, p.lon)
}

// wxCoords reads a pair of coordinates and refuses anything that is not a place
// on Earth.
func wxCoords(latS, lonS string) (float64, float64, bool) {
	lat, errLat := strconv.ParseFloat(strings.TrimSpace(latS), 64)
	lon, errLon := strconv.ParseFloat(strings.TrimSpace(lonS), 64)
	if errLat != nil || errLon != nil || !wxPointOK(lat, lon) {
		return 0, 0, false
	}
	return lat, lon, true
}

// wxPickPixels is how wide a press on the map is taken to be. A click is a
// finger, not a coordinate.
const wxPickPixels = 24

// wxPickMinKm keeps the reach of a press sane at the closest zooms: a village is
// still a village when the map is showing streets.
const wxPickMinKm = 3

// wxPickKm is how far a press may reach for a place, in kilometres, at the zoom
// the reader is looking at. With the whole country on the screen that same
// finger covers a hundred kilometres, and somebody who put it on the Almaty
// blob meant Almaty — not the village six kilometres from the pixel they hit.
// Zoomed in until streets show, it covers barely one, and then the village is
// exactly what they meant.
func wxPickKm(zoom int, lat float64) float64 {
	if zoom <= 0 {
		return 0
	}
	if zoom > 20 {
		zoom = 20
	}
	// Web Mercator: one pixel at the equator is 156.543 km at zoom zero, halved
	// with every step, and narrowed by the latitude.
	perPixel := 156.543 * math.Cos(lat*math.Pi/180) / math.Pow(2, float64(zoom))
	km := perPixel * wxPickPixels
	if km < wxPickMinKm {
		km = wxPickMinKm
	}
	if km > wxNearKm {
		km = wxNearKm
	}
	return km
}

// pointPlace turns coordinates into the place a page can answer about: the
// place the press was aimed at, the nearest settlement within wxNearKm when it
// was aimed at nothing, or the point itself when the steppe around it is empty.
// zoom is the map's, zero when the coordinates came from something other than a
// press and there is no finger to be wide.
func (m *Module) pointPlace(ctx context.Context, lang string, lat, lon float64, zoom int) wxPlace {
	// Rounded before anything else: two decimals is about a kilometre, which is
	// one forecast, and it keeps a dragged finger from asking for a thousand.
	lat = math.Round(lat*100) / 100
	lon = math.Round(lon*100) / 100
	where := fmt.Sprintf("%.2f, %.2f", lat, lon)
	p := wxPlace{lat: lat, lon: lon, name: where}
	if m.geo == nil {
		return p
	}
	var (
		node  GeoNode
		km    float64
		found bool
		err   error
	)
	if reach := wxPickKm(zoom, lat); reach > 0 {
		node, km, found, err = m.geo.Prominent(ctx, lang, lat, lon, reach)
	}
	// Nothing was aimed at: the press landed in open steppe, and the nearest
	// town — however far — is still a better answer than two numbers.
	if err == nil && !found {
		node, km, found, err = m.geo.Nearest(ctx, lang, lat, lon, wxNearKm)
	}
	if err != nil {
		m.rt.Logger.Warn("weather point place", zap.Error(err))
		return p
	}
	if !found || node.Lat == nil || node.Lng == nil {
		return p
	}
	// The forecast is taken at the town rather than at the pixel: this is that
	// town's page now, and a page must not disagree with its own address.
	p.node, p.slug, p.name = &node, node.Slug, node.Name
	p.lat, p.lon = *node.Lat, *node.Lng
	if km >= 1 {
		p.note = fmt.Sprintf(T(lang, "wx.point_from"), where,
			fmt.Sprintf("%.0f %s", km, T(lang, "wx.km")))
	} else {
		p.note = fmt.Sprintf(T(lang, "wx.point_here"), where)
	}
	return p
}

// weatherAt builds everything a weather page shows about a place: the forecast,
// the address it sits at, its population and its neighbours. The page and the
// map's answer are the same page, so both are built here.
func (m *Module) weatherAt(ctx context.Context, lang string, p wxPlace) WeatherPage {
	page, ok := m.weatherCached(ctx, p.key(), lang, p.lat, p.lon)
	if !ok {
		m.rt.Logger.Warn("weather unavailable", zap.String("place", p.key()))
	}
	page.PlaceName, page.Slug, page.FromPoint = p.name, p.slug, p.note
	page.Lat, page.Lng = p.lat, p.lon
	page.Radar = m.radarTiles(ctx)
	if p.node == nil || m.geo == nil {
		return page
	}
	id, err := uuid.Parse(p.node.ID)
	if err != nil {
		return page
	}
	if label, err := m.geo.PlaceLabel(ctx, id, lang); err == nil {
		page.PlaceLabel = label
	}
	if kids, err := m.geo.Children(ctx, id, lang); err == nil {
		page.Nearby = wxWithCoords(kids)
	}
	// A settlement has no children; its region's other places serve.
	if len(page.Nearby) == 0 {
		if sib, err := m.geo.Siblings(ctx, id, lang); err == nil {
			page.Nearby = wxWithCoords(sib)
		}
	}
	if len(page.Nearby) > wxNearbyMax {
		page.Nearby = page.Nearby[:wxNearbyMax]
	}
	if pop, year, err := m.geo.PlaceFacts(ctx, id); err == nil && pop > 0 {
		page.Population, page.PopYear = wxGroupDigits(pop), year
	}
	return page
}

// handleWeatherPoint answers a press on the map with the page for that place.
func (m *Module) handleWeatherPoint(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	lat, lon, ok := wxCoords(r.URL.Query().Get("lat"), r.URL.Query().Get("lon"))
	if !ok {
		http.Error(w, "bad point", http.StatusBadRequest)
		return
	}
	zoom, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("z")))
	place := m.pointPlace(r.Context(), lang, lat, lon, zoom)
	page := m.weatherAt(r.Context(), lang, place)
	// The parts rendered here are the page's own, and they ask the templates for
	// the reader's language the way the whole page does.
	page.Base = Base{Lang: lang}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// The answer is as fresh as the forecast behind it and no fresher.
	w.Header().Set("Cache-Control", "public, max-age=600")
	m.render(w, "wx_swap", WxSwap{
		Page:  page,
		URL:   place.url(),
		Title: fmt.Sprintf(T(lang, "wx.title_place"), place.name) + wxSiteTitle,
	})
}

// wxNearKm is how far the reference is searched for a name to call a point by.
// Beyond it a "nearest" settlement says nothing useful: the steppe is wide, and
// a hundred kilometres away is a different sky.
const wxNearKm = 80

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

	// The picker is a plain GET form, and a form cannot post into a path
	// segment — so it arrives as a query and is sent on to the real address.
	if p := strings.TrimSpace(r.URL.Query().Get("p")); p != "" && slug == "" {
		http.Redirect(w, r, "/weather/"+p+"?lang="+lang, http.StatusFound)
		return
	}

	var place wxPlace
	switch {
	case slug != "":
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
		place = wxPlace{slug: slug, name: n.Name, node: n, lat: *n.Lat, lon: *n.Lng}

	// A point in the address: somebody shared the answer the map gave them for
	// a spot the reference does not name. It is the same for everybody who
	// opens it, so unlike the bare address it is a page a cache may keep.
	case strings.Contains(r.URL.Query().Get("at"), ","):
		at := strings.SplitN(strings.TrimSpace(r.URL.Query().Get("at")), ",", 2)
		lat, lon, ok := wxCoords(at[0], at[1])
		if !ok {
			http.NotFound(w, r)
			return
		}
		place = m.pointPlace(r.Context(), lang, lat, lon, 0)
		// The reference does name it after all: the place has a page, and two
		// addresses for one forecast is one address too many.
		if place.slug != "" {
			http.Redirect(w, r, "/weather/"+place.slug+"?lang="+lang, http.StatusFound)
			return
		}

	default:
		// The bare address answers differently for different readers, so it is
		// nobody's to keep: a proxy that cached one reader's town would hand it
		// to the next city along.
		w.Header().Set("Cache-Control", "private, no-store")
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
		// A reader who told us nothing still came from somewhere. The address
		// gives a city — coarsely, and read for this answer only — and the
		// reference turns it into a place with a page of its own, which is
		// where the reader is sent. The strip above already shows that town's
		// temperature; landing on Almaty after pressing it was the page
		// contradicting its own header.
		pt, ok := m.readerPoint(r)
		if !ok {
			// Everyone else keeps the city the strip has always shown, so the
			// link in the header lands where its temperature came from.
			place = wxPlace{lat: wxDefaultLat, lon: wxDefaultLon, name: T(lang, "wx.default_city")}
			break
		}
		place = m.pointPlace(r.Context(), lang, pt.lat, pt.lon, 0)
		if place.slug != "" {
			http.Redirect(w, r, "/weather/"+place.slug+"?lang="+lang, http.StatusFound)
			return
		}
		// Nothing in the reference within eighty kilometres: the forecast is
		// still the reader's own, it simply has no page of its own to live at,
		// so it is rendered here under the name the address knows. The reader
		// pressed nothing, so there is nothing to explain about a press.
		place.note = ""
		if pt.city != "" {
			place.name = pt.city
		}
	}

	page := m.weatherAt(r.Context(), lang, place)
	name := place.name
	title := fmt.Sprintf(T(lang, "wx.title_place"), name)
	page.Base = m.base(r, title, lang)
	// The map needs Leaflet, and Leaflet is only shipped to the pages that draw
	// one: it is the heaviest asset on the site.
	page.NeedsMap = true
	page.Desc = fmt.Sprintf(T(lang, "wx.desc_place"), name)
	if slug != "" {
		page.Base.CanonURL = canonURL("/weather/"+slug, "", lang)
		page.Base.LangLinks = langLinks("/weather/"+slug, "")
	}
	m.render(w, "weather", page)
}

// wxNearbyMax caps the neighbours listed. A dozen is navigation; a hundred is
// the picker this replaced.
const wxNearbyMax = 12

// wxGroupDigits writes a population with thin gaps: 20 423, not 20423.
func wxGroupDigits(n int) string {
	s := fmt.Sprintf("%d", n)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteRune('\u00a0')
		}
		b.WriteRune(r)
	}
	return b.String()
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
			// No unit here: wxTemp already ends in a degree sign, and the
			// readout was printing "+30° °C".
			Hourly: true, Format: func(v float64) string { return wxTemp(v) },
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

// readerPlace is where a request came from, as far as an address database can
// tell: a city's middle, and its name in whatever language that database holds.
type readerPlace struct {
	lat, lon float64
	city     string
}

// readerPoint resolves the request's address to a place. False when there is no
// city database, or it has nothing for this address — in which case the page
// keeps its default city rather than guessing.
func (m *Module) readerPoint(r *http.Request) (readerPlace, bool) {
	if m.geoip == nil || !m.geoip.hasCity() {
		return readerPlace{}, false
	}
	lat, lon, city, ok := m.geoip.point(clientIP(r))
	if !ok || !wxPointOK(lat, lon) {
		return readerPlace{}, false
	}
	return readerPlace{lat: lat, lon: lon, city: city}, true
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

// The precipitation layer.
//
// RainViewer publishes radar tiles for the whole world without a key, but the
// address of the current frame changes every ten minutes and has to be looked
// up first. The browser cannot do that lookup — connect-src is 'self', and
// relaxing it to reach a weather service would open the page to every other
// one — so the server fetches it and hands the finished tile template down.
//
// The tiles themselves are ordinary images, which img-src already allows from
// any HTTPS host. So a live global weather layer costs no change to the content
// security policy and loads no third-party script.

// rvTTL is how long a frame address is reused. Frames appear about every ten
// minutes; asking more often spends requests on the same picture.
const rvTTL = 5 * time.Minute

var rvCache = struct {
	mu   sync.Mutex
	at   time.Time
	tmpl string
}{}

// radarTiles returns the Leaflet tile template for the newest radar frame, or
// "" when the service cannot be reached — the map then draws without the layer
// rather than not at all.
func (m *Module) radarTiles(ctx context.Context) string {
	rvCache.mu.Lock()
	if time.Since(rvCache.at) < rvTTL && rvCache.tmpl != "" {
		t := rvCache.tmpl
		rvCache.mu.Unlock()
		return t
	}
	rvCache.mu.Unlock()

	body, err := macroFetch(ctx, "https://api.rainviewer.com/public/weather-maps.json")
	if err != nil {
		m.rt.Logger.Warn("radar frames", zap.Error(err))
		return ""
	}
	var doc struct {
		Host  string `json:"host"`
		Radar struct {
			Past []struct {
				Path string `json:"path"`
			} `json:"past"`
		} `json:"radar"`
	}
	if err := json.Unmarshal(body, &doc); err != nil || len(doc.Radar.Past) == 0 {
		m.rt.Logger.Warn("radar frames unreadable", zap.Error(err))
		return ""
	}
	host := doc.Host
	if host == "" {
		host = "https://tilecache.rainviewer.com"
	}
	// 256-pixel tiles, colour scheme 2, smoothed with snow shown separately.
	tmpl := host + doc.Radar.Past[len(doc.Radar.Past)-1].Path + "/256/{z}/{x}/{y}/2/1_1.png"

	rvCache.mu.Lock()
	rvCache.at, rvCache.tmpl = time.Now(), tmpl
	rvCache.mu.Unlock()
	return tmpl
}
