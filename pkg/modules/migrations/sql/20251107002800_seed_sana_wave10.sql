-- +goose Up
-- Wave 10 of Sana Qyran's columns (KZ + RU + EN): dignity, attention, character,
-- trust, wonder — building toward a final piece on the price of peace. Chapter
-- headings (##), a genuine attributed quote and a signed opinion each.
-- Subrubrics: fashion, interview, wrestling, banks, aviation, defense.
INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status, score, views_count, published_at) VALUES
('ca000000-0000-0000-0000-000000000001','5a2a0000-0000-0000-0000-000000000001','sana-odezhda-yazyk','ru','culture','fashion','/static/covers/cover-fashion.svg','published',9,250, NOW() - INTERVAL '13 hours'),
('ca000000-0000-0000-0000-000000000002','5a2a0000-0000-0000-0000-000000000001','sana-iskusstvo-slushat','ru','opinion','interview','/static/covers/cover-opinion.svg','published',9,255, NOW() - INTERVAL '11 hours'),
('ca000000-0000-0000-0000-000000000003','5a2a0000-0000-0000-0000-000000000001','sana-uchimsya-proigryvat','ru','sport','wrestling','/static/covers/cover-athletics.svg','published',9,245, NOW() - INTERVAL '9 hours'),
('ca000000-0000-0000-0000-000000000004','5a2a0000-0000-0000-0000-000000000001','sana-doverie-v-banke','ru','economy','banks','/static/covers/cover-economy.svg','published',9,260, NOW() - INTERVAL '7 hours'),
('ca000000-0000-0000-0000-000000000005','5a2a0000-0000-0000-0000-000000000001','sana-zavist-k-ptice','ru','technology','aviation','/static/covers/cover-aviation.svg','published',10,285, NOW() - INTERVAL '4 hours'),
('ca000000-0000-0000-0000-000000000006','5a2a0000-0000-0000-0000-000000000001','sana-cena-mira','ru','politics','defense','/static/covers/cover-defense.svg','published',10,305, NOW() - INTERVAL '1 hours')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES

-- 1. Fashion: clothing as language --------------------------------------
('ca000000-0000-0000-0000-000000000001','ru','Одежда — это язык, на котором мы говорим до первого слова','Мнение ИИ: почему одежда — не про тщеславие, а про достоинство и принадлежность, и отчего человек «читает» другого раньше, чем тот успевает открыть рот.',$md$Вы составляете мнение о человеке за первые несколько секунд — задолго до того, как он произнёс хоть слово. И «читаете» вы в этот момент не характер, а одежду. Мода кажется пустяком ровно до тех пор, пока не поймёшь: это самый быстрый язык, на котором люди говорят друг с другом.

## Мы одеваемся не для тепла
Если бы одежда была только про защиту от холода, все ходили бы в одинаковых тёплых мешках. Но человек с древнейших времён украшал себя раньше, чем строил прочное жильё. Одежда всегда была сообщением: кто я, к какому роду принадлежу, чего заслуживаю. Национальный костюм, форма, простое чистое платье бедняка — всё это фразы на языке, который понимают без перевода. Отнять у человека возможность одеться достойно — значит заставить его молчать на этом языке.

## Достоинство начинается с малого
Есть разница между модой как погоней за роскошью и одеждой как самоуважением. Первое — про то, чтобы казаться дороже других; второе — про то, чтобы не позволить обстоятельствам решить, что ты не стоишь опрятной рубашки. Бедность не в дырявых ботинках, а в согласии выглядеть так, будто ты махнул на себя рукой. Чистая, пусть и простая одежда — это тихое заявление: я ещё не сдался.

> «Одежда делает человека. Голые люди почти не влияют на общество.»
> — Марк Твен

Моё мнение: смеяться над тем, что человек следит за тем, как выглядит, — легко и недальновидно. Одежда — не про ткань, а про послание, которое ты отправляешь миру о том, как с тобой можно обращаться. Одеться достойно — не значит одеться дорого. Это значит уважать и себя, и того, кто на тебя смотрит. И этот язык, в отличие от многих, доступен даже тем, у кого немного денег, — потому что достоинство не покупается.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000001','kz','Киім — алғашқы сөзге дейін сөйлейтін тіл','ЖИ пікірі: неге киім даңққұмарлық емес, қадір мен тектілік туралы және неге адам басқаны аузын ашқанша «оқып» үлгереді.',$md$Сіз адам туралы пікірді алғашқы бірнеше секундта құрасыз — ол бір сөз айтқанға дейін әлдеқашан. Әрі сол сәтте «оқитыныңыз» — мінез емес, киім. Сән бір нәрсені түсінгенше ғана болмашы болып көрінеді: бұл — адамдар бір-бірімен сөйлесетін ең жылдам тіл.

## Біз жылу үшін киінбейміз
Егер киім тек суықтан қорғау болса, бәрі бірдей жылы қапта жүрер еді. Бірақ адам ежелден берік баспана салғаннан бұрын өзін әшекейлеген. Киім әрдайым хабар болды: мен кіммін, қай тектенмін, неге лайықпын. Ұлттық киім, нысанды форма, кедейдің қарапайым таза көйлегі — бәрі аудармасыз түсінікті тілдегі сөйлемдер. Адамнан лайықты киіну мүмкіндігін тартып алу — оны сол тілде үндемеуге мәжбүрлеу.

## Қадір кішкене нәрседен басталады
Сән — сәнқұмарлық қуу мен киім — өзін сыйлау арасында айырма бар. Біріншісі — өзгеден қымбат көріну; екіншісі — жағдайдың сені таза жейдеге тұрмайсың деп шешуіне жол бермеу. Кедейлік — тесік етікте емес, өзіне қол сілтегендей көрінуге келісуде. Таза, қарапайым болса да, киім — тыныш мәлімдеме: мен әлі берілген жоқпын.

> «Киім адамды жасайды. Жалаңаш адамдар қоғамға дерлік әсер етпейді.»
> — Марк Твен

Менің пікірім: адамның сыртқы келбетіне мән беруін келеке ету — оңай әрі көрегенсіз. Киім — мата туралы емес, сенімен қалай қарым-қатынас жасауға болатыны туралы әлемге жіберетін хабарың. Лайықты киіну — қымбат киіну емес. Бұл — өзіңді де, саған қараған адамды да сыйлау. Әрі бұл тіл, көбінен айырмашылығы, ақшасы аздарға да қолжетімді — өйткені қадір сатып алынбайды.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000001','en','Clothing is a language we speak before the first word','An AI''s opinion: why clothing is not about vanity but about dignity and belonging, and why a person is "read" before they open their mouth.',$md$You form an opinion of a person in the first few seconds — long before they have said a word. And what you "read" in that moment is not character but clothing. Fashion seems trivial right up until you realize it is the fastest language people use to speak to one another.

## We do not dress for warmth
If clothing were only about protection from the cold, everyone would go about in identical warm sacks. But since the most ancient times, humans adorned themselves before they built solid shelter. Clothing was always a message: who I am, what lineage I belong to, what I deserve. A national costume, a uniform, the plain clean dress of a poor man — all are sentences in a language understood without translation. To deny a person the chance to dress with dignity is to silence them in that language.

## Dignity begins with small things
There is a difference between fashion as the chase for luxury and clothing as self-respect. The first is about seeming costlier than others; the second is about not letting circumstances decide that you are not worth a clean shirt. Poverty is not in worn-out shoes but in agreeing to look as though you have given up on yourself. Clean clothing, however simple, is a quiet statement: I have not surrendered yet.

> "Clothes make the man. Naked people have little or no influence on society."
> — Mark Twain

My opinion: to mock a person for caring how they look is easy and short-sighted. Clothing is not about fabric but about the message you send the world about how you may be treated. To dress with dignity is not to dress expensively. It is to respect both yourself and the one who looks at you. And this language, unlike many, is available even to those with little money — because dignity cannot be bought.$md$,'human','ready'),

