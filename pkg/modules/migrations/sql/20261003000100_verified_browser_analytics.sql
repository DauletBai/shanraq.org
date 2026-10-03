-- +goose Up
-- The request counter could exclude declared crawlers and known hosting ASNs,
-- but a residential-proxy scanner with a normal Chrome User-Agent was still
-- indistinguishable from a reader. On 2026-10-02 one such client generated 634
-- page requests in a single half-hour and entered every audience figure.
--
-- From this migration onward a human view is written only after a same-origin
-- browser beacon proves that the rendered page was visible. The old aggregates
-- cannot be repaired after the fact, so they are archived rather than mixed
-- with the verified series.
CREATE TABLE IF NOT EXISTS analytics_daily_unverified
    (LIKE analytics_daily INCLUDING ALL);

INSERT INTO analytics_daily_unverified
SELECT * FROM analytics_daily
WHERE kind IN ('page', 'source', 'device', 'os', 'browser', 'country', 'lang',
               'geolang', 'course_hub', 'course_lesson')
ON CONFLICT (day, kind, label, is_guest) DO UPDATE
SET n = EXCLUDED.n;

DELETE FROM analytics_daily
WHERE kind IN ('page', 'source', 'device', 'os', 'browser', 'country', 'lang',
               'geolang', 'course_hub', 'course_lesson');

CREATE TABLE IF NOT EXISTS analytics_slots_unverified
    (LIKE analytics_slots INCLUDING ALL);
INSERT INTO analytics_slots_unverified SELECT * FROM analytics_slots
ON CONFLICT (slot, vid) DO UPDATE SET views = EXCLUDED.views;
DELETE FROM analytics_slots;

-- sid is a daily salted HMAC of a short-lived first-party cookie. It makes a
-- visit a real rolling 30-minute browser session instead of an arbitrary fixed
-- half-hour bucket. The cookie value itself is never stored.
ALTER TABLE analytics_slots ADD COLUMN IF NOT EXISTS sid BYTEA;
ALTER TABLE analytics_slots ALTER COLUMN sid SET NOT NULL;

-- Reader-facing article and listing totals were fed by the same request rule.
-- Preserve them separately and restart the visible values from verified views.
ALTER TABLE articles ADD COLUMN IF NOT EXISTS views_unverified BIGINT NOT NULL DEFAULT 0;
UPDATE articles
SET views_unverified = views_unverified + views_count,
    views_count = 0;

CREATE TABLE IF NOT EXISTS article_views_daily_unverified
    (LIKE article_views_daily INCLUDING ALL);
INSERT INTO article_views_daily_unverified SELECT * FROM article_views_daily
ON CONFLICT (article_id, lang, day) DO UPDATE SET views = EXCLUDED.views;
DELETE FROM article_views_daily;

ALTER TABLE listings ADD COLUMN IF NOT EXISTS views_unverified BIGINT NOT NULL DEFAULT 0;
UPDATE listings
SET views_unverified = views_unverified + views_count,
    views_count = 0;

-- Keep the editable, deployed privacy page aligned with the new first-party
-- measurement rule. It still stores no path, IP address or reusable visitor
-- identity; the only added browser state is the short visit cookie described
-- here.
UPDATE content_pages
SET body_md = replace(
        replace(body_md, '_Редакция от 29.09.2026._', '_Редакция от 03.10.2026._'),
        '**Что именно считается.** Хосты, посетители, визиты и просмотры — по получасовым интервалам, с разбивкой на Казахстан и остальной мир, мобильные и прочие устройства. Эти же цифры открыто показаны всем на странице «Аналитика»: мы не держим для себя данных, которых не показываем вам.',
        '**Что именно считается.** Просмотр записывается только после того, как страница отобразилась в видимой вкладке браузера. Визит объединяет просмотры одного браузера, пока между ними нет 30 минут бездействия; для этого браузер хранит технический cookie визита, а база — только его необратимый суточный HMAC. Хосты, посетители, визиты и просмотры разбиваются на Казахстан и остальной мир, мобильные и прочие устройства. Эти же цифры открыто показаны всем на странице «Аналитика»: мы не держим для себя данных, которых не показываем вам.'
    )
