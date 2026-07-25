-- +goose Up
-- Runtime-switchable AI model settings, editable from the admin panel without a
-- redeploy. Single row (id = 1). NO API keys are stored here — secrets stay in
-- the server config/env per provider; this table only records which provider
-- and models are active.
CREATE TABLE IF NOT EXISTS ai_settings (
    id              SMALLINT     PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    enabled         BOOLEAN      NOT NULL DEFAULT FALSE,
    provider        TEXT         NOT NULL DEFAULT 'anthropic',
    editor_model    TEXT         NOT NULL DEFAULT 'claude-sonnet-5',
    translate_model TEXT         NOT NULL DEFAULT 'claude-haiku-4-5',
    max_tokens      INTEGER      NOT NULL DEFAULT 4096,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by      UUID
);

-- +goose Down
DROP TABLE IF EXISTS ai_settings;
