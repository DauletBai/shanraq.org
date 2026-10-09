#!/usr/bin/env python3
"""Draw 24 localized vector support maps for algorithm lessons 19–26."""

from html import escape
from pathlib import Path
import textwrap

OUT = Path(__file__).resolve().parents[2] / "web/static/course/informatics"

# Four short, exact stages. The last line is the observable result from the lesson.
MAPS = {
    19: {
        "ru": ("Большая просьба", ("найти t-01", "прочесть срок", "сравнить с сегодня", "показать ответ"), "Нет срока → не сравниваем"),
        "kz": ("Үлкен өтініш", ("t-01 табу", "мерзімді оқу", "бүгінмен салыстыру", "жауап көрсету"), "Мерзім жоқ → салыстырмаймыз"),
        "en": ("A large request", ("find t-01", "read due date", "compare with today", "show result"), "No date → no comparison"),
    },
    20: {
        "ru": ("Состояние меняется", ("начало: 3", "ловушка: 2", "бонус: 4", "проверка: ≥ 0"), "done: false → true"),
        "kz": ("Күй өзгереді", ("басы: 3", "тұзақ: 2", "сыйлық: 4", "тексеру: ≥ 0"), "done: false → true"),
        "en": ("State changes", ("start: 3", "trap: 2", "bonus: 4", "check: ≥ 0"), "done: false → true"),
    },
    21: {
        "ru": ("Порядок имеет значение", ("count = 0", "+1 → 1", "показать 1", "+1 → 2"), "Сначала два +1 → другой экран"),
        "kz": ("Рет маңызды", ("count = 0", "+1 → 1", "1 көрсету", "+1 → 2"), "Әуелі екі +1 → өзге экран"),
        "en": ("Order matters", ("count = 0", "+1 → 1", "display 1", "+1 → 2"), "Two increments first → different screen"),
    },
    22: {
        "ru": ("Проверьте границу", ("−1: просрочено", "0: напомнить", "2: напомнить", "3: рано"), "0 ≤ дни ≤ 2 и done=false"),
        "kz": ("Шекараны тексер", ("−1: мерзімі өтті", "0: еске сал", "2: еске сал", "3: әлі ерте"), "0 ≤ күн ≤ 2 және done=false"),
        "en": ("Test the boundary", ("−1: overdue", "0: remind", "2: remind", "3: too early"), "0 ≤ days ≤ 2 and done=false"),
    },
    23: {
        "ru": ("Повторяем и останавливаем", ("позиция 0", "false → 0", "true → 1", "позиция 3: стоп"), "false, true, false → 1"),
        "kz": ("Қайталау және тоқтау", ("орын 0", "false → 0", "true → 1", "орын 3: тоқта"), "false, true, false → 1"),
        "en": ("Repeat, then stop", ("position 0", "false → 0", "true → 1", "position 3: stop"), "false, true, false → 1"),
    },
    24: {
        "ru": ("Функция обещает результат", ("done?", "дата есть?", "дней 0–2?", "вернуть ответ"), "Без срока → NO_DATE"),
        "kz": ("Функция жауап береді", ("done?", "мерзім бар ма?", "күн 0–2 ме?", "жауап қайтару"), "Мерзімсіз → NO_DATE"),
        "en": ("A function promises a result", ("done?", "date present?", "days 0–2?", "return outcome"), "Missing date → NO_DATE"),
    },
    25: {
        "ru": ("Сначала верно, потом быстро", ("4 записи", "6 пар", "seen: 4 проверки", "повтор t-02"), "n(n−1)/2 пар против прохода"),
        "kz": ("Әуелі дұрыс, кейін тез", ("4 жазба", "6 жұп", "seen: 4 тексеріс", "t-02 қайталанды"), "n(n−1)/2 жұп немесе бір өту"),
        "en": ("Correct first, then fast", ("4 records", "6 pairs", "seen: 4 checks", "duplicate t-02"), "n(n−1)/2 pairs versus one scan"),
    },
    26: {
        "ru": ("Контрольная точка 0.3", ("договор", "трасса", "границы", "чужая проверка"), "8/10 сейчас · 7/10 позже"),
        "kz": ("0.3 бақылау нүктесі", ("келісім", "трасса", "шекара", "өзге тексеріс"), "қазір 8/10 · кейін 7/10"),
        "en": ("Checkpoint 0.3", ("contract", "trace", "boundaries", "peer check"), "8/10 now · 7/10 later"),
    },
}

