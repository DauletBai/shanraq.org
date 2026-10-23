#!/usr/bin/env python3
"""Render chemistry lessons 6–13 and original, scalable teaching diagrams."""
from html import escape
from pathlib import Path
from textwrap import wrap

try:
    from .chemistry_particles_data import LESSONS
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from .chemistry_data import LESSONS as FIRST_BLOCK
    from .chemistry_elements_data import LESSONS as ELEMENT_BLOCK
except ImportError:
    from chemistry_particles_data import LESSONS
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from chemistry_data import LESSONS as FIRST_BLOCK
    from chemistry_elements_data import LESSONS as ELEMENT_BLOCK

ROOT = Path(__file__).resolve().parents[2]
VISUAL_DIR = MAP_DIR / 'visuals'
MAP_STEPS = {
    'ru': (
        ('ВИДНО', 'Форма образца', 'МОДЕЛЬ', 'Движение частиц', 'СРАВНИТЬ', 'Лёд и вода', 'ПРЕДЕЛ', 'Не фото'),
        ('ВИДНО', 'Цвет расплылся', 'МОДЕЛЬ', 'Случайное движение', 'КОНТРОЛЬ', 'Возможен поток', 'ПРЕДЕЛ', 'Не видны молекулы'),
        ('КАРТОЧКА', 'Число протонов', 'ЭЛЕМЕНТ', 'Одинаковое Z', 'СРАВНИТЬ', '6 и 8 протонов', 'ПРЕДЕЛ', 'Цвет условен'),
        ('ЯДРО', 'Протоны и нейтроны', 'ОБЛАСТЬ', 'Электроны', 'ЗАРЯД', 'p минус e', 'ПРЕДЕЛ', 'Не орбиты'),
        ('ПРОТОНЫ', 'Элемент', 'НЕЙТРОНЫ', 'Изотоп', 'СУММА', 'Массовое число', 'ПРЕДЕЛ', 'Не средняя масса'),
        ('ВОДА', 'Молекулы H₂O', 'СОЛЬ', 'Ионная решётка', 'ЗАПИСЬ', 'NaCl: 1 к 1', 'ПРЕДЕЛ', 'Не пара'),
        ('ЗНАК', 'Элемент', 'ИНДЕКС', 'Внутри частицы', 'ЧИСЛО', 'Целые частицы', 'ПРОВЕРЬ', 'Счёт атомов'),
        ('ВИДНО', 'Соль не видна', 'МОДЕЛЬ', 'Ионы в воде', 'МАССА', 'Закрытый сосуд', 'ПРЕДЕЛ', 'Нужна проверка'),
    ),
    'kz': (
        ('КӨРІНЕДІ', 'Үлгі пішіні', 'МОДЕЛЬ', 'Бөлшек қозғалысы', 'САЛЫСТЫР', 'Мұз бен су', 'ШЕК', 'Фото емес'),
        ('КӨРІНЕДІ', 'Түс тарады', 'МОДЕЛЬ', 'Ретсіз қозғалыс', 'БАҚЫЛАУ', 'Ағын болуы мүмкін', 'ШЕК', 'Молекула көрінбейді'),
        ('КАРТОЧКА', 'Протон саны', 'ЭЛЕМЕНТ', 'Бірдей Z', 'САЛЫСТЫР', '6 және 8 протон', 'ШЕК', 'Түс шартты'),
        ('ЯДРО', 'Протон, нейтрон', 'АЙМАҚ', 'Электрондар', 'ЗАРЯД', 'p азайту e', 'ШЕК', 'Орбита емес'),
        ('ПРОТОН', 'Элемент', 'НЕЙТРОН', 'Изотоп', 'ҚОСЫНДЫ', 'Массалық сан', 'ШЕК', 'Орташа масса емес'),
        ('СУ', 'H₂O молекуласы', 'ТҰЗ', 'Иондық тор', 'ЖАЗУ', 'NaCl: 1 және 1', 'ШЕК', 'Жеке жұп емес'),
        ('ТАҢБА', 'Элемент', 'ИНДЕКС', 'Бөлшектің ішінде', 'САН', 'Тұтас бөлшектер', 'ТЕКСЕР', 'Атомдарды сана'),
        ('КӨРІНЕДІ', 'Тұз көрінбейді', 'МОДЕЛЬ', 'Судағы иондар', 'МАССА', 'Жабық ыдыс', 'ШЕК', 'Тексеру қажет'),
    ),
    'en': (
        ('OBSERVE', 'Sample shape', 'MODEL', 'Particle motion', 'COMPARE', 'Ice and water', 'LIMIT', 'Not a photo'),
        ('OBSERVE', 'Colour spread', 'MODEL', 'Random motion', 'CONTROL', 'Flow is possible', 'LIMIT', 'Molecules unseen'),
        ('CARD', 'Proton count', 'ELEMENT', 'Same Z', 'COMPARE', 'Six and eight', 'LIMIT', 'Colour is chosen'),
        ('NUCLEUS', 'Protons, neutrons', 'REGION', 'Electrons', 'CHARGE', 'p minus e', 'LIMIT', 'Not orbits'),
        ('PROTONS', 'Element', 'NEUTRONS', 'Isotope', 'SUM', 'Mass number', 'LIMIT', 'Not average mass'),
        ('WATER', 'H₂O molecules', 'SALT', 'Ionic lattice', 'NOTATION', 'NaCl: one to one', 'LIMIT', 'Not a pair'),
        ('SYMBOL', 'Element', 'SUBSCRIPT', 'Within particle', 'NUMBER', 'Whole particles', 'CHECK', 'Count atoms'),
        ('OBSERVE', 'Salt unseen', 'MODEL', 'Ions in water', 'MASS', 'Closed vessel', 'LIMIT', 'Needs a test'),
    ),
}
NOT_TO_SCALE = {'ru': 'Учебная модель · частицы и расстояния не в масштабе',
                'kz': 'Оқу моделі · бөлшектер мен арақашықтық масштабта емес',
                'en': 'Teaching model · particles and distances are not to scale'}
