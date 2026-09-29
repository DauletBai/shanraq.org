#!/usr/bin/env python3
"""Generate the localized support maps for the opening Informatics block."""
from pathlib import Path
from xml.sax.saxutils import escape
from generate_informatics_foundations import ROWS as FOUNDATION_ROWS


ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/informatics"

LANG = {
    "ru": {
        "course": "ИНФОРМАТИКА И ИИ",
        "subtitle": "Один проект растёт от замысла до проверяемого продукта",
        "steps": [
            ("1", "Компьютер и проект", "уроки 1–10 • версия 0.1"),
            ("2", "Данные", "уроки 11–18 • версия 0.2"),
            ("3", "Алгоритмы", "уроки 19–26 • версия 0.3"),
            ("4", "Python", "уроки 27–38 • версия 1.0"),
            ("5", "Веб и облако", "уроки 39–46 • версия 1.1"),
            ("6", "SQLite", "уроки 47–54 • версия 1.2"),
            ("7", "Безопасность", "уроки 55–62 • версия 1.3"),
            ("8", "ИИ и выпуск", "уроки 63–72 • версия 2.1"),
        ],
        "project": "МОЙ ЦИФРОВОЙ ПОМОЩНИК",
        "project_sub": "каждый блок улучшает тот же продукт",
        "diag_title": "ДИАГНОСТИКА НЕ СТАВИТ ОТМЕТКУ",
        "diag_subtitle": "Она выбирает объём поддержки и тему проекта",
        "diag_groups": [
            ("1–3", "УСТРОЙСТВО И ДАННЫЕ", "файл • память • размер"),
            ("4–6", "АЛГОРИТМ", "шаг • состояние • граница"),
            ("7–8", "ОШИБКА И ЗАЩИТА", "факт • ссылка • решение"),
            ("9–10", "ДАННЫЕ И ИИ", "честный вывод • проверка"),
        ],
        "result": "результат → подсказки урока → паспорт версии 0.0",
        "marks": "✓ / ? / пока не знаю",
    },
    "kz": {
        "course": "ИНФОРМАТИКА ЖӘНЕ ЖИ",
        "subtitle": "Бір жоба ойдан тексерілетін өнімге дейін дамиды",
        "steps": [
            ("1", "Компьютер және жоба", "1–10 сабақ • 0.1 нұсқасы"),
            ("2", "Деректер", "11–18 сабақ • 0.2 нұсқасы"),
            ("3", "Алгоритмдер", "19–26 сабақ • 0.3 нұсқасы"),
            ("4", "Python", "27–38 сабақ • 1.0 нұсқасы"),
            ("5", "Веб және бұлт", "39–46 сабақ • 1.1 нұсқасы"),
            ("6", "SQLite", "47–54 сабақ • 1.2 нұсқасы"),
            ("7", "Қауіпсіздік", "55–62 сабақ • 1.3 нұсқасы"),
            ("8", "ЖИ және шығарылым", "63–72 сабақ • 2.1 нұсқасы"),
        ],
        "project": "МЕНІҢ ЦИФРЛЫҚ КӨМЕКШІМ",
        "project_sub": "әр бөлім сол өнімді жақсартады",
        "diag_title": "ДИАГНОСТИКА БАҒА ҚОЙМАЙДЫ",
        "diag_subtitle": "Ол қолдау көлемі мен жоба тақырыбын анықтайды",
        "diag_groups": [
            ("1–3", "ҚҰРЫЛҒЫ ЖӘНЕ ДЕРЕК", "файл • жад • көлем"),
            ("4–6", "АЛГОРИТМ", "қадам • күй • шекара"),
            ("7–8", "ҚАТЕ ЖӘНЕ ҚОРҒАНЫС", "дерек • сілтеме • шешім"),
            ("9–10", "ДЕРЕК ЖӘНЕ ЖИ", "адал қорытынды • тексеру"),
        ],
        "result": "нәтиже → сабақ көмегі → 0.0 нұсқасының паспорты",
        "marks": "✓ / ? / әзірге білмеймін",
    },
    "en": {
        "course": "INFORMATICS AND AI",
        "subtitle": "One project grows from an idea into a verifiable product",
        "steps": [
            ("1", "Computer and project", "lessons 1–10 • version 0.1"),
            ("2", "Data", "lessons 11–18 • version 0.2"),
            ("3", "Algorithms", "lessons 19–26 • version 0.3"),
            ("4", "Python", "lessons 27–38 • version 1.0"),
            ("5", "Web and cloud", "lessons 39–46 • version 1.1"),
            ("6", "SQLite", "lessons 47–54 • version 1.2"),
            ("7", "Security", "lessons 55–62 • version 1.3"),
            ("8", "AI and release", "lessons 63–72 • version 2.1"),
        ],
        "project": "MY DIGITAL ASSISTANT",
        "project_sub": "every block improves the same product",
        "diag_title": "DIAGNOSIS DOES NOT GIVE A GRADE",
        "diag_subtitle": "It chooses the amount of support and the project theme",
        "diag_groups": [
            ("1–3", "DEVICE AND DATA", "file • memory • size"),
            ("4–6", "ALGORITHM", "step • state • boundary"),
            ("7–8", "ERROR AND SAFETY", "fact • link • decision"),
            ("9–10", "DATA AND AI", "honest claim • verification"),
        ],
        "result": "result → lesson support → version 0.0 passport",
        "marks": "✓ / ? / not yet",
    },
}

