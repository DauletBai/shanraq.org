package articles

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// ---- published content renders ----

func TestArticleAndHomeRender(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("writer@t.test", "Parol12345")
	_, slug := app.seedArticle(author, "published")

	if w := app.do(http.MethodGet, "/read/"+slug, nil); w.Code != http.StatusOK {
		t.Errorf("GET /read/%s = %d, want 200", slug, w.Code)
	}
	// The home feed now has a published article to render.
	if w := app.do(http.MethodGet, "/", nil); w.Code != http.StatusOK {
		t.Errorf("home = %d, want 200", w.Code)
	}
	// A draft is not publicly readable.
	_, draftSlug := app.seedArticle(author, "draft")
	if w := app.do(http.MethodGet, "/read/"+draftSlug, nil); w.Code != http.StatusNotFound {
		t.Errorf("draft /read = %d, want 404", w.Code)
	}
}

// ---- studio editor + article creation ----

func TestStudioCreateArticle(t *testing.T) {
	app := newTestApp(t)
	app.createUser("author2@t.test", "Parol12345")
	cookie := app.login("author2@t.test", "Parol12345")

	// The new-article editor renders.
	if w := app.do(http.MethodGet, "/studio/new", nil, withCookie(cookie)); w.Code != http.StatusOK {
		t.Errorf("GET /studio/new = %d, want 200", w.Code)
	}

	// Creating a draft redirects to its editor.
	w := app.do(http.MethodPost, "/studio/new", url.Values{
		"original_lang": {"ru"}, "category": {"economy"},
		"title_ru": {"Заголовок из теста"}, "body_ru": {"Тело статьи"},
	}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("create article = %d, want 303 (%s)", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/studio/a/") {
		t.Errorf("redirect = %q, want /studio/a/...", loc)
	}

	// Missing required fields re-render with an error.
	w = app.do(http.MethodPost, "/studio/new", url.Values{"original_lang": {"ru"}, "category": {"economy"}}, withCookie(cookie))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "alert--error") {
		t.Errorf("empty article should re-render with error, got %d", w.Code)
	}
}

// ---- comments ----

func TestCommentFlow(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("cauthor@t.test", "Parol12345")
	_, slug := app.seedArticle(author, "published")

	app.createUser("commenter@t.test", "Parol12345")
	cookie := app.login("commenter@t.test", "Parol12345")

	// A logged-in user can post a comment.
	w := app.do(http.MethodPost, "/read/"+slug+"/comment", url.Values{"body": {"Хороший разбор, спасибо."}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("comment = %d, want 303", w.Code)
	}
	// The comment now shows on the article page.
	if w := app.do(http.MethodGet, "/read/"+slug, nil); !strings.Contains(w.Body.String(), "Хороший разбор") {
		t.Error("comment should appear on the article page")
	}
	// Anonymous comment attempt is redirected to login.
	if w := app.do(http.MethodPost, "/read/"+slug+"/comment", url.Values{"body": {"x"}}); w.Code != http.StatusSeeOther {
		t.Errorf("anon comment = %d, want 303 (login)", w.Code)
	}
}

// ---- listings ----

func TestListingViewAndMine(t *testing.T) {
	app := newTestApp(t)
	owner := app.createUser("owner@t.test", "Parol12345")
	id := app.seedListing(owner)

	if w := app.do(http.MethodGet, "/listings/"+id.String(), nil); w.Code != http.StatusOK {
		t.Errorf("listing view = %d, want 200", w.Code)
	}
	cookie := app.login("owner@t.test", "Parol12345")
	if w := app.do(http.MethodGet, "/listings/my", nil, withCookie(cookie)); w.Code != http.StatusOK {
		t.Errorf("my listings = %d, want 200", w.Code)
	}
	if w := app.do(http.MethodGet, "/listings/new", nil, withCookie(cookie)); w.Code != http.StatusOK {
		t.Errorf("new listing form = %d, want 200", w.Code)
	}
}

// ---- agent cabinet ----

func TestAgentCabinet(t *testing.T) {
	app := newTestApp(t)
	app.createUser("agent@t.test", "Parol12345")
	cookie := app.login("agent@t.test", "Parol12345")
	if w := app.do(http.MethodGet, "/agent", nil, withCookie(cookie)); w.Code != http.StatusOK {
		t.Errorf("agent cabinet = %d, want 200", w.Code)
	}
	// Registration requires verified email+phone; an unverified user gets the
	// cabinet re-rendered with the "verify first" error, and no profile saved.
	w := app.do(http.MethodPost, "/agent", url.Values{"first_name": {"Асан"}, "last_name": {"Серіков"}}, withCookie(cookie))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "notice--warn") {
		t.Errorf("agent save (unverified) = %d, want 200 with error", w.Code)
	}
}

// ---- admin actions ----

func TestAdminServiceFlagToggle(t *testing.T) {
	app := newTestApp(t)
	app.createUser("admin2@t.test", "Parol12345")
	app.makeStaff("admin2@t.test", "admin")
	cookie := app.login("admin2@t.test", "Parol12345")

	// Put listing promotion into maintenance from the panel.
	w := app.do(http.MethodPost, "/admin/services", url.Values{
		"code": {"listing_promo"}, "status": {"maintenance"},
		"message_ru": {"Тех. работы"},
	}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("service toggle = %d, want 303", w.Code)
	}
	// A non-staff user is redirected to login by RequireSession — never reaching
	// the handler, so services are not touched.
	app.createUser("nonstaff@t.test", "Parol12345")
	pc := app.login("nonstaff@t.test", "Parol12345")
	w = app.do(http.MethodPost, "/admin/services", url.Values{"code": {"comments"}, "status": {"off"}}, withCookie(pc))
	if !strings.Contains(w.Header().Get("Location"), "login") {
		t.Errorf("non-staff toggle should redirect to login, got %d → %q", w.Code, w.Header().Get("Location"))
	}
}

// ---- author verification page ----

func TestAuthorVerifyPage(t *testing.T) {
	app := newTestApp(t)
	app.createUser("verify@t.test", "Parol12345")
	cookie := app.login("verify@t.test", "Parol12345")
	if w := app.do(http.MethodGet, "/studio/author", nil, withCookie(cookie)); w.Code != http.StatusOK {
		t.Errorf("author verify page = %d, want 200", w.Code)
	}
	// Setting a valid name redirects with a success flag.
	w := app.do(http.MethodPost, "/studio/author/name", url.Values{"first_name": {"Асан"}, "last_name": {"Серіков"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Errorf("set author name = %d, want 303", w.Code)
	}
}
