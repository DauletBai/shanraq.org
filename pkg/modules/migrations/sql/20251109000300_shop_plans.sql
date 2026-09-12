-- +goose Up
-- A product is sold in packages: the book alone, or the book with the code the
-- book builds. The two differ in price and in what the buyer downloads, and
-- everything else about them -- the page, the description, the cover -- is the
-- same product. So the package is its own row, and the price lives on it.
--
-- The price column on shop_products goes away in the same change: two places
-- that can both answer "how much" is the kind of pair that drifts apart, and
-- the answer belongs to the package a buyer actually chooses.
CREATE TABLE IF NOT EXISTS shop_plans (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  uuid NOT NULL REFERENCES shop_products(id) ON DELETE CASCADE,
    -- A short stable name used in code and in the payment ledger: 'book', 'kit'.
    code        text NOT NULL,
    -- Whole tenge. 0 means the price has not been announced yet, and the page
    -- says exactly that instead of offering something for nothing.
    price       bigint NOT NULL DEFAULT 0 CHECK (price >= 0),
    position    int    NOT NULL DEFAULT 0,
    title_kz    text NOT NULL DEFAULT '',
    title_ru    text NOT NULL DEFAULT '',
    title_en    text NOT NULL DEFAULT '',
    -- What the package includes, one line per item (Markdown list).
    includes_kz text NOT NULL DEFAULT '',
    includes_ru text NOT NULL DEFAULT '',
    includes_en text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, code)
);

CREATE INDEX IF NOT EXISTS shop_plans_product_idx ON shop_plans (product_id, position);

INSERT INTO shop_plans (product_id, code, position, title_kz, title_ru, title_en,
                        includes_kz, includes_ru, includes_en)
SELECT p.id, 'book', 10, 'Кітап', 'Книга', 'The book',
'- Кітап толығымен: PDF және EPUB
- Барлық болашақ редакция тегін
- DRM жоқ: файл сізде қалады',
'- Книга целиком: PDF и EPUB
- Все будущие редакции бесплатно
- Без DRM: файл остаётся у вас',
'- The whole book: PDF and EPUB
- Every future edition free
- No DRM: the file is yours to keep'
FROM shop_products p WHERE p.slug = 'go-book'
ON CONFLICT (product_id, code) DO NOTHING;

INSERT INTO shop_plans (product_id, code, position, title_kz, title_ru, title_en,
                        includes_kz, includes_ru, includes_en)
SELECT p.id, 'kit', 20, 'Кітап пен код', 'Книга и код', 'Book and code',
'- «Кітаптағының» бәрі, оған қоса браузерде оқуға арналған HTML
- Кітап құратын дүкеннің толық бастапқы коды
- Кодтың да, кітаптың да барлық болашақ редакциясы',
'- Всё из «Книги», плюс HTML для чтения в браузере
- Полный исходный код магазина, который строит книга
- Все будущие редакции и книги, и кода',
'- Everything in "The book", plus HTML for reading in a browser
- The complete source code of the shop the book builds
- Every future edition of both the book and the code'
FROM shop_products p WHERE p.slug = 'go-book'
ON CONFLICT (product_id, code) DO NOTHING;

ALTER TABLE shop_products DROP COLUMN IF EXISTS price;

-- +goose Down
ALTER TABLE shop_products ADD COLUMN IF NOT EXISTS price bigint NOT NULL DEFAULT 0;
DROP TABLE IF EXISTS shop_plans;
