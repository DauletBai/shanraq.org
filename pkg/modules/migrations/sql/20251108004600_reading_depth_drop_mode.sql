-- +goose Up
-- Fold the read/listen split back into one figure.
--
-- The column answered "how far, and by which route", because a reader who
-- scrolls away at 50% has stopped while a listener at 50% may be driving and
-- will finish. There is no second route now: reading aloud is gone, and the
-- column would sit collecting one value for ever.
--
-- The rows are kept, not dropped. They were about how far people got, and that
-- stays true however they got there. Where an article has both a read and a
-- listen row for the same milestone, the listen row is removed rather than
-- summed: a listener who reached 50% is one person who reached 50%, and adding
-- the two would count them twice.
--
-- This is separate from 20251108004500 because that migration had already run
-- by the time the column needed dropping; editing an applied migration changes
-- nothing on a database that has seen it.
DELETE FROM reading_depth a
      USING reading_depth b
      WHERE a.article_id = b.article_id AND a.depth = b.depth
        AND a.mode = 'listen' AND b.mode = 'read';
UPDATE reading_depth SET mode = 'read' WHERE mode <> 'read';
ALTER TABLE reading_depth DROP CONSTRAINT IF EXISTS reading_depth_pkey;
ALTER TABLE reading_depth DROP COLUMN IF EXISTS mode;
ALTER TABLE reading_depth ADD PRIMARY KEY (article_id, depth);

-- +goose Down
ALTER TABLE reading_depth DROP CONSTRAINT IF EXISTS reading_depth_pkey;
ALTER TABLE reading_depth
    ADD COLUMN IF NOT EXISTS mode text NOT NULL DEFAULT 'read'
        CHECK (mode IN ('read', 'listen'));
ALTER TABLE reading_depth ADD PRIMARY KEY (article_id, mode, depth);