def lines(value: str, width: int = 19) -> str:
    parts = textwrap.wrap(value, width=width, break_long_words=False, break_on_hyphens=False)
    return "".join(f'<tspan x="24" dy="{0 if i == 0 else 34}">{escape(part)}</tspan>' for i, part in enumerate(parts))

def render(number: int, lang: str) -> str:
    title, stages, outcome = MAPS[number][lang]
    cards = []
    for i, label in enumerate(stages):
        x = 66 + i * 377
        color = ("#4d80ae", "#7163aa", "#bd705b", "#408d81")[i]
        number_label = f"{i + 1:02d}"
        cards.append(f'''<g transform="translate({x},255)">
          <rect x="9" y="13" width="332" height="277" rx="18" fill="#162333" opacity=".20"/>
          <rect width="332" height="277" rx="18" fill="url(#paper)" stroke="#b9c8d1" stroke-width="2"/>
          <path d="M0 19 Q0 0 19 0 H313 Q332 0 332 19 V71 H0 Z" fill="{color}"/>
          <text x="24" y="49" class="card-num">{number_label}</text>
          <text class="card-text" x="24" y="126">{lines(label)}</text>
          <path d="M25 213 H304" stroke="#d7e2e7" stroke-width="2"/>
          <circle cx="45" cy="238" r="8" fill="{color}"/>
          <circle cx="68" cy="238" r="8" fill="{color}" opacity=".45"/>
        </g>''')
    arrows = ''.join(f'<path d="M{405+i*377} 393 H{428+i*377}" stroke="#285773" stroke-width="4" stroke-linecap="round" marker-end="url(#arrow)"/>' for i in range(3))
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">
<title id="title">{escape(title)}</title><desc id="desc">{escape(' → '.join(stages))}</desc>
<defs><linearGradient id="bg" x2="0" y2="1"><stop stop-color="#edf7f9"/><stop offset="1" stop-color="#d8e9ee"/></linearGradient><linearGradient id="paper" x2="0" y2="1"><stop stop-color="#ffffff"/><stop offset="1" stop-color="#e9f2f5"/></linearGradient><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0 H0 V40" fill="none" stroke="#80a6b6" stroke-opacity=".22" stroke-width="1"/></pattern><marker id="arrow" viewBox="0 0 12 12" refX="10" refY="6" markerWidth="12" markerHeight="12" markerUnits="userSpaceOnUse" orient="auto"><path d="M1 1 L11 6 L1 11 Z" fill="#285773"/></marker></defs>
<style>.title{{font:700 48px system-ui,Arial,sans-serif;fill:#1c3448}}.card-num{{font:700 27px system-ui,Arial,sans-serif;fill:white}}.card-text{{font:650 27px system-ui,Arial,sans-serif;fill:#1a3346}}.outcome{{font:650 35px system-ui,Arial,sans-serif;fill:#153d50}}</style>
<rect width="1600" height="900" fill="url(#bg)"/><rect width="1600" height="900" fill="url(#grid)"/>
<text class="title" x="72" y="117">{escape(title)}</text><path d="M72 147 H1528" stroke="#86a8b6" stroke-width="3"/>
{''.join(cards)}{arrows}
<rect x="68" y="655" width="1464" height="122" rx="18" fill="#163a50" opacity=".18"/>
<rect x="60" y="647" width="1464" height="122" rx="18" fill="white" stroke="#b4cbd4" stroke-width="2"/>
<text class="outcome" x="90" y="722">{escape(outcome)}</text>
</svg>'''

def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    for number in range(19, 27):
        for lang in ("ru", "kz", "en"):
            (OUT / f"map-{number:02d}-{('problem-decomposition','state-variables','sequence-tracing','conditions-boundaries','loops-invariants','functions-contracts','correctness-efficiency','algorithms-mastery')[number-19]}-{lang}.svg").write_text(render(number, lang), encoding="utf-8")
    print("Generated 24 localized algorithm maps")

if __name__ == "__main__":
    main()
