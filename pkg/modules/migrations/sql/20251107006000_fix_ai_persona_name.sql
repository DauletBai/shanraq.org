-- +goose Up
-- The AI columnist was renamed to "AI Dake" (same Latin form in all three
-- languages), but a handful of first-person passages inside article bodies
-- still introduced it by the old name. The same passages transliterated the
-- platform four different ways — «Шанырак», «Шанрак», «Шаңырақ», «Шанрақ» —
-- where the brand is written Shanraq.
--
-- Note the Kazakh words "сана" (mind) and "шаңырақ" (the yurt crown) are
-- ordinary vocabulary that appears in these articles, so these are
-- exact-phrase replacements rather than a blanket search for the tokens.

-- Persona name in Russian bodies (self-introductions and quote attributions).
UPDATE article_translations
   SET body_md = replace(body_md, 'Сана Кыран', 'AI Dake')
 WHERE body_md LIKE '%Сана Кыран%';

-- Platform name: normalise every declension to the Latin brand.
UPDATE article_translations
   SET body_md = replace(
                   replace(
                     replace(
                       replace(body_md, 'платформы «Шанырак»', 'платформы Shanraq'),
                       'ИИ-колумнист Шанрак', 'ИИ-колумнист Shanraq'),
                     'Шаңырақ платформасының', 'Shanraq платформасының'),
                   'Шанрақтың', 'Shanraq-тың')
 WHERE body_md LIKE '%«Шанырак»%'
    OR body_md LIKE '%ИИ-колумнист Шанрак%'
    OR body_md LIKE '%Шаңырақ платформасының%'
    OR body_md LIKE '%Шанрақтың%';

UPDATE article_translations
   SET body_md = replace(body_md, 'the Shanyraq platform', 'the Shanraq platform')
 WHERE body_md LIKE '%the Shanyraq platform%';

-- Two leftovers of the old feminine persona name: a Kazakh suffix glued to the
-- Latin brand without a hyphen, and a Russian sentence that referred to the
-- columnist in the feminine ("соврать самой себе"). The replacement avoids
-- grammatical gender altogether, which is what an AI byline should do.
UPDATE article_translations
   SET body_md = replace(body_md, 'Shanraqтың', 'Shanraq-тың')
 WHERE body_md LIKE '%Shanraqтың%';

UPDATE article_translations
   SET body_md = replace(body_md, 'не соврать самой себе', 'не врать себе')
 WHERE body_md LIKE '%не соврать самой себе%';

-- +goose Down
SELECT 1;
