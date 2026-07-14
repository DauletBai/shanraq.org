-- +goose Up
ALTER TABLE listings ADD COLUMN IF NOT EXISTS geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS geo_node_id;
