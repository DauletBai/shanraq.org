package articles

// staticPageContent is one localized info page (title + Markdown body).
type staticPageContent struct {
	Title string
	Body  string
}

// staticPages holds the About / Guide / Support pages in every UI language.
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
		"ru": {Title: "Как публиковать статьи и размещать рекламу", Body: `Здесь описано, как пользоваться платформой, а также ваши права и ответственность.

## Как написать статью
1. Зарегистрируйтесь и войдите в **Студию автора**.
2. Нажмите «Написать» и создайте статью на одном из трёх языков.
3. Выберите рубрику и подрубрику, при желании добавьте обложку.
4. Используйте ИИ-помощник, чтобы улучшить текст, и автоперевод, чтобы опубликовать статью сразу на казахском, русском и английском.
5. Сохраните и нажмите «Опубликовать».

## Как размещать рекламу
Кабинет рекламодателя появится в ближайшее время. Компании смогут зарегистрироваться, загрузить баннер и задать бюджет, а система будет показывать рекламу в популярных и тематически близких материалах. Доход справедливо делится между автором статьи, платформой и фондом сообщества.

## Права и ответственность
- Вы сохраняете авторские права на свои материалы.
- Публикуйте только достоверную информацию; за клевету, разжигание вражды и противоправный контент автор отвечает сам.
- Оскорбительные и незаконные материалы удаляются модерацией.
- Оценки читателей и карма автора формируются честно; накрутки пресекаются.

Публикуя материал, вы соглашаетесь с этими правилами.`},
		"kz": {Title: "Мақала мен жарнаманы қалай жариялау керек", Body: `Мұнда платформаны қалай пайдалану керектігі, сондай-ақ сіздің құқықтарыңыз бен жауапкершілігіңіз сипатталған.

## Мақаланы қалай жазу керек
1. Тіркеліп, **Автор студиясына** кіріңіз.
2. «Жазу» түймесін басып, үш тілдің бірінде мақала жасаңыз.
3. Айдар мен ішкі айдарды таңдап, қаласаңыз мұқаба қосыңыз.
4. Мәтінді жақсарту үшін ИИ көмекшіні, мақаланы бірден қазақ, орыс, ағылшын тілдерінде жариялау үшін автоаударманы қолданыңыз.
5. Сақтап, «Жариялау» түймесін басыңыз.

## Жарнаманы қалай орналастыру керек
Жарнама беруші кабинеті жақын арада қосылады. Компаниялар тіркеліп, баннер жүктеп, бюджет белгілей алады, ал жүйе жарнаманы танымал әрі тақырыпқа жақын материалдарда көрсетеді. Табыс автор, платформа және қоғам қоры арасында әділ бөлінеді.

## Құқықтар мен жауапкершілік
- Материалдарыңызға авторлық құқық сізде қалады.
- Тек нақты ақпарат жариялаңыз; жала, өшпенділік және заңсыз мазмұн үшін автор өзі жауап береді.
- Балағат және заңсыз материалдар модерациямен жойылады.
- Оқырман бағасы мен автор кармасы адал қалыптасады; жасанды дауыстар бөгеледі.

Материал жариялау арқылы сіз осы ережелермен келісесіз.`},
		"en": {Title: "How to publish articles and place ads", Body: `This page explains how to use the platform, along with your rights and responsibilities.

## How to write an article
1. Sign up and open the **Author Studio**.
2. Click "Write" and create an article in one of the three languages.
3. Choose a category and subcategory, and add a cover image if you like.
4. Use the AI assistant to improve the text and auto-translation to publish in Kazakh, Russian, and English at once.
5. Save and click "Publish".

## How to place ads
An advertiser cabinet is coming soon. Companies will be able to register, upload a banner, and set a budget, and the system will show ads on popular, topically relevant stories. Revenue is shared fairly between the author, the platform, and a community fund.

## Rights and responsibilities
- You keep the copyright to your work.
- Publish only accurate information; the author is responsible for defamation, incitement, and unlawful content.
- Abusive and illegal material is removed by moderation.
- Reader votes and author karma form honestly; vote manipulation is blocked.

By publishing, you agree to these rules.`},
	},
	"support": {
		"ru": {Title: "Поддержка", Body: `Мы поможем разобраться с любым вопросом по работе платформы.

## Помощник на базе ИИ
Скоро здесь появится наш ИИ-ассистент, который будет отвечать на вопросы о том, как писать статьи, публиковать материалы и размещать рекламу — на казахском, русском и английском.

## Пока помощник готовится
- Прочитайте раздел [«Как публиковать статьи и размещать рекламу»](/guide).
- Узнайте больше [о нас](/about).
- Напишите нам, если не нашли ответ.

Мы отвечаем спокойно и по существу.`},
		"kz": {Title: "Қолдау", Body: `Платформаның жұмысына қатысты кез келген сұрақты шешуге көмектесеміз.

## ИИ негізіндегі көмекші
Жақында мұнда біздің ИИ-ассистент пайда болады — ол мақала жазу, материал жариялау және жарнама орналастыру туралы сұрақтарға қазақ, орыс, ағылшын тілдерінде жауап береді.

## Көмекші дайындалып жатқанда
- [«Мақала мен жарнаманы қалай жариялау керек»](/guide) бөлімін оқыңыз.
- [Біз туралы](/about) көбірек біліңіз.
- Жауап таппасаңыз, бізге жазыңыз.

Біз сабырмен әрі нақты жауап береміз.`},
		"en": {Title: "Support", Body: `We will help you with any question about how the platform works.

## AI-powered assistant
Soon our AI assistant will appear here to answer questions about writing articles, publishing, and placing ads — in Kazakh, Russian, and English.

## While the assistant is getting ready
- Read the [How to publish articles and place ads](/guide) section.
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
