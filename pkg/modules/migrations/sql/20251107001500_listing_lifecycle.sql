-- +goose Up
-- Free tier: a listing is visible for 14 days; the owner can extend it, and
-- (later, for a fee) promote it to the top or highlight it. During the launch
-- period these actions are free.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS expires_at      TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '14 days');
ALTER TABLE listings ADD COLUMN IF NOT EXISTS promoted_until  TIMESTAMPTZ;
ALTER TABLE listings ADD COLUMN IF NOT EXISTS featured_until  TIMESTAMPTZ;
-- Set true once the "expiring soon" reminder e-mail has been sent; reset on extend.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS expiry_reminded BOOLEAN NOT NULL DEFAULT false;

-- Keep existing (demo) listings visible for 14 days from now.
UPDATE listings SET expires_at = NOW() + INTERVAL '14 days';

-- Illustrate promotion/highlight on a few demo listings.
UPDATE listings SET promoted_until = NOW() + INTERVAL '7 days'
WHERE id IN ('a0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000009');
UPDATE listings SET featured_until = NOW() + INTERVAL '7 days'
WHERE id IN ('a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000011');

CREATE INDEX IF NOT EXISTS idx_listings_active ON listings (status, expires_at);

-- +goose Down
DROP INDEX IF EXISTS idx_listings_active;
ALTER TABLE listings DROP COLUMN IF EXISTS expiry_reminded;
ALTER TABLE listings DROP COLUMN IF EXISTS featured_until;
ALTER TABLE listings DROP COLUMN IF EXISTS promoted_until;
ALTER TABLE listings DROP COLUMN IF EXISTS expires_at;
