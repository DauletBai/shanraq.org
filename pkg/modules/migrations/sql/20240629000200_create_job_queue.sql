-- +goose Up
CREATE TYPE job_status AS ENUM ('pending', 'running', 'retry', 'done', 'failed');

CREATE TABLE IF NOT EXISTS job_queue (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES auth_users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    run_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    status job_status NOT NULL DEFAULT 'pending',
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS job_queue_status_run_at_idx ON job_queue (status, run_at);
CREATE INDEX IF NOT EXISTS job_queue_user_idx ON job_queue (user_id);

-- +goose Down
DROP TABLE IF EXISTS job_queue;
DROP TYPE IF EXISTS job_status;
