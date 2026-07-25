-- +goose Up
-- Aggregate, privacy-respecting traffic counters. One row per (day, kind, label,
-- audience) — NO per-visitor identifiers are ever stored, so this stays within
-- the Privacy Policy's "minimal analytics, no behavioral profiling" promise.
--   kind='page'  → label is a coarse page kind (home, article, listing, ...)
--   kind='click' → label is a named UI event (show_contact, register_cta, ...)
--   is_guest     → TRUE for anonymous visitors, FALSE for signed-in users
CREATE TABLE IF NOT EXISTS analytics_daily (
    day      DATE    NOT NULL,
    kind     TEXT    NOT NULL,
    label    TEXT    NOT NULL,
    is_guest BOOLEAN NOT NULL,
    n        BIGINT  NOT NULL DEFAULT 0,
    PRIMARY KEY (day, kind, label, is_guest)
);

-- Fast range scans for the day/month/year roll-ups shown in the admin panel.
CREATE INDEX IF NOT EXISTS analytics_daily_kind_day_idx ON analytics_daily (kind, day);

-- +goose Down
DROP TABLE IF EXISTS analytics_daily;
