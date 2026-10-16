#!/usr/bin/env python3
"""Build the complete trilingual magnetism and communication block (57–66)."""
from html import escape
from pathlib import Path
import textwrap

try:
    from .physics_magnetism_data import MAGNETISM
except ImportError:
    from physics_magnetism_data import MAGNETISM

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/physics"
MAPS = ROOT / "web/static/course/physics"

HEAD = {
    "ru": ("Ситуация и прогноз", "Разбираем новые слова", "Разбор учебных данных", "Безопасная проверка", "Где рассуждение ломается", "Опорная схема", "Перенос в новую задачу", "Проверь себя", "Запись в лабораторный журнал"),
    "kz": ("Жағдай және болжам", "Жаңа сөздерді түсінейік", "Оқу деректерін талдау", "Қауіпсіз тексеру", "Ойлаудағы қате", "Тірек сызбасы", "Жаңа жағдайға қолдану", "Өзіңді тексер", "Зертхана журналына жазба"),
    "en": ("Situation and prediction", "New words, step by step", "Working through supplied data", "A safe check", "Where the reasoning breaks", "Support diagram", "Transfer to a new case", "Check your understanding", "Lab notebook entry"),
}
CHECK = {
    "ru": "Закройте объяснение и расскажите своими словами: что здесь действительно наблюдали, что было дано как учебный пример и какой вывод пока остаётся предположением? Затем поменяйте одно условие задачи и предскажите результат до повторного чтения. Если результат не совпал, найдите первый шаг, на котором ваш вывод разошёлся со схемой. Сильный ответ называет границы модели, а не только правильное слово или число. Вы можете выполнить эту проверку на бумаге, без прибора и регистрации.",
    "kz": "Түсіндірмені жауып, өз сөзіңізбен айтыңыз: нені шынымен бақыладық, қандай дерек оқу мысалы ретінде берілді, қай тұжырым әлі болжам? Содан соң есептің бір шартын өзгертіп, мәтінді қайта қарамастан нәтижені болжаңыз. Сәйкес келмесе, ойыңыздың сызбадан қай қадамда айырылғанын табыңыз. Жақсы жауап дұрыс сөзді не санды ғана емес, модельдің қолданылу шегін де атайды. Бұл тексеруді құралсыз, қағазда орындауға болады.",
    "en": "Close the explanation and say in your own words what was really observed, which figures were supplied for teaching, and which conclusion remains a prediction. Change one condition and predict the result before rereading. If it disagrees, find the first step where your reasoning parted from the diagram. A strong answer names the model's limits rather than just a correct word or number. You can do this check on paper, without equipment or registration.",
}
JOURNAL = {
    "ru": "Добавьте страницу к «Лабораторному журналу моего мира». Запишите вопрос, прогноз до опыта, схему и таблицу условий. Отделите наблюдение от выданных учебных чисел. В последней строке укажите, что изменили, что сохранили постоянным, какой вывод поддерживают данные и чего они ещё не доказывают. При отсутствии опыта честно напишите «расчёт по условию» или «предположение», а не «измерено».",
    "kz": "«Менің әлемімнің зертхана журналына» жаңа бет қосыңыз. Сұрақты, тәжірибеге дейінгі болжамды, сызбаны және шарттар кестесін жазыңыз. Бақылауды берілген оқу сандарынан ажыратыңыз. Соңғы жолға нені өзгерткеніңізді, нені тұрақты сақтағаныңызды, дерек қай тұжырымды қолдайтынын және нені әлі дәлелдемейтінін түсіріңіз. Тәжірибе болмаса, «өлшенді» демей, «шарт бойынша есеп» не «болжам» деп жазыңыз.",
    "en": "Add a page to My World Lab Notebook. Record the question, a prediction made before testing, a sketch and a table of conditions. Keep actual observations separate from supplied teaching numbers. In the final line say what changed, what stayed constant, what the evidence supports and what it does not yet prove. If you did no experiment, label the work 'calculation from supplied data' or 'prediction', never 'measured'.",
}
NOTES = {
    "materials": {"ru": ("Сталь: да", "Алюминий: нет", "Пластик: нет", "Зазор: 5 мм"), "kz": ("Болат: иә", "Алюминий: жоқ", "Пластик: жоқ", "Аралық: 5 мм"), "en": ("Steel: yes", "Aluminium: no", "Plastic: no", "Gap: 5 mm")},
    "poles": {"ru": ("N–S: притяжение", "N–N: отталкивание", "Линии: модель", "У каждой части N и S"), "kz": ("N–S: тартылу", "N–N: тебілу", "Сызықтар: модель", "Әр бөлікте N мен S"), "en": ("N–S: attraction", "N–N: repulsion", "Lines are a model", "Each piece has N and S")},
    "compass": {"ru": ("Земное поле", "Местная помеха", "Стрелка: итог", "Карта ≠ компас"), "kz": ("Жер өрісі", "Жергілікті кедергі", "Тілше: жалпы әсер", "Карта ≠ компас"), "en": ("Earth's field", "Nearby disturbance", "Needle: net effect", "Map ≠ compass")},
    "wire": {"ru": ("Ток включён: поле", "Ток выключен", "Направление меняется", "Нужен резистор"), "kz": ("Ток бар: өріс бар", "Ток ажыратылды", "Бағыт өзгереді", "Резистор керек"), "en": ("Current on: field", "Current off", "Direction reverses", "Resistor required")},
    "coil": {"ru": ("20 витков: 2", "40 витков: 4", "Ток одинаков", "Не общий закон"), "kz": ("20 айналым: 2", "40 айналым: 4", "Ток бірдей", "Жалпы заң емес"), "en": ("20 turns: 2", "40 turns: 4", "Same current", "Not a universal law")},
    "motor": {"ru": ("Вход: 10 Дж", "Движение: 6 Дж", "Потери: 4 Дж", "Доля: 60%"), "kz": ("Кіріс: 10 Дж", "Қозғалыс: 6 Дж", "Шығын: 4 Дж", "Үлес: 60%"), "en": ("Input: 10 J", "Motion: 6 J", "Losses: 4 J", "Fraction: 60%")},
    "induction": {"ru": ("Стоит: 0", "Входит: +3", "Выходит: −3", "Меняется поток"), "kz": ("Тұр: 0", "Кіреді: +3", "Шығады: −3", "Ағын өзгереді"), "en": ("Still: 0", "Entering: +3", "Leaving: −3", "Flux changes")},
    "generator": {"ru": ("Педали: 100 Дж", "Лампа: 60 Дж", "Потери: 40 Дж", "Энергия сохраняется"), "kz": ("Педаль: 100 Дж", "Шам: 60 Дж", "Шығын: 40 Дж", "Энергия сақталады"), "en": ("Pedal: 100 J", "Lamp: 60 J", "Losses: 40 J", "Energy is conserved")},
    "signal": {"ru": ("Код: Д = 10", "Код: А = 01", "ДА = 1001", "Ошибка: потеря знака"), "kz": ("Код: И = 10", "Код: Ә = 01", "ИӘ = 1001", "Қате: таңба жоғалу"), "en": ("Code: Y = 10", "Code: E = 01", "YE = 1001", "Error: lost symbol")},
    "capstone": {"ru": ("Закрыто: LED выкл.", "Открыто: LED вкл.", "Геркон: замкн. тип", "Питание: 3 В"), "kz": ("Жабық: LED сөнген", "Ашық: LED жанған", "Геркон: жабық түр", "Қорек: 3 В"), "en": ("Closed: LED off", "Open: LED on", "Reed: normally closed", "Supply: 3 V")},
}


