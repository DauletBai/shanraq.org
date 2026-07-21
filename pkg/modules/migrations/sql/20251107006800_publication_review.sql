-- +goose Up
-- Publishing used to be a single step: the author pressed the button and the
-- article went live unread. Now it passes a rules check first, so an article
-- either goes out or comes back with the specific rules it failed.
--
-- Two new statuses sit between draft and published:
--   review     — submitted, waiting on the checker (or on a human, if the
--                checker is unavailable). Never publicly visible.
--   needs_work — returned to the author with findings attached.
-- The original constraint is named ..._chk, not ..._check; dropping only the
-- latter would leave the old one enforcing the old set, and 'review' would be
-- rejected at runtime with the migration reporting success.
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_chk;
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_check;
ALTER TABLE articles ADD CONSTRAINT articles_status_check
    CHECK (status IN ('draft', 'review', 'needs_work', 'published', 'archived'));

ALTER TABLE articles ADD COLUMN IF NOT EXISTS submitted_at TIMESTAMPTZ;
ALTER TABLE articles ADD COLUMN IF NOT EXISTS reviewed_at  TIMESTAMPTZ;

-- The queue a reviewer works: oldest submission first.
CREATE INDEX IF NOT EXISTS idx_articles_review ON articles (submitted_at)
    WHERE status = 'review';

-- One decision can fail several rules at once, and the author is entitled to
-- see every one of them rather than the first. Findings hang off the ledger
-- entry, so an approval and a rejection are recorded the same way.
CREATE TABLE IF NOT EXISTS moderation_findings (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action_id UUID NOT NULL REFERENCES moderation_actions(id) ON DELETE CASCADE,
    rule_code TEXT NOT NULL,
    -- 'block' returns the article; 'warn' is advice that does not stop it.
    severity  TEXT NOT NULL DEFAULT 'block' CHECK (severity IN ('block', 'warn')),
    -- The passage the finding refers to, so "unsourced claim" is actionable
    -- rather than an accusation the author has to go hunting for.
    quote     TEXT NOT NULL DEFAULT '',
    note      TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_findings_action ON moderation_findings (action_id);

-- +goose Down
DROP TABLE IF EXISTS moderation_findings;
DROP INDEX IF EXISTS idx_articles_review;
ALTER TABLE articles DROP COLUMN IF EXISTS reviewed_at;
ALTER TABLE articles DROP COLUMN IF EXISTS submitted_at;
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_check;
ALTER TABLE articles ADD CONSTRAINT articles_status_check
    CHECK (status IN ('draft', 'published', 'archived'));
