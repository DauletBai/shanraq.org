-- +goose Up
-- Listings became editable, so "when was this last touched" is now a real
-- question a buyer asks: a price edited yesterday means something different
-- from one posted in spring and never revisited. Existing rows inherit their
-- creation time, which is the truth for them — they have never been edited.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS updated_at timestamptz;
UPDATE listings SET updated_at = created_at WHERE updated_at IS NULL;
ALTER TABLE listings ALTER COLUMN updated_at SET DEFAULT now();
ALTER TABLE listings ALTER COLUMN updated_at SET NOT NULL;

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS updated_at;
