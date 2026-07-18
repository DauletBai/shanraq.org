-- +goose Up
-- Email verification: a nullable verified-at stamp on the user plus a token
-- table. Actions that can be abused by throwaway accounts (posting listings,
-- reporting) require a verified email.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS email_verification_user_idx ON email_verification_tokens(user_id);

-- Grandfather every existing account (demo/showcase + any real users) as
-- verified so this rollout never locks anyone out. New accounts created after
-- this migration start unverified.
UPDATE auth_users SET email_verified_at = NOW() WHERE email_verified_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS email_verification_tokens;
ALTER TABLE auth_users DROP COLUMN IF EXISTS email_verified_at;
