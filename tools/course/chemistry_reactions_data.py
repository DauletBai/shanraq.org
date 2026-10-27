"""Authored lessons 32–40: reactions, evidence and conservation.

All procedures are paper analysis of supplied data, never home reactions.
Each locale has title, summary and eight teaching sections.
"""

LESSONS = [
{
 'stem': '32-change-and-evidence',
 'ru': (
  'Как отличить химическую реакцию от изменения формы?',
  'Отделим наблюдаемый признак от вывода о новых веществах; проверим пузыри, плавление и необратимость.',
  'В дневнике четыре случая: лёд стал водой; стакан разбился; из нагретой воды выходят пузырьки; железная деталь покрылась ржавчиной. Что действительно говорит об изменении состава? Ничего не нагревайте и не разбивайте: это готовые описания для анализа. Сначала запишите свой выбор и основание, не заглядывая в проверку.',
  '**Физическое изменение** меняет состояние, форму или распределение вещества без образования новых веществ. **Химическая реакция** приводит к появлению веществ другого химического состава или строения. **Признак** — наблюдение, например изменение цвета, выделение газа, появление осадка или изменение температуры. Сам по себе признак ещё не устанавливает причину: нужны условия и данные о составе.',
  'Предскажите: любые ли пузырьки доказывают реакцию? Всегда ли необратимое изменение химическое? Для каждого утверждения приведите возможное альтернативное объяснение. Важно не угадать ярлык, а указать, какая проверка позволила бы отличить объяснения. Не предлагайте пробу на вкус или запах неизвестного вещества.',
  'Составьте таблицу из четырёх описаний. Для льда дано, что до и после вещество H₂O; для осколков — тот же состав стекла; пузырьки при кипении состоят из водяного пара; в ржавчине обнаружены соединения железа с кислородом, отличающиеся от металла. Эти сведения даны условием, не получены вами из одного внешнего вида.',
  'Плавление льда, дробление стекла и кипение воды в этих случаях физические: химическая идентичность сохраняется. Образование ржавчины — химическое изменение. При этом ржавчина не обязана иметь одну простую формулу: её состав зависит от условий. Не записывайте всю ржавчину как чистый Fe₂O₃ без дополнительных сведений.',
  '**Ошибка:** «стакан не собрать обратно, значит была реакция» или «пузырьки всегда новое вещество». Обратимость не является определением химического изменения. При смешении двух окрашенных растворов оттенок может измениться просто из-за смешения цветов. Совокупность признаков полезна, но вывод о продуктах требует более подходящих данных.',
  'Новый случай: из охлаждённой газированной воды после открытия выходят пузырьки. По условию это ранее растворённый CO₂; новых веществ не обнаружено. Сравните с другим случаем, где анализ подтвердил образование CO₂ из исходных веществ. Может ли одинаковый внешний признак сопровождать разные процессы?',
  '**Проверка:** выход уже растворённого CO₂ можно объяснить физическим процессом, а появление нового CO₂ в результате превращения — химической реакцией. Пузырьки внешне похожи, но данные о происхождении газа разные. Восстановите цепочку «наблюдение → возможные объяснения → данные о составе → вывод с границей». Для неизвестного образца оставьте ответ открытым, если данных недостаточно.'
 ),
 'kz': (
  'Химиялық реакцияны пішін өзгерісінен қалай ажыратамыз?',
  'Бақыланатын белгіні жаңа зат туралы қорытындыдан ажыратып, көпіршік, балқу және қайтымсыздықты тексереміз.',
  'Күнделікте төрт жағдай бар: мұз суға айналды; стақан сынды; қыздырылған судан көпіршік шықты; темір бөлшек тот басты. Қайсысы құрам өзгергенін көрсетеді? Ештеңені қыздырмаңыз, сындырмаңыз: бұлар талдауға берілген дайын сипаттамалар. Тексеруді оқымай тұрып таңдауыңыз бен себебін жазыңыз.',
  '**Физикалық өзгеріс** жаңа зат түзбей, күйді, пішінді немесе заттың таралуын өзгертеді. **Химиялық реакция** химиялық құрамы не құрылысы басқа заттардың пайда болуына әкеледі. **Белгі** — түс өзгеруі, газ бөлінуі, тұнба пайда болуы немесе температура өзгеруі сияқты бақылау. Белгі жалғыз өзі себепті анықтамайды: жағдай мен құрам туралы дерек керек.',
  'Болжаңыз: кез келген көпіршік реакцияны дәлелдей ме? Әр қайтымсыз өзгеріс химиялық па? Әр пікірге басқа ықтимал түсіндіру ұсыныңыз. Атауды табу ғана емес, түсіндірулерді қандай тексеру ажырататынын көрсету маңызды. Белгісіз заттың дәмін татып не иіскеп тексеруді ұсынбаңыз.',
  'Төрт сипаттамаға кесте жасаңыз. Мұз үшін бастапқы және соңғы зат H₂O; сынықтарда шыны құрамы сақталған; қайнаудағы көпіршіктер су буынан тұрады; тоттан металдан өзгеше темірдің оттекті қосылыстары табылған. Бұл мәліметтер шартта берілген, оларды сіз сыртқы түрден ғана анықтаған жоқсыз.',
  'Осы жағдайларда мұздың балқуы, шынының ұсақталуы, судың қайнауы физикалық: заттың химиялық табиғаты сақталады. Тоттың түзілуі — химиялық өзгеріс. Тоттың бір ғана қарапайым формуласы болуы міндетті емес: құрамы жағдайға тәуелді. Қосымша дерексіз бүкіл тотты таза Fe₂O₃ деп жазбаңыз.',
  '**Қате:** «стақанды қайта жинау қиын, демек реакция болды» немесе «көпіршік әрқашан жаңа зат». Қайтымдылық химиялық өзгерістің анықтамасы емес. Екі боялған ерітінді араласқанда реңк түстер араласқандықтан да өзгереді. Белгілер жиыны пайдалы, бірақ өнімдер туралы қорытындыға тиісті дерек қажет.',
  'Жаңа жағдай: салқын газдалған суды ашқанда көпіршік шығады. Шарт бойынша бұл бұрын еріген CO₂, жаңа зат табылмаған. Оны талдау бастапқы заттардан жаңа CO₂ түзілгенін растаған басқа жағдаймен салыстырыңыз. Бір сыртқы белгі әртүрлі процеске сәйкес келе ала ма?',
  '**Тексеру:** бұрын еріген CO₂ шығуын физикалық процесс түсіндіреді, ал айналу нәтижесінде жаңа CO₂ пайда болуы — химиялық реакция. Көпіршік ұқсас, газдың шығу тегі туралы дерек бөлек. «Бақылау → ықтимал түсіндіру → құрам дерегі → шегі бар қорытынды» тізбегін қалпына келтіріңіз. Дерек жетпесе, белгісіз үлгі туралы жауапты ашық қалдырыңыз.'
 ),
 'en': (
  'How do we distinguish a reaction from a change of shape?',
  'Separate an observed clue from evidence of new substances, using bubbles, melting and irreversibility.',
  'Four notebook cases describe melting ice, a broken glass, bubbles in heated water and rust on an iron part. Which establishes a change in composition? Do not heat or break anything: these are supplied descriptions. Record your choice and reason before reading the check.',
  'A **physical change** alters state, shape or distribution without forming new substances. A **chemical reaction** produces substances with a different chemical composition or structure. A **clue** is an observation such as colour change, gas release, a precipitate or temperature change. A clue alone does not establish its cause; conditions and composition evidence matter.',
  'Predict whether every bubble proves a reaction and whether every irreversible change is chemical. Offer an alternative explanation for each claim. The aim is to identify evidence that would distinguish explanations, rather than guess a label. Never propose tasting or smelling an unidentified substance.',
  'Make a table for the four cases. The supplied evidence says ice and water are both H₂O, the fragments retain the glass composition, boiling bubbles contain water vapour, and rust contains iron–oxygen compounds different from the metal. Those facts are given, not established by appearance alone.',
  'Melting, breaking glass and boiling are physical changes in these cases: chemical identity remains. Rust formation is chemical. Rust need not have one simple formula because its composition depends on conditions. Do not label all rust as pure Fe₂O₃ without further evidence.',
  '**Mistake:** “the glass cannot easily be rebuilt, so it reacted” or “bubbles always mean a new substance”. Reversibility does not define chemical change. Mixing coloured solutions can alter the observed shade simply by mixing colours. Several clues help, but identifying products needs appropriate evidence.',
  'A new case describes bubbles when a chilled fizzy-water container is opened. The gas is specified as previously dissolved CO₂, with no new substance detected. Compare a second case in which analysis confirms new CO₂ formed from reactants. Can similar visible clues accompany different processes?',
  '**Check:** release of previously dissolved CO₂ can be a physical process; formation of new CO₂ through a transformation is chemical. Similar bubbles have different origins. Reconstruct “observation → possible explanations → composition evidence → limited conclusion”. Leave an unidentified sample unresolved when evidence is insufficient.'
 )
},
{
 'stem': '33-reactants-products-equation',
 'ru': (
  'Как записать, какие вещества вступили в реакцию и что получилось?',
  'Прочитаем словесную запись, формулы и обозначения состояния; отделим схему реакции от доказанного уравнения.',
  'В учебном отчёте сказано: магний взаимодействовал с кислородом и образовал оксид магния. Названия уже известны из условия. Как передать это короткой записью и не заменить продукты собственным предположением? Работа только с текстом и карточками; магний не поджигайте и кислород не получайте.',
  '**Реагенты** — исходные вещества, участвующие в реакции; их записывают слева. **Продукты** — образующиеся вещества; их записывают справа. Стрелка показывает направление рассматриваемого превращения. Знак плюс разделяет разные вещества. **Схема реакции** перечисляет формулы, а уравненная запись дополнительно согласует число атомов каждого элемента.',
  'Предскажите стороны записи для магния, кислорода и MgO. Должен ли кислород быть записан O или O₂, если условие говорит об обычном молекулярном кислороде? Можно ли поставить любую формулу продукта, а затем объявить её правильной только потому, что удалось подобрать коэффициенты?',
  'Запишите словами «магний + кислород → оксид магния». По данным формул составьте схему Mg + O₂ → MgO. Составьте две колонки атомов: слева Mg=1, O=2; справа Mg=1, O=1. Пока это не уравненная запись. Подпишите, что свойства и состав MgO даны в задаче, а не выведены из плюса и стрелки.',
  'Правильная уравненная запись для заданных веществ: 2Mg(s) + O₂(g) → 2MgO(s). Буква s означает твёрдое состояние, g — газ. В других примерах l обозначает жидкость, aq — водный раствор. Обозначение состояния не является ещё одним химическим элементом или индексом. Уравнение не показывает всех промежуточных стадий.',
  '**Ошибка:** читать плюс как образование смеси, а стрелку как знак равенства масс отдельных веществ. Плюс здесь перечисляет реагенты или продукты; массы каждого реагента не обязаны равняться массе каждого продукта. Уравнение с согласованными атомами не доказывает само по себе, что реакция идёт при любых условиях.',
  'Новый словесный случай: карбонат кальция при нагревании превращается в оксид кальция и углекислый газ. Формулы даны: CaCO₃, CaO, CO₂. Расположите их по сторонам, проверьте атомы и укажите, что обозначение условия нагревания относится к описанию реакции, а не разрешает домашний опыт.',
  '**Проверка:** CaCO₃(s) → CaO(s) + CO₂(g); слева и справа Ca=1, C=1, O=3. Нагревание — условие данного превращения, а не новый реагент с формулой «тепло». Закройте страницу и объясните четыре части записи: реагенты, продукты, коэффициенты, состояния. Для незнакомых веществ сначала потребуйте сведения о продуктах и условиях.'
 ),
 'kz': (
  'Реагенттер мен өнімдерді қалай жазамыз?',
  'Сөздік жазбаны, формула мен күй белгілерін оқып, реакция сұлбасын теңестірілген теңдеуден ажыратамыз.',
  'Оқу есебінде магний оттекпен әрекеттесіп, магний оксидін түзді делінген. Атаулар шарттан белгілі. Өнімді өз болжамымызбен алмастырмай, оны қысқа қалай жазамыз? Тек мәтін және карточкамен жұмыс істеңіз; магнийді жақпаңыз, оттек алмаңыз.',
  '**Реагенттер** — реакцияға қатысатын бастапқы заттар, олар сол жақта жазылады. **Өнімдер** — түзілетін заттар, олар оң жақта. Көрсеткі қарастырылған айналудың бағытын көрсетеді. Плюс әртүрлі заттарды бөледі. **Реакция сұлбасы** формулаларды атайды, ал теңестірілген жазба әр элемент атомдарының санын да сәйкестендіреді.',
  'Магний, оттек, MgO жазбаның қай жағында болатынын болжаңыз. Шарт кәдімгі молекулалық оттекті атаса, O ма, O₂ ме жазамыз? Кез келген өнім формуласын қойып, коэффициенттер табылғаны үшін ғана оны дұрыс деуге бола ма?',
  '«Магний + оттек → магний оксиді» деп жазыңыз. Берілген формулалармен Mg + O₂ → MgO сұлбасын құрыңыз. Атомдарды екі бағанға санаңыз: сол жақта Mg=1, O=2; оң жақта Mg=1, O=1. Бұл әлі теңестірілмеген. MgO қасиеті мен құрамы плюс пен көрсеткіден шығарылмай, есепте берілгенін белгілеңіз.',
  'Осы заттар үшін теңестірілген жазба: 2Mg(s) + O₂(g) → 2MgO(s). s — қатты күй, g — газ. Басқа мысалдарда l сұйықтықты, aq судағы ерітіндіні білдіреді. Күй белгісі жаңа элемент не индекс емес. Теңдеу барлық аралық сатыларды көрсетпейді.',
  '**Қате:** плюсті қоспаның түзілуі, көрсеткіні жеке заттар массасының теңдігі деп оқу. Плюс реагенттерді не өнімдерді тізеді; әр реагент массасы әр өнім массасына жеке-жеке тең болуға міндетті емес. Атомдары сәйкес теңдеу реакция кез келген жағдайда жүретінін өздігінен дәлелдемейді.',
  'Жаңа сөздік жағдай: кальций карбонаты қыздырғанда кальций оксиді мен көмірқышқыл газына айналады. CaCO₃, CaO, CO₂ формулалары берілген. Оларды тиісті жаққа қойып, атомдарды тексеріңіз. Қыздыру шарты реакция сипаттамасына қатысты, үйде тәжірибе жасауға рұқсат емес екенін айтыңыз.',
  '**Тексеру:** CaCO₃(s) → CaO(s) + CO₂(g); екі жақта Ca=1, C=1, O=3. Қыздыру — осы айналудың шарты, «жылу» формулалы жаңа реагент емес. Бетті жауып, реагент, өнім, коэффициент, күй белгілерін түсіндіріңіз. Бейтаныс заттар үшін алдымен өнім мен жағдай туралы дерек сұраңыз.'
 ),
 'en': (
  'How do we record reactants and products?',
  'Read word equations, formulas and state symbols, distinguishing a reaction scheme from a balanced equation.',
  'A teaching report says magnesium reacted with oxygen to form magnesium oxide. The substance names are supplied. How can we write this concisely without substituting our own guess about the product? Use text and cards only; do not burn magnesium or produce oxygen.',
  '**Reactants** are starting substances participating in a reaction, written on the left. **Products** are substances formed, written on the right. The arrow gives the direction of the transformation being considered. Plus signs separate different substances. A **reaction scheme** lists formulas; a balanced equation also reconciles each element’s atom count.',
  'Predict the sides for magnesium, oxygen and MgO. Should ordinary molecular oxygen be O or O₂? Could you choose any product formula and claim it is correct merely because coefficients can be found?',
  'Write “magnesium + oxygen → magnesium oxide”. With the supplied formulas this becomes Mg + O₂ → MgO. Count atoms in two columns: left Mg=1, O=2; right Mg=1, O=1. This scheme is not balanced. Note that MgO’s identity and properties were supplied, not deduced from a plus sign and arrow.',
  'The balanced equation for these substances is 2Mg(s) + O₂(g) → 2MgO(s). Here s denotes solid and g denotes gas. Other examples use l for liquid and aq for aqueous solution. State symbols are neither extra elements nor subscripts. The equation does not display every intermediate step.',
  '**Mistake:** reading plus as “a mixture is formed” or the arrow as equality between the masses of individual substances. Plus lists reactants or products; one reactant’s mass need not equal each product’s mass. Atom balance alone does not prove a reaction occurs under every condition.',
  'A new description says heated calcium carbonate forms calcium oxide and carbon dioxide. Given formulas are CaCO₃, CaO and CO₂. Place them on the correct sides, check atoms and explain that the heating condition describes the reaction rather than authorising a home experiment.',
  '**Check:** CaCO₃(s) → CaO(s) + CO₂(g), with Ca=1, C=1, O=3 on each side. Heating is a condition, not an additional reactant with a formula “heat”. Close the page and explain reactants, products, coefficients and states. For unfamiliar substances, first request evidence about products and conditions.'
 )
},
{
 'stem': '34-atoms-before-after',
 'ru': (
  'Что сохраняется, когда молекулы перестраиваются?',
  'Посчитаем атомы каждого элемента и покажем, почему число молекул может измениться.',
  'На карточках слева две молекулы H₂ и одна O₂, справа две H₂O. Кажется, трёх частиц стало две: не исчезло ли вещество? Нарисуйте исходные группы до объяснения. Это бумажная модель заданной реакции; она не требует получать, смешивать или поджигать газы.',
  'В обычной **химической реакции** атомы перегруппировываются, а ядра не превращаются в ядра других элементов. Поэтому сохраняется число атомов каждого элемента в рассматриваемой системе. **Число молекул** и число атомов — разные величины: изменение размера групп может менять первое при сохранении второго. Ядерные процессы в этом блоке не рассматриваются.',
  'Предскажите отдельно число H, число O, общее число атомов и число молекул с каждой стороны. Если общие суммы совпали, достаточно ли этого без проверки элементов? Можно ли обменять два атома H на два O только потому, что общее количество шариков останется тем же?',
  'Обозначьте H буквой H или светлой карточкой, O — буквой O или тёмной. Слева сгруппируйте H–H, H–H и O–O; справа составьте две группы H–O–H. Каждую карточку перенесите ровно один раз. Буква остаётся на ней при переносе: цвет не должен превращаться в название нового элемента.',
  'Для 2H₂ + O₂ → 2H₂O слева H=4, O=2, всего 6 атомов и 3 молекулы. Справа H=4, O=2, всего 6 атомов и 2 молекулы. Подсчёт подтверждает атомный баланс, но рисунок не показывает механизм реакции, реальные размеры, температуру или скорость превращения.',
  '**Ошибка:** требовать одинакового числа молекул либо проверять только общую сумму атомов. Пример с четырьмя H и двумя O нельзя заменить двумя H и четырьмя O: сумма шесть сохранится, состав — нет. В химических уравнениях нужно проверять каждый элемент, а для ионных записей позже добавится проверка суммарного заряда.',
  'Новый случай задан записью 2CO + O₂ → 2CO₂. Посчитайте C, O и молекулы. Затем объясните, почему правильный счёт не делает CO безопасным: его опасность определяется свойствами вещества, а не аккуратностью формулы. Никаких образцов газа для задачи не нужно.',
  '**Проверка:** с обеих сторон C=2, O=4; молекул слева 3, справа 2. Каждому C и O соответствует атом на другой стороне, но молекулы сгруппированы иначе. Восстановите из памяти две колонки элементного счёта. Для непонятной записи найдите первый элемент с несовпадением, а не пытайтесь исправлять её новым цветом шарика.'
 ),
 'kz': (
  'Молекулалар қайта құрылғанда не сақталады?',
  'Әр элемент атомын санап, молекула саны неліктен өзгеруі мүмкін екенін көрсетеміз.',
  'Карточкалардың сол жағында екі H₂ және бір O₂ молекуласы, оң жағында екі H₂O бар. Үш бөлшек екіге айналғандай: зат жоғалды ма? Түсіндіруден бұрын бастапқы топтарды салыңыз. Бұл — берілген реакцияның қағаз моделі; газ алуды, араластыруды, жағуды қажет етпейді.',
  'Кәдімгі **химиялық реакцияда** атомдар қайта топтасады, ядролар басқа элемент ядроларына айналмайды. Сондықтан қарастырылған жүйеде әр элемент атомдарының саны сақталады. **Молекула саны** мен атом саны бөлек: топ өлшемінің өзгеруі екіншісін сақтап, біріншісін өзгерте алады. Ядролық процестер бұл блокта қарастырылмайды.',
  'Екі жақтағы H, O, барлық атом және молекула санын бөлек болжаңыз. Барлық атомдардың жалпы саны тең болса, элементтерді тексермеуге бола ма? Шарлар саны өзгермесін деп екі H атомын екі O атомына алмастыруға бола ма?',
  'H-ты H әрпімен не ашық карточкамен, O-ны O әрпімен не күңгірт карточкамен белгілеңіз. Сол жақта H–H, H–H, O–O топтарын, оң жақта екі H–O–H тобын құрыңыз. Әр карточканы бір рет қана көшіріңіз. Әрпі сақталсын: түс жаңа элемент атауына айналмауы керек.',
  '2H₂ + O₂ → 2H₂O үшін сол жақта H=4, O=2, барлығы 6 атом және 3 молекула. Оң жақта H=4, O=2, барлығы 6 атом және 2 молекула. Санақ атомдық теңгерімді растайды, бірақ сурет механизмді, нақты өлшемді, температураны не жылдамдықты көрсетпейді.',
  '**Қате:** молекула саны бірдей болсын деу немесе тек атомдардың жалпы санын тексеру. Төрт H пен екі O-ны екі H пен төрт O-ға алмастыруға болмайды: алты саны сақталғанымен, құрам сақталмайды. Әр элементті тексеру керек; кейін иондық жазбаларда жалпы заряд та тексеріледі.',
  'Жаңа жағдай: 2CO + O₂ → 2CO₂. C, O және молекула санын табыңыз. Кейін дұрыс санақ неге CO-ны қауіпсіз етпейтінін түсіндіріңіз: қауіптілік формуланың ұқыптылығына емес, зат қасиетіне байланысты. Есепке газ үлгісі қажет емес.',
  '**Тексеру:** екі жақта C=2, O=4; сол жақта 3, оң жақта 2 молекула. Әр C мен O-ға екінші жақта атом сәйкес келеді, молекулалар ғана басқаша топтасқан. Элементтік санақтың екі бағанын есіңізден жасаңыз. Түсініксіз жазбада шар түсін ауыстырмай, саны сәйкес емес бірінші элементті табыңыз.'
 ),
 'en': (
  'What is conserved when molecules rearrange?',
  'Count each element’s atoms and explain why the number of molecules can change.',
  'Cards show two H₂ molecules and one O₂ on the left, and two H₂O on the right. Three particles appear to become two: has matter disappeared? Draw the starting groups first. This is a paper model of a specified reaction, not a reason to obtain, mix or ignite gases.',
  'In an ordinary **chemical reaction**, atoms regroup while nuclei do not turn into nuclei of different elements. The number of atoms of each element is therefore conserved within the system being considered. **Molecule count** differs from atom count: changing group sizes can change the former without changing the latter. Nuclear processes are outside this block.',
  'Predict H atoms, O atoms, total atoms and molecules on each side separately. Is an equal total enough without checking elements? Could two H atoms become two O atoms merely because the total number of model pieces would remain unchanged?',
  'Use H-labelled or light cards and O-labelled or dark cards. Make H–H, H–H and O–O on the left, then two H–O–H groups on the right. Transfer every card exactly once, retaining its letter. A model colour must not change an atom’s element identity.',
  'For 2H₂ + O₂ → 2H₂O, the left has H=4, O=2: six atoms in three molecules. The right has H=4, O=2: six atoms in two molecules. This confirms atom balance, not the reaction mechanism, real dimensions, temperature or speed.',
  '**Mistake:** demanding equal molecule counts or checking only total atoms. Four H plus two O cannot become two H plus four O: the total stays six but the composition changes. Check each element; later, ionic equations will also require a total-charge check.',
  'A new case gives 2CO + O₂ → 2CO₂. Count C, O and molecules. Explain why correct bookkeeping does not make CO safe: hazards depend on substance properties, not a neatly written formula. No gas sample is needed.',
  '**Check:** both sides have C=2, O=4; molecule counts are three and two. Every C and O has a counterpart, but grouping changes. Reconstruct two element-count columns from memory. For an unfamiliar equation, find the first mismatched element rather than changing a model piece’s colour.'
 )
},
{
 'stem': '35-balance-with-coefficients',
 'ru': (
  'Как уравнять реакцию, не меняя формулы веществ?',
  'Подберём целые коэффициенты, проверим каждый элемент и сократим общий множитель.',
  'В схеме Al + O₂ → Al₂O₃ кислорода слева два атома, справа три. Запись Al₂O₂ дала бы удобный счёт, но изменила бы заданный продукт. Как сохранить формулы? Используйте только бумажные карточки и счёт; работа с алюминием, кислородом и нагреванием не требуется.',
  '**Коэффициент** перед формулой умножает все атомы этой формулы. При уравнивании состав заданных реагентов и продуктов сохраняется: менять можно коэффициенты, а не индексы. Для школьной итоговой записи выбираем наименьшие положительные целые коэффициенты. Коэффициент 1 обычно не пишется; дроби допустимы как промежуточный способ счёта, но здесь приводятся к целым.',
  'Предскажите, сколько атомов O удобно взять, чтобы их можно было сгруппировать и по два, и по три. После согласования O сколько потребуется Al? Не считайте баланс окончательным, пока не проверили второй элемент. Сначала запишите свой ход, затем сравните его с примером.',
  'Для кислорода возьмите шесть атомов: это три O₂ и две формульные единицы Al₂O₃. В двух Al₂O₃ содержится четыре Al. Запишите коэффициенты 4 : 3 : 2 и составьте проверочную таблицу. Фраза «формульная единица» здесь не превращает твёрдый оксид в отдельные молекулы.',
  'Получаем 4Al + 3O₂ → 2Al₂O₃. Слева и справа Al=4, O=6. Если все коэффициенты удвоить, запись 8 : 6 : 4 тоже сохранит атомы, но общий делитель 2 позволяет вернуться к простейшему отношению 4: 3: 2. Коэффициенты передают отношение количеств, а не равные массы в граммах.',
  '**Ошибка:** исправлять недостаток O заменой Al₂O₃ на Al₂O₂; забывать, что коэффициент действует и на Al, и на O; проверять только один элемент. Алгоритм уравнивания работает с заданными продуктами и не выбирает за нас, какое вещество реально образуется при определённых условиях.',
  'Самостоятельно уравняйте Fe + O₂ → Fe₂O₃, считая Fe₂O₃ заданным продуктом, и N₂ + H₂ → NH₃. Не связывайте первый пример с утверждением, что вся ржавчина — чистый Fe₂O₃. Для каждой записи напишите отдельные суммы по элементам и проверьте общий делитель коэффициентов.',
  '**Проверка:** 4Fe + 3O₂ → 2Fe₂O₃: Fe=4, O=6; N₂ + 3H₂ → 2NH₃: N=2, H=6. Условия синтеза аммиака здесь не задаются для выполнения, а продукты не предсказываются одним балансом. Воспроизведите порядок «формулы фиксированы → коэффициенты → каждый элемент → простейшее отношение» и объясните ошибку в изменении индекса.'
 ),
 'kz': (
  'Зат формуласын өзгертпей реакцияны қалай теңестіреміз?',
  'Бүтін коэффициенттерді таңдап, әр элементті тексереміз және ортақ көбейткішті қысқартамыз.',
  'Al + O₂ → Al₂O₃ сұлбасында сол жақта екі, оң жақта үш O атомы бар. Al₂O₂ жазбасы санақты жеңілдетер еді, бірақ берілген өнімді өзгертеді. Формулаларды қалай сақтаймыз? Тек қағаз карточкасы мен санақты қолданыңыз; алюминий, оттек және қыздыру қажет емес.',
  'Формула алдындағы **коэффициент** оның барлық атомын көбейтеді. Теңестіргенде берілген реагент пен өнім құрамы сақталады: индексті емес, коэффициентті өзгертеміз. Мектептік соңғы жазба үшін ең кіші оң бүтін коэффициенттерді таңдаймыз. 1 әдетте жазылмайды; аралық санақта бөлшек қолданылуы мүмкін, бірақ мұнда бүтінге келтіреміз.',
  'O атомдарын екіден де, үштен де топтау үшін қанша алу ыңғайлы екенін болжаңыз. O теңескеннен кейін қанша Al керек? Екінші элементті тексермей теңгерім бітті демеңіз. Алдымен өз қадамыңызды жазыңыз, кейін мысалмен салыстырыңыз.',
  'Оттек үшін алты атом алыңыз: бұл үш O₂ және екі Al₂O₃ формулалық бірлігі. Екі Al₂O₃ ішінде төрт Al бар. 4 : 3 : 2 коэффициенттерін жазып, тексеру кестесін құрыңыз. «Формулалық бірлік» қатты оксидті жеке молекулаларға айналдырмайды.',
  '4Al + 3O₂ → 2Al₂O₃ шығады. Екі жақта Al=4, O=6. Коэффициенттерді екі еселесек, 8 : 6 : 4 жазбасы да атомдарды сақтайды, бірақ ортақ бөлгіш 2 арқылы ең қарапайым 4: 3: 2 қатынасына ораламыз. Коэффициенттер грамдағы тең массаны емес, мөлшерлер қатынасын береді.',
  '**Қате:** O жетпегендіктен Al₂O₃-ті Al₂O₂-ге ауыстыру; коэффициенттің Al-ға да, O-ға да әсерін ұмыту; бір элементті ғана тексеру. Алгоритм берілген өніммен жұмыс істейді, нақты жағдайда қандай зат түзілетінін өздігінен таңдамайды.',
  'Өнім Fe₂O₃ деп берілген Fe + O₂ → Fe₂O₃ және N₂ + H₂ → NH₃ жазбаларын теңестіріңіз. Бірінші мысалды барлық тот таза Fe₂O₃ деген пікірмен байланыстырмаңыз. Әр теңдеу үшін элементтер қосындысын бөлек жазып, коэффициенттердің ортақ бөлгішін тексеріңіз.',
  '**Тексеру:** 4Fe + 3O₂ → 2Fe₂O₃: Fe=4, O=6; N₂ + 3H₂ → 2NH₃: N=2, H=6. Аммиак синтезінің жағдайы орындауға берілмеген; өнімді жалғыз теңгерім болжамайды. «Формулалар бекітілген → коэффициенттер → әр элемент → қарапайым қатынас» ретін жатқа айтып, индекс өзгертудегі қатені түсіндіріңіз.'
 ),
 'en': (
  'How can we balance a reaction without changing formulas?',
  'Choose integer coefficients, check every element and reduce any common factor.',
  'In Al + O₂ → Al₂O₃, oxygen has two atoms on the left and three on the right. Writing Al₂O₂ would simplify counting but change the specified product. How do we preserve formulas? Use paper cards and arithmetic only; no aluminium, oxygen or heating is required.',
  'A **coefficient** before a formula multiplies every atom in it. Balancing preserves the compositions of the specified reactants and products: change coefficients, not subscripts. Our final school notation uses the smallest positive integers. One is normally omitted; fractions may help intermediate bookkeeping, but we convert them to integers here.',
  'Predict a convenient number of O atoms divisible into both pairs and threes. Once O is reconciled, how many Al atoms are required? Check the second element before declaring success. Record your method before comparing it with the worked example.',
  'Choose six O atoms: three O₂ and two formula units of Al₂O₃. Two Al₂O₃ units contain four Al atoms. Write coefficients 4 : 3 : 2 and build a check table. “Formula unit” does not turn the solid oxide into isolated molecules.',
  'The result is 4Al + 3O₂ → 2Al₂O₃, with Al=4, O=6 on both sides. Doubling all coefficients gives 8 : 6 : 4 and still balances atoms, but the common divisor two reduces it to 4: 3: 2. Coefficients express relative amounts, not equal masses in grams.',
  '**Mistake:** replacing Al₂O₃ with Al₂O₂ to repair an O shortage; forgetting the coefficient multiplies both Al and O; checking only one element. The method operates on given products, rather than determining which product actually forms under particular conditions.',
  'Balance Fe + O₂ → Fe₂O₃, taking Fe₂O₃ as the specified product, and N₂ + H₂ → NH₃. The first example does not claim all rust is pure Fe₂O₃. Show element totals and check common factors separately for each equation.',
  '**Check:** 4Fe + 3O₂ → 2Fe₂O₃ has Fe=4, O=6; N₂ + 3H₂ → 2NH₃ has N=2, H=6. Ammonia-production conditions are not instructions to perform a reaction, and balance alone does not predict products. Recall “fixed formulas → coefficients → each element → simplest ratio” and explain why changing a subscript is wrong.'
 )
},
{
 'stem': '36-brackets-and-state-symbols',
 'ru': (
  'Как уравнивать записи со скобками и не потерять атомы?',
  'Разберём Ca(OH)₂ и Fe₂(SO₄)₃; сохраним формулы и отличим водный раствор от отдельной молекулы.',
  'Схема Ca(OH)₂ + HCl → CaCl₂ + H₂O содержит знакомые скобки. Но если поставить коэффициент 2 перед водой, что случится с числом H и O? Не смешивайте вещества: это задача по данным. Названия реагентов не являются приглашением брать кислоты или гидроксиды дома.',
  'Сначала **раскрытие скобок** даёт состав одной формульной единицы. Затем коэффициент умножает всё целиком. Ca(OH)₂ содержит Ca=1, O=2, H=2. Группу можно временно считать как целое при уравнивании только тогда, когда она в неизменном составе встречается по обе стороны; после этого всё равно проверяют отдельные элементы.',
  'Предскажите, можно ли считать OH неизменной группой, если справа дана вода H₂O. Что означает aq после HCl? В водном растворе частицы могут быть ионами: условная молекулярная запись уравнения не доказывает, что там плавают отдельные молекулы всех перечисленных солей.',
  'Для схемы найдите сначала Ca и Cl. Одна CaCl₂ требует двух HCl. Слева теперь H=2+2=4, O=2, поэтому справа две H₂O. Составьте строки Ca, H, O, Cl с суммами по каждой стороне. Слова «водный раствор» и s, l, g, aq не включаются в арифметику атомов.',
  'Получаем Ca(OH)₂ + 2HCl → CaCl₂ + 2H₂O. С каждой стороны Ca=1, H=4, O=2, Cl=2. Мы не сохраняли группу OH как готовый блок: её H и O вошли в воду. Коэффициент 2 перед HCl относится и к H, и к Cl. Это школьная суммарная запись для данных продуктов, не описание всех частиц раствора.',
  '**Ошибка:** умножить только последний символ в скобках; прибавить индекс к количеству вместо умножения; записать воду H₂O₂ ради баланса. Состояния веществ и условия реакции задаются отдельно. Осадок — образовавшаяся твёрдая фаза; знак s не является самостоятельным доказательством, что осадок обязательно появится при любом разбавлении.',
  'Новая схема с заданными продуктами: Fe₂(SO₄)₃ + NaOH → Fe(OH)₃ + Na₂SO₄. По условию в растворе образуется осадок Fe(OH)₃. Попробуйте сначала считать неизменённые группы SO₄, затем Fe и Na, а в конце раскройте скобки и проверьте O и H. Работайте только на бумаге.',
  '**Проверка:** Fe₂(SO₄)₃ + 6NaOH → 2Fe(OH)₃ + 3Na₂SO₄. В обеих частях Fe=2, S=3, Na=6, O=18, H=6. Уравнение не предписывает получать осадок самостоятельно. По памяти объясните, почему SO₄ здесь удобно считать группой, а OH в первом примере нельзя переносить целиком в воду без пересчёта.'
 ),
 'kz': (
  'Жақшалы жазбаны атом жоғалтпай қалай теңестіреміз?',
  'Ca(OH)₂ мен Fe₂(SO₄)₃ талдап, формулаларды сақтаймыз және судағы ерітіндіні жеке молекуладан ажыратамыз.',
  'Ca(OH)₂ + HCl → CaCl₂ + H₂O сұлбасында таныс жақша бар. Су алдына 2 қойсақ, H пен O саны не болады? Заттарды араластырмаңыз: бұл дерекке сүйенген есеп. Реагент атауы үйде қышқыл не гидроксид алуға шақырмайды.',
  'Алдымен **жақшаны ашу** бір формулалық бірліктің құрамын береді. Кейін коэффициент бәрін көбейтеді. Ca(OH)₂ құрамында Ca=1, O=2, H=2. Топты уақытша тұтас санауға ол екі жақта да құрамы өзгермей кездескенде ғана болады; соңында жеке элементтер бәрібір тексеріледі.',
  'Оң жақта H₂O суы берілсе, OH-ты өзгермейтін топ деп санауға бола ма? HCl-дан кейінгі aq нені білдіреді? Судағы бөлшектер ион болуы мүмкін: теңдеудің шартты молекулалық жазбасы аталған барлық тұздың жеке молекулалары жүзіп жүргенін дәлелдемейді.',
  'Алдымен Ca мен Cl табыңыз. Бір CaCl₂ үшін екі HCl қажет. Сол жақта енді H=2+2=4, O=2, демек оң жақта екі H₂O. Ca, H, O, Cl жолдарын жасап, екі жақ қосындысын жазыңыз. «Судағы ерітінді» мен s, l, g, aq белгілері атом есебіне кірмейді.',
  'Ca(OH)₂ + 2HCl → CaCl₂ + 2H₂O шығады. Әр жақта Ca=1, H=4, O=2, Cl=2. OH дайын блок ретінде сақталған жоқ: оның H пен O атомдары суға кірді. HCl алдындағы 2 H-қа да, Cl-ға да қатысты. Бұл берілген өнімдердің мектептік жиынтық жазбасы, ерітіндідегі барлық бөлшектің сипаттамасы емес.',
  '**Қате:** жақшадағы соңғы таңбаны ғана көбейту; көбейтудің орнына индексті қосу; теңгерім үшін суды H₂O₂ деп жазу. Күй мен реакция жағдайы бөлек беріледі. Тұнба — түзілген қатты фаза; s белгісі кез келген сұйылтуда тұнба міндетті түрде шығатынын дәлелдемейді.',
  'Өнімдері берілген жаңа сұлба: Fe₂(SO₄)₃ + NaOH → Fe(OH)₃ + Na₂SO₄. Шарт бойынша ерітіндіде Fe(OH)₃ тұнбасы түзіледі. Алдымен өзгермейтін SO₄ топтарын, кейін Fe мен Na санаңыз; соңында жақшаны ашып, O және H тексеріңіз. Тек қағазда жұмыс істеңіз.',
  '**Тексеру:** Fe₂(SO₄)₃ + 6NaOH → 2Fe(OH)₃ + 3Na₂SO₄. Екі жақта Fe=2, S=3, Na=6, O=18, H=6. Теңдеу тұнбаны өз бетімен алуды тапсырмайды. Мұнда SO₄-ті топ ретінде санау неге ыңғайлы, ал бірінші мысалдағы OH-ты суға қайта санамай тұтас көшіруге неге болмайтынын жатқа түсіндіріңіз.'
 ),
 'en': (
  'How do we balance bracketed formulas without losing atoms?',
  'Analyse Ca(OH)₂ and Fe₂(SO₄)₃, preserving formulas and distinguishing aqueous notation from individual molecules.',
  'Ca(OH)₂ + HCl → CaCl₂ + H₂O contains familiar brackets. If we place two before water, what happens to H and O counts? Do not mix substances: this is a data exercise. Reagent names do not invite you to obtain acids or hydroxides at home.',
  'First **expand brackets** to find the composition of one formula unit. Then apply the coefficient to the whole formula. Ca(OH)₂ contains Ca=1, O=2, H=2. You may temporarily count a group as a unit only when it appears unchanged on both sides; individual elements must still be checked afterwards.',
  'Predict whether OH can be treated as unchanged when the product is H₂O. What does aq after HCl mean? Aqueous particles may be ions: a conventional molecular equation does not prove that separate molecules of every listed salt exist in solution.',
  'Start with Ca and Cl. One CaCl₂ requires two HCl. The left now has H=2+2=4, O=2, requiring two H₂O on the right. Make Ca, H, O, Cl rows and total both sides. State symbols s, l, g, aq do not enter atom arithmetic.',
  'The result is Ca(OH)₂ + 2HCl → CaCl₂ + 2H₂O. Each side has Ca=1, H=4, O=2, Cl=2. OH was not retained as a ready-made block: its atoms entered water. Two before HCl multiplies both H and Cl. This is a school-level overall equation for supplied products, not a description of every solution particle.',
  '**Mistake:** multiplying only the final bracketed symbol, adding rather than multiplying a subscript, or changing water to H₂O₂ to obtain balance. States and reaction conditions are separate information. A precipitate is a solid phase formed; writing s does not prove precipitation at every dilution.',
  'A new scheme supplies Fe₂(SO₄)₃ + NaOH → Fe(OH)₃ + Na₂SO₄, with Fe(OH)₃ precipitation specified. Try unchanged SO₄ groups first, then Fe and Na; finally expand brackets and check O and H. Use paper only.',
  '**Check:** Fe₂(SO₄)₃ + 6NaOH → 2Fe(OH)₃ + 3Na₂SO₄. Both sides have Fe=2, S=3, Na=6, O=18, H=6. The equation is not an instruction to make a precipitate. Explain why SO₄ can be counted as a group here, whereas OH in the first example must be recounted when its atoms enter water.'
 )
},
{
 'stem': '37-mass-and-system-boundary',
 'ru': (
  'Почему масса на весах может измениться, хотя атомы сохраняются?',
  'Определим границу системы, учтём входящие и выходящие вещества и точность измерения.',
  'В одном учебном отчёте масса установки до и после равна 100,00 г. В другом открытая установка стала легче: 100,00 → 97,00 г. Ученик решил, что вторая запись опровергает сохранение массы. Сначала нарисуйте границу взвешиваемого объекта. Никаких сосудов с газообразованием не закрывайте: закрытая система здесь только условие задачи.',
  '**Система** — выбранная для анализа часть мира; **окружение** — остальное. **Закрытая система** не обменивается веществом с окружением, хотя может обмениваться энергией. В обычных химических задачах масса всех учтённых веществ сохраняется в пределах точности измерения. Для открытой системы нужны массы того, что вошло и вышло за границу.',
  'Предскажите, может ли газ иметь массу, если его не видно. Должен ли рост массы металла при окислении означать создание атомов из ничего? Не путайте массу самого куска металла с массой металла вместе со всем исходным воздухом. Выберите объект до начала вычисления.',
  'Даны учебные балансы: закрытая система 100,00 → 100,00 г; открытая 100,00 → 97,00 г, при этом отдельно учтено 3,00 г вышедшего газа; кусок металла 10,00 → 12,00 г после присоединения 2,00 г кислорода. Рассчитайте итог каждого баланса. Это готовые числа, не инструкция нагревать металл или собирать газ.',
  'Для открытого случая 97,00+3,00=100,00 г: масса не исчезла, часть оказалась вне взвешиваемой границы. Для металла 10,00+2,00=12,00 г: дополнительная масса пришла с кислородом. Для закрытой системы сравнивается вся установка с её содержимым, а не один выбранный продукт.',
  '**Ошибка:** объявлять газ невесомым либо объяснять любую разницу химическим исчезновением массы. Реальные данные имеют погрешность: 100,00 и 100,01 г при неопределённости ±0,02 г каждого измерения не позволяют утверждать значимое изменение. Разница также может быть вызвана утечкой, испарением, переносом капель или иной ошибкой учёта.',
  'Новая открытая система: сначала 64,00 г, после 61,50 г; вход вещества исключён. Сколько массы вышло, если другие потери отсутствуют? Можно ли по одной разнице узнать формулу газа? Отдельно разберите случай поступления 1,20 г кислорода в исходный образец 8,80 г.',
  '**Проверка:** вышло 2,50 г, но формула газа из этой массы не следует; второй образец станет 10,00 г при отсутствии других потоков. По памяти запишите «начальная масса + вошло = конечная масса + вышло». Сначала примените границу, затем оцените точность. Уравнивание атомов и хороший массовый баланс не заменяют определения состава неизвестных продуктов.'
 ),
 'kz': (
  'Атомдар сақталса да, таразыдағы масса неге өзгеруі мүмкін?',
  'Жүйе шекарасын анықтап, кіретін және шығатын заттарды, өлшеу дәлдігін ескереміз.',
  'Бір оқу есебінде қондырғы массасы бұрын да, кейін де 100,00 г. Екінші ашық қондырғы жеңілдеген: 100,00 → 97,00 г. Оқушы мұны массаның сақталуына қайшы деді. Алдымен өлшенетін нысан шекарасын сызыңыз. Газ түзілетін ыдысты жаппаңыз: жабық жүйе мұнда тек есеп шарты.',
  '**Жүйе** — талдауға таңдалған дүние бөлігі, **қоршаған орта** — қалғаны. **Жабық жүйе** ортамен зат алмаспайды, бірақ энергия алмаса алады. Кәдімгі химиялық есепте барлық ескерілген зат массасы өлшеу дәлдігі шегінде сақталады. Ашық жүйе үшін шекарадан кірген және шыққан зат массалары қажет.',
  'Көрінбейтін газдың массасы болуы мүмкін бе? Металл тотығып ауырласа, бұл атомдар жоқтан жасалды деген сөз бе? Металл кесегінің массасын металл мен бүкіл бастапқы ауаның массасымен шатастырмаңыз. Есептеуден бұрын нысанды таңдаңыз.',
  'Оқу баланстары: жабық жүйе 100,00 → 100,00 г; ашық жүйе 100,00 → 97,00 г, бөлек 3,00 г шыққан газ ескерілген; металл 10,00 → 12,00 г, оған 2,00 г оттек қосылған. Әр балансты есептеңіз. Бұлар дайын сандар, металл қыздыру не газ жинау нұсқауы емес.',
  'Ашық жағдай үшін 97,00+3,00=100,00 г: масса жоғалмай, бір бөлігі өлшенетін шекарадан шыққан. Металл үшін 10,00+2,00=12,00 г: қосымша масса оттекпен келген. Жабық жүйеде бір өнім емес, ішіндегісімен бірге бүкіл қондырғы салыстырылады.',
  '**Қате:** газды массасыз деу немесе кез келген айырманы массаның химиялық жоғалуы деп түсіндіру. Нақты деректің белгісіздігі бар: әр өлшеуде ±0,02 г болса, 100,00 және 100,01 г мәндері маңызды өзгерісті дәлелдемейді. Айырма ағу, булану, тамшы тасымалы не есепке алу қатесінен де болуы мүмкін.',
  'Жаңа ашық жүйе: бастапқы 64,00 г, соңғы 61,50 г; зат кіруі жоқ. Басқа шығын болмаса, қанша масса шыққан? Осы айырмадан газ формуласын білуге бола ма? Бастапқы 8,80 г үлгіге 1,20 г оттек кіруін бөлек қарастырыңыз.',
  '**Тексеру:** 2,50 г шыққан, бірақ газ формуласы бұл массадан анықталмайды; екінші үлгі басқа ағынсыз 10,00 г болады. «Бастапқы масса + кірген = соңғы масса + шыққан» деп жатқа жазыңыз. Алдымен шекараны, кейін дәлдікті қолданыңыз. Атом теңгерімі мен жақсы масса балансы белгісіз өнім құрамын анықтауды алмастырмайды.'
 ),
 'en': (
  'Why can a balance reading change when atoms are conserved?',
  'Define the system boundary and account for material entering, leaving and measurement uncertainty.',
  'One teaching report gives a whole apparatus mass of 100.00 g before and after. Another open apparatus changes from 100.00 to 97.00 g. A learner calls the second result a violation of conservation. First draw the weighing boundary. Never seal a gas-producing reaction: a closed system here is only a problem condition.',
  'A **system** is the part of the world selected for analysis; its **surroundings** are everything else. A **closed system** exchanges no matter with the surroundings but may exchange energy. In ordinary chemical exercises, the mass of all accounted material is conserved within measurement accuracy. An open-system balance must include material crossing its boundary.',
  'Predict whether an invisible gas has mass. Does a metal gaining mass during oxidation mean atoms appeared from nothing? Distinguish the metal piece from metal plus all the original air. Select the object before calculating.',
  'Given balances are: closed system 100.00 → 100.00 g; open system 100.00 → 97.00 g with 3.00 g of escaped gas separately accounted for; metal 10.00 → 12.00 g after 2.00 g oxygen joins it. Calculate each balance. These supplied numbers are not instructions to heat metal or collect gas.',
  'For the open case, 97.00+3.00=100.00 g: mass has crossed the weighing boundary rather than vanished. For metal, 10.00+2.00=12.00 g: oxygen supplies the added mass. The closed-system comparison includes the entire apparatus and contents rather than one selected product.',
  '**Mistake:** calling gas weightless or interpreting every difference as chemical destruction of mass. Measurements have uncertainty: 100.00 and 100.01 g, each with uncertainty ±0.02 g, do not establish a significant change. Leaks, evaporation, transferred droplets or accounting errors can also cause differences.',
  'An open system starts at 64.00 g and ends at 61.50 g with no incoming material. How much left if other losses are absent? Can that difference identify the gas formula? Separately consider 1.20 g oxygen entering an 8.80 g sample.',
  '**Check:** 2.50 g left, but its gas formula is not established; the second sample becomes 10.00 g without other flows. Recall “initial mass + incoming = final mass + outgoing”. Apply the boundary before assessing uncertainty. Atom balance and good mass accounting do not replace identification of unknown products.'
 )
},
{
 'stem': '38-reaction-energy-and-heat',
 'ru': (
  'Почему одни реакции нагревают окружение, а другие охлаждают?',
  'Сравним затраты на разрыв связей и выделение энергии при образовании новых; не спутаем тепло и температуру.',
  'В учебных данных две реакции проведены в сопоставимых условиях: вокруг первой среда нагрелась, вокруг второй охладилась. Обе сопровождаются перестройкой связей. Как связать это с энергией, не утверждая «любая разорванная связь выделяет тепло»? Рассматривайте таблицу, а не проводите реакции дома.',
  '**Экзотермический процесс** передаёт энергию в виде тепла окружению, **эндотермический** получает её из окружения. **Температура** характеризует тепловое состояние; **тепло** — передача энергии из-за разности температур. При разрыве связи энергия требуется, при образовании — выделяется. Для итогового эффекта сравнивают обе части, а не выбирают один удобный шаг.',
  'Предскажите итог, если на разрыв в модели требуется 100 единиц, а при образовании выделяется 140. Куда направлены оставшиеся 40? А если требуется 120 и выделяется 80? Все числа учебные, в условных единицах; они не являются измеренными энергиями конкретных веществ.',
  'Сделайте две строки: «потрачено», «выделилось», «разница», «направление передачи». Для первой модели вычислите 100−140=−40, для второй 120−80=+40. Минус здесь означает уменьшение энергии реагирующей системы и передачу окружению, плюс — увеличение энергии за счёт получения извне. Пока это упрощённый энергетический учёт.',
  'Первая модель экзотермическая: система отдаёт 40 условных единиц. Вторая эндотермическая: получает 40. В сопоставимых условиях наблюдение температуры окружения помогает выявить направление передачи, но одной разницы градусов недостаточно для вычисления энергии: нужны также массы, теплоёмкости и сведения о потерях.',
  '**Ошибка:** считать всё, что нагрелось, доказательством химической реакции или полагать, что экзотермическая реакция обязательно начинается сама. Нагрев возможен при физическом процессе; для начала некоторых реакций требуется преодолеть энергетический барьер. Остывание при растворении тоже не доказывает появление нового вещества. Нельзя проверять это смешением бытовой химии.',
  'Новая модель: на разрыв требуется 90 единиц, при образовании выделяется 70. Определите итог и направление обмена. Затем оцените фразу «если термометр показал рост на 5 градусов, выделилось 5 единиц энергии». Каких данных в ней нет? Не заменяйте температуру энергией.',
  '**Проверка:** 90−70=+20, процесс получает 20 условных единиц и является эндотермическим в данной модели. Пять градусов не равны пяти единицам энергии; не заданы масса, теплоёмкость и условия обмена. Восстановите цепочку «разрыв требует → образование выделяет → сравни обе суммы → укажи систему и окружение». Для целой реакции одной связи недостаточно.'
 ),
 'kz': (
  'Кей реакция ортаны неге қыздырады, ал басқасы салқындатады?',
  'Байланыс үзу шығынын жаңа байланыс түзу энергиясымен салыстырып, жылу мен температураны шатастырмаймыз.',
  'Оқу дерегінде салыстырмалы жағдайдағы екі реакция бар: біріншісінің айналасындағы орта жылынды, екіншісінде салқындады. Екеуінде де байланыс қайта құрылады. Мұны «әр үзілген байланыс жылу бөледі» демей қалай түсіндіреміз? Үйде реакция жасамай, кестені талдаңыз.',
  '**Экзотермиялық процесс** энергияны жылу түрінде ортаға береді, **эндотермиялық процесс** ортадан алады. **Температура** жылулық күйді сипаттайды, **жылу** — температура айырмасынан энергияның берілуі. Байланысты үзу энергияны қажет етеді, түзу кезінде энергия бөлінеді. Қорытынды әсер үшін бір қолайлы сатыны емес, екі бөлікті салыстырамыз.',
  'Модельде үзуге 100 бірлік қажет болып, түзілгенде 140 бөлінсе, нәтижені болжаңыз. Қалған 40 қайда кетеді?120 қажет болып, 80 бөлінсе ше? Сандардың бәрі шартты оқу бірліктері, нақты заттардың өлшенген энергиялары емес.',
  '«Жұмсалды», «бөлінді», «айырма», «берілу бағыты» бағандарын жасаңыз. Бірінші модельде 100−140=−40, екіншісінде 120−80=+40. Мұнда минус реагенттер жүйесінің энергиясы азайып, ортаға берілгенін, плюс сырттан алу арқылы өскенін білдіреді. Бұл әзірге жеңілдетілген энергетикалық есеп.',
  'Біріншісі экзотермиялық: жүйе 40 шартты бірлік береді. Екіншісі эндотермиялық: 40 алады. Салыстырмалы жағдайда орта температурасы берілу бағытын түсінуге көмектеседі, бірақ энергияны есептеуге градус айырмасы жеткіліксіз: масса, жылусыйымдылық және шығын дерегі керек.',
  '**Қате:** кез келген қызуды химиялық реакция дәлелі деу немесе экзотермиялық реакция міндетті түрде өзінен басталады деу. Физикалық процесс те қыздыруы мүмкін; кей реакция басталу үшін энергетикалық тосқауылды еңсеруі керек. Ерігендегі салқындау да жаңа затты дәлелдемейді. Мұны тұрмыстық химияны араластырып тексеруге болмайды.',
  'Жаңа модель: үзуге 90 бірлік қажет, түзілгенде 70 бөлінеді. Нәтиже мен алмасу бағытын табыңыз. Кейін «термометр 5 градус өсті, демек 5 энергия бірлігі бөлінді» дегенді бағалаңыз. Қандай дерек жоқ? Температураны энергиямен алмастырмаңыз.',
  '**Тексеру:** 90−70=+20, процесс 20 шартты бірлік алады, осы модельде эндотермиялық. Бес градус бес энергия бірлігі емес; масса, жылусыйымдылық және алмасу жағдайы берілмеген. «Үзу қажет етеді → түзілу бөледі → екі қосындыны салыстыр → жүйе мен ортаны ата» тізбегін қалпына келтіріңіз. Тұтас реакцияны бір байланыспен түсіндіру жеткіліксіз.'
 ),
 'en': (
  'Why do some reactions warm their surroundings while others cool them?',
  'Compare bond-breaking costs with energy released by bond formation, keeping heat distinct from temperature.',
  'Supplied observations describe two reactions under comparable conditions: the surroundings warm in one and cool in the other. Both rearrange bonds. How does energy explain this without claiming that every broken bond releases heat? Analyse the table rather than performing reactions at home.',
  'An **exothermic process** transfers energy as heat to its surroundings; an **endothermic process** takes it from them. **Temperature** describes thermal state; **heat** is energy transfer due to a temperature difference. Breaking bonds requires energy, while forming bonds releases it. Compare both contributions to determine the overall effect.',
  'Predict the result when bond breaking requires 100 units and formation releases 140. Where do the remaining 40 go? What if 120 are required and 80 released? All values are invented teaching units, not measured bond energies for particular substances.',
  'Make columns for cost, release, difference and transfer direction. The first model gives 100−140=−40; the second 120−80=+40. Minus denotes energy leaving the reacting system for the surroundings; plus denotes energy gained from outside. This remains simplified energy bookkeeping.',
  'The first model is exothermic, releasing 40 arbitrary units. The second is endothermic, taking in 40. Under comparable conditions, surrounding temperature helps indicate transfer direction, but a degree change alone cannot determine energy: masses, heat capacities and losses also matter.',
  '**Mistake:** treating all warming as proof of a chemical reaction, or assuming an exothermic reaction starts on its own. Physical processes can warm surroundings; some reactions require an initial energy barrier to be overcome. Cooling during dissolution also does not prove a new substance formed. Do not test such ideas by mixing household chemicals.',
  'A new model requires 90 units to break bonds and releases 70 on formation. Determine the net transfer and direction. Then evaluate “a thermometer rose by 5 degrees, so 5 energy units were released”. Which information is missing? Temperature must not be substituted for energy.',
  '**Check:** 90−70=+20: this model takes in 20 units and is endothermic. Five degrees are not five energy units; mass, heat capacity and exchange conditions are missing. Recall “breaking costs → formation releases → compare totals → name system and surroundings”. One bond cannot determine a whole reaction’s energy balance.'
 )
},
{
 'stem': '39-reaction-patterns-and-limits',
 'ru': (
  'Зачем различать типы реакций, если формулы уже известны?',
  'Распознаем соединение, разложение, замещение и обмен; не превратим схему в предсказание любой реакции.',
  'В дневнике несколько уравнений. В одном два исходных вещества дали один продукт, в другом одно вещество дало два продукта. Можно ли сгруппировать записи по строению? Это классификация данных уравнений, а не способ самостоятельно смешивать неизвестные вещества или гарантировать результат.',
  '**Соединение** в нашей школьной классификации: несколько реагентов образуют один продукт. **Разложение**: один реагент даёт несколько продуктов. **Замещение**: простое вещество реагирует со сложным, образуя другое простое и другое сложное. **Обмен**: два сложных вещества обмениваются составными частями. Здесь простое вещество содержит один элемент, сложное — два или больше.',
  'Предскажите тип для 2Mg + O₂ → 2MgO и CaCO₃ → CaO + CO₂. Означает ли коэффициент 2 перед MgO, что получилось два разных вещества? Считайте виды веществ, а не только числа перед формулами. Во втором случае условие превращения — нагревание, но опыт не предлагается.',
  'Разберите ещё две данные записи: Zn + CuSO₄ → ZnSO₄ + Cu; Na₂CO₃ + CaCl₂ → CaCO₃ + 2NaCl. Для последней по условию водных растворов образуется осадок CaCO₃. Подчеркните простые вещества, сравните элементы в сложных и проверьте баланс. Не берите соли металлов, кислоты или неизвестные порошки.',
  'Магний и кислород дают реакцию соединения; карбонат кальция — разложения; цинк с сульфатом меди — замещения; данные растворы карбоната натрия и хлорида кальция — обмена. Для обмена с каждой стороны Na=2, C=1, O=3, Ca=1, Cl=2. Все четыре записи уравнены; тип и баланс отвечают на разные вопросы.',
  '**Ошибка:** считать, что любая пара солей обязательно вступит в обмен с наблюдаемым превращением или любой металл вытеснит любой другой. Нужны свойства, состав, условия и движущая причина конкретного процесса. Классификации могут пересекаться: реакцию можно описывать по схеме и по переносу электронов. Последний критерий будет изучаться отдельно.',
  'Новые записи: 2CO + O₂ → 2CO₂ и Ca(OH)₂ + 2HCl → CaCl₂ + 2H₂O. Определите тип по принятой схеме. Можно ли назвать первую ещё и окислением? Должен ли один ярлык отменять другой? Реакции даны для чтения, не для выполнения.',
  '**Проверка:** первая — соединение по числу видов реагентов и продуктов; она также связана с окислением CO. Вторая — обмен в школьной суммарной записи. Разные основания классификации не противоречат друг другу. По памяти приведите пример каждого типа, укажите условия как отдельную информацию и завершите фразой «узнать тип — не значит доказать возможность реакции».'
 ),
 'kz': (
  'Формула белгілі болса, реакция түрін ажырату не үшін керек?',
  'Қосылу, айырылу, орынбасу және алмасуды танып, сұлбаны кез келген реакция болжамына айналдырмаймыз.',
  'Күнделікте бірнеше теңдеу бар. Бірінде екі бастапқы зат бір өнімге, екіншісінде бір зат екі өнімге айналған. Жазбаларды құрылымы бойынша топтастыруға бола ма? Бұл — берілген теңдеулерді жіктеу, белгісіз заттарды өз бетімен араластыру не нәтижеге кепілдік беру тәсілі емес.',
  'Мектептік жіктеуде **қосылу** — бірнеше реагенттен бір өнім түзілуі. **Айырылу** — бір реагенттен бірнеше өнім түзілуі. **Орынбасу** — жай зат күрделі затпен әрекеттесіп, басқа жай және күрделі зат түзуі. **Алмасу** — екі күрделі заттың құрамдас бөліктерін алмастыруы. Жай зат бір, күрделі зат екі не одан көп элементтен тұрады.',
  '2Mg + O₂ → 2MgO және CaCO₃ → CaO + CO₂ түрін болжаңыз. MgO алдындағы 2 екі түрлі зат пайда болды деген сөз бе? Формула алдындағы санды емес, зат түрлерін санаңыз. Екінші айналудың шарты — қыздыру, бірақ тәжірибе ұсынылмайды.',
  'Тағы екі дайын жазба: Zn + CuSO₄ → ZnSO₄ + Cu; Na₂CO₃ + CaCl₂ → CaCO₃ + 2NaCl. Соңғысында судағы ерітінділер үшін CaCO₃ тұнбасы түзіледі деп берілген. Жай заттарды сызып, күрделі зат элементтерін салыстырыңыз, теңгерімді тексеріңіз. Металл тұзын, қышқылды, белгісіз ұнтақты алмаңыз.',
  'Магний мен оттек — қосылу; кальций карбонаты — айырылу; мыстың сульфатымен мырыш — орынбасу; берілген натрий карбонаты мен кальций хлориді ерітінділері — алмасу. Алмасуда екі жақта Na=2, C=1, O=3, Ca=1, Cl=2. Төрт жазба да теңестірілген; түр мен теңгерім әртүрлі сұраққа жауап береді.',
  '**Қате:** кез келген екі тұз міндетті түрде көрінетін алмасуға түседі немесе кез келген металл басқасын ығыстырады деу. Нақты процестің қасиеті, құрамы, жағдайы, қозғаушы себебі керек. Жіктеулер қиылысуы мүмкін: реакцияны сұлба және электрон берілуі бойынша сипаттауға болады. Соңғы белгі кейін оқытылады.',
  'Жаңа жазбалар: 2CO + O₂ → 2CO₂ және Ca(OH)₂ + 2HCl → CaCl₂ + 2H₂O. Қабылданған сұлба бойынша түрін анықтаңыз. Біріншісін тотығу деп те атауға бола ма? Бір атау екіншісін жоққа шығаруы керек пе? Реакциялар орындауға емес, оқуға берілген.',
  '**Тексеру:** біріншісі реагент пен өнім түрлерінің саны бойынша қосылу; CO тотығуына да қатысты. Екіншісі мектептік жиынтық жазбада алмасу. Әртүрлі жіктеу негізі қайшылық емес. Әр түрге мысал келтіріп, жағдайды бөлек дерек деп көрсетіңіз; «түрін тану — реакция мүмкіндігін дәлелдеу емес» деп аяқтаңыз.'
 ),
 'en': (
  'Why classify reactions when their formulas are already known?',
  'Recognise combination, decomposition, displacement and exchange without treating patterns as universal predictions.',
  'Several notebook equations show different patterns: two starting substances form one product in one, while one substance forms two products in another. Can we group the equations by structure? This classifies supplied reactions rather than authorising mixtures of unknown substances or guaranteeing results.',
  'In this school classification, **combination** joins several reactants into one product; **decomposition** forms several products from one reactant. **Single displacement** has an elemental substance reacting with a compound to form a different elemental substance and compound. **Double displacement**, or exchange, swaps components between two compounds. An elemental substance contains one element; a compound contains two or more.',
  'Predict the types of 2Mg + O₂ → 2MgO and CaCO₃ → CaO + CO₂. Does two before MgO mean two different substances? Count substance identities, not coefficients. Heating is a condition of the second transformation, not an experiment being assigned.',
  'Analyse two further supplied equations: Zn + CuSO₄ → ZnSO₄ + Cu; Na₂CO₃ + CaCl₂ → CaCO₃ + 2NaCl. The latter specifies aqueous conditions producing a CaCO₃ precipitate. Mark elemental substances, compare compound compositions and check balance. Do not obtain metal salts, acids or unidentified powders.',
  'The magnesium reaction is combination; calcium carbonate decomposition; zinc with copper sulfate single displacement; and the specified carbonate/chloride solutions double displacement. For the last equation each side has Na=2, C=1, O=3, Ca=1, Cl=2. All four are balanced; type and balance answer different questions.',
  '**Mistake:** assuming every salt pair reacts observably or every metal displaces every other metal. Properties, conditions and a reason for a particular transformation matter. Classifications can overlap: a pattern and an electron-transfer description may apply to the same reaction. Electron transfer is studied separately later.',
  'Classify 2CO + O₂ → 2CO₂ and Ca(OH)₂ + 2HCl → CaCl₂ + 2H₂O by the stated patterns. Can the first also be called oxidation? Must one label cancel another? These equations are for reading, not for practical performance.',
  '**Check:** the first is combination by substance identities and also involves oxidation of CO; the second is exchange in the school-level overall notation. Different classification criteria need not conflict. Recall an example of each type, distinguish conditions from pattern and finish: “recognising a type does not prove a reaction is possible”.'
 )
},
{
 'stem': '40-reaction-case-review',
 'ru': (
  'Дело о соли: как проверить превращение и не потерять вещество?',
  'Свяжем признаки, атомный счёт, массу и энергию; отличим проверенный вывод от догадки о составе.',
  'В готовом отчёте вещество A известно как CaCO₃, раствор B содержит HCl. Анализом продуктов установлены CaCl₂, H₂O и CO₂. Даны две модели измерения: вся закрытая система 150,00 → 150,01 г с неопределённостью ±0,02 г каждого значения; открытая установка 150,00 → 149,56 г, и единственный выход вещества — CO₂. Это задача на бумаге: реальные газообразующие смеси не закрывайте.',
  'Повторим четыре разных проверки. **Признаки** помогают поставить вопрос о реакции. **Атомный баланс** проверяет согласованность уравнения с сохранением элементов. **Массовый баланс** требует выбранной границы и учёта потоков. **Энергетический баланс** сравнивает затраты и выделение энергии. Ни одна из этих проверок сама по себе не определяет формулу неизвестного порошка.',
  'До разбора оцените три вывода ученика: «газ улетел, значит масса уничтожена»; «разность 0,01 г доказывает появление массы»; «для уравнивания можно поменять H₂O на H₂O₂». Исправьте каждую причину. Затем предскажите коэффициент перед HCl в данной схеме CaCO₃ + HCl → CaCl₂ + H₂O + CO₂.',
  'Составьте дневник с колонками «дано», «расчёт», «вывод», «граница». Сначала посчитайте Ca, C, O, H, Cl. Затем отдельно вычислите массу вышедшего CO₂ и оцените отличие закрытых измерений. Формулы продуктов даны по результату анализа в условии: не пытайтесь получать их из одной разности масс.',
  'Уравнение: CaCO₃ + 2HCl → CaCl₂ + H₂O + CO₂. В обеих частях Ca=1, C=1, O=3, H=2, Cl=2. В открытом случае 150,00−149,56=0,44 г CO₂ покинуло выбранную границу. Разность 0,01 г в закрытой системе меньше заданной неопределённости и не устанавливает значимое изменение массы.',
  '**Ошибка:** объединить все вопросы в один ярлык «реакция доказана». В этой задаче вывод о продуктах основан на заданном анализе, а не только на пузырьках. Стрелка не показывает скорость, температуру и все стадии. Нельзя считать выделившееся тепло из массы CO₂ без дополнительных энергетических сведений.',
  'Итоговый перенос: уравняйте Na₂CO₃ + HCl → NaCl + H₂O + CO₂. Для другой открытой системы даны 80,00 → 79,20 г, единственный выход — газ. Найдите его массу, но не угадывайте формулу. В независимой энергетической модели затраты 50, выделение 65 условных единиц: определите направление итоговой передачи.',
  '**Проверка:** Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂; вышло 0,80 г, формула газа из этих измерений неизвестна; 50−65=−15, система отдаёт 15 условных единиц. Оцените отдельно пять умений: различать признаки и доказательства, читать уравнение, уравнивать, выбирать границу, объяснять энергию. Через сутки и неделю решите новые случаи. Следующий блок добавит количественный язык частиц и расчётов.'
 ),
 'kz': (
  'Тұздың сыры: айналуды тексеріп, затты қалай жоғалтпаймыз?',
  'Белгі, атом санағы, масса мен энергияны байланыстырып, тексерілген қорытындыны құрам туралы жорамалдан ажыратамыз.',
  'Дайын есепте A заты CaCO₃, B ерітіндісінде HCl бары белгілі. Өнім талдауы CaCl₂, H₂O, CO₂ анықтаған. Екі өлшеу моделі берілген: бүкіл жабық жүйе 150,00 → 150,01 г, әр мәннің белгісіздігі ±0,02 г; ашық қондырғы 150,00 → 149,56 г, жалғыз шығатын зат — CO₂. Бұл қағаз есебі: нақты газ түзетін қоспаны жаппаңыз.',
  'Төрт бөлек тексеруді қайталаймыз. **Белгілер** реакция туралы сұрақ қоюға көмектеседі. **Атом теңгерімі** теңдеудің элементтер сақталуымен сәйкестігін тексереді. **Масса балансы** таңдалған шекара мен ағындарды ескеруді қажет етеді. **Энергия балансы** шығын мен бөлінуді салыстырады. Ешқайсысы белгісіз ұнтақ формуласын жеке өзі анықтамайды.',
  'Талдаудан бұрын оқушының үш қорытындысын бағалаңыз: «газ ұшты, демек масса жойылды»; «0,01 г айырма масса пайда болғанын дәлелдейді»; «теңестіру үшін H₂O-ны H₂O₂-ге ауыстыруға болады». Әр қатенің себебін түзетіңіз. Кейін CaCO₃ + HCl → CaCl₂ + H₂O + CO₂ сұлбасындағы HCl коэффициентін болжаңыз.',
  '«Берілген», «есеп», «қорытынды», «шек» бағанды күнделік жасаңыз. Алдымен Ca, C, O, H, Cl санаңыз. Кейін шыққан CO₂ массасын бөлек есептеп, жабық өлшеулер айырмасын бағалаңыз. Өнім формулалары шарттағы талдаудан белгілі: оларды жалғыз масса айырмасынан шығармаңыз.',
  'Теңдеу: CaCO₃ + 2HCl → CaCl₂ + H₂O + CO₂. Екі жақта Ca=1, C=1, O=3, H=2, Cl=2. Ашық жағдайда 150,00−149,56=0,44 г CO₂ шекарадан шыққан. Жабық жүйедегі 0,01 г айырма берілген белгісіздіктен кіші, сондықтан маңызды масса өзгерісін анықтамайды.',
  '**Қате:** барлық сұрақты «реакция дәлелденді» деген бір атауға біріктіру. Өнім қорытындысы мұнда көпіршікке ғана емес, шарттағы талдауға сүйенеді. Көрсеткі жылдамдықты, температураны, барлық сатыны көрсетпейді. Қосымша энергетикалық дерексіз CO₂ массасынан бөлінген жылуды есептеуге болмайды.',
  'Қорытынды көшіру: Na₂CO₃ + HCl → NaCl + H₂O + CO₂ теңестіріңіз. Басқа ашық жүйеде 80,00 → 79,20 г, жалғыз шығатын зат — газ. Массасын табыңыз, формуласын болжай салмаңыз. Тәуелсіз энергия моделінде шығын 50, бөліну 65 шартты бірлік: қорытынды берілу бағытын анықтаңыз.',
  '**Тексеру:** Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂; 0,80 г шыққан, формула өлшеуден белгісіз; 50−65=−15, жүйе 15 шартты бірлік береді. Бес дағдыны бөлек бағалаңыз: белгі мен дәлелді ажырату, теңдеуді оқу, теңестіру, шекара таңдау, энергияны түсіндіру. Бір тәулік және аптадан кейін жаңа жағдайларды шешіңіз. Келесі блок бөлшектер мен есептеудің сандық тілін қосады.'
 ),
 'en': (
  'The salt case: how do we check a transformation without losing matter?',
  'Connect observations, atoms, mass and energy, separating supported conclusions from guesses about composition.',
  'A supplied report identifies A as CaCO₃ and solution B as containing HCl. Product analysis identifies CaCl₂, H₂O, CO₂. Two weighing models are given: the whole closed system 150.00 → 150.01 g, each measurement uncertain by ±0.02 g; an open apparatus 150.00 → 149.56 g, with CO₂ its only material outflow. This is a paper problem: never seal real gas-producing mixtures.',
  'Recall four distinct checks. **Observations** suggest questions about reaction. **Atom balance** checks consistency with element conservation. **Mass accounting** needs a boundary and material flows. **Energy accounting** compares energy costs and releases. No one check independently determines an unidentified powder’s formula.',
  'Evaluate three claims before the explanation: “escaping gas means mass was destroyed”; “a 0.01 g difference proves mass appeared”; “we may change H₂O into H₂O₂ to balance”. Correct each reason. Then predict the HCl coefficient in CaCO₃ + HCl → CaCl₂ + H₂O + CO₂.',
  'Make notebook columns for evidence, calculation, conclusion and limit. Count Ca, C, O, H, Cl first. Separately calculate the escaped CO₂ mass and assess the closed readings. Product formulas are supplied from analysis; do not attempt to derive them from one mass difference.',
  'The equation is CaCO₃ + 2HCl → CaCl₂ + H₂O + CO₂. Both sides have Ca=1, C=1, O=3, H=2, Cl=2. In the open case, 150.00−149.56=0.44 g CO₂ crossed the boundary. The closed-system 0.01 g difference is smaller than the stated uncertainty and does not establish a significant mass change.',
  '**Mistake:** combining all questions into “a reaction is proved”. Here product identification relies on the supplied analysis, not bubbles alone. The arrow does not give rate, temperature or every stage. Escaped CO₂ mass alone cannot calculate heat without additional energy information.',
  'Final transfer: balance Na₂CO₃ + HCl → NaCl + H₂O + CO₂. Another open system changes 80.00 → 79.20 g with gas as its only outflow. Find the mass but do not guess its formula. An independent energy model has costs 50 and release 65 arbitrary units: determine the net transfer direction.',
  '**Check:** Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂; 0.80 g left but its formula is unknown from weighing; 50−65=−15, so the system releases 15 units. Assess five skills separately: observation versus evidence, equation reading, balancing, boundaries and energy explanations. Retry new cases after a day and a week. The next block adds quantitative language for particles and calculations.'
 )
}
]

