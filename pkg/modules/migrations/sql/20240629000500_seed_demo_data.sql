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

WITH admin_user AS (
    INSERT INTO auth_users (id, email, password_hash, role, password_reset_required)
    VALUES (
        gen_random_uuid(),
        'admin@shanraq.org',
        '$2a$10$9rh2mzsokAHtKJBw2DIO1.4P2ghZH0RGWvLeMhPw6W2IvKDlvGCZq',
        'admin',
        TRUE
    )
    ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
    RETURNING id
), operator_user AS (
    INSERT INTO auth_users (id, email, password_hash, role, password_reset_required)
    VALUES (
        gen_random_uuid(),
        'operator@shanraq.org',
        '$2a$10$.qzJRMYmQyaGP4a4msVXRObeTRnJLK4POhDI3NMHZqm5p0BRVdxj.',
        'operator',
        TRUE
    )
    ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
    RETURNING id
), seed_jobs AS (
    INSERT INTO job_queue (id, user_id, name, payload, run_at, max_attempts, status, attempts, updated_at)
    SELECT gen_random_uuid(), (SELECT id FROM admin_user), 'seed_email_digest', '{"email":"digest@shanraq.org"}'::jsonb, NOW() + INTERVAL '10 minutes', 3, 'pending', 0, NOW()
    UNION ALL
    SELECT gen_random_uuid(), (SELECT id FROM operator_user), 'seed_report', '{"report":"daily"}'::jsonb, NOW() - INTERVAL '5 minutes', 3, 'running', 1, NOW()
    UNION ALL
    SELECT gen_random_uuid(), (SELECT id FROM operator_user), 'seed_reconciliation', '{"batch":42}'::jsonb, NOW() + INTERVAL '2 minutes', 5, 'retry', 2, NOW()
    UNION ALL
    SELECT gen_random_uuid(), (SELECT id FROM operator_user), 'seed_failed_import', '{"source":"s3"}'::jsonb, NOW() - INTERVAL '20 minutes', 3, 'failed', 3, NOW()
    UNION ALL
    SELECT gen_random_uuid(), (SELECT id FROM admin_user), 'seed_cleanup', '{"scope":"logs"}'::jsonb, NOW() - INTERVAL '40 minutes', 3, 'done', 1, NOW()
    RETURNING 1
)
INSERT INTO auth_api_keys (user_id, key_hash, prefix, label)
SELECT id, '75409192143765d27dfe16765579fb68439a2a6db8b66f3635d4a75d2789fcb8', 'sk_demo_op', 'Operator demo key'
FROM operator_user
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM job_queue WHERE name LIKE 'seed_%';
DELETE FROM auth_users WHERE email IN ('admin@shanraq.org', 'operator@shanraq.org');
DELETE FROM auth_api_keys WHERE prefix = 'sk_demo_op';
DELETE FROM framework_about WHERE headline = 'Shanraq Console' AND subheadline = 'A Go-first framework for resilient backends.';