WHERE page_key = 'privacy' AND lang = 'ru';

UPDATE content_pages
SET body_md = replace(
        replace(body_md, '_29.09.2026 жағдайындағы редакция._', '_03.10.2026 жағдайындағы редакция._'),
        '**Не саналады.** Хосттар, келушілер, кірулер және қаралымдар — жарты сағаттық аралықтармен, Қазақстан мен қалған әлемге, мобильді және өзге құрылғыларға бөлініп. Дәл осы сандар «Аналитика» бетінде барлығына ашық көрсетілген: өзімізге көрсетпейтін дерек ұстамаймыз.',
        '**Не саналады.** Қаралым бет браузердің көрінетін қойындысында көрсетілгеннен кейін ғана жазылады. Бір браузердің қаралымдары арасында 30 минут әрекетсіздік болмаса, олар бір кіруге біріктіріледі; бұл үшін браузер кірудің техникалық cookie файлын, ал дерекқор оның тек қайтарымсыз тәуліктік HMAC мәнін сақтайды. Хосттар, келушілер, кірулер және қаралымдар Қазақстан мен қалған әлемге, мобильді және өзге құрылғыларға бөлінеді. Дәл осы сандар «Аналитика» бетінде барлығына ашық көрсетілген: өзімізге көрсетпейтін дерек ұстамаймыз.'
    )
WHERE page_key = 'privacy' AND lang = 'kz';

UPDATE content_pages
SET body_md = replace(
        replace(body_md, '_Revision of 29.09.2026._', '_Revision of 03.10.2026._'),
        '**What is counted.** Hosts, visitors, visits and views, in half-hour windows, split between Kazakhstan and the rest of the world and between mobile and other devices. The same figures are shown openly to everyone on the Analytics page: we keep no traffic data that we do not show you.',
        '**What is counted.** A view is recorded only after the page has appeared in a visible browser tab. Views from one browser remain one visit until there has been 30 minutes without activity; the browser keeps a technical visit cookie and the database stores only its irreversible daily HMAC. Hosts, visitors, visits and views are split between Kazakhstan and the rest of the world and between mobile and other devices. The same figures are shown openly to everyone on the Analytics page: we keep no traffic data that we do not show you.'
    )
WHERE page_key = 'privacy' AND lang = 'en';

-- +goose Down
DELETE FROM analytics_daily
WHERE kind IN ('page', 'source', 'device', 'os', 'browser', 'country', 'lang',
               'geolang', 'course_hub', 'course_lesson');
INSERT INTO analytics_daily SELECT * FROM analytics_daily_unverified
ON CONFLICT (day, kind, label, is_guest) DO UPDATE SET n = EXCLUDED.n;

DELETE FROM analytics_slots;
ALTER TABLE analytics_slots DROP COLUMN IF EXISTS sid;
INSERT INTO analytics_slots SELECT * FROM analytics_slots_unverified
ON CONFLICT (slot, vid) DO UPDATE SET views = EXCLUDED.views;

DELETE FROM article_views_daily;
INSERT INTO article_views_daily SELECT * FROM article_views_daily_unverified
ON CONFLICT (article_id, lang, day) DO UPDATE SET views = EXCLUDED.views;
UPDATE articles SET views_count = views_count + views_unverified;
UPDATE listings SET views_count = views_count + views_unverified;

DROP TABLE IF EXISTS article_views_daily_unverified;
DROP TABLE IF EXISTS analytics_slots_unverified;
DROP TABLE IF EXISTS analytics_daily_unverified;
ALTER TABLE articles DROP COLUMN IF EXISTS views_unverified;
ALTER TABLE listings DROP COLUMN IF EXISTS views_unverified;