# Atom-balance data drive displayed reference tables and independent checks.
# Every product here is specified, not predicted by the balancing algorithm.
EQUATIONS = {
 'water': (((2,'H2'),(1,'O2')), ((2,'H2O'),)),
 'co': (((2,'CO'),(1,'O2')), ((2,'CO2'),)),
 'magnesium': (((2,'Mg'),(1,'O2')), ((2,'MgO'),)),
 'aluminium': (((4,'Al'),(3,'O2')), ((2,'Al2O3'),)),
 'iron': (((4,'Fe'),(3,'O2')), ((2,'Fe2O3'),)),
 'ammonia': (((1,'N2'),(3,'H2')), ((2,'NH3'),)),
 'neutralisation': (((1,'Ca(OH)2'),(2,'HCl')), ((1,'CaCl2'),(2,'H2O'))),
 'brackets': (((1,'Fe2(SO4)3'),(6,'NaOH')), ((2,'Fe(OH)3'),(3,'Na2SO4'))),
 'decomposition': (((1,'CaCO3'),), ((1,'CaO'),(1,'CO2'))),
 'displacement': (((1,'Zn'),(1,'CuSO4')), ((1,'ZnSO4'),(1,'Cu'))),
 'exchange': (((1,'Na2CO3'),(1,'CaCl2')), ((1,'CaCO3'),(2,'NaCl'))),
 'calcium_carbonate': (((1,'CaCO3'),(2,'HCl')), ((1,'CaCl2'),(1,'H2O'),(1,'CO2'))),
 'sodium_carbonate': (((1,'Na2CO3'),(2,'HCl')), ((2,'NaCl'),(1,'H2O'),(1,'CO2'))),
}
