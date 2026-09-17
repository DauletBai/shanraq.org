-- +goose Up
-- Publish the checked free sample and its cover, without changing prices or sale status.
UPDATE shop_products
SET cover_url = '/static/shop/go-book-sample/0.21.0-sample/read/cover/go-book-cover-v3.png',
    preview_url = '/static/shop/go-book-sample/0.21.0-sample/index.html',
    updated_at = now()
WHERE slug = 'go-book';

-- +goose Down
-- Restore embedded legacy URLs only if this revision still owns the values.
UPDATE shop_products
SET cover_url = '/static/shop/go-book-cover.jpg', updated_at = now()
WHERE slug = 'go-book'
  AND cover_url = '/static/shop/go-book-sample/0.21.0-sample/read/cover/go-book-cover-v3.png';
UPDATE shop_products
SET preview_url = '/static/shop/go-book-preview.pdf', updated_at = now()
WHERE slug = 'go-book'
  AND preview_url = '/static/shop/go-book-sample/0.21.0-sample/index.html';
