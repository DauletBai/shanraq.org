package articles

// staticPageContent is one localized info page (title + Markdown body).
type staticPageContent struct {
	Title string
	Body  string
}

// staticPages holds the About / Guide / Support / Pricing pages in every UI language.
var staticPages = map[string]map[string]staticPageContent{
	"about": {
		"ru": {Title: "О нас", Body: `Shanraq — независимая площадка, где любой человек может писать и публиковать статьи на казахском, русском и английском языках. Мы верим, что доступ к информации и право высказываться — основа развития общества.

## Наша миссия
Дать голос не только профессиональным журналистам, но и обычным людям: врачам, учителям, предпринимателям, студентам. ИИ-помощник помогает превратить мысль в понятную и профессиональную статью, а автоматический перевод делает её доступной на трёх языках.

## Наши принципы
- **Проверяемость.** Мы за факты, которые можно проверить, а не за громкие лозунги.
- **Уважение.** Профессиональный тон вместо оскорблений.
- **Открытость.** Прозрачные правила и устойчивость к произвольным блокировкам.

Shanraq — это дом, где сходятся свободные голоса.`},
		"kz": {Title: "Біз туралы", Body: `Shanraq — кез келген адам қазақ, орыс және ағылшын тілдерінде мақала жазып, жариялай алатын тәуелсіз алаң. Ақпаратқа қолжетімділік пен пікір білдіру құқығы қоғам дамуының негізі деп сенеміз.

## Біздің миссиямыз
Кәсіби журналистерге ғана емес, қарапайым адамдарға — дәрігерлерге, мұғалімдерге, кәсіпкерлерге, студенттерге — дауыс беру. ИИ көмекші ойды түсінікті әрі кәсіби мақалаға айналдыруға көмектеседі, ал автоматты аударма оны үш тілде қолжетімді етеді.

## Біздің ұстанымдарымыз
- **Тексерілетіндік.** Айқайлаған ұрандар емес, тексеруге болатын фактілер.
- **Құрмет.** Балағаттың орнына кәсіби тон.
- **Ашықтық.** Мөлдір ережелер және еркін блоктауға төзімділік.

Shanraq — еркін дауыстар тоғысатын үй.`},
		"en": {Title: "About us", Body: `Shanraq is an independent platform where anyone can write and publish articles in Kazakh, Russian, and English. We believe that access to information and the right to speak are the foundation of a thriving society.

## Our mission
To give a voice not only to professional journalists but to ordinary people — doctors, teachers, entrepreneurs, students. An AI assistant helps turn a thought into a clear, professional article, and automatic translation makes it available in three languages.

## Our principles
- **Verifiability.** We stand for facts that can be checked, not loud slogans.
- **Respect.** A professional tone instead of insults.
- **Openness.** Transparent rules and resilience to arbitrary blocking.

Shanraq is a home where free voices meet.`},
	},
	"guide": {
		"ru": {Title: "Как публиковать статьи и размещать объявления", Body: `Здесь описано, как пользоваться платформой, что бесплатно, что платно, а также ваши права и ответственность.

## Как написать статью
1. Зарегистрируйтесь и войдите в **Студию автора**.
2. Нажмите «Написать» и создайте статью на одном из трёх языков.
3. Выберите рубрику и подрубрику, при желании добавьте обложку.
4. Используйте ИИ-помощник, чтобы улучшить текст, и автоперевод, чтобы опубликовать статью сразу на казахском, русском и английском.
5. Сохраните и нажмите «Опубликовать».

## Что бесплатно, а что платно
**Бесплатно и навсегда:** чтение всех материалов, подписка на рассылку, регистрация, ручная публикация статей (в том числе на трёх языках, если вы переводите сами), оценки, комментарии и карма.

**Платные услуги** — это дополнительные возможности на базе ИИ и продвижение: улучшение статьи ИИ-редактором, автоматический перевод на три языка, генерация обложки, а также размещение и продвижение объявлений о недвижимости. Актуальные цены — на странице [Тарифы](/pricing).

## Недвижимость и объявления
Наш символ — шаңырақ, а значит, дом. В разделе **«Недвижимость»** можно подать объявление о продаже или аренде: квартиры, дома, земельного участка или коммерческого помещения. Нажмите «Подать объявление», заполните поля (страна, область, город, село, тип, цена, площадь и описание) — и объявление появится в разделе.

Кабинет рекламодателя для баннерной рекламы появится позже: доход будет справедливо делиться между автором материала, платформой и фондом сообщества.

## Права и ответственность
- Вы сохраняете авторские права на свои материалы и отвечаете за их содержание.
- **Публикуйте только достоверную информацию.** За клевету, заведомо ложную информацию, разжигание межнациональной розни, вражды и другой противоправный контент автор несёт уголовную ответственность по законам Республики Казахстан.
- **Оскорбительные и незаконные материалы, неуважительный тон** модерируются нашими читателями (пользователями) и нашей моделью ИИ; такие материалы понижаются в выдаче, скрываются или удаляются.
- Не публикуйте чужие тексты и изображения без разрешения, не нарушайте авторские права.
- Уважайте частную жизнь: не раскрывайте персональные данные других людей без их согласия.
- Запрещены спам, мошенничество, реклама запрещённых товаров и услуг.
- Материалы 18+ и любой контент, вредящий детям, недопустимы.
- При нарушении правил мы вправе скрыть материал, ограничить или заблокировать аккаунт. Решение можно оспорить через раздел [Поддержка](/support).
- ИИ-инструменты помогают писать в уважительном тоне и подсказывают, где нужен источник, но окончательная ответственность за текст — на авторе.

Публикуя материал, вы соглашаетесь с этими правилами.`},
		"kz": {Title: "Мақала мен хабарландыруды қалай жариялау керек", Body: `Мұнда платформаны қалай пайдалану керектігі, не тегін, не ақылы екені, сондай-ақ сіздің құқықтарыңыз бен жауапкершілігіңіз сипатталған.

## Мақаланы қалай жазу керек
1. Тіркеліп, **Автор студиясына** кіріңіз.
2. «Жазу» түймесін басып, үш тілдің бірінде мақала жасаңыз.
3. Айдар мен ішкі айдарды таңдап, қаласаңыз мұқаба қосыңыз.
4. Мәтінді жақсарту үшін ИИ көмекшіні, мақаланы бірден қазақ, орыс, ағылшын тілдерінде жариялау үшін автоаударманы қолданыңыз.
5. Сақтап, «Жариялау» түймесін басыңыз.

## Не тегін, не ақылы
**Тегін әрі әрқашан:** барлық материалды оқу, жаңалық таратылымына жазылу, тіркелу, мақаланы қолмен жариялау (өзіңіз аударсаңыз, үш тілде де), бағалау, пікір және карма.

**Ақылы қызметтер** — ИИ негізіндегі қосымша мүмкіндіктер мен жылжыту: мәтінді ИИ-редактормен жақсарту, үш тілге автоаудару, мұқаба жасау, сондай-ақ жылжымайтын мүлік хабарландыруын орналастыру мен жылжыту. Ағымдағы бағалар — [Тарифтер](/pricing) бетінде.

## Жылжымайтын мүлік және хабарландырулар
Біздің нышанымыз — шаңырақ, яғни үй. **«Жылжымайтын мүлік»** бөлімінде сату немесе жалға беру туралы хабарландыру беруге болады: пәтер, үй, жер учаскесі немесе коммерциялық орын. «Хабарландыру беру» түймесін басып, өрістерді толтырыңыз (ел, облыс, қала, ауыл, түрі, бағасы, ауданы және сипаттама) — хабарландыру бөлімде пайда болады.

Баннерлік жарнама үшін жарнама беруші кабинеті кейінірек қосылады: табыс автор, платформа және қоғам қоры арасында әділ бөлінеді.

## Құқықтар мен жауапкершілік
- Материалдарыңызға авторлық құқық сізде қалады, мазмұны үшін өзіңіз жауап бересіз.
- **Тек нақты ақпарат жариялаңыз.** Жала, көрінеу жалған ақпарат, ұлтаралық араздық пен өшпенділікті қоздыру және өзге де құқыққа қайшы мазмұн үшін автор Қазақстан Республикасының заңдары бойынша қылмыстық жауаптылықта болады.
- **Балағат әрі заңсыз материалдар, құрметсіз тон** оқырмандар (пайдаланушылар) мен біздің ИИ моделіміз арқылы модерацияланады; ондай материалдар төмендетіледі, жасырылады немесе жойылады.
- Бөтен мәтін мен суретті рұқсатсыз жарияламаңыз, авторлық құқықты бұзбаңыз.
- Жеке өмірді құрметтеңіз: басқалардың дербес деректерін келісімінсіз ашпаңыз.
- Спам, алаяқтық, тыйым салынған тауар мен қызмет жарнамасына тыйым салынады.
- 18+ материал және балаларға зиян келтіретін кез келген мазмұн жол берілмейді.
- Ереже бұзылса, материалды жасыруға, аккаунтты шектеуге не бұғаттауға құқылымыз. Шешімге [Қолдау](/support) бөлімі арқылы шағымдана аласыз.
- ИИ құралдары құрметті тонмен жазуға көмектеседі әрі дереккөз қажет жерді ескертеді, бірақ мәтін үшін түпкілікті жауапкершілік авторда.

Материал жариялау арқылы сіз осы ережелермен келісесіз.`},
		"en": {Title: "How to publish articles and post listings", Body: `This page explains how to use the platform, what is free, what is paid, and your rights and responsibilities.

## How to write an article
1. Sign up and open the **Author Studio**.
2. Click "Write" and create an article in one of the three languages.
3. Choose a category and subcategory, and add a cover image if you like.
4. Use the AI assistant to improve the text and auto-translation to publish in Kazakh, Russian, and English at once.
5. Save and click "Publish".

## What is free and what is paid
**Free forever:** reading all content, the newsletter, registration, publishing articles manually (including in three languages if you translate them yourself), votes, comments, and karma.

**Paid services** are the optional AI features and promotion: improving an article with the AI editor, automatic translation into three languages, cover generation, and posting or promoting real-estate listings. Current prices are on the [Pricing](/pricing) page.

## Real estate and listings
Our symbol is the shanyraq — a home. In the **Real Estate** section you can post a sale or rental listing: apartment, house, land, or commercial space. Click "Post a listing", fill in the fields (country, region, city, village, type, price, area, and description), and it will appear in the section.

An advertiser cabinet for banner ads will come later; revenue will be shared fairly between the author, the platform, and a community fund.

## Rights and responsibilities
- You keep the copyright to your work and are responsible for its content.
- **Publish only accurate information.** For defamation, knowingly false information, incitement of ethnic hatred or hostility, and other unlawful content, the author bears criminal liability under the laws of the Republic of Kazakhstan.
- **Abusive and unlawful material and a disrespectful tone** are moderated by our readers (users) and our AI model; such material is down-ranked, hidden, or removed.
- Do not publish others' texts or images without permission; do not infringe copyright.
- Respect privacy: do not disclose other people's personal data without consent.
- Spam, fraud, and advertising of prohibited goods and services are forbidden.
- Adult (18+) material and any content harmful to children are not allowed.
- If the rules are broken, we may hide the material or restrict or block the account. You can appeal via [Support](/support).
- AI tools help you write in a respectful tone and flag where a source is needed, but final responsibility for the text lies with the author.

By publishing, you agree to these rules.`},
	},
	"pricing": {
		"ru": {Title: "Тарифы", Body: `Мы верим, что доступ к информации должен быть бесплатным. Поэтому чтение и подписка не стоят ничего. Платными являются только дополнительные услуги на базе ИИ и продвижение — по прозрачным ценам.

## Бесплатно и навсегда
- Чтение всех материалов на трёх языках
- Подписка на рассылку
- Регистрация и ручная публикация статей
- Оценки, комментарии, карма

## Услуги для авторов (ИИ)
| Услуга | Цена |
| --- | --- |
| ИИ-редактор (улучшить статью) | 199 ₸ / статья |
| Автоперевод KZ ↔ RU ↔ EN | 149 ₸ / статья |
| ИИ-обложка | 99 ₸ / изображение |
| Пакет «Автор Про» (до 15 ИИ-статей + переводы + 30 обложек в месяц) | 2 490 ₸ / мес |

Цены рассчитаны от нашей себестоимости (оплата ИИ-моделей) плюс небольшая наценка. Платите только за то, чем пользуетесь.

## Объявления о недвижимости
| Услуга | Цена |
| --- | --- |
| Размещение объявления (30 дней) | 490 ₸ |
| Продвижение «Топ раздела» (7 дней) | 1 990 ₸ |
| Выделение объявления (7 дней) | 990 ₸ |

## Реклама для компаний
Баннерная реклама — по модели CPM (оплата за показы). Доход делится: **50% автору материала, 40% платформе, 10% в фонд сообщества.** Кабинет рекламодателя откроется позже.

Все цены указаны в тенге и могут корректироваться. Об изменениях мы сообщаем заранее.`},
		"kz": {Title: "Тарифтер", Body: `Ақпаратқа қолжетімділік тегін болуы керек деп сенеміз. Сондықтан оқу мен жазылу тегін. Тек ИИ негізіндегі қосымша қызметтер мен жылжыту ақылы — ашық бағамен.

## Тегін әрі әрқашан
- Барлық материалды үш тілде оқу
- Жаңалық таратылымына жазылу
- Тіркелу және мақаланы қолмен жариялау
- Бағалау, пікір, карма

## Авторларға арналған қызметтер (ИИ)
| Қызмет | Бағасы |
| --- | --- |
| ИИ-редактор (мақаланы жақсарту) | 199 ₸ / мақала |
| Автоаударма KZ ↔ RU ↔ EN | 149 ₸ / мақала |
| ИИ-мұқаба | 99 ₸ / сурет |
| «Автор Про» пакеті (айына 15 ИИ-мақала + аударма + 30 мұқаба) | 2 490 ₸ / ай |

Бағалар біздің өзіндік құнымыздан (ИИ модельдеріне төлем) және шағын үстемеден құралады. Тек пайдаланғаныңызға төлейсіз.

## Жылжымайтын мүлік хабарландырулары
| Қызмет | Бағасы |
| --- | --- |
| Хабарландыру орналастыру (30 күн) | 490 ₸ |
| «Бөлім Топы» жылжыту (7 күн) | 1 990 ₸ |
| Хабарландыруды ерекшелеу (7 күн) | 990 ₸ |

## Компанияларға жарнама
Баннерлік жарнама — CPM моделі бойынша (көрсетілім үшін төлем). Табыс бөлінісі: **50% материал авторына, 40% платформаға, 10% қоғам қорына.** Жарнама беруші кабинеті кейінірек ашылады.

Барлық баға теңгемен және өзгеруі мүмкін. Өзгерістер туралы алдын ала хабарлаймыз.`},
		"en": {Title: "Pricing", Body: `We believe access to information should be free. Reading and the newsletter cost nothing. Only optional AI services and promotion are paid — at transparent prices.

## Free forever
- Reading all content in three languages
- The newsletter
- Registration and manual publishing
- Votes, comments, karma

## Author services (AI)
| Service | Price |
| --- | --- |
| AI editor (improve a story) | 199 ₸ / story |
| Auto-translation KZ ↔ RU ↔ EN | 149 ₸ / story |
| AI cover image | 99 ₸ / image |
| "Author Pro" bundle (up to 15 AI stories + translations + 30 covers per month) | 2,490 ₸ / mo |

Prices are based on our cost (AI model usage) plus a small margin. You pay only for what you use.

## Real-estate listings
| Service | Price |
| --- | --- |
| Post a listing (30 days) | 490 ₸ |
| "Top of section" promotion (7 days) | 1,990 ₸ |
| Highlight a listing (7 days) | 990 ₸ |

## Advertising for businesses
Banner ads use a CPM model (pay per impression). Revenue is shared: **50% to the author, 40% to the platform, 10% to a community fund.** The advertiser cabinet will open later.

All prices are in tenge and may change. We announce changes in advance.`},
	},
	"support": {
		"ru": {Title: "Поддержка", Body: `Мы поможем разобраться с любым вопросом по работе платформы.

## Помощник на базе ИИ
Скоро здесь появится наш ИИ-ассистент, который будет отвечать на вопросы о том, как писать статьи, публиковать материалы и размещать рекламу — на казахском, русском и английском.

## Пока помощник готовится
- Прочитайте раздел [«Как публиковать статьи и размещать объявления»](/guide).
- Посмотрите [Тарифы](/pricing).
- Узнайте больше [о нас](/about).
- Напишите нам, если не нашли ответ.

Мы отвечаем спокойно и по существу.`},
		"kz": {Title: "Қолдау", Body: `Платформаның жұмысына қатысты кез келген сұрақты шешуге көмектесеміз.

## ИИ негізіндегі көмекші
Жақында мұнда біздің ИИ-ассистент пайда болады — ол мақала жазу, материал жариялау және жарнама орналастыру туралы сұрақтарға қазақ, орыс, ағылшын тілдерінде жауап береді.

## Көмекші дайындалып жатқанда
- [«Мақала мен хабарландыруды қалай жариялау керек»](/guide) бөлімін оқыңыз.
- [Тарифтер](/pricing) бетін қараңыз.
- [Біз туралы](/about) көбірек біліңіз.
- Жауап таппасаңыз, бізге жазыңыз.

Біз сабырмен әрі нақты жауап береміз.`},
		"en": {Title: "Support", Body: `We will help you with any question about how the platform works.

## AI-powered assistant
Soon our AI assistant will appear here to answer questions about writing articles, publishing, and placing ads — in Kazakh, Russian, and English.

## While the assistant is getting ready
- Read the [How to publish articles and post listings](/guide) section.
- See the [Pricing](/pricing) page.
- Learn more [about us](/about).
- Write to us if you did not find an answer.

We reply calmly and to the point.`},
	},
}

// staticContent returns the page for key in lang, falling back to Russian.
func staticContent(key, lang string) staticPageContent {
	if m, ok := staticPages[key]; ok {
		if c, ok := m[lang]; ok {
			return c
		}
		return m[LangRU]
	}
	return staticPageContent{}
}
