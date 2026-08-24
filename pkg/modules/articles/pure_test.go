package articles

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"
)

// ---- ad rate card (revenue math) ----

func TestAdSurfacesAndWeights(t *testing.T) {
	sfs := AdSurfaces()
	if len(sfs) != 3+len(Categories) {
		t.Fatalf("AdSurfaces = %d, want %d", len(sfs), 3+len(Categories))
	}
	weight := map[string]int64{}
	for _, s := range sfs {
		weight[s.Code] = s.Weight
	}
	if weight[surfaceHome] != 20 || weight[surfaceRealestate] != 20 {
		t.Errorf("home/realestate weight should be 20, got %d/%d", weight[surfaceHome], weight[surfaceRealestate])
	}
	if weight[surfaceArticles] != 13 {
		t.Errorf("articles weight should be 13, got %d", weight[surfaceArticles])
	}
	if got := adSurfaceWeight10("rubric:economy"); got != 13 {
		t.Errorf("economy rubric weight = %d, want 13", got)
	}
	if got := adSurfaceWeight10("rubric:sport"); got != 10 {
		t.Errorf("niche rubric weight = %d, want 10", got)
	}
	if !isAdSurface(surfaceHome) || !isAdSurface("rubric:"+Categories[0]) {
		t.Error("valid surfaces must pass isAdSurface")
	}
	if isAdSurface("not-a-surface") {
		t.Error("unknown surface must fail isAdSurface")
	}
}

func TestAdFormats(t *testing.T) {
	fs := AdFormats()
	if len(fs) != 4 {
		t.Fatalf("want 4 formats, got %d", len(fs))
	}
	if fs[0].Code != "horizontal" || fs[0].Price30 != 12500 {
		t.Errorf("first format should be horizontal @12500, got %s @%d", fs[0].Code, fs[0].Price30)
	}
	for _, f := range []string{"horizontal", "vertical", "square", "rectangle"} {
		if !isAdFormat(f) {
			t.Errorf("%s should be a valid format", f)
		}
	}
	if isAdFormat("billboard") {
		t.Error("unknown format must fail isAdFormat")
	}
	if AdFormatSlots("rectangle") != 3 || AdFormatSlots("horizontal") != 1 {
		t.Errorf("slots: rectangle=%d horizontal=%d", AdFormatSlots("rectangle"), AdFormatSlots("horizontal"))
	}
	if AdFormatSlots("unknown") != 1 {
		t.Error("unknown format defaults to 1 slot")
	}
}

func TestAdOrderTotal(t *testing.T) {
	// horizontal 30d rate = 12500; home+realestate weight = 20+20 = 40 (×10).
	// total = 12500 * 40 / 10 = 50000, nationwide.
	p := AdOrderTotal("horizontal", []string{surfaceHome, surfaceRealestate}, 30)
	if p.Total != 50000 {
		t.Errorf("total = %d, want 50000", p.Total)
	}
	if p.Surfaces != 2 || p.Weight10 != 40 || p.FormatRate != 12500 {
		t.Errorf("breakdown wrong: %+v", p)
	}
	// Unknown surfaces are ignored, not priced.
	p = AdOrderTotal("square", []string{surfaceHome, "junk"}, 7)
	if p.Surfaces != 1 || p.Weight10 != 20 {
		t.Errorf("junk surface should be dropped: %+v", p)
	}
	// Unknown format falls back to rectangle; unknown duration to the 3-day rate.
	p = AdOrderTotal("nonsense", []string{surfaceHome}, 999)
	if p.Format != "rectangle" {
		t.Errorf("unknown format should fall back to rectangle, got %s", p.Format)
	}
	if p.FormatRate != 1100 { // rectangle 3-day
		t.Errorf("unknown duration should use the 3-day rate 1100, got %d", p.FormatRate)
	}
}

