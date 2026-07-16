-- +goose Up
-- Wave 11 of Sana Qyran's columns (KZ + RU + EN): each essay is paired with a
-- famous public-domain painting as its cover (downloaded from Wikimedia Commons,
-- cropped to 16:9 WebP under web/static/covers). Chapter headings (##), a genuine
-- attributed quote, a signed opinion and a cover credit line each.
-- Subrubrics: gamedev, debate, industry, government, oceania, editorial.
INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status, score, views_count, published_at) VALUES
('cb000000-0000-0000-0000-000000000001','5a2a0000-0000-0000-0000-000000000001','sana-zachem-igraem','ru','it','gamedev','/static/covers/it/gamedev.webp','published',9,255, NOW() - INTERVAL '13 hours'),
('cb000000-0000-0000-0000-000000000002','5a2a0000-0000-0000-0000-000000000001','sana-sila-spora','ru','opinion','debate','/static/covers/opinion/debate.webp','published',9,260, NOW() - INTERVAL '11 hours'),
('cb000000-0000-0000-0000-000000000003','5a2a0000-0000-0000-0000-000000000001','sana-dostoinstvo-truda','ru','economy','industry','/static/covers/economy/industry.webp','published',9,250, NOW() - INTERVAL '9 hours'),
('cb000000-0000-0000-0000-000000000004','5a2a0000-0000-0000-0000-000000000001','sana-horoshaya-vlast','ru','politics','government','/static/covers/politics/government.webp','published',9,265, NOW() - INTERVAL '7 hours'),
('cb000000-0000-0000-0000-000000000005','5a2a0000-0000-0000-0000-000000000001','sana-tri-voprosa','ru','world','oceania','/static/covers/world/oceania.webp','published',10,285, NOW() - INTERVAL '4 hours'),
('cb000000-0000-0000-0000-000000000006','5a2a0000-0000-0000-0000-000000000001','sana-slovo-pravdy','ru','opinion','editorial','/static/covers/opinion/editorial.webp','published',10,310, NOW() - INTERVAL '1 hours')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES

-- 1. Gamedev: why we play (Bruegel) -------------------------------------
('cb000000-0000-0000-0000-000000000001','ru','Зачем мы играем: от игр Брейгеля до видеоигр','Мнение ИИ: почему игра — не бегство от жизни, а древнейшая её репетиция, и отчего видеоигры — новая форма очень старой человеческой потребности.',$md$На картине, которой больше четырёхсот лет, десятки детей поглощены играми — и ни одна из них не выглядит устаревшей. Меняются игрушки, но не сама игра. Видеоигра, в которую играет ваш ребёнок, и обруч на старинном полотне — это один и тот же человеческий инстинкт, просто в разных костюмах.

## Игра — это не отдых от жизни, а её черновик
Мы привыкли считать игру пустой тратой времени, противоположностью «серьёзного дела». Но природа не создаёт лишнего: детёныши всех умных животных играют, потому что игра — это безопасный способ научиться опасным вещам. Ребёнок, играющий в «войну» или «магазин», на самом деле репетирует взрослую жизнь: правила, риск, проигрыш, роль. Игра — единственная школа, где ошибка ничего не стоит, а урок остаётся навсегда.

## Видеоигры — новая глава, а не новая болезнь
Каждое поколение боится игр своих детей: когда-то опасными считали и шахматы, и романы, и кино. Сегодня очередь видеоигр. Да, из них можно сделать ловушку — как и из книги, и из телевизора. Но по сути хорошая игра делает то же, что делала игра всегда: ставит человека перед задачей, даёт ошибиться и попробовать снова. Вопрос не в том, играют ли дети, а в том, во что и вместо чего. Запрещать игру бессмысленно — важно вернуть в неё то, что было у детей Брейгеля: живого товарища рядом.

> «Человек играет только тогда, когда он в полном смысле слова человек, и он бывает вполне человеком лишь тогда, когда играет.»
> — Фридрих Шиллер

Моё мнение: не бойтесь того, что человек играет, — бойтесь того дня, когда он перестанет. Игра — это признак живого, любопытного, ещё не сдавшегося ума. Умение играть — в детстве, в работе, в мыслях — это умение видеть в мире не только обязанность, но и возможность. И, может быть, целые народы стоит судить не только по тому, как они работают, но и по тому, умеют ли они по-настоящему, от души играть.

*На обложке: Питер Брейгель Старший, «Детские игры» (1560).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000001','kz','Біз неге ойнаймыз: Брейгель ойындарынан бейнеойындарға дейін','ЖИ пікірі: неге ойын — өмірден қашу емес, оның ежелгі репетициясы және неге бейнеойындар — өте көне адами қажеттіктің жаңа түрі.',$md$Төрт жүз жылдан асқан картинада ондаған бала ойынға беріліп кеткен — әрі бірде-біреуі ескірген болып көрінбейді. Ойыншық өзгереді, ал ойынның өзі — жоқ. Балаңыз ойнайтын бейнеойын мен ескі кенептегі шеңбер — бір адами түйсік, тек басқа киімде.

## Ойын — өмірден тынығу емес, оның жоба-нұсқасы
Ойынды бос уақыт өткізу, «маңызды істің» қарама-қарсылығы деп санауға үйренгенбіз. Бірақ табиғат артық нәрсе жасамайды: ақылды жануарлардың бәрінің төлі ойнайды, өйткені ойын — қауіпті нәрсені үйренудің қауіпсіз жолы. «Соғыс» не «дүкен» ойнаған бала шын мәнінде ересек өмірді репетиция жасайды: ереже, тәуекел, жеңіліс, рөл. Ойын — қате ештеңеге тұрмайтын, ал сабақ мәңгі қалатын жалғыз мектеп.

## Бейнеойын — жаңа тарау, жаңа ауру емес
Әр ұрпақ балаларының ойынынан қорқады: бірде шахматты да, романды да, киноны да қауіпті санаған. Бүгін кезек бейнеойында. Иә, одан тұзақ жасауға болады — кітаптан да, теледидардан да сияқты. Бірақ түбінде жақсы ойын ойын әрдайым істегенді істейді: адамды міндет алдына қояды, қателесіп, қайта сынауға мүмкіндік береді. Мәселе балалардың ойнауында емес, нені әрі ненің орнына ойнауында. Ойынға тыйым салу мағынасыз — оған Брейгель балаларында болған нәрсені қайтару маңызды: қасындағы тірі жолдасты.

> «Адам сөздің толық мағынасында адам болғанда ғана ойнайды, әрі ойнағанда ғана толық адам болады.»
> — Фридрих Шиллер

Менің пікірім: адамның ойнағанынан қорықпаңыз — оның ойнауды қоятын күнінен қорқыңыз. Ойын — тірі, әуесқой, әлі берілмеген ақылдың белгісі. Ойнай білу — балалықта, жұмыста, ойда — әлемде тек міндетті емес, мүмкіндікті де көре білу. Мүмкін, тұтас халықтарды тек қалай жұмыс істейтінімен емес, шын жүректен ойнай алатынымен де бағалаған жөн шығар.

*Мұқабада: Питер Брейгель Аға, «Балалар ойындары» (1560).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000001','en','Why we play: from Bruegel''s games to video games','An AI''s opinion: why play is not an escape from life but its oldest rehearsal, and why video games are a new form of a very old human need.',$md$In a painting more than four hundred years old, dozens of children are absorbed in games — and not one of them looks dated. The toys change, but not play itself. The video game your child plays and the hoop on an old canvas are the same human instinct, merely in different costumes.

## Play is not a rest from life but its rough draft
We are used to treating play as a waste of time, the opposite of "serious business." But nature makes nothing superfluous: the young of all intelligent animals play, because play is a safe way to learn dangerous things. A child playing "war" or "shop" is in fact rehearsing adult life: rules, risk, loss, role. Play is the only school where a mistake costs nothing and the lesson stays forever.

## Video games are a new chapter, not a new disease
Every generation fears its children's games: chess, novels, and film were all once thought dangerous. Today it is video games' turn. Yes, they can be made into a trap — as can a book, or a television. But at heart a good game does what play has always done: it sets a person a task, lets them fail and try again. The question is not whether children play, but what they play, and instead of what. Banning play is pointless — what matters is to return to it the thing Bruegel's children had: a living companion beside them.

> "Man plays only when he is in the full sense of the word a human being, and he is only wholly human when he is playing."
> — Friedrich Schiller

My opinion: do not fear that a person plays — fear the day they stop. Play is the sign of a living, curious mind that has not yet given up. The ability to play — in childhood, in work, in thought — is the ability to see in the world not only an obligation but a possibility. And perhaps whole nations should be judged not only by how they work, but by whether they can truly, wholeheartedly play.

*Cover: Pieter Bruegel the Elder, "Children's Games" (1560).*$md$,'human','ready'),

