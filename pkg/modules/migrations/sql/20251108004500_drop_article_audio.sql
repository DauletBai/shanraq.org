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
