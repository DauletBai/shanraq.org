"""Authored paper-only block 23–31: bonding, formulas and evidence.

Each translation supplies a title, summary and eight teaching sections.
Diagrams are symbolic; electron bookkeeping is not an observed reaction path.
"""

LESSONS = [
{
 'stem': '23-why-atoms-bond',
 'ru': (
  'Почему атомы образуют химические связи?',
  'Свяжем притяжение и отталкивание с устойчивостью; отделим объяснение от выражения «атом хочет».',
  'В дневнике есть H и H₂. Обе записи относятся к водороду, но обозначают разные частицы. Почему два нейтральных атома могут удерживаться вместе, если заряд каждого равен нулю? Запишите предположение. Работайте с рисунками: водород получать не нужно. Ни цвет шариков, ни линия между ними не являются наблюдением настоящих электронов.',
  '**Химическая связь** — взаимодействие, удерживающее атомы или ионы в составе вещества. В атомах есть положительные ядра и отрицательные электроны. Между зарядами действуют притяжение и отталкивание. **Энергия системы** помогает сравнить состояния: для H₂ при подходящем расстоянии она ниже, чем у двух далёких атомов H. Для разрыва этой связи нужно затратить энергию.',
  'Предскажите: будут ли ядра всё сильнее сближаться вплоть до слияния? Означает ли нейтральность отсутствие заряженных частиц внутри атома? Отметьте, где вам известен суммарный заряд, а где нужно рассмотреть устройство частицы. Проверьте ответы после работы с таблицей, а не заменяйте объяснение словом «устойчивость».',
  'Разберите условные данные для двух H: далеко — энергия 0; подходящее расстояние — −4; слишком близко — +6. Это учебные числа в условных единицах, не измеренные энергии водорода. Нарисуйте три положения и подпишите притяжение к электронам и отталкивание ядер. Сравнивайте только состояния одной и той же пары атомов.',
  'Из трёх заданных состояний самое низкое значение −4. При чрезмерном сближении отталкивание становится существенным: энергия возрастает до +6. Поэтому связь не означает слияния ядер. Для перехода от −4 к двум далёким атомам с энергией 0 в этой модели нужно сообщить 4 условные единицы энергии.',
  '**Ошибка:** «атомы хотят восемь электронов». У атомов нет желаний; правило заполненного слоя — полезная модель для некоторых простых случаев, а не универсальная причина. Для H заполнение связано с двумя электронами, не восемью. Образование связи выделяет энергию, но реакция целиком может требовать подвода энергии из-за других разрываемых связей.',
  'Новая модель пары частиц: далеко 0, в связанном состоянии −7, слишком близко +9. Выберите состояние с наименьшей энергией и рассчитайте затраты на разрыв связи до далёких частиц. Можно ли из этих чисел узнать цвет вещества или доказать, что любая пара элементов образует такое же соединение?',
  '**Проверка:** −7; для разрыва требуется 7 условных единиц. Цвет и состав других веществ неизвестны. Воспроизведите цепочку «заряды → взаимодействия → подходящее расстояние → более низкая энергия». Назовите ограничение: это качественная модель заданной пары, без реальных масштабов и без расчёта всех возможных веществ.'
 ),
 'kz': (
  'Атомдар неліктен химиялық байланыс түзеді?',
  'Тартылыс пен тебілісті тұрақтылықпен байланыстырамыз; түсіндіруді «атом қалайды» деген сөзден ажыратамыз.',
  'Күнделікте H және H₂ жазылған. Екеуі де сутекке қатысты, бірақ әртүрлі бөлшектерді білдіреді. Әр атомның жалпы заряды нөл болса, екі бейтарап атом қалай бірге ұсталады? Болжам жазыңыз. Тек суретпен жұмыс істеңіз: сутек алудың қажеті жоқ. Шардың түсі де, олардың арасындағы сызық та нақты электронды бақылау емес.',
  '**Химиялық байланыс** — зат құрамындағы атомдарды немесе иондарды бірге ұстайтын әрекеттесу. Атомда оң ядро мен теріс электрондар бар. Зарядтар арасында тартылыс пен тебіліс болады. **Жүйе энергиясы** күйлерді салыстыруға көмектеседі: H₂ үшін қолайлы аралықта энергия бір-бірінен алыс екі H атомының энергиясынан төмен. Осы байланысты үзуге энергия жұмсалады.',
  'Болжаңыз: ядролар қосылып кеткенше жақындай бере ме? Бейтарап атомның ішінде зарядталған бөлшектер жоқ па? Қай жерде жалпы заряд белгілі, қай жерде бөлшектің құрылысын ескеру керек екенін белгілеңіз. Жауапты кестеден кейін тексеріңіз; түсіндіру орнына тек «тұрақтылық» деген сөзді жазбаңыз.',
  'Екі H үшін шартты деректі талдаңыз: алыс — энергия 0; қолайлы аралықта — −4; тым жақын — +6. Бұлар шартты бірліктегі оқу сандары, сутектің өлшенген энергиялары емес. Үш орналасуды сызып, электрондарға тартылыс пен ядролардың тебілісін белгілеңіз. Бірдей екі атомның күйлерін ғана салыстырыңыз.',
  'Берілген үш күйдің ең төмен энергиясы −4. Тым жақындағанда тебіліс елеулі болып, энергия +6-ға өседі. Сондықтан байланыс ядролардың қосылуын білдірмейді. Бұл модельде −4 күйінен энергиясы 0 болатын алыс атомдарға өту үшін 4 шартты энергия бірлігін беру керек.',
  '**Қате:** «атом сегіз электрон қалайды». Атомның қалауы жоқ; толған қабат ережесі кейбір қарапайым жағдайларға жарайтын модель, әмбебап себеп емес. H үшін толған қабат сегіз емес, екі электронмен байланысты. Байланыс түзілгенде энергия бөлінеді, бірақ басқа байланыстарды үзу қажет болғандықтан тұтас реакция энергияны қажет етуі мүмкін.',
  'Бөлшектер жұбының жаңа моделі: алыста 0, байланысқан күйде −7, тым жақында +9. Энергиясы ең төмен күйді таңдап, алыс бөлшектерге дейін байланысты үзуге қажет энергияны есептеңіз. Осы сандардан зат түсін білуге немесе барлық элемент жұбы осындай қосылыс түзетінін дәлелдеуге бола ма?',
  '**Тексеру:** −7; байланысты үзуге 7 шартты бірлік қажет. Түс пен басқа заттардың құрамы белгісіз. «Зарядтар → әрекеттесу → қолайлы аралық → төменірек энергия» тізбегін есіңізден жазыңыз. Шегін атаңыз: бұл нақты масштабсыз, барлық затты есептемейтін, берілген жұптың сапалық моделі.'
 ),
 'en': (
  'Why do atoms form chemical bonds?',
  'Connect attraction and repulsion with stability without saying that an atom “wants” something.',
  'Your notebook contains H and H₂. Both refer to hydrogen but denote different particles. Why can two neutral atoms remain together when each has zero total charge? Record a prediction. Use drawings only; do not produce hydrogen. Neither the colour of a model sphere nor a line connecting spheres is an observation of actual electrons.',
  'A **chemical bond** is an interaction holding atoms or ions together in a substance. Atoms contain positive nuclei and negative electrons. Charged particles attract or repel one another. **System energy** helps compare states: H₂ at a suitable separation has lower energy than two distant H atoms. Breaking that bond requires an input of energy.',
  'Predict whether the nuclei keep approaching until they merge. Does a neutral atom contain no charged particles? Mark where you know the total charge and where the internal structure matters. Check your prediction after using the data; simply writing “stability” does not explain the interaction.',
  'Analyse these invented values for two H atoms: far apart, energy 0; suitable separation, −4; too close, +6. These are teaching values in arbitrary units, not measured hydrogen energies. Draw the three arrangements and label attraction to electrons and repulsion between nuclei. Compare states of the same pair of atoms.',
  'The lowest of the three values is −4. At excessively short separation, repulsion matters and the energy rises to +6. A bond therefore does not mean that nuclei merge. Moving from −4 to distant atoms at energy 0 requires 4 arbitrary energy units in this model.',
  '**Mistake:** “atoms want eight electrons”. Atoms have no wishes; a filled-shell rule is a useful model for some simple cases, not a universal cause. Hydrogen has a filled shell with two electrons, not eight. Bond formation releases energy, but a whole reaction may require energy because other bonds must also be broken.',
  'A new pair model has energy 0 far apart, −7 in its bonded state and +9 when too close. Select the lowest-energy state and calculate the energy required to separate the particles. Can these values tell you the substance’s colour or prove that every pair of elements forms the same sort of compound?',
  '**Check:** −7; separation requires 7 arbitrary units. Colour and the composition of other substances remain unknown. Reconstruct “charges → interactions → suitable separation → lower energy”. State the limit: this is a qualitative model of the given pair, without real scales or a calculation of every possible substance.'
 )
},
{
 'stem': '24-shared-electron-pairs',
 'ru': (
  'Как общая электронная пара связывает атомы?',
  'Построим модели H₂, Cl₂ и HCl; посчитаем общие и неподелённые пары.',
  'Одна линия H—H и одна линия Cl—Cl похожи. Значит ли это, что в обеих молекулах всего два электрона? Вспомните: у нейтрального H один электрон, у Cl — семнадцать, из них семь на внешнем слое. Нас интересуют внешние электроны, а не все заряды атома. Образцы газов для урока не нужны.',
  '**Ковалентная связь** возникает при образовании общей электронной пары, взаимодействующей с обоими ядрами. **Валентные электроны** в наших простых примерах — электроны внешнего слоя. В записи Льюиса точки показывают их, а линия — общую пару. **Неподелённая пара** принадлежит области одного атома и не изображается соединяющей линией. Это модель распределения электронов, не их траектории.',
  'Предскажите число внешних электронов в H₂ и Cl₂. После образования связи должны ли исчезнуть электроны, которые не попали в общую пару? Напишите отдельно число электронов и число пар. Закрытый слой из двух электронов у H нельзя заменять требованием восьми.',
  'На бумаге обозначьте H одной точкой, Cl — семью. Соедините два H одной общей парой. Для двух Cl выделите общую пару, затем разместите остальные точки около атомов парами. Не используйте проволоку, источники питания и реактивы: рисунок нужен для проверки счёта, а не для получения вещества.',
  'H₂: 1+1=2 валентных электрона, одна общая пара. Cl₂: 7+7=14; одна общая пара и по три неподелённые пары у каждого Cl. Вокруг каждого Cl в этой модели учитываются восемь электронов, но общую пару нельзя дважды добавлять к общему числу 14. Внутренние электроны точками здесь не показаны.',
  '**Ошибка:** считать линию одним электроном или полагать, что все электроны Cl₂ находятся на связи. Для одинаковых атомов общей паре не приписываем постоянное владение одним атомом. Рисунок не показывает размер молекулы, скорость электронов или безопасность газа. Схема H—Cl также не требует обязательного переноса целого электрона.',
  'Новый случай HCl: H даёт в общий счёт один внешний электрон, Cl — семь. Найдите суммарное число, число общих пар и число неподелённых пар на Cl в простой схеме. Сравните молекулярную запись HCl с ионной записью Na⁺Cl⁻: одинаково ли в них описаны частицы?',
  '**Проверка:** восемь валентных электронов, одна общая пара и три неподелённые пары на Cl. HCl обозначает молекулу; запись Na⁺Cl⁻ указывает ионы в отношении 1:1. В следующем уроке разберём ионную решётку. По памяти нарисуйте H₂, затем Cl₂ и проверьте, что ни одна точка не потерялась.'
 ),
 'kz': (
  'Ортақ электрон жұбы атомдарды қалай байланыстырады?',
  'H₂, Cl₂ және HCl моделін құрып, ортақ және бөлінбеген жұптарды санаймыз.',
  'H—H пен Cl—Cl жазбаларында бір-бір сызық бар. Екі молекулада да барлығы екі электрон ғана бар ма? Еске түсіріңіз: бейтарап H атомында бір, Cl атомында он жеті электрон бар; Cl сыртқы қабатында жетеуі орналасқан. Біз барлық зарядты емес, сыртқы электрондарды қараймыз. Газ үлгілерінің қажеті жоқ.',
  '**Ковалентті байланыс** екі ядромен де әрекеттесетін ортақ электрон жұбы арқылы түзіледі. Біздің қарапайым мысалдарда **валенттік электрондар** — сыртқы қабаттағы электрондар. Льюис жазбасындағы нүктелер соларды, сызық ортақ жұпты көрсетеді. **Бөлінбеген жұп** бір атом маңында орналасады және атомдарды қосатын сызықпен көрсетілмейді. Бұл — траектория емес, электрондардың таралу моделі.',
  'H₂ және Cl₂ құрамындағы сыртқы электрондардың санын болжаңыз. Байланыс түзілгенде ортақ жұпқа кірмеген электрондар жоғалуы керек пе? Электрон санын және жұп санын бөлек жазыңыз. H үшін екі электронмен толатын қабатты сегіз электрон талабымен алмастыруға болмайды.',
  'Қағазда H маңына бір, Cl маңына жеті нүкте салыңыз. Екі H атомын бір ортақ жұппен қосыңыз. Екі Cl үшін ортақ жұпты белгілеп, қалған нүктелерді әр атом маңына жұптап орналастырыңыз. Сым, ток көзі және реактив қолданбаңыз: сурет зат алу үшін емес, санақты тексеру үшін керек.',
  'H₂: 1+1=2 валенттік электрон, бір ортақ жұп. Cl₂: 7+7=14; бір ортақ жұп және әр Cl маңында үш бөлінбеген жұп бар. Бұл модельде әр Cl айналасында сегіз электрон есептеледі, бірақ жалпы 14 санын есептегенде ортақ жұпты екі рет қосуға болмайды. Ішкі электрондар мұнда көрсетілмеген.',
  '**Қате:** сызықты бір электрон деп санау немесе Cl₂ электрондарының бәрі байланыста деу. Бірдей атомдардағы ортақ жұпты тек бір атомға тұрақты тиесілі деп қарастырмаймыз. Сурет молекула өлшемін, электрон жылдамдығын немесе газ қауіпсіздігін көрсетпейді. H—Cl сызбасы да міндетті түрде бүтін электронның берілуін білдірмейді.',
  'Жаңа жағдай HCl: H жалпы есепке бір, Cl жеті сыртқы электрон береді. Қарапайым сызбадағы барлық электронды, ортақ жұпты және Cl маңындағы бөлінбеген жұптарды санаңыз. HCl молекулалық жазбасын Na⁺Cl⁻ иондық жазбасымен салыстырыңыз: бөлшектер бірдей сипатталған ба?',
  '**Тексеру:** сегіз валенттік электрон, бір ортақ жұп, Cl маңында үш бөлінбеген жұп. HCl молекуланы білдіреді; Na⁺Cl⁻ жазбасы 1:1 қатынасындағы иондарды көрсетеді. Келесі сабақта иондық торды талдаймыз. H₂, кейін Cl₂ сызбасын жатқа салып, еш нүктенің жоғалмағанын тексеріңіз.'
 ),
 'en': (
  'How does a shared electron pair connect atoms?',
  'Model H₂, Cl₂ and HCl and count bonding and lone pairs.',
  'H—H and Cl—Cl each show one line. Does that mean each molecule contains just two electrons? Recall that neutral H has one electron, whereas Cl has seventeen, seven in its outer shell. We are focusing on outer electrons rather than every charge inside an atom. No gas samples are needed.',
  'A **covalent bond** involves a shared electron pair interacting with both nuclei. In our simple examples, **valence electrons** are outer-shell electrons. A Lewis diagram uses dots for these electrons and a line for a shared pair. A **lone pair** occupies the region of one atom and is not a connecting line. The diagram models electron distribution, not electron paths.',
  'Predict the number of outer electrons in H₂ and Cl₂. Should electrons outside the shared pair disappear when a bond forms? Write electron counts separately from pair counts. Do not replace hydrogen’s filled shell of two electrons with a requirement for eight.',
  'On paper, place one dot by H and seven by Cl. Connect two H atoms with one shared pair. For two Cl atoms, mark the shared pair and arrange the remaining dots in pairs around the atoms. Do not use wires, electrical supplies or chemicals: the drawing checks bookkeeping rather than producing a substance.',
  'H₂ has 1+1=2 valence electrons, one shared pair. Cl₂ has 7+7=14: one shared pair and three lone pairs on each Cl. Each Cl counts eight surrounding electrons in this model, but the shared pair must not be added twice to the total of fourteen. Inner electrons are not drawn here.',
  '**Mistake:** treating a line as one electron, or assuming all Cl₂ electrons belong to its bond. For identical atoms, do not assign permanent ownership of the shared pair to one atom. The drawing gives neither molecular dimensions nor electron speeds nor gas safety. H—Cl also does not require transfer of a whole electron.',
  'Try HCl: H contributes one outer electron to the total and Cl contributes seven. Find the total, the number of shared pairs and the number of lone pairs on Cl in a simple diagram. Compare molecular HCl with ionic Na⁺Cl⁻: do these notations describe the particles in the same way?',
  '**Check:** eight valence electrons, one shared pair and three lone pairs on Cl. HCl denotes a molecule; Na⁺Cl⁻ indicates ions in a 1:1 ratio. The next lesson examines an ionic lattice. Draw H₂ and then Cl₂ from memory, checking that no dot has been lost.'
 )
},
{
 'stem': '25-ionic-lattice-charge',
 'ru': (
  'Почему ионный кристалл не состоит из молекул NaCl?',
  'Разделим образование ионов, их притяжение и соотношение в протяжённой решётке.',
  'В уроке 22 вы уравновесили заряды Na⁺ и Cl⁻. Значит ли это, что внутри крупинки соли есть отдельные пары, изолированные от соседей? Представьте шахматный узор, который продолжается за край страницы. Он помогает увидеть повторение, но настоящий кристалл трёхмерен, а не плоский.',
  '**Катион** имеет положительный заряд, **анион** — отрицательный. **Ионная связь** связана с электростатическим притяжением противоположно заряженных ионов. **Ионная решётка** — повторяющееся расположение ионов в кристалле. **Формульная единица** передаёт простейшее целочисленное отношение ионов: для NaCl это 1:1. Она не означает отдельную молекулу внутри кристалла.',
  'Предскажите, изменится ли число протонов Na при образовании Na⁺. Сколько электронов у Cl⁻, если Z=17? Должен ли каждый Na⁺ взаимодействовать только с одним Cl⁻? Не смешивайте счёт зарядов для формулы со всеми взаимодействиями в пространстве.',
  'По данным Na: 11p,11e; Cl:17p,17e составьте бумажную запись переноса одного электрона. Затем нарисуйте чередующиеся знаки + и − в трёх рядах. Не проводите реакцию натрия с хлором и не получайте газы. Плоская схема показывает лишь фрагмент, а не полное координационное окружение.',
  'После переноса Na⁺ имеет 11p и 10e, Cl⁻ — 17p и 18e. Всего электронов осталось 28, протонов тоже 28. Заряды +1 и −1 дают ноль при отношении 1:1. Притяжение действует по всему кристаллу. Сам перенос электрона — удобный учёт, но полную энергетику образования соли нельзя объяснить одним этим шагом.',
  '**Ошибка:** рисовать отдельные молекулы NaCl или заключать «соль сохраняет опасные свойства хлора». Ионы и простые вещества имеют разные свойства. Это не делает неизвестный белый порошок безопасным: вид не доказывает состав. На схеме не меняйте размер иона только ради подсчёта и не принимайте цвет за настоящий.',
  'Новое условие: Mg²⁺ и O²⁻ образуют нейтральное соединение. Найдите простейшее отношение. Объясните, почему Mg₂O₂ можно сократить до MgO именно при записи отношения ионов, тогда как молекулярную формулу H₂O₂ нельзя произвольно сокращать до HO.',
  '**Проверка:** 1:1, MgO; суммарный заряд +2−2=0. Ионная формула передаёт простейшее отношение, а молекулярная H₂O₂ — действительное число атомов в одной молекуле. По памяти восстановите: образование зарядов, притяжение, решётка, формульная единица. Граница вывода: состав неизвестного образца этими рисунками ещё не установлен.'
 ),
 'kz': (
  'Иондық кристалл неліктен NaCl молекулаларынан тұрмайды?',
  'Иондардың түзілуін, тартылысын және кеңістіктік тордағы қатынасын ажыратамыз.',
  '22-сабақта Na⁺ және Cl⁻ зарядтарын теңестірдіңіз. Бұл тұз түйіршігінде көршілерден оқшауланған жеке жұптар бар деген сөз бе? Қағаз шетінен әрі жалғасатын шахмат өрнегін елестетіңіз. Ол қайталануды көрсетеді, бірақ нақты кристалл жазық емес, үш өлшемді.',
  '**Катионның** заряды оң, **анионның** заряды теріс. **Иондық байланыс** қарама-қарсы зарядты иондардың электростатикалық тартылуымен байланысты. **Иондық тор** — кристалдағы иондардың қайталанатын орналасуы. **Формулалық бірлік** иондардың ең қарапайым бүтін сандық қатынасын береді: NaCl үшін 1:1. Ол кристалл ішіндегі жеке молекула емес.',
  'Na⁺ түзілгенде Na протондарының саны өзгере ме? Z=17 болса, Cl⁻ неше электронға ие? Әр Na⁺ тек бір Cl⁻ ионымен ғана әрекеттесуі керек пе? Формула үшін заряд санауды кеңістіктегі барлық әрекеттесумен шатастырмаңыз.',
  'Na:11p,11e; Cl:17p,17e деректерінен бір электронның берілуін қағазға жазыңыз. Кейін үш қатарда + және − таңбаларын кезектестіріп салыңыз. Натрийді хлормен әрекеттестірмеңіз, газ алмаңыз. Жазық сызба толық көршілес ортаны емес, тек үзіндіні көрсетеді.',
  'Электрон берілгеннен кейін Na⁺:11p,10e; Cl⁻:17p,18e. Электрондардың жалпы саны 28 болып қалды, протон да 28. +1 және −1 зарядтары 1:1 қатынасында нөл береді. Тартылыс бүкіл кристалда әрекет етеді. Электронның берілуі — ыңғайлы есеп, бірақ тұз түзілуінің толық энергетикасын тек осы қадаммен түсіндіру жеткіліксіз.',
  '**Қате:** жеке NaCl молекулаларын салу немесе «тұз хлордың қауіпті қасиеттерін сақтайды» деу. Ион мен жай заттың қасиеті әртүрлі. Бұл белгісіз ақ ұнтақты қауіпсіз етпейді: сыртқы түр құрамды дәлелдемейді. Сызбада санақ үшін ион өлшемін өзгертпеңіз, түсті нақты түс деп қабылдамаңыз.',
  'Жаңа шарт: Mg²⁺ және O²⁻ бейтарап қосылыс түзеді. Ең қарапайым қатынасты табыңыз. Иондар қатынасын жазғанда Mg₂O₂-ні MgO-ға қысқартуға неге болады, ал H₂O₂ молекулалық формуласын HO-ға өз бетімен қысқартуға неге болмайтынын түсіндіріңіз.',
  '**Тексеру:** 1:1, MgO; жалпы заряд +2−2=0. Иондық формула ең қарапайым қатынасты, H₂O₂ молекулалық формуласы бір молекуладағы нақты атом санын береді. Заряд түзілуі, тартылыс, тор және формулалық бірлік тізбегін есіңізден қалпына келтіріңіз. Шегі: суреттер белгісіз үлгінің құрамын әлі анықтамайды.'
 ),
 'en': (
  'Why is an ionic crystal not made of NaCl molecules?',
  'Separate ion formation, attraction and the ratio in an extended lattice.',
  'In lesson 22 you balanced Na⁺ and Cl⁻ charges. Does a salt grain therefore contain individual pairs isolated from their neighbours? Imagine a checkerboard pattern continuing beyond the paper. It illustrates repetition, but a real crystal is three-dimensional rather than flat.',
  'A **cation** has positive charge and an **anion** has negative charge. **Ionic bonding** involves electrostatic attraction between oppositely charged ions. An **ionic lattice** is their repeating arrangement in a crystal. A **formula unit** gives the simplest whole-number ratio of ions: 1:1 for NaCl. It does not denote a separate molecule within the crystal.',
  'Predict whether the proton count of Na changes when Na⁺ forms. How many electrons does Cl⁻ have if Z=17? Must each Na⁺ interact with only one Cl⁻? Do not confuse the charge bookkeeping used for a formula with all interactions in space.',
  'Use Na:11p,11e and Cl:17p,17e to write a paper model of one electron transferring. Then draw alternating + and − signs in three rows. Do not react sodium with chlorine or produce gases. The flat drawing shows a fragment, not the full coordination environment.',
  'After transfer, Na⁺ has 11p and 10e; Cl⁻ has 17p and 18e. The total remains 28 electrons and 28 protons. Charges +1 and −1 sum to zero in a 1:1 ratio. Attraction operates throughout the crystal. Electron transfer is useful bookkeeping, but this step alone does not explain the entire energy balance of salt formation.',
  '**Mistake:** drawing separate NaCl molecules or claiming salt retains chlorine’s hazardous properties. Ions and elemental substances have different properties. That does not make an unidentified white powder safe: appearance cannot establish composition. Do not treat model colours as real colours or alter ion sizes to do the counting.',
  'A new condition specifies Mg²⁺ and O²⁻ in a neutral compound. Find the simplest ratio. Explain why Mg₂O₂ can be reduced to MgO when expressing an ion ratio, while the molecular formula H₂O₂ cannot arbitrarily be shortened to HO.',
  '**Check:** 1:1, MgO; the total charge is +2−2=0. An ionic formula gives the simplest ratio; molecular H₂O₂ gives the actual atom count in one molecule. Reconstruct ion formation, attraction, lattice and formula unit from memory. Limit: these drawings have not identified an unknown sample.'
 )
},
{
 'stem': '26-bond-polarity',
 'ru': (
  'Может ли общая электронная пара распределяться неравномерно?',
  'Отличим полярную связь от заряда иона и полярности всей молекулы.',
  'В H₂ и HCl мы нарисовали одну общую пару. Означает ли одинаковая линия одинаковое распределение электронной плотности? Сравните двух одинаковых участников и двух разных. Аналогия с перетягиванием помогает задать вопрос, но электрон не является верёвкой и не движется по нарисованной линии.',
  '**Электроотрицательность** — способность связанного атома притягивать электронную плотность. В **полярной ковалентной связи** она распределена неравномерно. Знаки **δ⁺ и δ⁻** обозначают частичные заряды, а не потерю и получение целого электрона. Полярность молекулы зависит и от полярности связей, и от их пространственного расположения.',
  'Предскажите, на какой стороне H—Cl будет δ⁻, если Cl сильнее притягивает общую электронную плотность. Должна ли молекула с полярными связями всегда быть полярной целиком? Сначала запишите оба ответа: второй проверим на разных формах молекул.',
  'Даны учебные сведения: в H₂ атомы одинаковые; в HCl электронная плотность смещена к Cl; CO₂ линейна O=C=O; H₂O имеет угловую форму. Для C—O и O—H смещение направлено к O. Стрелками отметьте направления на бумаге. Эти рисунки не требуют газов, нагревания или приборов.',
  'H₂: одинаковые атомы, связь неполярная. HCl: Hδ⁺—Clδ⁻, молекула полярная. В линейной CO₂ два одинаковых смещения направлены противоположно и компенсируются: молекула неполярная. В угловой H₂O смещения не компенсируются: молекула полярная. Одна формула без сведений о форме не заменяет пространственную модель.',
  '**Ошибка:** читать Hδ⁺ как H⁺ или превращать любую полярную связь в ионную. Частичный заряд не равен заряду отдельного иона. Шкалы электроотрицательности и пороговые числа — приближения; здесь мы не проводим универсальную границу одним числом. Полярность не является универсальным тестом токсичности или растворимости.',
  'Дана молекула X—Y—X с одинаковыми связями: электронная плотность смещена от Y к каждому X. Сначала расположение строго линейное, затем угловое. Буквы условные. Сравните компенсацию направлений при этих формах. Можно ли изменить ответ, просто повернув весь рисунок на листе?',
  '**Проверка:** в симметричной линейной модели направления компенсируются, в угловой — нет. Поворот всей молекулы не меняет её полярность. Восстановите три различия: связь и молекула; частичный и полный заряд; формула и форма. Объясните их на CO₂ и H₂O, не пользуясь выражением «есть кислород — значит полярно».'
 ),
 'kz': (
  'Ортақ электрон жұбы біркелкі таралмауы мүмкін бе?',
  'Полюсті байланысты ион зарядынан және тұтас молекула полюстілігінен ажыратамыз.',
  'H₂ мен HCl үшін бір ортақ жұп салдық. Бірдей сызық электрон тығыздығы да бірдей таралғанын білдіре ме? Бірдей екі қатысушыны екі әртүрлі қатысушымен салыстырыңыз. Арқан тарту ұқсастығы сұрақ қоюға көмектеседі, бірақ электрон арқан емес және сызылған сызықпен қозғалмайды.',
  '**Электртерістілік** — байланысқан атомның электрон тығыздығын тарту қабілеті. **Полюсті ковалентті байланыста** ол біркелкі таралмайды. **δ⁺ және δ⁻** толық электрон беру мен алуды емес, ішінара зарядтарды білдіреді. Молекула полюстілігі байланыстардың полюстілігіне де, кеңістікте орналасуына да тәуелді.',
  'Cl ортақ электрон тығыздығын күштірек тартса, H—Cl қай жағында δ⁻ болады? Полюсті байланысы бар молекула міндетті түрде тұтас полюсті ме? Екі болжамды да алдын ала жазыңыз: екінші сұрақты молекулалардың әртүрлі пішінімен тексереміз.',
  'Оқу деректері: H₂ атомдары бірдей; HCl тығыздығы Cl жағына ығысқан; CO₂ сызықтық O=C=O; H₂O бұрыштық пішінді. C—O және O—H үшін ығысу O жағына бағытталған. Қағазға бағыттарды көрсеткімен белгілеңіз. Газ, қыздыру және құрал қажет емес.',
  'H₂: атомдар бірдей, байланыс полюссіз. HCl: Hδ⁺—Clδ⁻, молекула полюсті. Сызықтық CO₂-дегі екі бірдей ығысу қарама-қарсы бағытталып, бірін-бірі өтейді: молекула полюссіз. Бұрыштық H₂O-да ығысулар өтелмейді: молекула полюсті. Пішін туралы дерексіз формуланың өзі кеңістіктік модельді алмастырмайды.',
  '**Қате:** Hδ⁺ жазбасын H⁺ деп оқу немесе кез келген полюсті байланысты иондыққа айналдыру. Ішінара заряд жеке ион зарядына тең емес. Электртерістілік шкалалары мен шекті сандар жуықтаулар; мұнда бір санмен әмбебап шекара қоймаймыз. Полюстілік уыттылықты не ерігіштікті анықтайтын әмбебап сынақ емес.',
  'X—Y—X молекуласының байланыстары бірдей, тығыздық Y-ден әр X-ке ығысқан. Бірінші орналасу қатаң сызықтық, екіншісі бұрыштық болсын. Әріптер шартты. Осы пішіндердегі бағыттардың өтелуін салыстырыңыз. Суретті қағазда тұтас бұрып қана жауапты өзгертуге бола ма?',
  '**Тексеру:** симметриялы сызықтық модельде бағыттар өтеледі, бұрыштық модельде өтелмейді. Молекуланы тұтас бұру оның полюстілігін өзгертпейді. Үш айырманы қалпына келтіріңіз: байланыс пен молекула; ішінара және толық заряд; формула мен пішін. «Оттек бар — демек полюсті» демей, CO₂ мен H₂O арқылы түсіндіріңіз.'
 ),
 'en': (
  'Can a shared electron pair be distributed unequally?',
  'Distinguish a polar bond from an ion charge and from the polarity of a whole molecule.',
  'We drew one shared pair in H₂ and in HCl. Does the same line imply the same distribution of electron density? Compare identical partners with different partners. A tug-of-war analogy helps pose the question, but an electron is not a rope and does not travel along the drawn line.',
  '**Electronegativity** is a bonded atom’s ability to attract electron density. In a **polar covalent bond**, that density is unevenly distributed. **δ⁺ and δ⁻** denote partial charges, not the loss or gain of a whole electron. Molecular polarity depends on both bond polarity and the spatial arrangement of bonds.',
  'Predict which end of H—Cl receives δ⁻ if Cl attracts shared electron density more strongly. Must a molecule containing polar bonds be polar overall? Record both answers before testing the second against different molecular shapes.',
  'Use these supplied facts: H₂ contains identical atoms; HCl electron density shifts towards Cl; CO₂ is linear, O=C=O; H₂O is bent. In C—O and O—H bonds, density shifts towards O. Mark the directions on paper. No gases, heating or electrical apparatus are required.',
  'H₂ has a nonpolar bond between identical atoms. HCl is Hδ⁺—Clδ⁻ and is a polar molecule. In linear CO₂, two identical shifts point in opposite directions and cancel: the molecule is nonpolar. In bent H₂O, the shifts do not cancel: the molecule is polar. A formula without shape information cannot replace a spatial model.',
  '**Mistake:** reading Hδ⁺ as H⁺ or treating every polar bond as ionic. Partial charge differs from the charge of a separate ion. Electronegativity scales and numerical cutoffs are approximations; we do not impose a universal boundary using one number here. Polarity is not a universal test of toxicity or solubility.',
  'A model molecule X—Y—X has identical bonds with density shifting from Y towards each X. First arrange it strictly linearly, then in a bent shape. The letters are placeholders. Compare cancellation of the directions. Could rotating the entire drawing on the paper change the answer?',
  '**Check:** the shifts cancel in the symmetric linear model but not in the bent model. Rotating an entire molecule does not change its polarity. Reconstruct three distinctions: bond versus molecule, partial versus full charge, formula versus shape. Explain them with CO₂ and H₂O rather than saying “oxygen means polar”.'
 )
},
{
 'stem': '27-structure-and-properties',
 'ru': (
  'Как строение помогает объяснять свойства вещества?',
  'Сравним молекулярные, ионные, металлические и сетчатые структуры по данным, а не одному ярлыку.',
  'Перед нами четыре карточки: лёд, соль, медь и алмаз. Все могут быть твёрдыми, но это не означает одинаковое устройство. Почему лёд плавится легче соли, а медь проводит ток в твёрдом состоянии? Используйте готовые сведения: ничего не плавьте, не нагревайте и не подключайте к сети.',
  '**Молекулярное вещество** состоит из молекул; взаимодействия между ними отличаются от связей внутри. **Ковалентная сетка** — протяжённая сеть связанных атомов. В **металлической модели** положительные атомные остовы связаны с общими подвижными электронами. В ионном кристалле повторяются катионы и анионы. Это типы моделей, а не четыре универсальных набора свойств без исключений.',
  'Предскажите: при плавлении льда разрываются ли все O—H связи? Достаточно ли знать слово «ковалентный», чтобы объявить любое вещество мягким и легкоплавким? Сначала объясните предполагаемые движущиеся частицы. При проводимости важно не просто наличие зарядов, а их возможность перемещаться через образец.',
  'Данные карточек: лёд при плавлении сохраняет молекулы H₂O; твёрдый NaCl плохо проводит ток, расплав проводит; медь проводит в твёрдом состоянии; алмаз очень твёрд, обычный ток почти не проводит. Расплав соли здесь только условие задачи, не инструкция к опыту. Составьте таблицу «частицы — связи — подвижные носители».',
  'Лёд: перестраиваются взаимодействия между молекулами, молекулы остаются. NaCl: в твёрдом кристалле ионы ограничены местами, в расплаве подвижны. Медь: электроны могут перемещаться по металлу. Алмаз: протяжённая сеть C—C; отдельных молекул алмаза нет. Разные носители и связи дают разные объяснения, хотя образцы твёрдые.',
  '**Ошибка:** «ковалентное вещество всегда состоит из молекул» или «проводит только металл». Графит — неметаллическая форма углерода с подвижными электронами и проводимостью; он не становится металлом. Чистота, температура и направление измерения тоже влияют. Школьная классификация помогает объяснять данные, но не заменяет проверку конкретного образца.',
  'Неизвестный образец X твёрдый и не проводит ток. Достаточно ли этого, чтобы выбрать между льдом, солью и алмазом? Предложите дополнительное безопасное сведение из справочника: например, температуру плавления при указанном давлении. Не предлагайте самостоятельное нагревание неизвестного вещества.',
  '**Проверка:** данных недостаточно: в заданных условиях все три могут плохо проводить. Справочные температура плавления и состав помогут различить варианты. Восстановите четыре модели, назвав частицы и подвижные носители. Отметьте, что растворение и проводимость раствора требуют отдельного анализа: не всякое вещество растворяется одинаково и не всякий раствор — ионный.'
 ),
 'kz': (
  'Құрылыс зат қасиеттерін түсіндіруге қалай көмектеседі?',
  'Молекулалық, иондық, металдық және торлы құрылымды бір атаумен емес, дерекпен салыстырамыз.',
  'Төрт карточка берілген: мұз, тұз, мыс және алмаз. Бәрі қатты бола алады, бірақ құрылысы бірдей емес. Мұз неге тұздан оңай балқиды, мыс қатты күйде неге ток өткізеді? Дайын мәліметті қолданыңыз: ештеңені балқытпаңыз, қыздырмаңыз, электр желісіне қоспаңыз.',
  '**Молекулалық зат** молекулалардан тұрады; олардың арасындағы әрекеттесу молекула ішіндегі байланыстан бөлек. **Ковалентті тор** — байланысқан атомдардың кеңістікте жалғасатын жүйесі. **Металдық модельде** оң атомдық қаңқалар ортақ қозғалмалы электрондармен байланысқан. Иондық кристалда катион мен анион қайталанады. Бұл — ешбір ерекшелігі жоқ төрт қасиет жиыны емес, модель түрлері.',
  'Болжаңыз: мұз балқығанда барлық O—H байланысы үзіле ме? «Ковалентті» деген сөзден кез келген затты жұмсақ әрі оңай балқитын деуге бола ма? Алдымен қандай бөлшектер қозғалады деп ойлайтыныңызды түсіндіріңіз. Өткізгіштікте зарядтың болуы ғана емес, үлгі арқылы орын ауыстыра алуы маңызды.',
  'Карточка деректері: мұз балқығанда H₂O молекулалары сақталады; қатты NaCl токты нашар, балқымасы жақсы өткізеді; мыс қатты күйде өткізеді; алмаз өте қатты, қалыпты токты іс жүзінде өткізбейді. Тұз балқымасы — тәжірибе нұсқауы емес, есеп шарты. «Бөлшек — байланыс — қозғалмалы тасымалдаушы» кестесін жасаңыз.',
  'Мұзда молекулалар арасындағы әрекеттесу қайта құрылады, молекулалар қалады. Қатты NaCl иондары орындарымен шектелген, балқымада қозғалады. Мыста электрондар металл бойымен орын ауыстырады. Алмазда C—C байланыстары тор түзеді, жеке алмаз молекулалары жоқ. Бәрі қатты болса да, байланыс пен тасымалдаушы әртүрлі түсіндіру береді.',
  '**Қате:** «ковалентті зат міндетті түрде молекулалы» немесе «токты тек металл өткізеді». Графит — қозғалмалы электрондары бар, өткізгіш көміртектің бейметалдық түрі; ол металлға айналмайды. Тазалық, температура, өлшеу бағыты да әсер етеді. Мектептік жіктеу деректі түсіндіреді, бірақ нақты үлгіні тексеруді алмастырмайды.',
  'Белгісіз X қатты және ток өткізбейді. Мұз, тұз, алмаз арасынан таңдау үшін осы жеткілікті ме? Анықтамалықтан қосымша қауіпсіз дерек ұсыныңыз: мысалы, көрсетілген қысымдағы балқу температурасы. Белгісіз затты өз бетімен қыздыруды ұсынбаңыз.',
  '**Тексеру:** дерек жеткіліксіз: берілген жағдайда үшеуі де нашар өткізуі мүмкін. Анықтамалық балқу температурасы мен құрам нұсқаларды ажыратуға көмектеседі. Бөлшектер мен қозғалмалы тасымалдаушыларды атап, төрт модельді еске түсіріңіз. Еру мен ерітінді өткізгіштігі бөлек талдауды қажет етеді: барлық зат бірдей ерімейді, әр ерітінді иондық емес.'
 ),
 'en': (
  'How does structure help explain a substance’s properties?',
  'Compare molecular, ionic, metallic and network structures using data rather than one label.',
  'Four cards describe ice, salt, copper and diamond. All can be solids, but that does not give them the same structure. Why does ice melt more easily than salt, while copper conducts as a solid? Use supplied information: do not melt, heat or connect anything to an electrical supply.',
  'A **molecular substance** consists of molecules; interactions between molecules differ from bonds inside them. A **covalent network** is an extended structure of bonded atoms. A **metallic model** has positive atomic cores associated with shared mobile electrons. An ionic crystal contains repeating cations and anions. These are models, not four exception-free lists of properties.',
  'Predict whether melting ice breaks every O—H bond. Does the word “covalent” prove that any substance is soft and melts easily? Explain which particles you think move. Conductivity requires charges that can move through the sample, not merely the existence of charged particles.',
  'Use the supplied cards: melting ice retains H₂O molecules; solid NaCl conducts poorly but its melt conducts; copper conducts as a solid; diamond is very hard and conducts ordinary current extremely poorly. Molten salt is a problem condition, not an experiment instruction. Build a “particles — interactions — mobile carriers” table.',
  'In ice, interactions between molecules rearrange while molecules remain. In solid NaCl, ions are constrained to lattice sites; in its melt they are mobile. Copper has electrons able to move through the metal. Diamond has an extended C—C network, with no separate diamond molecules. Different carriers and interactions explain different properties even though all samples are solids.',
  '**Mistake:** “a covalent substance always consists of molecules” or “only metals conduct”. Graphite is a nonmetallic form of carbon with mobile electrons and conductivity; it does not become a metal. Purity, temperature and measurement direction also matter. School classifications help interpret data but do not replace testing a particular sample.',
  'An unknown solid X does not conduct. Is this enough to choose between ice, salt and diamond? Suggest further safe reference information, such as melting temperature at a specified pressure. Do not propose heating an unknown substance yourself.',
  '**Check:** the evidence is insufficient: all three can conduct poorly under the stated conditions. Reference melting temperatures and composition could help distinguish them. Recall the four models, naming particles and mobile carriers. Dissolving and solution conductivity need separate analysis: substances differ in solubility and not every solution is ionic.'
 )
},
{
 'stem': '28-balance-ion-formulas',
 'ru': (
  'Как составить формулу по заданным зарядам ионов?',
  'Найдём минимальное нейтральное отношение без слепого переноса цифр крест-накрест.',
  'В карточках записаны Mg²⁺ и Cl⁻. Формула MgCl выглядит коротко, но суммарный заряд равен +1. Как получить нейтральное соединение, не меняя сами ионы? Рисуйте карточки на бумаге. Здесь заряды даны условием; их нельзя угадывать по неизвестному образцу или менять ради удобного ответа.',
  '**Индекс** справа внизу показывает число атомов или групп в формуле; единицу не пишут. Для ионного соединения он задаёт отношение. **Электронейтральность** означает, что сумма положительных и отрицательных зарядов равна нулю. Берём наименьшее целочисленное отношение заданных ионов. Это относится к простым заданным составам, а не ко всем возможным соединениям элементов.',
  'Предскажите: сколько Cl⁻ нужно на один Mg²⁺? Для Al³⁺ и O²⁻ достаточно одного иона каждого вида? Запишите сумму зарядов до выбора формулы. Если сумма не ноль, исправляйте количество ионов, а не величину заряда отдельного иона.',
  'В таблице запишите пару, количества, сумму зарядов и итог. Пары: Na⁺/Cl⁻; Mg²⁺/Cl⁻; Al³⁺/O²⁻; Ca²⁺/O²⁻. Подбирайте небольшие целые числа. Проверяйте не только ноль, но и отсутствие общего делителя у количеств. Работа с числами не требует соли, металлов или кислорода в сосуде.',
  'Mg²⁺ и два Cl⁻: +2+2×(−1)=0, формула MgCl₂. Для двух Al³⁺ и трёх O²⁻: 2×3+3×(−2)=0, Al₂O₃. Для Ca²⁺ и O²⁻ отношение 1:1, CaO. Запись Ca₂O₂ нейтральна, но не является простейшей записью данного ионного отношения.',
  '**Ошибка:** переносить знаки зарядов в нейтральную формулу или путать MgCl₂ с Mg²Cl. Заряд относится к иону, индекс — к количеству. Для Fe возможны разные ионы; формулу нельзя однозначно получить из названия элемента без заряда или других условий. Даже правильная формула не удостоверяет чистоту порошка.',
  'Новые пары: Ca²⁺/N³⁻ и K⁺/S²⁻. Найдите простейшие нейтральные отношения и формулы. Затем сравните Fe²⁺/Cl⁻ с Fe³⁺/Cl⁻: получится ли одна формула? Условные карточки N³⁻ и S²⁻ задают ионы в соединении, не свободные частицы для получения дома.',
  '**Проверка:** Ca₃N₂: 3×2−2×3=0; K₂S: 2−2=0. Железо даёт FeCl₂ и FeCl₃ соответственно. По памяти восстановите алгоритм «заряды → количества → сумма ноль → сокращение отношения → индексы». Поясните, почему этот алгоритм не разрешает сокращать заданную молекулярную формулу H₂O₂.'
 ),
 'kz': (
  'Берілген ион зарядтары бойынша формуланы қалай құрамыз?',
  'Сандарды ойланбай айқастырмай, ең кіші бейтарап қатынасты табамыз.',
  'Карточкаларда Mg²⁺ және Cl⁻ берілген. MgCl қысқа көрінеді, бірақ жалпы заряды +1. Иондардың өзін өзгертпей бейтарап қосылысты қалай аламыз? Карточкаларды қағазға салыңыз. Зарядтар шартта берілген; белгісіз үлгіге қарап болжауға немесе жауапқа ыңғайлап өзгертуге болмайды.',
  'Оң жақ төменгі **индекс** формуладағы атом не топ санын көрсетеді; бір саны жазылмайды. Иондық қосылыста ол қатынасты береді. **Электрбейтараптық** оң және теріс зарядтар қосындысының нөл болуын білдіреді. Берілген иондардың ең кіші бүтін сандық қатынасын аламыз. Бұл — барлық мүмкін қосылыстың емес, берілген қарапайым құрамның ережесі.',
  'Бір Mg²⁺ үшін неше Cl⁻ қажет? Al³⁺ пен O²⁻ үшін әрқайсысынан біреуі жеткілікті ме? Формула таңдаудан бұрын зарядтар қосындысын жазыңыз. Қосынды нөл болмаса, жеке ион зарядын емес, иондар санын түзетіңіз.',
  'Кестеге жұпты, сандарды, заряд қосындысын және нәтижені жазыңыз. Жұптар: Na⁺/Cl⁻; Mg²⁺/Cl⁻; Al³⁺/O²⁻; Ca²⁺/O²⁻. Кіші бүтін сандарды таңдаңыз. Нөлді ғана емес, сандардың ортақ бөлгіші жоқтығын да тексеріңіз. Сандық жұмысқа тұз, металл немесе ыдыстағы оттек қажет емес.',
  'Mg²⁺ және екі Cl⁻: +2+2×(−1)=0, формула MgCl₂. Екі Al³⁺ және үш O²⁻: 2×3+3×(−2)=0, Al₂O₃. Ca²⁺ пен O²⁻ қатынасы 1:1, CaO. Ca₂O₂ бейтарап болғанымен, осы иондық қатынастың ең қарапайым жазбасы емес.',
  '**Қате:** заряд таңбаларын бейтарап формулаға көшіру немесе MgCl₂-ні Mg²Cl-мен шатастыру. Заряд ионға, индекс санға қатысты. Fe үшін әртүрлі иондар бар; заряд не басқа шартсыз элемент атауынан бір ғана формула шығару мүмкін емес. Дұрыс формула да ұнтақтың тазалығын растамайды.',
  'Жаңа жұптар: Ca²⁺/N³⁻ және K⁺/S²⁻. Ең қарапайым бейтарап қатынас пен формуланы табыңыз. Кейін Fe²⁺/Cl⁻ пен Fe³⁺/Cl⁻ салыстырыңыз: формула бірдей ме? N³⁻ және S²⁻ карточкалары үйде алатын бос бөлшектерді емес, қосылыстағы иондарды береді.',
  '**Тексеру:** Ca₃N₂:3×2−2×3=0; K₂S:2−2=0. Темір үшін тиісінше FeCl₂ және FeCl₃. «Зарядтар → сандар → қосынды нөл → қатынасты қысқарту → индекстер» алгоритмін есіңізден жазыңыз. Неге бұл алгоритм берілген H₂O₂ молекулалық формуласын қысқартуға рұқсат етпейтінін түсіндіріңіз.'
 ),
 'en': (
  'How do we build a formula from given ion charges?',
  'Find the smallest neutral ratio instead of blindly swapping charge numbers.',
  'Two cards say Mg²⁺ and Cl⁻. MgCl looks concise, but its charges sum to +1. How can you represent a neutral compound without changing the ions themselves? Draw cards on paper. Charges are supplied here; do not infer them from an unidentified sample or alter them to obtain a convenient answer.',
  'A **subscript** at the lower right gives the number of atoms or groups; one is omitted. For an ionic compound it expresses a ratio. **Electrical neutrality** means positive and negative charges sum to zero. Use the smallest whole-number ratio of the specified ions. This handles simple given compositions, not every possible compound of the elements.',
  'Predict how many Cl⁻ ions balance one Mg²⁺. Is one ion of each enough for Al³⁺ and O²⁻? Write the charge sum before selecting a formula. If it is not zero, change ion counts rather than the charge of an individual ion.',
  'Make columns for ion pair, counts, charge sum and result. Use Na⁺/Cl⁻, Mg²⁺/Cl⁻, Al³⁺/O²⁻ and Ca²⁺/O²⁻. Try small whole numbers. Check both neutrality and the absence of a common divisor of the counts. Numerical work requires no salts, metals or oxygen in a container.',
  'One Mg²⁺ and two Cl⁻ give +2+2×(−1)=0: MgCl₂. Two Al³⁺ and three O²⁻ give 2×3+3×(−2)=0: Al₂O₃. Ca²⁺ and O²⁻ need 1:1, so CaO. Ca₂O₂ is neutral but is not the simplest expression of this ionic ratio.',
  '**Mistake:** putting charge signs into a neutral formula or confusing MgCl₂ with Mg²Cl. A charge describes an ion; a subscript describes a count. Iron can form different ions, so an element name alone cannot fix a unique formula without charges or other conditions. A correct formula also does not establish a powder’s purity.',
  'Try Ca²⁺/N³⁻ and K⁺/S²⁻. Find the simplest neutral ratios and formulas. Then compare Fe²⁺/Cl⁻ with Fe³⁺/Cl⁻: do they give the same formula? N³⁻ and S²⁻ cards specify ions in compounds, not free particles to produce at home.',
  '**Check:** Ca₃N₂ gives 3×2−2×3=0; K₂S gives 2−2=0. Iron gives FeCl₂ and FeCl₃ respectively. Recall “charges → counts → zero sum → reduce ratio → subscripts”. Explain why this procedure does not allow you to reduce a given molecular formula H₂O₂.'
 )
},
{
 'stem': '29-valency-and-formula',
 'ru': (
  'Чем валентность отличается от заряда иона?',
  'Посчитаем связи в простых молекулах и не будем приписывать валентности знак заряда.',
  'Для воды нарисовано H—O—H. У O две линии, но нет надписи O²⁻. Почему две связи не делают атом кислорода отдельным ионом с зарядом −2? Сравните модель молекулы с ионными карточками прошлого урока. В обоих случаях встречается число два, но оно отвечает на разные вопросы.',
  '**Валентность** в используемой школьной модели показывает число образованных связей с учётом кратности. В H₂O у H валентность I, у O — II; в CH₄ у C — IV. Валентность записываем римской цифрой без плюса и минуса. **Заряд иона** описывает избыток или недостаток электронов у всей частицы. Степень окисления — другой вид учёта, не тема этого урока.',
  'Предскажите валентность O в O=O, если двойная связь учитывается как две связи. Должна ли валентность элемента быть одной и той же во всех соединениях? У нас есть заданные структуры, поэтому проверять надо по ним, а не превращать запомненное число в универсальный закон.',
  'Перерисуйте H—H, H—O—H, O=O и CH₄ с четырьмя одинарными C—H. Для каждого атома подсчитайте одинарные линии как одну, двойную как две. Не путайте число соседних атомов и число связей с учётом кратности. Геометрию метана плоский крест не передаёт: реальная форма пространственная.',
  'В O₂ каждый O имеет одного соседа, но двойную связь и валентность II. В воде O имеет двух соседей и две одинарные связи, тоже II. В метане C имеет четыре одинарные связи и валентность IV. Из этого не следует наличие C⁴⁺ в молекуле CH₄: ковалентная модель использует общие пары.',
  '**Ошибка:** «валентность кислорода −2» или «один сосед всегда означает I». Минус относится к иному учёту, а двойная связь опровергает второе утверждение. Наш простой подсчёт не охватывает все типы связи и исключения; для новых сложных соединений требуется более подробная модель. Формула без структуры не всегда показывает все связи.',
  'Новые схемы: H—Cl и N≡N. Тройная связь обозначает три общие пары. Определите валентность Cl и каждого N, число соседей N и суммарное число общих пар в N₂. Можно ли из записи III вывести, что свободная молекула N₂ имеет заряд +3?',
  '**Проверка:** Cl — I, каждый N — III; у каждого N один сосед, между ними три общие пары. N₂ нейтральна; валентность не является её зарядом. Составьте две колонки «сколько связей?» и «каков заряд частицы?». Дайте по одному примеру, где одинаковое число не означает одинаковый физический смысл.'
 ),
 'kz': (
  'Валенттілік ион зарядынан несімен өзгеше?',
  'Қарапайым молекуладағы байланыстарды санап, валенттілікке заряд таңбасын қоспаймыз.',
  'Су үшін H—O—H салынған. O жанында екі сызық бар, бірақ O²⁻ жазылмаған. Екі байланыс неге оттек атомын заряды −2 болатын жеке ион етпейді? Молекула моделін өткен сабақтағы ион карточкасымен салыстырыңыз. Екі жерде де екі саны кездеседі, бірақ ол әртүрлі сұраққа жауап береді.',
  'Қолданылатын мектептік модельде **валенттілік** еселікті ескергендегі түзілген байланыстар санын көрсетеді. H₂O-да H валенттілігі I, O — II; CH₄-те C — IV. Валенттілікті плюс, минуссыз рим цифрымен жазамыз. **Ион заряды** тұтас бөлшектегі электронның артығын не жетіспеуін сипаттайды. Тотығу дәрежесі — бөлек есеп түрі, бұл сабақтың тақырыбы емес.',
  'O=O қос байланысын екі байланыс деп есептесек, O валенттілігі қандай? Элемент валенттілігі барлық қосылыста міндетті түрде бірдей ме? Бізде дайын құрылымдар бар; жатталған санды әмбебап заң етпей, солар бойынша тексеру керек.',
  'H—H, H—O—H, O=O және төрт дара C—H байланысы бар CH₄ көшірмесін салыңыз. Әр атом үшін дара сызықты бір, қос сызықты екі деп санаңыз. Көрші атомдар санын еселігі ескерілген байланыс санымен шатастырмаңыз. Жазық айқыш метанның нақты геометриясын бермейді: оның пішіні кеңістіктік.',
  'O₂-де әр O бір көршіге ие, бірақ қос байланыс түзеді, валенттілігі II. Суда O екі көршімен екі дара байланыс түзеді, валенттілігі де II. Метанда C төрт дара байланыс түзеді, валенттілігі IV. Бұдан CH₄ ішінде C⁴⁺ бар деген қорытынды шықпайды: ковалентті модель ортақ жұптарды қолданады.',
  '**Қате:** «оттектің валенттілігі −2» немесе «бір көрші әрқашан I береді». Минус басқа есепке қатысты, ал қос байланыс екінші пікірді жоққа шығарады. Қарапайым санақ барлық байланыс түрі мен ерекшелікті қамтымайды; күрделі қосылысқа тереңірек модель керек. Құрылымсыз формула барлық байланысты көрсете бермейді.',
  'Жаңа сызбалар: H—Cl және N≡N. Үш байланыс үш ортақ жұпты білдіреді. Cl және әр N валенттілігін, N көршілерінің санын, N₂-дегі барлық ортақ жұпты анықтаңыз. III жазбасынан бос N₂ молекуласының заряды +3 деген тұжырым шыға ма?',
  '**Тексеру:** Cl — I, әр N — III; әр N бір көршіге ие, олардың арасында үш ортақ жұп бар. N₂ бейтарап; валенттілік оның заряды емес. «Неше байланыс?» және «Бөлшек заряды қандай?» деген екі баған жасаңыз. Бірдей санның физикалық мағынасы бірдей емес жағдайға мысал келтіріңіз.'
 ),
 'en': (
  'How does valency differ from an ion’s charge?',
  'Count bonds in simple molecules without assigning charge signs to valency.',
  'Water is drawn H—O—H. Oxygen has two lines but is not labelled O²⁻. Why do two bonds not turn it into a separate ion with charge −2? Compare the molecule model with the previous lesson’s ion cards. Both involve the number two, but the numbers answer different questions.',
  'In the school model used here, **valency** counts bonds with their multiplicity included. H has valency I and O has II in H₂O; C has IV in CH₄. Write valency as a Roman numeral without a plus or minus sign. **Ion charge** describes excess or missing electrons in a whole particle. Oxidation state is another bookkeeping system, outside this lesson.',
  'Predict O’s valency in O=O if a double bond counts as two. Must an element have the same valency in every compound? We have specified structures: check those rather than turning one memorised number into a universal rule.',
  'Draw H—H, H—O—H, O=O and CH₄ with four single C—H bonds. Count each single line as one and each double line as two for each atom. Distinguish the number of neighbouring atoms from bond count with multiplicity. A flat cross does not show methane’s actual three-dimensional geometry.',
  'Each O in O₂ has one neighbour but a double bond and valency II. Oxygen in water has two neighbours and two single bonds, also II. Carbon in methane has four single bonds and valency IV. This does not mean CH₄ contains C⁴⁺: the covalent model uses shared pairs.',
  '**Mistake:** saying “oxygen’s valency is −2” or “one neighbour always means I”. The minus sign belongs to another accounting system, while a double bond refutes the second claim. This simple count does not cover all bonding types and exceptions; complex compounds need a fuller model. A formula without a structure does not always reveal every bond.',
  'New diagrams are H—Cl and N≡N. A triple bond represents three shared pairs. Find Cl’s and each N’s valency, the number of neighbours of each N and the total shared pairs in N₂. Does the numeral III imply that a free N₂ molecule has charge +3?',
  '**Check:** Cl is I and each N is III; each N has one neighbour, with three shared pairs between them. N₂ is neutral; valency is not its charge. Make two columns: “how many bonds?” and “what particle charge?”. Give an example where equal numbers have different physical meanings.'
 )
},
{
 'stem': '30-subscripts-brackets-coefficients',
 'ru': (
  'Как читать индексы, скобки и коэффициенты?',
  'Посчитаем атомы в 3H₂O и Ca(OH)₂, не меняя химический состав ради удобного числа.',
  'Две записи 2H₂O и H₂O₂ различаются не только местом цифры. В первой — две молекулы воды, во второй — одна молекула другого вещества. Как сделать счёт прозрачным? Разделите вопросы «сколько частиц?» и «что внутри одной частицы?». Работаем с обозначениями и бумагой, не с перекисью водорода.',
  '**Коэффициент** перед формулой умножает число всех её атомов или формульных единиц. **Скобки** объединяют группу, а индекс после скобок умножает всё внутри. В Ca(OH)₂ один Ca, два O и два H. OH⁻ — гидроксид-ион с зарядом −1; здесь группа сохраняется при записи состава соединения с Ca²⁺.',
  'Предскажите число H и O в 3H₂O. Затем посчитайте Ca, O и H в 2Ca(OH)₂. Нужно ли удваивать только H, потому что он стоит ближе к индексу 2? Запишите промежуточный состав одной формульной единицы, а уже затем примените коэффициент.',
  'Создайте таблицу «запись — одна единица — коэффициент — всего». Строки: H₂O; 3H₂O; Ca(OH)₂; 2Ca(OH)₂; Al₂(SO₄)₃. Последняя запись дана как готовая формула: SO₄²⁻ — сульфат-ион, его внутренние связи сейчас не выводим. Ни гидроксиды, ни кислоты смешивать не требуется.',
  '3H₂O даёт 6 H и 3 O. Одна единица Ca(OH)₂:1 Ca,2 O,2 H; две:2 Ca,4 O,4 H. Для Al₂(SO₄)₃:2 Al,3 S,12 O. Заряды двух Al³⁺ и трёх SO₄²⁻ взаимно уравновешены. Скобки умножили и S, и O, но не Al перед ними.',
  '**Ошибка:** менять индекс при подсчёте количества. H₂O₂ не является сокращённой записью двух H₂O: состав отдельной молекулы другой. У ионного соединения не придумываем молекулы Ca(OH)₂; считаем формульные единицы. У молекулы коэффициент может обозначать счёт частиц, а позже — количество вещества в уравнении при согласованных единицах.',
  'Самостоятельно разберите 3Mg(OH)₂ и 2Al₂(SO₄)₃. Посчитайте атомы каждого элемента во всём выражении. Затем сравните 2O и O₂: одинаково ли общее число атомов и одинаково ли их объединение? Объясните ответ словами, а не только итоговым числом.',
  '**Проверка:** 3Mg,6 O,6 H; 4 Al,6 S,24 O. В 2O и O₂ по два атома кислорода, но 2O обозначает два отдельных атома, O₂ — одну двухатомную молекулу. По памяти сформулируйте порядок: внутренняя группа → индекс → состав единицы → коэффициент. Не выводите безопасность вещества из простоты формулы.'
 ),
 'kz': (
  'Индекс, жақша және коэффициентті қалай оқимыз?',
  'Химиялық құрамды өзгертпей, 3H₂O және Ca(OH)₂ атомдарын санаймыз.',
  '2H₂O мен H₂O₂ тек сан орны арқылы ғана ерекшеленбейді. Біріншісі — екі су молекуласы, екіншісі — басқа заттың бір молекуласы. Санақты қалай түсінікті етеміз? «Неше бөлшек?» пен «Бір бөлшектің ішінде не бар?» сұрағын бөліңіз. Сутек асқын тотығымен емес, таңбалар және қағазбен жұмыс істейміз.',
  'Формула алдындағы **коэффициент** оның барлық атомының не формулалық бірлігінің санын көбейтеді. **Жақша** топты біріктіреді, кейінгі индекс ішіндегінің бәрін көбейтеді. Ca(OH)₂-де бір Ca, екі O, екі H бар. OH⁻ — заряды −1 гидроксид-ионы; Ca²⁺ қосылысының құрамын жазғанда бұл топ сақталады.',
  '3H₂O-дағы H пен O санын болжаңыз. Кейін 2Ca(OH)₂ ішіндегі Ca, O, H санаңыз. H индекс 2-ге жақын тұрғандықтан тек оны екі еселеу керек пе? Алдымен бір формулалық бірлік құрамын жазыңыз, содан соң коэффициентті қолданыңыз.',
  '«Жазба — бір бірлік — коэффициент — барлығы» кестесін жасаңыз. Жолдар: H₂O;3H₂O;Ca(OH)₂;2Ca(OH)₂;Al₂(SO₄)₃. Соңғысы дайын формула ретінде берілген: SO₄²⁻ — сульфат-ионы, ішкі байланыстарын қазір шығармаймыз. Гидроксидтер мен қышқылдарды араластыру қажет емес.',
  '3H₂O:6 H және 3 O. Ca(OH)₂ бір бірлігі:1 Ca,2 O,2 H; екі бірлігі:2 Ca,4 O,4 H. Al₂(SO₄)₃ үшін:2 Al,3 S,12 O. Екі Al³⁺ және үш SO₄²⁻ зарядтары теңеседі. Жақша S пен O санын көбейтеді, бірақ алдындағы Al санын көбейтпейді.',
  '**Қате:** санақ кезінде индексті өзгерту. H₂O₂ — екі H₂O қысқартуы емес: жеке молекула құрамы өзгеше. Иондық қосылысқа Ca(OH)₂ молекулаларын ойдан қоспай, формулалық бірліктерін санаймыз. Молекула алдында коэффициент бөлшек санын, кейін келісілген бірліктермен теңдеуде зат мөлшерін де білдіре алады.',
  '3Mg(OH)₂ және 2Al₂(SO₄)₃ өз бетіңізбен талдаңыз. Тұтас өрнектегі әр элемент атомын санаңыз. Кейін 2O және O₂ салыстырыңыз: жалпы атом саны бірдей ме, олардың бірігуі де бірдей ме? Жауапты санмен ғана емес, сөзбен түсіндіріңіз.',
  '**Тексеру:** 3 Mg,6 O,6 H;4 Al,6 S,24 O. 2O мен O₂-де екі-екіден оттек атомы бар, бірақ 2O екі жеке атомды, O₂ бір екіатомды молекуланы білдіреді. «Ішкі топ → индекс → бірлік құрамы → коэффициент» ретін жатқа айтыңыз. Формуланың қарапайым болуынан зат қауіпсіздігін шығармаңыз.'
 ),
 'en': (
  'How do we read subscripts, brackets and coefficients?',
  'Count atoms in 3H₂O and Ca(OH)₂ without changing composition to obtain a convenient number.',
  '2H₂O and H₂O₂ differ in more than the position of a digit. The first means two water molecules; the second means one molecule of another substance. Make the counting explicit by separating “how many particles?” from “what is inside one particle?”. Work with paper and symbols, not hydrogen peroxide.',
  'A **coefficient** before a formula multiplies all its atom counts or formula units. **Brackets** group symbols; the subscript after them multiplies everything inside. Ca(OH)₂ has one Ca, two O and two H. OH⁻ is the hydroxide ion, charge −1; this group is retained when recording its compound with Ca²⁺.',
  'Predict the H and O counts in 3H₂O, then Ca, O and H in 2Ca(OH)₂. Should only H be doubled because it is nearest the subscript? First write the composition of one formula unit, then apply the coefficient.',
  'Make a table with “expression — one unit — coefficient — total”. Use H₂O, 3H₂O, Ca(OH)₂, 2Ca(OH)₂ and Al₂(SO₄)₃. The last formula is supplied: SO₄²⁻ is sulfate, whose internal bonds we are not deriving here. No hydroxides or acids need to be mixed.',
  '3H₂O gives 6 H and 3 O. One Ca(OH)₂ unit has 1 Ca,2 O,2 H; two give 2 Ca,4 O,4 H. Al₂(SO₄)₃ contains 2 Al,3 S,12 O. Two Al³⁺ and three SO₄²⁻ charges balance. Brackets multiply both S and O but not Al outside them.',
  '**Mistake:** changing subscripts while counting quantities. H₂O₂ is not shorthand for two H₂O molecules: the individual molecular composition differs. Count formula units, rather than invented Ca(OH)₂ molecules, for the ionic compound. For molecules, a coefficient can count particles; later it can express amounts in equations with consistent units.',
  'Analyse 3Mg(OH)₂ and 2Al₂(SO₄)₃ independently. Count each element’s atoms across the whole expression. Then compare 2O with O₂: is the total atom count the same, and is the grouping the same? Explain in words as well as numbers.',
  '**Check:** 3 Mg,6 O,6 H;4 Al,6 S,24 O. Both 2O and O₂ contain two oxygen atoms, but 2O denotes separate atoms while O₂ denotes one diatomic molecule. Recall “inner group → subscript → unit composition → coefficient”. A simple formula does not establish a substance’s safety.'
 )
},
{
 'stem': '31-bonds-formulas-case-review',
 'ru': (
  'Дело о соли: что доказывают формулы и модели связей?',
  'Соберём весь блок в одно расследование и отделим корректный расчёт от определения неизвестного вещества.',
  'В дневнике три подписанные карточки: A — NaCl, B — H₂O, C — CO₂. Четвёртая D — неизвестный белый порошок. Автор дневника написал: «D белый, значит NaCl; CO₂ содержит кислород, значит полярный; у H₂O два H, значит заряд +2». Найдите разные ошибки, не выполняя опытов с порошком.',
  '**Модель** объясняет выбранные свойства и имеет границы. **Доказательство состава** требует подходящих данных об образце; название карточки не является результатом анализа D. Повторим различия: молекула и формульная единица; индекс и коэффициент; валентность и заряд; полярность связи и молекулы. Эти различия помогают не превращать правильный счёт в лишний вывод.',
  'До разбора предскажите: какую запись можно объяснить решёткой, какую — отдельной молекулой? Может ли белый цвет подтвердить формулу? Для CO₂ используйте данную линейную форму, а не только буквы. Для воды вспомните угловую форму. Каждое утверждение снабдите причиной и условием применимости.',
  'Сделайте четыре строки таблицы A–D и колонки «что дано», «модель», «что можно заключить», «чего не знаем». Для A заряды Na⁺ и Cl⁻ заданы; для B даны две одинарные O—H; для C две полярные C=O в линейной форме. Для D известно лишь описание внешнего вида; не трогайте, не пробуйте и не нюхайте неизвестное.',
  'A: отношение ионов 1:1 и сумма зарядов ноль; отдельные молекулы соли не требуются. B: нейтральная молекула с двумя H и одним O, валентность O — II, но это не заряд O²⁻. C: полярные связи компенсируются, молекула неполярна. D: по белому цвету состав и безопасность не определяются.',
  '**Ошибка:** одно правильное наблюдение объявлять полным исследованием. Даже если формула A дана верно, это ничего не доказывает о D. Проводимость и температура плавления тоже зависят от условий. Ионное или ковалентное строение само по себе не разрешает контакт с веществом. В этом блоке все задания решаются по данным и на бумаге.',
  'Итоговый перенос: по Ca²⁺ и Cl⁻ составьте формулу; посчитайте атомы в 2CaCl₂; сравните H₂O₂ и 2H₂O; укажите, что происходит со связью при её разрыве с точки зрения энергии. Затем придумайте вопрос к следующему блоку: как записывать превращение веществ, сохраняя число атомов каждого элемента?',
  '**Проверка:** CaCl₂;2 Ca и 4 Cl; H₂O₂ — одна молекула состава 2:2, 2H₂O — две молекулы, всего 4 H и 2 O. Разрыв связи требует энергии. Оцените пять умений отдельно: различать частицы; балансировать заряды; читать формулу; объяснять свойства; ограничивать вывод. Повторите новые случаи завтра и через неделю; следующий блок начнётся с признаков и записи реакций.'
 ),
 'kz': (
  'Тұздың сыры: формула мен байланыс моделі нені дәлелдейді?',
  'Блокты бір зерттеуге біріктіріп, дұрыс есептеуді белгісіз затты анықтаудан ажыратамыз.',
  'Күнделікте үш жазулы карточка бар: A — NaCl, B — H₂O, C — CO₂. Төртінші D — белгісіз ақ ұнтақ. Автор «D ақ, демек NaCl; CO₂-де оттек бар, демек полюсті; H₂O-да екі H бар, демек заряды +2» деп жазған. Ұнтақпен тәжірибе жасамай, әртүрлі қатені табыңыз.',
  '**Модель** таңдалған қасиеттерді түсіндіреді және шегі болады. **Құрам дәлелі** үлгі туралы тиісті деректі қажет етеді; карточка атауы D талдауының нәтижесі емес. Айырмаларды қайталаймыз: молекула мен формулалық бірлік; индекс пен коэффициент; валенттілік пен заряд; байланыс пен молекула полюстілігі. Бұлар дұрыс санақтан артық қорытынды жасамауға көмектеседі.',
  'Талдаудан бұрын болжаңыз: қай жазбаны тормен, қайсын жеке молекуламен түсіндіруге болады? Ақ түс формуланы растай ала ма? CO₂ үшін тек әріптерді емес, берілген сызықтық пішінді қолданыңыз. Су үшін бұрыштық пішінді еске түсіріңіз. Әр пікірге себеп пен қолданылу шартын қосыңыз.',
  'A–D төрт жолын, «не берілген», «модель», «қандай қорытынды», «не белгісіз» бағандарын жасаңыз. A үшін Na⁺ және Cl⁻ зарядтары; B үшін екі дара O—H; C үшін сызықтық пішіндегі екі полюсті C=O берілген. D туралы сыртқы түрі ғана белгілі; белгісіз затты ұстамаңыз, дәмін татпаңыз, иіскемеңіз.',
  'A: иондар қатынасы 1:1, зарядтар қосындысы нөл; жеке тұз молекулалары қажет емес. B: екі H және бір O бар бейтарап молекула; O валенттілігі II, бірақ бұл O²⁻ заряды емес. C: полюсті байланыстар бірін-бірі өтейді, молекула полюссіз. D: ақ түстен құрам мен қауіпсіздік анықталмайды.',
  '**Қате:** бір дұрыс бақылауды толық зерттеу деп қабылдау. A формуласы дұрыс берілсе де, бұл D туралы ештеңе дәлелдемейді. Өткізгіштік пен балқу температурасы да жағдайға тәуелді. Иондық не ковалентті құрылыс өздігінен затқа жанасуға рұқсат етпейді. Блоктағы барлық тапсырма дерекпен және қағазда орындалады.',
  'Қорытынды көшіру: Ca²⁺ және Cl⁻ бойынша формула құрыңыз; 2CaCl₂ атомдарын санаңыз; H₂O₂ мен 2H₂O салыстырыңыз; байланыс үзілуін энергия тұрғысынан түсіндіріңіз. Келесі блокқа сұрақ қойыңыз: әр элемент атомдарының санын сақтап, заттардың өзгерісін қалай жазамыз?',
  '**Тексеру:** CaCl₂;2 Ca және 4 Cl; H₂O₂ — құрамы 2:2 бір молекула, 2H₂O — екі молекула, барлығы 4 H,2 O. Байланысты үзу энергияны қажет етеді. Бес дағдыны бөлек бағалаңыз: бөлшекті ажырату; зарядты теңестіру; формуланы оқу; қасиетті түсіндіру; қорытындыны шектеу. Жаңа жағдайларды ертең және бір аптадан соң қайталаңыз; келесі блок реакция белгілері мен жазбасынан басталады.'
 ),
 'en': (
  'The salt case: what do formulas and bonding models prove?',
  'Combine the block into one investigation and distinguish correct calculation from identifying an unknown substance.',
  'Three labelled cards say A:NaCl, B:H₂O and C:CO₂. A fourth, D, is an unknown white powder. The notebook claims: “D is white, so it is NaCl; CO₂ contains oxygen, so it is polar; H₂O has two H, so its charge is +2.” Identify the different errors without experimenting on the powder.',
  'A **model** explains selected properties and has limits. **Evidence of composition** requires suitable sample data; a card’s name is not an analysis of D. Recall molecule versus formula unit, subscript versus coefficient, valency versus charge and bond versus molecular polarity. These distinctions stop correct counting from becoming an unsupported conclusion.',
  'Predict which notation calls for a lattice and which for separate molecules. Can white colour establish a formula? Use the specified linear shape for CO₂ rather than its letters alone. Recall water’s bent shape. Give a reason and a condition for each claim.',
  'Make four rows A–D with “given evidence”, “model”, “allowed conclusion” and “unknown”. A supplies Na⁺ and Cl⁻ charges; B supplies two single O—H bonds; C supplies two polar C=O bonds in a linear arrangement. D supplies appearance only. Do not touch, taste or smell the unidentified substance.',
  'A has a 1:1 ion ratio and zero charge sum; no separate salt molecules are needed. B is a neutral molecule with two H and one O; O’s valency II does not mean an O²⁻ ion. C has cancelling polar bonds and is nonpolar overall. D’s composition and safety cannot be established from white colour.',
  '**Mistake:** treating one correct observation as a complete investigation. Even a correctly supplied A formula proves nothing about D. Conductivity and melting temperature also depend on conditions. Ionic or covalent structure alone never authorises contact with a substance. Every task in this block can be solved using paper and given data.',
  'Final transfer: use Ca²⁺ and Cl⁻ to build a formula; count atoms in 2CaCl₂; compare H₂O₂ with 2H₂O; state the energy requirement of bond breaking. Pose a question for the next block: how can we record a transformation while preserving the number of each element’s atoms?',
  '**Check:** CaCl₂;2 Ca and 4 Cl. H₂O₂ is one molecule containing 2:2 atoms; 2H₂O is two molecules totalling 4 H and 2 O. Breaking a bond requires energy. Assess five skills separately: identifying particles, balancing charge, reading formulas, explaining properties and limiting claims. Retry new cases tomorrow and a week later; the next block begins with evidence and notation for reactions.'
 )
}
]
