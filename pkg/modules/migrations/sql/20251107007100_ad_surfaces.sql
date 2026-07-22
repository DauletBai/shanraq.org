-- +goose Up
-- Move ad orders from a single placement to a set of surfaces. An advertiser
-- now buys the exact places they want (Sport rubric, Culture rubric, home…),
-- priced per surface, instead of one flat "everywhere" package that sold the
-- whole site for a single low fee.
ALTER TABLE ad_orders ADD COLUMN IF NOT EXISTS surfaces TEXT[] NOT NULL DEFAULT '{}';

-- Backfill the new column from the old placement so existing orders keep
-- serving. 'all' expands to every surface; 'rubric' becomes the specific
-- rubric surface it targeted.
UPDATE ad_orders SET surfaces = CASE
    WHEN placement = 'home'       THEN ARRAY['home']
    WHEN placement = 'realestate' THEN ARRAY['realestate']
    WHEN placement = 'articles'   THEN ARRAY['articles']
    WHEN placement = 'rubric'     THEN ARRAY['rubric:' || COALESCE(NULLIF(rubric,''), 'society')]
    WHEN placement = 'all'        THEN ARRAY['home','realestate','articles',
                                             'rubric:sport','rubric:society','rubric:politics',
                                             'rubric:economy','rubric:culture','rubric:technology',
                                             'rubric:it','rubric:opinion','rubric:world']
    ELSE ARRAY['home']
  END
 WHERE surfaces = '{}';

-- Slot-availability queries look up orders whose surfaces contain the one being
-- served, so the array needs a GIN index for @> / && to stay fast.
CREATE INDEX IF NOT EXISTS idx_ad_orders_surfaces ON ad_orders USING GIN (surfaces);

-- +goose Down
DROP INDEX IF EXISTS idx_ad_orders_surfaces;
ALTER TABLE ad_orders DROP COLUMN IF EXISTS surfaces;
