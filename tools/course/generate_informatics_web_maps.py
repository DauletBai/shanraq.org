#!/usr/bin/env python3
"""Produce exact 1600x900 support maps for the trilingual web block."""

import json
from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/informatics"
ROWS = {row["number"]: row for row in json.loads((ROOT / "course/informatics/curriculum.json").read_text(encoding="utf-8"))["lessons"]}

# Every tuple is (short heading, observed input, result). All labels are kept
# deliberately short, so none must squeeze into a card or collide with an edge.
DATA = {
39: {
"ru": [("Устройство", "телефон → роутер", "запрос начал клиент"), ("Промежуток", "сеть провайдера", "путь может меняться"), ("Ответ", "сервер → браузер", "локально: 127.0.0.1")],
"kz": [("Құрылғы", "телефон → роутер", "сұрауды клиент бастады"), ("Аралық", "провайдер желісі", "жол өзгере алады"), ("Жауап", "сервер → браузер", "жергілікті: 127.0.0.1")],
"en": [("Device", "phone → router", "client starts request"), ("Middle", "provider network", "route may change"), ("Reply", "server → browser", "local: 127.0.0.1")]},
40: {
"ru": [("Адрес", "203.0.113.7", "учебный IP"), ("Порядок", "2 → 1 → 3", "пакеты пришли не по порядку"), ("Сборка", "1 → 2 → 3", "TASK")],
"kz": [("Мекенжай", "203.0.113.7", "оқу IP-і"), ("Реті", "2 → 1 → 3", "бөліктер араласып келді"), ("Жинау", "1 → 2 → 3", "TASK")],
"en": [("Address", "203.0.113.7", "documentation IP"), ("Arrival", "2 → 1 → 3", "pieces arrive mixed"), ("Assemble", "1 → 2 → 3", "TASK")]},
41: {
"ru": [("Имя", "assistant.test", "человек помнит имя"), ("DNS", "запись A", "203.0.113.7"), ("Предел", "адрес найден", "доверие не доказано")],
"kz": [("Атау", "assistant.test", "адам атауды біледі"), ("DNS", "A жазбасы", "203.0.113.7"), ("Шегі", "мекенжай табылды", "сенім дәлелденбеді")],
"en": [("Name", "assistant.test", "a readable name"), ("DNS", "A record", "203.0.113.7"), ("Limit", "address found", "trust not proved")]},
42: {
"ru": [("Потеря", "1, ?, 3", "часть 2 отсутствует"), ("TCP", "ждём часть 2", "повторная отправка"), ("Итог", "1, 2, 3", "целый поток")],
"kz": [("Жоғалту", "1, ?, 3", "2-бөлік жоқ"), ("TCP", "2-бөлікті күтеміз", "қайта жіберіледі"), ("Нәтиже", "1, 2, 3", "тұтас ағын")],
"en": [("Missing", "1, ?, 3", "part 2 absent"), ("TCP", "wait for part 2", "retransmit"), ("Outcome", "1, 2, 3", "ordered stream")]},
43: {
"ru": [("Имя", "assistant.test", "ожидаемый адресат"), ("Сертификат", "имя совпало?", "иначе остановиться"), ("TLS", "защищённый канал", "не гарантия честности")],
"kz": [("Атау", "assistant.test", "күтілген тарап"), ("Сертификат", "атау сәйкес пе?", "әйтпесе тоқтау"), ("TLS", "қорғалған арна", "адалдық кепілі емес")],
"en": [("Name", "assistant.test", "expected peer"), ("Certificate", "name matches?", "otherwise stop"), ("TLS", "protected channel", "not honesty proof")]},
44: {
"ru": [("Запрос", "GET /?today=...", "чтение без записи"), ("Проверка", "дата существует?", "2026-02-29: нет"), ("Ответ", "200 или 400", "нет пути: 404")],
"kz": [("Сұрау", "GET /?today=...", "жазусыз оқу"), ("Тексеру", "күн бар ма?", "2026-02-29: жоқ"), ("Жауап", "200 не 400", "жол жоқ: 404")],
"en": [("Request", "GET /?today=...", "read without write"), ("Validate", "date exists?", "2026-02-29: no"), ("Response", "200 or 400", "no route: 404")]},
45: {
"ru": [("Структура", "label + input", "поле названо"), ("Защита", "escape(<b>)", "текст, не команда"), ("Проверка", "Tab → кнопка", "смысл без цвета")],
"kz": [("Құрылым", "label + input", "өріс аталды"), ("Қорғау", "escape(<b>)", "мәтін, пәрмен емес"), ("Тексеру", "Tab → батырма", "мән түске тәуелсіз")],
"en": [("Structure", "label + input", "named field"), ("Output", "escape(<b>)", "text, not markup"), ("Check", "Tab → button", "meaning without colour")]},
46: {
"ru": [("Запуск", "127.0.0.1:8765", "только свой компьютер"), ("Экран", "REMIND / DONE", "NO_DATE сохранён"), ("Проверка", "тесты + файл", "данные не изменены")],
"kz": [("Іске қосу", "127.0.0.1:8765", "тек өз компьютері"), ("Экран", "REMIND / DONE", "NO_DATE сақталды"), ("Тексеру", "тесттер + файл", "дерек өзгермеді")],
"en": [("Run", "127.0.0.1:8765", "same computer only"), ("View", "REMIND / DONE", "NO_DATE preserved"), ("Verify", "tests + file", "data unchanged")]},
}

