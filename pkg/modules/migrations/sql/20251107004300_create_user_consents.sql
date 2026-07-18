-- +goose Up
-- Consent history: proof that a user accepted the Terms and Privacy Policy at
-- registration, with the document version, channel, IP, and timestamp — as the
-- KZ online-platform law expects. Append-only; never updated.
CREATE TABLE IF NOT EXISTS user_consents (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    document    TEXT NOT NULL,            -- e.g. 'terms_privacy'
    version     TEXT NOT NULL,            -- document version accepted
    source      TEXT NOT NULL,            -- 'web' | 'api'
    ip          TEXT,
    accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS user_consents_user_idx ON user_consents(user_id);

-- +goose Down
DROP TABLE IF EXISTS user_consents;
