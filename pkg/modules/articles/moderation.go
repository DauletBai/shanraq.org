package articles

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Moderation reason codes. They are stable identifiers, not prose: the wording
// shown to the author is a translation, so the same decision reads the same in
// all three languages.
var ModerationReasons = []string{"spam", "insult", "fraud", "personal_data", "illegal", "misleading_photos", "off_topic", "duplicate"}

// ModAction is one recorded moderation decision.
type ModAction struct {
	ID         string
	TargetType string // article | listing | comment
	TargetID   string
	Title      string
	Action     string // hide | restore | reject | approve | warn
	ReasonCode string
	ReasonNote string
	ActorKind  string // agent | human
	ActorName  string
	Created    time.Time

	// Appeal state, filled when the row is read with its appeal.
	AppealStatus string // "" when never appealed
	AppealBody   string
	Resolution   string
}

// Appealable reports whether the author may still contest this decision. A
// restore or an approval is not a penalty, so there is nothing to appeal.
func (a ModAction) Appealable() bool {
	if a.AppealStatus != "" {
		return false
	}
	return a.Action == "hide" || a.Action == "reject" || a.Action == "warn"
}

// ModAppeal is an author contesting a decision, as shown to a moderator.
type ModAppeal struct {
	ID         string
	ActionID   string
	Body       string
	Status     string
	Created    time.Time
	AuthorMail string
	// The decision being contested.
	TargetType string
	Title      string
	Action     string
	ReasonCode string
	ActorKind  string
	ActorName  string
}

// ModStore is the moderation ledger: every decision and every appeal.
type ModStore struct{ db *pgxpool.Pool }

func NewModStore(db *pgxpool.Pool) *ModStore { return &ModStore{db: db} }

// Log records a decision. It is deliberately the only way to hide something:
// an action that leaves no row is an action nobody can appeal.
func (s *ModStore) Log(ctx context.Context, a ModAction, subject *uuid.UUID, actor *uuid.UUID) (string, error) {
	var id string
	err := s.db.QueryRow(ctx, `
		INSERT INTO moderation_actions
		    (target_type, target_id, subject_id, title, action, reason_code, reason_note, actor_kind, actor_id, actor_name)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`,
		a.TargetType, a.TargetID, subject, a.Title, a.Action, a.ReasonCode, a.ReasonNote,
		a.ActorKind, actor, a.ActorName).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("log moderation: %w", err)
	}
	return id, nil
}

const modCols = `m.id, m.target_type, m.target_id, m.title, m.action, m.reason_code, m.reason_note,
	m.actor_kind, m.actor_name, m.created_at,
	COALESCE(ap.status,''), COALESCE(ap.body,''), COALESCE(ap.resolution,'')`

