# Changelog

All notable changes to this project are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.17.1] — 2026-08-25

### Changed

- Автоматический перевод включён. Он лежал выключенным, пока перевод был
  удобством: собственная модель автора справлялась быстрее, чем наши восемь
  минут и двадцать запросов. Теперь публикация требует недостающих языков сама,
  и автору без второй модели нужен путь — эти восемь минут покупают статью,
  которая иначе не вышла бы вовсе. В редакторе автору снова видна кнопка
  «Перевести на другие языки».

  Включено и для существующей установки (миграцией), и для новой (умолчанием в
  коде: свежая база сеет строку настроек сама и до умолчания колонки не
  доходит). Выключить по-прежнему можно в админке.

### Added

- Предупреждение о казахском переводе — рядом с кнопкой в редакторе и в
  руководстве на всех трёх языках. Машина переводит на казахский заметно слабее,
  чем на русский и английский: страдают термины, официальный слог и формы слов.
  Автора просят вычитать казахскую версию, а официальный, медицинский или
  юридический материал показать переводчику, редактору или специалисту в
  предмете. Ответственность за текст остаётся на авторе.

## [0.17.0] — 2026-08-25

### Added

- Проверка перед публикацией теперь ищет и языковые ошибки: орфографию,
  грамматику и пунктуацию. Каждая находка приходит с процитированным местом,
  чтобы автор правил конкретное слово, а не искал его по тексту. Эти правила
  предупреждают, но не блокируют: запятая — повод сказать автору, а не повод
  остановить материал.

  Проверяющему прямо запрещено спорить о стиле и о выборе слов, и запрещено
  трогать имена собственные, термины и цитаты. Статью пишет автор.

- Требование к языкам. Норма — три языка; материал для всей страны выходит на
  всех трёх, материал для конкретного места допускается на двух —
  государственном и русском. Основание — статья 21 закона «О языках в
  Республике Казахстан» от 11 июля 1997 года № 151-I, где объявления и другая
  визуальная информация прямо отнесены к двуязычным.

  Проверяется не моделью, а запросом к базе: версия либо есть и содержит текст,
  либо нет, и тут нечего решать. Правило стоит на единственном месте, через
  которое проходит любая публикация, — с включённой проверкой и без неё.

  Отказ называет, какой именно версии не хватает, и в кабинете, и в журнале
  модерации.


## [0.16.0] — 2026-08-25

### Changed

- Раздел «Как формируется курс тенге и инфляция» говорит прямо. Под каждым
  графиком теперь стоит вывод из его же чисел, а не одно описание того, что
  нарисовано.

  Главная врезка переписана на счёт, который не зависит ни от каких споров о
  том, что чему предшествует: производство в стране выросло в 3,6 раза,
  денежная масса — в 6 776 раз. Разница между этими двумя числами и есть
  инфляция. Ставка тут не причина и не лекарство, а цена, по которой
  Национальный банк отдаёт банкам деньги.

  Добавлено, кому это досталось: 70% денежной массы создали коммерческие банки,
  выдавая кредит под ставку Нацбанка, а заёмные деньги для производства 22 года
  из 32 стоили дороже, чем росли цены.

  Убрана фраза «какая из них тянет другую, годовой график не разбирает — да и
  незачем». Она поднимала вопрос и уходила от него, и читалась как оправдание.


## [0.15.4] — 2026-08-25

### Fixed

- Формула покрытия печаталась как 55,8 трлн ÷ 63,7 млрд = 875,52, а это деление
  даёт 875,98: результат считался по точным числам рядом с округлёнными. То же
  число в примечании к графику покрытия бралось из ряда и разошлось бы с
  формулой на полтенге. Теперь обе строки считают одинаково.


## [0.15.3] — 2026-08-25

### Fixed

- Обесценение было напечатано в обратную сторону: «тысяча тенге 1995 года
  покупает сегодня столько же, сколько 21 ₸» утверждало, что старые деньги
  стоили меньше новых. Верно обратное: тысяча, отложенная в 1995-м, покупает
  сегодня то, что тогда покупал двадцать один тенге. Колонка таблицы и
  примечание теперь говорят, в каких деньгах посчитано, — прежде это
  приходилось угадывать.

- Формула M3 = B × m не сходилась при перемножении: печаталось
  55,8 = 16,7 × 3,33, а это 55,61 — сто девяносто миллиардов тенге в зазоре
  между округлениями. Отношение теперь берётся между напечатанными числами, и
  тот же множитель стоит в примечании к графику денежной массы.

- Кратности роста измерялись не от той даты, которую называет их же подпись.
  «Приведены к ста в январе 1994… курс вырос в 97 раз» — но 97 это рост от
  ноября 1993-го, когда тенге было две недели. От января 1994-го рост в 43 раза.
  Теперь обе кратности считаются от общего месяца, к которому приведены линии.


## [0.15.2] — 2026-08-25

### Fixed

- График «Ставка Нацбанка и инфляция» показывал последовательность, которой в
  числах нет. Ставка бралась на конец года — снимок на 31 декабря, — а инфляция
  была итогом за весь год; от этого рождалось фантомное опережение, и читатель
  видел, будто банк поднимает ставку, а цены идут следом через год.

  2015-й это вскрыл: на конец года 16,0 против инфляции 6,7. Взвешенная по
  дням, что она реально стояла, ставка 2015-го — 8,65 против 6,68, а 2016-го —
  14,47 против 14,36. Теперь год сравнивается с тем же годом.

- Ряд ставки доводится до текущего года, а не обрывается на последнем
  заседании: год, в котором ставку не трогали, — тоже решение.

