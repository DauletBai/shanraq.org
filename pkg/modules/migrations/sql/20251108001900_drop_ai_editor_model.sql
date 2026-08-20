-- +goose Up
-- The co-editor and the column drafter are gone: this site does not edit with a
-- model and does not write with one. Nothing reads editor_model any more, and a
-- settings column nobody reads is a setting somebody will one day fill in and
-- wonder why it has no effect.
--
-- Translation and moderation remain, and both are cheap-tier jobs served by
-- translate_model and moderate_model.
ALTER TABLE ai_settings DROP COLUMN IF EXISTS editor_model;

-- +goose Down
ALTER TABLE ai_settings ADD COLUMN IF NOT EXISTS editor_model TEXT NOT NULL DEFAULT 'claude-haiku-4-5';
