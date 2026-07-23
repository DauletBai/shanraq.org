-- +goose Up
-- Real-estate agents. Any registered user can become an agent for free: they
-- fill a short profile and their listings then carry an "Agent" badge that links
-- to a public page with all of their listings. One profile per user.
CREATE TABLE IF NOT EXISTS re_agents (
    user_id    UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    agency     TEXT NOT NULL DEFAULT '',
    phone      TEXT NOT NULL DEFAULT '',
    about      TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS re_agents;
