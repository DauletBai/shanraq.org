#!/usr/bin/env python3
"""Build hand-authored trilingual wave, sound and light lessons and diagrams."""
from html import escape
from pathlib import Path
import math
import textwrap

try:
    from .physics_waves_data import WAVES
except ImportError:
    from physics_waves_data import WAVES

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/physics"
MAPS = ROOT / "web/static/course/physics"

HEAD = {
    "ru": ("Вопрос и прогноз", "Слова без пропусков", "Разбираем учебный пример", "Проверяем безопасно", "Типичная ошибка", "Опорная схема", "Применяем в другой ситуации", "Проверь себя", "Страница лабораторного журнала"),
    "kz": ("Сұрақ және болжам", "Жаңа сөздерді түсінейік", "Оқу мысалын талдаймыз", "Қауіпсіз тексереміз", "Жиі кездесетін қате", "Тірек сызбасы", "Басқа жағдайға қолданамыз", "Өзіңді тексер", "Зертхана журналының беті"),
    "en": ("Question and prediction", "New words without gaps", "Work through supplied data", "Test safely", "A common mistake", "Support diagram", "Transfer to another case", "Check your understanding", "Lab notebook page"),
}
CHECK = {
    "ru": "Закройте текст и объясните новому ученику рисунок слева направо. Назовите наблюдение, учебное допущение и то, чего эти данные не доказывают. Затем решите вопрос переноса без подсказки и проверьте единицы или направление на схеме. Если ответ расходится, вернитесь к первому непонятному слову, а не заучивайте готовую фразу. Устный ответ, чертёж и измерение — разные виды доказательства; подпишите каждый честно.",
    "kz": "Мәтінді жауып, суретті жаңа оқушыға солдан оңға қарай түсіндіріңіз. Бақылауды, оқу жорамалын және бұл дерек әлі дәлелдемейтін нәрсені атаңыз. Содан кейін жаңа жағдайдың сұрағын көмексіз шешіп, бірлікті не бағытты сызбамен тексеріңіз. Жауап сәйкес келмесе, дайын сөйлемді жаттамай, алғашқы түсініксіз сөзге оралыңыз. Ауызша жауап, сызба және өлшеу — дәлелдің әртүрлі түрі; әрқайсын адал белгілеңіз.",
    "en": "Close the text and explain the left-to-right drawing to a new learner. Name an observation, a teaching assumption and something the evidence does not prove. Then answer the transfer question without a hint and check units or direction against the diagram. If you disagree, return to the first unfamiliar word rather than memorising a sentence. A spoken explanation, a drawing and a measurement are different kinds of evidence: label each honestly.",
}
JOURNAL = {
    "ru": "Добавьте к «Лабораторному журналу моего мира» строку: вопрос → прогноз → что оставили постоянным → что изменили → наблюдение или расчёт → вывод → ограничение. В рисунке поставьте масштаб или укажите «схематично». Если опыт не проводился, запишите «учебный расчёт», а не «мы измерили». Вернитесь к этой записи после следующего урока и отметьте, что стало понятнее.",
    "kz": "«Менің әлемімнің зертхана журналына» жол қосыңыз: сұрақ → болжам → нені тұрақты қалдырдық → нені өзгерттік → бақылау не есеп → қорытынды → шектеу. Суретте масштабты көрсетіңіз немесе «схемалық» деп белгілеңіз. Тәжірибе жасалмаса, «біз өлшедік» демей, «оқу есебі» деп жазыңыз. Келесі сабақтан соң бұл жазбаға оралып, ненің анық болғанын белгілеңіз.",
    "en": "Add a line to My World Lab Notebook: question → prediction → controlled conditions → changed condition → observation or calculation → conclusion → limit. Put a scale on a drawing or mark it 'schematic'. If you ran no experiment, write 'supplied calculation', not 'we measured'. Revisit the entry after the next lesson and note what became clearer.",
}


def txt(x, y, value, size=29, anchor="middle", colour="#16445e"):
    return (f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" '
            f'font-size="{size}" font-weight="700" fill="{colour}">{escape(value)}</text>')


def path(d, colour="#25718d", width=7, dash=""):
    extra = f' stroke-dasharray="{dash}"' if dash else ""
    return f'<path d="{d}" fill="none" stroke="{colour}" stroke-width="{width}" stroke-linecap="round" stroke-linejoin="round"{extra}/>'


