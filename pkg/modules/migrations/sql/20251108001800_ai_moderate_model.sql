-- +goose Up
-- One model setting served three jobs of different shapes. The moderation check
-- reads every article and answers with a short JSON verdict — a classification,
-- and the one call whose volume scales with publishing. Drafting a column and
-- improving an author's text are long-form writing, done rarely and on demand.
-- Tying them to one identifier means either paying frontier prices to classify,
-- or writing columns with the cheap model.
--
-- Empty falls back to the translation model, which is the other cheap-tier job:
-- the sensible default for a classifier, and overridable in the admin panel.
ALTER TABLE ai_settings ADD COLUMN IF NOT EXISTS moderate_model TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE ai_settings DROP COLUMN IF EXISTS moderate_model;
