#!/usr/bin/env python3
"""Render the authored chemistry block and its large, depth-styled SVG maps."""
from html import escape
from pathlib import Path
from textwrap import wrap

try:
    from .chemistry_data import LESSONS
except ImportError:
    from chemistry_data import LESSONS

ROOT = Path(__file__).resolve().parents[2]
LESSON_DIR = ROOT / "course/lessons/chemistry"
MAP_DIR = ROOT / "web/static/course/chemistry"
LANGS = ("ru", "kz", "en")
HEAD = {
    "ru": ("Опорная карта", "Ситуация и вопрос", "Новые слова без пропусков", "Предскажите", "Исследуйте или разберите данные", "Разобранный пример", "Где легко ошибиться", "Перенос в новый случай", "Ответ и проверка памяти", "Таблица исследователя", "Когда вернуться"),
    "kz": ("Тірек карта", "Жағдай және сұрақ", "Жаңа сөздер түсінікті тілмен", "Алдын ала болжаңыз", "Зерттеңіз не берілген деректі талдаңыз", "Талданған мысал", "Жиі кездесетін қате", "Жаңа жағдайға көшіру", "Жауап және есте сақтау", "Зерттеушінің кестесі", "Қашан қайта қарау"),
    "en": ("Support map", "Situation and question", "New terms, fully explained", "Predict first", "Investigate or analyse given data", "Worked example", "A common mistake", "Transfer to a new case", "Answer and recall", "Investigator's table", "When to revisit"),
}
LABELS = {
    "ru": [("ВОПРОС", "Что произошло?", "НАБЛЮДЕНИЕ", "Крупинок не видно", "МОДЕЛЬ", "Соль могла остаться", "ПРОВЕРКА", "Без вкуса и нагрева"),
           ("ПРЕДМЕТ", "Стакан", "МАТЕРИАЛ", "Стекло", "СВОЙСТВО", "Прозрачность", "ПРОВЕРКА", "Признак ≠ состав"),
           ("ВИДНО", "Один цвет", "СОСТАВ", "Один или несколько?", "СМЕСЬ", "Доли меняются", "ВЫВОД", "Вид ≠ чистота"),
           ("ВИДНО", "Кристалл исчез", "МОДЕЛЬ", "Ионы среди воды", "РАСТВОР", "Соль + вода", "ГРАНИЦА", "Шары не в масштабе"),
           ("НАЧАЛО", "Песок + соль + вода", "СУЗИТЬ", "Осадить и отфильтровать", "ФИЛЬТРАТ", "Вода + соль", "ГРАНИЦА", "Не пить")],
    "kz": [("СҰРАҚ", "Не болды?", "БАҚЫЛАУ", "Түйір көрінбейді", "МОДЕЛЬ", "Тұз қалуы мүмкін", "ТЕКСЕРУ", "Дәмсіз, қыздырмай"),
           ("БҰЙЫМ", "Стақан", "МАТЕРИАЛ", "Шыны", "ҚАСИЕТ", "Мөлдірлік", "ТЕКСЕРУ", "Белгі ≠ құрам"),
           ("КӨРІНЕДІ", "Бір түс", "ҚҰРАМ", "Бір ме, бірнеше ме?", "ҚОСПА", "Үлесі өзгереді", "ҚОРЫТЫНДЫ", "Түр ≠ тазалық"),
           ("КӨРІНЕДІ", "Кристалл жоқ", "МОДЕЛЬ", "Судағы иондар", "ЕРІТІНДІ", "Тұз + су", "ШЕК", "Шарлар масштабта емес"),
           ("БАСТАУ", "Құм + тұз + су", "БӨЛУ", "Тұндыру, сүзу", "СҮЗІНДІ", "Су + тұз", "ШЕК", "Ішпеңіз")],
    "en": [("QUESTION", "What happened?", "OBSERVE", "No visible grains", "MODEL", "Salt may remain", "TEST", "No taste or heat"),
           ("OBJECT", "Drinking glass", "MATERIAL", "Glass", "PROPERTY", "Transparency", "CHECK", "Feature ≠ composition"),
           ("APPEARANCE", "One colour", "COMPOSITION", "One kind or several?", "MIXTURE", "Ratios can change", "CLAIM", "Looks ≠ purity"),
           ("OBSERVE", "Crystal not seen", "MODEL", "Ions in water", "SOLUTION", "Salt + water", "LIMIT", "Spheres not to scale"),
           ("START", "Sand + salt + water", "SEPARATE", "Settle and filter", "FILTRATE", "Water + salt", "LIMIT", "Do not drink")]
}
TABLE_HEAD = {
    "ru": "| Шаг | Что записать | Чего это пока не доказывает |\n|---|---|---|\n| До | Предсказание и одинаковые условия | Что предсказание верно |\n| Во время | Только видимое или измеренное | Причину без проверки |\n| После | Вывод с границей и новым вопросом | Полную чистоту и безопасность |",
    "kz": "| Кезең | Не жазу керек | Бұл әлі нені дәлелдемейді |\n|---|---|---|\n| Дейін | Болжам мен бірдей жағдай | Болжамның дұрыстығын |\n| Кезінде | Көрген не өлшеген нәрсе | Тексерусіз себепті |\n| Кейін | Шегі көрсетілген қорытынды, жаңа сұрақ | Толық тазалық пен қауіпсіздікті |",
    "en": "| Stage | Record | This does not yet prove |\n|---|---|---|\n| Before | Prediction and matched conditions | That the prediction is right |\n| During | Only what was seen or measured | A cause without a test |\n| After | A claim with a limit and a new question | Complete purity or safety |"
}
RETURN = {
    "ru": "Закройте страницу и через 10 минут воспроизведите карту, определите термины своими словами и решите новый случай. Через сутки и через неделю повторите без подсказки. Если ошибка осталась, вернитесь к разобранному примеру и запишите, какой шаг был пропущен.",
    "kz": "Бетті жауып, 10 минуттан кейін картаны есіңізден салыңыз, ұғымдарды өз сөзіңізбен түсіндіріңіз және жаңа жағдайды шешіңіз. Бір күн және бір аптадан соң кеңессіз қайталаңыз. Қате қалса, талданған мысалға оралып, қай қадамды өткізіп алғаныңызды жазыңыз.",
    "en": "Close the page and, after 10 minutes, reconstruct the map, define the terms in your own words and solve the new case. Repeat without hints after a day and a week. If an error persists, return to the worked example and write down the missing step."
}
VISUAL = {
    "ru": ("Исследование воды и соли: фотографическая реконструкция учебной ситуации.", "Стакан, металлическая ложка и пластиковая бутылка: сравните предмет, материал и свойство.", "Образцы смеси и прозрачной воды: видимый облик не устанавливает химический состав.", "Увеличенная 3D-модель растворения: цвета и размеры частиц условны; это не микрофотография.", "Фильтрация песка: прозрачный фильтрат всё ещё может содержать растворённую соль."),
    "kz": ("Су мен тұзды зерттеу: оқу жағдайының фотосуретке ұқсас көрінісі.", "Стақан, металл қасық және пластик бөтелке: бұйым, материал, қасиетті салыстырыңыз.", "Қоспа мен мөлдір су үлгілері: сыртқы көрініс химиялық құрамды анықтамайды.", "Ерудің үлкейтілген 3D-моделі: бөлшектердің түсі мен өлшемі шартты; бұл микрофотосурет емес.", "Құмды сүзу: мөлдір сүзіндіде әлі еріген тұз болуы мүмкін."),
    "en": ("Investigating water and salt: a photorealistic reconstruction of a lesson scene.", "A glass, metal spoon and plastic bottle: compare object, material and property.", "Mixture and clear-water samples: appearance does not establish composition.", "Magnified 3D model of dissolving: particle colours and sizes are symbolic, not a micrograph.", "Filtering sand: clear filtrate may still contain dissolved salt.")
}
VISUAL_FILES = ('/static/covers/school/chemistry/first-block/01-investigate.webp',
                '/static/course/chemistry/visuals/02-materials.webp',
                '/static/course/chemistry/visuals/03-mixtures.webp',
                '/static/course/chemistry/visuals/04-dissolving.webp',
                '/static/course/chemistry/visuals/05-filter.webp')
