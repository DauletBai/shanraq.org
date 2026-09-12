-- +goose Up
-- The prices the owner settled on for the Go book, after comparing the nearest
-- analogue in the world (a self-published Go book selling the same two
-- packages at $34.95 and $59.95) with what a digital book costs here.
--
-- Only a package that has no price yet is touched: once the panel has set a
-- price, this migration must never reach back and overwrite it.
UPDATE shop_plans SET price = 9900, updated_at = NOW()
 WHERE code = 'book' AND price = 0
   AND product_id = (SELECT id FROM shop_products WHERE slug = 'go-book');

UPDATE shop_plans SET price = 19900, updated_at = NOW()
 WHERE code = 'kit' AND price = 0
   AND product_id = (SELECT id FROM shop_products WHERE slug = 'go-book');

-- +goose Down
UPDATE shop_plans SET price = 0
 WHERE code IN ('book', 'kit')
   AND product_id = (SELECT id FROM shop_products WHERE slug = 'go-book');
