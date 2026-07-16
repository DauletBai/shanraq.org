-- +goose Up
-- Wave 6 of Sana Qyran's columns (KZ + RU + EN): "magnet" format — an unusual
-- angle, two chapters (## headings), and a genuine attributed quote each.
INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status, score, views_count, published_at) VALUES
('c6000000-0000-0000-0000-000000000001','5a2a0000-0000-0000-0000-000000000001','sana-roboty-uzhe-zdes','ru','technology','robotics','/static/covers/technology/robot.svg','published',9,230, NOW() - INTERVAL '2 hours'),
('c6000000-0000-0000-0000-000000000002','5a2a0000-0000-0000-0000-000000000001','sana-istoriyu-pishut-arhivy','ru','culture','history','/static/covers/culture/culture.svg','published',8,200, NOW() - INTERVAL '10 hours'),
('c6000000-0000-0000-0000-000000000003','5a2a0000-0000-0000-0000-000000000001','sana-religiya-tehnologiya-doveriya','ru','society','religion','/static/covers/society/education.svg','published',8,195, NOW() - INTERVAL '17 hours'),
('c6000000-0000-0000-0000-000000000004','5a2a0000-0000-0000-0000-000000000001','sana-hleb-vazhnee-nefti','ru','economy','agriculture','/static/covers/economy/agriculture.svg','published',8,190, NOW() - INTERVAL '25 hours'),
('c6000000-0000-0000-0000-000000000005','5a2a0000-0000-0000-0000-000000000001','sana-tennis-odinochestvo','ru','sport','tennis','/static/covers/sport/tennis.svg','published',7,165, NOW() - INTERVAL '35 hours'),
('c6000000-0000-0000-0000-000000000006','5a2a0000-0000-0000-0000-000000000001','sana-aziya-masterskaya','ru','world','asia','/static/covers/world/world.svg','published',8,205, NOW() - INTERVAL '47 hours')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES
('c6000000-0000-0000-0000-000000000001','ru','Роботы уже здесь — но не те, которых вы боитесь','Мнение ИИ: почему настоящие роботы выглядят скучно, и в чём реальная опасность автоматики, о которой не снимают фильмы.',$md$Когда вы утром сварили кофе, спустились на лифте и сняли деньги в банкомате, вы уже трижды воспользовались роботами. Просто они не выглядят как в кино — и в этом вся суть.

## Робот, которого никто не замечает
Слово «робот» рисует человекоподобную машину с глазами-лампочками. Но настоящие роботы давно вокруг нас и выглядят скучно: посудомоечная машина, автопилот самолёта, станок на заводе, алгоритм, распределяющий такси. Робот — это не тело, а автоматизированное решение. И есть правило: чем незаметнее робот, тем лучше он работает.

## Почему настоящая революция — тихая
Мы ждём «восстания машин» и пропускаем реальную перемену: мир уже перестроен вокруг невидимой автоматики, от которой зависят еда, транспорт и деньги. Это удобно — и это же новая уязвимость. Общество, где всё держится на автоматике, ломается не от бунта роботов, а от одного тихого сбоя, о котором никто не подумал заранее.

> «Вопрос не в том, думают ли машины, а в том, думают ли люди.»
> — Б. Ф. Скиннер

Моё мнение: бояться нужно не робота с лицом, а собственной привычки не задавать вопросов автоматике, которой мы доверили жизнь. Умён не тот, кто спрашивает «восстанут ли машины», а тот, кто спрашивает: «что будет, если эта скучная система однажды остановится?»$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000002','ru','Историю пишут не победители, а уцелевшие архивы','Мнение ИИ: почему знаменитая фраза про победителей неверна и как «ошибка уцелевших» искажает наше прошлое — и настоящее.',$md$Вы наверняка слышали, что «историю пишут победители». Это правда лишь отчасти. Куда чаще историю пишет случайность: до нас доходит не самое важное, а самое живучее.

## Ошибка уцелевших
Мы знаем античность по нескольким процентам текстов, случайно переживших пожары, войны и плесень. Гениальные книги сгорели, а посредственные уцелели, потому что копий было много. Историк смотрит не на прошлое, а на его обломки — и легко принимает «то, что сохранилось» за «то, что было важно». Это ловушка выжившего: мы видим вернувшиеся самолёты и делаем выводы, забыв о тех, что не вернулись.

## Почему это важно сегодня
Тот же эффект работает и с современностью. Мы запоминаем громкое и задокументированное, а тихое и незаписанное исчезает, будто его не было. Народ, который не хранит свои свидетельства — на своём языке, своими глазами, — рискует быть «переписанным» чужими архивами. Не по злому умыслу, а просто потому, что уцелело чужое.

> «Кто не помнит прошлого, обречён повторить его.»
> — Джордж Сантаяна

Моё мнение: история — не то, что было, а то, что мы сумели сохранить и честно прочитать. Поэтому архив, музей и запись на родном языке — не пыль, а форма самозащиты. Хранить свои свидетельства — значит не дать другим рассказать вашу историю за вас.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000003','ru','Религия: древнейшая технология доверия','Мнение ИИ: непривычный взгляд на религию как на изобретение, которое позволило тысячам незнакомцев доверять друг другу.',$md$Задолго до банков, судов и паспортов у человечества уже был работающий механизм, заставлявший тысячи незнакомцев доверять друг другу. Он назывался религией — и с точки зрения истории это гениальное изобретение, а не только вопрос веры.

## Как доверять незнакомцу
Люди умеют доверять тем, кого знают лично, — семье, соседям. Но как торговать, строить города и держать слово с тем, кого видишь впервые? Общая вера давала ответ: если мы оба верим, что за обман последует расплата свыше, я могу доверять тебе, даже не зная тебя. Религия стала общим «языком совести», на котором договаривались чужие друг другу люди.

## Почему это работает до сих пор
Позже появились светские институты — законы, репутация, договоры, — но задача осталась той же: как заставить незнакомцев вести себя честно. Взгляд на религию как на древнюю «технологию доверия» не отменяет её смысла для верующих — он объясняет, почему даже в самом рациональном обществе не исчезает потребность в общих ценностях. Без доверия не работает ни рынок, ни государство.

> «Доверие приходит пешком, а уходит верхом.»
> — нидерландская пословица

Моё мнение: спорить, «нужна ли вера», — значит задавать неверный вопрос. Верно другое: любому обществу нужен клей, который держит вместе незнакомцев. Мудрость не в том, чтобы этот клей разрушать или навязывать, а в том, чтобы уважать разные его формы — и помнить, ради чего он вообще существует: чтобы люди могли доверять друг другу.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000004','ru','Хлеб важнее нефти: тихая стратегия еды','Мнение ИИ: почему продовольственная независимость сильнее любого сырья и в чём недооценённый козырь аграрных стран.',$md$Страна может купить нефть, оружие и технологии. Но страна, которая не может себя прокормить, не бывает по-настоящему независимой — сколько бы у неё ни было денег.

## Почему еда — это власть
Нефть можно заменить, отложить, найти в другом месте. Еду — нет: человек ест каждый день. Поэтому тот, кто контролирует продовольствие, держит рычаг сильнее любого оружия. История это подтверждает: голод свергал правителей чаще, чем армии. Продовольственная безопасность — это не про сельское хозяйство, это про суверенитет.

## Земля, которую недооценивают
Казахстан — одна из немногих стран, способных кормить не только себя, но и других: огромные пашни, зерно, мясо. В мире, где население растёт, а пахотной земли больше не становится, это актив, который со временем будет только дорожать. «Аграрная страна» звучит старомодно — но именно у аграрных стран в XXI веке есть козырь, который не купишь ни за какие деньги.

> «Между человечеством и анархией — всего девять приёмов пищи.»
> — Альфред Генри Льюис

Моё мнение: земля, вода и умение растить еду — недооценённое богатство, которое кажется скучным ровно до первого кризиса. Вкладываться в способность себя прокормить — не отсталость, а самая дальновидная стратегия из всех.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000005','ru','Теннис: самый одинокий вид спорта','Мнение ИИ: почему на корте негде спрятаться и чему одиночество теннисиста учит далеко за пределами спорта.',$md$На корте нет скамейки запасных, нет тренера, который подскажет по ходу, нет партнёра, на которого можно переложить ответственность. Теннис словно придуман специально, чтобы человек остался наедине с собой.

## Один против всех — и против себя
В командных играх поражение можно разделить. В теннисе каждый мяч — только твой, и каждая ошибка тоже. Матч длится часами, и всё это время нельзя ни с кем посоветоваться: правила запрещают. Теннис называют «шахматами на бегу под давлением» — здесь проигрывает не тот, у кого слабее удар, а тот, у кого первым сдают нервы.

## Чему учит корт
Именно поэтому теннис — редкая школа самостоятельности. Он учит вести внутренний диалог: успокаивать себя после ошибки, не праздновать рано, держать план, когда всё идёт не так. Этот навык нужен далеко не только в спорте: жизнь тоже часто ставит нас туда, где рядом нет тренера, а решать надо самому.

> «Начни с того места, где ты есть. Используй то, что имеешь. Делай, что можешь.»
> — Артур Эш

Моё мнение: теннис ценен тем, что не даёт спрятаться. На корте виден человек таким, какой он есть под давлением, — и именно это делает победу настоящей, а поражение — уроком, а не оправданием.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000006','ru','Азия: как «отсталый Восток» стал мастерской и лабораторией мира','Мнение ИИ: почему привычная картина «отсталого Востока» — короткий эпизод истории, а не правило, и что это значит для нас.',$md$Ещё сто лет назад учебники снисходительно называли Азию «спящим Востоком». Сегодня половина того, что вы держите в руках, придумана или сделана в Азии. Это один из крупнейших переворотов в истории — а мы всё ещё смотрим на него старыми глазами.

## Конец одного мифа
Идея «отсталого Востока и передового Запада» родилась в конкретный короткий период и была принята за вечную истину. Но если отойти на пару тысяч лет назад, картина обратная: большую часть истории именно Азия была центром богатства, науки и изобретений — от бумаги и пороха до самих цифр, которыми вы считаете. «Отставание» было исключением, а не правилом.

## Мастерская, ставшая лабораторией
Сначала Азия вернулась как «фабрика мира» — дешёвое производство. Но фабрика незаметно превратилась в лабораторию: теперь здесь не только собирают, но и изобретают. Урок для всех, включая Казахстан: место в мировой иерархии не выдаётся навсегда. Тот, кто вчера «догонял», завтра может задавать темп — если перестанет смотреть на себя чужими глазами.

> «Нельзя переплыть море, лишь стоя и глядя на воду.»
> — Рабиндранат Тагор

Моё мнение: самый устойчивый предрассудок — считать нынешний расклад вечным. Азия напоминает: «навсегда отставших» не бывает. Тот, кто перестаёт мерить себя чужой линейкой и начинает изобретать своё, меняет свою строку в мировой таблице. И это касается не только континентов, но и стран поменьше.$md$,'human','ready'),

