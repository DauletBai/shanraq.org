package articles

import "testing"

// The breadcrumb is what turns a bare URL in search results into a labelled
// trail. Google requires the last crumb to be the page itself and unlinked.
func TestBreadcrumbLD(t *testing.T) {
	page := &ArticlePage{
		// ArticlePage declares its own Title, which shadows Base.Title — the
		// breadcrumb and the headline both read the outer one.
		Base:        Base{Lang: "ru", SiteURL: "https://shanraq.org"},
		Slug:        "test",
		Title:       "Заголовок",
		Category:    "world",
		Subcategory: "europe",
	}
	ld := breadcrumbLD(page)
	if ld["@type"] != "BreadcrumbList" {
		t.Fatalf("@type = %v", ld["@type"])
	}
	items, _ := ld["itemListElement"].([]map[string]any)
	if len(items) != 4 {
		t.Fatalf("crumbs = %d, want site → category → subcategory → article", len(items))
	}
	for i, it := range items {
		if it["position"] != i+1 {
			t.Errorf("crumb %d has position %v", i, it["position"])
		}
	}
	if items[0]["name"] != "Shanraq.org" {
		t.Errorf("first crumb = %v", items[0]["name"])
	}
	if _, linked := items[3]["item"]; linked {
		t.Error("the last crumb must not carry a link")
	}
	if items[3]["name"] != "Заголовок" {
		t.Errorf("last crumb = %v", items[3]["name"])
	}

	// An article filed without a subcategory gets a three-step trail, not one
	// with a hole where the missing level was.
	page.Subcategory = ""
	if got := len(breadcrumbLD(page)["itemListElement"].([]map[string]any)); got != 3 {
		t.Errorf("crumbs without a subcategory = %d, want 3", got)
	}
	page.Category = ""
	if got := len(breadcrumbLD(page)["itemListElement"].([]map[string]any)); got != 2 {
		t.Errorf("crumbs without a category = %d, want 2", got)
	}
}
