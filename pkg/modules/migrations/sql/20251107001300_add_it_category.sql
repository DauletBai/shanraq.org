-- +goose Up
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_category_chk;
ALTER TABLE articles ADD CONSTRAINT articles_category_chk
    CHECK (category IN ('general', 'sport', 'society', 'politics', 'economy', 'culture', 'technology', 'it', 'opinion', 'world'));

-- +goose Down
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_category_chk;
ALTER TABLE articles ADD CONSTRAINT articles_category_chk
    CHECK (category IN ('general', 'sport', 'society', 'politics', 'economy', 'culture', 'technology', 'opinion', 'world'));
