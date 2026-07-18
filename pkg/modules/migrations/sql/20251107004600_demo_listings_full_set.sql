-- +goose Up
-- Rebuild the six demo listings as a full showcase covering every property
-- type (apartment, house, commercial, land, dacha) with rooms, amenities and
-- land area populated, so the listing cards and detail pages demonstrate the
-- whole feature set with icons.

-- 01 — apartment, Almaty (sale)
UPDATE listings SET
  deal_type='sale', property_type='apartment', country='Казахстан', region='Алматы', city='Алматы', village='',
  price=42000000, area=68, rooms=2, land_area=0,
  title='2-комнатная квартира в Бостандыкском районе',
  description='Светлая квартира с ремонтом, рядом парк и школа. Тихий двор, развитая инфраструктура, вся мебель и техника остаются.',
  amenities=ARRAY['air_conditioner','furniture','internet','elevator','heating','hot_water','plastic_windows'],
  room_specs='[{"type":"living","area":20},{"type":"bedroom","area":14},{"type":"kitchen","area":12},{"type":"bathroom","area":4},{"type":"hallway","area":6},{"type":"balcony","area":4}]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000001';

-- 02 — house, Almaty region (sale)
UPDATE listings SET
  deal_type='sale', property_type='house', country='Казахстан', region='Алматинская область', city='Алматы', village='Каменское плато',
  price=120000000, area=220, rooms=5, land_area=8,
  title='Дом у гор на Каменском плато',
  description='Просторный двухэтажный дом с участком 8 соток и панорамой Заилийского Алатау. Гараж на две машины, автономное отопление на газе, охрана посёлка.',
  amenities=ARRAY['garage','parking','heating','hot_water','gas','security','internet'],
  room_specs='[{"type":"living","area":40},{"type":"bedroom","area":20,"note":"master"},{"type":"bedroom","area":18},{"type":"bedroom","area":16},{"type":"kitchen","area":18},{"type":"bathroom","area":8},{"type":"wc","area":3},{"type":"hallway","area":12}]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000002';

-- 03 — commercial, Astana (rent)
UPDATE listings SET
  deal_type='rent', property_type='commercial', country='Казахстан', region='Астана', city='Астана', village='',
  price=700000, area=150, rooms=0, land_area=0,
  title='Коммерческое помещение на первой линии, Есильский район',
  description='Помещение свободного назначения: офис, магазин или салон. Первая линия, отдельный вход, витрины, парковка для клиентов. Готово к заезду.',
  amenities=ARRAY['parking','internet','security','air_conditioner','elevator'],
  room_specs='[]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000003';

-- 04 — land, Almaty region (sale)
UPDATE listings SET
  deal_type='sale', property_type='land', country='Казахстан', region='Алматинская область', city='Талгар', village='',
  price=18000000, area=0, rooms=0, land_area=12,
  title='Земельный участок 12 соток под ИЖС в Талгаре',
  description='Ровный участок правильной формы под индивидуальное жилищное строительство. Подведён газ и электричество, асфальтированный подъезд, все документы готовы.',
  amenities=ARRAY['gas'],
  room_specs='[]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000004';

-- 05 — apartment, Astana (rent, furnished)
UPDATE listings SET
  deal_type='rent', property_type='apartment', country='Казахстан', region='Астана', city='Астана', village='',
  price=250000, area=55, rooms=2, land_area=0,
  title='Аренда 2-комнатной квартиры в Есильском районе',
  description='Меблированная квартира со всей техникой, рядом ТРЦ и станции LRT. Свежий ремонт, только на длительный срок, без животных.',
  amenities=ARRAY['furniture','fridge','washer','tv','internet','air_conditioner','elevator'],
  room_specs='[{"type":"living","area":18},{"type":"bedroom","area":14},{"type":"kitchen","area":10},{"type":"bathroom","area":5},{"type":"hallway","area":5},{"type":"balcony","area":3}]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000005';

-- 06 — dacha, Almaty region (sale)
UPDATE listings SET
  deal_type='sale', property_type='dacha', country='Казахстан', region='Алматинская область', city='Талгар', village='',
  price=25000000, area=90, rooms=3, land_area=10,
  title='Дача в предгорье Талгара, 10 соток',
  description='Сад с плодовыми деревьями, скважина, баня. Круглогодичный подъезд, электричество 380В, газ по границе участка. Готова к проживанию.',
  amenities=ARRAY['garage','parking','gas','hot_water','internet'],
  room_specs='[{"type":"living","area":24},{"type":"bedroom","area":16},{"type":"kitchen","area":14},{"type":"bathroom","area":6},{"type":"hallway","area":8},{"type":"other","area":12,"note":"веранда"}]'::jsonb
 WHERE id='a0000000-0000-0000-0000-000000000006';

-- +goose Down
-- Content-only refresh; nothing to revert.
SELECT 1;