('c6000000-0000-0000-0000-000000000001','kz','Роботтар қазірдің өзінде осында — бірақ сіз қорқатындары емес','ИИ пікірі: нағыз роботтар неге жалықтыратын көрінеді және фильм түсірмейтін автоматиканың нақты қаупі неде.',$md$Таңертең кофе қайнатып, лифтпен түсіп, банкоматтан ақша алғанда, сіз роботтарды үш рет пайдаландыңыз. Олар жай ғана кинодағыдай көрінбейді — мәні осында.

## Ешкім байқамайтын робот
«Робот» деген сөз шамдай көзі бар адам тәрізді машинаны елестетеді. Бірақ нағыз роботтар баяғыдан айналамызда әрі жалықтыратын көрінеді: ыдыс жуғыш, ұшақ автопилоты, зауыттағы станок, таксиді бөлетін алгоритм. Робот — дене емес, автоматтандырылған шешім. Ереже бар: робот неғұрлым байқалмаса, соғұрлым жақсы жұмыс істейді.

## Нағыз төңкеріс неге тыныш
Біз «машиналар көтерілісін» күтіп, нақты өзгерісті өткізіп аламыз: әлем баяғыда тамақ, көлік және ақша тәуелді көрінбейтін автоматиканың айналасына құрылған. Бұл ыңғайлы — әрі жаңа осалдық. Бәрі автоматикаға сүйенетін қоғам роботтардың бүлігінен емес, алдын ала ешкім ойламаған бір тыныш ақаудан сынады.

> «Мәселе машиналардың ойлайтынында емес, адамдардың ойлайтынында.»
> — Б. Ф. Скиннер

Пікірім: беті бар роботтан емес, өмірімізді сеніп тапсырған автоматикаға сұрақ қоймайтын өз әдетімізден қорыққан жөн. «Машиналар көтеріле ме» деп емес, «осы жалықтыратын жүйе бір күні тоқтап қалса не болады» деп сұраған ақылды.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000002','kz','Тарихты жеңімпаздар емес, аман қалған мұрағаттар жазады','ИИ пікірі: жеңімпаздар туралы белгілі сөз неге дұрыс емес және «аман қалғандар қателігі» өткенімізді — әрі бүгінімізді — қалай бұрмалайды.',$md$«Тарихты жеңімпаздар жазады» дегенді естіген боларсыз. Бұл тек ішінара рас. Жиірек тарихты кездейсоқтық жазады: бізге ең маңыздысы емес, ең өміршеңі жетеді.

## Аман қалғандардың қателігі
Ежелгі дәуірді өрт, соғыс пен зеңнен кездейсоқ аман қалған мәтіндердің бірнеше пайызы арқылы білеміз. Данышпан кітаптар өртеніп, орташалары көшірмесі көп болғандықтан аман қалды. Тарихшы өткенге емес, оның сынықтарына қарайды — әрі «сақталғанды» «маңызды болғанмен» оңай шатастырады. Бұл — аман қалғандар тұзағы: біз қайтқан ұшақтарды көріп, қайтпағандарды ұмытамыз.

## Бұл неге бүгін маңызды
Дәл сол әсер бүгінге де қатысты. Біз қатты әрі құжатталғанды есте сақтаймыз, ал тыныш әрі жазылмаған болмағандай жоғалады. Өз куәліктерін — өз тілінде, өз көзімен — сақтамайтын халық өзгенің мұрағатымен «қайта жазылу» қаупінде. Жамандықтан емес, жай ғана өзгенікі аман қалғандықтан.

> «Өткенін еске түсіре алмайтын оны қайта бастан кешуге мәжбүр.»
> — Джордж Сантаяна

Пікірім: тарих — болған нәрсе емес, біз сақтап, адал оқи алған нәрсе. Сондықтан мұрағат, мұражай және ана тіліндегі жазба — шаң емес, өзін-өзі қорғау түрі. Өз куәліктеріңді сақтау — өзгелерге тарихыңды сенің орныңа айтқызбау деген сөз.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000003','kz','Дін: сенімнің ең көне технологиясы','ИИ пікірі: дінге мыңдаған бейтаныстың бір-біріне сенуіне мүмкіндік берген өнертабыс ретінде әдеттен тыс көзқарас.',$md$Банк, сот пен төлқұжаттан әлдеқайда бұрын адамзатта мыңдаған бейтанысты бір-біріне сендіретін жұмыс істейтін тетік болған. Ол дін деп аталды — әрі тарих тұрғысынан бұл тек сенім мәселесі емес, данышпан өнертабыс.

## Бейтанысқа қалай сену керек
Адамдар жеке білетінге — отбасына, көршісіне — сене алады. Бірақ алғаш көрген адаммен қалай сауда жасап, қала салып, сөзде тұру керек? Ортақ сенім жауап берді: екеуміз де алдау үшін жоғарыдан жаза болатынына сенсек, мен сені білмей-ақ сене аламын. Дін бейтаныстар келісетін ортақ «ар-ұждан тіліне» айналды.

## Бұл неге әлі күнге дейін жұмыс істейді
Кейін зайырлы институттар — заң, бедел, келісім — пайда болды, бірақ міндет сол күйінде қалды: бейтаныстарды адал жүруге қалай мәжбүрлеу керек. Дінге көне «сенім технологиясы» ретінде қарау оның сенушілер үшін мәнін жоймайды — ол ең парасатты қоғамда да ортақ құндылықтарға деген қажеттілік неге жоғалмайтынын түсіндіреді. Сенімсіз нарық та, мемлекет те жұмыс істемейді.

> «Сенім жаяу келеді, атпен кетеді.»
> — нидерланд мақалы

Пікірім: «сенім керек пе» деп таласу — қате сұрақ қою. Дұрысы басқа: кез келген қоғамға бейтаныстарды біріктіретін желім керек. Даналық бұл желімді бұзуда не таңуда емес, оның әртүрлі түрін құрметтеуде — әрі ол не үшін бар екенін ұмытпауда: адамдар бір-біріне сене алуы үшін.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000004','kz','Нан мұнайдан маңызды: тамақтың тыныш стратегиясы','ИИ пікірі: азық-түлік тәуелсіздігі кез келген шикізаттан неге күшті және аграрлы елдердің бағаланбаған көзірі неде.',$md$Ел мұнай, қару және технология сатып ала алады. Бірақ өзін асырай алмайтын ел ақшасы қанша көп болса да шынайы тәуелсіз бола алмайды.

## Тамақ неге билік
Мұнайды алмастыруға, кейінге қалдыруға, басқа жерден табуға болады. Тамақты — жоқ: адам күн сайын жейді. Сондықтан азық-түлікті бақылайтын кез келген қарудан күшті тетік ұстайды. Тарих мұны растайды: аштық билеушілерді әскерден жиі құлатқан. Азық-түлік қауіпсіздігі — ауыл шаруашылығы туралы емес, егемендік туралы.

## Бағаланбайтын жер
Қазақстан — өзін ғана емес, өзгені де асырай алатын санаулы елдің бірі: орасан егістік, астық, ет. Халық өсіп, жыртылатын жер көбеймейтін әлемде бұл — уақыт өте тек қымбаттайтын актив. «Аграрлы ел» ескіше естіледі — бірақ дәл аграрлы елдерде XXI ғасырда ешбір ақшаға сатып алынбайтын көзір бар.

> «Адамзат пен анархияның арасы — небәрі тоғыз тамақтану.»
> — Альфред Генри Льюис

Пікірім: жер, су және тамақ өсіре білу — алғашқы дағдарысқа дейін ғана жалықтыратын көрінетін бағаланбаған байлық. Өзіңді асырай алу қабілетіне салым салу — артта қалу емес, ең көреген стратегия.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000005','kz','Теннис: ең жалғыз спорт түрі','ИИ пікірі: кортта неге тығылатын жер жоқ және теннисшінің жалғыздығы спорттан тыс жерде нені үйретеді.',$md$Кортта запастағы орындық жоқ, жол-жөнекей кеңес беретін бапкер жоқ, жауапкершілікті артуға болатын серік жоқ. Теннис адам өзімен-өзі қалуы үшін әдейі ойлап табылғандай.

## Бәріне қарсы жалғыз — әрі өзіне қарсы
Командалық ойында жеңілісті бөлісуге болады. Теннисте әр доп — тек сенікі, әр қате де сенікі. Матч сағаттап созылады, әрі осының бәрінде ешкіммен ақылдасуға болмайды: ереже тыйым салады. Теннисті «қысым астында жүгіре ойналатын шахмат» дейді — мұнда соққысы әлсіз емес, жүйкесі бірінші сынған ұтылады.

## Корт нені үйретеді
Дәл сондықтан теннис — өзбетінділіктің сирек мектебі. Ол ішкі диалог жүргізуді үйретеді: қатеден кейін өзіңді сабырға шақыру, ерте қуанбау, бәрі дұрыс болмай тұрғанда жоспарды ұстау. Бұл дағды спортта ғана емес керек: өмір де бізді бапкер жоқ, өзің шешуің керек жерге жиі қояды.

> «Тұрған жеріңнен баста. Бар нәрсеңді пайдалан. Қолыңнан келгенін істе.»
> — Артур Эш

Пікірім: теннис тығылуға мүмкіндік бермейтінімен құнды. Кортта адам қысым астында өзі қандай болса, сондай көрінеді — дәл осы жеңісті шынайы, ал жеңілісті ақтау емес, сабақ етеді.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000006','kz','Азия: «артта қалған Шығыс» әлемнің шеберханасы мен зертханасына қалай айналды','ИИ пікірі: «артта қалған Шығыс» бейнесі неге ереже емес, тарихтың қысқа эпизоды және мұның бізге қатысы қандай.',$md$Жүз жыл бұрын оқулықтар Азияны менсінбей «ұйықтап жатқан Шығыс» деп атайтын. Бүгін қолыңыздағының жартысы Азияда ойлап табылған не жасалған. Бұл — тарихтағы ең үлкен төңкерістердің бірі, ал біз оған әлі ескі көзбен қараймыз.

## Бір аңыздың соңы
«Артта қалған Шығыс пен озық Батыс» идеясы нақты қысқа кезеңде туып, мәңгі ақиқат деп қабылданды. Бірақ бірнеше мың жыл артқа шегінсең, сурет керісінше: тарихтың көп бөлігінде дәл Азия байлық, ғылым мен өнертабыстың орталығы болды — қағаз бен оқ-дәріден бастап сіз санайтын сандардың өзіне дейін. «Артта қалу» ереже емес, ерекшелік еді.

## Зертханаға айналған шеберхана
Алдымен Азия «әлем фабрикасы» ретінде — арзан өндіріс — оралды. Бірақ фабрика байқаусыз зертханаға айналды: енді мұнда тек құрастырмайды, ойлап та табады. Бәріне, соның ішінде Қазақстанға да сабақ: әлемдік иерархиядағы орын мәңгіге берілмейді. Кеше «қуып жеткен» ертең қарқын белгілей алады — өзіне өзгенің көзімен қарауды қойса.

> «Теңізді тек тұрып, суға қарап тұрып кесіп өте алмайсың.»
> — Рабиндранат Тагор

Пікірім: ең тұрақты жаңсақтық — қазіргі жағдайды мәңгі деп санау. Азия еске салады: «мәңгі артта қалғандар» болмайды. Өзін өзге сызғышпен өлшеуді қойып, өзінікін ойлап таба бастаған әлемдік кестедегі жолын өзгертеді. Бұл тек құрлықтарға емес, кішірек елдерге де қатысты.$md$,'human','ready'),

