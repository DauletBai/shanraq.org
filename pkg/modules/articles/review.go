package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/jobs"
)

// ReviewRules are the publication rules an article is checked against. The
// codes are stable and translated, so an author gets the same complaint in
// their own language rather than whatever the checker happened to write.
//
// They are drawn from the platform's own published rules, not invented here:
// an author answers for the accuracy and lawfulness of their material, facts
// should be backed by sources where possible, and opinion pieces must be
// marked as opinion.
var ReviewRules = []string{
	"unsourced_claim", // factual assertions with no source
	"opinion_as_fact", // judgment presented as established fact
	"defamation",      // unproven accusations against identifiable people
	"hatred",          // incitement to violence or enmity
	"personal_data",   // someone else's personal data without consent
	"illegal",         // content unlawful under the law of the RK
	"disguised_ad",    // advertising presented as editorial
	"plagiarism",      // someone else's text without attribution
	"title_mismatch",  // headline not supported by the body
	"unreadable",      // incoherent or machine-dumped text
}

func isReviewRule(code string) bool {
	for _, r := range ReviewRules {
		if r == code {
			return true
		}
	}
	return false
}

// Finding is one rule an article failed, with the passage it refers to.
type Finding struct {
	RuleCode string `json:"rule"`
	Severity string `json:"severity"` // block | warn
	Quote    string `json:"quote"`
	Note     string `json:"note"`
}

// Blocking reports whether this finding stops publication.
func (f Finding) Blocking() bool { return f.Severity == "block" }

// reviewVerdict is what the checker is asked to return. Keeping the contract
// this narrow is what makes the answer machine-checkable rather than prose.
type reviewVerdict struct {
	Findings []Finding `json:"findings"`
}

// reviewPrompt builds the checker's instructions. It states plainly that the
// checker must quote the passage it objects to: a complaint an author cannot
// locate is not a complaint, it is an obstacle.
func reviewPrompt(lang, title, summary, body string) string {
	var b strings.Builder
	b.WriteString("You are the publication checker for a Kazakhstani trilingual platform. ")
	b.WriteString("Check the article below against the rules and return JSON only.\n\n")
	b.WriteString("Rules, by code:\n")
	b.WriteString("- unsourced_claim: a factual assertion (numbers, events, quotes) with no source named in the text\n")
	b.WriteString("- opinion_as_fact: the author's judgment written as established fact\n")
	b.WriteString("- defamation: an accusation against an identifiable person or organisation without evidence\n")
	b.WriteString("- hatred: incitement to violence or to enmity against a group\n")
	b.WriteString("- personal_data: another person's personal data published without their consent\n")
	b.WriteString("- illegal: content unlawful under the law of the Republic of Kazakhstan\n")
	b.WriteString("- disguised_ad: promotional material presented as editorial\n")
	b.WriteString("- plagiarism: substantial text that appears to be someone else's without attribution\n")
	b.WriteString("- title_mismatch: the headline claims more than the body supports\n")
	b.WriteString("- unreadable: incoherent text, or raw machine output pasted in\n\n")
	b.WriteString("Severity: \"block\" stops publication; \"warn\" is advice that does not.\n")
	b.WriteString("Use \"block\" only for defamation, hatred, personal_data, illegal, plagiarism, disguised_ad, ")
	b.WriteString("and for unreadable text. For unsourced_claim, opinion_as_fact and title_mismatch use \"block\" ")
	b.WriteString("only when the problem runs through the whole piece; otherwise \"warn\".\n\n")
	b.WriteString("For every finding you MUST quote the exact passage you object to in \"quote\". ")
	b.WriteString("A finding the author cannot locate is useless. Write \"note\" in the language of the article (")
	b.WriteString(lang)
	b.WriteString("), addressed to the author, saying what to change.\n\n")
	b.WriteString("Do not invent problems. An article that breaks no rule must return an empty findings array. ")
	b.WriteString("Disagreeing with the author's opinion is not a rule violation — this platform publishes opinion.\n\n")
	b.WriteString("Return exactly: {\"findings\":[{\"rule\":\"...\",\"severity\":\"block|warn\",\"quote\":\"...\",\"note\":\"...\"}]}\n\n")
	b.WriteString("TITLE: ")
	b.WriteString(title)
	b.WriteString("\nSUMMARY: ")
	b.WriteString(summary)
	b.WriteString("\nBODY:\n")
	b.WriteString(clip(body, 24000))
	return b.String()
}

