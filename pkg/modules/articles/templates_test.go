package articles

import (
	"html/template"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"
)

// buildTemplates mirrors Module.Init template wiring so we can validate the
// embedded templates without a running server or database.
func buildTemplates(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.New("articles").Funcs(templateFuncs()).ParseFS(templateFiles, "templates/*.html")
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
		base := Base{Title: "T", Lang: lang, Authed: true, ShowLangs: true, ActiveCat: "sport", ActiveSub: "football", LangLinks: langLinks("/", "cat=sport"), Ads: demoAds(lang)}
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
			{"form", FormPage{Base: base, Mode: "register", Email: "a@b.c", Last: "Баймурза", First: "Даулет", Middle: "Абаевич", Ref: "abc23", Error: "err"}},
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
				Documents: []string{"/media/plan.pdf"}, ContractURL: "/media/lease.pdf",
				Images: []string{"/static/demo/rooms/exterior.svg", "/static/demo/rooms/living.svg", "/static/demo/rooms/kitchen.svg"}}}},
			{"listing_view", ListingViewPage{Base: base, Owner: true, L: &Listing{ // cover-only fallback, owner controls
				ID: "id", DealType: "sale", PropertyType: "land", Title: "Участок", Contact: "+7 700 000 00 00", CoverURL: "http://x/y.jpg",
				ExpiresAt: now.Add(72 * time.Hour), PromotedUntil: &now, FeaturedUntil: &now}}},
			{"listing_my", MyListingsPage{Base: base, Listings: []*Listing{
				{ID: "id1", Title: "Активное", Price: 12000000, Region: "Алматы", ExpiresAt: now.Add(72 * time.Hour), FeaturedUntil: &now},
				{ID: "id2", Title: "Истёкшее", Price: 5000000, Region: "Астана", ExpiresAt: now.Add(-24 * time.Hour)},
			}}},
			{"listing_my", MyListingsPage{Base: base}}, // empty state
			{"maintenance", MaintenancePage{Lang: lang, Title: "Maintenance", Message: "Back soon"}},
			{"maintenance", MaintenancePage{Lang: lang, Title: "Maintenance", Message: "Back soon", UntilMilli: now.Add(time.Hour).UnixMilli()}}, // with countdown
			{"agent_cabinet", AgentCabinetPage{Base: base, Draft: Agent{FirstName: "Асан", LastName: "Серіков"}}},                                // registration form
			{"agent_cabinet", AgentCabinetPage{Base: base, Agent: &Agent{UserID: "u1", FirstName: "Асан", LastName: "Серіков", Name: "Асан Серіков", Agency: "Дом", Phone: "+7", Status: agentVerified}, Draft: Agent{FirstName: "Асан", LastName: "Серіков"}, Count: 3, Saved: true}},
			{"agent_cabinet", AgentCabinetPage{Base: base, Agent: &Agent{UserID: "u1", Name: "Асан", Status: agentPending}, Draft: Agent{FirstName: "Асан"}}},                        // pending
			{"agent_cabinet", AgentCabinetPage{Base: base, Agent: &Agent{UserID: "u1", Name: "Асан", Status: agentRejected, Reason: "нет данных"}, Draft: Agent{FirstName: "Асан"}}}, // rejected
			{"agent_public", AgentPublicPage{Base: base, Agent: &Agent{UserID: "u1", Name: "Асан Серіков", Agency: "Дом", Phone: "+7", About: "Опыт 10 лет", Status: agentVerified}, Listings: []*Listing{{
				ID: "id", DealType: "sale", PropertyType: "apartment", Title: "Квартира", Price: 18000000, AgentID: "u1", AgentName: "Асан Серіков",
				Images: []string{"/static/demo/rooms/living.svg"}}}}},
			{"agent_public", AgentPublicPage{Base: base, Agent: &Agent{UserID: "u1", Name: "Асан", Status: agentVerified}}}, // no listings
			{"admin", AdminPage{Base: base, Email: "a@b.c", Role: "admin", CanManageUsers: true, CanModerate: true, CanFinance: true,
				AssignRoles: assignableRoles, ServiceStates: []string{svcOn, svcMaintenance, svcOff},
				Services:      []ServiceFlag{{Code: SvcAdOrders, TitleKey: "svc.ad_orders", Status: svcMaintenance}},
				Site:          ServiceFlag{Code: SvcSite, TitleKey: "svc.site", Status: svcOn},
				PendingAgents: []Agent{{UserID: "u1", Name: "Асан Серіков", Agency: "Дом", Phone: "+7 700", Email: "a@b.c", Status: agentPending}},
				Payments:      paymentsAdminView{Enabled: true, Provider: PayProviderKaspi, ActiveReady: false, Providers: []paymentProviderStatus{{Code: PayProviderKaspi, Label: "Kaspi Pay", Implemented: false, IsActive: true}, {Code: PayProviderIoka, Label: "ioka", Implemented: false}}},
				Stats:         AdminStats{Users: 3, Articles: 2}}},
			{"admin_pages", adminPagesList{Base: base, Items: []adminPageItem{{Key: "privacy", Name: "Конфиденциальность"}, {Key: "terms", Name: "Условия"}}}},
			{"admin_page_edit", adminPageEditView{Base: base, Key: "privacy", Name: "Конфиденциальность", Notice: "N", LastEdited: "2026-07-28 10:00", LastEditor: "a@b.c", Langs: []adminPageLangView{
				{Code: "kz", Label: "Қазақша", Title: "T", Body: "# Hi"},
				{Code: "ru", Label: "Русский", Title: "T", Body: "# Hi"},
				{Code: "en", Label: "English", Title: "T", Body: "# Hi"},
			}}},
			{"admin_tariffs", tariffsAdminView{Base: base, Notice: "N", AdDays: []int{3, 7, 14, 30},
				AdFormats: []tariffAdRow{{Label: "970×250", Cells: []tariffCell{{"ad.horizontal.3", 18000}, {"ad.horizontal.7", 36000}, {"ad.horizontal.14", 60000}, {"ad.horizontal.30", 90000}}}},
				Weights:   []tariffField{{"weight.high", "High", 20}},
				Banner:    []tariffField{{"banner.1", "1", 990}},
				Services:  []tariffField{{"promote.price", "Top", 299}},
				Durations: []tariffField{{"listing.free_days", "Free", 21}}}},
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

