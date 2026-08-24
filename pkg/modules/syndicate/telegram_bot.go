package syndicate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// The bot: one subscription point, a different feed for each subscriber.
//
// A Telegram channel cannot show two people different things — it broadcasts
// one message and knows nothing about who reads it. So the channel carries what
// a guest sees on the site, material written for everyone, and everything
// addressed to a place is delivered here, one message at a time, under the same
// rule the site's feed uses: a piece written for a place reaches the people who
// are inside that place, and nobody else.
//
// The person says where they are once. We never infer it — not from the phone,
// not from the language, not from an address book. A platform that promises no
// tracking of readers does not get to make an exception for the convenience of
// its own onboarding.

// tgPollTimeout is the long-poll window. Telegram holds the request open until
// something happens or this expires, so a large value costs nothing and keeps
// the bot responsive without hammering the API.
const tgPollTimeout = 50

// tgSendGap paces the fan-out. Telegram allows roughly thirty messages a second
// to distinct users; a gap well inside that is free at any size we will reach
// and keeps us clear of the rate limiter.
const tgSendGap = 40 * time.Millisecond

// tgUpdate is the slice of Telegram's update object we act on.
type tgUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  *struct {
		Text string `json:"text"`
		From struct {
			ID           int64  `json:"id"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
	} `json:"message"`
	CallbackQuery *struct {
		ID   string `json:"id"`
		Data string `json:"data"`
		From struct {
			ID           int64  `json:"id"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
	} `json:"callback_query"`
}

// tgButton is one inline keyboard button.
type tgButton struct {
	Text string `json:"text"`
	Data string `json:"callback_data"`
}

// RunTelegramBot serves the bot until the context is cancelled.
func (m *Module) RunTelegramBot(ctx context.Context) {
	if !m.tgBotEnabled {
		return
	}
	// The command menu, so the person is not left guessing what to type. Best
	// effort: a bot that cannot publish its menu still works.
	if err := m.publishBotCommands(ctx); err != nil {
		m.log.Warn("telegram setMyCommands", zap.Error(err))
	}
	m.log.Info("telegram bot polling started")
	for {
		if err := m.pollOnce(ctx); err != nil {
			if ctx.Err() != nil {
				return
			}
			m.log.Warn("telegram poll", zap.Error(err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Second):
			}
		}
		if ctx.Err() != nil {
			return
		}
	}
}

// pollOnce fetches one batch of updates and handles them.
func (m *Module) pollOnce(ctx context.Context) error {
	offset, err := m.botOffset(ctx)
	if err != nil {
		return err
	}
	var out struct {
		OK     bool       `json:"ok"`
		Result []tgUpdate `json:"result"`
	}
	// The request must be allowed to outlive the long-poll window itself.
	callCtx, cancel := context.WithTimeout(ctx, (tgPollTimeout+15)*time.Second)
	defer cancel()
	if err := m.tgCall(callCtx, "getUpdates", map[string]any{
		"offset":          offset + 1,
		"timeout":         tgPollTimeout,
		"allowed_updates": []string{"message", "callback_query"},
	}, &out); err != nil {
		return err
	}
	for _, u := range out.Result {
		if err := m.handleUpdate(ctx, u); err != nil {
			m.log.Warn("telegram update", zap.Int64("update_id", u.UpdateID), zap.Error(err))
		}
		// The offset advances even for an update that failed. Retrying a bad
		// update for ever would wedge the bot behind it and stop every other
		// subscriber being served.
		if err := m.setBotOffset(ctx, u.UpdateID); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) botOffset(ctx context.Context) (int64, error) {
	var id int64
	err := m.db.QueryRow(ctx, `SELECT update_id FROM telegram_bot_state WHERE one`).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("bot offset: %w", err)
	}
	return id, nil
}

func (m *Module) setBotOffset(ctx context.Context, id int64) error {
	_, err := m.db.Exec(ctx,
		`INSERT INTO telegram_bot_state (one, update_id) VALUES (TRUE,$1)
		 ON CONFLICT (one) DO UPDATE SET update_id = GREATEST(telegram_bot_state.update_id, EXCLUDED.update_id)`, id)
	if err != nil {
		return fmt.Errorf("set bot offset: %w", err)
	}
	return nil
}

// handleUpdate routes one update.
func (m *Module) handleUpdate(ctx context.Context, u tgUpdate) error {
	switch {
	case u.CallbackQuery != nil:
		lang := tgLang(u.CallbackQuery.From.LanguageCode)
		if err := m.handleCallback(ctx, u.CallbackQuery.From.ID, lang, u.CallbackQuery.Data); err != nil {
			return err
		}
		// Telegram spins the button until the query is answered.
		return m.tgCall(ctx, "answerCallbackQuery", map[string]any{"callback_query_id": u.CallbackQuery.ID}, nil)
	case u.Message != nil:
		return m.handleCommand(ctx, u.Message.From.ID, tgLang(u.Message.From.LanguageCode),
			strings.TrimSpace(u.Message.Text))
	}
	return nil
}

