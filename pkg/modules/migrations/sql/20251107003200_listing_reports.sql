-- +goose Up
-- Reader reports on listings — chiefly for photos edited with filters that
-- distort real dimensions (misleading the buyer). One report per user per
-- listing; enough distinct reports auto-hides the listing for review.
CREATE TABLE IF NOT EXISTS listing_reports (
    listing_id  uuid        NOT NULL,
    reporter_id uuid        NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    reason      text        NOT NULL DEFAULT 'misleading_photos',
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (listing_id, reporter_id)
);
CREATE INDEX IF NOT EXISTS idx_listing_reports_listing ON listing_reports (listing_id);

-- +goose Down
DROP TABLE IF EXISTS listing_reports;
