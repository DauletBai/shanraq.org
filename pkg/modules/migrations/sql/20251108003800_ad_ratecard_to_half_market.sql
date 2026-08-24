-- +goose Up
-- The ad card, re-anchored at half the market instead of a fraction of it.
--
-- The previous migration cut the card to a base of 12 500 ₸. That was an
-- overcorrection: a Kostanay region booking came to 5 000 ₸ against ng.kz's
-- 54 000 for the equivalent month, eleven times under. The intent was to
-- undercut the market, not to leave it.
--
-- The error was applying a geography discount on top of a base that had already
-- been set below market, so the two reductions compounded. The rule now: set the
-- base so the COMPARABLE targeted product costs half what the competitor asks,
-- and let the nationwide figure fall out of it.
--
--   Kostanay region, top banner, 30 days:  27 000 ₸  = half of ng.kz's 54 000
--   Nationwide, same banner and period:   135 000 ₸  ≈ 47% under cifrum.kz
--
-- As before, only rows still holding the value they were seeded with are
-- rewritten; anything an operator set by hand is theirs and is left alone.
UPDATE tariffs SET value = 13500, updated_at = NOW() WHERE key = 'ad.horizontal.3'  AND value = 2500;
UPDATE tariffs SET value = 27000, updated_at = NOW() WHERE key = 'ad.horizontal.7'  AND value = 5000;
UPDATE tariffs SET value = 45000, updated_at = NOW() WHERE key = 'ad.horizontal.14' AND value = 8300;
UPDATE tariffs SET value = 67500, updated_at = NOW() WHERE key = 'ad.horizontal.30' AND value = 12500;
UPDATE tariffs SET value = 10500, updated_at = NOW() WHERE key = 'ad.vertical.3'    AND value = 1900;
UPDATE tariffs SET value = 21000, updated_at = NOW() WHERE key = 'ad.vertical.7'    AND value = 3900;
UPDATE tariffs SET value = 35300, updated_at = NOW() WHERE key = 'ad.vertical.14'   AND value = 6500;
UPDATE tariffs SET value = 52500, updated_at = NOW() WHERE key = 'ad.vertical.30'   AND value = 9800;
UPDATE tariffs SET value =  6800, updated_at = NOW() WHERE key = 'ad.square.3'      AND value = 1300;
UPDATE tariffs SET value = 13500, updated_at = NOW() WHERE key = 'ad.square.7'      AND value = 2500;
UPDATE tariffs SET value = 22500, updated_at = NOW() WHERE key = 'ad.square.14'     AND value = 4200;
UPDATE tariffs SET value = 33800, updated_at = NOW() WHERE key = 'ad.square.30'     AND value = 6300;
UPDATE tariffs SET value =  6000, updated_at = NOW() WHERE key = 'ad.rectangle.3'   AND value = 1100;
UPDATE tariffs SET value = 12000, updated_at = NOW() WHERE key = 'ad.rectangle.7'   AND value = 2200;
UPDATE tariffs SET value = 20300, updated_at = NOW() WHERE key = 'ad.rectangle.14'  AND value = 3700;
UPDATE tariffs SET value = 30000, updated_at = NOW() WHERE key = 'ad.rectangle.30'  AND value = 5500;

-- The floor rises with the card: at 1 000 ₸ it sat below what a village booking
-- now costs on its own, so it never applied and the smallest orders were priced
-- under the cost of handling them.
UPDATE tariffs SET value = 2000, updated_at = NOW() WHERE key = 'geo.min_price' AND value = 1000;

-- +goose Down
UPDATE tariffs SET value = 1000 WHERE key = 'geo.min_price' AND value = 2000;
UPDATE tariffs SET value = 2500  WHERE key = 'ad.horizontal.3'  AND value = 13500;
UPDATE tariffs SET value = 5000  WHERE key = 'ad.horizontal.7'  AND value = 27000;
UPDATE tariffs SET value = 8300  WHERE key = 'ad.horizontal.14' AND value = 45000;
UPDATE tariffs SET value = 12500 WHERE key = 'ad.horizontal.30' AND value = 67500;
UPDATE tariffs SET value = 1900  WHERE key = 'ad.vertical.3'    AND value = 10500;
UPDATE tariffs SET value = 3900  WHERE key = 'ad.vertical.7'    AND value = 21000;
UPDATE tariffs SET value = 6500  WHERE key = 'ad.vertical.14'   AND value = 35300;
UPDATE tariffs SET value = 9800  WHERE key = 'ad.vertical.30'   AND value = 52500;
UPDATE tariffs SET value = 1300  WHERE key = 'ad.square.3'      AND value = 6800;
UPDATE tariffs SET value = 2500  WHERE key = 'ad.square.7'      AND value = 13500;
UPDATE tariffs SET value = 4200  WHERE key = 'ad.square.14'     AND value = 22500;
UPDATE tariffs SET value = 6300  WHERE key = 'ad.square.30'     AND value = 33800;
UPDATE tariffs SET value = 1100  WHERE key = 'ad.rectangle.3'   AND value = 6000;
UPDATE tariffs SET value = 2200  WHERE key = 'ad.rectangle.7'   AND value = 12000;
UPDATE tariffs SET value = 3700  WHERE key = 'ad.rectangle.14'  AND value = 20300;
UPDATE tariffs SET value = 5500  WHERE key = 'ad.rectangle.30'  AND value = 30000;
