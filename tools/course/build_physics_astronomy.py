#!/usr/bin/env python3
"""Build the manually authored astronomy block and exact trilingual support maps."""
from html import escape
from pathlib import Path
import textwrap
try:
    from .physics_astronomy_data import ASTRONOMY
    from .build_physics_waves import HEAD, CHECK, JOURNAL
except ImportError:
    from physics_astronomy_data import ASTRONOMY
    from build_physics_waves import HEAD, CHECK, JOURNAL

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/physics'
MAPS = ROOT / 'web/static/course/physics'
SOURCES = {
    '75': 'https://science.nasa.gov/moon/facts/',
    '76': 'https://science.nasa.gov/earth/facts/',
    '77': 'https://science.nasa.gov/earth/facts/',
    '78': 'https://science.nasa.gov/helio-and-you-seasons-on-earth-mars-and-beyond/',
    '79': 'https://science.nasa.gov/moon/moon-phases/',
    '80': 'https://science.nasa.gov/eclipses/geometry/',
    '81': 'https://science.nasa.gov/learn/basics-of-space-flight/chapter3-4/',
    '82': 'https://science.nasa.gov/solar-system/planet-sizes-and-locations-in-our-solar-system/',
    '83': 'https://science.nasa.gov/universe/cosmic-distances/',
    '84': 'https://science.nasa.gov/earth/facts/',
}

def text(x, y, value, size=28, anchor='middle', color='#123b58'):
    return f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-family="Arial,sans-serif" font-size="{size}" font-weight="700" fill="{color}">{escape(value)}</text>'

def line(d, color='#267fa5', width=6, dash=''):
    extra=f' stroke-dasharray="{dash}"' if dash else ''
    return f'<path d="{d}" fill="none" stroke="{color}" stroke-width="{width}" stroke-linecap="round" stroke-linejoin="round"{extra}/>'

def circle(x,y,r,fill,stroke='#305e77'):
    earth={'#71b9d5','#6eb3cd','#69b4d3','#70bad5','#6eb7d3','#6bb4d0','#69b4d2','#69b7d5','#6eb4d0'}
    sun={'#f7c56a','#f5be67','#f3bd69','#f7c46e','#f5c06d','#f3bf69','#f5bd65'}
    moon={'#e8b77a','#e7c985','#f6d790','#e5c88d','#dfa364','#ddb680'}
    shade='url(#earthSphere)' if fill in earth else ('url(#sunSphere)' if fill in sun else ('url(#moonSphere)' if fill in moon else fill))
    shadow=' filter="url(#sphereShadow)"' if r>10 and fill!='none' else ''
    stroke_width='1' if r<=10 else '4'
    return f'<circle cx="{x}" cy="{y}" r="{r}" fill="{shade}" stroke="{stroke}" stroke-width="{stroke_width}"{shadow}/>'

def arrow(d,color='#df8652'):
    # All arrows are explicitly drawn; no oversized marker heads.
    return line(d,color,5)

