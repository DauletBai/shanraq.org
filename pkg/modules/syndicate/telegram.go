package syndicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
)

// JobTelegram is the queue job that posts a published article to Telegram.
const JobTelegram = "syndicate_telegram"

// TelegramJobPayload carries the article to announce.
type TelegramJobPayload struct {
	ArticleID string `json:"article_id"`
}

// TelegramPayload builds a job payload for an article.
func TelegramPayload(articleID uuid.UUID) (json.RawMessage, error) {
	return json.Marshal(TelegramJobPayload{ArticleID: articleID.String()})
}

func (m *Module) handleTelegramJob(ctx context.Context, _ *shanraq.Runtime, job jobs.Job) error {
	if !m.tgEnabled {
		return nil
	}
	var payload TelegramJobPayload
	if err := job.Decode(&payload); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	id, err := uuid.Parse(payload.ArticleID)
	if err != nil {
		return fmt.Errorf("bad article id: %w", err)
	}

	slug, title, summary, lang, err := m.loadAnnouncement(ctx, id)
	if err != nil {
		return err
	}
	text := buildTelegramMessage(title, summary, m.articleURL(slug, lang))
	if err := m.sendTelegram(ctx, text); err != nil {
		return err
	}
	m.log.Info("telegram announced article", zap.String("article_id", payload.ArticleID))
	return nil
}

func (m *Module) loadAnnouncement(ctx context.Context, articleID uuid.UUID) (slug, title, summary, lang string, err error) {
	err = m.db.QueryRow(ctx, `
		SELECT a.slug, a.original_lang, t.title, t.summary
		FROM articles a
		JOIN article_translations t ON t.article_id = a.id AND t.lang = a.original_lang
		WHERE a.id = $1 AND a.status = 'published'
	`, articleID).Scan(&slug, &lang, &title, &summary)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", "", fmt.Errorf("article %s not published", articleID)
		}
		return "", "", "", "", fmt.Errorf("load announcement: %w", err)
	}
	return slug, title, summary, lang, nil
}

// buildTelegramMessage formats an HTML-safe Telegram announcement.
func buildTelegramMessage(title, summary, url string) string {
	var b strings.Builder
	b.WriteString("📰 <b>")
	b.WriteString(html.EscapeString(strings.TrimSpace(title)))
	b.WriteString("</b>")
	if s := strings.TrimSpace(summary); s != "" {
		b.WriteString("\n\n")
		b.WriteString(html.EscapeString(s))
	}
	b.WriteString("\n\n🔗 <a href=\"")
	b.WriteString(html.EscapeString(url))
	b.WriteString("\">Оқу · Читать</a>")
	return b.String()
}

func (m *Module) sendTelegram(ctx context.Context, text string) error {
	body, _ := json.Marshal(map[string]any{
		"chat_id":                  m.tgChatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": false,
	})
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", m.tgBotToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.http.Do(req)
	if err != nil {
		return fmt.Errorf("telegram request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return fmt.Errorf("telegram api status %d: %s", resp.StatusCode, strings.TrimSpace(buf.String()))
	}
	return nil
}
