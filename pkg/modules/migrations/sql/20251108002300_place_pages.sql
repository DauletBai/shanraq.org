-- +goose Up
-- Place pages: every place gets an address, every article gets a place.
--
-- The idea: the chair of a housing service publishes for their settlement, a mayor's
-- office for its region, and that acquires a page a search engine can find. The same
-- geo_nodes reference is already used by the listings, so a place on the page brings
-- one area's articles and listings into a single feed.
--
-- The address does not come from code: that column is not uniform — Kachar has
-- kz-kostanay-kachar while the Kostanay region has g65, and /place/g65 is a cipher,
-- not an address. The slug is filled from the name by transliteration at application
-- start-up: it cannot be done in SQL, because "ж" becomes "zh" while translate()
-- can only substitute one character for one.
ALTER TABLE geo_nodes ADD COLUMN IF NOT EXISTS slug TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_geo_nodes_slug ON geo_nodes (slug) WHERE slug IS NOT NULL;

-- An article's place. NULL means "for everyone", and it stays the default: an
-- article with no place behaves exactly as every article has behaved until now.
ALTER TABLE articles ADD COLUMN IF NOT EXISTS geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_articles_geo ON articles (geo_node_id) WHERE geo_node_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_articles_geo;
ALTER TABLE articles DROP COLUMN IF EXISTS geo_node_id;
DROP INDEX IF EXISTS idx_geo_nodes_slug;
ALTER TABLE geo_nodes DROP COLUMN IF EXISTS slug;
