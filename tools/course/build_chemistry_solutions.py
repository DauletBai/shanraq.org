#!/usr/bin/env python3
"""Generate chemistry block seven with localized original diagrams."""
from textwrap import wrap
try:
    from .chemistry_solutions_data import LESSONS
    from .chemistry_quantities_data import LESSONS as PREVIOUS
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from .build_chemistry_elements import start_svg, text
    from .build_chemistry_particles import slug
    from .build_chemistry_bonds import LABELS
except ImportError:
    from chemistry_solutions_data import LESSONS
    from chemistry_quantities_data import LESSONS as PREVIOUS
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from build_chemistry_elements import start_svg, text
    from build_chemistry_particles import slug
    from build_chemistry_bonds import LABELS

MODEL = {'ru':'Учебная модель · частицы и размеры не в масштабе',
         'kz':'Оқу моделі · бөлшектер мен өлшемдер масштабта емес',
         'en':'Teaching model · particles and sizes are not to scale'}
SAFETY = {
'ru':'Работа только с бумагой, калькулятором и данными. Не нагревайте, не пробуйте и не нюхайте образцы, не смешивайте бытовую химию и не проводите электрические опыты. Прозрачность и обычная фильтрация не доказывают, что вода пригодна для питья.',
'kz':'Тек қағаз, калькулятор және деректермен жұмыс. Үлгілерді қыздырмаңыз, дәмін татпаңыз, иіскемеңіз, тұрмыстық химияны араластырмаңыз және электрлік тәжірибе жасамаңыз. Мөлдірлік пен қарапайым сүзу судың ішуге жарамдылығын дәлелдемейді.',
'en':'Use paper, a calculator and supplied data only. Do not heat, taste or smell samples, mix household chemicals or perform electrical experiments. Clarity and ordinary filtration do not establish drinking safety.'}
SOURCES = [
('OpenStax: geometry and polarity','https://openstax.org/books/chemistry-2e/pages/7-6-molecular-structure-and-polarity'),
('OpenStax: dissolution','https://openstax.org/books/chemistry-2e/pages/11-1-the-dissolution-process'),
('OpenStax: solubility','https://openstax.org/books/chemistry-2e/pages/11-3-solubility'),
('OpenStax: mass percentage','https://openstax.org/books/chemistry-2e/pages/3-4-other-units-for-solution-concentrations'),
('OpenStax: molarity and dilution','https://openstax.org/books/chemistry-2e/pages/3-3-molarity')]
SOURCE_KEYS=[(0,1),(1,),(1,),(2,),(3,),(3,),(4,),(3,4),(2,3,4)]
STEPS={
'ru':[
('СВЯЗЬ','Частичный заряд','ФОРМА','Угловая H₂O','МОЛЕКУЛА','Полярна','ИОН','Ориентация воды'),
('СОСТАВ','Смесь веществ','РАСТВОР','Однородный','ВЗВЕСЬ','Крупные частицы','ФИЛЬТР','Соль проходит'),
('САХАР','Молекулы','NaCl','Ионы в воде','СОХРАНЕНИЕ','Масса и заряд','МОДЕЛЬ','Не фотография'),
('УСЛОВИЯ','Температура','S','На 100 г воды','ПРЕДЕЛ','Равновесие','СКОРОСТЬ','Другой вопрос'),
('ВЕЩЕСТВО','Масса компонента','РАСТВОР','Полная масса','ДОЛЯ','w = m / m','ПРОЦЕНТ','w × 100%'),
('ЦЕЛЬ','Масса раствора','ДОЛЯ','Процент / 100','ВЕЩЕСТВО','w × m раствора','ВОДА','Разность масс'),
('МАССА','n = m / M','ОБЪЁМ','Литры раствора','КОНЦЕНТРАЦИЯ','c = n / V','ПРОВЕРКА','mol/L'),
('БАЛАНС','n сохраняется','РАЗБАВЛЕНИЕ','c₁V₁ = c₂V₂','СМЕШИВАНИЕ','Сложи массы','ГРАНИЦА','Нет потерь'),
('ДАННЫЕ','Что известно?','БАЛАНС','Компонент и смесь','ЕДИНИЦЫ','Процент или mol/L','ВЫВОД','Чего не знаем?')],
'kz':[
('БАЙЛАНЫС','Ішінара заряд','ПІШІН','Бұрышты H₂O','МОЛЕКУЛА','Полярлы','ИОН','Судың бағыты'),
('ҚҰРАМ','Заттар қоспасы','ЕРІТІНДІ','Біртекті','СУСПЕНЗИЯ','Ірі бөлшектер','СҮЗГІ','Тұз өтеді'),
('ҚАНТ','Молекулалар','NaCl','Судағы иондар','САҚТАЛУ','Масса мен заряд','МОДЕЛЬ','Фото емес'),
('ЖАҒДАЙ','Температура','S','100 г суға','ШЕК','Тепе-теңдік','ЖЫЛДАМДЫҚ','Бөлек сұрақ'),
('ЗАТ','Құрамдас массасы','ЕРІТІНДІ','Толық масса','ҮЛЕС','w = m / m','ПАЙЫЗ','w × 100%'),
('МАҚСАТ','Ерітінді массасы','ҮЛЕС','Пайыз / 100','ЗАТ','w × ерітінді m','СУ','Масса айырмасы'),
('МАССА','n = m / M','КӨЛЕМ','Ерітінді литрі','КОНЦЕНТРАЦИЯ','c = n / V','ТЕКСЕРУ','mol/L'),
('БАЛАНС','n сақталады','СҰЙЫЛТУ','c₁V₁ = c₂V₂','АРАЛАСТЫРУ','Массаларды қос','ШЕК','Шығын жоқ'),
('ДЕРЕК','Не белгілі?','БАЛАНС','Құрамдас, қоспа','БІРЛІК','Пайыз не mol/L','ҚОРЫТЫНДЫ','Не белгісіз?')],
'en':[
('BOND','Partial charge','SHAPE','Bent H₂O','MOLECULE','Polar','ION','Water orientation'),
('COMPOSITION','Mixture','SOLUTION','Homogeneous','SUSPENSION','Larger particles','FILTER','Salt passes'),
('SUGAR','Molecules','NaCl','Ions in water','CONSERVE','Mass and charge','MODEL','Not a photograph'),
('CONDITIONS','Temperature','S','Per 100 g water','LIMIT','Equilibrium','RATE','A different question'),
('SOLUTE','Component mass','SOLUTION','Total mass','FRACTION','w = m / m','PERCENT','w × 100%'),
('TARGET','Solution mass','FRACTION','Percent / 100','SOLUTE','w × solution m','WATER','Mass difference'),
('MASS','n = m / M','VOLUME','Litres of solution','CONCENTRATION','c = n / V','CHECK','mol/L'),
('BALANCE','n conserved','DILUTION','c₁V₁ = c₂V₂','MIXING','Add masses','LIMIT','No losses'),
('DATA','What is known?','BALANCE','Solute and mixture','UNITS','Percent or mol/L','CLAIM','What is unknown?')]
}
ROWS = {
50:[('H₂O','δ−: O; δ+: H','q = 0'),('CO₂','O=C=O','μ = 0'),('Na⁺ / Cl⁻','H₂O','δ− / δ+')],
51:[('NaCl','10 g + 90 g','100 g'),('20 g','20 / 100','2 g NaCl'),('15 g','15 / 60','1.5 g')],
52:[('NaCl','0.10 mol','0.10 mol Na⁺ + 0.10 mol Cl⁻'),('NaCl','0.20 mol','0.40 mol ions'),('5 g + 95 g','100 g','q = 0')],
53:[('20 °C','S = 30 g / 100 g H₂O','15 g / 50 g H₂O'),('20 g X','15 g + 5 g','65 g + 5 g'),('40 → 20 °C','50 g → 30 g','20 g crystals')],
54:[('12 g + 188 g','12 / 200','6%'),('12 g + 88 g','12 / 100','12%'),('7.5 g + 142.5 g','7.5 / 150','5%')],
55:[('250 g; 8%','0.08 × 250','20 g + 230 g'),('400 g; 5%','0.05 × 400','20 g + 380 g'),('8 g; 5%','8 / 0.05','160 g')],
56:[('5.85 g NaCl','5.85 / 58.5','0.100 mol'),('0.500 L','0.100 / 0.500','0.200 mol/L'),('0.30 mol/L; 0.20 L','0.060 × 58.5','3.51 g')],
57:[('100 → 250 mL','0.50 × 100 / 250','0.20 mol/L'),('100 g; 10% + 100 g; 4%','14 / 200','7%'),('50 g; 10% + 150 g; 4%','11 / 200','5.5%')],
58:[('10 g + 90 g','10 / 100','10%'),('+ 100 g H₂O','10 / 200','5%'),('V = 0.195 L','(10 / 58.5) / 0.195','≈ 0.877 mol/L')]
}
# Localize textual data labels; numerical notation stays consistent across versions.
WORDS={
'ru':{'0.40 mol ions':'0.40 mol ионов','20 g crystals':'20 g кристаллов'},
'kz':{'0.40 mol ions':'0.40 mol ион','20 g crystals':'20 g кристалл'},
'en':{}
}
QUIZ={
'ru':[
('Полярная связь доказывает полярность молекулы? Какой конец H₂O обращён к Na⁺?', 'Нет, важна геометрия; кислородный конец δ−. Молекула воды нейтральна.'),
('В 30 г раствора с составом 10 г соли на 100 г раствора сколько соли? Удалит её бумажный фильтр?', '3 г; обычный бумажный фильтр растворённую соль не удаляет.'),
('Для 0,30 mol NaCl при полной диссоциации сколько mol Na⁺ и Cl⁻? Суммарный заряд?', 'По 0,30 mol каждого вида, всего 0,60 mol ионов; суммарный заряд нулевой.'),
('При S=30 г на 100 г воды и 20 °C в 200 г воды добавили 70 г X. Равновесие достигнуто: сколько растворено и осталось?', '60 г растворены, 10 г твёрдого остатка; жидкий раствор весит 260 г.'),
('9 г вещества + 171 г воды: массовая доля? Что в знаменателе?', '9/180=0,05=5%; масса всего раствора.'),
('Найдите массы для 300 г 4%-го раствора сахара. Затем проверьте ответ.', '12 г сахара и 288 г воды; 12/300=0,04.'),
('0,040 mol в 200 мл готового раствора: c? Достаточно ли вместо этого знать только массу воды?', '0,20 mol/L; нет, нужен объём раствора или дополнительные данные.'),
('К 50 г 12%-го раствора добавили 100 г воды: новая доля? Что сохранилось?', '6 г компонента в 150 г раствора: 4%; масса компонента при отсутствии потерь и реакции.'),
('Почему 5% нельзя заменить 5 mol/L? Какие данные отсутствуют для перехода между ними?', 'Разные величины; нужны молярная масса и объём либо плотность раствора вместе с его массой.')
],
'kz':[
('Полярлы байланыс молекула полярлығын дәлелдей ме? H₂O Na⁺-қа қай жағымен бағытталады?', 'Жоқ, пішін маңызды; δ− бар оттек жағымен. Су молекуласы бейтарап.'),
('100 г ерітіндіде 10 г тұз болса, 30 грамда қанша тұз бар? Қағаз сүзгі оны жоя ма?', '3 г; қарапайым қағаз сүзгі еріген тұзды жоймайды.'),
('0,30 mol NaCl толық диссоциацияланса, Na⁺ пен Cl⁻ мөлшері қандай? Жалпы заряд?', 'Әрқайсысы 0,30 mol, барлығы 0,60 mol ион; жалпы заряд нөл.'),
('20 °C-та S=100 г суға 30 г; 200 г суға 70 г X қосылды. Тепе-теңдікте қаншасы ериді және қалады?', '60 г ериді, 10 г қатты қалдық; сұйық ерітінді массасы 260 г.'),
('9 г зат + 171 г су: массалық үлес? Бөлімде не бар?', '9/180=0,05=5%; бүкіл ерітінді массасы.'),
('300 г 4%-дық қант ерітіндісінің құрамдастар массасын табыңыз. Жауапты тексеріңіз.', '12 г қант және 288 г су; 12/300=0,04.'),
('Дайын ерітіндінің 200 миллилитрінде 0,040 mol бар: c қандай? Тек су массасын білу жеткілікті ме?', '0,20 mol/L; жоқ, ерітінді көлемі не қосымша дерек керек.'),
('50 г 12%-дық ерітіндіге 100 г су қосылды: жаңа үлес? Не сақталды?', '150 г ерітіндідегі 6 г құрамдас: 4%; шығынсыз, реакциясыз құрамдас массасы сақталды.'),
('Неге 5%-ды 5 mol/L-мен алмастырмаймыз? Ауыстыру үшін қандай дерек керек?', 'Шамалар әртүрлі; молярлық масса және көлем не ерітінді массасымен бірге тығыздық керек.')
],
'en':[
('Do polar bonds prove molecular polarity? Which end of H₂O faces Na⁺?', 'No; shape matters. The oxygen end with δ− faces Na⁺. Water is neutral overall.'),
('A solution contains 10 g of salt per 100 g of solution. How much is in 30 g? Does paper filtration remove it?', '3 g; ordinary filter paper does not remove dissolved salt.'),
('For complete dissociation of 0.30 mol of NaCl, find both ion amounts and total charge.', '0.30 mol of each ion, 0.60 mol of ions total; total charge is zero.'),
('At 20 °C, S=30 g per 100 g water. Add 70 g of X to 200 g water. At equilibrium how much is dissolved and solid?', '60 g dissolved and 10 g solid; liquid solution mass is 260 g.'),
('Find mass fraction for 9 g of solute plus 171 g of water. What is the denominator?', '9/180=0.05=5%; total solution mass.'),
('Find component masses for 300 g of 4% sugar solution and check.', '12 g of sugar and 288 g of water; 12/300=0.04.'),
('Find c for 0.040 mol in 200 mL of final solution. Would knowing water mass alone suffice?', '0.20 mol/L; no, final solution volume or additional information is needed.'),
('Add 100 g of water to 50 g of a 12% solution. Find mass fraction and the conserved quantity.', '6 g of component in 150 g gives 4%; component mass is conserved without losses or reaction.'),
('Why cannot 5% be replaced by 5 mol/L? What extra data would allow conversion?', 'Different quantities; molar mass and solution volume, or density together with solution mass, are needed.')
]}
ALT={
'ru':[
'Угловая H₂O: O с δ−, два H с δ+; суммарный заряд ноль. Вода обращена кислородом к Na⁺.',
'10 г соли и 90 г воды дают 100 г однородного раствора; 20 г порции содержат 2 г соли.',
'Модель четырёх Na⁺ и четырёх Cl⁻: числа равны, заряд ноль. Сахароза остаётся молекулами.',
'Для вымышленного X при 20 °C предел 30 г на 100 г воды, при 40 °C — 50 г. Точки не измерены учеником.',
'12 г вещества и 188 г воды дают 200 г раствора и 6% вещества по массе.',
'250 г раствора с долей 8%: 20 г вещества и 230 г воды, сумма 250 г.',
'5,85 г NaCl соответствуют 0,100 моль; в 0,500 л готового раствора концентрация 0,200 моль/л.',
'Количество вещества сохраняется: 100 мл при 0,50 моль/л доведены до 250 мл при 0,20 моль/л.',
'В цепочке 10 г соли, 100 г начального и 200 г конечного раствора доля уменьшается с 10% до 5%.'],
'kz':[
'Бұрышты H₂O: O-да δ−, екі H-та δ+; жалпы заряд нөл. Су Na⁺-қа оттек жағымен бағытталған.',
'10 г тұз және 90 г су 100 г біртекті ерітінді береді; 20 г порцияда 2 г тұз бар.',
'Төрт Na⁺ және төрт Cl⁻ моделі: сандар тең, заряд нөл. Сахароза молекула болып қалады.',
'Ойдан алынған X үшін 20 °C-та шек 100 г суға 30 г, 40 °C-та 50 г. Нүктелерді оқушы өлшемеген.',
'12 г зат және 188 г су 200 г ерітінді мен 6% массалық үлес береді.',
'8% үлесі бар 250 г ерітінді: 20 г зат және 230 г су, қосынды 250 г.',
'5,85 г NaCl 0,100 мольге сәйкес; дайын ерітіндінің 0,500 литрінде концентрация 0,200 моль/л.',
'Зат мөлшері сақталады: 0,50 моль/л концентрациядағы 100 мл соңғы 250 мл-де 0,20 моль/л береді.',
'10 г тұз, бастапқы 100 г және соңғы 200 г ерітінді тізбегінде үлес 10%-дан 5%-ға кемиді.'],
'en':[
'Bent H₂O: oxygen δ−, two hydrogens δ+; total charge zero. Water faces Na⁺ with oxygen.',
'10 g salt plus 90 g water gives 100 g homogeneous solution; a 20 g portion contains 2 g salt.',
'Four Na⁺ and four Cl⁻ ions: equal counts and zero total charge. Sucrose remains molecular.',
'Imaginary X: 30 g per 100 g water at 20 °C and 50 g at 40 °C. These are supplied points, not student measurements.',
'12 g solute plus 188 g water gives 200 g solution and 6% solute by mass.',
'250 g at 8%: 20 g solute and 230 g water, total 250 g.',
'5.85 g NaCl corresponds to 0.100 mol; in 0.500 L of final solution c is 0.200 mol/L.',
'Solute amount conserved: 100 mL at 0.50 mol/L becomes a final 250 mL at 0.20 mol/L.',
'10 g salt in initially 100 g and finally 200 g of solution gives a change from 10% to 5%.']
}

