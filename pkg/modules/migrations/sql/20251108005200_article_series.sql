-- +goose Up
-- Article series: an ordered course made out of ordinary articles.
--
-- A course is not a new kind of content. Every lesson is a normal article with a
-- normal URL, so it keeps its translations, cover, comments, votes and its place
-- in the sitemap -- which is the whole point, because the lessons have to be
-- indexable on their own to bring anyone in. What a course adds is order: which
-- lesson comes first, what follows this one, and a hub page linking all of them.
--
-- Hence a join table rather than a column on articles: an article may later
-- belong to more than one series (a lesson on templates fits both a Go course
-- and a web-basics course), and a column would have to be widened into this
-- table the first time that happens.
CREATE TABLE IF NOT EXISTS article_series (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug       text        NOT NULL UNIQUE,
    cover_url  text        NOT NULL DEFAULT '',
    status     text        NOT NULL DEFAULT 'draft',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Title and blurb per language, same shape as article translations. A series
-- with no row for the reader's language falls back in Go, not in SQL.
CREATE TABLE IF NOT EXISTS article_series_i18n (
    series_id uuid NOT NULL REFERENCES article_series (id) ON DELETE CASCADE,
    lang      text NOT NULL,
    title     text NOT NULL DEFAULT '',
    summary   text NOT NULL DEFAULT '',
    PRIMARY KEY (series_id, lang)
);

-- position is deliberately NOT unique. Reordering a forty-lesson course by
-- swapping two positions would need a temporary value to dodge a unique
-- violation on every move; ties are broken by created_at instead, which is
-- stable and costs nothing.
CREATE TABLE IF NOT EXISTS article_series_items (
    series_id  uuid NOT NULL REFERENCES article_series (id) ON DELETE CASCADE,
    article_id uuid NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    position   int  NOT NULL DEFAULT 0,
    PRIMARY KEY (series_id, article_id)
);
CREATE INDEX IF NOT EXISTS idx_series_items_order ON article_series_items (series_id, position);
CREATE INDEX IF NOT EXISTS idx_series_items_article ON article_series_items (article_id);

-- +goose Down
DROP TABLE IF EXISTS article_series_items;
DROP TABLE IF EXISTS article_series_i18n;
DROP TABLE IF EXISTS article_series;