def scene(kind,lang):
    names={
      'ru': {'earth':'Земля','moon':'Луна','sun':'Солнце','june':'Июнь','dec':'Декабрь','new':'Новая','quarter':'Четверть','full':'Полная','star':'Звезда','speed':'скорость','force':'сила','question':'Вопрос','model':'Модель','check':'Проверка','limit':'Граница'},
      'kz': {'earth':'Жер','moon':'Ай','sun':'Күн','june':'Маусым','dec':'Желтоқсан','new':'Жаңа','quarter':'Ширек','full':'Толған','star':'Жұлдыз','speed':'жылдамдық','force':'күш','question':'Сұрақ','model':'Үлгі','check':'Тексеру','limit':'Шектеу'},
      'en': {'earth':'Earth','moon':'Moon','sun':'Sun','june':'June','dec':'December','new':'New','quarter':'Quarter','full':'Full','star':'Star','speed':'velocity','force':'force','question':'Question','model':'Model','check':'Check','limit':'Limit'},
    }[lang]
    p=[]
    if kind=='scale':
        decimal='0.27 cm' if lang=='en' else '0,27 см'
        # 20 screen pixels represent one model centimetre in both diameter and separation.
        p += [line('M215 520H815','#849fab',3),circle(215,470,10,'#71b9d5'),circle(815,470,2.7,'#e8b77a'),line('M215 540H815','#db8858',4),line('M215 528V553 M815 528V553','#db8858',3),text(515,590,'30 cm' if lang=='en' else '30 см'),text(215,425,names['earth']),text(815,425,names['moon']),text(215,635,'1 cm' if lang=='en' else '1 см'),text(815,635,decimal)]
    elif kind=='rotation':
        east={'ru':'на восток','kz':'шығысқа','en':'eastwards'}[lang]
        p += [circle(260,480,58,'#f7c56a'),circle(550,480,145,'#6eb3cd'),'<path d="M550 335 A145 145 0 0 1 550 625Z" fill="#244963" opacity=".78"/>',line('M330 410H395 M330 480H395 M330 550H395','#edb960',4),circle(510,395,9,'#e79c68'),text(260,590,names['sun']),text(550,675,names['earth']),arrow('M415 300C505 250 605 255 660 325 M648 311L660 325L642 326'),text(550,235,east,24)]
    elif kind=='orbit':
        orbit_label={'ru':'1 обход ≈ 365,24 сут','kz':'1 айналым ≈ 365,24 тәулік','en':'1 orbit ≈ 365.24 d'}[lang]
        p += [circle(515,475,210,'none','#7eafbf'),circle(515,475,50,'#f5be67'),circle(725,475,23,'#69b4d3'),text(515,485,names['sun'],21),text(725,535,names['earth'],24),arrow('M657 319L679 343L649 352'),text(515,730,orbit_label,27)]
    elif kind=='seasons':
        angle='23.4°' if lang=='en' else '23,4°'
        p += [circle(515,480,70,'#f3bd69'),circle(270,480,75,'#70bad5'),circle(760,480,75,'#70bad5'),line('M270 395V565 M760 395V565','#b6d4dc',2,'7 7'),line('M233.3 565L306.7 395 M723.3 565L796.7 395','#dd8d60',5),line('M355 425H415 M355 480H415 M355 535H415 M615 425H675 M615 480H675 M615 535H675','#f0bb6b',4),text(515,490,names['sun'],22),text(270,640,names['june']),text(760,640,names['dec']),text(360,358,angle,25)]
    elif kind=='phases':
        p += [circle(185,500,54,'#f7c46e'),circle(535,500,60,'#6eb7d3'),'<ellipse cx="535" cy="500" rx="247" ry="145" fill="none" stroke="#85b3c3" stroke-width="4"/>']
        # Top view: the Sun is on the left, so the left half of each Moon is lit.
        for x,y in ((288,500),(535,355),(782,500)):
            p += [circle(x,y,19,'#445c72'),f'<path d="M{x} {y-19} A19 19 0 0 0 {x} {y+19}Z" fill="#f2d68f"/>']
        p += [text(185,610,names['sun']),text(535,590,names['earth']),text(288,660,names['new']),text(535,305,names['quarter']),text(782,660,names['full']),line('M245 435H272 M245 500H265 M245 565H272','#efb85f',3)]
    elif kind=='eclipse':
        plane={'ru':'Плоскости орбит','kz':'Орбита жазықтықтары','en':'Orbital planes'}[lang]
        miss={'ru':'Тень проходит мимо','kz':'Көлеңке Жерге тимейді','en':'Shadow misses Earth'}[lang]
        p += [line('M250 290H820','#df9b65',4,'9 10'),line('M250 315L820 265','#7eaac0',4),text(780,245,'≈5°',25),text(490,355,plane,22),circle(225,550,53,'#f5c06d'),circle(725,550,57,'#6bb4d0'),circle(490,435,18,'#e5c88d'),line('M285 550H652','#df9b65',4,'9 10'),line('M505 435H665','#cf8c68',3,'8 9'),text(225,650,names['sun']),text(725,650,names['earth']),text(490,405,names['moon']),text(495,720,miss,23)]
    elif kind=='gravity':
        p += [circle(520,520,185,'none','#85b6c8'),circle(520,520,115,'#69b4d2'),circle(520,335,15,'#dfa364'),arrow('M540 335H690 M672 323L690 335L672 347','#dd8755'),arrow('M520 355V400 M508 381L520 400L532 381','#5a91ad'),text(720,300,names['speed'],24),text(435,390,names['force'],24),text(520,540,names['earth'])]
    elif kind=='planets':
        planets={'ru':('Марс','Юпитер','1 а.е.','1,5 а.е.','5,2 а.е.'),'kz':('Марс','Юпитер','1 а.б.','1,5 а.б.','5,2 а.б.'),'en':('Mars','Jupiter','1 AU','1.5 AU','5.2 AU')}[lang]
        sizes={'ru':'Размеры планет схематичны','kz':'Планета өлшемдері схемалық','en':'Planet sizes are schematic'}[lang]
        p += [line('M180 500H860','#789ead',4),circle(180,500,42,'#f3bf69'),circle(310,500,13,'#69b7d5'),circle(375,500,11,'#cf9168'),circle(856,500,22,'#ddb680'),text(180,590,names['sun'],22),text(310,590,names['earth'],21),text(375,635,planets[0],21),text(820,590,planets[1],21),text(310,420,planets[2],22),text(390,420,planets[3],22),text(815,420,planets[4],22),text(515,720,sizes,22)]
    elif kind=='starlight':
        scale_note={'ru':'Расстояния без масштаба','kz':'Қашықтықтар масштабсыз','en':'Not to distance scale'}[lang]
        star_time='≈4.2 y' if lang=='en' else ('≈4,2 года' if lang=='ru' else '≈4,2 жыл')
        p += [circle(215,445,52,'#f5bd65'),circle(515,445,28,'#6eb4d0'),'<path d="M805 383L820 430L870 445L820 460L805 507L790 460L745 445L790 430Z" fill="#f7dd94" stroke="#ca9a69" stroke-width="3"/>',line('M280 445H470 M550 445H745','#dfab67',4,'10 9'),text(215,565,names['sun']),text(515,565,names['earth']),text(805,565,names['star']),text(365,390,'≈8 min',25),text(665,390,star_time,25),text(515,710,scale_note,22)]
    elif kind=='checkpoint':
        labels=[names['question'],names['model'],names['check'],names['limit']]
        for i,label in enumerate(labels):
            x=122+i*195
            p += [f'<rect x="{x}" y="390" width="175" height="160" rx="12" fill="{("#e6f5fa","#fff0e4","#eaf4e9","#edf0fa")[i]}" stroke="#93bac9" stroke-width="3"/>',text(x+87,480,label,23)]
            if i<3:p += [arrow(f'M{x+181} 470h12 m-10-10 10 10-10 10','#cf8b63')]
        summary={'ru':'сутки · сезон · фаза','kz':'тәулік · маусым · фаза','en':'day · season · phase'}[lang]
        p += [text(515,680,summary,25)]
    return p

