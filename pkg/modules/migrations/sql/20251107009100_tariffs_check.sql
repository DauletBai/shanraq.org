-- +goose Up
-- Belt-and-suspenders: enforce the tariff range at the DB level too, so a bad
-- value can never reach the pricing math even if a future code path bypasses the
-- app-side clamp. Mirrors maxTariffValue (10 000 000).
ALTER TABLE tariffs ADD CONSTRAINT tariffs_value_range CHECK (value >= 0 AND value <= 10000000);

-- +goose Down
ALTER TABLE tariffs DROP CONSTRAINT IF EXISTS tariffs_value_range;
