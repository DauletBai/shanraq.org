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

// subscriber is one confirmed newsletter recipient.
type subscriber struct {
	Email string
	Lang  string
	Token string
}

// digestStrings holds localized email + landing-page copy. Kept local so
// syndicate needn't import the articles package (which would cycle).
var digestStrings = map[string]map[string]string{
	"subject": {
		"kz": "Shanraq.org: апталық шолу",
		"ru": "Shanraq.org: обзор недели",
		"en": "Shanraq.org: weekly digest",
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

	// ---- confirmation email (double opt-in) ----
	"confirm_subject": {
		"kz": "Shanraq.org: жазылымды растаңыз",
		"ru": "Shanraq.org: подтвердите подписку",
		"en": "Shanraq.org: please confirm your subscription",
	},
	"confirm_body": {
		"kz": "Сіз Shanraq.org апталық шолуына жазылдыңыз.\n\n" +
			"Жазылымды растаңыз — оған дейін бірде-бір хат жібермейміз:\n%s\n\n" +
			"Егер бұл сіз болмасаңыз, сілтемеге өтпей-ақ қойыңыз: растаусыз\n" +
			"мекенжай таратылымға түспейді.\n\n—\nShanraq.org · %s",
		"ru": "Вы подписались на еженедельный обзор Shanraq.org.\n\n" +
			"Подтвердите подписку — до этого мы не отправим ни одного письма:\n%s\n\n" +
			"Если это были не вы, просто не переходите по ссылке: без подтверждения\n" +
			"адрес в рассылку не попадёт.\n\n—\nShanraq.org · %s",
		"en": "You have signed up for the Shanraq.org weekly digest.\n\n" +
			"Please confirm — until you do, we will not send you a single email:\n%s\n\n" +
			"If this was not you, simply ignore this message: without confirmation\n" +
			"the address never enters the list.\n\n—\nShanraq.org · %s",
	},

	// ---- "really unsubscribe?" page ----
	"unsub_ask_title": {
		"kz": "Жазылудан бас тартасыз ба?",
		"ru": "Отписаться от рассылки?",
		"en": "Unsubscribe from the newsletter?",
	},
	"unsub_ask_lead": {
		"kz": "Хаттарды бұдан әрі жіберудің қажеті жоқ екенін растаңыз.",
		"ru": "Подтвердите, что письма больше присылать не нужно.",
		"en": "Please confirm that you no longer want these emails.",
	},
	"unsub_ask_p1": {
		"kz": "Апталық шолу — аптасына бір хат, жаңа материалдарға сілтемелермен.",
		"ru": "Обзор недели — одно письмо в неделю со ссылками на новые материалы.",
		"en": "The weekly digest is one email a week, linking to the new stories.",
	},
	"unsub_ask_p2": {
		"kz": "Мекенжайыңызды үшінші тұлғаларға бермейміз және жарнама жібермейміз.",
		"ru": "Мы не передаём адрес третьим лицам и не присылаем рекламу.",
		"en": "We never pass your address to anyone and never send advertising.",
	},
	"unsub_ask_p3": {
		"kz": "Ойыңыз өзгерсе — сайттан кез келген уақытта қайта жазыла аласыз.",
		"ru": "Передумаете — подписаться снова можно на сайте в любой момент.",
		"en": "Change your mind and you can subscribe again on the site at any time.",
	},
	"unsub_ask_btn": {
		"kz": "Бас тарту",
		"ru": "Отписаться",
		"en": "Unsubscribe",
	},
	"unsub_keep_btn": {
		"kz": "Жазылымды қалдыру",
		"ru": "Оставить подписку",
		"en": "Keep my subscription",
	},
	"unsub_ask_foot": {
		"kz": "Бұл бет хаттағы сілтеме арқылы ашылды. Түймені баспағаныңызша жазылым белсенді күйде қалады.",
		"ru": "Эта страница открыта по ссылке из письма. Пока вы не нажали кнопку, подписка остаётся активной.",
		"en": "You reached this page from a link in an email. Nothing has changed yet — your subscription stays active until you press the button.",
	},

	// ---- "you have unsubscribed" page ----
	"unsub_done_title": {
		"kz": "Сіз жазылудан бас тарттыңыз",
		"ru": "Вы отписались",
		"en": "You have unsubscribed",
	},
	"unsub_done_lead": {
		"kz": "Мекенжай тізімнен жойылды — хаттар бұдан әрі келмейді. Оқығаныңызға рахмет.",
		"ru": "Адрес удалён из списка — писем больше не будет. Спасибо, что читали.",
		"en": "Your address has been removed from the list. Thank you for reading.",
	},
	"unsub_done_p1": {
		"kz": "Жаңа материалдар сайтта бұрынғысынша шығып тұрады.",
		"ru": "Новые материалы по-прежнему выходят на сайте.",
		"en": "New stories keep going up on the site.",
	},
	"unsub_done_p2": {
		"kz": "Сайт қолжетімсіз болса, мақалаларды Telegram-нан оқуға немесе RSS арқылы алуға болады.",
		"ru": "Если сайт окажется недоступен, статьи можно читать в Telegram или получать по RSS.",
		"en": "If the site is ever unreachable, the articles are also on Telegram and in the RSS feed.",
	},
	"unsub_done_p3": {
		"kz": "Жазылымға кез келген сәтте орала аласыз — форма сайттың оң бағанында.",
		"ru": "Вернуться в рассылку можно в любой момент — форма в правой колонке сайта.",
		"en": "You can rejoin at any time — the form is in the site's right-hand column.",
	},

	// ---- "subscription confirmed" page ----
	"confirmed_title": {
		"kz": "Жазылым расталды",
		"ru": "Подписка подтверждена",
		"en": "Subscription confirmed",
	},
	"confirmed_lead": {
		"kz": "Дайын. Ең жақын апталық шолу осы мекенжайға келеді.",
		"ru": "Готово. Ближайший обзор недели придёт на этот адрес.",
		"en": "Done. The next weekly digest will arrive at this address.",
	},
	"confirmed_p1": {
		"kz": "Аптасына бір хат — жаңа талдаулар сілтемелерімен, жарнамасыз.",
		"ru": "Одно письмо в неделю — новые разборы со ссылками, без рекламы.",
		"en": "One email a week — new analysis with links, no advertising.",
	},
	"confirmed_p2": {
		"kz": "Мекенжай үшінші тұлғаларға берілмейді.",
		"ru": "Адрес не передаётся третьим лицам.",
		"en": "Your address is never passed to third parties.",
	},
	"confirmed_p3": {
		"kz": "Бас тарту бір басумен — сілтеме әр хатта бар.",
		"ru": "Отписаться можно одним нажатием — ссылка есть в каждом письме.",
		"en": "One click to unsubscribe — the link is in every email.",
	},

	// ---- unknown / spent token ----
	"bad_title": {
		"kz": "Сілтеме жарамсыз",
		"ru": "Ссылка недействительна",
		"en": "This link is no longer valid",
	},
	"bad_lead": {
		"kz": "Мүмкін, сіз бұрын бас тартқансыз немесе сілтеме ескірген. Ештеңе істеудің қажеті жоқ.",
		"ru": "Возможно, вы уже отписались или ссылка устарела. Ничего делать не нужно.",
		"en": "You may have unsubscribed already, or the link has expired. There is nothing you need to do.",
	},
	"to_site": {
		"kz": "Сайтқа",
		"ru": "На сайт",
		"en": "Go to the site",
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

// htmlLang maps our language codes to BCP 47 tags for the <html lang> attribute.
func htmlLang(lang string) string {
	if lang == "kz" {
		return "kk"
	}
	return lang
}

// ---------- subscription HTTP ----------

// handleSubscribe records a subscription *request* and emails a confirmation
// link. Nothing is ever sent to an address that has not clicked that link, so
// signing up somebody else's address achieves nothing, and the person who did
// sign up gets an immediate reply instead of up to a week of silence.
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
	back := safeBack(r.FormValue("back"), lang)
	if !plausibleEmail(email) {
		http.Redirect(w, r, back+"&subscribed=bad", http.StatusSeeOther)
		return
	}
	token, err := m.subscribe(r.Context(), email, lang)
	if err != nil {
		// Never claim success we did not achieve: the old handler logged the
		// error and still showed "thanks, you're subscribed".
		m.log.Warn("subscribe failed", zap.Error(err))
		http.Redirect(w, r, back+"&subscribed=err", http.StatusSeeOther)
		return
	}
	if err := m.sendConfirmation(r.Context(), email, lang, token); err != nil {
		m.log.Warn("confirmation email failed", zap.String("to", email), zap.Error(err))
		http.Redirect(w, r, back+"&subscribed=err", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, back+"&subscribed=pending", http.StatusSeeOther)
}

// handleConfirm activates a pending subscription.
func (m *Module) handleConfirm(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	lang, ok := m.confirmSubscriber(r.Context(), token)
	if !ok {
		m.renderNotice(w, badTokenNotice(lang))
		return
	}
	m.renderNotice(w, noticePage{
		Lang:        lang,
		Title:       ds(lang, "confirmed_title"),
		Lead:        ds(lang, "confirmed_lead"),
		Points:      []string{ds(lang, "confirmed_p1"), ds(lang, "confirmed_p2"), ds(lang, "confirmed_p3")},
		CancelHref:  "/?lang=" + lang,
		CancelLabel: ds(lang, "to_site"),
	})
}

// handleUnsubscribePage answers the GET behind the link in every digest. It
// only *asks*. Gmail, corporate mail gateways and link prefetchers fetch every
// URL in a message, so a GET that deleted the row would unsubscribe readers who
// never clicked anything — the removal lives on POST below.
func (m *Module) handleUnsubscribePage(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	email, lang, ok := m.lookupSubscriber(r.Context(), token)
	if !ok {
		m.renderNotice(w, badTokenNotice(lang))
		return
	}
	m.renderNotice(w, noticePage{
		Lang:         lang,
		Title:        ds(lang, "unsub_ask_title"),
		Lead:         ds(lang, "unsub_ask_lead"),
		Email:        maskEmail(email),
		Points:       []string{ds(lang, "unsub_ask_p1"), ds(lang, "unsub_ask_p2"), ds(lang, "unsub_ask_p3")},
		ConfirmPost:  "/unsubscribe?token=" + token,
		ConfirmLabel: ds(lang, "unsub_ask_btn"),
		CancelHref:   "/?lang=" + lang,
		CancelLabel:  ds(lang, "unsub_keep_btn"),
		Foot:         ds(lang, "unsub_ask_foot"),
	})
}

// handleUnsubscribe performs the removal. It serves both the button on the page
// above and RFC 8058 one-click unsubscribe, where the mail client POSTs to the
// List-Unsubscribe URL by itself — which is why it must not be CSRF-guarded:
// the 24-byte token in the URL is the authorization.
func (m *Module) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	lang, ok := m.unsubscribe(r.Context(), token)
	if !ok {
		m.renderNotice(w, badTokenNotice(lang))
		return
	}
	m.renderNotice(w, noticePage{
		Lang:        lang,
		Title:       ds(lang, "unsub_done_title"),
		Lead:        ds(lang, "unsub_done_lead"),
		Points:      []string{ds(lang, "unsub_done_p1"), ds(lang, "unsub_done_p2"), ds(lang, "unsub_done_p3")},
		CancelHref:  "/?lang=" + lang,
		CancelLabel: ds(lang, "to_site"),
	})
}

func badTokenNotice(lang string) noticePage {
	return noticePage{
		Lang:        lang,
		Title:       ds(lang, "bad_title"),
		Lead:        ds(lang, "bad_lead"),
		CancelHref:  "/?lang=" + lang,
		CancelLabel: ds(lang, "to_site"),
	}
}

// ---------- notice pages ----------

// noticePage is the context for templates/notice.html.
type noticePage struct {
	Lang         string
	HTMLLang     string
	Title        string
	Lead         string
	Email        string   // masked, shown only when we know whose link it is
	Points       []string // what the reader gets / loses, in plain words
	ConfirmPost  string   // POST target for a destructive action ("" = none)
	ConfirmLabel string
	CancelHref   string
	CancelLabel  string
	Foot         string
}

func (m *Module) renderNotice(w http.ResponseWriter, p noticePage) {
	if !rssLangs[p.Lang] {
		p.Lang = "ru"
	}
	p.HTMLLang = htmlLang(p.Lang)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Robots-Tag", "noindex")
	// ExecuteTemplate by file name, not Execute: ParseFS names each template
	// after its file, so the root "notice" template is empty and Execute on it
	// would render a blank page.
	if err := noticeTmpl.ExecuteTemplate(w, "notice.html", p); err != nil {
		m.log.Warn("notice render", zap.Error(err))
	}
}

// ---------- subscriber store ----------

// subscribe upserts a pending subscription and returns the confirmation token.
// An address that is already confirmed keeps its confirmation: re-submitting
// the form must not silently drop somebody off the list.
func (m *Module) subscribe(ctx context.Context, email, lang string) (string, error) {
	confirm, err := randomToken()
	if err != nil {
		return "", err
	}
	unsub, err := randomToken()
	if err != nil {
		return "", err
	}
	_, err = m.db.Exec(ctx, `
		INSERT INTO subscribers (email, lang, unsubscribe_token, confirm_token)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE
		   SET lang = EXCLUDED.lang, confirm_token = EXCLUDED.confirm_token
	`, email, lang, unsub, confirm)
	if err != nil {
		return "", err
	}
	return confirm, nil
}

// confirmSubscriber activates a pending row. Confirming twice is harmless: the
// second click finds no token and reports the link as spent, which is also what
// a stranger guessing tokens would see.
func (m *Module) confirmSubscriber(ctx context.Context, token string) (string, bool) {
	if token == "" {
		return "ru", false
	}
	var lang string
	err := m.db.QueryRow(ctx, `
		UPDATE subscribers
		   SET confirmed_at = COALESCE(confirmed_at, NOW()), confirm_token = NULL
		 WHERE confirm_token = $1
		RETURNING lang
	`, token).Scan(&lang)
	if err != nil {
		return "ru", false
	}
	return lang, true
}

// lookupSubscriber resolves an unsubscribe token without changing anything.
func (m *Module) lookupSubscriber(ctx context.Context, token string) (email, lang string, ok bool) {
	if token == "" {
		return "", "ru", false
	}
	err := m.db.QueryRow(ctx,
		`SELECT email, lang FROM subscribers WHERE unsubscribe_token = $1`, token).Scan(&email, &lang)
	if err != nil {
		return "", "ru", false
	}
	return email, lang, true
}

func (m *Module) unsubscribe(ctx context.Context, token string) (string, bool) {
	if token == "" {
		return "ru", false
	}
	var lang string
	err := m.db.QueryRow(ctx,
		`DELETE FROM subscribers WHERE unsubscribe_token = $1 RETURNING lang`, token).Scan(&lang)
	if err != nil {
		return "ru", false
	}
	return lang, true
}

// listSubscribers returns confirmed recipients only. Pending rows are requests,
// not subscribers, and must never receive mail.
func (m *Module) listSubscribers(ctx context.Context) ([]subscriber, error) {
	rows, err := m.db.Query(ctx,
		`SELECT email, lang, unsubscribe_token FROM subscribers WHERE confirmed_at IS NOT NULL`)
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
	b.WriteString(m.unsubURL(token))
	return ds(lang, "subject"), b.String()
}

// sendConfirmation emails the double opt-in link. Without SMTP there is nobody
// to confirm to, so the request is reported as failed rather than left pending
// forever.
func (m *Module) sendConfirmation(ctx context.Context, email, lang, token string) error {
	if !m.emailEnabled || m.mailer == nil {
		return fmt.Errorf("email not configured")
	}
	link := m.baseURL + "/subscribe/confirm?token=" + token
	body := fmt.Sprintf(ds(lang, "confirm_body"), link, m.baseURL)
	return m.mailer.Send(ctx, email, ds(lang, "confirm_subject"), body)
}

// SendDigest emails the weekly digest to every confirmed subscriber in their
// language. Returns how many messages were sent. A no-op without SMTP.
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
		if err := m.sendBulk(ctx, s.Email, subject, bodyText, s.Token); err != nil {
			m.log.Warn("digest send failed", zap.String("to", s.Email), zap.Error(err))
			continue
		}
		sent++
	}
	return sent, nil
}

