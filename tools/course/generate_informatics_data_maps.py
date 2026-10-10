#!/usr/bin/env python3
"""Generate exact, trilingual 1600×900 support diagrams for lessons 47–54."""

import json
from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/informatics"
ROWS = {row["number"]: row for row in json.loads(
    (ROOT / "course/informatics/curriculum.json").read_text(encoding="utf-8"))["lessons"]}

# Each tuple is (stage label, concrete input, verified output). The same
# fictional IDs and values appear in the CSV, database, examples and reports.
DATA = {
47: {
"ru": [("Событие", "t-01 · 7 октября", "20 минут"), ("Строка", "s-01 → t-01", "2026-10-07 · paper"), ("Проверка", "4 строки ≠ 4 человека", "единица: минуты")],
"kz": [("Оқиға", "t-01 · 7 қазан", "20 минут"), ("Жол", "s-01 → t-01", "2026-10-07 · paper"), ("Тексеру", "4 жол ≠ 4 адам", "өлшем: минут")],
"en": [("Event", "t-01 · 7 October", "20 minutes"), ("Row", "s-01 → t-01", "2026-10-07 · paper"), ("Check", "4 rows ≠ 4 people", "unit: minutes")]},
48: {
"ru": [("Четыре строки", "20 + 25 + 15 + 10", "минуты, не часы"), ("Диапазон", "D2:D5", "обе границы включены"), ("Сумма", "70 минут", "45 + 25 = 70")],
"kz": [("Төрт жол", "20 + 25 + 15 + 10", "минут, сағат емес"), ("Ауқым", "D2:D5", "екі шеті де қамтылады"), ("Қосынды", "70 минут", "45 + 25 = 70")],
"en": [("Four rows", "20 + 25 + 15 + 10", "minutes, not hours"), ("Range", "D2:D5", "both endpoints included"), ("Total", "70 minutes", "45 + 25 = 70")]},
49: {
"ru": [("Сырой файл", "4 исходные строки", "оставить без изменений"), ("Проверка", "повтор ID · плохая дата", "ошибки в карантин"), ("Импорт", "только 4 годные", "ошибка → нет базы")],
"kz": [("Бастапқы файл", "4 бастапқы жол", "өзгеріссіз сақтау"), ("Тексеру", "ID қайталау · қате күн", "қате карантинге"), ("Импорт", "4 дұрыс жол ғана", "қате → база жоқ")],
"en": [("Raw file", "4 original rows", "preserve unchanged"), ("Validation", "duplicate ID · bad date", "quarantine errors"), ("Import", "4 valid rows only", "error → no database")]},
50: {
"ru": [("Источник", "4 вымышленные сессии", "7–9 октября"), ("Шкала", "# = 5 минут", "ось начинается с 0"), ("Диаграмма", "45 · 0 · 25", "нет записи ≠ 0 минут")],
"kz": [("Дереккөз", "4 ойдан алынған сабақ", "7–9 қазан"), ("Шкала", "# = 5 минут", "ось 0-ден басталады"), ("Диаграмма", "45 · 0 · 25", "жазба жоқ ≠ 0 минут")],
"en": [("Source", "4 fictional sessions", "7–9 October"), ("Scale", "# = 5 minutes", "axis starts at 0"), ("Chart", "45 · 0 · 25", "no row ≠ zero time")]},
51: {
"ru": [("Задачи", "t-01 · t-02 · t-03", "task_id — ключ"), ("Занятия", "s-01 … s-04", "session_id — ключ"), ("Связь", "s-01 → t-01", "t-99: отказ")],
"kz": [("Тапсырмалар", "t-01 · t-02 · t-03", "task_id — кілт"), ("Сабақтар", "s-01 … s-04", "session_id — кілт"), ("Байланыс", "s-01 → t-01", "t-99: бас тарту")],
"en": [("Tasks", "t-01 · t-02 · t-03", "task_id is key"), ("Sessions", "s-01 … s-04", "session_id is key"), ("Relation", "s-01 → t-01", "t-99: reject")]},
52: {
"ru": [("Условие", "task_id = ?", "параметр: t-01"), ("Результат", "20 → 25", "ORDER BY date"), ("Сверка", "COUNT=2 · SUM=45", "t-02: 0 строк")],
"kz": [("Шарт", "task_id = ?", "параметр: t-01"), ("Нәтиже", "20 → 25", "ORDER BY date"), ("Салыстыру", "COUNT=2 · SUM=45", "t-02: 0 жол")],
"en": [("Condition", "task_id = ?", "parameter: t-01"), ("Rows", "20 → 25", "ORDER BY date"), ("Verify", "COUNT=2 · SUM=45", "t-02: 0 rows")]},
53: {
"ru": [("Левая таблица", "3 задачи", "включая t-02"), ("LEFT JOIN", "4 занятия + пустая связь", "сохранить t-02"), ("Отчёт", "45 · 0 · 25", "три строки, сумма 70")],
"kz": [("Сол кесте", "3 тапсырма", "t-02 ішінде"), ("LEFT JOIN", "4 сабақ + бос байланыс", "t-02 сақталады"), ("Есеп", "45 · 0 · 25", "үш жол, қосынды 70")],
"en": [("Left table", "3 tasks", "including t-02"), ("LEFT JOIN", "4 sessions + empty match", "keep t-02"), ("Report", "45 · 0 · 25", "3 rows, total 70")]},
54: {
"ru": [("Исходник", "JSON: 3 задачи", "CSV: 4 занятия"), ("SQLite", "новый локальный файл", "ошибка → нет базы"), ("Проверка", "REMIND · DONE · NO_DATE", "отчёт: 45 · 0 · 25")],
"kz": [("Бастапқы", "JSON: 3 тапсырма", "CSV: 4 сабақ"), ("SQLite", "жаңа жергілікті файл", "қате → база жоқ"), ("Тексеру", "REMIND · DONE · NO_DATE", "есеп: 45 · 0 · 25")],
"en": [("Sources", "JSON: 3 tasks", "CSV: 4 sessions"), ("SQLite", "new local file", "error → no database"), ("Check", "REMIND · DONE · NO_DATE", "report: 45 · 0 · 25")]},
}

