-- +goose Up
ALTER TABLE listings ADD COLUMN IF NOT EXISTS images TEXT[] NOT NULL DEFAULT '{}';

-- Seed demo galleries (room illustrations) for the 12 demo listings, and set
-- each cover to the first image so cards show it too.
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000001';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bathroom.svg','/static/demo/rooms/exterior.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000002';
UPDATE listings SET images = ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/exterior.svg' WHERE id = 'a0000000-0000-0000-0000-000000000003';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000004';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000005';
UPDATE listings SET images = ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg'], cover_url = '/static/demo/rooms/exterior.svg' WHERE id = 'a0000000-0000-0000-0000-000000000006';
UPDATE listings SET images = ARRAY['/static/demo/rooms/plot.svg','/static/demo/rooms/exterior.svg'], cover_url = '/static/demo/rooms/plot.svg' WHERE id = 'a0000000-0000-0000-0000-000000000007';
UPDATE listings SET images = ARRAY['/static/demo/rooms/office.svg','/static/demo/rooms/exterior.svg'], cover_url = '/static/demo/rooms/office.svg' WHERE id = 'a0000000-0000-0000-0000-000000000008';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000009';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000010';
UPDATE listings SET images = ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/exterior.svg' WHERE id = 'a0000000-0000-0000-0000-000000000011';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bathroom.svg'], cover_url = '/static/demo/rooms/living.svg' WHERE id = 'a0000000-0000-0000-0000-000000000012';

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS images;
