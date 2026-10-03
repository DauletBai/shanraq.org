-- +goose Up
-- Keep the pre-verification aggregates visible without putting them back into
-- the tables that receive verified browser events.  These views are the read
-- model for dashboards and counters; the underlying archive remains immutable.
CREATE OR REPLACE VIEW analytics_daily_display AS
SELECT day, kind, label, is_guest, n, FALSE AS verified
  FROM analytics_daily_unverified
UNION ALL
SELECT day, kind, label, is_guest, n, TRUE AS verified
  FROM analytics_daily;

CREATE OR REPLACE VIEW article_views_daily_display AS
SELECT article_id, lang, day, views, FALSE AS verified
  FROM article_views_daily_unverified
UNION ALL
SELECT article_id, lang, day, views, TRUE AS verified
  FROM article_views_daily;

-- Historical rows predate the rolling-session id.  One old row represented
-- one visitor in one fixed half-hour visit; current rows carry sid and use the
-- rolling 30-minute rule.  Keeping the origin flag lets reports apply the
-- correct definition to each half of the series.
CREATE OR REPLACE VIEW analytics_slots_display AS
SELECT slot, vid, host, is_kz, is_mobile, views,
       NULL::BYTEA AS sid, FALSE AS verified
  FROM analytics_slots_unverified
UNION ALL
SELECT slot, vid, host, is_kz, is_mobile, views,
       sid, TRUE AS verified
  FROM analytics_slots;

-- +goose Down
DROP VIEW IF EXISTS analytics_slots_display;
DROP VIEW IF EXISTS article_views_daily_display;
DROP VIEW IF EXISTS analytics_daily_display;
