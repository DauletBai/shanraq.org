#!/usr/bin/env python3
"""Render chemistry block six: quantities, original diagrams and retrieval maps."""
from textwrap import wrap
try:
    from .chemistry_quantities_data import LESSONS
    from .chemistry_reactions_data import LESSONS as PREVIOUS
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from .build_chemistry_bonds import LABELS
    from .build_chemistry_elements import start_svg, text
    from .build_chemistry_particles import slug
except ImportError:
    from chemistry_quantities_data import LESSONS
    from chemistry_reactions_data import LESSONS as PREVIOUS
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from build_chemistry_bonds import LABELS
    from build_chemistry_elements import start_svg, text
    from build_chemistry_particles import slug

MODEL = {'ru':'Учебная расчётная модель · величины и частицы не в масштабе',
         'kz':'Оқу есептік моделі · шамалар мен бөлшектер масштабта емес',
         'en':'Teaching calculation model · quantities and particles not to scale'}
SAFETY = {'ru':'Только бумага, калькулятор и готовые данные. Не проводите реакции, не получайте газы, не нагревайте вещества и не смешивайте бытовую химию.',
          'kz':'Тек қағаз, калькулятор және дайын деректер. Реакция жүргізбеңіз, газ алмаңыз, заттарды қыздырмаңыз және тұрмыстық химикаттарды араластырмаңыз.',
          'en':'Use only paper, a calculator and supplied data. Do not carry out reactions, generate gases, heat substances or mix household chemicals.'}
