package syndicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

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

	slug, title, summary, lang, place, err := m.loadAnnouncement(ctx, id)
	if err != nil {
		return err
	}
	// Tag the auto-posted link so visits are attributed to Telegram in analytics
	// even when the messenger strips the referrer.
	url := m.articleURL(slug, lang) + "&utm_source=telegram"

	// A piece written for everybody goes to the channel. A piece written for a
	// place does not: a channel sends one identical message to every subscriber
	// in the country, and a village's news arriving on twenty thousand phones is
	// not reach, it is noise. Those go through the bot, to the people the piece
	// is actually addressed to.
	if place == uuid.Nil {
		if !m.tgEnabled {
			return nil
		}
		if err := m.sendTelegram(ctx, buildTelegramMessage(title, summary, url)); err != nil {
			return err
		}
		m.log.Info("telegram announced article to the channel", zap.String("article_id", payload.ArticleID))
		return nil
	}
	return m.announceToSubscribers(ctx, payload.ArticleID, place, title, summary, url)
}

// announceToSubscribers delivers a placed article to the people inside that
// place, one message each.
func (m *Module) announceToSubscribers(ctx context.Context, articleID string, place uuid.UUID, title, summary, url string) error {
	if !m.tgBotEnabled {
		return nil
	}
	ids, langs, err := m.subscribersInside(ctx, place)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		m.log.Info("telegram: nobody subscribed to this place yet",
			zap.String("article_id", articleID), zap.String("place", place.String()))
		return nil
	}
	sent := 0
	for i, id := range ids {
		label, err := m.placeLabel(ctx, place, langs[i])
		if err != nil {
			label = ""
		}
		text := buildPlacedMessage(label, title, summary, url)
		if err := m.sendTo(ctx, id, text); err != nil {
			// One unreachable person must not stop the rest of the delivery.
			m.log.Warn("telegram send", zap.Int64("user", id), zap.Error(err))
		} else {
			sent++
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(tgSendGap):
		}
	}
	m.log.Info("telegram announced article to subscribers",
		zap.String("article_id", articleID), zap.Int("sent", sent), zap.Int("addressed", len(ids)))
	return nil
}

func (m *Module) loadAnnouncement(ctx context.Context, articleID uuid.UUID) (slug, title, summary, lang string, place uuid.UUID, err error) {
	var node *uuid.UUID
	err = m.db.QueryRow(ctx, `
		SELECT a.slug, a.original_lang, t.title, t.summary, a.geo_node_id
		FROM articles a
		JOIN article_translations t ON t.article_id = a.id AND t.lang = a.original_lang
		WHERE a.id = $1 AND a.status = 'published'
	`, articleID).Scan(&slug, &lang, &title, &summary, &node)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", "", uuid.Nil, fmt.Errorf("article %s not published", articleID)
		}
		return "", "", "", "", uuid.Nil, fmt.Errorf("load announcement: %w", err)
	}
	if node != nil {
		place = *node
	}
	return slug, title, summary, lang, place, nil
}

// buildPlacedMessage is the announcement a subscriber gets, headed by the place
// it was written for. The line answers the question a personal message provokes
// and a channel post does not: why am I being told this.
func buildPlacedMessage(place, title, summary, url string) string {
	var b strings.Builder
	if p := strings.TrimSpace(place); p != "" {
		b.WriteString("📍 <b>")
		b.WriteString(html.EscapeString(p))
		b.WriteString("</b>\n\n")
	}
	b.WriteString(buildTelegramMessage(title, summary, url))
	return b.String()
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
