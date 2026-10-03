#!/usr/bin/env python3
"""Generate the localized support-signal maps for Kazakh block one."""
from html import escape
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/kazakh-language"

MAPS = [
    ("01-first-contact", {
        "ru": ("Первые 30 секунд", "не переводим всё — держим маршрут разговора", [("1 · контакт", ["Сәлем!", "взгляд + приветствие"]), ("2 · мост", ["Менің атым…", "Сіздің атыңыз кім?"]), ("3 · ответ", ["Менің атым…", "Танысқаныма", "қуаныштымын"])]),
        "kz": ("Алғашқы 30 секунд", "бәрін аудармаймыз — әңгіменің бағытын ұстаймыз", [("1 · байланыс", ["Сәлем!", "көзқарас + амандасу"]), ("2 · көпір", ["Менің атым…", "Сіздің атыңыз кім?"]), ("3 · жауап", ["Менің атым…", "Танысқаныма", "қуаныштымын"])]),
        "en": ("The first 30 seconds", "do not translate everything — hold the route", [("1 · contact", ["Сәлем!", "look + greeting"]), ("2 · bridge", ["Менің атым…", "Сіздің атыңыз кім?"]), ("3 · response", ["Менің атым…", "Танысқаныма", "қуаныштымын"])]),
    }),
    ("02-nine-sounds", {
        "ru": ("Девять букв — девять ориентиров", "сначала положение рта и звук, затем буква", [("гласные", ["Ә · Ө · Ү", "Ұ · І"]), ("согласные", ["Қ · Ғ · Ң · Һ", "не русские замены"]), ("цикл", ["слушай → смотри", "скажи → сравни"])]),
        "kz": ("Тоғыз әріп — тоғыз бағдар", "алдымен ауыз қалпы мен дыбыс, содан кейін әріп", [("дауыстылар", ["Ә · Ө · Ү", "Ұ · І"]), ("дауыссыздар", ["Қ · Ғ · Ң · Һ", "орысша баламасы емес"]), ("айналым", ["тыңда → қара", "айт → салыстыр"])]),
        "en": ("Nine letters — nine anchors", "mouth position and sound first, then the letter", [("vowels", ["Ә · Ө · Ү", "Ұ · І"]), ("consonants", ["Қ · Ғ · Ң · Һ", "not Russian substitutes"]), ("cycle", ["listen → watch", "say → compare"])]),
    }),
    ("03-harmony", {
        "ru": ("Две звуковые семьи", "гласная основы выбирает гласную окончания", [("задний ряд", ["а · о · ұ · ы", "бала → балалар"]), ("передний ряд", ["ә · ө · ү · і · е", "үй → үйлер"]), ("проверка", ["услышь основу", "выбери семью"])]),
        "kz": ("Екі дыбыстық отбасы", "түбір дауыстысы қосымша дауыстысын таңдайды", [("жуан", ["а · о · ұ · ы", "бала → балалар"]), ("жіңішке", ["ә · ө · ү · і · е", "үй → үйлер"]), ("тексеру", ["түбірді тыңда", "отбасын таңда"])]),
        "en": ("Two sound families", "the stem vowel selects the suffix vowel", [("back vowels", ["а · о · ұ · ы", "бала → балалар"]), ("front vowels", ["ә · ө · ү · і · е", "үй → үйлер"]), ("check", ["hear the stem", "choose the family"])]),
    }),
    ("04-sentence", {
        "ru": ("Скелет предложения", "сначала участники и детали — действие завершает мысль", [("кто?", ["Мен", "Аружан"]), ("что? где?", ["кітап", "мектепке"]), ("что делает?", ["оқимын", "барады"])]),
        "kz": ("Сөйлем қаңқасы", "алдымен қатысушылар мен дерек — әрекет ойды аяқтайды", [("кім?", ["Мен", "Аружан"]), ("не? қайда?", ["кітап", "мектепке"]), ("не істейді?", ["оқимын", "барады"])]),
        "en": ("Sentence skeleton", "participants and detail first — the action closes the thought", [("who?", ["Мен", "Аружан"]), ("what? where?", ["кітап", "мектепке"]), ("does what?", ["оқимын", "барады"])]),
    }),
    ("05-introduction", {
        "ru": ("Собираем «я»", "смысл → основа → личное окончание", [("имя", ["Менің атым — Алина", "готовая рамка"]), ("роль", ["Мен оқушымын", "оқушы + мын"]), ("откуда", ["Мен Қостанайданмын", "место + дан + мын"])]),
        "kz": ("«Мен» үлгісін жинаймыз", "мағына → түбір → жіктік жалғауы", [("аты", ["Менің атым — Алина", "дайын қалып"]), ("рөлі", ["Мен оқушымын", "оқушы + мын"]), ("қайдан", ["Мен Қостанайданмын", "орын + дан + мын"])]),
        "en": ("Build the ‘I’ pattern", "meaning → stem → personal ending", [("name", ["Менің атым — Алина", "a fixed frame"]), ("role", ["Мен оқушымын", "оқушы + мын"]), ("from", ["Мен Қостанайданмын", "place + дан + мын"])]),
    }),
    ("06-questions", {
        "ru": ("Вопрос открывает место", "не меняем весь порядок — заменяем неизвестное", [("неизвестное", ["Кім? Не?", "Қайда? Қайдан?"]), ("да / нет", ["… ба? … бе?", "частица пишется отдельно"]), ("ответ", ["Иә, …", "Жоқ, … емес"])]),
        "kz": ("Сұрақ орын ашады", "бүкіл ретті өзгертпейміз — белгісізді алмастырамыз", [("белгісіз", ["Кім? Не?", "Қайда? Қайдан?"]), ("иә / жоқ", ["… ба? … бе?", "шылау бөлек жазылады"]), ("жауап", ["Иә, …", "Жоқ, … емес"])]),
        "en": ("A question opens a slot", "keep the route — replace the unknown part", [("unknown", ["Кім? Не?", "Қайда? Қайдан?"]), ("yes / no", ["… ба? … бе?", "the particle is separate"]), ("response", ["Иә, …", "Жоқ, … емес"])]),
    }),
    ("07-repair", {
        "ru": ("Не молчим — ремонтируем", "непонятная реплика не завершает разговор", [("останови", ["Кешіріңіз", "Түсінбедім"]), ("измени", ["Қайталаңызшы", "Баяу айтыңызшы"]), ("проверь", ["Бұл нені білдіреді?", "Дұрыс айттым ба?"])]),
        "kz": ("Үндемей қалмаймыз — түзетеміз", "түсініксіз реплика әңгімені аяқтамайды", [("тоқтат", ["Кешіріңіз", "Түсінбедім"]), ("өзгерт", ["Қайталаңызшы", "Баяу айтыңызшы"]), ("тексер", ["Бұл нені білдіреді?", "Дұрыс айттым ба?"])]),
        "en": ("Do not freeze — repair", "an unclear line does not end the conversation", [("pause", ["Кешіріңіз", "Түсінбедім"]), ("change", ["Қайталаңызшы", "Баяу айтыңызшы"]), ("check", ["Бұл нені білдіреді?", "Дұрыс айттым ба?"])]),
    }),
    ("08-mastery", {
        "ru": ("Первое знакомство без сценария", "маршрут известен — слова и порядок выбирает ученик", [("начни", ["приветствие", "имя + роль + место"]), ("поддержи", ["3 разных вопроса", "ответ + встречный вопрос"]), ("восстанови", ["1 фраза ремонта", "заверши разговор"])]),
        "kz": ("Дайын мәтінсіз алғашқы танысу", "бағыт белгілі — сөз бен ретті оқушы таңдайды", [("баста", ["амандасу", "аты + рөлі + орны"]), ("жалғастыр", ["3 түрлі сұрақ", "жауап + қарсы сұрақ"]), ("түзет", ["1 түзету тіркесі", "әңгімені аяқта"])]),
        "en": ("A first meeting without a script", "the route is known — the learner chooses words and order", [("open", ["greeting", "name + role + place"]), ("sustain", ["3 different questions", "answer + return question"]), ("repair", ["1 repair phrase", "close the exchange"])]),
    }),
]


