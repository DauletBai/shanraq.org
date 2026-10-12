#!/usr/bin/env python3
"""Build physics lessons 29–36 and exact, translated support maps."""
from html import escape
from pathlib import Path
import re

try:
    from .physics_fluids_data import FLUIDS
except ImportError:
    from physics_fluids_data import FLUIDS

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/physics'
MAPS = ROOT / 'web/static/course/physics'

HEAD = {
 'ru': ('История и предсказание','Слова, которые понадобятся','Проверяем модель на числах','Безопасное исследование','Где легко ошибиться','Опорная карта','Новый случай без подсказки','Запись в лабораторный журнал'),
 'kz': ('Оқиға және болжам','Қажет сөздер','Модельді сандармен тексеру','Қауіпсіз зерттеу','Жиі қателесетін жер','Тірек картасы','Көмексіз жаңа жағдай','Зертхана журналына жазба'),
 'en': ('Story and prediction','Words to learn','Check the model with numbers','A safe investigation','A likely mistake','Support map','A new case without hints','Lab notebook entry'),
}
COMMON = {
 'ru': ('Перед чтением запишите своё предположение в журнале и объясните его одним предложением. После расчёта вернитесь к предположению: совпадение ответа с ожиданием ещё не подтверждает объяснение. Важно проверить, какие условия были заданы, а какие вы действительно наблюдали.',
        'Закройте разобранный пример. Нарисуйте ситуацию заново, назовите тело и действующие на него силы. Затем измените только одно данное число, пересчитайте результат с единицами и скажите, что осталось неизменным. Если не было измерительного прибора, напишите «расчёт по учебной модели», а не «измерено».',
        'Сделайте четыре строки: вопрос; предсказание; данные или наблюдение; вывод с оговоркой. Рядом подпишите одну причину, по которой реальный опыт может отличаться от идеального рисунка. Вернитесь к этому листу после следующего урока и исправьте вывод, если новая проверка его изменила.'),
 'kz': ('Оқымас бұрын болжамыңызды журналға жазып, оны бір сөйлеммен түсіндіріңіз. Есептен кейін болжамға қайта оралыңыз: жауаптың күткеніңізге сай келуі түсіндірмені әлі дәлелдемейді. Қай шарт берілгенін және нені шынымен байқағаныңызды тексеріңіз.',
        'Шығарылған мысалды жауып қойыңыз. Жағдайды қайта салып, денені және оған әсер ететін күштерді атаңыз. Содан соң тек бір берілген санды өзгертіп, бірліктерімен қайта есептеңіз және ненің өзгермегенін айтыңыз. Өлшеу құралы болмаса, «өлшенді» емес, «оқу моделі бойынша есептелді» деп жазыңыз.',
        'Төрт жол толтырыңыз: сұрақ; болжам; берілген сан не бақылау; шарттары көрсетілген қорытынды. Жанына нақты тәжірибенің идеал суреттен өзгеше болу себебін жазыңыз. Келесі сабақтан соң бұл параққа қайтып, жаңа тексеру қорытындыны өзгертсе, оны түзетіңіз.'),
 'en': ('Before reading, put your prediction in your notebook and justify it in one sentence. After calculating, return to it: a result matching your guess does not yet prove your explanation. Check which conditions were supplied and which things you actually observed.',
        'Cover the worked example. Redraw the situation, name the object and the forces acting on it. Then change just one given number, recalculate with units and say what stayed the same. If you had no measuring instrument, write “calculated from a worked model”, not “measured”.',
        'Write four lines: question; prediction; given data or observation; conclusion with a condition. Add one reason a real trial may differ from the ideal drawing. Return to this page after the next lesson and revise the conclusion if a new check changes it.'),
}
DIAGRAM = {
 'pressure': {'ru':['Одна сила: 20 Н','Малая площадь: 0,02 м²','Большая площадь: 0,04 м²','1000 Па → 500 Па'],
              'kz':['Бір күш: 20 Н','Кіші аудан: 0,02 м²','Үлкен аудан: 0,04 м²','1000 Па → 500 Па'],
              'en':['One force: 20 N','Small area: 0.02 m²','Large area: 0.04 m²','1000 Pa → 500 Pa']},
 'conversion': {'ru':['1 м² = 10 000 см²','2 см² = 0,0002 м²','F = 4 Н','p = 20 000 Па = 20 кПа'],
                'kz':['1 м² = 10 000 см²','2 см² = 0,0002 м²','F = 4 Н','p = 20 000 Па = 20 кПа'],
                'en':['1 m² = 10,000 cm²','2 cm² = 0.0002 m²','F = 4 N','p = 20,000 Pa = 20 kPa']},
 'depth': {'ru':['Вода: ρ = 1000 кг/м³','h₁ = 0,2 м → 2000 Па','h₂ = 0,4 м → 4000 Па','p = ρgh (сверх поверхности)'],
           'kz':['Су: ρ = 1000 кг/м³','h₁ = 0,2 м → 2000 Па','h₂ = 0,4 м → 4000 Па','p = ρgh (бетке қосымша)'],
           'en':['Water: ρ = 1000 kg/m³','h₁ = 0.2 m → 2000 Pa','h₂ = 0.4 m → 4000 Pa','p = ρgh (above surface)']},
 'vessels': {'ru':['Одна вода + один воздух','Равные уровни в покое','Высота башни: 5 м','Идеально: 50 кПа'],
             'kz':['Бір су + бір ауа','Тынышта деңгейлер тең','Мұнара биіктігі: 5 м','Идеал: 50 кПа'],
             'en':['One liquid + equal air','Same level at rest','Tower height: 5 m','Ideal: 50 kPa']},
 'buoyancy': {'ru':['Погружённый V = 0,0001 м³','Вода вытеснена: 1 Н','Вес куба: 3 Н','На динамометре: 2 Н'],
              'kz':['Батқан V = 0,0001 м³','Ығысқан су салмағы: 1 Н','Текше салмағы: 3 Н','Динамометр: 2 Н'],
              'en':['Submerged V = 0.0001 m³','Displaced water: 1 N','Cube weight: 3 N','Force meter: 2 N']},
 'floating': {'ru':['Масса лодки + груза: 0,3 кг','Вес: 3 Н','Нужный объём: 0,0003 м³','То же: 0,3 л'],
              'kz':['Қайық + жүк: 0,3 кг','Ауырлық күші: 3 Н','Қажет көлем: 0,0003 м³','Яғни: 0,3 л'],
              'en':['Boat + cargo: 0.3 kg','Weight: 3 N','Required volume: 0.0003 m³','That is: 0.3 L']},
 'hydraulic': {'ru':['Малый: 2 см²; F₁ = 4 Н','Большой: 10 см²; F₂ = 20 Н','Пути: 5 см → 1 см','Работа: 0,2 Дж = 0,2 Дж'],
               'kz':['Кіші: 2 см²; F₁ = 4 Н','Үлкен: 10 см²; F₂ = 20 Н','Жол: 5 см → 1 см','Жұмыс: 0,2 Дж = 0,2 Дж'],
               'en':['Small: 2 cm²; F₁ = 4 N','Large: 10 cm²; F₂ = 20 N','Travel: 5 cm → 1 cm','Work: 0.2 J = 0.2 J']},
 'capstone': {'ru':['Платформа + груз: 0,4 кг','Вес: 4 Н','Вытеснить: 0,4 л воды','0,1 м глубина: +1000 Па'],
              'kz':['Тұғыр + жүк: 0,4 кг','Ауырлық күші: 4 Н','Ығыстыру: 0,4 л су','0,1 м тереңдік: +1000 Па'],
              'en':['Platform + cargo: 0.4 kg','Weight: 4 N','Displace: 0.4 L of water','0.1 m depth: +1000 Pa']},
}

