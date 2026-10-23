#!/usr/bin/env python3
"""Render the third chemistry block and its original teaching diagrams."""
from html import escape
from textwrap import wrap

try:
    from .chemistry_elements_data import ELEMENTS, LESSONS
    from .chemistry_particles_data import LESSONS as PREVIOUS
    from .build_chemistry_particles import slug, NOT_TO_SCALE
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
except ImportError:
    from chemistry_elements_data import ELEMENTS, LESSONS
    from chemistry_particles_data import LESSONS as PREVIOUS
    from build_chemistry_particles import slug, NOT_TO_SCALE
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR

VISUAL_DIR = MAP_DIR / 'visuals'
STEPS = {
    'ru': (
        ('СИМВОЛ', 'Латинские буквы', 'Z', 'Число протонов', 'Aᵣ', 'Средняя масса', 'ПРОВЕРЬ', 'Нейтроны неизвестны'),
        ('ПЕРИОД', 'Горизонталь', 'ГРУППА', 'Вертикаль 1–18', 'НАЙДИ', 'Na: 3 / 1', 'ПРЕДЕЛ', 'Не полный каталог'),
        ('УСЛОВИЕ', 'Нейтральный атом', 'СЧЁТ', 'Сумма равна Z', 'Ca', '2 | 8 | 8 | 2', 'ПРЕДЕЛ', 'Первые 20 элементов'),
        ('СЕМЕЙСТВО', 'Сходство свойств', 'СЛОЙ', 'Сравни внешние e', 'He', 'Заполнен двумя e', 'H', 'Не щелочной металл'),
        ('ДАННЫЕ', 'Несколько свойств', 'Al', 'Проводит, пластичен', 'ГРАФИТ', 'Проводит, хрупок', 'ПРЕДЕЛ', 'Не один признак'),
        ('ВНИЗ', 'Li < Na < K', 'ВПРАВО', 'Na > Mg > Al', 'ПРИЧИНА', 'Слои и ядро', 'ПРЕДЕЛ', 'Одно свойство'),
        ('СМЕСЬ', 'Доли дают 100%', 'РАСЧЁТ', '0,75×20 + 0,25×22', 'Aᵣ', '20,5 без единицы', 'ПРЕДЕЛ', 'Не один атом'),
        ('Z', 'Сохраняется', 'ИОН', 'Меняется число e', 'ЗАРЯД', 'p − e', 'ПОЛОЖЕНИЕ', 'Определяет Z'),
        ('КАРТОЧКИ', 'Na и Cl', 'ИОНЫ', 'Na⁺ и Cl⁻', 'БАЛАНС', '+1 − 1 = 0', 'ПРЕДЕЛ', 'Состав дан условием'),
    ),
    'kz': (
        ('ТАҢБА', 'Латын әріптері', 'Z', 'Протон саны', 'Aᵣ', 'Орташа масса', 'ТЕКСЕР', 'Нейтрон белгісіз'),
        ('ПЕРИОД', 'Көлденең жол', 'ТОП', 'Тік баған 1–18', 'ТАП', 'Na: 3 / 1', 'ШЕК', 'Толық тізім емес'),
        ('ШАРТ', 'Бейтарап атом', 'САНАУ', 'Қосынды Z-ке тең', 'Ca', '2 | 8 | 8 | 2', 'ШЕК', 'Алғашқы 20 элемент'),
        ('ОТБАСЫ', 'Қасиет ұқсастығы', 'ҚАБАТ', 'Сыртқы e салыстыр', 'He', 'Екі e-мен толған', 'H', 'Сілтілік металл емес'),
        ('ДЕРЕК', 'Бірнеше қасиет', 'Al', 'Өткізгіш, пластикалы', 'ГРАФИТ', 'Өткізгіш, морт', 'ШЕК', 'Жалғыз белгі емес'),
        ('ТӨМЕН', 'Li < Na < K', 'ОҢҒА', 'Na > Mg > Al', 'СЕБЕП', 'Қабат пен ядро', 'ШЕК', 'Бір қасиет'),
        ('ҚОСПА', 'Үлестер: 100%', 'ЕСЕП', '0,75×20 + 0,25×22', 'Aᵣ', '20,5 бірліксіз', 'ШЕК', 'Жеке атом емес'),
        ('Z', 'Сақталады', 'ИОН', 'e саны өзгереді', 'ЗАРЯД', 'p − e', 'ОРНЫ', 'Z анықтайды'),
        ('КАРТОЧКА', 'Na және Cl', 'ИОНДАР', 'Na⁺ және Cl⁻', 'БАЛАНС', '+1 − 1 = 0', 'ШЕК', 'Құрам шартта берілген'),
    ),
    'en': (
        ('SYMBOL', 'Latin letters', 'Z', 'Proton count', 'Aᵣ', 'Average mass', 'CHECK', 'Neutrons unknown'),
        ('PERIOD', 'Horizontal row', 'GROUP', 'Vertical, 1–18', 'LOCATE', 'Na: 3 / 1', 'LIMIT', 'Not a full catalogue'),
        ('CONDITION', 'Neutral atom', 'COUNT', 'Sum equals Z', 'Ca', '2 | 8 | 8 | 2', 'LIMIT', 'First 20 elements'),
        ('FAMILY', 'Similar properties', 'SHELL', 'Compare outer e', 'He', 'Filled with two e', 'H', 'Not an alkali metal'),
        ('DATA', 'Several properties', 'Al', 'Conducts, malleable', 'GRAPHITE', 'Conducts, brittle', 'LIMIT', 'Not one property'),
        ('DOWN', 'Li < Na < K', 'ACROSS', 'Na > Mg > Al', 'EXPLAIN', 'Shells and nucleus', 'LIMIT', 'One property'),
        ('MIXTURE', 'Fractions total 100%', 'CALCULATE', '0.75×20 + 0.25×22', 'Aᵣ', '20.5, no unit', 'LIMIT', 'Not a single atom'),
        ('Z', 'Unchanged', 'ION', 'Electron count varies', 'CHARGE', 'p − e', 'POSITION', 'Determined by Z'),
        ('CARDS', 'Na and Cl', 'IONS', 'Na⁺ and Cl⁻', 'BALANCE', '+1 − 1 = 0', 'LIMIT', 'Composition supplied'),
    ),
}
QUIZ = {
    'ru': (
        ('Что сообщает Z=17? Можно ли узнать нейтроны из Aᵣ≈35,45? Чем Co отличается от CO?', '17 протонов; нет, нужен изотоп; Co — кобальт, CO — формула угарного газа.'),
        ('Назовите период и группу Cl. Где He? Есть ли в таблице группы 3–12, хотя наш фрагмент их не заполняет?', 'Cl: 3 и 17; He: 1 и 18; да, они есть в полной таблице.'),
        ('Как распределены электроны Ca? Как проверить сумму? Ограничена ли третья оболочка восемью электронами вообще?', '2,8,8,2; сумма 20=Z; нет, это не её полная вместимость.'),
        ('Является ли H щелочным металлом? Сколько внешних электронов у He и Ne? Зачем сравнивать Li и Na?', 'Нет; два и восемь; у Li и Na сходство внешнего слоя — по одному электрону.'),
        ('Какой изученный неметалл проводит ток? Достаточно ли блеска для определения Al? Одинаковы ли графит и алмаз?', 'Графит; нет; оба углерод, но строение и свойства различаются.'),
        ('Куда растёт размер Li → Na → K? Каков порядок Na, Mg, Al? Можно ли так вычислить температуру плавления?', 'Растёт вниз; Na > Mg > Al при согласованном определении радиуса; нет.'),
        ('Найдите Aᵣ при 75% атомов с 20 u и 25% с 22 u. Какая единица у Aᵣ? Приемлемы ли доли 75% и 40%?', '20,5; единицы нет; нет, сумма 115%.'),
        ('Каков заряд при 12p и 10e? Становится ли Mg²⁺ элементом второго периода? Равен ли Na⁺ неону?', '+2; нет, Mg остаётся в третьем; нет, различается число протонов.'),
        ('Каков заряд набора 2Na⁺ и Cl⁻? Какое отношение нужно для NaCl? Сколько Cl⁻ уравновешивает Mg²⁺? Доказывает ли прозрачность состав?', '+1; 1:1; два; нет. Объясните каждый ответ через данные и предел вывода.'),
    ),
    'kz': (
        ('Z=17 нені білдіреді? Aᵣ≈35,45 бойынша нейтронды табамыз ба? Co мен CO айырмасы қандай?', '17 протон; жоқ, изотоп керек; Co — кобальт, CO — көміртек монооксидінің формуласы.'),
        ('Cl периоды мен тобы қандай? He қайда? Үзіндіде толтырылмаған 3–12 топтар толық кестеде бар ма?', 'Cl: 3 және 17; He: 1 және 18; иә, толық кестеде бар.'),
        ('Ca электрондары қалай таралады? Қосындыны қалай тексереміз? Үшінші қабаттың толық сыйымдылығы сегіз бе?', '2,8,8,2; қосынды 20=Z; жоқ, бұл оның толық сыйымдылығы емес.'),
        ('H сілтілік металл ма? He және Ne сыртқы электрон саны қандай? Li мен Na неге салыстырылады?', 'Жоқ; екі және сегіз; Li мен Na сыртқы қабатында бір электроннан бар.'),
        ('Қай оқылған бейметалл ток өткізеді? Al анықтауға жылтыр жеткілікті ме? Графит пен алмаз бірдей ме?', 'Графит; жоқ; екеуі көміртек, бірақ құрылымы мен қасиеті өзгеше.'),
        ('Li → Na → K өлшемі қалай өзгереді? Na, Mg, Al реті қандай? Балқу температурасын осылай есептейміз бе?', 'Төмен қарай өседі; бірдей радиус анықтамасында Na > Mg > Al; жоқ.'),
        ('20 u атомдар 75%, 22 u атомдар 25% болса, Aᵣ қандай? Оның бірлігі бар ма? 75% және 40% жарамды ма?', '20,5; бірліксіз; жоқ, қосынды 115%.'),
        ('12p және 10e заряды қандай? Mg²⁺ екінші период элементіне айнала ма? Na⁺ неонмен бірдей ме?', '+2; жоқ, Mg үшінші периодта қалады; жоқ, протон саны өзгеше.'),
        ('2Na⁺ пен Cl⁻ қосынды заряды қандай? NaCl қатынасы қандай? Mg²⁺ үшін неше Cl⁻ керек? Мөлдірлік құрамды дәлелдей ме?', '+1; 1:1; екі; жоқ. Әр жауапты дерек және қорытынды шегімен түсіндіріңіз.'),
    ),
    'en': (
        ('What does Z=17 tell you? Can Aᵣ≈35.45 give neutron count? How do Co and CO differ?', '17 protons; no, an isotope is needed; Co is cobalt, CO is carbon monoxide’s formula.'),
        ('Give Cl’s period and group. Where is He? Do groups 3–12 exist although our extract leaves them unfilled?', 'Cl: 3 and 17; He: 1 and 18; yes, in the full table.'),
        ('Give Ca’s shell arrangement. How do you check the sum? Is eight the third shell’s full capacity?', '2,8,8,2; total 20=Z; no, it is not the shell’s full capacity.'),
        ('Is H an alkali metal? How many outer electrons do He and Ne have? Why compare Li and Na?', 'No; two and eight; Li and Na each have one outer electron.'),
        ('Which studied nonmetal conducts? Does shine identify Al? Are graphite and diamond identical?', 'Graphite; no; both are carbon but their structures and properties differ.'),
        ('How does size change Li → Na → K? Order Na, Mg and Al. Does this calculate melting points?', 'It increases downwards; Na > Mg > Al under a consistent radius definition; no.'),
        ('Find Aᵣ for 75% atoms at 20 u and 25% at 22 u. What unit does Aᵣ have? Are 75% and 40% valid fractions?', '20.5; no unit; no, the total is 115%.'),
        ('What charge comes from 12p and 10e? Does Mg²⁺ move to period 2? Is Na⁺ neon?', '+2; no, Mg remains in period 3; no, proton counts differ.'),
        ('Find the total charge of 2Na⁺ and Cl⁻. What ratio does NaCl need? How many Cl⁻ balance Mg²⁺? Does clarity establish composition?', '+1; 1:1; two; no. Explain each using evidence and the limit of the conclusion.'),
    ),
}
LABELS = {
    'ru': ('Первые 20 элементов: фрагмент таблицы', 'Группы 1–18 →; периоды 1–4 ↓', 'Пустые места здесь не означают отсутствие элементов в полной таблице.', 'Проверка знаний', 'Ответы', 'Запись в дневнике', 'Запишите первоначальный ответ, использованные данные, исправленное объяснение и границу вывода. Затем закройте пример и объясните новый случай своими словами. Если ошиблись, укажите конкретный шаг: чтение карточки, счёт частиц, выбор модели или вывод.', 'Источники и границы модели', 'Z | символ | период | группа | электроны по слоям'),
    'kz': ('Алғашқы 20 элемент: кесте үзіндісі', '1–18 топтар →; 1–4 периодтар ↓', 'Мұндағы бос орындар толық кестеде элементтер жоқ дегенді білдірмейді.', 'Білімді тексеру', 'Жауаптар', 'Күнделік жазбасы', 'Алғашқы жауапты, қолданылған деректі, түзетілген түсіндірмені және қорытынды шегін жазыңыз. Кейін мысалды жауып, жаңа жағдайды өз сөзіңізбен түсіндіріңіз. Қате болса, нақты қадамды көрсетіңіз: карточканы оқу, бөлшектерді санау, модельді таңдау немесе қорытынды жасау.', 'Дереккөздер мен модель шегі', 'Z | таңба | период | топ | қабаттағы электрондар'),
    'en': ('First 20 elements: a table extract', 'Groups 1–18 →; periods 1–4 ↓', 'Empty spaces here do not mean elements are absent from the full table.', 'Knowledge check', 'Answers', 'Notebook record', 'Record your initial answer, the evidence used, the corrected explanation and the limit of your conclusion. Close the example and explain a new case in your own words. If you made an error, identify the step: reading the card, counting particles, choosing a model or drawing a conclusion.', 'Sources and model limits', 'Z | symbol | period | group | electrons by shell'),
}
SOURCES = (
    ('IUPAC: periodic table', 'https://iupac.org/what-we-do/periodic-table-of-elements/'),
    ('ACS: energy-level models', 'https://www.acs.org/middleschoolchemistry/lessonplans/chapter4/lesson3.html'),
    ('RSC: carbon', 'https://periodic-table.rsc.org/element/6/carbon'),
    ('CIAAW: atomic weights', 'https://ciaaw.org/atomic-weights.htm'),
    ('RSC: periodic trends', 'https://periodic-table.rsc.org/trends'),
    ('CIAAW: mass and relative atomic mass', 'https://ciaaw.org/info.htm'),
    ('RSC: sodium radii', 'https://periodic-table.rsc.org/element/11/sodium'),
    ('RSC: magnesium radii', 'https://periodic-table.rsc.org/element/12/magnesium'),
    ('RSC: aluminium radii', 'https://periodic-table.rsc.org/element/13/aluminium'),
)