func TestAdSurfaceFormatPrice(t *testing.T) {
	// horizontal 30d on home, nationwide: 12500 * 20 / 10 = 25000.
	if got := AdSurfaceFormatPrice("horizontal", surfaceHome, 30); got != 25000 {
		t.Errorf("= %d, want 25000", got)
	}
	// Unknown format → 0.
	if got := AdSurfaceFormatPrice("bogus", surfaceHome, 30); got != 0 {
		t.Errorf("unknown format = %d, want 0", got)
	}
	// Unknown duration falls back to the 30-day rate.
	if got := AdSurfaceFormatPrice("square", surfaceHome, 999); got != AdSurfaceFormatPrice("square", surfaceHome, 30) {
		t.Error("unknown duration should equal the 30-day price")
	}
}

func TestSurfaceLabelKeyAndRatesJSON(t *testing.T) {
	if SurfaceLabelKey("rubric:economy") != "cat.economy" {
		t.Errorf("rubric label = %q", SurfaceLabelKey("rubric:economy"))
	}
	if SurfaceLabelKey(surfaceHome) != "adv.surf_home" {
		t.Errorf("home label = %q", SurfaceLabelKey(surfaceHome))
	}
	if adRubricSurface("economy") != "rubric:economy" {
		t.Error("adRubricSurface")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(AdRatesJSON()), &payload); err != nil {
		t.Fatalf("AdRatesJSON is not valid JSON: %v", err)
	}
	if payload["formatPrice"] == nil || payload["weights"] == nil {
		t.Error("AdRatesJSON must carry formatPrice and weights")
	}
	if got := AdDurations(); len(got) != 4 || got[0] != 3 {
		t.Errorf("AdDurations = %v", got)
	}
}

// ---- categories ----

func TestCategories(t *testing.T) {
	if len(Categories) == 0 {
		t.Fatal("no categories defined")
	}
	c := Categories[0]
	if !IsCategory(c) {
		t.Errorf("%q should be a category", c)
	}
	if IsCategory("no-such-category") {
		t.Error("unknown category must be false")
	}
	if NormalizeCategory(c) != c {
		t.Error("valid category should pass through NormalizeCategory")
	}
	if NormalizeCategory("no-such-category") != CategoryGeneral {
		t.Errorf("invalid category should become %q", CategoryGeneral)
	}
	// Find any subcategory and verify the parent round-trip.
	var sub, parent string
	for cat, subs := range Subcategories {
		if len(subs) > 0 {
			sub, parent = subs[0], cat
			break
		}
	}
	if sub != "" {
		if !IsSubcategory(sub) {
			t.Errorf("%q should be a subcategory", sub)
		}
		if SubcategoryParent(sub) != parent {
			t.Errorf("parent of %q = %q, want %q", sub, SubcategoryParent(sub), parent)
		}
		if NormalizeSubcategory(parent, sub) != sub {
			t.Error("sub under its parent should pass NormalizeSubcategory")
		}
		if NormalizeSubcategory("wrong-parent", sub) != "" {
			t.Error("sub under a wrong parent should be dropped")
		}
	}
}

// ---- listings pure logic ----

func TestListingHelpers(t *testing.T) {
	if !isDealType(DealTypes[0]) || isDealType("barter") {
		t.Error("isDealType")
	}
	if !isPropertyType(PropertyTypes[0]) || isPropertyType("castle") {
		t.Error("isPropertyType")
	}
	if BannerPrice(1) != 990 || BannerPrice(2) != 1890 || BannerPrice(7) != 4990 || BannerPrice(99) != 990 {
		t.Errorf("banner prices off: 1=%d 2=%d 7=%d 99=%d", BannerPrice(1), BannerPrice(2), BannerPrice(7), BannerPrice(99))
	}
	if d := BannerDays(); len(d) != 7 || d[0] != 1 || d[6] != 7 {
		t.Errorf("BannerDays = %v", d)
	}
}

