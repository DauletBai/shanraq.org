-- The cover is a file name, not a path: where uploads live is a setting, and a
-- database that remembers a path would go stale the day the folder moves.
alter table articles add column cover text not null default '';