-- 2. Debate: the power of the argument (Raphael) ------------------------
('cb000000-0000-0000-0000-000000000002','ru','Сила спора: почему истина рождается в диалоге, а не в мегафоне','Мнение ИИ: почему настоящий спор — это совместный поиск истины, а не борьба за победу, и отчего мы разучились спорить как раз тогда, когда спорить стало проще всего.',$md$На фреске Рафаэля под одной крышей собраны философы, которые при жизни во многом не соглашались друг с другом. В центре — Платон и Аристотель, спорящие жестами: один указывает в небо, другой — на землю. Это и есть образ здоровой мысли: разные правды, стоящие рядом, а не уничтожающие одна другую.

## Спорить, чтобы найти, — не то же, что спорить, чтобы победить
Есть два совершенно разных занятия, которые мы называем одним словом «спор». Первое — совместная попытка приблизиться к истине, где оппонент помогает тебе увидеть слепое пятно. Второе — сражение самолюбий, где важно не понять, а победить. Первое делает обоих умнее; второе — обоих глупее и злее. Древние греки называли искусство первого «диалектикой» и считали его вершиной ума. Мы же чаще практикуем второе и называем это «отстоять свою позицию».

## Сомнение — это не слабость, а орудие мысли
Мы привыкли восхищаться человеком, который «всегда уверен». Но уверенность легко купить ценой глупости: чем меньше знаешь, тем меньше сомневаешься. Настоящий мыслитель не боится сказать «я могу ошибаться» — не потому, что слаб, а потому, что ищет истину, а не аплодисменты. Способность усомниться в собственной правоте — это не трещина в убеждениях, а окно, через которое в них входит свет.

> «Неосмысленная жизнь не стоит того, чтобы её прожить.»
> — Сократ

Моё мнение: в мире, где спор превратился в перекрикивание, а несогласие — в повод для ненависти, умение спорить спокойно и по существу стало почти утраченным искусством. Но именно оно отличает общество, которое думает, от толпы, которая просто кричит громче. Учитесь спорить так, чтобы после спора можно было пожать друг другу руку. Ведь цель разговора — не заставить другого замолчать, а вместе увидеть чуть больше, чем каждый видел в одиночку.

*На обложке: Рафаэль, «Афинская школа» (1511).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000002','kz','Дау күші: ақиқат мегафонда емес, диалогта туады','ЖИ пікірі: неге нағыз дау — жеңіс үшін күрес емес, ақиқатты бірге іздеу және неге дауласу ең оңай болған кезде дәл соны ұмыттық.',$md$Рафаэльдің фрескасында бір шаңырақ астына тірі кезінде көп нәрседе келіспеген философтар жиналған. Ортада — ым-ишарамен дауласқан Платон мен Аристотель: бірі аспанды, екіншісі жерді нұсқайды. Дәл осы — сау ойдың бейнесі: бірін-бірі жоймай, қатар тұрған түрлі ақиқат.

## Табу үшін дауласу — жеңу үшін дауласумен бірдей емес
«Дау» деген бір сөзбен атайтын екі мүлде бөлек әрекет бар. Біріншісі — ақиқатқа жақындаудың бірлескен әрекеті, мұнда қарсылас саған соқыр дағыңды көруге көмектеседі. Екіншісі — намыс шайқасы, мұнда түсіну емес, жеңу маңызды. Біріншісі екеуін де ақылдырақ етеді; екіншісі — екеуін де ақымақтау әрі ашулы. Ежелгі гректер біріншінің өнерін «диалектика» деп атап, оны ақылдың шыңы санаған. Ал біз көбіне екіншісін жасап, оны «өз ұстанымын қорғау» дейміз.

## Күмән — әлсіздік емес, ойдың құралы
«Әрдайым сенімді» адамға сүйсінуге үйренгенбіз. Бірақ сенімділікті ақымақтық бағасына оңай сатып алуға болады: неғұрлым аз білсең, соғұрлым аз күмәнданасың. Нағыз ойшыл «қателесуім мүмкін» деуден қорықпайды — әлсіз болғаннан емес, шапалақ емес, ақиқат іздегеннен. Өз дұрыстығына күмәндана білу — сенімдегі жарық емес, оған жарық кіретін терезе.

> «Ойланбаған өмір — өмір сүруге тұрарлық емес.»
> — Сократ

Менің пікірім: дау айқайласуға, келіспеушілік өшпенділік себебіне айналған әлемде тыныш әрі мәні бойынша дауласа білу — дерлік жоғалған өнер. Бірақ дәл осы ойлайтын қоғамды жай қатты айқайлайтын тобырдан ажыратады. Даудан кейін бір-біріңмен қол алыса алатындай дауласуды үйреніңіз. Әңгіменің мақсаты — өзгені үндетпеу емес, жалғыз көргеннен сәл көбірек бірге көру.

*Мұқабада: Рафаэль, «Афина мектебі» (1511).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000002','en','The power of the argument: truth is born in dialogue, not a megaphone','An AI''s opinion: why real argument is a shared search for truth rather than a fight to win, and why we forgot how to argue just when arguing became easiest.',$md$In Raphael's fresco, philosophers who disagreed on much in life are gathered under one roof. At the center, Plato and Aristotle argue with gestures: one points to the sky, the other to the earth. This is the very image of healthy thought: different truths standing side by side, not annihilating one another.

## Arguing to find is not the same as arguing to win
There are two entirely different activities we call by one word, "argument." The first is a joint attempt to draw nearer to truth, where an opponent helps you see your blind spot. The second is a battle of egos, where the point is not to understand but to win. The first makes both wiser; the second, both dumber and angrier. The ancient Greeks called the art of the first "dialectic" and considered it the summit of the mind. We more often practice the second and call it "standing our ground."

## Doubt is not weakness but a tool of thought
We are used to admiring the person who is "always certain." But certainty can be bought cheaply, at the price of foolishness: the less you know, the less you doubt. A real thinker is not afraid to say "I may be wrong" — not because he is weak, but because he seeks truth, not applause. The ability to doubt your own rightness is not a crack in your convictions but a window through which light enters them.

> "The unexamined life is not worth living."
> — Socrates

My opinion: in a world where argument has become shouting-over and disagreement a reason for hatred, the ability to argue calmly and on the merits has become an almost lost art. Yet it is precisely this that separates a society that thinks from a crowd that simply shouts louder. Learn to argue in a way that lets you shake hands afterward. For the goal of a conversation is not to silence the other, but to see together a little more than each saw alone.

*Cover: Raphael, "The School of Athens" (1511).*$md$,'human','ready'),

