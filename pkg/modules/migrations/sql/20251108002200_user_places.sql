-- +goose Up
-- Where the reader lives. One field, not four.
--
-- The form asks for country, region and city as a cascade, but only the chosen node
-- is stored: the republic and the region are derived from it by the geo_nodes tree,
-- whereas four separate columns would sooner or later start contradicting one
-- another — a city from one region beside a region from another.
--
-- This is for delivery: an author will be able to publish for their own settlement,
-- district or region, and it will reach the people it concerns. The place cannot be
-- determined by IP: the DB-IP Lite database we run tells only countries apart, and
-- city-level databases in Kazakhstan put half the country in Almaty because mobile
-- traffic leaves through a handful of gateways. So we ask the person and let them
-- change it in one click.
--
-- A separate table rather than a column in auth_users: geography is the articles
-- module's business, and the auth module has no need to know about it.
CREATE TABLE IF NOT EXISTS user_places (
    user_id UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    geo_node_id UUID NOT NULL REFERENCES geo_nodes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_places_node ON user_places (geo_node_id);

-- +goose Down
DROP TABLE IF EXISTS user_places;
