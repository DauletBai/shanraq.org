#!/usr/bin/env python3
"""Build trilingual electricity lessons and dimensioned support maps (47–56)."""
from html import escape
from pathlib import Path

try:
    from .physics_electricity_data import ELECTRICITY
except ImportError:
    from physics_electricity_data import ELECTRICITY

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/physics"
MAPS = ROOT / "web/static/course/physics"

HEAD = {
    "ru": ("Ситуация и прогноз", "Новые слова без пропусков", "Разбор учебных данных", "Безопасная проверка", "Ловушка в рассуждении", "Опорная карта", "Перенос в новую задачу", "Запись в лабораторный журнал"),
    "kz": ("Жағдай және болжам", "Жаңа сөздерді түсінейік", "Оқу деректерін талдау", "Қауіпсіз тексеру", "Ойлаудағы қате", "Тірек картасы", "Жаңа жағдайға қолдану", "Зертхана журналына жазба"),
    "en": ("Situation and prediction", "New words, step by step", "Working through supplied data", "A safe check", "A reasoning trap", "Support map", "Transfer to a new case", "Lab notebook entry"),
}
JOURNAL = {
    "ru": "Добавьте к сквозному «Лабораторному журналу моего мира» страницу с четырьмя строками: вопрос; прогноз до чтения ответа; наблюдение или расчёт по заданным числам; уточнённый вывод и условие его применимости. Если опыта не было, не пишите «измерено». После схемы объясните словами, где переносится заряд, где передаётся энергия и где модель идеализирует настоящий прибор.",
    "kz": "«Менің әлемімнің зертхана журналы» жобасына төрт жолды бет қосыңыз: сұрақ; жауапты оқымастан бұрынғы болжам; бақылау немесе берілген сандар бойынша есеп; нақтыланған қорытынды және оның қолданылу шарты. Тәжірибе болмаса, «өлшенді» деп жазбаңыз. Сызбадан кейін заряд қайда тасымалданатынын, энергия қайда берілетінін және модель нақты құралды қай жерде жеңілдететінін сөзбен түсіндіріңіз.",
    "en": "Add a page to My World Lab Notebook with four lines: question; prediction before reading the answer; observation or calculation from supplied numbers; revised conclusion and its conditions. If you did not do an experiment, do not write ‘measured’. After the diagram, explain where charge moves, where energy is transferred, and where the model simplifies a real device.",
}
MAP_WORDS = {
    "ru": [
        ("Заряд не создаётся", "Нейтральное тоже тянется"),
        ("Один путь туда и обратно", "Открытый ключ: тока нет"),
        ("Ток = заряд / время", "Амперметр: последовательно"),
        ("Вольт = джоуль / кулон", "Вольтметр: параллельно"),
        ("Ом: сопротивление", "Для омического резистора"),
        ("Одна ветвь: ток одинаков", "Напряжения складываются"),
        ("Ветви: напряжение одинаково", "Токи ветвей складываются"),
        ("Обрыв: пути нет", "Проверьте режим прибора"),
        ("Ватт = джоуль / секунда", "Энергия = мощность × время"),
        ("Два независимых ключа", "Модель: только резисторы"),
    ],
    "kz": [
        ("Заряд жоқтан шықпайды", "Бейтарап зат та тартылады"),
        ("Барып-қайтатын бір жол", "Ашық кілт: ток жоқ"),
        ("Ток = заряд / уақыт", "Амперметр: тізбектей"),
        ("Вольт = джоуль / кулон", "Вольтметр: параллель"),
        ("Ом: кедергі", "Омдық резистор үшін"),
        ("Бір тармақ: ток бірдей", "Кернеулер қосылады"),
        ("Тармақтарда кернеу бірдей", "Тармақ токтары қосылады"),
        ("Үзіліс: жол жоқ", "Құрал режимін тексеріңіз"),
        ("Ватт = джоуль / секунд", "Энергия = қуат × уақыт"),
        ("Екі тәуелсіз кілт", "Модель: тек резисторлар"),
    ],
    "en": [
        ("Charge is not created", "Neutral objects can attract"),
        ("One complete return route", "Open switch: no current"),
        ("Current = charge / time", "Ammeter: in series"),
        ("Volt = joule / coulomb", "Voltmeter: in parallel"),
        ("Ohm: resistance", "For an ohmic resistor"),
        ("One branch: equal current", "Voltage drops add"),
        ("Branches share voltage", "Branch currents add"),
        ("Open gap: no path", "Check meter mode"),
        ("Watt = joule / second", "Energy = power × time"),
        ("Two independent switches", "Model: resistors only"),
    ],
}
MAP_VALUES = {
    47: ("+ и −" , "Σq = 0"),
    48: ("U = 3 V", "Open: I = 0 A"),
    49: ("Δq = 6 C; Δt = 3 s", "I = 2 A"),
    50: ("ΔE = 12 J; Δq = 4 C", "U = 3 V"),
    51: ("U = 3 V; R = 300 Ω", "I = 10 mA"),
    52: ("R₁ = 100 Ω; R₂ = 200 Ω", "I = 10 mA; U₁+U₂ = 3 V"),
    53: ("I₁ = 10 mA; I₂ = 5 mA", "Iобщ = 15 mA; U = 3 V"),
    54: ("U = 3 V", "Open: I = 0 A"),
    55: ("P = 0.06 W; t = 60 s", "E = 3.6 J"),
    56: ("I₁ = 10 mA; I₂ = 5 mA", "Iобщ = 15 mA; E = 2.7 J"),
}


