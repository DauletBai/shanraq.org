"""Original trilingual lessons 41–49; calculations use supplied, rounded data.

No lesson instructs a learner to carry out a reaction. Each locale contains
title, summary and eight teaching sections. Entity type is always specified.
"""

LESSONS = [
{
 'stem': '41-mole-and-specified-entities',
 'ru': (
  'Как считать частицы, которые нельзя пересчитать по одной?',
  'Введём моль, постоянную Авогадро и обязательный вопрос: какие именно частицы считаем?',
  'В отчёте написано «взяли 1 моль». Этого мало: 1 моль молекул воды и 1 моль атомов водорода — разные наборы. Как связать огромный счёт частиц с короткой записью? Работайте с числами и бумажными карточками, без образцов и весов.',
  '**Количество вещества n** характеризует число указанных частиц. Его единица **моль (mol)** содержит ровно 6,02214076 × 10²³ указанных частиц. **Постоянная Авогадро Nₐ** равна ровно 6,02214076 × 10²³ mol⁻¹. **N** — число частиц, безразмерное число. Связь: N = nNₐ. Для приближённых расчётов здесь Nₐ ≈ 6,02 × 10²³ mol⁻¹.',
  'Сначала сравните 0,5 моль молекул H₂O и 0,5 моль молекул CO₂: где молекул больше? Затем сравните число всех атомов в этих наборах. Объясните, почему равное число молекул ещё не означает одинаковую массу или одинаковое число атомов каждого элемента.',
  'Постройте таблицу: 1 моль молекул H₂O содержит 1 моль молекул, 2 моль атомов H и 1 моль атомов O. Всего это 3 моль атомов. Для NaCl считайте формульные единицы состава кристалла, а не отдельные молекулы NaCl. При описании ионов называйте их знак и вид.',
  'Для 0,5 моль молекул H₂O: N ≈ 0,5 × 6,02 × 10²³ = 3,01 × 10²³ молекул. Атомов H вдвое больше: 6,02 × 10²³; атомов O столько же, сколько молекул. Эти оценки используют округлённую постоянную; определение моля остаётся точным.',
  '**Ошибка:** «моль — масса в граммах» или «в воде столько же атомов H, сколько молекул». Моль — единица количества вещества, грамм — массы. Индекс 2 в H₂O умножает число атомов H. Само число 6,02 × 10²³ не является массой частицы.',
  'Новый случай: дано 0,25 моль молекул O₂. Найдите количество вещества атомов O и приблизительное число молекул. Ещё один вопрос: в 1 моль формульных единиц CaCl₂ сколько моль ионов Cl⁻ соответствует составу? Реального растворения не требуется.',
  '**Проверка:** 0,50 моль атомов O; около 1,51 × 10²³ молекул O₂; 2 моль ионов Cl⁻ по составу CaCl₂. Восстановите маршрут «назови частицу → найди n → умножь на Nₐ → учти индексы». Проверяйте единицы и округляйте только результат.'
 ),
 'kz': (
  'Бір-бірлеп санауға келмейтін бөлшектерді қалай санаймыз?',
  'Мольді, Авогадро тұрақтысын және нақты қандай бөлшек саналатынын анықтаймыз.',
  'Есепте «1 моль алынды» деп жазылған. Бұл жеткіліксіз: 1 моль су молекуласы мен 1 моль сутек атомы — әртүрлі жиын. Бөлшектердің орасан көп санын қысқа жазбамен қалай береміз? Үлгі мен таразысыз, сандар және қағаз карточкалармен жұмыс істеңіз.',
  '**Зат мөлшері n** көрсетілген бөлшектер санын сипаттайды. Оның **моль (mol)** бірлігінде дәл 6,02214076 × 10²³ көрсетілген бөлшек бар. **Авогадро тұрақтысы Nₐ** дәл 6,02214076 × 10²³ mol⁻¹. **N** — өлшемсіз бөлшектер саны. Байланысы: N = nNₐ. Жуық есептерде Nₐ ≈ 6,02 × 10²³ mol⁻¹ қолданамыз.',
  'Алдымен 0,5 моль H₂O молекуласы мен 0,5 моль CO₂ молекуласын салыстырыңыз: қайсысында молекула көп? Содан кейін барлық атом санын салыстырыңыз. Молекула санының теңдігі неліктен массаның немесе әр элемент атомдары санының теңдігін білдірмейді?',
  'Кесте құрыңыз: 1 моль H₂O молекуласында 1 моль молекула, 2 моль H атомы және 1 моль O атомы бар. Барлығы 3 моль атом. NaCl үшін жеке молекулаларды емес, кристалл құрамының формулалық бірліктерін санаңыз. Иондарды сипаттағанда олардың түрі мен зарядын атаңыз.',
  '0,5 моль H₂O молекуласы үшін N ≈ 0,5 × 6,02 × 10²³ = 3,01 × 10²³ молекула. H атомдары екі есе көп: 6,02 × 10²³; O атомдарының саны молекула санымен бірдей. Бұл бағаларда тұрақты дөңгелектелген, ал мольдің анықтамасы дәл болып қалады.',
  '**Қате:** «моль — граммен берілген масса» немесе «судағы H атомы молекула санымен бірдей». Моль — зат мөлшерінің, грамм — массаның бірлігі. H₂O ішіндегі 2 индексі H атомдарының санын көбейтеді. 6,02 × 10²³ саны бөлшек массасы емес.',
  'Жаңа жағдай: 0,25 моль O₂ молекуласы берілген. O атомдарының зат мөлшерін және молекулалардың жуық санын табыңыз. 1 моль CaCl₂ формулалық бірлігінің құрамына неше моль Cl⁻ ионы сәйкес келеді? Затты шын еріту қажет емес.',
  '**Тексеру:** 0,50 моль O атомы; шамамен 1,51 × 10²³ O₂ молекуласы; CaCl₂ құрамы бойынша 2 моль Cl⁻ ионы. «Бөлшекті ата → n-ді тап → Nₐ-ға көбейт → индекстерді ескер» жолын қалпына келтіріңіз. Бірліктерді тексеріп, тек соңғы нәтижені дөңгелектеңіз.'
 ),
 'en': (
  'How can we count particles too numerous to count one by one?',
  'Introduce the mole and the Avogadro constant, always naming the entities being counted.',
  'A report says “we used 1 mole”. That is incomplete: 1 mole of water molecules and 1 mole of hydrogen atoms describe different sets. How can a short quantity describe an enormous count? Work with numbers and paper cards, without samples or balances.',
  '**Amount of substance n** measures the number of specified entities. Its unit, the **mole (mol)**, contains exactly 6.02214076 × 10²³ specified entities. The **Avogadro constant Nₐ** is exactly 6.02214076 × 10²³ mol⁻¹. **N** is the dimensionless entity count. The relation is N = nNₐ. Approximate calculations here use Nₐ ≈ 6.02 × 10²³ mol⁻¹.',
  'Compare 0.5 mol of H₂O molecules with 0.5 mol of CO₂ molecules: which has more molecules? Then compare the total atom counts. Explain why equal molecule counts need not mean equal masses or equal numbers of atoms of each element.',
  'Build a table: 1 mol of H₂O molecules contains 1 mol of molecules, 2 mol of H atoms and 1 mol of O atoms, or 3 mol of atoms in total. For NaCl, count formula units describing the crystal composition rather than separate NaCl molecules. When counting ions, name their species and charge.',
  'For 0.5 mol of H₂O molecules, N ≈ 0.5 × 6.02 × 10²³ = 3.01 × 10²³ molecules. There are twice as many H atoms, 6.02 × 10²³, and as many O atoms as molecules. These are estimates using a rounded constant; the definition of the mole remains exact.',
  '**Mistake:** “a mole is a mass in grams”, or “water has as many H atoms as molecules”. The mole measures amount; the gram measures mass. The subscript 2 in H₂O multiplies the H atom count. The number 6.02 × 10²³ is not a particle mass.',
  'Transfer: a set contains 0.25 mol of O₂ molecules. Find the amount of O atoms and the approximate molecule count. Also find how many moles of Cl⁻ ions correspond to the composition of 1 mol of CaCl₂ formula units. No actual dissolution is needed.',
  '**Check:** 0.50 mol of O atoms; approximately 1.51 × 10²³ O₂ molecules; 2 mol of Cl⁻ ions in the CaCl₂ composition. Reconstruct “name the entity → find n → multiply by Nₐ → account for subscripts”. Check units and round only the final result.'
 )
},
{
 'stem': '42-molar-mass-from-formula',
 'ru': (
  'Почему один моль разных веществ имеет разную массу?',
  'Рассчитаем молярную массу по формуле и отличим её от относительной формульной массы.',
  'В двух строках отчёта указано по 1 моль: H₂O и CO₂. Количество молекул одинаково, но массы 18 г и 44 г в данной округлённой модели. Откуда берётся разница? Вернитесь к составу каждой молекулы и массам её атомов.',
  '**Молярная масса M** — масса, делённая на количество вещества: M = m/n. Здесь единица g/mol, в СИ kg/mol. **Относительная формульная масса Mᵣ** — безразмерная сумма относительных атомных масс по формуле. Для школьной точности её числовое значение используем для M в g/mol; это не равенство величин с разными единицами.',
  'Предскажите, что тяжелее: 1 моль атомов O или 1 моль молекул O₂. Затем решите, надо ли учитывать индекс 2 за скобкой в Ca(OH)₂. Сначала объясните структуру формулы словами, а уже после складывайте числа.',
  'Даны округлённые значения Aᵣ: H=1, C=12, O=16, Na=23, Cl=35,5, Ca=40. Используйте именно их во всех задачах блока. Для Ca(OH)₂ запишите один Ca, два O и два H. Коэффициент перед формулой не входит в молярную массу вещества.',
  'M(H₂O) ≈ 2×1+16 = 18 g/mol; M(CO₂) ≈ 12+2×16 = 44 g/mol. Для Ca(OH)₂: 40+2×(16+1) = 74 g/mol. NaCl даёт 58,5 g/mol и относится к формульным единицам ионного вещества. Это расчётные приближения, а не точность лабораторного измерения.',
  '**Ошибка:** записать M(2H₂O)=36 g/mol как свойство воды. Две порции воды меняют массу и количество вещества одинаково, а отношение m/n сохраняется. Другая ошибка — дать Mᵣ единицу g/mol. Различайте символ величины, число и единицу.',
  'Рассчитайте M(CaCO₃), M(O₂) и M(H₂) по заданным массам. Затем определите, сколько граммов составляет 1 моль каждого вещества. Объясните, почему сама формула не сообщает массу произвольного образца без знания его количества.',
  '**Проверка:** около 100, 32 и 2 g/mol; 1 моль имеет соответственно 100, 32 и 2 г в принятой модели. Повторите подсчёт атомов до сложения. Для точных задач нужны согласованные справочные массы и сведения об изотопном составе; здесь используем учебное округление.'
 ),
 'kz': (
  'Неліктен әртүрлі заттың бір молінің массасы әртүрлі?',
  'Формуладан молярлық массаны есептеп, оны салыстырмалы формулалық массадан ажыратамыз.',
  'Есептің екі жолында да 1 моль берілген: H₂O және CO₂. Молекула саны бірдей, бірақ осы дөңгелектелген модельде массалары 18 г және 44 г. Айырма қайдан шығады? Әр молекуланың құрамы мен атомдарының массасына оралыңыз.',
  '**Молярлық масса M** — массаның зат мөлшеріне қатынасы: M = m/n. Мұнда бірлігі g/mol, СИ-де kg/mol. **Салыстырмалы формулалық масса Mᵣ** — формула бойынша салыстырмалы атомдық массалардың өлшемсіз қосындысы. Мектептік дәлдікте оның сандық мәнін M үшін g/mol бірлігімен қолданамыз; бірліктері бөлек шамалар тең емес.',
  'Болжаңыз: 1 моль O атомы ауыр ма, әлде 1 моль O₂ молекуласы ма? Ca(OH)₂ формуласындағы жақшадан кейінгі 2 индексін ескеру керек пе? Сандарды қоспай тұрып, формуланың құрамын сөзбен түсіндіріңіз.',
  'Дөңгелектелген Aᵣ мәндері берілген: H=1, C=12, O=16, Na=23, Cl=35,5, Ca=40. Блоктағы есептерде осы мәндерді қолданыңыз. Ca(OH)₂ үшін бір Ca, екі O және екі H жазыңыз. Формула алдындағы коэффициент заттың молярлық массасына кірмейді.',
  'M(H₂O) ≈ 2×1+16 = 18 g/mol; M(CO₂) ≈ 12+2×16 = 44 g/mol. Ca(OH)₂ үшін 40+2×(16+1) = 74 g/mol. NaCl үшін 58,5 g/mol шығады; мұнда иондық заттың формулалық бірлігі қаралады. Бұлар зертханалық өлшеудің дәлдігі емес, есептік жуықтаулар.',
  '**Қате:** M(2H₂O)=36 g/mol деп судың қасиеті ретінде жазу. Су порциясы екі еселенсе, масса мен зат мөлшері бірдей өсіп, m/n қатынасы сақталады. Тағы бір қате — Mᵣ-ге g/mol бірлігін беру. Шама белгісін, санды және бірлікті ажыратыңыз.',
  'Берілген массалармен M(CaCO₃), M(O₂), M(H₂) есептеңіз. Әр заттың 1 молі неше грамм екенін анықтаңыз. Формула неліктен зат мөлшері белгісіз кез келген үлгінің массасын өздігінен хабарламайтынын түсіндіріңіз.',
  '**Тексеру:** шамамен 100, 32 және 2 g/mol; берілген модельде 1 моль сәйкесінше 100, 32 және 2 г. Қоспай тұрып атомдарды санаңыз. Дәл есептерге келісілген анықтамалық массалар мен изотоптық құрам қажет; мұнда оқу үшін дөңгелектеу қолданылған.'
 ),
 'en': (
  'Why do equal mole amounts have different masses?',
  'Calculate molar mass from a formula and distinguish it from relative formula mass.',
  'Two report entries each specify 1 mol: H₂O and CO₂. Their molecule counts match, but their masses are 18 g and 44 g in the supplied rounded model. Where does the difference come from? Return to the composition and atomic masses in each molecule.',
  '**Molar mass M** is mass divided by amount: M = m/n. We use g/mol; its SI unit is kg/mol. **Relative formula mass Mᵣ** is a dimensionless sum of relative atomic masses according to a formula. At school precision, we use its numerical value for M in g/mol; quantities with different units are not identical.',
  'Predict whether 1 mol of O atoms or 1 mol of O₂ molecules is heavier. Then decide how to account for the subscript outside the brackets in Ca(OH)₂. Explain the composition in words before adding any masses.',
  'Use these supplied rounded Aᵣ values throughout this block: H=1, C=12, O=16, Na=23, Cl=35.5, Ca=40. For Ca(OH)₂, count one Ca, two O and two H. A coefficient before a formula does not enter the molar mass of that substance.',
  'M(H₂O) ≈ 2×1+16 = 18 g/mol; M(CO₂) ≈ 12+2×16 = 44 g/mol. For Ca(OH)₂, 40+2×(16+1) = 74 g/mol. NaCl gives 58.5 g/mol and describes formula units of an ionic substance. These are calculation approximations rather than claims about measurement precision.',
  '**Mistake:** writing M(2H₂O)=36 g/mol as a property of water. Doubling a portion doubles mass and amount equally, leaving m/n unchanged. Another mistake is assigning g/mol to Mᵣ. Separate the quantity symbol, numerical value and unit.',
  'Calculate M(CaCO₃), M(O₂) and M(H₂) using the supplied values. Find the mass of 1 mol of each. Explain why a formula alone cannot give the mass of an arbitrary sample without knowing its amount.',
  '**Check:** approximately 100, 32 and 2 g/mol; 1 mol has masses 100, 32 and 2 g respectively in this model. Count atoms before adding. Precise work needs consistent reference masses and isotope information; these exercises use school-level rounding.'
 )
},
{
 'stem': '43-mass-amount-particle-conversions',
 'ru': (
  'Как пройти от массы к числу частиц без путаницы?',
  'Свяжем m, n и N через молярную массу и проверим единицы на каждом шаге.',
  'В задаче даны 9,0 г чистой воды и M(H₂O)=18 g/mol. Требуется число молекул. Почему нельзя просто умножить 9 на постоянную Авогадро? Нарисуйте две стрелки: масса → количество вещества → число молекул.',
  'Между массой и количеством работает n = m/M, обратно m = nM. Между количеством и числом частиц: N = nNₐ и n = N/Nₐ. **Размерностная проверка** означает сокращение единиц: g ÷ (g/mol) = mol. Только после этого mol × mol⁻¹ даёт безразмерное число частиц.',
  'Предскажите: в 18 г воды молекул вдвое больше, чем в 9 г, или вдвое меньше? Если масса дана в mg, можно ли сразу делить её число на M в g/mol? Сначала приведите единицы к одной системе.',
  'Даны три строки: 9,0 г H₂O; 11 г CO₂ при M=44 g/mol; 500 mg H₂O. Переведите 500 mg в 0,500 г. Для каждого образца выпишите m, M, n, N и тип частицы. Используйте Nₐ ≈ 6,02 × 10²³ mol⁻¹.',
  'Для 9,0 г воды n = 9,0/18 = 0,50 mol; N ≈ 3,01 × 10²³ молекул. Для 11 г CO₂ n = 0,25 mol, N ≈ 1,51 × 10²³ молекул. Половина массы одного и того же чистого вещества соответствует половине числа его частиц при неизменном составе.',
  '**Ошибка:** делить массу воды на 6,02 × 10²³ и называть ответ молями. Постоянная переводит число частиц в количество вещества, а масса требует M. Проверяйте и чистоту: общая масса смеси не равна массе выбранного компонента.',
  'Найдите массу 0,20 mol CO₂ и количество вещества 1,204 × 10²³ молекул H₂O. Для 500 mg H₂O оцените число молекул. До вычисления укажите, какой результат должен быть меньше результата для 9,0 г воды.',
  '**Проверка:** m(CO₂)=8,8 г; n(H₂O)=0,20 mol при принятом округлении Nₐ. Для 0,500 г воды n≈0,0278 mol, N≈1,67 × 10²² молекул. Вернитесь к двум стрелкам и объясните, какой коэффициент относится к каждой.'
 ),
 'kz': (
  'Массадан бөлшек санына шатаспай қалай өтеміз?',
  'm, n және N шамаларын молярлық масса арқылы байланыстырып, әр қадамда бірліктерді тексереміз.',
  'Есепте 9,0 г таза су және M(H₂O)=18 g/mol берілген. Молекула санын табу керек. Неліктен 9 санын Авогадро тұрақтысына бірден көбейтуге болмайды? Екі бағыттауыш салыңыз: масса → зат мөлшері → молекула саны.',
  'Масса мен зат мөлшерін n = m/M, кері бағытта m = nM байланыстырады. Зат мөлшері мен бөлшек саны үшін N = nNₐ, n = N/Nₐ. **Өлшемді тексеру** — бірліктерді қысқарту: g ÷ (g/mol) = mol. Содан кейін ғана mol × mol⁻¹ өлшемсіз бөлшек санын береді.',
  'Болжаңыз: 18 г судағы молекула саны 9 г-дағыдан екі есе көп пе, аз ба? Масса mg арқылы берілсе, оны M-нің g/mol сандық мәніне бірден бөлуге бола ма? Алдымен бірліктерді сәйкестендіріңіз.',
  'Үш жол берілген: 9,0 г H₂O; M=44 g/mol болғандағы 11 г CO₂; 500 mg H₂O. 500 mg-ды 0,500 г-ға айналдырыңыз. Әр үлгі үшін m, M, n, N және бөлшек түрін жазыңыз. Nₐ ≈ 6,02 × 10²³ mol⁻¹ қолданыңыз.',
  '9,0 г су үшін n = 9,0/18 = 0,50 mol; N ≈ 3,01 × 10²³ молекула. 11 г CO₂ үшін n = 0,25 mol, N ≈ 1,51 × 10²³ молекула. Бір таза заттың құрамы өзгермесе, массасының жартысына бөлшек санының да жартысы сәйкес келеді.',
  '**Қате:** су массасын 6,02 × 10²³ санына бөліп, нәтижені моль деу. Тұрақты бөлшек санын зат мөлшеріне айналдырады, ал массаға M керек. Тазалықты да ескеріңіз: қоспаның толық массасы таңдалған компонент массасына тең емес.',
  '0,20 mol CO₂ массасын және 1,204 × 10²³ H₂O молекуласының зат мөлшерін табыңыз. 500 mg H₂O молекула санын бағалаңыз. Есептеуге дейін қай нәтиже 9,0 г су үшін нәтижеден кіші болатынын айтыңыз.',
  '**Тексеру:** m(CO₂)=8,8 г; қабылданған Nₐ жуықтауымен n(H₂O)=0,20 mol. 0,500 г су үшін n≈0,0278 mol, N≈1,67 × 10²² молекула. Екі бағыттауышқа оралып, әрқайсысында қандай көбейткіш қолданылатынын түсіндіріңіз.'
 ),
 'en': (
  'How do we move from mass to particle count without confusion?',
  'Connect m, n and N through molar mass, checking units at every step.',
  'A problem supplies 9.0 g of pure water and M(H₂O)=18 g/mol and asks for the molecule count. Why can we not simply multiply 9 by the Avogadro constant? Draw two arrows: mass → amount → number of molecules.',
  'Mass and amount are connected by n = m/M, or m = nM in reverse. Amount and entity count use N = nNₐ and n = N/Nₐ. A **dimensional check** cancels units: g ÷ (g/mol) = mol. Only then does mol × mol⁻¹ give a dimensionless entity count.',
  'Predict whether 18 g of water contains twice as many or half as many molecules as 9 g. If a mass is given in mg, can its number be divided directly by M in g/mol? Convert the units to a consistent system first.',
  'Use three entries: 9.0 g of H₂O; 11 g of CO₂ with M=44 g/mol; and 500 mg of H₂O. Convert 500 mg to 0.500 g. List m, M, n, N and the entity type for each. Use Nₐ ≈ 6.02 × 10²³ mol⁻¹.',
  'For 9.0 g of water, n = 9.0/18 = 0.50 mol and N ≈ 3.01 × 10²³ molecules. For 11 g of CO₂, n = 0.25 mol and N ≈ 1.51 × 10²³ molecules. Halving the mass of the same pure substance halves its particle count, provided its composition stays unchanged.',
  '**Mistake:** dividing the mass of water by 6.02 × 10²³ and calling the answer moles. The constant converts entity counts to amounts; mass requires M. Also check purity: the total mass of a mixture is not the mass of the selected component.',
  'Find the mass of 0.20 mol of CO₂ and the amount of 1.204 × 10²³ H₂O molecules. Estimate the molecule count in 500 mg of water. Before calculating, say which result should be smaller than the count for 9.0 g of water.',
  '**Check:** m(CO₂)=8.8 g; n(H₂O)=0.20 mol using the chosen rounded Nₐ. For 0.500 g of water, n≈0.0278 mol and N≈1.67 × 10²² molecules. Return to the two arrows and explain which conversion factor belongs to each.'
 )
},
{
 'stem': '44-element-mass-fractions',
 'ru': (
  'Почему доля атомов и доля массы — разные ответы?',
  'Найдём массовую долю элемента в известном чистом соединении и проверим сумму долей.',
  'В молекуле воды два атома H и один O. Ученик решил, что водород составляет две трети массы. Верно ли это? Сравните число атомов с их вкладами в массу, используя H=1 и O=16 в принятой учебной модели.',
  '**Массовая доля элемента w** — масса этого элемента в соединении, делённая на массу соединения. Для известной формулы w(E) = kAᵣ(E)/Mᵣ, где k — число атомов E в формуле. Доля безразмерна; для процентов умножьте её на 100%. Это состав чистого соединения, а не концентрация раствора.',
  'До расчёта предскажите, какой элемент даёт большую часть массы H₂O. Если взять 18 г или 36 г чистой воды с тем же составом, изменится ли массовая доля? Объясните разницу между неизменной долей и растущей массой элемента.',
  'Для воды составьте вклады: H даёт 2×1=2, O даёт 16, сумма 18. Для CO₂: C даёт 12, O даёт 32, сумма 44. Сначала найдите безразмерные доли, затем проценты. При округлении сумма процентов может слегка отличаться от 100%.',
  'В воде w(H)=2/18≈0,111, то есть 11,1%; w(O)=16/18≈88,9%. В 36 г воды масса водорода по составу равна 36×2/18=4 г. Это расчёт состава, а не инструкция выделять водород из воды.',
  '**Ошибка:** подставить число атомов вместо их массовых вкладов либо смешать 11,1% и число 11,1. Для умножения на массу используйте 0,111… или делите процент на 100. Масса кислорода в воде не означает, что внутри неё есть свободный газ O₂.',
  'Найдите массовые доли C и O в CO₂. Затем вычислите массу C, соответствующую 22 г чистого CO₂. Дополнительный случай: для CaCO₃ с Mᵣ≈100 найдите долю Ca при Aᵣ(Ca)=40.',
  '**Проверка:** CO₂ содержит около 27,3% C и 72,7% O; в 22 г CO₂ по составу 6 г C. В CaCO₃ доля Ca около 40%. Проверьте сумму долей и подпишите, идёт ли речь о составе соединения или о массе всего образца.'
 ),
 'kz': (
  'Атомдар үлесі мен масса үлесі неге бөлек жауап береді?',
  'Белгілі таза қосылыстағы элементтің массалық үлесін тауып, үлестер қосындысын тексереміз.',
  'Су молекуласында екі H атомы және бір O атомы бар. Оқушы сутек массаның үштен екісін құрайды деп шешті. Бұл дұрыс па? Қабылданған оқу моделінде H=1, O=16 мәндерін қолданып, атом саны мен массаға қосатын үлесін салыстырыңыз.',
  '**Элементтің массалық үлесі w** — қосылыстағы сол элемент массасының қосылыс массасына қатынасы. Белгілі формула үшін w(E) = kAᵣ(E)/Mᵣ, мұндағы k — формуладағы E атомдарының саны. Үлес өлшемсіз; пайыз алу үшін 100%-ға көбейтіңіз. Бұл ерітінді концентрациясы емес, таза қосылыс құрамы.',
  'Есептеуге дейін H₂O массасының көп бөлігін қай элемент беретінін болжаңыз. Құрамы бірдей 18 г немесе 36 г таза су алсақ, массалық үлес өзгере ме? Тұрақты үлес пен өсетін элемент массасының айырмасын түсіндіріңіз.',
  'Су үшін үлестерді құрыңыз: H массалық үлесі үшін 2×1=2, O үшін 16, қосындысы 18. CO₂ үшін C — 12, O — 32, қосындысы 44. Алдымен өлшемсіз үлестерді, кейін пайыздарды табыңыз. Дөңгелектеуден пайыздар қосындысы 100%-дан сәл ауытқуы мүмкін.',
  'Суда w(H)=2/18≈0,111, яғни 11,1%; w(O)=16/18≈88,9%. 36 г су құрамындағы сутек массасы 36×2/18=4 г. Бұл судан сутек бөліп алу нұсқауы емес, құрам есебі.',
  '**Қате:** массалық үлестер орнына атом санын қою немесе 11,1% пен 11,1 санын шатастыру. Массаға көбейту үшін 0,111… қолданыңыз не пайызды 100-ге бөліңіз. Судағы оттек массасы оның ішінде бос O₂ газы бар дегенді білдірмейді.',
  'CO₂ ішіндегі C пен O массалық үлестерін табыңыз. Содан кейін 22 г таза CO₂ құрамына сәйкес C массасын есептеңіз. Қосымша жағдай: Mᵣ≈100 болатын CaCO₃ үшін Aᵣ(Ca)=40 болса, Ca үлесін табыңыз.',
  '**Тексеру:** CO₂ құрамында шамамен 27,3% C және 72,7% O; 22 г CO₂ құрамына 6 г C сәйкес келеді. CaCO₃ ішіндегі Ca үлесі шамамен 40%. Үлестер қосындысын тексеріп, қосылыс құрамы ма, бүкіл үлгі массасы ма — белгілеңіз.'
 ),
 'en': (
  'Why are atom fractions and mass fractions different answers?',
  'Find an element’s mass fraction in a known pure compound and check that the fractions add up.',
  'A water molecule contains two H atoms and one O atom. A learner concludes that hydrogen contributes two thirds of the mass. Is that correct? Compare atom counts with mass contributions using H=1 and O=16 in the supplied teaching model.',
  'An **element’s mass fraction w** is its mass in a compound divided by the mass of the compound. For a known formula, w(E) = kAᵣ(E)/Mᵣ, where k counts E atoms in the formula. The fraction is dimensionless; multiply by 100% for a percentage. This describes a pure compound’s composition, not solution concentration.',
  'Predict which element contributes most of the mass of H₂O. Would the mass fraction change between 18 g and 36 g of pure water with the same composition? Explain the difference between a constant fraction and an increasing mass of an element.',
  'For water, list the contributions: H contributes 2×1=2, O contributes 16, total 18. For CO₂, C contributes 12 and O contributes 32, total 44. Find fractions first and percentages second. Rounding may make the percentages add to slightly more or less than 100%.',
  'In water, w(H)=2/18≈0.111, or 11.1%, and w(O)=16/18≈88.9%. In 36 g of water, the hydrogen mass represented by its composition is 36×2/18=4 g. This is a composition calculation, not an instruction to extract hydrogen from water.',
  '**Mistake:** using atom counts instead of mass contributions, or confusing 11.1% with 11.1. To multiply by a mass, use 0.111… or divide the percentage by 100. Oxygen mass in water does not mean free O₂ gas is present inside it.',
  'Find the mass fractions of C and O in CO₂. Then find the carbon mass corresponding to 22 g of pure CO₂. As an additional case, find the Ca fraction in CaCO₃ with Mᵣ≈100 and Aᵣ(Ca)=40.',
  '**Check:** CO₂ contains approximately 27.3% C and 72.7% O by mass; 22 g corresponds to 6 g of C in its composition. The Ca fraction in CaCO₃ is about 40%. Check the fraction sum and label whether you describe composition or the whole sample mass.'
 )
},
{
 'stem': '45-equation-mole-ratios',
 'ru': (
  'Как коэффициенты уравнения связывают количества веществ?',
  'Прочитаем коэффициенты как отношения молей, сохранив различие между молями и граммами.',
  'Дана уже проверенная реакция 2H₂ + O₂ → 2H₂O. Ученик считает, что 2 г водорода требуют 1 г кислорода. Что именно сравнивают числа 2:1:2? Это бумажная задача; газы не получайте, не смешивайте и не поджигайте.',
  '**Стехиометрическое отношение** — соотношение количеств веществ по коэффициентам уравнения. Для aA → bB: n(B)/n(A)=b/a для превращённых A и образовавшихся B в данной реакции. Индексы задают состав, коэффициенты — отношение количеств. Расчёт предполагает, что протекает именно указанная реакция.',
  'Предскажите расход O₂ для 4 mol H₂ при достаточном кислороде и полном превращении водорода. Потребуется ли 4 mol O₂? Затем решите, изменится ли отношение, если все коэффициенты умножить на два. Сократите его до простого вида.',
  'Заполните строки для превращения 2, 4 и 0,5 mol H₂: сколько O₂ расходуется и сколько H₂O образуется? Отделите «дано в начале» от «израсходовано»: исходный избыток кислорода не обязан целиком участвовать в реакции.',
  'Для 4 mol H₂ нужно n(O₂)=4×1/2=2 mol; образуется n(H₂O)=4×2/2=4 mol при полном превращении. В принятой модели массы составят 8 г H₂ и 64 г O₂, дающие 72 г H₂O. Отношение масс 1:8:9, а не 2:1:2.',
  '**Ошибка:** переносить коэффициенты напрямую на граммы или использовать неуравненную запись. Сначала проверьте атомный баланс, затем переведите массу в моли. Коэффициенты не доказывают, что реакция пойдёт сама или закончится полностью; это отдельные условия задачи.',
  'Дана 2CO + O₂ → 2CO₂. Если превратилось 0,60 mol CO, найдите расход O₂ и образование CO₂. Затем для CaCO₃ → CaO + CO₂ найдите количество CO₂ из 0,20 mol CaCO₃ при полном превращении. Это только расчёты, без опытов.',
  '**Проверка:** расход O₂ 0,30 mol; CO₂ образуется 0,60 mol. В случае CaCO₃ образуется 0,20 mol CO₂. Запишите для каждого перехода дробь из коэффициентов и подпишите, какое вещество дано, какое ищется, что израсходовано и что остаётся неизвестным.'
 ),
 'kz': (
  'Теңдеу коэффициенттері зат мөлшерлерін қалай байланыстырады?',
  'Коэффициенттерді моль қатынасы ретінде оқып, моль мен грамның айырмасын сақтаймыз.',
  'Тексерілген 2H₂ + O₂ → 2H₂O реакциясы берілген. Оқушы 2 г сутекке 1 г оттек қажет деп ойлайды. 2:1:2 сандары нақты нені салыстырады? Бұл — қағаздағы есеп; газдарды алмаңыз, араластырмаңыз және жақпаңыз.',
  '**Стехиометриялық қатынас** — теңдеу коэффициенттері бойынша зат мөлшерлерінің қатынасы. aA → bB үшін осы реакцияда айналған A мен түзілген B арасында n(B)/n(A)=b/a. Индекстер құрамды, коэффициенттер зат мөлшерінің қатынасын береді. Есепте дәл көрсетілген реакция жүреді деп қабылданады.',
  'Оттек жеткілікті әрі сутек толық айналса, 4 mol H₂ үшін O₂ шығынын болжаңыз. 4 mol O₂ керек пе? Барлық коэффициентті екі есе арттырсақ, қатынас өзгере ме? Оны ең қарапайым түрге қысқартыңыз.',
  '2, 4 және 0,5 mol H₂ айналуы үшін жолдарды толтырыңыз: қанша O₂ жұмсалып, қанша H₂O түзіледі? «Бастапқыда берілген» мен «жұмсалған» шамаларды ажыратыңыз: бастапқы артық оттек түгел реакцияға қатысуға міндетті емес.',
  '4 mol H₂ үшін n(O₂)=4×1/2=2 mol қажет; толық айналғанда n(H₂O)=4×2/2=4 mol түзіледі. Қабылданған модельде 8 г H₂ мен 64 г O₂ 72 г H₂O береді. Масса қатынасы 2:1:2 емес, 1:8:9.',
  '**Қате:** коэффициенттерді бірден грамға қолдану немесе теңестірілмеген жазбаны пайдалану. Алдымен атом теңгерімін тексеріп, кейін массаны мольге айналдырыңыз. Коэффициенттер реакцияның өздігінен жүруін не толық аяқталуын дәлелдемейді; бұлар — есептің бөлек шарттары.',
  '2CO + O₂ → 2CO₂ берілген. 0,60 mol CO айналса, O₂ шығыны мен CO₂ түзілуін табыңыз. CaCO₃ → CaO + CO₂ үшін 0,20 mol CaCO₃ толық айналғанда CO₂ мөлшерін табыңыз. Бұлар тәжірибесіз орындалатын есептер.',
  '**Тексеру:** O₂ шығыны 0,30 mol; CO₂ мөлшері 0,60 mol. CaCO₃ жағдайында 0,20 mol CO₂ түзіледі. Әр өтуге коэффициенттер бөлшегін жазып, берілген және ізделетін затты, жұмсалған және белгісіз қалған шамаларды белгілеңіз.'
 ),
 'en': (
  'How do equation coefficients connect amounts of substances?',
  'Read coefficients as mole ratios while keeping moles distinct from grams.',
  'The verified reaction is 2H₂ + O₂ → 2H₂O. A learner thinks 2 g of hydrogen requires 1 g of oxygen. What do the numbers 2:1:2 actually compare? This is a paper problem: do not generate, mix or ignite gases.',
  'A **stoichiometric ratio** relates amounts using balanced equation coefficients. For aA → bB, n(B)/n(A)=b/a for A consumed and B formed in that reaction. Subscripts specify composition; coefficients specify amount ratios. The calculation assumes that the stated reaction takes place.',
  'Predict the O₂ consumed by 4 mol of H₂ when oxygen is sufficient and hydrogen reacts completely. Would 4 mol of O₂ be needed? Decide whether doubling every coefficient changes the ratio, and reduce it to its simplest form.',
  'Complete rows for 2, 4 and 0.5 mol of H₂ consumed: how much O₂ is consumed and H₂O formed? Distinguish initial quantities from consumed quantities. An initial oxygen excess does not have to participate completely.',
  'For 4 mol of H₂, n(O₂)=4×1/2=2 mol is required and n(H₂O)=4×2/2=4 mol forms with complete conversion. In the supplied mass model, 8 g of H₂ and 64 g of O₂ give 72 g of H₂O. The mass ratio is 1:8:9, not 2:1:2.',
  '**Mistake:** applying coefficients directly to grams or using an unbalanced equation. Check atom balance first and convert mass to moles. Coefficients do not establish that a reaction proceeds spontaneously or goes to completion; those are separate conditions.',
  'For 2CO + O₂ → 2CO₂, find O₂ consumption and CO₂ formation when 0.60 mol of CO reacts. Then find the CO₂ amount from complete conversion of 0.20 mol of CaCO₃ in CaCO₃ → CaO + CO₂. These are calculations only, without experiments.',
  '**Check:** 0.30 mol of O₂ is consumed and 0.60 mol of CO₂ forms in the first case. The carbonate case gives 0.20 mol of CO₂. Write each coefficient fraction and label the given substance, target substance, consumed quantities and what remains unknown.'
 )
},
{
 'stem': '46-stoichiometric-mass-calculations',
 'ru': (
  'Как рассчитать массу продукта по уравнению реакции?',
  'Пройдём маршрут масса → моли → отношение коэффициентов → моли → масса.',
  'В расчётной модели 10,0 г чистого CaCO₃ полностью превращаются по уравнению CaCO₃ → CaO + CO₂. Найдите массу CO₂. Условие описывает уже заданное превращение; нагревать карбонаты и получать газы дома не требуется.',
  '**Теоретическая масса продукта** — расчёт по уравнению при заданном составе и полном превращении ограничивающего реагента. Сначала n(исходного)=m/M, затем n(продукта)=n(исходного)×b/a, в конце m(продукта)=nM. **Допущение** — условие модели, которое нужно назвать: чистота, выбранная реакция и степень превращения.',
  'Предскажите: из 10 г CaCO₃ может получиться 44 г CO₂ только потому, что M(CO₂)=44 g/mol? Сравните массу части продукта с массой исходного вещества для данной реакции. Затем решите, какая часть массы остаётся в CaO.',
  'Даны M(CaCO₃)=100, M(CaO)=56, M(CO₂)=44 g/mol. Коэффициенты 1:1:1. Запишите единицы во всех трёх переходах и отдельно проверьте баланс 100=56+44 для одного моля. Из этих данных нельзя вычислить время реакции или расход энергии.',
  'n(CaCO₃)=10,0/100=0,100 mol. Значит, n(CO₂)=0,100 mol и m(CO₂)=0,100×44=4,40 г. Остаток CaO имеет массу 0,100×56=5,60 г. Масса продуктов 4,40+5,60=10,00 г соответствует сохранению вещества в модели.',
  '**Ошибка:** остановиться на 0,100 и подписать «г», хотя получены моли. Другая ошибка — считать выход по всей массе загрязнённого образца как по чистому CaCO₃. Для смеси сначала нужна масса реагирующего вещества, а остальные компоненты учитываются отдельно.',
  'Новый расчёт: 5,30 г чистого Na₂CO₃ полностью реагируют с достаточным количеством HCl: Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂. Даны M(Na₂CO₃)=106, M(NaCl)=58,5 g/mol. Найдите теоретическую массу NaCl. Газ собирать не нужно.',
  '**Проверка:** n(Na₂CO₃)=0,0500 mol, n(NaCl)=0,100 mol, m(NaCl)=5,85 г. Масса соли больше массы карбоната, потому что часть её состава приходит из HCl. При проверке общей массы включайте все реагенты, а не только указанный первым.'
 ),
 'kz': (
  'Реакция теңдеуінен өнім массасын қалай есептейміз?',
  'Масса → моль → коэффициенттер қатынасы → моль → масса жолынан өтеміз.',
  'Есептік модельде 10,0 г таза CaCO₃ CaCO₃ → CaO + CO₂ теңдеуі бойынша толық айналады. CO₂ массасын табыңыз. Шартта айналу берілген; үйде карбонаттарды қыздыру не газ алу қажет емес.',
  '**Өнімнің теориялық массасы** — берілген құрам және шектеуші реагенттің толық айналуы жағдайында теңдеуден есептелетін масса. Алдымен n(бастапқы)=m/M, кейін n(өнім)=n(бастапқы)×b/a, соңында m(өнім)=nM. **Жорамал** — модельдің аталуы тиіс шарты: тазалық, таңдалған реакция және айналу дәрежесі.',
  'Болжаңыз: M(CO₂)=44 g/mol болғандықтан ғана 10 г CaCO₃-тен 44 г CO₂ шығуы мүмкін бе? Осы реакциядағы өнім бөлігінің массасын бастапқы зат массасымен салыстырыңыз. Массаның қай бөлігі CaO құрамында қалатынын анықтаңыз.',
  'M(CaCO₃)=100, M(CaO)=56, M(CO₂)=44 g/mol берілген. Коэффициенттер 1:1:1. Үш өтудің әрқайсысында бірліктерді жазыңыз; бір моль үшін 100=56+44 теңгерімін бөлек тексеріңіз. Бұл деректерден реакция уақытын не энергия шығынын табуға болмайды.',
  'n(CaCO₃)=10,0/100=0,100 mol. Демек, n(CO₂)=0,100 mol, m(CO₂)=0,100×44=4,40 г. Қалған CaO массасы 0,100×56=5,60 г. Өнімдер массасы 4,40+5,60=10,00 г модельдегі зат сақталуына сәйкес келеді.',
  '**Қате:** 0,100 нәтижесіне тоқтап, мольдің орнына «г» жазу. Тағы бір қате — қоспалы үлгінің бүкіл массасын таза CaCO₃ ретінде қолдану. Қоспа үшін алдымен әрекеттесетін зат массасы қажет; басқа компоненттер бөлек есепке алынады.',
  'Жаңа есеп: 5,30 г таза Na₂CO₃ жеткілікті HCl-мен толық әрекеттеседі: Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂. M(Na₂CO₃)=106, M(NaCl)=58,5 g/mol берілген. NaCl теориялық массасын табыңыз. Газ жинау қажет емес.',
  '**Тексеру:** n(Na₂CO₃)=0,0500 mol, n(NaCl)=0,100 mol, m(NaCl)=5,85 г. Тұз массасы карбонат массасынан көп, себебі оның құрамының бір бөлігі HCl-ден келеді. Жалпы массаны тексергенде тек бірінші затты емес, барлық реагентті қосыңыз.'
 ),
 'en': (
  'How can we calculate product mass from a reaction equation?',
  'Follow mass → moles → coefficient ratio → moles → mass.',
  'A calculation model states that 10.0 g of pure CaCO₃ converts completely by CaCO₃ → CaO + CO₂. Find the CO₂ mass. The conversion is supplied as a condition; no carbonate heating or gas generation at home is required.',
  'A **theoretical product mass** is calculated from an equation for the given composition and complete conversion of the limiting reactant. First n(start)=m/M, then n(product)=n(start)×b/a, finally m(product)=nM. An **assumption** is a model condition to state explicitly: purity, selected reaction and degree of conversion.',
  'Predict whether 10 g of CaCO₃ could give 44 g of CO₂ merely because M(CO₂)=44 g/mol. Compare this product’s mass with the starting mass for this particular reaction. Decide which portion of the mass remains in CaO.',
  'Use M(CaCO₃)=100, M(CaO)=56 and M(CO₂)=44 g/mol. The coefficients are 1:1:1. Write units for all three conversions, and separately check 100=56+44 for one mole. These data do not determine reaction duration or energy consumption.',
  'n(CaCO₃)=10.0/100=0.100 mol, so n(CO₂)=0.100 mol and m(CO₂)=0.100×44=4.40 g. CaO has mass 0.100×56=5.60 g. Product masses total 4.40+5.60=10.00 g, consistent with conservation in the model.',
  '**Mistake:** stopping at 0.100 and labelling it grams when the calculated unit is mol. Another mistake is treating the full mass of an impure sample as pure CaCO₃. For a mixture, first find the reacting substance’s mass and account separately for the other components.',
  'Transfer: 5.30 g of pure Na₂CO₃ reacts completely with sufficient HCl: Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂. Use M(Na₂CO₃)=106 and M(NaCl)=58.5 g/mol. Find the theoretical NaCl mass. No gas collection is needed.',
  '**Check:** n(Na₂CO₃)=0.0500 mol, n(NaCl)=0.100 mol and m(NaCl)=5.85 g. The salt mass exceeds the carbonate mass because some of its composition comes from HCl. Include all reactants in a total mass check, not just the first substance listed.'
 )
},
{
 'stem': '47-limiting-reactant-and-excess',
 'ru': (
  'Как понять, какого реагента не хватит?',
  'Сравним количества с коэффициентами, найдём теоретический продукт и остаток избытка.',
  'В бумажной модели 2H₂ + O₂ → 2H₂O даны 3 mol H₂ и 2 mol O₂. Оба числа положительны: закончится ли всё одновременно? Сначала составьте карточки в отношении 2:1 и посмотрите, какой набор кончается раньше. Реальные газы не используйте.',
  '**Ограничивающий реагент** задаёт максимальное количество продукта в принятой модели полного превращения. **Избыток** — часть другого реагента, остающаяся после этого. Для одной заданной реакции сравните n/ν для каждого реагента, где ν — его положительный коэффициент. Наименьшее отношение задаёт предел. При равенстве смесь стехиометрическая.',
  'Предскажите, можно ли выбрать ограничивающий реагент только по меньшей массе или меньшему числу молей. Для 3 mol H₂ и 2 mol O₂ сравните не 3 с 2, а 3/2 с 2/1. Объясните, зачем делить на коэффициент.',
  'Запишите начальные количества, расход и остаток отдельными строками. Для H₂ коэффициент 2, для O₂ — 1. Отношения 1,5 mol и 2 mol описывают одинаково нормированные возможности реакции; число «1,5» не является количеством любого продукта без учёта его коэффициента.',
  'Минимум n/ν равен 1,5 mol, поэтому H₂ ограничивает. Израсходуются 3 mol H₂ и 1,5 mol O₂, образуется 3 mol H₂O. Кислорода останется 2−1,5=0,5 mol. Его остаток не прибавляется к количеству воды.',
  '**Ошибка:** объявить O₂ ограничивающим, потому что 2 меньше 3, или перевести весь исходный O₂ в продукт. Выход ограничен одновременно обоими исходными количествами. Отрицательный остаток означает ошибку модели или расчёта. Здесь предполагаются одна реакция и полное превращение ограничивающего реагента.',
  'Новый случай той же модели: 4 mol H₂ и 1 mol O₂. Найдите ограничивающий реагент, количество H₂O и остаток H₂. Затем проверьте набор 2 mol H₂ и 1 mol O₂: есть ли избыток в идеальной модели?',
  '**Проверка:** в первом случае ограничивает O₂; образуется 2 mol H₂O, остаётся 2 mol H₂. Во втором количества стехиометрические, избытка нет. Проверьте атомы не только в продуктах, но и в оставшихся реагентах; остаток тоже часть итоговой системы.'
 ),
 'kz': (
  'Қай реагенттің жетпейтінін қалай анықтаймыз?',
  'Зат мөлшерлерін коэффициенттермен салыстырып, теориялық өнім мен артық зат қалдығын табамыз.',
  '2H₂ + O₂ → 2H₂O қағаз моделінде 3 mol H₂ және 2 mol O₂ берілген. Екі сан да оң: бәрі бір уақытта таусыла ма? Карточкаларды 2:1 қатынасымен топтап, қай жиын ертерек аяқталатынын қараңыз. Нақты газдарды пайдаланбаңыз.',
  '**Шектеуші реагент** толық айналу моделіндегі өнімнің ең көп мөлшерін анықтайды. **Артық реагент** — содан кейін басқа реагенттің қалатын бөлігі. Бір берілген реакцияда әр реагент үшін n/ν салыстырыңыз; ν — оның оң коэффициенті. Ең кіші қатынас шекті береді. Қатынастар тең болса, қоспа стехиометриялық.',
  'Шектеушіні тек кіші масса не аз моль саны бойынша таңдауға бола ма? 3 mol H₂ және 2 mol O₂ үшін 3 пен 2-ні емес, 3/2 мен 2/1-ді салыстырыңыз. Коэффициентке неліктен бөлетінімізді түсіндіріңіз.',
  'Бастапқы мөлшерді, шығынды және қалдықты бөлек жолдарға жазыңыз. H₂ коэффициенті 2, O₂ коэффициенті 1. 1,5 mol және 2 mol қатынастары реакция мүмкіндігін бірдей негізге келтіреді; «1,5» саны өнім коэффициентін ескермей кез келген өнім мөлшері бола алмайды.',
  'Ең кіші n/ν — 1,5 mol, сондықтан H₂ шектейді. 3 mol H₂ мен 1,5 mol O₂ жұмсалып, 3 mol H₂O түзіледі. Оттек қалдығы 2−1,5=0,5 mol. Оның қалдығы су мөлшеріне қосылмайды.',
  '**Қате:** 2 саны 3-тен аз болғандықтан O₂-ні шектеуші деу немесе бастапқы O₂-нің бәрін өнімге айналдыру. Өнімді екі реагент мөлшері де шектейді. Теріс қалдық модель не есеп қатесін көрсетеді. Мұнда бір реакция және шектеуші реагенттің толық айналуы қабылданған.',
  'Сол модельдің жаңа жағдайы: 4 mol H₂ және 1 mol O₂. Шектеуші реагентті, H₂O мөлшерін, H₂ қалдығын табыңыз. Кейін 2 mol H₂ мен 1 mol O₂ жиынын тексеріңіз: идеал модельде артық зат бар ма?',
  '**Тексеру:** бірінші жағдайда O₂ шектейді; 2 mol H₂O түзіліп, 2 mol H₂ қалады. Екінші жағдайда қатынас стехиометриялық, артық зат жоқ. Атомдарды өнімде ғана емес, қалған реагенттерде де тексеріңіз: қалдық та соңғы жүйеге кіреді.'
 ),
 'en': (
  'How do we identify the reactant that runs out first?',
  'Compare amounts with coefficients to find theoretical product and excess remaining.',
  'A paper model for 2H₂ + O₂ → 2H₂O supplies 3 mol of H₂ and 2 mol of O₂. Both numbers are positive: will everything run out together? Group cards in a 2:1 ratio and see which supply ends first. Do not use actual gases.',
  'The **limiting reactant** sets the maximum product amount in this complete-conversion model. An **excess** is the portion of another reactant left afterward. For one specified reaction, compare n/ν for each reactant, where ν is its positive coefficient. The smallest ratio sets the limit. Equal ratios describe a stoichiometric mixture.',
  'Can the limiting reactant be selected from the smaller mass or smaller mole number alone? For 3 mol of H₂ and 2 mol of O₂, compare 3/2 with 2/1, not simply 3 with 2. Explain why division by the coefficient is needed.',
  'Write initial amount, consumption and remainder on separate rows. H₂ has coefficient 2 and O₂ coefficient 1. The ratios 1.5 mol and 2 mol normalize each supply to the same reaction basis. The number 1.5 is not automatically the amount of every product; its coefficient must also be used.',
  'The minimum n/ν is 1.5 mol, so H₂ limits the reaction. It consumes 3 mol of H₂ and 1.5 mol of O₂ and forms 3 mol of H₂O. Oxygen remaining is 2−1.5=0.5 mol. That remainder is not added to the amount of water.',
  '**Mistake:** choosing O₂ because 2 is smaller than 3, or converting all the initial O₂ into product. Both reactant supplies constrain the outcome. A negative remainder signals a calculation or model error. We assume one reaction and complete conversion of the limiting reactant.',
  'Transfer: the same model now supplies 4 mol of H₂ and 1 mol of O₂. Find the limiting reactant, H₂O formed and H₂ left. Then check 2 mol of H₂ with 1 mol of O₂: is any excess left in the ideal model?',
  '**Check:** O₂ limits the first case; 2 mol of H₂O forms and 2 mol of H₂ remains. The second mixture is stoichiometric, with no excess. Check atoms in the remaining reactants as well as in the products: leftovers are part of the final system.'
 )
},
{
 'stem': '48-theoretical-and-actual-yield',
 'ru': (
  'Почему полученная масса отличается от расчётной?',
  'Определим выход продукта, учтём чистоту образца и проверим правдоподобие результата.',
  'По расчёту можно получить 5,85 г NaCl, а в готовом отчёте указано 4,68 г чистого сухого NaCl. Куда делись остальные 1,17 г? Можно ли из этих двух чисел установить одну точную причину? Это анализ чужих данных, а не задание проводить синтез.',
  '**Теоретический выход** задаёт максимум продукта по выбранной реакции и ограничивающему реагенту. **Фактический выход** — полученное количество выбранного продукта, определённое измерением с учётом его чистоты. **Процентный выход η** = m(факт)/m(теория)×100% для одного и того же продукта. Можно сравнивать и количества вещества, если основания согласованы.',
  'Предскажите выход для половины теоретической массы. Должен ли выход 80% означать, что ровно 20% исходного вещества не прореагировало? Перечислите другие возможные объяснения: потери при отделении, побочные реакции, неполное превращение, измерительная ошибка.',
  'Даны 5,85 г теоретически и 4,68 г фактически для чистого сухого NaCl. Делите массу продукта на массу того же продукта, а не на массу исходного Na₂CO₃. Отдельный отчёт содержит 6,20 г влажного остатка: его масса ещё не является массой чистого NaCl.',
  'η = 4,68/5,85×100% = 80%. Два числа показывают выход, но не различают все причины потерь. Для влажного остатка 6,20/5,85×100%≈106%: это сигнал проверить воду, примеси, расчёт и измерение, а не доказательство получения вещества сверх корректного теоретического предела.',
  '**Ошибка:** объяснить любой недостаток продукта одной только потерей массы «из мира». Вещество может остаться в растворе или перейти в другие продукты. Выход не равен чистоте: 80% выхода и 80% NaCl в смеси отвечают на разные вопросы.',
  'Если теоретически ожидается 10,0 г продукта, а выход составляет 75%, сколько чистого продукта получено? В другом отчёте масса смеси 5,00 г, массовая доля целевого продукта 90%, теоретически 6,00 г: сначала найдите чистую массу, затем выход.',
  '**Проверка:** 7,50 г в первом случае. Во втором m(продукта)=5,00×0,90=4,50 г, η=4,50/6,00×100%=75%. Подпишите, какая доля описывает чистоту, а какая — выход. Если данных о чистоте нет, не выдавайте массу смеси за доказанную массу продукта.'
 ),
 'kz': (
  'Алынған масса неліктен есептік массадан өзгеше?',
  'Өнім шығымын анықтап, үлгі тазалығын және нәтиженің қисындылығын тексереміз.',
  'Есеп бойынша 5,85 г NaCl алуға болады, ал дайын есепте 4,68 г таза құрғақ NaCl көрсетілген. Қалған 1,17 г қайда кетті? Осы екі саннан жалғыз нақты себепті табуға бола ма? Бұл — дайын деректі талдау, синтез жасау тапсырмасы емес.',
  '**Теориялық шығым** таңдалған реакция мен шектеуші реагент бойынша ең көп өнімді береді. **Нақты шығым** — тазалығы ескеріліп өлшенген мақсатты өнім мөлшері. Бір өнім үшін **пайыздық шығым η** = m(нақты)/m(теория)×100%. Негіздері сәйкестенсе, зат мөлшерлерін де салыстыруға болады.',
  'Теориялық массаның жартысы алынғанда шығымды болжаңыз. 80% шығым бастапқы заттың дәл 20%-ы әрекеттеспегенін білдіре ме? Бөлу кезіндегі шығын, жанама реакциялар, толық айналмау, өлшеу қатесі сияқты басқа түсіндірулерді атаңыз.',
  'Таза құрғақ NaCl үшін теориялық 5,85 г және нақты 4,68 г берілген. Өнім массасын бастапқы Na₂CO₃ массасына емес, сол өнімнің теориялық массасына бөліңіз. Бөлек есепте 6,20 г ылғал қалдық бар: бұл әлі таза NaCl массасы емес.',
  'η = 4,68/5,85×100% = 80%. Екі сан шығымды береді, бірақ шығынның барлық себебін ажыратпайды. Ылғал қалдық үшін 6,20/5,85×100%≈106%: бұл теориялық шектен артық зат алуды дәлелдемейді, суды, қоспаларды, есеп пен өлшемді тексеру белгісі.',
  '**Қате:** өнім тапшылығын массаның «әлемнен жоғалуымен» түсіндіру. Зат ерітіндіде қалуы не басқа өнімге айналуы мүмкін. Шығым тазалыққа тең емес: 80% шығым мен қоспадағы 80% NaCl бөлек сұрақтарға жауап береді.',
  'Теориялық 10,0 г өнім күтіліп, шығым 75% болса, қанша таза өнім алынды? Басқа есепте қоспа массасы 5,00 г, мақсатты өнімнің массалық үлесі 90%, теориялық масса 6,00 г: алдымен таза массаны, кейін шығымды табыңыз.',
  '**Тексеру:** бірінші жағдайда 7,50 г. Екінші жағдайда m(өнім)=5,00×0,90=4,50 г, η=4,50/6,00×100%=75%. Қай үлес тазалықты, қайсысы шығымды сипаттайтынын белгілеңіз. Тазалық дерегі болмаса, қоспа массасын расталған өнім массасы деп алмаңыз.'
 ),
 'en': (
  'Why can the collected mass differ from the calculated mass?',
  'Determine product yield, account for purity and check whether a result is plausible.',
  'A calculation predicts 5.85 g of NaCl, while a supplied report lists 4.68 g of pure dry NaCl. Where did the other 1.17 g go? Can those two numbers identify one definite cause? Analyse supplied data rather than carrying out a synthesis.',
  '**Theoretical yield** is the maximum product from the selected reaction and limiting reactant. **Actual yield** is the quantity of the selected product obtained, measured with its purity accounted for. **Percentage yield η** = m(actual)/m(theoretical)×100% for the same product. Amounts in moles may also be compared on a consistent basis.',
  'Predict the yield when half the theoretical mass is collected. Does 80% yield mean exactly 20% of the starting material did not react? List alternatives such as separation losses, side reactions, incomplete conversion and measurement error.',
  'The supplied masses are 5.85 g theoretical and 4.68 g actual for pure dry NaCl. Divide product mass by the theoretical mass of that same product, not by the initial Na₂CO₃ mass. A separate report gives 6.20 g of wet residue; this is not yet a measurement of pure NaCl mass.',
  'η = 4.68/5.85×100% = 80%. The two numbers establish yield but cannot distinguish every cause of loss. The wet residue would give 6.20/5.85×100%≈106%, signalling a need to check water, impurities, calculation and measurement rather than proving production beyond a correct theoretical limit.',
  '**Mistake:** explaining every shortage as mass disappearing from the world. Material may remain in solution or form other products. Yield is different from purity: 80% yield and 80% NaCl in a mixture answer different questions.',
  'If the theoretical product mass is 10.0 g and yield is 75%, how much pure product was obtained? In another report, a 5.00 g mixture has a target-product mass fraction of 90%, with a theoretical mass of 6.00 g. Find pure-product mass first, then yield.',
  '**Check:** 7.50 g in the first case. In the second, m(product)=5.00×0.90=4.50 g and η=4.50/6.00×100%=75%. Label which fraction describes purity and which describes yield. Without purity evidence, do not treat mixture mass as established product mass.'
 )
},
{
 'stem': '49-salt-quantitative-case-review',
 'ru': (
  'Дело о соли: как проверить всю цепочку расчёта?',
  'Объединим состав, моли, ограничивающий реагент, чистоту и выход в одном отчёте.',
  'Вымышленный отчёт: образец 13,25 г содержит 80% Na₂CO₃ по массе; остальная часть в данной модели инертна. Добавлено 0,150 mol HCl. Реакция Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂. Получено 8,00 г сухой смеси с 87,75% NaCl. Можно ли вычислить выход соли? Никаких реальных реакций выполнять не нужно.',
  'Сначала различите три доли: чистоту исходного образца, долю NaCl в полученной смеси и процентный выход. Количество реагирующего вещества считаем по его собственной массе. Теоретический продукт ограничивается обоими реагентами. Фактическую массу выбранного продукта получаем после поправки на его содержание в смеси.',
  'Предскажите, можно ли принять все 13,25 г за Na₂CO₃ и все 8,00 г за NaCl. Достаточно ли 0,150 mol HCl для полного превращения 0,100 mol карбоната? Назовите нужное отношение, ещё не вычисляя массу соли.',
  'Используйте M(Na₂CO₃)=106 и M(NaCl)=58,5 g/mol. Даны именно 0,150 mol вещества HCl, а не масса или объём его раствора. Инертность примеси — условие вымышленной модели, не утверждение о любой неизвестной примеси. Состав отчёта уже установлен независимым анализом.',
  'Масса Na₂CO₃ равна 13,25×0,80=10,60 г, n=0,100 mol. Сравниваем 0,100/1 и 0,150/2=0,0750 mol: ограничивает HCl. Теоретически n(NaCl)=2×0,0750=0,150 mol, m=8,775 г. Фактически NaCl: 8,00×0,8775=7,020 г. Выход 7,020/8,775×100%=80,0%.',
  '**Ошибка:** считать выход 8,00/8,775 или ограничивать его всем карбонатом. В принятой модели остаётся 0,100−0,0750=0,0250 mol Na₂CO₃, то есть 2,65 г, помимо исходной инертной примеси. Эти остатки не нужно называть массой чистой соли. По выходу нельзя определить единственную причину потерь.',
  'Новый отчёт: 6,625 г исходного образца с теми же 80% Na₂CO₃; HCl теперь 0,200 mol. Получено 5,00 г смеси, содержащей 93,6% NaCl. Найдите ограничивающий реагент, теоретическую и фактическую массу NaCl и процентный выход. Для дополнительной проверки найдите количество оставшегося HCl.',
  '**Проверка:** Na₂CO₃ 5,30 г = 0,0500 mol, он ограничивает; теоретически NaCl 0,100 mol = 5,85 г. Фактически 5,00×0,936=4,68 г; выход 80%. HCl остаётся 0,100 mol в идеальной модели. Закройте ответы и повторите каждый переход с единицами; следующий блок перейдёт к воде и растворам.'
 ),
 'kz': (
  'Тұздың сыры: есептеудің бүкіл тізбегін қалай тексереміз?',
  'Құрамды, мольді, шектеуші реагентті, тазалықты және шығымды бір есепке біріктіреміз.',
  'Ойдан алынған есеп: 13,25 г үлгіде масса бойынша 80% Na₂CO₃ бар; қалған бөлігі осы модельде инертті. 0,150 mol HCl қосылған. Реакция: Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂. Алынған 8,00 г құрғақ қоспада 87,75% NaCl бар. Тұз шығымын есептеуге бола ма? Нақты реакция орындау қажет емес.',
  'Алдымен үш үлесті ажыратыңыз: бастапқы үлгі тазалығы, алынған қоспадағы NaCl үлесі және пайыздық шығым. Әрекеттесетін зат мөлшерін оның өз массасынан табамыз. Теориялық өнімді екі реагент те шектейді. Таңдалған өнімнің нақты массасын қоспадағы құрамына түзету жасап аламыз.',
  '13,25 г-ның бәрін Na₂CO₃, ал 8,00 г-ның бәрін NaCl деп алуға бола ма? 0,100 mol карбонатты толық айналдыруға 0,150 mol HCl жеткілікті ме? Тұз массасын есептемей тұрып керекті қатынасты атаңыз.',
  'M(Na₂CO₃)=106 және M(NaCl)=58,5 g/mol қолданыңыз. Дәл 0,150 mol HCl заты берілген, оның ерітіндісінің массасы не көлемі емес. Қоспаның инерттілігі — ойдан берілген модель шарты, кез келген белгісіз қоспаға ортақ тұжырым емес. Есеп құрамы бөлек талдаумен анықталған.',
  'Na₂CO₃ массасы 13,25×0,80=10,60 г, n=0,100 mol. 0,100/1 мен 0,150/2=0,0750 mol салыстырылады: HCl шектейді. Теориялық n(NaCl)=2×0,0750=0,150 mol, m=8,775 г. Нақты NaCl: 8,00×0,8775=7,020 г. Шығым 7,020/8,775×100%=80,0%.',
  '**Қате:** шығымды 8,00/8,775 деп есептеу немесе бүкіл карбонатпен шектеу. Қабылданған модельде бастапқы инертті қоспадан бөлек 0,100−0,0750=0,0250 mol, яғни 2,65 г Na₂CO₃ қалады. Бұл қалдықтар таза тұз массасы емес. Шығымнан шығынның жалғыз себебін анықтау мүмкін емес.',
  'Жаңа есеп: бастапқы үлгі 6,625 г, Na₂CO₃ үлесі сол 80%; енді HCl 0,200 mol. 93,6% NaCl бар 5,00 г қоспа алынған. Шектеушіні, NaCl теориялық және нақты массасын, пайыздық шығымды табыңыз. Қосымша тексеру үшін қалған HCl зат мөлшерін есептеңіз.',
  '**Тексеру:** Na₂CO₃ 5,30 г = 0,0500 mol, сол шектейді; теориялық NaCl 0,100 mol = 5,85 г. Нақты масса 5,00×0,936=4,68 г; шығым 80%. Идеал модельде HCl 0,100 mol қалады. Жауаптарды жауып, әр өтуді бірлігімен қайталаңыз; келесі блок су мен ерітінділерге көшеді.'
 ),
 'en': (
  'The salt case: how can we check the entire calculation chain?',
  'Combine composition, moles, limiting reactants, purity and yield in one report.',
  'An invented report describes a 13.25 g sample containing 80% Na₂CO₃ by mass; the remainder is inert in this model. It receives 0.150 mol of HCl. The reaction is Na₂CO₃ + 2HCl → 2NaCl + H₂O + CO₂. The dry mixture obtained weighs 8.00 g and contains 87.75% NaCl. Can we calculate salt yield? No actual reaction is required.',
  'Separate three fractions: starting-sample purity, NaCl content of the collected mixture and percentage yield. Calculate reacting amount from the reacting substance’s own mass. Both reactants constrain theoretical product. Correct collected-mixture mass for its composition to find actual target-product mass.',
  'May all 13.25 g be treated as Na₂CO₃ and all 8.00 g as NaCl? Is 0.150 mol of HCl enough to convert 0.100 mol of carbonate completely? State the required ratio before calculating any salt mass.',
  'Use M(Na₂CO₃)=106 and M(NaCl)=58.5 g/mol. The data specify 0.150 mol of HCl substance, not the mass or volume of its solution. Inertness is a supplied model assumption, not a claim about every unidentified impurity. Composition is given from separate analysis.',
  'Na₂CO₃ mass is 13.25×0.80=10.60 g, or 0.100 mol. Compare 0.100/1 with 0.150/2=0.0750 mol: HCl limits. Theoretical n(NaCl)=2×0.0750=0.150 mol and m=8.775 g. Actual NaCl is 8.00×0.8775=7.020 g. Yield is 7.020/8.775×100%=80.0%.',
  '**Mistake:** calculating yield as 8.00/8.775 or using all the carbonate as the limit. The ideal model leaves 0.100−0.0750=0.0250 mol, or 2.65 g, of Na₂CO₃ in addition to the original inert impurity. These leftovers are not pure salt mass. Yield alone cannot establish a unique cause of losses.',
  'Transfer: the starting sample is 6.625 g with the same 80% Na₂CO₃, and HCl is now 0.200 mol. The collected mixture is 5.00 g containing 93.6% NaCl. Find the limiting reactant, theoretical and actual NaCl masses and percentage yield. Also calculate the remaining HCl amount.',
  '**Check:** Na₂CO₃ is 5.30 g = 0.0500 mol and limits; theoretical NaCl is 0.100 mol = 5.85 g. Actual mass is 5.00×0.936=4.68 g, giving 80% yield. The ideal model leaves 0.100 mol of HCl. Hide the answers and reconstruct each step with units; the next block moves to water and solutions.'
 )
},
]
