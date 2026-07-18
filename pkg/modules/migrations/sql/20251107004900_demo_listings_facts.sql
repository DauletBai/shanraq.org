-- +goose Up
-- Fill the facts row (area · rooms · areas) on the two demo cards that were
-- still bare: the commercial space and the land plot. After this every one of
-- the six demo listings shows a populated facts summary under its price.

-- 03 — commercial: give it an open trading hall, a utility kitchenette and two
-- restrooms so the card shows "150 м² · санузел 2 · зал 110 · кухня 8".
UPDATE listings SET
  area=150,
  room_specs='[{"type":"living","area":110,"note":"торговый зал"},{"type":"kitchen","area":8,"note":"подсобное"},{"type":"wc","area":4},{"type":"wc","area":3}]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000003';

-- 04 — land: a plot has no rooms, so its defining fact is the plot size. Keep
-- the land area (12 соток) and make sure the built-up area stays 0 so the card
-- shows the plot size rather than an empty floor area.
UPDATE listings SET
  area=0, land_area=12, rooms=0, room_specs='[]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000004';

-- +goose Down
SELECT 1;
