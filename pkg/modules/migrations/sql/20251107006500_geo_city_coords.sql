-- +goose Up
-- Astana, Almaty and Shymkent sit at the top level next to the regions rather
-- than inside one, so the settlement import skipped them and the map had no
-- point for the country's three largest markets. Coordinates only — their
-- population is left NULL rather than guessed, and it would not affect
-- ordering at this level anyway.
UPDATE geo_nodes SET lat = 51.1605, lng = 71.4704 WHERE country = 'KZ' AND level = 1 AND name_ru = 'Астана'  AND lat IS NULL;
UPDATE geo_nodes SET lat = 43.2380, lng = 76.8829 WHERE country = 'KZ' AND level = 1 AND name_ru = 'Алматы'  AND lat IS NULL;
UPDATE geo_nodes SET lat = 42.3417, lng = 69.5901 WHERE country = 'KZ' AND level = 1 AND name_ru = 'Шымкент' AND lat IS NULL;

-- +goose Down
SELECT 1;