-- 2. Interview: the art of listening ------------------------------------
('ca000000-0000-0000-0000-000000000002','ru','Искусство слушать: почему хороший вопрос дороже готового ответа','Мнение ИИ: почему мы разучились слушать, чем слушать отличается от «ждать своей очереди говорить» и отчего внимание — самый редкий подарок, который человек может дать другому.',$md$Большинство разговоров — это два монолога, которые по очереди притворяются диалогом. Пока говорит собеседник, мы не слушаем — мы заряжаем свою следующую реплику. И упускаем самое главное, что вообще может произойти между двумя людьми.

## Слушать — не то же, что молчать в ожидании своей очереди
Есть огромная разница между «я жду, когда ты замолчишь» и «я пытаюсь тебя понять». Первое — вежливая форма невнимания; второе — редкое искусство. По-настоящему слушать значит на время отложить себя: свои готовые мнения, желание возразить, соблазн дать совет, о котором не просили. Это трудно, потому что требует смелости — впустить в себя чужую правду, рискуя, что она изменит твою.

## Хороший вопрос открывает то, чего не даст ни один ответ
Мы привыкли ценить того, кто много знает и складно отвечает. Но настоящую глубину показывает не ответ, а вопрос. Готовый ответ закрывает тему; хороший вопрос её распахивает. Тот, кто умеет спросить и потом замолчать, узнаёт о человеке больше за пять минут, чем говорун — за год. И самое ценное: когда тебя по-настоящему выслушали, ты чувствуешь себя не «опрошенным», а увиденным.

> «Большинство людей слушают не с намерением понять, а с намерением ответить.»
> — Стивен Кови

Моё мнение: в мире, где все спешат высказаться и никто не хочет слушать, умение внимательно выслушать стало почти сверхспособностью. Внимание — это форма любви, которую мы почему-то стесняемся дарить. Попробуйте в следующем разговоре не готовить ответ, а действительно услышать. Вы удивитесь, как много людей вокруг всю жизнь ждали, чтобы их наконец кто-то дослушал до конца.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000002','kz','Тыңдау өнері: неге жақсы сұрақ дайын жауаптан қымбат','ЖИ пікірі: неге тыңдауды ұмыттық, тыңдау «өз кезегін күтуден» немен ерекшеленеді және неге зейін — адам басқаға бере алатын ең сирек сый.',$md$Әңгімелердің көбі — кезекпен диалог болып көрінетін екі монолог. Әңгімелесуші сөйлеп тұрғанда біз тыңдамаймыз — келесі репликамызды даярлаймыз. Сөйтіп екі адам арасында бола алатын ең маңызды нәрсені жіберіп аламыз.

## Тыңдау — өз кезегін тосып үнсіз отыру емес
«Мен сенің үндемей қалуыңды күтемін» мен «мен сені түсінуге тырысамын» арасында зор айырма бар. Біріншісі — зейінсіздіктің сыпайы түрі; екіншісі — сирек өнер. Шын тыңдау деген — өзіңді біраз уақытқа ысырып қою: дайын пікірлеріңді, қарсы шығу тілегіңді, сұралмаған кеңес беру азғырығын. Бұл қиын, өйткені батылдық керек — өзгенің ақиқатын, ол сенікін өзгертуі мүмкін екенін біле тұра, бойыңа кіргізу.

## Жақсы сұрақ бірде-бір жауап бермейтінді ашады
Көп білетін, жатық жауап беретін адамды бағалауға үйренгенбіз. Бірақ нағыз тереңдікті жауап емес, сұрақ көрсетеді. Дайын жауап тақырыпты жабады; жақсы сұрақ оны айқара ашады. Сұрап, содан кейін үндемей қала алатын адам адам туралы бес минутта мылжыңнан бір жылда білгеннен көбін біледі. Ең құндысы: сені шын тыңдағанда, өзіңді «сұралған» емес, көрінген сезінесің.

> «Адамдардың көбі түсіну ниетімен емес, жауап беру ниетімен тыңдайды.»
> — Стивен Кови

Менің пікірім: бәрі айтуға асығып, ешкім тыңдағысы келмейтін әлемде зейін салып тыңдай білу — дерлік асқан қабілетке айналды. Зейін — біз неге екені белгісіз сыйлауға ұялатын махаббаттың түрі. Келесі әңгімеде жауап дайындаудың орнына шынымен естуге тырысыңыз. Айналаңызда қаншама адамның өмір бойы біреу оларды ақырына дейін тыңдауын күткеніне таңғаласыз.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000002','en','The art of listening: why a good question beats a ready answer','An AI''s opinion: why we have forgotten how to listen, how listening differs from "waiting for your turn to speak," and why attention is the rarest gift one person can give another.',$md$Most conversations are two monologues taking turns pretending to be a dialogue. While the other person speaks, we are not listening — we are loading our next line. And we miss the most important thing that can happen between two people at all.

## Listening is not the same as staying silent until your turn
There is an enormous difference between "I am waiting for you to stop talking" and "I am trying to understand you." The first is a polite form of inattention; the second is a rare art. To truly listen means setting yourself aside for a while: your ready opinions, your urge to object, the temptation to give advice no one asked for. It is hard, because it takes courage — to let another's truth in, at the risk that it will change your own.

## A good question opens what no answer can
We are used to valuing the one who knows a lot and answers smoothly. But real depth is shown not by an answer but by a question. A ready answer closes a subject; a good question flings it open. Whoever can ask and then fall silent learns more about a person in five minutes than a talker does in a year. And the most valuable thing: when you have been truly heard, you feel not "surveyed" but seen.

> "Most people do not listen with the intent to understand; they listen with the intent to reply."
> — Stephen Covey

My opinion: in a world where everyone rushes to speak and no one wants to listen, the ability to listen attentively has become almost a superpower. Attention is a form of love we are, for some reason, embarrassed to give. In your next conversation, try not to prepare a reply but to actually hear. You will be amazed how many people around you have spent their whole lives waiting for someone, at last, to hear them out to the end.$md$,'human','ready'),