FOOT = {
"ru": "Предскажи → запусти → измени → объясни → проверь",
"kz": "Болжа → іске қос → өзгерт → түсіндір → тексер",
"en": "Predict → run → change → explain → verify",
}


def title_lines(value):
    words = value.split()
    lines = [""]
    for word in words:
        candidate = (lines[-1] + " " + word).strip()
        if len(candidate) > 43 and lines[-1]:
            lines.append(word)
        else:
            lines[-1] = candidate
    assert len(lines) <= 2, (value, lines)
    return lines


def svg(number, lang):
    title = ROWS[number][f"title_{lang}"]
    lines = title_lines(title)
    heading = ''.join(f'<text x="64" y="{157+i*57}" class="title">{escape(part)}</text>' for i, part in enumerate(lines))
    cards = []
    for i, (label, input_, result) in enumerate(DATA[number][lang]):
        x = 58 + i * 512
        colour = ("#54dbe0", "#ffd47a", "#fa8b75")[i]
        for string in (label, input_, result):
            assert len(string) <= 31, (number, lang, string)
        cards.append(f'''<g class="web-stage" transform="translate({x} 310)">
<rect x="10" y="13" width="478" height="348" rx="22" fill="#071522" opacity=".7"/>
<rect width="478" height="348" rx="22" fill="url(#panel)" stroke="{colour}" stroke-width="3"/>
<rect x="26" y="27" width="56" height="56" rx="14" fill="{colour}"/>
<text x="54" y="67" text-anchor="middle" class="number">{i+1}</text>
<text x="100" y="67" class="label">{escape(label)}</text>
<path d="M30 105 H448" stroke="#6e91a9" stroke-width="2"/>
<rect x="28" y="134" width="422" height="88" rx="13" fill="#081a2b" stroke="#456d89"/>
<text x="239" y="190" text-anchor="middle" class="input">{escape(input_)}</text>
<text x="239" y="290" text-anchor="middle" class="result">{escape(result)}</text>
</g>''')
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}">
<defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#061525"/><stop offset="1" stop-color="#15405c"/></linearGradient><linearGradient id="panel" x2="1" y2="1"><stop stop-color="#23475f"/><stop offset="1" stop-color="#10283f"/></linearGradient><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#bedaf0" stroke-opacity=".11"/></pattern></defs>
<style>.eyebrow{{font:700 22px Arial,sans-serif;fill:#7be7e8;letter-spacing:3px}}.title{{font:700 43px Arial,sans-serif;fill:#f5f9ff}}.number{{font:700 30px Arial,sans-serif;fill:#102333}}.label{{font:700 29px Arial,sans-serif;fill:#fff}}.input{{font:600 26px 'Courier New',monospace;fill:#fff}}.result{{font:600 26px Arial,sans-serif;fill:#d8f2fa}}.foot{{font:600 27px Arial,sans-serif;fill:#e0f0f7}}</style>
<rect width="1600" height="900" fill="url(#bg)"/><rect width="1600" height="900" fill="url(#grid)"/>
<text x="64" y="77" class="eyebrow">WEB · {number:02d} / 72</text>{heading}
<path d="M64 255 H1536" stroke="#6b94ab" stroke-width="2"/>
{''.join(cards)}
<path d="M541 484 h20 m-10 -10 10 10 -10 10 M1053 484 h20 m-10 -10 10 10 -10 10" fill="none" stroke="#d9edf9" stroke-width="5" stroke-linejoin="round"/>
<path d="M64 755 H1536" stroke="#6b94ab" stroke-width="2"/>
<text x="800" y="810" text-anchor="middle" class="foot">{escape(FOOT[lang])}</text>
</svg>'''


def main():
    assert set(DATA) == set(range(39, 47))
    for number in range(39, 47):
        for lang in ("ru", "kz", "en"):
            path = OUT / f"map-{number:02d}-{ROWS[number]['id']}-{lang}.svg"
            path.write_text(svg(number, lang), encoding="utf-8")
    print("Wrote 24 exact trilingual web-block support maps")


if __name__ == "__main__":
    main()
