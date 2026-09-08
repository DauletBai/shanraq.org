-- +goose Up
-- Which language a course's exercises are written in.
--
-- The exercise checker was built for the Go course and assumed it everywhere:
-- a submission went through gofmt before anything else, so a correct Python
-- solution was refused by a Go parser before the reviewer ever saw it. The
-- second course made that assumption wrong, and there was nowhere to say so.
--
-- The course is the right place to keep it: a lesson belongs to one, and every
-- part of the check -- the tidier, the syntax gate, the colours, the brief the
-- reviewer is given -- needs the same answer.
ALTER TABLE article_series
    ADD COLUMN IF NOT EXISTS code_lang text NOT NULL DEFAULT 'go';

ALTER TABLE article_series
    DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series
    ADD CONSTRAINT article_series_code_lang_chk CHECK (code_lang IN ('go', 'python'));

UPDATE article_series SET code_lang = 'python' WHERE slug = 'python';

-- +goose Down
ALTER TABLE article_series DROP CONSTRAINT IF EXISTS article_series_code_lang_chk;
ALTER TABLE article_series DROP COLUMN IF EXISTS code_lang;
