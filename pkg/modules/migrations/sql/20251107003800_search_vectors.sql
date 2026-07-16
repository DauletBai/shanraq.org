-- +goose Up
-- Full-text search over articles (all language translations) and listings,
-- using the 'simple' config (tokenize + lowercase, no language-specific
-- stemming) so KZ/RU/EN all match consistently.
ALTER TABLE article_translations ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple',
            coalesce(title, '') || ' ' || coalesce(summary, '') || ' ' || coalesce(body_md, ''))
    ) STORED;
CREATE INDEX IF NOT EXISTS idx_at_search ON article_translations USING GIN (search_vector);

ALTER TABLE listings ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple',
            coalesce(title, '') || ' ' || coalesce(description, '') || ' ' ||
            coalesce(region, '') || ' ' || coalesce(city, '') || ' ' || coalesce(country, ''))
    ) STORED;
CREATE INDEX IF NOT EXISTS idx_listings_search ON listings USING GIN (search_vector);

-- +goose Down
DROP INDEX IF EXISTS idx_listings_search;
ALTER TABLE listings DROP COLUMN IF EXISTS search_vector;
DROP INDEX IF EXISTS idx_at_search;
ALTER TABLE article_translations DROP COLUMN IF EXISTS search_vector;