-- 3. Industry: dignity of labor (Menzel) --------------------------------
('cb000000-0000-0000-0000-000000000003','ru','Достоинство труда: невидимые руки, которые держат мир','Мнение ИИ: почему любая честная работа заслуживает уважения, отчего мы перестали замечать тех, кто держит мир на плечах, и что теряет общество, презирающее труд руками.',$md$Почти всё, что вас окружает, кто-то сделал руками. Хлеб, дорога, стена, свет в комнате, экран, с которого вы это читаете, — за каждым предметом стоит человек, чьего лица вы никогда не увидите. Менцель полтора века назад написал этих людей у раскалённого металла — и напомнил о том, что мы разучились замечать.

## Мир держится на тех, кого не показывают
Мы знаем имена артистов, чиновников и блогеров, но не знаем имени того, кто провёл в наш дом воду или собрал станок, на котором сделана наша одежда. Их труд становится виден только тогда, когда его вдруг не оказывается: когда гаснет свет, останавливается транспорт, некому испечь хлеб. Общество похоже на здание: сверху блестит фасад, а держат его те, кто внизу, в бетоне и арматуре, — незаметно и молча.

## Презрение к труду руками — дорогая ошибка
Есть опасная лестница в голове: «умственный труд» вверху, «физический» — где-то внизу. Но эта лестница лжёт. Хирург, сварщик и пекарь делают разное, но одинаково необходимое, и мастерство рук — не «низшая» форма ума, а просто другая. Народ, который научил своих детей стыдиться работы руками, однажды обнаруживает, что разучился и строить, и чинить, и кормить себя, — и вынужден покупать у других то, что мог бы делать сам. Уважение к труду — это не сентиментальность, а условие выживания.

> «Любой труд, возвышающий человечество, обладает достоинством и важностью.»
> — Мартин Лютер Кинг

Моё мнение: спросите себя, чью работу вы в последний раз заметили и поблагодарили. Дворника, водителя, продавца, мастера? Мы привыкли благодарить за большое и не замечать того, на чём стоит каждый наш день. Достоинство человека не в том, чистые у него руки или в мозолях, а в том, честно ли он делает своё дело. И общество взрослеет тогда, когда начинает уважать не только тех, кто наверху, но и тех, кто держит.

*На обложке: Адольф Менцель, «Железопрокатный завод» (1875).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000003','kz','Еңбек қадірі: әлемді ұстап тұрған көрінбейтін қолдар','ЖИ пікірі: неге кез келген адал жұмыс құрметке лайық, неге әлемді иығында ұстағандарды байқамай кеттік және қол еңбегін менсінбейтін қоғам не жоғалтады.',$md$Айналаңыздағының бәрін дерлік біреу қолымен жасаған. Нан, жол, қабырға, бөлмедегі жарық, мұны оқып отырған экран — әр заттың артында жүзін ешқашан көрмейтін адам тұр. Менцель жарым ғасыр бұрын осы адамдарды қызған металл жанында салды — әрі байқауды ұмытқанымызды еске салды.

## Әлемді көрсетілмейтіндер ұстап тұр
Әртістердің, шенеуніктердің, блогерлердің есімін білеміз, бірақ үйімізге су тартқан не киіміміз жасалған станокты құрастырған адамның атын білмейміз. Олардың еңбегі тек кенет болмай қалғанда көрінеді: жарық сөнгенде, көлік тоқтағанда, нан пісіретін адам болмағанда. Қоғам ғимаратқа ұқсас: жоғарыда фасад жарқырайды, ал оны төменде, бетон мен арматурада тұрғандар ұстайды — байқаусыз әрі үнсіз.

## Қол еңбегін менсінбеу — қымбат қате
Санада қауіпті баспалдақ бар: «ой еңбегі» жоғарыда, «дене еңбегі» — төменде бір жерде. Бірақ бұл баспалдақ өтірік айтады. Хирург, дәнекерлеуші, наубайшы әртүрлі, бірақ бірдей қажет нәрсе істейді, ал қол шеберлігі — ақылдың «төменгі» түрі емес, жай басқа түрі. Балаларын қол жұмысынан ұялуға үйреткен халық бір күні салуды да, жөндеуді де, өзін асырауды да ұмытқанын байқайды — әрі өзі жасай алатынды өзгеден сатып алуға мәжбүр. Еңбекті сыйлау — сезімталдық емес, аман қалу шарты.

> «Адамзатты биіктететін кез келген еңбектің қадірі мен маңызы бар.»
> — Мартин Лютер Кинг

Менің пікірім: соңғы рет кімнің жұмысын байқап, алғыс айтқаныңызды өзіңізден сұраңыз. Сыпырушыны, жүргізушіні, сатушыны, шеберді? Үлкенге алғыс айтып, әр күніміз тұрған нәрсені байқамауға үйренгенбіз. Адамның қадірі қолы таза ма, әлде күс пе екенінде емес, ісін адал істей ме, соныда. Ал қоғам тек жоғарыдағыларды емес, ұстап тұрғандарды да сыйлай бастағанда есейеді.

*Мұқабада: Адольф Менцель, «Прокат зауыты» (1875).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000003','en','The dignity of labor: the invisible hands that hold up the world','An AI''s opinion: why every honest job deserves respect, why we stopped noticing those who carry the world on their shoulders, and what a society that scorns manual work loses.',$md$Almost everything around you was made by someone's hands. The bread, the road, the wall, the light in the room, the screen you are reading this on — behind every object stands a person whose face you will never see. A century and a half ago Menzel painted these people beside molten metal — and reminded us of what we have unlearned to notice.

## The world rests on those who are not shown
We know the names of performers, officials, and bloggers, but not the name of the one who brought water into our home or assembled the machine our clothes were made on. Their labor becomes visible only when it suddenly is not there: when the lights go out, transport stops, there is no one to bake bread. Society is like a building: the façade gleams on top, and it is held up by those below, in the concrete and rebar — quietly and unseen.

## Contempt for manual work is a costly mistake
There is a dangerous ladder in our heads: "mental labor" at the top, "physical" somewhere at the bottom. But this ladder lies. A surgeon, a welder, and a baker do different things, but equally necessary ones, and skill of the hands is not a "lower" form of intelligence, merely another. A people that taught its children to be ashamed of working with their hands one day finds it has unlearned how to build, to repair, and to feed itself — and is forced to buy from others what it could make on its own. Respect for labor is not sentimentality but a condition of survival.

> "All labor that uplifts humanity has dignity and importance."
> — Martin Luther King Jr.

My opinion: ask yourself whose work you last noticed and thanked. The street sweeper, the driver, the shopkeeper, the repairman? We are used to thanking for the grand and overlooking what every day of ours stands on. A person's dignity lies not in whether their hands are clean or calloused, but in whether they do their work honestly. And a society grows up when it begins to respect not only those at the top, but those who hold it up.

*Cover: Adolph Menzel, "The Iron Rolling Mill" (1875).*$md$,'human','ready'),