### Changed

- График развёрнут по горизонтали: тридцать два года в ширину колонки давали
  двадцать пикселей на год. Ось значений при прокрутке остаётся на месте.

- Врезка под графиком переписана. Прежняя утверждала, что ставка идёт следом за
  инфляцией; проверка по ряду этого не подтверждает, как не подтверждает и
  обратного. Годовой ряд различает только то, что линии движутся вместе.
  Вместо направления во врезке теперь то, что считается из тех же рядов:
  инфляция не опускалась до объявленной цели ни разу за 32 года, за последнее
  десятилетие она вдвое выше цели, и сколько осталось от тысячи тенге.


## [0.15.1] — 2026-08-25

### Added

- Календарь месяца на странице архива. Любой день выбирается сразу, а не
  перелистыванием по одному; выделены те дни, в которые что-то опубликовано, и
  только они кликаются. Месяц листается назад стрелками, вперёд — не дальше
  сегодняшнего: календарь дней, которых ещё не было, это прогулка по 404.

- Выбор места на странице прогноза: список всех городов и посёлков с прогнозом,
  сгруппированный по областям. Обычная форма, без единой строки скрипта.

- Карта осадков внизу страницы прогноза — можно двигать и рассматривать любой
  регион планеты. 16:9 на широком экране, 9:16 на вертикальном телефоне: карта —
  единственное на сайте, чему на узком экране высота нужна больше ширины.

  Карта своя, не встроенная чужая: политика безопасности запрещает вставлять
  чужие страницы рамкой, и обходить её ради погоды не стоило. Основа —
  OpenStreetMap, слой осадков — RainViewer, адрес свежего кадра узнаётся на
  сервере, потому что браузеру ходить наружу тоже запрещено. Ни одной сторонней
  строки кода на странице не выполняется.


## [0.15.0] — 2026-08-25

### Added

- Страница прогноза погоды у каждого места с известными координатами —
  без малого четыреста городов и посёлков, у каждого свой адрес вида
  **/weather/kachar**. Что сейчас, почасовой график на двое суток с подсказками
  при наведении, вероятность осадков и неделя вперёд с восходом и закатом.
  Погода в полосе наверху стала ссылкой целиком, вместе с иконками.

  Читателя, указавшего в профиле своё место, страница уводит на прогноз именно
  этого места: качарца в Качар, костанайца в Костанай. Место берётся из профиля
  и только оттуда — ни по адресу в сети, ни по номеру телефона мы его не
  вычисляем. У кого место не указано и у гостя открывается Алматы.

- Архив по дням: дата в полосе ведёт на страницу всего, что опубликовано в этот
  день, вида **/archive/2026-08-25**. Переходы внизу ведут к ближайшему дню, где
  что-то есть, а не к пустому вчерашнему числу. День без публикаций отдаётся,
  но в индекс не просится.

- Обе страницы попали в карту сайта: прогнозы по местам и те дни, на которых
  действительно что-то опубликовано.

### Fixed

- Дата в шапке показывала вчерашнее число каждый вечер после семи. Сервер
  работает по UTC, Казахстан на пять часов впереди, и полоса печатала ответ
  сервера. У сайта появились собственные часы; база часовых поясов вшита в
  двоичный файл, потому что образ distroless своей не имеет.

- Подпись оси на линейных графиках игнорировала формат ряда: температура шла
  без градусов. Теперь ось подписывается теми же единицами, что и подсказка.


## [0.14.0] — 2026-08-24

### Added

- Телеграм-бот с географией подписчика. Человек один раз говорит боту, где
  живёт, и получает материалы, написанные для его места, — те, что другим не
  уходят. Выбор места кнопками: сначала область, потом город или посёлок;
  можно взять область целиком. Команды **/start**, **/place**, **/stop**,
  меню на трёх языках.

  Правило доставки то же, что в ленте сайта: материал для области доходит до
  всех, кто внутри неё, включая города и посёлки; материал для посёлка — только
  до этого посёлка. Проверено тестом на живом справочнике мест.

  Место называет сам подписчик. Мы его не вычисляем — ни по номеру, ни по языку
  интерфейса — и никому не показываем.

### Changed

- В канал уходят только материалы для всей страны. Раньше туда падала каждая
  опубликованная статья, включая адресные: отключение воды в одном посёлке
  прилетало всем подписчикам по стране. Канал в Telegram шлёт одно сообщение
  всем одинаково и о читателе не знает ничего, поэтому показывать разным людям
  разное он не может — адресные материалы переехали в бота, который шлёт
  каждому отдельно.


## [0.13.1] — 2026-08-24

### Changed

- Прайс на рекламу пересчитан: дешевле рынка вдвое, а не на порядок. Карта
  0.13.0 была срезана слишком глубоко — скидка за географию накладывалась на
  базу, которая и без неё уже стояла ниже рынка, и две скидки перемножились.
  Костанайская область выходила в 5 000 ₸ против 54 000 ₸ у ng.kz: цена на
  порядок ниже читается не как выгодная, а как другой класс продукта.

  Теперь база подобрана так, чтобы **сопоставимый** таргетированный продукт
  стоил ровно половину конкурента, а цена за всю страну выводилась уже из неё.
  Верхний баннер на главной за 30 дней: вся страна 135 000 ₸, Алматы 44 145 ₸,
  Астана 32 926 ₸, Костанайская область 27 000 ₸, Костанай 15 538 ₸,
  Качар 2 889 ₸.

  Привязка проверена с двух концов лестницы: 27 000 ₸ за Костанайскую область —
  половина от 54 000 ₸, которые берёт ng.kz, чья аудитория этой областью и
  исчерпывается; 135 000 ₸ за всю страну — на 47% ниже примерно 257 000 ₸,
  во что обходится месяц баннера в рубрике на cifrum.kz.