SOURCES = (
    ('BIPM: SI definition of the mole','https://www.bipm.org/en/si-base-units/mole'),
    ('OpenStax: formula mass and the mole','https://openstax.org/books/chemistry-2e/pages/3-1-formula-mass-and-the-mole-concept'),
    ('OpenStax: composition and formulas','https://openstax.org/books/chemistry-2e/pages/3-2-determining-empirical-and-molecular-formulas'),
    ('OpenStax: reaction stoichiometry','https://openstax.org/books/chemistry-2e/pages/4-3-reaction-stoichiometry'),
    ('OpenStax: limiting reactants and yields','https://openstax.org/books/chemistry-2e/pages/4-4-reaction-yields'),
)
SOURCE_KEYS = ((0,1),(1,),(0,1),(2,),(3,),(3,),(4,),(4,),(3,4))
STEPS = {
 'ru': [('ЧАСТИЦА','Назови вид','МОЛЬ','n = N / Nₐ','СЧЁТ','N = nNₐ','ИНДЕКС','Учти состав'),
        ('ФОРМУЛА','Посчитай атомы','ВКЛАДЫ','Сложи массы','ЕДИНИЦЫ','M: g/mol','ПРОВЕРКА','Mᵣ без единицы'),
        ('МАССА','Переведи в g','КОЛИЧЕСТВО','n = m / M','ЧАСТИЦЫ','N = nNₐ','ПРОВЕРКА','Назови частицу'),
        ('СОСТАВ','k × Aᵣ(E)','ДОЛЯ','w = kAᵣ / Mᵣ','ПРОЦЕНТ','w × 100%','ПРОВЕРКА','Сумма долей = 1'),
        ('УРАВНЕНИЕ','Проверь атомы','КОЭФФИЦИЕНТЫ','a : b','КОЛИЧЕСТВА','n(B) = n(A)b/a','ГРАНИЦА','Не граммы'),
        ('МАССА','m исходного','МОЛИ','n = m / M','ОТНОШЕНИЕ','n продукта','МАССА','m = nM'),
        ('ИСХОДНО','Все реагенты','СРАВНИ','n / ν','МИНИМУМ','Ограничивает','ОСТАТОК','Начало − расход'),
        ('ТЕОРИЯ','Уравнение и предел','ФАКТ','Чистый продукт','ВЫХОД','η = факт/теория','ПРОВЕРКА','Вода и примеси'),
        ('ЧИСТОТА','Масса реагента','ПРЕДЕЛ','Минимум n / ν','ПРОДУКТ','Поправка на состав','ВЫХОД','Согласуй основания')],
 'kz': [('БӨЛШЕК','Түрін ата','МОЛЬ','n = N / Nₐ','САНАҚ','N = nNₐ','ИНДЕКС','Құрамды ескер'),
        ('ФОРМУЛА','Атомдарды сана','ҮЛЕСТЕР','Массаларды қос','БІРЛІК','M: g/mol','ТЕКСЕРУ','Mᵣ өлшемсіз'),
        ('МАССА','g-ға айналдыр','ЗАТ МӨЛШЕРІ','n = m / M','БӨЛШЕКТЕР','N = nNₐ','ТЕКСЕРУ','Бөлшекті ата'),
        ('ҚҰРАМ','k × Aᵣ(E)','ҮЛЕС','w = kAᵣ / Mᵣ','ПАЙЫЗ','w × 100%','ТЕКСЕРУ','Үлестер қосындысы 1'),
        ('ТЕҢДЕУ','Атомды тексер','КОЭФФИЦИЕНТ','a : b','ЗАТ МӨЛШЕРІ','n(B) = n(A)b/a','ШЕК','Грамм емес'),
        ('МАССА','Бастапқы m','МОЛЬ','n = m / M','ҚАТЫНАС','Өнімнің n мәні','МАССА','m = nM'),
        ('БАСТАПҚЫ','Барлық реагент','САЛЫСТЫР','n / ν','ЕҢ КІШІ','Шектейді','ҚАЛДЫҚ','Бастапқы − шығын'),
        ('ТЕОРИЯ','Теңдеу мен шек','НАҚТЫ','Таза өнім','ШЫҒЫМ','η = нақты/теория','ТЕКСЕРУ','Су мен қоспалар'),
        ('ТАЗАЛЫҚ','Реагент массасы','ШЕК','Ең кіші n / ν','ӨНІМ','Құрамды ескер','ШЫҒЫМ','Негізді сәйкестендір')],
 'en': [('ENTITY','Name its kind','MOLE','n = N / Nₐ','COUNT','N = nNₐ','SUBSCRIPT','Account for formula'),
        ('FORMULA','Count atoms','CONTRIBUTIONS','Add masses','UNITS','M: g/mol','CHECK','Mᵣ is dimensionless'),
        ('MASS','Convert to g','AMOUNT','n = m / M','ENTITIES','N = nNₐ','CHECK','Name the entity'),
        ('COMPOSITION','k × Aᵣ(E)','FRACTION','w = kAᵣ / Mᵣ','PERCENT','w × 100%','CHECK','Fractions sum to 1'),
        ('EQUATION','Check atoms','COEFFICIENTS','a : b','AMOUNTS','n(B) = n(A)b/a','LIMIT','Not grams'),
        ('MASS','Starting m','MOLES','n = m / M','RATIO','Product n','MASS','m = nM'),
        ('INITIAL','All reactants','COMPARE','n / ν','MINIMUM','Sets the limit','REMAINDER','Initial − consumed'),
        ('THEORY','Equation and limit','ACTUAL','Pure product','YIELD','η = actual/theory','CHECK','Water and impurities'),
        ('PURITY','Reactant mass','LIMIT','Minimum n / ν','PRODUCT','Correct composition','YIELD','Consistent basis')],
}

