-- +goose Up
-- geo_nodes held only names, so the city list could not be ordered by size and
-- the map had to keep region centroids hardcoded in Go. Give a node its own
-- population and coordinates; both stay optional because districts and
-- countries have no meaningful value here.
ALTER TABLE geo_nodes ADD COLUMN IF NOT EXISTS population      INTEGER;
ALTER TABLE geo_nodes ADD COLUMN IF NOT EXISTS population_year SMALLINT;
ALTER TABLE geo_nodes ADD COLUMN IF NOT EXISTS lat             DOUBLE PRECISION;
ALTER TABLE geo_nodes ADD COLUMN IF NOT EXISTS lng             DOUBLE PRECISION;

-- The settlement dropdown lists one parent's children ordered by size, so the
-- index carries the sort key with the parent.
CREATE INDEX IF NOT EXISTS idx_geo_children_by_size
    ON geo_nodes (parent_id, population DESC NULLS LAST);

-- +goose Down
DROP INDEX IF EXISTS idx_geo_children_by_size;
ALTER TABLE geo_nodes DROP COLUMN IF EXISTS lng;
ALTER TABLE geo_nodes DROP COLUMN IF EXISTS lat;
ALTER TABLE geo_nodes DROP COLUMN IF EXISTS population_year;
ALTER TABLE geo_nodes DROP COLUMN IF EXISTS population;
