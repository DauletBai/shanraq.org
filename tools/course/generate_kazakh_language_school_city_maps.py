#!/usr/bin/env python3
"""Generate precise localized support maps for Kazakh lessons 17–24."""
from __future__ import annotations

from html import escape
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/kazakh-language"


TEXT = {
    "17-school-map": {
        "ru": (["Школа как карта"], "место отвечает на Қайда?", ["1 этаж", "асхана · кіреберіс"], ["2 этаж", "сынып · кітапхана"], ["точка", "спортзал оң жақта"]),
        "kz": (["Мектеп карта ретінде"], "орын Қайда? сұрағына жауап береді", ["1-қабат", "асхана · кіреберіс"], ["2-қабат", "сынып · кітапхана"], ["координат", "спортзал оң жақта"]),
        "en": (["School as a map"], "a place answers Қайда?", ["floor 1", "асхана · кіреберіс"], ["floor 2", "сынып · кітапхана"], ["coordinate", "спортзал оң жақта"]),
    },
    "18-days-subjects": {
        "ru": (["Неделя и предметы"], "день становится координатой времени", ["дүйсенбіде", "математика"], ["сәрсенбіде", "дене шынықтыру"], ["бейсенбіде", "информатика"]),
        "kz": (["Апта және пәндер"], "күн уақыт координатына айналады", ["дүйсенбіде", "математика"], ["сәрсенбіде", "дене шынықтыру"], ["бейсенбіде", "информатика"]),
        "en": (["Week and subjects"], "a day becomes a time coordinate", ["дүйсенбіде", "математика"], ["сәрсенбіде", "дене шынықтыру"], ["бейсенбіде", "информатика"]),
    },
    "19-clock-time": {
        "ru": (["Время без путаницы"], "слово, цифры и стрелки показывают одно", ["08:00", "сегіз"], ["08:30", "сегіз жарым"], ["10:15", "оннан 15 минут кетті"]),
        "kz": (["Уақытты шатастырмаймыз"], "сөз, цифр және тілдер бір уақытты көрсетеді", ["08:00", "сегіз"], ["08:30", "сегіз жарым"], ["10:15", "оннан 15 минут кетті"]),
        "en": (["Time without confusion"], "words, digits, and hands show the same time", ["08:00", "сегіз"], ["08:30", "сегіз жарым"], ["10:15", "оннан 15 минут кетті"]),
    },
    "20-directions": {
        "ru": (["Маршрут голосом"], "каждая стрелка точно совпадает с командой", ["1", "Тура жүріңіз"], ["2", "Оңға бұрылыңыз"], ["3", "Солға бұрылыңыз"]),
        "kz": (["Дауыспен бағыттау"], "әр көрсеткі өз нұсқауымен дәл сәйкес", ["1", "Тура жүріңіз"], ["2", "Оңға бұрылыңыз"], ["3", "Солға бұрылыңыз"]),
        "en": (["A spoken route"], "every arrow matches its instruction", ["1", "Тура жүріңіз"], ["2", "Оңға бұрылыңыз"], ["3", "Солға бұрылыңыз"]),
    },
    "21-from-to": {
        "ru": (["Весь путь"], "сначала роль, затем окончание", ["старт · Қайдан?", "мектептен"], ["ориентир", "саябақтан кейін"], ["цель · Қайда?", "кітапханаға дейін"]),
        "kz": (["Толық бағыт"], "алдымен рөл, кейін жалғау", ["бастау · Қайдан?", "мектептен"], ["бағдар", "саябақтан кейін"], ["мақсат · Қайда?", "кітапханаға дейін"]),
        "en": (["The whole route"], "choose the role before the ending", ["start · Қайдан?", "мектептен"], ["landmark", "саябақтан кейін"], ["goal · Қайда?", "кітапханаға дейін"]),
    },
    "22-transport": {
        "ru": (["Транспорт и цель"], "номер, направление и остановка проверяются отдельно", ["средство", "12 · автобуспен"], ["цель", "орталыққа"], ["выход", "келесі аялдама"]),
        "kz": (["Көлік және мақсат"], "нөмір, бағыт және аялдама бөлек тексеріледі", ["құрал", "12 · автобуспен"], ["мақсат", "орталыққа"], ["түсу", "келесі аялдама"]),
        "en": (["Transport and goal"], "check route number, direction, and stop separately", ["means", "12 · автобуспен"], ["goal", "орталыққа"], ["exit", "келесі аялдама"]),
    },
    "23-rules": {
        "ru": (["Три статуса правила"], "действие одно — результат различается", ["можно", "кіруге болады"], ["нельзя", "жүгіруге болмайды"], ["нужно", "өшіру керек"]),
        "kz": (["Ереженің үш мәртебесі"], "әрекет бір, нәтиже әртүрлі", ["рұқсат", "кіруге болады"], ["тыйым", "жүгіруге болмайды"], ["қажет", "өшіру керек"]),
        "en": (["Three rule states"], "the action stays clear; its status changes", ["allowed", "кіруге болады"], ["not allowed", "жүгіруге болмайды"], ["required", "өшіру керек"]),
    },
    "24-mastery": {
        "ru": (["Проект «Моя среда 0.3»"], "учебный день и маршрут без сценария", ["школа", "место · день · время"], ["город", "старт · путь · транспорт"], ["доказательство", "правило · вопрос · 8/10"]),
        "kz": (["«Менің ортам 0.3» жобасы"], "дайын мәтінсіз оқу күні мен бағыт", ["мектеп", "орын · күн · уақыт"], ["қала", "бастау · жол · көлік"], ["дәлел", "ереже · сұрақ · 8/10"]),
        "en": (["My World 0.3 project"], "a school day and route without a script", ["school", "place · day · time"], ["city", "start · route · transport"], ["evidence", "rule · question · 8/10"]),
    },
}


