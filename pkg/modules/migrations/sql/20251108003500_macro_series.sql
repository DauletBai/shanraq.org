-- +goose Up
-- Макроряды: то, из чего складываются курс и инфляция.
--
-- Одна таблица на все показатели, а не таблица на источник: рядов немного, они
-- одинаковой формы — период и число, — и сравнивать их между собой приходится
-- постоянно. Раскладывать их по отдельным таблицам значило бы писать join там,
-- где хватает WHERE.
--
-- period — первое число месяца для месячных рядов и первое января для годовых.
CREATE TABLE IF NOT EXISTS macro_series (
    code   text          NOT NULL,
    period date          NOT NULL,
    value  numeric(20,4) NOT NULL,
    PRIMARY KEY (code, period)
);
CREATE INDEX IF NOT EXISTS idx_macro_series_code ON macro_series (code, period);

-- +goose Down
DROP TABLE IF EXISTS macro_series;