QUIZ = {
    'ru': (
        ('Почему лёд сохраняет форму?', 'Что остаётся тем же при плавлении?', 'Частицы колеблются около устойчивых положений; молекулы H₂O сохраняются.'),
        ('Что видно в опыте с красителем?', 'Все ли отдельные частицы идут в одну сторону?', 'Расширение цветной области; нет, отдельные движения случайны.'),
        ('Атомы с 6 и 8 протонами одного элемента?', 'Что задаёт элемент?', 'Нет; число протонов.'),
        ('Сколько электронов у нейтрального атома с 3 протонами?', 'Каков заряд при 2 электронах?', 'Три; +1.'),
        ('Каковы массовые числа 6p+6n и 6p+7n?', 'Меняется ли элемент?', '12 и 13; нет, оба углерод.'),
        ('Отдельна ли молекула NaCl в обычном кристалле?', 'Что означает H₂O?', 'Нет, NaCl задаёт отношение ионов; одна молекула воды имеет два H и один O.'),
        ('Сколько O в 3CO₂?', 'Сколько молекул в 2H₂O?', 'Шесть O; две молекулы.'),
        ('Какова масса закрытой системы 50+100+5 г?', 'Доказывает ли прозрачность состав?', '155 г; нет.'),
    ),
    'kz': (
        ('Мұз неге пішінін сақтайды?', 'Ерігенде не сақталады?', 'Бөлшектер тұрақты орын маңында тербеледі; H₂O молекулалары сақталады.'),
        ('Бояу мысалында не көрінеді?', 'Барлық бөлшек бір бағытқа бара ма?', 'Түсті аймақ кеңейеді; жоқ, жеке қозғалыс ретсіз.'),
        ('6 және 8 протоны бар атомдар бір элемент пе?', 'Элементті не анықтайды?', 'Жоқ; протон саны.'),
        ('3 протоны бар бейтарап атомда неше электрон?', '2 электрон болса заряды қандай?', 'Үш; +1.'),
        ('6p+6n және 6p+7n массалық сандары қандай?', 'Элемент өзгере ме?', '12 және 13; жоқ, екеуі де көміртек.'),
        ('Кәдімгі кристалда жеке NaCl молекуласы бар ма?', 'H₂O нені білдіреді?', 'Жоқ, NaCl иондар қатынасы; су молекуласында екі H, бір O бар.'),
        ('3CO₂-де қанша O бар?', '2H₂O-де неше молекула?', 'Алты O; екі молекула.'),
        ('Жабық жүйе массасы 50+100+5 г болса, нәтижесі?', 'Мөлдірлік құрамды дәлелдей ме?', '155 г; жоқ.'),
    ),
    'en': (
        ('Why does ice keep its shape?', 'What survives melting?', 'Particles vibrate near stable positions; H₂O molecules remain.'),
        ('What is visible in the dye case?', 'Does every particle move one way?', 'The coloured region expands; no, individual motion is random.'),
        ('Are atoms with six and eight protons one element?', 'What fixes element identity?', 'No; proton count.'),
        ('How many electrons in a neutral atom with three protons?', 'Charge with two electrons?', 'Three; +1.'),
        ('Mass numbers for 6p+6n and 6p+7n?', 'Does the element change?', '12 and 13; no, both are carbon.'),
        ('Is there a separate NaCl molecule in an ordinary crystal?', 'What does H₂O mean?', 'No, NaCl is an ion ratio; a water molecule has two H and one O.'),
        ('How many O atoms in 3CO₂?', 'How many molecules in 2H₂O?', 'Six O; two molecules.'),
        ('Closed-system mass of 50+100+5 g?', 'Does clarity establish composition?', '155 g; no.'),
    ),
}


