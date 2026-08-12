package syndicate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func testModule() *Module {
	return &Module{baseURL: "https://shanraq.org", log: zap.NewNop()}
}

func TestRenderRSS(t *testing.T) {
	m := testModule()
	when := time.Date(2026, 7, 13, 10, 0, 0, 0, time.UTC)
	out, err := m.renderRSS("ru", []feedEntry{
		{Slug: "ekonomika", Title: "Экономика 2026", Summary: "Разбор цифр", Lang: "ru", Modified: when},
	})
	if err != nil {
		t.Fatalf("renderRSS: %v", err)
	}
	s := string(out)
	for _, want := range []string{
		`<?xml version="1.0"`,
		`<rss version="2.0">`,
		"<title>Экономика 2026</title>",
		"https://shanraq.org/read/ekonomika?lang=ru",
		"<language>ru</language>",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("RSS missing %q\n---\n%s", want, s)
		}
	}
}

func TestBuildTelegramMessageEscapes(t *testing.T) {
	msg := buildTelegramMessage("Цены <выросли> & упали", "Кратко про <тэги>", "https://shanraq.org/read/x?lang=ru")
	if strings.Contains(msg, "<выросли>") {
		t.Errorf("title not HTML-escaped: %s", msg)
	}
	for _, want := range []string{"&lt;выросли&gt;", "<b>", "🔗", "https://shanraq.org/read/x?lang=ru"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q: %s", want, msg)
		}
	}
}

func TestTelegramDisabledByDefault(t *testing.T) {
	m := testModule()
	if m.TelegramEnabled() {
		t.Fatal("telegram should be disabled without config")
	}
	// EnqueuePublish must be a safe no-op when disabled (nil store tolerated).
	if err := m.EnqueuePublish(context.Background(), nil, uuid.New()); err != nil {
		t.Fatalf("EnqueuePublish no-op: %v", err)
	}
}

func TestRenderDigest(t *testing.T) {
	m := testModule()
	when := time.Date(2026, 7, 13, 10, 0, 0, 0, time.UTC)
	subject, body := m.renderDigest("ru", []feedEntry{
		{Slug: "ekonomika", Title: "Экономика 2026", Lang: "ru", Modified: when},
	}, "tok123")
	if subject != "Shanraq.org: обзор недели" {
		t.Errorf("subject = %q", subject)
	}
	for _, want := range []string{"Экономика 2026", "https://shanraq.org/read/ekonomika?lang=ru", "Отписаться", "/unsubscribe?token=tok123"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q:\n%s", want, body)
		}
	}
}

type sentMail struct {
	To      string
	Subject string
	Body    string
	Headers map[string]string
}

type fakeMailer struct {
	sent []string
	mail []sentMail
}

func (f *fakeMailer) Send(ctx context.Context, to, subject, body string) error {
	return f.SendWithHeaders(ctx, to, subject, body, nil)
}

func (f *fakeMailer) SendWithHeaders(_ context.Context, to, subject, body string, h map[string]string) error {
	f.sent = append(f.sent, to+"|"+subject)
	f.mail = append(f.mail, sentMail{To: to, Subject: subject, Body: body, Headers: h})
	return nil
}