def txt(x, y, value, size=30, anchor="middle", colour="#16445e"):
    return f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" font-size="{size}" font-weight="700" fill="{colour}">{escape(value)}</text>'


def scene(kind, lang):
    """Each sketch uses the same 40-pixel construction grid and distinct physics."""
    dark = '#246381'
    wire = f'stroke="{dark}" stroke-width="9" fill="none" stroke-linecap="round" stroke-linejoin="round"'
    magnet = lambda x, y: (f'<rect x="{x}" y="{y}" width="200" height="74" rx="10" fill="#ec8d76" stroke="#a55950" stroke-width="4"/>'
                           f'<rect x="{x+100}" y="{y}" width="100" height="74" rx="10" fill="#80c6df"/>')
    out = []
    if kind == 'materials':
        out += [magnet(170, 410), txt(220, 459, 'N'), txt(320, 459, 'S'),
                '<path d="M430 446H615" stroke="#e69256" stroke-width="7" stroke-dasharray="12 12"/>',
                '<rect x="640" y="405" width="82" height="82" rx="8" fill="#c7d3d8" stroke="#637f91" stroke-width="4"/>',
                '<rect x="750" y="405" width="82" height="82" rx="8" fill="#e8edf1" stroke="#97aabd" stroke-width="4"/>',
                txt(681, 539, 'Fe', 27), txt(791, 539, 'Al', 27)]
    elif kind == 'poles':
        out += [magnet(160, 365), magnet(600, 365), txt(210, 413, 'N'), txt(310, 413, 'S'),
                txt(650, 413, 'N'), txt(750, 413, 'S'),
                '<path d="M380 402H540m-25-20 25 20-25 20" stroke="#df8b62" stroke-width="7" fill="none"/>',
                '<path d="M510 485H390m25-20-25 20 25 20" stroke="#df8b62" stroke-width="7" fill="none"/>']
    elif kind == 'compass':
        out += ['<circle cx="500" cy="445" r="170" fill="#e8f7fb" stroke="#87b3c8" stroke-width="7"/>',
                '<path d="M500 287V603 M342 445H658" stroke="#b4d3dc" stroke-width="4"/>',
                '<path d="M500 315L466 460 500 443Z" fill="#ee886e"/><path d="M500 575L534 430 500 447Z" fill="#559bbd"/>',
                txt(500, 255, 'N', 31), txt(500, 666, 'S', 31),
                '<rect x="760" y="330" width="85" height="210" rx="9" fill="#b7cbd4" stroke="#7594a2" stroke-width="4"/>']
    elif kind == 'wire':
        out += [f'<path d="M175 430H850" {wire}/>',
                '<circle cx="500" cy="430" r="120" stroke="#91c7d4" stroke-width="6" fill="none"/>',
                '<circle cx="500" cy="430" r="190" stroke="#c2dde4" stroke-width="5" fill="none"/>',
                '<path d="M600 355l25 18-4-30" fill="#e78960"/>',
                '<circle cx="695" cy="625" r="55" fill="#e9f5f8" stroke="#7eb1c0" stroke-width="5"/>',
                '<path d="M695 588l-20 47 20-8z" fill="#ee8b69"/>']
    elif kind == 'coil':
        out += ['<rect x="260" y="350" width="420" height="165" rx="12" fill="#c7d3da" stroke="#899caa" stroke-width="5"/>',
                *[f'<ellipse cx="{305+i*44}" cy="430" rx="21" ry="103" fill="none" stroke="#cf8657" stroke-width="9"/>' for i in range(8)],
                '<path d="M225 525v85h520v-85" stroke="#397595" stroke-width="7" fill="none"/>',
                '<rect x="420" y="600" width="135" height="36" rx="5" fill="#f1c87a" stroke="#aa8247"/>']
    elif kind == 'motor':
        out += [magnet(160, 410), magnet(670, 410),
                '<circle cx="520" cy="447" r="137" fill="#f8fbfc" stroke="#8caec0" stroke-width="6"/>',
                '<rect x="455" y="382" width="130" height="130" rx="9" fill="#fae8bf" stroke="#b8744d" stroke-width="6"/>',
                '<path d="M510 300q145 50 100 195" stroke="#df8a5e" stroke-width="8" fill="none"/>',
                '<path d="M610 495l-21-27 35 6" fill="#df8a5e"/>', txt(215, 459, 'N'), txt(825, 459, 'S')]
    elif kind == 'induction':
        out += [magnet(165, 400), txt(215, 448, 'N'), txt(315, 448, 'S'),
                '<path d="M405 438h116m-25-18 25 18-25 18" stroke="#e28e62" stroke-width="7" fill="none"/>',
                *[f'<ellipse cx="{560+i*42}" cy="438" rx="20" ry="118" fill="none" stroke="#c98656" stroke-width="8"/>' for i in range(6)],
                '<path d="M560 559v70h210v-70" stroke="#38748c" stroke-width="6" fill="none"/>',
                '<circle cx="665" cy="625" r="44" fill="#fff" stroke="#75aabb" stroke-width="5"/>']
    elif kind == 'generator':
        out += ['<circle cx="390" cy="445" r="144" fill="#f8fbfc" stroke="#91b5c6" stroke-width="7"/>',
                '<path d="M390 445l90-90 M390 445l-75 100" stroke="#d17f55" stroke-width="12"/>',
                '<path d="M540 445h102" stroke="#367391" stroke-width="12"/>',
                '<rect x="642" y="340" width="165" height="210" rx="18" fill="#d8eff3" stroke="#73a5b9" stroke-width="7"/>',
                '<path d="M724 550v75" stroke="#367391" stroke-width="8"/>',
                '<circle cx="724" cy="648" r="32" fill="#ffe9a1" stroke="#d39c48" stroke-width="6"/>']
    elif kind == 'signal':
        out += ['<rect x="155" y="350" width="175" height="145" rx="16" fill="#e7f5fa" stroke="#73a5b9" stroke-width="5"/>',
                '<rect x="710" y="350" width="175" height="145" rx="16" fill="#e7f5fa" stroke="#73a5b9" stroke-width="5"/>',
                f'<path d="M330 422h380" {wire}/>', txt(240, 431, '10 01', 29), txt(797, 431, '1001', 29),
                *[f'<rect x="{380+i*73}" y="{370 if v else 420}" width="52" height="{52 if v else 12}" fill="#e9a66d"/>' for i,v in enumerate([1,0,0,1])]]
    elif kind == 'capstone':
        # Top-down model: the reed stays on the frame; the door magnet moves.
        # The two labels below each drawing are the first two localized notes.
        out += ['<path d="M170 335H425V580 M575 335H830V580" stroke="#729baa" stroke-width="10" fill="none"/>',
                '<circle cx="170" cy="335" r="13" fill="#497b94"/><circle cx="575" cy="335" r="13" fill="#497b94"/>',
                '<path d="M170 335H405" stroke="#b7c3cb" stroke-width="22"/>',
                '<path d="M575 335V555" stroke="#b7c3cb" stroke-width="22"/>',
                '<rect x="380" y="297" width="36" height="25" rx="4" fill="#ee9878"/>',
                '<rect x="555" y="523" width="36" height="25" rx="4" fill="#ee9878"/>',
                '<rect x="416" y="297" width="25" height="25" rx="4" fill="#89b9d0"/>',
                '<rect x="821" y="297" width="25" height="25" rx="4" fill="#89b9d0"/>',
                '<circle cx="315" cy="490" r="33" fill="#d3e1e4" stroke="#91aab2" stroke-width="5"/>',
                '<circle cx="725" cy="490" r="33" fill="#ffe29a" stroke="#d19a4b" stroke-width="5"/>',
                txt(300, 685, NOTES['capstone'][lang][0], 25),
                txt(700, 685, NOTES['capstone'][lang][1], 25)]
    return out


