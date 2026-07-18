-- +goose Up
-- Rename the AI columnist from "Сана Қыран" / "Sana Qyran" to "AI Dake" inside
-- the persona's already-published columns, so the article prose matches the new
-- byline. Replace the Kazakh predicate-suffixed form ("Сана Қыранмын" = "I am
-- Sana Qyran") first, then the standalone Cyrillic, then the Latin spelling, so
-- no dangling suffix is left behind.
UPDATE article_translations t SET
  title = replace(replace(replace(title,
    'Сана Қыранмын', 'AI Dake'), 'Сана Қыран', 'AI Dake'), 'Sana Qyran', 'AI Dake'),
  summary = replace(replace(replace(summary,
    'Сана Қыранмын', 'AI Dake'), 'Сана Қыран', 'AI Dake'), 'Sana Qyran', 'AI Dake'),
  body_md = replace(replace(replace(body_md,
    'Сана Қыранмын', 'AI Dake'), 'Сана Қыран', 'AI Dake'), 'Sana Qyran', 'AI Dake')
FROM articles a
WHERE t.article_id = a.id
  AND a.author_id = '5a2a0000-0000-0000-0000-000000000001';

-- +goose Down
-- Content rename; not reverted.
SELECT 1;