def scene(kind, lang):
    p = []
    if kind == "wave":
        # One fixed source, not three circles whose centres travel with time.
        # The rings are at radii 2 and 4 cm using 40 px per teaching cm.
        p += [path("M170 470H870", "#a3c2cb", 3),
              '<circle cx="440" cy="470" r="160" fill="none" stroke="#a7d3df" stroke-width="6"/>',
              '<circle cx="440" cy="470" r="80" fill="none" stroke="#499ab8" stroke-width="8"/>',
              '<circle cx="440" cy="470" r="12" fill="#f8c87c" stroke="#b77f50" stroke-width="4"/>',
              '<circle cx="520" cy="470" r="12" fill="#e8a269" stroke="#ad714b" stroke-width="3"/>',
              path("M520 423V392m-10 10 10-10 10 10", "#d9835b", 4),
              txt(345, 480, "t = 0"), txt(440, 375, "t = 1"), txt(440, 290, "t = 2"),
              path("M440 520H520", "#468aa4", 4), path("M440 570H600", "#92bbca", 4),
              txt(480, 554, "2 cm", 25), txt(520, 604, "4 cm", 25)]
    elif kind == "period":
        pts = " ".join(f"{x},{int(475-110*math.sin(2*math.pi*(x-240)/240))}" for x in range(240, 841, 8))
        p += [f'<polyline points="{pts}" fill="none" stroke="#2784a4" stroke-width="8"/>',
              path("M220 475H870", "#849fab", 3),
              *[path(f"M{x} 335V640", "#d39c67", 3, "8 9") for x in (300, 540, 780)],
              txt(300, 680, "0 m"), txt(540, 680, "1.5 m"), txt(780, 680, "3 m"),
              path("M300 320H540", "#d98754", 5)]
    elif kind == "sound":
        p += ['<rect x="175" y="320" width="42" height="260" rx="9" fill="#d58d62"/>',
              path("M220 450H860", "#a5c8d2", 3)]
        for cx in (340, 560, 780):
            for dx in (-22, -8, 6, 20):
                for dy in (-35, 0, 35):
                    p.append(f'<circle cx="{cx+dx}" cy="{450+dy}" r="7" fill="#3489a3"/>')
        for x in (450, 670):
            for dy in (-27, 27):
                p.append(f'<circle cx="{x}" cy="{450+dy}" r="7" fill="#3489a3"/>')
        p += [path("M395 315H505m-18-15 18 15-18 15", "#db8d5c", 5),
              path("M615 315H725m-18-15 18 15-18 15", "#db8d5c", 5)]
    elif kind == "pitch":
        for y, cycles in ((365, 2), (600, 4)):
            pts = " ".join(f"{x},{int(y-58*math.sin(2*math.pi*cycles*(x-230)/630))}" for x in range(230, 861, 7))
            p += [path(f"M215 {y}H875", "#9cb9c6", 3),
                  f'<polyline points="{pts}" fill="none" stroke="#368ca8" stroke-width="7"/>']
        p += [txt(176, 372, "A"), txt(176, 607, "B")]
    elif kind == "reflection":
        p += [path("M170 550H865", "#658796", 10), path("M470 260V670", "#90bbc9", 4, "10 10"),
              path("M300 380L470 550L640 380", "#db8555", 8),
              txt(335, 490, "45°", 31), txt(608, 490, "45°", 31),
              '<circle cx="470" cy="550" r="10" fill="#1d647c"/>']
    elif kind == "refraction":
        p += ['<rect x="170" y="500" width="695" height="230" fill="#cce8f2" opacity=".8"/>',
              path("M170 500H865", "#5e9bb1", 7), path("M500 260V735", "#91b9c8", 4, "9 10"),
              path("M320 320L500 500L650 750", "#d98857", 8),
              txt(392, 420, "45°", 31), txt(582, 620, "≈32°", 31)]
    elif kind == "lens":
        p += [path("M170 500H865", "#8bb0bf", 4),
              '<ellipse cx="500" cy="500" rx="23" ry="230" fill="#cbeef3" stroke="#69a2b5" stroke-width="6"/>',
              path("M260 500V380m-13 18 13-18 13 18", "#d48358", 7),
              path("M740 500V620m-13-18 13 18 13-18", "#d48358", 7),
              path("M260 380H500L740 620", "#277995", 6),
              path("M260 380L740 620", "#e2a05d", 6),
              '<circle cx="380" cy="500" r="7" fill="#476f80"/><circle cx="620" cy="500" r="7" fill="#476f80"/>',
              txt(380, 540, "F", 26), txt(620, 540, "F", 26),
              txt(260, 695, "20 cm", 27), txt(740, 695, "20 cm", 27)]
    elif kind == "capstone":
        # Incoming rays have the same 30° direction. Turning the mirror 10°
        # changes the reflected direction by 20° in the plane-mirror model.
        p += [path("M175 580H425", "#557f91", 9), path("M555 602L805 558", "#557f91", 9),
              path("M160 337L300 580L440 337", "#df8a54", 7),
              path("M540 337L680 580L723 337", "#df8a54", 7),
              path("M300 305V650", "#9cbac6", 3, "8 9"),
              path("M633 313L708 740", "#9cbac6", 3, "8 9"),
              txt(300, 690, "30° → 30°", 28), txt(680, 690, "10° → 20°", 28),
              path("M460 445H520m-15-13 15 13-15 13", "#36859e", 5)]
    return p


