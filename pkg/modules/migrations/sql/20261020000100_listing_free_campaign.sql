-- +goose Up
-- Preserve every currently live, published listing through the six-month
-- introductory period. Expired or moderated listings are not revived.
UPDATE listings
SET expires_at = GREATEST(expires_at, TIMESTAMPTZ '2027-04-20 00:00:00+05'),
    expiry_reminded = false
WHERE status = 'published' AND expires_at > NOW();

-- +goose Down
-- Never shorten owners' already announced expiry dates on rollback.
SELECT 1;