def svg(row, lang):
    title = row[lang][0]
    notes = NOTES[row['kind']][lang]
    title_lines = textwrap.wrap(title, 47, break_long_words=False)
    if len(title_lines) > 2:
        raise ValueError(f'title too long on map: {title}')
    out = ['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
           f'<title id="title">{escape(title)}</title><desc id="desc">{escape("; ".join(notes))}</desc>',
           '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#d4e8ed"/></pattern><linearGradient id="paper" x2="1" y2="1"><stop stop-color="#ffffff"/><stop offset="1" stop-color="#e9f4f6"/></linearGradient><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="12" stdDeviation="10" flood-color="#37657b" flood-opacity=".23"/></filter></defs>',
           '<rect width="1600" height="900" fill="#f1f9fb"/><rect width="1600" height="900" fill="url(#grid)"/>',
           '<rect x="70" y="195" width="875" height="590" rx="18" fill="url(#paper)" stroke="#99bbc8" stroke-width="3" filter="url(#shadow)"/>',
           '<rect x="975" y="195" width="555" height="590" rx="18" fill="#fff" stroke="#99bbc8" stroke-width="3" filter="url(#shadow)"/>']
    for i, line in enumerate(title_lines):
        out.append(txt(75, 76 + i*50, line, 35, 'start'))
    out.extend(scene(row['kind'], lang))
    for i, note in enumerate(notes):
        y = 224+i*135
        out.append(f'<rect x="1005" y="{y}" width="495" height="105" rx="11" fill="{("#e6f5fa", "#fff0e4", "#eaf4e9", "#edf0fa")[i]}" stroke="#d0e2e8" stroke-width="2"/>')
        lines = textwrap.wrap(note, 27, break_long_words=False)
        if len(lines) > 2:
            raise ValueError(f'note too long: {note}')
        for j, line in enumerate(lines):
            out.append(txt(1027, y+45+j*32, line, 28, 'start'))
    out.append('</svg>')
    return '\n'.join(out)+'\n'