func TestShortAuthor(t *testing.T) {
	cases := map[string]string{
		"Daulet Baimurza":       "D. Baimurza",
		"Даулет Баймурза":       "Д. Баймурза",
		"Иван Петрович Сидоров": "И. П. Сидоров",
		"AI Dake": "A. Dake",
		"Сана":    "Сана", // one word — nothing to abbreviate
		"  Айгерим  Нурланова   ": "А. Нурланова",
		"": "",
	}
	for in, want := range cases {
		if got := shortAuthor(in); got != want {
			t.Errorf("shortAuthor(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCompactNum(t *testing.T) {
	cases := []struct {
		lang string
		n    int64
		want string
	}{
		{LangRU, 0, "0"},
		{LangRU, 115, "115"},
		{LangRU, 999, "999"},
		{LangRU, 1000, "1 тыс."},
		{LangRU, 1999, "1,9 тыс."}, // truncated, never rounded up
		{LangRU, 12345, "12 тыс."},
		{LangRU, 999999, "999 тыс."},
		{LangRU, 2_500_000, "2,5 млн"},
		{LangKZ, 1500, "1,5 мың"},
		{LangEN, 1500, "1.5k"},
		{LangEN, 3_000_000, "3M"},
	}
	for _, c := range cases {
		if got := compactNum(c.lang, c.n); got != c.want {
			t.Errorf("compactNum(%q, %d) = %q, want %q", c.lang, c.n, got, c.want)
		}
	}
}

// The admin panel's lists grow forever as data accumulates. Left unbounded they
// stretch their card until one column runs several screens past the other, so
// every growing list must scroll inside a fixed height instead.
func TestAdminGrowingListsScroll(t *testing.T) {
	tmpl := buildTemplates(t)
	now := time.Now()
	rows := []GuestSimpleRow{{Title: "Прямые", N: 1345, Pct: 100}, {Title: "Facebook", N: 62, Pct: 5}}
	var cs []AdminComment
	var ml []ModAction
	for i := 0; i < 40; i++ {
		cs = append(cs, AdminComment{ID: "c", AuthorName: "A", Body: "текст", Slug: "s", CreatedAt: now})
		ml = append(ml, ModAction{Created: now, Action: "hide", TargetType: "comment", ReasonCode: "spam", ActorKind: "human"})
	}
	page := AdminPage{Base: Base{Lang: LangRU, Title: "T"}, Email: "a@b.c", Role: "admin",
		CanManageUsers: true, CanModerate: true, CanFinance: true,
		AssignRoles: assignableRoles, ServiceStates: []string{svcOn},
		Stats:  AdminStats{Users: 3, Articles: 2, RecentComments: cs},
		ModLog: ml,
		Guests: GuestAnalytics{HasData: true, Sources: rows, Bots: rows, Devices: rows, OS: rows,
			Browsers: rows, Countries: rows, Langs: rows, EnglishBy: rows, VPNLangs: rows},
	}
	var sb strings.Builder
	if err := tmpl.ExecuteTemplate(&sb, "admin", page); err != nil {
		t.Fatal(err)
	}
	out := sb.String()

	if !strings.Contains(out, `class="adm-comments adm-scroll"`) {
		t.Error("recent comments must scroll")
	}
	if !strings.Contains(out, `class="table-wrap adm-scroll adm-scroll--tall"`) {
		t.Error("the moderation log must scroll")
	}
	if n := strings.Count(out, "adm-scroll"); n < 5 {
		t.Errorf("scrollable blocks: %d, want at least 5 (comments, log, three queues)", n)
	}

	// All nine guest-analytics panels must sit in ONE grid, otherwise they cannot
	// line up as 3×3 and the layout falls back to stacked pairs.
	g := regexp.MustCompile(`(?s)<div class="adm-cols adm-cols--3">(.*?)\n      </div>`).FindStringSubmatch(out)
	if g == nil {
		t.Fatal("three-column analytics grid not rendered")
	}
	if n := strings.Count(g[1], `<div class="adm-panel">`); n != 9 {
		t.Errorf("panels inside the grid: %d, want 9", n)
	}
	if strings.Contains(g[1], `<div class="adm-panel"></div>`) {
		t.Error("the empty filler panel should be gone — nine cells fill 3×3 exactly")
	}
}

func TestStudioDeleteIsDraftOnly(t *testing.T) {
	tmpl := buildTemplates(t)
	now := time.Now()
	var sb strings.Builder
	page := StudioPage{Base: Base{Lang: "ru", Title: "T"}, Articles: []StudioRow{
		{ID: "d1", Slug: "s1", Title: "Черновик", Status: "draft", Updated: now, Langs: []string{LangRU}},
		{ID: "p1", Slug: "s2", Title: "Опубликована", Status: "published", Updated: now, Langs: []string{LangRU}},
	}}
	if err := tmpl.ExecuteTemplate(&sb, "studio_dashboard", page); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	if !strings.Contains(out, `action="/studio/a/d1/delete"`) {
		t.Error("a draft must offer delete")
	}
	// A published article carries views, votes, comments and a shared link;
	// deleting it must cost an extra deliberate step (unpublish first).
	if strings.Contains(out, `action="/studio/a/p1/delete"`) {
		t.Error("a published article must not offer delete")
	}
	if !strings.Contains(out, "onsubmit=\"return confirm(") {
		t.Error("delete must sit behind a confirmation")
	}
}

func TestAgentKind(t *testing.T) {
	// An unknown or empty kind must fall back to the least-privileged shape, so
	// a bad value can never grant a company badge.
	for _, in := range []string{"", "company", "COMPANY", "администратор"} {
		if got := NormalizeAgentKind(in); got != "private" {
			t.Errorf("NormalizeAgentKind(%q) = %q, want private", in, got)
		}
	}
	for _, k := range agentKinds {
		if NormalizeAgentKind(k) != k {
			t.Errorf("NormalizeAgentKind(%q) changed a valid kind", k)
		}
		if key := (Agent{Kind: k}).KindLabelKey(); messages[key] == nil {
			t.Errorf("kind %q has no translation key %q", k, key)
		}
	}
	if !(Agent{Kind: "agency"}).IsCompany() || !(Agent{Kind: "developer"}).IsCompany() {
		t.Error("agency and developer must count as companies")
	}
	if (Agent{Kind: "private"}).IsCompany() {
		t.Error("a private realtor is not a company")
	}
	// A profile written before the kind column existed still renders a label.
	if key := (Agent{}).KindLabelKey(); messages[key] == nil {
		t.Errorf("empty kind produced unusable key %q", key)
	}
}

func TestValidBIN(t *testing.T) {
	if !validBIN("123456789012") {
		t.Error("12 digits must be accepted")
	}
	for _, bad := range []string{"", "12345678901", "1234567890123", "12345678901a", "1234 5678 9012"} {
		if validBIN(bad) {
			t.Errorf("validBIN(%q) = true, want false", bad)
		}
	}
}

func TestLiveSocial(t *testing.T) {
	in := []SocialLink{
		{Name: "telegram", URL: "https://t.me/shanraq_org"},
		{Name: "youtube", URL: "#"}, // not opened yet — must not become a dead click
		{Name: "facebook", URL: ""},
		{Name: "instagram", URL: "https://instagram.com/shanraq_org"},
	}
	got := liveSocial(in)
	if len(got) != 2 || got[0].Name != "telegram" || got[1].Name != "instagram" {
		t.Errorf("liveSocial() = %+v, want telegram + instagram only", got)
	}
	if n := len(liveSocial(nil)); n != 0 {
		t.Errorf("liveSocial(nil) returned %d links, want 0", n)
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

// TestShareRowLinks pins the sharing block, which is easy to break silently:
// a share link that is double-escaped still renders fine and only fails in the
// target app, where nobody on the team would notice.
func TestShareRowLinks(t *testing.T) {
	tmpl := buildTemplates(t)

	const url = "https://shanraq.org/read/ne-ta-tablica?lang=ru"
	var b strings.Builder
	if err := tmpl.ExecuteTemplate(&b, "share_row", map[string]any{
		"Lang": LangRU, "Title": "Не та таблица & спорт", "URL": url,
	}); err != nil {
		t.Fatalf("execute share_row: %v", err)
	}
	out := b.String()

	// The article URL must reach every target exactly once-escaped. "%2520" is
	// the signature of double escaping (a space encoded, then the % encoded).
	if strings.Contains(out, "%2520") || strings.Contains(out, "%253A") || strings.Contains(out, "%253a") {
		t.Error("share URLs are double-escaped")
	}
	const escaped = "https%3a%2f%2fshanraq.org%2fread%2fne-ta-tablica%3flang%3dru"
	for _, host := range []string{"wa.me", "t.me/share/url", "facebook.com/sharer", "linkedin.com/sharing"} {
		if !strings.Contains(out, host) {
			t.Errorf("share row is missing %s", host)
		}
	}
	if n := strings.Count(out, escaped); n != 4 {
		t.Errorf("escaped article URL appears %d times, want 4 (one per network)", n)
	}
	// The ampersand in the title must be encoded, never emitted raw into the
	// query, or Telegram would drop everything after it.
	if !strings.Contains(out, "%26") {
		t.Error("title ampersand was not percent-encoded")
	}
	// No Instagram button: Instagram has no web share endpoint, so one could
	// only pretend to work. See the comment on the partial.
	if strings.Contains(strings.ToLower(out), "instagram") {
		t.Error("share row must not offer Instagram — there is no web share endpoint for it")
	}
	// The OS share sheet button ships hidden and is revealed by JS only where
	// navigator.share exists.
	if !strings.Contains(out, "data-share-native") || !strings.Contains(out, "hidden") {
		t.Error("native share button must render hidden by default")
	}
}

// TestArticleShowsShareToGuests guards the point of the feature: the reader most
// likely to pass an article on is the one who is not logged in.
func TestArticleShowsShareToGuests(t *testing.T) {
	tmpl := buildTemplates(t)
	now := time.Now()
	page := ArticlePage{
		Base: Base{Title: "T", Lang: LangRU, Authed: false,
			SiteURL: "https://shanraq.org", CanonURL: "/read/s?lang=ru"},
		Slug: "s", Title: "T", AuthorName: "A", Category: "sport",
		Published: &now, Body: template.HTML("<p>x</p>"),
	}
	var b strings.Builder
	if err := tmpl.ExecuteTemplate(&b, "article", page); err != nil {
		t.Fatalf("execute article: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, `data-share`) {
		t.Fatal("guest article page has no share row")
	}
	if !strings.Contains(out, "shanraq.org%2fread%2fs%3flang%3dru") {
		t.Error("share row did not receive the absolute canonical article URL")
	}
	if strings.Contains(out, "favbtn") {
		t.Error("favourites button must stay hidden from guests")
	}
}
