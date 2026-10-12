#!/usr/bin/env python3
"""Build trilingual, concept-first forces/energy lessons and exact SVG models."""
from html import escape
try:
    from .physics_forces_data import FORCES
    from .physics_fluids_data import FLUIDS
    from .build_physics_motion import LESSONS, MAPS, txt, line, base
except ImportError:
    from physics_forces_data import FORCES
    from physics_fluids_data import FLUIDS
    from build_physics_motion import LESSONS, MAPS, txt, line, base

HEAD = {
 'ru':['История одной тележки','Что означают слова','Проверяем числа и границы','Попробуйте сами','Типичная ловушка','Опорная карта','Проверка без подсказки','Запись в журнал'],
 'kz':['Бір арбаның тарихы','Сөздердің мағынасы','Сандар мен шектерді тексерейік','Өзіңіз байқап көріңіз','Жиі кездесетін қате','Тірек картасы','Көмексіз тексеру','Журналға жазба'],
 'en':['One cart, one story','What the words mean','Check the numbers and limits','Try it yourself','A common trap','Support map','Check without hints','Notebook entry']}
INTRO = {
 'ru':'Мы продолжаем историю тележки из блока «Движение и графики». Теперь спрашиваем не только «как она двигалась?», но и «какие тела изменили её движение?». Каждый новый термин сначала свяжем с наблюдаемой сценой. На опорной карте показаны заданные учебные числа: они не заменяют ваше измерение.',
 'kz':'«Қозғалыс және графиктер» блогындағы арбаның тарихын жалғастырамыз. Енді «ол қалай қозғалды?» дегенге қоса «оның қозғалысын қандай денелер өзгертті?» деп сұраймыз. Әр жаңа сөзді алдымен көзге көрінетін көрініспен байланыстырамыз. Тірек картадағы сандар — берілген оқу мәндері, сіздің өлшеміңіз емес.',
 'en':'We continue the cart story from Motion and Graphs. Now ask not only “how did it move?” but “which objects changed its motion?” Link each new word to a visible scene before using an equation. Numbers on the support map are given teaching values, not a replacement for your own measurements.'}
CHECK = {
 'ru':'Закройте расчёт рукой. Назовите тело, к которому приложены силы, их источники и выбранное положительное направление. Потом восстановите арифметику с единицами. Смените одно исходное число и расскажите, какой вывод изменится. Если реального опыта не было, в журнале пишите «учебная модель».',
 'kz':'Есепті жауып қойыңыз. Күш түсетін денені, күштердің көзін және таңдалған оң бағытты атаңыз. Содан кейін бірліктерімен есепті қайта шығарыңыз. Бір бастапқы санды ауыстырып, қай қорытынды өзгеретінін айтыңыз. Нақты тәжірибе болмаса, журналға «оқу моделі» деп жазыңыз.',
 'en':'Cover the calculation. Name the object receiving forces, their sources, and the chosen positive direction. Rebuild the arithmetic with units. Change one given number and explain which conclusion changes. If you did not run a real trial, write “worked model” in your notebook.'}
NOTE = {
 'ru':'Запишите четыре строки: «Что я видел?», «Что было задано?», «Что вычислил?», «Какие условия нужны для такого вывода?». Так вы отделите наблюдение от модели и сможете позже проверить собственную ошибку.',
 'kz':'Төрт жол жазыңыз: «Не көрдім?», «Не берілді?», «Не есептедім?», «Бұл қорытынды үшін қандай шарт керек?». Осылайша бақылауды модельден ажыратып, кейін өз қатеңізді тексере аласыз.',
 'en':'Write four lines: “What did I see?”, “What was given?”, “What did I calculate?”, and “Which assumptions does this conclusion need?” This separates observation from model and lets you revisit an error later.'}
