-- +goose Up
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust'));

-- +goose Down
-- Refuse rollback while Rust courses exist instead of silently labelling them Go.
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql'));
