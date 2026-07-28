-- +goose Up
-- Editable info & legal pages (About, Guide, Pricing, Support, Formatting,
-- Privacy, Terms) — one row per (page, language). Content moves out of Go so
-- operators can update policies and tariffs without a code change or redeploy;
-- the app seeds the current built-in text on boot and never overwrites a later
-- edit (ON CONFLICT DO NOTHING on the seed path).
CREATE TABLE IF NOT EXISTS content_pages (
    page_key   TEXT NOT NULL,
    lang       TEXT NOT NULL,
    title      TEXT NOT NULL DEFAULT '',
    body_md    TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    PRIMARY KEY (page_key, lang)
);

-- +goose Down
DROP TABLE IF EXISTS content_pages;
