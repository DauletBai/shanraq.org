-- +goose Up
INSERT INTO framework_about (headline, subheadline, feature_one, feature_two, feature_three)
VALUES (
    'Shanraq Console',
    'A Go-first framework for resilient backends.',
    'PostgreSQL-native data layer with migrations built-in.',
    'Composable module system for HTTP, workers, and observability.',
    'Cloud-ready tooling: Docker, Prometheus telemetry, structured logging.'
)
ON CONFLICT DO NOTHING;

INSERT INTO auth_users (id, email, password_hash, role, password_reset_required)
VALUES (
    gen_random_uuid(),
    'admin@shanraq.org',
    '$2a$10$9rh2mzsokAHtKJBw2DIO1.4P2ghZH0RGWvLeMhPw6W2IvKDlvGCZq',
    'admin',
    TRUE
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO auth_users (id, email, password_hash, role, password_reset_required)
VALUES (
    gen_random_uuid(),
    'operator@shanraq.org',
    '$2a$10$.qzJRMYmQyaGP4a4msVXRObeTRnJLK4POhDI3NMHZqm5p0BRVdxj.',
    'operator',
    TRUE
)
ON CONFLICT (email) DO NOTHING;

ALTER TABLE job_queue
    ADD COLUMN IF NOT EXISTS user_id UUID;

ALTER TABLE job_queue DROP CONSTRAINT IF EXISTS job_queue_user_id_fkey;
ALTER TABLE job_queue
    ADD CONSTRAINT job_queue_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES auth_users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS job_queue_user_idx ON job_queue(user_id);

INSERT INTO job_queue (id, user_id, name, payload, run_at, max_attempts, status, attempts, updated_at)
VALUES
    (gen_random_uuid(), (SELECT id FROM auth_users WHERE email = 'admin@shanraq.org'), 'seed_email_digest', '{"email":"digest@shanraq.org"}'::jsonb, NOW() + INTERVAL '10 minutes', 3, 'pending', 0, NOW()),
    (gen_random_uuid(), (SELECT id FROM auth_users WHERE email = 'operator@shanraq.org'), 'seed_report', '{"report":"daily"}'::jsonb, NOW() - INTERVAL '5 minutes', 3, 'running', 1, NOW()),
    (gen_random_uuid(), (SELECT id FROM auth_users WHERE email = 'operator@shanraq.org'), 'seed_reconciliation', '{"batch":42}'::jsonb, NOW() + INTERVAL '2 minutes', 5, 'retry', 2, NOW()),
    (gen_random_uuid(), (SELECT id FROM auth_users WHERE email = 'operator@shanraq.org'), 'seed_failed_import', '{"source":"s3"}'::jsonb, NOW() - INTERVAL '20 minutes', 3, 'failed', 3, NOW()),
    (gen_random_uuid(), (SELECT id FROM auth_users WHERE email = 'admin@shanraq.org'), 'seed_cleanup', '{"scope":"logs"}'::jsonb, NOW() - INTERVAL '40 minutes', 3, 'done', 1, NOW())
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM job_queue WHERE name LIKE 'seed_%';
DELETE FROM auth_users WHERE email IN ('admin@shanraq.org', 'operator@shanraq.org');
DELETE FROM framework_about WHERE headline = 'Shanraq Console' AND subheadline = 'A Go-first framework for resilient backends.';
