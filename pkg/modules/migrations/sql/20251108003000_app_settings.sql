-- +goose Up
-- Мелкие настройки, которые сайт назначает себе сам.
--
-- Заводится ради ключа IndexNow: протокол требует, чтобы ключ лежал файлом на
-- домене и не менялся, иначе поисковики отвергают заявки. Держать его в
-- конфиге значит требовать ручного шага при разворачивании, а генерировать
-- заново при каждом запуске — ломать уже отданный файл. Поэтому ключ
-- выписывается один раз и живёт в базе.
CREATE TABLE IF NOT EXISTS app_settings (
    name       text        NOT NULL PRIMARY KEY,
    value      text        NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS app_settings;
