-- +goose Up
-- A listing could only say which city it was in. Buyers choose by district and
-- street long before they call, so store the street address alongside the geo
-- node, plus optional coordinates for a per-listing pin on the map (the map
-- currently only bubbles counts per region).
ALTER TABLE listings ADD COLUMN IF NOT EXISTS microdistrict TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS street        TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS house         TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS lat           DOUBLE PRECISION;
ALTER TABLE listings ADD COLUMN IF NOT EXISTS lng           DOUBLE PRECISION;

-- Coordinates are optional, so only index the rows that actually carry a pin.
CREATE INDEX IF NOT EXISTS idx_listings_coords ON listings (lat, lng)
    WHERE lat IS NOT NULL AND lng IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_listings_coords;
ALTER TABLE listings DROP COLUMN IF EXISTS lng;
ALTER TABLE listings DROP COLUMN IF EXISTS lat;
ALTER TABLE listings DROP COLUMN IF EXISTS house;
ALTER TABLE listings DROP COLUMN IF EXISTS street;
ALTER TABLE listings DROP COLUMN IF EXISTS microdistrict;
