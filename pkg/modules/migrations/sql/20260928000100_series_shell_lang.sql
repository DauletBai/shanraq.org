-- +goose Up
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust', 'shell'));

-- +goose Down
-- Refuse rollback while shell courses exist instead of silently relabelling them.
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust'));
