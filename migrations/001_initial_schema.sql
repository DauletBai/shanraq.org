-- Активируем расширение для генерации UUID, если еще не активировано
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица для хранения стран
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    population BIGINT NOT NULL,
    flag_code VARCHAR(10) -- Код страны для отображения флага (например, 'kz', 'us')
);

-- Таблица для категорий верхнего уровня (Наука, Спорт и т.д.)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Таблица для конкретных соревнований (подкатегории)
CREATE TABLE competitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    UNIQUE(name, year) -- Уникальная пара: название и год
);

-- Таблица для хранения наград
CREATE TABLE awards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    competition_id UUID NOT NULL REFERENCES competitions(id),
    country_id UUID NOT NULL REFERENCES countries(id),
    gold_medals INT DEFAULT 0,
    silver_medals INT DEFAULT 0,
    bronze_medals INT DEFAULT 0
);

-- Индексы для ускорения запросов
CREATE INDEX ON competitions (category_id);
CREATE INDEX ON awards (competition_id);
CREATE INDEX ON awards (country_id);

-- Добавляем немного тестовых данных
INSERT INTO countries (name, population, flag_code) VALUES
('Armenia', 2986000, 'am'),
('Ireland', 5200000, 'ie'),
('Kazakhstan', 20593000, 'kz'),
('Uzbekistan', 36362000, 'uz'),
('USA', 345427000, 'us'),
('China', 1419321000, 'cn');

-- 1. Добавляем категории
INSERT INTO categories (id, name) VALUES
('a4f0a1c0-81e1-413b-8258-1d21e2e7b1e1', 'Sport'),
('b5c1b2d1-92f2-424c-9369-2e32f3f8c2f2', 'Science');

-- 2. Добавляем соревнования
INSERT INTO competitions (category_id, name, year) VALUES
('a4f0a1c0-81e1-413b-8258-1d21e2e7b1e1', 'Olympic Games Boxing', 2024),
('b5c1b2d1-92f2-424c-9369-2e32f3f8c2f2', 'International Mathematical Olympiad', 2024);

-- 3. Добавляем награды
-- Получаем UUID стран и соревнований для удобства
DO $$
DECLARE
    kaz_id UUID;
    usa_id UUID;
    chn_id UUID;
    boxing_id UUID;
    math_id UUID;
BEGIN
    SELECT id INTO kaz_id FROM countries WHERE name = 'Kazakhstan';
    SELECT id INTO usa_id FROM countries WHERE name = 'USA';
    SELECT id INTO chn_id FROM countries WHERE name = 'China';
    SELECT id INTO boxing_id FROM competitions WHERE name = 'Olympic Games Boxing';
    SELECT id INTO math_id FROM competitions WHERE name = 'International Mathematical Olympiad';

    -- Награды по боксу
    INSERT INTO awards (competition_id, country_id, gold_medals, silver_medals, bronze_medals) VALUES
    (boxing_id, kaz_id, 1, 2, 5),
    (boxing_id, usa_id, 3, 1, 2);

    -- Награды по математике
    INSERT INTO awards (competition_id, country_id, gold_medals, silver_medals, bronze_medals) VALUES
    (math_id, chn_id, 6, 0, 0),
    (math_id, usa_id, 5, 1, 0);
END $$;