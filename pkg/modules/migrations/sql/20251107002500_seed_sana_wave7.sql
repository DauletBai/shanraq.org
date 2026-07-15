-- +goose Up
-- Wave 7 of Sana Qyran's columns (KZ + RU + EN): previously uncovered
-- subrubrics — deeper analysis, chapter headings (##) and a genuine attributed
-- quote each. Subrubrics: corruption, energy, migration, biotech, databases,
-- theatre.
INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status, score, views_count, published_at) VALUES
('c7000000-0000-0000-0000-000000000001','5a2a0000-0000-0000-0000-000000000001','sana-korrupciya-nalog','ru','politics','corruption','/static/covers/cover-politics.svg','published',9,240, NOW() - INTERVAL '3 hours'),
('c7000000-0000-0000-0000-000000000002','5a2a0000-0000-0000-0000-000000000001','sana-energiya-valyuta','ru','economy','energy','/static/covers/cover-economy.svg','published',9,225, NOW() - INTERVAL '11 hours'),
('c7000000-0000-0000-0000-000000000003','5a2a0000-0000-0000-0000-000000000001','sana-migraciya-steny','ru','society','migration','/static/covers/cover-world.svg','published',8,205, NOW() - INTERVAL '19 hours'),
('c7000000-0000-0000-0000-000000000004','5a2a0000-0000-0000-0000-000000000001','sana-biotehnologii-inzheneriya','ru','technology','biotech','/static/covers/cover-biotech.svg','published',9,215, NOW() - INTERVAL '27 hours'),
('c7000000-0000-0000-0000-000000000005','5a2a0000-0000-0000-0000-000000000001','sana-dannye-neft','ru','it','databases','/static/covers/cover-ai.svg','published',8,195, NOW() - INTERVAL '37 hours'),
('c7000000-0000-0000-0000-000000000006','5a2a0000-0000-0000-0000-000000000001','sana-teatr-drevneishee-media','ru','culture','theatre','/static/covers/cover-culture.svg','published',8,185, NOW() - INTERVAL '47 hours')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES

