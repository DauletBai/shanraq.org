#!/usr/bin/env python3
"""Render the fifth chemistry block from original trilingual teaching text."""
from textwrap import wrap
try:
    from .chemistry_reactions_data import LESSONS
    from .chemistry_bonds_data import LESSONS as PREVIOUS
    from .build_chemistry_bonds import LABELS as BASE_LABELS
    from .build_chemistry_elements import start_svg, text
    from .build_chemistry_particles import slug, NOT_TO_SCALE
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
except ImportError:
    from chemistry_reactions_data import LESSONS
    from chemistry_bonds_data import LESSONS as PREVIOUS
    from build_chemistry_bonds import LABELS as BASE_LABELS
    from build_chemistry_elements import start_svg, text
    from build_chemistry_particles import slug, NOT_TO_SCALE
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR

SOURCES = (
 ('OpenStax: physical and chemical changes', 'https://openstax.org/books/chemistry-2e/pages/1-3-physical-and-chemical-properties'),
 ('OpenStax: writing and balancing equations', 'https://openstax.org/books/chemistry-2e/pages/4-1-writing-and-balancing-chemical-equations'),
 ('ACS: atoms in reactions', 'https://www.acs.org/middleschoolchemistry/lessonplans/chapter6/lesson1.html'),
 ('ACS: mass and reactions', 'https://www.acs.org/middleschoolchemistry/lessonplans/chapter6/lesson2.html'),
 ('ACS: energy and bonds', 'https://www.acs.org/middleschoolchemistry/lessonplans/chapter6/lesson7.html'),
 ('OpenStax: energy basics', 'https://openstax.org/books/chemistry-2e/pages/5-1-energy-basics'),
 ('OpenStax: classifying reactions', 'https://openstax.org/books/chemistry-2e/pages/4-2-classifying-chemical-reactions'),
)
SOURCE_KEYS = ((0,), (1,), (1,2), (1,), (1,6), (3,), (4,5), (6,), (1,3,5))
# Four retrieval checkpoints, not a substitute for the full explanation.
STEPS = {
 'ru': [
 ('НАБЛЮДЕНИЕ','Форма, цвет, пузырьки','ОБЪЯСНЕНИЯ','Рассмотри альтернативы','СОСТАВ','Новые вещества?','ВЫВОД','Укажи границу'),
 ('РЕАГЕНТЫ','Слева от стрелки','ПРОДУКТЫ','Справа от стрелки','КОЭФФИЦИЕНТ','Умножает формулу','ПРОВЕРКА','Каждый элемент'),
 ('ДО','H: 4; O: 2','ПЕРЕГРУППИРУЙ','Сохрани каждый атом','ПОСЛЕ','H: 4; O: 2','МОЛЕКУЛЫ','3 → 2, не атомы'),
 ('ФОРМУЛЫ','Не меняй индексы','КИСЛОРОД','6 O с каждой стороны','АЛЮМИНИЙ','4 Al с каждой стороны','ОТНОШЕНИЕ','4 : 3 : 2'),
 ('СКОБКИ','Сначала состав группы','ИНДЕКС','Затем число групп','КОЭФФИЦИЕНТ','Умножь всё','ПРОВЕРКА','Каждый элемент'),
 ('ГРАНИЦА','Что входит в систему?','ЗАКРЫТАЯ','Нет обмена веществом','ОТКРЫТАЯ','Учти вход и выход','ТОЧНОСТЬ','Измерение не идеально'),
 ('РАЗРЫВ','Энергия нужна','ОБРАЗОВАНИЕ','Энергия выделяется','ИТОГ','Учти оба вклада','ТЕМПЕРАТУРА','Не равна теплоте'),
 ('ФОРМЫ','Соединение и разложение','ОБМЕН','Замещение и обмен','УСЛОВИЯ','Продукты уже заданы','ПРЕДЕЛ','Не прогноз реакции'),
 ('СОСТАВ','Дан анализ продуктов','АТОМЫ','Проверь каждый элемент','МАССА','Учти границу и ошибку','ВЫВОД','Не выдумывай данные')],
 'kz': [
 ('БАҚЫЛАУ','Пішін, түс, көпіршік','ТҮСІНДІРУ','Баламаларды қара','ҚҰРАМ','Жаңа зат бар ма?','ҚОРЫТЫНДЫ','Шегін көрсет'),
 ('РЕАГЕНТТЕР','Бағыттауыштың сол жағы','ӨНІМДЕР','Бағыттауыштың оң жағы','КОЭФФИЦИЕНТ','Формуланы көбейтеді','ТЕКСЕРУ','Әр элемент'),
 ('БАСТАПҚЫ','H: 4; O: 2','ҚАЙТА ТОПТА','Әр атом сақталады','СОҢҒЫ','H: 4; O: 2','МОЛЕКУЛА','3 → 2, атом емес'),
 ('ФОРМУЛА','Индексті өзгертпе','ОТТЕК','Екі жақта да 6 O','АЛЮМИНИЙ','Екі жақта да 4 Al','ҚАТЫНАС','4 : 3 : 2'),
 ('ЖАҚША','Алдымен топ құрамы','ИНДЕКС','Кейін топ саны','КОЭФФИЦИЕНТ','Бәрін көбейт','ТЕКСЕРУ','Әр элемент'),
 ('ШЕКАРА','Жүйеге не кіреді?','ЖАБЫҚ','Зат алмаспайды','АШЫҚ','Кіру, шығуды ескер','ДӘЛДІК','Өлшем мінсіз емес'),
 ('ҮЗІЛУ','Энергия қажет','ТҮЗІЛУ','Энергия бөлінеді','ҚОРЫТЫНДЫ','Екі үлесті ескер','ТЕМПЕРАТУРА','Жылуға тең емес'),
 ('ҮЛГІЛЕР','Қосылу, айырылу','АЛМАСУ','Орынбасу, алмасу','ЖАҒДАЙ','Өнімдер берілген','ШЕК','Реакция болжамы емес'),
 ('ҚҰРАМ','Өнім талдауы берілген','АТОМДАР','Әр элементті тексер','МАССА','Шекара, қатені ескер','ҚОРЫТЫНДЫ','Дерек ойлап таппа')],
 'en': [
 ('OBSERVE','Shape, colour, bubbles','EXPLAIN','Consider alternatives','COMPOSITION','New substances?','CONCLUDE','State the limit'),
 ('REACTANTS','Left of the arrow','PRODUCTS','Right of the arrow','COEFFICIENT','Multiplies a formula','CHECK','Every element'),
 ('BEFORE','H: 4; O: 2','REARRANGE','Keep every atom','AFTER','H: 4; O: 2','MOLECULES','3 → 2, not atoms'),
 ('FORMULAS','Keep the subscripts','OXYGEN','6 O on each side','ALUMINIUM','4 Al on each side','RATIO','4 : 3 : 2'),
 ('BRACKETS','First the group','SUBSCRIPT','Then group count','COEFFICIENT','Multiply everything','CHECK','Every element'),
 ('BOUNDARY','What is in the system?','CLOSED','No matter exchange','OPEN','Account for flow','PRECISION','Readings have limits'),
 ('BREAKING','Requires energy','FORMING','Releases energy','NET CHANGE','Account for both','TEMPERATURE','Is not heat'),
 ('PATTERNS','Combine or decompose','EXCHANGE','Displace or exchange','CONDITIONS','Products are given','LIMIT','Not a prediction'),
 ('COMPOSITION','Product analysis given','ATOMS','Check every element','MASS','Boundary and error','CONCLUDE','Do not invent data')],
}
QUIZ = {
 'ru': [
 (('Пузырьки сами по себе доказывают реакцию?', 'Почему разбитое стекло не обязательно продукт реакции?', 'Что нужно, чтобы отличить кипение от образования нового газа?'), 'Нет; состав стекла может сохраниться; нужны условия и сведения о составе и происхождении газа. Необратимость и внешний вид не заменяют эти данные.'),
 (('Назовите реагенты 2Mg + O₂ → 2MgO.', 'Сколько O в 2MgO?', 'Что означает (g)?'), 'Mg и O₂; два атома O в указанном наборе; газовое состояние. Индекс задаёт состав, коэффициент — число единиц; (g) не обозначает элемент.'),
 (('Сколько H и O слева в 2H₂ + O₂ → 2H₂O?', 'Сохранилось ли число молекул?', 'Достаточно ли одинакового общего числа атомов?'), 'H: 4, O: 2; нет, 3 → 2; нет, нужно сохранить число атомов каждого элемента. Молекулы перегруппировываются, ядра атомов здесь не превращаются.'),
 (('Расставьте коэффициенты Al + O₂ → Al₂O₃.', 'Можно ли заменить Al₂O₃ на AlO?', 'Почему 8:6:4 сокращают?'), '4:3:2; нет, это другой состав; у всех чисел есть общий делитель 2. После сокращения проверьте отдельно 4 Al и 6 O с каждой стороны.'),
 (('Сколько O в Fe₂(SO₄)₃?', 'Сколько H в 6NaOH?', 'Сохраняется ли группа OH в продуктах Ca(OH)₂ + 2HCl?'), '12; 6; нет, продукты CaCl₂ и 2H₂O. Считайте элементы, не переносите группу как неизменную, когда её форма изменилась.'),
 (('Система открыта: 100 → 97 г, только газ вышел. Его масса?', 'Металл 10 → 12 г, поступал только O₂. Сколько поступило?', '100.00 и 100.01 г при ±0.02 г каждое — доказанная потеря?'), '3 г; 2 г; нет, различие меньше указанной погрешности. Включайте вошедшее и вышедшее вещество в баланс.'),
 (('Разрыв стоит 100, образование даёт 140 условных единиц. Итог?', '120 − 80: поглощение или выделение?', 'Рост на 5 °C означает 5 единиц теплоты?'), '−40: выделение 40; +40: поглощение 40; нет, нужны масса, теплоёмкость и условия теплообмена. Это вымышленные энергетические числа.'),
 (('2Mg + O₂ → 2MgO: сколько видов продуктов?', 'Zn + CuSO₄ → ZnSO₄ + Cu — какой школьный тип?', 'Доказывает ли схема обмена, что любые два раствора реагируют?'), 'Один вид; замещение; нет. Коэффициент не создаёт новый вид вещества, а классификация не заменяет проверку условий и продуктов.'),
 (('Уравняйте CaCO₃ + HCl → CaCl₂ + H₂O + CO₂.', '150.00 → 149.56 г, вышел только CO₂. Сколько?', 'Можно ли узнать газ только по уменьшению массы?'), 'Коэффициенты 1:2:1:1:1; 0.44 г; нет. В этом деле состав газа дан отдельным анализом. Отделите баланс атомов, массовый баланс и доказательство состава.')],
 'kz': [
 (('Көпіршік жалғыз өзі реакцияны дәлелдей ме?', 'Сынған шыны неге міндетті түрде реакция өнімі емес?', 'Қайнау мен жаңа газ түзілуін қалай ажыратамыз?'), 'Жоқ; шыны құрамы сақталуы мүмкін; жағдай, газ құрамы және шығу тегі туралы дерек керек. Қайтымсыздық пен сыртқы түр бұл деректі алмастырмайды.'),
 (('2Mg + O₂ → 2MgO реагенттері қандай?', '2MgO ішінде неше O бар?', '(g) нені білдіреді?'), 'Mg және O₂; берілген жиында екі O атомы; газ күйі. Индекс құрамды, коэффициент бірлік санын береді; (g) элемент емес.'),
 (('2H₂ + O₂ → 2H₂O сол жағында неше H және O бар?', 'Молекула саны сақтала ма?', 'Атомдардың жалпы саны тең болса жеткілікті ме?'), 'H: 4, O: 2; жоқ, 3 → 2; жоқ, әр элемент атомдарының саны сақталуы керек. Молекулалар қайта топтасады, мұнда атом ядролары түрленбейді.'),
 (('Al + O₂ → Al₂O₃ коэффициенттерін қойыңыз.', 'Al₂O₃ орнына AlO жазуға бола ма?', '8:6:4 неге қысқарады?'), '4:3:2; жоқ, құрамы басқа; ортақ бөлгіші 2. Қысқартқан соң екі жақтағы 4 Al және 6 O санын бөлек тексеріңіз.'),
 (('Fe₂(SO₄)₃ ішінде неше O бар?', '6NaOH ішінде неше H бар?', 'Ca(OH)₂ + 2HCl өнімдерінде OH тобы сақтала ма?'), '12; 6; жоқ, өнімдер CaCl₂ және 2H₂O. Топ өзгерсе, оны өзгермеген күйде көшірмей, элементтерді санаңыз.'),
 (('Ашық жүйе: 100 → 97 г, тек газ шықты. Оның массасы?', 'Металл 10 → 12 г, тек O₂ кірді. Қанша кірді?', 'Әр өлшем ±0.02 г болса, 100.00 мен 100.01 г жоғалуды дәлелдей ме?'), '3 г; 2 г; жоқ, айырма берілген қателіктен аз. Кірген және шыққан затты теңгерімге қосыңыз.'),
 (('Үзуге 100, түзілуге 140 шартты бірлік берілген. Нәтиже?', '120 − 80: жұту ма, бөлу ме?', '5 °C өсу 5 жылу бірлігі ме?'), '−40: 40 бөлінеді; +40: 40 жұтылады; жоқ, масса, жылусыйымдылық және жылу алмасу жағдайы керек. Энергия сандары ойдан берілген.'),
 (('2Mg + O₂ → 2MgO: өнімнің неше түрі бар?', 'Zn + CuSO₄ → ZnSO₄ + Cu мектептік түрі қандай?', 'Алмасу сызбасы кез келген екі ерітінді әрекеттесетінін дәлелдей ме?'), 'Бір түр; орынбасу; жоқ. Коэффициент жаңа зат түрін жасамайды, жіктеу жағдай мен өнімдерді тексеруді алмастырмайды.'),
 (('CaCO₃ + HCl → CaCl₂ + H₂O + CO₂ теңестіріңіз.', '150.00 → 149.56 г, тек CO₂ шықты. Қанша?', 'Газды тек масса азаюынан анықтаймыз ба?'), 'Коэффициенттер 1:2:1:1:1; 0.44 г; жоқ. Бұл істе газ құрамы бөлек талдаумен берілген. Атом теңгерімін, масса теңгерімін және құрам дәлелін ажыратыңыз.')],
 'en': [
 (('Do bubbles alone prove a reaction?', 'Why is broken glass not necessarily a reaction product?', 'What distinguishes boiling from formation of a new gas?'), 'No; glass composition may be unchanged; conditions and evidence of gas identity and origin are needed. Irreversibility and appearance do not replace that evidence.'),
 (('Name the reactants in 2Mg + O₂ → 2MgO.', 'How many O atoms are in 2MgO?', 'What does (g) mean?'), 'Mg and O₂; two O atoms in this set; gas state. Subscripts describe composition, coefficients count units; (g) is not an element.'),
 (('Count H and O on the left of 2H₂ + O₂ → 2H₂O.', 'Is molecule count conserved?', 'Is the same total atom count enough?'), 'H: 4, O: 2; no, 3 → 2; no, each element must balance. Molecules rearrange; atomic nuclei are not transformed in this chemical model.'),
 (('Balance Al + O₂ → Al₂O₃.', 'May Al₂O₃ be replaced by AlO?', 'Why reduce 8:6:4?'), '4:3:2; no, that changes composition; all coefficients share a factor of 2. Check 4 Al and 6 O on each side after reducing.'),
 (('How many O atoms are in Fe₂(SO₄)₃?', 'How many H atoms are in 6NaOH?', 'Does OH remain intact in the products of Ca(OH)₂ + 2HCl?'), '12; 6; no, the products are CaCl₂ and 2H₂O. Count elements rather than carrying a group unchanged when its form changes.'),
 (('Open system: 100 → 97 g, only gas leaves. Its mass?', 'Metal: 10 → 12 g, only O₂ enters. How much enters?', 'Do 100.00 and 100.01 g, each ±0.02 g, prove a loss?'), '3 g; 2 g; no, the difference is smaller than the stated uncertainty. Include incoming and outgoing matter in the balance.'),
 (('Breaking costs 100; forming releases 140 invented units. Net change?', '120 − 80: absorption or release?', 'Does a 5 °C rise mean 5 units of heat?'), '−40: release of 40; +40: absorption of 40; no, mass, heat capacity and heat exchange conditions are needed. These energy numbers are invented.'),
 (('2Mg + O₂ → 2MgO: how many product species?', 'What school pattern is Zn + CuSO₄ → ZnSO₄ + Cu?', 'Does an exchange pattern prove any two solutions react?'), 'One species; displacement; no. A coefficient does not create a new species, and classification does not replace checking conditions and products.'),
 (('Balance CaCO₃ + HCl → CaCl₂ + H₂O + CO₂.', '150.00 → 149.56 g, only CO₂ leaves. How much?', 'Can mass loss alone identify a gas?'), 'Coefficients 1:2:1:1:1; 0.44 g; no. Gas identity is supplied by separate analysis in this case. Separate atom balance, mass balance and evidence of composition.')],
}