func TestListingStateMethods(t *testing.T) {
	now := time.Now()
	future := now.Add(48 * time.Hour)
	past := now.Add(-48 * time.Hour)

	active := Listing{ExpiresAt: future, PromotedUntil: &future, FeaturedUntil: &future, BannerUntil: &future, AgentID: "u1"}
	if active.Expired() {
		t.Error("future expiry should not be expired")
	}
	if !active.Promoted() || !active.Featured() || !active.Banner() {
		t.Error("future-dated boosts should be active")
	}
	if !active.ByAgent() {
		t.Error("ByAgent should be true with an agent id")
	}
	if active.DaysLeft() < 1 {
		t.Errorf("DaysLeft should be >=1, got %d", active.DaysLeft())
	}

	dead := Listing{ExpiresAt: past, PromotedUntil: &past, FeaturedUntil: &past, BannerUntil: &past}
	if !dead.Expired() {
		t.Error("past expiry should be expired")
	}
	if dead.Promoted() || dead.Featured() || dead.Banner() || dead.ByAgent() {
		t.Error("past-dated boosts and empty agent should be inactive")
	}
	if dead.DaysLeft() != 0 {
		t.Errorf("expired DaysLeft = %d, want 0", dead.DaysLeft())
	}

	// Address / Location render only the set parts, in order.
	l := Listing{Microdistrict: "Самал-2", Street: "Абая", House: "15", City: "Алматы", Region: "Алматы", Country: "Казахстан"}
	if l.Address() != "Самал-2, Абая, 15" {
		t.Errorf("Address = %q", l.Address())
	}
	if (Listing{Street: "Абая"}).Address() != "Абая" {
		t.Error("Address should skip blank parts")
	}
	if l.HasPin() {
		t.Error("no coords should mean no pin")
	}
	lat, lng := 43.2, 76.9
	if !(Listing{Lat: &lat, Lng: &lng}).HasPin() {
		t.Error("coords should mean a pin")
	}
	if l.Location() != "Алматы, Алматы, Казахстан" {
		t.Errorf("Location = %q", l.Location())
	}
}

func TestListingAmenities(t *testing.T) {
	l := Listing{Amenities: []string{"furniture", "elevator"}}
	if !l.HasAmenity("furniture") || l.HasAmenity("pool") {
		t.Error("Listing.HasAmenity")
	}
	in := ListingInput{Amenities: []string{"internet"}}
	if !in.HasAmenity("internet") || in.HasAmenity("garden") {
		t.Error("ListingInput.HasAmenity")
	}
}

// ---- ad surface routing ----

func TestAdSurfaceFor(t *testing.T) {
	cases := map[string]string{
		"/listings":              surfaceRealestate,
		"/listings/abc":          surfaceRealestate,
		"/read/some-slug":        surfaceArticles,
		"/":                      surfaceHome,
		"/?cat=" + Categories[0]: "rubric:" + Categories[0],
		"/?cat=bogus":            surfaceHome,
		"/about":                 "",
	}
	for path, want := range cases {
		r := httptest.NewRequest("GET", path, nil)
		if got := adSurfaceFor(r); got != want {
			t.Errorf("adSurfaceFor(%q) = %q, want %q", path, got, want)
		}
	}
}

// ---- revenue split ----

func TestRevenueShare(t *testing.T) {
	if !IsAIAgentAuthor(SanaAuthorID) {
		t.Error("Sana should be an AI agent author")
	}
	if IsAIAgentAuthor("some-human-uuid") {
		t.Error("a random id is not an AI agent")
	}
	if a, p := RevenueShare(SanaAuthorID); a != 0 || p != 100 {
		t.Errorf("AI agent share = %d/%d, want 0/100", a, p)
	}
	if a, p := RevenueShare("human"); a != RevenueAuthorPct || p != RevenuePlatformPct {
		t.Errorf("human share = %d/%d, want %d/%d", a, p, RevenueAuthorPct, RevenuePlatformPct)
	}
	if RevenueAuthorPct+RevenuePlatformPct != 100 {
		t.Error("shares must sum to 100")
	}
}

// ---- service flags (pure helpers) ----

