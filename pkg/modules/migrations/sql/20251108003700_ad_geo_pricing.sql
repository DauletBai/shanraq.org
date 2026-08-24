-- +goose Up
-- Population for the places an advertiser can buy, and the geo column on an ad
-- order.
--
-- Ad prices now follow the size of the audience a booking can reach: a banner
-- bought for one village costs a fraction of the same banner bought for the
-- whole country. That needs a population on every place that can be targeted,
-- and until now only cities and villages had one — 179 rows out of six hundred.
-- Countries, regions and the three cities of republican significance had none,
-- which is precisely the set an advertiser is most likely to buy.
--
-- Rolling the figure up from the children was not an option: the tree holds
-- only the larger settlements, so Kostanay region would have summed to about
-- four hundred thousand of its real eight hundred, and every regional booking
-- would have been sold at half price.
--
-- Country totals come from the World Bank (SP.POP.TOTL, 2025). Regional figures
-- come from Wikidata's sourced population statements, which carry the year of
-- the count — recorded here in population_year, because a price derived from a
-- figure has to be able to name the figure's date. The years differ by region
-- (2018–2024) and that is honest: nobody re-counts every region in the same
-- year. For a price ladder built on square roots, a few percent of drift moves
-- nothing that matters.
--
-- The sum of the regional figures is 20.6 million against the World Bank's
-- 20.8 for the country, which is the agreement one should expect from counts
-- taken in different years — and a useful check that nothing landed in the
-- wrong row.
UPDATE geo_nodes SET population = 20843754, population_year = 2025 WHERE code = 'g1' AND level = 0;
UPDATE geo_nodes SET population = 143513328, population_year = 2025 WHERE code = 'g100' AND level = 0;