def lines(items, x, y, size=27, gap=43, weight=500):
    spans = []
    for i, item in enumerate(items):
        spans.append(f'<tspan x="{x}" y="{y + i * gap}">{escape(item)}</tspan>')
    return f'<text text-anchor="middle" font-size="{size}" font-weight="{weight}" fill="#24455f">' + "".join(spans) + "</text>"


def render(title, subtitle, cards):
    card_svg = []
    xs = [90, 590, 1090]
    colors = [("#fff4d7", "#e6a83a"), ("#e7f6ff", "#2992bd"), ("#eaf8eb", "#46a260")]
    for index, ((heading, body), x, (fill, stroke)) in enumerate(zip(cards, xs, colors)):
        card_svg.append(f'''<g filter="url(#shadow)">
          <rect x="{x}" y="278" width="420" height="360" rx="24" fill="{fill}" stroke="{stroke}" stroke-width="3"/>
          <rect x="{x + 18}" y="296" width="384" height="70" rx="17" fill="{stroke}" opacity=".16"/>
          <text x="{x + 210}" y="342" text-anchor="middle" font-size="29" font-weight="750" fill="#17364e">{escape(heading)}</text>
          {lines(body, x + 210, 420, 28, 52, 620)}
          <circle cx="{x + 210}" cy="585" r="15" fill="{stroke}" opacity=".78"/>
        </g>''')
        if index < 2:
            card_svg.append(f'''<path d="M{x + 438} 458 H{x + 478}" stroke="#c52b2f" stroke-width="9" stroke-linecap="round"/>
              <path d="M{x + 469} 443 L{x + 486} 458 L{x + 469} 473" fill="none" stroke="#c52b2f" stroke-width="9" stroke-linecap="round" stroke-linejoin="round"/>''')
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}">
      <defs>
        <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#a9cbe0" stroke-width="1" opacity=".35"/></pattern>
        <linearGradient id="paper" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#fbfdff"/><stop offset="1" stop-color="#eef7fb"/></linearGradient>
        <filter id="shadow" x="-15%" y="-15%" width="140%" height="145%"><feDropShadow dx="0" dy="13" stdDeviation="12" flood-color="#17364e" flood-opacity=".18"/></filter>
      </defs>
      <rect width="1600" height="900" rx="36" fill="url(#paper)"/>
      <rect width="1600" height="900" rx="36" fill="url(#grid)"/>
      <path d="M62 222 H1538" stroke="#c52b2f" stroke-width="4" opacity=".76"/>
      <text x="800" y="104" text-anchor="middle" font-family="Inter, Arial, sans-serif" font-size="54" font-weight="800" fill="#17364e">{escape(title)}</text>
      <text x="800" y="166" text-anchor="middle" font-family="Inter, Arial, sans-serif" font-size="25" font-weight="520" fill="#4b6a7e">{escape(subtitle)}</text>
      <g font-family="Inter, Arial, sans-serif">{''.join(card_svg)}</g>
      <g transform="translate(326 724)" filter="url(#shadow)">
        <rect width="948" height="88" rx="22" fill="#17364e"/>
        <text x="474" y="55" text-anchor="middle" font-family="Inter, Arial, sans-serif" font-size="29" font-weight="720" fill="#fff">КӨР → ТЫҢДА → АЙТ → ҚОЛДАН</text>
      </g>
    </svg>'''


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    for stem, localized in MAPS:
        for lang, spec in localized.items():
            (OUT / f"map-{stem}-{lang}.svg").write_text(render(*spec), encoding="utf-8")
    print(f"Generated {len(MAPS) * 3} localized Kazakh support maps")


if __name__ == "__main__":
    main()