def svg(row, lang):
    title = row[lang][0]
    labels = row["map"][lang]
    title_lines = textwrap.wrap(title, 53, break_long_words=False)
    if len(title_lines) > 2:
        raise ValueError((row["stem"], lang, "title too long"))
    out = ['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
           f'<title id="title">{escape(title)}</title><desc id="desc">{escape("; ".join(labels))}</desc>',
           '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#d4e8ed"/></pattern><linearGradient id="paper" x2="1" y2="1"><stop stop-color="#fff"/><stop offset="1" stop-color="#e8f3f7"/></linearGradient><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="10" stdDeviation="8" flood-color="#37657b" flood-opacity=".2"/></filter></defs>',
           '<rect width="1600" height="900" fill="#f2f9fb"/><rect width="1600" height="900" fill="url(#grid)"/>',
           '<rect x="70" y="195" width="875" height="590" rx="16" fill="url(#paper)" stroke="#99bac8" stroke-width="3" filter="url(#shadow)"/>',
           '<rect x="975" y="195" width="555" height="590" rx="16" fill="#fff" stroke="#99bac8" stroke-width="3" filter="url(#shadow)"/>']
    for i, line in enumerate(title_lines):
        out.append(txt(75, 76+i*50, line, 35, "start"))
    out.extend(scene(row["kind"], lang))
    for i, label in enumerate(labels):
        y = 220+i*138
        out.append(f'<rect x="1005" y="{y}" width="495" height="110" rx="10" fill="{("#e6f5fa", "#fff0e4", "#eaf4e9", "#edf0fa")[i]}" stroke="#d0e2e8" stroke-width="2"/>')
        lines = textwrap.wrap(label, 28, break_long_words=False)
        if len(lines) > 2:
            raise ValueError((row["stem"], lang, label))
        for j, line in enumerate(lines):
            out.append(txt(1027, y+46+j*32, line, 27, "start"))
    out.append('</svg>')
    return '\n'.join(out)+'\n'


def markdown(row, number, lang):
    title, summary, hook, terms, worked, activity, trap, transfer = row[lang]
    contents = [hook, terms, worked, activity, trap,
                f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)',
                transfer, CHECK[lang], JOURNAL[lang]]
    body = '\n\n'.join(f'## {head}\n\n{value}' for head, value in zip(HEAD[lang], contents))
    if number < 74:
        next_row = WAVES[number-66]
        label = {"ru": "Следующий урок", "kz": "Келесі сабақ", "en": "Next lesson"}[lang]
        body += f'\n\n[{label}: {next_row[lang][0]}](/read/physics-{next_row["stem"]}?lang={lang})'
    else:
        label = {"ru": "Вернуться к курсу «Физика»", "kz": "«Физика» курсына оралу", "en": "Return to the Physics course"}[lang]
        body += f'\n\n[{label}](/course/physics?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'


def build():
    if len(WAVES) != 8:
        raise ValueError("waves block requires eight lessons")
    for number, row in enumerate(WAVES, 67):
        assert row["stem"].startswith(f"{number}-"), row["stem"]
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            (LESSONS / f'{row["stem"]}{suffix}.md').write_text(markdown(row, number, lang), encoding="utf-8")
            (MAPS / f'map-{row["stem"]}-{lang}.svg').write_text(svg(row, lang), encoding="utf-8")


if __name__ == "__main__":
    build()
