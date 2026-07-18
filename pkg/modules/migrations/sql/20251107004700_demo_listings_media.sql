-- +goose Up
-- Give each rebuilt demo listing a coherent photo gallery matching its type,
-- so the card carousels and detail galleries look full and demonstrate the
-- feature. Uses the bundled demo room illustrations.
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url='/static/demo/rooms/living.svg'
 WHERE id='a0000000-0000-0000-0000-000000000001';
UPDATE listings SET images = ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url='/static/demo/rooms/exterior.svg'
 WHERE id='a0000000-0000-0000-0000-000000000002';
UPDATE listings SET images = ARRAY['/static/demo/rooms/office.svg','/static/demo/rooms/exterior.svg'], cover_url='/static/demo/rooms/office.svg'
 WHERE id='a0000000-0000-0000-0000-000000000003';
UPDATE listings SET images = ARRAY['/static/demo/rooms/plot.svg','/static/demo/rooms/exterior.svg'], cover_url='/static/demo/rooms/plot.svg'
 WHERE id='a0000000-0000-0000-0000-000000000004';
UPDATE listings SET images = ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'], cover_url='/static/demo/rooms/living.svg'
 WHERE id='a0000000-0000-0000-0000-000000000005';
UPDATE listings SET images = ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg'], cover_url='/static/demo/rooms/exterior.svg'
 WHERE id='a0000000-0000-0000-0000-000000000006';

-- +goose Down
SELECT 1;