def wrap(value, limit):
    lines, line = [], ""
    for word in value.split():
        candidate = (line + " " + word).strip()
        if line and len(candidate) > limit:
            lines.append(line)
            line = word
        else:
            line = candidate
    if line:
        lines.append(line)
    return lines


def label(x, y, value, size=30, colour="#17415a", anchor="start"):
    return f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" font-size="{size}" font-weight="650" fill="{colour}">{escape(value)}</text>'


def battery(x=250, y=465):
    return [f'<path d="M{x} {y-62}v42m-18 0h36m-10 14h20m-28 0v62" stroke="#194863" stroke-width="8" fill="none" stroke-linecap="round"/>',
            label(x-30, y-24, "+", 30), label(x-30, y+34, "−", 30)]


def resistor(x, y, name):
    return [f'<rect x="{x}" y="{y-29}" width="148" height="58" rx="8" fill="#f7c778" stroke="#b66b35" stroke-width="6"/>',
            label(x+74, y+9, name, 27, anchor="middle")]


def circuit(number):
    """Circuit paths are aligned to the 40 px grid; crossings are never joints."""
    if number == 47:
        out = ['<ellipse cx="365" cy="475" rx="110" ry="145" fill="#f5a384" stroke="#bd6853" stroke-width="7"/>',
               '<path d="M365 620v65" stroke="#bd6853" stroke-width="5"/>',
               '<rect x="650" y="390" width="140" height="140" rx="8" fill="#f9f4e7" stroke="#ad9874" stroke-width="6"/>']
        out += [label(340,445,"−",48),label(400,485,"−",48),label(610,435,"+",48),label(735,435,"−",48),
                '<path d="M490 475h115m-24-20 24 20-24 20" stroke="#e38458" stroke-width="8" fill="none"/>']
        return out
    if number == 49:
        return ['<path d="M225 465h570" stroke="#1d6385" stroke-width="16" stroke-linecap="round"/>',
                '<path d="M495 340v250" stroke="#ec906a" stroke-width="8" stroke-dasharray="13 13"/>',
                '<path d="M595 415h105m-25-22 25 22-25 22" fill="none" stroke="#e8895d" stroke-width="8"/>',
                label(505,630,"Δq / Δt",38,anchor="middle")]
    if number == 50:
        out = ['<path d="M250 400V340H770V409 M770 521V590H250V520" fill="none" stroke="#1d6385" stroke-width="9"/>']
        out += battery()
        out += ['<circle cx="770" cy="465" r="56" fill="#fff5d5" stroke="#c29a4b" stroke-width="7"/>',
                '<path d="M735 430l70 70m0-70-70 70" stroke="#d79346" stroke-width="6"/>',
                '<path d="M770 409h115v22 M770 521h115v-22" fill="none" stroke="#7388b0" stroke-width="6"/>',
                '<circle cx="885" cy="465" r="34" fill="#eff3fb" stroke="#7388b0" stroke-width="6"/>',
                label(885,478,"V",34,anchor="middle")]
        return out
    if number in (53, 56):
        branches = 'M480 320V345 M480 385V408 M700 320V345 M700 385V408' if number==56 else 'M480 320V408 M700 320V408'
        out = [f'<path d="M250 400V320H800 M250 520V610H800 {branches} M480 522V610 M700 522V610" fill="none" stroke="#1d6385" stroke-width="9"/>']
        out += battery()
        out += ['<rect x="440" y="408" width="80" height="114" rx="8" fill="#f7c778" stroke="#b66b35" stroke-width="6"/>',
                '<rect x="660" y="408" width="80" height="114" rx="8" fill="#f7c778" stroke="#b66b35" stroke-width="6"/>',
                label(480,460,"300",25,anchor="middle"),label(700,460,"600",25,anchor="middle"),
                label(480,500,"Ω",28,anchor="middle"),label(700,500,"Ω",28,anchor="middle")]
        if number == 56:
            out += ['<path d="M480 345v40 M700 345v40" stroke="#e48859" stroke-width="8"/>',
                    '<circle cx="480" cy="345" r="9" fill="#e48859"/><circle cx="480" cy="385" r="9" fill="#e48859"/>',
                    '<circle cx="700" cy="345" r="9" fill="#e48859"/><circle cx="700" cy="385" r="9" fill="#e48859"/>',
                    label(480,292,"S₁",27,anchor="middle"),label(700,292,"S₂",27,anchor="middle")]
        return out
    top = 'M250 400V320H490 M548 320H795' if number in (48,54) else ('M250 400V320H420 M570 320h68 M786 320h9' if number==52 else ('M250 400V320H795' if number==55 else 'M250 400V320H480 M628 320h167'))
    right = 'M795 320V407 M795 523V610H250V520' if number in (48,54,55) else 'M795 320V610H250V520'
    out = [f'<path d="{top} {right}" fill="none" stroke="#1d6385" stroke-width="9"/>']
    out += battery()
    if number in (48, 54):
        out += ['<circle cx="795" cy="465" r="58" fill="#fff3cb" stroke="#c09146" stroke-width="7"/>',
                '<path d="M758 428l74 74m0-74-74 74" stroke="#ce904a" stroke-width="6"/>']
        if number == 54:
            out += ['<path d="M490 320h58" stroke="#f1f8fa" stroke-width="18"/>',
                    '<circle cx="490" cy="320" r="9" fill="#e57d62"/><circle cx="548" cy="320" r="9" fill="#e57d62"/>']
        else:
            out += ['<path d="M490 320h52 M490 320l48-38" stroke="#e48859" stroke-width="8"/>']
    elif number == 52:
        out += resistor(422,320,"100 Ω") + resistor(638,320,"200 Ω")
    elif number != 55:
        out += resistor(480,320,"300 Ω")
    else:
        out += ['<circle cx="795" cy="465" r="52" fill="#fff3cb" stroke="#c09146" stroke-width="7"/>',
                '<path d="M760 430l70 70m0-70-70 70" stroke="#ce904a" stroke-width="6"/>']
    return out