// safeBack must keep the reader on the page they subscribed from without ever
// becoming an open redirect — the value arrives in a form field.
func TestSafeBack(t *testing.T) {
	cases := []struct{ in, want string }{
		{"/read/kurultai", "/read/kurultai?lang=ru"},
		{"/", "/?lang=ru"},
		{"", "/?lang=ru"},
		{"/read/x?lang=en&foo=1", "/read/x?lang=ru"}, // our own query wins
		{"//evil.example.com", "/?lang=ru"},          // protocol-relative
		{"https://evil.example.com", "/?lang=ru"},
		{"  /listings  ", "/listings?lang=ru"},
	}
	for _, c := range cases {
		if got := safeBack(c.in, "ru"); got != c.want {
			t.Errorf("safeBack(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaskEmail(t *testing.T) {
	cases := []struct{ in, want string }{
		{"baimurza.daulet@gmail.com", "bai***@gmail.com"},
		{"ab@mail.kz", "a***@mail.kz"},
		{"notanemail", ""},
	}
	for _, c := range cases {
		if got := maskEmail(c.in); got != c.want {
			t.Errorf("maskEmail(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPlausibleEmail(t *testing.T) {
	for _, ok := range []string{"a@b.kz", "user.name+tag@mail.example.com"} {
		if !plausibleEmail(ok) {
			t.Errorf("%q should be accepted", ok)
		}
	}
	for _, bad := range []string{"", "nope", "@b.kz", "a@", "a@b", "a b@c.kz", "a@b.kz\r\nBcc: x@y.z", "a@@b.kz"} {
		if plausibleEmail(bad) {
			t.Errorf("%q should be rejected", bad)
		}
	}
}

// TestDigestIntegration exercises subscribe → SendDigest → unsubscribe against a
// real DB with a fake mailer. Skipped unless SHANRAQ_TEST_DB is set.
func TestDigestIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the digest integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	// A published article within the last 7 days, and a subscriber.
	authorID := uuid.New()
	articleID := uuid.New()
	email := "digest-" + articleID.String()[:8] + "@t.test"
	_, _ = pool.Exec(ctx, `INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`, authorID, "dg-"+authorID.String()+"@t.test")
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, authorID)
		_, _ = pool.Exec(ctx, `DELETE FROM subscribers WHERE email=$1`, email)
	})
	if _, err := pool.Exec(ctx, `INSERT INTO articles (id, author_id, slug, original_lang, status, published_at) VALUES ($1,$2,$3,'ru','published',NOW())`, articleID, authorID, "dg-"+articleID.String()[:8]); err != nil {
		t.Fatalf("insert article: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES ($1,'ru','Дайджест тест','Аннотация','Тело','human','ready')`, articleID); err != nil {
		t.Fatalf("insert translation: %v", err)
	}

	fm := &fakeMailer{}
	m := &Module{db: pool, baseURL: "https://shanraq.org", log: zap.NewNop(), mailer: fm, emailEnabled: true}

	confirmTok, err := m.subscribe(ctx, email, "ru")
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	// The whole point of double opt-in: an unconfirmed address is a request,
	// not a subscriber, and must never receive mail.
	if _, err := m.SendDigest(ctx); err != nil {
		t.Fatalf("SendDigest (pending): %v", err)
	}
	for _, s := range fm.sent {
		if strings.HasPrefix(s, email+"|") {
			t.Fatalf("digest sent to an unconfirmed address")
		}
	}

	if _, ok := m.confirmSubscriber(ctx, confirmTok); !ok {
		t.Fatal("confirmSubscriber rejected a fresh token")
	}
	// Confirming twice must not resurrect a spent token.
	if _, ok := m.confirmSubscriber(ctx, confirmTok); ok {
		t.Error("a spent confirm token was accepted again")
	}

	sent, err := m.SendDigest(ctx)
	if err != nil {
		t.Fatalf("SendDigest: %v", err)
	}
	if sent < 1 {
		t.Fatalf("expected ≥1 sent, got %d", sent)
	}
	var got *sentMail
	for i := range fm.mail {
		if fm.mail[i].To == email {
			got = &fm.mail[i]
		}
	}
	if got == nil {
		t.Fatalf("digest not sent to %s (sent: %v)", email, fm.sent)
	}
	// Bulk mail has to carry one-click unsubscribe or Gmail scores it as spam.
	if u := got.Headers["List-Unsubscribe"]; !strings.Contains(u, "/unsubscribe?token=") {
		t.Errorf("List-Unsubscribe header = %q", u)
	}
	if p := got.Headers["List-Unsubscribe-Post"]; p != "List-Unsubscribe=One-Click" {
		t.Errorf("List-Unsubscribe-Post = %q", p)
	}

	var token string
	if err := pool.QueryRow(ctx, `SELECT unsubscribe_token FROM subscribers WHERE email=$1`, email).Scan(&token); err != nil {
		t.Fatalf("read token: %v", err)
	}
	// Resolving the token must not delete anything: the GET behind the link in
	// every email only asks, because mail scanners follow links on their own.
	gotEmail, _, ok := m.lookupSubscriber(ctx, token)
	if !ok || gotEmail != email {
		t.Fatalf("lookupSubscriber = %q, %v", gotEmail, ok)
	}
	var cnt int
	_ = pool.QueryRow(ctx, `SELECT COUNT(*) FROM subscribers WHERE email=$1`, email).Scan(&cnt)
	if cnt != 1 {
		t.Fatalf("lookup deleted the subscriber (count=%d)", cnt)
	}

	// Only the POST removes it.
	if _, ok := m.unsubscribe(ctx, token); !ok {
		t.Fatal("unsubscribe rejected a valid token")
	}
	_ = pool.QueryRow(ctx, `SELECT COUNT(*) FROM subscribers WHERE email=$1`, email).Scan(&cnt)
	if cnt != 0 {
		t.Fatalf("subscriber not removed after unsubscribe")
	}
	if _, ok := m.unsubscribe(ctx, token); ok {
		t.Error("a spent unsubscribe token was accepted again")
	}
}

// TestFetchFeedIntegration checks the RSS query against a real DB (schema from
// migrations). Skipped unless SHANRAQ_TEST_DB is set.
func TestFetchFeedIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the RSS integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	authorID := uuid.New()
	articleID := uuid.New()
	slug := "rss-itest-" + articleID.String()[:8]
	_, _ = pool.Exec(ctx, `INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`, authorID, "rss-"+authorID.String()+"@t.test")
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, authorID) })
	if _, err := pool.Exec(ctx, `INSERT INTO articles (id, author_id, slug, original_lang, status, published_at) VALUES ($1,$2,$3,'ru','published',NOW())`, articleID, authorID, slug); err != nil {
		t.Fatalf("insert article: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES ($1,'ru','РСС Тест','Аннотация','Тело','human','ready')`, articleID); err != nil {
		t.Fatalf("insert translation: %v", err)
	}

	m := &Module{db: pool, baseURL: "https://shanraq.org", log: zap.NewNop()}
	entries, err := m.fetchFeed(ctx, "ru", 30)
	if err != nil {
		t.Fatalf("fetchFeed: %v", err)
	}
	var found bool
	for _, e := range entries {
		if e.Slug == slug {
			found = true
			if e.Title != "РСС Тест" {
				t.Errorf("title = %q", e.Title)
			}
		}
	}
	if !found {
		t.Fatalf("published article %s not present in feed", slug)
	}
}