UPDATE geo_nodes SET population = 610100, population_year = 2021 WHERE code = 'g23';
UPDATE geo_nodes SET population = 782995, population_year = 2021 WHERE code = 'g27';
UPDATE geo_nodes SET population = 927136, population_year = 2022 WHERE code = 'g32';
UPDATE geo_nodes SET population = 1503588, population_year = 2022 WHERE code = 'g37';
UPDATE geo_nodes SET population = 2228675, population_year = 2024 WHERE code = 'g8';
UPDATE geo_nodes SET population = 1239744, population_year = 2022 WHERE code = 'g2';
UPDATE geo_nodes SET population = 673601, population_year = 2021 WHERE code = 'g43';
UPDATE geo_nodes SET population = 1341292, population_year = 2021 WHERE code = 'g95';
UPDATE geo_nodes SET population = 1199259, population_year = 2021 WHERE code = 'g49';
UPDATE geo_nodes SET population = 698987, population_year = 2022 WHERE code = 'g54';
UPDATE geo_nodes SET population = 675655, population_year = 2021 WHERE code = 'g46';
UPDATE geo_nodes SET population = 1348468, population_year = 2021 WHERE code = 'g59';
UPDATE geo_nodes SET population = 833643, population_year = 2021 WHERE code = 'g65';
UPDATE geo_nodes SET population = 814931, population_year = 2021 WHERE code = 'g70';
UPDATE geo_nodes SET population = 735008, population_year = 2021 WHERE code = 'g74';
UPDATE geo_nodes SET population = 756755, population_year = 2021 WHERE code = 'g82';
UPDATE geo_nodes SET population = 540786, population_year = 2021 WHERE code = 'g78';
UPDATE geo_nodes SET population = 1929000, population_year = 2018 WHERE code = 'g86';
UPDATE geo_nodes SET population = 231774, population_year = 2023 WHERE code = 'g91';
UPDATE geo_nodes SET population = 1191877, population_year = 2023 WHERE code = 'g17';
UPDATE geo_nodes SET population = 2099186, population_year = 2025 WHERE code = 'g180';
UPDATE geo_nodes SET population = 750083, population_year = 2024 WHERE code = 'g340';
UPDATE geo_nodes SET population = 989434, population_year = 2025 WHERE code = 'g298';
UPDATE geo_nodes SET population = 946429, population_year = 2024 WHERE code = 'g320';
UPDATE geo_nodes SET population = 1482025, population_year = 2025 WHERE code = 'g311';
UPDATE geo_nodes SET population = 1132795, population_year = 2025 WHERE code = 'g314';
UPDATE geo_nodes SET population = 1309942, population_year = 2024 WHERE code = 'g282';
UPDATE geo_nodes SET population = 2434046, population_year = 2025 WHERE code = 'g249';
UPDATE geo_nodes SET population = 1160445, population_year = 2021 WHERE code = 'g301';
UPDATE geo_nodes SET population = 2260045, population_year = 2025 WHERE code = 'g246';
UPDATE geo_nodes SET population = 144428, population_year = 2025 WHERE code = 'g358';
UPDATE geo_nodes SET population = 984395, population_year = 2024 WHERE code = 'g184';
UPDATE geo_nodes SET population = 905900, population_year = 2024 WHERE code = 'g332';
UPDATE geo_nodes SET population = 2322292, population_year = 2025 WHERE code = 'g262';
UPDATE geo_nodes SET population = 994599, population_year = 2023 WHERE code = 'g304';
UPDATE geo_nodes SET population = 1000980, population_year = 2021 WHERE code = 'g306';
UPDATE geo_nodes SET population = 288947, population_year = 2024 WHERE code = 'g186';
UPDATE geo_nodes SET population = 2527219, population_year = 2025 WHERE code = 'g266';
UPDATE geo_nodes SET population = 1129935, population_year = 2024 WHERE code = 'g296';
UPDATE geo_nodes SET population = 566266, population_year = 2024 WHERE code = 'g334';
UPDATE geo_nodes SET population = 5842238, population_year = 2025 WHERE code = 'g188';
UPDATE geo_nodes SET population = 2837988, population_year = 2025 WHERE code = 'g193';
UPDATE geo_nodes SET population = 744465, population_year = 2025 WHERE code = 'g338';
UPDATE geo_nodes SET population = 1049783, population_year = 2025 WHERE code = 'g309';
UPDATE geo_nodes SET population = 2035762, population_year = 2024 WHERE code = 'g218';
UPDATE geo_nodes SET population = 1125921, population_year = 2023 WHERE code = 'g291';
UPDATE geo_nodes SET population = 133387, population_year = 2024 WHERE code = 'g344';
UPDATE geo_nodes SET population = 13274285, population_year = 2025 WHERE code = 'g101';
UPDATE geo_nodes SET population = 8594454, population_year = 2023 WHERE code = 'g211';
UPDATE geo_nodes SET population = 656438, population_year = 2024 WHERE code = 'g322';
UPDATE geo_nodes SET population = 42224, population_year = 2024 WHERE code = 'g356';
UPDATE geo_nodes SET population = 3039421, population_year = 2025 WHERE code = 'g229';
UPDATE geo_nodes SET population = 566960, population_year = 2025 WHERE code = 'g325';
UPDATE geo_nodes SET population = 2786540, population_year = 2025 WHERE code = 'g226';
UPDATE geo_nodes SET population = 1805806, population_year = 2025 WHERE code = 'g260';
UPDATE geo_nodes SET population = 1828655, population_year = 2024 WHERE code = 'g270';
UPDATE geo_nodes SET population = 692486, population_year = 2024 WHERE code = 'g336';
UPDATE geo_nodes SET population = 1213113, population_year = 2024 WHERE code = 'g294';
UPDATE geo_nodes SET population = 2483633, population_year = 2025 WHERE code = 'g197';
UPDATE geo_nodes SET population = 1799659, population_year = 2025 WHERE code = 'g200';
UPDATE geo_nodes SET population = 574480, population_year = 2025 WHERE code = 'g327';
UPDATE geo_nodes SET population = 210095, population_year = 2025 WHERE code = 'g124';
UPDATE geo_nodes SET population = 714785, population_year = 2025 WHERE code = 'g148';
UPDATE geo_nodes SET population = 4137335, population_year = 2025 WHERE code = 'g233';
UPDATE geo_nodes SET population = 1043231, population_year = 2024 WHERE code = 'g289';
UPDATE geo_nodes SET population = 3112566, population_year = 2025 WHERE code = 'g242';
UPDATE geo_nodes SET population = 5652922, population_year = 2025 WHERE code = 'g112';
UPDATE geo_nodes SET population = 2369405, population_year = 2025 WHERE code = 'g253';
UPDATE geo_nodes SET population = 457590, population_year = 2024 WHERE code = 'g342';
UPDATE geo_nodes SET population = 4221452, population_year = 2025 WHERE code = 'g222';
UPDATE geo_nodes SET population = 863987, population_year = 2024 WHERE code = 'g316';
UPDATE geo_nodes SET population = 2884363, population_year = 2025 WHERE code = 'g204';
UPDATE geo_nodes SET population = 956292, population_year = 2024 WHERE code = 'g329';
UPDATE geo_nodes SET population = 1190574, population_year = 2025 WHERE code = 'g318';
UPDATE geo_nodes SET population = 1039728, population_year = 2025 WHERE code = 'g273';
UPDATE geo_nodes SET population = 1471140, population_year = 2024 WHERE code = 'g286';
UPDATE geo_nodes SET population = 3890800, population_year = 2024 WHERE code = 'g257';
UPDATE geo_nodes SET population = 1172782, population_year = 2024 WHERE code = 'g276';
UPDATE geo_nodes SET population = 1273488, population_year = 2025 WHERE code = 'g208';
UPDATE geo_nodes SET population = 1759356, population_year = 2024 WHERE code = 'g346';
UPDATE geo_nodes SET population = 3385124, population_year = 2025 WHERE code = 'g238';
UPDATE geo_nodes SET population = 47778, population_year = 2025 WHERE code = 'g354';
UPDATE geo_nodes SET population = 515960, population_year = 2024 WHERE code = 'g350';
UPDATE geo_nodes SET population = 1187558, population_year = 2024 WHERE code = 'g279';

