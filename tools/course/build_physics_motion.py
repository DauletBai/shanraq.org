#!/usr/bin/env python3
"""Build manually written trilingual motion lessons and numeric SVG maps."""
from html import escape
from pathlib import Path
import math
try:
    from .physics_motion_data import MOTION
except ImportError:
    from physics_motion_data import MOTION
ROOT=Path(__file__).resolve().parents[2]
LESSONS=ROOT/'course/lessons/physics'
MAPS=ROOT/'web/static/course/physics'
COVER='/static/covers/school/physics/motion/02-motion-graphs.webp'
TEXT={
'ru':{'heads':['Где мы на карте','Сцена из жизни','Новые слова без пропусков','Разберём числа вместе','Ваш опыт на бумаге','Ошибка, которую легко сделать','Проверьте, что поняли','Запись в журнал'], 'intro':'Это второй блок курса. Сохраняйте в «Лабораторном журнале моего мира» вопрос, данные и источник каждого числа. Числа в этой схеме — учебная модель, а не результат вашего измерения.', 'check':'Закройте пример и объясните себе вслух: что задано, что вычислено и в каких единицах? Затем измените одно число в условии и скажите, какое действие пришлось повторить. Если объяснение требует слова, которого вы ещё не понимаете, вернитесь к его определению выше.', 'note':'Запишите ответ вместе с единицами и короткое ограничение: что этот пример ещё не доказывает. Если опыт не проводили, пометьте числа «учебные данные».', 'next':'Следующий урок'},
'kz':{'heads':['Картадағы орнымыз','Өмірден алынған көрініс','Жаңа сөздерді түсінейік','Сандарды бірге талдайық','Қағаздағы өз тапсырмаңыз','Жиі кездесетін қате','Түсінгеніңізді тексеріңіз','Журналға жазба'], 'intro':'Бұл — курстың екінші блогы. «Менің әлемімнің зертхана журналында» сұрақты, деректерді және әр санның көзін сақтаңыз. Сызбадағы сандар — өз өлшеміңіз емес, оқу моделі.', 'check':'Мысалды жауып, дауыстап түсіндіріңіз: не берілген, не есептелген және қандай бірлікпен? Сосын шарттағы бір санды өзгертіп, қай амалды қайталау керегін айтыңыз. Түсінбейтін жаңа сөз қажет болса, жоғарыдағы анықтамасына қайта оралыңыз.', 'note':'Жауапты бірлігімен және бұл мысал әлі нені дәлелдемейтінін көрсететін қысқа шектеумен жазыңыз. Тәжірибе жасамасаңыз, сандарды «оқу деректері» деп белгілеңіз.', 'next':'Келесі сабақ'},
'en':{'heads':['Where we are on the map','A scene from everyday life','New words, step by step','Work through the numbers','Your paper investigation','An easy mistake to make','Check your understanding','Notebook entry'], 'intro':'This is the second block. Keep the question, data, and source of each number in My World Lab Notebook. The numbers in the map are a worked model, not a measurement you performed.', 'check':'Close the worked example and explain aloud: what was given, what was calculated, and in which units? Change one number in the question and say which calculation must be repeated. If you need a word you do not yet understand, return to its definition above.', 'note':'Record the answer with units and one short limit: what does this example not yet prove? If you did not perform a trial, label the numbers “worked data.”', 'next':'Next lesson'}}

def slug(row): return 'physics-'+row['stem']
def txt(x,y,s,size=28,fill='#173b5e',anchor='start',weight=600):
 return f'<text x="{x}" y="{y}" font-family="Arial,sans-serif" font-size="{size}" fill="{fill}" font-weight="{weight}" text-anchor="{anchor}">{escape(s)}</text>'
def line(x1,y1,x2,y2,color='#397397',width=4):
 return f'<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" stroke-width="{width}" stroke-linecap="round"/>'
def circle(x,y,r=13,color='#e1794e'):
 return f'<circle cx="{x}" cy="{y}" r="{r}" fill="{color}" stroke="white" stroke-width="3"/>'
def base(title,subtitle):
 title_size=min(44,max(27,int(1800/max(len(title),1))))
 return [f'''<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}"><defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#e9f4fa"/><stop offset="1" stop-color="#fff3e6"/></linearGradient><filter id="sh"><feDropShadow dx="0" dy="14" stdDeviation="13" flood-color="#315c78" flood-opacity=".22"/></filter></defs><rect width="1600" height="900" fill="url(#bg)"/><path d="M0 190H1600M0 390H1600M0 590H1600M0 790H1600M200 0V900M400 0V900M600 0V900M800 0V900M1000 0V900M1200 0V900M1400 0V900" fill="none" stroke="#cbdce6" stroke-opacity=".4" stroke-width="2"/>''', txt(800,104,title,title_size,'#143853','middle',800),txt(800,151,subtitle,24,'#547083','middle',500),'<g filter="url(#sh)"><rect x="130" y="205" width="1340" height="580" rx="26" fill="white"/></g>']