NEXT={'ru':'Следующий урок','kz':'Келесі сабақ','en':'Next lesson'}
# Diagram labels are hand-localized and tied to the worked examples, never guessed from body text.
LABELS={
'interaction':{'ru':('Рука','Арба','Стол','Земля','Действия разных тел'),'kz':('Қол','Арба','Үстел','Жер','Әртүрлі денелердің әсері'),'en':('Hand','Cart','Table','Earth','Actions of different objects')},
'arrows':{'ru':('2 Н вправо','3 Н влево','1 клетка = 1 Н','Длина стрелки зависит от масштаба'),'kz':('Оңға 2 Н','Солға 3 Н','1 торкөз = 1 Н','Бағыттауыш ұзындығы масштабқа тәуелді'),'en':('2 N right','3 N left','1 square = 1 N','Arrow length depends on the scale')},
'weight':{'ru':('Сила Земли: 5 Н вниз','Сила стола: 5 Н вверх','m = 0,5 кг · g ≈ 10 Н/кг','На арбе показаны две силы; вес действует на стол'),'kz':('Жер күші: төмен 5 Н','Үстел күші: жоғары 5 Н','m = 0,5 кг · g ≈ 10 Н/кг','Арбада екі күш; салмақ үстелге әсер етеді'),'en':('Earth: 5 N down','Table: 5 N up','m = 0.5 kg · g ≈ 10 N/kg','Two forces on cart; weight acts on table')},
'balance':{'ru':('5 Н вправо','5 Н влево','FΣ = 0 Н','Скорость может быть ненулевой'),'kz':('Оңға 5 Н','Солға 5 Н','FΣ = 0 Н','Жылдамдық нөлден өзге болуы мүмкін'),'en':('5 N right','5 N left','FΣ = 0 N','Velocity can be nonzero')},
'acceleration':{'ru':('FΣ = 2 Н','m = 0,5 кг','a = 4 м/с²','v: 0 → 4 → 8 м/с за 0, 1, 2 с'),'kz':('FΣ = 2 Н','m = 0,5 кг','a = 4 м/с²','v: 0 → 4 → 8 м/с, t: 0, 1, 2 с'),'en':('FΣ = 2 N','m = 0.5 kg','a = 4 m/s²','v: 0 → 4 → 8 m/s at 0, 1, 2 s')},
'friction':{'ru':('Тяга 3 Н вправо','Кедергі 1 Н влево','FΣ = 2 Н вправо','Для 0,5 кг: a = 4 м/с²'),'kz':('Тарту оңға 3 Н','Кедергі солға 1 Н','FΣ = оңға 2 Н','0,5 кг үшін: a = 4 м/с²'),'en':('Pull 3 N right','Resistance 1 N left','FΣ = 2 N right','For 0.5 kg: a = 4 m/s²')},
'work':{'ru':('Тяга 3 Н вправо','Кедергі 1 Н влево','Путь 2 м вправо','Работы: +6 Дж, −2 Дж; сумма +4 Дж'),'kz':('Тарту оңға 3 Н','Кедергі солға 1 Н','Оңға жол 2 м','Жұмыс: +6 Дж, −2 Дж; қосынды +4 Дж'),'en':('Pull 3 N right','Resistance 1 N left','Distance 2 m right','Work: +6 J, −2 J; total +4 J')},
'energy':{'ru':('Движение: 4 Дж','Высота: 4 Дж','0,5 кг; 4 м/с; 0,8 м','Два разных учебных состояния'),'kz':('Қозғалыс: 4 Дж','Биіктік: 4 Дж','0,5 кг; 4 м/с; 0,8 м','Екі бөлек оқу жағдайы'),'en':('Motion: 4 J','Height: 4 J','0.5 kg; 4 m/s; 0.8 m','Two different worked states')},
'transfer':{'ru':('До: 4 Дж','После: 0 + 3 + 1 Дж','Тележка + стол + воздух','3 и 1 Дж — условные доли, не измерения'),'kz':('Алдында: 4 Дж','Соңында: 0 + 3 + 1 Дж','Арба + үстел + ауа','3 және 1 Дж — шартты үлестер'),'en':('Before: 4 J','After: 0 + 3 + 1 J','Cart + table + air','3 and 1 J are assumed, not measured')},
'capstone':{'ru':('Тяга 3 Н','Кедергі 1 Н','Путь 2 м','FΣ=2 Н; a=4 м/с²; W=4 Дж; v=4 м/с'),'kz':('Тарту 3 Н','Кедергі 1 Н','Жол 2 м','FΣ=2 Н; a=4 м/с²; W=4 Дж; v=4 м/с'),'en':('Pull 3 N','Resistance 1 N','Distance 2 m','FΣ=2 N; a=4 m/s²; W=4 J; v=4 m/s')}}

def arrow(x1,y1,x2,y2,color):
 import math
 a=math.atan2(y2-y1,x2-x1);d=22;w=12
 p1=(x2-d*math.cos(a)+w*math.sin(a),y2-d*math.sin(a)-w*math.cos(a))
 p2=(x2-d*math.cos(a)-w*math.sin(a),y2-d*math.sin(a)+w*math.cos(a))
 return line(x1,y1,x2,y2,color,8)+f'<path d="M{x2},{y2}L{p1[0]:.1f},{p1[1]:.1f}L{p2[0]:.1f},{p2[1]:.1f}Z" fill="{color}"/>'