FOOT = {
    "ru": "Предскажи → запусти → измени → объясни → проверь",
    "kz": "Болжа → іске қос → өзгерт → түсіндір → тексер",
    "en": "Predict → run → change → explain → verify",
}


def wrap_title(title):
    lines = [""]
    for word in title.split():
        candidate = (lines[-1] + " " + word).strip()
        if len(candidate) > 43 and lines[-1]:
            lines.append(word)
        else:
            lines[-1] = candidate
    if len(lines) > 2:
        raise ValueError(f"title exceeds two lines: {title}")
    return lines


def svg(number, lang):
    title = ROWS[number][f"title_{lang}"]
    heading = "".join(f'<text x="64" y="{147 + i*55}" class="title">{escape(line)}</text>'
                      for i, line in enumerate(wrap_title(title)))
    cards = []
    for i, (label, value, explanation) in enumerate(DATA[number][lang]):
        for string in (label, value, explanation):
            if len(string) > 31:
                raise ValueError(f"map {number}/{lang} text too wide: {string}")
        x = 54 + i*514
        colour = ("#66e6df", "#ffd27d", "#ff9584")[i]
        value_size = ' style="font-size:21px"' if len(value) > 22 else ""
        value_markup = f'<text x="241" y="192" text-anchor="middle" class="value"{value_size}>{escape(value)}</text>'
        explanation_y = 291
        if number == 50 and i == 2:
            # Actual zero-based bars: 45, no recorded session, 25 minutes.
            value_markup = '''<path d="M87 209 H395" stroke="#c4e6ef" stroke-width="2"/>
<rect x="123" y="145" width="47" height="64" rx="5" fill="#66e6df"/>
<rect x="227" y="207" width="47" height="2" fill="#ffd27d"/>
<rect x="331" y="173" width="47" height="36" rx="5" fill="#ff9584"/>
<text x="146" y="137" text-anchor="middle" class="bar">45</text>
<text x="250" y="190" text-anchor="middle" class="bar">0</text>
<text x="354" y="165" text-anchor="middle" class="bar">25</text>
<text x="146" y="249" text-anchor="middle" class="bar">t-01</text>
<text x="250" y="249" text-anchor="middle" class="bar">t-02</text>
<text x="354" y="249" text-anchor="middle" class="bar">t-03</text>'''
            explanation_y = 316
        if number == 51 and i == 2:
            value_markup = '''<text x="241" y="177" text-anchor="middle" class="relation">s-01, s-02 → t-01</text>
<text x="241" y="210" text-anchor="middle" class="relation">s-03, s-04 → t-03</text>'''
        if number == 53 and i == 2:
            value_markup = '''<text x="241" y="169" text-anchor="middle" class="relation">t-01 → 45</text>
<text x="241" y="191" text-anchor="middle" class="relation">t-02 → 0</text>
<text x="241" y="213" text-anchor="middle" class="relation">t-03 → 25</text>'''
        cards.append(f'''<g class="data-stage" transform="translate({x} 307)">
<rect x="9" y="13" width="482" height="354" rx="22" fill="#020d17" opacity=".67"/>
<rect width="482" height="354" rx="22" fill="url(#panel)" stroke="{colour}" stroke-width="3"/>
<path d="M22 22 H460" stroke="#b7dcea" opacity=".4"/>
<rect x="25" y="24" width="59" height="57" rx="14" fill="{colour}"/>
<text x="54" y="65" text-anchor="middle" class="number">{i+1}</text>
<text x="99" y="65" class="label">{escape(label)}</text>
<path d="M27 107 H454" stroke="#6e91a9" stroke-width="2"/>
<rect x="27" y="136" width="428" height="89" rx="12" fill="#061b2b" stroke="#4b758b"/>
{value_markup}
<text x="241" y="{explanation_y}" text-anchor="middle" class="explain">{escape(explanation)}</text>
</g>''')
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}">
<defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#071827"/><stop offset="1" stop-color="#164c60"/></linearGradient><linearGradient id="panel" x2="1" y2="1"><stop stop-color="#286076"/><stop offset="1" stop-color="#102b3f"/></linearGradient><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#c7eaf2" stroke-opacity=".12"/></pattern></defs>
<style>.eyebrow{{font:700 22px Arial,sans-serif;fill:#86eee7;letter-spacing:3px}}.title{{font:700 43px Arial,sans-serif;fill:#f5f9ff}}.number{{font:700 30px Arial,sans-serif;fill:#102333}}.label{{font:700 28px Arial,sans-serif;fill:#fff}}.value{{font:600 27px 'Courier New',monospace;fill:#fff}}.relation{{font:700 21px 'Courier New',monospace;fill:#fff}}.bar{{font:700 20px Arial,sans-serif;fill:#fff}}.explain{{font:600 25px Arial,sans-serif;fill:#e0f4f8}}.foot{{font:600 27px Arial,sans-serif;fill:#e0f0f7}}</style>
<rect width="1600" height="900" fill="url(#bg)"/><rect width="1600" height="900" fill="url(#grid)"/>
<text x="64" y="75" class="eyebrow">DATA · {number:02d} / 72</text>{heading}
<path d="M64 253 H1536" stroke="#6b94ab" stroke-width="2"/>
{''.join(cards)}
<path d="M544 482 h18 m-9 -9 9 9 -9 9 M1058 482 h18 m-9 -9 9 9 -9 9" fill="none" stroke="#e2f6f9" stroke-width="5" stroke-linejoin="round"/>
<path d="M64 755 H1536" stroke="#6b94ab" stroke-width="2"/>
<text x="800" y="812" text-anchor="middle" class="foot">{escape(FOOT[lang])}</text>
</svg>'''


def main():
    assert set(DATA) == set(range(47, 55))
    for number in range(47, 55):
        for lang in ("ru", "kz", "en"):
            path = OUT / f"map-{number:02d}-{ROWS[number]['id']}-{lang}.svg"
            path.write_text(svg(number, lang), encoding="utf-8")
    print("Wrote 24 exact trilingual data-block support maps")


if __name__ == "__main__":
    main()