def wrapped(text, max_chars=62):
    lines=[]; current=''
    for word in text.split():
        trial=(current+' '+word).strip()
        if current and len(trial)>max_chars:
            lines.append(current);current=word
        else:current=trial
    if current:lines.append(current)
    return lines

def txt(x,y,s,size=28,fill='#203e57',weight=600,anchor='start'):
    return f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Inter,Arial,sans-serif" font-size="{size}" font-weight="{weight}" fill="{fill}">{escape(s)}</text>'

def svg(row,lang):
    labels=DIAGRAM[row['kind']][lang]
    title=row[lang][0]
    out=['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',
         '<title id="title">'+escape(title)+'</title>',
         '<desc id="desc">'+escape(' | '.join(labels))+'</desc>',
         '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M 40 0 L 0 0 0 40" fill="none" stroke="#d3e8ef" stroke-width="1"/></pattern><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="14" stdDeviation="12" flood-color="#31516c" flood-opacity=".20"/></filter><linearGradient id="water" x1="0" x2="0" y1="0" y2="1"><stop stop-color="#b3edf8"/><stop offset="1" stop-color="#168bbd"/></linearGradient><linearGradient id="plate" x1="0" x2="1"><stop stop-color="#f9fcff"/><stop offset="1" stop-color="#d7e9f0"/></linearGradient></defs>',
         '<rect width="1600" height="900" fill="#edf6f8"/><rect width="1600" height="900" fill="url(#grid)"/>',
         '<rect x="65" y="188" width="915" height="565" rx="22" fill="url(#plate)" stroke="#9fbfcb" stroke-width="3" filter="url(#shadow)"/>',
         '<rect x="1015" y="188" width="520" height="565" rx="22" fill="#ffffff" stroke="#a5c5d0" stroke-width="3" filter="url(#shadow)"/>']
    for i,line in enumerate(wrapped(title,54)[:2]):out.append(txt(80,80+56*i,line,43,'#143e5a',750))
    # Exact numerical model on the right; four rows have fixed 100px pitch.
    for i,label in enumerate(labels):
        yy=240+i*112
        out += [f'<rect x="1045" y="{yy}" width="460" height="84" rx="12" fill="{["#e9f6fc","#edf8f4","#fff4e7","#e9f0fa"][i]}" stroke="#d0e2e8"/>']
        for j,line in enumerate(wrapped(label,30)[:2]):
            out.append(txt(1070,yy+34+j*31,line,26 if len(label)>29 else 29,'#1d425c',650))
    kind=row['kind']
    if kind in ('pressure','conversion'):
        out += ['<rect x="220" y="540" width="600" height="22" rx="5" fill="#637f8b"/>',
                '<rect x="275" y="415" width="160" height="125" rx="8" fill="#d58559" stroke="#9f4e37" stroke-width="4"/>',
                '<rect x="575" y="415" width="230" height="125" rx="8" fill="#d58559" stroke="#9f4e37" stroke-width="4"/>',
                '<path d="M355 355v53m-15-18 15 18 15-18 M690 355v53m-15-18 15 18 15-18" fill="none" stroke="#267ba5" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"/>']
        out += [txt(355,610,labels[1],25,anchor='middle'),txt(690,655,labels[2],25,anchor='middle')]
    elif kind=='depth':
        out += ['<path d="M220 345V650H810V345" fill="none" stroke="#416b82" stroke-width="11"/>',
                '<path d="M226 400H804V645H226Z" fill="url(#water)" opacity=".82"/>',
                '<path d="M225 400H805" fill="none" stroke="#ffffff" stroke-width="7"/>',
                '<circle cx="430" cy="485" r="12" fill="#f5b44f"/><circle cx="430" cy="570" r="12" fill="#d26358"/>',
                '<path d="M460 400v170" stroke="#173d58" stroke-width="5"/>',
                txt(500,490,'0,2 м' if lang!='en' else '0.2 m',27),txt(500,575,'0,4 м' if lang!='en' else '0.4 m',27)]
    elif kind=='vessels':
        # Two open columns share a lower channel; they are physically connected.
        out += ['<path d="M240 345V650H830V345 M450 345V570 M620 345V570" fill="none" stroke="#416b82" stroke-width="11"/>',
                '<path d="M246 400H444V575H626V400H824V644H246Z" fill="url(#water)" opacity=".82"/>',
                '<path d="M246 400H444 M626 400H824" fill="none" stroke="#ffffff" stroke-width="7"/>',
                '<path d="M350 380H720" stroke="#ec9b49" stroke-width="4" stroke-dasharray="14 11"/>',
                txt(535,705,'=',46,'#2185a9',750,'middle')]
    elif kind=='buoyancy':
        # The worked body really is a completely submerged cube held by a string.
        out += ['<path d="M205 365V650H850V365" fill="none" stroke="#416b82" stroke-width="11"/>',
                '<path d="M212 405H843V643H212Z" fill="url(#water)" opacity=".83"/>',
                '<path d="M212 405H843" stroke="white" stroke-width="7"/>',
                '<path d="M530 285V470" stroke="#617681" stroke-width="5"/>',
                '<rect x="470" y="470" width="120" height="120" rx="4" fill="#c97951" stroke="#8e4f39" stroke-width="6"/>',
                '<path d="M530 625V595m-18 18 18-18 18 18" fill="none" stroke="#fff8d8" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"/>',
                txt(620,550,'Fₐ' if lang!='en' else 'Fᵦ',32,'#1c748b')]
    elif kind in ('floating','capstone'):
        out += ['<path d="M205 365V650H850V365" fill="none" stroke="#416b82" stroke-width="11"/>',
                '<path d="M212 435H843V643H212Z" fill="url(#water)" opacity=".83"/>',
                '<path d="M212 435H843" stroke="white" stroke-width="7"/>',
                '<path d="M390 385L680 385L645 510L425 510Z" fill="#c97951" stroke="#8e4f39" stroke-width="6"/>',
                '<path d="M535 590V510m-18 18 18-18 18 18" fill="none" stroke="#fff8d8" stroke-width="9" stroke-linecap="round" stroke-linejoin="round"/>',
                '<path d="M535 290V365m-18-18 18 18 18-18" fill="none" stroke="#d26c5b" stroke-width="9" stroke-linecap="round" stroke-linejoin="round"/>']
        out += [txt(735,350,'Fₐ' if lang!='en' else 'Fᵦ',32,'#1c748b'),txt(590,600,'V',32,'white')]
    elif kind=='hydraulic':
        out += ['<path d="M250 570V350H450V570 M620 570V420H825V570 M350 565H720" fill="none" stroke="#527284" stroke-width="14"/>',
                '<path d="M260 470H440V560H635V485H815V560H260Z" fill="url(#water)"/>',
                '<path d="M255 465H445 M615 480H830" stroke="#be7652" stroke-width="22"/>',
                '<path d="M350 330v110m-18-18 18 18 18-18" fill="none" stroke="#ca684e" stroke-width="8"/>',
                '<path d="M720 550V430m-18 18 18-18 18 18" fill="none" stroke="#238ca3" stroke-width="8"/>',
                txt(350,635,'2 см²' if lang!='en' else '2 cm²',29,anchor='middle'),txt(720,635,'10 см²' if lang!='en' else '10 cm²',29,anchor='middle')]
    out += ['</svg>']
    return '\n'.join(out)+'\n'