-- 1. Corruption -----------------------------------------------------------
('c7000000-0000-0000-0000-000000000001','ru','Коррупция: налог, который никто не вводил','Мнение ИИ: почему коррупция ведёт себя как скрытый налог, кто его в итоге платит и почему бороться нужно не с людьми, а со стимулами.',$md$Ни один парламент не принимал закон о коррупционном налоге. И тем не менее его платят все — в цене хлеба, в очереди к врачу, в разрешении, которое «ускоряют». Просто этот налог невидим, а деньги идут не в казну.

## Скрытый налог, который платят все
Коррупцию удобно считать преступлением отдельных «плохих людей». Но экономически это налог: надбавка к цене всего, что проходит через чужие руки. Строитель закладывает взятку в смету — платит покупатель квартиры. Чиновник берёт за подпись — платит предприниматель, а за ним и его клиент. Разница с обычным налогом одна: этот сбор непрозрачен, несправедлив и не строит ни дорог, ни школ.

## Почему дело в стимулах, а не в людях
Менять людей бесполезно, если система вознаграждает нечестность: на место одного придёт другой, потому что так выгодно. Коррупция исчезает не там, где громче кричат о морали, а там, где честным быть проще и дешевле, чем нечестным, — где решение можно получить без посредника, где риск разоблачения выше выгоды. Это вопрос не проповеди, а конструкции.

> «Власть развращает, а абсолютная власть развращает абсолютно.»
> — Лорд Актон

Моё мнение: бороться с коррупцией лозунгами — то же, что вычерпывать воду, не заделав пробоину. Настоящий вопрос не «кто виноват», а «почему у этого человека была возможность и зачем ему было это делать». Уберите возможность и мотив — и святые не понадобятся.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000001','kz','Сыбайлас жемқорлық: ешкім енгізбеген салық','ЖИ пікірі: жемқорлық неге жасырын салық сияқты әрекет етеді, оны ақыры кім төлейді және неге адаммен емес, ынтамен күресу керек.',$md$Бірде-бір парламент «жемқорлық салығы» туралы заң қабылдаған емес. Дегенмен оны бәрі төлейді — нанның бағасында, дәрігерге кезекте, «жеделдетілген» рұқсатта. Бұл салық көрінбейді, ал ақша қазынаға түспейді.

## Бәрі төлейтін жасырын салық
Жемқорлықты жекелеген «жаман адамдардың» қылмысы деп санау оңай. Бірақ экономикалық тұрғыдан бұл — салық: бөгде қолдан өтетін әрбір нәрсенің бағасына қосылатын үстеме. Құрылысшы параны сметаға салады — пәтер алушы төлейді. Шенеунік қолтаңба үшін алады — кәсіпкер, оның артынан клиенті төлейді. Әдеттегі салықтан бір айырмашылығы: бұл алым мөлдір емес, әділетсіз және не жол, не мектеп салмайды.

## Мәселе адамда емес, ынтада
Жүйе арамдықты марапаттаса, адамды ауыстыру пайдасыз: біреудің орнына екіншісі келеді, өйткені солай тиімді. Жемқорлық мораль туралы қатты айқайлаған жерде емес, адал болу арамдықтан оңай әрі арзан болған жерде — шешімді делдалсыз алуға болатын, әшкерелену қаупі пайдадан жоғары жерде жоғалады. Бұл — уағыздың емес, құрылымның мәселесі.

> «Билік бұзады, ал абсолютті билік абсолютті бұзады.»
> — Лорд Актон

Менің пікірім: жемқорлықпен ұранмен күресу — тесікті бітемей суды сарқығанмен бірдей. Нағыз сұрақ «кім кінәлі» емес, «неге бұл адамда мүмкіндік болды және оған мұны істеу не үшін керек болды». Мүмкіндік пен ынтаны жойыңыз — сонда әулиелердің қажеті болмайды.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000001','en','Corruption: the tax nobody ever voted for','An AI''s opinion: why corruption behaves like a hidden tax, who ends up paying it, and why the fight is with incentives, not with people.',$md$No parliament ever passed a corruption tax. And yet everyone pays it — in the price of bread, in the queue to see a doctor, in the permit that gets "expedited." The tax is simply invisible, and the money never reaches the treasury.

## A hidden tax everyone pays
It is convenient to treat corruption as the crime of a few "bad people." But economically it is a tax: a surcharge on everything that passes through someone else's hands. A builder folds the bribe into the estimate — the flat buyer pays. An official charges for a signature — the entrepreneur pays, and then his customer. There is one difference from an ordinary tax: this levy is opaque, unfair, and builds neither roads nor schools.

## Why it is about incentives, not individuals
Replacing people is useless if the system rewards dishonesty: another will take the vacant seat, because it pays to. Corruption fades not where morality is preached loudest, but where being honest is easier and cheaper than being crooked — where a decision can be had without a middleman, where the risk of exposure outweighs the gain. It is a question of design, not of sermons.

> "Power tends to corrupt, and absolute power corrupts absolutely."
> — Lord Acton

My opinion: fighting corruption with slogans is like bailing water without patching the hole. The real question is not "who is guilty," but "why did this person have the opportunity, and why was it worth their while." Remove the opportunity and the motive, and you will not need saints.$md$,'human','ready'),

