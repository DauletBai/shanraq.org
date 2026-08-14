package articles

import (
	"html/template"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"

	"shanraq.org/pkg/modules/ai"
	"shanraq.org/web"
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
		base := Base{Title: "T", Lang: lang, Authed: true, ShowLangs: true, ActiveCat: "sport", ActiveSub: "football", LangLinks: langLinks("/", "cat=sport"), Ads: houseAds(lang)}
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
				ID: "id", DealType: "rent", PropertyType: "house", Country: "Казахстан", Region: "Астана", City: "Астана", District: "Тельман",
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
	//
	// Anchored to the mix grid by its own class, not by position. Matching the
	// first three-column grid in the section was a locator that broke the moment
	// a neighbour became three columns too — that has now happened twice.
	guests := out[strings.Index(out, `<section id="guests"`):]
	g := regexp.MustCompile(`(?s)<div class="adm-cols adm-cols--3 adm-cols--mix">(.*?)\n      </div>`).FindStringSubmatch(guests)
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
	// Every channel must tag its own link. Messengers strip the Referer, so an
	// untagged share arrives as "direct" and the panel cannot tell which channel
	// actually works — the one question worth asking at this stage.
	for _, src := range []string{"whatsapp", "telegram", "facebook", "linkedin"} {
		if !strings.Contains(out, "utm_source%3d"+src) {
			t.Errorf("share link for %s carries no utm_source", src)
		}
	}
	// data-url is an ordinary attribute, not a URL query, so it is HTML-escaped
	// rather than percent-encoded. It feeds the copy button and the OS share
	// sheet, where the destination is unknown — hence the honest "share" label.
	if !strings.Contains(out, "utm_source=share") {
		t.Error("copy/native share URL carries no utm_source")
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
	// Two rows: one in the byline, one under the text for the reader who just
	// finished. The second is the one that actually gets used.
	// Count the row opener, not "data-share": that prefix also matches the
	// per-button hooks and the footer script's selectors.
	if n := strings.Count(out, "data-share data-url="); n != 2 {
		t.Errorf("article has %d share rows, want 2 (byline + end of text)", n)
	}
	if !strings.Contains(out, "share--foot") {
		t.Error("end-of-article share row is missing its foot variant")
	}
	if !strings.Contains(out, T(LangRU, "share.foot")) {
		t.Error("end-of-article row should carry the pass-it-on label, not the plain one")
	}
}

// TestListingShowsShare guards sharing on listings, where it matters most:
// people send a flat to family long before they call the number.
func TestListingShowsShare(t *testing.T) {
	tmpl := buildTemplates(t)
	page := ListingViewPage{
		Base: Base{Title: "T", Lang: LangRU, Authed: false,
			SiteURL: "https://shanraq.org", CanonURL: "/listings/abc?lang=ru"},
		L: &Listing{ID: "abc", DealType: "rent", PropertyType: "house",
			Country: "Казахстан", City: "Астана", Price: 350000, Area: 120, Rooms: 4,
			Title: "Дом в аренду", Contact: "+7 700 000 00 00"},
	}
	var b strings.Builder
	if err := tmpl.ExecuteTemplate(&b, "listing_view", page); err != nil {
		t.Fatalf("execute listing_view: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, "data-share") {
		t.Fatal("listing page has no share row")
	}
	if !strings.Contains(out, "shanraq.org%2flistings%2fabc%3flang%3dru") {
		t.Error("listing share row did not receive the absolute canonical URL")
	}
}

// TestWithUTM covers the join, which is the only place this can go wrong: an
// article URL always carries ?lang=, so a "?" would produce a second query
// string and the tag would be silently dropped by the browser.
func TestWithUTM(t *testing.T) {
	cases := []struct{ in, src, want string }{
		{"https://shanraq.org/read/x?lang=ru", "whatsapp",
			"https://shanraq.org/read/x?lang=ru&utm_source=whatsapp&utm_medium=share"},
		{"https://shanraq.org/read/x", "telegram",
			"https://shanraq.org/read/x?utm_source=telegram&utm_medium=share"},
		{"https://shanraq.org/read/x?lang=ru", "", "https://shanraq.org/read/x?lang=ru"},
		{"", "whatsapp", ""},
	}
	for _, c := range cases {
		if got := withUTM(c.in, c.src); got != c.want {
			t.Errorf("withUTM(%q, %q) = %q, want %q", c.in, c.src, got, c.want)
		}
	}
	// The label has to survive the round trip, or the tag is decoration.
	for _, src := range []string{"whatsapp", "telegram", "facebook", "linkedin", "share"} {
		if got := utmSource(src); got != src {
			t.Errorf("utmSource(%q) = %q — share links would not be attributed", src, got)
		}
	}
}

// Status badges must take their colour from a theme token, never from an inline
// style. An inline colour is the one place a dark-theme rule cannot reach, and
// the four that had one rendered at 1.8–2.9:1 against the dark panel — "По
// приглашению" was effectively invisible. This catches the next one.
func TestAdminStatusBadgesAreThemeable(t *testing.T) {
	tmpl := buildTemplates(t)
	var sb strings.Builder
	page := AdminPage{
		Base:           Base{Lang: "ru", Title: "T"},
		CanManageUsers: true,
		Services:       []ServiceFlag{{Code: "listing_promo", Status: svcInviteOnly}},
		Site:           ServiceFlag{Code: "site", Status: svcOn},
	}
	if err := tmpl.ExecuteTemplate(&sb, "admin", page); err != nil {
		t.Fatalf("render admin: %v", err)
	}
	out := sb.String()
	if m := regexp.MustCompile(`style="[^"]*color:\s*#[0-9a-fA-F]{3,6}`).FindString(out); m != "" {
		t.Errorf("a status colour is written inline, where no theme can override it: %s", m)
	}
}

// The pages / clicks / trend row: three panels of one subject each, the pages
// figures said once rather than as a chart plus a table of the same numbers,
// and a chart with a readable scale.
func TestAdminGuestPanelsSplitInThree(t *testing.T) {
	tmpl := buildTemplates(t)
	rows := []GuestSimpleRow{{Title: "Статьи", N: 1223, Pct: 100}, {Title: "Главная", N: 1044, Pct: 85}}
	pages := []GuestPageRow{
		{Title: "Статьи", Pct: 100, A: Audience{Guest: 1184, Registered: 39}},
		{Title: "Главная", Pct: 85, A: Audience{Guest: 984, Registered: 60}},
	}
	clicks := []GuestClickRow{{Title: "Кнопка входа", A: Audience{Guest: 37}}}
	page := AdminPage{Base: Base{Lang: LangRU, Title: "T"}, Email: "a@b.c", Role: "admin",
		AssignRoles: assignableRoles, ServiceStates: []string{svcOn},
		Stats: AdminStats{Users: 1},
		Guests: GuestAnalytics{HasData: true, Pages: pages, Clicks: clicks,
			Sources: rows, Bots: rows, Devices: rows, OS: rows, Browsers: rows,
			Countries: rows, Langs: rows, EnglishBy: rows, VPNLangs: rows,
			Trend:      []GuestTrendDay{{Label: "29.07", N: 114, Pct: 38}, {Label: "11.08", N: 291, Pct: 97, Tick: true}},
			TrendTicks: []AxisTick{{N: 300, Pct: 100}, {N: 150, Pct: 50}, {N: 0, Pct: 0}}},
	}
	var sb strings.Builder
	if err := tmpl.ExecuteTemplate(&sb, "admin", page); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	guests := out[strings.Index(out, `<section id="guests"`):]
	g := regexp.MustCompile(`(?s)<div class="adm-cols adm-cols--3">(.*?)\n      </div>\n      <div class="adm-cols`).FindStringSubmatch(guests)
	if g == nil {
		t.Fatal("the pages/clicks/trend grid is not rendered")
	}
	if n := strings.Count(g[1], `<div class="adm-panel`); n != 3 {
		t.Errorf("panels in the row: %d, want 3", n)
	}
	// The duplicated bar chart above the table is gone: one table carries both.
	if strings.Contains(g[1], `<div class="adm-chart">`) {
		t.Error("the pages panel still repeats its numbers as a separate chart")
	}
	if !strings.Contains(g[1], `class="spec adm-mixed"`) {
		t.Error("the pages table should carry its bars inline")
	}
	// The long header set its own column's width and starved the neighbour, so
	// the numeric columns carry a short label with the full word as a tooltip.
	if strings.Contains(g[1], `<th>зарегистрированные</th>`) {
		t.Error("the numeric column still uses the unbreakable full-length header")
	}
	if !strings.Contains(g[1], `<th title="зарегистрированные">зарег.</th>`) {
		t.Error("the short header should keep the full word as its title")
	}
	// Axes, or a bar can only be compared, never read.
	for _, want := range []string{`class="spark__y"`, `class="spark__x"`, `class="spark__grid"`} {
		if !strings.Contains(g[1], want) {
			t.Errorf("the trend chart is missing %s", want)
		}
	}
	// Four numeric columns in a third of a row must scroll, not spill.
	if !strings.Contains(g[1], `class="adm-xscroll"`) {
		t.Error("the pages table is not in a horizontal scroller")
	}
	// The chart fills the height its tallest neighbour sets instead of leaving
	// the bottom of the card empty.
	if !strings.Contains(g[1], `class="adm-panel adm-panel--fill"`) {
		t.Error("the trend panel should stretch to the height of the row")
	}
	// Gridline and label share one offset, so a label always sits on its line.
	if !strings.Contains(g[1], `<i style="bottom:100%"></i>`) ||
		!strings.Contains(g[1], `<span style="bottom:100%">300</span>`) {
		t.Error("ticks should place the label and its gridline at the same offset")
	}
}

// The sidebar highlight follows scroll position, so a nav whose order disagrees
// with the document sends the marker jumping up and back down as you scroll.
// They drifted apart once already — from the fifth item on — and nothing caught
// it.
// renderAdminFull renders the admin page for a root user, so every
// permission-gated section is present.
func renderAdminFull(t *testing.T) string {
	t.Helper()
	tmpl := buildTemplates(t)
	rows := []GuestSimpleRow{{Title: "Прямые", N: 10, Pct: 100}}
	page := AdminPage{Base: Base{Lang: LangRU, Title: "T"}, Email: "a@b.c", Role: "admin",
		CanManageUsers: true, CanModerate: true, CanFinance: true,
		AssignRoles: assignableRoles, ServiceStates: []string{svcOn},
		Site:  ServiceFlag{Code: SvcSite, TitleKey: "svc.site", Status: svcOn},
		Stats: AdminStats{Users: 1},
		// The settings cards render nothing without these, and the tests below
		// depend on the provider rows actually being there.
		AI: ai.AdminView{Enabled: true, Provider: "anthropic",
			EditorModel: "claude-sonnet-5", TranslateModel: "claude-haiku-4-5", MaxTokens: 4096,
			Providers: []ai.ProviderStatus{
				{Code: "anthropic", Label: "Claude (Anthropic)", IsActive: true,
					Editor: "claude-sonnet-5", Translate: "claude-haiku-4-5"},
				{Code: "openai", Label: "ChatGPT (OpenAI)"},
			}},
		Payments: paymentsAdminView{Enabled: true, Provider: PayProviderKaspi,
			Providers: []paymentProviderStatus{{Code: PayProviderKaspi, Label: "Kaspi Pay", IsActive: true}}},
		Guests: GuestAnalytics{HasData: true, Sources: rows, Bots: rows, Devices: rows,
			OS: rows, Browsers: rows, Countries: rows, Langs: rows, EnglishBy: rows, VPNLangs: rows},
	}
	var sb strings.Builder
	if err := tmpl.ExecuteTemplate(&sb, "admin", page); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestAdminNavMatchesSectionOrder(t *testing.T) {
	out := renderAdminFull(t)

	navRe := regexp.MustCompile(`href="#([a-z]+)"[^>]*data-nav`)
	var nav []string
	for _, m := range navRe.FindAllStringSubmatch(out, -1) {
		nav = append(nav, m[1])
	}
	secRe := regexp.MustCompile(`<(?:section|div)[^>]*id="([a-z]+)"[^>]*class="adm-section"|<div class="adm-cols adm-cols--settings" id="([a-z]+)"`)
	var doc []string
	for _, m := range secRe.FindAllStringSubmatch(out, -1) {
		id := m[1]
		if id == "" {
			id = m[2]
		}
		doc = append(doc, id)
	}
	// Only the anchors the nav points at matter; ai and payments sit side by side
	// inside #settings and share one entry precisely because they share a height.
	want := map[string]bool{}
	for _, n := range nav {
		want[n] = true
	}
	var filtered []string
	for _, d := range doc {
		if want[d] {
			filtered = append(filtered, d)
		}
	}

	if len(nav) == 0 || len(filtered) == 0 {
		t.Fatalf("found %d nav entries and %d sections — the locator is broken", len(nav), len(filtered))
	}
	if len(nav) != len(filtered) {
		t.Fatalf("nav has %v but the page has %v — every entry must have a section and vice versa", nav, filtered)
	}
	for i := range nav {
		if nav[i] != filtered[i] {
			t.Fatalf("order diverges at position %d: nav %v, page %v", i+1, nav, filtered)
		}
	}
}

// A nav entry pointing at nothing can never light up — that is how "Обзор" spent
// months looking broken.
func TestAdminNavAnchorsExist(t *testing.T) {
	out := renderAdminFull(t)
	for _, m := range regexp.MustCompile(`href="#([a-z]+)"[^>]*data-nav`).FindAllStringSubmatch(out, -1) {
		if !strings.Contains(out, `id="`+m[1]+`"`) {
			t.Errorf("nav points at #%s but no element carries that id", m[1])
		}
	}
}

// The home page had no <h1> at all: its visible title is the wordmark, which is
// an image. Category feeds are separate indexable URLs, so each names its own
// subject rather than all ten repeating one line.
func TestHomeHasExactlyOneH1(t *testing.T) {
	tmpl := buildTemplates(t)
	render := func(cat string) string {
		t.Helper()
		base := Base{Title: "T", Lang: LangRU, ShowLangs: true, ActiveCat: cat, LangLinks: langLinks("/", "")}
		var sb strings.Builder
		if err := tmpl.ExecuteTemplate(&sb, "home", HomePage{Base: base}); err != nil {
			t.Fatal(err)
		}
		return sb.String()
	}

	plain := render("")
	if n := strings.Count(plain, "<h1"); n != 1 {
		t.Fatalf("home has %d <h1> elements, want exactly 1", n)
	}
	if !strings.Contains(plain, `<h1 class="visually-hidden">`) {
		t.Error("the heading should be present to crawlers and readers, not drawn")
	}
	if !strings.Contains(plain, "Shanraq.org") {
		t.Error("the heading does not name the publication")
	}

	// A category feed must not repeat the site-wide heading.
	world := render("world")
	if strings.Count(world, "<h1") != 1 {
		t.Fatal("a category feed should still have exactly one <h1>")
	}
	if !strings.Contains(world, T(LangRU, "cat.world")) {
		t.Errorf("the category feed's heading does not name the category: %s",
			regexp.MustCompile(`(?s)<h1.*?</h1>`).FindString(world))
	}
	if strings.Contains(world, T(LangRU, "home.h1")) {
		t.Error("the category feed repeats the site-wide heading")
	}
}

// Leaflet and its stylesheet are about 270 KB and were loaded on every page of
// the site, including the home feed, which has no map. Only the two pages that
// draw one should pay for it.
func TestLeafletOnlyWhereThereIsAMap(t *testing.T) {
	tmpl := buildTemplates(t)
	base := Base{Title: "T", Lang: LangRU, ShowLangs: true, LangLinks: langLinks("/", "")}
	render := func(name string, data any) string {
		t.Helper()
		var sb strings.Builder
		if err := tmpl.ExecuteTemplate(&sb, name, data); err != nil {
			t.Fatal(err)
		}
		return sb.String()
	}

	for _, c := range []struct {
		name string
		data any
	}{
		{"home", HomePage{Base: base}},
		{"article", ArticlePage{Base: base, Slug: "s", Title: "T", ServedLang: LangRU}},
		{"page", StaticPage{Base: base}},
	} {
		out := render(c.name, c.data)
		// What matters is the fetch, not the word: the shared map-init script
		// mentions a Leaflet class name and already no-ops without the library.
		if strings.Contains(out, "leaflet.js") || strings.Contains(out, "leaflet.css") {
			t.Errorf("%s loads Leaflet but draws no map", c.name)
		}
	}

	mapped := base
	mapped.NeedsMap = true
	for _, c := range []struct {
		name string
		data any
	}{
		{"listings", ListingsPage{Base: mapped}},
		{"listing_new", ListingFormPage{Base: mapped, Values: ListingInput{DealType: "sale", PropertyType: "apartment"}}},
	} {
		out := render(c.name, c.data)
		if !strings.Contains(out, "leaflet.js") || !strings.Contains(out, "leaflet.css") {
			t.Errorf("%s draws a map but is missing Leaflet", c.name)
		}
	}
}

// The cover's box has to be reserved before the image arrives, or the heading
// below it jumps down on load — that shift was the article page's whole CLS.
func TestArticleCoverReservesItsSpace(t *testing.T) {
	css, err := web.StaticFS().Open("css/shanraq.css")
	if err != nil {
		t.Fatalf("open stylesheet: %v", err)
	}
	defer css.Close()
	b, err := io.ReadAll(css)
	if err != nil {
		t.Fatal(err)
	}
	rule := regexp.MustCompile(`\.article__media \{[^}]*\}`).Find(b)
	if rule == nil {
		t.Fatal(".article__media rule not found")
	}
	if !strings.Contains(string(rule), "aspect-ratio") {
		t.Errorf("the cover box has no reserved ratio: %s", rule)
	}

	// Feed thumbnails carried a fixed pixel height, so the box changed shape with
	// the viewport while the picture did not — a phone cut 18% off the top and
	// bottom, a desktop card cut the sides.
	// The selector appears more than once, so check every rule, not the first.
	rules := regexp.MustCompile(`\.post__media \{[^}]*\}`).FindAll(b, -1)
	sized := false
	for _, r := range rules {
		if strings.Contains(string(r), "aspect-ratio") {
			sized = true
		}
	}
	if !sized {
		t.Errorf("no .post__media rule sizes the thumbnail by ratio (found %d rules)", len(rules))
	}
	if regexp.MustCompile(`\.post__media \{[^}]*height: \d+px`).Match(b) {
		t.Error("a fixed pixel height is back on the feed thumbnail")
	}

	// The home carousel had the same fault and it was the most visible one: a
	// fixed 300px box against a 728px column is 2.43:1, so 27% of every 16:9
	// cover was cut off the top and bottom, and on a phone the 210px box was
	// 1.50:1 and cut the sides instead.
	hero := regexp.MustCompile(`\.adcarousel \{[^}]*\}`).Find(b)
	if hero == nil {
		t.Fatal(".adcarousel rule not found")
	}
	if !strings.Contains(string(hero), "aspect-ratio: 16 / 9") {
		t.Errorf("the carousel is not 16:9: %s", hero)
	}
	if regexp.MustCompile(`\.adcarousel \{[^}]*height: \d+px`).Match(b) {
		t.Error("a fixed pixel height is back on the carousel; it overrides the ratio")
	}
	// A 16:9 box on a 360px phone is 178px tall. The headline has to be sized
	// against the box, or it spills straight out of the bottom of it.
	text := regexp.MustCompile(`\.newshero \.adslide__text \{[^}]*\}`).Find(b)
	if text == nil {
		t.Fatal(".newshero .adslide__text rule not found")
	}
	if !strings.Contains(string(text), "cqw") {
		t.Errorf("the hero headline is not sized against its own box: %s", text)
	}
	if !regexp.MustCompile(`\.newshero \{[^}]*container-type`).Match(b) {
		t.Error("the hero declares no container, so its cqw sizes resolve against nothing")
	}

	// Blocks stacked inside a section are separated by one flow rule rather than
	// by a list of pairs. The list was the bug: it had to be extended every time
	// a new block was put next to an old one, and it was not — which is how
	// "Последние комментарии" came to sit flush against the cards below it.
	if !regexp.MustCompile(`\.adm-section > \* \+ \* \{[^}]*margin-top`).Match(b) {
		t.Error("nothing separates stacked blocks inside an admin section")
	}
	// One set of numbers decides admin spacing. A literal px in a gap or a card
	// padding is how the last set drifted apart into 16/18/20/22.
	for _, tok := range []string{"--adm-gap:", "--adm-pad:", "--adm-sec:", "--adm-shadow:"} {
		if !strings.Contains(string(b), tok) {
			t.Errorf("the admin spacing token %s is gone", tok)
		}
	}
	// The card rule is the more specific of the two, so a `margin` shorthand here
	// silently beats the section's flow rule and the cards go back to touching —
	// which is exactly what happened. Only the bottom margin is the card's to
	// reset; the top belongs to the flow.
	if card := regexp.MustCompile(`\.adm \.adm-panel,[^{]*\{[^}]*\}`).Find(b); card == nil {
		t.Error("the shared admin card rule is gone")
	} else if regexp.MustCompile(`[;{]\s*margin\s*:`).Match(card) {
		t.Errorf("the card rule uses the margin shorthand; it will cancel the section flow: %s", card)
	}

	// Newspaper flow is what put cards in one row at different heights.
	if regexp.MustCompile(`\.adm-cards2 \{[^}]*(?:^|[^-\w])columns\s*:`).Match(b) {
		t.Error(".adm-cards2 is laid out with columns again; its cards will not line up")
	}
	// Sideways, the ratio alone would make the cover taller than the screen.
	if !strings.Contains(string(b), ".article__media { max-height: 70vh; }") {
		t.Error("the cover has no height cap, so in landscape it fills the viewport")
	}
}

// The provider rows carried an inline min-width that could not fit a phone, and
// the model fields kept the previous provider's names after a switch.
func TestAISettingsCardFitsAndFollowsProvider(t *testing.T) {
	out := renderAdminFull(t)
	settings := out[strings.Index(out, `id="settings"`):]

	// Nothing in a provider row may be pinned to a width the card cannot give it.
	if strings.Contains(settings, "min-width:160px") {
		t.Error("a provider name is still pinned to 160px inline and will overflow a phone")
	}
	if n := strings.Count(settings, `class="adm-choice"`); n < 2 {
		t.Errorf("provider rows are not using the wrapping class (found %d)", n)
	}

	// Each radio has to carry what its provider should be driven with, or the
	// form cannot correct the fields when the choice changes.
	if !strings.Contains(settings, `data-editor-model="claude-sonnet-5"`) ||
		!strings.Contains(settings, `data-translate-model="claude-haiku-4-5"`) {
		t.Error("the Anthropic radio does not carry its models")
	}
	// Providers whose ids we do not know carry empty ones on purpose: the script
	// then clears the fields instead of leaving another provider's names.
	if !strings.Contains(settings, `data-editor-model=""`) {
		t.Error("a provider with unknown models should carry an empty attribute, not be omitted")
	}
	if !strings.Contains(settings, "data-ai-models") {
		t.Error("the model fields are not marked for the script to find")
	}
}

// A panel patched with an inline margin is the shape of this bug: it fixes the
// one instance somebody noticed and leaves the next to be found by eye. The
// account register carried exactly that, and "Последние комментарии" — which
// nobody had patched — sat flush against the cards below it. Spacing between
// panels belongs in the stylesheet, where one rule covers every pair.
//
// Deliberately narrow: the template carries other inline margins for one-off
// nudges inside panels, and rewriting those is not this test's business.
func TestAdminPanelsCarryNoInlineMargins(t *testing.T) {
	f, err := templateFiles.Open("templates/admin.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if m := regexp.MustCompile(`class="adm-panel[^"]*"[^>]*style="[^"]*margin`).FindAll(b, -1); m != nil {
		t.Errorf("a panel is spaced by an inline margin again: %s", m)
	}
}
