-- +goose Up
-- Listings can now be posted for property in Russia as well as Kazakhstan, so a
-- listing carries its own currency: tenge (KZT) for Kazakh addresses, ruble
-- (RUB) for Russian ones. Chosen from the selected location's country on submit.
ALTER TABLE listings ADD COLUMN IF NOT EXISTS currency CHAR(3) NOT NULL DEFAULT 'KZT';

-- Backfill existing rows: any listing whose location node sits in Russia becomes
-- RUB; the free-text fallback covers rows without a resolved geo node.
UPDATE listings l SET currency = 'RUB'
  FROM geo_nodes g
 WHERE l.geo_node_id = g.id AND g.country = 'RU';

UPDATE listings SET currency = 'RUB'
 WHERE geo_node_id IS NULL AND currency = 'KZT'
   AND country IN ('Россия', 'Ресей', 'Russia');

-- +goose Down
ALTER TABLE listings DROP COLUMN IF EXISTS currency;
