-- +goose Up
-- Agents want to attach plans, schemes and documents (floor plans, technical
-- passports, contracts) to illustrate an object better than the photo feed.
-- Stored like images: an array of URLs (PDF stored as-is, image plans via the
-- normal pipeline).
ALTER TABLE listings ADD COLUMN IF NOT EXISTS documents TEXT[] NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS documents;
