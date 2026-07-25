-- +goose Up
-- Tidy the seeded content so a fresh database looks like a real publication.
--
-- 1) Drop the Redaksiya demo articles. They shipped without covers or
--    translations and pushed the AI Dake columns down the feed. The redaksiya
--    account and its showcase real-estate listings are kept.
-- 2) Spread AI Dake's publish dates. The seed stamps every column with NOW(),
--    so all 90 land on one day; restage them descending (newest first) so the
--    feed reads as an ongoing publication rather than a single dump.

DELETE FROM articles
 WHERE author_id = '11111111-1111-1111-1111-111111111111';

WITH ordered AS (
    SELECT id, (row_number() OVER (ORDER BY published_at DESC, id) - 1) AS rn
    FROM articles
    WHERE author_id = '5a2a0000-0000-0000-0000-000000000001'
      AND status = 'published'
)
UPDATE articles a
   SET published_at = NOW() - (o.rn * INTERVAL '11 hours'),
       created_at   = NOW() - (o.rn * INTERVAL '11 hours') - INTERVAL '30 minutes'
  FROM ordered o
 WHERE a.id = o.id;

-- +goose Down
-- One-way cleanup: deleted demo articles and the original single-day timestamps
-- are not restored.
SELECT 1;
