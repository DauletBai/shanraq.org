package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"shanraq.org/pkg/modules/auth"
)

// The policy is the whole balance the feature rests on, and it is the part that
// can be got wrong quietly: too loose and three annoyed readers bury a piece,
// too tight and reader moderation is decoration.
func TestArticleHidePolicy(t *testing.T) {
	cases := []struct {
		name                      string
		weighted, distinct, views int
		want                      bool
	}{
		{"nobody has read it yet", 9, 9, 20, false},
		{"two people, however loudly", 40, 2, 5000, false},
		{"three fresh accounts on a hundred readers", 3, 3, 100, false},
		{"five fresh accounts on a hundred readers", 5, 3, 100, true},
		{"three readers with standing on a hundred", 6, 3, 100, true},
		{"five percent of a wide audience", 25, 12, 500, true},
		{"a handful against a wide audience", 6, 6, 5000, false},
		{"the hard cap on a very wide audience", 25, 25, 100000, true},
		{"exactly at the audience floor", 5, 3, 100, true},
		{"one below the audience floor", 5, 3, 99, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := shouldHideArticle(c.weighted, c.distinct, c.views); got != c.want {
				t.Errorf("shouldHideArticle(%d, %d, %d) = %v, want %v",
					c.weighted, c.distinct, c.views, got, c.want)
			}
		})
	}
}

// Being wrong has to cost something, or reporting is free and the loudest group
// decides what the site publishes.
func TestReporterWeightFallsWithDismissedReports(t *testing.T) {
	cases := []struct {
		karma, dismissed, want int
	}{
		{0, 0, 1},    // a new reader still counts, once
		{500, 0, 5},  // standing earns weight
		{500, 1, 4},  // one mistake costs a step
		{500, 2, 3},  // two, another
		{500, 3, 0},  // a pattern stops counting automatically
		{0, 1, 1},    // a newcomer cannot fall below one…
		{0, 3, 0},    // …until the pattern is established
		{5000, 9, 0}, // no amount of karma buys back a habit
	}
	for _, c := range cases {
		if got := reporterWeight(c.karma, c.dismissed); got != c.want {
			t.Errorf("reporterWeight(karma=%d, dismissed=%d) = %d, want %d",
				c.karma, c.dismissed, got, c.want)
		}
	}
}

// verifiedReader creates an account that may report: registered and confirmed.
//
// The session is minted through the module rather than the login form. The form
// is rate-limited per address — correctly, it is the sign-in door — and a test
// that needs five distinct readers would otherwise be testing the limiter.
func (a *testApp) verifiedReader(email string) *http.Cookie {
	a.t.Helper()
	a.createUser(email, "Sup3r-Secret-Pass!")
	a.exec(`UPDATE auth_users SET email_verified_at = now() WHERE lower(email) = lower($1)`, email)
	_, token, err := a.auth.LoginPassword(context.Background(), email, "Sup3r-Secret-Pass!")
	if err != nil {
		a.t.Fatalf("session for %s: %v", email, err)
	}
	return &http.Cookie{Name: auth.SessionCookieName, Value: token, Path: "/"}
}

// unverifiedReader is an account that exists but has not confirmed its e-mail.
func (a *testApp) unverifiedReader(email string) *http.Cookie {
	a.t.Helper()
	a.createUser(email, "Sup3r-Secret-Pass!")
	return a.sessionFor(email)
}

// sessionFor mints a session cookie for an existing account.
func (a *testApp) sessionFor(email string) *http.Cookie {
	a.t.Helper()
	_, token, err := a.auth.LoginPassword(context.Background(), email, "Sup3r-Secret-Pass!")
	if err != nil {
		a.t.Fatalf("session for %s: %v", email, err)
	}
	return &http.Cookie{Name: auth.SessionCookieName, Value: token, Path: "/"}
}

