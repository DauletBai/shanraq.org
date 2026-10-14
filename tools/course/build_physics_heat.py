#!/usr/bin/env python3
"""Build the trilingual heat block (37–46) and its exact-data support maps."""
from html import escape
from pathlib import Path

try:
    from .physics_heat_data import HEAT
except ImportError:
    from physics_heat_data import HEAT

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/physics"
MAPS = ROOT / "web/static/course/physics"

HEAD = {
    "ru": ("Вопрос из жизни и прогноз", "Разбираем новые слова", "Разбор на числах", "Проверьте безопасно", "Ошибка, которую легко сделать", "Опорная карта", "Теперь сами", "Запись в лабораторный журнал"),
    "kz": ("Өмірден алынған сұрақ және болжам", "Жаңа сөздерді түсінейік", "Сандармен талдау", "Қауіпсіз тексеріңіз", "Жиі кездесетін қате", "Тірек картасы", "Енді өзіңіз", "Зертхана журналына жазба"),
    "en": ("A real-life question and prediction", "Words to understand", "A worked example", "Test it safely", "An easy mistake to make", "Support map", "Your turn", "Lab notebook entry"),
}
COMMON = {
    "ru": (
        "Перед объяснением остановитесь на минуту: что именно вы ожидаете увидеть и почему? Запишите прогноз до расчёта. После урока отдельно отметьте, что вы наблюдали, что получили из заданных чисел и что лишь предполагаете о невидимых частицах. Такая запись позволяет исправить модель, не подменяя опыт красивым ответом.",
        "Закройте готовое решение. Нарисуйте систему: что находится внутри мысленной границы, а что снаружи? Подпишите стрелки передачи энергии и единицы величин. Измените в условии только одно число и объясните словами направление изменения результата. Если опыта не было, не пишите «измерено»: пишите «вычислено по учебным данным».",
        "Оформите четыре строки: вопрос; первоначальный прогноз; наблюдение или заданные числа; исправленный вывод с условиями применимости. Назовите один возможный источник ошибки: неточный термометр, потерю энергии в комнате, изменение количества воды или другое конкретное условие. В следующем уроке вернитесь к этому выводу и проверьте, выдерживает ли он новую ситуацию."),
    "kz": (
        "Түсіндіруге дейін бір минут тоқтаңыз: нені көресіз деп күтесіз және неліктен? Есептемей тұрып болжамды жазыңыз. Сабақтан кейін нені байқағаныңызды, берілген сандардан нені есептегеніңізді және көзге көрінбейтін бөлшектер жайында нені тек болжағаныңызды бөлек көрсетіңіз. Осылайша әдемі жауаппен тәжірибені алмастырмай, модельді түзете аласыз.",
        "Дайын шешімді жауып қойыңыз. Жүйенің шекарасын сызыңыз: ішінде не, сыртында не бар? Энергия берілуінің бағыттарын және шамалардың бірліктерін белгілеңіз. Шарттағы тек бір санды өзгертіп, нәтиженің қай бағытта өзгеретінін сөзбен түсіндіріңіз. Тәжірибе жасалмаса, «өлшенді» демей, «оқу деректері бойынша есептелді» деп жазыңыз.",
        "Төрт жол толтырыңыз: сұрақ; бастапқы болжам; бақылау не берілген сандар; қолданылу шарттары бар түзетілген қорытынды. Қатенің бір нақты көзін атаңыз: дәл емес термометр, бөлмеге кеткен энергия, су мөлшерінің өзгеруі немесе басқа шарт. Келесі сабақта осы қорытындыға қайта оралып, жаңа жағдайға жарай ма, тексеріңіз."),
    "en": (
        "Pause before the explanation: what exactly do you expect to see, and why? Write a prediction before calculating. After the lesson, separate what you observed, what you calculated from supplied numbers, and what you only inferred about invisible particles. That record lets you correct a model without replacing evidence with a neat answer.",
        "Cover the worked solution. Draw a boundary around the system: what is inside and what is outside? Label energy-flow arrows and the units of each quantity. Change just one supplied number and explain in words which way the result changes. If there was no experiment, write “calculated from teaching data”, not “measured”.",
        "Write four lines: question; initial prediction; observation or supplied values; revised conclusion with its conditions. Identify one concrete source of error, such as a coarse thermometer, energy escaping into the room, or a changed amount of water. Return to the conclusion after the next lesson and see whether it survives a new situation."),
}

