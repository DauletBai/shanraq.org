-- +goose Up
-- Editorial covers live in the database so a database restore preserves them.
CREATE TABLE IF NOT EXISTS article_cover_blobs (
    article_id uuid PRIMARY KEY REFERENCES articles(id) ON DELETE CASCADE,
    content_type text NOT NULL CHECK (content_type IN ('image/jpeg', 'image/png', 'image/webp')),
    body bytea NOT NULL CHECK (octet_length(body) BETWEEN 1 AND 10000000)
);

-- +goose Down
DROP TABLE IF EXISTS article_cover_blobs;