// The end-to-end path: readers report, the article goes out of sight, the
// decision is in the ledger under the readers' name, and the author can find it.
func TestReadersCanHideAnArticleAndStaffCanPutItBack(t *testing.T) {
	app := newTestApp(t)
	author := app.createUser("rm-author-"+uuid.NewString()[:6]+"@t.test", "Sup3r-Secret-Pass!")
	id, slug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = now() WHERE id = $1`, id)

	if w := app.do(http.MethodGet, "/read/"+slug, nil); w.Code != http.StatusOK {
		t.Fatalf("article should be readable before reports: %d", w.Code)
	}

	// A guest is sent to sign in rather than counted.
	if w := app.do(http.MethodPost, "/read/"+slug+"/report", nil); w.Code != http.StatusSeeOther ||
		!strings.Contains(w.Header().Get("Location"), "/studio/login") {
		t.Fatalf("a guest report should ask for a sign-in: %d %s", w.Code, w.Header().Get("Location"))
	}

	// An account whose e-mail is not confirmed is told why, not counted.
	unverified := app.unverifiedReader("rm-unverified-" + uuid.NewString()[:6] + "@t.test")
	w := app.do(http.MethodPost, "/read/"+slug+"/report", nil, withCookie(unverified))
	if !strings.Contains(w.Header().Get("Location"), "notice=verify") {
		t.Errorf("unverified reporter should be sent to verify: %s", w.Header().Get("Location"))
	}

	// Fix the audience after the reads above, which move the counter: at a
	// hundred readers the policy asks for five percent, and five readers with no
	// standing yet weigh one each.
	app.exec(`UPDATE articles SET views_count = 100 WHERE id = $1`, id)

	// Five verified readers, each weighing one.
	var last *http.Cookie
	for i := 0; i < 5; i++ {
		c := app.verifiedReader("rm-reader-" + uuid.NewString()[:8] + "@t.test")
		last = c
		app.do(http.MethodPost, "/read/"+slug+"/report", url.Values{"reason": {"defamation"}}, withCookie(c))
	}

	if st := app.articleStatus(id); st != "flagged" {
		t.Fatalf("article status is %q after five reports, want flagged", st)
	}
	if w := app.do(http.MethodGet, "/read/"+slug, nil); w.Code != http.StatusNotFound {
		t.Errorf("a hidden article is still public: %d", w.Code)
	}

	// The ledger has to say who decided, and readers are neither staff nor a
	// machine.
	var kind, reason string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT actor_kind, reason_code FROM moderation_actions
		  WHERE target_type='article' AND target_id=$1 AND action='hide'
		  ORDER BY created_at DESC LIMIT 1`, id.String()).Scan(&kind, &reason); err != nil {
		t.Fatalf("no ledger entry for the reader hide: %v", err)
	}
	if kind != "readers" || reason != "reader_reports" {
		t.Errorf("ledger says actor=%q reason=%q, want readers/reader_reports", kind, reason)
	}

	// Reporting again after the hide must not error out on the reader.
	if w := app.do(http.MethodPost, "/read/"+slug+"/report", nil, withCookie(last)); w.Code >= 500 {
		t.Errorf("reporting a hidden article returned %d", w.Code)
	}

	// Staff overrule the readers, through the route a moderator actually uses.
	staffEmail := "rm-staff-" + uuid.NewString()[:6] + "@t.test"
	app.createUser(staffEmail, "Sup3r-Secret-Pass!")
	app.makeStaff(staffEmail, "admin")
	staff := app.sessionFor(staffEmail)
	if w := app.do(http.MethodPost, "/admin/articles/"+id.String()+"/decide",
		url.Values{"decision": {"approve"}, "note": {"checked"}}, withCookie(staff)); w.Code != http.StatusSeeOther {
		t.Fatalf("staff restore: %d %s", w.Code, w.Body.String())
	}
	if st := app.articleStatus(id); st != "published" {
		t.Fatalf("article status is %q after a staff restore, want published", st)
	}
	var live, dismissed int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FILTER (WHERE NOT dismissed), count(*) FILTER (WHERE dismissed)
		   FROM article_reports WHERE article_id = $1`, id).Scan(&live, &dismissed); err != nil {
		t.Fatalf("count reports: %v", err)
	}
	if live != 0 || dismissed != 5 {
		t.Errorf("after the restore %d reports still count and %d are dismissed; want 0 and 5", live, dismissed)
	}
}

// An author cannot report their own piece, and cannot press publish to undo a
// hide readers have applied.
func TestAuthorCannotReportOrUnhideTheirOwnArticle(t *testing.T) {
	app := newTestApp(t)
	email := "rm-self-" + uuid.NewString()[:6] + "@t.test"
	author := app.createUser(email, "Sup3r-Secret-Pass!")
	app.exec(`UPDATE auth_users SET email_verified_at = now() WHERE id = $1`, author)
	cookie := app.sessionFor(email)
	id, slug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET views_count = 500 WHERE id = $1`, id)

	app.do(http.MethodPost, "/read/"+slug+"/report", url.Values{"reason": {"illegal"}}, withCookie(cookie))
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM article_reports WHERE article_id = $1`, id).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 0 {
		t.Errorf("an author reported their own article %d times", n)
	}

	// Hidden by readers, the publish button must not be a way out.
	app.exec(`UPDATE articles SET status = 'flagged' WHERE id = $1`, id)
	app.do(http.MethodPost, "/studio/a/"+id.String()+"/publish", nil, withCookie(cookie))
	if st := app.articleStatus(id); st != "flagged" {
		t.Errorf("the author republished a reader-hidden article: status %q", st)
	}
}
