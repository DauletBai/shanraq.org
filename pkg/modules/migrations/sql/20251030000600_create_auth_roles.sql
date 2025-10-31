-- +goose Up
CREATE TABLE IF NOT EXISTS auth_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth_user_roles (
    user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES auth_roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS auth_user_roles_role_idx ON auth_user_roles(role_id);

INSERT INTO auth_roles (name, description)
VALUES
    ('admin', 'Full access across Shanraq control plane.'),
    ('operator', 'Operational management of jobs, telemetry, and console.'),
    ('user', 'Default application role with minimal privileges.')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;

INSERT INTO auth_user_roles (user_id, role_id, assigned_at)
SELECT u.id, r.id, NOW()
FROM auth_users u
JOIN auth_roles r ON LOWER(COALESCE(u.role, '')) = r.name
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO auth_user_roles (user_id, role_id, assigned_at)
SELECT u.id, r.id, NOW()
FROM auth_users u
JOIN auth_roles r ON r.name = 'user'
WHERE NOT EXISTS (
    SELECT 1
    FROM auth_user_roles aur
    WHERE aur.user_id = u.id AND aur.role_id = r.id
);

-- +goose Down
DROP TABLE IF EXISTS auth_user_roles;
DROP TABLE IF EXISTS auth_roles;