def lines(parts,x,y,value,width=55,size=26,color='#f1f7f2',anchor='start'):
    for j,line in enumerate(wrap(value,width,break_long_words=False,break_on_hyphens=False)):
        parts.append(text(x,y+size*1.35*j,line,size,color,anchor))

def map_svg(number,lang):
    row=LESSONS[number-50]
    p=start_svg(row[lang][0],MODEL[lang],height=600)
    lines(p,40,55,row[lang][0],58,29)
    p.append(text(40,155,MODEL[lang],20,'#bedde4'))
    for i in range(4):
        x=40+i*310
        p.append(f'<rect x="{x}" y="210" width="275" height="300" rx="18" fill="#e6f4ef"/>')
        lines(p,x+18,265,STEPS[lang][number-50][i*2],18,19,'#143242')
        lines(p,x+18,345,STEPS[lang][number-50][i*2+1],17,23,'#143242')
        if i<3:p.append(text(x+282,385,'→',27))
    return '\n'.join(p+['</svg>'])+'\n'

def water(p,x,y,oxygen_toward='left'):
    # Oxygen faces the ion on the left; schematic bent geometry only.
    p.append(f'<line x1="{x}" y1="{y}" x2="{x+65}" y2="{y-45}" stroke="#aecbd4" stroke-width="7"/>')
    p.append(f'<line x1="{x}" y1="{y}" x2="{x+65}" y2="{y+45}" stroke="#aecbd4" stroke-width="7"/>')
    for e,dx,dy,r in [('O',0,0,32),('H',65,-45,22),('H',65,45,22)]:
        p.append(f'<circle data-element="{e}" cx="{x+dx}" cy="{y+dy}" r="{r}" fill="#a1d9dd"/>')
        p.append(text(x+dx,y+dy+8,e,24,'#102b3b','middle'))
    p.append(text(x-13,y-45,'δ−',24))
    p.append(text(x+88,y-45,'δ+',24))

