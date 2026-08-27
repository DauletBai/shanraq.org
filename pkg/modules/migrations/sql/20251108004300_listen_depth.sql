-- +goose Up
-- Split the depth funnel by how the article was taken in: read or listened to.
--
-- The funnel until now answered "how far did they get". With the read-aloud
-- button it has to answer "how far, and by which route", because the two
-- behave nothing alike: a reader who scrolls away at 50% has stopped, while a
-- listener at 50% may well be driving and will finish. Averaging them together
-- would hide both.
--
-- Every existing row is reading -- listening did not exist when they were
-- written -- so the column defaults to 'read' and the back-fill is the default
-- itself. The primary key grows to include the mode, which is what lets one
-- article hold both funnels side by side.
ALTER TABLE reading_depth
    ADD COLUMN IF NOT EXISTS mode text NOT NULL DEFAULT 'read'
        CHECK (mode IN ('read', 'listen'));

ALTER TABLE reading_depth DROP CONSTRAINT IF EXISTS reading_depth_pkey;
ALTER TABLE reading_depth ADD PRIMARY KEY (article_id, mode, depth);

-- +goose Down
ALTER TABLE reading_depth DROP CONSTRAINT IF EXISTS reading_depth_pkey;
DELETE FROM reading_depth WHERE mode <> 'read';
ALTER TABLE reading_depth DROP COLUMN IF EXISTS mode;
ALTER TABLE reading_depth ADD PRIMARY KEY (article_id, depth);
