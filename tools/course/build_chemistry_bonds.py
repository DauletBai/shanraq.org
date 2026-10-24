#!/usr/bin/env python3
"""Render block four and original localized bonding diagrams."""
from textwrap import wrap

try:
    from .chemistry_bonds_data import LESSONS
    from .chemistry_elements_data import LESSONS as PREVIOUS
    from .build_chemistry_elements import start_svg, text
    from .build_chemistry_particles import slug, NOT_TO_SCALE
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
except ImportError:
    from chemistry_bonds_data import LESSONS
    from chemistry_elements_data import LESSONS as PREVIOUS
    from build_chemistry_elements import start_svg, text
    from build_chemistry_particles import slug, NOT_TO_SCALE
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR

SOURCES = {
    'energy': ('OpenStax: covalent bonding and energy', 'https://openstax.org/books/chemistry-2e/pages/7-2-covalent-bonding'),
    'covalent': ('ACS: covalent bonding', 'https://www.acs.org/middleschoolchemistry/lessonplans/chapter4/lesson4.html'),
    'ionic': ('OpenStax: ionic bonding', 'https://openstax.org/books/chemistry-2e/pages/7-1-ionic-bonding'),
    'lewis': ('OpenStax: Lewis structures', 'https://openstax.org/books/chemistry-2e/pages/7-3-lewis-symbols-and-structures'),
    'shape': ('OpenStax: geometry and polarity', 'https://openstax.org/books/chemistry-2e/pages/7-6-molecular-structure-and-polarity'),
    'solids': ('OpenStax: solid structures', 'https://openstax.org/books/chemistry-2e/pages/10-5-the-solid-state-of-matter'),
    'formulas': ('OpenStax: chemical formulas', 'https://openstax.org/books/chemistry-2e/pages/2-4-chemical-formulas'),
    'naming': ('OpenStax: ionic formulas', 'https://openstax.org/books/chemistry-2e/pages/2-7-chemical-nomenclature'),
}
SOURCE_KEYS = [('energy',), ('covalent','lewis'), ('ionic',), ('shape',), ('solids',),
               ('naming',), ('lewis',), ('formulas','naming'), ('formulas','ionic','shape')]