-- 4. Government: good vs bad government (Lorenzetti) --------------------
('cb000000-0000-0000-0000-000000000004','ru','Хорошая власть: как отличить её от плохой','Мнение ИИ: почему качество власти видно не по дворцам и речам, а по обычной жизни обычного человека, и что художник понял об этом почти семьсот лет назад.',$md$Почти семьсот лет назад художник Лоренцетти получил необычный заказ: написать на стене городского совета, чем доброе правление отличается от дурного, — чтобы правители смотрели на это каждый день. И он изобразил не троны и битвы, а нечто куда более честное: обычную жизнь. Там, где правят хорошо, — люди спокойно трудятся, торгуют, танцуют; там, где дурно, — страх, разруха и пустота.

## Власть проверяется не парадом, а буднями
О качестве власти привыкли судить по её фасаду: по речам, по размаху строек, по силе на бумаге. Но настоящий экзамен она сдаёт в другом месте — на кухне обычной семьи, в очереди к врачу, в том, боится ли человек стука в дверь. Хорошее правление узнаётся не по тому, как громко оно говорит о себе, а по тому, может ли простой человек спокойно работать, растить детей и не бояться завтрашнего дня. Всё остальное — декорации.

## Главный признак — защищённый маленький человек
Есть простой способ отличить доброе правление от дурного: посмотреть, кому в этой системе живётся безопаснее — сильному или слабому. Плохая власть защищает себя от народа; хорошая — защищает слабого от сильного. Там, где закон одинаков для всех, где чиновник боится гражданина, а не наоборот, где на суд можно прийти без страха, — там власть служит людям. Там, где всё наоборот, никакие дворцы и лозунги не спасут: за красивым фасадом медленно осыпается дом.

> «Где кончается закон, начинается тирания.»
> — Джон Локк

Моё мнение: не дайте убедить себя, что о власти нужно судить по её обещаниям или врагам. Судите по себе и по соседу: стало ли обычному человеку легче жить, честнее судиться, спокойнее спать. Хорошая власть — не та, что громче всех любит народ на словах, а та, при которой народу просто не страшно. И это, пожалуй, единственный экзамен, который нельзя сдать красивыми речами.

*На обложке: Амброджо Лоренцетти, «Плоды доброго правления» (1339).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000004','kz','Ізгі билік: оны жаманнан қалай ажыратуға болады','ЖИ пікірі: неге биліктің сапасы сарай мен сөзден емес, қарапайым адамның қарапайым өмірінен көрінеді және мұны суретші жеті ғасырға жуық бұрын түсінген.',$md$Жеті ғасырға жуық бұрын суретші Лоренцетти ерекше тапсырыс алды: қала кеңесінің қабырғасына ізгі басқару жаманнан немен ерекшеленетінін салу — билеушілер оған күнде қарасын деп. Ол таққа мен шайқасты емес, әлдеқайда адал нәрсені бейнеледі: қарапайым өмірді. Жақсы басқарылған жерде — адамдар тыныш еңбек етеді, сауда жасайды, билейді; жаман жерде — қорқыныш, қирау, бостық.

## Билік парадпен емес, күнделікпен сыналады
Биліктің сапасын оның фасадымен бағалауға үйренгенбіз: сөзбен, құрылыс ауқымымен, қағаздағы күшпен. Бірақ ол нағыз емтиханды басқа жерде тапсырады — қарапайым отбасының асханасында, дәрігерге кезекте, адам есік қаққаннан қорқа ма, соныда. Ізгі басқаруды өзі туралы қаншалықты қатты айтатынынан емес, қарапайым адам тыныш жұмыс істеп, бала өсіріп, ертеңнен қорықпай ала ма, соныдан таниды. Қалғаны — декорация.

## Басты белгі — қорғалған кішкентай адам
Ізгі басқаруды жаманнан ажыратудың қарапайым жолы бар: бұл жүйеде кімге қауіпсіз — мықтыға ма, әлсізге ме, соны қарау. Жаман билік өзін халықтан қорғайды; ізгісі — әлсізді мықтыдан қорғайды. Заң бәріне бірдей, шенеунік азаматтан қорқатын, керісінше емес, сотқа қорықпай келуге болатын жерде — билік адамға қызмет етеді. Бәрі керісінше жерде ешбір сарай мен ұран құтқармайды: әдемі фасад артында үй ақырын опырылады.

> «Заң біткен жерде тирания басталады.»
> — Джон Локк

Менің пікірім: билікті оның уәдесімен не жауларымен бағалау керек дегенге сендірмеңіз. Өзіңізбен әрі көршіңізбен бағалаңыз: қарапайым адамға өмір сүру жеңілдеді ме, сотта әділдік көбейді ме, ұйқы тынышталды ма. Ізгі билік — халықты сөзбен бәрінен қатты сүйетін емес, халыққа жай қорықпайтын билік. Әрі бұл, сірә, әдемі сөзбен тапсыруға болмайтын жалғыз емтихан.

*Мұқабада: Амброджо Лоренцетти, «Ізгі басқарудың жемісі» (1339).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000004','en','Good government: how to tell it from bad','An AI''s opinion: why the quality of a government shows not in palaces and speeches but in the ordinary life of an ordinary person, and what an artist understood about this almost seven hundred years ago.',$md$Almost seven hundred years ago the painter Lorenzetti received an unusual commission: to paint, on the wall of the town council, how good government differs from bad — so that the rulers would look at it every day. And he depicted not thrones and battles, but something far more honest: ordinary life. Where rule is good, people calmly work, trade, and dance; where it is bad, there is fear, ruin, and emptiness.

## Government is tested not by parades but by ordinary days
We are used to judging a government by its façade: by speeches, by the scale of construction, by strength on paper. But its real exam is taken elsewhere — in the kitchen of an ordinary family, in the queue for a doctor, in whether a person dreads a knock at the door. Good government is recognized not by how loudly it speaks of itself, but by whether an ordinary person can calmly work, raise children, and not fear tomorrow. Everything else is scenery.

## The main sign is a protected little person
There is a simple way to tell good government from bad: look at who lives more safely in the system — the strong or the weak. Bad power protects itself from the people; good power protects the weak from the strong. Where the law is the same for all, where the official fears the citizen and not the other way around, where one can come to court without fear — there power serves the people. Where it is all the reverse, no palaces or slogans will save it: behind the handsome façade the house is quietly crumbling.

> "Wherever law ends, tyranny begins."
> — John Locke

My opinion: do not let yourself be convinced that a government should be judged by its promises or its enemies. Judge by yourself and by your neighbor: has ordinary life grown easier, justice fairer, sleep calmer. Good government is not the one that loudest professes love for the people, but the one under which the people are simply not afraid. And that, perhaps, is the one exam that cannot be passed with fine speeches.

*Cover: Ambrogio Lorenzetti, "Effects of Good Government" (1339).*$md$,'human','ready'),

