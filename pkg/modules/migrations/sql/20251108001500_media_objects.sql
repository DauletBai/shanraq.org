-- +goose Up
-- Uploads were anonymous the moment they landed. The file went to disk under a
-- content hash, the URL went into a listing or an article, and nothing recorded
-- who had put it there or how much of the disk they were now using. So there
-- was no quota to enforce, no way to answer "whose file is this", and no way to
-- find the uploads that were never referenced by anything — the ones from a
-- listing form somebody abandoned halfway. On a single VPS that is a disk that
-- only fills.
--
-- media_objects is the file: one row per stored key, with its size. Keys are
-- content hashes, so two people uploading the same photo share one object and
-- one copy on disk — which is why ownership is a separate table rather than a
-- column here.
CREATE TABLE IF NOT EXISTS media_objects (
    key          text PRIMARY KEY,
    bytes        bigint NOT NULL CHECK (bytes >= 0),
    content_type text NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS media_objects_created_idx ON media_objects (created_at);

-- media_owners is who uploaded it. Many-to-many for the shared-hash case above,
-- and it cascades from both sides: deleting the account releases its claim on
-- the file, which is what lets the sweep collect files nobody is left to own.
CREATE TABLE IF NOT EXISTS media_owners (
    key        text NOT NULL REFERENCES media_objects (key) ON DELETE CASCADE,
    user_id    uuid NOT NULL REFERENCES auth_users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (key, user_id)
);

CREATE INDEX IF NOT EXISTS media_owners_user_idx ON media_owners (user_id);

-- +goose Down
DROP TABLE IF EXISTS media_owners;
DROP TABLE IF EXISTS media_objects;
