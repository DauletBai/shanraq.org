-- +goose Up
-- SECURITY: an earlier migration (20251030000700) seeded a demonstration API
-- key ('sk_demo_op…', label "Operator demo key") whose hash shipped in the
-- public repository — anyone could use it to reach /jobs/ in any deployment
-- that ran that migration. Revoke and delete it. Fresh installs run this right
-- after the seed, so the demo key is never usable post-migration.
UPDATE auth_api_keys SET revoked_at = NOW()
    WHERE prefix = 'sk_demo_op' AND revoked_at IS NULL;
DELETE FROM auth_api_keys WHERE prefix = 'sk_demo_op';

-- +goose Down
-- Intentionally NOT recreated: a leaked demo credential must not come back.
SELECT 1;
