-- +goose Up
-- The admin panel could show how many accounts exist by role and nothing else —
-- not a name, not an e-mail, not where anyone signed up from. Country is the one
-- of those the database did not already hold.
--
-- It is the ISO code resolved from the request IP at registration and nothing
-- more: no continuing location tracking, no per-visit history. The privacy
-- policy already declares the IP address itself as technical data kept for
-- security and diagnostics, and a two-letter country is a coarsening of what is
-- in the access log anyway. Empty for everyone who registered before this, and
-- for anyone whose address cannot be resolved.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS signup_country char(2);

-- +goose Down
ALTER TABLE auth_users DROP COLUMN IF EXISTS signup_country;
