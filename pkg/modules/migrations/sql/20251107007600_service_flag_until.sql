-- +goose Up
-- Optional planned end time for a maintenance window. When set on the global
-- 'site' switch, the maintenance page shows a live countdown to this moment and
-- the site auto-restores itself once it passes, so we never leave the site down
-- by forgetting to switch it back.
ALTER TABLE service_flags ADD COLUMN IF NOT EXISTS until TIMESTAMPTZ;

-- +goose Down
ALTER TABLE service_flags DROP COLUMN IF EXISTS until;
