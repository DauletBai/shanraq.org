-- +goose Up
-- Reader moderation for articles.
--
-- The published rules have always said that offensive and unlawful material is
-- "moderated by our readers and by our AI model". For listings that was true —
-- a report button, a threshold, an auto-hide. For articles neither half existed:
-- no report route, no report button, nothing. With the pre-publication checker
-- off and authors publishing directly, readers are now the only moderation
-- there is, so the promise has to become code.
--
-- weight is stored on the report rather than derived at read time, because it
-- is the weight *at the moment of reporting*: a reader whose credibility later
-- changes should not retroactively rewrite decisions already taken.
CREATE TABLE IF NOT EXISTS article_reports (
    article_id  UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    reporter_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    reason      TEXT NOT NULL DEFAULT '',
    weight      INTEGER NOT NULL DEFAULT 1 CHECK (weight >= 0),
    -- Set when a hidden article is restored: the reports that hid it were wrong,
    -- and a reader with a habit of being wrong stops moving the threshold.
    dismissed   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (article_id, reporter_id)
);

-- The reporter index answers one question on every report: how often has this
-- person been wrong before.
CREATE INDEX IF NOT EXISTS idx_article_reports_reporter ON article_reports (reporter_id) WHERE dismissed;

-- 'flagged' is hidden-pending-review, the same state a reported listing enters.
-- It is not deletion: the article keeps its row, its author sees it, and it
-- comes back if the decision is overturned.
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_chk;
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_check;
ALTER TABLE articles ADD CONSTRAINT articles_status_check
    CHECK (status IN ('draft', 'review', 'needs_work', 'published', 'flagged', 'archived'));

-- A decision taken by readers is not a staff decision and not an AI decision,
-- and the ledger's own rule is that these are never presented as the same
-- thing. It needs its own name.
ALTER TABLE moderation_actions DROP CONSTRAINT IF EXISTS moderation_actions_actor_kind_check;
ALTER TABLE moderation_actions ADD CONSTRAINT moderation_actions_actor_kind_check
    CHECK (actor_kind IN ('agent', 'human', 'readers'));

-- +goose Down
ALTER TABLE moderation_actions DROP CONSTRAINT IF EXISTS moderation_actions_actor_kind_check;
ALTER TABLE moderation_actions ADD CONSTRAINT moderation_actions_actor_kind_check
    CHECK (actor_kind IN ('agent', 'human'));
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_check;
ALTER TABLE articles ADD CONSTRAINT articles_status_check
    CHECK (status IN ('draft', 'review', 'needs_work', 'published', 'archived'));
DROP TABLE IF EXISTS article_reports;
