-- +goose Up
-- Start the view counters from zero, keeping the old numbers.
--
-- Until 2026-08-07 every crawler hit incremented views_count: of 2 284 counted
-- article views, roughly 1 474 were Googlebot, the Facebook link scraper and AI
-- crawlers (the bot-filtered analytics panel counted 810 real ones over the same
-- period). Those totals cannot be cleaned after the fact — the database holds
-- one number per article, with no record of who caused it.
--
-- Leaving them means the owner divides every figure by three in his head for
-- months, and the studio's views column stays permanently incomparable with its
-- reading-depth column, which was always bot-free. The site is two weeks old, so
-- the history being set aside is worth very little; a trustworthy number from
-- today is worth a lot more.
--
-- Nothing is destroyed. The old per-article totals move to views_legacy and the
-- old daily rows to article_views_daily_legacy, so the decision is reversible
-- with one UPDATE and one INSERT.
ALTER TABLE articles ADD COLUMN IF NOT EXISTS views_legacy BIGINT NOT NULL DEFAULT 0;

UPDATE articles SET views_legacy = views_count WHERE views_legacy = 0;
UPDATE articles SET views_count = 0;

CREATE TABLE IF NOT EXISTS article_views_daily_legacy (LIKE article_views_daily INCLUDING ALL);
INSERT INTO article_views_daily_legacy SELECT * FROM article_views_daily
    ON CONFLICT DO NOTHING;
DELETE FROM article_views_daily;

-- +goose Down
INSERT INTO article_views_daily SELECT * FROM article_views_daily_legacy
    ON CONFLICT DO NOTHING;
UPDATE articles SET views_count = views_legacy WHERE views_legacy > 0;
DROP TABLE IF EXISTS article_views_daily_legacy;
ALTER TABLE articles DROP COLUMN IF EXISTS views_legacy;
