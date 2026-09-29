-- +goose Up
-- The privacy page is editable and may have been saved by an administrator,
-- in which case the code-owned seed deliberately leaves it alone. Replace the
-- obsolete ZERO.kz paragraph in either kind of row so the deployed legal text
-- describes the JavaScript counter before the counter starts running.
UPDATE content_pages
SET body_md = replace(
        replace(body_md, '_Редакция от 25.07.2026._', '_Редакция от 29.09.2026._'),
        '**Чего мы не делаем.** Не записываем ваш путь по сайту, не строим профиль интересов, не передаём эти данные третьим лицам и не используем их для подбора рекламы под конкретного человека.' || E'\n\n' ||
        'В подвале сайта установлен счётчик посещаемости ZERO.kz: он считает просмотры страниц независимо от нас, чтобы цифры посещаемости, которые мы показываем рекламодателям, не были посчитаны только нами самими. Счётчик — это картинка: при открытии страницы ваш браузер запрашивает её у ZERO.kz и тем самым сообщает туда ваш IP-адрес и тип браузера. Код ZERO.kz на наших страницах не выполняется — политика безопасности этого не позволяет. Таргетированная реклама на основе расы, национальности, политических взглядов, биометрических данных или данных о здоровье не ведётся.',
        '**Чего мы не делаем в собственной аналитике.** Не записываем ваш путь по сайту, не строим профиль интересов, не передаём данные собственной аналитики третьим лицам и не используем их для подбора рекламы под конкретного человека.' || E'\n\n' ||
        'На публичных страницах работает официальный JavaScript-счётчик ZERO.kz: он измеряет посещаемость независимо от нас, чтобы цифры для читателей и рекламодателей можно было проверить во внешнем сервисе. При открытии страницы браузер передаёт ZERO.kz полный адрес и заголовок страницы, источник перехода, время посещения, IP-адрес, сведения о браузере, устройстве, размере экрана, языке и часовом поясе, а также анонимные идентификаторы посетителя и сеанса; ZERO.kz может сохранять свои cookie. Эти данные позволяют отдельно считать страницы, посетителей, визиты и источники переходов. Счётчик отключён в админ-панели, личном кабинете, редакторе, на страницах входа, регистрации, восстановления доступа и в служебных формах. ZERO.kz обрабатывает полученные данные как самостоятельный внешний сервис. Таргетированная реклама на основе расы, национальности, политических взглядов, биометрических данных или данных о здоровье не ведётся.'
    )
WHERE page_key = 'privacy' AND lang = 'ru';

UPDATE content_pages
SET body_md = replace(
        replace(body_md, '_25.07.2026 жағдайындағы редакция._', '_29.09.2026 жағдайындағы редакция._'),
        '**Не істемейміз.** Сайттағы жолыңызды жазбаймыз, қызығушылық профилін құрмаймыз, бұл деректерді үшінші тұлғаларға бермейміз және нақты адамға жарнама таңдау үшін пайдаланбаймыз.' || E'\n\n' ||
        'Сайттың төменгі бөлігінде ZERO.kz келушілер санағышы орнатылған: ол бет қаралымдарын бізден тәуелсіз санайды, сондықтан жарнама берушілерге көрсететін сандарымызды тек өзіміз санаған болып шықпаймыз. Санағыш — сурет: бет ашылғанда браузеріңіз оны ZERO.kz-тен сұрайды және сол арқылы IP-мекенжайыңыз бен браузер түріңізді хабарлайды. ZERO.kz коды біздің беттерімізде орындалмайды — қауіпсіздік саясаты оған жол бермейді. Нәсіл, ұлт, саяси көзқарас, биометриялық деректер немесе денсаулық туралы деректер негізінде мақсатты жарнама жүргізілмейді.',
        '**Өз аналитикамызда не істемейміз.** Сайттағы жолыңызды жазбаймыз, қызығушылық профилін құрмаймыз, өз аналитикамыздың деректерін үшінші тұлғаларға бермейміз және нақты адамға жарнама таңдау үшін пайдаланбаймыз.' || E'\n\n' ||
        'Жария беттерде ZERO.kz ресми JavaScript санағышы жұмыс істейді: ол оқырмандар мен жарнама берушілерге көрсетілетін сандарды сыртқы қызмет арқылы тексеруге мүмкіндік беру үшін сайтқа кіруді бізден тәуелсіз өлшейді. Бет ашылған кезде браузер ZERO.kz қызметіне беттің толық мекенжайы мен тақырыбын, ауысу көзін, кіру уақытын, IP-мекенжайды, браузер, құрылғы, экран өлшемі, тіл және сағат белдеуі туралы мәліметтерді, сондай-ақ келуші мен сеанстың анонимді идентификаторларын жібереді; ZERO.kz өз cookie файлдарын сақтай алады. Бұл деректер беттерді, келушілерді, кірулерді және ауысу көздерін бөлек санауға мүмкіндік береді. Санағыш әкімшілік панельде, жеке кабинетте, редакторда, кіру, тіркелу, қолжетімділікті қалпына келтіру беттерінде және қызметтік нысандарда өшірілген. ZERO.kz алынған деректерді дербес сыртқы қызмет ретінде өңдейді. Нәсіл, ұлт, саяси көзқарас, биометриялық деректер немесе денсаулық туралы деректер негізінде мақсатты жарнама жүргізілмейді.'
    )
WHERE page_key = 'privacy' AND lang = 'kz';

UPDATE content_pages
SET body_md = replace(
        replace(body_md, '_Revision of 25.07.2026._', '_Revision of 29.09.2026._'),
        '**What we do not do.** We do not record your path through the site, do not build an interest profile, do not pass this data to third parties, and do not use it to select advertising for a particular person.' || E'\n\n' ||
        'A ZERO.kz visitor counter runs in the site footer: it counts page views independently of us, so that the traffic figures we show advertisers are not figures only we have counted. The counter is an image: when a page opens, your browser requests it from ZERO.kz and thereby reports your IP address and browser type there. ZERO.kz code does not run on our pages — our security policy does not permit it. We do not conduct targeted advertising based on race, ethnicity, political views, biometric data, or health data.',
        '**What we do not do in our own analytics.** We do not record your path through the site, build an interest profile, pass our own analytics data to third parties, or use it to select advertising for a particular person.' || E'\n\n' ||
        'The official ZERO.kz JavaScript counter runs on public pages. It measures traffic independently of us so that readers and advertisers can verify our figures through an external service. When a page opens, the browser sends ZERO.kz the full page address and title, referrer, visit time, IP address, browser and device details, screen size, language and time zone, together with anonymous visitor and session identifiers; ZERO.kz may store its own cookies. These data allow pages, visitors, visits and traffic sources to be counted separately. The counter is disabled in the administration panel, account area, editor, sign-in, registration and account-recovery pages, and service forms. ZERO.kz processes the information it receives as an independent external service. We do not conduct targeted advertising based on race, ethnicity, political views, biometric data, or health data.'
    )
WHERE page_key = 'privacy' AND lang = 'en';

-- +goose Down
-- Legal disclosures are not rolled back: after the external counter has run,
-- restoring a statement that it never ran would make the policy inaccurate.
SELECT 1;
