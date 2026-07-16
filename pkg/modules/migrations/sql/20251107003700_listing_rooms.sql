-- +goose Up
-- Per-room breakdown: kind (kitchen/bedroom/living/…), floor area and a note,
-- stored as a JSON array so a listing can describe each room.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS room_specs jsonb NOT NULL DEFAULT '[]'::jsonb;

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS room_specs;