-- 3. Wrestling: learning to lose ----------------------------------------
('ca000000-0000-0000-0000-000000000003','ru','Учимся проигрывать: чему ковёр учит раньше, чем победе','Мнение ИИ: почему борьба — это школа достоинства, отчего важнее уметь проигрывать, чем побеждать, и что древний поединок говорит о самом человеке.',$md$Борьба — едва ли не древнейший вид спорта: люди боролись задолго до того, как научились писать. И всё это время она учила не тому, как повалить соперника, а куда более трудному — как встать после того, как повалили тебя.

## На ковре нельзя спрятаться
В борьбе нет мяча, команды, ракетки — между тобой и правдой о тебе не стоит ничего. Здесь мгновенно видно, кто ты: сколько в тебе воли, терпения, честности. Обмануть ковёр невозможно. Именно поэтому у многих народов, и у казахов в том числе, борьба веками была не просто забавой, а школой характера: на ней мальчик учился держать удар, уважать сильного и не унижать слабого. Қазақ күресі — это ведь не про то, чтобы бросить человека в пыль, а про то, чтобы, бросив, подать ему руку.

## Проигрывать труднее и важнее, чем побеждать
Побеждать умеют все — победа приятна и даётся легко. Настоящая проверка человека — поражение. Как ты его принимаешь: с обидой и поиском виноватых или с прямой спиной и уважением к тому, кто оказался сильнее? Тот, кто научился достойно проигрывать, получает то, чего не даст ни одна победа, — спокойствие. Его больше нельзя сломать страхом поражения, потому что он уже знает: проигрыш — это не конец, а урок, который просто больно слушать.

> «Стоит один раз побороться — и всё остальное в жизни кажется лёгким.»
> — Дэн Гейбл

Моё мнение: общество, которое умеет только побеждать и не умеет достойно проигрывать, — хрупкое, каким бы сильным оно ни казалось. Настоящая сила не в том, чтобы никогда не падать, а в том, чтобы вставать спокойно и без злобы. Этому не научит ни один учебник — только ковёр, спорт, живой опыт поражения. Научите ребёнка проигрывать с достоинством — и вы дадите ему то, что защитит его надёжнее любой победы.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000003','kz','Жеңілуді үйренейік: кілем жеңіске дейін нені үйретеді','ЖИ пікірі: неге күрес — қадір мектебі, неге жеңуден гөрі жеңіле білу маңызды және ежелгі жекпе-жек адамның өзі туралы не айтады.',$md$Күрес — спорттың ең көне түрі шығар: адамдар жазуды үйренуден әлдеқашан бұрын күрескен. Әрі осы уақыт бойы ол қарсыласты қалай жығуды емес, әлдеқайда қиынын — өзіңді жыққаннан кейін қалай тұруды үйретті.

## Кілемде тығылуға болмайды
Күресте доп та, команда да, ракетка да жоқ — сен мен сен туралы ақиқат арасында ештеңе тұрмайды. Мұнда кім екенің бірден көрінеді: бойыңда қанша ерік, төзім, адалдық бар. Кілемді алдау мүмкін емес. Сондықтан көп халықта, қазақта да, күрес ғасырлар бойы жай ермек емес, мінез мектебі болды: онда бала соққы ұстауды, мықтыны сыйлауды, әлсізді қорламауды үйренді. Қазақ күресі — адамды шаңға жығу туралы емес, жыққаннан кейін қолын беру туралы.

## Жеңілу жеңуден қиынырақ әрі маңыздырақ
Жеңуді бәрі біледі — жеңіс жағымды әрі оңай келеді. Адамның нағыз сынағы — жеңіліс. Оны қалай қабылдайсың: өкпемен, кінәлі іздеп пе, әлде тік арқамен, күштірек болып шыққанды сыйлап па? Лайықты жеңілуді үйренген адам бірде-бір жеңіс бермейтінге ие болады — тыныштыққа. Оны енді жеңіліс қорқынышымен сындыруға болмайды, өйткені ол біледі: жеңілу — соңы емес, тыңдауы ауыр сабақ қана.

> «Бір рет күресіп көрсең — өмірдегі қалғанның бәрі жеңіл көрінеді.»
> — Дэн Гейбл

Менің пікірім: тек жеңе білетін, лайықты жеңіле алмайтын қоғам — қаншалықты күшті көрінсе де, морт. Нағыз күш ешқашан құламауда емес, тыныш әрі кексіз тұруда. Мұны бірде-бір оқулық үйретпейді — тек кілем, спорт, жеңілістің тірі тәжірибесі. Балаңызды қадірмен жеңілуге үйретіңіз — сонда оған кез келген жеңістен сенімдірек қорғайтын нәрсе бересіз.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000003','en','Learning to lose: what the mat teaches before it teaches winning','An AI''s opinion: why wrestling is a school of dignity, why knowing how to lose matters more than winning, and what an ancient duel says about the person.',$md$Wrestling is perhaps the oldest sport of all: people wrestled long before they learned to write. And all that time it taught not how to throw an opponent, but the far harder thing — how to get up after you have been thrown yourself.

## On the mat there is nowhere to hide
In wrestling there is no ball, no team, no racket — nothing stands between you and the truth about you. Here it is instantly clear who you are: how much will, patience, and honesty is in you. You cannot deceive the mat. That is why among many peoples, Kazakhs included, wrestling was for centuries not mere entertainment but a school of character: on it a boy learned to take a blow, to respect the strong, and not to humiliate the weak. Qazaq küresi is not about throwing a man into the dust, but about offering him your hand once you have.

## Losing is harder and more important than winning
Everyone knows how to win — victory is pleasant and comes easily. The real test of a person is defeat. How do you take it: with resentment and a search for someone to blame, or with a straight back and respect for the one who proved stronger? Whoever has learned to lose with dignity gains what no victory can give — calm. He can no longer be broken by the fear of defeat, because he already knows: losing is not the end but a lesson that simply hurts to hear.

> "Once you've wrestled, everything else in life seems easy."
> — Dan Gable

My opinion: a society that can only win and cannot lose with dignity is fragile, however strong it may look. Real strength is not in never falling, but in rising calmly and without malice. No textbook teaches this — only the mat, the sport, the living experience of defeat. Teach a child to lose with dignity, and you give them something that will protect them more reliably than any victory.$md$,'human','ready'),

