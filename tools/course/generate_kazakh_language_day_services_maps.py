#!/usr/bin/env python3
"""Generate precise, localized support maps for Kazakh lessons 25–32."""
from __future__ import annotations

from html import escape
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/kazakh-language"


TEXT = {
    "25-daily-routine": {
        "ru": ("Мой обычный день", "Время отвечает «когда?», личная форма — «что делаю я?»", [("07:00", "тұрамын", "встаю"), ("08:00", "үйден шығамын", "выхожу из дома"), ("КЕШКЕ", "үйге қайтамын", "возвращаюсь домой")]),
        "kz": ("Менің күн тәртібім", "Уақыт «қашан?», жіктік форма «мен не істеймін?» дейді", [("07:00", "тұрамын", "ұйқыдан тұрамын"), ("08:00", "үйден шығамын", "сыртқа шығамын"), ("КЕШКЕ", "үйге қайтамын", "үйге ораламын")]),
        "en": ("My usual day", "Time answers ‘when?’; the personal form answers ‘what do I do?’", [("07:00", "тұрамын", "I get up"), ("08:00", "үйден шығамын", "I leave home"), ("EVENING", "үйге қайтамын", "I return home")]),
    },
    "26-sequence": {
        "ru": ("Из списка — в понятный рассказ", "Положение карточки определяет связку", [("1", "АЛДЫМЕН", "сначала"), ("2–3", "СОДАН КЕЙІН", "после этого"), ("4", "СОҢЫНДА", "в конце")]),
        "kz": ("Тізімнен түсінікті әңгімеге", "Карточканың орны байланыстырушыны анықтайды", [("1", "АЛДЫМЕН", "бастапқы қадам"), ("2–3", "СОДАН КЕЙІН", "келесі қадам"), ("4", "СОҢЫНДА", "соңғы қадам")]),
        "en": ("From a list to a clear account", "A card's position determines its connector", [("1", "АЛДЫМЕН", "first"), ("2–3", "СОДАН КЕЙІН", "after that"), ("4", "СОҢЫНДА", "finally")]),
    },
    "27-shop": {
        "ru": ("Покупка без недоразумения", "Назовите товар, меру и цену; затем сверьте чек", [("НЕ?", "бір бөлке нан", "одна буханка хлеба"), ("ҚАНША?", "бір келі алма", "один килограмм яблок"), ("БАҒА + ЧЕК", "қанша тұрады?", "цена и итог совпали?")]),
        "kz": ("Түсінікті сатып алу", "Тауарды, өлшемді, бағаны атап, чекті тексеріңіз", [("НЕ?", "бір бөлке нан", "тауар және дана"), ("ҚАНША?", "бір келі алма", "тауар және өлшем"), ("БАҒА + ЧЕК", "қанша тұрады?", "баға мен қорытынды")]),
        "en": ("A purchase without confusion", "Name the product, measure, and price; then check the receipt", [("WHAT?", "бір бөлке нан", "one loaf of bread"), ("HOW MUCH?", "бір келі алма", "one kilogram of apples"), ("PRICE + RECEIPT", "қанша тұрады?", "do price and total match?")]),
    },
    "28-cafe": {
        "ru": ("Заказ: сначала понять, потом подтвердить", "Фотография не сообщает состав блюда", [("1", "ІШІНДЕ НЕ БАР?", "уточните состав"), ("2", "БЕРІҢІЗШІ", "сделайте выбор"), ("3", "ҚОСПАҢЫЗШЫ", "назовите ограничение")]),
        "kz": ("Тапсырыс: алдымен түсіну, кейін растау", "Сурет тағамның құрамын көрсетпейді", [("1", "ІШІНДЕ НЕ БАР?", "құрамды нақтылаңыз"), ("2", "БЕРІҢІЗШІ", "таңдау жасаңыз"), ("3", "ҚОСПАҢЫЗШЫ", "шектеуді айтыңыз")]),
        "en": ("Order: understand first, confirm second", "A photograph does not establish a dish's ingredients", [("1", "ІШІНДЕ НЕ БАР?", "clarify ingredients"), ("2", "БЕРІҢІЗШІ", "make a choice"), ("3", "ҚОСПАҢЫЗШЫ", "state a restriction")]),
    },
    "29-appointment": {
        "ru": ("Запись завершает подтверждение", "Предложенное время и подтверждённое время — разные этапы", [("ЗАПРОС", "бос уақыт бар ма?", "есть свободное время?"), ("ВАРИАНТ", "сейсенбі, 15:00", "предложение"), ("ПРОВЕРКА", "дұрыс па?", "явное подтверждение")]),
        "kz": ("Жазылу растаумен аяқталады", "Ұсынылған және расталған уақыт — екі бөлек кезең", [("СҰРАУ", "бос уақыт бар ма?", "бос терезені табу"), ("НҰСҚА", "сейсенбі, 15:00", "ұсынылған уақыт"), ("ТЕКСЕРУ", "дұрыс па?", "анық растау")]),
        "en": ("A booking ends with confirmation", "A proposed time and a confirmed time are separate stages", [("REQUEST", "бос уақыт бар ма?", "find an open slot"), ("OPTION", "сейсенбі, 15:00", "proposed time"), ("CHECK", "дұрыс па?", "explicit confirmation")]),
    },
    "30-health": {
        "ru": ("Сообщить наблюдение — обратиться к специалисту", "Урок помогает говорить; диагноз ставит не схема", [("1", "ҚАЛАЙ?", "плохо себя чувствую"), ("2", "ҚАШАН?", "началось вчера вечером"), ("3", "МАМАН", "нужна помощь врача")]),
        "kz": ("Белгіні айту — маманға жүгіну", "Сабақ сөйлеуге көмектеседі; сызба диагноз қоймайды", [("1", "ҚАЛАЙ?", "өзімді нашар сезінемін"), ("2", "ҚАШАН?", "кеше кешке басталды"), ("3", "МАМАН", "дәрігер көмегі керек")]),
        "en": ("Report an observation — contact a professional", "The lesson supports communication; the map does not diagnose", [("1", "ҚАЛАЙ?", "I feel unwell"), ("2", "ҚАШАН?", "started yesterday evening"), ("3", "МАМАН", "a doctor's help is needed")]),
    },
    "31-service-error": {
        "ru": ("Ошибка становится исправимой, когда она точна", "Сравните одно поле и попросите следующий шаг", [("ГДЕ?", "тегім", "поле «фамилия»"), ("ЧТО НЕ ТАК?", "дұрыс жазылмаған", "написано неверно"), ("ЧТО ДЕЛАТЬ?", "түзетуге бола ма?", "можно исправить?")]),
        "kz": ("Нақты қате түзетіледі", "Бір өрісті салыстырып, келесі қадамды сұраңыз", [("ҚАЙДА?", "тегім", "«тегі» өрісі"), ("НЕ ДҰРЫС ЕМЕС?", "дұрыс жазылмаған", "қате жазылған"), ("НЕ ІСТЕЙМІЗ?", "түзетуге бола ма?", "түзету мүмкіндігі")]),
        "en": ("A precise error can be corrected", "Compare one field and request the next step", [("WHERE?", "тегім", "the surname field"), ("WHAT IS WRONG?", "дұрыс жазылмаған", "written incorrectly"), ("WHAT NEXT?", "түзетуге бола ма?", "can it be corrected?")]),
    },
    "32-mastery": {
        "ru": ("Моя среда 1.0: один день, семь задач", "Новая карточка проверяет понимание, а не память", [("1–2", "КҮН + РЕТ", "режим и порядок"), ("3–5", "ДҮКЕН + КАФЕ + ЖАЗЫЛУ", "три бытовые задачи"), ("6–7", "КӨМЕК + ҚАТЕ", "помощь и исправление")]),
        "kz": ("Менің ортам 1.0: бір күн, жеті міндет", "Жаңа карточка жаттауды емес, түсінуді тексереді", [("1–2", "КҮН + РЕТ", "күн тәртібі мен рет"), ("3–5", "ДҮКЕН + КАФЕ + ЖАЗЫЛУ", "үш тұрмыстық міндет"), ("6–7", "КӨМЕК + ҚАТЕ", "көмек және түзету")]),
        "en": ("My World 1.0: one day, seven tasks", "An unseen card tests understanding rather than memory", [("1–2", "КҮН + РЕТ", "routine and sequence"), ("3–5", "ДҮКЕН + КАФЕ + ЖАЗЫЛУ", "three everyday tasks"), ("6–7", "КӨМЕК + ҚАТЕ", "help and correction")]),
    },
}


