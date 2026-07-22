-- +goose Up
-- Add the banner format to ad orders. Price is now format × surface, and the
-- rotation capacity is per format (a big sidebar half-page is one exclusive
-- slot; the small rectangle rotates three), so serving and slot-counting need
-- to know it.
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS format TEXT NOT NULL DEFAULT 'rectangle';

-- Old orders were the single sidebar unit, which is the rectangle.
UPDATE ad_orders SET format = 'rectangle' WHERE format = '' OR format IS NULL;

CREATE INDEX IF NOT EXISTS idx_ad_orders_format ON ad_orders (format);

-- +goose Down
DROP INDEX IF EXISTS idx_ad_orders_format;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS format;