// The landing pages render without the site header, so the mark and the copy on
// them are the whole page — a template slip would ship a blank card.
func TestRenderNoticePages(t *testing.T) {
	m := &Module{log: zap.NewNop()}

	t.Run("unsubscribe asks before removing", func(t *testing.T) {
		w := httptest.NewRecorder()
		m.renderNotice(w, noticePage{
			Lang:         "ru",
			Title:        ds("ru", "unsub_ask_title"),
			Lead:         ds("ru", "unsub_ask_lead"),
			Email:        maskEmail("baimurza.daulet@gmail.com"),
			Points:       []string{ds("ru", "unsub_ask_p1"), ds("ru", "unsub_ask_p2"), ds("ru", "unsub_ask_p3")},
			ConfirmPost:  "/unsubscribe?token=abc",
			ConfirmLabel: ds("ru", "unsub_ask_btn"),
			CancelHref:   "/?lang=ru",
			CancelLabel:  ds("ru", "unsub_keep_btn"),
			Foot:         ds("ru", "unsub_ask_foot"),
		})
		body := w.Body.String()
		for _, want := range []string{
			`class="notice__brand"`,
			`/static/brand/shanraq.svg`,       // light-theme mark
			`/static/brand/shanraq-light.svg`, // dark-theme mark
			`Shanraq.org`,
			"Отписаться от рассылки?",
			"bai***@gmail.com",
			"Мы не передаём адрес третьим лицам",
			`method="post" action="/unsubscribe?token=abc"`,
			"Оставить подписку",
			`content="noindex, nofollow"`,
		} {
			if !strings.Contains(body, want) {
				t.Errorf("missing %q", want)
			}
		}
		if w.Header().Get("X-Robots-Tag") != "noindex" {
			t.Error("token page is missing X-Robots-Tag: noindex")
		}
	})

	t.Run("confirmed page offers no destructive button", func(t *testing.T) {
		w := httptest.NewRecorder()
		m.renderNotice(w, noticePage{
			Lang:        "kz",
			Title:       ds("kz", "confirmed_title"),
			Lead:        ds("kz", "confirmed_lead"),
			Points:      []string{ds("kz", "confirmed_p1")},
			CancelHref:  "/?lang=kz",
			CancelLabel: ds("kz", "to_site"),
		})
		body := w.Body.String()
		if strings.Contains(body, "<form") {
			t.Error("confirmation page should not carry a form")
		}
		if !strings.Contains(body, `<html lang="kk"`) {
			t.Error("kz must render as BCP 47 kk")
		}
		if !strings.Contains(body, "Жазылым расталды") {
			t.Error("missing Kazakh title")
		}
	})

	t.Run("unknown language falls back", func(t *testing.T) {
		w := httptest.NewRecorder()
		m.renderNotice(w, noticePage{Lang: "zz", Title: "x", CancelHref: "/", CancelLabel: "y"})
		if !strings.Contains(w.Body.String(), `<html lang="ru"`) {
			t.Error("unknown lang should fall back to ru")
		}
	})
}

