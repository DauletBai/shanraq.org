-- +goose Up
-- Small settings the site assigns to itself.
--
-- Created for the IndexNow key: the protocol requires the key to sit as a file on
-- the domain and not to change, or search engines reject the submissions. Keeping it
-- in the config would demand a manual step at deployment, and generating it afresh
-- on every start would break the file already served. So the key is issued once and
-- lives in the database.
CREATE TABLE IF NOT EXISTS app_settings (
    name       text        NOT NULL PRIMARY KEY,
    value      text        NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS app_settings;
