package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// between returns the markup from the first occurrence of start up to the next
// end, or "" when the pair is not there. Lets a test assert about one table on
// a page that holds a dozen of them.
func between(s, start, end string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return ""
	}
	rest := s[i:]
	j := strings.Index(rest, end)
	if j < 0 {
		return ""
	}
	return rest[:j]
}

// adminSession promotes an account and returns its session cookie and id.
func adminSession(t *testing.T, app *testApp, email string) (*http.Cookie, uuid.UUID) {
	t.Helper()
	id := app.createUser(email, "Parol12345")
	app.makeStaff(email, "admin")
	return app.login(email, "Parol12345"), id
}

// The register is the whole point: before it, the panel could say how many
// accounts held each role and nothing else — not a name, not an address — so a
// complaint about a person could not even be looked up.
func TestAdminSeesEveryUser(t *testing.T) {
	app := newTestApp(t)
	cookie, _ := adminSession(t, app, "root@t.test")
	app.createUser("reader@t.test", "Parol12345")
	app.exec(`UPDATE auth_users SET first_name='Мария', last_name='Иванова', signup_country='KZ'
	          WHERE lower(email)='reader@t.test'`)

	w := app.do(http.MethodGet, "/admin", nil, withCookie(cookie))
	if w.Code != http.StatusOK {
		t.Fatalf("admin page = %d", w.Code)
	}
	body := w.Body.String()
	// Country is shown as flag + ISO code, never a spelled-out country name: the
	// column has to stay narrow, and the code is what is actually stored.
	for _, want := range []string{"reader@t.test", "Иванова", "Мария", "🇰🇿 KZ"} {
		if !strings.Contains(body, want) {
			t.Errorf("the register does not show %q", want)
		}
	}
	// Scoped to the register's own markup: the page elsewhere lists listings and
	// audience panels that legitimately spell country names out.
	if table := between(body, `class="spec spec--sticky adm-users"`, "</table>"); table == "" {
		t.Error("the register table is not on the page")
	} else if strings.Contains(table, "Казахстан") {
		t.Error("the country column spelled the country out instead of showing its code")
	}

	// Search narrows it, and a non-match drops the row.
	w = app.do(http.MethodGet, "/admin", url.Values{"uq": {"иванов"}}, withCookie(cookie))
	if !strings.Contains(w.Body.String(), "reader@t.test") {
		t.Error("search by surname did not find the account")
	}
	w = app.do(http.MethodGet, "/admin", url.Values{"uq": {"неттакого"}}, withCookie(cookie))
	if strings.Contains(w.Body.String(), "reader@t.test") {
		t.Error("search returned an account that does not match it")
	}
}

