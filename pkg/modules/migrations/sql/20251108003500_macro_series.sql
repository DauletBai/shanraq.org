-- +goose Up
-- Macro series: what the exchange rate and inflation are made of.
--
-- One table for every indicator rather than a table per source: there are few series,
-- they are all the same shape — a period and a number — and they constantly have to
-- be compared with one another. Spreading them across separate tables would mean
-- writing a join where a WHERE does.
--
-- period is the first of the month for monthly series and the first of January for
-- annual ones.
CREATE TABLE IF NOT EXISTS macro_series (
    code   text          NOT NULL,
    period date          NOT NULL,
    value  numeric(20,4) NOT NULL,
    PRIMARY KEY (code, period)
);
CREATE INDEX IF NOT EXISTS idx_macro_series_code ON macro_series (code, period);

-- +goose Down
DROP TABLE IF EXISTS macro_series;
