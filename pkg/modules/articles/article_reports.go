package articles

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"shanraq.org/pkg/modules/ratings"
)

// Reader moderation of articles. With no checker standing between an author and
// publication, this is the whole of moderation until a person intervenes — so
// it has to be hard to abuse in both directions: hard for a handful of annoyed
// readers to bury a piece they disagree with, and hard for a piece that breaks
// the rules to stay up because nobody with standing noticed.
//
// Four rules do that work:
//
//  1. A report is not a downvote. Disagreement has a button of its own; this one
//     claims a rule was broken, and it is counted separately.
//  2. Reports are weighted by the reporter's standing, so an article cannot be
//     hidden by accounts created for the purpose. A new account weighs one.
//  3. The threshold is a share of the actual readership, not a fixed number, and
//     it does not apply at all until enough people have read the piece. Three
//     readers cannot hide an article nobody has seen yet.
//  4. Being wrong costs something. When a hidden article is restored, every
//     report that hid it is marked dismissed, and a reader with a habit of
//     dismissed reports stops moving the threshold at all.
const (
	// articleReportFloor is the minimum number of distinct people. No weighting
	// and no share can hide an article below it: brigading starts with one
	// person and several accounts, and this is the floor under that.
	articleReportFloor = 3

	// articleReportMinViews is the audience an article must have reached before
	// readers can hide it automatically. Below it reports are still recorded and
	// still visible to staff — they simply do not fire the automatic hide, so a
	// piece cannot be buried in the minutes after publication, before anyone
	// who might defend it has read it. Anything urgent below this line is a
	// staff hide, which is immediate and does not wait for an audience.
	articleReportMinViews = 100

	// articleReportPercent is the weighted report total as a share of readers.
	// At five percent the floor of three people is not on its own enough to hide
	// anything: three accounts made for the purpose weigh three, and a hundred
	// readers require five.
	articleReportPercent = 5

	// articleReportHard hides regardless of audience share, for the case where
	// a very widely read article is genuinely and obviously in breach.
	articleReportHard = 25

	// falseReportsIgnored is how many dismissed reports a reader may accumulate
	// before their reports stop counting toward the automatic threshold. They
	// can still report, and staff still see it; it just no longer moves alone.
	falseReportsIgnored = 3
)

// reporterWeight is the standing a report carries: the reader's karma weight,
// reduced by each previous report of theirs that was dismissed.
//
// Two mistakes are tolerated — anyone can misread a piece — and a pattern is
// not: at the third, the reader's reports weigh nothing automatically.
func reporterWeight(karma, dismissed int) int {
	if dismissed >= falseReportsIgnored {
		return 0
	}
	w := ratings.Weight(karma) - dismissed
	if w < 1 {
		w = 1
	}
	return w
}

// shouldHideArticle applies the policy above. weighted is the summed standing of
// the reports, distinct the number of people behind them, views the readership.
func shouldHideArticle(weighted, distinct, views int) bool {
	if distinct < articleReportFloor || views < articleReportMinViews {
		return false
	}
	if weighted >= articleReportHard {
		return true
	}
	return weighted*100 >= articleReportPercent*views
}

// ErrOwnArticle is returned when an author reports their own piece.
var ErrOwnArticle = errors.New("cannot report your own article")

// ArticleReportResult is what one report did.
type ArticleReportResult struct {
	Distinct int  // people who have reported this article
	Weighted int  // their summed standing
	Hidden   bool // this report is the one that crossed the threshold
	Author   uuid.UUID
	Title    string
}

// ReportArticle records one reader's report and hides the article if that
// tipped it past the policy. Everything happens in one transaction: the count
// the decision rests on cannot change between reading it and acting on it.
func (s *Store) ReportArticle(ctx context.Context, articleID, reporterID uuid.UUID, reason, lang string) (ArticleReportResult, error) {
	var out ArticleReportResult
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return out, fmt.Errorf("report: begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var status string
	var views int
	if err := tx.QueryRow(ctx,
		`SELECT author_id, status, COALESCE(views_count, 0) FROM articles WHERE id = $1`,
		articleID).Scan(&out.Author, &status, &views); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, ErrNotFound
		}
		return out, fmt.Errorf("report: load article: %w", err)
	}
	if out.Author == reporterID {
		return out, ErrOwnArticle
	}

	var karma, dismissed int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(karma, 0) FROM author_reputation WHERE user_id = $1`,
		reporterID).Scan(&karma); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return out, fmt.Errorf("report: reporter karma: %w", err)
	}
	if err := tx.QueryRow(ctx,
		`SELECT count(*) FROM article_reports WHERE reporter_id = $1 AND dismissed`,
		reporterID).Scan(&dismissed); err != nil {
		return out, fmt.Errorf("report: reporter history: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO article_reports (article_id, reporter_id, reason, weight)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (article_id, reporter_id) DO UPDATE SET reason = EXCLUDED.reason`,
		articleID, reporterID, clip(reason, 60), reporterWeight(karma, dismissed)); err != nil {
		return out, fmt.Errorf("report: insert: %w", err)
	}

	// Dismissed reports are excluded from both totals: a decision already found
	// to be wrong must not help make the next one.
	if err := tx.QueryRow(ctx, `
		SELECT count(*), COALESCE(sum(weight), 0)
		  FROM article_reports WHERE article_id = $1 AND NOT dismissed`,
		articleID).Scan(&out.Distinct, &out.Weighted); err != nil {
		return out, fmt.Errorf("report: count: %w", err)
	}

	if status == "published" && shouldHideArticle(out.Weighted, out.Distinct, views) {
		ct, err := tx.Exec(ctx,
			`UPDATE articles SET status = 'flagged', updated_at = NOW() WHERE id = $1 AND status = 'published'`,
			articleID)
		if err != nil {
			return out, fmt.Errorf("report: flag: %w", err)
		}
		out.Hidden = ct.RowsAffected() > 0
	}

	if err := tx.QueryRow(ctx,
		`SELECT COALESCE(title, '') FROM article_translations WHERE article_id = $1 AND lang = $2`,
		articleID, lang).Scan(&out.Title); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return out, fmt.Errorf("report: title: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return out, fmt.Errorf("report: commit: %w", err)
	}
	return out, nil
}

// HasReportedArticle reports whether this reader has already reported it, so the
// page can say so instead of offering the button again.
func (s *Store) HasReportedArticle(ctx context.Context, articleID, readerID uuid.UUID) bool {
	var ok bool
	if err := s.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM article_reports WHERE article_id = $1 AND reporter_id = $2)`,
		articleID, readerID).Scan(&ok); err != nil {
		return false
	}
	return ok
}

// RestoreArticle brings a reader-hidden article back and marks the reports that
// hid it dismissed — which is what makes a false report cost its author
// something. Returns false when the article was not hidden.
//
// This is the other half of the policy. Without it, hiding is free to attempt
// and free to get wrong, and a determined group can keep trying forever.
func (s *Store) RestoreArticle(ctx context.Context, articleID uuid.UUID) (bool, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, fmt.Errorf("restore: begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	ct, err := tx.Exec(ctx,
		`UPDATE articles SET status = 'published', updated_at = NOW() WHERE id = $1 AND status = 'flagged'`,
		articleID)
	if err != nil {
		return false, fmt.Errorf("restore: status: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return false, tx.Commit(ctx)
	}
	if _, err := tx.Exec(ctx,
		`UPDATE article_reports SET dismissed = TRUE WHERE article_id = $1 AND NOT dismissed`,
		articleID); err != nil {
		return false, fmt.Errorf("restore: dismiss reports: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("restore: commit: %w", err)
	}
	return true, nil
}