QUIZ = {
 'ru': [
  ('В 0,10 mol H₂O сколько mol атомов H? В 0,30 mol CaCl₂ сколько mol Cl⁻ по составу?', '0,20 mol атомов H; 0,60 mol ионов Cl⁻. Всегда называйте частицу.'),
  ('Найдите M(Na₂O) и M(CaCl₂) по данным блока. Меняется ли M воды при удвоении порции?', '62 и 111 g/mol; нет. Коэффициент меняет количество, а не состав вещества.'),
  ('Найдите n для 4,4 г CO₂ и m для 0,25 mol H₂O. Как перейти от mg к g?', '0,10 mol; 4,5 г; разделить число миллиграммов на 1000.'),
  ('Какова массовая доля O в CaCO₃? Сколько O по составу в 25 г CaCO₃?', '48%; 12 г. Это масса элемента в соединении, а не свободного O₂.'),
  ('Для 2H₂ + O₂ → 2H₂O сколько O₂ требуется для 0,80 mol H₂? Что сравнивают коэффициенты?', '0,40 mol; количества веществ в молях или числа соответствующих частиц, а не массы в граммах.'),
  ('Сколько CO₂ теоретически образуется из 25 г CaCO₃ при полном превращении? Какая масса CaO остаётся?', '11 г CO₂ и 14 г CaO; 25=11+14 в данной модели.'),
  ('Для 2H₂ + O₂ → 2H₂O даны 1 mol H₂ и 1 mol O₂. Что ограничивает? Сколько O₂ останется?', 'H₂; остаётся 0,50 mol O₂, образуется 1 mol H₂O.'),
  ('Теория 12 г, чистый продукт 9 г: выход? Смесь 10 г с 90% продукта при той же теории: выход?', '75% в обоих случаях: во втором сначала 10×0,90=9 г чистого продукта.'),
  ('Почему в основном деле HCl ограничивает? Являются ли 8,00 г чистой солью? Что означает 80%?', '0,150/2 меньше 0,100/1; нет, нужна доля 87,75%; это выход чистого NaCl относительно теоретического NaCl.')],
 'kz': [
  ('0,10 mol H₂O ішінде неше mol H атомы бар? 0,30 mol CaCl₂ құрамына неше mol Cl⁻ сәйкес келеді?', '0,20 mol H атомы; 0,60 mol Cl⁻ ионы. Бөлшек түрін әрқашан атаңыз.'),
  ('Блок деректерімен M(Na₂O), M(CaCl₂) табыңыз. Су порциясы екі еселенсе, M өзгере ме?', '62 және 111 g/mol; жоқ. Коэффициент зат құрамын емес, мөлшерін өзгертеді.'),
  ('4,4 г CO₂ үшін n, 0,25 mol H₂O үшін m табыңыз. mg-нан g-ға қалай өтеміз?', '0,10 mol; 4,5 г; миллиграмм санын 1000-ға бөлеміз.'),
  ('CaCO₃ ішіндегі O массалық үлесі қандай? 25 г CaCO₃ құрамына қанша O сәйкес келеді?', '48%; 12 г. Бұл бос O₂ емес, қосылыстағы элемент массасы.'),
  ('2H₂ + O₂ → 2H₂O үшін 0,80 mol H₂-ге қанша O₂ керек? Коэффициент нені салыстырады?', '0,40 mol; граммен массаны емес, мольмен зат мөлшерін не сәйкес бөлшек сандарын.'),
  ('25 г CaCO₃ толық айналғанда теориялық қанша CO₂ түзіледі? Қанша CaO қалады?', '11 г CO₂ және 14 г CaO; осы модельде 25=11+14.'),
  ('2H₂ + O₂ → 2H₂O үшін 1 mol H₂ және 1 mol O₂ бар. Қайсысы шектейді? Қанша O₂ қалады?', 'H₂; 0,50 mol O₂ қалады, 1 mol H₂O түзіледі.'),
  ('Теория 12 г, таза өнім 9 г: шығым? Сол теориямен 90% өнімі бар 10 г қоспаның шығымы?', 'Екеуінде де 75%: екіншісінде алдымен 10×0,90=9 г таза өнім.'),
  ('Негізгі істе HCl неге шектейді? 8,00 г таза тұз ба? 80% нені білдіреді?', '0,150/2 саны 0,100/1-ден кіші; жоқ, 87,75% үлесі керек; бұл таза NaCl-дың теориялық NaCl-ға қатысты шығымы.')],
 'en': [
  ('How many mol of H atoms are in 0.10 mol H₂O? How many mol of Cl⁻ correspond to 0.30 mol CaCl₂?', '0.20 mol of H atoms; 0.60 mol of Cl⁻ ions. Always specify the entity.'),
  ('Find M(Na₂O) and M(CaCl₂) from the supplied data. Does doubling a water portion change M?', '62 and 111 g/mol; no. A coefficient changes amount, not composition.'),
  ('Find n for 4.4 g CO₂ and m for 0.25 mol H₂O. How do we convert mg to g?', '0.10 mol; 4.5 g; divide the milligram number by 1000.'),
  ('What is the O mass fraction in CaCO₃? What O mass corresponds to 25 g CaCO₃?', '48%; 12 g. This is element mass in the compound, not free O₂.'),
  ('For 2H₂ + O₂ → 2H₂O, how much O₂ is needed for 0.80 mol H₂? What do coefficients compare?', '0.40 mol; mole amounts or corresponding entity counts, not masses in grams.'),
  ('What CO₂ mass theoretically forms from complete conversion of 25 g CaCO₃? What CaO mass remains?', '11 g CO₂ and 14 g CaO; 25=11+14 in this model.'),
  ('For 2H₂ + O₂ → 2H₂O, we have 1 mol H₂ and 1 mol O₂. Which limits? How much O₂ remains?', 'H₂ limits; 0.50 mol O₂ remains and 1 mol H₂O forms.'),
  ('Theory is 12 g and pure product 9 g: yield? What about a 10 g mixture with 90% product and the same theory?', '75% in both: first correct the second sample to 10×0.90=9 g of pure product.'),
  ('Why does HCl limit the main case? Is 8.00 g pure salt? What does 80% mean?', '0.150/2 is less than 0.100/1; no, apply the 87.75% fraction; 80% is pure NaCl yield relative to theoretical NaCl.')],
}

