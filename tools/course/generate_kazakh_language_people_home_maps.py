#!/usr/bin/env python3
"""Generate precise localized support maps for Kazakh lessons 9–16."""
from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/kazakh-language"

MAPS = [
    ("09-short-vowels", {
        "ru": ("Ы и І: кратко, но слышно", "отделяем звук, редукцию и настоящее выпадение", [("ы · задний ряд", ["қыс · ыдыс", "короткий импульс ●"]), ("і · передний ряд", ["тіс · іні", "короткий импульс ●"]), ("выпадение", ["ауыз + ы → аузы", "орын + ы → орны"])]),
        "kz": ("Ы мен І: қысқа, бірақ естіледі", "дыбыс, редукция және нақты түсуді ажыратамыз", [("ы · жуан", ["қыс · ыдыс", "қысқа серпін ●"]), ("і · жіңішке", ["тіс · іні", "қысқа серпін ●"]), ("дауыстының түсуі", ["ауыз + ы → аузы", "орын + ы → орны"])]),
        "en": ("Ы and І: short but audible", "separate sound, reduction, and genuine deletion", [("ы · back lane", ["қыс · ыдыс", "one short pulse ●"]), ("і · front lane", ["тіс · іні", "one short pulse ●"]), ("deletion", ["ауыз + ы → аузы", "орын + ы → орны"])]),
    }),
    ("10-family", {
        "ru": ("Кто это?", "не угадываем человека — спрашиваем и связываем", [("внимание", ["смотрим на значок", "личные данные не нужны"]), ("вопрос", ["Бұл кім?", "кім? → человек"]), ("ответ", ["Бұл — менің анам", "әке · әже · ата · іні"])]),
        "kz": ("Бұл кім?", "адамды болжамаймыз — сұрап, байланыстырамыз", [("назар", ["белгіге қараймыз", "жеке дерек қажет емес"]), ("сұрақ", ["Бұл кім?", "кім? → адам"]), ("жауап", ["Бұл — менің анам", "әке · әже · ата · іні"])]),
        "en": ("Who is this?", "do not guess — ask and connect the relationship", [("focus", ["look at the icon", "no private detail"]), ("question", ["Бұл кім?", "кім? → person"]), ("answer", ["Бұл — менің анам", "әке · әже · ата · іні"])]),
    }),
    ("11-possessives", {
        "ru": ("Мост принадлежности", "владелец и предмет указывают на одно лицо", [("владелец", ["менің · сіздің", "оның · Аружанның"]), ("связь", ["кімнің? ↔ несі?", "два берега"]), ("предмет", ["кітабым · кітабыңыз", "кітабы"])]),
        "kz": ("Иелік көпірі", "иесі мен заты бір жақты көрсетеді", [("иесі", ["менің · сіздің", "оның · Аружанның"]), ("байланыс", ["кімнің? ↔ несі?", "екі жаға"]), ("заты", ["кітабым · кітабыңыз", "кітабы"])]),
        "en": ("The possession bridge", "possessor and item point to the same person", [("possessor", ["менің · сіздің", "оның · Аружанның"]), ("link", ["кімнің? ↔ несі?", "both banks agree"]), ("item", ["кітабым · кітабыңыз", "кітабы"])]),
    }),
    ("12-plurals", {
        "ru": ("Два сигнала — одно окончание", "последняя гласная выбирает а/е, последний звук — л/д/т", [("1 · гармония", ["жуан → а", "жіңішке → е"]), ("2 · последний звук", ["гласн., й/р/у → л", "л/м/н/ң/з/ж → д", "глухой → т"]), ("результат", ["балалар · адамдар", "кітаптар · үйлер"])]),
        "kz": ("Екі белгі — бір жалғау", "соңғы дауысты а/е, соңғы дыбыс л/д/т таңдайды", [("1 · үндестік", ["жуан → а", "жіңішке → е"]), ("2 · соңғы дыбыс", ["дауысты, й/р/у → л", "л/м/н/ң/з/ж → д", "қатаң → т"]), ("нәтиже", ["балалар · адамдар", "кітаптар · үйлер"])]),
        "en": ("Two signals, one ending", "the final vowel selects а/е; final sound selects л/д/т", [("1 · harmony", ["back → а", "front → е"]), ("2 · final sound", ["vowel, й/р/у → л", "л/м/н/ң/з/ж → д", "voiceless → т"]), ("result", ["балалар · адамдар", "кітаптар · үйлер"])]),
    }),
    ("13-numbers", {
        "ru": ("Число уже показывает количество", "после точного числа предмет остаётся в единственном числе", [("посчитай", ["бір · екі · үш", "төрт · бес · алты"]), ("назови", ["үш бөлме", "төрт кітап"]), ("возраст", ["Ол жеті жаста", "Мен он жастамын"])]),
        "kz": ("Мөлшерді сан көрсетеді", "нақты саннан кейін зат есім жекеше тұрады", [("сана", ["бір · екі · үш", "төрт · бес · алты"]), ("ата", ["үш бөлме", "төрт кітап"]), ("жасы", ["Ол жеті жаста", "Мен он жастамын"])]),
        "en": ("The number already marks quantity", "after an exact number, the noun stays singular", [("count", ["бір · екі · үш", "төрт · бес · алты"]), ("name", ["үш бөлме", "төрт кітап"]), ("age", ["Ол жеті жаста", "Мен он жастамын"])]),
    }),
    ("14-home", {
        "ru": ("Координата отвечает на Қайда?", "гармония выбирает а/е, последний звук выбирает д/т", [("место", ["үй · бөлме · үстел", "шкаф · мектеп"]), ("выбор", ["жуан/жіңішке → а/е", "звонкий/глухой → д/т"]), ("координата", ["үйде · үстелде", "шкафта · мектепте"])]),
        "kz": ("Координат Қайда? сұрағына жауап береді", "үндестік а/е, соңғы дыбыс д/т таңдайды", [("орын", ["үй · бөлме · үстел", "шкаф · мектеп"]), ("таңдау", ["жуан/жіңішке → а/е", "ұяң/қатаң → д/т"]), ("координат", ["үйде · үстелде", "шкафта · мектепте"])]),
        "en": ("A coordinate answers Қайда?", "harmony selects а/е and final sound selects д/т", [("place", ["үй · бөлме · үстел", "шкаф · мектеп"]), ("select", ["back/front → а/е", "voiced/voiceless → д/т"]), ("coordinate", ["үйде · үстелде", "шкафта · мектепте"])]),
    }),
    ("15-existence", {
        "ru": ("Сцена → предмет → результат", "бар сообщает наличие, жоқ — отсутствие", [("сцена", ["Бөлмеде…", "Менің…"]), ("предмет", ["үстел · кітап", "інім · балкон"]), ("результат", ["бар · жоқ", "вопрос: бар ма?"])]),
        "kz": ("Орын → зат → нәтиже", "бар болуды, жоқ болмауды білдіреді", [("орын/ие", ["Бөлмеде…", "Менің…"]), ("зат", ["үстел · кітап", "інім · балкон"]), ("нәтиже", ["бар · жоқ", "сұрақ: бар ма?"])]),
        "en": ("Scene → item → result", "бар states presence; жоқ states absence", [("scene", ["Бөлмеде…", "Менің…"]), ("item", ["үстел · кітап", "інім · балкон"]), ("result", ["бар · жоқ", "question: бар ма?"])]),
    }),
    ("16-mastery", {
        "ru": ("Экскурсия без сценария", "семь опор соединяются в один новый разговор", [("люди", ["ы/і · Бұл кім?", "кімнің? · сан"]), ("дом", ["қайда?", "бар · жоқ · бар ма?"]), ("доказательство", ["8/10 сейчас", "7/10 через 7 дней"])]),
        "kz": ("Дайын мәтінсіз экскурсия", "жеті тірек бір жаңа әңгімеге бірігеді", [("адамдар", ["ы/і · Бұл кім?", "кімнің? · сан"]), ("үй", ["қайда?", "бар · жоқ · бар ма?"]), ("дәлел", ["қазір 8/10", "7 күннен кейін 7/10"])]),
        "en": ("An unscripted home tour", "seven supports join in one new conversation", [("people", ["ы/і · Бұл кім?", "кімнің? · number"]), ("home", ["қайда?", "бар · жоқ · бар ма?"]), ("evidence", ["8/10 now", "7/10 after 7 days"])]),
    }),
]