def visual(number,lang):
    row=LESSONS[number-50]
    p=start_svg(row[lang][0],ALT[lang][number-50],height=820)
    lines(p,40,55,f'{number} · '+row[lang][0],60,29)
    p.append(text(40,155,MODEL[lang],20,'#bedde4'))
    if number==50:
        p.append('<circle cx="265" cy="350" r="60" fill="#e9bb80"/>')
        p.append(text(265,362,'Na⁺',38,'#102b3b','middle'))
        water(p,440,350)
        p.append(text(850,350,'H₂O',44,anchor='middle'))
        p.append(text(850,415,'q = 0',36,anchor='middle'))
        lines(p,40,610,ALT[lang][0],76,25)
    elif number==52:
        for i in range(8):
            x=100+(i%4)*155;y=260+(i//4)*160
            ion='Na⁺' if i%2==0 else 'Cl⁻'
            p.append(f'<circle data-ion="{ion}" cx="{x}" cy="{y}" r="48" fill="#a1d9dd"/>')
            p.append(text(x,y+9,ion,27,'#102b3b','middle'))
        p.append(text(910,300,'C₁₂H₂₂O₁₁',36,anchor='middle'))
        p.append(text(910,400,'q = 0',36,anchor='middle'))
        lines(p,40,610,ALT[lang][2],76,25)
    elif number==53:
        p.append('<path d="M180 540 H1120 M180 540 V220" fill="none" stroke="#d3e8eb" stroke-width="4"/>')
        p.append(text(660,605,'T / °C',27,anchor='middle'))
        p.append(text(180,200,'S / (g X / 100 g H₂O)',24))
        for temp,s,x,y in [(20,30,430,400),(40,50,860,280)]:
            p.append(f'<circle data-temperature="{temp}" data-solubility="{s}" cx="{x}" cy="{y}" r="12" fill="#efbc7e"/>')
            p.append(text(x,575,temp,26,anchor='middle'))
            p.append(text(x,y-30,f'({temp}; {s})',30,anchor='middle'))
        lines(p,40,665,ALT[lang][3],80,23)
    else:
        for i,rowdata in enumerate(ROWS[number]):
            y=205+i*150
            p.append(f'<rect x="40" y="{y}" width="1200" height="130" rx="18" fill="#173e4d"/>')
            parts=[WORDS[lang].get(s,s) for s in rowdata]
            lines(p,65,y+42,parts[0]+'  →  '+parts[1],65,26)
            lines(p,65,y+91,parts[2],65,30)
        lines(p,40,725,ALT[lang][number-50],90,21)
    return '\n'.join(p+['</svg>'])+'\n'

def table(number,lang):
    header={'ru':('Объект / шаг','Дано / расчёт','Результат'),
            'kz':('Нысан / қадам','Берілгені / есеп','Нәтиже'),
            'en':('Object / step','Given / calculation','Result')}[lang]
    rows=[[WORDS[lang].get(s,s) for s in r] for r in ROWS[number]]
    return '| '+' | '.join(header)+' |\n|---|---|---|\n'+'\n'.join('| '+' | '.join(r)+' |' for r in rows)

def render(number,lang):
    row=LESSONS[number-50];title,summary,*sections=row[lang];labels=LABELS[lang]
    body=[f'# {title}',f'_Лид (summary):_ **{summary}**',
          f'## {HEAD[lang][0]}',f'![{title}](/static/course/chemistry/map-{number:02d}-{lang}.svg)',
          f'## {labels[0]}',f'![{ALT[lang][number-50]}](/static/course/chemistry/visuals/{number:02d}-solutions-{lang}.svg)',
          SAFETY[lang]]
    for i,(heading,paragraph) in enumerate(zip(HEAD[lang][1:9],sections,strict=True)):
        body.extend([f'## {heading}',paragraph])
        if i==3:body.append(table(number,lang))
    q,a=QUIZ[lang][number-50]
    body.extend([f'## {labels[1]}',q,f'## {labels[2]}',a,
                 f'## {HEAD[lang][10]}',RETURN[lang],f'## {labels[4]}',
                 '\n'.join(f'- [{SOURCES[k][0]}]({SOURCES[k][1]})' for k in SOURCE_KEYS[number-50])])
    prev=PREVIOUS[-1] if number==50 else LESSONS[number-51]
    body.append(f'[{prev[lang][0]}](/read/{slug(number-1,prev["stem"])}?lang={lang})')
    if number<58:
        nxt=LESSONS[number-49]
        body.append(f'[{nxt[lang][0]}](/read/{slug(number+1,nxt["stem"])}?lang={lang})')
    return '\n\n'.join(body)+'\n'

def bridge_body(lang):
    try:from .build_chemistry_quantities import render as previous_render
    except ImportError:from build_chemistry_quantities import render as previous_render
    row=LESSONS[0]
    return previous_render(49,lang).split('\n\n',2)[2].rstrip()+f'\n\n[{row[lang][0]}](/read/{slug(50,row["stem"])}?lang={lang})\n'

def main():
    (MAP_DIR/'visuals').mkdir(parents=True,exist_ok=True)
    for number,row in enumerate(LESSONS,50):
        for lang in LANGS:
            suffix='' if lang=='ru' else '-'+lang
            (LESSON_DIR/f'{row["stem"]}{suffix}.md').write_text(render(number,lang),encoding='utf8')
            (MAP_DIR/f'map-{number:02d}-{lang}.svg').write_text(map_svg(number,lang),encoding='utf8')
            (MAP_DIR/f'visuals/{number:02d}-solutions-{lang}.svg').write_text(visual(number,lang),encoding='utf8')

if __name__=='__main__':main()