func scanActions(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]ModAction, error) {
	out := []ModAction{}
	for rows.Next() {
		var a ModAction
		if err := rows.Scan(&a.ID, &a.TargetType, &a.TargetID, &a.Title, &a.Action, &a.ReasonCode,
			&a.ReasonNote, &a.ActorKind, &a.ActorName, &a.Created,
			&a.AppealStatus, &a.AppealBody, &a.Resolution); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ForSubject returns the decisions taken about one author's content, newest
// first — this is what the author sees in their own cabinet.
func (s *ModStore) ForSubject(ctx context.Context, subject uuid.UUID, limit int) ([]ModAction, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.Query(ctx, fmt.Sprintf(`
		SELECT %s FROM moderation_actions m
		LEFT JOIN moderation_appeals ap ON ap.action_id = m.id
		WHERE m.subject_id = $1 ORDER BY m.created_at DESC LIMIT $2`, modCols), subject, limit)
	if err != nil {
		return nil, fmt.Errorf("moderation for subject: %w", err)
	}
	defer rows.Close()
	return scanActions(rows)
}

// Recent returns the newest decisions across the platform, for the admin log.
func (s *ModStore) Recent(ctx context.Context, limit int) ([]ModAction, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.Query(ctx, fmt.Sprintf(`
		SELECT %s FROM moderation_actions m
		LEFT JOIN moderation_appeals ap ON ap.action_id = m.id
		ORDER BY m.created_at DESC LIMIT $1`, modCols), limit)
	if err != nil {
		return nil, fmt.Errorf("recent moderation: %w", err)
	}
	defer rows.Close()
	return scanActions(rows)
}

// OpenAppeals is the moderator's actual work queue: contested decisions,
// oldest first, so nothing sits forgotten at the bottom.
func (s *ModStore) OpenAppeals(ctx context.Context, limit int) ([]ModAppeal, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.Query(ctx, `
		SELECT ap.id, ap.action_id, ap.body, ap.status, ap.created_at, COALESCE(u.email,''),
		       m.target_type, m.title, m.action, m.reason_code, m.actor_kind, m.actor_name
		  FROM moderation_appeals ap
		  JOIN moderation_actions m ON m.id = ap.action_id
		  LEFT JOIN auth_users u ON u.id = ap.author_id
		 WHERE ap.status = 'open'
		 ORDER BY ap.created_at ASC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("open appeals: %w", err)
	}
	defer rows.Close()
	out := []ModAppeal{}
	for rows.Next() {
		var a ModAppeal
		if err := rows.Scan(&a.ID, &a.ActionID, &a.Body, &a.Status, &a.Created, &a.AuthorMail,
			&a.TargetType, &a.Title, &a.Action, &a.ReasonCode, &a.ActorKind, &a.ActorName); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Appeal files an author's objection. The subject check is in the statement
// itself, so one user cannot contest a decision about another's content.
func (s *ModStore) Appeal(ctx context.Context, actionID, author uuid.UUID, body string) error {
	tag, err := s.db.Exec(ctx, `
		INSERT INTO moderation_appeals (action_id, author_id, body)
		SELECT $1, $2, $3 FROM moderation_actions m
		 WHERE m.id = $1 AND m.subject_id = $2
		   AND m.action IN ('hide','reject','warn')
		ON CONFLICT (action_id) DO NOTHING`, actionID, author, body)
	if err != nil {
		return fmt.Errorf("file appeal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("appeal not accepted")
	}
	return nil
}

// Resolve closes an appeal. Overturning one records the reversal as its own
// action, so the ledger shows the correction rather than quietly rewriting the
// original decision.
func (s *ModStore) Resolve(ctx context.Context, appealID, moderator uuid.UUID, uphold bool, note string) error {
	status := "overturned"
	if uphold {
		status = "upheld"
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var actionID, targetType, targetID, title string
	var subject *uuid.UUID
	err = tx.QueryRow(ctx, `
		UPDATE moderation_appeals SET status = $2, resolution = $3, resolved_by = $4, resolved_at = NOW()
		 WHERE id = $1 AND status = 'open'
		RETURNING action_id`, appealID, status, note, moderator).Scan(&actionID)
	if err != nil {
		return fmt.Errorf("resolve appeal: %w", err)
	}
	if err := tx.QueryRow(ctx, `
		SELECT target_type, target_id, title, subject_id FROM moderation_actions WHERE id = $1`,
		actionID).Scan(&targetType, &targetID, &title, &subject); err != nil {
		return err
	}
	if !uphold {
		if _, err := tx.Exec(ctx, `
			INSERT INTO moderation_actions
			    (target_type, target_id, subject_id, title, action, reason_code, reason_note, actor_kind, actor_id, actor_name)
			VALUES ($1,$2,$3,$4,'restore','appeal_upheld',$5,'human',$6,'')`,
			targetType, targetID, subject, title, note, moderator); err != nil {
			return err
		}
		switch targetType {
		case "comment":
			if _, err := tx.Exec(ctx, `UPDATE comments SET status = 'published' WHERE id = $1::uuid`, targetID); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

// isModerationReason guards the reason code so the ledger cannot be filled
// with free-text that no translation exists for.
func isModerationReason(code string) bool {
	for _, r := range ModerationReasons {
		if r == code {
			return true
		}
	}
	return false
}
