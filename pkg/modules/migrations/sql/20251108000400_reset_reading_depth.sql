-- +goose Up
-- Finish the reset that 20251108000300 started.
--
-- That migration zeroed the view counters and left reading_depth alone, on the
-- reasoning that the depth data was clean and worth keeping. It was clean. It
-- was also displayed in the same table row, so the studio ended up showing an
-- article with 0 views and 11 readers who finished it — an impossible state,
-- and a worse lie than the one being fixed.
--
-- reading_depth has no timestamp (article_id, depth, count), so the rows cannot
-- be split into before and after the reset. Either both counters start together
-- or neither does. They start together.
--
-- The numbers are preserved, not destroyed: 245 readers reached a quarter, 202
-- half, 194 three quarters, 181 the end. That is the most trustworthy data the
-- site has produced — a browser and a scroll are needed to record it, so no
-- crawler could ever inflate it — and it stays available in the legacy table.
CREATE TABLE IF NOT EXISTS reading_depth_legacy (LIKE reading_depth INCLUDING ALL);
INSERT INTO reading_depth_legacy SELECT * FROM reading_depth
    ON CONFLICT DO NOTHING;
DELETE FROM reading_depth;

-- +goose Down
INSERT INTO reading_depth SELECT * FROM reading_depth_legacy
    ON CONFLICT (article_id, depth) DO UPDATE SET count = EXCLUDED.count;
DROP TABLE IF EXISTS reading_depth_legacy;
