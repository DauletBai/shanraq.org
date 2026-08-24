package articles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/ai"
)

// The reader's half of proof-reading: a page with three fields, and the outcome
// on the same page a moment later.
//
// The outcome is synchronous on purpose. A reader who has just told us about a
// typo is owed an answer while they are still looking at the screen; a queue would
// mean they never learn whether they were right, and the promise that a refusal
// comes with a reason would be a promise to send an e-mail nobody reads.

// CorrectionPage is the form, and its result once submitted.
type CorrectionPage struct {
	Base
	Slug  string
	Title string
	// Chapters offers the article's own headings as suggestions, so the chapter
	// field can be filled from a list instead of retyped from memory.
	Chapters []TOCItem

	// The three fields, echoed back so nothing is retyped after a refusal.
	Chapter  string
	Sentence string
	Word     string

	// Authed is false for a guest: the outcome has to reach a person, and an
	// anonymous claim has nobody to reach.
	Authed bool

	// Done marks a submitted claim; Status is what became of it.
	Done   bool
	Status string
	// Fixed and WordWas are the two sides of an applied edit.
	Fixed   string
	WordWas string
	// Reason explains a refusal, in the reader's language.
	Reason string
	// Err names a field-level problem with the form itself.
	Err string
}

// handleCorrectionForm serves the empty form.
func (m *Module) handleCorrectionForm(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	slug := chi.URLParam(r, "slug")
	a, err := m.store.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tr, served := a.Translation(lang)
	if tr == nil {
		http.NotFound(w, r)
		return
	}
	_, authed := m.authorID(r)
	page := CorrectionPage{
		Slug:   slug,
		Title:  tr.Title,
		Authed: authed,
	}
	_, page.Chapters = RenderMarkdownTOC(tr.BodyMD)
	page.Base = m.base(r, T(lang, "corr.title"), lang)
	page.Base.CanonURL = "/read/" + slug + "/typo"
	// A form is not a page a search engine has any use for, and it must not
	// compete with the article it belongs to.
	page.Base.NoIndex = true
	_ = served
	m.render(w, "correction", page)
}

// handleCorrectionSubmit takes a claim, decides it, and shows the outcome.
func (m *Module) handleCorrectionSubmit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	slug := chi.URLParam(r, "slug")

	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	a, err := m.store.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tr, served := a.Translation(lang)
	if tr == nil {
		http.NotFound(w, r)
		return
	}

	_ = r.ParseForm()
	page := CorrectionPage{
		Slug:     slug,
		Title:    tr.Title,
		Authed:   true,
		Chapter:  clipField(r.FormValue("chapter"), correctionChapterMax),
		Sentence: clipField(r.FormValue("sentence"), correctionSentenceMax),
		Word:     clipField(r.FormValue("word"), correctionWordMax),
	}
	_, page.Chapters = RenderMarkdownTOC(tr.BodyMD)
	page.Base = m.base(r, T(lang, "corr.title"), lang)
	page.Base.CanonURL = "/read/" + slug + "/typo"
	page.Base.NoIndex = true

	switch {
	case page.Sentence == "":
		page.Err = "corr.e_sentence"
	case page.Word == "":
		page.Err = "corr.e_word"
	case len(strings.Fields(page.Word)) > 1:
		// One word means one word. A phrase here is the commonest way a claim
		// arrives unusable, and saying so is more use than a silent refusal.
		page.Err = "corr.e_one_word"
	case !wordSet(page.Sentence)[foldWord(page.Word)]:
		page.Err = "corr.e_not_in_sentence"
	}
	if page.Err != "" {
		m.render(w, "correction", page)
		return
	}

	c := Correction{
		ArticleID: a.ID,
		Lang:      served,
		Reporter:  uid,
		Chapter:   page.Chapter,
		Sentence:  page.Sentence,
		Word:      page.Word,
	}
	page.Done = true
	page.WordWas = page.Word
	page.Status, page.Fixed, page.Reason = m.decideCorrection(r.Context(), a, tr, served, c)
	m.render(w, "correction", page)
}

// fmtT fills a translated template with its arguments.
func fmtT(lang, key string, args ...any) string {
	return fmt.Sprintf(T(lang, key), args...)
}

// clipField trims a submitted field and bounds its length in runes.
func clipField(s string, max int) string {
	s = normalizeSpace(s)
	if r := []rune(s); len(r) > max {
		s = strings.TrimSpace(string(r[:max]))
	}
	return s
}