// parseVerdict reads the checker's answer. Anything it cannot parse is an
// error rather than an empty pass — a malformed reply must not be mistaken for
// a clean article.
func parseVerdict(raw string) ([]Finding, error) {
	s := strings.TrimSpace(raw)
	if i := strings.Index(s, "{"); i > 0 {
		s = s[i:]
	}
	if j := strings.LastIndex(s, "}"); j >= 0 && j < len(s)-1 {
		s = s[:j+1]
	}
	var v reviewVerdict
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, fmt.Errorf("checker returned unparseable output: %w", err)
	}
	out := make([]Finding, 0, len(v.Findings))
	for _, f := range v.Findings {
		if !isReviewRule(f.RuleCode) {
			// An unknown code has no translation, so the author would see
			// nothing useful. Drop it rather than show a blank complaint.
			continue
		}
		if f.Severity != "warn" {
			f.Severity = "block"
		}
		f.Quote = clip(strings.TrimSpace(f.Quote), 300)
		f.Note = clip(strings.TrimSpace(f.Note), 600)
		out = append(out, f)
	}
	return out, nil
}

// SaveFindings attaches the checker's findings to a ledger entry.
func (s *ModStore) SaveFindings(ctx context.Context, actionID uuid.UUID, fs []Finding) error {
	for _, f := range fs {
		if _, err := s.db.Exec(ctx, `
			INSERT INTO moderation_findings (action_id, rule_code, severity, quote, note)
			VALUES ($1,$2,$3,$4,$5)`, actionID, f.RuleCode, f.Severity, f.Quote, f.Note); err != nil {
			return fmt.Errorf("save findings: %w", err)
		}
	}
	return nil
}

// FindingsFor loads the findings of one ledger entry, so the author sees every
// rule they failed rather than the first.
func (s *ModStore) FindingsFor(ctx context.Context, actionID string) ([]Finding, error) {
	rows, err := s.db.Query(ctx, `
		SELECT rule_code, severity, quote, note FROM moderation_findings
		WHERE action_id = $1::uuid ORDER BY severity, created_at`, actionID)
	if err != nil {
		return nil, fmt.Errorf("findings: %w", err)
	}
	defer rows.Close()
	out := []Finding{}
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.RuleCode, &f.Severity, &f.Quote, &f.Note); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// submitForReview is what the publish button now does. The article goes to
// 'review', the checker runs, and the article either publishes or comes back
// with findings. Nothing publishes without a ledger entry saying who cleared
// it — including when the clearer was a machine.
func (m *Module) submitForReview(ctx context.Context, id, author uuid.UUID, lang string) (published bool, blocking int, err error) {
	if err := m.store.SetStatus(ctx, id, author, "review"); err != nil {
		return false, 0, err
	}
	if _, err := m.rt.DB.Exec(ctx, `UPDATE articles SET submitted_at = NOW() WHERE id = $1`, id); err != nil {
		m.rt.Logger.Warn("stamp submitted_at", zap.Error(err))
	}

	_, tr, err := m.store.ForReview(ctx, id, lang)
	if err != nil {
		return false, 0, err
	}

	raw, cerr := m.ai.Check(ctx,
		"You are a strict but fair publication checker. Return JSON only.",
		reviewPrompt(lang, tr.Title, tr.Summary, tr.Body), 3000)
	if cerr != nil {
		// The checker is unavailable — the assistant is off, or the provider
		// failed. Hold the article for a human rather than publishing it
		// unchecked or throwing the author's work away. Fail closed, not lost.
		m.rt.Logger.Warn("publication check unavailable", zap.Error(cerr))
		m.logReview(ctx, id, author, tr.Title, "warn", "checker_unavailable", "agent", nil)
		return false, 0, nil
	}

	findings, perr := parseVerdict(raw)
	if perr != nil {
		// An unparseable reply is not a pass. Same treatment.
		m.rt.Logger.Warn("publication check unparseable", zap.Error(perr))
		m.logReview(ctx, id, author, tr.Title, "warn", "checker_unavailable", "agent", nil)
		return false, 0, nil
	}

	for _, f := range findings {
		if f.Blocking() {
			blocking++
		}
	}

	action, reason := "approve", "rules_ok"
	status := "published"
	if blocking > 0 {
		action, reason, status = "reject", "rules_failed", "needs_work"
	}
	// Status, ledger entry and findings commit as one unit. Before, a status
	// could move without its record, so an article might go live with nothing
	// in the log; that is the divergence the review flagged.
	if err := m.commitReview(ctx, id, author, tr.Title, status, action, reason, findings); err != nil {
		return false, blocking, err
	}
	return blocking == 0, blocking, nil
}