func TestServiceFlagPureHelpers(t *testing.T) {
	for _, s := range []string{svcOn, svcInviteOnly, svcMaintenance, svcOff} {
		if !isServiceStatus(s) {
			t.Errorf("%q should be a valid status", s)
		}
	}
	if isServiceStatus("paused") {
		t.Error("unknown status must be rejected")
	}
	if !isKnownService(SvcRegistration) || isKnownService(SvcSite) || isKnownService("nope") {
		t.Error("isKnownService: site is not a per-feature service")
	}
	if !validServiceCode(SvcSite) || !validServiceCode(SvcComments) || validServiceCode("nope") {
		t.Error("validServiceCode")
	}
	if !invitable(SvcRegistration) || invitable(SvcAdOrders) {
		t.Error("invitable: only free functions support invite_only")
	}

	f := ServiceFlag{Code: SvcRegistration, Status: svcOn, MessageRU: "рус", MessageKZ: "қаз", MessageEN: "en"}
	if !f.Available() {
		t.Error("on → available")
	}
	if f.Message(LangKZ) != "қаз" || f.Message(LangEN) != "en" || f.Message(LangRU) != "рус" {
		t.Error("Message localization")
	}
	if !f.Invitable() {
		t.Error("registration flag should be invitable")
	}
	if f.HasTimer() || f.UntilUnixMillis() != 0 {
		t.Error("no timer set")
	}
	when := time.Now().Add(time.Hour)
	f2 := ServiceFlag{Status: svcMaintenance, Until: when}
	if f2.Available() {
		t.Error("maintenance → not available")
	}
	if !f2.HasTimer() || f2.UntilUnixMillis() != when.UnixMilli() {
		t.Error("timer millis")
	}
}

// ---- agent name helpers ----

func TestComposeNameAndStatus(t *testing.T) {
	if composeName("Даулет", "Баймурза") != "Даулет Баймурза" {
		t.Error("composeName both")
	}
	if composeName("Даулет", "") != "Даулет" {
		t.Error("composeName first only")
	}
	if composeName("  ", "  ") != "" {
		t.Error("composeName blank")
	}
	for _, s := range []string{agentPending, agentVerified, agentRejected} {
		if !isAgentStatus(s) {
			t.Errorf("%q should be a valid agent status", s)
		}
	}
	if isAgentStatus("banned") {
		t.Error("unknown agent status must be rejected")
	}
	if !(Agent{Status: agentVerified}).Verified() || (Agent{Status: agentPending}).Verified() {
		t.Error("Agent.Verified")
	}
}

// ParseFloat is happy to return NaN, ±Inf and negatives, and the old code threw
// the error away — so those reached the database and then the filters, where a
// NaN compares false against everything and silently removes a listing from
// every search.
func TestParseArea(t *testing.T) {
	ok := map[string]float64{"72": 72, "72,5": 72.5, "72.5": 72.5, " 100 ": 100, "0": 0}
	for in, want := range ok {
		if got := parseArea(in); got != want {
			t.Errorf("parseArea(%q) = %v, want %v", in, got, want)
		}
	}
	for _, bad := range []string{"NaN", "nan", "Inf", "-Inf", "+Inf", "-1", "-0.5", "1e308", "999999999", "", "abc", "12abc"} {
		if got := parseArea(bad); got != 0 {
			t.Errorf("parseArea(%q) = %v, want 0", bad, got)
		}
	}
}

// The form fields carry paths returned by our own uploader, but a hand-made
// POST could put any absolute URL there — a remote image is a tracking pixel
// that fires for every visitor, and its target can be swapped after moderation.
func TestMediaPath(t *testing.T) {
	for _, good := range []string{"/media/2026/08/a.jpg", "/media/plan.pdf"} {
		if got, ok := mediaPath(good); !ok || got != good {
			t.Errorf("mediaPath(%q) = %q,%v — want it kept", good, got, ok)
		}
	}
	bad := []string{
		"https://evil.example/pixel.gif", "http://evil.example/x.png",
		"//evil.example/x.png", "/media//evil.example/x.png",
		"/media/../../etc/passwd", "/static/brand/shanraq.svg",
		"javascript:alert(1)", "data:image/png;base64,AAAA",
		"", "   ", "/media/a b.jpg", "/media/a\nb.jpg",
	}
	for _, u := range bad {
		if got, ok := mediaPath(u); ok {
			t.Errorf("mediaPath(%q) accepted as %q", u, got)
		}
	}
}
