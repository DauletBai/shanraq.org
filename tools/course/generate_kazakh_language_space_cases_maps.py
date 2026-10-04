#!/usr/bin/env python3
"""Generate localized technical support maps for Kazakh lessons 33–40."""
from __future__ import annotations

from html import escape
from pathlib import Path
import math


ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/kazakh-language"


TEXT = {
    "33-case-map": {
        "ru": ("Семь ролей — одна живая сцена", "Сначала смысл и вопрос, затем окончание", ["кто?", "чей?", "что именно?", "кому / куда?", "где?", "откуда?", "с кем / чем?"]),
        "kz": ("Жеті рөл — бір тірі көрініс", "Алдымен мағына мен сұрақ, кейін жалғау", ["кім?", "кімнің?", "нені?", "кімге / қайда?", "қайда?", "қайдан?", "кіммен / немен?"]),
        "en": ("Seven roles — one living scene", "Meaning and question first, ending second", ["who?", "whose?", "what exactly?", "to whom / where?", "where?", "from where?", "with whom / what?"]),
    },
    "34-genitive": {
        "ru": ("Родительный падеж — двусторонний мост", "Оба имени показывают одну связь", ["ЧЬЁ ЦЕЛОЕ?", "ШКОЛЫ", "ЧЬЯ ЧАСТЬ?", "ДВОР"]),
        "kz": ("Ілік септік — екіжақты көпір", "Екі атау бір байланысты көрсетеді", ["НЕНІҢ?", "МЕКТЕПТІҢ", "НЕСІ?", "АУЛАСЫ"]),
        "en": ("Genitive is a two-sided bridge", "Both nouns display one relationship", ["WHOSE WHOLE?", "THE SCHOOL'S", "WHOSE PART?", "ITS YARD"]),
    },
    "35-accusative": {
        "ru": ("Общее действие или конкретный объект?", "Видят ли собеседники одну и ту же книгу?", ["ОБЩЕЕ", "читаю книги", "КОНКРЕТНОЕ", "читаю эту книгу"]),
        "kz": ("Жалпы әрекет пе, нақты нысан ба?", "Әңгімелесушілер бір кітапты көріп тұр ма?", ["ЖАЛПЫ", "кітап оқимын", "НАҚТЫ", "бұл кітапты оқимын"]),
        "en": ("General activity or specific object?", "Can both speakers identify the same book?", ["GENERAL", "кітап оқимын", "SPECIFIC", "бұл кітапты оқимын"]),
    },
    "36-dative": {
        "ru": ("Дательный ведёт к конечной точке", "Получатель, место и назначенное время", ["ПОЛУЧАТЕЛЬ", "мұғалімге", "МЕСТО", "кітапханаға", "ВРЕМЯ", "сағат үшке"]),
        "kz": ("Барыс септік соңғы нүктеге апарады", "Алушы, орын және белгіленген уақыт", ["АЛУШЫ", "мұғалімге", "ОРЫН", "кітапханаға", "УАҚЫТ", "сағат үшке"]),
        "en": ("Dative leads to an endpoint", "Recipient, place, and appointed time", ["RECIPIENT", "мұғалімге", "PLACE", "кітапханаға", "TIME", "сағат үшке"]),
    },
    "37-locative": {
        "ru": ("Местный закрепляет координату", "Движение остановилось: где или когда?", ["МЕСТО", "кітапханада", "ДЕНЬ", "сәрсенбіде", "ЧАС", "сағат онда"]),
        "kz": ("Жатыс септік координатты бекітеді", "Қозғалыс тоқтады: қайда немесе қашан?", ["ОРЫН", "кітапханада", "КҮН", "сәрсенбіде", "САҒАТ", "сағат онда"]),
        "en": ("Locative fixes a coordinate", "Movement has stopped: where or when?", ["PLACE", "кітапханада", "DAY", "сәрсенбіде", "HOUR", "сағат онда"]),
    },
    "38-ablative": {
        "ru": ("Исходный показывает начало", "Место, человек, материал и время смотрят от источника", ["МЕСТО", "мектептен", "ЧЕЛОВЕК", "досымнан", "МАТЕРИАЛ", "ағаштан"]),
        "kz": ("Шығыс септік бастауды көрсетеді", "Орын, адам, материал және уақыт бастауға қарайды", ["ОРЫН", "мектептен", "АДАМ", "досымнан", "МАТЕРИАЛ", "ағаштан"]),
        "en": ("Ablative identifies a source", "Place, person, material, and time look back to a start", ["PLACE", "мектептен", "PERSON", "досымнан", "MATERIAL", "ағаштан"]),
    },
    "39-instrumental": {
        "ru": ("Слитно — спутник или средство; отдельно — «и»", "Пробел меняет грамматическую работу", ["С КЕМ?", "досыммен", "ЧЕМ?", "қарындашпен", "И", "Аружан мен Мария"]),
        "kz": ("Бірге — серік не құрал; бөлек — жалғаулық", "Бос орын грамматикалық қызметті өзгертеді", ["КІММЕН?", "досыммен", "НЕМЕН?", "қарындашпен", "ЖӘНЕ", "Аружан мен Мария"]),
        "en": ("Attached means companion or tool; separate means ‘and’", "A space changes the grammatical job", ["WITH WHOM?", "досыммен", "WITH WHAT?", "қарындашпен", "AND", "Аружан мен Мария"]),
    },
    "40-mastery": {
        "ru": ("Моя среда 1.1: роль меняется — форма отвечает", "Семь вопросов вокруг одного проекта", ["кто?", "чей?", "что?", "кому?", "где?", "откуда?", "с кем?"]),
        "kz": ("Менің ортам 1.1: рөл өзгерсе, форма жауап береді", "Бір жобаның айналасындағы жеті сұрақ", ["кім?", "кімнің?", "нені?", "кімге?", "қайда?", "қайдан?", "кіммен?"]),
        "en": ("My World 1.1: role changes, form responds", "Seven questions around one project", ["who?", "whose?", "what?", "to whom?", "where?", "from where?", "with whom?"]),
    },
}


