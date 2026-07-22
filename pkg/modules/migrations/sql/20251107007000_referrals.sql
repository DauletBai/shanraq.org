-- +goose Up
-- Referral loop: the cheapest acquisition channel there is, because a user
-- brings a user at no ad cost. The reward is deliberately not cash and not a
-- points economy — on a shoestring budget those invite fraud and fake accounts,
-- which would poison a platform whose whole value is honesty. Instead the
-- reward is promotion inventory the platform already sells, granted only when
-- the invited user does something real and costly to fake: posts a genuine
-- listing.

-- Each user's own invite code. Generated lazily the first time they open the
-- invite page.
CREATE TABLE IF NOT EXISTS referral_codes (
    user_id    UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    code       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Who invited whom. One row per invited user (a person is referred once), so a
-- referred user cannot be farmed for repeat rewards.
CREATE TABLE IF NOT EXISTS referrals (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_id  UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    referred_id  UUID NOT NULL UNIQUE REFERENCES auth_users(id) ON DELETE CASCADE,
    status       TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'qualified')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    qualified_at TIMESTAMPTZ,
    -- A user cannot refer themselves.
    CONSTRAINT referral_not_self CHECK (referrer_id <> referred_id)
);

CREATE INDEX IF NOT EXISTS idx_referrals_referrer ON referrals (referrer_id);

-- Promotion-day credit as an append-only ledger: a grant is a positive delta,
-- spending it is a negative one, and the balance is their sum. Append-only so
-- the history of how credit was earned and used is auditable, matching the
-- moderation ledger's approach.
CREATE TABLE IF NOT EXISTS promo_credit_ledger (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    delta_days INTEGER NOT NULL,
    reason     TEXT NOT NULL,
    ref_id     UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_promo_credit_user ON promo_credit_ledger (user_id);

-- +goose Down
DROP TABLE IF EXISTS promo_credit_ledger;
DROP TABLE IF EXISTS referrals;
DROP TABLE IF EXISTS referral_codes;
