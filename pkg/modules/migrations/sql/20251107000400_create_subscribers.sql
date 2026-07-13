-- +goose Up
CREATE TABLE IF NOT EXISTS subscribers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    lang TEXT NOT NULL DEFAULT 'ru',
    unsubscribe_token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT subscribers_lang_chk CHECK (lang IN ('kz', 'ru', 'en'))
);

-- Singleton row tracking when the last weekly digest went out.
CREATE TABLE IF NOT EXISTS digest_state (
    id INT PRIMARY KEY DEFAULT 1,
    last_sent_at TIMESTAMPTZ,
    CONSTRAINT digest_state_singleton CHECK (id = 1)
);
INSERT INTO digest_state (id, last_sent_at) VALUES (1, NULL) ON CONFLICT (id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS digest_state;
DROP TABLE IF EXISTS subscribers;