// Two layers stand between somebody and the register, which carries every name
// and e-mail on the site. A plain reader never reaches the admin group at all
// and is bounced to the login page. An editor does reach it — moderation lives
// there — and must still be refused: reviewing an article is not the same
// authority as deleting the person who wrote it.
func TestUserAdminIsRestricted(t *testing.T) {
	app := newTestApp(t)
	victim := app.createUser("victim@t.test", "Parol12345")

	app.createUser("reader@t.test", "Parol12345")
	reader := app.login("reader@t.test", "Parol12345")
	app.createUser("editor@t.test", "Parol12345")
	app.makeStaff("editor@t.test", "editor")
	editor := app.login("editor@t.test", "Parol12345")

	paths := []string{
		"/admin/users/" + victim.String(),
		"/admin/users/" + victim.String() + "/role",
		"/admin/users/" + victim.String() + "/delete",
	}
	for _, path := range paths {
		if w := app.do(http.MethodPost, path, url.Values{"role": {"admin"}}, withCookie(reader)); w.Code != http.StatusSeeOther {
			t.Errorf("POST %s by a reader = %d, want 303 to the login page", path, w.Code)
		}
		if w := app.do(http.MethodPost, path, url.Values{"role": {"admin"}}, withCookie(editor)); w.Code != http.StatusForbidden {
			t.Errorf("POST %s by an editor = %d, want 403", path, w.Code)
		}
	}
	var alive int
	var role string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*), COALESCE(max(role), '') FROM auth_users WHERE id = $1`, victim).Scan(&alive, &role); err != nil {
		t.Fatal(err)
	}
	if alive != 1 {
		t.Fatal("somebody without the authority deleted an account")
	}
	if role == "admin" {
		t.Fatal("somebody without the authority granted an account admin")
	}
}

// Editing and deleting, and the two ways an administrator could lock themselves
// out of the site: deleting their own account, or demoting it.
func TestAdminEditsAndDeletesAccounts(t *testing.T) {
	app := newTestApp(t)
	cookie, me := adminSession(t, app, "boss@t.test")
	target := app.createUser("target@t.test", "Parol12345")

	// Correct a name.
	form := url.Values{"first_name": {"Даулет"}, "last_name": {"Баймурза"}, "middle_name": {""}, "verified": {"on"}}
	if w := app.do(http.MethodPost, "/admin/users/"+target.String(), form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("edit = %d (%s)", w.Code, w.Body.String())
	}
	var first, last string
	var verified bool
	if err := app.pool.QueryRow(context.Background(),
		`SELECT first_name, last_name, email_verified_at IS NOT NULL FROM auth_users WHERE id = $1`, target).
		Scan(&first, &last, &verified); err != nil {
		t.Fatal(err)
	}
	if first != "Даулет" || last != "Баймурза" || !verified {
		t.Errorf("after the edit: %q %q verified=%v", first, last, verified)
	}

	// Rubbish in a byline must be refused here exactly as at registration.
	bad := url.Values{"first_name": {"<script>"}, "last_name": {"Баймурза"}}
	app.do(http.MethodPost, "/admin/users/"+target.String(), bad, withCookie(cookie))
	if err := app.pool.QueryRow(context.Background(),
		`SELECT first_name FROM auth_users WHERE id = $1`, target).Scan(&first); err != nil {
		t.Fatal(err)
	}
	if first != "Даулет" {
		t.Errorf("a non-name was written into a byline: %q", first)
	}

	// Self-demotion and self-deletion are refused: both end administrative
	// access to the site, and the only way back is a database console.
	app.do(http.MethodPost, "/admin/users/"+me.String()+"/role", url.Values{"role": {"user"}}, withCookie(cookie))
	var myRole string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT role FROM auth_users WHERE id = $1`, me).Scan(&myRole); err != nil {
		t.Fatal(err)
	}
	if myRole != "admin" {
		t.Errorf("the administrator demoted themselves to %q", myRole)
	}
	app.do(http.MethodPost, "/admin/users/"+me.String()+"/delete", url.Values{}, withCookie(cookie))
	var meAlive int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM auth_users WHERE id = $1`, me).Scan(&meAlive); err != nil {
		t.Fatal(err)
	}
	if meAlive != 1 {
		t.Fatal("the administrator deleted themselves out of their own panel")
	}

	// Someone else can be deleted, and their content goes with them.
	app.exec(`INSERT INTO comments (article_id, user_id, body, status)
	          SELECT id, $1, 'x', 'published' FROM articles LIMIT 1`, target)
	if w := app.do(http.MethodPost, "/admin/users/"+target.String()+"/delete", url.Values{}, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("delete = %d", w.Code)
	}
	var gone, orphanComments int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT (SELECT count(*) FROM auth_users WHERE id = $1),
		        (SELECT count(*) FROM comments WHERE user_id = $1)`, target).
		Scan(&gone, &orphanComments); err != nil {
		t.Fatal(err)
	}
	if gone != 0 {
		t.Error("the account survived deletion")
	}
	if orphanComments != 0 {
		t.Error("comments were left behind by a deleted account")
	}
}

// Deleting the last administrator would leave the site with nobody who can
// administer it, and no way back short of a database console.
func TestLastAdminCannotBeRemoved(t *testing.T) {
	app := newTestApp(t)
	cookie, me := adminSession(t, app, "only@t.test")
	// The test database is shared, so other administrators may exist from the
	// bootstrap account or a previous run. Stand them down for the duration —
	// otherwise the guard has nothing to guard and the test passes vacuously.
	app.exec(`UPDATE auth_users SET role = 'user' WHERE role IN ('admin','director') AND id <> $1`, me)

	app.do(http.MethodPost, "/admin/users/"+me.String()+"/delete", url.Values{}, withCookie(cookie))
	var alive int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM auth_users WHERE id = $1`, me).Scan(&alive); err != nil {
		t.Fatal(err)
	}
	if alive != 1 {
		t.Fatal("the site deleted its only administrator")
	}
}
