-- +goose Up
-- Split the agent's single name into first + last name. `name` is kept as the
-- composed display name (maintained by the app on save) so the listing badge
-- and the join keep working unchanged.
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS first_name TEXT NOT NULL DEFAULT '';
ALTER TABLE re_agents ADD COLUMN IF NOT EXISTS last_name  TEXT NOT NULL DEFAULT '';
UPDATE re_agents SET first_name = name WHERE first_name = '' AND name <> '';

-- +goose Down
ALTER TABLE re_agents DROP COLUMN IF EXISTS first_name;
ALTER TABLE re_agents DROP COLUMN IF EXISTS last_name;
