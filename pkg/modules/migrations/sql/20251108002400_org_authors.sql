-- +goose Up
-- An organisation as author: the "Kacharets" housing service, "Vodokanal" LLP, a
-- city mayor's office.
--
-- An organisation is not an account but a right an account holds. A person stays
-- behind every publication: our rule is that everyone uses their real name so that
-- somebody specific stands behind a text, and a shared departmental login cancels
-- that rule — afterwards there is no finding who wrote it. Staff change, an account
-- outlives a person, and by law a responsible individual stands behind a mayoral
-- office's message, not a login.
--
-- Hence: the byline shows the organisation, the moderation log shows the person.
--
-- Verification is mandatory and it is the whole substance of this. A name anyone can
-- type in turns the platform into a machine for impersonating a mayor's office, and
-- one forged publication in its name costs more than the good of a hundred genuine
-- ones. Until status is 'verified' the name is shown nowhere.
--
-- The shape deliberately mirrors re_agents: the business number, the moderation
-- queue and the verified badge already work there, and a second entity with the same
-- needs would cost more than a shared one.
CREATE TABLE IF NOT EXISTS org_authors (
    user_id UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    kind TEXT NOT NULL,
    bin TEXT NOT NULL DEFAULT '',
    -- The organisation's territory. The Kostanay mayor's office publishes for
    -- Kostanay and what lies inside it; otherwise a verified badge becomes a pass
    -- to the whole country.
    geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL,
    about TEXT NOT NULL DEFAULT '',
    contact TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    reject_reason TEXT NOT NULL DEFAULT '',
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT org_authors_status_chk CHECK (status IN ('pending','verified','rejected'))
);

CREATE INDEX IF NOT EXISTS idx_org_authors_status ON org_authors (status);
CREATE INDEX IF NOT EXISTS idx_org_authors_place ON org_authors (geo_node_id) WHERE geo_node_id IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS org_authors;
