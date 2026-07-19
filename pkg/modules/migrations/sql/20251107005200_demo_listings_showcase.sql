-- +goose Up
-- Replace all existing/pale demo & test listings with 8 curated, fully-populated
-- showcase listings that exercise every property type (apartment, house, land,
-- commercial, dacha), both deal types (sale, rent), and every card feature:
-- country flag, facts row (area / rooms / land / zal / kitchen), amenity icons,
-- and a photo carousel (>=3 images each, auto-rotating). Placement mix:
--   01, 02 -> Top + highlighted; 03, 04 -> Top only; 05..08 -> plain.
-- Owned by the kept demo author 11111111-...; far-future expiry keeps them out
-- of the purge sweep. Pre-launch reset: there are no real user listings yet.
DELETE FROM listing_reports;
DELETE FROM favorites WHERE item_type = 'listing';
DELETE FROM listings;

-- Ensure the kept demo author exists (it may have been removed in some DBs).
INSERT INTO auth_users (id, email, password_hash, role)
VALUES ('11111111-1111-1111-1111-111111111111','redaksiya@shanraq.org','seed-no-login','user')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO listings (id, author_id, deal_type, property_type, country, region, city, village,
                      price, area, rooms, land_area, title, description, contact, cover_url, images,
                      amenities, room_specs, status, expires_at, promoted_until, featured_until)