def lines(items, x, y, size=28, gap=42, weight=650, color="#17364e"):
    body = "".join(f'<tspan x="{x}" y="{y + i * gap}">{escape(item)}</tspan>' for i, item in enumerate(items))
    return f'<text x="{x}" y="{y}" text-anchor="middle" font-size="{size}" font-weight="{weight}" fill="{color}">{body}</text>'


def base(title, subtitle, body):
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(' '.join(title))}">
<defs>
  <linearGradient id="paper" x1="0" y1="0" x2="0" y2="1"><stop stop-color="#fff"/><stop offset="1" stop-color="#edf7fb"/></linearGradient>
  <linearGradient id="blue" x1="0" y1="0" x2="0" y2="1"><stop stop-color="#73c8ee"/><stop offset="1" stop-color="#2f8fbd"/></linearGradient>
  <linearGradient id="gold" x1="0" y1="0" x2="0" y2="1"><stop stop-color="#ffe497"/><stop offset="1" stop-color="#e2a63e"/></linearGradient>
  <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#84afc8" stroke-width="1" opacity=".25"/></pattern>
  <filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="13" stdDeviation="11" flood-color="#17364e" flood-opacity=".22"/></filter>
  <marker id="arrow" markerWidth="10" markerHeight="10" refX="8" refY="5" orient="auto"><path d="M0 0L10 5L0 10Z" fill="#c83236"/></marker>
