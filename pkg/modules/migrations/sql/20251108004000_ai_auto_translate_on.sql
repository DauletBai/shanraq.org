-- +goose Up
-- Turn automatic translation on.
--
-- It was left off while translating was a convenience an author's own model
-- did faster. Publication now requires the missing languages outright -- three
-- for a countrywide article, Kazakh and Russian at the least for a local one --
-- so an author without a second model needs a way through, and the eight
-- minutes the pipeline costs buy an article that would otherwise not run.
--
-- This turns it on for an install that already has the row; a fresh one seeds
-- itself from the defaults in ai.go, which say the same. Operators can still
-- switch it off in the admin panel.
UPDATE ai_settings SET auto_translate = TRUE WHERE id = 1;

-- +goose Down
UPDATE ai_settings SET auto_translate = FALSE WHERE id = 1;
