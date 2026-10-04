-- +goose Up
-- Account-linked learning history. Aggregate audience analytics remains
-- anonymous; this table is written only for a signed-in account after the
-- browser has proved that the page was rendered in a visible tab. It contains
-- no IP address, visitor cookie, user-agent or location.
CREATE TABLE IF NOT EXISTS user_activity_events (
    id               bigserial   PRIMARY KEY,
    user_id          uuid        NOT NULL REFERENCES auth_users (id) ON DELETE CASCADE,
    occurred_at      timestamptz NOT NULL DEFAULT now(),
    event_type       text        NOT NULL,
    content_kind     text        NOT NULL DEFAULT '',
    path             text        NOT NULL DEFAULT '',
    lang             text        NOT NULL DEFAULT '',
    article_id       uuid        REFERENCES articles (id) ON DELETE SET NULL,
    series_id        uuid        REFERENCES article_series (id) ON DELETE SET NULL,
    depth            smallint    NOT NULL DEFAULT 0,
    duration_seconds integer     NOT NULL DEFAULT 0,
    passed           boolean     NOT NULL DEFAULT false,
    CONSTRAINT user_activity_event_type_chk
        CHECK (event_type IN ('view', 'read', 'check')),
    CONSTRAINT user_activity_depth_chk CHECK (depth BETWEEN 0 AND 100),
    CONSTRAINT user_activity_duration_chk CHECK (duration_seconds >= 0),
    CONSTRAINT user_activity_lang_chk CHECK (lang IN ('', 'kz', 'ru', 'en'))
);