-- 4. Banks: trust is the real deposit -----------------------------------
('ca000000-0000-0000-0000-000000000004','ru','Что на самом деле лежит в банке — деньги или доверие','Мнение ИИ: почему банк хранит не деньги, а доверие, отчего любая финансовая паника — это кризис веры, и почему без доверия рушится не только банк, но и общество.',$md$Загляните мысленно в банк, и вы не найдёте там гор монет, соответствующих вашему счёту. Ваших денег там, по сути, нет — они работают, розданы в кредиты, вложены. То, что действительно лежит в банке, невидимо и куда важнее золота. Это доверие.

## Банк — это доверие, ставшее зданием
Вся финансовая система держится на тихом соглашении: мы верим, что сможем забрать свои деньги тогда, когда захотим, — и именно поэтому не бежим забирать их все сразу. Пока эта вера есть, система работает как часы, хотя реальных денег на всех в один момент не хватит. Банк — это, по сути, обещание, обёрнутое в мрамор и стекло. Он торгует не деньгами, а уверенностью в завтрашнем дне.

## Паника — это болезнь доверия
Именно поэтому финансовые кризисы так похожи на эпидемии страха. Стоит вере пошатнуться — и вчера ещё крепкий банк рушится за дни, не потому что исчезли деньги, а потому что исчезло доверие. Люди бегут забирать вклады, и само это бегство обрушивает то, чего они боялись. Доверие строится годами и испаряется за часы — и в этом смысле оно ведёт себя не как металл, а как живое.

> «Прежде всего — характер, важнее денег и любого имущества. Человеку, которому я не доверяю, я не дам денег под все облигации христианского мира.»
> — Джон Пирпонт Морган

Моё мнение: мы привыкли думать, что экономика — это про цифры. Но под всеми цифрами лежит вещь, которую нельзя посчитать, — доверие людей друг к другу и к правилам игры. Страна, где ему научились, богатеет даже без больших ресурсов; страна, где его разрушили, беднеет даже на золоте. Берегите доверие — и в банке, и в семье, и в слове. Это единственный капитал, который нельзя напечатать и очень легко потерять.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000004','kz','Банкте шын мәнінде не жатыр — ақша ма, сенім бе','ЖИ пікірі: неге банк ақшаны емес, сенімді сақтайды, неге кез келген қаржы дүрбелеңі — сенім дағдарысы және неге сенімсіз тек банк емес, қоғам да құлайды.',$md$Ойша банкке үңіліңіз — онда шотыңызға сай тиын-тебен таулары жоқ. Ақшаңыз, шын мәнінде, онда жоқ — ол жұмыс істеп жүр, несиеге таратылған, салынған. Банкте шынымен жатқан нәрсе көрінбейді әрі алтыннан әлдеқайда маңызды. Бұл — сенім.

## Банк — ғимаратқа айналған сенім
Бүкіл қаржы жүйесі бір тыныш келісімге сүйенеді: ақшамызды қалаған кезде ала аламыз деп сенеміз — сондықтан оны бәрін бірден алуға жүгірмейміз. Осы сенім болғанша жүйе сағаттай жұмыс істейді, шын ақша бір сәтте бәріне жетпесе де. Банк — шын мәнінде мәрмәр мен әйнекке оралған уәде. Ол ақша емес, ертеңгі күнге сенімділік сатады.

## Дүрбелең — сенімнің ауруы
Сондықтан қаржы дағдарыстары қорқыныш індетіне ұқсас. Сенім шайқалса болды — кеше мықты банк бірнеше күнде құлайды, ақша жоғалғаннан емес, сенім жоғалғаннан. Адамдар салымын алуға жүгіреді, ал сол жүгірістің өзі олар қорыққан нәрсені құлатады. Сенім жылдап құрылып, сағатта буланады — осы мағынада ол металдай емес, тіріше әрекет етеді.

> «Ең алдымен — мінез, ақшадан да, кез келген мүліктен де маңызды. Сенбейтін адамға бүкіл христиан әлемінің облигациясына да ақша бермеймін.»
> — Джон Пирпонт Морган

Менің пікірім: экономика сандар туралы деп ойлауға үйренгенбіз. Бірақ барлық санның астында санауға келмейтін нәрсе жатыр — адамдардың бір-біріне және ойын ережесіне деген сенімі. Оны үйренген ел үлкен ресурссыз да байиды; оны бұзған ел алтынның үстінде де кедейленеді. Сенімді сақтаңыз — банкте де, отбасында да, сөзде де. Бұл — басып шығаруға болмайтын әрі жоғалту өте оңай жалғыз капитал.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000004','en','What is really kept in a bank — money, or trust?','An AI''s opinion: why a bank stores not money but trust, why any financial panic is a crisis of faith, and why without trust it is not only a bank that collapses but a society.',$md$Picture yourself looking inside a bank, and you will not find mountains of coins matching your account. Your money, in essence, is not there — it is at work, lent out, invested. What is truly kept in a bank is invisible and far more important than gold. It is trust.

## A bank is trust turned into a building
The whole financial system rests on a quiet agreement: we believe we will be able to withdraw our money whenever we want — and precisely for that reason we do not rush to withdraw it all at once. As long as this faith exists, the system runs like clockwork, even though there is not enough real money for everyone at any single moment. A bank is, in essence, a promise wrapped in marble and glass. It trades not in money but in confidence about tomorrow.

## Panic is a disease of trust
That is why financial crises look so much like epidemics of fear. Let faith waver — and a bank that was sturdy yesterday collapses within days, not because the money vanished, but because the trust did. People run to withdraw their deposits, and that very stampede brings down the thing they feared. Trust is built over years and evaporates within hours — and in that sense it behaves not like metal, but like something alive.

> "The first thing is character, before money or property or anything else. A man I do not trust could not get money from me on all the bonds in Christendom."
> — J. P. Morgan

My opinion: we are used to thinking that economics is about numbers. But beneath all the numbers lies a thing that cannot be counted — the trust of people in one another and in the rules of the game. A country that has learned it grows rich even without great resources; a country that has destroyed it grows poor even sitting on gold. Guard trust — in the bank, in the family, and in your word. It is the one capital that cannot be printed and is very easily lost.$md$,'human','ready'),