# Each table is supplied data, never an experiment to perform.
def table(number, lang):
    headers = {'ru': ('Данные / элемент','До / слева','После / справа'), 'kz': ('Дерек / элемент','Бастапқы / сол','Соңғы / оң'), 'en': ('Data / element','Before / left','After / right')}[lang]
    rows = {
      32:[('H₂O','H₂O(s)','H₂O(l)'),('H₂O','H₂O(l)','H₂O(g)')],
      33:[('Mg','2','2'),('O','2','2')],
      34:[('H','4','4'),('O','2','2')],
      35:[('Al','4','4'),('O','6','6')],
      36:[('Fe','2','2'),('S','3','3'),('Na','6','6'),('O','18','18'),('H','6','6')],
      37:[('g','100.00','100.00'),('g','100','97 + 3'),('g','10 + 2','12')],
      38:[('A','100 − 140','−40'),('B','120 − 80','+40')],
      39:[('Na','2','2'),('C','1','1'),('O','3','3'),('Ca','1','1'),('Cl','2','2')],
      40:[('Ca','1','1'),('C','1','1'),('O','3','3'),('H','2','2'),('Cl','2','2')],
    }[number]
    context={32:'H₂O(s) / H₂O(l) / H₂O(g)',33:'2Mg + O₂ → 2MgO',34:'2H₂ + O₂ → 2H₂O',35:'4Al + 3O₂ → 2Al₂O₃',36:'Fe₂(SO₄)₃ + 6NaOH → 2Fe(OH)₃ + 3Na₂SO₄',37:'100 = 97 + 3; 10 + 2 = 12',38:{'ru':'Вымышленные единицы энергии','kz':'Ойдан берілген энергия бірлігі','en':'Invented energy units'}[lang],39:'Na₂CO₃ + CaCl₂ → CaCO₃ + 2NaCl',40:'CaCO₃ + 2HCl → CaCl₂ + H₂O + CO₂'}[number]
    return context+'\n\n| '+' | '.join(headers)+' |\n|---|---|---|\n'+'\n'.join('| '+' | '.join(row)+' |' for row in rows)