LABELS = {
 'ru': ('Схема модели', 'Задание', 'Ответы', 'Дневник исследователя', 'Источники и границы',
        'Только бумага и готовые данные. Не получайте газы, не нагревайте вещества, не используйте электрическое оборудование; неизвестные вещества не трогайте, не пробуйте и не нюхайте.',
        'Запишите свой первый ответ до чтения разбора. Для каждой ошибки отметьте конкретный шаг: модель частиц, счёт электронов, баланс зарядов, скобки или граница вывода. Исправьте объяснение, затем придумайте новый пример и решите его без подсказки.'),
 'kz': ('Модель сызбасы', 'Тапсырма', 'Жауаптар', 'Зерттеуші күнделігі', 'Дереккөздер мен шектер',
        'Тек қағаз бен дайын деректер. Газ алмаңыз, заттарды қыздырмаңыз, электр құралдарын қолданбаңыз; белгісіз заттарды ұстамаңыз, дәмін татпаңыз, иіскемеңіз.',
        'Талдауды оқымай тұрып алғашқы жауабыңызды жазыңыз. Әр қатенің қадамын белгілеңіз: бөлшек моделі, электрон санағы, заряд теңгерімі, жақша немесе қорытынды шегі. Түсіндіруді түзетіп, жаңа мысал ойлап тауып, оны кеңессіз шешіңіз.'),
 'en': ('Model diagram', 'Exercise', 'Answers', 'Investigator’s notebook', 'Sources and limits',
        'Use paper and supplied data only. Do not produce gases, heat substances or use electrical equipment; do not touch, taste or smell unknown substances.',
        'Record your first answer before reading the explanation. For each error, identify the step: particle model, electron count, charge balance, brackets or limits of inference. Correct the explanation, invent a new example and solve it without hints.'),
}
# Four recall checkpoints per lesson, localized independently of body text.
STEPS = {
 'ru': [
  ('ЗАРЯДЫ','Ядра и электроны','ВЗАИМОДЕЙСТВИЯ','Притяжение и отталкивание','СВЯЗЬ','Ниже энергия','ПРЕДЕЛ','Не слияние ядер'),
  ('ВНЕШНИЕ e','H: 1, Cl: 7','ОБЩАЯ ПАРА','Линия = два e','Cl₂','14 внешних e','ПРЕДЕЛ','Точки не траектории'),
  ('ИОНЫ','Na⁺ и Cl⁻','ПРИТЯЖЕНИЕ','Противоположные заряды','РЕШЁТКА','Не отдельные пары','ФОРМУЛА','NaCl: 1 : 1'),
  ('СВЯЗЬ','Неравномерная плотность','δ','Частичный заряд','МОЛЕКУЛА','Важна форма','CO₂ / H₂O','Линейная / угловая'),
  ('ЧАСТИЦЫ','Молекулы или решётка','СВЯЗИ','Внутри и между','ПРОВОДИМОСТЬ','Подвижные носители','ПРЕДЕЛ','Одного признака мало'),
  ('ДАНО','Заряды ионов','ПОДБОР','Целые количества','ПРОВЕРКА','Сумма зарядов ноль','ИНДЕКСЫ','Простейшее отношение'),
  ('ВАЛЕНТНОСТЬ','Число связей','КРАТНОСТЬ','Двойная = две','ЗАРЯД','Другой вопрос','ПРЕДЕЛ','Не всегда постоянна'),
  ('СКОБКИ','Умножь группу','ИНДЕКС','Состав единицы','КОЭФФИЦИЕНТ','Умножь всё','ПРЕДЕЛ','2H₂O ≠ H₂O₂'),
  ('ДАНО','Формула или внешний вид','ВЫБЕРИ','Молекула или решётка','ПРОВЕРЬ','Счёт и заряды','ОГРАНИЧЬ','Цвет не состав'),
 ],
 'kz': [
  ('ЗАРЯДТАР','Ядро мен электрон','ӘРЕКЕТТЕСУ','Тартылыс, тебіліс','БАЙЛАНЫС','Энергия төмен','ШЕК','Ядро қосылмайды'),
  ('СЫРТҚЫ e','H: 1, Cl: 7','ОРТАҚ ЖҰП','Сызық = екі e','Cl₂','14 сыртқы e','ШЕК','Нүкте жол емес'),
  ('ИОНДАР','Na⁺ және Cl⁻','ТАРТЫЛЫС','Қарама-қарсы заряд','ТОР','Жеке жұп емес','ФОРМУЛА','NaCl: 1 : 1'),
  ('БАЙЛАНЫС','Тығыздық біркелкі емес','δ','Ішінара заряд','МОЛЕКУЛА','Пішін маңызды','CO₂ / H₂O','Сызықтық / бұрыштық'),
  ('БӨЛШЕКТЕР','Молекула не тор','БАЙЛАНЫС','Ішкі және аралық','ӨТКІЗГІШТІК','Қозғалмалы заряд','ШЕК','Бір белгі аз'),
  ('БЕРІЛГЕН','Ион зарядтары','ТАҢДАУ','Бүтін сандар','ТЕКСЕРУ','Заряд қосындысы нөл','ИНДЕКСТЕР','Қарапайым қатынас'),
  ('ВАЛЕНТТІЛІК','Байланыс саны','ЕСЕЛІК','Қос байланыс = екі','ЗАРЯД','Басқа сұрақ','ШЕК','Әрдайым тұрақты емес'),
  ('ЖАҚША','Топты көбейт','ИНДЕКС','Бірлік құрамы','КОЭФФИЦИЕНТ','Бәрін көбейт','ШЕК','2H₂O ≠ H₂O₂'),
  ('БЕРІЛГЕН','Формула не сыртқы түр','ТАҢДАУ','Молекула не тор','ТЕКСЕРУ','Санақ пен заряд','ШЕКТЕУ','Түс құрам емес'),
 ],
 'en': [
  ('CHARGES','Nuclei and electrons','INTERACTIONS','Attraction and repulsion','BOND','Lower energy','LIMIT','No merging nuclei'),
  ('OUTER e','H: 1, Cl: 7','SHARED PAIR','Line = two e','Cl₂','14 valence e','LIMIT','Dots are not paths'),
  ('IONS','Na⁺ and Cl⁻','ATTRACTION','Opposite charges','LATTICE','Not isolated pairs','FORMULA','NaCl: 1 : 1'),
  ('BOND','Unequal density','δ','Partial charge','MOLECULE','Shape matters','CO₂ / H₂O','Linear / bent'),
  ('PARTICLES','Molecules or lattice','INTERACTIONS','Within and between','CONDUCTIVITY','Mobile carriers','LIMIT','One clue is not enough'),
  ('GIVEN','Ion charges','CHOOSE','Whole-number counts','CHECK','Charge sum is zero','SUBSCRIPTS','Simplest ratio'),
  ('VALENCY','Bond count','MULTIPLICITY','Double counts as two','CHARGE','A different question','LIMIT','Not always constant'),
  ('BRACKETS','Multiply the group','SUBSCRIPT','Unit composition','COEFFICIENT','Multiply everything','LIMIT','2H₂O ≠ H₂O₂'),
  ('GIVEN','Formula or appearance','SELECT','Molecule or lattice','CHECK','Counts and charges','LIMIT','Colour is not identity'),
 ],
}
QUIZ = {
 'ru': [
  (('Связь разорвали до далёких атомов: энергия нужна или выделяется?', 'У H заполненный слой содержит два или восемь электронов?', 'В модели 0 / −5 / +8 какое состояние ниже по энергии?'), 'Энергия нужна; два; −5. Понижение энергии при образовании одной связи не определяет знак энергии целой реакции.'),
  (('Сколько общих пар в H₂?', 'Сколько неподелённых пар у каждого Cl в Cl₂?', 'Сколько валентных электронов всего в HCl?'), 'Одна; три у каждого; восемь. Линия показывает два электрона, а не весь электронный состав атомов.'),
  (('Сколько электронов в Na⁺ при Z=11?', 'Является ли NaCl в кристалле отдельной молекулой?', 'Каково отношение Mg²⁺ и O²⁻ в простейшей нейтральной формуле?'), '10; нет; 1:1. Протоны не изменились. Формульная единица задаёт отношение в решётке.'),
  (('Означает ли δ⁻ целый полученный электрон?', 'Полярна ли линейная симметричная CO₂?', 'Почему полярна угловая H₂O?'), 'Нет; нет; направления смещения электронной плотности не компенсируются. Полярную связь не отождествляем с полярной молекулой.'),
  (('При плавлении льда исчезают ли молекулы H₂O?', 'Какие носители подвижны в расплаве NaCl?', 'Обязательно ли проводник является металлом?'), 'Нет; ионы; нет, графит — контрпример. Ответы относятся к указанным состояниям и условиям.'),
  (('Какая формула у Al³⁺ и Cl⁻?', 'Какова сумма зарядов для MgCl₂?', 'Нужно ли сокращать Ca₂O₂ при заданных Ca²⁺ и O²⁻?'), 'AlCl₃; +2−2=0; да, до CaO. Сокращается отношение ионов, а не произвольная молекулярная формула.'),
  (('Какова валентность O в O₂ с двойной связью?', 'Сколько соседей у каждого N в N₂?', 'Означает ли IV у C в CH₄ заряд +4?'), 'II; один; нет. Тройная связь даёт три общие пары, а валентность не равна заряду.'),
  (('Сколько O в Al₂(SO₄)₃?', 'Сколько H в 4H₂O?', 'Одинаковы ли частицы 2H и H₂?'), '12;8; нет: два отдельных атома и одна молекула. Сначала применяем индекс группы, затем коэффициент ко всей формуле.'),
  (('Что можно узнать о D только по белому цвету?', 'Какая формула у Ca²⁺ и Cl⁻?', 'Сколько Ca и Cl в 3CaCl₂?', 'Даёт ли полярная связь обязательную полярность молекулы?'), 'Лишь внешний вид, состав неизвестен; CaCl₂;3 Ca и 6 Cl; нет, нужна форма. Оцените каждый ответ отдельно и исправьте причину ошибки.'),
 ],
 'kz': [
  (('Байланысты алыс атомдарға дейін үзгенде энергия қажет пе, бөлінеді ме?', 'H толған қабатында екі ме, сегіз электрон ба?', '0 / −5 / +8 моделінде ең төмен энергия қайсы?'), 'Энергия қажет; екі; −5. Бір байланыс түзілгендегі төмендеу тұтас реакция энергиясының таңбасын анықтамайды.'),
  (('H₂-де неше ортақ жұп бар?', 'Cl₂-дегі әр Cl маңында неше бөлінбеген жұп бар?', 'HCl-де барлығы неше валенттік электрон бар?'), 'Бір; әрқайсысында үш; сегіз. Сызық атомның барлық электронын емес, екі электронды көрсетеді.'),
  (('Z=11 болса, Na⁺ неше электронға ие?', 'Кристалдағы NaCl жеке молекула ма?', 'Mg²⁺ пен O²⁻ ең қарапайым бейтарап қатынасы қандай?'), '10; жоқ;1:1. Протондар өзгермейді. Формулалық бірлік тордағы қатынасты береді.'),
  (('δ⁻ бір толық электрон алғанын білдіре ме?', 'Симметриялы сызықтық CO₂ полюсті ме?', 'Бұрыштық H₂O неге полюсті?'), 'Жоқ; жоқ; электрон тығыздығы ығысуының бағыттары өтелмейді. Полюсті байланыс пен полюсті молекуланы теңестірмейміз.'),
  (('Мұз балқығанда H₂O молекулалары жоғала ма?', 'NaCl балқымасында қандай тасымалдаушылар қозғалады?', 'Өткізгіш міндетті түрде металл ма?'), 'Жоқ; иондар; жоқ, графит қарсы мысал. Жауаптар берілген күй мен жағдайға қатысты.'),
  (('Al³⁺ және Cl⁻ формуласы қандай?', 'MgCl₂ зарядтарының қосындысы неге тең?', 'Ca²⁺ және O²⁻ үшін Ca₂O₂ қысқара ма?'), 'AlCl₃; +2−2=0; иә, CaO-ға дейін. Кез келген молекулалық формула емес, иондар қатынасы қысқарады.'),
  (('Қос байланысы бар O₂ ішіндегі O валенттілігі қандай?', 'N₂-де әр N неше көршіге ие?', 'CH₄-тегі C үшін IV заряды +4 деген сөз бе?'), 'II; бір; жоқ. Үш байланыс үш ортақ жұп береді, валенттілік зарядқа тең емес.'),
  (('Al₂(SO₄)₃ ішінде неше O бар?', '4H₂O ішінде неше H бар?', '2H және H₂ бөлшектері бірдей ме?'), '12;8; жоқ: екі жеке атом және бір молекула. Алдымен топ индексін, кейін тұтас формула коэффициентін қолданамыз.'),
  (('Тек ақ түсінен D туралы нені білеміз?', 'Ca²⁺ және Cl⁻ формуласы қандай?', '3CaCl₂ ішінде неше Ca мен Cl бар?', 'Полюсті байланыс молекуланы міндетті түрде полюсті ете ме?'), 'Тек сыртқы түрін, құрам белгісіз; CaCl₂;3 Ca,6 Cl; жоқ, пішін қажет. Әр жауапты бөлек бағалап, қате себебін түзетіңіз.'),
 ],
 'en': [
  (('Does breaking a bond into distant atoms require or release energy?', 'Does H have a filled shell with two or eight electrons?', 'Which state has lowest energy in a 0 / −5 / +8 model?'), 'Requires energy; two; −5. Energy lowering for one bond does not determine the energy change of a whole reaction.'),
  (('How many shared pairs are in H₂?', 'How many lone pairs are on each Cl in Cl₂?', 'How many valence electrons are in HCl in total?'), 'One; three on each; eight. A line represents two electrons, not the complete electron content of the atoms.'),
  (('How many electrons are in Na⁺ if Z=11?', 'Is NaCl a separate molecule in its crystal?', 'What is the simplest neutral ratio of Mg²⁺ to O²⁻?'), '10; no;1:1. Proton counts remain unchanged. A formula unit gives a lattice ratio.'),
  (('Does δ⁻ mean one whole electron has been gained?', 'Is symmetric linear CO₂ polar?', 'Why is bent H₂O polar?'), 'No; no; the electron-density shifts do not cancel. Distinguish a polar bond from a polar molecule.'),
  (('Do H₂O molecules disappear when ice melts?', 'Which carriers are mobile in molten NaCl?', 'Must a conductor be a metal?'), 'No; ions; no, graphite is a counterexample. These answers apply to the specified states and conditions.'),
  (('What formula combines Al³⁺ and Cl⁻?', 'What is the charge sum for MgCl₂?', 'Should Ca₂O₂ be reduced for specified Ca²⁺ and O²⁻?'), 'AlCl₃; +2−2=0; yes, to CaO. Reduce the ion ratio, not an arbitrary molecular formula.'),
  (('What is O’s valency in double-bonded O₂?', 'How many neighbours does each N have in N₂?', 'Does IV for C in CH₄ imply charge +4?'), 'II; one; no. A triple bond has three shared pairs; valency is not charge.'),
  (('How many O atoms are in Al₂(SO₄)₃?', 'How many H atoms are in 4H₂O?', 'Do 2H and H₂ denote identical particles?'), '12;8; no: two separate atoms versus one molecule. Apply group subscripts first, then coefficients to the whole formula.'),
  (('What can white colour alone tell us about D?', 'What formula combines Ca²⁺ and Cl⁻?', 'How many Ca and Cl atoms are in 3CaCl₂?', 'Does a polar bond guarantee molecular polarity?'), 'Appearance only, not composition; CaCl₂;3 Ca and 6 Cl; no, shape matters. Assess each answer separately and correct the reason for any error.'),
 ],
}


