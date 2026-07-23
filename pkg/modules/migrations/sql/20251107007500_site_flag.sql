-- +goose Up
-- The global site switch. Status 'on' means the site serves normally; any other
-- status makes every page (except staff sessions and the admin/login recovery
-- routes) return a maintenance page. Seeded 'on' with a ready trilingual notice
-- so a nightly full reload / major change can be announced with one click.
INSERT INTO service_flags (code, status, message_kz, message_ru, message_en) VALUES
  ('site', 'on',
   'Сайтта техникалық жұмыстар жүріп жатыр. Жақын арада қайта ораламыз. Түсіністікпен қарағаныңызға рахмет.',
   'На сайте идут технические работы. Мы скоро вернёмся. Спасибо за понимание.',
   'The site is undergoing maintenance. We will be back shortly. Thank you for your patience.')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DELETE FROM service_flags WHERE code = 'site';