// decideCorrection runs a claim through the four tiers and returns the outcome:
// a status, the word's new form when one was applied, and the reason when it was
// not.
func (m *Module) decideCorrection(ctx context.Context, a *Article, tr *Translation, lang string, c Correction) (string, string, string) {
	reject := func(id, key string, args ...any) (string, string, string) {
		reason := T(lang, key)
		if id != "" {
			m.recordCorrection(ctx, id, CorrFailed, "", reason)
		}
		return CorrFailed, "", reason
	}

	// The rate limit is checked before the row is written, so a client hammering
	// the endpoint does not fill the table with its own refusals.
	if n, err := m.corrections.CountSince(ctx, c.Reporter, time.Now().Add(-24*time.Hour)); err == nil && n >= correctionsPerDay {
		return CorrFailed, "", T(lang, "corr.e_rate")
	}
	// The same word reported twice by the same reader is answered from the first
	// decision instead of costing a second model call.
	if st, found, _ := m.corrections.Duplicate(ctx, a.ID, c.Reporter, lang, c.Word); found && st == CorrApplied {
		return CorrApplied, "", T(lang, "corr.already")
	}

	id, err := m.corrections.Insert(ctx, c)
	if err != nil {
		m.rt.Logger.Error("correction insert", zap.Error(err))
		return CorrFailed, "", T(lang, "corr.e_internal")
	}
	sid := id.String()

	// Tier 1 — structure. Where in the article did they mean?
	site, ok := findCorrectionSite(tr.BodyMD, c.Sentence, c.Word)
	if !ok {
		return reject(sid, "corr.e_not_found")
	}

	// Tier 2 — mechanics. A Latin letter wearing a Cyrillic face needs no model.
	if fixed, ok := repairHomoglyphs(tr.BodyMD[site.Start:site.End]); ok {
		if err := safeFix(tr.BodyMD[site.Start:site.End], fixed); err == nil {
			return m.commitCorrection(ctx, sid, a, tr, lang, site, fixed)
		}
	}

	// Tier 3 — judgement.
	if m.ai == nil {
		return m.holdCorrection(ctx, sid, lang)
	}
	word := tr.BodyMD[site.Start:site.End]
	v, err := m.ai.Proofread(ctx, lang, site.Para, c.Sentence, word)
	switch {
	case errors.Is(err, ai.ErrDisabled):
		return m.holdCorrection(ctx, sid, lang)
	case err != nil:
		m.rt.Logger.Warn("proofread", zap.Error(err))
		return m.holdCorrection(ctx, sid, lang)
	}
	if !v.Fix {
		reason := v.Reason
		if reason == "" {
			reason = T(lang, "corr.r_no_error")
		}
		m.recordCorrection(ctx, sid, CorrRejected, "", reason)
		return CorrRejected, "", reason
	}

	fixed := matchCase(word, v.Fixed)
	if err := safeFix(word, fixed); err != nil {
		// The model answered with something that is not a typo fix. That is our
		// problem, not the reader's, so it is logged as such and they are told
		// the claim needs a person.
		m.rt.Logger.Warn("proofread fix rejected by guard",
			zap.String("word", word), zap.String("fixed", v.Fixed), zap.Error(err))
		return m.holdCorrection(ctx, sid, lang)
	}

	// The second reading, and only on the path that would change the text.
	stands, why, err := m.ai.ProofreadRefute(ctx, lang, site.Para, c.Sentence, word, fixed)
	if err != nil {
		m.rt.Logger.Warn("proofread refute", zap.Error(err))
		return m.holdCorrection(ctx, sid, lang)
	}
	if !stands {
		reason := why
		if reason == "" {
			reason = T(lang, "corr.r_no_error")
		}
		m.recordCorrection(ctx, sid, CorrRejected, "", reason)
		return CorrRejected, "", reason
	}

	return m.commitCorrection(ctx, sid, a, tr, lang, site, fixed)
}

// commitCorrection writes the one-word edit and tells the author about it.
func (m *Module) commitCorrection(ctx context.Context, id string, a *Article, tr *Translation, lang string, site correctionSite, fixed string) (string, string, string) {
	was := tr.BodyMD
	word := was[site.Start:site.End]
	updated := applyAt(was, site, fixed)

	okWrite, err := m.corrections.ApplyToBody(ctx, a.ID, lang, was, updated)
	if err != nil {
		m.rt.Logger.Error("correction apply", zap.Error(err))
		m.recordCorrection(ctx, id, CorrPending, fixed, "")
		return CorrPending, fixed, T(lang, "corr.r_held")
	}
	if !okWrite {
		// The article moved under us. Nothing was written, and saying "done"
		// would be a lie in the one place this feature cannot afford one.
		m.recordCorrection(ctx, id, CorrPending, fixed, "")
		return CorrPending, fixed, T(lang, "corr.r_stale")
	}
	m.recordCorrection(ctx, id, CorrApplied, fixed, "")
	m.notifyAuthorOfCorrection(ctx, a, lang, word, fixed, site)
	return CorrApplied, fixed, ""
}

// holdCorrection leaves a claim undecided and says so plainly. A reader told
// "we'll look at it" and meaning it is better served than one told "rejected" by
// a checker that never ran.
func (m *Module) holdCorrection(ctx context.Context, id, lang string) (string, string, string) {
	m.recordCorrection(ctx, id, CorrPending, "", "")
	return CorrPending, "", T(lang, "corr.r_held")
}

// recordCorrection stores an outcome, logging rather than surfacing a failure to
// write it: the decision has already been acted on, and the reader is owed the
// answer either way.
func (m *Module) recordCorrection(ctx context.Context, id, status, fixed, reason string) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return
	}
	if err := m.corrections.Decide(ctx, uid, status, fixed, reason); err != nil {
		m.rt.Logger.Error("correction decide", zap.Error(err))
	}
}

// notifyAuthorOfCorrection tells the author what a reader changed in their text.
//
// Not a courtesy. Their name is on the piece, the edit happened without them, and
// an author who cannot see what was changed cannot undo it.
func (m *Module) notifyAuthorOfCorrection(ctx context.Context, a *Article, lang, word, fixed string, site correctionSite) {
	if m.mailer == nil || strings.TrimSpace(a.AuthorEmail) == "" {
		return
	}
	url := m.rt.Config.PublicBaseURL + "/read/" + a.Slug + "?lang=" + lang
	subject := T(a.OriginalLang, "corr.mail_subject")
	body := strings.Join([]string{
		fmtT(a.OriginalLang, "corr.mail_intro", url),
		"",
		fmtT(a.OriginalLang, "corr.mail_change", word, fixed),
		"",
		T(a.OriginalLang, "corr.mail_where") + " " + site.Para,
		"",
		T(a.OriginalLang, "corr.mail_undo"),
	}, "\n")
	if err := m.mailer.Send(ctx, a.AuthorEmail, subject, body); err != nil {
		m.rt.Logger.Warn("correction mail", zap.Error(err))
	}
}
