-- +goose Up
-- Settlements around the ERG production sites, so an enterprise can address its
-- own workers' town rather than the nearest oblast centre.
--
-- The earlier sweep (20251107006300) took every settlement above roughly ten
-- thousand residents. That threshold is exactly why these are missing: a mine's
-- village is small, and it is small because the mine is the only thing there.
-- Those are the places least covered by anyone else and the ones a works
-- newsletter actually reaches.
--
-- Each is taken from an address published by the enterprise itself:
--
--   Kerege­tas limestone mine   -- Pavlodar Region, Bayanaul District,
--                                 Maikain settlement, Ushkulun village
--   Krasnooktyabrskoye bauxite  -- Kostanay Region, Lisakovsk, Oktyabrsky
--                                 settlement
--
-- Bayanaul is added as the district seat: it is where the district's offices
-- and much of its housing are, and a Maikain reader looks there next.
--
-- Coordinates and Kazakh names come from Wikipedia/Wikidata. Population is left
-- empty where no recent figure was found rather than filled with a guess: the
-- column is nullable, and an absent number is honest where an invented one is
-- not.
--
-- Slugs are normally left to the application, which transliterates them at
-- start-up. Oktyabrsky is the exception. Two dozen city districts of that name
-- already hold the short forms, so the automatic fallback would name this one
-- "kostanaiskaya-oblast-oktyabr-skii" -- correct, unambiguous and unusable in a
-- letter. It is given "oktyabrskiy-lisakovsk" here instead, which says which
-- Oktyabrsky it is. The application only fills slugs that are empty, so this
-- one survives.

CREATE TEMP TABLE erg_settlements (
    code TEXT, region TEXT, kind TEXT,
    name_ru TEXT, name_kk TEXT, name_en TEXT,
    population INTEGER, population_year SMALLINT,
    lat DOUBLE PRECISION, lng DOUBLE PRECISION,
    slug TEXT
) ON COMMIT DROP;

INSERT INTO erg_settlements VALUES
('kz-pavlodar-maikain','Павлодарская область','town','Майкаин','Майқайың','Maikain',7459,2021,51.4625,75.79861,NULL),
('kz-pavlodar-ushkulun','Павлодарская область','village','Ушкулун','Үшқұлын','Ushkulun',NULL,NULL,51.449845,75.707973,NULL),
('kz-pavlodar-bayanaul','Павлодарская область','village','Баянаул','Баянауыл','Bayanaul',NULL,NULL,50.78889,75.69556,NULL),
('kz-kostanay-oktyabrskiy','Костанайская область','town','Октябрьский','Октябрьский','Oktyabrskiy',NULL,NULL,52.63556,62.70389,'oktyabrskiy-lisakovsk');

-- Existing rows keep their identity: a settlement already present is completed,
-- never duplicated, so no article or listing loses the place it points at.
UPDATE geo_nodes c
   SET name_kk         = s.name_kk,
       name_en         = s.name_en,
       kind            = s.kind,
       population      = COALESCE(s.population, c.population),
       population_year = COALESCE(s.population_year, c.population_year),
       lat             = COALESCE(s.lat, c.lat),
       lng             = COALESCE(s.lng, c.lng),
       slug            = COALESCE(s.slug, c.slug)
  FROM erg_settlements s
  JOIN geo_nodes p ON p.country = 'KZ' AND p.level = 1 AND p.name_ru = s.region
 WHERE c.parent_id = p.id AND c.name_ru = s.name_ru;

INSERT INTO geo_nodes (code, parent_id, parent_code, country, level, kind,
                       name_ru, name_kk, name_en, population, population_year, lat, lng, slug, sort)
SELECT s.code, p.id, p.code, 'KZ', 2, s.kind,
       s.name_ru, s.name_kk, s.name_en, s.population, s.population_year, s.lat, s.lng, s.slug, 0
  FROM erg_settlements s
  JOIN geo_nodes p ON p.country = 'KZ' AND p.level = 1 AND p.name_ru = s.region
 WHERE NOT EXISTS (
       SELECT 1 FROM geo_nodes c WHERE c.parent_id = p.id AND c.name_ru = s.name_ru);

-- +goose Down
DELETE FROM geo_nodes WHERE code IN (
    'kz-pavlodar-maikain','kz-pavlodar-ushkulun',
    'kz-pavlodar-bayanaul','kz-kostanay-oktyabrskiy');
