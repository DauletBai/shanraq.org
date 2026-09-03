-- +goose Up
-- A cover per language, not per article.
--
-- An article is one thing in three languages, but its cover was a single
-- column, so every reader saw the same picture. That is harmless for a
-- photograph and wrong for a diagram: the Go course draws its lesson maps in
-- all three languages, and a Kazakh reader was opening a lesson under an
-- English illustration -- the first thing on the page, before a word of the
-- lesson.
--
-- Empty means "use the article's own cover", so nothing changes for the
-- articles that have one picture and want one picture.
ALTER TABLE article_translations
    ADD COLUMN IF NOT EXISTS cover_url text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE article_translations DROP COLUMN IF EXISTS cover_url;
