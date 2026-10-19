#!/usr/bin/env python3
"""Render the final eight physics lessons and matching three-language support maps."""

from html import escape
from pathlib import Path
import textwrap

try:
    from .physics_final_data import FINAL
    from .build_physics_waves import HEAD, CHECK, JOURNAL
except ImportError:
    from physics_final_data import FINAL
    from build_physics_waves import HEAD, CHECK, JOURNAL

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/physics'
MAPS = ROOT / 'web/static/course/physics'
SOURCE = 'https://openstax.org/books/college-physics-ap-courses-2e/pages/5-1-friction'
NOTES = {
    'question': {
        'ru': ('Меняем только покрытие', 'Измеряем путь до остановки', 'Числа пока учебные'),
        'kz': ('Тек төсеніш өзгереді', 'Тоқтағанға дейінгі жол', 'Сандар әзірше оқу дерегі'),
        'en': ('Change only the covering', 'Measure distance to rest', 'Values are supplied data')},
    'measure': {
        'ru': ('Ноль у выхода со ската', 'Путь = конец − начало', 'Не забудьте единицу'),
        'kz': ('Нөл еңістен шыққанда', 'Жол = соңы − басы', 'Бірлікті ұмытпаңыз'),
        'en': ('Zero at ramp exit', 'Distance = end − start', 'Always include a unit')},
    'repeats': {
        'ru': ('Три отдельных запуска', 'Среднее: 83 и 45 см', 'Размах не равен всей ошибке'),
        'kz': ('Үш бөлек жүріс', 'Орташа: 83 және 45 см', 'Ауқым барлық қате емес'),
        'en': ('Three separate runs', 'Means: 83 and 45 cm', 'Range is not total error')},
    'compare': {
        'ru': ('Сохраняем условия', 'Разность: 38 см', 'Причину проверяем отдельно'),
        'kz': ('Шарттарды бірдей ұстаймыз', 'Айырма: 38 см', 'Себепті бөлек тексереміз'),
        'en': ('Keep conditions alike', 'Difference: 38 cm', 'Check rival causes')},
    'graph': {
        'ru': ('Ось начинается с нуля', 'Бағандар: 83 и 45 см', 'Путь — не скорость'),
        'kz': ('Ось нөлден басталады', 'Бағандар: 83 және 45 см', 'Жол — жылдамдық емес'),
        'en': ('Axis begins at zero', 'Bars: 83 and 45 cm', 'Distance is not speed')},
    'evidence': {
        'ru': ('Данные: учебная таблица', 'Источник: объяснение', 'Вывод имеет границы'),
        'kz': ('Дерек: оқу кестесі', 'Дереккөз: түсіндірме', 'Қорытындының шегі бар'),
        'en': ('Data: supplied table', 'Source: explanation', 'Claim has limits')},
    'design': {
        'ru': ('Критерий: не более 60 см', 'Ткань даёт средние 45 см', 'Нужен новый реальный тест'),
        'kz': ('Шек: 60 см-ден аспау', 'Мата: орташа 45 см', 'Жаңа нақты сынақ керек'),
        'en': ('Criterion: at most 60 cm', 'Fabric mean: 45 cm', 'A fresh real test is needed')},
    'defend': {
        'ru': ('Вопрос и метод открыты', 'Все шесть чисел показаны', 'Ограничение названо'),
        'kz': ('Сұрақ пен әдіс ашық', 'Алты сан да көрсетілген', 'Шектеу аталған'),
        'en': ('Question and method clear', 'All six values visible', 'Limitation stated')},
}


def label(x, y, value, size=24, anchor='middle', color='#173c57'):
    return (f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" '
            f'font-size="{size}" font-weight="700" fill="{color}">{escape(str(value))}</text>')


def line(d, color='#397e9b', width=4, dash=''):
    dash_attr = f' stroke-dasharray="{dash}"' if dash else ''
    return f'<path d="{d}" fill="none" stroke="{color}" stroke-width="{width}" stroke-linecap="round" stroke-linejoin="round"{dash_attr}/>'


def card(x, y, width, height, fill='#f8fcfd'):
    return (f'<rect x="{x}" y="{y}" width="{width}" height="{height}" rx="15" '
            f'fill="{fill}" stroke="#9ebfca" stroke-width="2" filter="url(#shadow)"/>')


