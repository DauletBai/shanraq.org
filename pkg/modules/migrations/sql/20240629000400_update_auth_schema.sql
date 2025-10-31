-- +goose Up
ALTER TABLE auth_users
    ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

ALTER TABLE auth_users
    ADD COLUMN IF NOT EXISTS password_reset_required BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE auth_users
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE TABLE IF NOT EXISTS auth_refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS auth_refresh_tokens_user_idx ON auth_refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS auth_refresh_tokens_expires_idx ON auth_refresh_tokens(expires_at);

CREATE TABLE IF NOT EXISTS auth_password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS auth_password_resets_user_idx ON auth_password_resets(user_id);
CREATE INDEX IF NOT EXISTS auth_password_resets_expires_idx ON auth_password_resets(expires_at);

-- +goose Down
DROP TABLE IF EXISTS auth_password_resets;
DROP TABLE IF EXISTS auth_refresh_tokens;
ALTER TABLE auth_users
    DROP COLUMN IF EXISTS password_reset_required,
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS updated_at;
