-- +goose Up
-- Hosts, visitors, visits and views over time.
--
-- analytics_daily counts events and deliberately cannot tell one visitor from
-- another, so it can report views and nothing else. Reporting the other three
-- means telling people apart, which is what this table does -- and the shape
-- below is chosen so that telling them apart stays possible only inside one
-- day, and only in aggregate.
--
-- vid and host are truncated HMACs, not addresses. The key is a salt that is
-- regenerated every day and never stored beside the rows it hashed, so the same
-- reader on Tuesday and Wednesday produces two unrelated values: the table can
-- count a day's distinct visitors and cannot follow anyone across days. Nothing
-- here identifies a person, records a path, or survives the salt it was made
-- with.
--
-- A row is one visitor in one half-hour slot, which is how a visit is counted:
-- the same reader returning after a break lands in a new slot and is a new
-- visit, while a reader clicking through five pages stays one. That makes every
-- figure a plain query over stored rows, with no session state to keep in
-- memory and nothing to reconstruct after a restart.
--
--   Views      = SUM(views)
--   Visits     = COUNT(*)                 -- one row is one visitor-slot
--   Visitors   = COUNT(DISTINCT vid)
--   Hosts      = COUNT(DISTINCT host)
--
-- is_kz and is_mobile are the audience switches, stored as dimensions rather
-- than four pre-summed audiences so a query can combine them freely.
CREATE TABLE IF NOT EXISTS analytics_slots (
    -- Half past or on the hour, UTC. Rolls up to hours, days, weeks, months.
    slot      TIMESTAMPTZ NOT NULL,
    vid       BYTEA       NOT NULL,
    host      BYTEA       NOT NULL,
    is_kz     BOOLEAN     NOT NULL DEFAULT FALSE,
    is_mobile BOOLEAN     NOT NULL DEFAULT FALSE,
    views     BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (slot, vid)
);

CREATE INDEX IF NOT EXISTS analytics_slots_slot_idx ON analytics_slots (slot);

-- The salt each day's identifiers were made with. Kept only so a restart does
-- not split one day into two populations of unrelated hashes; the sweeper drops
-- it with the rows, and once dropped no stored value can be traced back to an
-- address even by us.
CREATE TABLE IF NOT EXISTS analytics_salt (
    day  DATE  NOT NULL PRIMARY KEY,
    salt BYTEA NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS analytics_salt;
DROP TABLE IF EXISTS analytics_slots;
