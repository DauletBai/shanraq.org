-- +goose Up
-- Accumulated reputation ("karma") an author has earned from readers' votes.
CREATE TABLE IF NOT EXISTS author_reputation (
    user_id UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    karma INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One vote per user per article. The voter's weight is snapshotted at vote
-- time (weight grows with the voter's own karma → anti-brigading).
CREATE TABLE IF NOT EXISTS article_votes (
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    value SMALLINT NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (article_id, user_id),
    CONSTRAINT article_votes_value_chk CHECK (value IN (-1, 1))
);

CREATE INDEX IF NOT EXISTS idx_article_votes_article ON article_votes (article_id);

-- +goose Down
DROP TABLE IF EXISTS article_votes;
DROP TABLE IF EXISTS author_reputation;
