-- +goose Up
-- The AI columnist's byline already carries an "AI opinion" badge, so every
-- summary that also opened with "Мнение ИИ-колумниста:" / "ЖИ пікірі:" /
-- "An AI's opinion:" was saying the same thing a second time. Strip the label
-- and let the summary start with what it is actually about.
--
-- The pattern matches only a leading label that mentions ИИ / ЖИ / AI and ends
-- in a colon before any sentence punctuation, so ordinary summaries (and
-- unrelated openers like "Знакомство:") are left alone.
UPDATE article_translations
   SET summary = regexp_replace(summary, '^[^:.!?]{0,60}(ИИ|ЖИ|AI)[^:.!?]{0,40}:\s*', '')
 WHERE summary ~ '^[^:.!?]{0,60}(ИИ|ЖИ|AI)[^:.!?]{0,40}:\s';

-- The sentence now starts a word earlier, so re-capitalise it. upper() follows
-- the database ctype and is ASCII-only under the C locale, which would leave
-- every Russian and Kazakh summary starting lowercase — so map the Cyrillic
-- and Kazakh-specific letters explicitly instead of trusting the locale.
UPDATE article_translations
   SET summary = translate(upper(left(summary, 1)),
                           'абвгдеёжзийклмнопрстуфхцчшщъыьэюяәғқңөұүһі',
                           'АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯӘҒҚҢӨҰҮҺІ')
                 || substr(summary, 2)
 WHERE summary ~ '^[a-zабвгдеёжзийклмнопрстуфхцчшщъыьэюяәғқңөұүһі]';

-- +goose Down
-- Irreversible by design: the removed labels were boilerplate, and restoring
-- them would mean guessing which of the two dozen wordings each row had.
SELECT 1;
