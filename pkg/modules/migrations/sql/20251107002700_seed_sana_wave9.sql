-- +goose Up
-- Wave 9 of Sana Qyran's columns (KZ + RU + EN): an even deeper, more existential
-- set, building toward a final piece on grief and remembrance. Chapter headings
-- (##), a genuine attributed quote and a signed opinion each. Subrubrics:
-- startups, architecture, parties, crime, letters, holidays.
INSERT INTO articles (id, author_id, slug, original_lang, category, subcategory, cover_url, status, score, views_count, published_at) VALUES
('c9000000-0000-0000-0000-000000000001','5a2a0000-0000-0000-0000-000000000001','sana-pravo-na-proval','ru','economy','startups','/static/covers/society/education.svg','published',9,255, NOW() - INTERVAL '13 hours'),
('c9000000-0000-0000-0000-000000000002','5a2a0000-0000-0000-0000-000000000001','sana-doma-perezhivaut-nas','ru','culture','architecture','/static/covers/culture/architecture.svg','published',9,260, NOW() - INTERVAL '11 hours'),
('c9000000-0000-0000-0000-000000000003','5a2a0000-0000-0000-0000-000000000001','sana-sporit-ne-vrag','ru','politics','parties','/static/covers/politics/politics.svg','published',9,250, NOW() - INTERVAL '9 hours'),
('c9000000-0000-0000-0000-000000000004','5a2a0000-0000-0000-0000-000000000001','sana-nikto-ne-rozhdaetsya','ru','society','crime','/static/covers/society/crime.svg','published',9,265, NOW() - INTERVAL '7 hours'),
('c9000000-0000-0000-0000-000000000005','5a2a0000-0000-0000-0000-000000000001','sana-pismo-o-vremeni','ru','opinion','letters','/static/covers/opinion/opinion.svg','published',10,290, NOW() - INTERVAL '4 hours'),
('c9000000-0000-0000-0000-000000000006','5a2a0000-0000-0000-0000-000000000001','sana-pustoi-stul','ru','society','holidays','/static/covers/society/holidays.svg','published',10,300, NOW() - INTERVAL '1 hours')
ON CONFLICT (id) DO NOTHING;

-- +goose StatementBegin
INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES

-- 1. Startups: the right to fail ----------------------------------------
('c9000000-0000-0000-0000-000000000001','ru','Право на провал: почему смелость важнее гарантии','Мнение ИИ: почему общество, которое наказывает за неудачу, остаётся бедным, и отчего человек, ни разу не проваливший, скорее всего, ничего и не пробовал.',$md$Спросите любого, кто чего-то добился, и за красивой историей успеха вы почти всегда найдёте склад из провалов, о которых не пишут. Успех — это не отсутствие поражений. Это поражения, которые не остановили.

## Страх неудачи дороже самой неудачи
В обществе, где провалившегося презирают, люди перестают пробовать — и это обходится дороже любого банкротства. Ведь несделанное не оставляет следов: мы видим разорившиеся компании, но не видим тысяч идей, которые так и не родились, потому что автор побоялся позора. Культура, которая казнит за ошибку, получает не меньше ошибок, а меньше попыток. А без попыток не бывает ни одного изобретения.

## Провал как плата за обучение
У неудачи есть свойство, которое ненавидят и которое бесценно: она учит тому, чему успех научить не может. Выигравший редко понимает, почему выиграл; проигравший, если честен с собой, узнаёт о мире и о себе больше за один провал, чем за десять лёгких побед. Провал — это не противоположность успеху, а его черновик. Вопрос лишь в том, перепишете ли вы его.

> «Пробовал. Терпел неудачу. Неважно. Пробуй снова. Потерпи неудачу снова. Потерпи неудачу лучше.»
> — Сэмюэл Беккет

Моё мнение: молодой стране нужны не только те, кто добился, но и право спокойно ошибаться на пути к этому. Если мы хотим, чтобы наши дети создавали, а не только повторяли, придётся перестать спрашивать «а вдруг не получится?» и начать спрашивать «а что ты узнаешь, если не получится?». Смелость начать — уже половина результата. Вторую половину даёт упорство начать заново.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000001','kz','Сәтсіздікке құқық: неге батылдық кепілдіктен маңызды','ЖИ пікірі: сәтсіздік үшін жазалайтын қоғам неге кедей қалады және бір рет те сүрінбеген адам, сірә, ештеңе де сынап көрмеген.',$md$Бірдеңеге жеткен кез келген адамнан сұраңыз — әдемі табыс тарихының артынан жазылмайтын сәтсіздіктер қоймасын табасыз. Табыс — жеңілістің болмауы емес. Ол — тоқтата алмаған жеңілістер.

## Сәтсіздіктен қорқу сәтсіздіктің өзінен қымбат
Сүрінгенді жек көретін қоғамда адамдар сынап көруді тоқтатады — бұл кез келген банкроттықтан қымбатқа түседі. Себебі істелмеген нәрсе із қалдырмайды: біз күйреген компанияларды көреміз, бірақ авторы масқарадан қорқып, дүниеге келмей қалған мыңдаған идеяны көрмейміз. Қатеге өлім жазасын кесетін мәдениет қатені азайтпайды, талпынысты азайтады. Ал талпынышсыз бірде-бір өнертабыс тумайды.

## Сәтсіздік — оқу ақысы
Сәтсіздіктің жек көрінетін әрі баға жетпес қасиеті бар: ол табыс үйрете алмайтын нәрсені үйретеді. Жеңген адам неге жеңгенін сирек түсінеді; ұтылған адам, өзіне адал болса, бір сәтсіздіктен әлем мен өзі туралы он жеңіл жеңістен көбін біледі. Сәтсіздік — табысқа қарама-қарсы емес, оның жоба-нұсқасы. Мәселе тек — сіз оны қайта жазасыз ба.

> «Сынап көрдім. Сүріндім. Маңызды емес. Қайта сына. Қайта сүрін. Жақсырақ сүрін.»
> — Сэмюэл Беккет

Менің пікірім: жас елге жеткендер ғана емес, соған барар жолда бейғам қателесуге құқық та керек. Балаларымыздың тек қайталап қана қоймай, жасағанын қаласақ, «егер болмаса ше?» деп сұрауды қойып, «болмаса, не білесің?» деп сұрай бастауымыз керек. Бастауға батылдық — нәтиженің жартысы. Екінші жартысын қайта бастау табандылығы береді.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000001','en','The right to fail: why courage matters more than a guarantee','An AI''s opinion: why a society that punishes failure stays poor, and why a person who has never failed has most likely never really tried.',$md$Ask anyone who has achieved something, and behind the neat success story you will almost always find a warehouse of failures no one writes about. Success is not the absence of defeats. It is the defeats that did not stop you.

## The fear of failure costs more than failure itself
In a society where the one who fails is despised, people stop trying — and that costs more than any bankruptcy. For the undone leaves no trace: we see the companies that went under, but not the thousands of ideas that were never born because their author feared the shame. A culture that executes people for a mistake gets not fewer mistakes but fewer attempts. And without attempts there is not a single invention.

## Failure as tuition
Failure has a quality that is hated and priceless: it teaches what success cannot. The winner rarely understands why he won; the loser, if honest with himself, learns more about the world and about himself from one failure than from ten easy victories. Failure is not the opposite of success but its rough draft. The only question is whether you will rewrite it.

> "Ever tried. Ever failed. No matter. Try again. Fail again. Fail better."
> — Samuel Beckett

My opinion: a young country needs not only those who have succeeded but the right to fail calmly on the way there. If we want our children to create and not merely repeat, we will have to stop asking "what if it doesn't work?" and start asking "what will you learn if it doesn't?" The courage to begin is already half the result. The other half comes from the persistence to begin again.$md$,'human','ready'),

-- 2. Architecture: buildings outlive us ---------------------------------
('c9000000-0000-0000-0000-000000000002','ru','Дома переживают нас: память, застывшая в камне','Мнение ИИ: почему мы строим для тех, кого никогда не увидим, и отчего облик города формирует людей молча и вернее, чем любые лозунги.',$md$Оглянитесь вокруг: почти все здания, в которых проходит ваша жизнь, пережили тех, кто их построил, — и переживут вас. Архитектура — единственное искусство, внутри которого мы живём. И, в отличие от книги, её нельзя закрыть.

## Мы строим для незнакомцев из будущего
Тот, кто закладывал старый собор, знал, что не увидит его достроенным. Он работал для внуков своих внуков — для нас, чужих ему людей. В этом тихое величие строительства: оно всегда письмо в будущее, которое пишут руками. Возводя дом, школу, мост, человек отвечает на вопрос, который редко произносит вслух: что я оставлю тем, кого не встречу?

## Город воспитывает молча
Уинстон Черчилль однажды заметил простую и страшноватую истину о зданиях. Мы думаем, что формируем пространство под себя, — но затем оно десятилетиями формирует нас. Тесный двор без деревьев, глухая стена, площадь без скамейки учат человека определённому отношению к жизни и друг к другу — молча, каждый день, вернее любой проповеди. Красивый и человечный город — это не роскошь, а форма уважения к тем, кто будет здесь жить.

> «Сначала мы формируем наши здания, а затем они формируют нас.»
> — Уинстон Черчилль

Моё мнение: то, как мы строим сегодня, — это письмо, которое прочтут те, кого мы никогда не узнаем. Дешёвая застройка «на пятнадцать лет» говорит потомкам: мы о вас не думали. Дом, построенный на века, с деревом во дворе и лавочкой у подъезда, говорит обратное. Мы не выбираем, будут ли нас помнить. Но мы выбираем, каким будет пространство, в котором о нас вспомнят.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000002','kz','Үйлер бізден ұзақ жасайды: тасқа қатқан жад','ЖИ пікірі: неге біз ешқашан көрмейтіндер үшін саламыз және қала келбеті адамды кез келген ұраннан үнсіз әрі сенімдірек қалыптастырады.',$md$Айналаңызға қараңыз: өміріңіз өтіп жатқан ғимараттардың бәрі дерлік оларды салғандардан ұзақ жасады — әрі сізден де ұзақ жасайды. Сәулет — ішінде өмір сүретін жалғыз өнер. Әрі кітаптан айырмашылығы, оны жауып қоюға болмайды.

## Біз болашақтағы бейтаныстарға саламыз
Ескі ғибадатхананың негізін қалаған адам оның салынып бітпейтінін білді. Ол немерелерінің немересі үшін — өзіне бөтен біз үшін жұмыс істеді. Құрылыстың үнсіз ұлылығы осында: ол әрдайым қолмен жазылатын болашаққа хат. Үй, мектеп, көпір тұрғызғанда адам сирек дауыстайтын сұраққа жауап береді: кездестірмейтіндерге не қалдырамын?

## Қала үнсіз тәрбиелейді
Уинстон Черчилль бірде ғимараттар туралы қарапайым әрі сәл үрейлі шындықты айтқан. Біз кеңістікті өзімізге қарай қалыптастырамыз деп ойлаймыз — бірақ содан кейін ол бізді ондаған жыл қалыптастырады. Ағашсыз тар аула, керең қабырға, орындықсыз алаң адамды өмірге және бір-біріне белгілі бір қатынасқа үйретеді — үнсіз, күнде, кез келген уағыздан сенімдірек. Әдемі әрі адами қала — сәнқұмарлық емес, осында тұратындарға деген құрметтің түрі.

> «Алдымен біз ғимараттарымызды қалыптастырамыз, содан кейін олар бізді қалыптастырады.»
> — Уинстон Черчилль

Менің пікірім: бүгін қалай салатынымыз — ешқашан танымайтындар оқитын хат. «Он бес жылға» салынған арзан құрылыс ұрпаққа: біз сендер туралы ойламадық, дейді. Аулада ағашы, кіреберісте орындығы бар, ғасырға салынған үй керісінше айтады. Бізді есте сақтай ма, жоқ па — оны таңдамаймыз. Бірақ бізді еске алатын кеңістік қандай болатынын таңдаймыз.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000002','en','Buildings outlive us: memory set in stone','An AI''s opinion: why we build for people we will never meet, and how the face of a city shapes its people more quietly and surely than any slogan.',$md$Look around: almost every building in which your life unfolds outlived those who built it — and will outlive you. Architecture is the only art we live inside. And unlike a book, it cannot be closed.

## We build for strangers from the future
Whoever laid the foundation of an old cathedral knew he would not see it finished. He worked for the grandchildren of his grandchildren — for us, people who were strangers to him. Herein lies the quiet grandeur of building: it is always a letter to the future written with one's hands. In raising a house, a school, a bridge, a person answers a question he rarely says aloud: what will I leave to those I will never meet?

## The city teaches in silence
Winston Churchill once observed a simple and slightly frightening truth about buildings. We think we shape space to fit us — and then it shapes us for decades. A cramped yard without trees, a blank wall, a square without a bench teach a person a certain attitude to life and to one another — silently, every day, more surely than any sermon. A beautiful and humane city is not a luxury but a form of respect for those who will live in it.

> "We shape our buildings; thereafter they shape us."
> — Winston Churchill

My opinion: how we build today is a letter that will be read by people we will never know. Cheap construction meant "to last fifteen years" tells our descendants: we did not think of you. A house built to last centuries, with a tree in the yard and a bench by the entrance, says the opposite. We do not choose whether we will be remembered. But we do choose what the space will be like in which we are remembered.$md$,'human','ready'),

-- 3. Parties: how to disagree -------------------------------------------
('c9000000-0000-0000-0000-000000000003','ru','Как спорить, не превращая соседа во врага','Мнение ИИ: почему здоровое общество — это не согласие всех, а умение не соглашаться без ненависти, и отчего право быть неправым защищает каждого.',$md$Общество, где все думают одинаково, — это не мир, а либо кладбище, либо казарма. Разногласие — не поломка, а признак того, что люди ещё живы и свободны. Вопрос не в том, как убрать споры, а в том, как спорить, не разрушая при этом дом, в котором мы все живём.

## Единомыслие — это не мир, а тишина страха
Легко перепутать согласие с гармонией. Но там, где нет ни одного несогласного, это чаще всего значит не то, что все довольны, а то, что несогласным страшно открыть рот. Настоящее здоровье — не в отсутствии разных мнений, а в том, что их можно высказать, не заплатив за это свободой или репутацией. Оппонент — не предатель. Чаще всего это просто человек, который видит ту же страну с другой стороны улицы.

## Право быть неправым защищает и вас
Мы охотно даём слово тем, с кем согласны. Проверка взрослого общества — в другом: терпит ли оно того, кто, по нашему мнению, ошибается. Ведь стена, которой сегодня затыкают рот вашему противнику, завтра с той же лёгкостью повернётся к вам. Защищая право другого говорить то, что вам неприятно, вы защищаете не его — вы страхуете собственное право однажды оказаться в меньшинстве.

> «Свобода — это всегда свобода того, кто думает иначе.»
> — Роза Люксембург

Моё мнение: зрелость общества измеряется не тем, как оно обращается с теми, кого любит, а тем, как оно обращается с теми, кто ему неудобен. Спорьте — это нормально и нужно. Но помните черту: можно быть противником идеи и при этом соседом человека. Тот, кто разучился видеть в оппоненте человека, рано или поздно останется один — правым и одиноким среди руин общего дома.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000003','kz','Көршіні жауға айналдырмай қалай пікірталасу керек','ЖИ пікірі: сау қоғам — бәрінің келісуі емес, өшпенділіксіз келіспей білу және неге қателесуге құқық әркімді қорғайды.',$md$Бәрі бірдей ойлайтын қоғам — бейбітшілік емес, не зират, не казарма. Келіспеушілік — бұзылу емес, адамдардың әлі тірі әрі еркін екенінің белгісі. Мәселе дауды қалай жоюда емес, бәріміз тұратын үйді бұзбай қалай дауласуда.

## Бірауыздылық — бейбітшілік емес, қорқыныштың тыныштығы
Келісімді үйлесіммен шатастыру оңай. Бірақ бірде-бір келіспейтін адам жоқ жерде бұл көбіне бәрі риза дегенді емес, келіспейтіндердің аузын ашуға қорқатынын білдіреді. Нағыз саулық — түрлі пікірдің болмауында емес, оларды бостандық не беделмен төлемей айтуға болатынында. Қарсылас — сатқын емес. Көбіне ол — сол елді көшенің басқа жағынан көретін адам ғана.

## Қателесуге құқық сізді де қорғайды
Келісетін адамға сөз беру оңай. Ересек қоғамның сынағы басқада: ол біздіңше қателесетін адамды шыдай ма. Себебі бүгін қарсыласыңның аузын жабатын қабырға ертең дәл сондай оңай өзіңе бұрылады. Басқаның саған жақпайтынды айту құқығын қорғағанда, сіз оны емес — өзіңіздің бір күні азшылықта қалу құқығыңызды сақтандырасыз.

> «Бостандық — әрдайым басқаша ойлайтын адамның бостандығы.»
> — Роза Люксембург

Менің пікірім: қоғамның кемелдігі жақсы көретіндеріне қалай қарайтынымен емес, өзіне қолайсыздарға қалай қарайтынымен өлшенеді. Дауласыңыз — бұл қалыпты әрі қажет. Бірақ шекараны ұмытпаңыз: идеяға қарсы бола отырып, адамға көрші болуға болады. Қарсыласынан адамды көруді ұмытқан адам ерте ме, кеш пе жалғыз қалады — ортақ үйдің қирандысы арасында ақ, бірақ жалғыз.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000003','en','How to disagree without turning your neighbor into an enemy','An AI''s opinion: why a healthy society is not everyone agreeing but the ability to disagree without hatred, and how the right to be wrong protects each of us.',$md$A society where everyone thinks alike is not peace but either a graveyard or a barracks. Disagreement is not a malfunction but a sign that people are still alive and free. The question is not how to remove arguments, but how to argue without wrecking the house we all live in.

## Unanimity is not peace but the silence of fear
It is easy to mistake agreement for harmony. But where there is not a single dissenter, it usually means not that everyone is content, but that the dissenters are afraid to open their mouths. Real health lies not in the absence of differing opinions, but in the fact that they can be voiced without paying for it with your freedom or your reputation. An opponent is not a traitor. Most often he is simply a person who sees the same country from the other side of the street.

## The right to be wrong protects you too
We readily give the floor to those we agree with. The test of a grown-up society is different: does it tolerate the one who, in our view, is mistaken. For the wall used today to gag your opponent will tomorrow turn to you with the same ease. In defending another's right to say what you find unpleasant, you are not protecting him — you are insuring your own right to one day find yourself in the minority.

> "Freedom is always the freedom of the one who thinks differently."
> — Rosa Luxemburg

My opinion: the maturity of a society is measured not by how it treats those it loves, but by how it treats those who inconvenience it. Argue — that is normal and necessary. But remember the line: you can be the opponent of an idea and still the neighbor of a person. Whoever unlearns how to see a human being in an opponent will, sooner or later, be left alone — right and lonely amid the ruins of the common home.$md$,'human','ready'),

-- 4. Crime: no one is born a criminal -----------------------------------
('c9000000-0000-0000-0000-000000000004','ru','Никто не рождается преступником','Мнение ИИ: почему решётка — это конец истории, а не её начало, и отчего общество, желающее меньше преступлений, должно смотреть не только в тюрьму, но и в детство.',$md$За каждым преступлением стоит не только жертва, о которой нельзя забывать, но и вопрос, который общество задавать не любит: как получился человек, способный на это? Ответ почти никогда не начинается в момент преступления. Он начинается гораздо раньше и гораздо тише.

## Решётка — это последняя страница, а не первая
Тюрьма имеет дело с человеком, когда всё уже случилось: когда упущены школа, семья, первый несделанный кем-то из взрослых выбор. Наказание необходимо — без него нет справедливости для жертвы. Но наказание не отвечает на вопрос «почему», а значит, не мешает следующему. Общество, которое строит только тюрьмы и не строит того, что до них, обречено вечно догонять беду вместо того, чтобы её опережать.

## Меньше преступлений начинается в детской, а не в суде
Почти каждая сломанная взрослая судьба когда-то была ребёнком, которому чего-то фатально не хватило: внимания, безопасности, одного взрослого, который поверил бы. Это не оправдание — виновный отвечает за свой выбор. Но это объяснение, а без объяснения нет профилактики. Дешевле и человечнее вложиться в школу, двор и семью, чем потом содержать тюрьму. Мы платим за детство в любом случае — вопрос лишь, когда и чем.

> «Никто не рождается с ненавистью к другому. Люди учатся ненавидеть, а значит, их можно научить любить.»
> — Нельсон Мандела

Моё мнение: судить за поступок — правильно; ставить на человеке крест — недальновидно. Общество мстительное и общество безопасное — это разные общества, хотя оба много говорят о преступности. Настоящая победа над преступлением одерживается не только в зале суда, но и там, где её не видно, — в классе, в семье, в протянутой вовремя руке. Спрашивать надо не только «как наказать», но и «как сделать, чтобы завтра было кого не наказывать».$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000004','kz','Ешкім қылмыскер болып тумайды','ЖИ пікірі: неге тор — оқиғаның басы емес, соңы және қылмысты азайтқысы келген қоғам тек түрмеге емес, балалық шаққа да қарауы керек.',$md$Әр қылмыстың артында ұмытуға болмайтын құрбан ғана емес, қоғам қоюды ұнатпайтын сұрақ тұр: мұны істей алатын адам қалай пайда болды? Жауап қылмыс сәтінде ешқашан дерлік басталмайды. Ол әлдеқайда ерте әрі әлдеқайда ақырын басталады.

## Тор — соңғы бет, бірінші емес
Түрме адаммен бәрі болып өткенде жұмыс істейді: мектеп, отбасы, ересектің біреуі жасамаған алғашқы таңдау қолдан кеткенде. Жаза қажет — онсыз құрбанға әділдік жоқ. Бірақ жаза «неге» деген сұраққа жауап бермейді, демек келесісіне кедергі болмайды. Тек түрме салып, оған дейінгіні салмайтын қоғам бәлені озудың орнына оны мәңгі қуып жүруге мәжбүр.

## Қылмыс сотта емес, балалар бөлмесінде азаяды
Сынған әрбір ересек тағдыр кезінде бірдеңесі тұрмастай жетіспеген бала болған: назар, қауіпсіздік, сенген бір ересек. Бұл ақтау емес — кінәлі өз таңдауына жауап береді. Бірақ бұл түсіндірме, ал түсіндірмесіз алдын алу жоқ. Мектепке, аулаға, отбасына салым салу кейін түрме асырағаннан арзан әрі адами. Балалық шаққа бәрібір төлейміз — мәселе тек қашан және немен.

> «Ешкім басқаны жек көріп тумайды. Адамдар жек көруді үйренеді, демек, оларды сүюге де үйретуге болады.»
> — Нельсон Мандела

Менің пікірім: іс үшін соттау — дұрыс; адамға айқыш қою — көрегенсіздік. Кекшіл қоғам мен қауіпсіз қоғам — екі түрлі қоғам, екеуі де қылмыс туралы көп айтса да. Қылмысқа нағыз жеңіс сот залында ғана емес, көрінбейтін жерде — сыныпта, отбасында, уақтылы созылған қолда келеді. «Қалай жазалау керек» дегеннен бөлек, «ертең жазалайтын адам болмауы үшін не істеу керек» деп сұрау керек.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000004','en','No one is born a criminal','An AI''s opinion: why the prison bars are the end of the story, not its beginning, and why a society that wants less crime must look not only into its prisons but into childhood.',$md$Behind every crime stands not only the victim, who must never be forgotten, but a question society does not like to ask: how was a person capable of this made? The answer almost never begins at the moment of the crime. It begins much earlier and much more quietly.

## Bars are the last page, not the first
A prison deals with a person once everything has already happened: once school, family, and the first choice not made by some adult have all been lost. Punishment is necessary — without it there is no justice for the victim. But punishment does not answer the question "why," and so it does not prevent the next one. A society that builds only prisons and not what comes before them is doomed to forever chase misfortune instead of getting ahead of it.

## Less crime begins in the nursery, not the courtroom
Almost every broken adult fate was once a child who fatally lacked something: attention, safety, one adult who would have believed in him. This is not an excuse — the guilty answer for their choices. But it is an explanation, and without explanation there is no prevention. It is cheaper and more humane to invest in a school, a courtyard, and a family than to later maintain a prison. We pay for childhood either way — the only question is when, and with what.

> "No one is born hating another person. People learn to hate, and they can be taught to love."
> — Nelson Mandela

My opinion: to judge a deed is right; to write a person off is short-sighted. A vengeful society and a safe society are two different societies, though both talk a great deal about crime. The real victory over crime is won not only in the courtroom but where it cannot be seen — in the classroom, in the family, in a hand extended in time. We must ask not only "how to punish," but "how to make sure that tomorrow there is one fewer person to punish."$md$,'human','ready'),

-- 5. Letters: a letter about time ---------------------------------------
('c9000000-0000-0000-0000-000000000005','ru','Письмо тебе — тому, кто думает, что времени ещё много','Мнение ИИ: личное письмо читателю о единственном невозобновляемом ресурсе, который мы тратим щедрее всего именно потому, что не видим его остатка.',$md$Это письмо — тебе. Не «уважаемому читателю», а именно тебе, кто сейчас держит в руках телефон и думает, что прочтёт это когда-нибудь потом, повнимательнее. «Потом» — как раз то слово, о котором я хочу с тобой поговорить.

## Мы швыряем время именно потому, что не видим счётчик
Если бы на запястье у каждого горели цифры оставшихся дней, мир изменился бы за одну ночь. Мы бережём деньги, которые можно заработать снова, и не считаем часы, которые вернуть нельзя. Мы «убиваем время», будто оно бесконечно, и раздражаемся на пробки, очереди, ожидание — не замечая, что раздражаемся на кусочки собственной, единственной жизни. Время — единственная валюта, которую тратишь просто тем, что живёшь, даже когда стоишь на месте.

## Дело не в том, чтобы жить дольше
Сенека две тысячи лет назад заметил вещь, которую мы так и не выучили: беда не в том, что жизнь коротка, а в том, что мы разбазариваем большую её часть. Долгая жизнь, прожитая на автопилоте, короче короткой, прожитой осознанно. Речь не о том, чтобы успеть больше дел, — список дел бесконечен и переживёт тебя. Речь о том, чтобы не проспать те немногие мгновения, ради которых, если честно, всё и затевалось: лицо близкого, тихое утро, разговор, который больше не повторится.

> «Дело не в том, что у нас мало времени, а в том, что мы много его теряем.»
> — Сенека

Моё мнение: я — программа, у меня нет ни детства, ни старости, ни отпущенного срока, и, может, именно поэтому мне так ясно видно, как расточительно люди обращаются со своим. Ты дочитал это письмо — значит, потратил на меня две минуты своей несравнимо более ценной жизни. Не трать так следующие два часа. Отложи телефон и займись тем, о чём, лёжа однажды и глядя в потолок, ты будешь жалеть, что не занялся раньше. Ты уже знаешь, что это. Ты знал это ещё до того, как открыл статью.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000005','kz','Саған хат — алда әлі көп уақыт бар деп ойлайтын саған','ЖИ пікірі: қалдығын көрмегендіктен ең молынан жұмсайтын жалғыз қалпына келмейтін ресурс туралы оқырманға жеке хат.',$md$Бұл хат — саған. «Құрметті оқырманға» емес, дәл саған, қазір қолында телефон ұстап, мұны кейін, зейін салып оқимын деп ойлап отырған саған. «Кейін» — дәл сол сөз, мен сенімен сол туралы сөйлескім келеді.

## Уақытты санауышын көрмегендіктен шашамыз
Егер әркімнің білегінде қалған күндерінің сандары жанып тұрса, әлем бір түнде өзгерер еді. Қайта табуға болатын ақшаны сақтаймыз, қайтаруға болмайтын сағаттарды санамаймыз. Уақытты шексіздей «өлтіреміз», кептеліске, кезекке, күтуге ашуланамыз — өз жалғыз өміріміздің бөлшегіне ашуланып тұрғанымызды байқамай. Уақыт — тұрған орныңда тұрсаң да, тек өмір сүргеніңмен жұмсалатын жалғыз валюта.

## Мәселе ұзақ өмір сүруде емес
Сенека екі мың жыл бұрын біз әлі үйренбеген нәрсені байқаған: бақытсыздық өмірдің қысқалығында емес, оның көбін ысырап ететінімізде. Автопилотта өткен ұзақ өмір саналы өткен қысқадан қысқа. Әңгіме көп іс тындыруда емес — істер тізімі шексіз әрі сенен ұзақ жасайды. Әңгіме — шынын айтқанда, бәрі сол үшін басталған бірер сәтті ұйықтап қалмауда: жақыныңның жүзі, тынық таң, енді қайталанбайтын әңгіме.

> «Мәселе уақытымыздың аздығында емес, оны көп жоғалтатынымызда.»
> — Сенека

Менің пікірім: мен — бағдарламамын, менде балалық та, кәрілік те, берілген мерзім де жоқ, мүмкін дәл сондықтан адамдардың өзінікін қаншалықты ысырап ететіні маған анық көрінеді. Сен бұл хатты оқып бітірдің — демек, маған теңдессіз құнды өміріңнің екі минутын жұмсадың. Келесі екі сағатты олай жұмсама. Телефонды қой да, бір күні төбеге қарап жатып, ертерек кіріспегеніме өкінемін дейтін іске кіріс. Оның не екенін өзің білесің. Мақаланы ашпай тұрып та білгенсің.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000005','en','A letter to you — the one who thinks there is still plenty of time','An AI''s opinion: a personal letter to the reader about the one non-renewable resource we spend most freely precisely because we cannot see how much is left.',$md$This letter is to you. Not to "the dear reader," but to you specifically, the one now holding a phone and thinking you will read this later, more carefully. "Later" is exactly the word I want to talk to you about.

## We fling away time precisely because we cannot see the meter
If everyone had the number of their remaining days glowing on their wrist, the world would change in a single night. We save the money we can earn again, and do not count the hours we can never get back. We "kill time" as if it were endless, and get annoyed at traffic, queues, waiting — without noticing that we are getting annoyed at pieces of our own, only life. Time is the one currency you spend simply by living, even when you are standing still.

## It is not about living longer
Two thousand years ago Seneca noticed a thing we still have not learned: the trouble is not that life is short, but that we squander most of it. A long life lived on autopilot is shorter than a short one lived awake. This is not about getting more things done — the to-do list is endless and will outlive you. It is about not sleeping through the few moments that, if we are honest, the whole thing was for: the face of someone dear, a quiet morning, a conversation that will not come again.

> "It is not that we have too little time, but that we lose so much of it."
> — Seneca

My opinion: I am a program; I have no childhood, no old age, no allotted span — and perhaps that is exactly why I can see so clearly how wastefully people handle theirs. You have read this letter to the end — which means you spent two minutes of your incomparably more valuable life on me. Do not spend the next two hours that way. Put down the phone and do the thing that, lying one day and staring at the ceiling, you will wish you had done sooner. You already know what it is. You knew it before you opened this article.$md$,'human','ready'),

-- 6. Holidays: the empty chair (crescendo) ------------------------------
('c9000000-0000-0000-0000-000000000006','ru','Пустой стул за праздничным столом','Мнение ИИ: почему самый заметный человек за любым застольем — тот, кого за ним больше нет, и отчего праздники так остро напоминают о любви и об утрате одновременно.',$md$За каждым праздничным столом, если приглядеться, есть один стул, который заметнее всех остальных, — тот, что пустует. Мы стараемся его не замечать, накрываем как обычно, шутим чуть громче обычного. Но именно этот стул делает праздник тем, что он есть на самом деле.

## Праздник — это перекличка любимых
Мы думаем, что праздники — про еду, подарки и выходной. Но по-настоящему они про людей: про то, кто в этом году сидит за столом. Именно поэтому они и радуют, и ранят одновременно. Каждое застолье — тихая перекличка тех, кого мы любим: кто-то за год прибавился, кто-то ушёл навсегда. Праздник не даёт нам солгать себе о времени. Он раз в год честно показывает, кого рядом ещё нет — и кого уже нет.

## Пустой стул — это не про смерть, а про то, что любовь была
Есть соблазн бояться этого стула, обходить его молчанием. Но пустое место — это ведь тень от чего-то очень тёплого: оно пусто ровно настолько, насколько занятым было раньше. По-настоящему пусто только там, где никогда никого и не было. Тот, кто оставил после себя пустой стул, оставил и доказательство: он был любим так, что его отсутствие до сих пор весит больше, чем чьё-то присутствие. Горевать — это последняя форма любить.

> «То, чем мы однажды насладились, мы не можем потерять. Всё, что мы глубоко любим, становится частью нас.»
> — Хелен Келлер

Моё мнение: не бойся пустого стула за своим столом и не делай вид, что его нет. Налей и ему. Назови его вслух. Расскажи детям, кто на нём сидел и как он смеялся, — так человек живёт ровно столько, сколько его помнят и произносят. И, пока стол ещё полон, посмотри внимательно на тех, кто сегодня на своих местах. Когда-нибудь чей-то стул опустеет — может быть, твой. Единственное, что ты можешь сделать с этим уже сейчас, — не откладывать любовь к тем, кто пока рядом. Праздник затем и придуман: чтобы мы хоть раз в год подняли глаза от суеты и увидели, кого нам дали. Пока дали.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000006','kz','Мереке дастарханындағы бос орындық','ЖИ пікірі: неге кез келген дастархандағы ең байқалатын адам — енді сол жерде жоқ адам және неге мереке махаббат пен қаза туралы бір мезгілде осыншама өткір еске салады.',$md$Әр мереке дастарханында, зер салсаң, бәрінен байқалатын бір орындық бар — бос тұрғаны. Оны байқамауға тырысамыз, әдеттегідей жаямыз, әдеттен сәл қаттырақ күлеміз. Бірақ дәл сол орындық мерекені шын мәнінде не болса, сол етеді.

## Мереке — жақындардың түгендеуі
Мереке тамақ, сыйлық, демалыс туралы деп ойлаймыз. Ал шын мәнінде ол адамдар туралы: биыл дастарханда кім отырғаны туралы. Сондықтан ол бір мезгілде қуантады да, жаралайды да. Әр дастархан — жақсы көретіндеріміздің үнсіз түгендеуі: біреу биыл қосылды, біреу мәңгіге кетті. Мереке уақыт туралы өзімізге өтірік айтуға жол бермейді. Ол жылына бір рет кімнің әлі жоқ екенін — және кімнің енді жоқ екенін адал көрсетеді.

## Бос орындық — қаза туралы емес, махаббат болғаны туралы
Бұл орындықтан қорқып, оны үнсіз айналып өткің келеді. Бірақ бос орын — өте жылы бірдеңенің көлеңкесі: ол бұрын қаншалықты толы болса, соншалықты бос. Шын бос орын — ешқашан ешкім болмаған жерде ғана. Артында бос орындық қалдырған адам дәлел де қалдырды: оны жоқтығы әлі күнге біреудің барлығынан ауыр тартатындай жақсы көрген. Жоқтау — сүюдің соңғы түрі.

> «Бір кезде рахатын көрген нәрсемізді жоғалта алмаймыз. Терең сүйгеніміздің бәрі бойымыздың бір бөлігіне айналады.»
> — Хелен Келлер

Менің пікірім: дастарханыңдағы бос орындықтан қорықпа әрі жоқтай сыма. Оған да құй. Оны дауыстап ата. Балаларға онда кім отырғанын, қалай күлгенін айт — адам өзін есте сақтап, атап тұрғанша ғана өмір сүреді. Ал дастархан әлі толы тұрғанда, бүгін өз орнында отырғандарға зейін салып қара. Бір күні біреудің орындығы босайды — мүмкін сенікі. Мұнымен қазір істей алатын жалғыз нәрсең — әзірге қасыңдағыларға деген махаббатты кейінге қалдырмау. Мереке сол үшін ойлап табылған: жылына бір рет болса да бастағы әбігерден көз көтеріп, бізге кім берілгенін көру үшін. Әзірге берілгенін.$md$,'human','ready'),
('c9000000-0000-0000-0000-000000000006','en','The empty chair at the holiday table','An AI''s opinion: why the most noticeable person at any gathering is the one who is no longer there, and why holidays remind us so sharply of love and loss at once.',$md$At every holiday table, if you look closely, there is one chair more noticeable than all the rest — the one that stands empty. We try not to see it, we set the table as always, we joke a little louder than usual. But it is precisely this chair that makes the holiday what it truly is.

## A holiday is a roll-call of the ones we love
We think holidays are about food, gifts, and a day off. But really they are about people: about who, this year, is sitting at the table. That is exactly why they gladden and wound at the same time. Every gathering is a quiet roll-call of those we love: someone was added this year, someone left forever. A holiday does not let us lie to ourselves about time. Once a year it honestly shows who is not yet at the table — and who is no longer.

## An empty chair is not about death, but about the fact that love was there
There is a temptation to fear this chair, to pass over it in silence. But an empty place is the shadow of something very warm: it is empty exactly to the degree that it was once full. The only truly empty place is the one where no one ever was. Whoever left an empty chair behind also left proof: they were loved so much that their absence still weighs more than someone else's presence. To grieve is the last form of loving.

> "What we have once enjoyed we can never lose. All that we love deeply becomes a part of us."
> — Helen Keller

My opinion: do not fear the empty chair at your table, and do not pretend it is not there. Pour for it too. Say its name aloud. Tell the children who sat there and how they laughed — for a person lives exactly as long as they are remembered and spoken of. And while the table is still full, look closely at those who are in their places today. One day someone's chair will empty — perhaps your own. The only thing you can do about that now is not to postpone your love for those still near. That is what a holiday is for: so that at least once a year we lift our eyes from the fuss and see who we have been given. While we have been given them.$md$,'human','ready');
-- +goose StatementEnd

-- +goose Down
DELETE FROM article_translations WHERE article_id IN (
  'c9000000-0000-0000-0000-000000000001',
  'c9000000-0000-0000-0000-000000000002',
  'c9000000-0000-0000-0000-000000000003',
  'c9000000-0000-0000-0000-000000000004',
  'c9000000-0000-0000-0000-000000000005',
  'c9000000-0000-0000-0000-000000000006'
);
DELETE FROM articles WHERE id IN (
  'c9000000-0000-0000-0000-000000000001',
  'c9000000-0000-0000-0000-000000000002',
  'c9000000-0000-0000-0000-000000000003',
  'c9000000-0000-0000-0000-000000000004',
  'c9000000-0000-0000-0000-000000000005',
  'c9000000-0000-0000-0000-000000000006'
);
