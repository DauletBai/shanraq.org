-- +goose Up
-- Richer listing details: land-plot area (for houses/plots/dachas) and a set of
-- amenity flags (air conditioner, pool, parking, …).
ALTER TABLE listings ADD COLUMN IF NOT EXISTS land_area numeric NOT NULL DEFAULT 0;
ALTER TABLE listings ADD COLUMN IF NOT EXISTS amenities text[] NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS land_area;
ALTER TABLE listings DROP COLUMN IF EXISTS amenities;