def text(x, y, value, size=25, color='#f1f7f2', anchor='start'):
    return (f'<text x="{x}" y="{y}" font-family="Arial,sans-serif" '
            f'font-size="{size}" fill="{color}" text-anchor="{anchor}">{escape(str(value))}</text>')


def start_svg(title, description, width=1280, height=720):
    return [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {width} {height}" role="img" aria-labelledby="title desc">',
            f'<title id="title">{escape(title)}</title><desc id="desc">{escape(description)}</desc>',
            f'<rect width="{width}" height="{height}" rx="24" fill="#102b3b"/>']


def table_svg(lang):
    labels = LABELS[lang]
    parts = start_svg(labels[0], labels[2])
    parts.extend([text(40, 60, labels[0], 30), text(40, 108, labels[1], 23)])
    for group in range(1, 19):
        parts.append(text(110+(group-1)*63, 180, group, 21, anchor='middle'))
    for period in range(1, 5):
        parts.append(text(40, 245+(period-1)*95, period, 27))
    for z, symbol, period, group, shells in ELEMENTS:
        x, y = 80+(group-1)*63, 200+(period-1)*95
        parts.extend([f'<rect data-z="{z}" data-period="{period}" data-group="{group}" x="{x}" y="{y}" width="58" height="85" rx="6" fill="#d8efec"/>',
                      text(x+8, y+23, z, 17, '#19394b'), text(x+29, y+60, symbol, 28, '#19394b', 'middle')])
    for i, line in enumerate(wrap(labels[2], 80)):
        parts.append(text(40, 625+i*30, line, 23))
    parts.append('</svg>')
    return '\n'.join(parts)+'\n'


