-- +goose Up
-- Editable tariffs (ad rate card, listing service prices, durations) as a
-- key→value table, seeded from the built-in defaults on boot. Prices and
-- periods move out of Go so an operator can retune them from the admin panel
-- without a code change; the reader falls back to the built-in default if a row
-- is ever missing, so pricing never breaks.
CREATE TABLE IF NOT EXISTS tariffs (
    key        TEXT PRIMARY KEY,
    value      BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth_users(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS tariffs;
