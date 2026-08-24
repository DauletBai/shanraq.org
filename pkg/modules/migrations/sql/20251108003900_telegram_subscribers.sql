-- +goose Up
-- Telegram subscribers, each with the place they asked to hear about.
--
-- The channel and the bot answer two different questions, and the split mirrors
-- the one the site already makes. A channel broadcasts one identical message to
-- everybody and cannot do otherwise — Telegram gives a channel no knowledge of
-- who is reading it — so the channel carries exactly what a guest sees on the
-- site: material written for everyone. Anything addressed to a place goes
-- through the bot, which sends individual messages and can therefore apply the
-- same rule placeClause applies in the feed.
--
-- Without this, a power cut in one village reached every subscriber in the
-- country, and the answer to that is not a quieter channel but a channel that
-- was never the right carrier for local news.
CREATE TABLE IF NOT EXISTS telegram_subscribers (
    tg_user_id  BIGINT PRIMARY KEY,
    -- The place this person asked about. NULL means they started the bot but
    -- have not finished choosing; they receive nothing until they do, because
    -- guessing someone's town from their phone is exactly the tracking this
    -- platform promises not to do.
    geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL,
    lang        TEXT NOT NULL DEFAULT 'ru',
    -- Cleared when Telegram reports the person blocked the bot. The row stays
    -- so that coming back does not look like a new subscriber, and so the count
    -- of people who left is knowable.
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- The one question asked on every publish: who is inside this place.
CREATE INDEX IF NOT EXISTS idx_tg_subs_place
    ON telegram_subscribers (geo_node_id) WHERE active AND geo_node_id IS NOT NULL;

-- Where the long poll left off. A single row: restarting must not replay the
-- updates already handled, which would answer every /start twice.
CREATE TABLE IF NOT EXISTS telegram_bot_state (
    one       BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (one),
    update_id BIGINT NOT NULL DEFAULT 0
);
INSERT INTO telegram_bot_state (one, update_id) VALUES (TRUE, 0) ON CONFLICT DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS telegram_bot_state;
DROP INDEX IF EXISTS idx_tg_subs_place;
DROP TABLE IF EXISTS telegram_subscribers;
