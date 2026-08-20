-- +goose Up
-- Readers judge comments, the way they already judge articles.
--
-- What stood here before was a model: every comment went to kimi-k2.6, which
-- decided whether to hide it. It never hid anything — there were twenty-one
-- comments and it passed all of them — but the same model had by then invented
-- a number in a published article and dropped a sentence from another, and the
-- appetite for letting it decide what readers may see ran out.
--
-- A vote is a better instrument for this than a classifier. It needs no prompt,
-- it cannot hallucinate, it costs nothing per comment, and when it is wrong the
-- people who were wrong are the readers themselves — which is the only kind of
-- moderation this platform can honestly promise, since it already promises that
-- readers moderate articles.
--
-- Nothing is deleted and nothing is hidden outright: a comment voted far enough
-- down is folded away behind one line, and one click still opens it. Being
-- unpopular is not the same as being unpublishable.

-- One vote per reader per comment. The voter's weight is snapshotted at vote
-- time and grows with their own karma, so a fresh account cannot bury a comment
-- and a crowd of them cannot either.
CREATE TABLE IF NOT EXISTS comment_votes (
    comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    value SMALLINT NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (comment_id, user_id),
    CONSTRAINT comment_votes_value_chk CHECK (value IN (-1, 1))
);

CREATE INDEX IF NOT EXISTS idx_comment_votes_comment ON comment_votes (comment_id);

-- The running total, cached on the comment so a page of them costs one query.
ALTER TABLE comments ADD COLUMN IF NOT EXISTS score INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE comments DROP COLUMN IF EXISTS score;
DROP TABLE IF EXISTS comment_votes;
