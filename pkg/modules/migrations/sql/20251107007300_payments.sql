-- +goose Up
-- Payment scaffolding, provider-agnostic. A payment records the intent to
-- charge for something (an ad order, a listing promotion) and its lifecycle.
-- The actual provider — an aggregator or bank acquirer — is wired in later
-- behind an interface; this table is what every provider writes through, so
-- swapping providers never touches the order flow.
CREATE TABLE IF NOT EXISTS payments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- What is being paid for. target_id points at the ad_order / listing.
    kind         TEXT NOT NULL CHECK (kind IN ('ad_order', 'listing_promo')),
    target_id    UUID NOT NULL,
    amount       BIGINT NOT NULL CHECK (amount >= 0),
    currency     TEXT NOT NULL DEFAULT 'KZT',
    -- Which provider is handling it, and its reference in that provider. Empty
    -- until a provider is engaged.
    provider     TEXT NOT NULL DEFAULT '',
    provider_ref TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'paid', 'failed', 'expired')),
    -- An unpaid payment holds its order's slot only until this moment, so a
    -- checkout nobody completes cannot occupy inventory forever.
    expires_at   TIMESTAMPTZ,
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_target  ON payments (kind, target_id);
CREATE INDEX IF NOT EXISTS idx_payments_pending ON payments (expires_at) WHERE status = 'pending';

-- An ad order whose payment expires or is cancelled must stop occupying its
-- slot, so add the terminal states the booking index already needs to exclude.
ALTER TABLE ad_orders DROP CONSTRAINT IF EXISTS ad_orders_status_chk;
ALTER TABLE ad_orders DROP CONSTRAINT IF EXISTS ad_orders_status_check;
ALTER TABLE ad_orders ADD CONSTRAINT ad_orders_status_chk
    CHECK (status IN ('pending_payment', 'active', 'paused', 'finished', 'rejected', 'expired', 'cancelled'));

-- +goose Down
ALTER TABLE ad_orders DROP CONSTRAINT IF EXISTS ad_orders_status_chk;
ALTER TABLE ad_orders ADD CONSTRAINT ad_orders_status_chk
    CHECK (status IN ('pending_payment', 'active', 'paused', 'finished', 'rejected'));
DROP TABLE IF EXISTS payments;
