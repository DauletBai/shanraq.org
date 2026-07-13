-- +goose Up
-- The "story" — one logical article, language-independent metadata.
CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    slug TEXT UNIQUE NOT NULL,
    original_lang TEXT NOT NULL DEFAULT 'ru',
    status TEXT NOT NULL DEFAULT 'draft',
    cover_url TEXT NOT NULL DEFAULT '',
    score INTEGER NOT NULL DEFAULT 0,
    views_count BIGINT NOT NULL DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT articles_status_chk CHECK (status IN ('draft', 'published', 'archived')),
    CONSTRAINT articles_orig_lang_chk CHECK (original_lang IN ('kz', 'ru', 'en'))
);

CREATE INDEX IF NOT EXISTS idx_articles_author ON articles (author_id);
CREATE INDEX IF NOT EXISTS idx_articles_feed ON articles (status, published_at DESC);

-- One row per language version (kz | ru | en). Written by a human or the AI.
CREATE TABLE IF NOT EXISTS article_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    lang TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    body_md TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL DEFAULT 'human',
    status TEXT NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT article_tr_lang_chk CHECK (lang IN ('kz', 'ru', 'en')),
    CONSTRAINT article_tr_source_chk CHECK (source IN ('human', 'ai')),
    CONSTRAINT article_tr_status_chk CHECK (status IN ('draft', 'pending', 'ready')),
    UNIQUE (article_id, lang)
);

CREATE INDEX IF NOT EXISTS idx_article_tr_article ON article_translations (article_id);

-- Lightweight per-day view analytics for the author cabinet dashboard.
CREATE TABLE IF NOT EXISTS article_views_daily (
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    lang TEXT NOT NULL,
    day DATE NOT NULL,
    views BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (article_id, lang, day)
);

-- +goose Down
DROP TABLE IF EXISTS article_views_daily;
DROP TABLE IF EXISTS article_translations;
DROP TABLE IF EXISTS articles;
