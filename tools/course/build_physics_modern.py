#!/usr/bin/env python3
"""Render the manually authored modern-physics block in three languages."""

from html import escape
from pathlib import Path
import textwrap

try:
    from .physics_modern_data import MODERN
    from .build_physics_waves import HEAD, CHECK, JOURNAL
except ImportError:
    from physics_modern_data import MODERN
    from build_physics_waves import HEAD, CHECK, JOURNAL

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/physics"
MAPS = ROOT / "web/static/course/physics"


def t(x, y, value, size=26, anchor="middle", colour="#153d59"):
    return (f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" '
            f'font-size="{size}" font-weight="700" fill="{colour}">{escape(str(value))}</text>')


def path(d, colour="#367d9e", width=5, dash=""):
    extra = f' stroke-dasharray="{dash}"' if dash else ""
    return f'<path d="{d}" fill="none" stroke="{colour}" stroke-width="{width}" stroke-linecap="round" stroke-linejoin="round"{extra}/>'


def box(x, y, w, h, fill="#e9f6fa"):
    return f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="13" fill="{fill}" stroke="#88aec0" stroke-width="3" filter="url(#shadow)"/>'


def scene(row, lang):
    kind = row["kind"]
    out = []
    if kind == "atom":
        # The translucent cloud is a probability cue, expressly not an orbit.
        out += ['<ellipse cx="515" cy="480" rx="280" ry="198" fill="url(#cloud)" opacity=".68"/>',
                '<circle cx="515" cy="480" r="78" fill="url(#nucleus)" stroke="#a75f65" stroke-width="4"/>',
                t(515, 490, "6p + 6n", 29),
                '<circle cx="315" cy="424" r="17" fill="#4c83b3"/>', t(315, 433, "−", 25, colour="#fff"),
                '<circle cx="694" cy="536" r="17" fill="#4c83b3"/>', t(694, 545, "−", 25, colour="#fff")]
        for x,y in ((610,340),(510,305),(365,585),(675,630)):
            out += [f'<circle cx="{x}" cy="{y}" r="17" fill="#4c83b3"/>', t(x,y+9,"−",25,colour="#fff")]
        label={"ru":"Облако вероятности, не орбита","kz":"Ықтималдық бұлты, орбита емес","en":"Probability cloud, not an orbit"}[lang]
        out += [t(515, 736, label, 22)]
    elif kind == "isotopes":
        for x, name, neutrons in ((145, "C-12", 6), (555, "C-14", 8)):
            out += [box(x, 322, 330, 345), t(x+165, 375, name, 35),
                    '<circle cx="%d" cy="%d" r="86" fill="url(#nucleus)" stroke="#a75f65" stroke-width="3"/>' % (x+165, 500),
                    t(x+165, 492, "6p", 29), t(x+165, 533, f"{neutrons}n", 29)]
        out += [t(515, 735, {"ru":"Элемент задают 6 протонов","kz":"Элементті 6 протон анықтайды","en":"Six protons define the element"}[lang], 23)]
    elif kind == "half-life":
        labels=("0", "1", "2", "3")
        for i, count in enumerate((80, 40, 20, 10)):
            x=180+i*180
            out += [box(x, 665-count*4, 96, count*4, ("#74b9d2", "#83c4d1", "#a6d0b4", "#e7c78d")[i]),
                    t(x+48, 700, labels[i], 23), t(x+48, 630-count*4, str(count), 26)]
        out += [path("M150 666H830", "#7e9eac", 3), t(515, 750, {"ru":"Число периодов","kz":"Кезең саны","en":"Number of half-lives"}[lang], 23)]
    elif kind == "radiation":
        out += [box(145, 395, 160, 140, "#fff0dc"), t(225, 472, {"ru":"Источник","kz":"Көз","en":"Source"}[lang], 22),
                path("M325 445H445 M325 470H445 M325 495H445", "#dc9674", 4, "8 8"),
                box(450, 340, 75, 250, "#d8e5ea"),
                path("M550 445H695 M550 470H695 M550 495H695", "#88adbd", 3, "6 9"),
                '<circle cx="752" cy="421" r="26" fill="#d7a17e"/>', path("M752 450V550 M752 480L700 515 M752 480L802 515 M752 550L715 605 M752 550L789 605", "#3f7291", 9),
                t(490, 652, {"ru":"Барьер","kz":"Бөгет","en":"Shield"}[lang], 22)]
    elif kind == "photon":
        # One straight frequency-energy line; threshold at x corresponding to f=5.
        out += [path("M195 668H818 M195 668V300", "#688ba1", 4),
                path("M245 615L745 315", "#5793b4", 6),
                path("M545 668V435", "#d58a63", 3, "8 8"),
                '<circle cx="445" cy="495" r="12" fill="#e39b6e"/>',
                '<circle cx="645" cy="375" r="12" fill="#58a88c"/>',
                t(445, 462, "4", 27), t(545, 704, "5", 27), t(645, 342, "6", 27),
                t(784, 710, "f", 27), t(171, 305, "E", 27)]
    elif kind == "spectra":
        # Each energy unit is 60 px; 2, 5 and 8 are 180 px apart.
        for energy, y in ((8, 270), (5, 450), (2, 630)):
            out += [path(f"M260 {y}H745", "#5a92af", 5), t(210, y+8, str(energy), 28)]
        out += [path("M570 470V608 M558 590L570 608L582 590", "#d17d66", 6),
                t(650, 545, "5 − 2 = 3", 26), t(190, 710, "E", 28)]
    elif kind == "clocks":
        # Common scale: 10 px per microsecond; 0 at x=300.
        out += [path("M215 675H810 M300 290V690", "#718fa2", 3), t(300, 714, "0", 24)]
        for y, value, color in ((340, -7, "#dc9577"), (470, 45, "#70afc9"), (600, 38, "#88c6af")):
            x=300+value*10
            out += [f'<rect x="{min(x,300)}" y="{y}" width="{abs(value)*10}" height="55" rx="6" fill="{color}"/>',
                    t(x+(24 if value>=0 else -24), y+37, f"{value:+d}", 23, "start" if value>=0 else "end")]
        out += [t(545, 755, "−7 + 45 = +38", 27)]
    elif kind == "energy":
        labels={"ru":("Ядро", "Тепло", "Турбина", "Генератор"),
                "kz":("Ядро", "Жылу", "Турбина", "Генератор"),
                "en":("Nucleus", "Heat", "Turbine", "Generator")}[lang]
        for i,label in enumerate(labels):
            x=122+i*195
            out += [box(x, 393, 175, 160, ("#f8e1d7", "#fff0d5", "#e5f5ed", "#e4eff9")[i]), t(x+87, 478, label, 22)]
            if i<3:out += [path(f"M{x+180} 475h12 m-11-9 11 9-11 9", "#ca8664", 4)]
        out += [t(515, 680, {"ru":"Часть энергии → в среду","kz":"Энергияның бір бөлігі → ортаға","en":"Some energy → surroundings"}[lang], 23)]
    return out