def slug(number, stem):
    return f'chemistry-{number:02d}-{stem.split("-", 1)[1]}'


def diagram(number, lang):
    title = escape(LESSONS[number - 6][lang][0])
    labels = MAP_STEPS[lang][number - 6]
    parts = [f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1280 490" role="img" aria-labelledby="title desc">',
             f'<title id="title">{title}</title><desc id="desc">{escape(NOT_TO_SCALE[lang])}</desc>',
             '<rect width="1280" height="490" rx="24" fill="#102b3b"/>',
             f'<text x="45" y="68" fill="white" font-family="Arial,sans-serif" font-size="29" font-weight="bold">{title}</text>',
             f'<text x="45" y="105" fill="#bedde4" font-family="Arial,sans-serif" font-size="18">{escape(NOT_TO_SCALE[lang])}</text>']
    colors = ('#2db7c4', '#f0ad58', '#9d8bdb', '#75c39c')
    for index in range(4):
        x = 42 + index * 308
        top, bottom = labels[index * 2:index * 2 + 2]
        parts += [f'<rect x="{x}" y="157" width="275" height="265" rx="22" fill="#f7faf7" stroke="{colors[index]}" stroke-width="6"/>',
                  f'<circle cx="{x+41}" cy="205" r="22" fill="{colors[index]}"/>',
                  f'<text x="{x+33}" y="213" font-family="Arial,sans-serif" font-size="22" font-weight="bold" fill="#143242">{index+1}</text>',
                  f'<text x="{x+20}" y="274" font-family="Arial,sans-serif" font-size="19" font-weight="bold" fill="#173749">{escape(top)}</text>']
        for line_no, line in enumerate(wrap(bottom, 21, break_long_words=False)):
            parts.append(f'<text x="{x+20}" y="{319+line_no*30}" font-family="Arial,sans-serif" font-size="20" fill="#173749">{escape(line)}</text>')
        if index < 3:
            parts.append(f'<path d="M{x+279} 290 h24 m-10 -10 10 10 -10 10" fill="none" stroke="#d8e9e6" stroke-width="4"/>')
    parts.append('</svg>')
    return '\n'.join(parts) + '\n'


