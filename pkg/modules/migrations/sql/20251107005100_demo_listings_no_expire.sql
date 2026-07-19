-- +goose Up
-- The new background purge permanently deletes listings past their expiry. The
-- six demo listings are a permanent showcase (like the AI columnist's starter
-- columns), so push their expiry far into the future to keep them out of the
-- purge in development. In production the demo fixtures are stripped separately.
UPDATE listings
   SET expires_at = NOW() + INTERVAL '100 years'
 WHERE id::text LIKE 'a0000000-0000-0000-0000-%';

-- +goose Down
SELECT 1;
