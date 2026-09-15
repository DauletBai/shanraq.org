-- +goose Up
ALTER TABLE article_series
    DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series
    ADD CONSTRAINT article_series_code_lang_chk CHECK (code_lang IN ('go', 'python', 'sql'));

-- +goose Down
UPDATE article_series SET code_lang = 'go' WHERE code_lang = 'sql';
ALTER TABLE article_series
    DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series
    ADD CONSTRAINT article_series_code_lang_chk CHECK (code_lang IN ('go', 'python'));
