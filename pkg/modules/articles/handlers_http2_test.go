package articles

import (
	"context"
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

// A comment is published under the writer's real name. Until this existed the
// only way to unsay one was to write to an administrator — and the delete has to
// be theirs alone, or one reader could erase another's words.
func TestCommentDeleteIsOwnerOnly(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("cdauthor@t.test", "Parol12345")
	articleID, slug := app.seedArticle(author, "published")

	app.createUser("mine@t.test", "Parol12345")
	mine := app.login("mine@t.test", "Parol12345")
	if w := app.do(http.MethodPost, "/read/"+slug+"/comment",
		url.Values{"body": {"Сказал и передумал."}}, withCookie(mine)); w.Code != http.StatusSeeOther {
		t.Fatalf("post comment = %d", w.Code)
	}
	var id string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM comments WHERE article_id = $1`, articleID).Scan(&id); err != nil {
		t.Fatal(err)
	}

	// A stranger must not be able to delete it, even knowing the id.
	app.createUser("notmine@t.test", "Parol12345")
	other := app.login("notmine@t.test", "Parol12345")
	if w := app.do(http.MethodPost, "/read/"+slug+"/comment/"+id+"/delete", url.Values{}, withCookie(other)); w.Code != http.StatusNotFound {
		t.Errorf("a stranger deleting a comment = %d, want 404", w.Code)
	}
	if w := app.do(http.MethodGet, "/read/"+slug, nil); !strings.Contains(w.Body.String(), "Сказал и передумал") {
		t.Fatal("the comment was destroyed by someone who did not write it")
	}

	// Its author can.
	if w := app.do(http.MethodPost, "/read/"+slug+"/comment/"+id+"/delete", url.Values{}, withCookie(mine)); w.Code != http.StatusSeeOther {
		t.Fatalf("the author deleting their own comment = %d, want 303", w.Code)
	}
	if w := app.do(http.MethodGet, "/read/"+slug, nil); strings.Contains(w.Body.String(), "Сказал и передумал") {
		t.Error("the comment is still on the page after its author deleted it")
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

// Exercises the real browser form and handler: a Russian agency must be able
// to submit its 10-digit INN once both account contacts are verified.
func TestRussianAgencyRegistrationJourney(t *testing.T) {
	app := newTestApp(t)
	uid := app.createUser("ru-agency-journey@t.test", "Parol12345")
	cookie := app.login("ru-agency-journey@t.test", "Parol12345")
	form := url.Values{
		"lang": {"ru"}, "kind": {"agency"}, "first_name": {"Ирина"},
		"last_name": {"Тестова"}, "agency": {"Тестовое агентство"},
		"bin": {"1234567890"},
	}
	page := app.do(http.MethodGet, "/agent?lang=ru", nil, withCookie(cookie))
	if page.Code != http.StatusOK || !strings.Contains(page.Body.String(), `pattern="[0-9]{10}([0-9]{2})?"`) {
		t.Fatalf("agency form does not accept 10-digit INN: %d", page.Code)
	}
	blocked := app.do(http.MethodPost, "/agent", form, withCookie(cookie))
	if blocked.Code != http.StatusOK || !strings.Contains(blocked.Body.String(), "подтвердите email и телефон") || !strings.Contains(blocked.Body.String(), "1234567890") {
		t.Fatalf("unverified agency must see reason and keep entered INN: %d", blocked.Code)
	}
	app.exec(`UPDATE auth_users SET email_verified_at = now(), phone_verified_at = now() WHERE id = $1`, uid)
	accepted := app.do(http.MethodPost, "/agent", form, withCookie(cookie))
	if accepted.Code != http.StatusSeeOther || !strings.Contains(accepted.Header().Get("Location"), "saved=1") {
		t.Fatalf("verified agency registration: %d, location %q", accepted.Code, accepted.Header().Get("Location"))
	}
	var status, bin string
	if err := app.pool.QueryRow(t.Context(), `SELECT status, bin FROM re_agents WHERE user_id = $1`, uid).Scan(&status, &bin); err != nil {
		t.Fatal(err)
	}
	if status != "pending" || bin != "1234567890" {
		t.Fatalf("saved agency = status %q, INN %q", status, bin)
	}
}

func TestAdvertiserValidationKeepsDraftWithoutPretendingRegistration(t *testing.T) {
	app := newTestApp(t)
	uid := app.createUser("advertiser-journey@t.test", "Parol12345")
	cookie := app.login("advertiser-journey@t.test", "Parol12345")
	bad := app.do(http.MethodPost, "/advertise/company", url.Values{
		"lang": {"ru"}, "company_name": {"Компания Тест"},
		"bin": {"1234567890"}, "contact_name": {"Ирина"},
	}, withCookie(cookie))
	if bad.Code != http.StatusOK || !strings.Contains(bad.Body.String(), "телефон обязательны") || !strings.Contains(bad.Body.String(), `value="Компания Тест"`) {
		t.Fatalf("invalid company must explain missing phone and retain input: %d", bad.Code)
	}
	if strings.Contains(bad.Body.String(), "adv.new_order") || strings.Contains(bad.Body.String(), `action="/advertise/order"`) {
		t.Fatal("invalid company is displayed as a registered advertiser")
	}
	good := app.do(http.MethodPost, "/advertise/company", url.Values{
		"lang": {"ru"}, "company_name": {"Компания Тест"},
		"bin": {"1234567890"}, "contact_name": {"Ирина"},
		"contact_phone": {"+79030000000"},
	}, withCookie(cookie))
	if good.Code != http.StatusSeeOther || !strings.Contains(good.Header().Get("Location"), "saved=company") {
		t.Fatalf("valid advertiser registration: %d", good.Code)
	}
	// The paid-order service is in maintenance. A direct HTTP POST must not
	// bypass the disabled browser button and silently create an order.
	blocked := app.do(http.MethodPost, "/advertise/order", url.Values{
		"lang": {"ru"}, "title": {"Тест"}, "target_url": {"https://example.org"},
		"format": {"rectangle"}, "surface": {"home"}, "duration": {"7"},
	}, withCookie(cookie))
	if blocked.Code != http.StatusSeeOther || blocked.Header().Get("Location") != "/advertise?blocked=orders" {
		t.Fatalf("paused paid order must explain status in cabinet: %d, %q", blocked.Code, blocked.Header().Get("Location"))
	}
	message := app.do(http.MethodGet, blocked.Header().Get("Location")+"&lang=ru", nil, withCookie(cookie))
	if message.Code != http.StatusOK || !strings.Contains(message.Body.String(), "Приём платных заказов сейчас приостановлен") {
		t.Fatalf("paused paid order explanation missing: %d", message.Code)
	}
	var count int
	if err := app.pool.QueryRow(t.Context(), `SELECT count(*) FROM ad_orders o JOIN advertisers a ON a.id = o.advertiser_id WHERE a.owner_id = $1`, uid).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("paused order unexpectedly created %d rows", count)
	}
}