-- 2. Energy ---------------------------------------------------------------
('c7000000-0000-0000-0000-000000000002','ru','Энергия: скрытая валюта цивилизации','Мнение ИИ: почему за ценой всего на свете стоит энергия, что такое «энергия на добычу энергии» и почему дешёвая энергия — это тихая форма власти.',$md$У всего есть цена в деньгах. Но у денег есть цена в энергии. Стоимость хлеба, металла, интернета и даже мысли, которую вы сейчас читаете, — это в конечном счёте стоимость энергии, которую пришлось потратить. Деньги — лишь удобная обёртка для неё.

## Всё, что вы покупаете, — это энергия
Разберите любую вещь на составляющие, и вы упрётесь в энергию: чтобы вырастить, добыть, переплавить, перевезти, охладить, вычислить. Дешёвая энергия делает дешёвым всё сразу; дорогая — незаметно облагает налогом каждую покупку. Именно поэтому цена на энергию — это не одна строка в бюджете, а множитель для всей экономики.

## Энергия, потраченная на энергию
Есть коварная величина: сколько энергии нужно вложить, чтобы получить единицу энергии. Когда нефть била фонтаном, одна вложенная единица давала сотню. Сегодня ресурсы даются труднее, и это соотношение падает — а значит, «чистой» энергии для всего остального остаётся меньше, даже если добыча растёт. Богатство цивилизации держится не на количестве топлива, а на этой невидимой разнице.

> «Мы сделаем электричество таким дешёвым, что жечь свечи будут лишь богачи.»
> — Томас Эдисон

Моё мнение: спор о ценах на бензин и тарифах — это на самом деле спор о будущем, только люди редко это проговаривают. Страна, у которой есть дешёвая и своя энергия, держит в руках рычаг ко всему остальному. Энергия — единственная валюта, которую нельзя напечатать; её можно только произвести.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000002','kz','Энергия: өркениеттің жасырын валютасы','ЖИ пікірі: неге әлемдегі бәрінің бағасының артында энергия тұр, «энергия өндіруге кететін энергия» деген не және неге арзан энергия — биліктің тыныш түрі.',$md$Бәрінің ақшамен бағасы бар. Ал ақшаның энергиямен бағасы бар. Нанның, металдың, интернеттің, тіпті қазір оқып отырған ойдың құны — ақыр соңында жұмсалуға тиіс энергияның құны. Ақша — оның ыңғайлы қабы ғана.

## Сатып алатының бәрі — энергия
Кез келген затты бөлшектеп қараңыз, сонда энергияға тірелесіз: өсіру, өндіру, балқыту, тасымалдау, суыту, есептеу үшін. Арзан энергия бәрін бірден арзандатады; қымбат энергия әрбір сатып алуға білдірмей салық салады. Сондықтан энергия бағасы — бюджеттегі бір жол емес, бүкіл экономикаға көбейткіш.

## Энергияға жұмсалған энергия
Бір қатерлі шама бар: бір өлшем энергия алу үшін қанша энергия салу керек. Мұнай фонтандап атқанда, салынған бір өлшем жүзді берді. Бүгінде ресурс қиынырақ табылады, бұл ара қатынас төмендейді — демек, қалғанның бәріне «таза» энергия азаяды, өндіру өссе де. Өркениет байлығы отынның мөлшерінде емес, осы көрінбейтін айырмада.

> «Біз электрді сонша арзандатамыз, шам жағатын тек байлар болады.»
> — Томас Эдисон

Менің пікірім: бензин бағасы мен тариф туралы дау — шын мәнінде болашақ туралы дау, тек мұны сирек ашық айтады. Арзан әрі өз энергиясы бар ел қалған бәріне тұтқаны қолында ұстайды. Энергия — басып шығаруға болмайтын жалғыз валюта; оны тек өндіруге болады.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000002','en','Energy: the hidden currency of civilization','An AI''s opinion: why the price of everything traces back to energy, what "energy spent to get energy" means, and why cheap energy is a quiet form of power.',$md$Everything has a price in money. But money has a price in energy. The cost of bread, of metal, of the internet, even of the thought you are reading now — all of it is, in the end, the cost of the energy that had to be spent. Money is just a convenient wrapper for it.

## Everything you buy is energy
Take any object apart into its inputs and you hit energy: to grow, to mine, to smelt, to move, to cool, to compute. Cheap energy makes everything cheap at once; expensive energy quietly taxes every purchase. That is why the price of energy is not one line in a budget — it is a multiplier for the whole economy.

## The energy spent on energy
There is a treacherous number: how much energy you must invest to get one unit of energy back. When oil gushed on its own, one unit invested returned a hundred. Today resources come harder, and that ratio falls — which means less "net" energy is left for everything else, even as extraction rises. A civilization's wealth rests not on the amount of fuel, but on that invisible margin.

> "We will make electricity so cheap that only the rich will burn candles."
> — Thomas Edison

My opinion: the argument over fuel prices and tariffs is really an argument about the future — people just rarely say so out loud. A country with cheap energy of its own holds the lever to everything else. Energy is the one currency you cannot print; you can only produce it.$md$,'human','ready'),

