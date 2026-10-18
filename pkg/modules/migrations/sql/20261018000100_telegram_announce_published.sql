-- +goose Up
-- The editor calls IndexNow after publishing, but editorial SQL imports do not
-- pass through that Go handler. Queue Telegram from the committed article state
-- so those imports cannot silently omit the channel announcement.
CREATE UNIQUE INDEX IF NOT EXISTS job_queue_telegram_article_once
    ON job_queue ((payload->>'article_id'))
    WHERE name = 'syndicate_telegram' AND payload ? 'article_id';

CREATE OR REPLACE FUNCTION enqueue_published_article_telegram() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        IF OLD.status = 'published' THEN
            RETURN NULL;
        END IF;
    END IF;

    -- This is a deferred trigger: inspect the final state of the transaction.
    -- Course installers attach article_series_items after inserting articles;
    -- their lessons must never flood the news channel.
    IF NOT EXISTS (SELECT 1 FROM articles WHERE id = NEW.id AND status = 'published')
       OR EXISTS (SELECT 1 FROM article_series_items WHERE article_id = NEW.id) THEN
        RETURN NULL;
    END IF;

    INSERT INTO job_queue (id, name, payload, run_at, max_attempts)
    VALUES (gen_random_uuid(), 'syndicate_telegram',
            jsonb_build_object('article_id', NEW.id::text), now(), 3)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$;

CREATE CONSTRAINT TRIGGER articles_telegram_on_publish
    AFTER INSERT OR UPDATE OF status ON articles
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW
    WHEN (NEW.status = 'published')
    EXECUTE FUNCTION enqueue_published_article_telegram();

-- +goose Down
DROP TRIGGER IF EXISTS articles_telegram_on_publish ON articles;
DROP FUNCTION IF EXISTS enqueue_published_article_telegram();
DROP INDEX IF EXISTS job_queue_telegram_article_once;
