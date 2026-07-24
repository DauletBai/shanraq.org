-- +goose Up
-- Patronymic ("отчество"). Optional, but stored for everyone so two users with
-- the same first + last name can be told apart. first_name / last_name already
-- exist; they become required at registration.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS middle_name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE auth_users DROP COLUMN IF EXISTS middle_name;