def svg(row,lang):
    title=row[lang][0]
    # Keep the heading short in the visual. The full localized title is in <title>.
    caption={'ru':'Опорная карта: проверяем модель','kz':'Тірек карта: үлгіні тексереміз','en':'Support map: test the model'}[lang]
    translated_values={
      'scale': {'ru':('Земля Ø 1 см','Луна Ø 0,27 см','Центры: 30 см'),'kz':('Жер Ø 1 см','Ай Ø 0,27 см','Орталықтар: 30 см'),'en':('Earth Ø 1 cm','Moon Ø 0.27 cm','Centres: 30 cm')},
      'rotation': {'ru':('Сутки ≈ 24 ч','Поворот ≈ 15°/ч','Земля: на восток'),'kz':('Тәулік ≈ 24 сағ','Бұрылу ≈ 15°/сағ','Жер: шығысқа'),'en':('Day ≈ 24 h','Turn ≈ 15°/h','Earth: eastwards')},
      'orbit': {'ru':('Один обход','Год ≈ 365,24 сут','Вращение ≠ обход'),'kz':('Бір айналып шығу','Жыл ≈ 365,24 тәулік','Ось ≠ орбита'),'en':('One circuit','Year ≈ 365.24 d','Spin ≠ orbit')},
      'seasons': {'ru':('Ось: ≈ 23,4°','Июнь: лето на севере','Декабрь: лето на юге'),'kz':('Ось: ≈ 23,4°','Маусым: солтүстікте жаз','Желтоқсан: оңтүстікте жаз'),'en':('Axis: ≈ 23.4°','June: northern summer','December: southern summer')},
      'phases': {'ru':('Новая: Луна слева','Четверть: Луна сверху','Полная: Луна справа'),'kz':('Жаңа Ай: сол жақта','Ширек: Ай жоғарыда','Толған Ай: оң жақта'),'en':('New: Moon at left','Quarter: Moon above','Full: Moon at right')},
      'eclipse': {'ru':('Новая + узел','Полная + узел','Наклон ≈ 5°'),'kz':('Жаңа Ай + түйін','Толған Ай + түйін','Еңкею ≈ 5°'),'en':('New + node','Full + node','Tilt ≈ 5°')},
      'gravity': {'ru':('Скорость: по касательной','Сила: к центру','Орбита: кривой путь'),'kz':('Жылдамдық: жанама','Күш: орталыққа','Орбита: қисық жол'),'en':('Velocity: tangent','Force: to centre','Orbit: curved path')},
      'planets': {'ru':('Земля: 1 а.е. = 10 см','Марс: 1,5 а.е. = 15 см','Юпитер: 5,2 а.е. = 52 см'),'kz':('Жер: 1 а.б. = 10 см','Марс: 1,5 а.б. = 15 см','Юпитер: 5,2 а.б. = 52 см'),'en':('Earth: 1 AU = 10 cm','Mars: 1.5 AU = 15 cm','Jupiter: 5.2 AU = 52 cm')},
      'starlight': {'ru':('Күн → Жер ≈ 8 мин','Звезда → ≈ 4,2 года','Световой год: длина'),'kz':('Күн → Жер ≈ 8 мин','Жұлдыз → ≈ 4,2 жыл','Жарық жылы: ұзындық'),'en':('Sun → Earth ≈ 8 min','Star → ≈ 4.2 years','Light-year: distance')},
      'checkpoint': {'ru':('Сутки: вращение','Сезон: наклон оси','Фаза: освещение'),'kz':('Тәулік: айналу','Маусым: ось еңкеюі','Фаза: жарық түсуі'),'en':('Day: rotation','Season: axis tilt','Phase: illumination')},
    }
    vals=translated_values[row['kind']][lang]
    out=['<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-labelledby="title desc">',f'<title id="title">{escape(title)}</title><desc id="desc">{escape(row[lang][1])}</desc>',
      '<defs><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#d9eaf0"/></pattern><linearGradient id="paper" x2="1" y2="1"><stop stop-color="#fff"/><stop offset="1" stop-color="#e8f3f7"/></linearGradient><radialGradient id="earthSphere" cx="35%" cy="28%" r="75%"><stop stop-color="#dcf5f7"/><stop offset=".45" stop-color="#72bcd6"/><stop offset="1" stop-color="#2d719d"/></radialGradient><radialGradient id="sunSphere" cx="35%" cy="28%" r="75%"><stop stop-color="#fff7c5"/><stop offset=".5" stop-color="#fac778"/><stop offset="1" stop-color="#d77d4c"/></radialGradient><radialGradient id="moonSphere" cx="35%" cy="28%" r="75%"><stop stop-color="#fff4dc"/><stop offset=".5" stop-color="#e6c592"/><stop offset="1" stop-color="#a98c75"/></radialGradient><filter id="shadow" x="-20%" y="-20%" width="150%" height="160%"><feDropShadow dx="0" dy="10" stdDeviation="8" flood-color="#37657b" flood-opacity=".18"/></filter><filter id="sphereShadow" x="-25%" y="-25%" width="160%" height="160%"><feDropShadow dx="3" dy="6" stdDeviation="4" flood-color="#23465a" flood-opacity=".22"/></filter></defs>',
      '<rect width="1600" height="900" fill="#f2f9fb"/><rect width="1600" height="900" fill="url(#grid)"/>',
      '<rect x="70" y="190" width="875" height="590" rx="14" fill="url(#paper)" stroke="#93bac9" stroke-width="3" filter="url(#shadow)"/>',
      '<rect x="975" y="190" width="555" height="590" rx="14" fill="#fff" stroke="#93bac9" stroke-width="3" filter="url(#shadow)"/>',
      text(75,100,caption,36,'start')]
    out+=scene(row['kind'],lang)
    for i,value in enumerate(vals):
        y=235+i*165
        out += [f'<rect x="1005" y="{y}" width="495" height="125" rx="10" fill="{("#e6f5fa","#fff0e4","#eaf4e9")[i]}" stroke="#d0e2e8" stroke-width="2"/>']
        wrapped=textwrap.wrap(value,29,break_long_words=False)
        for j,part in enumerate(wrapped):out += [text(1025,y+58+j*35,part,27,'start')]
    out += ['</svg>']
    return '\n'.join(out)+'\n'

