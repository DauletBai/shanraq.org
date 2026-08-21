-- +goose Up
-- Где живёт читатель. Одно поле, а не четыре.
--
-- Форма спрашивает страну, область и город каскадом, но хранится только
-- выбранный узел: республику и область из него выводит дерево geo_nodes, а
-- четыре отдельные колонки рано или поздно начали бы противоречить друг другу —
-- город из одной области рядом с областью из другой.
--
-- Это нужно для доставки: автор сможет опубликовать материал для своего
-- посёлка, района или области, и он дойдёт до тех, кого касается. Определить
-- место по IP нельзя: установленная у нас база DB-IP Lite различает только
-- страну, а городские базы в Казахстане показывают половину страны в Алматы,
-- потому что мобильный трафик выходит через несколько шлюзов. Поэтому
-- спрашиваем у человека и даём поменять в один клик.
--
-- Отдельная таблица, а не колонка в auth_users: география — предмет модуля
-- статей, и модулю авторизации о ней знать незачем.
CREATE TABLE IF NOT EXISTS user_places (
    user_id UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    geo_node_id UUID NOT NULL REFERENCES geo_nodes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_places_node ON user_places (geo_node_id);

-- +goose Down
DROP TABLE IF EXISTS user_places;
