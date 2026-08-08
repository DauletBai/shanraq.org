-- +goose Up
-- Reset the traffic-source counters only, and keep everything else.
--
-- The guest panel records ten metrics. Exactly one of them was ever wrong.
--
-- Sources were hit twice. Until 2026-08-07 the source was recorded on every
-- page view, so a reader who arrived from Facebook and opened four more pages
-- scored Facebook 1 and Direct 4 — internal navigation carries a same-host
-- referrer, which classified as "direct". And until 2026-08-08 every host
-- containing "google" counted as the search engine, so links opened from Gmail
-- and Google Translate were filed as organic search. Both are fixed; the
-- accumulated rows are not repairable, because the damage is inside the totals.
--
-- The other nine metrics — pages, clicks, bots, devices, OS, browsers,
-- countries, reading language, country×language — were correct throughout.
-- Crawlers were never counted as guests: that separation is what produced the
-- honest 810 real article views against the counter's 2 284, and it is the
-- yardstick the whole clean-up was measured with. Wiping that data would
-- destroy the audience breakdown, the bot statistics and the only record of how
-- the site launched, in exchange for nothing but a matching date.
--
-- Preserved, not deleted, like the view and depth counters before it.
CREATE TABLE IF NOT EXISTS analytics_daily_legacy (LIKE analytics_daily INCLUDING ALL);

INSERT INTO analytics_daily_legacy
SELECT * FROM analytics_daily WHERE kind = 'source'
    ON CONFLICT DO NOTHING;

DELETE FROM analytics_daily WHERE kind = 'source';

-- +goose Down
INSERT INTO analytics_daily
SELECT * FROM analytics_daily_legacy WHERE kind = 'source'
    ON CONFLICT (day, kind, label, is_guest) DO UPDATE SET n = EXCLUDED.n;
DELETE FROM analytics_daily_legacy WHERE kind = 'source';
