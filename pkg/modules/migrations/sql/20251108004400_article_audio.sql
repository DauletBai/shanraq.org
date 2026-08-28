-- +goose Up
-- Pre-rendered narration for an article, one row per (article, language).
--
-- Kazakh is the reason this exists. The read-aloud button leans on the
-- browser's own speech synthesis, and browsers ship no Kazakh voice: the
-- player detects that, offers a Russian or English substitute, and a Kazakh
-- reader ends up hearing Kazakh read by a Russian voice. That is worse than
-- silence, because it sounds like the site does not know the difference.
--
-- The audio is synthesised outside this server -- there is no text-to-speech
-- engine here and no reason to put one on a small VPS -- and uploaded. So the
-- row records what arrived rather than what to produce: the storage key, the
-- public URL, how long it plays, how big it is, and which voice made it. The
-- voice is kept because the audio outlives the tool that made it, and in two
-- years the only way to answer "what read this" will be this column.
--
-- The audio deliberately does not go through the media ledger. That ledger
-- exists to meter what people upload and to sweep files nothing refers to, and
-- it decides "refers to" by searching article bodies for the key. Narration is
-- referenced by this table and never by the prose, so a ledger entry would be
-- swept as an orphan on the next pass. Cleanup is instead tied to the article:
-- delete the article, the row goes with it.
CREATE TABLE IF NOT EXISTS article_audio (
    article_id   UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    lang         TEXT NOT NULL,
    storage_key  TEXT NOT NULL,
    url          TEXT NOT NULL,
    duration_sec INTEGER NOT NULL DEFAULT 0,
    bytes        BIGINT NOT NULL DEFAULT 0,
    voice        TEXT NOT NULL DEFAULT '',
    -- The hash of the text the audio was made from. An article edited after
    -- narration leaves the two out of step, and without this there is no way to
    -- notice: the audio still plays, just saying something the page no longer
    -- says. With it the page can mark the narration stale and the generator can
    -- skip everything that has not changed.
    text_sha256  TEXT NOT NULL DEFAULT '',
    -- Where each block of the article starts and ends in the recording, as
    -- [{"i":0,"a":0.0,"b":4.2}, ...] against the same block sequence the page
    -- reads aloud.
    --
    -- Browser speech synthesis gave this away for free: it fires a boundary
    -- event as it crosses each word, which is what the old highlight followed.
    -- A finished audio file says nothing about itself, so the timings have to
    -- be recorded while it is made -- one block synthesised at a time, its
    -- length noted -- or the listener gets sound with no idea where on the page
    -- it is coming from. That is the whole point of the feature: you can be
    -- doing something else and still see where the reading has got to.
    cues         JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (article_id, lang),
    CONSTRAINT article_audio_lang_chk CHECK (lang IN ('kz', 'ru', 'en'))
);

-- +goose Down
DROP TABLE IF EXISTS article_audio;