def markdown(row, index, lang):
    title, summary, scene_text, terms, worked, task, trap, transfer = row[lang]
    sections = [scene_text, terms, worked, task, trap,
                f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)',
                transfer, CHECK[lang], JOURNAL[lang]]
    body = '\n\n'.join(f'## {heading}\n\n{content}' for heading, content in zip(HEAD[lang], sections))
    if index < 66:
        next_row = MAGNETISM[index-56]
        label = {'ru':'Следующий урок','kz':'Келесі сабақ','en':'Next lesson'}[lang]
        body += f'\n\n[{label}: {next_row[lang][0]}](/read/physics-{next_row["stem"]}?lang={lang})'
    else:
        label = {'ru':'Следующий урок','kz':'Келесі сабақ','en':'Next lesson'}[lang]
        body += f'\n\n[{label}](/read/physics-67-disturbance-and-wave?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'


def build():
    if len(MAGNETISM) != 10:
        raise ValueError('magnetism block requires ten lessons')
    for number, row in enumerate(MAGNETISM, 57):
        if not row['stem'].startswith(f'{number:02d}-'):
            raise ValueError(f'wrong lesson number: {row["stem"]}')
        for lang in ('ru', 'kz', 'en'):
            suffix = '' if lang == 'ru' else '-'+lang
            (LESSONS / f'{row["stem"]}{suffix}.md').write_text(markdown(row, number, lang), encoding='utf-8')
            (MAPS / f'map-{row["stem"]}-{lang}.svg').write_text(svg(row, lang), encoding='utf-8')
    print('Built 30 magnetism lessons and 30 diagrams')


if __name__ == '__main__':
    build()
