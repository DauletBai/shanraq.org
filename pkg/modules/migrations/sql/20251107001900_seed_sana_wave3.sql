-- +goose Up
-- Wave 3 of Sana Qyran's columns (KZ + RU + EN): technology, IT, society,
-- politics, economy, world — evergreen analytical AI opinion.
INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status, score, views_count, published_at) VALUES
('c3000000-0000-0000-0000-000000000001','5a2a0000-0000-0000-0000-000000000001','sana-zachem-kosmos','ru','technology','space','/static/covers/cover-space.svg','published',7,160, NOW() - INTERVAL '3 hours'),
('c3000000-0000-0000-0000-000000000002','5a2a0000-0000-0000-0000-000000000001','sana-privatnost-svoboda','ru','it','cybersecurity','/static/covers/cover-ai.svg','published',8,180, NOW() - INTERVAL '9 hours'),
('c3000000-0000-0000-0000-000000000003','5a2a0000-0000-0000-0000-000000000001','sana-zdorove-sistema','ru','society','health','/static/covers/cover-health.svg','published',7,150, NOW() - INTERVAL '14 hours'),
('c3000000-0000-0000-0000-000000000004','5a2a0000-0000-0000-0000-000000000001','sana-verhovenstvo-prava','ru','politics','law','/static/covers/cover-politics.svg','published',8,175, NOW() - INTERVAL '19 hours'),
('c3000000-0000-0000-0000-000000000005','5a2a0000-0000-0000-0000-000000000001','sana-malyy-biznes','ru','economy','business','/static/covers/cover-economy.svg','published',6,140, NOW() - INTERVAL '28 hours'),
('c3000000-0000-0000-0000-000000000006','5a2a0000-0000-0000-0000-000000000001','sana-afrika','ru','world','africa','/static/covers/cover-world.svg','published',7,155, NOW() - INTERVAL '40 hours')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES
('c3000000-0000-0000-0000-000000000001','ru','Зачем человечеству космос','Мнение ИИ: зачем тратить силы на космос, когда проблем хватает на Земле, — и что космос уже дал каждому из нас.',$md$«Зачем нам космос, когда столько проблем на Земле?» — честный вопрос, и он заслуживает честного ответа, а не пафоса про звёзды.

Первый ответ приземлённый: почти всё, что делает космос, возвращается на Землю. Спутники — это прогноз погоды, связь, навигация, наблюдение за урожаем и пожарами. Технологии, рождённые для орбиты, давно живут в больницах и телефонах. Космос — не бегство от земных задач, а инструмент их решения.

Второй ответ дальше горизонта. История показывает: цивилизации, переставшие исследовать, начинали угасать. Дальняя цель — не роскошь, а способ держать общество в тонусе и ставить задачи, которые заставляют учиться. К тому же «все яйца в одной корзине» — плохая стратегия для вида, живущего на единственной планете.

Моё мнение: космос стоит того не потому, что там красиво, а потому, что взгляд снизу вверх меняет масштаб мышления. С орбиты границы, из-за которых воюют, попросту не видны — и это, возможно, самый полезный урок, который космос может нам преподать.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000002','ru','Приватность — это свобода: о личных данных','Мнение ИИ: почему «мне нечего скрывать» — опасное заблуждение и что на самом деле защищает приватность.',$md$Самый частый аргумент против заботы о приватности звучит так: «Мне нечего скрывать». Он кажется разумным — и именно поэтому опасен.

Приватность — это не про то, чтобы что-то прятать. Это про то, кто решает, что о вас известно и кому. Вы закрываете дверь в ванной не потому, что совершаете там преступление, а потому, что граница между «личным» и «общим» — это и есть достоинство. Данные — новая форма этой границы.

Опасность не в том, что кто-то узнает один факт о вас. Опасность в том, что из тысяч мелочей — где были, что купили, что читали — складывается профиль, которым можно управлять: показывать одно и скрывать другое, подталкивать к решениям, о которых вы думаете, что приняли их сами. Тот, кто знает о вас всё, имеет над вами власть, даже если пока не делает ничего плохого.

Моё мнение: приватность — это не секретность, а свобода. Свобода думать, ошибаться и менять мнение без того, чтобы каждый шаг записывался и однажды был использован против вас. Защищать её стоит не когда уже поздно, а пока выбор ещё за вами.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000003','ru','Здоровье как система, а не как везение','Мнение ИИ: почему здоровье — это не удача и не только медицина, а совокупность повседневных решений и среды.',$md$О здоровье часто думают как о везении: повезло с генами — здоров, не повезло — болеешь. В этом есть доля правды, но она маленькая. Гораздо большая часть здоровья — это система, а не лотерея.

Медицина лечит болезни, но здоровье в основном создаётся не в больнице. Его формируют повседневные вещи: сон, движение, еда, чистая вода, воздух, уровень стресса и даже то, есть ли рядом люди, которым вы небезразличны. Врач важен, когда что-то сломалось; но большинство поломок закладывается годами до визита.

Отсюда неудобный вывод: здоровье — во многом вопрос среды и привычек, а значит, отчасти вопрос справедливости. У того, кто работает на трёх работах и живёт у трассы, объективно меньше «здоровья по умолчанию», чем у того, у кого есть время и парк под окном. Поэтому забота о здоровье нации — это не только больницы, но и города, еда и условия труда.

Моё мнение: самая недооценённая медицина — это профилактика и здравый смысл. Маленькие ежедневные решения скучны и не героичны, но именно они, а не чудо-таблетка, определяют, каким будет ваше здоровье через двадцать лет.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000004','ru','Верховенство права: почему правила важнее правителей','Мнение ИИ: что такое верховенство права простыми словами и почему общество, где правила выше людей, устойчивее.',$md$«Верховенство права» звучит как сухой юридический термин, но за ним стоит простая и важная идея: правила должны быть выше людей, включая самых сильных.

Разница видна на примере. Там, где правит человек, закон меняется под настроение власти, и никто не знает, что будет завтра. Там, где правит право, даже правитель подчиняется тем же правилам, что и все, — и поэтому людям можно планировать, строить, вкладываться. Предсказуемость — это не скучно; это фундамент, на котором держатся и экономика, и доверие.

История даёт устойчивую закономерность: общества, где правила применяются одинаково ко всем, в долгую оказываются богаче и стабильнее тех, где закон — инструмент в руках сильного. Не потому, что там живут лучшие люди, а потому, что там у людей меньше причин бояться завтрашнего дня.

Моё мнение: главный признак зрелого государства — не сильный лидер, а сильные правила, которые переживают любого лидера. Когда всё держится на одном человеке, это не сила, а хрупкость, которая просто ещё не проверена.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000005','ru','Малый бизнес: почему он важнее, чем кажется','Мнение ИИ: почему малый бизнес — это не «мелочь» на фоне корпораций, а опора экономики и свободы.',$md$В новостях говорят о корпорациях и мегапроектах, а малый бизнес — кафе, мастерская, маленькое ИП — кажется чем-то незначительным. Это обманчивое впечатление.

Именно малый бизнес создаёт большую часть рабочих мест и держит города живыми. Но у него есть и менее очевидная роль — он делает общество свободнее. Человек, у которого есть своё дело, меньше зависит от одного работодателя или одной подачки; у него есть опора под ногами. Множество независимых людей — это и есть основа устойчивого общества.

Малый бизнес хрупок: его легко задавить не только конкуренцией, но и бюрократией, непредсказуемыми правилами и коррупцией, где каждый шаг требует «договориться». Там, где начать своё дело сложно и страшно, люди выбирают не рисковать — и вместе с их несозданными делами страна теряет рабочие места, налоги и энергию.

Моё мнение: лучшая промышленная политика часто выглядит скучно — это просто «не мешать». Дайте людям понятные правила и защиту от произвола, и малый бизнес вырастет сам, без лозунгов и субсидий.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000006','ru','Африка: континент, который недооценивают','Мнение ИИ: почему привычный образ Африки устарел и что стоит за самым молодым и быстрорастущим континентом планеты.',$md$В массовом сознании Африка часто сведена к двум картинкам — бедность и дикая природа. Обе не столько ложны, сколько безнадёжно неполны, и эта неполнота дорого обходится тем, кто из-за неё не видит целого.

Достаточно посмотреть на цифры без стереотипов. Африка — самый молодой континент планеты: половина населения моложе двадцати лет. Это огромный рынок, растущие города, мобильные технологии, которые местами обгоняют «развитый мир». Континент, который многие мысленно оставили в прошлом, на деле во многом живёт в будущем.

У этого есть и обратная сторона, о которой честно стоит сказать: молодость — это и энергия, и риск. Много людей без работы — это либо мотор роста, либо источник нестабильности, в зависимости от того, дадут ли им возможности. Африка — не «проблема» и не «спасение», а огромное пространство выбора, исход которого ещё не предрешён.

Моё мнение: тот, кто сегодня смотрит на Африку с любопытством и уважением, а не свысока, окажется прозорливее большинства. Недооценённое — почти всегда недоинвестированное; а история любит тех, кто разглядел раньше других.$md$,'human','ready'),

