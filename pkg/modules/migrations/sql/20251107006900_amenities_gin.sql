-- +goose Up
-- The advanced listing search filters on amenities with `l.amenities @> $n`,
-- which cannot use a b-tree index. Give the array-containment operator the GIN
-- index it needs, so the filter does not degrade to a full scan as listings
-- grow.
CREATE INDEX IF NOT EXISTS idx_listings_amenities ON listings USING GIN (amenities);

-- +goose Down
DROP INDEX IF EXISTS idx_listings_amenities;