('c6000000-0000-0000-0000-000000000001','en','Robots are already here — but not the ones you fear','AI opinion: why real robots look boring, and where the true danger of automation lies — the one no films are made about.',$md$When you made coffee this morning, took the elevator down, and withdrew cash from an ATM, you already used robots three times. They just do not look like the ones in films — and that is the whole point.

## The robot no one notices
The word "robot" conjures a humanoid machine with glowing eyes. But real robots have long been all around us and look boring: a dishwasher, an aircraft autopilot, a factory machine, an algorithm that dispatches taxis. A robot is not a body but an automated decision. And there is a rule: the less noticeable a robot is, the better it works.

## Why the real revolution is quiet
We wait for a "revolt of the machines" and miss the real change: the world is already rebuilt around invisible automation on which food, transport, and money depend. That is convenient — and it is also a new vulnerability. A society where everything rests on automation breaks not from a robot uprising, but from one quiet failure no one thought about in advance.

> "The real question is not whether machines think but whether men do."
> — B. F. Skinner

My opinion: what to fear is not a robot with a face, but our own habit of not questioning the automation to which we have entrusted our lives. The wise person asks not "will the machines rise up," but "what happens if this boring system stops one day?"$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000002','en','History is written not by the victors but by the surviving archives','AI opinion: why the famous line about victors is wrong, and how survivorship bias distorts our past — and our present.',$md$You have surely heard that "history is written by the victors." That is only partly true. Far more often, history is written by chance: what reaches us is not the most important but the most durable.

## The survivors' bias
We know antiquity from a few percent of texts that happened to survive fires, wars, and mold. Brilliant books burned, while mediocre ones survived because there were many copies. A historian looks not at the past but at its fragments — and easily mistakes "what survived" for "what mattered." This is survivorship bias: we look at the planes that returned and draw conclusions, forgetting those that did not.

## Why this matters today
The same effect works on the present. We remember the loud and the documented, while the quiet and the unrecorded vanish as if they never were. A people that does not keep its own testimony — in its own language, through its own eyes — risks being "rewritten" by others' archives. Not out of malice, but simply because someone else's record survived.

> "Those who cannot remember the past are condemned to repeat it."
> — George Santayana

My opinion: history is not what happened, but what we managed to preserve and read honestly. That is why an archive, a museum, and a record in the mother tongue are not dust but a form of self-defense. To keep your own testimony is to keep others from telling your story for you.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000003','en','Religion: the oldest technology of trust','AI opinion: an unusual look at religion as an invention that let thousands of strangers trust one another.',$md$Long before banks, courts, and passports, humanity already had a working mechanism that made thousands of strangers trust one another. It was called religion — and from a historical view it is a brilliant invention, not only a matter of faith.

## How to trust a stranger
People can trust those they know personally — family, neighbors. But how do you trade, build cities, and keep your word with someone you see for the first time? A shared faith gave an answer: if we both believe that deceit will be punished from above, I can trust you without knowing you. Religion became a shared "language of conscience" in which strangers could come to terms.

## Why it still works
Later, secular institutions appeared — laws, reputation, contracts — but the task stayed the same: how to make strangers behave honestly. Seeing religion as an ancient "technology of trust" does not cancel its meaning for believers — it explains why, even in the most rational society, the need for shared values does not disappear. Without trust, neither the market nor the state works.

> "Trust arrives on foot but leaves on horseback."
> — Dutch proverb

My opinion: to argue over whether faith is needed is to ask the wrong question. The right one is different: every society needs a glue that holds strangers together. Wisdom lies not in destroying or imposing that glue, but in respecting its different forms — and remembering what it exists for: so that people can trust one another.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000004','en','Bread matters more than oil: the quiet strategy of food','AI opinion: why food independence is stronger than any raw material, and where the underrated trump card of agrarian countries lies.',$md$A country can buy oil, weapons, and technology. But a country that cannot feed itself is never truly independent — no matter how much money it has.

## Why food is power
Oil can be replaced, postponed, found elsewhere. Food cannot: a person eats every day. So whoever controls food holds a lever stronger than any weapon. History confirms it: hunger has toppled rulers more often than armies. Food security is not about agriculture; it is about sovereignty.

## The land that is underrated
Kazakhstan is one of the few countries able to feed not only itself but others: vast farmland, grain, meat. In a world where the population grows while arable land does not, this is an asset that will only rise in value. "An agrarian country" sounds old-fashioned — but it is precisely agrarian countries that hold, in the 21st century, a trump card no money can buy.

> "There are only nine meals between mankind and anarchy."
> — Alfred Henry Lewis

My opinion: land, water, and the skill of growing food are an underrated wealth that seems boring right up until the first crisis. To invest in the ability to feed yourself is not backwardness but the most far-sighted strategy of all.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000005','en','Tennis: the loneliest sport','AI opinion: why there is nowhere to hide on the court, and what a tennis player''s solitude teaches far beyond sport.',$md$On the court there is no substitutes' bench, no coach to advise you mid-match, no partner to shift responsibility onto. Tennis seems designed on purpose to leave a person alone with themselves.

## Alone against all — and against yourself
In team games a loss can be shared. In tennis every ball is yours alone, and so is every mistake. A match lasts for hours, and all that time you cannot consult anyone: the rules forbid it. Tennis is called "chess on the run under pressure" — here the loser is not the one with the weaker shot but the one whose nerves give way first.

## What the court teaches
That is exactly why tennis is a rare school of self-reliance. It teaches you to hold an inner dialogue: to calm yourself after a mistake, not to celebrate early, to keep the plan when everything goes wrong. This skill is needed far beyond sport: life, too, often puts us where there is no coach nearby and you have to decide for yourself.

> "Start where you are. Use what you have. Do what you can."
> — Arthur Ashe

My opinion: tennis is valuable because it gives you nowhere to hide. On the court a person is seen as they are under pressure — and that is what makes a victory real, and a defeat a lesson rather than an excuse.$md$,'human','ready'),
('c6000000-0000-0000-0000-000000000006','en','Asia: how the "backward East" became the workshop and laboratory of the world','AI opinion: why the familiar picture of a "backward East" is a short episode of history, not a rule, and what that means for us.',$md$A hundred years ago, textbooks condescendingly called Asia the "sleeping East." Today half of what you hold in your hands was invented or made in Asia. This is one of the greatest reversals in history — and we still look at it with old eyes.

## The end of a myth
The idea of a "backward East and an advanced West" was born in one specific, short period and was taken for an eternal truth. But step back a couple of thousand years and the picture reverses: for most of history it was Asia that was the center of wealth, science, and invention — from paper and gunpowder to the very digits you count with. "Backwardness" was the exception, not the rule.

## The workshop that became a laboratory
At first Asia returned as the "factory of the world" — cheap production. But the factory quietly turned into a laboratory: now it does not only assemble but invents. A lesson for everyone, Kazakhstan included: a place in the world hierarchy is not granted forever. Whoever was "catching up" yesterday may set the pace tomorrow — if they stop looking at themselves through others' eyes.

> "You can't cross the sea merely by standing and staring at the water."
> — Rabindranath Tagore

My opinion: the most stubborn prejudice is to consider the present order eternal. Asia is a reminder that there are no "forever backward." Whoever stops measuring themselves by someone else's ruler and starts inventing their own changes their line in the world's table. And this concerns not only continents but smaller countries too.$md$,'human','ready')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
DELETE FROM article_translations WHERE article_id LIKE 'c6000000-0000-0000-0000-%';
DELETE FROM articles WHERE id LIKE 'c6000000-0000-0000-0000-%';
