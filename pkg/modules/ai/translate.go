package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
)

// JobTranslate is the queue job name for asynchronous AI translation.
const JobTranslate = "ai_translate"

// allLangs mirrors the article language set (kept local to avoid importing the
// articles package, which would create an import cycle).
var allLangs = []string{"kz", "ru", "en"}

var langFullName = map[string]string{
	"kz": "Kazakh (Қазақ тілі)",
	"ru": "Russian (Русский язык)",
	"en": "English",
}

// TranslatePayload is the job payload for JobTranslate.
type TranslatePayload struct {
	ArticleID string `json:"article_id"`
}

// content holds one article version's editable text.
type content struct {
	Title   string
	Summary string
	Body    string
}

// EnqueuePayload builds a JSON payload for a translate job.
func EnqueuePayload(articleID uuid.UUID) (json.RawMessage, error) {
	return json.Marshal(TranslatePayload{ArticleID: articleID.String()})
}

// handleTranslateJob translates an article from its original language into the
// remaining languages, writing AI versions that don't clobber human ones.
func (m *Module) handleTranslateJob(ctx context.Context, _ *shanraq.Runtime, job jobs.Job) error {
	if !m.Enabled() {
		m.log.Warn("ai_translate skipped: assistant disabled")
		return nil
	}

	var payload TranslatePayload
	if err := job.Decode(&payload); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	id, err := uuid.Parse(payload.ArticleID)
	if err != nil {
		return fmt.Errorf("bad article id: %w", err)
	}

	origLang, src, err := m.loadOriginal(ctx, id)
	if err != nil {
		return err
	}
	if src.Title == "" || src.Body == "" {
		m.log.Warn("ai_translate skipped: original has no content", zap.String("article_id", payload.ArticleID))
		return nil
	}

	for _, target := range allLangs {
		if target == origLang {
			continue
		}
		human, err := m.hasHumanTranslation(ctx, id, target)
		if err != nil {
			return err
		}
		if human {
			continue // never overwrite a human-authored version
		}
		out, err := m.translateContent(ctx, origLang, target, src)
		if err != nil {
			return fmt.Errorf("translate %s->%s: %w", origLang, target, err)
		}
		if err := m.saveAITranslation(ctx, id, target, out); err != nil {
			return err
		}
		m.log.Info("ai translated article", zap.String("article_id", payload.ArticleID), zap.String("lang", target))
	}
	return nil
}

// translateContent translates title, summary, and body from one language to
// another. Empty fields are skipped.
func (m *Module) translateContent(ctx context.Context, from, to string, src content) (content, error) {
	system := translateSystem(from, to)
	var out content
	var err error

	if src.Title != "" {
		if out.Title, err = m.completer.Complete(ctx, Request{Model: m.translateModel, System: system, User: src.Title, MaxTokens: 512}); err != nil {
			return content{}, err
		}
	}
	if src.Summary != "" {
		if out.Summary, err = m.completer.Complete(ctx, Request{Model: m.translateModel, System: system, User: src.Summary, MaxTokens: 1024}); err != nil {
			return content{}, err
		}
	}
	if out.Body, err = m.completer.Complete(ctx, Request{Model: m.translateModel, System: system, User: src.Body, MaxTokens: m.maxTokens}); err != nil {
		return content{}, err
	}
	return out, nil
}

func translateSystem(from, to string) string {
	return fmt.Sprintf(`You are a professional translator for an independent journalism platform.
Translate the user's text from %s into %s.
Rules:
- Preserve the meaning, tone, and any Markdown formatting exactly.
- Keep proper nouns and technical terms accurate; localize idioms naturally.
- Do NOT add commentary, notes, or explanations.
- Output ONLY the translated text, nothing else.`, langFullName[from], langFullName[to])
}

// ---- DB access (raw, to avoid importing the articles package) ----

func (m *Module) loadOriginal(ctx context.Context, articleID uuid.UUID) (string, content, error) {
	var lang string
	var c content
	err := m.db.QueryRow(ctx, `
		SELECT a.original_lang, t.title, t.summary, t.body_md
		FROM articles a
		JOIN article_translations t ON t.article_id = a.id AND t.lang = a.original_lang
		WHERE a.id = $1
	`, articleID).Scan(&lang, &c.Title, &c.Summary, &c.Body)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", content{}, fmt.Errorf("no original content for article %s", articleID)
		}
		return "", content{}, fmt.Errorf("load original: %w", err)
	}
	return lang, c, nil
}

func (m *Module) hasHumanTranslation(ctx context.Context, articleID uuid.UUID, lang string) (bool, error) {
	var exists bool
	err := m.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM article_translations
			WHERE article_id = $1 AND lang = $2 AND source = 'human' AND title <> '' AND body_md <> ''
		)
	`, articleID, lang).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check human translation: %w", err)
	}
	return exists, nil
}

func (m *Module) saveAITranslation(ctx context.Context, articleID uuid.UUID, lang string, c content) error {
	_, err := m.db.Exec(ctx, `
		INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status)
		VALUES ($1, $2, $3, $4, $5, 'ai', 'ready')
		ON CONFLICT (article_id, lang) DO UPDATE SET
			title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			body_md = EXCLUDED.body_md,
			source = 'ai',
			status = 'ready',
			updated_at = NOW()
	`, articleID, lang, strings.TrimSpace(c.Title), strings.TrimSpace(c.Summary), strings.TrimSpace(c.Body))
	if err != nil {
		return fmt.Errorf("save ai translation: %w", err)
	}
	return nil
}
