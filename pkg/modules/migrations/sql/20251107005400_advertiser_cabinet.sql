-- +goose Up
-- Advertiser cabinet (Phase 0b MVP): a company account plus ad-placement orders.
-- No billing yet — orders are captured as 'pending_payment' until a payment
-- provider (Kaspi) is wired. Owned by a regular auth user (the responsible
-- person); one company per user for the MVP.
CREATE TABLE IF NOT EXISTS advertisers (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id      UUID NOT NULL UNIQUE REFERENCES auth_users(id) ON DELETE CASCADE,
    company_name  TEXT NOT NULL,
    bin           TEXT NOT NULL DEFAULT '',
    legal_form    TEXT NOT NULL DEFAULT '',
    address       TEXT NOT NULL DEFAULT '',
    website       TEXT NOT NULL DEFAULT '',
    industry      TEXT NOT NULL DEFAULT '',
    contact_name  TEXT NOT NULL DEFAULT '',
    contact_role  TEXT NOT NULL DEFAULT '',
    contact_phone TEXT NOT NULL DEFAULT '',
    contact_email TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ad_orders (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advertiser_id  UUID NOT NULL REFERENCES advertisers(id) ON DELETE CASCADE,
    title          TEXT NOT NULL,
    body           TEXT NOT NULL DEFAULT '',
    image_url      TEXT NOT NULL DEFAULT '',
    target_url     TEXT NOT NULL DEFAULT '',
    cta            TEXT NOT NULL DEFAULT '',
    placement      TEXT NOT NULL DEFAULT 'all',   -- all | articles | listings
    geo_region     TEXT NOT NULL DEFAULT '',      -- contextual geo (optional)
    rubric         TEXT NOT NULL DEFAULT '',      -- contextual rubric (optional)
    duration_days  INT  NOT NULL DEFAULT 7,
    pay_model      TEXT NOT NULL DEFAULT 'flat',  -- flat | cpm | cpc | cpa (future)
    price          BIGINT NOT NULL DEFAULT 0,     -- estimated total, tenge
    payment_method TEXT NOT NULL DEFAULT 'kaspi', -- kaspi | card | invoice
    status         TEXT NOT NULL DEFAULT 'pending_payment',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ad_orders_placement_chk CHECK (placement IN ('all','articles','listings')),
    CONSTRAINT ad_orders_status_chk CHECK (status IN ('pending_payment','active','paused','finished','rejected'))
);
CREATE INDEX IF NOT EXISTS idx_ad_orders_advertiser ON ad_orders (advertiser_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS ad_orders;
DROP TABLE IF EXISTS advertisers;
