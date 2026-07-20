-- +goose Up
-- Turn ad_orders into a bookable placement: a zone (home / rubric / articles /
-- realestate / all), a concrete date range, and an optional exclusive hold on
-- the zone's rotation. This is what lets the cabinet auto-check free slots.
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS starts_at DATE;
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS ends_at   DATE;
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS exclusive BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS lang      TEXT NOT NULL DEFAULT '';

-- Re-map legacy placement values before tightening the constraint. The old
-- CHECK has to go first, otherwise it rejects the new zone names mid-remap.
ALTER TABLE ad_orders DROP CONSTRAINT IF EXISTS ad_orders_placement_chk;
UPDATE ad_orders SET placement = 'realestate' WHERE placement = 'listings';
UPDATE ad_orders SET placement = 'articles'   WHERE placement NOT IN ('home','rubric','articles','realestate','all');

-- Backfill a date range for rows created before booking existed.
UPDATE ad_orders
   SET starts_at = COALESCE(starts_at, created_at::date),
       ends_at   = COALESCE(ends_at, created_at::date + (GREATEST(duration_days,1) || ' days')::interval)
 WHERE starts_at IS NULL OR ends_at IS NULL;

ALTER TABLE ad_orders ADD CONSTRAINT ad_orders_placement_chk
    CHECK (placement IN ('home','rubric','articles','realestate','all'));

-- Availability lookups scan booked rows by zone and overlapping dates.
CREATE INDEX IF NOT EXISTS idx_ad_orders_booking
    ON ad_orders (placement, starts_at, ends_at)
    WHERE status IN ('pending_payment', 'active');

-- +goose Down
DROP INDEX IF EXISTS idx_ad_orders_booking;
ALTER TABLE ad_orders DROP CONSTRAINT IF EXISTS ad_orders_placement_chk;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS lang;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS exclusive;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS ends_at;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS starts_at;
