-- +goose Up
-- Profile avatar for a user's cabinet. Stored as a URL to a processed image
-- (square, EXIF-stripped) served by the media module. Empty = no avatar yet.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS avatar_url TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE auth_users DROP COLUMN IF EXISTS avatar_url;