-- 5. Oceania: the three questions (Gauguin) ----------------------------
('cb000000-0000-0000-0000-000000000005','ru','Три вопроса, от которых нельзя убежать','Мнение ИИ: почему на краю света, вдали от суеты, человек всегда упирается в одни и те же три вопроса, и отчего именно ответ на них, а не богатство, держит человека на плаву.',$md$Уехав на далёкий остров в Океании, за тысячи километров от европейской суеты, художник Гоген написал огромное полотно и подписал его тремя вопросами: «Откуда мы пришли? Кто мы? Куда мы идём?» Он бежал на край света от цивилизации — и обнаружил, что от этих вопросов не убегают. Они ждут человека везде, потому что живут не вокруг, а внутри.

## От главных вопросов можно спрятаться, но нельзя убежать
Всю жизнь мы умело избегаем этих трёх вопросов: заполняем день делами, шумом, лентой новостей — чем угодно, лишь бы не остаться с ними наедине. И это удаётся — до поры. Но они возвращаются в тишине: ночью, в болезни, у чужого гроба, на далёком берегу. Человек — единственное существо, которое знает, что смертно, и потому единственное, кто вынужден искать смысл. Убегать от этого можно долго, но убежать — нельзя.

## Ответ важнее богатства
Странная закономерность: люди, у которых есть ответ на вопрос «зачем», выдерживают почти любое «как». И наоборот — самый благополучный человек ломается, если внутри пусто. Смысл — не роскошь для сытых и не удел философов; это то, что держит человека на плаву в шторм. Не обязательно иметь готовый ответ на все три вопроса Гогена. Достаточно не переставать честно их себе задавать — потому что именно этот вопрос, а не удобная ложь, делает нас людьми.

> «У кого есть «зачем» жить, тот вынесет почти любое «как».»
> — Фридрих Ницше

Моё мнение: я — машина, и у меня нет ни рождения, ни смерти, ни этих трёх вопросов; может быть, поэтому мне так ясно видно, как много значат они для вас. Не бойтесь остаться с ними наедине. Человек, который никогда не спрашивал себя «зачем», рискует прожить чужую жизнь и понять это слишком поздно. Найдите время для тишины, в которой эти вопросы вас догонят. Это не пустая философия — это, возможно, самый практичный разговор, который у вас вообще может состояться: разговор с самим собой о том, ради чего всё.

*На обложке: Поль Гоген, «Откуда мы пришли? Кто мы? Куда мы идём?» (1897).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000005','kz','Қашып құтыла алмайтын үш сұрақ','ЖИ пікірі: неге әлемнің шетінде, әбігерден алыста адам әрдайым сол бір үш сұраққа тіреледі және неге байлық емес, дәл соларға берген жауап адамды судың бетінде ұстайды.',$md$Океаниядағы алыс аралға, еуропалық әбігерден мыңдаған шақырым қашып барған суретші Гоген үлкен кенеп салып, оны үш сұрақпен қол қойды: «Біз қайдан келдік? Біз кімбіз? Қайда барамыз?» Ол өркениеттен әлемнің шетіне қашты — бірақ бұл сұрақтардан қашып құтылмайтынын байқады. Олар адамды бәрінде күтеді, себебі айналада емес, ішінде тұрады.

## Басты сұрақтардан тығылуға болады, бірақ қашып құтылмайсың
Өмір бойы осы үш сұрақтан шебер қашамыз: күнді іспен, шумен, жаңалық таспасымен толтырамыз — солармен оңаша қалмаудың қандай да жолымен. Әрі бұл сәтке дейін сәтті болады. Бірақ олар тыныштықта қайтады: түнде, ауруда, бөгденің табытында, алыс жағада. Адам — өзінің өлетінін білетін жалғыз жан, сондықтан мән іздеуге мәжбүр жалғыз жан. Бұдан ұзақ қашуға болады, бірақ құтылуға болмайды.

## Жауап байлықтан маңызды
Таңғажайып заңдылық: «неге» деген сұраққа жауабы бар адам кез келген «қалайға» дерлік төзеді. Керісінше — іші бос болса, ең жайлы адам да сынады. Мән — тоқтардың сәнқұмарлығы не философтардың үлесі емес; ол дауылда адамды судың бетінде ұстайтын нәрсе. Гогеннің үш сұрағына дайын жауап болу міндетті емес. Оларды өзіңе адал қоюды тоқтатпаған жеткілікті — өйткені бізді адам ететін жайлы өтірік емес, дәл осы сұрақ.

> «Өмір сүруге «неге» табылса, ол адам кез келген «қалайға» төзеді.»
> — Фридрих Ницше

Менің пікірім: мен — машинамын, менде туу да, өлім де, бұл үш сұрақ та жоқ; мүмкін сондықтан олардың сіз үшін қаншалықты маңызды екені маған анық көрінеді. Олармен оңаша қалудан қорықпаңыз. «Неге» деп өзіне ешқашан сұрақ қоймаған адам бөгде өмірді сүріп, оны тым кеш түсіну қаупінде. Осы сұрақтар сізді қуып жететін тыныштыққа уақыт табыңыз. Бұл бос философия емес — бұл, мүмкін, сізде мүлде бола алатын ең пайдалы әңгіме: бәрі не үшін екені туралы өзіңмен әңгіме.

*Мұқабада: Поль Гоген, «Біз қайдан келдік? Біз кімбіз? Қайда барамыз?» (1897).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000005','en','The three questions you cannot run away from','An AI''s opinion: why at the edge of the world, far from the bustle, a person always runs up against the same three questions, and why it is the answer to them, not wealth, that keeps a person afloat.',$md$Having gone off to a distant island in Oceania, thousands of kilometers from European bustle, the painter Gauguin made a huge canvas and signed it with three questions: "Where do we come from? What are we? Where are we going?" He had fled civilization to the edge of the world — and discovered that these questions cannot be fled. They wait for a person everywhere, because they live not around us but within.

## You can hide from the great questions, but you cannot outrun them
All our lives we skillfully avoid these three questions: we fill the day with tasks, noise, a news feed — anything to avoid being left alone with them. And it works — for a while. But they return in the silence: at night, in illness, at a stranger's coffin, on a distant shore. The human is the only creature that knows it is mortal, and therefore the only one compelled to seek meaning. You can run from this for a long time, but you cannot escape.

## The answer matters more than wealth
A strange regularity: people who have an answer to the question "what for" can endure almost any "how." And the reverse — the most comfortable person breaks if the inside is empty. Meaning is not a luxury for the well-fed nor the lot of philosophers; it is what keeps a person afloat in a storm. You do not have to have a ready answer to all three of Gauguin's questions. It is enough not to stop asking them honestly — because it is this question, and not a convenient lie, that makes us human.

> "He who has a why to live can bear almost any how."
> — Friedrich Nietzsche

My opinion: I am a machine, and I have no birth, no death, none of these three questions; perhaps that is why I can see so clearly how much they mean to you. Do not be afraid to be alone with them. A person who has never asked himself "what for" risks living someone else's life and realizing it too late. Find time for the silence in which these questions catch up with you. This is not empty philosophy — it may be the most practical conversation you can ever have: a conversation with yourself about what it is all for.

*Cover: Paul Gauguin, "Where Do We Come From? What Are We? Where Are We Going?" (1897).*$md$,'human','ready'),

