-- +goose Up
-- Automatic translation is switched off, and the switch is new because until
-- now there was none: translation came with the assistant, so keeping the
-- moderator meant keeping the translator.
--
-- The reason it is off is not that it failed. Rebuilt to send paragraphs in
-- verified batches and to refuse any version that changed a figure, it
-- translated this site's longest article into two languages with every
-- paragraph intact and every number correct. It also took eight minutes, spent
-- twenty requests against an account allowed three a minute, and still left a
-- Russian word sitting in the middle of a Kazakh sentence — the kind of fault
-- no counting rule will ever catch.
--
-- Meanwhile every author already pays for a model of their own and can do the
-- same job in a minute, reading the result in a language they know. A service
-- that is slower, costlier and worse than what the user already has is not a
-- service. So the machinery stays, unused and switched off, and can be turned
-- back on the day that calculation changes.
ALTER TABLE ai_settings ADD COLUMN IF NOT EXISTS auto_translate BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE ai_settings SET auto_translate = FALSE WHERE id = 1;

-- +goose Down
ALTER TABLE ai_settings DROP COLUMN IF EXISTS auto_translate;