def visual(number):
    """A concept-specific cover built from scalable vector geometry."""
    parts = ['<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1280 720" role="img" aria-labelledby="title desc">',
             f'<title id="title">Chemistry lesson {number}: symbolic model</title>',
             '<desc id="desc">Teaching diagram, not a photograph. Particle sizes and distances are not to scale.</desc>',
             '<defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#081b2d"/><stop offset="1" stop-color="#145467"/></linearGradient></defs>',
             '<rect width="1280" height="720" fill="url(#bg)"/>',
             f'<text x="86" y="104" font-family="Arial,sans-serif" font-size="38" font-weight="bold" fill="#e9f5f2">CHEMISTRY · {number:02d}</text>']
    def circle(x, y, r=20, color='#5bc5d1'):
        parts.append(f'<circle cx="{x}" cy="{y}" r="{r}" fill="{color}" stroke="#ecffff" stroke-width="3"/>')
    def label(x, y, value, size=34):
        parts.append(f'<text x="{x}" y="{y}" text-anchor="middle" font-family="Arial,sans-serif" font-size="{size}" font-weight="bold" fill="#f5f8ef">{escape(value)}</text>')
    if number == 6:
        for x in (95, 478, 861):
            parts.append(f'<rect x="{x}" y="205" width="310" height="355" rx="24" fill="#f3fcf8" fill-opacity=".08" stroke="#7fd5d4" stroke-width="6"/>')
        for row in range(4):
            for col in range(4): circle(145+col*69, 360+row*55, 21, '#f7bd73')
        for x,y in ((540,350),(601,375),(670,345),(565,445),(650,460),(528,515),(703,507)):
            circle(x,y,22,'#5bc5d1')
        for x,y in ((912,265),(1110,280),(1018,350),(1140,416),(930,475),(1065,520)):
            circle(x,y,20,'#b7a2ee')
        for x,v in ((250,'■'),(632,'~'),(1015,'↗')): label(x,640,v,55)
    elif number == 7:
        parts.append('<rect x="118" y="260" width="1045" height="320" rx="30" fill="#ebffff" fill-opacity=".12" stroke="#81d5d7" stroke-width="6"/>')
        for i in range(32):
            x=170+(i*53)%770; y=310+(i*71)%220
            circle(x,y,13,'#f1aa69')
        for i in range(14):
            circle(960+(i*29)%150,320+(i*83)%220,12,'#f1aa69')
        parts.append('<path d="M270 225 H1000 m-30 -20 30 20 -30 20" fill="none" stroke="#fff1ba" stroke-width="9"/>')
    elif number == 8:
        for x,count,mark in ((385,6,'C'),(885,8,'O')):
            parts.append(f'<circle cx="{x}" cy="395" r="146" fill="#dcefeb" fill-opacity=".13" stroke="#8bd2d5" stroke-width="7"/>')
            for i in range(count): circle(x-65+(i%4)*43,360+(i//4)*44,16,'#f1af70')
            label(x,620,f'{mark} · {count} p',47)
        parts.append('<path d="M560 395 H710" stroke="#f6eac9" stroke-width="7" stroke-dasharray="18 14"/>')
    elif number == 9:
        parts.append('<circle cx="640" cy="395" r="234" fill="none" stroke="#70bcca" stroke-width="6" stroke-dasharray="10 15"/>')
        parts.append('<circle cx="640" cy="395" r="96" fill="#efd4ab"/>')
        for x,y,v in ((605,371,'+'),(665,371,'0'),(605,431,'+'),(665,431,'0')):
            label(x,y,v,45)
        for x,y in ((640,160),(875,395)):circle(x,y,22,'#79d9e4')
        label(1110,422,'−',70)
    elif number == 10:
        for x,neutrons,mass in ((370,6,'¹²C'),(910,7,'¹³C')):
            parts.append(f'<circle cx="{x}" cy="405" r="151" fill="#cfe9e8" fill-opacity=".15" stroke="#85d5d8" stroke-width="7"/>')
            for i in range(6):circle(x-78+(i%3)*67,352+(i//3)*64,25,'#f2ad6d')
            for i in range(neutrons):circle(x-85+(i%4)*57,475+(i//4)*44,19,'#b8a2ea')
            label(x,630,mass,54)
        parts.append('<path d="M535 401 H744" stroke="#e5f4e9" stroke-width="7" stroke-dasharray="15 12"/>')
    elif number == 11:
        for x in (180,490):
            circle(x,390,48,'#f1ad70'); circle(x+85,330,30,'#7bcdda'); circle(x+85,450,30,'#7bcdda')
        label(370,605,'H₂O',62)
        for row in range(3):
            for col in range(4):
                x=820+col*78; y=302+row*78
                circle(x,y,32,'#8bd1df' if (row+col)%2 else '#eab275')
        label(945,605,'NaCl',62)
        parts.append('<path d="M645 245 V575" stroke="#f5f5e8" stroke-opacity=".6" stroke-width="4"/>')
    elif number == 12:
        for x,y,value,size in ((355,340,'H₂O',94),(355,500,'2H₂O',70),(910,340,'CO₂',94),(910,500,'3CO₂',70)):
            label(x,y,value,size)
        parts.append('<path d="M620 205 V570" stroke="#7bced2" stroke-width="5"/>')
    else:
        parts.append('<path d="M280 275 H1000 V545 H280 Z" fill="#ccebf0" fill-opacity=".15" stroke="#85d4d8" stroke-width="7"/>')
        for x,y,color in ((400,390,'#f0ad6b'),(495,445,'#f0ad6b'),(615,365,'#78d0db'),(760,455,'#78d0db'),(910,385,'#f0ad6b')):
            circle(x,y,29,color)
        label(640,650,'50 g + 100 g + 5 g = 155 g',51)
    parts.append('</svg>')
    return '\n'.join(parts) + '\n'


def render(number, lang):
    row = LESSONS[number - 6]
    title, summary, *sections = row[lang]
    headings = HEAD[lang]
    base = f'/static/course/chemistry/map-{number:02d}-{lang}.svg'
    cover = f'/static/course/chemistry/visuals/{number:02d}-particles.svg'
    out = [f'# {title}', f'_Лид (summary):_ **{summary}**', f'## {headings[0]}',
           f'![{title}]({base})', f'![{escape(NOT_TO_SCALE[lang])}]({cover})',
           f'*{NOT_TO_SCALE[lang]}*']
    for heading, section in zip(headings[1:9], sections):
        out.extend([f'## {heading}', section])
    q1, q2, answer = QUIZ[lang][number-6]
    out.extend([f'## {"Проверка знаний" if lang == "ru" else "Білімді тексеру" if lang == "kz" else "Knowledge check"}',
                f'1. {q1}\n2. {q2}\n\n**{"Ответы" if lang == "ru" else "Жауаптар" if lang == "kz" else "Answers"}:** {answer}',
                f'## {headings[10]}', RETURN[lang]])
    if number == 6:
        previous = FIRST_BLOCK[4]
        out.append(f'[{previous[lang][0]}](/read/{slug(5, previous["stem"])}?lang={lang})')
    else:
        previous = LESSONS[number - 7]
        out.append(f'[{previous[lang][0]}](/read/{slug(number-1, previous["stem"])}?lang={lang})')
    if number < 13:
        next_row = LESSONS[number - 5]
        out.append(f'[{next_row[lang][0]}](/read/{slug(number+1, next_row["stem"])}?lang={lang})')
    else:
        next_row = ELEMENT_BLOCK[0]
        out.append(f'[{next_row[lang][0]}](/read/{slug(14, next_row["stem"])}?lang={lang})')
    return '\n\n'.join(out) + '\n'


def main():
    LESSON_DIR.mkdir(parents=True, exist_ok=True)
    VISUAL_DIR.mkdir(parents=True, exist_ok=True)
    for number in range(6, 14):
        (VISUAL_DIR / f'{number:02d}-particles.svg').write_text(visual(number), encoding='utf-8')
        for lang in LANGS:
            suffix = '' if lang == 'ru' else f'-{lang}'
            (LESSON_DIR / f'{LESSONS[number-6]["stem"]}{suffix}.md').write_text(render(number, lang), encoding='utf-8')
            (MAP_DIR / f'map-{number:02d}-{lang}.svg').write_text(diagram(number, lang), encoding='utf-8')


if __name__ == '__main__':
    main()