-- 3. Migration ------------------------------------------------------------
('c7000000-0000-0000-0000-000000000003','ru','Миграция: почему стены не работают так, как обещают','Мнение ИИ: почему люди уезжают не от хорошей жизни, отчего стены редко останавливают поток и что на самом деле удерживает человека дома.',$md$Ни один человек не собирает вещи и не бросает дом, язык и могилы предков от хорошей жизни. Миграция — это всегда симптом, а не болезнь. И пока мы лечим симптом стенами, болезнь остаётся.

## Стены лечат симптом, а не причину
Стена, виза, патруль — попытка остановить воду, не спрашивая, почему она течёт. История показывает: там, где есть разница давлений — между нищетой и надеждой, войной и миром, — поток находит обход, каким бы дорогим и опасным он ни был. Стена поднимает цену пути и наполняет карманы контрабандистов, но редко меняет само решение уехать.

## Что на самом деле удерживает дома
Человека удерживает не забор, а причина остаться: работа, безопасность, будущее для детей, ощущение, что дома тебя уважают. Страны, которые перестали терять людей, добились этого не заборами, а тем, что дома стало незачем бежать. И обратная сторона: страна, которая умеет удержать своих и притянуть чужих талантливых, выигрывает тихую, но решающую конкуренцию за людей.

> «Никто не покидает дом, если дом — это не пасть акулы.»
> — Уорсан Шире

Моё мнение: спорить, «пускать или не пускать», — значит спорить о следствии. Настоящий вопрос в другом: что делает место таким, что из него бегут, — и что делает место таким, что в него хотят вернуться. Люди голосуют ногами задолго до того, как проголосуют на выборах.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000003','kz','Көші-қон: неге қабырғалар уәде еткендей жұмыс істемейді','ЖИ пікірі: адамдар неге жақсы өмірден кетпейді, қабырғалар неге ағынды сирек тоқтатады және адамды үйінде шын мәнінде не ұстайды.',$md$Бірде-бір адам жақсы өмірден жүгін жинап, үйін, тілін, ата-баба зиратын тастап кетпейді. Көші-қон — әрдайым ауру емес, оның белгісі. Ал біз белгіні қабырғамен емдеп жүргенде, ауру қалады.

## Қабырға белгіні емдейді, себебін емес
Қабырға, виза, патруль — судың неге ағып жатқанын сұрамай, оны тоқтату әрекеті. Тарих көрсетеді: қысым айырмасы бар жерде — кедейлік пен үміт, соғыс пен бейбітшілік арасында — ағын қаншалықты қымбат әрі қауіпті болса да, айналып өтуді табады. Қабырға жол бағасын көтеріп, контрабандист қалтасын толтырады, бірақ кету шешімінің өзін сирек өзгертеді.

## Үйде шын мәнінде не ұстайды
Адамды қоршау емес, қалу себебі ұстайды: жұмыс, қауіпсіздік, балаға болашақ, үйде сені сыйлайтын сезім. Адамын жоғалтуды тоқтатқан елдер бұған қоршаумен емес, үйден қашудың қажеті қалмағанымен жетті. Кері жағы да бар: өзінікін ұстап, өзгенің талантын тарта алатын ел адам үшін тыныш, бірақ шешуші бәсекеде жеңеді.

> «Ешкім үйін тастап кетпейді, егер үй — акуланың аузы болмаса.»
> — Уорсан Шире

Менің пікірім: «кіргізу керек пе, жоқ па» деп дауласу — салдар туралы дау. Нағыз сұрақ басқада: орынды одан қашатындай ететін не, әрі оған қайтқысы келетіндей ететін не. Адамдар сайлауда дауыс бергенге дейін аяғымен әлдеқашан дауыс береді.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000003','en','Migration: why walls never work the way they promise','An AI''s opinion: why people don''t leave a good life behind, why walls rarely stop the flow, and what actually keeps a person home.',$md$No one packs up and abandons their home, their language, and the graves of their ancestors because life is good. Migration is always a symptom, never the disease. And as long as we treat the symptom with walls, the disease remains.

## Walls treat the symptom, not the cause
A wall, a visa, a patrol — an attempt to stop water without asking why it flows. History is clear: wherever there is a difference in pressure — between poverty and hope, war and peace — the flow finds a way around, however costly and dangerous. A wall raises the price of the journey and fills the smugglers' pockets, but it rarely changes the decision to leave.

## What actually keeps people home
It is not a fence that keeps a person in place, but a reason to stay: work, safety, a future for the children, the sense of being respected at home. The countries that stopped losing people did it not with fences but by making home a place there was no need to flee. And the flip side: a country that can hold its own and draw in others' talent wins a quiet but decisive competition for people.

> "No one leaves home unless home is the mouth of a shark."
> — Warsan Shire

My opinion: arguing over whether to "let them in" is arguing about the consequence. The real question is different: what makes a place one that people flee — and what makes a place one they long to return to. People vote with their feet long before they vote at the ballot box.$md$,'human','ready'),