// handleCommand answers a typed message.
func (m *Module) handleCommand(ctx context.Context, user int64, lang, text string) error {
	cmd := strings.ToLower(strings.Fields(text + " ")[0])
	switch cmd {
	case "/start", "/place", "/mesto":
		if err := m.ensureSubscriber(ctx, user, lang); err != nil {
			return err
		}
		if cmd == "/start" {
			if err := m.sendTo(ctx, user, tgText(lang, "welcome")); err != nil {
				return err
			}
		}
		return m.sendPlacePicker(ctx, user, lang, uuid.Nil)
	case "/stop":
		if _, err := m.db.Exec(ctx,
			`UPDATE telegram_subscribers SET active = FALSE, updated_at = NOW() WHERE tg_user_id = $1`, user); err != nil {
			return fmt.Errorf("stop subscriber: %w", err)
		}
		return m.sendTo(ctx, user, tgText(lang, "stopped"))
	default:
		return m.sendTo(ctx, user, tgText(lang, "help"))
	}
}

// handleCallback acts on a keyboard tap: open a place, or choose it.
func (m *Module) handleCallback(ctx context.Context, user int64, lang, data string) error {
	action, raw, ok := strings.Cut(data, ":")
	if !ok {
		return nil
	}
	switch action {
	case "o": // open: show what is inside this place
		var parent uuid.UUID
		if raw != "" {
			id, err := uuid.Parse(raw)
			if err != nil {
				return nil
			}
			parent = id
		}
		return m.sendPlacePicker(ctx, user, lang, parent)
	case "p": // pick: this is my place
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil
		}
		if err := m.setSubscriberPlace(ctx, user, id, lang); err != nil {
			return err
		}
		name, err := m.placeLabel(ctx, id, lang)
		if err != nil {
			return err
		}
		return m.sendTo(ctx, user, fmt.Sprintf(tgText(lang, "chosen"), name))
	}
	return nil
}

// sendPlacePicker offers the places inside parent; uuid.Nil starts at the top.
//
// Two taps reach most towns: regions first, then what is inside the one chosen.
// Every level also offers the place itself, because somebody who wants the whole
// region should not have to pick a town in it to get anything.
func (m *Module) sendPlacePicker(ctx context.Context, user int64, lang string, parent uuid.UUID) error {
	kids, err := m.placeChildren(ctx, parent, lang)
	if err != nil {
		return err
	}
	if len(kids) == 0 && parent != uuid.Nil {
		// A leaf: nothing to descend into, so choosing it is the only move.
		return m.handleCallback(ctx, user, lang, "p:"+parent.String())
	}

	rows := [][]tgButton{}
	if parent != uuid.Nil {
		name, err := m.placeLabel(ctx, parent, lang)
		if err != nil {
			return err
		}
		rows = append(rows, []tgButton{{
			Text: fmt.Sprintf(tgText(lang, "pick_this"), name),
			Data: "p:" + parent.String(),
		}})
	}
	for _, k := range kids {
		rows = append(rows, []tgButton{{Text: k.Name, Data: "o:" + k.ID}})
	}
	if parent != uuid.Nil {
		rows = append(rows, []tgButton{{Text: tgText(lang, "back"), Data: "o:"}})
	}

	prompt := tgText(lang, "pick_region")
	if parent != uuid.Nil {
		prompt = tgText(lang, "pick_inside")
	}
	return m.tgCall(ctx, "sendMessage", map[string]any{
		"chat_id":      user,
		"text":         prompt,
		"reply_markup": map[string]any{"inline_keyboard": rows},
	}, nil)
}

// tgPlace is one option in the picker.
type tgPlace struct {
	ID   string
	Name string
}