def map_svg(index,row,lang):
 title=row[lang][0]; kind=row['kind']; parts=base(title,{'ru':'Учебная модель · проверяемые числа','kz':'Оқу моделі · тексерілетін сандар','en':'Worked model · checkable numbers'}[lang]); k=lang
 if kind in ('position','route'):
  y=505; x0=270; scale=105
  parts += [line(270,y,1320,y,'#234d70',6)]
  for n in range(0,9):
   x=x0+n*scale;parts += [line(x,y-16,x,y+16,'#234d70',3),txt(x,y+55,str(n),25,'#274d69','middle')]
  parts += [txt(800,690,{'ru':'Ось положения x, м','kz':'Орын осі x, м','en':'Position axis x, m'}[k],30,'#315d76','middle')]
  if kind=='position':
   for n,c in ((2,'#258db9'),(6,'#e68250')): parts += [circle(x0+n*scale,y-60,30,c),line(x0+n*scale,y-35,x0+n*scale,y-5,c,3)]
   parts += [txt(800,310,{'ru':'2 м → 6 м: изменение +4 м','kz':'2 м → 6 м: өзгеріс +4 м','en':'2 m → 6 m: change +4 m'}[k],35,'#a34a43','middle',700)]
  else:
   parts += [f'<path d="M{x0} 425 H{x0+6*scale-18}" fill="none" stroke="#258db9" stroke-width="12" stroke-linecap="round"/>',f'<path d="M{x0+6*scale} 425 l-28 -18 v36z" fill="#258db9"/>',f'<path d="M{x0+6*scale} 590 H{x0+2*scale+18}" fill="none" stroke="#e68250" stroke-width="12" stroke-linecap="round"/>',f'<path d="M{x0+2*scale} 590 l28 -18 v36z" fill="#e68250"/>',txt(800,315,{'ru':'0 → 6 → 2: путь 10 м, перемещение +2 м','kz':'0 → 6 → 2: жол 10 м, орын ауыстыру +2 м','en':'0 → 6 → 2: distance 10 m, displacement +2 m'}[k],30,'#a34a43','middle',700)]
 elif kind=='time':
  parts += [line(270,500,1330,500,'#234d70',6)]
  for n in range(0,10):
   x=270+n*110;parts += [line(x,485,x,515,'#234d70',3),txt(x,565,str(n),23,'#234d70','middle')]
  parts += [f'<path d="M490 400V350H1040V400" fill="none" stroke="#e68250" stroke-width="7"/>',txt(765,330,{'ru':'7 с − 2 с = 5 с','kz':'7 с − 2 с = 5 с','en':'7 s − 2 s = 5 s'}[k],40,'#a34a43','middle',800),txt(800,690,{'ru':'Показания времени t, с','kz':'Уақыт көрсеткіші t, с','en':'Time reading t, s'}[k],30,'#315d76','middle')]
 elif kind in ('speed','units','average'):
  cards={
   'speed':{'ru':('Путь','12 м','Время','4 с','Скорость','3 м/с'),'kz':('Жол','12 м','Уақыт','4 с','Жылдамдық','3 м/с'),'en':('Distance','12 m','Time','4 s','Speed','3 m/s')},
   'units':{'ru':('За час','36 км','За секунду','10 м','Тот же темп','36 км/ч = 10 м/с'),'kz':('Бір сағатта','36 км','Бір секундта','10 м','Сол қарқын','36 км/сағ = 10 м/с'),'en':('In one hour','36 km','In one second','10 m','Same rate','36 km/h = 10 m/s')},
   'average':{'ru':('Участок 1','30 м / 10 с','Участок 2','10 м / 10 с','Весь путь','40 м / 20 с = 2 м/с'),'kz':('1-бөлік','30 м / 10 с','2-бөлік','10 м / 10 с','Бүкіл жол','40 м / 20 с = 2 м/с'),'en':('Leg 1','30 m / 10 s','Leg 2','10 m / 10 s','Whole trip','40 m / 20 s = 2 m/s')}
  }[kind][k]
  for j,x in enumerate((365,800,1235)):
   lab,val=cards[2*j:2*j+2];parts += [f'<rect x="{x-170}" y="345" width="350" height="260" rx="18" fill="{("#e5f3fa","#e5f8f5","#fff1df")[j]}" stroke="#a8cad9" stroke-width="2"/>',txt(x,415,lab,26,'#315d76','middle'),txt(x,515,val,25 if len(val)>16 else 34,'#173b5e','middle',800)]
  parts += [txt(800,690,{'ru':'Всегда записывайте единицы','kz':'Бірліктерді әрдайым жазыңыз','en':'Always write the units'}[k],30,'#a34a43','middle')]
 else:
  # Numeric position–time diagrams: axes share exact scales (2 s and 2 m per division).
  x0,y0=300,690; sx,sy=150,55
  parts += [line(x0,y0,1360,y0,'#234d70',5),line(x0,y0,x0,290,'#234d70',5)]
  for n in (0,2,4,6):
   x=x0+n*sx;parts += [line(x,y0-9,x,y0+9,'#234d70',2),txt(x,y0+44,str(n),22,'#234d70','middle')]
  for n in (0,2,4,6,8):
   y=y0-n*sy;parts += [line(x0-9,y,x0+9,y,'#234d70',2),txt(x0-30,y+8,str(n),22,'#234d70','end')]
  parts += [txt(1370,735,{'ru':'t, с','kz':'t, с','en':'t, s'}[k],25,'#234d70','middle'),txt(210,305,{'ru':'x, м','kz':'x, м','en':'x, m'}[k],25,'#234d70','middle')]
  pts={'xt':[(0,0),(2,4),(4,8),(6,8)],'slope':[(0,0),(2,4)],'reversal':[(0,0),(2,4),(4,4),(6,2)],'mastery':[(0,0),(2,2),(4,4),(6,6)]}[kind]
  coords=' '.join(f'{x0+t*sx},{y0-x*sy}' for t,x in pts)
  parts += [f'<polyline points="{coords}" fill="none" stroke="#278eaf" stroke-width="7" stroke-linejoin="round"/>']
  for t,x in pts:parts.append(circle(x0+t*sx,y0-x*sy,13,'#e68250'))
  if kind=='slope':
   coords2=' '.join(f'{x0+t*sx},{y0-x*sy}' for t,x in [(0,0),(2,2)])
   parts += [f'<polyline points="{coords2}" fill="none" stroke="#50a989" stroke-width="7"/>',circle(x0+2*sx,y0-2*sy,13,'#50a989'),txt(900,380,{'ru':'А: 2 м/с · Б: 1 м/с','kz':'А: 2 м/с · Ә: 1 м/с','en':'A: 2 m/s · B: 1 m/s'}[k],32,'#a34a43','middle')]
  else:
   subtitles={
    'xt':{'ru':'Пары (0;0), (2;4), (4;8), (6;8)','kz':'Жұптар (0;0), (2;4), (4;8), (6;8)','en':'Pairs (0,0), (2,4), (4,8), (6,8)'},
    'reversal':{'ru':'Вперёд → пауза → назад','kz':'Алға → аялдау → артқа','en':'Forward → stop → back'},
    'mastery':{'ru':'(0;0), (2;2), (4;4), (6;6)','kz':'(0;0), (2;2), (4;4), (6;6)','en':'(0,0), (2,2), (4,4), (6,6)'}}[kind][k]
   parts.append(txt(800,760 if kind in ('xt','mastery') else 260,subtitles,23,'#a34a43','middle'))
 parts.append('</svg>');return '\n'.join(parts)+'\n'