def svg(row, number, lang):
    title = row[lang][0]
    notes = MAP_WORDS[lang][number-47]
    figures = MAP_VALUES[number]
    notes = [*notes, *figures]
    if lang != "en":
        notes = [note.replace("Open:", "Үзіліс:" if lang=="kz" else "Обрыв:")
                 .replace("Iобщ", "Iжалпы" if lang=="kz" else "Iобщ")
                 .replace("mA", "мА").replace(" W", " Вт")
                 .replace(" V", " В").replace(" A", " А")
                 .replace(" J", " Дж").replace(" C", " Кл").replace(" s", " с")
                 .replace("0.06", "0,06").replace("3.6", "3,6").replace("2.7", "2,7") for note in notes]
    else:
        notes = [note.replace("Iобщ", "Itotal").replace("и −", "and −") for note in notes]
    if any(len(wrap(note, 25)) > 2 for note in notes):
        raise ValueError(f"Map note too long: {number} {lang}: {notes}")
    title_lines = wrap(title, 35)
    if len(title_lines)>3:
        raise ValueError(f"Map title too long: {number} {lang}")
    out = ['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
           f'<title id="title">{escape(title)}</title><desc id="desc">{escape("; ".join(notes))}</desc>',
           '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#d3e9ef"/></pattern><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="12" stdDeviation="10" flood-color="#345d73" flood-opacity=".22"/></filter><linearGradient id="panel" x2="1" y2="1"><stop stop-color="#fff"/><stop offset="1" stop-color="#e6f2f5"/></linearGradient></defs>',
           '<rect width="1600" height="900" fill="#eff8fa"/><rect width="1600" height="900" fill="url(#grid)"/>',
           '<rect x="74" y="205" width="860" height="555" rx="20" fill="url(#panel)" stroke="#95b8c7" stroke-width="3" filter="url(#shadow)"/>',
           '<rect x="975" y="205" width="550" height="555" rx="20" fill="#fff" stroke="#95b8c7" stroke-width="3" filter="url(#shadow)"/>']
    for index, line in enumerate(title_lines):
        out.append(label(78, 72 + 45*index, line, 33))
    out.extend(circuit(number))
    for index, note in enumerate(notes):
        y=238+index*126
        out.append(f'<rect x="1002" y="{y}" width="494" height="99" rx="12" fill="{("#e8f5fb","#fff0e2","#eaf5ed","#edf0fa")[index]}" stroke="#d1e2e9"/>')
        for line_index,line in enumerate(wrap(note,25)):
            out.append(label(1020,y+40+line_index*31,line,26))
    out.append('</svg>')
    return '\n'.join(out)+'\n'


