package articles

import (
	"html/template"
	"io"
	"testing"
	"time"
)

// buildTemplates mirrors Module.Init template wiring so we can validate the
// embedded templates without a running server or database.
func buildTemplates(t *testing.T) *template.Template {
	t.Helper()
	funcs := template.FuncMap{
		"t":                T,
		"label":            func(l string) string { return LangLabels[l] },
		"langName":         func(l string) string { return LangNames[l] },
		"langs":            func() []string { return Langs },
		"categories":       func() []string { return Categories },
		"editorCategories": func() []string { return append([]string{CategoryGeneral}, Categories...) },
		"subcats":          func(cat string) []string { return Subcats(cat) },
		"dealTypes":        func() []string { return DealTypes },
		"propertyTypes":    func() []string { return PropertyTypes },
		"amenities":        AmenityKeys,
		"roomTypes":        RoomTypeKeys,
		"bannerDays":       BannerDays,
		"bannerPrice":      BannerPrice,
		"adSurfaces":       AdSurfaces,
		"adDurations":      AdDurations,
		"adSurfacePrice":   AdSurfacePriceOf,
		"adRatesJSON":      AdRatesJSON,
		"surfaceLabel":     SurfaceLabelKey,
		"adSlotCapacity":   AdSlotCapacity,
		"money":            money,
		"ogLocale":         ogLocale,
		"htmlLang":         htmlLang,
		"curSymbol":        curSymbol,
		"icon":             icon,
		"roomIcon":         roomIcon,
		"amenityIcon":      amenityIcon,
		"catIcon":          catIcon,
		"firstN":           firstStrings,
		"countryFlag":      countryFlag,
		"dict":             dict,
		"year":             func() int { return time.Now().Year() },
		"markdown":         RenderMarkdown,
		"fmtDate": func(tm time.Time) string {
			if tm.IsZero() {
				return "—"
			}
			return tm.Format("02.01.06")
		},
		"fmtDatePtr": func(tm *time.Time) string {
			if tm == nil || tm.IsZero() {
				return "—"
			}
			return tm.Format("02.01.06")
		},
	}
	tmpl, err := template.New("articles").Funcs(funcs).ParseFS(templateFiles, "templates/*.html")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}
	return tmpl
}

