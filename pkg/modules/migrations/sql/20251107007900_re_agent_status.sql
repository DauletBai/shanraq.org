-- +goose Up
-- Agent verification. A real-estate agent profile starts as 'pending' and only
-- becomes a public, badged agent after an admin approves it. Any rows created
-- before this migration (self-serve, unvetted) default to 'pending', so no
-- unverified profile keeps a trust badge.
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'verified', 'rejected'));
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS reject_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ;
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS reviewed_by UUID;

CREATE INDEX IF NOT EXISTS idx_re_agents_status ON re_agents (status);

-- +goose Down
DROP INDEX IF EXISTS idx_re_agents_status;
ALTER TABLE re_agents DROP COLUMN IF EXISTS status;
ALTER TABLE re_agents DROP COLUMN IF EXISTS reject_reason;
ALTER TABLE re_agents DROP COLUMN IF EXISTS reviewed_at;
ALTER TABLE re_agents DROP COLUMN IF EXISTS reviewed_by;
