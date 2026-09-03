-- +goose Up
-- One attempt back when the reviewer was wrong.
--
-- A lesson gives three checks, and the reviewer is a language model: it can
-- refuse a solution that does the job. Losing a third of the budget to that is
-- unfair enough that a reader stops trusting the mark -- and an unfair mark
-- teaches nothing at all.
--
-- Once per exercise, no argument required: the reader says the check was wrong
-- and the attempt comes back. Once, so it cannot be farmed, and recorded, so we
-- can count how often the reviewer is actually wrong instead of guessing.
ALTER TABLE course_progress
    ADD COLUMN IF NOT EXISTS appealed boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE course_progress DROP COLUMN IF EXISTS appealed;
