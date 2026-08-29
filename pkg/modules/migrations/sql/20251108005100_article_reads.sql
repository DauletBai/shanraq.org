-- +goose Up
-- Reads, as distinct from scrolls.
--
-- reading_depth records how far down the text a reader got, which is a good
-- signal and an incomplete one: flicking to the bottom of an article in two
-- seconds reports the same 100% as reading it. The milestone says the end of
-- the prose was on screen; it cannot say anyone was there.
--
-- So the other half is time. Every article already carries an estimate -- its
-- word count at 180 a minute, the figure printed under the title -- and a
-- reader who reached the end and stayed for a serious fraction of that estimate
-- did something a scraper and a flick cannot imitate together.
--
-- Time here is ENGAGED time, not time since the page opened: it stops while the
-- tab is in the background and while nothing at all is happening, so a page
-- left open in a forgotten tab overnight counts as the minutes it was read for.
CREATE TABLE IF NOT EXISTS article_reads (
    article_id uuid PRIMARY KEY,
    -- Sessions that reached the end of the prose and spent at least half the
    -- estimated reading time engaged with it.
    finished   bigint NOT NULL DEFAULT 0,
    -- Sessions that reported at all, and the engaged seconds they add up to, so
    -- an average can be shown beside the estimate the article claims.
    samples    bigint NOT NULL DEFAULT 0,
    seconds    bigint NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS article_reads;