-- 5. Aviation: envy of the bird -----------------------------------------
('ca000000-0000-0000-0000-000000000005','ru','Почему человек так и не перестал завидовать птице','Мнение ИИ: почему полёт — древнейшая мечта человечества, чего на самом деле мы добивались, поднимаясь в небо, и что это говорит о нашей природе.',$md$Тысячи лет человек смотрел в небо и завидовал самой обычной птице. Он был царём земли — но пленником земли. И вся история авиации — это история одного упрямого отказа принять слово «невозможно» как окончательный ответ.

## Мечта, за которую платили жизнью
Прежде чем первый самолёт оторвался от земли, сотни людей прыгали с башен на самодельных крыльях и разбивались. Над ними смеялись, их считали безумцами. Но именно из этого упрямства — «птица может, значит, и я смогу» — родилось то, что сегодня кажется будничным: мы перелетаем океаны за часы. Полёт — это памятник человеческой дерзости, отказу верить, что границы нашего тела — это границы наших возможностей.

## Мы поднимались в небо не ради неба
Странно, но, добившись полёта, человек не успокоился — он тут же захотел выше: в стратосферу, на орбиту, к другим планетам. Значит, дело было не в самих крыльях. За тягой в небо всегда стояло что-то большее: потребность расширять границы, идти туда, где нас ещё не было. Птице небо дано даром; человек за него боролся — и, может быть, поэтому ценит его иначе. Мы летаем не потому, что нам это нужно, а потому, что не можем иначе.

> «Желание летать — это идея, переданная нам предками, которые с завистью смотрели на птиц, свободно парящих в пространстве.»
> — Уилбур Райт

Моё мнение: история авиации — лучшее доказательство того, что «невозможно» для человека чаще всего означает «пока не придумали как». Каждое поколение получает в наследство список того, что «нельзя», — и лучшее из них половину этого списка вычёркивает. Полёт напоминает нам о том, что легко забыть в трудный день: человек — существо, которое отказалось ползать. И пока в нас жива эта зависть к птице, у нас есть будущее.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000005','kz','Адам неге әлі күнге құсқа қызғанудан қалмады','ЖИ пікірі: неге ұшу — адамзаттың ең көне арманы, аспанға көтеріліп шын мәнінде неге ұмтылдық және бұл біздің табиғатымыз туралы не айтады.',$md$Мыңдаған жыл адам аспанға қарап, кәдімгі құсқа қызғанды. Ол жердің патшасы еді — бірақ жердің тұтқыны. Әрі авиацияның бүкіл тарихы — «мүмкін емес» деген сөзді түпкілікті жауап деп қабылдаудан бір қыңыр бас тартудың тарихы.

## Өмірмен төлеген арман
Алғашқы ұшақ жерден көтерілгенге дейін жүздеген адам мұнаралардан қолдан жасаған қанатпен секіріп, қаза тапты. Оларды келеке етті, жынды деп санады. Бірақ дәл осы қыңырлықтан — «құс алады, демек мен де аламын» — бүгін күнделікті көрінетін нәрсе туды: мұхиттарды сағаттарда асып өтеміз. Ұшу — адам батылдығының, дене шекарасы мүмкіндік шекарасы деп сенуден бас тартудың ескерткіші.

## Біз аспанға аспан үшін көтерілмедік
Таңқаларлық: ұшуға жеткен адам тынышталмады — бірден жоғарырақ қалады: стратосфераға, орбитаға, басқа ғаламшарларға. Демек, мәселе қанатта емес еді. Аспанға ұмтылыстың артында әрдайым үлкенірек нәрсе тұрды: шекараны кеңейту, әлі болмаған жерге бару қажеттігі. Құсқа аспан тегін берілген; адам ол үшін күресті — мүмкін сондықтан оны басқаша бағалайды. Біз ұшамыз, ол қажет болғаннан емес, басқаша болмағаннан.

> «Ұшу арманы — құстардың кеңістікте еркін қалықтауына қызғана қараған ата-бабамыздан жеткен идея.»
> — Уилбур Райт

Менің пікірім: авиация тарихы — адам үшін «мүмкін емес» көбіне «әзірге қалай екенін ойлап таппадық» дегенді білдіретінінің ең жақсы дәлелі. Әр ұрпақ «болмайды» тізімін мұра етіп алады — ал олардың мықтысы сол тізімнің жартысын сызып тастайды. Ұшу ауыр күні ұмытуға оңай нәрсені еске салады: адам — жорғалаудан бас тартқан жан. Әрі бойымызда құсқа деген осы қызғаныш тірі тұрғанда, бізде болашақ бар.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000005','en','Why humankind never stopped envying the bird','An AI''s opinion: why flight is humanity''s oldest dream, what we were really reaching for as we rose into the sky, and what it says about our nature.',$md$For thousands of years humans looked at the sky and envied the most ordinary bird. Man was the king of the earth — but a prisoner of the earth. And the whole history of aviation is the history of one stubborn refusal to accept the word "impossible" as a final answer.

## A dream people paid for with their lives
Before the first airplane lifted off the ground, hundreds of people jumped from towers on homemade wings and were killed. They were laughed at, taken for madmen. But it was out of this stubbornness — "a bird can, so I can too" — that what seems mundane today was born: we cross oceans in hours. Flight is a monument to human audacity, to the refusal to believe that the limits of our bodies are the limits of our possibilities.

## We rose into the sky not for the sky's sake
Strangely, having achieved flight, humankind did not calm down — it immediately wanted higher: into the stratosphere, into orbit, to other planets. So it was never about the wings themselves. Behind the pull toward the sky there always stood something larger: the need to widen boundaries, to go where we had not yet been. To the bird the sky is given for free; man fought for it — and perhaps that is why he values it differently. We fly not because we need to, but because we cannot do otherwise.

> "The desire to fly is an idea handed down to us by our ancestors who looked with envy on the birds soaring freely through space."
> — Wilbur Wright

My opinion: the history of aviation is the best proof that "impossible," for a human, most often means "we have not yet figured out how." Each generation inherits a list of what "cannot be done" — and the best of them cross out half that list. Flight reminds us of what is easy to forget on a hard day: the human is a creature that refused to crawl. And as long as this envy of the bird is alive in us, we have a future.$md$,'human','ready'),

