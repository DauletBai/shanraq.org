-- +goose Up
-- The three things a buyer asks before agreeing to a viewing: how old the
-- building is, what it is built of, and how high the ceilings are. All three
-- are optional — a plot of land has none of them — so zero and the empty
-- string mean "not stated" rather than "unknown to the seller".
ALTER TABLE listings ADD COLUMN IF NOT EXISTS build_year     integer NOT NULL DEFAULT 0;
ALTER TABLE listings ADD COLUMN IF NOT EXISTS wall_material  text    NOT NULL DEFAULT '';
ALTER TABLE listings ADD COLUMN IF NOT EXISTS ceiling_height numeric NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS build_year;
ALTER TABLE listings DROP COLUMN IF EXISTS wall_material;
ALTER TABLE listings DROP COLUMN IF EXISTS ceiling_height;