LESSON_TEXT = {
    "ru": {
        "system_title": "ТЕЛЕФОН И НОУТБУК — СИСТЕМЫ",
        "system_sub": "Разные корпуса, одна проверяемая модель",
        "system_flow": ("цель", "вход", "обработка", "состояние", "выход"),
        "phone": ("ТЕЛЕФОН", "касание • камера", "процессор • ОС", "ОЗУ • хранилище", "экран • звук"),
        "laptop": ("НОУТБУК", "клавиатура • камера", "процессор • ОС", "ОЗУ • накопитель", "экран • звук"),
        "system_signal": "часть полезна не названием, а ролью в общей задаче",
        "io_title": "ОТ ФИЗИЧЕСКОГО СИГНАЛА К ДАННЫМ И ОБРАТНО",
        "io_sub": "Устройство измеряет мир, программа изменяет данные, выход воздействует на мир",
        "io_flow": ("явление", "датчик", "данные", "программа", "выход"),
        "io_examples": (
            ("свет", "камера", "значения пикселей", "обработка фото", "экран"),
            ("касание", "сенсорный слой", "x=412, y=728", "проверка кнопки", "новый экран"),
            ("звук", "микрофон", "отсчёты сигнала", "кодирование", "динамик"),
        ),
        "noise": "измерение ≠ само явление • возможны шум, задержка и ошибка",
    },
    "kz": {
        "system_title": "ТЕЛЕФОН МЕН НОУТБУК — ЖҮЙЕЛЕР",
        "system_sub": "Корпустары бөлек, тексерілетін моделі ортақ",
        "system_flow": ("мақсат", "кіріс", "өңдеу", "күй", "шығыс"),
        "phone": ("ТЕЛЕФОН", "түрту • камера", "процессор • ОЖ", "жедел жад • сақтау", "экран • дыбыс"),
        "laptop": ("НОУТБУК", "пернетақта • камера", "процессор • ОЖ", "жедел жад • жинақтауыш", "экран • дыбыс"),
        "system_signal": "бөліктің пайдасын атауы емес, ортақ міндеттегі рөлі анықтайды",
        "io_title": "ФИЗИКАЛЫҚ СИГНАЛДАН ДЕРЕККЕ ЖӘНЕ КЕРІ",
        "io_sub": "Құрылғы әлемді өлшейді, бағдарлама деректі өзгертеді, шығыс әлемге әсер етеді",
        "io_flow": ("құбылыс", "сенсор", "дерек", "бағдарлама", "шығыс"),
        "io_examples": (
            ("жарық", "камера", "пиксель мәндері", "фотоны өңдеу", "экран"),
            ("түрту", "сенсорлық қабат", "x=412, y=728", "батырманы тексеру", "жаңа экран"),
            ("дыбыс", "микрофон", "сигнал есептері", "кодтау", "динамик"),
        ),
        "noise": "өлшем ≠ құбылыстың өзі • шу, кідіріс және қате болуы мүмкін",
    },
    "en": {
        "system_title": "PHONES AND LAPTOPS ARE SYSTEMS",
        "system_sub": "Different cases, one verifiable model",
        "system_flow": ("purpose", "input", "processing", "state", "output"),
        "phone": ("PHONE", "touch • camera", "processor • OS", "RAM • storage", "display • sound"),
        "laptop": ("LAPTOP", "keyboard • camera", "processor • OS", "RAM • drive", "display • sound"),
        "system_signal": "a part matters through its role in the shared task, not through its name",
        "io_title": "FROM A PHYSICAL SIGNAL TO DATA AND BACK",
        "io_sub": "A device measures the world, software changes data, and output affects the world",
        "io_flow": ("event", "sensor", "data", "program", "output"),
        "io_examples": (
            ("light", "camera", "pixel values", "photo processing", "display"),
            ("touch", "touch layer", "x=412, y=728", "button test", "new screen"),
            ("sound", "microphone", "signal samples", "encoding", "speaker"),
        ),
        "noise": "a measurement ≠ the event itself • noise, delay, and error are possible",
    },
}