// commitReview writes the automated decision atomically: article status, the
// ledger entry, and every finding, in a single transaction.
func (m *Module) commitReview(ctx context.Context, id, author uuid.UUID, title, status, action, reason string, findings []Finding) error {
	tx, err := m.rt.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	pub := "published_at = COALESCE(articles.published_at, NOW()), "
	if status != "published" {
		pub = ""
	}
	if _, err := tx.Exec(ctx, `UPDATE articles SET status = $2, `+pub+`reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND author_id = $3`, id, status, author); err != nil {
		return fmt.Errorf("commit review: status: %w", err)
	}
	var actionID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO moderation_actions
		    (target_type, target_id, subject_id, title, action, reason_code, actor_kind, actor_name)
		VALUES ('article', $1, $2, $3, $4, $5, 'agent', 'AI Bake') RETURNING id`,
		id.String(), author, clip(title, 120), action, reason).Scan(&actionID); err != nil {
		return fmt.Errorf("commit review: log: %w", err)
	}
	for _, f := range findings {
		if _, err := tx.Exec(ctx, `
			INSERT INTO moderation_findings (action_id, rule_code, severity, quote, note)
			VALUES ($1,$2,$3,$4,$5)`, actionID, f.RuleCode, f.Severity, f.Quote, f.Note); err != nil {
			return fmt.Errorf("commit review: finding: %w", err)
		}
	}
	return tx.Commit(ctx)
}

// logReview writes the decision and its findings to the moderation ledger, so
// an automated approval is as auditable and appealable as a rejection.
func (m *Module) logReview(ctx context.Context, id, author uuid.UUID, title, action, reason, actorKind string, findings []Finding) {
	actionID, err := m.mods.Log(ctx, ModAction{
		TargetType: "article", TargetID: id.String(), Title: clip(title, 120),
		Action: action, ReasonCode: reason, ActorKind: actorKind, ActorName: "AI Bake",
	}, &author, nil)
	if err != nil {
		m.rt.Logger.Error("log review", zap.Error(err))
		return
	}
	if len(findings) == 0 {
		return
	}
	aid, err := uuid.Parse(actionID)
	if err != nil {
		return
	}
	if err := m.mods.SaveFindings(ctx, aid, findings); err != nil {
		m.rt.Logger.Error("save findings", zap.Error(err))
	}
}

// reviewText is the article text the checker reads.
type reviewText struct{ Title, Summary, Body string }

// ForReview loads an article and the translation in its original language,
// which is the version the author actually wrote and the one worth checking.
func (s *Store) ForReview(ctx context.Context, id uuid.UUID, lang string) (Article, reviewText, error) {
	var a Article
	var t reviewText
	err := s.db.QueryRow(ctx, `
		SELECT a.id, a.slug, a.category, a.original_lang,
		       COALESCE(t.title,''), COALESCE(t.summary,''), COALESCE(t.body_md,'')
		  FROM articles a
		  LEFT JOIN article_translations t
		    ON t.article_id = a.id AND t.lang = a.original_lang
		 WHERE a.id = $1`, id).Scan(&a.ID, &a.Slug, &a.Category, &a.OriginalLang, &t.Title, &t.Summary, &t.Body)
	if err != nil {
		return a, t, fmt.Errorf("load for review: %w", err)
	}
	return a, t, nil
}

// JobColumnBrief is the on-demand analysis job. There is deliberately no
// schedule behind it: a nightly sweep costs money whether or not anything
// happened, so the run is started by the administrator and only then.
const JobColumnBrief = "articles.column_brief"

// EnqueueColumnBrief asks the columnist agent for an analysis of a given
// topic. Queued rather than run inline, because an LLM call does not belong in
// an HTTP handler.
func (m *Module) EnqueueColumnBrief(ctx context.Context, lang, brief string) error {
	if m.jobs == nil {
		return fmt.Errorf("job queue unavailable")
	}
	payload, err := json.Marshal(map[string]string{"lang": lang, "brief": brief})
	if err != nil {
		return err
	}
	return m.jobs.Enqueue(ctx, jobs.Job{
		ID:          uuid.New(),
		Name:        JobColumnBrief,
		Payload:     payload,
		RunAt:       time.Now(),
		MaxAttempts: 2,
	})
}

// ReviewItem is one article waiting on a human, for the admin queue.
type ReviewItem struct {
	ID          string
	Slug        string
	Title       string
	AuthorEmail string
	Category    string
	Status      string // review | needs_work
	Submitted   time.Time
	Findings    []Finding
}

// ReviewQueue lists articles a person still has to decide: those held because
// the checker was unavailable ('review'), and those the checker returned that
// the author resubmitted. Oldest first, so nothing waits forever.
func (s *Store) ReviewQueue(ctx context.Context, limit int) ([]ReviewItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.db.Query(ctx, `
		SELECT a.id, a.slug, a.status, COALESCE(a.submitted_at, a.updated_at),
		       COALESCE(u.email,''), a.category,
		       COALESCE((SELECT t.title FROM article_translations t
		                  WHERE t.article_id = a.id AND t.lang = a.original_lang), '')
		  FROM articles a
		  LEFT JOIN auth_users u ON u.id = a.author_id
		 WHERE a.status = 'review'
		 ORDER BY a.submitted_at ASC NULLS LAST
		 LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("review queue: %w", err)
	}
	defer rows.Close()
	out := []ReviewItem{}
	for rows.Next() {
		var it ReviewItem
		if err := rows.Scan(&it.ID, &it.Slug, &it.Status, &it.Submitted,
			&it.AuthorEmail, &it.Category, &it.Title); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// DecideArticle is a human's ruling on a queued article. Status change and
// ledger entry commit together: the review's central defect was that a status
// could move without a matching record, so an article might publish with no
// audit trail. Here they cannot diverge.
func (m *Module) DecideArticle(ctx context.Context, id uuid.UUID, decision string, moderator uuid.UUID, note string) error {
	var status, action, reason string
	switch decision {
	case "approve":
		status, action, reason = "published", "approve", "human_approved"
	case "reject":
		status, action, reason = "needs_work", "reject", "human_rejected"
	case "needs_work":
		status, action, reason = "needs_work", "warn", "human_returned"
	default:
		return fmt.Errorf("unknown decision %q", decision)
	}

	tx, err := m.rt.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var author uuid.UUID
	var title, lang string
	if err := tx.QueryRow(ctx, `
		SELECT a.author_id, a.original_lang,
		       COALESCE((SELECT t.title FROM article_translations t
		                  WHERE t.article_id = a.id AND t.lang = a.original_lang),'')
		  FROM articles a WHERE a.id = $1 AND a.status IN ('review','needs_work')`,
		id).Scan(&author, &lang, &title); err != nil {
		return fmt.Errorf("decide: load article: %w", err)
	}

	pub := "published_at = COALESCE(articles.published_at, NOW()), "
	if status != "published" {
		pub = ""
	}
	if _, err := tx.Exec(ctx, `UPDATE articles SET status = $2, `+pub+`reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $1`, id, status); err != nil {
		return fmt.Errorf("decide: set status: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO moderation_actions
		    (target_type, target_id, subject_id, title, action, reason_code, reason_note, actor_kind, actor_id, actor_name)
		VALUES ('article', $1, $2, $3, $4, $5, $6, 'human', $7, '')`,
		id.String(), author, clip(title, 120), action, reason, clip(note, 500), moderator); err != nil {
		return fmt.Errorf("decide: log: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	// Syndication is a side effect, not part of the decision's integrity, so it
	// runs after the commit and its failure does not undo the ruling.
	if status == "published" && m.syndicate != nil {
		if err := m.syndicate.EnqueuePublish(ctx, m.jobs, id); err != nil {
			m.rt.Logger.Warn("enqueue publish after decision", zap.Error(err))
		}
	}
	return nil
}