</defs>
<rect width="1600" height="900" rx="34" fill="url(#paper)"/><rect width="1600" height="900" rx="34" fill="url(#grid)"/>
<g font-family="Inter,Arial,sans-serif">{lines(title, 800, 86, 50, 54, 830)}<text x="800" y="174" text-anchor="middle" font-size="25" font-weight="560" fill="#4d6d80">{escape(subtitle)}</text>{body}
<g transform="translate(356 792)" filter="url(#shadow)"><rect width="888" height="64" rx="19" fill="#17364e"/><text x="444" y="42" text-anchor="middle" font-size="25" font-weight="760" fill="#fff">КӨР → ТЫҢДА → АЙТ → ҚОЛДАН → ТЕКСЕР</text></g></g></svg>'''


def card(x, title, detail, color, index):
    return f'''<g transform="translate({x} 294)" filter="url(#shadow)"><rect width="410" height="356" rx="24" fill="#fff" stroke="{color}" stroke-width="3"/><rect x="12" y="12" width="386" height="82" rx="16" fill="{color}" opacity=".17"/><circle cx="205" cy="172" r="52" fill="{color}" opacity=".18" stroke="{color}" stroke-width="3"/><text x="205" y="185" text-anchor="middle" font-size="34" font-weight="850" fill="#17364e">{escape(index)}</text>{lines([title],205,246,27,40,780)}{lines([detail],205,304,23,38,620,"#486779")}</g>'''


def three_cards(spec):
    title, subtitle, a, b, c = spec
    body = card(80, a[0], a[1], "#e1a53a", "1") + card(595, b[0], b[1], "#3598c3", "2") + card(1110, c[0], c[1], "#49a468", "3")
    body += '<path d="M505 472H570" fill="none" stroke="#c83236" stroke-width="7" stroke-linecap="round" marker-end="url(#arrow)"/><path d="M1020 472H1085" fill="none" stroke="#c83236" stroke-width="7" stroke-linecap="round" marker-end="url(#arrow)"/>'
    return base(title, subtitle, body)


def clock(cx, cy, hour, minute, label, phrase):
    # Angles are measured clockwise from twelve. Exact hands prevent the diagram
    # from contradicting the time named below it.
    import math
    ma = math.radians(minute * 6 - 90)
    ha = math.radians((hour % 12) * 30 + minute * .5 - 90)
    mx, my = cx + 100 * math.cos(ma), cy + 100 * math.sin(ma)
    hx, hy = cx + 68 * math.cos(ha), cy + 68 * math.sin(ha)
    ticks = []
    for n in range(12):
        a = math.radians(n * 30 - 90)
        x1, y1 = cx + 116 * math.cos(a), cy + 116 * math.sin(a)
        x2, y2 = cx + 128 * math.cos(a), cy + 128 * math.sin(a)
        ticks.append(f'<path d="M{x1:.1f} {y1:.1f}L{x2:.1f} {y2:.1f}" stroke="#17364e" stroke-width="5" stroke-linecap="round"/>')
    return f'''<g filter="url(#shadow)"><circle cx="{cx}" cy="{cy}" r="148" fill="#fff" stroke="#3998c2" stroke-width="6"/>{''.join(ticks)}<path d="M{cx} {cy}L{hx:.1f} {hy:.1f}" stroke="#17364e" stroke-width="11" stroke-linecap="round"/><path d="M{cx} {cy}L{mx:.1f} {my:.1f}" stroke="#c83236" stroke-width="7" stroke-linecap="round"/><circle cx="{cx}" cy="{cy}" r="11" fill="#e1a53a"/></g><text x="{cx}" y="{cy+205}" text-anchor="middle" font-size="31" font-weight="820" fill="#17364e">{escape(label)}</text><text x="{cx}" y="{cy+248}" text-anchor="middle" font-size="23" font-weight="620" fill="#486779">{escape(phrase)}</text>'''


def render(stem, spec):
    if stem == "19-clock-time":
        title, subtitle, a, b, c = spec
        body = clock(310, 420, 8, 0, *a) + clock(800, 420, 8, 30, *b) + clock(1290, 420, 10, 15, *c)
        return base(title, subtitle, body)
    if stem == "20-directions":
        title, subtitle, a, b, c = spec
        body = '''<g filter="url(#shadow)"><rect x="150" y="270" width="1300" height="410" rx="30" fill="#fff" stroke="#63a7c8" stroke-width="3"/>
        <path d="M270 620V420H1050V300" fill="none" stroke="#c83236" stroke-width="15" stroke-linecap="round" stroke-linejoin="round" marker-end="url(#arrow)"/>
        <circle cx="270" cy="620" r="28" fill="#e1a53a"/><circle cx="1050" cy="300" r="28" fill="#49a468"/></g>'''
        for x, y, item in ((445,590,a),(650,365,b),(1240,350,c)):
            body += f'<g filter="url(#shadow)"><rect x="{x-160}" y="{y-44}" width="320" height="88" rx="18" fill="#fff" stroke="#3998c2" stroke-width="3"/><text x="{x}" y="{y-5}" text-anchor="middle" font-size="23" font-weight="780" fill="#17364e">{escape(item[0])}</text><text x="{x}" y="{y+27}" text-anchor="middle" font-size="21" font-weight="620" fill="#486779">{escape(item[1])}</text></g>'
        return base(title, subtitle, body)
    if stem in {"21-from-to", "22-transport"}:
        title, subtitle, a, b, c = spec
        colors = ("#e1a53a", "#3998c2", "#49a468")
        xs = (260,800,1340)
        body = '<path d="M300 470H760M840 470H1300" stroke="#c83236" stroke-width="9" stroke-linecap="round" marker-end="url(#arrow)"/>'
        for x, item, color in zip(xs,(a,b,c),colors):
            body += f'<g filter="url(#shadow)"><circle cx="{x}" cy="470" r="138" fill="#fff" stroke="{color}" stroke-width="7"/><circle cx="{x}" cy="470" r="92" fill="{color}" opacity=".16"/><text x="{x}" y="455" text-anchor="middle" font-size="25" font-weight="790" fill="#17364e">{escape(item[0])}</text><text x="{x}" y="500" text-anchor="middle" font-size="23" font-weight="650" fill="#486779">{escape(item[1])}</text></g>'
        return base(title, subtitle, body)
    return three_cards(spec)


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    for stem, localized in TEXT.items():
        for lang, spec in localized.items():
            (OUT / f"map-{stem}-{lang}.svg").write_text(render(stem, spec), encoding="utf-8")
    print(f"Generated {len(TEXT) * 3} localized School and City maps")


if __name__ == "__main__":
    main()