// A malformed key must be refused rather than sent and silently rejected by the
// endpoint, and an absent one must read as "disabled", not as an error.
func TestNormalizeIndexNowKey(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantOK  bool
		comment string
	}{
		{"", "", true, "unset means disabled, not broken"},
		{"   ", "", true, "whitespace only is also unset"},
		{"a1b2c3d4", "a1b2c3d4", true, "minimum length"},
		{"  a1b2c3d4  ", "a1b2c3d4", true, "trimmed"},
		{"A1B2C3D4-e5f6", "A1B2C3D4-e5f6", true, "uppercase hex and dashes allowed"},
		{"short", "", false, "under eight characters"},
		{"zzzzzzzz", "", false, "not hex"},
		{"a1b2c3d4 e5", "", false, "spaces inside"},
		{strings.Repeat("a", 129), "", false, "over the limit"},
	}
	for _, c := range cases {
		got, ok := normalizeIndexNowKey(c.in)
		if got != c.want || ok != c.wantOK {
			t.Errorf("normalizeIndexNowKey(%q) = %q,%v — want %q,%v (%s)", c.in, got, ok, c.want, c.wantOK, c.comment)
		}
	}
}

// The key file proves domain ownership; without it every submission is refused.
func TestIndexNowKeyEndpoint(t *testing.T) {
	m := &Module{log: zap.NewNop(), indexNowKey: "a1b2c3d4e5f6"}
	w := httptest.NewRecorder()
	m.handleIndexNowKey(w, httptest.NewRequest(http.MethodGet, indexNowKeyPath, nil))
	if w.Code != http.StatusOK || w.Body.String() != "a1b2c3d4e5f6" {
		t.Errorf("key file = %d %q", w.Code, w.Body.String())
	}

	off := &Module{log: zap.NewNop()}
	w = httptest.NewRecorder()
	off.handleIndexNowKey(w, httptest.NewRequest(http.MethodGet, indexNowKeyPath, nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("unconfigured key file = %d, want 404", w.Code)
	}
}

// All three languages go out per publish: submitting only Russian would leave
// two thirds of the catalogue waiting to be crawled.
func TestIndexNowCoversEveryLanguage(t *testing.T) {
	if len(rssLangOrder) != len(rssLangs) {
		t.Fatalf("rssLangOrder has %d entries, rssLangs %d", len(rssLangOrder), len(rssLangs))
	}
	for _, l := range rssLangOrder {
		if !rssLangs[l] {
			t.Errorf("%q is not a feed language", l)
		}
	}
}
