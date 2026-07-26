-- +goose Up
-- Short public bio shown at the top of a user's author page (e.g. the founder's
-- "CEO and creator of Shanraq…"). Empty = no bio.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS bio TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE auth_users DROP COLUMN IF EXISTS bio;