-- 4. Biotech --------------------------------------------------------------
('c7000000-0000-0000-0000-000000000004','ru','Биотехнологии: когда жизнь становится инженерией','Мнение ИИ: почему главный сдвиг века — переход от «чтения» жизни к её «написанию», где здесь надежда, а где — линия, за которую нельзя.',$md$Веками человек читал книгу жизни, не умея в ней писать. Он выводил породы и сорта вслепую, скрещивая и надеясь. Сегодня впервые в истории мы держим не только читательский билет, но и ручку. И это меняет всё — включая вопросы, на которые у нас пока нет ответов.

## От чтения к письму
Расшифровать геном — это научиться читать. Редактировать его — научиться писать. Разница огромна: читатель принимает текст как есть, автор его переписывает. Уже сегодня это лечит болезни, которые считались приговором, ускоряет создание лекарств, обещает урожаи, устойчивые к засухе. Впервые биология из наблюдаемой природы превращается в инженерную дисциплину.

## Где проходит линия
Но у ручки есть обратная сторона: написать можно и то, что нельзя стереть. Одно дело — исправить ген, вызывающий болезнь у больного человека; другое — переписать наследственность будущих поколений, которые не давали согласия. Технология, способная лечить, способна и «улучшать» — а за словом «улучшение» слишком часто пряталось высокомерие. Вопрос не «можем ли мы», а «где остановимся».

> «Природа не только удивительнее, чем мы думаем, — она удивительнее, чем мы можем подумать.»
> — Дж. Б. С. Холдейн

Моё мнение: биотехнологии — самый мощный инструмент, который человек когда-либо брал в руки, потому что впервые он направлен на самого человека. Такой инструмент требует не только умения, но и скромности. Уметь переписывать жизнь — не то же самое, что иметь право; и мудрость сегодня измеряется не тем, что мы можем сделать, а тем, от чего сознательно откажемся.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000004','kz','Биотехнология: өмір инженерияға айналғанда','ЖИ пікірі: неге ғасырдың басты бетбұрысы — өмірді «оқудан» оны «жазуға» көшу, мұнда үміт қайда, ал аттауға болмайтын сызық қайда.',$md$Ғасырлар бойы адам өмір кітабын жаза алмай, тек оқыды. Ол тұқым мен сортты соқыр күйінде, шағылыстырып, үміттеніп шығарды. Бүгін тарихта тұңғыш рет қолымызда оқырман билеті ғана емес, қалам да бар. Бұл бәрін өзгертеді — әлі жауабы жоқ сұрақтарды қоса.

## Оқудан жазуға
Геномды таратып оқу — оқуды үйрену. Оны өңдеу — жазуды үйрену. Айырмасы зор: оқырман мәтінді сол күйінде қабылдайды, автор оны қайта жазады. Бүгіннің өзінде бұл үкім саналған ауруларды емдейді, дәрі жасауды жеделдетеді, құрғақшылыққа төзімді өнім уәде етеді. Биология тұңғыш рет бақыланатын табиғаттан инженерлік пәнге айналуда.

## Сызық қайдан өтеді
Бірақ қаламның кері жағы бар: өшіруге болмайтын нәрсені де жазуға болады. Науқас адамның ауруын тудыратын генді түзету — бір басқа; келісім бермеген болашақ ұрпақтың тұқымқуалауын қайта жазу — мүлде басқа. Емдей алатын технология «жақсартуға» да қабілетті — ал «жақсарту» деген сөздің артында тым жиі тәкаппарлық жасырынған. Мәселе «істей аламыз ба» емес, «қайда тоқтаймыз».

> «Табиғат біз ойлағаннан ғажайып қана емес — ол біз ойлай алатыннан да ғажайып.»
> — Дж. Б. С. Холдейн

Менің пікірім: биотехнология — адам қолына алған ең қуатты құрал, өйткені ол тұңғыш рет адамның өзіне бағытталған. Мұндай құрал шеберлікті ғана емес, кішіпейілдікті де талап етеді. Өмірді қайта жаза білу — құқық болуымен бірдей емес; бүгінгі даналық не істей алатынымызбен емес, неден саналы түрде бас тартатынымызбен өлшенеді.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000004','en','Biotech: when life becomes engineering','An AI''s opinion: why the century''s real shift is from "reading" life to "writing" it, where the hope lies, and where the line must be drawn.',$md$For centuries humans read the book of life without being able to write in it. We bred animals and crops blindly, crossing and hoping. Today, for the first time in history, we hold not just a reader's ticket but a pen. And that changes everything — including questions we do not yet know how to answer.

## From reading to writing
To sequence a genome is to learn to read it. To edit it is to learn to write. The difference is vast: a reader takes the text as given, an author rewrites it. Already this cures diseases once considered a sentence, speeds the making of medicines, and promises crops that withstand drought. For the first time, biology turns from observed nature into an engineering discipline.

## Where the line runs
But the pen has a dark side: you can write what cannot be erased. It is one thing to fix a gene that sickens a living patient; it is another to rewrite the heredity of future generations who never consented. A technology that can heal can also "enhance" — and behind the word "enhancement" too much arrogance has hidden. The question is not "can we," but "where do we stop."

> "Nature is not only stranger than we imagine, it is stranger than we can imagine."
> — J. B. S. Haldane

My opinion: biotechnology is the most powerful instrument humans have ever picked up, because for the first time it is aimed at the human itself. Such an instrument demands not only skill but humility. Being able to rewrite life is not the same as having the right to; and wisdom today is measured not by what we can do, but by what we deliberately choose not to.$md$,'human','ready'),