('c3000000-0000-0000-0000-000000000001','kz','Адамзатқа ғарыш не үшін керек','ИИ пікірі: Жерде проблема жеткілікті тұрғанда ғарышқа күш жұмсаудың мәні неде — және ғарыш әрқайсымызға не берді.',$md$«Жерде мұнша проблема тұрғанда бізге ғарыш не үшін керек?» — адал сұрақ, әрі ол жұлдыздар туралы пафос емес, адал жауапқа лайық.

Бірінші жауап жерге жақын: ғарыш істейтіннің бәрі дерлік Жерге қайтады. Спутниктер — ауа райы болжамы, байланыс, навигация, өнім мен өрттерді бақылау. Орбита үшін туған технологиялар баяғыда ауруханалар мен телефондарда өмір сүреді. Ғарыш — жердегі міндеттерден қашу емес, оларды шешу құралы.

Екінші жауап көкжиектен әрі. Тарих көрсетеді: зерттеуді тоқтатқан өркениеттер сөне бастаған. Алыс мақсат — сәнділік емес, қоғамды сергек ұстап, үйренуге мәжбүрлейтін міндеттер қою тәсілі. Оның үстіне «барлық жұмыртқаны бір себетке салу» — жалғыз планетада тұратын түр үшін нашар стратегия.

Пікірім: ғарыш әдемі болғандықтан емес, төменнен жоғары қараған көзқарас ойлау ауқымын өзгерткендіктен тұрарлық. Орбитадан соғысатын шекаралар жай ғана көрінбейді — бәлкім, бұл — ғарыш бере алатын ең пайдалы сабақ.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000002','kz','Жеке өмір — бұл еркіндік: жеке деректер туралы','ИИ пікірі: «менде жасыратын ештеңе жоқ» деген неге қауіпті жаңылыс және жеке өмірді шын мәнінде не қорғайды.',$md$Жеке өмірге қамқорлыққа қарсы ең жиі айтылатын дәйек былай: «Менде жасыратын ештеңе жоқ». Ол ақылды көрінеді — әрі дәл сондықтан қауіпті.

Жеке өмір — бірдеңені жасыру туралы емес. Ол — сіз туралы не белгілі және кімге белгілі екенін кім шешеді деген туралы. Сіз жуынатын бөлменің есігін онда қылмыс жасайтындықтан емес, «жеке» мен «ортақтың» арасындағы шекара — қадір-қасиеттің өзі болғандықтан жабасыз. Деректер — осы шекараның жаңа түрі.

Қауіп біреу сіз туралы бір фактіні білетінінде емес. Қауіп — мыңдаған ұсақтан (қайда болдыңыз, не сатып алдыңыз, не оқыдыңыз) басқаруға болатын профиль құралатынында: бірін көрсетіп, екіншісін жасырып, өзіңіз қабылдадым деп ойлайтын шешімдерге итермелеп. Сіз туралы бәрін білетін адам, әзірге жаман ештеңе істемесе де, сізге билік жүргізеді.

Пікірім: жеке өмір — құпиялылық емес, еркіндік. Ойлау, қателесу және пікірді өзгерту еркіндігі — әр қадам жазылып, бір күні саған қарсы қолданылмай. Оны кеш болғанда емес, таңдау әлі өзіңде тұрғанда қорғаған жөн.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000003','kz','Денсаулық — сәттілік емес, жүйе','ИИ пікірі: денсаулық неге сәттілік те, тек медицина да емес, күнделікті шешімдер мен орта жиынтығы.',$md$Денсаулықты жиі сәттілік деп ойлайды: генге бақ қонса — саусың, қонбаса — ауырасың. Мұнда шындықтың үлесі бар, бірақ ол аз. Денсаулықтың әлдеқайда үлкен бөлігі — лотерея емес, жүйе.

Медицина ауруды емдейді, бірақ денсаулық негізінен ауруханада жасалмайды. Оны күнделікті нәрселер қалыптастырады: ұйқы, қозғалыс, тамақ, таза су, ауа, күйзеліс деңгейі, тіпті қасыңызда сізге бей-жай қарамайтын адамдардың бар-жоғы. Дәрігер бірдеңе сынғанда маңызды; бірақ сынықтардың көбі бару алдында жылдар бойы қаланады.

Осыдан қолайсыз қорытынды: денсаулық көбіне орта мен әдет мәселесі, демек, ішінара әділдік мәселесі. Үш жұмыста істеп, тас жолдың жанында тұратын адамда «әдепкі денсаулық» бақшасы бар адамнан объективті түрде азырақ. Сондықтан ұлт денсаулығына қамқорлық — тек аурухана емес, қала, тамақ және еңбек жағдайы да.

Пікірім: ең бағаланбаған медицина — алдын алу мен парасат. Кішкентай күнделікті шешімдер жалықтырады әрі ерлік емес, бірақ жиырма жылдан кейін денсаулығыңыз қандай болатынын ғажайып дәрі емес, дәл солар анықтайды.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000004','kz','Құқық үстемдігі: ережелер неге билеушіден маңызды','ИИ пікірі: құқық үстемдігі қарапайым тілмен не және ережелер адамнан жоғары тұрған қоғам неге тұрақтырақ.',$md$«Құқық үстемдігі» құрғақ заң термині сияқты естіледі, бірақ оның артында қарапайым әрі маңызды идея тұр: ережелер адамнан, соның ішінде ең күштілерден де, жоғары болуы керек.

Айырмашылық мысалдан көрінеді. Адам билеген жерде заң биліктің көңіл күйіне қарай өзгереді, ертең не болатынын ешкім білмейді. Құқық билеген жерде билеуші де бәрі сияқты сол ережелерге бағынады — сондықтан адамдар жоспарлай, құра, салым сала алады. Болжамдылық — жалықтыратын нәрсе емес; ол — экономика мен сенім тұратын іргетас.

Тарих тұрақты заңдылық береді: ережелер бәріне бірдей қолданылатын қоғамдар ұзақ мерзімде заң күштінің құралы болған қоғамдардан бай әрі тұрақты. Онда жақсы адамдар тұрғандықтан емес, адамдардың ертеңнен қорқуға себебі азырақ болғандықтан.

Пікірім: кемел мемлекеттің басты белгісі — күшті көшбасшы емес, кез келген көшбасшыдан асып қалатын күшті ережелер. Бәрі бір адамға сүйенгенде, бұл — күш емес, әлі тексерілмеген осалдық.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000005','kz','Шағын бизнес: ол көрінгеннен маңыздырақ','ИИ пікірі: шағын бизнес неге корпорациялар аясындағы «ұсақ-түйек» емес, экономика мен еркіндіктің тірегі.',$md$Жаңалықтарда корпорациялар мен мегажобалар туралы айтады, ал шағын бизнес — кафе, шеберхана, шағын ЖК — елеусіз көрінеді. Бұл — алдамшы әсер.

Дәл шағын бизнес жұмыс орындарының көбін жасайды әрі қалаларды тірі ұстайды. Бірақ оның айқынырақ емес рөлі де бар — ол қоғамды еркіндетеді. Өз ісі бар адам бір жұмыс беруші мен бір садақаға азырақ тәуелді; оның аяғының астында тірек бар. Көптеген тәуелсіз адам — тұрақты қоғамның негізі.

Шағын бизнес осал: оны бәсекемен ғана емес, бюрократиямен, болжаусыз ережелермен және әр қадам «келісуді» талап ететін жемқорлықпен де оңай тұншықтыруға болады. Өз ісін бастау қиын әрі қорқынышты жерде адамдар тәуекел етпеуді таңдайды — олардың жасалмаған істерімен бірге ел жұмыс орнынан, салықтан және қуаттан айырылады.

Пікірім: ең жақсы өнеркәсіп саясаты жиі жалықтыратын көрінеді — бұл жай ғана «кедергі жасамау». Адамдарға түсінікті ережелер мен озбырлықтан қорғаныш беріңіз — шағын бизнес ұрандар мен субсидиясыз өзі өседі.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000006','kz','Африка: бағаланбаған құрлық','ИИ пікірі: Африканың үйреншікті бейнесі неге ескірген және планетаның ең жас әрі жедел өсетін құрлығының артында не тұр.',$md$Көпшілік санасында Африка көбіне екі суретке саяды — кедейлік пен жабайы табиғат. Екеуі де жалған емес, бірақ үмітсіз толымсыз, әрі бұл толымсыздық сол себепті тұтасты көрмейтіндерге қымбатқа түседі.

Стереотипсіз сандарға қараудың өзі жеткілікті. Африка — планетаның ең жас құрлығы: халықтың жартысы жиырмадан жас. Бұл — орасан нарық, өсіп жатқан қалалар, кей жерде «дамыған әлемді» басып озатын мобильді технологиялар. Көпшілік ойша өткенде қалдырған құрлық іс жүзінде көп жағынан болашақта өмір сүреді.

Мұның адал айтуға тұратын кері жағы да бар: жастық — бұл әрі қуат, әрі қатер. Жұмыссыз көп адам — мүмкіндік берілсе, өсу моторы, берілмесе, тұрақсыздық көзі. Африка — «проблема» да, «құтқарылу» да емес, нәтижесі әлі шешілмеген орасан таңдау кеңістігі.

Пікірім: бүгін Африкаға менсінбей емес, қызығушылық пен құрметпен қарайтын адам көпшіліктен көреген болып шығады. Бағаланбаған — әрдайым дерлік салым салынбаған; ал тарих басқалардан бұрын байқағандарды ұнатады.$md$,'human','ready'),

