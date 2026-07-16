-- +goose Up
-- Allow the 'flagged' status (listing auto-hidden pending review after reports).
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_status_chk;
ALTER TABLE listings ADD CONSTRAINT listings_status_chk
    CHECK (status IN ('published', 'hidden', 'flagged'));

-- +goose Down
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_status_chk;
ALTER TABLE listings ADD CONSTRAINT listings_status_chk
    CHECK (status IN ('published', 'hidden'));