def map_svg(number, lang):
    parts = start_svg(LESSONS[number-23][lang][0], NOT_TO_SCALE[lang], height=560)
    for i, line in enumerate(wrap(LESSONS[number-23][lang][0], 62)):
        parts.append(text(40,55+i*37,line,30))
    parts.append(text(40,145,NOT_TO_SCALE[lang],20,'#bedde4'))
    steps = STEPS[lang][number-23]
    for i in range(4):
        x=40+310*i
        parts.append(f'<rect x="{x}" y="195" width="275" height="295" rx="18" fill="#e6f4ef"/>')
        parts.append(text(x+18,244,steps[2*i],21,'#143242'))
        for j,line in enumerate(wrap(steps[2*i+1],17,break_long_words=False)):
            parts.append(text(x+18,320+36*j,line,24,'#143242'))
        if i<3: parts.append(text(x+283,365,'→',28))
    return '\n'.join(parts+['</svg>'])+'\n'


def visual(number, lang):
    labels = {
      'ru': ('Условная модель, не в масштабе','Вымышленные энергии','далеко','связаны','слишком близко','Внешние электроны','общая пара','Решётка продолжается в 3D','Форма и смещения','Линейная: компенсация','Угловая: нет компенсации','Частицы и носители','молекулы','ионы','электроны','сеть атомов','Заряды заданы','Связи ≠ заряд','Скобки → индекс → коэффициент','Данные → модель → предел','Цвет ≠ состав'),
      'kz': ('Шартты модель, масштабсыз','Ойдан берілген энергия','алыс','байланысқан','тым жақын','Сыртқы электрондар','ортақ жұп','Тор 3D кеңістікте жалғасады','Пішін мен ығысулар','Сызықтық: өтеледі','Бұрыштық: өтелмейді','Бөлшек пен тасымалдаушы','молекулалар','иондар','электрондар','атомдар торы','Зарядтар берілген','Байланыс ≠ заряд','Жақша → индекс → коэффициент','Дерек → модель → шек','Түс ≠ құрам'),
      'en': ('Symbolic model, not to scale','Invented energies','far apart','bonded','too close','Valence electrons','shared pair','Lattice continues in 3D','Shape and density shifts','Linear: cancellation','Bent: no cancellation','Particles and carriers','molecules','ions','electrons','atom network','Given ion charges','Bonds ≠ charge','Brackets → subscript → coefficient','Evidence → model → limit','Colour ≠ composition'),
    }[lang]
    parts=start_svg(LESSONS[number-23][lang][0],labels[0])
    parts.append(text(45,55,f'{number:02d} · '+LESSONS[number-23][lang][0],27))
    parts.append(text(45,110,labels[0],22,'#bedde4'))
    def label(x,y,s,size=38): parts.append(text(x,y,s,size,anchor='middle'))
    def line(x,y,x2,y2): parts.append(f'<path d="M{x},{y} L{x2},{y2}" stroke="#8dd8d4" stroke-width="5"/>')
    def atom(x,y,s,color='#d4eee7',size=35):
        parts.append(f'<circle cx="{x}" cy="{y}" r="51" fill="{color}"/>')
        parts.append(text(x,y+12,s,size,'#143242','middle'))
    if number==23:
        label(640,190,labels[1],33)
        for x,phrase,val in ((240,labels[2],'0'),(640,labels[3],'−4'),(1040,labels[4],'+6')):
            label(x,310,phrase,30);label(x,440,val,80)
        label(640,610,'0 − (−4) = 4',48)
    elif number==24:
        label(640,190,labels[5],34)
        label(640,280,'H : H    →    H—H',55)
        # Two Cl valence shells: one shared pair plus six lone pairs.
        for x in (480,800):
            label(x,460,'Cl',64)
            for dx,dy in ((-14,-92),(14,-92),(-14,48),(14,48),(-105,-20),(-105,8)) if x==480 else ((-14,-92),(14,-92),(-14,48),(14,48),(105,-20),(105,8)):
                parts.append(f'<circle data-electron="lone" cx="{x+dx}" cy="{460+dy}" r="7" fill="#8dd8d4"/>')
        for y in (422,450): parts.append(f'<circle data-electron="shared" cx="640" cy="{y}" r="7" fill="#ffd49c"/>')
        label(640,610,'2 + 6 + 6 = 14 e',43)
        label(640,668,labels[6]+' = 2 e',28)
    elif number==25:
        for row in range(3):
            for col in range(6):
                positive=(row+col)%2==0
                atom(240+col*160,245+row*120,'Na⁺' if positive else 'Cl⁻','#f0c18b' if positive else '#98d5d8',29)
        label(640,610,'Na⁺ : Cl⁻ = 1 : 1',43)
        label(640,670,labels[7],27)
    elif number==26:
        label(640,180,labels[8],34)
        label(350,300,'O ← C → O',55)
        label(350,405,labels[9],29)
        atom(950,285,'O');atom(820,415,'H');atom(1080,415,'H')
        line(840,390,915,315);line(1060,390,985,315)
        label(950,495,labels[10],28)
        label(350,565,'CO₂',58);label(950,565,'H₂O',58)
        label(640,665,'δ⁺ ≠ +1     δ⁻ ≠ −1',39)
    elif number==27:
        label(640,185,labels[11],32)
        for x,s,name,carrier in ((190,'H₂O',labels[12],'—'),(490,'NaCl',labels[13],labels[13]),(790,'Cu','Cu',labels[14]),(1090,'C',labels[15],'—')):
            label(x,300,s,53)
            for j,t in enumerate(wrap(name,15)):label(x,395+j*32,t,26)
            label(x,525,carrier,26)
        # Distinguish state-specific carriers from the formula labels above.
        state={'ru':'NaCl: расплав · Cu: твёрдое','kz':'NaCl: балқыма · Cu: қатты','en':'NaCl: molten · Cu: solid'}[lang]
        label(640,640,state,30)
    elif number==28:
        label(640,185,labels[16],34)
        label(640,285,'Al³⁺  +  O²⁻',62)
        label(640,400,'2 × (+3) + 3 × (−2) = 0',45)
        label(640,535,'Al₂O₃',85)
        label(640,650,'2 : 3',49)
    elif number==29:
        label(640,185,labels[17],36)
        for x,s,count in ((230,'H—H','I'),(640,'O=O','II'),(1050,'N≡N','III')):
            label(x,330,s,55);label(x,475,count,75)
        label(640,650,'H₂ · O₂ · N₂ : 0',42)
    elif number==30:
        label(640,185,labels[18],32)
        label(640,300,'2Ca(OH)₂',78)
        label(640,420,'1 Ca · 2 O · 2 H',47)
        label(640,505,'↓ × 2',44)
        label(640,610,'2 Ca · 4 O · 4 H',55)
    else:
        label(640,185,labels[19],34)
        for x,s in ((240,'NaCl'),(640,'H₂O'),(1040,'CO₂')):
            label(x,325,s,58)
        label(640,465,'CaCl₂: (+2) + 2 × (−1) = 0',43)
        label(640,535,'2CaCl₂: Ca × 2; Cl × 4',42)
        label(640,650,labels[20],40)
    return '\n'.join(parts+['</svg>'])+'\n'


