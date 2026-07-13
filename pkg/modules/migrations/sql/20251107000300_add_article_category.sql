-- +goose Up
ALTER TABLE articles ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'general';
ALTER TABLE articles ADD CONSTRAINT articles_category_chk
    CHECK (category IN ('general', 'society', 'politics', 'economy', 'culture', 'technology', 'opinion', 'world'));
CREATE INDEX IF NOT EXISTS idx_articles_category ON articles (category, status, published_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_articles_category;
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_category_chk;
ALTER TABLE articles DROP COLUMN IF EXISTS category;