def wrap(text: str, limit: int) -> list[str]:
    words, rows, current = text.split(), [], []
    for word in words:
        if current and len(" ".join(current + [word])) > limit:
            rows.append(" ".join(current)); current = [word]
        else: current.append(word)
    if current: rows.append(" ".join(current))
    return rows


def text(rows, x, y, size=28, gap=34, weight=700, fill="#17364e", anchor="middle"):
    rows = [rows] if isinstance(rows, str) else rows
    return f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-size="{size}" font-weight="{weight}" fill="{fill}">' + "".join(f'<tspan x="{x}" dy="{0 if i == 0 else gap}">{escape(row)}</tspan>' for i, row in enumerate(rows)) + "</text>"


def shell(title: str, subtitle: str, body: str) -> str:
    tr, sr = wrap(title, 51), wrap(subtitle, 78)
    ty = 70 if len(tr) == 1 else 48
    sy = 145 if len(tr) == 1 else 160
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}"><defs>
    <linearGradient id="paper" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#fbf7eb"/><stop offset=".55" stop-color="#f3faf9"/><stop offset="1" stop-color="#e8f3f7"/></linearGradient>
    <linearGradient id="gold" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#ffd26d"/><stop offset="1" stop-color="#d78b22"/></linearGradient>
    <linearGradient id="blue" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#69c4e6"/><stop offset="1" stop-color="#197da8"/></linearGradient>
    <linearGradient id="green" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#79d398"/><stop offset="1" stop-color="#238050"/></linearGradient>
    <pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#709aac" stroke-opacity=".16" stroke-width="1.5"/></pattern>
    <filter id="shadow" x="-25%" y="-25%" width="160%" height="170%"><feDropShadow dx="0" dy="14" stdDeviation="11" flood-color="#17364e" flood-opacity=".25"/></filter>
    <marker id="redArrow" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto"><path d="M0 0L8 4L0 8Z" fill="#c9343a"/></marker>
    </defs><rect width="1600" height="900" rx="34" fill="url(#paper)"/><rect width="1600" height="900" rx="34" fill="url(#grid)"/><g font-family="Inter,Arial,sans-serif">{text(tr,800,ty,47,51,870)}{text(sr,800,sy,24,31,600,'#4e6c7d')}{body}<g transform="translate(310 790)" filter="url(#shadow)"><rect width="980" height="62" rx="19" fill="#17364e"/>{text('МАҒЫНА → СҰРАҚ → ФОРМА → ТЕКСЕРУ',490,41,24,28,800,'#fff')}</g></g></svg>'''


def bubble(cx, cy, label, detail, color, width=250):
    return f'''<g filter="url(#shadow)"><rect x="{cx-width/2}" y="{cy-58}" width="{width}" height="116" rx="25" fill="#fff" stroke="{color}" stroke-width="4"/><rect x="{cx-width/2+10}" y="{cy-48}" width="{width-20}" height="42" rx="14" fill="{color}" opacity=".18"/>{text(wrap(label,18),cx,cy-20,21,24,850)}{text(wrap(detail,23),cx,cy+28,19,24,620,'#4e6c7d')}</g>'''


def radial(spec, mastery=False):
    title, subtitle, labels = spec
    center = "ЖОБА" if mastery else "КІТАП"
    roles = ["АТАУ", "ІЛІК", "ТАБЫС", "БАРЫС", "ЖАТЫС", "ШЫҒЫС", "КӨМЕКТЕС"]
    colors = ["#d69627", "#2e9cc7", "#3d9b65", "#c14c78", "#7967bd", "#dd7143", "#378c8d"]
    cx, cy, radius = 800, 485, 305
    body = f'<g filter="url(#shadow)"><circle cx="{cx}" cy="{cy}" r="104" fill="url(#gold)" stroke="#fff" stroke-width="9"/>{text(center,cx,cy+11,31,34,900)}</g>'
    for i,(role,detail,color) in enumerate(zip(roles,labels,colors)):
        a=math.radians(-90+i*360/7); x=cx+radius*math.cos(a); y=cy+radius*.72*math.sin(a)
        dx,dy=x-cx,y-cy; distance=math.hypot(dx,dy); ux,uy=dx/distance,dy/distance
        sx,sy=cx+112*ux,cy+112*uy; ex,ey=x-128*ux,y-128*uy
        body += f'<path d="M{sx:.1f} {sy:.1f}L{ex:.1f} {ey:.1f}" stroke="{color}" stroke-width="4.5" stroke-linecap="round" marker-end="url(#redArrow)"/>'
        body += bubble(x,y,role,detail,color,230)
    return shell(title,subtitle,body)


def bridge(spec):
    title, subtitle, a = spec
    body = bubble(360,465,a[0],a[1],"#d69627",390)+bubble(1240,465,a[2],a[3],"#3d9b65",390)
    body += '<path d="M565 465C700 330 900 330 1035 465" fill="none" stroke="#2e9cc7" stroke-width="18" stroke-linecap="round"/><path d="M565 475C700 610 900 610 1035 475" fill="none" stroke="#2e9cc7" stroke-width="8" stroke-linecap="round"/>'
    body += '<rect x="560" y="445" width="480" height="76" rx="18" fill="#fff" opacity=".94"/>'+text("мектептің  ═══  ауласы",800,490,32,36,850)
    return shell(title,subtitle,body)


def split(spec):
    title, subtitle, a=spec
    body='<g filter="url(#shadow)"><circle cx="800" cy="325" r="70" fill="url(#gold)" stroke="#fff" stroke-width="8"/>'+text("КІТАП",800,336,27,30,850)+'</g>'
    body+='<path d="M770 390L430 520" stroke="#2e9cc7" stroke-width="6" marker-end="url(#redArrow)"/><path d="M830 390L1170 520" stroke="#c9343a" stroke-width="6" marker-end="url(#redArrow)"/>'
    body+=bubble(390,585,a[0],a[1],"#2e9cc7",480)+bubble(1210,585,a[2],a[3],"#c9343a",480)
    return shell(title,subtitle,body)


def three_targets(spec, reverse=False, coordinate=False):
    title,subtitle,a=spec; colors=("#d69627","#2e9cc7","#3d9b65"); xs=(310,800,1290); body=""
    if coordinate:
        body+='<g filter="url(#shadow)"><rect x="120" y="275" width="1360" height="400" rx="28" fill="#fff" stroke="#6ba6bf" stroke-width="4"/><path d="M160 610H1440M240 650V310" stroke="#17364e" stroke-width="5"/><path d="M240 485H1440" stroke="#8fb7c8" stroke-width="3" stroke-dasharray="14 12"/></g>'
    for i,x in enumerate(xs):
        label,detail=a[i*2:i*2+2]
        body+=bubble(x,500,label,detail,colors[i],350)
        if reverse:
            body+=f'<path d="M{x} 625L800 710" stroke="#c9343a" stroke-width="5" marker-end="url(#redArrow)"/>'
        elif not coordinate:
            body+=f'<path d="M800 700L{x} 625" stroke="#c9343a" stroke-width="5" marker-end="url(#redArrow)"/>'
    if not coordinate:
        cap="БАСТАУ" if not reverse else "НӘТИЖЕ"
        body+=f'<g filter="url(#shadow)"><rect x="660" y="680" width="280" height="70" rx="22" fill="#17364e"/>{text(cap,800,725,24,28,850,"#fff")}</g>'
    return shell(title,subtitle,body)


def render(stem, spec):
    if stem in {"33-case-map","40-mastery"}: return radial(spec,stem=="40-mastery")
    if stem=="34-genitive": return bridge(spec)
    if stem=="35-accusative": return split(spec)
    if stem=="37-locative": return three_targets(spec,coordinate=True)
    if stem=="38-ablative": return three_targets(spec,reverse=True)
    return three_targets(spec)


def main():
    OUT.mkdir(parents=True,exist_ok=True)
    for stem,localized in TEXT.items():
        for lang,spec in localized.items():
            (OUT/f"map-{stem}-{lang}.svg").write_text(render(stem,spec),encoding="utf-8")
    print(f"Generated {len(TEXT)*3} localized Space and Cases maps")


if __name__=="__main__": main()
