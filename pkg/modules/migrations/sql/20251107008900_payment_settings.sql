-- +goose Up
-- Runtime-switchable payment acquirer selection (which provider is live, and
-- whether payments are on) in a single row, mirroring ai_settings. Secrets
-- (merchant IDs, API keys) stay in the server config — only the choice and the
-- on/off live here, so an operator connects/disconnects Kaspi from the admin
-- panel without a code change or redeploy.
CREATE TABLE IF NOT EXISTS payment_settings (
    id         INT PRIMARY KEY DEFAULT 1,
    enabled    BOOLEAN NOT NULL DEFAULT FALSE,
    provider   TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    CONSTRAINT payment_settings_singleton CHECK (id = 1)
);

-- +goose Down
DROP TABLE IF EXISTS payment_settings;