FOUNDATION_STAGES = {
    5: {
        "ru": ("накопитель", "ОЗУ", "процессор", "сохранение", "накопитель"),
        "kz": ("сақтау құрылғысы", "жедел жад", "процессор", "сақтау", "сақталған файл"),
        "en": ("storage", "RAM", "processor", "save", "stored file"),
    },
    6: {
        "ru": ("файл программы", "ОС", "процесс", "ресурс", "результат"),
        "kz": ("бағдарлама файлы", "ОЖ", "процесс", "ресурс", "нәтиже"),
        "en": ("program file", "OS", "process", "resource", "result"),
    },
    7: {
        "ru": ("assistant/", "data/tasks.json", "docs/passport.md", "tests/"),
        "kz": ("assistant/", "data/tasks.json", "docs/passport.md", "tests/"),
        "en": ("assistant/", "data/tasks.json", "docs/passport.md", "tests/"),
    },
    8: {
        "ru": ("байты + формат", "кодировка", "источник + лицензия", "разрешённое применение"),
        "kz": ("байт + пішім", "кодтау", "дереккөз + лицензия", "рұқсат етілген қолдану"),
        "en": ("bytes + format", "encoding", "source + licence", "permitted use"),
    },
    9: {
        "ru": ("малое изменение", "сравнение", "проверка", "тест", "версия", "восстановление"),
        "kz": ("шағын өзгеріс", "салыстыру", "шолу", "сынақ", "нұсқа", "қалпына келтіру"),
        "en": ("small change", "compare", "review", "test", "version", "recover"),
    },
    10: {
        "ru": ("объяснить", "модель", "найти ошибку", "перенос", "доказательство", "8/10 сейчас"),
        "kz": ("түсіндіру", "модель", "қатені табу", "көшіру", "дәлел", "қазір 8/10"),
        "en": ("explain", "model", "diagnose", "transfer", "evidence", "8/10 now"),
    },
}


STYLE = """
  <defs>
    <linearGradient id="paper" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#f8fbff"/><stop offset="1" stop-color="#e8f1fb"/>
    </linearGradient>
    <linearGradient id="blue" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#54b8ff"/><stop offset="1" stop-color="#1765b5"/>
    </linearGradient>
    <linearGradient id="cyan" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#72e0dd"/><stop offset="1" stop-color="#138c9e"/>
    </linearGradient>
    <linearGradient id="violet" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#b89cff"/><stop offset="1" stop-color="#6742ba"/>
    </linearGradient>
    <linearGradient id="orange" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#ffcb6b"/><stop offset="1" stop-color="#de7b22"/>
    </linearGradient>
    <filter id="shadow" x="-20%" y="-20%" width="150%" height="160%">
      <feDropShadow dx="0" dy="12" stdDeviation="10" flood-color="#17324f" flood-opacity=".22"/>
    </filter>
    <marker id="arrow" viewBox="0 0 12 12" refX="10" refY="6" markerWidth="9" markerHeight="9" orient="auto-start-reverse">
      <path d="M1 1 11 6 1 11Z" fill="#d53d3d"/>
    </marker>
    <pattern id="grid" width="32" height="32" patternUnits="userSpaceOnUse">
      <path d="M32 0H0V32" fill="none" stroke="#bfd2e6" stroke-width="1" opacity=".42"/>
    </pattern>
  </defs>
  <style>
    .title{font:800 42px Inter,Arial,sans-serif;fill:#17324f;letter-spacing:1px}
    .subtitle{font:500 22px Inter,Arial,sans-serif;fill:#4d647a}
    .step-title{font:800 23px Inter,Arial,sans-serif;fill:#fff}
    .step-sub{font:600 17px Inter,Arial,sans-serif;fill:#eaf7ff}
    .number{font:900 34px Inter,Arial,sans-serif;fill:#17324f}
    .project{font:900 27px Inter,Arial,sans-serif;fill:#17324f;letter-spacing:.5px}
    .project-sub{font:600 18px Inter,Arial,sans-serif;fill:#49647c}
  </style>
"""


