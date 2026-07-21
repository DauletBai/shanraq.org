-- +goose Up
-- Zharkent was filed under Almaty Region and stayed there after the 2022
-- reform moved it to Zhetysu, so the town existed twice. Repoint anything
-- attached to the stale row, then drop it.
--
-- The other repeated names are NOT duplicates and must survive: Karabulak is
-- three different villages (Almaty 13k, Zhetysu 19k, Turkistan 52k), Abay is
-- four places, and Karatau/Esil/Alatau each name both a town and a city
-- district.
UPDATE listings l
   SET geo_node_id = (SELECT id FROM geo_nodes WHERE code = 'kz-zhetysu-zharkent')
 WHERE l.geo_node_id IN (
       SELECT c.id FROM geo_nodes c JOIN geo_nodes p ON c.parent_id = p.id
        WHERE c.name_ru = 'Жаркент' AND p.name_ru = 'Алматинская область');

DELETE FROM geo_nodes c
 USING geo_nodes p
 WHERE c.parent_id = p.id AND c.name_ru = 'Жаркент' AND p.name_ru = 'Алматинская область';

-- The city districts and the sub-10k towns never got their Kazakh or English
-- names, so a Kazakh reader saw a Russian label inside an otherwise Kazakh
-- interface. Fill them in; the store already falls back to name_ru, which is
-- exactly why this went unnoticed.
UPDATE geo_nodes c SET name_kk = v.kk, name_en = v.en
  FROM (VALUES
    ('Алатау','Алатау','Alatau'),          ('Алмалы','Алмалы','Almaly'),
    ('Ауэзов','Әуезов','Auezov'),          ('Бостандык','Бостандық','Bostandyk'),
    ('Жетысу','Жетісу','Zhetysu'),         ('Медеу','Медеу','Medeu'),
    ('Наурызбай','Наурызбай','Nauryzbai'), ('Турксиб','Түрксіб','Turksib')
  ) AS v(ru, kk, en)
  JOIN geo_nodes p ON p.name_ru = 'Алматы' AND p.level = 1
 WHERE c.parent_id = p.id AND c.name_ru = v.ru;

UPDATE geo_nodes c SET name_kk = v.kk, name_en = v.en
  FROM (VALUES
    ('Алматы','Алматы','Almaty'),   ('Байконур','Байқоңыр','Baikonyr'),
    ('Есиль','Есіл','Esil'),        ('Нура','Нұра','Nura'),
    ('Сарыарка','Сарыарқа','Saryarka')
  ) AS v(ru, kk, en)
  JOIN geo_nodes p ON p.name_ru = 'Астана' AND p.level = 1
 WHERE c.parent_id = p.id AND c.name_ru = v.ru;

UPDATE geo_nodes c SET name_kk = v.kk, name_en = v.en
  FROM (VALUES
    ('Абай','Абай','Abay'),            ('Аль-Фараби','Әл-Фараби','Al-Farabi'),
    ('Енбекши','Еңбекші','Enbekshi'),  ('Каратау','Қаратау','Karatau'),
    ('Туран','Тұран','Turan')
  ) AS v(ru, kk, en)
  JOIN geo_nodes p ON p.name_ru = 'Шымкент' AND p.level = 1
 WHERE c.parent_id = p.id AND c.name_ru = v.ru;

-- Towns that sit just below the 10,000 line but were already in the table.
UPDATE geo_nodes c SET name_kk = v.kk, name_en = v.en
  FROM (VALUES
    ('Серебрянск','Серебрянск','Serebryansk'), ('Казалинск','Қазалы','Kazaly'),
    ('Форт-Шевченко','Форт-Шевченко','Fort-Shevchenko'),
    ('Булаево','Бұлаев','Bulayevo'),           ('Каражал','Қаражал','Karazhal')
  ) AS v(ru, kk, en)
 WHERE c.country = 'KZ' AND c.level = 2 AND c.name_ru = v.ru AND c.name_kk = '';

-- The four legacy Kostanay towns carried a hand-numbered sort (0..3) which,
-- being the first ORDER BY key, split the list: everything imported with
-- sort = 0 came first and the hand-numbered rows trailed behind regardless of
-- size. Population is a better order for a settlement picker, so clear the
-- manual sort wherever a population now exists. Roots and regions keep theirs.
UPDATE geo_nodes SET sort = 0
 WHERE country = 'KZ' AND level = 2 AND population IS NOT NULL AND sort <> 0;

-- +goose Down
SELECT 1;
