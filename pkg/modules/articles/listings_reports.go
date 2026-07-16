package articles

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// reportAutoHideThreshold is how many distinct users must report a listing
// before it is auto-hidden for review.
const reportAutoHideThreshold = 3

// Flagged reports whether a listing was hidden pending review (e.g. after
// enough reports of misleading photos).
func (l Listing) Flagged() bool { return l.Status == "flagged" }

// RecordView bumps a listing's view counter (best-effort; errors are ignored
// by callers so a display path never fails on a counter write).
func (s *ListingStore) RecordView(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `UPDATE listings SET views_count = views_count + 1 WHERE id = $1`, id)
	return err
}

// RecordContact bumps the counter of how many times the seller's contact was
// revealed on a listing.
func (s *ListingStore) RecordContact(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `UPDATE listings SET contacts_count = contacts_count + 1 WHERE id = $1`, id)
	return err
}

// Report records one user's report of a listing and, once enough distinct
// users have reported it, hides the still-published listing for review. It
// returns the current report count and whether this call just auto-hid it.
func (s *ListingStore) Report(ctx context.Context, listingID, reporterID uuid.UUID, reason string) (count int, hidden bool, err error) {
	if reason == "" {
		reason = "misleading_photos"
	}
	if _, err = s.db.Exec(ctx,
		`INSERT INTO listing_reports (listing_id, reporter_id, reason) VALUES ($1,$2,$3)
		 ON CONFLICT (listing_id, reporter_id) DO NOTHING`,
		listingID, reporterID, reason); err != nil {
		return 0, false, fmt.Errorf("insert report: %w", err)
	}
	if err = s.db.QueryRow(ctx,
		`SELECT count(*) FROM listing_reports WHERE listing_id=$1`, listingID).Scan(&count); err != nil {
		return 0, false, fmt.Errorf("count reports: %w", err)
	}
	if count >= reportAutoHideThreshold {
		tag, uerr := s.db.Exec(ctx,
			`UPDATE listings SET status='flagged' WHERE id=$1 AND status='published'`, listingID)
		if uerr != nil {
			return count, false, fmt.Errorf("flag listing: %w", uerr)
		}
		hidden = tag.RowsAffected() > 0
	}
	return count, hidden, nil
}
