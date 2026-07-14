-- +goose Up
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_category_chk;
ALTER TABLE articles ADD CONSTRAINT articles_category_chk
    CHECK (category IN ('general', 'sport', 'society', 'politics', 'economy', 'culture', 'technology', 'opinion', 'world'));

ALTER TABLE articles ADD COLUMN IF NOT EXISTS subcategory TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_articles_subcategory ON articles (subcategory, status, published_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_articles_subcategory;
ALTER TABLE articles DROP COLUMN IF EXISTS subcategory;
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_category_chk;
ALTER TABLE articles ADD CONSTRAINT articles_category_chk
    CHECK (category IN ('general', 'society', 'politics', 'economy', 'culture', 'technology', 'opinion', 'world'));
