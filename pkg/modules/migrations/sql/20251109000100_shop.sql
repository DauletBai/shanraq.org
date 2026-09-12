-- +goose Up
-- The shop: digital goods the site sells itself.
--
-- One table for the product and one for the people waiting for it. A product
-- is trilingual like everything else here, and carries its own page copy: the
-- book's page is not an article, and squeezing it into the article tables
-- would give it a byline, a category and a comment thread it has no use for.
--
-- Price and payment are deliberately absent from this migration beyond a
-- column: nothing can be bought until an acquirer is connected, and a product
-- with no price is a page that says what is coming, which is what the first
-- one is for.
CREATE TABLE IF NOT EXISTS shop_products (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        text NOT NULL UNIQUE,
    kind        text NOT NULL DEFAULT 'book',
    -- draft: visible to staff only. announced: public page, no purchase.
    -- selling: purchase available. paused: page stays, purchase stops.
    status      text NOT NULL DEFAULT 'draft'
                CHECK (status IN ('draft', 'announced', 'selling', 'paused')),
    -- Whole tenge, as everywhere else on this site. 0 while unannounced.
    price       bigint NOT NULL DEFAULT 0 CHECK (price >= 0),
    currency    text NOT NULL DEFAULT 'KZT',
    -- What the buyer would receive today: the edition, and how far it has got.
    edition     text NOT NULL DEFAULT '',
    cover_url   text NOT NULL DEFAULT '',
    -- The free fragment. Empty hides the offer rather than showing a dead link.
    preview_url text NOT NULL DEFAULT '',
    title_kz    text NOT NULL DEFAULT '',
    title_ru    text NOT NULL DEFAULT '',
    title_en    text NOT NULL DEFAULT '',
    summary_kz  text NOT NULL DEFAULT '',
    summary_ru  text NOT NULL DEFAULT '',
    summary_en  text NOT NULL DEFAULT '',
    body_kz     text NOT NULL DEFAULT '',
    body_ru     text NOT NULL DEFAULT '',
    body_en     text NOT NULL DEFAULT '',
    position    int  NOT NULL DEFAULT 0,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW()
);

-- People who asked to be told when a product is ready. This is a queue of
-- consents, not a mailing list: one row per product per address, and the
-- address is kept only for that announcement.
CREATE TABLE IF NOT EXISTS shop_leads (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  uuid NOT NULL REFERENCES shop_products(id) ON DELETE CASCADE,
    email       text NOT NULL,
    lang        text NOT NULL DEFAULT 'ru',
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    notified_at timestamptz,
    UNIQUE (product_id, email)
);

CREATE INDEX IF NOT EXISTS shop_products_public_idx
    ON shop_products (position, created_at)
    WHERE status <> 'draft';

-- +goose Down
DROP TABLE IF EXISTS shop_leads;
DROP TABLE IF EXISTS shop_products;
