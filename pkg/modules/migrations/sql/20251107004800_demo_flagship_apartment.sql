-- +goose Up
-- Promote demo listing 01 to the flagship showcase card: a fully-detailed
-- three-room apartment ("Демо объявление трехкомнатной квартиры") that fills
-- the facts row (flag · area · bedrooms · bathrooms, then living / kitchen
-- areas), carries the "Топ" corner ribbon, and has a three-photo gallery.
UPDATE listings SET
  deal_type='sale', property_type='apartment', country='Казахстан', region='Алматы', city='Алматы', village='',
  price=18000000, area=72, rooms=3, land_area=0,
  title='Демо объявление трехкомнатной квартиры',
  description='Демонстрационное объявление: показывает все возможности карточки — флаг страны, общую площадь, число спален и санузлов, площади зала и кухни, галерею фотографий и удобства. Светлая трёхкомнатная квартира с ремонтом в тихом дворе.',
  amenities=ARRAY['air_conditioner','furniture','fridge','washer','internet','elevator','heating','hot_water','plastic_windows'],
  room_specs='[{"type":"living","area":20},{"type":"bedroom","area":15},{"type":"bedroom","area":12},{"type":"kitchen","area":7},{"type":"bathroom","area":4},{"type":"hallway","area":8},{"type":"balcony","area":6}]'::jsonb,
  images=ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg'],
  cover_url='/static/demo/rooms/living.svg'
 WHERE id='a0000000-0000-0000-0000-000000000001';

-- Give it the "Топ" ribbon (promoted) alongside its existing highlight.
UPDATE listings SET promoted_until = NOW() + INTERVAL '30 days'
 WHERE id='a0000000-0000-0000-0000-000000000001';

-- +goose Down
SELECT 1;
