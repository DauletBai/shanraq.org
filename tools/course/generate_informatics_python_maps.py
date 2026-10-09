#!/usr/bin/env python3
"""Render the Python block's concrete, bilingual-safe 1600×900 support maps."""

import json
from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/informatics"
ROWS = {x["number"]: x for x in json.loads((ROOT / "course/informatics/curriculum.json").read_text(encoding="utf-8"))["lessons"]}

# Each triple is an exact, short observable transition, rather than a decorative slogan.
DATA = {
  27: {"ru": [("До", "done = False", "экран: False"), ("Команда", "done = True", "имя получило True"), ("После", "print(done)", "экран: True")], "kz": [("Бұрын", "done = False", "экран: False"), ("Команда", "done = True", "атау True алды"), ("Кейін", "print(done)", "экран: True")], "en": [("Before", "done = False", "screen: False"), ("Command", "done = True", "name now holds True"), ("After", "print(done)", "screen: True")]},
  28: {"ru": [("Ввод", 'raw = "2"', "тип: str"), ("Преобразование", "days = int(raw)", "тип: int"), ("Операция", "days + 1", "результат: 3")], "kz": [("Енгізу", 'raw = "2"', "түрі: str"), ("Түрлендіру", "days = int(raw)", "түрі: int"), ("Әрекет", "days + 1", "нәтиже: 3")], "en": [("Input", 'raw = "2"', "type: str"), ("Convert", "days = int(raw)", "type: int"), ("Operate", "days + 1", "result: 3")]},
  29: {"ru": [("Число", "2 + 1", "результат: 3"), ("Текст", '"2" + "1"', 'результат: "21"'), ("Правило", "int(raw) + 1", "тип задан явно")], "kz": [("Сан", "2 + 1", "нәтиже: 3"), ("Мәтін", '"2" + "1"', 'нәтиже: "21"'), ("Ереже", "int(raw) + 1", "түрі ашық берілді")], "en": [("Number", "2 + 1", "result: 3"), ("Text", '"2" + "1"', 'result: "21"'), ("Rule", "int(raw) + 1", "type is explicit")]},
  30: {"ru": [("Просрочено", "days = -1", "OVERDUE"), ("Граница", "days = 2", "REMIND"), ("Позже", "days = 3", "NOT_YET")], "kz": [("Мерзімі өтті", "days = -1", "OVERDUE"), ("Шекара", "days = 2", "REMIND"), ("Кейін", "days = 3", "NOT_YET")], "en": [("Overdue", "days = -1", "OVERDUE"), ("Boundary", "days = 2", "REMIND"), ("Later", "days = 3", "NOT_YET")]},
  31: {"ru": [("Начало", "count = 0", "3 задачи"), ("Обход", "F → T → T", "+0 → +1 → +1"), ("Конец", "count = 2", "цикл завершён")], "kz": [("Басы", "count = 0", "3 тапсырма"), ("Қарау", "F → T → T", "+0 → +1 → +1"), ("Соңы", "count = 2", "цикл бітті")], "en": [("Start", "count = 0", "3 tasks"), ("Visit", "F → T → T", "+0 → +1 → +1"), ("Finish", "count = 2", "loop stops")]},
  32: {"ru": [("Список", "tasks = […, …]", "две записи"), ("Словарь", 'task["id"]', "поле записи"), ("Проверка", "t-01 ≠ t-02", "id различаются")], "kz": [("Тізім", "tasks = […, …]", "екі жазба"), ("Сөздік", 'task["id"]', "жазба өрісі"), ("Тексеру", "t-01 ≠ t-02", "id бөлек")], "en": [("List", "tasks = […, …]", "two records"), ("Dictionary", 'task["id"]', "record field"), ("Check", "t-01 ≠ t-02", "IDs differ")]},
  33: {"ru": [("Текст UTF-8", "Әліппе оқу", "Ә сохраняется"), ("Даты", "09.10 → 11.10", "2 перехода"), ("Итог", "days = 2", "не 48 часов")], "kz": [("UTF-8 мәтіні", "Әліппе оқу", "Ә сақталады"), ("Күндер", "09.10 → 11.10", "2 өту"), ("Нәтиже", "days = 2", "48 сағат емес")], "en": [("UTF-8 text", "Әліппе оқу", "Ә survives"), ("Dates", "09.10 → 11.10", "2 transitions"), ("Result", "days = 2", "not 48 hours")]},
  34: {"ru": [("Вход", "done = False", "due: 11.10"), ("Функция", "reminder_status", "today: 09.10"), ("Выход", "REMIND", "вход не меняется")], "kz": [("Кіріс", "done = False", "мерзім: 11.10"), ("Функция", "reminder_status", "бүгін: 09.10"), ("Шығыс", "REMIND", "кіріс өзгермейді")], "en": [("Input", "done = False", "due: 11 Oct"), ("Function", "reminder_status", "today: 09 Oct"), ("Output", "REMIND", "input unchanged")]},
  35: {"ru": [("Версия 0.2", "id, title, done", "срока нет"), ("Переход", "due_date = null", "старое сохранено"), ("Версия 1.0", "JSON UTF-8", "читать повторно")], "kz": [("0.2 нұсқасы", "id, title, done", "мерзім жоқ"), ("Көшу", "due_date = null", "ескісі сақталды"), ("1.0 нұсқасы", "JSON UTF-8", "қайта оқу")], "en": [("Version 0.2", "id, title, done", "no due date"), ("Migrate", "due_date = null", "old data survives"), ("Version 1.0", "JSON UTF-8", "read back") ]},
  36: {"ru": [("Ввод", "2026-02-29", "день не существует"), ("Проверка", "ValueError", "запись запрещена"), ("Состояние", "файл прежний", "сообщить ошибку")], "kz": [("Енгізу", "2026-02-29", "ондай күн жоқ"), ("Тексеру", "ValueError", "жазуға болмайды"), ("Күй", "файл өзгермеді", "қатені хабарлау")], "en": [("Input", "2026-02-29", "date does not exist"), ("Validate", "ValueError", "do not save"), ("State", "file unchanged", "report the error")]},
  37: {"ru": [("Договор", "days = 2", "ожидаем REMIND"), ("Тест", "assert actual", "сравнить с REMIND"), ("Регрессия", "days < 2", "тест падает")], "kz": [("Келісім", "days = 2", "күтіледі REMIND"), ("Тест", "assert actual", "REMIND-пен салыстыру"), ("Регрессия", "days < 2", "тест құлайды")], "en": [("Contract", "days = 2", "expect REMIND"), ("Test", "assert actual", "compare to REMIND"), ("Regression", "days < 2", "test fails")]},
  38: {"ru": [("Команда", "add → reminders", "срок: 11.10"), ("Сохранить", "done t-04", "повторно открыть"), ("Проверить", "DONE", "тесты проходят")], "kz": [("Команда", "add → reminders", "мерзім: 11.10"), ("Сақтау", "done t-04", "қайта ашу"), ("Тексеру", "DONE", "тесттер өтеді")], "en": [("Command", "add → reminders", "due: 11 Oct"), ("Persist", "done t-04", "reopen the file"), ("Verify", "DONE", "tests pass")]},
}


