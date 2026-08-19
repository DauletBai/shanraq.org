package articles

import (
	"html/template"
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
		"inc":              func(i int) int { return i + 1 },
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
		"countryFlagEmoji": countryFlagEmoji,
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