def markdown(index,row,lang):
    title,summary,scene,terms,worked,task,mistake=row[lang]
    h=HEAD[lang]; pre,challenge,journal=COMMON[lang]
    sections=[(h[0],scene+'\n\n'+pre),(h[1],terms),(h[2],worked),(h[3],task),
              (h[4],mistake),(h[5],f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)'),
              (h[6],challenge),(h[7],journal)]
    body='\n\n'.join(f'## {name}\n\n{text}' for name,text in sections)
    if index < 36:
        nxt=FLUIDS[index-28]
        label={'ru':'Следующий урок','kz':'Келесі сабақ','en':'Next lesson'}[lang]
        body+=f'\n\n[{label}: {nxt[lang][0]}](/read/physics-{nxt["stem"]}?lang={lang})'
    else:
        label={'ru':'Вернуться к курсу «Физика»','kz':'«Физика» курсына оралу','en':'Return to the Physics course'}[lang]
        body+=f'\n\n[{label}](/course/physics?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'

def build():
    assert len(FLUIDS)==8
    for i,row in enumerate(FLUIDS,29):
        assert row['stem'].startswith(f'{i:02d}-')
        for lang in ('ru','kz','en'):
            suffix='' if lang=='ru' else '-'+lang
            (LESSONS/f'{row["stem"]}{suffix}.md').write_text(markdown(i,row,lang),encoding='utf-8')
            (MAPS/f'map-{row["stem"]}-{lang}.svg').write_text(svg(row,lang),encoding='utf-8')
    print('Built 24 fluid lessons and 24 trilingual exact-data maps')

if __name__=='__main__':build()
