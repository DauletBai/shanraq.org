package syndicate

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func newSynFixture(t *testing.T) (*Module, *pgxpool.Pool, context.Context) {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the syndication handler tests")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return &Module{db: pool, baseURL: "https://shanraq.org", log: zap.NewNop()}, pool, ctx
}

// A mail client that follows List-Unsubscribe POSTs by itself, and a reader who
// clicks the link gets a page first. Both paths run on the same 24-byte token,
// which is the only authorisation either of them carries.
func TestUnsubscribeAsksBeforeItActs(t *testing.T) {
	m, pool, ctx := newSynFixture(t)
	email := "unsub-" + uuid.NewString()[:8] + "@t.test"
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM subscribers WHERE email=$1`, email) })

	if _, err := m.subscribe(ctx, email, "ru"); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	var token string
	if err := pool.QueryRow(ctx,
		`SELECT unsubscribe_token FROM subscribers WHERE email=$1`, email).Scan(&token); err != nil {
		t.Fatalf("read token: %v", err)
	}

	// GET asks. It must not remove anybody: mail clients and link scanners
	// fetch every URL in a message, and a GET that unsubscribed would empty the
	// list by itself.
	rec := httptest.NewRecorder()
	m.handleUnsubscribePage(rec, httptest.NewRequest(http.MethodGet, "/unsubscribe?token="+token, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("the confirmation page answered %d", rec.Code)
	}
	var still int
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM subscribers WHERE email=$1`, email).Scan(&still)
	if still != 1 {
		t.Fatal("merely opening the unsubscribe link removed the subscriber")
	}
	// The page shows the address masked, so a forwarded link does not disclose
	// who was subscribed.
	if body := rec.Body.String(); strings.Contains(body, email) {
		t.Error("the confirmation page printed the full address")
	}

	// POST performs.
	rec = httptest.NewRecorder()
	m.handleUnsubscribe(rec, httptest.NewRequest(http.MethodPost, "/unsubscribe?token="+token, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("unsubscribe answered %d", rec.Code)
	}
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM subscribers WHERE email=$1`, email).Scan(&still)
	if still != 0 {
		t.Error("the subscriber survived their own unsubscribe")
	}
}

// A token nobody issued removes nobody, and says so without an error page that
// would tell a prober they had guessed wrong in an interesting way.
func TestAnUnknownUnsubscribeTokenIsRefusedQuietly(t *testing.T) {
	m, _, _ := newSynFixture(t)
	for _, token := range []string{"", "not-a-token", strings.Repeat("a", 48)} {
		rec := httptest.NewRecorder()
		m.handleUnsubscribe(rec, httptest.NewRequest(http.MethodPost, "/unsubscribe?token="+url.QueryEscape(token), nil))
		if rec.Code != http.StatusOK {
			t.Errorf("token %q gave status %d; the notice page is the answer", token, rec.Code)
		}
	}
}

// Double opt-in: a subscription counts only once the address has confirmed it,
// which is what makes the list lawful to send to.
func TestASubscriptionCountsOnlyAfterItIsConfirmed(t *testing.T) {
	m, pool, ctx := newSynFixture(t)
	email := "conf-" + uuid.NewString()[:8] + "@t.test"
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM subscribers WHERE email=$1`, email) })

	confirmTok, err := m.subscribe(ctx, email, "kz")
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	var confirmedAt *string
	_ = pool.QueryRow(ctx, `SELECT confirmed_at::text FROM subscribers WHERE email=$1`, email).Scan(&confirmedAt)
	if confirmedAt != nil {
		t.Fatal("a subscription was confirmed before anyone confirmed it")
	}

	rec := httptest.NewRecorder()
	m.handleConfirm(rec, httptest.NewRequest(http.MethodGet, "/subscribe/confirm?token="+confirmTok, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("confirm answered %d", rec.Code)
	}
	_ = pool.QueryRow(ctx, `SELECT confirmed_at::text FROM subscribers WHERE email=$1`, email).Scan(&confirmedAt)
	if confirmedAt == nil {
		t.Error("the address confirmed and the row still says otherwise")
	}

	// The confirmation link is spent: a second visit cannot re-confirm, so a
	// leaked link is not a standing key to the list.
	rec = httptest.NewRecorder()
	m.handleConfirm(rec, httptest.NewRequest(http.MethodGet, "/subscribe/confirm?token="+confirmTok, nil))
	if rec.Code != http.StatusOK {
		t.Errorf("a spent confirmation gave status %d; the notice page is the answer", rec.Code)
	}
}

// The feed is how other people's readers find us, so it has to be XML an
// aggregator will accept, in the language it was asked for.
func TestTheFeedIsWellFormedInEveryLanguage(t *testing.T) {
	m, _, _ := newSynFixture(t)
	for _, lang := range []string{"ru", "kz", "en", "", "nonsense"} {
		rec := httptest.NewRecorder()
		m.handleRSS(rec, httptest.NewRequest(http.MethodGet, "/feed.xml?lang="+lang, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("lang %q: status %d", lang, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/rss+xml") {
			t.Errorf("lang %q: content type %q, which no reader will treat as a feed", lang, ct)
		}
		var doc struct {
			XMLName xml.Name `xml:"rss"`
			Version string   `xml:"version,attr"`
			Channel struct {
				Title string `xml:"title"`
				Link  string `xml:"link"`
				Items []struct {
					Title string `xml:"title"`
					Link  string `xml:"link"`
					GUID  string `xml:"guid"`
				} `xml:"item"`
			} `xml:"channel"`
		}
		if err := xml.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
			t.Fatalf("lang %q: the feed is not parseable XML: %v", lang, err)
		}
		if doc.Version != "2.0" || doc.Channel.Title == "" {
			t.Errorf("lang %q: rss version %q, channel title %q", lang, doc.Version, doc.Channel.Title)
		}
		for _, it := range doc.Channel.Items {
			if !strings.HasPrefix(it.Link, "https://shanraq.org/") {
				t.Errorf("lang %q: item links to %q, which is not this site", lang, it.Link)
			}
			if it.GUID == "" || it.Title == "" {
				t.Errorf("lang %q: an item has no guid or no title", lang)
			}
		}
	}
}
