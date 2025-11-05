-- +goose Up
CREATE TABLE IF NOT EXISTS auth_mfa_totp (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    secret TEXT NOT NULL,
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auth_mfa_totp_user_unique UNIQUE (user_id)
);

CREATE INDEX IF NOT EXISTS auth_mfa_totp_user_id_idx ON auth_mfa_totp (user_id);

-- +goose Down
DROP TABLE IF EXISTS auth_mfa_totp;
