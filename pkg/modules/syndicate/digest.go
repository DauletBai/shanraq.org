package syndicate

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// digestInterval is how often the weekly digest may go out.
const digestInterval = 7 * 24 * time.Hour

// subscriber is one newsletter recipient.
type subscriber struct {
	Email string
	Lang  string
	Token string
}

// digestStrings holds localized email + unsubscribe-page copy. Kept local so
// syndicate needn't import the articles package (which would cycle).
var digestStrings = map[string]map[string]string{
	"subject": {
		"kz": "Shanraq: апталық шолу",
		"ru": "Shanraq: обзор недели",
		"en": "Shanraq: weekly digest",
	},
	"intro": {
		"kz": "Осы аптадағы жаңа жарияланымдар:",
		"ru": "Новые публикации этой недели:",
		"en": "New stories this week:",
	},
	"unsub": {
		"kz": "Жазылудан бас тарту",
		"ru": "Отписаться от рассылки",
		"en": "Unsubscribe",
	},
	"unsub_done": {
		"kz": "Сіз жазылудан бас тарттыңыз. Қош келдіңіз кез келген уақытта!",
		"ru": "Вы отписались от рассылки. Возвращайтесь в любое время!",
		"en": "You have unsubscribed. Come back anytime!",
	},
}

func ds(lang, key string) string {
	if m, ok := digestStrings[key]; ok {
		if v, ok := m[lang]; ok && v != "" {
			return v
		}
		return m["ru"]
	}
	return key
}

// ---------- subscription HTTP ----------

func (m *Module) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	lang := r.FormValue("lang")
	if !rssLangs[lang] {
		lang = "ru"
	}
	back := "/?lang=" + lang
	if email == "" || !strings.Contains(email, "@") {
		http.Redirect(w, r, back, http.StatusSeeOther)
		return
	}
	if err := m.subscribe(r.Context(), email, lang); err != nil {
		m.log.Warn("subscribe failed", zap.Error(err))
	}
	http.Redirect(w, r, back+"&sub=ok", http.StatusSeeOther)
}

func (m *Module) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	lang, _ := m.unsubscribe(r.Context(), token)
	if !rssLangs[lang] {
		lang = "ru"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!doctype html><html lang="%s"><head><meta charset="utf-8">`+
		`<meta name="viewport" content="width=device-width,initial-scale=1">`+
		`<link rel="stylesheet" href="/static/css/shanraq.css"><title>Shanraq</title></head>`+
		`<body><main class="container"><div class="auth-wrap"><h1>Shanraq</h1>`+
		`<p class="sub">%s</p><a class="btn btn--primary" href="/?lang=%s" style="width:100%%">→</a></div></main></body></html>`,
		lang, ds(lang, "unsub_done"), lang)
}

// ---------- subscriber store ----------

func (m *Module) subscribe(ctx context.Context, email, lang string) error {
	token, err := randomToken()
	if err != nil {
		return err
	}
	_, err = m.db.Exec(ctx, `
		INSERT INTO subscribers (email, lang, unsubscribe_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET lang = EXCLUDED.lang
	`, email, lang, token)
	return err
}

func (m *Module) unsubscribe(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "ru", nil
	}
	var lang string
	err := m.db.QueryRow(ctx, `DELETE FROM subscribers WHERE unsubscribe_token = $1 RETURNING lang`, token).Scan(&lang)
	if err != nil {
		return "ru", nil // unknown token: still show a friendly page
	}
	return lang, nil
}

func (m *Module) listSubscribers(ctx context.Context) ([]subscriber, error) {
	rows, err := m.db.Query(ctx, `SELECT email, lang, unsubscribe_token FROM subscribers`)
	if err != nil {
		return nil, fmt.Errorf("list subscribers: %w", err)
	}
	defer rows.Close()
	var subs []subscriber
	for rows.Next() {
		var s subscriber
		if err := rows.Scan(&s.Email, &s.Lang, &s.Token); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

// ---------- digest build + send ----------

// fetchRecent returns articles published within the last 7 days, resolved to lang.
func (m *Module) fetchRecent(ctx context.Context, lang string, limit int) ([]feedEntry, error) {
	rows, err := m.db.Query(ctx, `
		SELECT a.slug,
		       COALESCE(NULLIF(tl.title, ''), torig.title),
		       COALESCE(NULLIF(tl.summary, ''), torig.summary),
		       CASE WHEN tl.title IS NOT NULL AND tl.title <> '' THEN $1 ELSE a.original_lang END,
		       COALESCE(a.published_at, a.updated_at)
		FROM articles a
		JOIN article_translations torig ON torig.article_id = a.id AND torig.lang = a.original_lang
		LEFT JOIN article_translations tl ON tl.article_id = a.id AND tl.lang = $1 AND tl.title <> '' AND tl.body_md <> ''
		WHERE a.status = 'published' AND a.published_at >= NOW() - INTERVAL '7 days'
		ORDER BY a.published_at DESC NULLS LAST
		LIMIT $2
	`, lang, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch recent: %w", err)
	}
	defer rows.Close()
	var entries []feedEntry
	for rows.Next() {
		var e feedEntry
		if err := rows.Scan(&e.Slug, &e.Title, &e.Summary, &e.Lang, &e.Modified); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// renderDigest builds the plain-text email body for one subscriber.
func (m *Module) renderDigest(lang string, entries []feedEntry, token string) (subject, body string) {
	var b strings.Builder
	b.WriteString(ds(lang, "intro"))
	b.WriteString("\n\n")
	for _, e := range entries {
		b.WriteString("• ")
		b.WriteString(strings.TrimSpace(e.Title))
		b.WriteString("\n  ")
		b.WriteString(m.articleURL(e.Slug, e.Lang))
		b.WriteString("\n\n")
	}
	b.WriteString("—\n")
	b.WriteString(ds(lang, "unsub"))
	b.WriteString(": ")
	b.WriteString(m.baseURL + "/unsubscribe?token=" + token)
	return ds(lang, "subject"), b.String()
}

// SendDigest emails the weekly digest to every subscriber in their language.
// Returns how many messages were sent. A no-op when email is not configured.
func (m *Module) SendDigest(ctx context.Context) (int, error) {
	if !m.emailEnabled || m.mailer == nil {
		return 0, nil
	}
	subs, err := m.listSubscribers(ctx)
	if err != nil {
		return 0, err
	}
	cache := map[string][]feedEntry{}
	sent := 0
	for _, s := range subs {
		entries, ok := cache[s.Lang]
		if !ok {
			entries, err = m.fetchRecent(ctx, s.Lang, 15)
			if err != nil {
				return sent, err
			}
			cache[s.Lang] = entries
		}
		if len(entries) == 0 {
			continue
		}
		subject, bodyText := m.renderDigest(s.Lang, entries, s.Token)
		if err := m.mailer.Send(ctx, s.Email, subject, bodyText); err != nil {
			m.log.Warn("digest send failed", zap.String("to", s.Email), zap.Error(err))
			continue
		}
		sent++
	}
	return sent, nil
}

func (m *Module) digestDue(ctx context.Context) bool {
	var last *time.Time
	if err := m.db.QueryRow(ctx, `SELECT last_sent_at FROM digest_state WHERE id = 1`).Scan(&last); err != nil {
		return false
	}
	return last == nil || time.Since(*last) >= digestInterval
}

func (m *Module) markDigestSent(ctx context.Context) {
	if _, err := m.db.Exec(ctx, `UPDATE digest_state SET last_sent_at = NOW() WHERE id = 1`); err != nil {
		m.log.Warn("mark digest sent", zap.Error(err))
	}
}

func randomToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
