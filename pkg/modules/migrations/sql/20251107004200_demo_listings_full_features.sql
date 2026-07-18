-- +goose Up
-- Keep six demo listings as fully-featured showcase samples (rooms with specs,
-- amenities, land area) and drop the rest. These are intentionally retained in
-- every environment as starter content.

-- Drop demo listings 07..12 (keep 01..06).
DELETE FROM listings
 WHERE id::text LIKE 'a0000000-0000-0000-0000-0000000000%'
   AND id NOT IN (
     'a0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000002',
     'a0000000-0000-0000-0000-000000000003',
     'a0000000-0000-0000-0000-000000000004',
     'a0000000-0000-0000-0000-000000000005',
     'a0000000-0000-0000-0000-000000000006'
   );

-- 01 — 2-room apartment, Almaty (68 m2)
UPDATE listings SET
  amenities = ARRAY['air_conditioner','furniture','internet','elevator','heating','hot_water','plastic_windows'],
  room_specs = '[{"type":"living","area":20},{"type":"bedroom","area":14},{"type":"kitchen","area":12},{"type":"bathroom","area":4},{"type":"hallway","area":6},{"type":"balcony","area":4}]'::jsonb
 WHERE id = 'a0000000-0000-0000-0000-000000000001';

-- 02 — 1-room apartment, Astana new build (42 m2)
UPDATE listings SET
  amenities = ARRAY['parking','internet','elevator','heating','hot_water','plastic_windows','security'],
  room_specs = '[{"type":"living","area":20},{"type":"kitchen","area":10},{"type":"bathroom","area":5},{"type":"hallway","area":4},{"type":"balcony","area":3}]'::jsonb
 WHERE id = 'a0000000-0000-0000-0000-000000000002';

-- 03 — house on Kamenskoe plato (220 m2, 8 ares)
UPDATE listings SET
  land_area = 8,
  amenities = ARRAY['garage','parking','heating','hot_water','gas','security','internet'],
  room_specs = '[{"type":"living","area":40},{"type":"bedroom","area":20,"note":"master"},{"type":"bedroom","area":18},{"type":"bedroom","area":16},{"type":"kitchen","area":18},{"type":"bathroom","area":8},{"type":"wc","area":3},{"type":"hallway","area":12}]'::jsonb
 WHERE id = 'a0000000-0000-0000-0000-000000000003';

-- 04 — 3-room apartment, Shymkent (90 m2)
UPDATE listings SET
  amenities = ARRAY['heating','playground','security','parking','plastic_windows','internet'],
  room_specs = '[{"type":"living","area":24},{"type":"bedroom","area":16},{"type":"bedroom","area":14},{"type":"kitchen","area":14},{"type":"bathroom","area":6},{"type":"wc","area":3},{"type":"hallway","area":8},{"type":"balcony","area":5}]'::jsonb
 WHERE id = 'a0000000-0000-0000-0000-000000000004';

-- 05 — 2-room apartment for rent, Astana (55 m2, furnished)
UPDATE listings SET
  amenities = ARRAY['furniture','fridge','washer','tv','internet','air_conditioner','elevator'],
  room_specs = '[{"type":"living","area":18},{"type":"bedroom","area":14},{"type":"kitchen","area":10},{"type":"bathroom","area":5},{"type":"hallway","area":5},{"type":"balcony","area":3}]'::jsonb
 WHERE id = 'a0000000-0000-0000-0000-000000000005';

-- 06 — dacha near Talgar (80 m2, 10 ares)
UPDATE listings SET
  land_area = 10,
  amenities = ARRAY['garage','parking','gas','hot_water','internet'],
  room_specs = '[{"type":"living","area":24},{"type":"bedroom","area":16},{"type":"kitchen","area":14},{"type":"bathroom","area":6},{"type":"hallway","area":8},{"type":"other","area":12,"note":"veranda"}]'::jsonb
 WHERE id = 'a0000000-0000-0000-0000-000000000006';

-- +goose Down
-- Best-effort: clear the enriched fields (the deleted 07..12 are not restored).
UPDATE listings SET amenities = '{}', room_specs = '[]'::jsonb, land_area = 0
 WHERE id::text LIKE 'a0000000-0000-0000-0000-0000000000%';
