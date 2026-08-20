-- +goose Up
-- The AI switch was one switch for three different things: the writing
-- assistant, the trilingual translation, and the pre-publication rules check.
-- Turning it off to stop standing between an author and publication also turned
-- off the translation that makes this a trilingual publication at all — and
-- turning it on to get the translation back reinstated the gate.
--
-- They are separate questions. review_check governs only the gate, and it is off
-- by default: the site's model is that the author answers for the text and
-- readers moderate it afterwards. An administrator who wants a machine to read
-- every article before anyone else does can say so here.
ALTER TABLE ai_settings ADD COLUMN IF NOT EXISTS review_check BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE ai_settings DROP COLUMN IF EXISTS review_check;