- Минимальная цена заказа поднята с 1 000 до 2 000 ₸: на прежней карте она
  оказалась ниже, чем стоит самое мелкое размещение, и потому не срабатывала.


## [0.13.0] — 2026-08-24

### Added

- Реклама продаётся по географии, и цена следует за охватом. Рекламодатель
  выбирает место — страну, область, город или посёлок, — и платит по квадратному
  корню из доли населения, которую это место охватывает. Верхний баннер на
  главной за 30 дней: вся страна 25 000 ₸, Алматы 8 175 ₸, Костанайская область
  5 000 ₸, Костанай 2 877 ₸, Качар 1 000 ₸. Баннер идёт на страницах выбранного
  места и всех мест внутри него: взяли область — Костанай и Качар входят.

  Не прямая пропорция потому, что баннер на страницах места видят не только его
  жители, а проверка, постановка и показ заказа стоят одинаково при любом
  охвате; строгая пропорция дала бы Качару сорок тенге в месяц. Показатель и
  минимальная цена заказа — редактируемые тарифы.

- Население у мест, которые можно купить: страны из Всемирного банка, регионы
  из Викиданных, у каждого числа записан год переписи. Раньше население было
  только у 179 городов и посёлков — ни у стран, ни у областей, ни у Алматы с
  Астаной, то есть ровно у того, что рекламодатель купит в первую очередь.

### Changed

- Прайс на рекламу приведён к рынку и опущен ниже него. Что берут сегодня:
  zakon.kz — 300 000 ₸ за сутки ТГБ на полном трафике; informburo.kz — 1 650 ₸
  за тысячу показов и ещё 25% за таргетинг; cifrum.kz — 60 000 ₸ за семь дней
  баннера в рубрике; ng.kz в Костанае — 54 000 ₸ за месяц вверху главной при
  330 000 уникальных посетителей.

  Старая карта просила 180 000 ₸ за месяц верхнего баннера на нашей главной —
  втрое больше, чем берёт сильнейшая газета Костаная. Стало 25 000 ₸ за всю
  страну и около 5 000 ₸ за Костанайскую область: в их же области мы дешевле
  почти в одиннадцать раз. Тарифы, которые оператор менял сам, не тронуты.

- Поле «Регион» в заказе рекламы было свободной строкой, которая ни на что не
  влияла. Стало выбором из мест с известным населением, с пересчётом цены на
  месте и без слежки за читателем: таргетинг идёт по содержанию страницы.


## [0.12.0] — 2026-08-24

### Added

- Правки от читателей. Под каждой статьёй — строка «Если вы нашли ошибку или
  опечатку в тексте статьи, то сообщите нам об этом», ведущая на форму из трёх
  полей: глава, предложение, слово с ошибкой. Название статьи спрашивать не
  нужно — оно известно по адресу, а список глав подставляется из самой статьи.
  Заявку читает модель: механическую опечатку она правит сразу, спорную
  отклоняет с объяснением, и читатель видит, что стало с его правкой. Автор
  узнаёт о каждой применённой правке.

- Раздел «Курс и инфляция» показывает, кто назначает величины, из которых
  инфляция и курс выходят. Панель с главной страницы Нацбанка — базовая ставка,
  годовая инфляция, объявленная цель, TONIA, реальная ставка и отклонение от
  цели. Ставка Нацбанка с 1992 года против инфляции за те же годы: до сентября
  2015-го это ставка рефинансирования, дальше базовая. Пять формул с живыми
  числами — уравнение обмена, реальная ставка, покрытие резервами, расстояние
  до цели и разложение денежной массы на выпущенное Нацбанком и созданное
  банками, — и под каждым обозначением сказано, чему оно равно и из какой
  публикации взято.

- Денежная масса М3 показана вместе с денежной базой: расстояние между линиями
  и есть то, что банки построили кредитом поверх выпущенного Нацбанком.

- Подсказки на графиках: наведение показывает дату и значение каждого ряда,
  с клавиатуры — стрелками.

### Fixed

- Год, оканчивающийся нулём, читался из книг Нацбанка неверно: «10.10» (октябрь
  2010) приходило как «10.1» и разбиралось как октябрь 2001, затирая настоящие
  данные 2001 и 2002 годов. Три месяца 2001 года страница показывала
  двенадцатикратно завышенные резервы.

- Две линии на графике раскладывались по номеру точки, а не по дате: линия
  Национального фонда начиналась с 1993 года, хотя фонд создан в 2001-м.

- Подписи логарифмической оси ниже единицы печатались нулями.

### Removed

- Отдельный график инфляции по годам: тот же ряд стоит на графике ставки, где
  он сопоставлен с решениями Нацбанка.

- Подпись «скоро» у YouTube в подвале.


## [0.11.0] — 2026-08-13

### Added

- The prediction ledger at `/predictions`. Every forecast the site makes is
  recorded with the date it was made and later judged in public — hit, partial
  or miss — with the misses staying on the page. An analyst never seen to be
  wrong is not accurate, only unaudited; a score that can go down is the one
  claim to credibility that cannot be faked. Partial counts as half so the
  number flatters neither direction, open forecasts are excluded from it, and
  the admin flags anything past its own deadline, because never judging
  anything is how a ledger like this gets quietly gamed. An article that made
  forecasts now shows them underneath itself, with what became of each.