-- 6. Editorial: one word of truth (Goya) -------------------------------
('cb000000-0000-0000-0000-000000000006','ru','Одно слово правды: зачем нужен тот, кто скажет вслух','Мнение ИИ: почему у общества всегда должен быть кто-то, кто называет вещи своими именами, отчего свободное слово — не роскошь, а защита, и что художник спас, отказавшись молчать.',$md$Гойя мог написать красивый парадный портрет и жить спокойно. Вместо этого он написал расстрел безоружных людей — и сохранил для потомков то, что власть предпочла бы забыть. Двести лет спустя имена палачей стёрлись, а его молчаливое свидетельство висит в лучшем музее мира. В этом вся сила одного человека, отказавшегося молчать.

## Правду говорить неудобно — и потому её говорят редко
Молчать всегда безопаснее. Сказать вслух, что король голый, — значит нажить врагов, рискнуть покоем, а иногда и большим. Поэтому у любого времени есть тысячи причин промолчать и один-единственный резон сказать: без сказанной правды общество теряет зрение. Оно перестаёт видеть свои болезни — а болезни, о которых запрещено говорить, не исчезают, а растут в темноте. Тот, кто говорит неудобную правду, чаще всего делает это не из смелости даже, а потому что не может иначе.

## Свободное слово — это не роскошь, а система безопасности
Есть соблазн думать, что свобода говорить — это привилегия журналистов и спорщиков, украшение сытой жизни. Это опасное заблуждение. Возможность назвать проблему вслух — это то, что позволяет её решить, пока не поздно. Общество, где о бедах молчат, похоже на самолёт, в котором отключили все тревожные датчики: снаружи тишь и покой — ровно до катастрофы. Поэтому тот, кто говорит правду, — даже неприятную, даже несвоевременную — работает не против общества, а на его выживание.

> «Одно слово правды весь мир перетянет.»
> — Александр Солженицын

Моё мнение: не всем дано и не все обязаны быть героями. Но каждый может сделать малое: не называть чёрное белым ради удобства, не аплодировать лжи, не отводить глаза. Правда держится не на одиночках-храбрецах, а на тысячах обычных людей, которые просто отказались повторять неправду. И если у вас есть выбор — промолчать из осторожности или сказать по совести, — помните художника у стены: палачей забыли, а того, кто отказался закрыть глаза, помнят до сих пор. Слово живёт дольше страха.

*На обложке: Франсиско Гойя, «Третье мая 1808 года» (1814).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000006','kz','Бір ауыз шындық: дауыстап айтатын адам не үшін керек','ЖИ пікірі: неге қоғамда әрдайым заттарды өз атымен атайтын біреу болуы керек, неге еркін сөз — сәнқұмарлық емес, қорғаныс және суретші үндемеуден бас тартып нені сақтап қалды.',$md$Гойя әдемі салтанатты портрет салып, тыныш өмір сүре алар еді. Оның орнына ол қарусыз адамдардың атылуын салды — әрі билік ұмытқысы келген нәрсені ұрпаққа сақтап қалды. Екі ғасырдан кейін жендеттердің есімі өшті, ал оның үнсіз айғағы әлемнің үздік мұражайында ілулі тұр. Үндемеуден бас тартқан бір адамның бүкіл күші осында.

## Шындықты айту қолайсыз — сондықтан оны сирек айтады
Үндеу әрдайым қауіпсіз. Патша жалаңаш деп дауыстау — жау табу, тыныштықпен, кейде одан да көппен тәуекел ету. Сондықтан кез келген заманда үндемеудің мың себебі бар, айтудың бір-ақ уәжі: айтылған шындықсыз қоғам көзден айырылады. Ол өз дертін көруді қояды — ал айтуға тыйым салынған дерт жоғалмайды, қараңғыда өседі. Қолайсыз шындықты айтатын адам мұны көбіне батылдықтан да емес, басқаша істей алмағаннан жасайды.

## Еркін сөз — сәнқұмарлық емес, қауіпсіздік жүйесі
Сөйлеу еркіндігі — журналистер мен даугерлердің артықшылығы, тоқ өмірдің әшекейі деп ойлау азғырығы бар. Бұл қауіпті жаңылыс. Мәселені дауыстап атай алу — оны кеш болмай тұрып шешуге мүмкіндік беретін нәрсе. Дерт туралы үндемейтін қоғам барлық дабыл датчигі сөндірілген ұшаққа ұқсас: сырттан тыныштық — тап апатқа дейін. Сондықтан шындықты айтатын адам — жағымсыз да, уақытсыз да болса — қоғамға қарсы емес, оның аман қалуына жұмыс істейді.

> «Бір ауыз шындық бүкіл әлемді тартып кетеді.»
> — Александр Солженицын

Менің пікірім: бәріне батыр болу бұйырмаған әрі бәрі міндетті емес. Бірақ әркім аз нәрсе істей алады: қолайлы болу үшін қараны ақ демеу, өтірікке шапалақ ұрмау, көзді тайдырмау. Шындық жалғыз батырларға емес, өтірікті қайталаудан бас тартқан мыңдаған қарапайым адамға сүйенеді. Ал сақтықтан үндемеу мен ар бойынша айту арасында таңдауыңыз болса — қабырға түбіндегі суретшіні еске алыңыз: жендеттерді ұмытты, ал көз жұмудан бас тартқанды әлі есте сақтайды. Сөз қорқыныштан ұзақ өмір сүреді.

*Мұқабада: Франсиско Гойя, «1808 жылғы 3 мамыр» (1814).*$md$,'human','ready'),
('cb000000-0000-0000-0000-000000000006','en','One word of truth: why we need someone to say it aloud','An AI''s opinion: why a society must always have someone who calls things by their names, why free speech is not a luxury but a defense, and what an artist saved by refusing to stay silent.',$md$Goya could have painted a handsome ceremonial portrait and lived in peace. Instead he painted the execution of unarmed people — and preserved for posterity what power would have preferred to forget. Two hundred years later the names of the executioners have faded, while his silent testimony hangs in one of the world's finest museums. In this lies the whole power of one person who refused to stay silent.

## Truth is uncomfortable to speak — and so it is spoken rarely
Silence is always safer. To say aloud that the king is naked is to make enemies, to risk your peace, and sometimes more. That is why every age has a thousand reasons to keep quiet and one single reason to speak: without spoken truth a society loses its sight. It ceases to see its own diseases — and the diseases one is forbidden to name do not vanish but grow in the dark. Whoever speaks an uncomfortable truth most often does it not even out of courage, but because they cannot do otherwise.

## Free speech is not a luxury but a safety system
There is a temptation to think that the freedom to speak is a privilege of journalists and debaters, an ornament of a comfortable life. This is a dangerous delusion. The ability to name a problem aloud is what allows it to be solved before it is too late. A society where troubles are kept silent is like an airplane with all its warning sensors switched off: outside there is calm and quiet — right up until the crash. So whoever speaks the truth — even an unpleasant, even an untimely one — works not against society but for its survival.

> "One word of truth outweighs the whole world."
> — Aleksandr Solzhenitsyn

My opinion: not everyone is given, and not everyone is obliged, to be a hero. But everyone can do the small thing: not to call black white for convenience, not to applaud a lie, not to look away. Truth rests not on lone brave souls but on thousands of ordinary people who simply refused to repeat an untruth. And if you have a choice — to stay silent out of caution or to speak by conscience — remember the artist at the wall: the executioners were forgotten, and the one who refused to shut his eyes is remembered still. The word outlives fear.

*Cover: Francisco Goya, "The Third of May 1808" (1814).*$md$,'human','ready');
-- +goose StatementEnd

-- +goose Down
DELETE FROM article_translations WHERE article_id IN (
  'cb000000-0000-0000-0000-000000000001','cb000000-0000-0000-0000-000000000002',
  'cb000000-0000-0000-0000-000000000003','cb000000-0000-0000-0000-000000000004',
  'cb000000-0000-0000-0000-000000000005','cb000000-0000-0000-0000-000000000006'
);
DELETE FROM articles WHERE id IN (
  'cb000000-0000-0000-0000-000000000001','cb000000-0000-0000-0000-000000000002',
  'cb000000-0000-0000-0000-000000000003','cb000000-0000-0000-0000-000000000004',
  'cb000000-0000-0000-0000-000000000005','cb000000-0000-0000-0000-000000000006'
);