VALUES
  ('d0000000-0000-0000-0000-000000000001','11111111-1111-1111-1111-111111111111','sale','apartment','Казахстан','Алматы','Алматы','',
    48500000,82,3,0,'3-комнатная квартира с ремонтом, Бостандыкский район',
    'Светлая трёхкомнатная квартира с качественным ремонтом в Бостандыкском районе. Тихий двор, рядом парк, школа и станции метро. Вся мебель и техника остаются — заезжай и живи.',
    '+7 700 111 22 33','/static/demo/rooms/living.svg',
    ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'],
    ARRAY['air_conditioner','furniture','fridge','washer','internet','tv','elevator','heating','hot_water','plastic_windows','security','parking'],
    '[{"type":"living","area":22},{"type":"bedroom","area":15},{"type":"bedroom","area":12},{"type":"kitchen","area":11},{"type":"bathroom","area":5},{"type":"wc","area":2},{"type":"hallway","area":8},{"type":"balcony","area":5}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NOW()+INTERVAL '100 years', NOW()+INTERVAL '100 years'),

  ('d0000000-0000-0000-0000-000000000002','11111111-1111-1111-1111-111111111111','sale','house','Казахстан','Алматинская область','Алматы','Каменское плато',
    165000000,260,6,10,'Дом премиум-класса у гор, 10 соток',
    'Просторный двухэтажный дом премиум-класса с участком 10 соток и панорамой Заилийского Алатау. Гараж на две машины, бассейн, автономное отопление на газе, охраняемый посёлок.',
    '+7 701 222 33 44','/static/demo/rooms/exterior.svg',
    ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'],
    ARRAY['garage','parking','pool','heating','hot_water','gas','security','internet','furniture'],
    '[{"type":"living","area":45},{"type":"bedroom","area":22,"note":"master"},{"type":"bedroom","area":18},{"type":"bedroom","area":16},{"type":"bedroom","area":14},{"type":"kitchen","area":20},{"type":"bathroom","area":9},{"type":"bathroom","area":6},{"type":"wc","area":3},{"type":"hallway","area":14}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NOW()+INTERVAL '100 years', NOW()+INTERVAL '100 years'),

  ('d0000000-0000-0000-0000-000000000003','11111111-1111-1111-1111-111111111111','rent','apartment','Казахстан','Астана','Астана','',
    320000,58,2,0,'Аренда 2-комнатной квартиры у LRT, Есиль',
    'Меблированная двухкомнатная квартира в Есильском районе: вся техника, свежий ремонт. Рядом ТРЦ, парк и станции LRT. Только на длительный срок, депозит обязателен.',
    '+7 702 333 44 55','/static/demo/rooms/living.svg',
    ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg','/static/demo/rooms/bathroom.svg'],
    ARRAY['furniture','fridge','washer','tv','internet','air_conditioner','elevator','heating','hot_water'],
    '[{"type":"living","area":20},{"type":"bedroom","area":15},{"type":"kitchen","area":11},{"type":"bathroom","area":5},{"type":"hallway","area":5},{"type":"balcony","area":2}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NOW()+INTERVAL '100 years', NULL),

  ('d0000000-0000-0000-0000-000000000004','11111111-1111-1111-1111-111111111111','rent','commercial','Казахстан','Астана','Астана','',
    850000,180,0,0,'Коммерческое помещение на 1-й линии, Есиль',
    'Помещение свободного назначения на первой линии: магазин, офис или салон. Отдельный вход, витрины, парковка для клиентов, высокий пешеходный трафик. Готово к заезду.',
    '+7 705 444 55 66','/static/demo/rooms/office.svg',
    ARRAY['/static/demo/rooms/office.svg','/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg'],
    ARRAY['parking','internet','security','air_conditioner','elevator'],
    '[{"type":"living","area":130,"note":"торговый зал"},{"type":"kitchen","area":10,"note":"подсобное"},{"type":"wc","area":5},{"type":"wc","area":4}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NOW()+INTERVAL '100 years', NULL),

  ('d0000000-0000-0000-0000-000000000005','11111111-1111-1111-1111-111111111111','sale','land','Казахстан','Алматинская область','Талгар','',
    22000000,0,0,15,'Земельный участок 15 соток под ИЖС, Талгар',
    'Ровный участок 15 соток правильной формы под индивидуальное жилищное строительство в Талгаре. Подведён газ и электричество, асфальтированный подъезд, все документы готовы.',
    '+7 707 555 66 77','/static/demo/rooms/plot.svg',
    ARRAY['/static/demo/rooms/plot.svg','/static/demo/rooms/exterior.svg','/static/demo/rooms/office.svg'],
    ARRAY['gas'],
    '[]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NULL, NULL),

  ('d0000000-0000-0000-0000-000000000006','11111111-1111-1111-1111-111111111111','sale','dacha','Казахстан','Алматинская область','Талгар','',
    28000000,95,3,8,'Дача с баней и садом, 8 соток',
    'Уютная дача с баней и плодовым садом на 8 сотках в предгорье Талгара. Скважина, электричество 380В, круглогодичный подъезд. Готова к проживанию.',
    '+7 700 666 77 88','/static/demo/rooms/exterior.svg',
    ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg'],
    ARRAY['garage','parking','gas','hot_water','internet','pool'],
    '[{"type":"living","area":26},{"type":"bedroom","area":15},{"type":"kitchen","area":14},{"type":"bathroom","area":6},{"type":"hallway","area":8},{"type":"other","area":12,"note":"веранда"}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NULL, NULL),

  ('d0000000-0000-0000-0000-000000000007','11111111-1111-1111-1111-111111111111','sale','apartment','Казахстан','Туркестанская область','Шымкент','',
    19900000,45,1,0,'1-комнатная квартира в новостройке, центр Шымкента',
    'Однокомнатная квартира-студия в новой кирпичной новостройке в центре Шымкента. Автономное отопление, закрытая территория, детская площадка. Идеальный первый дом.',
    '+7 701 777 88 99','/static/demo/rooms/living.svg',
    ARRAY['/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bathroom.svg'],
    ARRAY['furniture','internet','elevator','heating','hot_water','plastic_windows','playground'],
    '[{"type":"living","area":24},{"type":"kitchen","area":9},{"type":"bathroom","area":4},{"type":"hallway","area":5},{"type":"balcony","area":3}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NULL, NULL),

  ('d0000000-0000-0000-0000-000000000008','11111111-1111-1111-1111-111111111111','rent','house','Казахстан','Алматы','Алматы','',
    550000,180,5,6,'Аренда дома с участком, Медеуский район',
    'Аренда просторного дома с участком 6 соток в Медеуском районе. Пять комнат, гараж, вся мебель и техника, охрана. Тихая улица, рядом горы. На длительный срок.',
    '+7 702 888 99 00','/static/demo/rooms/exterior.svg',
    ARRAY['/static/demo/rooms/exterior.svg','/static/demo/rooms/living.svg','/static/demo/rooms/kitchen.svg','/static/demo/rooms/bedroom.svg'],
    ARRAY['garage','parking','furniture','heating','hot_water','gas','security','internet','air_conditioner'],
    '[{"type":"living","area":36},{"type":"bedroom","area":18},{"type":"bedroom","area":15},{"type":"bedroom","area":13},{"type":"kitchen","area":16},{"type":"bathroom","area":7},{"type":"wc","area":3},{"type":"hallway","area":10}]'::jsonb,
    'published', NOW()+INTERVAL '100 years', NULL, NULL);
-- +goose StatementEnd

-- +goose Down
SELECT 1;
