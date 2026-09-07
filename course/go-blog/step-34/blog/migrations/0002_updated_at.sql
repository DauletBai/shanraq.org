-- not null always comes with a default: existing rows have to be filled with
-- something, and on a database that already holds articles the database
-- refuses the column outright without one.
alter table articles add column updated_at text not null default '';
