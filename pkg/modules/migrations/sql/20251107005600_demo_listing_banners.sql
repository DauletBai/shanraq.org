-- +goose Up
-- Showcase the paid sidebar banner slot on a fresh install: give the two
-- flagship demo listings a long-running banner. Demo listings are deleted in
-- production by stripDemoFixtures, so this is dev/showcase only.
UPDATE listings
   SET banner_until = NOW() + INTERVAL '100 years'
 WHERE id IN ('d0000000-0000-0000-0000-000000000001',
              'd0000000-0000-0000-0000-000000000002');

-- +goose Down
UPDATE listings SET banner_until = NULL
 WHERE id IN ('d0000000-0000-0000-0000-000000000001',
              'd0000000-0000-0000-0000-000000000002');