BRIDGES = {
    'ru': (
        '**Масса** — сколько вещества содержится в образце; её сравнивают взвешиванием и здесь записывают в граммах (г). Закрытый сосуд при взвешивании включает и посуду, и содержимое. Если открывали его, нужно учесть всё внесённое и унесённое. Впишите в журнал две колонки: «я действительно заметил» и «я пока предполагаю».',
        '**Материя** — всё, что занимает место и имеет массу. **Состав** отвечает на вопрос, какие вещества есть в образце и сколько их. **Растворимость** — способность вещества распределяться в выбранном растворителе при оговорённых условиях; условия важны, поскольку разные вещества ведут себя по-разному. Не определяйте материал только по одному свойству.',
        '**H₂O** — короткая химическая запись воды; её подробный смысл разберём после введения атомов. **Дистиллированная вода** — вода, полученная после испарения и последующего превращения пара обратно в жидкость; это уменьшает многие примеси, но название на этикетке само по себе не доказывает абсолютную чистоту конкретной пробы. **Фаза** — область образца, в которой состояние и состав примерно одинаковы.',
        '**Атом** — маленькая частица химического элемента; элемент и строение атома изучим дальше. **Молекула** — группа связанных атомов; вода состоит из молекул. **Ион** — частица с электрическим зарядом; у поваренной соли есть положительные и отрицательные ионы. Когда соль растворяется, эти ионы расходятся в воде. Шарики на рисунке обозначают частицы условно: их реальные размеры, расстояния и движение здесь не показаны.',
        '**Пора** — маленькое отверстие в фильтре. Бумага **пористая**, потому что между её волокнами есть проходы для жидкости. Размер частиц и пор определяет, что задержится; растворённая соль обычно проходит через простой бумажный фильтр. **Микроорганизмы** — крошечные живые организмы, невидимые без специальных приборов; прозрачный фильтрат от этого не становится питьевой водой.'),
    'kz': (
        '**Масса** — үлгідегі заттың мөлшерін сипаттайтын шама; оны таразыда салыстырып, мұнда граммен (г) жазамыз. Жабық ыдысты өлшегенде ыдыс та, ішіндегісі де есепке кіреді. Ашылса, қосылған не шыққанның бәрін ескеру керек. Күнделікте «шын байқадым» және «әзірше болжаймын» деген екі баған жасаңыз.',
        '**Материя** — орын алатын әрі массасы бар нәрсенің бәрі. **Құрам** үлгіде қандай заттар және қанша мөлшерде бар деген сұраққа жауап береді. **Ерігіштік** — заттың белгілі жағдайдағы таңдалған еріткішке тарала алу қасиеті; әр зат әртүрлі әрекет етеді. Материалды жалғыз қасиетпен анықтамаңыз.',
        '**H₂O** — судың қысқа химиялық жазылуы; атомдарды таныған соң оның мәнін толық талдаймыз. **Дистилденген су** — суды буландырып, буды қайта сұйыққа айналдырғанда алынған су; бұл көптеген қоспаларды азайтады, бірақ жазудың өзі нақты үлгінің мүлдем таза екенін дәлелдемейді. **Фаза** — үлгінің күйі мен құрамы шамамен бірдей аймағы.',
        '**Атом** — химиялық элементтің кішкене бөлшегі; элемент пен атом құрылысын кейін үйренеміз. **Молекула** — байланысқан атомдар тобы; су молекулалардан тұрады. **Ион** — электр заряды бар бөлшек; ас тұзында оң және теріс иондар бар. Тұз ерігенде олар суда таралады. Суреттегі шарлар шартты белгі ғана: нақты өлшем, арақашықтық және қозғалыс көрсетілмейді.',
        '**Саңылау** — сүзгідегі кішкентай тесік. Қағаз **кеуекті**, өйткені талшықтарының арасында сұйық өтетін жолдар бар. Бөлшек пен саңылау өлшемі не бөгелетінін анықтайды; еріген тұз қарапайым қағаз сүзгіден әдетте өтеді. **Микроағзалар** — арнайы құралсыз көрінбейтін өте ұсақ тіршілік иелері; мөлдір сүзінді ішуге жарамды суға айналмайды.'),
    'en': (
        '**Mass** measures how much matter a sample contains; here compare it by weighing and record it in grams (g). Weighing a closed vessel includes both vessel and contents. If it was opened, account for everything added or removed. Make two notebook columns: "I actually observed" and "I currently infer".',
        '**Matter** is anything that occupies space and has mass. **Composition** tells which substances are present and how much of each. **Solubility** describes whether a substance can distribute in a chosen solvent under stated conditions; different substances behave differently. One property alone cannot identify a material.',
        '**H₂O** is shorthand for water; we will explain the formula in full after introducing atoms. **Distilled water** is water obtained by evaporating water and turning its vapour back into liquid; this reduces many impurities, but a label alone does not prove any sample is absolutely pure. A **phase** is a region of roughly uniform state and composition.',
        'An **atom** is a small unit of a chemical element; elements and atomic structure come later. A **molecule** is a group of bonded atoms; water consists of molecules. An **ion** is an electrically charged particle; table salt contains positive and negative ions. When salt dissolves, those ions spread through water. Spheres in the image are symbols: actual sizes, distances and movement are not shown.',
        'A **pore** is a tiny hole in a filter. Paper is **porous** because liquid can pass between its fibres. Particle and pore size determine what is retained; dissolved salt usually passes through ordinary filter paper. **Microorganisms** are tiny living things unseen without special equipment; clear filtrate is not thereby safe to drink.')
}


