-- +goose Up
-- Third paid promotion for listings (besides "Top" and "highlight"): a banner
-- slot in the real-estate sidebar that stays visible while the user scrolls.
-- Sold by the day (1..7), priced above Top/highlight.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS banner_until TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_listings_banner ON listings (banner_until) WHERE banner_until IS NOT NULL;

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS banner_until;