def scene(row, lang):
    kind = row['kind']
    unit = 'cm' if lang == 'en' else 'см'
    out = []
    if kind == 'question':
        out += [line('M145 360L345 540H840', '#6c98a8', 10),
                '<circle cx="340" cy="515" r="22" fill="#e9ad70"/>',
                '<rect x="305" y="492" width="100" height="23" rx="5" fill="#5c91ad"/>',
                line('M490 595H822', '#5d95ad', 7), line('M490 710H670', '#dc9271', 7),
                label(500, 575, {'ru':'бумага','kz':'қағаз','en':'paper'}[lang], 23, 'start'),
                label(500, 690, {'ru':'ткань','kz':'мата','en':'fabric'}[lang], 23, 'start'),
                label(822, 575, f'83 {unit}', 21, 'end'), label(670, 690, f'45 {unit}', 21, 'end'),
                label(175, 330, {'ru':'один старт','kz':'бір бастау','en':'one start'}[lang], 23, 'start')]
    elif kind == 'measure':
        out += [line('M190 520H865', '#678fa3', 6)]
        for cm in range(0, 91, 10):
            x = 200 + cm * 7
            out += [line(f'M{x} 510V545', '#658fa2', 2), label(x, 580, str(cm), 19)]
        out += [line('M214 415V635', '#d99270', 4, '8 8'),
                line('M788 415V635', '#d99270', 4, '8 8'),
                line('M214 420H788', '#e3a267', 5),
                label(215, 394, f'2 {unit}', 24), label(788, 394, f'84 {unit}', 24),
                label(500, 690, f'84 − 2 = 82 {unit}', 32)]
    elif kind == 'repeats':
        values = [82, 84, 83, 44, 46, 45]
        out += [line('M175 685H870', '#698f9f', 4)]
        for i, value in enumerate(values):
            x = 210 + i * 108
            h = value * 3.8
            out += [f'<rect x="{x}" y="{685-h:.1f}" width="63" height="{h:.1f}" rx="5" fill="{("url(#barBlue)" if i<3 else "url(#barPeach)")}" filter="url(#shadow)"/>',
                    label(x+31, 670-h, str(value), 22)]
        out += [label(350, 745, {'ru':'бумага','kz':'қағаз','en':'paper'}[lang], 23),
                label(675, 745, {'ru':'ткань','kz':'мата','en':'fabric'}[lang], 23),
                label(190, 322, unit, 21)]
    elif kind in ('compare', 'graph'):
        out += [line('M220 700H850', '#668a9c', 4), line('M220 700V300', '#668a9c', 4)]
        for tick in (0, 30, 60, 90):
            y = 700 - tick * 4
            out += [line(f'M213 {y}H850', '#d3e4e9', 1), label(185, y+8, str(tick), 20)]
        for x, value, gradient in ((350, 83, 'barBlue'), (620, 45, 'barPeach')):
            out += [f'<rect x="{x}" y="{700-value*4}" width="120" height="{value*4}" rx="6" fill="url(#{gradient})" filter="url(#shadow)"/>',
                    label(x+60, 680-value*4, f'{value} {unit}', 25)]
        out += [label(410, 747, {'ru':'бумага','kz':'қағаз','en':'paper'}[lang], 23),
                label(680, 747, {'ru':'ткань','kz':'мата','en':'fabric'}[lang], 23)]
    elif kind == 'evidence':
        names = {'ru':('ДАННЫЕ','ИСТОЧНИК','ВЫВОД'), 'kz':('ДЕРЕК','ДЕРЕККӨЗ','ҚОРЫТЫНДЫ'),
                 'en':('DATA','SOURCE','CLAIM')}[lang]
        for i, name in enumerate(names):
            x = 140 + i * 250
            out += [card(x, 350, 220, 215, ('#e2f2fa','#f9eadb','#eaf5e9')[i]), label(x+110, 420, name, 21)]
        out += [label(250, 505, f'83 / 45 {unit}', 22), label(500, 505, 'OpenStax', 22),
                label(750, 505, {'ru':'при условиях','kz':'осы шартта','en':'under conditions'}[lang], 19)]
    elif kind == 'design':
        out += [line('M205 490H840', '#7b9aa8', 5)]
        for value, color, y in ((45, '#d99371', 400), (60, '#cfb772', 490), (83, '#68acc6', 580)):
            x = 205 + value * 7
            out += [line(f'M{x} {y-35}V{y+35}', color, 7), label(x, y-54, f'{value} {unit}', 24)]
        out += [label(525, 710, f'60 − 45 = 15 {unit}', 29)]
    elif kind == 'defend':
        words = {'ru':('ВОПРОС','МЕТОД','ДАННЫЕ','ВЫВОД'),
                 'kz':('СҰРАҚ','ӘДІС','ДЕРЕК','ҚОРЫТЫНДЫ'),
                 'en':('QUESTION','METHOD','DATA','CLAIM')}[lang]
        for i, word in enumerate(words):
            x = 135 + i * 190
            out += [card(x, 420, 175, 125, ('#e5f4fb','#f5e9d8','#e9f3e4','#eaeaf8')[i]),
                    label(x+87, 493, word, 20)]
            if i < 3:
                out += [line(f'M{x+178} 482H{x+189}', '#c58566', 4)]
        out += [label(510, 665, f'83 − 45 = 38 {unit}', 29)]
    return out


