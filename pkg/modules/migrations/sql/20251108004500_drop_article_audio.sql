-- +goose Up
-- Remove the narration table. The feature it backed is gone.
--
-- Reading articles aloud was built twice: first on the browser's own speech
-- synthesis, then on recordings made by a synthesiser on this server. Both
-- worked. Neither survived the question of what they cost. Synthesis in the
-- browser has no Kazakh voice anywhere and dies on a locked phone; recordings
-- fix that and put four megabytes per article per language on the disk, growing
-- every day the site publishes, for a feature nobody had asked for.
--
-- Readers who want an article read to them already have it: Safari, Edge and
-- every screen reader will speak a page on request, using voices the reader
-- chose and already has. Rebuilding that badly, at our expense, was the mistake.
--
-- The files are already gone. This drops the rows that pointed at them.
--
-- Deliberately kept: the mode column on reading_depth and the listening
-- milestones recorded through it. Nothing new will be written there, but what
-- was measured while the feature existed is still true, and dropping a column
-- to tidy up is how measurements get lost.
DROP TABLE IF EXISTS article_audio;

-- +goose Down
-- Nothing to restore: the audio these rows addressed was deleted with them.
SELECT 1;

-- +goose StatementBegin
-- The read/listen split goes with it.
--
-- The column existed to answer "how far, and by which route", because a reader
-- who scrolls away at 50% has stopped while a listener at 50% may be driving
-- and will finish. With nothing left to listen with, the second half of that
-- question has no answer to give, and the column would sit collecting a single
-- value for ever.
--
-- The rows are kept and folded together: the figures were about how far people
-- got, and that is still true however they got there.
DELETE FROM reading_depth a
      USING reading_depth b
      WHERE a.article_id = b.article_id AND a.depth = b.depth
        AND a.mode = 'listen' AND b.mode = 'read';
UPDATE reading_depth SET mode = 'read' WHERE mode <> 'read';
ALTER TABLE reading_depth DROP CONSTRAINT IF EXISTS reading_depth_pkey;
ALTER TABLE reading_depth DROP COLUMN IF EXISTS mode;
ALTER TABLE reading_depth ADD PRIMARY KEY (article_id, depth);
-- +goose StatementEnd