def diagram(row,lang):
 kind=row['kind']; title=row[lang][0]; labels=LABELS[kind][lang]
 p=base(title,{'ru':'Опорная карта · учебные данные','kz':'Тірек карта · оқу деректері','en':'Support map · worked values'}[lang])
 p += ['<rect x="185" y="275" width="800" height="425" rx="18" fill="#f2f9fc" stroke="#c8dce7" stroke-width="2"/>',
       '<rect x="1025" y="275" width="390" height="425" rx="18" fill="#fff7e9" stroke="#ecd6af" stroke-width="2"/>']
 if kind in ('energy','transfer'):
  vals=(4,4) if kind=='energy' else (4,4)
  for i,(yy,c) in enumerate(((355,'#258db9'),(525,'#e68a4b'))):
   p += [f'<rect x="305" y="{yy}" width="520" height="95" rx="12" fill="{c}"/>',txt(565,yy+61,labels[i],35,'white','middle',700)]
  if kind=='transfer':p += [arrow(845,407,845,565,'#a64c4d')]
 elif kind=='interaction':
  p += ['<rect x="460" y="445" width="240" height="130" rx="14" fill="#4b9bb9"/>',txt(580,525,labels[1],37,'white','middle'),txt(315,510,labels[0],31,'#245472','middle'),txt(580,645,labels[2],31,'#245472','middle'),txt(580,340,labels[3],31,'#245472','middle'),arrow(365,500,445,500,'#e1784a'),arrow(580,625,580,580,'#50a580'),arrow(680,355,680,435,'#a66b9d')]
 elif kind=='weight':
  p += ['<rect x="440" y="470" width="300" height="110" rx="12" fill="#4b9bb9"/>',line(315,610,860,610,'#365c72',8),arrow(525,475,525,345,'#50a580'),arrow(655,470,655,585,'#e1784a'),txt(580,665,labels[2],27,'#284960','middle')]
 else:
  p += [line(260,590,915,590,'#426179',6),'<rect x="505" y="455" width="175" height="110" rx="14" fill="#4b9bb9"/>',
        '<circle cx="540" cy="580" r="24" fill="#294964"/><circle cx="645" cy="580" r="24" fill="#294964"/>']
  if kind=='arrows':
   p += [arrow(700,420,820,420,'#258db9'),arrow(485,420,305,420,'#e1784a'),txt(580,655,labels[2],29,'#245472','middle')]
  else:
   rightlen=180
   leftlen=180 if kind=='balance' else 60
   p += [arrow(690,415,690+rightlen,415,'#258db9')]
   if kind!='acceleration':p.append(arrow(495,415,495-leftlen,415,'#e1784a'))
   if kind in ('work','capstone'):
    p += [arrow(345,640,790,640,'#50a580'),txt(580,680,labels[2],25,'#245472','middle')]
 # Three/four localized lines on right card, sized to avoid edge collisions.
 if kind=='interaction': right=[labels[4],'Рука → арба' if lang=='ru' else ('Қол → арба' if lang=='kz' else 'Hand → cart')]
 elif kind=='arrows': right=[labels[0],labels[1],labels[2],labels[3]]
 else:right=list(labels)
 for j,s in enumerate(right):
  words=s.split(); lines=[]; current=''
  for word in words:
   trial=(current+' '+word).strip()
   if len(trial)>28 and current:lines.append(current);current=word
   else:current=trial
  if current:lines.append(current)
  size=24 if len(lines)>1 else (27 if len(s)>23 else 31)
  for k,part in enumerate(lines):
   p.append(txt(1220,345+j*90+k*31,part,size,'#244a65','middle',700))
 p += ['</svg>']
 return '\n'.join(p)+'\n'

def markdown(index,row,lang):
 title,summary,scene,terms,example,task,error=row[lang]
 heads=HEAD[lang]
 parts=[f'# {title}','','_Лид (summary):_ **'+summary+'**','',f'## {heads[0]}','',INTRO[lang],'',scene,'',
        f'## {heads[1]}','',terms,'',f'## {heads[2]}','',example,'',f'## {heads[3]}','',task,'',
        f'## {heads[4]}','',error,'',f'## {heads[5]}','',f'![{title}](/static/course/physics/map-{row["stem"]}-{lang}.svg)','',
        f'## {heads[6]}','',CHECK[lang],'',f'## {heads[7]}','',NOTE[lang],'']
 if index<28:
  nxt=FORCES[index-18]
  parts += [f'[{NEXT[lang]}: {nxt[lang][0]}](/read/physics-{nxt["stem"]}?lang={lang})','']
 else:
  nxt=FLUIDS[0]
  parts += [f'[{NEXT[lang]}: {nxt[lang][0]}](/read/physics-{nxt["stem"]}?lang={lang})','']
 return '\n'.join(parts).rstrip()+'\n'

def build():
 assert len(FORCES)==10
 for i,row in enumerate(FORCES,19):
  assert row['stem'].startswith(f'{i:02d}-')
  for lang in ('ru','kz','en'):
   suffix='' if lang=='ru' else '-'+lang
   (LESSONS/f'{row["stem"]}{suffix}.md').write_text(markdown(i,row,lang),encoding='utf-8')
   (MAPS/f'map-{row["stem"]}-{lang}.svg').write_text(diagram(row,lang),encoding='utf-8')
 print('Built 30 forces and energy lessons and 30 support maps')
if __name__=='__main__':build()
