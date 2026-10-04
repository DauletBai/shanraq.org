#!/usr/bin/env python3
"""Generate clean, localized support maps for Kazakh lessons 41–48.

The maps keep one visual grammar throughout the block: a time line, a scene
card, and a check card. Text is wrapped before it reaches the cards so that
the diagrams stay readable on phones as well as on the lesson page.
"""
from __future__ import annotations

from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/kazakh-language"

TEXT = {
    "41-time-map": {
        "ru": ("Время видно в сцене", "маяк → состояние действия → форма", ["ВЧЕРА", "ОБЫЧНО", "СЕЙЧАС", "ЗАВТРА"]),
        "kz": ("Уақыт көріністен байқалады", "белгі → әрекет күйі → форма", ["КЕШЕ", "ӘДЕТТЕ", "ҚАЗІР", "ЕРТЕҢ"]),
        "en": ("Time is visible in a scene", "beacon → action state → form", ["YESTERDAY", "USUALLY", "NOW", "TOMORROW"]),
    },
    "42-habit": {
        "ru": ("Повторение собирает форму", "участник + основа + личное окончание", ["МЕН", "СЕН", "ОЛ", "БІЗ"]),
        "kz": ("Қайталану форманы құрайды", "қатысушы + түбір + жіктік жалғауы", ["МЕН", "СЕН", "ОЛ", "БІЗ"]),
        "en": ("Repetition builds the form", "participant + stem + personal ending", ["МЕН", "СЕН", "ОЛ", "БІЗ"]),
    },
    "43-progressive": {
        "ru": ("Сейчас действие ещё движется", "сцена открыта → -ып/-іп жатыр", ["СЦЕНА", "ҚАЗІР", "ҮДЕРІС", "ТЕКСЕРКА"]),
        "kz": ("Қазір әрекет әлі жүріп жатыр", "көрініс ашық → -ып/-іп жатыр", ["КӨРІНІС", "ҚАЗІР", "ҮДЕРІС", "ТЕКСЕРУ"]),
        "en": ("The action is still moving now", "open scene → -ып/-іп жатыр", ["SCENE", "NOW", "PROCESS", "CHECK"]),
    },
    "44-past": {
        "ru": ("Прошлое оставляет результат", "событие завершилось → итог виден", ["СОБЫТИЕ", "ВЧЕРА", "РЕЗУЛЬТАТ", "РАССКАЗ"]),
        "kz": ("Өткен шақ нәтиже қалдырады", "әрекет аяқталды → нәтиже көрінеді", ["ӘРЕКЕТ", "КЕШЕ", "НӘТИЖЕ", "ӘҢГІМЕ"]),
        "en": ("The past leaves a result", "event completed → result is visible", ["EVENT", "YESTERDAY", "RESULT", "ACCOUNT"]),
    },
    "45-future": {
        "ru": ("План ведёт к следующему кадру", "намерение + время → договорённое действие", ["НАМЕРЕНИЕ", "КҮН", "КЕЛЕСІ ҚАДАМ", "ПЛАН"]),
        "kz": ("Жоспар келесі кадрға апарады", "ниет + уақыт → келісілген әрекет", ["НИЕТ", "КҮН", "КЕЛЕСІ ҚАДАМ", "ЖОСПАР"]),
        "en": ("A plan leads to the next frame", "intention + time → agreed action", ["INTENTION", "DAY", "NEXT STEP", "PLAN"]),
    },
    "46-questions": {
        "ru": ("Вопрос проверяет, отрицание уточняет", "ма/ме отдельно; отрицание внутри действия", ["СЦЕНА", "СҰРАҚ", "ТЕРІС", "ЖАУАП"]),
        "kz": ("Сұрақ тексереді, болымсыздық нақтылайды", "ма/ме бөлек; болымсыздық әрекет ішінде", ["КӨРІНІС", "СҰРАҚ", "БОЛЫМСЫЗ", "ЖАУАП"]),
        "en": ("A question checks; negation clarifies", "ма/ме separate; negation inside the action", ["SCENE", "QUESTION", "NEGATIVE", "REPLY"]),
    },
    "47-story": {
        "ru": ("Шесть кадров становятся рассказом", "время + связка + причина/следствие", ["КЕШЕ", "СОДАН КЕЙІН", "ҚАЗІР", "СЕБЕБІ", "СОНДЫҚТАН", "ЕРТЕҢ"]),
        "kz": ("Алты кадр әңгімеге айналады", "уақыт + байланыс + себеп/салдар", ["КЕШЕ", "СОДАН КЕЙІН", "ҚАЗІР", "СЕБЕБІ", "СОНДЫҚТАН", "ЕРТЕҢ"]),
        "en": ("Six frames become a story", "time + connector + cause/consequence", ["YESTERDAY", "AFTER THAT", "NOW", "BECAUSE", "THEREFORE", "TOMORROW"]),
    },
    "48-mastery": {
        "ru": ("Контрольная точка 1.2", "вчера → обычно → сейчас → план → перенос", ["ПРОШЛОЕ", "ПРИВЫЧКА", "ПРОЦЕСС", "ПЛАН", "8/10"]),
        "kz": ("1.2 бақылау нүктесі", "өткен → әдет → қазір → жоспар → тасымал", ["ӨТКЕН", "ӘДЕТ", "ҮДЕРІС", "ЖОСПАР", "8/10"]),
        "en": ("Mastery checkpoint 1.2", "past → habit → now → plan → transfer", ["PAST", "HABIT", "PROCESS", "PLAN", "8/10"]),
    },
}


