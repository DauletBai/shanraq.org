package articles

import (
	"html/template"
	"strconv"
	"strings"

	"shanraq.org/pkg/site"
)

// templateFuncs are the helpers that know what this module's pages are about.
// Everything the frame itself needs -- the dictionary, the languages, the icon
// set, the small formatting verbs -- comes from site.BaseFuncs, which the
// renderer already holds. Both the live module (Init) and the template tests
// go through the same renderer, so a helper can never be available in one
// place and missing in the other.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"wallMaterials":    func() []string { return WallMaterials },
		"maxPhotos":        func() int { return maxListingPhotos },
		"maxDocs":          func() int { return maxListingDocs },
		"wallKey":          WallMaterialKey,
		"editorCategories": func() []string { return append([]string{CategoryGeneral}, Categories...) },
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
		// The reader's report names one of the site's published rules — the same
		// list the checker used — so a report is a claim about a rule, not a
		// second opinion about the topic.
		"reviewRules":      func() []string { return ReviewRules },
		"compactNum":       compactNum, // 1234 → "1,2 тыс." for tight meta rows
		"shortAuthor":      shortAuthor,
		"countryFlag":      countryFlag,
		"countryMark":      countryMark,
		"countryFlagEmoji": countryFlagEmoji,
		"kilo":             kilo,
		"markdown":         RenderMarkdown,
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
	return s + site.T(lang, "tc.kilo")
}