// headerMailer is the optional Mailer extension for senders that can add
// headers. Type-asserted rather than required, so a plain Mailer still works.
type headerMailer interface {
	SendWithHeaders(ctx context.Context, to, subject, body string, headers map[string]string) error
}

// sendBulk sends one digest with the headers bulk mail is expected to carry.
// Gmail and Yahoo have required one-click unsubscribe from bulk senders since
// February 2024; without these headers the digest is scored as spam far more
// often, which for a newsletter is the same as not sending it.
func (m *Module) sendBulk(ctx context.Context, to, subject, body, token string) error {
	hm, ok := m.mailer.(headerMailer)
	if !ok {
		return m.mailer.Send(ctx, to, subject, body)
	}
	return hm.SendWithHeaders(ctx, to, subject, body, map[string]string{
		"List-Unsubscribe":      "<" + m.unsubURL(token) + ">",
		"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
		"Precedence":            "bulk",
		"Auto-Submitted":        "auto-generated",
	})
}

func (m *Module) unsubURL(token string) string {
	return m.baseURL + "/unsubscribe?token=" + token
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

// ---------- small helpers ----------

// plausibleEmail is a shape check, not validation: the confirmation link is
// what actually proves the address exists and belongs to the sender.
func plausibleEmail(email string) bool {
	at := strings.IndexByte(email, '@')
	if at <= 0 || at == len(email)-1 || strings.Count(email, "@") != 1 {
		return false
	}
	domain := email[at+1:]
	return strings.Contains(domain, ".") && !strings.HasPrefix(domain, ".") &&
		!strings.HasSuffix(domain, ".") && !strings.ContainsAny(email, " \t\r\n")
}

// maskEmail shortens an address for display on the unsubscribe page. The reader
// arrived with a token from their own inbox, but a forwarded link should not
// hand a third party the full address.
func maskEmail(email string) string {
	at := strings.IndexByte(email, '@')
	if at <= 0 {
		return ""
	}
	local, domain := email[:at], email[at:]
	keep := 3
	if len(local) < keep {
		keep = 1
	}
	return local[:keep] + "***" + domain
}

// safeBack keeps the reader on the page they subscribed from. Only same-site
// paths are honoured, so the form cannot be turned into an open redirect.
func safeBack(raw, lang string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || !strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") {
		return "/?lang=" + lang
	}
	if i := strings.IndexAny(raw, "?#"); i >= 0 {
		raw = raw[:i]
	}
	return raw + "?lang=" + lang
}
