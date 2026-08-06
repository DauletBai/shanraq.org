-- +goose Up
-- The handler now rejects NaN, ±Inf and negative areas, but the handler is one
-- way in and not the last one: a future import, an admin tool or a second
-- endpoint would find the column just as permissive as before. These are the
-- same rules stated where they cannot be bypassed.
--
-- NaN needs saying explicitly: in SQL it is not caught by a range test, because
-- every comparison against it is false, so `area >= 0` passes it straight
-- through. That is exactly how it got into the filters and quietly removed a
-- listing from every search.
UPDATE listings SET area = 0 WHERE area IS NOT NULL AND (area <> area OR area < 0 OR area > 100000);
UPDATE listings SET land_area = 0 WHERE land_area IS NOT NULL AND (land_area <> land_area OR land_area < 0 OR land_area > 100000);
UPDATE listings SET rooms = 0 WHERE rooms IS NOT NULL AND (rooms < 0 OR rooms > 100);
UPDATE listings SET price = 0 WHERE price IS NOT NULL AND price < 0;

ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_area_check;
ALTER TABLE listings ADD CONSTRAINT listings_area_check
	CHECK (area IS NULL OR (area = area AND area >= 0 AND area <= 100000));

ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_land_area_check;
ALTER TABLE listings ADD CONSTRAINT listings_land_area_check
	CHECK (land_area IS NULL OR (land_area = land_area AND land_area >= 0 AND land_area <= 100000));

ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_rooms_check;
ALTER TABLE listings ADD CONSTRAINT listings_rooms_check
	CHECK (rooms IS NULL OR (rooms >= 0 AND rooms <= 100));

ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_price_check;
ALTER TABLE listings ADD CONSTRAINT listings_price_check
	CHECK (price IS NULL OR price >= 0);

-- +goose Down
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_price_check;
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_rooms_check;
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_land_area_check;
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_area_check;
