-- +goose Up
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust', 'shell', 'math', 'informatics', 'kazakh'));

-- +goose Down
DO $guard$
BEGIN
  IF EXISTS (SELECT 1 FROM article_series WHERE code_lang = 'kazakh') THEN
    RAISE EXCEPTION 'Cannot remove Kazakh-language course support while such a course exists';
  END IF;
END $guard$;
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series ADD CONSTRAINT article_series_code_lang_chk
    CHECK (code_lang IN ('go', 'python', 'sql', 'rust', 'shell', 'math', 'informatics'));