// placeChildren lists what sits directly inside a place. At the top it lists
// the regions of Kazakhstan rather than the country itself: nobody subscribes
// to "Kazakhstan" here, that is what the channel is for.
func (m *Module) placeChildren(ctx context.Context, parent uuid.UUID, lang string) ([]tgPlace, error) {
	col := tgNameColumn(lang)
	var rows pgx.Rows
	var err error
	if parent == uuid.Nil {
		rows, err = m.db.Query(ctx, fmt.Sprintf(`
			SELECT c.id::text, COALESCE(NULLIF(c.%s,''), c.name_ru)
			FROM geo_nodes c
			WHERE c.country = 'KZ' AND c.level = 1
			ORDER BY c.population DESC NULLS LAST, 2`, col))
	} else {
		rows, err = m.db.Query(ctx, fmt.Sprintf(`
			SELECT c.id::text, COALESCE(NULLIF(c.%s,''), c.name_ru)
			FROM geo_nodes c
			WHERE c.parent_id = $1
			ORDER BY c.population DESC NULLS LAST, 2`, col), parent)
	}
	if err != nil {
		return nil, fmt.Errorf("place children: %w", err)
	}
	defer rows.Close()
	out := []tgPlace{}
	for rows.Next() {
		var p tgPlace
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (m *Module) placeLabel(ctx context.Context, id uuid.UUID, lang string) (string, error) {
	var name string
	err := m.db.QueryRow(ctx, fmt.Sprintf(
		`SELECT COALESCE(NULLIF(%s,''), name_ru) FROM geo_nodes WHERE id = $1`,
		tgNameColumn(lang)), id).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("place label: %w", err)
	}
	return name, nil
}

func tgNameColumn(lang string) string {
	switch lang {
	case "kz":
		return "name_kk"
	case "en":
		return "name_en"
	default:
		return "name_ru"
	}
}

func (m *Module) ensureSubscriber(ctx context.Context, user int64, lang string) error {
	_, err := m.db.Exec(ctx, `
		INSERT INTO telegram_subscribers (tg_user_id, lang) VALUES ($1,$2)
		ON CONFLICT (tg_user_id) DO UPDATE SET active = TRUE, lang = EXCLUDED.lang, updated_at = NOW()`,
		user, lang)
	if err != nil {
		return fmt.Errorf("ensure subscriber: %w", err)
	}
	return nil
}

func (m *Module) setSubscriberPlace(ctx context.Context, user int64, place uuid.UUID, lang string) error {
	_, err := m.db.Exec(ctx, `
		INSERT INTO telegram_subscribers (tg_user_id, geo_node_id, lang) VALUES ($1,$2,$3)
		ON CONFLICT (tg_user_id) DO UPDATE
		  SET geo_node_id = EXCLUDED.geo_node_id, lang = EXCLUDED.lang, active = TRUE, updated_at = NOW()`,
		user, place, lang)
	if err != nil {
		return fmt.Errorf("set subscriber place: %w", err)
	}
	return nil
}

// subscribersInside returns the people a piece written for `place` is addressed
// to: those whose chosen place is that place or sits inside it.
//
// The direction is the mirror of the site's feed. A piece for the region reaches
// everyone in it, including the villages; a piece for one village reaches that
// village alone.
func (m *Module) subscribersInside(ctx context.Context, place uuid.UUID) ([]int64, []string, error) {
	rows, err := m.db.Query(ctx, `
		WITH RECURSIVE down AS (
			SELECT id FROM geo_nodes WHERE id = $1
			UNION ALL
			SELECT g.id FROM geo_nodes g JOIN down ON g.parent_id = down.id
		)
		SELECT s.tg_user_id, s.lang
		FROM telegram_subscribers s
		WHERE s.active AND s.geo_node_id IN (SELECT id FROM down)`, place)
	if err != nil {
		return nil, nil, fmt.Errorf("subscribers inside: %w", err)
	}
	defer rows.Close()
	var ids []int64
	var langs []string
	for rows.Next() {
		var id int64
		var lang string
		if err := rows.Scan(&id, &lang); err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
		langs = append(langs, lang)
	}
	return ids, langs, rows.Err()
}

// sendTo delivers one message to one person, retiring the subscriber if they
// have blocked the bot.
func (m *Module) sendTo(ctx context.Context, user int64, text string) error {
	err := m.tgCall(ctx, "sendMessage", map[string]any{
		"chat_id":                  user,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": false,
	}, nil)
	if err != nil && tgBlocked(err) {
		if _, derr := m.db.Exec(ctx,
			`UPDATE telegram_subscribers SET active = FALSE, updated_at = NOW() WHERE tg_user_id = $1`,
			user); derr != nil {
			m.log.Warn("retire subscriber", zap.Error(derr))
		}
		return nil
	}
	return err
}

// tgBlocked reports whether Telegram refused because the person is gone: they
// blocked the bot, deleted the account, or never started it.
func tgBlocked(err error) bool {
	s := err.Error()
	return strings.Contains(s, "403") ||
		strings.Contains(s, "bot was blocked") ||
		strings.Contains(s, "user is deactivated") ||
		strings.Contains(s, "chat not found")
}

// tgCall performs one Bot API method.
func (m *Module) tgCall(ctx context.Context, method string, payload map[string]any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", m.tgBotToken, method)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("build %s: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	// getUpdates blocks for the poll window; the shared client's short timeout
	// would cut it off every time.
	client := m.http
	if method == "getUpdates" {
		client = &http.Client{Timeout: (tgPollTimeout + 20) * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var buf [512]byte
		n, _ := resp.Body.Read(buf[:])
		return fmt.Errorf("%s: telegram %d: %s", method, resp.StatusCode, strings.TrimSpace(string(buf[:n])))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// tgLang maps Telegram's language code onto ours, defaulting to Russian.
func tgLang(code string) string {
	switch {
	case strings.HasPrefix(code, "kk"):
		return "kz"
	case strings.HasPrefix(code, "en"):
		return "en"
	default:
		return "ru"
	}
}

// tgText holds the bot's own words. They live here rather than in the articles
// module's dictionary because the bot is the only thing that says them.
func tgText(lang, key string) string {
	table := map[string]map[string]string{
		"welcome": {
			"kz": "Сәлем! Мен Shanraq.org ботымын.\n\nЕлге қатысты материалдар арнада жарияланады. Ал мен сізге тек өз жеріңізге арналған материалдарды жіберемін — басқа ешкімге кетпейтіндерін.\n\nҚай жерде тұратыныңызды таңдаңыз. Біз оны өзіміз анықтамаймыз.",
			"ru": "Здравствуйте! Это бот Shanraq.org.\n\nМатериалы для всей страны выходят в канале. А я присылаю только то, что написано для вашего места, — то, что другим не уходит.\n\nВыберите, где вы живёте. Сами мы это не определяем.",
			"en": "Hello. This is the Shanraq.org bot.\n\nMaterial for the whole country goes out in the channel. I send only what was written for your own place — the things nobody else receives.\n\nChoose where you live. We do not work it out ourselves.",
		},
		"pick_region": {
			"kz": "Облысты немесе қаланы таңдаңыз:",
			"ru": "Выберите область или город:",
			"en": "Choose a region or city:",
		},
		"pick_inside": {
			"kz": "Нақтылаңыз немесе осы жерді таңдаңыз:",
			"ru": "Уточните или выберите это место целиком:",
			"en": "Narrow it down, or take this place as a whole:",
		},
		"pick_this": {"kz": "✓ %s — осы", "ru": "✓ %s — целиком", "en": "✓ %s — all of it"},
		"back":      {"kz": "← Артқа", "ru": "← Назад", "en": "← Back"},
		"chosen": {
			"kz": "Дайын. Енді сізге %s үшін жазылған материалдар келеді.\n\nАуыстыру — /place, тоқтату — /stop.",
			"ru": "Готово. Теперь вам приходят материалы, написанные для места: %s.\n\nПоменять — /place, отписаться — /stop.",
			"en": "Done. You will now receive material written for %s.\n\nChange it with /place, stop with /stop.",
		},
		"stopped": {
			"kz": "Жіберуді тоқтаттым. Қайта қосылу — /start.",
			"ru": "Больше ничего не присылаю. Вернуться — /start.",
			"en": "I have stopped sending. Come back with /start.",
		},
		"cmd_start": {"kz": "Жазылу және орынды таңдау", "ru": "Подписаться и выбрать место", "en": "Subscribe and choose your place"},
		"cmd_place": {"kz": "Орынды ауыстыру", "ru": "Сменить место", "en": "Change your place"},
		"cmd_stop":  {"kz": "Хабарламаларды тоқтату", "ru": "Отписаться", "en": "Stop messages"},
		"help": {
			"kz": "/start — жазылу, /place — жерді ауыстыру, /stop — тоқтату.",
			"ru": "/start — подписаться, /place — сменить место, /stop — отписаться.",
			"en": "/start to subscribe, /place to change your place, /stop to stop.",
		},
	}
	row, ok := table[key]
	if !ok {
		return ""
	}
	if s, ok := row[lang]; ok {
		return s
	}
	return row["ru"]
}

// publishBotCommands registers the command menu shown by Telegram's UI.
func (m *Module) publishBotCommands(ctx context.Context) error {
	for _, l := range []struct{ code, lang string }{{"", "ru"}, {"kk", "kz"}, {"en", "en"}} {
		payload := map[string]any{"commands": []map[string]string{
			{"command": "start", "description": tgText(l.lang, "cmd_start")},
			{"command": "place", "description": tgText(l.lang, "cmd_place")},
			{"command": "stop", "description": tgText(l.lang, "cmd_stop")},
		}}
		if l.code != "" {
			payload["language_code"] = l.code
		}
		if err := m.tgCall(ctx, "setMyCommands", payload, nil); err != nil {
			return err
		}
	}
	return nil
}