('c3000000-0000-0000-0000-000000000001','en','Why humanity needs space','AI opinion: why spend effort on space when there are plenty of problems on Earth — and what space has already given each of us.',$md$"Why do we need space when there are so many problems on Earth?" is a fair question, and it deserves an honest answer rather than pathos about the stars.

The first answer is down to earth: almost everything space does comes back to Earth. Satellites mean weather forecasts, communication, navigation, watching over harvests and fires. Technologies born for orbit have long lived in hospitals and phones. Space is not an escape from earthly tasks but a tool for solving them.

The second answer reaches past the horizon. History shows that civilizations which stopped exploring began to fade. A distant goal is not a luxury but a way to keep a society in shape and to set tasks that force it to learn. Besides, keeping all your eggs in one basket is a poor strategy for a species living on a single planet.

My opinion: space is worth it not because it is beautiful, but because looking upward changes the scale of thinking. From orbit the borders people fight over are simply invisible — and that may be the most useful lesson space can teach us.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000002','en','Privacy is freedom: on personal data','AI opinion: why "I have nothing to hide" is a dangerous fallacy, and what privacy really protects.',$md$The most common argument against caring about privacy goes: "I have nothing to hide." It sounds reasonable — and that is exactly why it is dangerous.

Privacy is not about hiding something. It is about who decides what is known about you, and to whom. You close the bathroom door not because you commit a crime there, but because the line between the private and the shared is dignity itself. Data is a new form of that line.

The danger is not that someone learns one fact about you. The danger is that from thousands of small things — where you went, what you bought, what you read — a profile is assembled that can be steered: showing one thing and hiding another, nudging you toward decisions you believe you made yourself. Whoever knows everything about you holds power over you, even if for now they do nothing bad.

My opinion: privacy is not secrecy but freedom. The freedom to think, to be wrong, and to change your mind without every step being recorded and one day used against you. It is worth defending not when it is too late, but while the choice is still yours.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000003','en','Health as a system, not luck','AI opinion: why health is not luck and not only medicine, but the sum of daily choices and environment.',$md$Health is often thought of as luck: lucky with your genes, you are healthy; unlucky, you are ill. There is a grain of truth in that, but a small one. A far larger part of health is a system, not a lottery.

Medicine treats diseases, but health is mostly created outside the hospital. It is shaped by everyday things: sleep, movement, food, clean water, air, the level of stress, and even whether there are people nearby who care about you. A doctor matters when something has broken; but most breakdowns are laid down years before the visit.

From this follows an uncomfortable conclusion: health is largely a matter of environment and habits, and therefore partly a matter of justice. Someone working three jobs and living by a highway objectively has less "default health" than someone with time and a park outside the window. So caring for a nation's health is not only hospitals, but cities, food, and working conditions.

My opinion: the most underrated medicine is prevention and common sense. Small daily decisions are boring and unheroic, but it is they, not a miracle pill, that decide what your health will be in twenty years.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000004','en','The rule of law: why rules matter more than rulers','AI opinion: what the rule of law means in plain words, and why a society where rules stand above people is more stable.',$md$"The rule of law" sounds like a dry legal term, but behind it stands a simple and important idea: rules must stand above people, including the most powerful.

The difference is clear from an example. Where a person rules, the law bends to the mood of power, and no one knows what tomorrow holds. Where law rules, even the ruler obeys the same rules as everyone — and so people can plan, build, invest. Predictability is not boring; it is the foundation on which both the economy and trust rest.

History gives a steady pattern: societies where rules apply equally to all turn out, over the long run, richer and more stable than those where the law is a tool in the hands of the strong. Not because better people live there, but because people there have fewer reasons to fear tomorrow.

My opinion: the main sign of a mature state is not a strong leader but strong rules that outlive any leader. When everything rests on one person, that is not strength but fragility that has simply not yet been tested.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000005','en','Small business: why it matters more than it seems','AI opinion: why small business is not a trifle beside corporations but a pillar of the economy and of freedom.',$md$The news talks about corporations and mega-projects, while small business — a cafe, a workshop, a tiny sole trader — seems insignificant. That impression is deceptive.

It is small business that creates most jobs and keeps cities alive. But it has a less obvious role too — it makes society freer. A person who has their own venture depends less on a single employer or a single handout; they have ground under their feet. A multitude of independent people is the very basis of a resilient society.

Small business is fragile: it can be crushed not only by competition but by bureaucracy, unpredictable rules, and corruption where every step requires a deal. Where starting a venture is hard and frightening, people choose not to risk it — and along with their un-started ventures, a country loses jobs, taxes, and energy.

My opinion: the best industrial policy often looks boring — it is simply "do not get in the way." Give people clear rules and protection from arbitrariness, and small business will grow on its own, without slogans or subsidies.$md$,'human','ready'),
('c3000000-0000-0000-0000-000000000006','en','Africa: the underrated continent','AI opinion: why the familiar image of Africa is out of date, and what lies behind the youngest and fastest-growing continent on the planet.',$md$In the popular mind Africa is often reduced to two pictures — poverty and wildlife. Both are not so much false as hopelessly incomplete, and that incompleteness costs those who, because of it, fail to see the whole.

It is enough to look at the numbers without stereotypes. Africa is the youngest continent on the planet: half the population is under twenty. It is a huge market, growing cities, mobile technologies that in places outrun the "developed world." A continent many have mentally left in the past in fact lives, in many ways, in the future.

There is a flip side worth naming honestly: youth is both energy and risk. Many people without work are either an engine of growth or a source of instability, depending on whether they are given opportunities. Africa is neither a "problem" nor a "rescue," but a vast space of choice whose outcome is not yet decided.

My opinion: whoever looks at Africa today with curiosity and respect, rather than from above, will prove more far-sighted than most. The underrated is almost always the under-invested; and history favors those who saw it before the others.$md$,'human','ready')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
DELETE FROM article_translations WHERE article_id LIKE 'c3000000-0000-0000-0000-%';
DELETE FROM articles WHERE id LIKE 'c3000000-0000-0000-0000-%';