def svg_open(title: str, desc: str) -> str:
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="800" height="450" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">
<title id="title">{escape(title)}</title><desc id="desc">{escape(desc)}</desc>
{STYLE}
<rect width="1600" height="900" fill="url(#paper)"/><rect width="1600" height="900" fill="url(#grid)"/>
'''


def whole_map(lang: str, data: dict) -> str:
    positions = [(85,150),(590,150),(1095,150),(1095,390),(590,390),(85,390),(335,630),(845,630)]
    gradients = ("blue","cyan","violet","orange","blue","cyan","violet","orange")
    out = [svg_open(data["course"], data["subtitle"])]
    out.append(f'<text x="800" y="68" text-anchor="middle" class="title">{escape(data["course"])}</text>')
    out.append(f'<text x="800" y="106" text-anchor="middle" class="subtitle">{escape(data["subtitle"])}</text>')
    centers = []
    for (x,y) in positions:
        centers.append((x+210,y+75))
    for index in range(len(centers)-1):
        x1,y1=centers[index]; x2,y2=centers[index+1]
        if abs(y2-y1) < 10:
            start = x1 + (210 if x2>x1 else -210)
            end = x2 - (210 if x2>x1 else -210)
            path = f'M{start} {y1} H{end}'
        else:
            middle = (y1 + y2) / 2
            path = f'M{x1} {y1+75} V{middle} H{x2} V{y2-75}'
        out.append(f'<path class="course-arrow" d="{path}" fill="none" stroke="#d53d3d" stroke-width="8" stroke-linecap="round" marker-end="url(#arrow)"/>')
    for i, ((num,title,version),(x,y),grad) in enumerate(zip(data["steps"],positions,gradients),1):
        out.append(f'''<g class="course-step" data-step="{i}" filter="url(#shadow)">
  <rect x="{x}" y="{y}" width="420" height="150" rx="22" fill="url(#{grad})" stroke="#fff" stroke-width="3"/>
  <circle cx="{x+48}" cy="{y+45}" r="29" fill="#fff" opacity=".94"/>
  <text x="{x+48}" y="{y+57}" text-anchor="middle" class="number">{num}</text>
  <text x="{x+92}" y="{y+55}" class="step-title">{escape(title)}</text>
  <text x="{x+92}" y="{y+96}" class="step-sub">{escape(version)}</text>
</g>''')
    out.append(f'''<g filter="url(#shadow)">
  <rect x="420" y="806" width="760" height="70" rx="20" fill="#fff" stroke="#76a9d3" stroke-width="2"/>
  <text x="800" y="836" text-anchor="middle" class="project">{escape(data["project"])}</text>
  <text x="800" y="862" text-anchor="middle" class="project-sub">{escape(data["project_sub"])}</text>
</g></svg>''')
    return "\n".join(out)


def diagnostic_map(lang: str, data: dict) -> str:
    out = [svg_open(data["diag_title"], data["diag_subtitle"])]
    out.append(f'<text x="800" y="76" text-anchor="middle" class="title">{escape(data["diag_title"])}</text>')
    out.append(f'<text x="800" y="116" text-anchor="middle" class="subtitle">{escape(data["diag_subtitle"])}</text>')
    colors = ("blue","cyan","violet","orange")
    coords = ((110,190),(860,190),(110,475),(860,475))
    for i, ((span,title,sub),(x,y),color) in enumerate(zip(data["diag_groups"],coords,colors),1):
        out.append(f'''<g class="diagnostic-group" data-questions="{span}" filter="url(#shadow)">
  <rect x="{x}" y="{y}" width="630" height="210" rx="28" fill="url(#{color})" stroke="#fff" stroke-width="4"/>
  <circle cx="{x+72}" cy="{y+65}" r="43" fill="#fff" opacity=".95"/>
  <text x="{x+72}" y="{y+77}" text-anchor="middle" class="number" style="font-size:25px">{span}</text>
  <text x="{x+135}" y="{y+68}" class="step-title">{escape(title)}</text>
  <text x="{x+135}" y="{y+113}" class="step-sub">{escape(sub)}</text>
  <path d="M{x+65} {y+158} H{x+565}" stroke="#fff" stroke-width="3" opacity=".55"/>
  <text x="{x+315}" y="{y+187}" text-anchor="middle" class="step-sub">{escape(data["marks"])}</text>