def svg(row, lang):
    title, summary = row[lang][:2]
    cap = {'ru':'ФИЗИКА · ИТОГОВОЕ ИССЛЕДОВАНИЕ', 'kz':'ФИЗИКА · ҚОРЫТЫНДЫ ЗЕРТТЕУ',
           'en':'PHYSICS · FINAL INVESTIGATION'}[lang]
    out = ['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
           f'<title id="title">{escape(title)}</title><desc id="desc">{escape(summary)}</desc>',
           '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#dcecf1"/></pattern><linearGradient id="barBlue" x2="1" y2=".2"><stop stop-color="#9bd1e3"/><stop offset=".55" stop-color="#65abc3"/><stop offset="1" stop-color="#3e819e"/></linearGradient><linearGradient id="barPeach" x2="1" y2=".2"><stop stop-color="#f3c4a0"/><stop offset=".55" stop-color="#db9975"/><stop offset="1" stop-color="#af6b55"/></linearGradient><filter id="shadow" x="-20%" y="-20%" width="150%" height="150%"><feDropShadow dx="0" dy="5" stdDeviation="6" flood-color="#234c62" flood-opacity=".18"/></filter></defs>',
           '<rect width="1600" height="900" fill="#f5f9fa"/><rect width="1600" height="900" fill="url(#grid)"/>',
           card(70, 190, 875, 590, '#fff'), card(975, 190, 555, 590, '#fff'),
           label(74, 105, cap, 34, 'start')]
    out.extend(scene(row, lang))
    for i, note in enumerate(NOTES[row['kind']][lang]):
        words = textwrap.wrap(note, 30, break_long_words=False, break_on_hyphens=False)
        y = 226 + i * 171
        out += [card(1004, y, 497, 146, ('#e6f5fa','#fff1df','#edf6e8')[i])]
        for j, word in enumerate(words):
            out += [label(1020, y+37+j*29, word, 21, 'start')]
    out.append('</svg>')
    return '\n'.join(out)+'\n'


def markdown(row, number, lang):
    title, summary, question, terms, example, activity, trap = row[lang]
    image = f'/static/course/physics/map-{row["stem"]}-{lang}.svg'
    headings = HEAD[lang]
    sections = (question, terms, example, activity, trap, f'![{title}]({image})',
                CHECK[lang], JOURNAL[lang])
    chosen = (headings[0], headings[1], headings[2], headings[3], headings[4],
              headings[5], headings[7], headings[8])
    body = '\n\n'.join(f'## {heading}\n\n{section}' for heading, section in zip(chosen, sections))
    if number == 98:
        source_label = {'ru':'Проверяемый источник объяснения', 'kz':'Тексерілетін түсіндірме көзі',
                        'en':'Checkable source for the explanation'}[lang]
        body += f'\n\n{source_label}: [OpenStax: Friction]({SOURCE}).'
    if number < 100:
        next_row = FINAL[number - 92]
        next_title = next_row[lang][0]
        link_label = {'ru':'Следующий урок', 'kz':'Келесі сабақ', 'en':'Next lesson'}[lang]
        body += f'\n\n[{link_label}: {next_title}](/read/physics-{next_row["stem"]}?lang={lang})'
    else:
        link_label = {'ru':'Вернуться к курсу «Физика»', 'kz':'«Физика» курсына оралу',
                      'en':'Return to Physics course'}[lang]
        body += f'\n\n[{link_label}](/course/physics?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'


def build():
    assert len(FINAL) == 8
    for number, row in enumerate(FINAL, 93):
        assert row['stem'].startswith(str(number) + '-')
        for lang in ('ru', 'kz', 'en'):
            suffix = '' if lang == 'ru' else '-' + lang
            (LESSONS / f'{row["stem"]}{suffix}.md').write_text(markdown(row, number, lang), encoding='utf-8')
            (MAPS / f'map-{row["stem"]}-{lang}.svg').write_text(svg(row, lang), encoding='utf-8')


if __name__ == '__main__':
    build()
