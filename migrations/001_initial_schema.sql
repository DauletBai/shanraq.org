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
