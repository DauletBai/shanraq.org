-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    body       TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'published',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT comments_status_chk CHECK (status IN ('published', 'hidden'))
);
CREATE INDEX IF NOT EXISTS idx_comments_article ON comments (article_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS comments;
