-- +goose Up
-- Staff roles for the admin dashboard. Claims carry auth_users.role verbatim,
-- so setting a user's role to one of these grants the matching access.
-- director = superadmin (equivalent to admin), then manager and editor.
INSERT INTO auth_roles (name, description) VALUES
    ('director', 'Superadmin: full access — users, roles, finance, content, settings.'),
    ('manager',  'Analytics, advertising/listings, finance (view). No role management.'),
    ('editor',   'Content and comment moderation, rubrics. No finance.')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM auth_roles WHERE name IN ('director', 'manager', 'editor');
