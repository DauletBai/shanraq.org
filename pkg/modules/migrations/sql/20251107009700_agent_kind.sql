-- +goose Up
-- Who the professional actually is. A private realtor, an agency and a
-- developer all need the same things from us — a verified badge, a public page,
-- listings attributed to them — so they stay one entity rather than three
-- parallel tables with three moderation queues. Only the label, the required
-- fields and the badge wording differ.
--
-- 'private' is the default so every profile created before this migration keeps
-- working unchanged.
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS kind TEXT NOT NULL DEFAULT 'private';
-- БИН/ИИН: identifies a legal entity, so a company claiming a brand can be
-- checked against the public register before it gets the badge.
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS bin TEXT NOT NULL DEFAULT '';

ALTER TABLE re_agents DROP CONSTRAINT IF EXISTS re_agents_kind_check;
ALTER TABLE re_agents ADD CONSTRAINT re_agents_kind_check
	CHECK (kind IN ('private', 'agency', 'developer'));

-- +goose Down
ALTER TABLE re_agents DROP CONSTRAINT IF EXISTS re_agents_kind_check;
ALTER TABLE re_agents DROP COLUMN IF EXISTS bin;
ALTER TABLE re_agents DROP COLUMN IF EXISTS kind;
