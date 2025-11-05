-- +goose Up
CREATE TABLE IF NOT EXISTS auth_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    key_hash TEXT NOT NULL UNIQUE,
    prefix TEXT NOT NULL,
    label TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS auth_api_keys_user_idx ON auth_api_keys(user_id);
CREATE INDEX IF NOT EXISTS auth_api_keys_prefix_idx ON auth_api_keys(prefix);

ALTER TABLE job_queue
    ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES auth_users(id) ON DELETE CASCADE;

UPDATE job_queue jq
SET user_id = u.id
FROM (
    SELECT id
    FROM auth_users
    ORDER BY created_at
    LIMIT 1
) AS u
WHERE jq.user_id IS NULL;

CREATE INDEX IF NOT EXISTS job_queue_user_idx ON job_queue(user_id);

INSERT INTO auth_api_keys (user_id, key_hash, prefix, label)
SELECT id, '75409192143765d27dfe16765579fb68439a2a6db8b66f3635d4a75d2789fcb8', 'sk_demo_op', 'Operator demo key'
FROM auth_users
WHERE email = 'operator@shanraq.org'
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM auth_api_keys WHERE prefix = 'sk_demo_op';
DROP TABLE IF EXISTS auth_api_keys;
DROP INDEX IF EXISTS job_queue_user_idx;
ALTER TABLE job_queue DROP COLUMN IF EXISTS user_id;
