-- +goose Up
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust', 'shell', 'math', 'informatics'));

-- +goose Down
-- Refuse rollback while Informatics courses exist instead of silently
-- treating written observations and system explanations as source code.
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust', 'shell', 'math'));
