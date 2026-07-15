-- +goose Up
-- The "internet" and "ai" subcategories moved from Technology to the new IT
-- rubric; realign the two demo articles that used them (science stays in
-- Technology, so demo-sub-science is left untouched).
UPDATE articles SET category = 'it'
WHERE slug IN ('demo-sub-internet', 'demo-sub-ai') AND category = 'technology';

-- +goose Down
UPDATE articles SET category = 'technology'
WHERE slug IN ('demo-sub-internet', 'demo-sub-ai') AND category = 'it';
