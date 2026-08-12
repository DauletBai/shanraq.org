-- +goose Up
-- A revocation counter for issued tokens.
--
-- Roles live inside the JWT and the guards read them from there, so nothing was
-- consulting the database on a privileged request. A demoted administrator kept
-- administrator rights, and a deleted account kept working, until the token
-- expired on its own — two hours in production. There is no session table to
-- delete a row from: the cookie is the credential.
--
-- The counter is stamped into every token at issue and compared on privileged
-- requests. Bumping it retires every token that account holds, at once.
ALTER TABLE auth_users ADD COLUMN IF NOT EXISTS auth_version INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE auth_users DROP COLUMN IF EXISTS auth_version;