def markdown(index,row,lang):
 title,summary,scene,terms,example,task,error=row[lang];d=TEXT[lang];heads=d['heads'];slugname=slug(row)
 lines=[f'# {title}','','_Лид (summary):_ **'+summary+'**','',f'## {heads[0]}','',d['intro'],'',f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)','']
 for h,body in zip(heads[1:7],[scene,terms,example,task,error,d['check']]):lines += [f'## {h}','',body,'']
 lines += [f'## {heads[7]}','',d['note'],'']
 if index<18:
  other=MOTION[index-8]
  lines += [f'[{d["next"]}: {other[lang][0]}](/read/{slug(other)}?lang={lang})','']
 else: lines += [f'[{d["next"]}: '+{'ru':'Почему тележка меняет движение','kz':'Арба қозғалысы неге өзгереді','en':'Why a cart changes motion'}[lang]+f'](/read/physics-19-interactions-force?lang={lang})','']
 return '\n'.join(lines).rstrip()+'\n'

def build():
 assert len(MOTION)==10
 for index,row in enumerate(MOTION,9):
  assert row['stem'].startswith(f'{index:02d}-')
  for lang in ('ru','kz','en'):
   suffix='' if lang=='ru' else '-'+lang
   (LESSONS/f'{row["stem"]}{suffix}.md').write_text(markdown(index,row,lang),encoding='utf-8')
   (MAPS/f'map-{row["stem"]}-{lang}.svg').write_text(map_svg(index,row,lang),encoding='utf-8')
 print('Built 30 motion lessons and 30 exact maps')
if __name__=='__main__':build()