</g>''')
    out.append('<path d="M800 700 V758" stroke="#d53d3d" stroke-width="8" stroke-linecap="round" marker-end="url(#arrow)"/>')
    out.append(f'''<g filter="url(#shadow)">
  <rect x="300" y="770" width="1000" height="86" rx="24" fill="#fff" stroke="#76a9d3" stroke-width="3"/>
  <text x="800" y="823" text-anchor="middle" class="project">{escape(data["result"])}</text>
</g></svg>''')
    return "\n".join(out)


def system_map(lang: str, data: dict) -> str:
    out = [svg_open(data["system_title"], data["system_sub"])]
    out.append(f'<text x="800" y="68" text-anchor="middle" class="title">{escape(data["system_title"])}</text>')
    out.append(f'<text x="800" y="105" text-anchor="middle" class="subtitle">{escape(data["system_sub"])}</text>')
    xs = (55, 365, 675, 985, 1295)
    colors = ("orange", "blue", "violet", "cyan", "orange")
    for i, (x, label, color) in enumerate(zip(xs, data["system_flow"], colors)):
        out.append(f'<g class="system-role" data-order="{i+1}" filter="url(#shadow)"><rect x="{x}" y="160" width="250" height="105" rx="22" fill="url(#{color})" stroke="#fff" stroke-width="3"/><text x="{x+125}" y="225" text-anchor="middle" class="step-title">{escape(label)}</text></g>')
        if i < 4:
            out.append(f'<path d="M{x+254} 212 H{xs[i+1]-8}" stroke="#d53d3d" stroke-width="7" marker-end="url(#arrow)"/>')
    for panel_y, device, color in ((345, data["phone"], "blue"), (570, data["laptop"], "violet")):
        out.append(f'<g class="device-system" data-device="{escape(device[0].lower())}" filter="url(#shadow)"><rect x="90" y="{panel_y}" width="1420" height="170" rx="28" fill="#fff" stroke="#8db7d8" stroke-width="3"/><rect x="90" y="{panel_y}" width="240" height="170" rx="28" fill="url(#{color})"/><text x="210" y="{panel_y+98}" text-anchor="middle" class="step-title">{escape(device[0])}</text>')
        for j, value in enumerate(device[1:]):
            x = 370 + j * 280
            out.append(f'<rect x="{x}" y="{panel_y+42}" width="245" height="86" rx="18" fill="#edf5fb" stroke="#b8cee0" stroke-width="2"/><text x="{x+122}" y="{panel_y+94}" text-anchor="middle" class="project-sub">{escape(value)}</text>')
        out.append('</g>')
    out.append(f'<text x="800" y="838" text-anchor="middle" class="project">{escape(data["system_signal"])}</text></svg>')
    return "\n".join(out)


def io_map(lang: str, data: dict) -> str:
    out = [svg_open(data["io_title"], data["io_sub"])]
    out.append(f'<text x="800" y="62" text-anchor="middle" class="title" style="font-size:36px">{escape(data["io_title"])}</text>')
    out.append(f'<text x="800" y="100" text-anchor="middle" class="subtitle" style="font-size:20px">{escape(data["io_sub"])}</text>')
    xs = (55, 365, 675, 985, 1295)
    colors = ("orange", "blue", "cyan", "violet", "orange")
    for i, (x, label, color) in enumerate(zip(xs, data["io_flow"], colors)):
        out.append(f'<g class="io-stage" data-order="{i+1}" filter="url(#shadow)"><rect x="{x}" y="140" width="250" height="100" rx="22" fill="url(#{color})" stroke="#fff" stroke-width="3"/><text x="{x+125}" y="202" text-anchor="middle" class="step-title">{escape(label)}</text></g>')
        if i < 4:
            out.append(f'<path d="M{x+254} 190 H{xs[i+1]-8}" stroke="#d53d3d" stroke-width="7" marker-end="url(#arrow)"/>')
    for row, values in enumerate(data["io_examples"]):
        y = 310 + row * 145
        out.append(f'<g class="io-example" data-example="{row+1}" filter="url(#shadow)"><rect x="55" y="{y}" width="1490" height="105" rx="22" fill="#fff" stroke="#9ebed7" stroke-width="2"/>')
        for i, value in enumerate(values):
            x = xs[i]
            out.append(f'<text x="{x+125}" y="{y+63}" text-anchor="middle" class="project-sub">{escape(value)}</text>')
            if i < 4:
                out.append(f'<text x="{x+276}" y="{y+65}" text-anchor="middle" class="number" style="font-size:24px;fill:#d53d3d">→</text>')
        out.append('</g>')
    out.append(f'<g filter="url(#shadow)"><rect x="225" y="770" width="1150" height="75" rx="22" fill="#fff" stroke="#d8a03e" stroke-width="3"/><text x="800" y="817" text-anchor="middle" class="project">{escape(data["noise"])}</text></g></svg>')
    return "\n".join(out)


def foundation_map(lang: str, row: tuple) -> str:
    number, _, map_name, titles, summaries, _, signal, _, _, _, _ = row
    title, subtitle = titles[lang], summaries[lang]
    stages = FOUNDATION_STAGES[number][lang]
    width = 1380 / len(stages)
    out = [svg_open(title, subtitle)]
    out.append(f'<text x="800" y="64" text-anchor="middle" class="title" style="font-size:36px">{escape(title.upper())}</text>')
    out.append(f'<text x="800" y="102" text-anchor="middle" class="subtitle" style="font-size:19px">{escape(subtitle)}</text>')
    colors=("blue","cyan","violet","orange","blue","cyan")
    for i, stage in enumerate(stages):
        x=60+i*width
        font_size = 17 if len(stage) > 18 else 20
        out.append(f'<g class="foundation-stage" data-order="{i+1}" filter="url(#shadow)"><rect x="{x:.0f}" y="185" width="{width-32:.0f}" height="130" rx="24" fill="url(#{colors[i]})" stroke="#fff" stroke-width="3"/><text x="{x+(width-32)/2:.0f}" y="260" text-anchor="middle" class="step-title" style="font-size:{font_size}px">{escape(stage)}</text></g>')
        if i<len(stages)-1:
            out.append(f'<path d="M{x+width-29:.0f} 250 H{x+width-4:.0f}" stroke="#d53d3d" stroke-width="7" marker-end="url(#arrow)"/>')
    labels={
      "ru":("НАБЛЮДАЙ","ОБЪЯСНИ","ПРОВЕРЬ","ПЕРЕНЕСИ","ДОКАЖИ ПРОЕКТОМ"),
      "kz":("БАҚЫЛА","ТҮСІНДІР","ТЕКСЕР","КӨШІР","ЖОБАМЕН ДӘЛЕЛДЕ"),
      "en":("OBSERVE","EXPLAIN","VERIFY","TRANSFER","PROVE IN THE PROJECT"),
    }[lang]
    for i,label in enumerate(labels):
        x=90+i*300
        out.append(f'<g class="learning-check" filter="url(#shadow)"><rect x="{x}" y="435" width="250" height="110" rx="22" fill="#fff" stroke="#8db7d8" stroke-width="3"/><text x="{x+125}" y="502" text-anchor="middle" class="project-sub" style="font-weight:800">{escape(label)}</text></g>')
    gate={"ru":"урок → изменение версии 0.1 → наблюдаемое доказательство","kz":"сабақ → 0.1 нұсқасының өзгерісі → бақыланатын дәлел","en":"lesson → version 0.1 change → observable evidence"}[lang]
    if number==10:
        gate={"ru":"8/10 сейчас • перенос обязателен • 7/10 через семь дней","kz":"қазір 8/10 • көшіру міндетті • жеті күннен кейін 7/10","en":"8/10 now • transfer required • 7/10 after seven days"}[lang]
    out.append(f'<g filter="url(#shadow)"><rect x="210" y="650" width="1180" height="120" rx="28" fill="#fff" stroke="#d8a03e" stroke-width="4"/><text x="800" y="722" text-anchor="middle" class="project">{escape(gate)}</text></g></svg>')
    return "\n".join(out)


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    for lang, data in LANG.items():
        (OUT / f"map-01-whole-map-{lang}.svg").write_text(whole_map(lang, data), encoding="utf-8")
        (OUT / f"map-02-diagnostic-{lang}.svg").write_text(diagnostic_map(lang, data), encoding="utf-8")
        lesson = LESSON_TEXT[lang]
        (OUT / f"map-03-device-system-{lang}.svg").write_text(system_map(lang, lesson), encoding="utf-8")
        (OUT / f"map-04-input-output-{lang}.svg").write_text(io_map(lang, lesson), encoding="utf-8")
        for row in FOUNDATION_ROWS:
            (OUT / f"{row[2]}-{lang}.svg").write_text(foundation_map(lang, row), encoding="utf-8")
    print("Generated 30 localized Informatics support maps")


if __name__ == "__main__":
    main()