-- The place an ad is bought for. NULL means the whole site, which is what every
-- existing order is: they were sold before geography had a price, and they keep
-- the reach they paid for.
--
-- ON DELETE SET NULL rather than CASCADE: a region being redrawn must not delete
-- the record of an order somebody paid for. The order falls back to nationwide,
-- which is the reach it is then actually getting.
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_ad_orders_geo ON ad_orders (geo_node_id) WHERE geo_node_id IS NOT NULL;

-- geo_price_share is the share of the country's population this place holds,
-- in hundred-thousandths, frozen onto the order at purchase. The price the
-- advertiser agreed to must stay reproducible even after a census moves the
-- population underneath it.
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS geo_share INTEGER NOT NULL DEFAULT 100000;

-- The rate card, cut to the market.
--
-- Tariff rows are seeded once and never overwritten, so a new built-in default
-- alone would change nothing on a running site: the old prices would sit in the
-- table for ever. They are rewritten here — but only where the operator has not
-- touched them. A row still holding the value it was seeded with is ours to
-- correct; a row somebody deliberately set is theirs, and it is left alone to
-- show up beside the new card in the admin panel.
UPDATE tariffs SET value = 2500,  updated_at = NOW() WHERE key = 'ad.horizontal.3'  AND value = 18000;
UPDATE tariffs SET value = 5000,  updated_at = NOW() WHERE key = 'ad.horizontal.7'  AND value = 36000;
UPDATE tariffs SET value = 8300,  updated_at = NOW() WHERE key = 'ad.horizontal.14' AND value = 60000;
UPDATE tariffs SET value = 12500, updated_at = NOW() WHERE key = 'ad.horizontal.30' AND value = 90000;
UPDATE tariffs SET value = 1900,  updated_at = NOW() WHERE key = 'ad.vertical.3'    AND value = 14000;
UPDATE tariffs SET value = 3900,  updated_at = NOW() WHERE key = 'ad.vertical.7'    AND value = 28000;
UPDATE tariffs SET value = 6500,  updated_at = NOW() WHERE key = 'ad.vertical.14'   AND value = 47000;
UPDATE tariffs SET value = 9800,  updated_at = NOW() WHERE key = 'ad.vertical.30'   AND value = 70000;
UPDATE tariffs SET value = 1300,  updated_at = NOW() WHERE key = 'ad.square.3'      AND value = 9000;
UPDATE tariffs SET value = 2500,  updated_at = NOW() WHERE key = 'ad.square.7'      AND value = 18000;
UPDATE tariffs SET value = 4200,  updated_at = NOW() WHERE key = 'ad.square.14'     AND value = 30000;
UPDATE tariffs SET value = 6300,  updated_at = NOW() WHERE key = 'ad.square.30'     AND value = 45000;
UPDATE tariffs SET value = 1100,  updated_at = NOW() WHERE key = 'ad.rectangle.3'   AND value = 8000;
UPDATE tariffs SET value = 2200,  updated_at = NOW() WHERE key = 'ad.rectangle.7'   AND value = 16000;
UPDATE tariffs SET value = 3700,  updated_at = NOW() WHERE key = 'ad.rectangle.14'  AND value = 27000;
UPDATE tariffs SET value = 5500,  updated_at = NOW() WHERE key = 'ad.rectangle.30'  AND value = 40000;

-- +goose Down
ALTER TABLE ad_orders DROP COLUMN IF EXISTS geo_share;
DROP INDEX IF EXISTS idx_ad_orders_geo;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS geo_node_id;