def wrap(value: str, limit: int) -> list[str]:
    words, rows, line = value.split(), [], []
    for word in words:
        if line and len(" ".join(line + [word])) > limit:
            rows.append(" ".join(line))
            line = [word]
        else:
            line.append(word)
    if line:
        rows.append(" ".join(line))
    return rows


def text(value: str | list[str], x: int, y: int, size: int, gap: int, weight: int, fill: str = "#17364e") -> str:
    rows = [value] if isinstance(value, str) else value
    return f'<text x="{x}" y="{y}" text-anchor="middle" font-size="{size}" font-weight="{weight}" fill="{fill}">' + "".join(
        f'<tspan x="{x}" dy="{0 if index == 0 else gap}">{escape(row)}</tspan>' for index, row in enumerate(rows)
    ) + "</text>"


def render(title: str, subtitle: str, labels: list[str]) -> str:
    title_rows = wrap(title, 42)
    subtitle_rows = wrap(subtitle, 72)
    title_y = 72 if len(title_rows) == 1 else 48
    subtitle_y = 142 if len(title_rows) == 1 else 158
    colors = ["#d28b24", "#2b91b8", "#438f64", "#b84d74", "#7161ae", "#d0653d"]
    n = len(labels)
    left, right = 130, 1470
    span = (right - left) / max(n - 1, 1)
    points = []
    for index, label in enumerate(labels):
        x = left + index * span
        y = 420 + (index % 2) * 26
        color = colors[index % len(colors)]
        points.append(f'<g filter="url(#shadow)"><circle cx="{x:.1f}" cy="{y}" r="46" fill="{color}" stroke="#fff" stroke-width="8"/><rect x="{x-105:.1f}" y="{y+72}" width="210" height="76" rx="18" fill="#fff" stroke="{color}" stroke-width="3"/>{text(wrap(label,16),int(x),y+8,20 if len(label) < 13 else 17,22,850)}{text("сцена",int(x),y+103,18,22,650,"#557383")}</g>')
        if index:
            previous = left + (index - 1) * span
            previous_y = 420 + ((index - 1) % 2) * 26
            points.insert(-1, f'<path d="M{previous+52:.1f} {previous_y}L{x-52:.1f} {y}" fill="none" stroke="#c7353a" stroke-width="7" stroke-linecap="round" marker-end="url(#arrow)"/>')
    body = "".join(points)
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}"><defs>
<linearGradient id="paper" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#fbf8ef"/><stop offset="1" stop-color="#e8f5f7"/></linearGradient>
<pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#7ba6b8" stroke-opacity=".16" stroke-width="1.5"/></pattern>
<filter id="shadow" x="-25%" y="-25%" width="160%" height="180%"><feDropShadow dx="0" dy="13" stdDeviation="10" flood-color="#17364e" flood-opacity=".24"/></filter>
<marker id="arrow" markerWidth="12" markerHeight="12" refX="10" refY="6" orient="auto"><path d="M0 0L12 6L0 12Z" fill="#c7353a"/></marker></defs>
<rect width="1600" height="900" rx="34" fill="url(#paper)"/><rect width="1600" height="900" rx="34" fill="url(#grid)"/>
<g font-family="Inter,Arial,sans-serif">{text(title_rows,800,title_y,48,52,870)}{text(subtitle_rows,800,subtitle_y,24,31,600,"#4f6f80")}
<path d="M{left} 420H{right}" stroke="#89aab7" stroke-width="8" stroke-linecap="round" opacity=".45"/>{body}
<g filter="url(#shadow)"><rect x="270" y="748" width="1060" height="70" rx="20" fill="#17364e"/>{text("КӨРІНІС → МАҒЫНА → ФОРМА → АЙТ → ТЕКСЕР",800,794,25,30,800,"#fff")}</g></g></svg>'''


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    for stem, localized in TEXT.items():
        for lang, spec in localized.items():
            (OUT / f"map-{stem}-{lang}.svg").write_text(render(*spec), encoding="utf-8")
    print(f"Generated {len(TEXT) * 3} localized Action and Time maps")


if __name__ == "__main__":
    main()
