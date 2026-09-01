package articles

import (
	"html/template"
	"strconv"
	"strings"
	"time"

	"shanraq.org/web"
)

// templateFuncs is the single source of the template function map. Both the
// live module (Init) and the template tests use it, so a new helper can never
// be available in one place and missing in the other.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"t": T,
		// Every stylesheet and script must go through this: a bare path keeps
		// serving from cache for a day after a deploy, so new markup lands on
		// old CSS. See web.AssetURL.
		"asset": web.AssetURL,
		// Brand mark for an analytics row, "" for rows that name no brand.
		"brandicon":        brandIcon,
		"svcOff":           serviceLinkOff, // is a service's entry link disabled?
		"svcMsg":           serviceLinkMsg, // its localized "unavailable" tooltip
		"label":            func(l string) string { return LangLabels[l] },
		"langName":         func(l string) string { return LangNames[l] },
		"langs":            func() []string { return Langs },
		"categories":       func() []string { return Categories },
		"wallMaterials":    func() []string { return WallMaterials },
		"maxPhotos":        func() int { return maxListingPhotos },
		"maxDocs":          func() int { return maxListingDocs },
		"wallKey":          WallMaterialKey,
		"editorCategories": func() []string { return append([]string{CategoryGeneral}, Categories...) },
		"subcats":          func(cat string) []string { return Subcats(cat) },
		"dealTypes":        func() []string { return DealTypes },
		"propertyTypes":    func() []string { return PropertyTypes },
		"amenities":        AmenityKeys,
		"roomTypes":        RoomTypeKeys,
		"bannerDays":       BannerDays,
		// The report count that hides a listing, so the seller's warning quotes
		// the real threshold instead of a number typed into a translation.
		"reportHideAt":      func() int { return reportMinReports },
		"bannerPrice":       BannerPrice,
		"adSurfaces":        AdSurfaces,
		"adDurations":       AdDurations,
		"adFormats":         AdFormats,
		"adSurfaceFmtPrice": AdSurfaceFormatPrice,
		"adRatesJSON":       AdRatesJSON,
		"surfaceLabel":      SurfaceLabelKey,
		"adFormatSlots":     AdFormatSlots,
		"money":             money,
		// Templates count from zero and people count from one; screen-reader
		// labels are read by people.
		"inc": func(i int) int { return i + 1 },
		"mul": func(a, b int) int { return a * b },
		// A cover that is a drawing rather than a photograph. The hero prints a
		// headline across the picture, and a diagram's own labels fight it.
		"isvector": func(s string) bool { return strings.HasSuffix(strings.ToLower(s), ".svg") },
		// The reader's report names one of the site's published rules — the same
		// list the checker used — so a report is a claim about a rule, not a
		// second opinion about the topic.
		"reviewRules": func() []string { return ReviewRules },
		// The tabs and the translate button must name languages the same way,
		// or the reader has to work out that "3 языка" and "Қазақша" are about
		// the same thing.
		"langNames": func() map[string]string { return LangNames },
		"langList": func(ls []string) string {
			out := make([]string, 0, len(ls))
			for _, l := range ls {
				out = append(out, LangNames[l])
			}
			return strings.Join(out, ", ")
		},
		"compactNum":       compactNum, // 1234 → "1,2 тыс." for tight meta rows
		"shortAuthor":      shortAuthor,
		"hasSuffix":        strings.HasSuffix,
		"ogLocale":         ogLocale,
		"htmlLang":         htmlLang,
		"curSymbol":        curSymbol,
		"icon":             icon,
		"roomIcon":         roomIcon,
		"amenityIcon":      amenityIcon,
		"catIcon":          catIcon,
		"firstN":           firstStrings,
		"countryFlag":      countryFlag,
		"countryMark":      countryMark,
		"countryFlagEmoji": countryFlagEmoji,
		"kilo":             kilo,
		"liveSocial":       liveSocial, // social profiles that aren't "#" placeholders
		"withUTM":          withUTM,
		"dict":             dict,
		"year":             func() int { return time.Now().Year() },
		"markdown":         RenderMarkdown,
		"fmtDate": func(t time.Time) string {
			if t.IsZero() {
				return "—"
			}
			return t.Format("02.01.06")
		},
		"fmtDatePtr": func(t *time.Time) string {
			if t == nil || t.IsZero() {
				return "—"
			}
			return t.Format("02.01.06")
		},
	}
}

// kilo shortens a figure to thousands so the scale beside a chart needs no
// room held open for digits that are not there yet. A gutter wide enough for
// six of them is mostly empty at every reading that has ever been taken, and
// widening it the day traffic grows is a change nobody would remember to make.
//
// The fraction is kept only while it separates two readings -- 1,5 tells you
// something 2 does not, 12,3 does not -- and the decimal mark follows the
// language rather than the machine.
func kilo(lang string, n int64) string {
	if n > -1000 && n < 1000 {
		return strconv.FormatInt(n, 10)
	}
	var s string
	if v := float64(n) / 1000; v > -10 && v < 10 {
		s = strings.TrimSuffix(strconv.FormatFloat(v, 'f', 1, 64), ".0")
	} else {
		s = strconv.FormatInt(n/1000, 10)
	}
	if lang != "en" {
		s = strings.Replace(s, ".", ",", 1)
	}
	return s + T(lang, "tc.kilo")
}
