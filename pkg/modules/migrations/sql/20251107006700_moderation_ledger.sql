-- +goose Up
-- Moderation had no memory: a comment could be hidden, but nothing recorded who
-- hid it, why, or on what grounds — and the author was never told. That is the
-- gap the production review called a launch blocker.
--
-- Every decision now leaves a row here, whether a person or an agent made it,
-- and the author can see the decision that concerns their own content and
-- contest it once.
CREATE TABLE IF NOT EXISTS moderation_actions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type TEXT NOT NULL CHECK (target_type IN ('article', 'listing', 'comment')),
    target_id   UUID NOT NULL,
    -- Whose content it was. Kept even if the target is later deleted, so the
    -- author's own history stays readable.
    subject_id  UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    title       TEXT NOT NULL DEFAULT '',
    action      TEXT NOT NULL CHECK (action IN ('hide', 'restore', 'reject', 'approve', 'warn')),
    reason_code TEXT NOT NULL,
    reason_note TEXT NOT NULL DEFAULT '',
    -- An agent decision and a human decision are not the same thing and must
    -- never be presented as if they were.
    actor_kind  TEXT NOT NULL CHECK (actor_kind IN ('agent', 'human')),
    actor_id    UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    actor_name  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_moderation_subject ON moderation_actions (subject_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_moderation_recent  ON moderation_actions (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_moderation_target  ON moderation_actions (target_type, target_id);

-- One appeal per decision. Unlimited appeals would turn the queue into a
-- shouting match; none at all would make an agent's mistake final.
CREATE TABLE IF NOT EXISTS moderation_appeals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action_id   UUID NOT NULL UNIQUE REFERENCES moderation_actions(id) ON DELETE CASCADE,
    author_id   UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    body        TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'upheld', 'overturned')),
    resolution  TEXT NOT NULL DEFAULT '',
    resolved_by UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- The queue a moderator actually works: open appeals, oldest first.
CREATE INDEX IF NOT EXISTS idx_appeals_open ON moderation_appeals (created_at)
    WHERE status = 'open';

-- +goose Down
DROP TABLE IF EXISTS moderation_appeals;
DROP TABLE IF EXISTS moderation_actions;
