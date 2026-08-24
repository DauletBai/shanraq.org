-- +goose Up
-- Reader-submitted typo corrections, and what became of each one.
--
-- A published article is its author's text, and nothing here rewrites it wholesale:
-- a correction names one chapter, one sentence and one word, and the only edit the
-- server will make is that one word inside that one sentence. The row keeps the
-- whole claim so an applied fix can be explained afterwards — and reversed by
-- hand if it was wrong.
--
-- lang is stored because an article has three language versions and a typo lives
-- in exactly one of them. Without it a correction to the Russian text could be
-- applied to the Kazakh.
CREATE TABLE IF NOT EXISTS article_corrections (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id  UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    lang        TEXT NOT NULL,
    reporter_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    -- What the reader pointed at. chapter may be blank: a short piece has no
    -- headings, and demanding one would block a correction to its first line.
    chapter  TEXT NOT NULL DEFAULT '',
    sentence TEXT NOT NULL,
    word     TEXT NOT NULL,
    -- pending  — recorded, not yet decided (the checker was unreachable)
    -- applied  — the word was replaced in the article
    -- rejected — the checker found nothing to fix
    -- failed   — the claim did not survive the server's own checks
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'applied', 'rejected', 'failed')),
    -- What the word became, on an applied fix. Kept so the edit is reversible
    -- and so the author's notification can quote both sides of it.
    fixed TEXT NOT NULL DEFAULT '',
    -- Why it was refused, in the reader's own language. A refusal without a
    -- reason is indistinguishable from the site being broken.
    reason     TEXT NOT NULL DEFAULT '',
    decided_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- The reporter index answers the rate-limit question on every submission: how
-- many has this person sent today.
CREATE INDEX IF NOT EXISTS idx_article_corrections_reporter
    ON article_corrections (reporter_id, created_at DESC);
-- And this one answers the author's question: what has been changed in my piece.
CREATE INDEX IF NOT EXISTS idx_article_corrections_article
    ON article_corrections (article_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS article_corrections;