def map_svg(number, lang):
    labels = LABELS[lang][number - 1]
    title = escape(LESSONS[number - 1][lang][0])
    chunks = [f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1280 490" role="img" aria-labelledby="title desc">
<title id="title">{title}</title><desc id="desc">Four-stage teaching model, enlarged and not to scale. Lesson {number}.</desc>
<defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#0b1e35"/><stop offset="1" stop-color="#143b58"/></linearGradient>
<linearGradient id="face" x2=".9" y2="1"><stop stop-color="#fffdf1"/><stop offset="1" stop-color="#d9eaf1"/></linearGradient>
<radialGradient id="orb"><stop stop-color="#fff9c6"/><stop offset=".55" stop-color="#6bd4e3"/><stop offset="1" stop-color="#147d9b"/></radialGradient>
<filter id="shadow" x="-30%" y="-30%" width="160%" height="180%"><feDropShadow dx="8" dy="12" stdDeviation="9" flood-opacity=".42"/></filter></defs>
<rect width="1280" height="490" rx="28" fill="url(#bg)"/>
<text x="56" y="65" fill="#fff" font-family="Arial,sans-serif" font-size="31" font-weight="700">{title}</text>
<text x="58" y="104" fill="#b5d4e4" font-family="Arial,sans-serif" font-size="18">{escape({'ru':'Схема модели · объекты и частицы не в масштабе','kz':'Модель сызбасы · заттар мен бөлшектер масштабта емес','en':'Model diagram · objects and particles not to scale'}[lang])}</text>''']
    colors = ("#00a7bd", "#eba85c", "#a16ae3", "#56bb83")
    for i in range(4):
        x = 44 + i * 306
        top, bottom = labels[2*i:2*i+2]
        chunks.append(f'''<g filter="url(#shadow)"><path d="M{x+9} 155 L{x+254} 155 L{x+273} 175 L{x+28} 175 Z" fill="{colors[i]}" opacity=".65"/>
<rect x="{x}" y="168" width="266" height="244" rx="20" fill="url(#face)" stroke="{colors[i]}" stroke-width="4"/></g>
<circle cx="{x+43}" cy="210" r="25" fill="url(#orb)" stroke="#fff" stroke-width="2"/>
<text x="{x+22}" y="219" fill="#073047" font-family="Arial,sans-serif" font-size="25" font-weight="700">{i+1}</text>
<text x="{x+20}" y="287" fill="#26405a" font-family="Arial,sans-serif" font-size="20" font-weight="700">{escape(top)}</text>
<circle cx="{x+204}" cy="210" r="13" fill="url(#orb)"/><circle cx="{x+225}" cy="225" r="8" fill="url(#orb)"/>
''')
        # Preserve natural glyph proportions. Wrap long labels within the card.
        lines = wrap(bottom, width=20, break_long_words=False, break_on_hyphens=False)
        for line_number, line in enumerate(lines):
            chunks.append(f'<text x="{x+20}" y="{329+line_number*30}" fill="#102d42" '
                          f'font-family="Arial,sans-serif" font-size="19">{escape(line)}</text>')
        if i < 3:
            chunks.append(f'<path d="M{x+271} 289 h28 m-11 -12 12 12 -12 12" stroke="#d7f5e8" stroke-width="4" fill="none"/>')
    chunks.append('</svg>')
    return '\n'.join(chunks) + '\n'


def render(number, lang):
    title, summary, *sections = LESSONS[number - 1][lang]
    headings = HEAD[lang]
    image = f'/static/course/chemistry/map-{number:02d}-{lang}.svg'
    body = [f'# {title}', f'_Лид (summary):_ **{summary}**',
            f'## {headings[0]}', f'![{title}]({image})',
            f'![{VISUAL[lang][number-1]}]({VISUAL_FILES[number-1]})', f'*{VISUAL[lang][number-1]}*']
    for section_number, (heading, paragraph) in enumerate(zip(headings[1:9], sections)):
        body.extend([f'## {heading}', paragraph])
        if section_number == 1:
            body.append(BRIDGES[lang][number-1])
    body.extend([f'## {headings[9]}', TABLE_HEAD[lang], f'## {headings[10]}', RETURN[lang]])
    if number > 1:
        previous = LESSONS[number - 2][lang][0]
        body.append(f'[{previous}](/read/chemistry-{number-1:02d}-{LESSONS[number-2]["stem"].split("-",1)[1]}?lang={lang})')
    if number < len(LESSONS):
        following = LESSONS[number][lang][0]
        body.append(f'[{following}](/read/chemistry-{number+1:02d}-{LESSONS[number]["stem"].split("-",1)[1]}?lang={lang})')
    return '\n\n'.join(body) + '\n'


def main():
    LESSON_DIR.mkdir(parents=True, exist_ok=True)
    MAP_DIR.mkdir(parents=True, exist_ok=True)
    for number, lesson in enumerate(LESSONS, 1):
        for lang in LANGS:
            suffix = '' if lang == 'ru' else f'-{lang}'
            (LESSON_DIR / f'{lesson["stem"]}{suffix}.md').write_text(render(number, lang), encoding='utf-8')
            (MAP_DIR / f'map-{number:02d}-{lang}.svg').write_text(map_svg(number, lang), encoding='utf-8')


if __name__ == '__main__':
    main()