def wrap(value, limit):
    result, line = [], ""
    for word in value.split():
        next_line = f"{line} {word}".strip()
        if line and len(next_line) > limit:
            result.append(line)
            line = word
        else:
            line = next_line
    if line:
        result.append(line)
    return result

def text(x, y, value, size=28, color="#173d58", anchor="start"):
    return f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" font-weight="650" font-size="{size}" fill="{color}">{escape(value)}</text>'

def svg(row, lang):
    title, labels = row[lang][0], row["map"][lang]
    if len(labels) != 4 or any(len(wrap(label, 24)) > 2 for label in labels):
        raise ValueError(f"Map label will not fit: {row['stem']} {lang}")
    out = [
        '<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
        f'<title id="title">{escape(title)}</title>',
        f'<desc id="desc">{escape(" | ".join(labels))}</desc>',
        '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#cfebf0"/></pattern><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="13" stdDeviation="12" flood-color="#375b74" flood-opacity=".2"/></filter><linearGradient id="card" x2="1" y2="1"><stop stop-color="#fff"/><stop offset="1" stop-color="#e0eef3"/></linearGradient></defs>',
        '<rect width="1600" height="900" fill="#eef8f8"/><rect width="1600" height="900" fill="url(#grid)"/>',
        '<rect x="75" y="200" width="900" height="550" rx="24" fill="url(#card)" stroke="#9bbecb" stroke-width="3" filter="url(#shadow)"/>',
        '<rect x="1015" y="200" width="510" height="550" rx="24" fill="#fff" stroke="#9bbecb" stroke-width="3" filter="url(#shadow)"/>',
    ]
    title_lines = wrap(title, 32)
    if len(title_lines) > 3:
        raise ValueError(f"Title will not fit: {row['stem']} {lang}")
    for i, line in enumerate(title_lines):
        out.append(text(76, 72 + i * 48, line, 34))
    for i, label in enumerate(labels):
        y = 235 + i * 122
        out.append(f'<rect x="1044" y="{y}" width="450" height="94" rx="14" fill="{("#e4f5fa", "#fff0df", "#e9f5ee", "#e9edfa")[i]}" stroke="#d3e4e9"/>')
        for j, line in enumerate(wrap(label, 24)):
            out.append(text(1064, y + 37 + j * 32, line, 28 if len(line) < 25 else 26))
    # The left drawing represents the physical process, while the right cards
    # carry the exact teaching data. Keep labels outside drawn objects.
    kind = row["kind"]
    if kind == "particles":
        out += ['<path d="M195 320V575H455V320 M590 320V575H850V320" fill="none" stroke="#487f96" stroke-width="10"/>',
                '<path d="M201 400H449V569H201Z M596 400H844V569H596Z" fill="#92d1e4" opacity=".65"/>',
                '<path d="M201 400H449 M596 400H844" stroke="#fff" stroke-width="6"/>']
        for i in range(18):
            x = 220 + (i % 6) * 39 + (i % 2) * 5
            y = 430 + (i // 6) * 55
            out.append(f'<circle cx="{x}" cy="{y}" r="9" fill="#3285a6"/>')
            out.append(f'<circle cx="{x+395}" cy="{y}" r="9" fill="#ee8758"/>')
        out += ['<path d="M645 355l33-30m40 65 35-23m-122 123 25-21" stroke="#e88255" stroke-width="5" stroke-linecap="round"/>']
    elif kind == "temperature":
        out += ['<rect x="383" y="305" width="170" height="294" rx="72" fill="#f9ffff" stroke="#4c829a" stroke-width="8"/>',
                '<rect x="445" y="355" width="50" height="204" rx="25" fill="#ffe6d0" stroke="#4c829a" stroke-width="5"/>',
                '<rect x="456" y="452" width="28" height="96" rx="14" fill="#e87859"/>',
                '<circle cx="470" cy="551" r="44" fill="#e87859" stroke="#fff" stroke-width="5"/>']
        for i in range(9):
            y = 380 + i * 23
            out.append(f'<path d="M507 {y}h{45 if i%2==0 else 26}" stroke="#4c829a" stroke-width="4"/>')
        out += ['<path d="M575 459h115m-28-23 28 23-28 23" stroke="#e87859" fill="none" stroke-width="7" stroke-linejoin="round"/>']
    elif kind == "internal":
        out += ['<rect x="360" y="375" width="335" height="195" rx="22" fill="#cbeaf1" stroke="#4c829a" stroke-width="8"/>',
                '<path d="M190 427h145m-28-24 28 24-28 24 M190 521h145m-28-24 28 24-28 24 M725 475h145m-28-24 28 24-28 24" fill="none" stroke="#e98656" stroke-width="9" stroke-linejoin="round"/>']
        out += [text(260,395,'50 J' if lang=='en' else '50 Дж',28,anchor='middle'),
                text(260,492,'20 J' if lang=='en' else '20 Дж',28,anchor='middle'),
                text(800,445,'30 J' if lang=='en' else '30 Дж',28,anchor='middle'),
                text(525,486,'ΔU = +40 J' if lang=='en' else 'ΔU = +40 Дж',36,anchor='middle')]
    elif kind == "conduction":
        out += ['<rect x="230" y="370" width="620" height="145" rx="24" fill="#a8c7d5" stroke="#507b91" stroke-width="7"/>',
                '<rect x="230" y="370" width="160" height="145" rx="24" fill="#ec986a"/>',
                '<rect x="690" y="370" width="160" height="145" rx="24" fill="#78b7d9"/>',
                '<path d="M408 443h250m-28-24 28 24-28 24" fill="none" stroke="#f0794d" stroke-width="10" stroke-linejoin="round"/>']
        out += [text(310,340,'30 °C',31,anchor='middle'),text(770,340,'20 °C',31,anchor='middle')]
    elif kind == "convection":
        out += ['<path d="M205 310V605H860V310" fill="none" stroke="#4e7183" stroke-width="12"/>',
                '<rect x="240" y="480" width="120" height="97" rx="12" fill="#ec9469" stroke="#a95c47" stroke-width="5"/>',
                '<path d="M400 525c-5-95 80-170 190-160s175 80 145 150c-17 45-70 58-144 55H415" fill="none" stroke="#e88356" stroke-width="10" stroke-linecap="round"/>',
                '<path d="m442 454-42 71-18-77 M493 548l-78 22 55 52" fill="none" stroke="#e88356" stroke-width="8" stroke-linejoin="round"/>',
                '<path d="M585 357c75 8 119 38 145 84" fill="none" stroke="#53a7bd" stroke-width="5" stroke-dasharray="11 11"/>']
        out += [text(305,640,'25 °C',30,anchor='middle'),text(760,640,'18 °C',30,anchor='middle')]
    elif kind == "radiation":
        out += ['<circle cx="315" cy="418" r="84" fill="#ffd36b" stroke="#e9a246" stroke-width="8"/>',
                '<path d="M485 355l315-75 M487 415h300 M485 475l315 75" fill="none" stroke="#f1b351" stroke-width="11" stroke-linecap="round"/>',
                '<path d="M720 300v280" stroke="#607e94" stroke-width="17"/>',
                '<rect x="775" y="500" width="90" height="60" rx="9" fill="#9ccadb" stroke="#4a8599" stroke-width="5"/>']
    elif kind == "capacity":
        out += ['<path d="M210 350V565H440V350 M570 350V565H830V350" fill="none" stroke="#4a849a" stroke-width="10"/>',
                '<path d="M216 430H434V559H216Z M576 390H824V559H576Z" fill="#8ccde2"/>',
                '<path d="M216 430H434 M576 390H824" stroke="#fff" stroke-width="6"/>']
        out += [text(325,615,'0,2 кг' if lang!='en' else '0.2 kg',30,anchor='middle'),text(700,615,'0,4 кг' if lang!='en' else '0.4 kg',30,anchor='middle')]
    elif kind == "phase":
        out += ['<rect x="230" y="390" width="170" height="170" rx="8" fill="#d7f2fa" stroke="#64a7be" stroke-width="7"/>',
                '<path d="M435 475h190m-28-25 28 25-28 25" fill="none" stroke="#e88e62" stroke-width="10" stroke-linejoin="round"/>',
                '<path d="M660 370V565H840V370 M666 455H834V559H666Z" fill="#8ccde2" stroke="#4a849a" stroke-width="8"/>']
        out += [text(315,355,'0 °C',30,anchor='middle'),text(750,345,'0 °C',30,anchor='middle')]
    elif kind in {"budget", "capstone"}:
        out += ['<rect x="340" y="330" width="420" height="278" rx="8" fill="#f3f9fb" stroke="#547f94" stroke-width="9"/>',
                '<rect x="592" y="360" width="118" height="94" fill="#a4d5e6" stroke="#547f94" stroke-width="6"/>',
                '<path d="M190 480h135m-26-24 26 24-26 24 M765 405h135m-26-24 26 24-26 24 M765 550h135m-26-24 26 24-26 24" fill="none" stroke="#df895f" stroke-width="9" stroke-linejoin="round"/>',
                '<path d="M355 595h380" stroke="#769aaa" stroke-width="6"/>']
        out += [text(257,449,'600',31,anchor='middle'),text(837,376,'150' if kind=='capstone' else '200',31,anchor='middle'),text(837,519,'100',31,anchor='middle')]
    out.append(text(525, 680, row["caption"][lang], 28, "#315b72", "middle"))
    out.append('</svg>')
    return '\n'.join(out) + '\n'

def markdown(number, row, lang):
    title, summary, story, terms, worked, activity, mistake = row[lang]
    headings = HEAD[lang]
    prediction, transfer, notebook = COMMON[lang]
    body = [
        (headings[0], f"{story}\n\n{prediction}"),
        (headings[1], terms),
        (headings[2], worked),
        (headings[3], activity),
        (headings[4], mistake),
        (headings[5], f"![{title}](/static/course/physics/map-{row['stem']}-{lang}.svg)"),
        (headings[6], transfer),
        (headings[7], notebook),
    ]
    result = '\n\n'.join(f'## {heading}\n\n{content}' for heading, content in body)
    if number < 46:
        next_row = HEAT[number - 36]
        label = {"ru": "Следующий урок", "kz": "Келесі сабақ", "en": "Next lesson"}[lang]
        result += f'\n\n[{label}: {next_row[lang][0]}](/read/physics-{next_row["stem"]}?lang={lang})'
    else:
        label = {"ru": "Следующий урок", "kz": "Келесі сабақ", "en": "Next lesson"}[lang]
        title = {"ru": "Почему шарик прилипает к стене?", "kz": "Шар неге қабырғаға жабысады?", "en": "Why does a balloon stick to a wall?"}[lang]
        result += f'\n\n[{label}: {title}](/read/physics-47-electric-charge?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{result}\n'

def build():
    if len(HEAT) != 10:
        raise ValueError("Heat block needs exactly ten lessons")
    for number, row in enumerate(HEAT, 37):
        if not row["stem"].startswith(f"{number:02d}-"):
            raise ValueError(f"Wrong lesson number: {row['stem']}")
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else "-" + lang
            (LESSONS / f"{row['stem']}{suffix}.md").write_text(markdown(number, row, lang), encoding="utf-8")
            (MAPS / f"map-{row['stem']}-{lang}.svg").write_text(svg(row, lang), encoding="utf-8")
    print("Built 30 heat lessons and 30 localized exact-data maps")

if __name__ == "__main__":
    build()
