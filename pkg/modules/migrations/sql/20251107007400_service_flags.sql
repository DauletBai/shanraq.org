-- +goose Up
-- Operational switches for individual product services, flippable from the
-- admin panel WITHOUT a redeploy. A service can be:
--   on          — fully available;
--   maintenance — visible, but its paid action is disabled and an apology /
--                 "technical works" message is shown in its place;
--   off         — hidden entirely.
-- This lets us launch a beta where paid flows (ad orders, listing promotion)
-- sit in maintenance until Kaspi Pay / Freedom Pay are wired, while every free
-- function keeps working, and lets us take any service down for maintenance
-- later with a localized notice.
CREATE TABLE IF NOT EXISTS service_flags (
    code       TEXT PRIMARY KEY,
    status     TEXT NOT NULL DEFAULT 'on'
               CHECK (status IN ('on', 'maintenance', 'off')),
    message_kz TEXT NOT NULL DEFAULT '',
    message_ru TEXT NOT NULL DEFAULT '',
    message_en TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID
);

-- Seed the two payment-dependent services in maintenance for the beta, with a
-- ready trilingual apology. Admin flips them to 'on' once payments go live.
INSERT INTO service_flags (code, status, message_kz, message_ru, message_en) VALUES
  ('ad_orders', 'maintenance',
   'Жарнама төлемі уақытша қолжетімсіз: төлем жүйесін қосу жүріп жатыр. Кешірім сұраймыз.',
   'Оплата рекламы временно недоступна: подключаем платёжную систему. Приносим извинения.',
   'Paid advertising is temporarily unavailable: we are connecting the payment system. Sorry for the inconvenience.'),
  ('listing_promo', 'maintenance',
   'Хабарландыруды ақылы жылжыту уақытша қолжетімсіз: төлем жүйесін қосу жүріп жатыр. Кешірім сұраймыз.',
   'Платное продвижение объявлений временно недоступно: подключаем платёжную систему. Приносим извинения.',
   'Paid listing promotion is temporarily unavailable: we are connecting the payment system. Sorry for the inconvenience.')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS service_flags;