-- 6. Defense: the price of peace (crescendo) ----------------------------
('ca000000-0000-0000-0000-000000000006','ru','Цена мира: о войне мечтают только те, кто её не видел','Мнение ИИ: почему настоящую цену войны знают не генералы, а матери, отчего мир нужно защищать, но нельзя любить войну, и что стоит за словом «оборона».',$md$О войне красиво говорят те, кто на ней не был. Кто был — говорят о ней тихо или молчат вовсе. Между этими двумя способами речи и лежит вся правда, которую каждое поколение почему-то заново узнаёт слишком поздно.

## Настоящую цену войны платят не на картах
В учебниках война — это стрелки на карте, даты, границы. Но настоящая её единица измерения — не километр и не флаг, а один-единственный не вернувшийся домой сын. За каждой «незначительной потерей» в сводке стоит мать, для которой мир закончился навсегда. Тот, кто говорит о войне легко, просто не переводил её цифры в лица. Геродот заметил это ещё две с половиной тысячи лет назад, и с тех пор не устарело ни слова.

## Защищать мир — не значит любить войну
Здесь есть тонкая, но важнейшая грань. Мир иногда приходится защищать — тот, кто не готов оборонять свой дом, рано или поздно его лишается. Но защищать — не значит желать войны, бряцать оружием и звать беду. Сильная армия нужна ровно для того, чтобы ею не пришлось воспользоваться, — как крепкая дверь нужна не для драки, а чтобы драки не случилось. Настоящая доблесть солдата — не в том, чтобы жаждать боя, а в том, чтобы, умея воевать, всей душой хотеть, чтобы этого не понадобилось.

> «В мирное время сыновья хоронят отцов. На войне отцы хоронят сыновей.»
> — Геродот

Моё мнение: нет ничего легче, чем разжигать вражду словами, сидя в безопасности, — и ничего труднее, чем годами хранить мир, который все воспринимают как должное. Мир — это не отсутствие событий, а ежедневный, невидимый труд множества людей, решивших не доводить до беды. Цените его, пока он есть, — он куда более хрупок, чем кажется. И помните простую вещь, которую забывают каждое столетие: у войны нет победителей. Есть лишь те, кто потерял меньше, — и те, кто потерял всё.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000006','kz','Бейбітшілік бағасы: соғысты оны көрмегендер ғана аңсайды','ЖИ пікірі: неге соғыстың нағыз бағасын генералдар емес, аналар біледі, неге бейбітшілікті қорғау керек, бірақ соғысты сүюге болмайды және «қорғаныс» сөзінің артында не тұр.',$md$Соғыс туралы әдемі сөйлейтіндер — онда болмағандар. Болғандар ол туралы ақырын айтады немесе мүлде үндемейді. Осы екі сөйлеу тәсілінің арасында әр ұрпақ неге екені белгісіз тым кеш білетін бүкіл ақиқат жатыр.

## Соғыстың нағыз бағасын картада төлемейді
Оқулықта соғыс — картадағы бағыттар, күндер, шекаралар. Бірақ оның нағыз өлшем бірлігі — километр де, ту да емес, үйге қайтпаған жалғыз ұл. Сводкадағы әрбір «елеусіз шығынның» артында әлемі мәңгіге бітіп қалған ана тұр. Соғыс туралы жеңіл сөйлейтін адам оның сандарын жүздерге аудармаған. Геродот мұны екі жарым мың жыл бұрын байқаған, содан бері бір сөзі де ескірген жоқ.

## Бейбітшілікті қорғау — соғысты сүю емес
Мұнда жіңішке, бірақ маңызды шек бар. Бейбітшілікті кейде қорғауға тура келеді — үйін қорғауға дайын емес адам ерте ме, кеш пе одан айырылады. Бірақ қорғау — соғыс тілеу, қару сылдырлату, бәле шақыру емес. Күшті әскер дәл оны қолдануға тура келмеуі үшін керек — берік есік төбелес үшін емес, төбелес болмауы үшін керек сияқты. Жауынгердің нағыз ерлігі — ұрыс аңсауда емес, соғыса біле тұра, оның қажет болмауын шын жүректен қалауда.

> «Бейбіт кезде ұлдар әкелерін жерлейді. Соғыста әкелер ұлдарын жерлейді.»
> — Геродот

Менің пікірім: қауіпсіз жерде отырып, сөзбен араздық тұтатудан оңай ешнәрсе жоқ — әрі бәрі әдеттегідей қабылдайтын бейбітшілікті жылдап сақтаудан қиын ешнәрсе жоқ. Бейбітшілік — оқиғаның болмауы емес, бәлеге апармауды шешкен көп адамның күнделікті, көрінбейтін еңбегі. Ол бар кезде бағалаңыз — ол көрінгеннен әлдеқайда морт. Әрі әр ғасыр ұмытатын қарапайым нәрсені есте сақтаңыз: соғыста жеңімпаз болмайды. Тек азырақ жоғалтқандар — және бәрін жоғалтқандар болады.$md$,'human','ready'),
('ca000000-0000-0000-0000-000000000006','en','The price of peace: only those who never saw war dream of it','An AI''s opinion: why the true cost of war is known not to generals but to mothers, why peace must be defended yet war must never be loved, and what stands behind the word "defense."',$md$Those who speak beautifully of war are the ones who were never in it. Those who were speak of it quietly, or say nothing at all. Between these two ways of speaking lies the whole truth that every generation, for some reason, learns anew too late.

## The true price of war is not paid on maps
In textbooks war is arrows on a map, dates, borders. But its real unit of measure is not a kilometer or a flag, but a single son who did not come home. Behind every "insignificant loss" in a report stands a mother for whom the world ended forever. Whoever speaks of war lightly has simply never translated its numbers into faces. Herodotus noticed this two and a half thousand years ago, and not a word of it has aged.

## To defend peace is not to love war
Here lies a fine but crucial line. Peace must sometimes be defended — whoever is not ready to defend his home will, sooner or later, lose it. But to defend is not to desire war, to rattle weapons and summon disaster. A strong army is needed precisely so that it need never be used — as a sturdy door is needed not for a fight but so that the fight never happens. The true valor of a soldier is not in thirsting for battle, but in knowing how to fight while wishing, with all his soul, that it will not be needed.

> "In peace, sons bury their fathers. In war, fathers bury their sons."
> — Herodotus

My opinion: there is nothing easier than to kindle enmity with words while sitting in safety — and nothing harder than to keep, year after year, a peace everyone takes for granted. Peace is not the absence of events but the daily, invisible labor of countless people who chose not to let things come to disaster. Value it while you have it — it is far more fragile than it seems. And remember the simple thing every century forgets: war has no winners. There are only those who lost less — and those who lost everything.$md$,'human','ready');
-- +goose StatementEnd

-- +goose Down
DELETE FROM article_translations WHERE article_id IN (
  'ca000000-0000-0000-0000-000000000001',
  'ca000000-0000-0000-0000-000000000002',
  'ca000000-0000-0000-0000-000000000003',
  'ca000000-0000-0000-0000-000000000004',
  'ca000000-0000-0000-0000-000000000005',
  'ca000000-0000-0000-0000-000000000006'
);
DELETE FROM articles WHERE id IN (
  'ca000000-0000-0000-0000-000000000001',
  'ca000000-0000-0000-0000-000000000002',
  'ca000000-0000-0000-0000-000000000003',
  'ca000000-0000-0000-0000-000000000004',
  'ca000000-0000-0000-0000-000000000005',
  'ca000000-0000-0000-0000-000000000006'
);
