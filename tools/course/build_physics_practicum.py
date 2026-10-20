#!/usr/bin/env python3
"""Build four trilingual practical physics lessons and their localized support maps."""

from html import escape
from pathlib import Path

try:
    from .physics_practicum_data import PRACTICUM
except ImportError:
    from physics_practicum_data import PRACTICUM

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/physics'
MAPS = ROOT / 'web/static/course/physics'

HEAD = {
    'ru': ('Задача', 'Опорная схема понятий', 'Разбор на числах', 'Практика без риска', 'Проверьте себя', 'Карта опыта', 'Журнал исследования', 'Следующий урок'),
    'kz': ('Мәселе', 'Тірек ұғымдар сызбасы', 'Сандармен талдау', 'Қауіпсіз тәжірибе', 'Өзіңізді тексеріңіз', 'Тәжірибе картасы', 'Зерттеу журналы', 'Келесі сабақ'),
    'en': ('The question', 'Concept map', 'Worked numbers', 'Safe practice', 'Check your reasoning', 'Investigation map', 'Lab notebook', 'Next lesson'),
}
JOURNAL = {
    'ru': 'Запишите в журнал три колонки: «Что предсказал», «Что показывают данные» и «Чего данные не доказывают». Сначала заполните первую колонку без подсказки; после расчёта исправьте объяснение, если оно не выдержало проверки. Если выполняли настоящий опыт, отделите свои измерения от учебных чисел.',
    'kz': 'Журналға үш баған жазыңыз: «Бастапқы болжамым», «Дерек нені көрсетеді» және «Дерек нені дәлелдемейді». Алдымен алғашқы бағанды көмексіз толтырыңыз; есептен кейін түсіндірме сәйкес келмесе, түзетіңіз. Нақты тәжірибе жасасаңыз, өз өлшеуіңізді оқу сандарынан бөлек көрсетіңіз.',
    'en': 'Give your notebook three columns: “My prediction”, “What the data show”, and “What the data cannot prove”. Fill the first before looking at the worked answer; revise an explanation that fails the check. If you perform a real experiment, keep your own readings distinct from the supplied practice values.',
}


def svg(row, lang):
    title, summary = row[lang][:2]
    labels = row['map'][lang]
    tops = {'ru': 'ФИЗИКА · ПРАКТИКУМ', 'kz': 'ФИЗИКА · ПРАКТИКУМ', 'en': 'PHYSICS · PRACTICUM'}
    parts = [
        '<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
        f'<title id="title">{escape(title)}</title><desc id="desc">{escape(summary)}</desc>',
        '<defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#e7f5f5"/><stop offset="1" stop-color="#f7efe2"/></linearGradient><filter id="shadow"><feDropShadow dy="6" stdDeviation="9" flood-color="#24516a" flood-opacity=".13"/></filter></defs>',
        '<rect width="1600" height="900" fill="url(#bg)"/>',
        f'<text x="90" y="112" font-size="37" font-weight="700" fill="#20516a" font-family="Arial,sans-serif">{escape(tops[lang])}</text>',
        '<path d="M110 660H1450" stroke="#8eafbb" stroke-width="5"/><path d="M1450 660l-20-15v30z" fill="#8eafbb"/>',
    ]
    colors = ('#ddf0f4', '#e4f3df', '#fff0d8')
    for i, label in enumerate(labels):
        x = 110 + 490 * i
        parts.extend([
            f'<rect x="{x}" y="270" width="430" height="290" rx="26" fill="{colors[i]}" stroke="#80aeba" stroke-width="3" filter="url(#shadow)"/>',
            f'<circle cx="{x+215}" cy="336" r="28" fill="#245d78"/>',
            f'<text x="{x+215}" y="346" text-anchor="middle" fill="white" font-size="27" font-weight="700" font-family="Arial,sans-serif">{i+1}</text>',
            f'<text x="{x+215}" y="432" text-anchor="middle" fill="#183d55" font-size="26" font-weight="700" font-family="Arial,sans-serif" textLength="380" lengthAdjust="spacingAndGlyphs">{escape(label)}</text>',
        ])
        if i < 2:
            parts.append(f'<path d="M{x+438} 415h44m-15-15 15 15-15 15" fill="none" stroke="#c57858" stroke-width="6"/>')
    foot = {'ru': 'Учебные данные · проверьте условия опыта', 'kz': 'Оқу деректері · тәжірибе шарттарын тексеріңіз',
            'en': 'Supplied practice data · check the conditions'}[lang]
    parts += [f'<text x="800" y="754" text-anchor="middle" fill="#3b6272" font-size="30" font-family="Arial,sans-serif">{escape(foot)}</text>', '</svg>']
    return '\n'.join(parts) + '\n'


def markdown(row, number, lang):
    title, summary, opening, concepts, example, practice, check = row[lang]
    h = HEAD[lang]
    sections = (opening, concepts, example, practice, check,
                f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)', JOURNAL[lang])
    body = '\n\n'.join(f'## {heading}\n\n{content}' for heading, content in zip(h[:7], sections))
    if number < 104:
        following = PRACTICUM[number-100]
        slug = f'physics-{following["stem"]}'
        body += f'\n\n[{h[7]}: {following[lang][0]}](/read/{slug}?lang={lang})'
    else:
        back = {'ru': 'Вернуться к курсу «Физика»', 'kz': '«Физика» курсына оралу',
                'en': 'Return to Physics course'}[lang]
        body += f'\n\n[{back}](/course/physics?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'


def build():
    assert len(PRACTICUM) == 4
    for number, row in enumerate(PRACTICUM, 101):
        assert row['stem'].startswith(str(number) + '-')
        for lang in ('ru', 'kz', 'en'):
            suffix = '' if lang == 'ru' else '-' + lang
            (LESSONS / f'{row["stem"]}{suffix}.md').write_text(markdown(row, number, lang), encoding='utf-8')
            (MAPS / f'map-{row["stem"]}-{lang}.svg').write_text(svg(row, lang), encoding='utf-8')


if __name__ == '__main__':
    build()
