-- +goose Up
CREATE TABLE IF NOT EXISTS framework_about (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    headline TEXT NOT NULL,
    subheadline TEXT NOT NULL,
    feature_one TEXT NOT NULL,
    feature_two TEXT NOT NULL,
    feature_three TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS framework_about;
