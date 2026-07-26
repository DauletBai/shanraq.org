-- +goose Up
-- Make the author bio trilingual (KZ/RU/EN), like articles. The previous single
-- `bio` was Russian by convention, so it moves into bio_ru.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS bio_kz TEXT NOT NULL DEFAULT '';
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS bio_ru TEXT NOT NULL DEFAULT '';
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS bio_en TEXT NOT NULL DEFAULT '';
UPDATE auth_users SET bio_ru = bio WHERE bio <> '' AND bio_ru = '';
ALTER TABLE auth_users DROP COLUMN IF EXISTS bio;

-- +goose Down
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS bio TEXT NOT NULL DEFAULT '';
UPDATE auth_users SET bio = bio_ru WHERE bio_ru <> '';
ALTER TABLE auth_users DROP COLUMN IF EXISTS bio_kz;
ALTER TABLE auth_users DROP COLUMN IF EXISTS bio_ru;
ALTER TABLE auth_users DROP COLUMN IF EXISTS bio_en;