# Numeric content uses the same units across all languages; prose is localized.
ROWS = {
 41:[('H₂O','1 mol','2 mol H + 1 mol O'),('H₂O','0.5 mol','N ≈ 3.01 × 10²³'),('CaCl₂','1 mol','2 mol Cl⁻')],
 42:[('H₂O','2 × 1 + 16','18 g/mol'),('CO₂','12 + 2 × 16','44 g/mol'),('Ca(OH)₂','40 + 2 × (16 + 1)','74 g/mol')],
 43:[('H₂O','9.0 g / (18 g/mol)','0.50 mol'),('H₂O','0.50 mol × Nₐ','3.01 × 10²³'),('CO₂','0.20 mol × 44 g/mol','8.8 g')],
 44:[('H₂O: H','2 / 18','≈ 11.1%'),('H₂O: O','16 / 18','≈ 88.9%'),('CO₂: C','12 / 44','≈ 27.3%')],
 45:[('H₂','2 mol','4 mol'),('O₂','1 mol','2 mol'),('H₂O','2 mol','4 mol')],
 46:[('CaCO₃','10.0 g / 100 g/mol','0.100 mol'),('CO₂','0.100 mol × 44 g/mol','4.40 g'),('CaO','0.100 mol × 56 g/mol','5.60 g')],
 47:[('H₂','3 / 2','1.5 mol'),('O₂','2 / 1','2 mol'),('O₂','2 − 1.5','0.5 mol')],
 48:[('NaCl','4.68 g / 5.85 g','80%'),('10.0 g × 0.75','7.50 g','75%'),('5.00 g × 0.90','4.50 g','4.50 / 6.00 = 75%')],
 49:[('Na₂CO₃','13.25 g × 0.80','10.60 g = 0.100 mol'),('HCl','0.150 mol / 2','0.0750 mol'),('NaCl','7.020 g / 8.775 g','80.0%')],
}


def map_svg(number, lang):
    parts=start_svg(LESSONS[number-41][lang][0],MODEL[lang],height=570)
    for i,line in enumerate(wrap(LESSONS[number-41][lang][0],60)):
        parts.append(text(40,55+i*37,line,30))
    parts.append(text(40,155,MODEL[lang],20,'#bedde4'))
    labels=STEPS[lang][number-41]
    for i in range(4):
        x=40+310*i
        parts.append(f'<rect x="{x}" y="200" width="275" height="290" rx="18" fill="#e6f4ef"/>')
        parts.append(text(x+18,250,labels[2*i],20,'#143242'))
        for j,line in enumerate(wrap(labels[2*i+1],18,break_long_words=False)):
            parts.append(text(x+18,320+35*j,line,23,'#143242'))
        if i<3:parts.append(text(x+283,365,'→',28))
    return '\n'.join(parts+['</svg>'])+'\n'


