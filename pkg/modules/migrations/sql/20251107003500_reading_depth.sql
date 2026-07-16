-- +goose Up
-- Aggregate reading-depth funnel per article: how many readers reached 25/50/75
-- /100% of the article. "Start" (0%) is the article view count already tracked.
CREATE TABLE IF NOT EXISTS reading_depth (
    article_id uuid     NOT NULL,
    depth      smallint NOT NULL CHECK (depth IN (25, 50, 75, 100)),
    count      bigint   NOT NULL DEFAULT 0,
    PRIMARY KEY (article_id, depth)
);

-- +goose Down
DROP TABLE IF EXISTS reading_depth;