def multiline(items, x, y, size=26, gap=44, weight=600):
    spans = "".join(f'<tspan x="{x}" y="{y + i * gap}">{escape(item)}</tspan>' for i, item in enumerate(items))
    return f'<text text-anchor="middle" font-size="{size}" font-weight="{weight}" fill="#183b52">{spans}</text>'


def render(title, subtitle, cards):
    xs = (88, 590, 1092)
    palettes = (("#fff7dd", "#e4a83e"), ("#e8f6ff", "#2c91bd"), ("#eaf8ed", "#3e9c61"))
    parts = []
    for idx, ((heading, body), x, (fill, stroke)) in enumerate(zip(cards, xs, palettes)):
        parts.append(f'''<g filter="url(#cardShadow)">
          <path d="M{x+24} 276 H{x+390} Q{x+420} 276 {x+420} 306 V626 Q{x+420} 654 {x+392} 654 H{x+24} Q{x} 654 {x} 630 V300 Q{x} 276 {x+24} 276Z" fill="{fill}" stroke="{stroke}" stroke-width="3"/>
          <path d="M{x+20} 296 H{x+400} V374 H{x+20}Z" rx="14" fill="{stroke}" opacity=".14"/>
          <path d="M{x+18} 630 H{x+402}" stroke="{stroke}" stroke-width="10" opacity=".18" stroke-linecap="round"/>
          <text x="{x+210}" y="346" text-anchor="middle" font-size="29" font-weight="800" fill="#17364e">{escape(heading)}</text>
          {multiline(body, x+210, 432, 26 if len(body) < 3 else 24, 48)}
        </g>''')
        if idx < 2:
            parts.append(f'''<g filter="url(#arrowShadow)">
              <path d="M{x+442} 463 H{x+477}" stroke="#c83236" stroke-width="8" stroke-linecap="round"/>
              <path d="M{x+468} 449 L{x+486} 463 L{x+468} 477" fill="none" stroke="#c83236" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"/>
            </g>''')
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}">
      <defs>
        <linearGradient id="paper" x1="0" y1="0" x2="0" y2="1"><stop stop-color="#ffffff"/><stop offset="1" stop-color="#eaf4f9"/></linearGradient>
        <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#8fb8d0" stroke-width="1" opacity=".28"/></pattern>
        <filter id="cardShadow" x="-15%" y="-15%" width="140%" height="150%"><feDropShadow dx="0" dy="15" stdDeviation="13" flood-color="#17364e" flood-opacity=".20"/></filter>
        <filter id="arrowShadow" x="-30%" y="-50%" width="180%" height="200%"><feDropShadow dx="0" dy="4" stdDeviation="3" flood-color="#69171a" flood-opacity=".25"/></filter>
      </defs>
      <rect width="1600" height="900" rx="34" fill="url(#paper)"/>
      <rect width="1600" height="900" rx="34" fill="url(#grid)"/>
      <path d="M66 218 H1534" stroke="#c83236" stroke-width="4" opacity=".75"/>
      <text x="800" y="102" text-anchor="middle" font-family="Inter,Arial,sans-serif" font-size="52" font-weight="820" fill="#17364e">{escape(title)}</text>
      <text x="800" y="164" text-anchor="middle" font-family="Inter,Arial,sans-serif" font-size="25" font-weight="540" fill="#4b697d">{escape(subtitle)}</text>
      <g font-family="Inter,Arial,sans-serif">{''.join(parts)}</g>
      <g transform="translate(300 732)" filter="url(#cardShadow)">
        <rect width="1000" height="92" rx="22" fill="#17364e"/>
        <path d="M24 18 H976" stroke="#5fa4c7" stroke-width="2" opacity=".55"/>
        <text x="500" y="59" text-anchor="middle" font-family="Inter,Arial,sans-serif" font-size="29" font-weight="760" fill="#fff">КӨР → ТЫҢДА → АЙТ → ҚОЛДАН → ТЕКСЕР</text>
      </g>
    </svg>'''


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    for stem, localized in MAPS:
        for lang, spec in localized.items():
            (OUT / f"map-{stem}-{lang}.svg").write_text(render(*spec), encoding="utf-8")
    print(f"Generated {len(MAPS) * 3} localized People and Home maps")


if __name__ == "__main__":
    main()
