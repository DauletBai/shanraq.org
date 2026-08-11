-- +goose Up
-- Double opt-in for the weekly digest. Until confirmed_at is set the row is a
-- pending request, not a subscriber: nothing is ever sent to it. This is what
-- stops one person signing up somebody else's address, and it gives the sender
-- an immediate reply instead of up to a week of silence.
ALTER TABLE subscribers ADD COLUMN IF NOT EXISTS confirmed_at TIMESTAMPTZ;
ALTER TABLE subscribers ADD COLUMN IF NOT EXISTS confirm_token TEXT;

-- Rows collected before this migration came from an explicit form submit and
-- are already receiving the digest. Dropping them silently would be worse than
-- grandfathering them in, so they are treated as confirmed at signup time.
UPDATE subscribers SET confirmed_at = created_at WHERE confirmed_at IS NULL;

CREATE INDEX IF NOT EXISTS subscribers_confirm_token_idx
    ON subscribers (confirm_token) WHERE confirm_token IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS subscribers_confirm_token_idx;
ALTER TABLE subscribers DROP COLUMN IF EXISTS confirm_token;
ALTER TABLE subscribers DROP COLUMN IF EXISTS confirmed_at;
