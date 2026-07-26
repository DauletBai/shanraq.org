package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func (a *testApp) articleStatus(id uuid.UUID) string {
	var s string
	_ = a.pool.QueryRow(context.Background(), `SELECT status FROM articles WHERE id = $1`, id).Scan(&s)
	return s
}

// A regular author with no verified phone cannot publish: they are sent to the
// author-verification page.
func TestNonStaffPublishRequiresVerification(t *testing.T) {
	app := newTestApp(t)
	email := "auth-" + uuid.NewString()[:8] + "@example.com"
	pass := "Str0ng-Pass-123"
	uid := app.createUser(email, pass)
	cookie := app.login(email, pass)
	id, _ := app.seedArticle(uid, "draft")

	w := app.do(http.MethodPost, "/studio/a/"+id.String()+"/publish", nil, withCookie(cookie))
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "/studio/author") {
		t.Fatalf("non-staff publish should require verification: code=%d loc=%s", w.Code, w.Header().Get("Location"))
	}
	if st := app.articleStatus(id); st != "draft" {
		t.Errorf("article should stay draft, got %q", st)
	}
}

// Leadership (admin) publishes immediately without email/phone verification,
// even when public article submission is switched off.
func TestStaffPublishesWithoutVerification(t *testing.T) {
	app := newTestApp(t)
	email := "ceo-" + uuid.NewString()[:8] + "@example.com"
	pass := "Str0ng-Pass-123"
	uid := app.createUser(email, pass)
	app.makeStaff(email, "admin")
	cookie := app.login(email, pass)

	// Close public submission — staff must still get through.
	app.exec(`INSERT INTO service_flags (code, status, updated_at) VALUES ('article_submission','off',NOW())
		ON CONFLICT (code) DO UPDATE SET status='off', updated_at=NOW()`)
	t.Cleanup(func() { app.exec(`DELETE FROM service_flags WHERE code='article_submission'`) })

	id, _ := app.seedArticle(uid, "draft")
	w := app.do(http.MethodPost, "/studio/a/"+id.String()+"/publish", nil, withCookie(cookie))
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "ok=published") {
		t.Fatalf("staff publish: code=%d loc=%s", w.Code, w.Header().Get("Location"))
	}
	if st := app.articleStatus(id); st != "published" {
		t.Errorf("staff article should be published immediately, got %q", st)
	}
}

// The bio saves and then appears on the author page together with the team badge.
func TestBioAndAuthorPageTeamBadge(t *testing.T) {
	app := newTestApp(t)
	email := "bio-" + uuid.NewString()[:8] + "@example.com"
	pass := "Str0ng-Pass-123"
	uid := app.createUser(email, pass)
	app.makeStaff(email, "admin")
	cookie := app.login(email, pass)

	bioRU := "Основатель и CEO проекта Shanraq."
	bioKZ := "Shanraq жобасының негізін қалаушы."
	w := app.do(http.MethodPost, "/studio/bio",
		url.Values{"bio_ru": {bioRU}, "bio_kz": {bioKZ}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || !strings.Contains(w.Header().Get("Location"), "bio_set") {
		t.Fatalf("bio save: code=%d loc=%s", w.Code, w.Header().Get("Location"))
	}

	// Publish one article so the author page renders, then check it.
	id, _ := app.seedArticle(uid, "draft")
	_ = app.do(http.MethodPost, "/studio/a/"+id.String()+"/publish", nil, withCookie(cookie))

	// Russian author page shows the RU bio + the team badge.
	w = app.do(http.MethodGet, "/author/"+uid.String(), url.Values{"lang": {"ru"}})
	if w.Code != http.StatusOK {
		t.Fatalf("author page = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, bioRU) {
		t.Error("RU author page should show the RU bio")
	}
	if !strings.Contains(body, "team-badge") {
		t.Error("author page should show the team badge for a staff author")
	}

	// Kazakh author page shows the KZ bio (localization works).
	w = app.do(http.MethodGet, "/author/"+uid.String(), url.Values{"lang": {"kz"}})
	if !strings.Contains(w.Body.String(), bioKZ) {
		t.Error("KZ author page should show the KZ bio")
	}
}