CREATE INDEX IF NOT EXISTS idx_user_activity_user_time
    ON user_activity_events (user_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_activity_user_article
    ON user_activity_events (user_id, article_id, occurred_at DESC)
    WHERE article_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_user_activity_user_series
    ON user_activity_events (user_id, series_id, occurred_at DESC)
    WHERE series_id IS NOT NULL;

-- The previous policy truthfully described anonymous aggregate analytics, but
-- it also said no account-linked path was retained. The registered-reader
-- learning journal changes that statement, so update all three published
-- versions in the same migration that starts collecting it.
UPDATE content_pages
SET body_md = replace(replace(replace(
        replace(body_md, '_Редакция от 03.10.2026._', '_Редакция от 04.10.2026._'),
        '**Что именно считается.** Просмотр записывается только после того, как страница отобразилась в видимой вкладке браузера. Визит объединяет просмотры одного браузера, пока между ними нет 30 минут бездействия; для этого браузер хранит технический cookie визита, а база — только его необратимый суточный HMAC. Хосты, посетители, визиты и просмотры разбиваются на Казахстан и остальной мир, мобильные и прочие устройства. Эти же цифры открыто показаны всем на странице «Аналитика»: мы не держим для себя данных, которых не показываем вам.',
        '**Что именно считается.** Просмотр записывается только после того, как страница отобразилась в видимой вкладке браузера. Визит объединяет просмотры одного браузера, пока между ними нет 30 минут бездействия; для этого браузер хранит технический cookie визита, а база — только его необратимый суточный HMAC. Хосты, посетители, визиты и просмотры разбиваются на Казахстан и остальной мир, мобильные и прочие устройства. Эти агрегированные цифры открыто показаны всем на странице «Аналитика».'),
        '**Как это устроено, чтобы не стать слежкой.** Посетитель считается по анонимному идентификатору. Он вычисляется из IP-адреса и строки браузера ключом, который создаётся заново каждые сутки; сам адрес нигде не сохраняется. Ключ вчерашнего дня не восстанавливается, поэтому один и тот же человек в понедельник и во вторник даёт два несвязанных значения. Проследить за кем-то дольше суток нельзя технически — в том числе нам.',
        '**Как устроена анонимная часть.** Посетитель, который не вошёл в аккаунт, считается по анонимному идентификатору. Он вычисляется из IP-адреса и строки браузера ключом, который создаётся заново каждые сутки; сам адрес нигде не сохраняется. Ключ вчерашнего дня не восстанавливается, поэтому один и тот же незарегистрированный посетитель в понедельник и во вторник даёт два несвязанных значения.'),
        '**Чего мы не делаем в собственной аналитике.** Не записываем ваш путь по сайту, не строим профиль интересов, не передаём данные собственной аналитики третьим лицам и не используем их для подбора рекламы под конкретного человека.',
        '**Для зарегистрированных пользователей** сохраняется история просмотренных статей, курсов и уроков, времени чтения и результатов заданий. Она доступна администраторам, используется для анализа интересов и обучения, не передаётся рекламодателям и удаляется вместе с аккаунтом.'
    )
WHERE page_key = 'privacy' AND lang = 'ru';

UPDATE content_pages
SET body_md = replace(replace(replace(
        replace(body_md, '_03.10.2026 жағдайындағы редакция._', '_04.10.2026 жағдайындағы редакция._'),
        '**Не саналады.** Қаралым бет браузердің көрінетін қойындысында көрсетілгеннен кейін ғана жазылады. Бір браузердің қаралымдары арасында 30 минут әрекетсіздік болмаса, олар бір кіруге біріктіріледі; бұл үшін браузер кірудің техникалық cookie файлын, ал дерекқор оның тек қайтарымсыз тәуліктік HMAC мәнін сақтайды. Хосттар, келушілер, кірулер және қаралымдар Қазақстан мен қалған әлемге, мобильді және өзге құрылғыларға бөлінеді. Дәл осы сандар «Аналитика» бетінде барлығына ашық көрсетілген: өзімізге көрсетпейтін дерек ұстамаймыз.',
        '**Не саналады.** Қаралым бет браузердің көрінетін қойындысында көрсетілгеннен кейін ғана жазылады. Бір браузердің қаралымдары арасында 30 минут әрекетсіздік болмаса, олар бір кіруге біріктіріледі; бұл үшін браузер кірудің техникалық cookie файлын, ал дерекқор оның тек қайтарымсыз тәуліктік HMAC мәнін сақтайды. Хосттар, келушілер, кірулер және қаралымдар Қазақстан мен қалған әлемге, мобильді және өзге құрылғыларға бөлінеді. Бұл жинақталған сандар «Аналитика» бетінде барлығына ашық көрсетілген.'),
        '**Бұл қалай бақылауға айналмайды.** Келуші анонимді идентификатормен саналады. Ол IP-мекенжай мен браузер жолынан күн сайын жаңадан жасалатын кілтпен есептеледі; мекенжайдың өзі еш жерде сақталмайды. Кешегі кілт қалпына келтірілмейді, сондықтан бір адам дүйсенбіде де, сейсенбіде де екі байланыссыз мән береді. Біреуді бір тәуліктен ұзақ қадағалау техникалық тұрғыдан мүмкін емес — бізге де.',
        '**Анонимді бөлік қалай құрылған.** Есептік жазбаға кірмеген келуші анонимді идентификатормен саналады. Ол IP-мекенжай мен браузер жолынан күн сайын жаңадан жасалатын кілтпен есептеледі; мекенжайдың өзі еш жерде сақталмайды. Кешегі кілт қалпына келтірілмейді, сондықтан бір тіркелмеген келуші дүйсенбіде де, сейсенбіде де екі байланыссыз мән береді.'),
        '**Өз аналитикамызда не істемейміз.** Сайттағы жолыңызды жазбаймыз, қызығушылық профилін құрмаймыз, өз аналитикамыздың деректерін үшінші тұлғаларға бермейміз және нақты адамға жарнама таңдау үшін пайдаланбаймыз.',
        '**Тіркелген қолданушылар үшін** қаралған мақалалар, курстар мен сабақтар, оқу уақыты және тапсырма нәтижелері сақталады. Бұл тарих әкімшілерге қолжетімді, қызығушылық пен оқу барысын талдау үшін қолданылады, жарнама берушілерге берілмейді және есептік жазбамен бірге жойылады.'
    )
WHERE page_key = 'privacy' AND lang = 'kz';

UPDATE content_pages
SET body_md = replace(replace(replace(
        replace(body_md, '_Revision of 03.10.2026._', '_Revision of 04.10.2026._'),
        '**What is counted.** A view is recorded only after the page has appeared in a visible browser tab. Views from one browser remain one visit until there has been 30 minutes without activity; the browser keeps a technical visit cookie and the database stores only its irreversible daily HMAC. Hosts, visitors, visits and views are split between Kazakhstan and the rest of the world and between mobile and other devices. The same figures are shown openly to everyone on the Analytics page: we keep no traffic data that we do not show you.',
        '**What is counted.** A view is recorded only after the page has appeared in a visible browser tab. Views from one browser remain one visit until there has been 30 minutes without activity; the browser keeps a technical visit cookie and the database stores only its irreversible daily HMAC. Hosts, visitors, visits and views are split between Kazakhstan and the rest of the world and between mobile and other devices. These aggregate figures are shown openly to everyone on the Analytics page.'),
        '**How it is kept from becoming surveillance.** A visitor is counted under an anonymous identifier, derived from the address and the browser string with a key that is generated afresh every day; the address itself is never stored anywhere. Yesterday''s key cannot be recovered, so the same person on Monday and on Tuesday yields two unrelated values. Following anyone for longer than a day is not technically possible -- for us either.',
        '**How the anonymous part works.** A visitor who is not signed in is counted under an anonymous identifier, derived from the address and the browser string with a key that is generated afresh every day; the address itself is never stored anywhere. Yesterday''s key cannot be recovered, so the same signed-out visitor on Monday and on Tuesday yields two unrelated values.'),
        '**What we do not do in our own analytics.** We do not record your path through the site, build an interest profile, pass our own analytics data to third parties, or use it to select advertising for a particular person.',
        '**For registered users,** we retain the articles, courses and lessons viewed, reading time and exercise results. Administrators can use this history to analyze interests and learning progress; it is not shared with advertisers and is deleted with the account.'
    )
WHERE page_key = 'privacy' AND lang = 'en';

-- +goose Down
DROP TABLE IF EXISTS user_activity_events;