def svg(row, lang):
    title, summary = row[lang][:2]
    cap={"ru":"Опорная карта: современная физика", "kz":"Тірек карта: қазіргі физика", "en":"Support map: modern physics"}[lang]
    out=['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
         f'<title id="title">{escape(title)}</title><desc id="desc">{escape(summary)}</desc>',
         '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#daeaf0"/></pattern><linearGradient id="paper" x2="1" y2="1"><stop stop-color="#fff"/><stop offset="1" stop-color="#e9f4f7"/></linearGradient><radialGradient id="nucleus"><stop stop-color="#fff5e7"/><stop offset="1" stop-color="#dc9b84"/></radialGradient><radialGradient id="cloud"><stop stop-color="#c6e6f3" stop-opacity=".9"/><stop offset="1" stop-color="#c6e6f3" stop-opacity="0"/></radialGradient><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="7" stdDeviation="7" flood-color="#315b74" flood-opacity=".16"/></filter></defs>',
         '<rect width="1600" height="900" fill="#f3f9fb"/><rect width="1600" height="900" fill="url(#grid)"/>',
         box(70, 190, 875, 590, "url(#paper)"), box(975, 190, 555, 590, "#fff"), t(75, 100, cap, 35, "start")]
    out += scene(row, lang)
    for i,value in enumerate(row["map"][lang]):
        y=235+i*165
        out += [box(1005,y,495,125,("#e6f5fa","#fff0e4","#eaf4e9")[i])]
        for j,part in enumerate(textwrap.wrap(value,29,break_long_words=False)):
            out += [t(1025,y+61+j*34,part,25,"start")]
    out += ['</svg>']
    return '\n'.join(out)+'\n'


def markdown(row, number, lang):
    title, summary, question, terms, example, activity, trap = row[lang]
    sections=(question,terms,example,activity,trap,
              f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)',
              CHECK[lang],JOURNAL[lang])
    heads=(HEAD[lang][0],HEAD[lang][1],HEAD[lang][2],HEAD[lang][3],HEAD[lang][4],HEAD[lang][5],HEAD[lang][7],HEAD[lang][8])
    body='\n\n'.join(f'## {h}\n\n{s}' for h,s in zip(heads,sections))
    source_label={"ru":"Источник сведений", "kz":"Мәлімет көзі", "en":"Evidence source"}[lang]
    body+=f'\n\n{source_label}: [{row["source"].split("/")[2]}]({row["source"]}).'
    if number<92:
        nxt=MODERN[number-84]
        label={"ru":"Следующий урок", "kz":"Келесі сабақ", "en":"Next lesson"}[lang]
        body+=f'\n\n[{label}: {nxt[lang][0]}](/read/physics-{nxt["stem"]}?lang={lang})'
    else:
        label={"ru":"Следующий урок", "kz":"Келесі сабақ", "en":"Next lesson"}[lang]
        next_title={"ru":"Можно ли ответить на вопрос опытом?", "kz":"Сұраққа тәжірибемен жауап беруге бола ма?", "en":"Can an experiment answer this question?"}[lang]
        body+=f'\n\n[{label}: {next_title}](/read/physics-93-testable-question?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'


def build():
    assert len(MODERN)==8
    for number,row in enumerate(MODERN,85):
        assert row['stem'].startswith(f'{number}-')
        for lang in ('ru','kz','en'):
            suffix='' if lang=='ru' else '-'+lang
            (LESSONS/f'{row["stem"]}{suffix}.md').write_text(markdown(row,number,lang),encoding='utf-8')
            (MAPS/f'map-{row["stem"]}-{lang}.svg').write_text(svg(row,lang),encoding='utf-8')


if __name__=='__main__':build()