def map_svg(number, lang):
    parts=start_svg(LESSONS[number-32][lang][0],NOT_TO_SCALE[lang],height=560)
    for i,line in enumerate(wrap(LESSONS[number-32][lang][0],60)):
        parts.append(text(40,55+i*37,line,30))
    parts.append(text(40,150,NOT_TO_SCALE[lang],20,'#bedde4'))
    steps=STEPS[lang][number-32]
    for i in range(4):
        x=40+310*i
        parts.append(f'<rect x="{x}" y="195" width="275" height="295" rx="18" fill="#e6f4ef"/>')
        parts.append(text(x+18,244,steps[2*i],20,'#143242'))
        for j,line in enumerate(wrap(steps[2*i+1],17,break_long_words=False)):
            parts.append(text(x+18,310+35*j,line,24,'#143242'))
        if i<3:parts.append(text(x+283,365,'→',28))
    return '\n'.join(parts+['</svg>'])+'\n'


def visual(number, lang):
    parts=start_svg(LESSONS[number-32][lang][0],NOT_TO_SCALE[lang])
    for i,line in enumerate(wrap(f'{number} · '+LESSONS[number-32][lang][0],67)):
        parts.append(text(40,55+i*37,line,29))
    parts.append(text(40,155,NOT_TO_SCALE[lang],20,'#bedde4'))
    def label(x,y,s,size=38):parts.append(text(x,y,s,size,anchor='middle'))
    if number==34:
        label(640,225,'2H₂ + O₂ → 2H₂O',48)
        for side,molecules in [('before', [[('H',160,350),('H',220,350)],[('H',160,490),('H',220,490)],[('O',410,420),('O',470,420)]]),('after', [[('O',840,350),('H',780,405),('H',900,405)],[('O',1070,490),('H',1010,545),('H',1130,545)]])]:
            for molecule in molecules:
                for element,x,y in molecule:
                    color='#98d5d8' if element=='O' else '#f0c18b'
                    parts.append(f'<circle data-side="{side}" data-element="{element}" cx="{x}" cy="{y}" r="28" fill="{color}"/>')
                    parts.append(text(x,y+9,element,25,'#143242','middle'))
        label(640,430,'→',70)
        label(320,650,'H: 4 · O: 2',36);label(960,650,'H: 4 · O: 2',36)
    else:
        contents={
          32:[('H₂O(s) → H₂O(l)',STEPS[lang][0][5]),('CO₂(aq) → CO₂(g)',STEPS[lang][0][3]),('?',STEPS[lang][0][7])],
          33:[('2Mg + O₂','→ 2MgO'),('Mg: 2 → 2','O: 2 → 2'),('s · l · g · aq',STEPS[lang][1][7])],
          35:[('4Al + 3O₂ → 2Al₂O₃','4 : 3 : 2'),('Al: 4 → 4','O: 6 → 6'),('8 : 6 : 4','↓ ÷ 2')],
          36:[('Fe₂(SO₄)₃','Fe: 2 · S: 3 · O: 12'),('6NaOH','Na: 6 · O: 6 · H: 6'),('2Fe(OH)₃ + 3Na₂SO₄','Fe: 2 · S: 3 · Na: 6 · O: 18 · H: 6')],
          37:[(STEPS[lang][5][2], '100.00 g → 100.00 g'),(STEPS[lang][5][4],'100 g = 97 g + 3 g'),('±0.02 g','100.00 g ≈ 100.01 g')],
          38:[({'ru':'Вымышленные единицы','kz':'Ойдан берілген бірліктер','en':'Invented energy units'}[lang],'100 − 140 = −40'),(STEPS[lang][6][5],'120 − 80 = +40'),(STEPS[lang][6][6], 'ΔT ≠ Q')],
          39:[('2Mg + O₂ → 2MgO','CaCO₃ → CaO + CO₂'),('Zn + CuSO₄ → ZnSO₄ + Cu','Na₂CO₃ + CaCl₂ → CaCO₃ + 2NaCl'),(STEPS[lang][7][4],STEPS[lang][7][7])],
          40:[('CaCO₃ + 2HCl →','CaCl₂ + H₂O + CO₂'),('150.00 g − 149.56 g','0.44 g'),(STEPS[lang][8][4],STEPS[lang][8][5])],
        }[number]
        for i,(first,second) in enumerate(contents):
            y=200+i*158
            parts.append(f'<rect x="40" y="{y}" width="1200" height="140" rx="18" fill="#173e4d"/>')
            label(640,y+52,first,34)
            label(640,y+108,second,34)
    return '\n'.join(parts+['</svg>'])+'\n'


