-- +goose Up
-- Per-listing engagement counters: how many times the listing was viewed and
-- how many times a viewer revealed the seller's contact.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS views_count    integer NOT NULL DEFAULT 0;
ALTER TABLE listings ADD COLUMN IF NOT EXISTS contacts_count integer NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS views_count;
ALTER TABLE listings DROP COLUMN IF EXISTS contacts_count;
