-- +goose Up
-- Favorites / bookmarks: a logged-in user can save articles and listings.
-- One generic table keyed by (user, item_type, item_id) so both content types
-- share the same mechanism. No FK to articles/listings (two different tables);
-- orphaned rows are harmless and filtered out by the status join on read.
CREATE TABLE IF NOT EXISTS favorites (
    user_id    uuid        NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    item_type  text        NOT NULL CHECK (item_type IN ('article', 'listing')),
    item_id    uuid        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, item_type, item_id)
);

CREATE INDEX IF NOT EXISTS idx_favorites_user
    ON favorites (user_id, item_type, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS favorites;