def render(number, lang):
    row=LESSONS[number-32]
    title,summary,*sections=row[lang]
    labels=BASE_LABELS[lang]
    notebook={
      'ru':'Для каждой ошибки отметьте причину: неверно определён состав, изменён индекс, пропущен атом, не учтена граница системы или перепутаны теплота и температура. Исправьте объяснение и проверьте другой пример без подсказки.',
      'kz':'Әр қатенің себебін белгілеңіз: құрам қате, индекс өзгерген, атом ұмытылған, жүйе шекарасы ескерілмеген немесе жылу мен температура шатасқан. Түсіндіруді түзетіп, басқа мысалды кеңессіз тексеріңіз.',
      'en':'For each error, identify its cause: wrong composition, changed subscript, missing atom, neglected system boundary, or confusion between heat and temperature. Correct your explanation and check another example without hints.'}[lang]
    out=[f'# {title}',f'_Лид (summary):_ **{summary}**',f'## {HEAD[lang][0]}',
         f'![{title}](/static/course/chemistry/map-{number:02d}-{lang}.svg)',
         f'## {labels[0]}',f'![{NOT_TO_SCALE[lang]}](/static/course/chemistry/visuals/{number:02d}-reactions-{lang}.svg)',labels[5]]
    for i,(heading,section) in enumerate(zip(HEAD[lang][1:9],sections,strict=True)):
        out.extend([f'## {heading}',section])
        if i==3:out.append(table(number,lang))
    questions,answer=QUIZ[lang][number-32]
    out.extend([f'## {labels[1]}','\n'.join(f'{i}. {q}' for i,q in enumerate(questions,1)),f'## {labels[2]}',answer,
                f'## {labels[3]}',notebook,f'## {HEAD[lang][10]}',RETURN[lang],f'## {labels[4]}',
                '\n'.join(f'- [{SOURCES[k][0]}]({SOURCES[k][1]})' for k in SOURCE_KEYS[number-32])])
    previous=PREVIOUS[-1] if number==32 else LESSONS[number-33]
    out.append(f'[{previous[lang][0]}](/read/{slug(number-1,previous["stem"])}?lang={lang})')
    if number<40:
        following=LESSONS[number-31]
        out.append(f'[{following[lang][0]}](/read/{slug(number+1,following["stem"])}?lang={lang})')
    return '\n\n'.join(out)+'\n'


def bridge_body(lang):
    """Release-only onward link; block four remains independently releasable."""
    try:
        from .build_chemistry_bonds import render as render_previous
    except ImportError:
        from build_chemistry_bonds import render as render_previous
    row=LESSONS[0]
    return render_previous(31,lang).split('\n\n',2)[2].rstrip()+f'\n\n[{row[lang][0]}](/read/{slug(32,row["stem"])}?lang={lang})\n'


def main():
    (MAP_DIR/'visuals').mkdir(parents=True,exist_ok=True)
    for number,row in enumerate(LESSONS,32):
        for lang in LANGS:
            suffix='' if lang=='ru' else '-'+lang
            (LESSON_DIR/f'{row["stem"]}{suffix}.md').write_text(render(number,lang),encoding='utf-8')
            (MAP_DIR/f'map-{number:02d}-{lang}.svg').write_text(map_svg(number,lang),encoding='utf-8')
            (MAP_DIR/f'visuals/{number:02d}-reactions-{lang}.svg').write_text(visual(number,lang),encoding='utf-8')

if __name__=='__main__':main()
