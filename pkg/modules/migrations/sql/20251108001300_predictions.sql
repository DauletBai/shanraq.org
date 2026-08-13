-- +goose Up
-- The prediction ledger: every forecast the site makes, recorded with the date
-- it was made and later marked as it actually turned out — including wrong.
--
-- The point is the failures. An analyst who is never seen to be wrong is not
-- accurate, only unaudited, and a reader knows it. Publishing a score that can
-- go down is the one claim to credibility that cannot be faked, and it is free.
CREATE TABLE IF NOT EXISTS predictions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- The article the forecast was made in. ON DELETE SET NULL, not CASCADE:
    -- a prediction outlives the piece it came from, and deleting the article
    -- must never quietly delete a miss from the record.
    article_id  UUID REFERENCES articles(id) ON DELETE SET NULL,
    made_on     DATE NOT NULL,
    -- When it becomes fair to judge. Null means "no deadline stated", which is
    -- honest for a structural forecast and is shown as such.
    horizon     DATE,
    status      TEXT NOT NULL DEFAULT 'open',
    resolved_on DATE,
    -- Evidence for the verdict: the report or filing that settled it.
    source_url  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT predictions_status_chk
        CHECK (status IN ('open', 'hit', 'miss', 'partial')),
    -- An open forecast has no verdict date, and a settled one must have it.
    -- Without this the scoreboard could count a resolution nobody dated.
    CONSTRAINT predictions_resolved_chk
        CHECK ((status = 'open') = (resolved_on IS NULL))
);

CREATE INDEX IF NOT EXISTS predictions_status_idx ON predictions (status, made_on DESC);
CREATE INDEX IF NOT EXISTS predictions_article_idx ON predictions (article_id)
    WHERE article_id IS NOT NULL;

-- One row per language, the same shape as content_pages: the site is trilingual
-- and a forecast readable in only one of them is a forecast two thirds of the
-- readers cannot check.
CREATE TABLE IF NOT EXISTS prediction_texts (
    prediction_id UUID NOT NULL REFERENCES predictions(id) ON DELETE CASCADE,
    lang          TEXT NOT NULL,
    statement     TEXT NOT NULL DEFAULT '',
    -- What actually happened, written when the forecast is settled.
    verdict       TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (prediction_id, lang)
);

-- +goose Down
DROP TABLE IF EXISTS prediction_texts;
DROP TABLE IF EXISTS predictions;
