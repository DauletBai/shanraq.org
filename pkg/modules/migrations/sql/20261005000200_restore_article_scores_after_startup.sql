-- +goose Up
-- The old production startup cleanup overwrote the score cache after the
-- previous reconciliation migration had restored it. It no longer does so.
UPDATE articles a
SET score = COALESCE((
    SELECT SUM(v.value * v.weight) FROM article_votes v WHERE v.article_id = a.id
), 0)
WHERE a.score IS DISTINCT FROM COALESCE((
    SELECT SUM(v.value * v.weight) FROM article_votes v WHERE v.article_id = a.id
), 0);

-- +goose Down
-- Votes remain authoritative; rolling back must not erase the real score.
UPDATE articles a
SET score = COALESCE((
    SELECT SUM(v.value * v.weight) FROM article_votes v WHERE v.article_id = a.id
), 0)
WHERE a.score IS DISTINCT FROM COALESCE((
    SELECT SUM(v.value * v.weight) FROM article_votes v WHERE v.article_id = a.id
), 0);