def svg(number, lang):
    row = ROWS[number]
    title = escape(row[f"title_{lang}"])
    cards = []
    colors = ("#00c0ca", "#ffd36e", "#fa7268")
    for index, (heading, code, outcome) in enumerate(DATA[number][lang]):
        x = 58 + index * 512
        color = colors[index]
        cards.append(f'''<g class="python-stage" transform="translate({x} 300)">
          <rect x="7" y="11" width="478" height="337" rx="24" fill="#061324" opacity=".54"/>
          <rect width="478" height="337" rx="24" fill="url(#card)" stroke="{color}" stroke-width="3"/>
          <rect x="28" y="28" width="54" height="54" rx="15" fill="{color}"/>
          <text x="55" y="66" text-anchor="middle" class="number">{index+1}</text>
          <text x="99" y="70" class="heading">{escape(heading)}</text>
          <path d="M28 106 H450" stroke="#314861" stroke-width="2"/>
          <rect x="28" y="132" width="422" height="91" rx="13" fill="#071727"/>
          <text x="239" y="190" text-anchor="middle" class="code">{escape(code)}</text>
          <text x="239" y="283" text-anchor="middle" class="result">{escape(outcome)}</text>
        </g>''')
    arrows = '<path d="M542 469 h20 m-10 -10 10 10 -10 10 M1054 469 h20 m-10 -10 10 10 -10 10" fill="none" stroke="#d3e6f3" stroke-width="5" stroke-linecap="round" stroke-linejoin="round"/>'
    footer = {"ru": "Предскажи → запусти → измени → объясни → проверь", "kz": "Болжа → іске қос → өзгерт → түсіндір → тексер", "en": "Predict → run → change → explain → verify"}[lang]
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{title}">
<defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#07182a"/><stop offset="1" stop-color="#153954"/></linearGradient><linearGradient id="card" x2="1" y2="1"><stop stop-color="#203953"/><stop offset="1" stop-color="#10263c"/></linearGradient><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#a9d2eb" stroke-opacity=".09"/></pattern></defs>
<style>.title{{font:700 46px Arial,sans-serif;fill:#f1f8ff}}.eyebrow{{font:700 22px Arial,sans-serif;fill:#73dce2;letter-spacing:4px}}.number{{font:700 32px Arial,sans-serif;fill:#061b2d}}.heading{{font:700 30px Arial,sans-serif;fill:#f4fbff}}.code{{font:600 29px 'Courier New',monospace;fill:#ffffff}}.result{{font:600 27px Arial,sans-serif;fill:#d7effc}}.foot{{font:600 27px Arial,sans-serif;fill:#c1d9e8}}</style>
<rect width="1600" height="900" fill="url(#bg)"/><rect width="1600" height="900" fill="url(#grid)"/>
<text x="64" y="100" class="eyebrow">PYTHON · {number:02d} / 38</text>
<text x="64" y="183" class="title">{title}</text>
<path d="M64 220 H1536" stroke="#648aa3" stroke-width="2"/>
{''.join(cards)}{arrows}
<path d="M64 726 H1536" stroke="#648aa3" stroke-width="2"/>
<text x="800" y="799" text-anchor="middle" class="foot">{escape(footer)}</text>
</svg>'''


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    for number in range(27, 39):
        for lang in ("ru", "kz", "en"):
            path = OUT / f"map-{number:02d}-{ROWS[number]['id']}-{lang}.svg"
            path.write_text(svg(number, lang), encoding="utf-8")
    print("Wrote 36 exact Python-block support maps")


if __name__ == "__main__":
    main()