func TestTemplatesExecute(t *testing.T) {
	tmpl := buildTemplates(t)
	now := time.Now()

	// Exercise every UI language so a missing translation key surfaces.
	for _, lang := range Langs {
		base := Base{Title: "T", Lang: lang, Authed: true, ShowLangs: true, ActiveCat: "sport", ActiveSub: "football", LangLinks: langLinks("/", lang), Ads: demoAds(lang)}
		item := FeedItem{Slug: "s", Title: "Заголовок", Summary: "Краткое", AuthorName: "Автор",
			ServedLang: LangRU, Category: "politics", Subcategory: "elections", Published: &now, Views: 5, Score: 12, AvailableLangs: []string{LangRU, LangKZ}}

		cases := []struct {
			name string
			data any
		}{
			{"home", HomePage{Base: base, Featured: &item, Posts: []FeedItem{item}, Recent: []FeedItem{item}}},
			{"home", HomePage{Base: base, Featured: &item, Subscribed: true}}, // subscribe success
			{"home", HomePage{Base: base}}, // empty state
			{"article", ArticlePage{Base: base, Slug: "s", Title: "T", AuthorName: "A",
				ServedLang: LangRU, Category: "society", Body: RenderMarkdown("# Hi\n\nText"), Published: &now, Views: 1,
				Translated: true, IsAI: true, AvailableLangs: []string{LangRU},
				Score: 3, UserVote: 1, AuthorKarma: 42, CanVote: true, Recent: []FeedItem{item}, Subscribed: false}},
			{"page", StaticPage{Base: base, Body: RenderMarkdown("# Hi\n\nText [guide](/guide)")}},
			{"form", FormPage{Base: base, Mode: "login", Email: "a@b.c", Error: "err"}},
			{"form", FormPage{Base: base, Mode: "register"}},
			{"studio_dashboard", StudioPage{Base: base, Karma: 42, Stats: AuthorStats{
				TotalArticles: 2, Published: 1, Drafts: 1, TotalViews: 10,
				ViewsByLang: map[string]int64{LangRU: 7, LangKZ: 3, LangEN: 0},
			}, Articles: []StudioRow{{ID: "id", Slug: "s", Title: "T", Status: "published", Updated: now, Views: 4, Langs: []string{LangRU}}}}},
			{"studio_editor", EditorPage{Base: base, IsNew: true, OriginalLang: LangRU, Category: "society", Status: "draft", Fields: emptyFields()}},
			{"studio_editor", EditorPage{Base: base, IsNew: false, ArticleID: "id", OriginalLang: LangKZ, Category: "politics", Status: "published", Fields: emptyFields(), AIEnabled: true, Notice: "N"}},
			{"listings", ListingsPage{Base: base, ActiveDeal: "sale", ActiveType: "apartment",
				Facets: ListingFacets{Total: 6, Deal: map[string]int{"sale": 4, "rent": 2}, Type: map[string]int{"apartment": 2, "house": 1, "land": 1, "commercial": 1, "dacha": 1}},
				Listings: []*Listing{{
					ID: "id", DealType: "sale", PropertyType: "apartment", Country: "Казахстан", Region: "Алматы", City: "Алматы",
					Price: 18000000, Area: 72, Rooms: 3, Title: "Демо объявление трехкомнатной квартиры", PromotedUntil: &now,
					Images:    []string{"/static/demo/rooms/living.svg", "/static/demo/rooms/kitchen.svg", "/static/demo/rooms/bedroom.svg"},
					Amenities: []string{"furniture", "elevator", "internet"},
					RoomSpecs: []RoomSpec{{Type: "living", Area: 20}, {Type: "bedroom", Area: 15}, {Type: "bedroom", Area: 12}, {Type: "kitchen", Area: 7}, {Type: "bathroom", Area: 4}}}}}},
			{"listings", ListingsPage{Base: base}}, // empty state
			{"listing_new", ListingFormPage{Base: base, Values: ListingInput{DealType: "rent", PropertyType: "house", Country: "Казахстан"}, Error: "err"}},
			{"listing_view", ListingViewPage{Base: base, L: &Listing{
				ID: "id", DealType: "rent", PropertyType: "house", Country: "Казахстан", Region: "Астана", City: "Астана", Village: "Тельман",
				Price: 350000, Area: 120, Rooms: 4, Title: "Дом в аренду", Description: "Line1\nLine2", Contact: "+7 700 000 00 00", CoverURL: "http://x/y.jpg",
				Images: []string{"/static/demo/rooms/exterior.svg", "/static/demo/rooms/living.svg", "/static/demo/rooms/kitchen.svg"}}}},
			{"listing_view", ListingViewPage{Base: base, Owner: true, L: &Listing{ // cover-only fallback, owner controls
				ID: "id", DealType: "sale", PropertyType: "land", Title: "Участок", Contact: "+7 700 000 00 00", CoverURL: "http://x/y.jpg",
				ExpiresAt: now.Add(72 * time.Hour), PromotedUntil: &now, FeaturedUntil: &now}}},
			{"listing_my", MyListingsPage{Base: base, Listings: []*Listing{
				{ID: "id1", Title: "Активное", Price: 12000000, Region: "Алматы", ExpiresAt: now.Add(72 * time.Hour), FeaturedUntil: &now},
				{ID: "id2", Title: "Истёкшее", Price: 5000000, Region: "Астана", ExpiresAt: now.Add(-24 * time.Hour)},
			}}},
			{"listing_my", MyListingsPage{Base: base}}, // empty state
		}
		for _, c := range cases {
			if err := tmpl.ExecuteTemplate(io.Discard, c.name, c.data); err != nil {
				t.Errorf("execute %q (lang %s): %v", c.name, lang, err)
			}
		}
	}
}

func TestTranslationsCoverAllLangs(t *testing.T) {
	for key, m := range messages {
		for _, lang := range Langs {
			if v, ok := m[lang]; !ok || v == "" {
				t.Errorf("translation key %q missing %s", key, lang)
			}
		}
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Привет мир":     "privet-mir",
		"Қазақстан 2026": "qazaqstan-2026",
		"Hello, World!":  "hello-world",
		"   ":            "article",
		"Экономика Казахстана": "ekonomika-kazahstana",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}
