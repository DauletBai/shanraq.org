-- +goose Up
-- Per-article search-indexing switch.
--
-- Why: 90 of the site's 104 articles are AI columns, and between them they have
-- collected 2 views. Google's spam policy does not target AI writing as such —
-- it targets many pages produced at scale that add little value, and a site
-- where 87% of the pages are unread thin content is exactly that pattern. The
-- 14 human articles carry all of the traffic and all of the reputation, and
-- they are the ones that pay for it.
--
-- Deliberately a flag, not a delete: nothing is destroyed, no URL starts
-- returning 404, the pages stay readable and linkable, and a column judged good
-- later can be put back into the index by flipping one boolean.
ALTER TABLE articles ADD COLUMN IF NOT EXISTS indexable BOOLEAN NOT NULL DEFAULT TRUE;

-- Take the AI columnist's back catalogue out of the index. The author id is the
-- fixed seed id of the AI columnist (persona.go SanaAuthorID).
UPDATE articles
   SET indexable = FALSE
 WHERE author_id = '5a2a0000-0000-0000-0000-000000000001';

-- Sitemap and article rendering both filter on this, so keep the lookup cheap.
CREATE INDEX IF NOT EXISTS articles_indexable_idx ON articles (indexable) WHERE indexable;

-- +goose Down
DROP INDEX IF EXISTS articles_indexable_idx;
ALTER TABLE articles DROP COLUMN IF EXISTS indexable;
