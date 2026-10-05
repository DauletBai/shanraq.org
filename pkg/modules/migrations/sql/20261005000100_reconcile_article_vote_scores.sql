-- +goose Up
-- Older score caches can disagree with the votes that readers actually cast.
-- Reconcile existing rows, then keep the cache correct even when a vote is
-- removed by a cascading account deletion rather than the voting handler.
UPDATE articles a
SET score = COALESCE((
    SELECT SUM(v.value * v.weight) FROM article_votes v WHERE v.article_id = a.id
), 0)
WHERE a.score IS DISTINCT FROM COALESCE((
    SELECT SUM(v.value * v.weight) FROM article_votes v WHERE v.article_id = a.id
), 0);

-- +goose StatementBegin
CREATE FUNCTION refresh_article_vote_score() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
    affected_article uuid;
BEGIN
    IF TG_OP = 'INSERT' THEN
        affected_article := NEW.article_id;
    ELSE
        affected_article := OLD.article_id;
    END IF;

    UPDATE articles a
    SET score = COALESCE((
        SELECT SUM(v.value * v.weight)
        FROM article_votes v WHERE v.article_id = affected_article
    ), 0)
    WHERE a.id = affected_article;

    IF TG_OP = 'UPDATE' AND NEW.article_id IS DISTINCT FROM OLD.article_id THEN
        UPDATE articles a
        SET score = COALESCE((
            SELECT SUM(v.value * v.weight)
            FROM article_votes v WHERE v.article_id = NEW.article_id
        ), 0)
        WHERE a.id = NEW.article_id;
    END IF;
    RETURN NULL;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER article_vote_score_refresh
AFTER INSERT OR UPDATE OR DELETE ON article_votes
FOR EACH ROW EXECUTE FUNCTION refresh_article_vote_score();

-- +goose Down
DROP TRIGGER IF EXISTS article_vote_score_refresh ON article_votes;
DROP FUNCTION IF EXISTS refresh_article_vote_score();