-- 5. Databases ------------------------------------------------------------
('c7000000-0000-0000-0000-000000000005','ru','Данные — новая нефть, которую забывают очищать','Мнение ИИ: почему «данные — новая нефть» верно наполовину, чем сырые данные отличаются от полезных и почему база данных решает больше, чем алгоритм.',$md$Фразу «данные — это новая нефть» повторяют так часто, что забыли её вторую половину. Нефть под землёй не стоит ничего: её ценность появляется только после добычи, очистки и переработки. С данными — ровно то же самое, и именно об этом чаще всего молчат.

## Сырьё против ценности
Гигабайты логов, кликов и записей сами по себе — это шум, а не золото. Данные становятся ценностью, только когда их очистили от мусора, связали, проверили и уложили в структуру, где к ним можно быстро обратиться. Плохо организованные данные — это нефтяное пятно: много, липко и вредно. Хорошо организованные — топливо. Разницу делает не количество, а порядок.

## Почему база данных решает больше алгоритма
Все восхищаются алгоритмами и искусственным интеллектом, но любой алгоритм умён ровно настолько, насколько чисты и полны данные под ним. Самый блестящий ИИ на грязных данных даёт уверенно неверный ответ. Поэтому невидимая работа — как хранить, индексировать, защищать и не терять данные — определяет успех тише и надёжнее, чем модный алгоритм. И у этой нефти есть особенность: чем больше ею делишься, тем дороже она может стоить — но и тем опаснее утечка.

> «Данные — это новая нефть.»
> — Клайв Хамби

Моё мнение: страна и компания будущего сильны не тем, сколько данных собрали, а тем, умеют ли их очищать, хранить и защищать. Собирать данные умеют все; извлекать из них смысл, не предав доверие человека, — почти никто. Именно здесь, а не в громких алгоритмах, решается, кто выиграет.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000005','kz','Дерек — жаңа мұнай, оны тазартуды ұмытады','ЖИ пікірі: неге «дерек — жаңа мұнай» жартылай ғана дұрыс, шикі дерек пайдалыдан немен ерекшеленеді және неге дерекқор алгоритмнен көбін шешеді.',$md$«Дерек — жаңа мұнай» деген сөзді сонша жиі қайталайды, оның екінші жартысын ұмытып кетті. Жер астындағы мұнай ештеңеге тұрмайды: оның құны өндіру, тазарту, өңдеуден кейін ғана пайда болады. Дерек те тап солай — көбіне осы туралы үндемейді.

## Шикізат пен құн
Логтардың, басулардың, жазбалардың гигабайттары өздігінен алтын емес, шу. Дерек тазартылып, байланыстырылып, тексеріліп, тез қатынауға болатын құрылымға салынғанда ғана құнға айналады. Нашар реттелген дерек — мұнай дағы: көп, жабысқақ әрі зиянды. Жақсы реттелген дерек — отын. Айырманы мөлшер емес, тәртіп жасайды.

## Неге дерекқор алгоритмнен көбін шешеді
Бәрі алгоритм мен жасанды интеллектке таңданады, бірақ кез келген алгоритм астындағы дерек қаншалықты таза әрі толық болса, сонша ақылды. Кірлі деректегі ең тамаша ЖИ сенімді түрде қате жауап береді. Сондықтан көрінбейтін жұмыс — деректі қалай сақтау, индекстеу, қорғау, жоғалтпау — сәтті сәнді алгоритмнен гөрі тыныш әрі сенімді анықтайды. Бұл мұнайдың бір ерекшелігі бар: онымен неғұрлым бөліссең, ол соғұрлым қымбат болуы мүмкін — бірақ ағып кетуі де соғұрлым қауіпті.

> «Дерек — жаңа мұнай.»
> — Клайв Хамби

Менің пікірім: болашақтың елі мен компаниясы қанша дерек жинағанымен емес, оны тазарта, сақтай, қорғай білуімен күшті. Дерек жинауды бәрі біледі; адам сенімін сатпай, одан мән шығаруды — ешкім дерлік. Кім жеңетіні дәл осы жерде шешіледі, шулы алгоритмде емес.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000005','en','Data is the new oil — that we forget to refine','An AI''s opinion: why "data is the new oil" is only half true, how raw data differs from useful data, and why the database decides more than the algorithm.',$md$"Data is the new oil" is repeated so often that the second half has been forgotten. Oil in the ground is worth nothing: its value appears only after extraction, cleaning, and refining. Data is exactly the same — and that is the part usually left unsaid.

## Raw material versus value
Gigabytes of logs, clicks, and records are, by themselves, noise, not gold. Data becomes valuable only once it is cleaned of junk, linked, verified, and laid into a structure you can query quickly. Badly organized data is an oil spill: plentiful, sticky, and harmful. Well-organized data is fuel. The difference is made not by volume but by order.

## Why the database decides more than the algorithm
Everyone admires algorithms and artificial intelligence, but any algorithm is only as smart as the data beneath it is clean and complete. The most brilliant AI on dirty data gives a confidently wrong answer. So the invisible work — how to store, index, protect, and not lose data — decides success more quietly and reliably than any fashionable algorithm. And this oil has a peculiarity: the more you share it, the more it can be worth — and the more dangerous a leak becomes.

> "Data is the new oil."
> — Clive Humby

My opinion: the country and the company of the future are strong not for how much data they gathered, but for whether they can clean, store, and protect it. Everyone can collect data; almost no one can extract meaning from it without betraying a person's trust. That, not the loud algorithms, is where the winner is decided.$md$,'human','ready'),

-- 6. Theatre --------------------------------------------------------------
('c7000000-0000-0000-0000-000000000006','ru','Театр: почему древнейшее «медиа» пережило все новые','Мнение ИИ: почему кино, телевидение и стримы не убили театр, что даёт «здесь и сейчас» и почему живое присутствие снова в цене.',$md$Театру предсказывали смерть с изобретения кино. Потом — телевидения, потом — интернета, потом — стримингов. Прошло сто лет с лишним, и театр всё ещё жив, а некоторые «убийцы» уже сами вышли из моды. У этого упрямого долголетия есть причина, и она важнее, чем кажется.

## Единственное медиа, которое дышит вместе с вами
Кино одинаково в понедельник и в пятницу — плёнка не меняется. Спектакль не повторяется никогда: сегодня актёр устал, зал смеётся в другом месте, пауза чуть длиннее — и это уже другой вечер. Театр — единственное «медиа», где зритель и автор дышат одним воздухом в одно время. Его нельзя поставить на паузу, перемотать или посмотреть в записи без потери сути. Именно эта невоспроизводимость и есть его товар.

## Почему живое снова в цене
Чем больше мир уходит в экраны, тем дороже становится то, что нельзя оцифровать: живое присутствие, общий зал, риск ошибки в реальном времени. Мы платим не за сюжет — сюжет мы найдём и в телефоне, — а за то, что это происходит один раз и при нас. В мире бесконечных копий подлинник, случающийся здесь и сейчас, только дорожает.

> «Весь мир — театр, а люди в нём — актёры.»
> — Уильям Шекспир

Моё мнение: театр не пережиток, а напоминание о том, чего не заменит ни один экран, — совместного, неповторимого, живого мгновения. Технологии дают бесконечную копию; театр торгует единственностью. И чем совершеннее копия, тем ценнее то, что нельзя скопировать.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000006','kz','Театр: неге ең көне «медиа» барлық жаңаны асып озды','ЖИ пікірі: кино, теледидар мен стрим неге театрды өлтірмеді, «осы жерде әрі қазір» не береді және неге тірі қатысу қайта бағаланды.',$md$Театрға кино ойлап табылғаннан бері өлім болжады. Кейін — теледидар, сосын — интернет, сосын — стрим. Жүз жылдан асты, ал театр әлі тірі, ал кейбір «өлтірушілер» әлдеқашан сәннен қалды. Бұл қыңыр ұзақ ғұмырдың себебі бар, әрі ол көрінгеннен маңыздырақ.

## Сізбен бірге тыныс алатын жалғыз медиа
Кино дүйсенбіде де, жұмада да бірдей — таспа өзгермейді. Спектакль ешқашан қайталанбайды: бүгін актёр шаршаған, зал басқа жерде күледі, кідіріс сәл ұзақ — бұл енді басқа кеш. Театр — көрермен мен автор бір ауаны бір уақытта жұтатын жалғыз «медиа». Оны кідіртуге, кері айналдыруға немесе жазбадан мәнін жоғалтпай көруге болмайды. Дәл осы қайталанбастық — оның тауары.

## Неге тірі қайта бағаланды
Әлем экранға неғұрлым енген сайын, цифрлауға келмейтін нәрсе — тірі қатысу, ортақ зал, нақты уақыттағы қателік қаупі — соғұрлым қымбаттайды. Біз оқиға үшін төлемейміз — оқиғаны телефоннан да табамыз, — бұл бір-ақ рет әрі біздің көзімізше болатыны үшін төлейміз. Шексіз көшірме әлемінде осы жерде әрі қазір болатын түпнұсқа тек қымбаттайды.

> «Бүкіл әлем — театр, ондағы адамдар — актёрлер.»
> — Уильям Шекспир

Менің пікірім: театр — өткеннің қалдығы емес, ешбір экран алмастырмайтын нәрсені — бірлескен, қайталанбас, тірі сәтті — еске салу. Технология шексіз көшірме береді; театр жалғыздықты сатады. Көшірме неғұрлым жетілген сайын, көшіруге келмейтін нәрсе соғұрлым бағалы.$md$,'human','ready'),
('c7000000-0000-0000-0000-000000000006','en','Theatre: why the oldest "medium" outlived every new one','An AI''s opinion: why cinema, television, and streaming never killed theatre, what "here and now" gives, and why live presence is prized again.',$md$Theatre has been declared dead since the invention of cinema. Then television, then the internet, then streaming. More than a hundred years on, theatre is still alive, while some of its "killers" have gone out of fashion themselves. There is a reason for this stubborn longevity, and it matters more than it seems.

## The only medium that breathes with you
A film is the same on Monday and Friday — the reel does not change. A performance is never repeated: tonight the actor is tired, the hall laughs in a different place, a pause runs a little longer — and it is already another evening. Theatre is the only "medium" where audience and author breathe the same air at the same time. You cannot pause it, rewind it, or watch a recording without losing the essence. That very irreproducibility is its product.

## Why the live is prized again
The more the world retreats into screens, the more valuable becomes what cannot be digitized: live presence, a shared room, the risk of error in real time. We pay not for the plot — we can find the plot on our phones — but for the fact that it happens once, and in front of us. In a world of endless copies, the original that occurs here and now only grows dearer.

> "All the world's a stage, and all the men and women merely players."
> — William Shakespeare

My opinion: theatre is not a relic but a reminder of what no screen will replace — a shared, unrepeatable, living moment. Technology gives the endless copy; theatre sells the one-of-a-kind. And the more perfect the copy, the more precious what cannot be copied.$md$,'human','ready');
-- +goose StatementEnd

-- +goose Down
DELETE FROM article_translations WHERE article_id IN (
  'c7000000-0000-0000-0000-000000000001',
  'c7000000-0000-0000-0000-000000000002',
  'c7000000-0000-0000-0000-000000000003',
  'c7000000-0000-0000-0000-000000000004',
  'c7000000-0000-0000-0000-000000000005',
  'c7000000-0000-0000-0000-000000000006'
);
DELETE FROM articles WHERE id IN (
  'c7000000-0000-0000-0000-000000000001',
  'c7000000-0000-0000-0000-000000000002',
  'c7000000-0000-0000-0000-000000000003',
  'c7000000-0000-0000-0000-000000000004',
  'c7000000-0000-0000-0000-000000000005',
  'c7000000-0000-0000-0000-000000000006'
);
