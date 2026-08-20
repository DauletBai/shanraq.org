package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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
	c, model, tok := m.translateClient()
	if c == nil {
		return content{}, ErrDisabled
	}
	system := translateSystem(from, to)
	var out content
	var err error

	// The body goes first, and everything else is translated against it.
	//
	// A headline and a summary are ambiguous on their own: "корь" came back as
	// "қызамық" — rubella — in the summary of an article whose body said
	// "қызылша" throughout. Showing them the Russian source as context did not
	// help, because the question is not what the article is about but which
	// word to use in the target language. Showing them the finished translation
	// does: the term has already been chosen, and it is right there to copy.
	if out.Body, err = c.Complete(ctx, Request{Model: model, System: system, User: src.Body,
		MaxTokens: translationBudget(src.Body, tok)}); err != nil {
		return content{}, err
	}
	// The one thing a reader cannot verify and the model got wrong: where the
	// links point.
	if restored, ok := restoreLinkTargets(src.Body, out.Body); ok {
		out.Body = restored
	} else {
		m.log.Warn("translation changed the number of links; URLs left as the model wrote them",
			zap.String("from", from), zap.String("to", to))
	}

	glossary := excerptForContext(out.Body)
	if src.Title != "" {
		if out.Title, err = c.Complete(ctx, Request{Model: model, System: system,
			User:      withGlossary(glossary, src.Title, langFullName[to]),
			MaxTokens: translationBudget(src.Title, 512)}); err != nil {
			return content{}, err
		}
		out.Title = afterMarker(out.Title)
	}
	if src.Summary != "" {
		if out.Summary, err = c.Complete(ctx, Request{Model: model, System: system,
			User:      withGlossary(glossary, src.Summary, langFullName[to]),
			MaxTokens: translationBudget(src.Summary, 1024)}); err != nil {
			return content{}, err
		}
		out.Summary = afterMarker(out.Summary)
	}
	return out, nil
}

// excerptForContext takes the opening of an article, enough to fix the subject
// and its terminology without paying for the whole text on every short field.
func excerptForContext(body string) string {
	r := []rune(body)
	if len(r) > 700 {
		r = r[:700]
	}
	return strings.TrimSpace(string(r))
}

// translateMarker separates the glossary from the field to be translated. It is
// a constant because two things must agree on it: the prompt that writes it and
// the code that cuts on it.
const translateMarker = "TRANSLATE ONLY WHAT FOLLOWS"

// withGlossary frames a short field against the already-translated article, so
// the wording chosen there is the wording used here.
func withGlossary(translated, field, targetName string) string {
	if translated == "" {
		return field
	}
	return "The same article, already translated into " + targetName +
		". Use exactly its terminology and wording — do NOT translate this part:\n\n" +
		translated + "\n\n---" + translateMarker + ", matching the terminology above---\n\n" + field
}

// afterMarker keeps only what the model wrote after the marker.
//
// The instruction not to translate the glossary is an instruction, and a model
// handed two texts translated both — returning the context, the marker line,
// and only then the answer. The summary of an article came out as 1,275
// characters against an original of 461, with the real translation buried at
// the end. Asking more politely is not a fix; cutting on the marker is.
func afterMarker(reply string) string {
	i := strings.LastIndex(reply, translateMarker)
	if i < 0 {
		return reply
	}
	rest := reply[i+len(translateMarker):]
	// Step over the tail of the marker line, whatever the model made of it.
	if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
		rest = rest[nl+1:]
	}
	return strings.TrimSpace(rest)
}

// mdLinkTarget matches the URL inside a Markdown link.
var mdLinkTarget = regexp.MustCompile(`\]\((https?://[^)\s]+)\)`)

// restoreLinkTargets puts the original URLs back into a translation.
//
// The prompt tells the model to keep URLs byte for byte, and it mostly does —
// but on the first long article one of nine came back swapped: the WHO fact
// sheet pointed at a Bangladeshi newspaper. A prompt cannot guarantee this and
// a reader cannot check it, so the guarantee is made here instead.
//
// Links are restored by position: the nth link of the translation gets the nth
// URL of the source. If the counts differ the translation has changed the
// document's shape, and guessing which link is which would be worse than
// leaving it — the caller logs that and moves on.
func restoreLinkTargets(src, translated string) (string, bool) {
	want := mdLinkTarget.FindAllStringSubmatch(src, -1)
	if len(want) == 0 || len(want) != len(mdLinkTarget.FindAllString(translated, -1)) {
		return translated, len(want) == 0
	}
	i := 0
	out := mdLinkTarget.ReplaceAllStringFunc(translated, func(string) string {
		u := want[i][1]
		i++
		return "](" + u + ")"
	})
	return out, true
}

// translationBudget sizes the output cap to the text being translated.
//
// A translation is about as long as its source, so a fixed ceiling is the wrong
// shape: it is wasteful for a headline and fatal for an article. The configured
// max_tokens — 4096 by default — cut off the first long piece this site tried to
// translate, and the job failed three times with nothing to show the author.
//
// Cyrillic runs roughly two characters to the token, and a Kazakh rendering of
// a Russian text can come out longer, so the estimate is halved characters with
// half again on top. The floor keeps short texts from getting an absurdly small
// allowance; the ceiling is there because a cap is still a cap.
//
// Nothing is spent by asking for headroom: max_tokens is a limit, not a
// reservation, and billing follows the tokens actually produced.
func translationBudget(src string, floor int) int {
	est := len([]rune(src)) / 2 * 3 / 2
	if est < floor {
		est = floor
	}
	if est > maxTranslationTokens {
		est = maxTranslationTokens
	}
	return est
}

// maxTranslationTokens bounds a single translation call. Beyond this a text
// should be split rather than asked for in one breath.
const maxTranslationTokens = 16000

func translateSystem(from, to string) string {
	// "Preserve the Markdown formatting exactly" read, to a model not given room
	// to think it over, as "leave everything inside the markup alone" — and a
	// link came back with its label still in the source language, sitting in the
	// middle of a translated sentence. The rule now separates the two things
	// that got conflated: the markup and the URLs are structure and are kept
	// byte for byte; the words a reader sees, wherever they sit, are prose and
	// are translated.
	return fmt.Sprintf(`You are a professional translator for an independent journalism platform.
Translate the user's text from %s into %s.
Rules:
- Translate ALL text a reader sees, including link labels, headings, list items,
  table cells, blockquotes, and image alt text. A word inside [brackets] is prose
  and must be translated like any other.
- Keep the Markdown structure byte for byte: the markup characters themselves,
  the URLs inside (parentheses), and anything in code blocks or `+"`inline code`"+`.
- Preserve the meaning and tone. Write natural, idiomatic %[2]s — a stiff literal
  rendering is a worse translation than a fluent one.
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