def markdown(row, number, lang):
    title, summary, scene, terms, worked, task, mistake, transfer = row[lang]
    headings=HEAD[lang]
    sections=[scene,terms,worked,task,mistake,
              f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)',transfer,JOURNAL[lang]]
    body='\n\n'.join(f'## {heading}\n\n{content}' for heading,content in zip(headings,sections))
    if number<56:
        next_row=ELECTRICITY[number-46]
        text={"ru":"Следующий урок","kz":"Келесі сабақ","en":"Next lesson"}[lang]
        body+=f'\n\n[{text}: {next_row[lang][0]}](/read/physics-{next_row["stem"]}?lang={lang})'
    else:
        text={"ru":"Вернуться к курсу «Физика»","kz":"«Физика» курсына оралу","en":"Return to the Physics course"}[lang]
        body+=f'\n\n[{text}](/course/physics?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'


def build():
    if len(ELECTRICITY)!=10:
        raise ValueError("Electricity block needs ten lessons")
    for number,row in enumerate(ELECTRICITY,47):
        if not row["stem"].startswith(f"{number:02d}-"):
            raise ValueError(f"Wrong lesson number: {row['stem']}")
        for lang in ("ru","kz","en"):
            suffix="" if lang=="ru" else "-"+lang
            (LESSONS/f'{row["stem"]}{suffix}.md').write_text(markdown(row,number,lang),encoding="utf-8")
            (MAPS/f'map-{row["stem"]}-{lang}.svg').write_text(svg(row,number,lang),encoding="utf-8")
    print("Built 30 electricity lessons and 30 trilingual maps")


if __name__ == "__main__":
    build()
