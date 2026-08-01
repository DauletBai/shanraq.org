-- +goose Up
-- Trilingual listings (KZ/RU/EN), like articles — the platform's flagship,
-- mandatory feature. The old single title/description stay as a base/fallback
-- so existing readers and code paths never see a blank.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS title_kz       TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS title_ru       TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS title_en       TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS description_kz TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS description_ru TEXT NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS description_en TEXT NOT NULL DEFAULT '';

-- Backfill existing single-language listings into all three so nothing is blank.
UPDATE listings SET title_ru = title, title_kz = title, title_en = title
    WHERE title <> '' AND title_ru = '';
UPDATE listings SET description_ru = description, description_kz = description, description_en = description
    WHERE description <> '' AND description_ru = '';

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS title_kz;
ALTER TABLE listings DROP COLUMN IF EXISTS title_ru;
ALTER TABLE listings DROP COLUMN IF EXISTS title_en;
ALTER TABLE listings DROP COLUMN IF EXISTS description_kz;
ALTER TABLE listings DROP COLUMN IF EXISTS description_ru;
ALTER TABLE listings DROP COLUMN IF EXISTS description_en;
