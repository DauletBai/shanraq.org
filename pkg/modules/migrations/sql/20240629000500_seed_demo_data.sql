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

INSERT INTO auth_users (id, email, password_hash, role)
VALUES
    (gen_random_uuid(), 'admin@shanraq.org', crypt('admin123', gen_salt('bf')), 'admin'),
    (gen_random_uuid(), 'operator@shanraq.org', crypt('operator123', gen_salt('bf')), 'operator')
ON CONFLICT (email) DO NOTHING;

INSERT INTO job_queue (id, name, payload, run_at, max_attempts, status, attempts, updated_at)
VALUES
    (gen_random_uuid(), 'seed_email_digest', '{"email":"digest@shanraq.org"}'::jsonb, NOW() + INTERVAL '10 minutes', 3, 'pending', 0, NOW()),
    (gen_random_uuid(), 'seed_report', '{"report":"daily"}'::jsonb, NOW() - INTERVAL '5 minutes', 3, 'running', 1, NOW()),
    (gen_random_uuid(), 'seed_reconciliation', '{"batch":42}'::jsonb, NOW() + INTERVAL '2 minutes', 5, 'retry', 2, NOW()),
    (gen_random_uuid(), 'seed_failed_import', '{"source":"s3"}'::jsonb, NOW() - INTERVAL '20 minutes', 3, 'failed', 3, NOW()),
    (gen_random_uuid(), 'seed_cleanup', '{"scope":"logs"}'::jsonb, NOW() - INTERVAL '40 minutes', 3, 'done', 1, NOW());

-- +goose Down
DELETE FROM job_queue WHERE name LIKE 'seed_%';
DELETE FROM auth_users WHERE email IN ('admin@shanraq.org', 'operator@shanraq.org');
DELETE FROM framework_about WHERE headline = 'Shanraq Console' AND subheadline = 'A Go-first framework for resilient backends.';
