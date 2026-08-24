-- +goose Up
-- Exchange rates we keep ourselves.
--
-- The National Bank serves historical rates for roughly the last five years: a
-- request for 2020 or earlier comes back empty. So the source is not an archive but
-- a window, and it moves forward with today. Hence we file every day away for
-- ourselves: a year from now our table will be deeper than what the bank is willing
-- to give, and "all time" will stop running into somebody else's window.
--
-- quant is how many units of the currency the rate is quoted for (per 100 yen, per
-- 10 hryvnia). It is stored beside the value because without it currencies cannot
-- be compared.
CREATE TABLE IF NOT EXISTS fx_rates (
    day   date          NOT NULL,
    code  text          NOT NULL,
    value numeric(14,4) NOT NULL,
    quant integer       NOT NULL DEFAULT 1,
    name  text          NOT NULL DEFAULT '',
    PRIMARY KEY (day, code)
);

-- One currency's series over a period — the page's main query.
CREATE INDEX IF NOT EXISTS idx_fx_rates_code_day ON fx_rates (code, day);

-- The probe log: which day we have already asked about and how many rates came
-- back.
--
-- Without it the backfill would ask the bank about every weekend of the last five
-- years again after each restart — a thousand pointless requests. A zero in found is
-- the answer "no rates for this day", and it is as much a result as a number: it
-- marks both the weekends and the boundary past which the source is silent.
CREATE TABLE IF NOT EXISTS fx_probed (
    day        date        NOT NULL PRIMARY KEY,
    found      integer     NOT NULL DEFAULT 0,
    probed_at  timestamptz NOT NULL DEFAULT now()
);

-- The monthly archive from November 1993 — from the tenge's birth.
--
-- The National Bank serves daily rates about five years back and is silent about
-- everything earlier. So the depth comes from the Bank for International
-- Settlements: it holds a monthly series of rates against the dollar from 1993-11,
-- and a cross rate off it gives tenge per any of our currencies. Tenge per ONE unit
-- is what is stored; the familiar "per 100 yen" is done by the display.
CREATE TABLE IF NOT EXISTS fx_monthly (
    month date          NOT NULL,
    code  text          NOT NULL,
    value numeric(18,6) NOT NULL,
    PRIMARY KEY (month, code)
);
CREATE INDEX IF NOT EXISTS idx_fx_monthly_code_month ON fx_monthly (code, month);

-- +goose Down
DROP TABLE IF EXISTS fx_monthly;
DROP TABLE IF EXISTS fx_probed;
DROP TABLE IF EXISTS fx_rates;