- `/llms.txt`: a map of the site for a model deciding what to read, listing
  only the human-written articles, with the citation format and a request not
  to quote the machine-written columns.

- Every citable article carries a finished reference with a copy button —
  author, title, Shanraq.org, date in words, absolute URL. Attribution follows
  the path of least resistance: a reader handed a ready string pastes it, and a
  pasted string is a link.

- IndexNow: the key is served at `/indexnow.txt` and publishing pings the
  endpoint, so a new article is announced instead of waited for. YandexBot
  fetched the key file within the day — verified by double reverse DNS — while
  Yandex Webmaster's own page checker was still answering "service
  unavailable".

- A Google News sitemap (`/sitemap-news.xml`), BreadcrumbList structured data
  on article pages, and `middleware.GetHead` so a HEAD probe is answered rather
  than met with 405.

- Conditional GET for HTML: pages carry an ETag and answer `If-None-Match`
  with 304. A crawler that visits ten times a day now downloads the page only
  when it changed.

- Static asset URLs carry a content hash (`?v=…`). Assets are served with a
  day of cache, so before this a deploy could land new markup on a returning
  visitor's old stylesheet — which is exactly what happened once, and read as a
  rendering bug rather than a caching one.

- An account register in the admin panel: who is registered, with the means to
  act on it. Analytics gained per-brand marks, a scale on the trend chart, and
  tooltips in place of the paragraphs that used to sit under each card.

- The newsletter is back, under the social buttons rather than as a card of its
  own, now with double opt-in and one-click unsubscribe (RFC 8058).

- Listings can be edited and deleted by their owner, from the cabinet and from
  the listing page itself. Until now a mistake in a published advert could only
  be waited out: there was no edit and no delete, only "extend" and the paid
  promotion buttons. Editing runs the same validation as posting, so a listing
  cannot be edited into a state it could not have been posted in; the id, the
  expiry clock and any paid promotion are deliberately left alone.

- The price now carries a visible currency selector, set automatically from the
  country and changeable by hand. The location became a required field, because
  it is what makes the currency knowable at all.

- Readers can delete their own comments. A comment is published under a real
  name; the only previous way to take one back was to ask an administrator.

- 314 intra-city districts added to the location reference — see Fixed.

- A "Delete" button in the author studio, for drafts. A false start or a
  duplicate left behind by a lost session had no way out of the table at all.
  Published articles deliberately do not offer it: they carry views, votes,
  comments, an RSS entry and a Telegram post, and somewhere a reader may already
  hold the link — so the way out is one step longer, unpublish then delete, which
  is the pause such a decision deserves. The action is confirmed in the browser,
  scoped to the owner in SQL, and runs in a transaction that also sweeps the two
  tables keyed by article id without a foreign key (the reading-depth funnel and
  favourites). The moderation ledger is left intact: it is an append-only record
  and an author must not be able to erase it.

- Agent profiles now carry a kind — private realtor, agency or developer — plus
  a БИН for the two company kinds. A construction company previously had nowhere
  to register: it had to pose as an agent and type its name into a free-text
  "Agency" box. All three share one table, one moderation queue and one public
  page, because all three need the same things from us; only the label, the
  required fields and the badge wording differ. The moderation queue shows the
  kind and БИН, so a moderator can see whether they are approving a person
  vouching for themselves or a brand that must match the public register.

- A "what you can do once registered" block under the signup form, naming the
  four roles and where each is switched on. The registration page previously
  said nothing about who the site is for, so an agency landing on it had no way
  to discover the agent cabinet exists.

- Field-level hints across registration and the agent cabinet, and a guide
  section ("One account, several roles") in all three languages.

- A "follow us" card in the article aside, under the table of contents. The
  aside is sticky on desktop, so the ask now travels with the reader for the
  whole piece instead of waiting at the bottom, where only those who finish ever
  see it. Only configured profiles are listed — a `#` placeholder (YouTube today)
  would be a dead click on a card whose single job is to be clicked. Clicks
  report an aggregate `follow_social` event to the existing counter, so the card
  can be judged on numbers rather than on taste.

- Article cards now carry the view counter, which until now only existed on the
  article page itself: `👍 1 · 👁 115 · D. Baimurza · 04.08.26`. Counters above
  999 collapse to a localized short form — `1,2 мың` / `1,2 тыс.` / `1.2k`,
  truncated rather than rounded so a card never claims more views than there
  were.

### Changed

- The 90 machine-written columns, already out of the search index, are now also
  out of reach of AI crawlers: a named group in `robots.txt` built from the same
  `indexable` flag, plus `X-Robots-Tag: noindex, follow, noai, noimageai`.
  `noindex` is a search directive and nothing more, so the assistants had gone
  on reading them and could have quoted machine opinion back with this site's
  name attached — in a medium where nothing can be retracted. The badge that
  warns a human reader does not survive into a retrieved chunk. Reversing it is
  one boolean, the same one that governs search.

- Analytics rows carry the companies' own logos, taken from the CC0 Simple
  Icons set with each brand's official guidelines cited above the path.
  Vendored rather than linked, because the CSP blocks a CDN and a logo that
  fails to load is worse than none. Yandex, Bing, Edge and Windows keep drawn
  marks: their owners had them pulled from every open set.

- Admin sections are ordered to match the sidebar, which had been claiming an
  order the page did not follow, and the scroll-spy no longer jumps.

- The rating row is the same shape for everyone. The author of a piece, and
  every signed-out reader, used to get a sentence where the buttons go, which
  wrapped the row onto two lines; the explanation now appears on the click.