def split_lines(text: str, limit: int = 24) -> list[str]:
    words, rows, line = text.split(), [], []
    for word in words:
        if line and len(" ".join(line + [word])) > limit:
            rows.append(" ".join(line)); line = [word]
        else:
            line.append(word)
    if line: rows.append(" ".join(line))
    return rows


def tspans(rows: list[str], x: int, y: int, size: int, gap: int, weight: int, fill: str = "#17364e") -> str:
    return f'<text x="{x}" y="{y}" text-anchor="middle" font-size="{size}" font-weight="{weight}" fill="{fill}">' + "".join(f'<tspan x="{x}" dy="{0 if i == 0 else gap}">{escape(row)}</tspan>' for i, row in enumerate(rows)) + '</text>'


def render(title: str, subtitle: str, cards: list[tuple[str, str, str]]) -> str:
    title_rows = split_lines(title, 42)
    subtitle_rows = split_lines(subtitle, 76)
    title_y = 78 if len(title_rows) == 1 else 58
    subtitle_y = 150 if len(title_rows) == 1 else 166
    chunks = []
    colors = ("#e2a13b", "#35a0c8", "#48a36b")
    for i, (badge, main, detail) in enumerate(cards):
        x = 75 + i * 510
        chunks.append(f'''<g transform="translate({x} 292)" filter="url(#cardShadow)">
          <path d="M26 0H414L460 46V352A28 28 0 0 1 432 380H28A28 28 0 0 1 0 352V28A28 28 0 0 1 28 0Z" fill="#fff" stroke="{colors[i]}" stroke-width="4"/>
          <path d="M414 0V46H460" fill="none" stroke="{colors[i]}" stroke-width="4"/>
          <rect x="24" y="24" width="190" height="58" rx="16" fill="{colors[i]}"/>
          {tspans(split_lines(badge, 20), 119, 62, 19, 22, 850, '#fff')}
          <circle cx="230" cy="151" r="48" fill="{colors[i]}" opacity=".17"/>
          <circle cx="230" cy="151" r="34" fill="none" stroke="{colors[i]}" stroke-width="6"/>
          {tspans(split_lines(main, 22), 230, 242 if len(split_lines(main,22)) == 1 else 222, 27, 36, 850)}
          {tspans(split_lines(detail, 30), 230, 307 if len(split_lines(detail,30)) == 1 else 287, 22, 31, 620, '#4f6f80')}
        </g>''')
    arrows = '<path d="M540 482H568" stroke="#c7353a" stroke-width="9" stroke-linecap="round" marker-end="url(#arrow)"/><path d="M1050 482H1078" stroke="#c7353a" stroke-width="9" stroke-linecap="round" marker-end="url(#arrow)"/>'
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}">
<defs>
 <linearGradient id="paper" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#fbf8ef"/><stop offset="1" stop-color="#edf7f8"/></linearGradient>
 <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#7ba6b8" stroke-opacity=".16" stroke-width="1.5"/></pattern>
 <filter id="cardShadow" x="-15%" y="-15%" width="140%" height="150%"><feDropShadow dx="0" dy="15" stdDeviation="12" flood-color="#17364e" flood-opacity=".24"/></filter>
 <filter id="barShadow" x="-10%" y="-30%" width="120%" height="180%"><feDropShadow dx="0" dy="9" stdDeviation="8" flood-color="#17364e" flood-opacity=".23"/></filter>
 <marker id="arrow" markerWidth="12" markerHeight="12" refX="9" refY="6" orient="auto"><path d="M0 0L12 6L0 12Z" fill="#c7353a"/></marker>
</defs>
<rect width="1600" height="900" rx="34" fill="url(#paper)"/><rect width="1600" height="900" rx="34" fill="url(#grid)"/>
<g font-family="Inter,Arial,sans-serif">{tspans(title_rows,800,title_y,48,52,870)}{tspans(subtitle_rows,800,subtitle_y,24,31,580,'#4f6f80')}
{''.join(chunks)}{arrows}
<g transform="translate(284 774)" filter="url(#barShadow)"><rect width="1032" height="70" rx="20" fill="#17364e"/>{tspans(['КӨР → ТЫҢДА → АЙТ → ҚОЛДАН → ТЕКСЕР'],516,45,25,30,800,'#fff')}</g>
</g></svg>'''


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    for stem, localized in TEXT.items():
        for lang, spec in localized.items():
            (OUT / f"map-{stem}-{lang}.svg").write_text(render(*spec), encoding="utf-8")
    print(f"Generated {len(TEXT) * 3} localized Daily Life and Services maps")


if __name__ == "__main__":
    main()