def diagram(number, lang):
    title = LESSONS[number-14][lang][0]
    parts = start_svg(title, NOT_TO_SCALE[lang], height=560)
    for i, line in enumerate(wrap(title, 64)):
        parts.append(text(40, 55+i*37, line, 30))
    parts.append(text(40, 145, NOT_TO_SCALE[lang], 20, '#bedde4'))
    labels = STEPS[lang][number-14]
    for i in range(4):
        x=40+i*310
        parts.append(f'<rect x="{x}" y="195" width="275" height="295" rx="18" fill="#e6f4ef"/>')
        parts.append(text(x+20, 244, f'{i+1}. {labels[i*2]}', 22, '#143242'))
        for j, line in enumerate(wrap(labels[i*2+1], 18, break_long_words=False)):
            parts.append(text(x+20, 320+j*36, line, 24, '#143242'))
        if i < 3:
            parts.append(text(x+282, 365, '→', 29))
    return '\n'.join(parts+['</svg>'])+'\n'


def visual(number):
    parts = start_svg(f'Chemistry {number}: symbolic teaching diagram',
                      'Original diagram. Sizes, colours and distances are symbolic, not to scale.')
    parts.append(text(65, 88, f'CHEMISTRY · {number:02d}', 36))

    def box(x, y, w=285, h=300):
        parts.append(f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="22" fill="#e5f3ed" fill-opacity=".12" stroke="#70c7cb" stroke-width="4"/>')

    def label(x, y, value, size=42):
        parts.append(text(x, y, value, size, anchor='middle'))

    if number == 14:
        for x, sym, z, mass in ((175, 'Na', 11, '22.99'), (770, 'Cl', 17, '35.45')):
            box(x, 190, 335, 365)
            label(x+168, 260, f'Z = {z}', 40)
            label(x+168, 400, sym, 110)
            label(x+168, 505, f'Aᵣ ≈ {mass}', 37)
    elif number == 15:
        # The cover uses the same real coordinates as the localized table.
        for z, sym, period, group, shells in ELEMENTS:
            x, y = 75+(group-1)*63, 200+(period-1)*100
            box(x, y, 57, 86)
            label(x+28, y+25, z, 17)
            label(x+28, y+62, sym, 26)
        for group in range(1, 19):
            label(103+(group-1)*63, 177, group, 21)
        label(640, 675, 'Z ≤ 20 · groups 1–18 · periods 1–4', 28)
    elif number == 16:
        for y, sym, values in ((285, 'Na', (2,8,1)), (480, 'Ca', (2,8,8,2))):
            label(155, y+30, sym, 52)
            for i, count in enumerate(values):
                box(270+i*215, y-50, 175, 130)
                label(357+i*215, y+30, count, 58)
        label(640, 645, 'neutral · ground state · Z ≤ 20', 31)
    elif number == 17:
        for x, head, vals in ((230, '1', ('Li','Na','K')), (640,'17',('F','Cl')), (1050,'18',('He','Ne','Ar'))):
            box(x-130, 180, 260, 390)
            label(x, 235, f'GROUP {head}', 30)
            for i, value in enumerate(vals): label(x, 325+i*90, value, 57)
        label(640, 660, 'H ≠ alkali metal · He: 2 outer e', 32)
    elif number == 18:
        for x, sym, word1, word2 in ((235,'Al','conducts','malleable'), (640,'C · graphite','conducts','brittle'), (1045,'S','poor conductor','brittle')):
            box(x-180, 210, 360, 330)
            label(x, 300, sym, 43)
            label(x, 410, word1, 29)
            label(x, 475, word2, 29)
        label(640, 650, 'One property alone is insufficient', 33)
    elif number == 19:
        for x, r, sym in ((240,95,'Na'),(640,75,'Mg'),(1040,55,'Al')):
            parts.append(f'<circle cx="{x}" cy="350" r="{r}" fill="#83ccd0" fill-opacity=".35" stroke="#e2f1e8" stroke-width="4"/>')
            label(x, 360, sym, 42)
        label(640, 535, 'Na > Mg > Al       Li < Na < K', 42)
        label(640, 635, 'Qualitative atomic size · not to scale', 30)
    elif number == 20:
        for i, val in enumerate((20,20,20,22)):
            box(80+i*300, 200, 220, 180)
            label(190+i*300, 310, f'{val} u', 50)
        label(640, 480, '(20 + 20 + 20 + 22) / 4 = 20.5 u', 40)
        label(640, 585, 'Aᵣ = 20.5', 65)
        label(640, 665, 'Invented mixture · given teaching masses', 29)
    elif number == 21:
        for x, sym, p, e, layers in ((350,'Mg',12,12,'2 | 8 | 2'),(950,'Mg²⁺',12,10,'2 | 8')):
            box(x-210, 190, 420, 360)
            label(x, 275, sym, 67)
            label(x, 375, f'{p} p · {e} e', 42)
            label(x, 465, layers, 43)
        label(650, 375, '→', 65)
        label(640, 655, 'Z = 12 · period 3 · group 2', 40)
    else:
        for row in range(3):
            for col in range(6):
                x, y = 240+col*160, 240+row*125
                value, color = ('Na⁺','#efba7c') if (row+col)%2 == 0 else ('Cl⁻','#88d0d4')
                parts.append(f'<circle cx="{x}" cy="{y}" r="52" fill="{color}"/>')
                parts.append(text(x, y+10, value, 31, '#143242', 'middle'))
        label(640, 655, 'Na⁺ : Cl⁻ = 1 : 1      +1 − 1 = 0', 40)
    return '\n'.join(parts+['</svg>'])+'\n'


def render(number, lang):
    row = LESSONS[number-14]
    title, summary, *sections = row[lang]
    labels = LABELS[lang]
    headings = HEAD[lang]
    out = [f'# {title}', f'_Лид (summary):_ **{summary}**', f'## {headings[0]}',
           f'![{title}](/static/course/chemistry/map-{number:02d}-{lang}.svg)',
           f'![{NOT_TO_SCALE[lang]}](/static/course/chemistry/visuals/{number:02d}-elements.svg)',
           f'*{NOT_TO_SCALE[lang]}*']
    for heading, section in zip(headings[1:9], sections, strict=True):
        out.extend([f'## {heading}', section])
    if number in (15, 16):
        out.extend([f'![{labels[0]}](/static/course/chemistry/first-20-{lang}.svg)',
                    f'*{labels[2]}*'])
        rows = [f'| {labels[8]} |', '| --- | --- | --- | --- | --- |']
        rows.extend(f'| {z} | {symbol} | {period} | {group} | {", ".join(map(str, shells))} |'
                    for z, symbol, period, group, shells in ELEMENTS)
        out.append('\n'.join(rows))
    questions, answer = QUIZ[lang][number-14]
    exercise = {'ru': 'Задание', 'kz': 'Тапсырма', 'en': 'Exercise'}[lang]
    numbered = '\n'.join(f'{i}. {q.strip()}?' for i, q in enumerate(questions.rstrip('?').split('? '), 1))
    out.extend([f'## {exercise}', numbered, f'## {labels[4]}', answer,
                f'## {labels[5]}', labels[6], f'## {headings[10]}', RETURN[lang], f'## {labels[7]}'])
    source_ids = ((0, 3, 5) if number in (14, 20) else (0, 2) if number == 18
                  else (0, 4, 6, 7, 8) if number == 19 else (0, 1))
    out.append('\n'.join(f'- [{SOURCES[i][0]}]({SOURCES[i][1]})' for i in source_ids))
    previous = PREVIOUS[-1] if number == 14 else LESSONS[number-15]
    out.append(f'[{previous[lang][0]}](/read/{slug(number-1, previous["stem"])}?lang={lang})')
    if number < 22:
        following = LESSONS[number-13]
        out.append(f'[{following[lang][0]}](/read/{slug(number+1, following["stem"])}?lang={lang})')
    return '\n\n'.join(out)+'\n'


def main():
    VISUAL_DIR.mkdir(parents=True, exist_ok=True)
    for lang in LANGS:
        (MAP_DIR / f'first-20-{lang}.svg').write_text(table_svg(lang), encoding='utf-8')
    for number, row in enumerate(LESSONS, 14):
        (VISUAL_DIR / f'{number:02d}-elements.svg').write_text(visual(number), encoding='utf-8')
        for lang in LANGS:
            suffix = '' if lang == 'ru' else f'-{lang}'
            (LESSON_DIR / f'{row["stem"]}{suffix}.md').write_text(render(number, lang), encoding='utf-8')
            (MAP_DIR / f'map-{number:02d}-{lang}.svg').write_text(diagram(number, lang), encoding='utf-8')


if __name__ == '__main__':
    main()