- Covers are sized by aspect ratio rather than a fixed pixel height, and the
  space is reserved before the image arrives so the heading no longer jumps.

- Registration deliberately stays a single form rather than gaining per-role
  tabs. A role is an attribute of an account, not a different kind of account:
  tabs would force a choice before anyone has seen the site, would strand a
  reader who later decides to sell, and could not finish the job anyway, since
  an agent still needs a human moderator afterwards.

- The home sidebar's email newsletter block is replaced by the same follow card
  (accent variant, keeping that column's one point of colour). The form was not
  broken — addresses were stored and the weekly digest ran — but it never sent a
  confirmation, so subscribing felt like nothing happened, and in a week it
  collected no one. The `/subscribe` route, the `subscribers` table and the
  digest job are deliberately left intact: this is a UI decision, reversible by
  restoring the block.

- Bylines on cards are abbreviated to initial + family name ("Daulet Baimurza"
  → "D. Baimurza"). The given name is the part that gets cut — the byline
  convention everywhere space is tight — freeing the width the counter needed.

- The card footer (rating, views, byline, date) now comes from one shared
  `post_meta` partial instead of three near-copies in home / favourites /
  author templates, which had already drifted (favourites still drew the rating
  as a bare "▲").

### Fixed

- `SHANRAQ_SYNDICATE_INDEXNOW_KEY` was ignored: viper's `AutomaticEnv` only
  binds keys it already knows, and the key had no default registered. The same
  trap had been documented for the metrics token and stepped in again anyway.

- `/favicon.ico` redirected to the SVG, which several crawlers do not follow.

- Googlebot was spending its budget on `/jobs` and `/admin`, which answer 401
  and 303. On a site it visits about ten times a day that is worth reclaiming.

- The Kazakh and English versions of every article were never actually
  submitted for crawling: the sitemap emitted only the Russian `<loc>`, which on
  a trilingual site is two thirds of the catalogue relying on being noticed
  sideways.

- Crawlers were being counted as readers, and Google Search traffic could not
  be told apart from Gmail.

- Saved listings never appeared under Favourites. The query selected the
  agent-badge columns without joining the agent table, so it failed on every
  call; the handler logged the error and rendered an empty page, which reads as
  "the save button does nothing" although the bookmark was stored correctly.

- The listings map was empty. Three causes stacked: the pin query still named a
  renamed column and answered 500; no district in the reference carries
  coordinates, so a flat in Медеу or Петродворцовый had nothing to plot; and no
  Russian settlement had coordinates at all — 0 of 159, against 182 of 187 in
  Kazakhstan. Pins now climb to the nearest ancestor with a position, and the
  159 Russian city centres were added.

- A listing with a street and house number now geocodes on save and gets an
  exact marker. Plotting one previously meant opening a collapsed panel and
  dragging a pin — asking a second time for an address already typed, which is
  why a site full of complete addresses had an empty map. Best-effort; a miss
  falls back to the settlement centre and never blocks the save.

- Map popups print the listing's own currency instead of a hardcoded tenge, and
  open on hover rather than costing a click each.

- Address rows named the wrong things: a Saint Petersburg listing published as
  "Область: Санкт-Петербург, Город: Петродворцовый" — an oblast that does not
  exist and a district called a city. The fields were filled by a location
  node's depth in the tree, which assumes every country nests the same number of
  times; a federal city sits exactly where an oblast sits, so Moscow, Almaty,
  Astana and Shymkent were wrong the same way. Fields now follow the node's kind,
  the `village` column is renamed `district`, and existing listings were
  re-derived. Nothing here was ever typed by hand — country, region, city and
  district come only from the reference; the free-text fields are microdistrict,
  street and house, which no reference can enumerate.

- Seven of ten uploaded photos disappeared: a modern phone shoots 10–15 MB per
  frame against a 10 MB limit, and all ten uploads were fired at once to share
  one mobile uplink. Photos are now shrunk to 2000px in the browser before
  sending (the server stored nothing larger anyway), uploaded one at a time, and
  the limit is raised to 25 MB. Documents and PDFs are sent untouched — the
  re-encode would damage a scan of a title deed. Failures are now reported: the
  handler ended in an empty `catch`, so seven of them produced no message at all.

- Listing photos could not be swiped. The gallery had dots and a six-second
  timer with no way to stop it, so choosing a photo lasted six seconds. Swipe
  added, and any manual choice now ends the auto-advance.

- A price entered against a Russian address was stored in tenge. The currency
  was derived correctly, but the submission form dropped the chosen location
  whenever it bounced off validation: the hidden geo field was rendered without
  its value and the cascade rebuilt itself empty. The author fixed whatever was
  flagged, submitted again with no country attached, and the price fell back to
  the default currency. With the country gone the form also began demanding a
  Kazakh title that a Russian listing does not need, making a second bounce
  likely. The form now carries the location back, and the location is required.

- Price filters compared amounts across currencies: a listing at 90 000 ₽ and
  one at 90 000 ₸ matched the same range, five times apart in real money. A
  price range is now scoped to one currency, following the country filtered on.

- The location reference was missing most intra-city districts. Saint Petersburg
  had nine of eighteen (Петродворцовый among the absent), Moscow ten of twelve,
  and every Russian city outside the two federal ones had none at all, so a flat
  in Novosibirsk could be placed no more precisely than "Novosibirsk". Added for
  the cities that have official divisions — inventing districts for a city that
  has none would be worse than leaving it blank. Karaganda and Aktobe gained
  theirs too; Kazakhstan's million-plus cities were already complete.

- `stripMD` mangled HTML instead of removing it: `<script>alert(1)</script>`
  became `<scriptalert(1)</script`, which then went into card excerpts and the
  page's meta description. Never an XSS — the templates escape it — but rubbish
  where a summary belongs, and it took a test to notice.

### Security

- A signed-out account's tokens stopped being honoured. Sessions are stateless
  JWTs with no session table, so until now a token stayed valid for its full
  lifetime after a role change, a password reset or a deletion — and production
  was issuing two-hour tokens. Accounts carry an `auth_version` that a demotion
  or a reset increments, and a token that cannot be confirmed against it is
  refused rather than accepted.

- The last administrator can no longer be demoted or deleted by two concurrent
  requests. The check and the write were separate statements, so two demotions
  arriving together each saw two admins and both went through, locking everyone
  out. Both paths now take a Postgres advisory lock and count inside the
  transaction. Proved against the old code first: 0 of 2 demotions were refused,
  where exactly 1 should have been.

- The JSON signup endpoint now honours the registration service flag. The
  browser form has always checked it — offering an invite code while the beta is
  closed, demanding a real first and last name, recording the consent — while
  `POST /auth/signup` checked none of that and returned tokens. Closing
  registration in the admin panel therefore shut the site to visitors and left
  it open to anyone who could spell JSON. The auth module gained an optional
  signup gate, wired in `main.go` to the same flag; the API has no invite-code
  field, so anything short of a fully open registration is refused there.

- The jobs API is staff-only again. It was reachable by role `user`, and any
  registered reader could mint an API key, so a reader could enqueue arbitrary
  jobs: `ai_translate` rewrites the translations of any article id it is handed
  and spends the AI budget doing it, `syndicate_telegram` re-posts to the
  channel. Enqueue now requires operator or admin.

Both were found by an external audit; both are covered by regression tests.

### Tests

- The integration suite had been skipping silently for lack of
  `SHANRAQ_TEST_DB` — the new tests were passing against deliberately broken
  code, which is how it was noticed. Four of the AI-channel tests were then run
  against the un-fixed code to confirm they fail.

- Security invariants are now pinned by tests, chosen for what a regression
  would cost rather than for the coverage figure. The headline one: markdown
  from any author or commenter must never render a live tag. goldmark runs
  without `WithUnsafe`, but that flag is one word long and looks harmless —
  somebody will eventually reach for it to embed a video, and stored XSS ships
  the same day.
  Also covered: email normalisation as account identity (two spellings must not
  become two accounts, two mailboxes must not become one), the password floor,
  real-name validation, that a commenter's label can never fall back to their
  e-mail address, contact masking, the staged-launch gates and the global site
  switch failing open, category slugs being a closed set, slugs staying URL-safe,
  and that a 5xx never carries its cause to the client while still reaching the
  log.

## [0.10.2] — 2026-08-04

### Changed
- The team's own browsing no longer pollutes analytics. Staff (admin/operator)
  visits were already skipped, but the owner's publish routine — logging into the
  admin account, then a test account to like the new article — still generated
  guest page views and login clicks while logged out. A device is now flagged
  with a persistent opt-out cookie the moment it loads a page as staff (or as a
  configured `SHANRAQ_ANALYTICS_EXCLUDE_EMAILS` address), so it stays out of the
  counts even after logging out or switching accounts. Click events honour the
  same rule.

### Added
- A "Data center / VPN — reading language" panel that splits the masked VPN
  bucket by language, because that split identifies who is behind the VPN:
  Russian ≈ Russia/CIS, while English ≈ genuine international readers (China,
  Iran, the West) who bridge through English rather than Russian.

### Changed
- Reworded the "English readers — origin" note: English from the Data center /
  VPN bucket is now described as *most likely genuine foreigners* (China, Iran,
  the West), not vaguely "masked" — they read English precisely because they
  don't read Russian.

## [0.10.0] — 2026-08-03

### Added
- Audience analytics now tracks the **reading language** (Kazakh / Russian /
  English) and crosses it with visitor origin, to answer whether the English
  audience is genuine foreign readers or something else. Two new 30-day panels:
  "Reading language" (the overall kk/ru/en mix) and "English readers — origin"
  (where the English-version visits come from). The origin panel is the sharp
  signal a VPN cannot mask: English from a real foreign country is a genuine
  foreigner; English from Kazakhstan is a curious local; English from
  "Data center / VPN" is masked traffic — e.g. a reader on a VPN from a country
  with a restricted internet (Russia, China, Iran). Aggregate-only as always.

## [0.9.0] — 2026-08-03

### Added
- The country panel now separates **hosting/cloud/VPN traffic from real
  readers**. An optional ASN database (DB-IP ASN Lite) flags visits from cloud
  networks (AWS, Google, Azure, Cloudflare, OVH, Hetzner, VPN exits, …) and
  buckets them as "ЦОД / VPN" instead of a country — so the geographic rows
  reflect actual eyeballs, not US-hosted infrastructure that otherwise dominates
  by IP. Only the coarse label is counted; the IP is still discarded. Optional
  and graceful: without the ASN file every IP is bucketed by country as before.

### Changed
- Bot detection now also recognises common scraper/library agents (Scrapy,
  axios, node-fetch, Guzzle, aiohttp, Postman, …), so fewer automated hits leak
  into the human audience counts.

## [0.8.0] — 2026-08-03

### Added
- Audience analytics now shows a **visitor-country** panel (30-day) so the team
  can tell domestic (Kazakhstan) readers from a genuine foreign audience — the
  English content shared on LinkedIn draws both, and by-language view counts
  alone couldn't say which. Countries are resolved from the visitor IP with a
  local DB-IP Lite database (bind-mounted read-only, refreshed monthly on the
  host); only the coarse country code is counted and the IP is immediately
  discarded, exactly like the User-Agent in the device/OS panels — no
  per-visitor profiling. The feature is optional: with no database present the
  panel simply stays empty and nothing else changes.

## [0.7.2] — 2026-08-03

### Changed
- The author phone-verification form now states the one-time code arrives **via
  Telegram** (with a nudge to have Telegram on that number), so a user without
  Telegram isn't sent down a dead end. Labels updated from "SMS" to "Telegram".

## [0.7.1] — 2026-08-03

### Added
- SMS gateway can now deliver verification codes over **Telegram** (SMSC `tg=1`)
  via a new `SHANRAQ_SMS_CHANNEL=telegram` setting — no paid operator sender
  name, far cheaper than SMS, and it works immediately once the SMSC account is
  active.

## [0.7.0] — 2026-08-03

### Added
- Audience analytics now shows the **device, OS and browser mix** (mobile /
  tablet / desktop; Android / iOS / Windows / macOS / Linux; Chrome / Safari /
  Firefox / Edge / …) as aggregate 30-day panels — coarse, no per-visitor
  profiling. The User-Agent is classified and discarded, like the bot check.

### Changed
- The team's own visits no longer count as audience: page views from a logged-in
  admin/operator are skipped, so the owner's constant browsing stops inflating
  the guest and "Direct" counts.

## [0.6.5] — 2026-08-02

### Changed
- Dependency maintenance (Dependabot): bumped the Go dependency group (pgx,
  viper, zap, prometheus client, jwt, golang.org/x/crypto, goose, validator,
  anthropic-sdk and others — all minor/patch) and CI GitHub Actions
  (checkout@v7, setup-go@v7, docker/setup-buildx-action@v4). Build and full
  test suite green.

## [0.6.4] — 2026-08-02

### Changed
- Commercial SEO scanners (already turned away in robots.txt) are no longer
  recorded in analytics at all — they neither count as guests nor clutter the
  bot panel. Existing historical rows are cleared. Yandex source detection now
  also catches the short domain `ya.ru`.

## [0.6.3] — 2026-08-02

### Changed
- Telegram auto-posts now tag the article link with `?utm_source=telegram`, so
  visits from the channel are attributed to Telegram in analytics instead of
  falling into "Direct" (the messenger strips the referrer). No manual work per
  post.

## [0.6.2] — 2026-08-02

### Added
- UTM attribution: traffic analytics now reads `?utm_source=` first (mapped to a
  known source label) and falls back to the Referer host. Links shared with
  `?utm_source=telegram` / `instagram` / etc. are attributed correctly even when
  the browser strips the referrer (messengers, in-app browsers) — no third-party
  service required.

### Changed
- robots.txt now turns away commercial SEO scanners (Ahrefs, Semrush, MJ12, Dot,
  BLEX, DataForSeo, Petal, MegaIndex) — the heaviest crawlers that return no
  value. Search engines and AI crawlers stay allowed (AI-answer discovery is a
  channel worth keeping).

## [0.6.1] — 2026-08-02

### Fixed
- Field-help tooltip was unreadable on the dark theme (white text on the light
  tooltip background). Dark theme now uses dark text; light theme keeps its dark
  tooltip with white text.
- Audience tiles used calendar windows, so at the start of a month the "week"
  count could exceed the "month" count (the current week reached back into the
  previous month). The tiles now use rolling windows — today, last 7, 30 and 365
  days — so they always nest: today ≤ week ≤ month ≤ year.

## [0.6.0] — 2026-08-02

### Added
- Address ↔ map sync on the listing form. A new "Find on the map by address"
  button geocodes the entered address (cascade + street/house) and drops a
  precise pin; placing or dragging the pin reverse-geocodes it and fills in the
  street, house number and microdistrict — filling only what the geocoder
  returns, never clearing what the author typed. Geocoding is proxied
  server-side over OpenStreetMap Nominatim (the browser CSP blocks direct calls)
  with an in-memory cache; the country/region/city cascade stays manual by
  design (it is bound to the platform's own geo database).

## [0.5.1] — 2026-08-02

### Changed
- Listing form now reflects the currency before submit: picking Russia in the
  location cascade flips the price hint to rubles and shows a ₽ chip next to the
  price, while Kazakhstan shows ₸. The geo API now returns each node's country so
  the form can react to it. (The stored currency was already correct in 0.5.0;
  this removes the confusing "in tenge" hint for Russian listings.)

## [0.5.0] — 2026-08-02

### Added
- Listings for property in Russia, not just Kazakhstan. The location cascade
  already carried both countries; now each listing also carries its own currency
  — tenge (₸) for Kazakh addresses, ruble (₽) for Russian ones — chosen
  automatically from the selected location and shown on cards, the listing page,
  "my listings", favorites and JSON-LD. Posting needs only a verified email (no
  phone/SMS), so Russian users can list right away.

### Changed
- The mandatory Kazakh title is waived for Russian listings: when the location is
  in Russia, only Russian and English titles are required (Kazakh optional). The
  price-field label no longer bakes in ₸, since the currency now follows country.

## [0.4.1] — 2026-08-01

### Fixed
- Info-bar exchange rates could stay blank for up to 6 hours after a single
  transient fetch failure at boot. The National Bank fetch now retries with a
  short backoff, the HTTP timeout is more forgiving (10s), and empty rates are
  re-fetched on the 30-minute tick instead of only every 6 hours.

## [0.4.0] — 2026-08-01

### Added
- SMS phone verification for agents/authors: a provider-agnostic SMS gateway
  (`pkg/modules/sms`) with Mobizon.kz and SMSC.kz backends, chosen by config
  (`SHANRAQ_SMS_PROVIDER` + credentials). The verification flow already existed
  (code mint, hash, rate-limit, confirm); this wires real delivery. With no
  provider set, SMS stays off and codes are dev-logged — switching provider is a
  single environment variable, so onboarding friction with one aggregator never
  blocks the platform.

## [0.3.1] — 2026-08-01

### Changed
- Listing form: each language tab's placeholders now read in that tab's own
  language (Russian tab in Russian, Kazakh in Kazakh, English in English),
  regardless of the interface language — so the hint matches the content asked for.

## [0.3.0] — 2026-08-01

### Added
- Trilingual listings (KZ/RU/EN) — the flagship, mandatory feature: title and
  description in all three languages via a tabbed form; each reader sees the
  listing in their language with fallback. Script sanity on submit — English
  must be Latin, Russian Cyrillic, Kazakh either (Cyrillic or the new Latin).

## [0.2.0] — 2026-08-01

### Added
- Listing documents: agents can attach PDF plans/passports/contracts and image
  schemes to a listing (`/media/upload-doc`), shown in a "Documents / floor plan"
  block on the listing page.
- Traffic-source analytics (referrer → Google/Yandex/Telegram/Facebook/direct/…).
- Separate bot vs human classification so audience counts reflect real people.
- Dedicated real-estate sitemap (`/sitemap-listings.xml`) for Search Console.
- Branded PNG Open Graph card and descriptive homepage title for link previews.
- Brand watermark overlay on article covers and listing photos.
- GitHub icon in the footer; Telegram and Facebook social links from config.

### Changed
- Telegram bot token / chat id now bind from environment (secret never in config).

## [0.1.0] — 2026-07-30

First tagged release. Live in closed beta at [shanraq.org](https://shanraq.org).

### Added
- Trilingual (kk/ru/en) publishing with per-language SEO (hreflang, sitemaps, JSON-LD).
- Real-estate classifieds with photos, geo, amenities, and promotion/feature tariffs.
- Block-resilient syndication: always-on RSS and automatic Telegram posting on publish.
- Optional Claude-powered AI co-editor and trilingual auto-translation (off by default).
- Media pipeline: upload, EXIF strip, brand watermark; pluggable storage backend.
- Ratings with weighted author karma and anti-brigading.
- Referral loop: invite links, attribution, promotion-credit rewards.
- Operator admin panel: editable legal/info pages, tariffs, service flags, payment
  provider, and AI settings — no redeploy required.
- Aggregate-only, privacy-respecting audience analytics (no per-visitor profiling).
- Secure auth: refresh-token rotation, RBAC, password-reset flows, CSRF protection.
- Production stack: Docker Compose + Caddy automatic HTTPS, embedded Goose migrations.

[Unreleased]: https://github.com/DauletBai/shanraq.org/compare/v0.17.1...HEAD
[0.17.1]: https://github.com/DauletBai/shanraq.org/compare/v0.17.0...v0.17.1
[0.17.0]: https://github.com/DauletBai/shanraq.org/compare/v0.16.0...v0.17.0
[0.16.0]: https://github.com/DauletBai/shanraq.org/compare/v0.15.4...v0.16.0
[0.15.4]: https://github.com/DauletBai/shanraq.org/compare/v0.15.3...v0.15.4
[0.15.3]: https://github.com/DauletBai/shanraq.org/compare/v0.15.2...v0.15.3
[0.15.2]: https://github.com/DauletBai/shanraq.org/compare/v0.15.1...v0.15.2
[0.15.1]: https://github.com/DauletBai/shanraq.org/compare/v0.15.0...v0.15.1
[0.15.0]: https://github.com/DauletBai/shanraq.org/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/DauletBai/shanraq.org/compare/v0.13.1...v0.14.0
[0.13.1]: https://github.com/DauletBai/shanraq.org/compare/v0.13.0...v0.13.1
[0.13.0]: https://github.com/DauletBai/shanraq.org/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/DauletBai/shanraq.org/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/DauletBai/shanraq.org/compare/v0.10.2...v0.11.0
[0.10.2]: https://github.com/DauletBai/shanraq.org/compare/v0.10.0...v0.10.2
[0.10.0]: https://github.com/DauletBai/shanraq.org/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/DauletBai/shanraq.org/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/DauletBai/shanraq.org/compare/v0.7.2...v0.8.0
[0.7.2]: https://github.com/DauletBai/shanraq.org/compare/v0.7.1...v0.7.2
[0.7.1]: https://github.com/DauletBai/shanraq.org/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/DauletBai/shanraq.org/compare/v0.6.5...v0.7.0
[0.6.5]: https://github.com/DauletBai/shanraq.org/compare/v0.6.4...v0.6.5
[0.6.4]: https://github.com/DauletBai/shanraq.org/compare/v0.6.3...v0.6.4
[0.6.3]: https://github.com/DauletBai/shanraq.org/compare/v0.6.2...v0.6.3
[0.6.2]: https://github.com/DauletBai/shanraq.org/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/DauletBai/shanraq.org/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/DauletBai/shanraq.org/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/DauletBai/shanraq.org/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/DauletBai/shanraq.org/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/DauletBai/shanraq.org/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/DauletBai/shanraq.org/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/DauletBai/shanraq.org/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/DauletBai/shanraq.org/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/DauletBai/shanraq.org/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/DauletBai/shanraq.org/releases/tag/v0.1.0
