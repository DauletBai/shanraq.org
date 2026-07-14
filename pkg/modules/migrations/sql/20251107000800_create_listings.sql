-- +goose Up
CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    deal_type TEXT NOT NULL DEFAULT 'sale',
    property_type TEXT NOT NULL DEFAULT 'apartment',
    country TEXT NOT NULL DEFAULT '',
    region TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL DEFAULT '',
    village TEXT NOT NULL DEFAULT '',
    price BIGINT NOT NULL DEFAULT 0,
    area NUMERIC NOT NULL DEFAULT 0,
    rooms INT NOT NULL DEFAULT 0,
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    contact TEXT NOT NULL DEFAULT '',
    cover_url TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'published',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT listings_deal_chk CHECK (deal_type IN ('sale', 'rent')),
    CONSTRAINT listings_type_chk CHECK (property_type IN ('apartment', 'house', 'land', 'commercial', 'dacha')),
    CONSTRAINT listings_status_chk CHECK (status IN ('published', 'hidden'))
);

CREATE INDEX IF NOT EXISTS idx_listings_feed ON listings (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_listings_author ON listings (author_id);

-- +goose Down
DROP TABLE IF EXISTS listings;