def markdown(row,number,lang):
    title,summary,question,terms,example,activity,trap=row[lang]
    sections=(question,terms,example,activity,trap,
              f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)',
              CHECK[lang],JOURNAL[lang])
    heads=(HEAD[lang][0],HEAD[lang][1],HEAD[lang][2],HEAD[lang][3],HEAD[lang][4],HEAD[lang][5],HEAD[lang][7],HEAD[lang][8])
    body='\n\n'.join(f'## {h}\n\n{s}' for h,s in zip(heads,sections))
    source_label={'ru':'Источник чисел и модели','kz':'Сандар мен үлгінің дерек көзі','en':'Source for values and model'}[lang]
    body+=f'\n\n{source_label}: [NASA Science]({SOURCES[str(number)]}).'
    if number<84:
        n=ASTRONOMY[number-74]
        label={'ru':'Следующий урок','kz':'Келесі сабақ','en':'Next lesson'}[lang]
        body+=f'\n\n[{label}: {n[lang][0]}](/read/physics-{n["stem"]}?lang={lang})'
    else:
        label={'ru':'Вернуться к курсу «Физика»','kz':'«Физика» курсына оралу','en':'Return to Physics course'}[lang]
        body+=f'\n\n[{label}](/course/physics?lang={lang})'
    return f'# {title}\n\n_Лид (summary):_ **{summary}**\n\n{body}\n'

def build():
    assert len(ASTRONOMY)==10
    for number,row in enumerate(ASTRONOMY,75):
        assert row['stem'].startswith(f'{number}-')
        for lang in ('ru','kz','en'):
            suffix='' if lang=='ru' else '-'+lang
            (LESSONS/f'{row["stem"]}{suffix}.md').write_text(markdown(row,number,lang),encoding='utf-8')
            (MAPS/f'map-{row["stem"]}-{lang}.svg').write_text(svg(row,lang),encoding='utf-8')
if __name__=='__main__':build()
