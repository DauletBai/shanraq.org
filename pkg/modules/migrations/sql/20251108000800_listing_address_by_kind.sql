-- +goose Up
-- Address fields were filled by a node's DEPTH in the location tree, not by what
-- the node actually is. That works only where every country is shaped the same
-- way, and none are: Saint Petersburg is a federal city sitting where an oblast
-- sits, so a listing there was published as "Область: Санкт-Петербург, Город:
-- Петродворцовый" — an oblast that does not exist and a district called a city.
-- Almaty, Astana and Shymkent have the same shape at home and were wrong in the
-- same way.
--
-- The fix is to fill each field from the node's kind. That also means the fourth
-- slot only ever holds an intra-city district now: settlements (city, town,
-- village) all land in `city`, so the column called `village` was misnamed.
ALTER TABLE listings RENAME COLUMN village TO district;

-- Re-derive the address of every listing that has a location node, so the ones
-- already published stop showing the wrong labels. Names come from the Russian
-- column, which is what the reference is keyed on.
WITH RECURSIVE up AS (
    SELECT l.id AS listing_id, g.id, g.parent_id, g.kind, g.name_ru
    FROM listings l
    JOIN geo_nodes g ON g.id = l.geo_node_id
    UNION ALL
    SELECT u.listing_id, p.id, p.parent_id, p.kind, p.name_ru
    FROM up u
    JOIN geo_nodes p ON p.id = u.parent_id
), agg AS (
    SELECT listing_id,
           max(name_ru) FILTER (WHERE kind = 'country') AS country,
           max(name_ru) FILTER (WHERE kind IN ('region', 'oblast', 'krai', 'republic', 'okrug')) AS region,
           max(name_ru) FILTER (WHERE kind IN ('city', 'town', 'village')) AS city,
           max(name_ru) FILTER (WHERE kind = 'district') AS district
    FROM up GROUP BY listing_id
)
UPDATE listings l SET
    country  = COALESCE(a.country, ''),
    region   = COALESCE(a.region, ''),
    city     = COALESCE(a.city, ''),
    district = COALESCE(a.district, '')
FROM agg a WHERE l.id = a.listing_id;

-- +goose Down
ALTER TABLE listings RENAME COLUMN district TO village;
