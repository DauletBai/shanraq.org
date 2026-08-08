-- +goose Up
-- Withdraw the exception made in 20251108000100.
--
-- That migration put 8 AI columns back into the index, and the argument rested
-- on the site's single organic click having landed on one of them. The owner
-- then said the click was almost certainly his own, made while checking whether
-- the site turns up in search at all.
--
-- What is left without it: 16 impressions across those 8 pages in three months,
-- at average position 13, and zero clicks from anyone. Google offers the pages
-- and nobody takes them. That is not a case for keeping 8 exceptions and it is
-- not worth muddying the next thirty days of measurement, when the question is
-- whether a clean index of 14 human articles performs better than 104 pages of
-- which 90 are unread.
--
-- Kept as a second migration rather than deleting the first: the first never
-- reached production, but the reasoning is worth keeping in the record, and so
-- is the correction.
UPDATE articles SET indexable = FALSE
 WHERE author_id = '5a2a0000-0000-0000-0000-000000000001';

-- +goose Down
-- Nothing to undo: 20251107009900 already sets exactly this state.
SELECT 1;
