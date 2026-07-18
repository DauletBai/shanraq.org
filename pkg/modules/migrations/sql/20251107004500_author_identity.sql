-- +goose Up
-- Author identity (level A): a verified phone number ties an account to a real
-- person (in KZ a SIM is registered to an IIN), and a real first/last name is
-- the byline shown on articles. Pseudonyms are not allowed for article authors.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS first_name TEXT NOT NULL DEFAULT '';
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS last_name  TEXT NOT NULL DEFAULT '';
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS phone      TEXT;
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS phone_verified_at TIMESTAMPTZ;

-- Short-lived SMS one-time codes for phone verification.
CREATE TABLE IF NOT EXISTS phone_verification_codes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    phone       TEXT NOT NULL,
    code_hash   TEXT NOT NULL,
    attempts    INT  NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS phone_codes_user_idx ON phone_verification_codes(user_id);

-- +goose Down
DROP TABLE IF EXISTS phone_verification_codes;
ALTER TABLE auth_users DROP COLUMN IF EXISTS phone_verified_at;
ALTER TABLE auth_users DROP COLUMN IF EXISTS phone;
ALTER TABLE auth_users DROP COLUMN IF EXISTS last_name;
ALTER TABLE auth_users DROP COLUMN IF EXISTS first_name;
