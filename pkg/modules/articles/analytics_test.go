package articles

import "testing"

func TestPageKind(t *testing.T) {
	cases := map[string]string{
		"/":                   "home",
		"/read":               "article",
		"/read/some-slug":     "article",
		"/listings":           "listings",
		"/listings/new":       "listings",
		"/listings/my":        "listings",
		"/listings/abc-123":   "listing",
		"/author/42":          "author",
		"/agent/7":            "agent",
		"/favorites":          "favorites",
		"/about":              "static",
		"/privacy":            "static",
		"/terms":              "static",
		"/api/geo/roots":      "",        // not counted
		"/static/css/x.css":   "",        // not counted
		"/studio/new":         "",        // staff area, not a public page
		"/admin":              "",        // not counted
		"/read/slug/progress": "article", // still under /read/ prefix (POST is filtered by method upstream)
	}
	for path, want := range cases {
		if got := pageKind(path); got != want {
			t.Errorf("pageKind(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestTrackedEventsClosedSet(t *testing.T) {
	if !trackedEvents["show_contact"] {
		t.Error("show_contact should be a tracked event")
	}
	if trackedEvents["arbitrary_injected"] {
		t.Error("unknown events must not be tracked")
	}
}