def render(number, lang):
    row=LESSONS[number-23]
    title,summary,*sections=row[lang]
    labels=LABELS[lang]
    out=[f'# {title}',f'_Лид (summary):_ **{summary}**',f'## {HEAD[lang][0]}',
         f'![{title}](/static/course/chemistry/map-{number:02d}-{lang}.svg)',
         f'## {labels[0]}',f'![{NOT_TO_SCALE[lang]}](/static/course/chemistry/visuals/{number:02d}-bonds-{lang}.svg)',
         f'*{NOT_TO_SCALE[lang]}*', labels[5]]
    for heading,section in zip(HEAD[lang][1:9],sections,strict=True):
        out.extend([f'## {heading}',section])
    questions,answer=QUIZ[lang][number-23]
    out.extend([f'## {labels[1]}','\n'.join(f'{i}. {q}' for i,q in enumerate(questions,1)),
                f'## {labels[2]}',answer,f'## {labels[3]}',labels[6],
                f'## {HEAD[lang][10]}',RETURN[lang],f'## {labels[4]}',
                '\n'.join(f'- [{SOURCES[key][0]}]({SOURCES[key][1]})' for key in SOURCE_KEYS[number-23])])
    previous=PREVIOUS[-1] if number==23 else LESSONS[number-24]
    out.append(f'[{previous[lang][0]}](/read/{slug(number-1,previous["stem"])}?lang={lang})')
    if number<31:
        following=LESSONS[number-22]
        out.append(f'[{following[lang][0]}](/read/{slug(number+1,following["stem"])}?lang={lang})')
    return '\n\n'.join(out)+'\n'


def main():
    (MAP_DIR/'visuals').mkdir(parents=True,exist_ok=True)
    for number,row in enumerate(LESSONS,23):
        for lang in LANGS:
            suffix='' if lang=='ru' else '-'+lang
            (LESSON_DIR/f'{row["stem"]}{suffix}.md').write_text(render(number,lang),encoding='utf-8')
            (MAP_DIR/f'map-{number:02d}-{lang}.svg').write_text(map_svg(number,lang),encoding='utf-8')
            (MAP_DIR/f'visuals/{number:02d}-bonds-{lang}.svg').write_text(visual(number,lang),encoding='utf-8')


if __name__=='__main__': main()
