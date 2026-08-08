-- +goose Up
-- Put back the AI columns that Google has actually surfaced.
--
-- Migration 20251107009900 took all 90 out of the index on the grounds that
-- they had 2 views between them. Search Console then showed what the site's own
-- analytics could not: 8 of the 10 pages Google shows at all are those columns,
-- and the site's single organic click in three months landed on one of them
-- (povorot-na-plato-solason).
--
-- The blanket call was too broad. It stays right for the other 80-odd columns —
-- Google has never surfaced them, and 90 unread pages against 14 read ones is
-- the scaled-content pattern whatever wrote them. But there is no case for
-- hiding pages that demonstrably reach readers, and the quality signal cannot
-- meaningfully differ between 8 columns and 0.
--
-- Slugs, not a rule: this is a record of a specific decision about specific
-- pages, taken from a specific report. A rule would be a claim we cannot make.
UPDATE articles SET indexable = TRUE
 WHERE slug IN (
    'povorot-na-plato-solason',          -- 1 impression, 1 click — the only one
    'sana-migraciya-steny',              -- 3 impressions
    'sana-nikto-ne-rozhdaetsya',         -- 3
    'sana-istoriyu-pishut-arhivy',       -- 2
    'rost-kotorogo-ne-pochuvstvovali',   -- 2
    'sana-vselennaya-poznaet-sebya',     -- 2
    'sana-zvonok-kotoryi-otkladyvaete',  -- 2
    'sana-zavist-k-ptice'                -- 1
 );

-- +goose Down
UPDATE articles SET indexable = FALSE
 WHERE author_id = '5a2a0000-0000-0000-0000-000000000001';