def visual(number, lang):
    parts=start_svg(LESSONS[number-41][lang][0],MODEL[lang],height=760)
    for i,line in enumerate(wrap(f'{number} · '+LESSONS[number-41][lang][0],65)):
        parts.append(text(40,55+i*37,line,29))
    parts.append(text(40,155,MODEL[lang],20,'#bedde4'))
    if number==41:
        # One symbolic H2O unit and its atom inventory, never a literal mole.
        for element,x,y,r,color in [('O',300,365,55,'#98d5d8'),('H',220,440,34,'#f0c18b'),('H',380,440,34,'#f0c18b')]:
            parts.append(f'<circle data-element="{element}" cx="{x}" cy="{y}" r="{r}" fill="{color}"/>')
            parts.append(text(x,y+10,element,30,'#143242','middle'))
        labels={'ru':('Одна условная молекула','1 mol молекул H₂O','2 mol атомов H','1 mol атомов O'),
                'kz':('Бір шартты молекула','1 mol H₂O молекуласы','2 mol H атомы','1 mol O атомы'),
                'en':('One symbolic molecule','1 mol of H₂O molecules','2 mol of H atoms','1 mol of O atoms')}[lang]
        parts.append(text(300,550,labels[0],26,anchor='middle'))
        for i,label in enumerate(labels[1:]):parts.append(text(650,310+90*i,label,32))
        parts.append(text(640,680,'Nₐ = 6.02214076 × 10²³ mol⁻¹',36,anchor='middle'))
    else:
        for i,row in enumerate(ROWS[number]):
            y=200+i*168
            parts.append(f'<rect x="40" y="{y}" width="1200" height="148" rx="18" fill="#173e4d"/>')
            parts.append(text(75,y+54,row[0],30))
            parts.append(text(450,y+54,row[1],30))
            parts.append(text(640,y+112,'→ '+row[2],32,anchor='middle'))
    return '\n'.join(parts+['</svg>'])+'\n'


def table(number, lang):
    headers={'ru':('Объект / шаг','Дано / расчёт','Результат'),
             'kz':('Нысан / қадам','Берілгені / есеп','Нәтиже'),
             'en':('Entity / step','Given / calculation','Result')}[lang]
    if number==45:
        headers={'ru':('Вещество','Отношение: пример A','Удвоено: пример B'),
                 'kz':('Зат','Қатынас: A мысалы','Екі еселенген: B'),
                 'en':('Substance','Ratio: example A','Doubled: example B')}[lang]
    return '| '+' | '.join(headers)+' |\n|---|---|---|\n'+'\n'.join('| '+' | '.join(row)+' |' for row in ROWS[number])


def render(number, lang):
    row=LESSONS[number-41];title,summary,*sections=row[lang];labels=LABELS[lang]
    out=[f'# {title}',f'_Лид (summary):_ **{summary}**',f'## {HEAD[lang][0]}',
         f'![{title}](/static/course/chemistry/map-{number:02d}-{lang}.svg)',
         f'## {labels[0]}',f'![{MODEL[lang]}](/static/course/chemistry/visuals/{number:02d}-quantities-{lang}.svg)',SAFETY[lang]]
    for i,(heading,section) in enumerate(zip(HEAD[lang][1:9],sections,strict=True)):
        out.extend([f'## {heading}',section])
        if i==3:out.append(table(number,lang))
    question,answer=QUIZ[lang][number-41]
    out.extend([f'## {labels[1]}',question,f'## {labels[2]}',answer,
                f'## {HEAD[lang][10]}',RETURN[lang],f'## {labels[4]}',
                '\n'.join(f'- [{SOURCES[k][0]}]({SOURCES[k][1]})' for k in SOURCE_KEYS[number-41])])
    previous=PREVIOUS[-1] if number==41 else LESSONS[number-42]
    out.append(f'[{previous[lang][0]}](/read/{slug(number-1,previous["stem"])}?lang={lang})')
    if number<49:
        following=LESSONS[number-40]
        out.append(f'[{following[lang][0]}](/read/{slug(number+1,following["stem"])}?lang={lang})')
    return '\n\n'.join(out)+'\n'


def bridge_body(lang):
    try:
        from .build_chemistry_reactions import render as previous_render
    except ImportError:
        from build_chemistry_reactions import render as previous_render
    row=LESSONS[0]
    return previous_render(40,lang).split('\n\n',2)[2].rstrip()+f'\n\n[{row[lang][0]}](/read/{slug(41,row["stem"])}?lang={lang})\n'


def main():
    (MAP_DIR/'visuals').mkdir(parents=True,exist_ok=True)
    for number,row in enumerate(LESSONS,41):
        for lang in LANGS:
            suffix='' if lang=='ru' else '-'+lang
            (LESSON_DIR/f'{row["stem"]}{suffix}.md').write_text(render(number,lang),encoding='utf-8')
            (MAP_DIR/f'map-{number:02d}-{lang}.svg').write_text(map_svg(number,lang),encoding='utf-8')
            (MAP_DIR/f'visuals/{number:02d}-quantities-{lang}.svg').write_text(visual(number,lang),encoding='utf-8')


if __name__=='__main__':main()
