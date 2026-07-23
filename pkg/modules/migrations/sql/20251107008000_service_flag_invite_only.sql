-- +goose Up
-- Allow the invite_only status on service flags (free-function staged-launch
-- gate: available to staff + invited users only).
ALTER TABLE service_flags DROP CONSTRAINT IF EXISTS service_flags_status_check;
ALTER TABLE service_flags ADD CONSTRAINT service_flags_status_check
    CHECK (status IN ('on', 'invite_only', 'maintenance', 'off'));

-- +goose Down
ALTER TABLE service_flags DROP CONSTRAINT IF EXISTS service_flags_status_check;
ALTER TABLE service_flags ADD CONSTRAINT service_flags_status_check
    CHECK (status IN ('on', 'maintenance', 'off'));
