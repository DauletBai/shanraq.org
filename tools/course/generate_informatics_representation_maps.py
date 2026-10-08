#!/usr/bin/env python3
"""Render exact, localized support diagrams for Informatics lessons 11–18."""
from pathlib import Path
from html import escape

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / 'web/static/course/informatics'

def txt(x,y,s,size=31,color='#18334b',weight=500,anchor='start'):
    return f'<text x="{x}" y="{y}" fill="{color}" font-family="DejaVu Sans,Arial,sans-serif" font-size="{size}" font-weight="{weight}" text-anchor="{anchor}">{escape(str(s))}</text>'

def rect(x,y,w,h,fill='#ffffff',r=24,stroke='#d6e4ed',sw=2,shadow=False):
    effect=' filter="url(#softShadow)"' if shadow else ''
    return f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="{r}" fill="{fill}" stroke="{stroke}" stroke-width="{sw}"{effect}/>'

def line(x1,y1,x2,y2,color='#9db4c5',sw=3):
    return f'<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" stroke-width="{sw}"/>'

def panel(x,y,w,h,title):
    return rect(x,y,w,h,shadow=True)+txt(x+32,y+55,title,31,'#136789',700)

LABELS={
'ru':{
 11:('Два состояния — один бит','Один разряд','Два разряда','00','01','10','11','8 бит = 256 сочетаний'),
 12:('Вес двоичных разрядов','Разряд','Вес','Вклад','Итог: 8 + 4 + 0 + 1 = 13₁₀'),
 13:('Символ, кодовая точка и байты UTF-8','Символ','Кодовая точка','Байты UTF-8','Байтов'),
 14:('Фотография как числовая сетка','Сетка 2 × 2','Каналы RGB','Сырые данные: 2 × 2 × 3 = 12 байтов'),
 15:('Два способа измерять время','Звук: 4 отсчёта/с','Видео: 2 кадра/с','Одна секунда: [0, 1)'),
 16:('Сжатие: что можно восстановить?','Без потерь','С потерями','точно','не точно','ААААБББ','4А3Б','201','200'),
 17:('Контрольная сумма: польза и предел','Исходные данные','Изменение замечено','Коллизия: один остаток','остаток ='),
 18:('Формат помощника 0.2','UTF-8 JSON','Три задачи','Точная копия','разные id','один true','сравнить байты'),
},
'kz':{
 11:('Екі күй — бір бит','Бір разряд','Екі разряд','00','01','10','11','8 бит = 256 үйлесім'),
 12:('Екілік разряд салмағы','Разряд','Салмақ','Үлес','Нәтиже: 8 + 4 + 0 + 1 = 13₁₀'),
 13:('Таңба, код нүктесі және UTF-8 байттары','Таңба','Код нүктесі','UTF-8 байттары','Байт'),
 14:('Фотосурет сандық тор ретінде','2 × 2 тор','RGB арналары','Шикі дерек: 2 × 2 × 3 = 12 байт'),
 15:('Уақытты өлшеудің екі тәсілі','Дыбыс: 4 өлшем/с','Бейне: 2 кадр/с','Бір секунд: [0, 1)'),
 16:('Сығу: нені қайтарамыз?','Шығынсыз','Шығынды','дәл','дәл емес','ААААБББ','4А3Б','201','200'),
 17:('Бақылау қосындысының пайдасы мен шегі','Бастапқы дерек','Өзгеріс байқалды','Коллизия: бір қалдық','қалдық ='),
 18:('Көмекші пішімі 0.2','UTF-8 JSON','Үш тапсырма','Дәл көшірме','бөлек id','бір true','байтты салыстыру'),
},
'en':{
 11:('Two states make one bit','One position','Two positions','00','01','10','11','8 bits = 256 patterns'),
 12:('Binary place weights','Digit','Weight','Contribution','Result: 8 + 4 + 0 + 1 = 13₁₀'),
 13:('Character, code point, and UTF-8 bytes','Character','Code point','UTF-8 bytes','Bytes'),
 14:('A photograph as a numeric grid','2 × 2 grid','RGB channels','Raw data: 2 × 2 × 3 = 12 bytes'),
 15:('Two ways to measure time','Audio: 4 samples/s','Video: 2 frames/s','One second: [0, 1)'),
 16:('Compression: what can be restored?','Lossless','Lossy','exact','not exact','ААААБББ','4А3Б','201','200'),
 17:('Checksum: useful but limited','Original data','Change detected','Collision: same remainder','remainder ='),
 18:('Assistant format 0.2','UTF-8 JSON','Three tasks','Exact copy','distinct ids','one true','compare bytes'),
}}

def draw(n,lang):
    l=LABELS[lang][n]
    a=['<defs><filter id="softShadow" x="-15%" y="-15%" width="130%" height="140%"><feDropShadow dx="0" dy="10" stdDeviation="13" flood-color="#416d89" flood-opacity="0.18"/></filter></defs>',rect(0,0,1600,900,'#f4f9fc',0,'#f4f9fc',0),txt(80,105,f'{n:02d}  {l[0]}',47,'#163c59',700)]
    if n==11:
        a += [panel(80,170,520,545,l[1]),panel(640,170,880,545,l[2])]
        for i,(v,name) in enumerate([('0','OFF'),('1','ON')]):
            x=160+i*205
            a += [rect(x,320,155,175,'#e8f6fb' if i==0 else '#fff0d8',18),txt(x+77,430,v,76,'#186588',700,'middle'),txt(x+77,550,name,30,'#647c8d',600,'middle')]
        for i,v in enumerate(l[3:7]):
            x=720+(i%2)*350;y=300+(i//2)*170
            a += [rect(x,y,250,120,'#edf5ff',18),txt(x+125,y+80,v,62,'#186588',700,'middle')]
        a += [txt(800,810,l[7],39,'#17485f',700)]
    elif n==12:
        a += [panel(80,170,1440,540,l[1]+'  1101₂')]
        for i,(bit,weight,contrib) in enumerate(zip('1101',('8','4','2','1'),('8','4','0','1'))):
            x=180+i*330
            a += [rect(x,270,250,105,'#e3f4fb',16),txt(x+125,340,bit,62,'#165f82',700,'middle'),txt(x+125,455,weight,46,'#477086',600,'middle'),txt(x+125,570,contrib,50,'#b25033',700,'middle')]
        a += [txt(110,455,l[2],26),txt(110,570,l[3],26),txt(800,805,l[4],39,'#17485f',700,'middle')]
    elif n==13:
        a += [panel(80,175,1440,580,l[0])]
        xs=[170,500,830,1200]
        for x,label in zip(xs,l[1:]):a.append(txt(x,285,label,27,'#136789',700))
        for j,(char,point,byte,count) in enumerate((('A','U+0041','41','1'),('Я','U+042F','D0 AF','2'),('Ә','U+04D8','D3 98','2'))):
            y=385+j*125
            a += [rect(140,y-60,1330,94,'#f4f9fc',12)]
            for x,val in zip(xs,(char,point,byte,count)):a.append(txt(x,y,val,37,'#203d54',600))
    elif n==14:
        a += [panel(80,170,690,560,l[1]),panel(810,170,710,560,l[2])]
        colors=[('#f04c45','255, 0, 0'),('#37ba69','0, 255, 0'),('#427ded','0, 0, 255'),('#ffffff','255, 255, 255')]
        for i,(color,val) in enumerate(colors):
            x=180+(i%2)*230;y=280+(i//2)*195
            a += [rect(x,y,165,150,color,14,'#afbfcc',3),txt(900+(i%2)*305,345+(i//2)*195,val,27,'#234056',600)]
        a += [txt(800,810,l[3],37,'#17485f',700,'middle')]
    elif n==15:
        a += [panel(80,170,1440,245,l[1]),panel(80,445,1440,245,l[2])]
        for y,times,color in ((315,('0','0.25','0.5','0.75'),'#197db2'),(590,('0','0.5'),'#c97639')):
            a += [line(180,y,1400,y,'#a9bbc8',5)]
            for j,t in enumerate(times):
                x=210+j*(295 if len(times)==4 else 590)
                a += [f'<circle cx="{x}" cy="{y}" r="17" fill="{color}"/>',txt(x,y-42,t,30,'#234056',600,'middle')]
        a += [txt(800,795,l[3],36,'#17485f',700,'middle')]
    elif n==16:
        a += [panel(80,170,690,560,l[1]),panel(810,170,710,560,l[2])]
        a += [txt(420,365,l[5],58,'#235579',700,'middle'),line(420,408,420,470,'#417e92',7),'<path d="M403 455 L420 478 L437 455" fill="none" stroke="#417e92" stroke-width="7"/>',txt(420,555,l[6],58,'#235579',700,'middle'),txt(420,655,l[3],29,'#28704c',700,'middle')]
        a += [txt(1150,365,l[7],58,'#235579',700,'middle'),line(1150,408,1150,470,'#b46636',7),'<path d="M1133 455 L1150 478 L1167 455" fill="none" stroke="#b46636" stroke-width="7"/>',txt(1150,555,l[8],58,'#235579',700,'middle'),txt(1150,655,l[4],29,'#a64e36',700,'middle')]
        a += [txt(800,810,'4А3Б = ААААБББ  |  200 ≠ 201',34,'#17485f',600,'middle')]
    elif n==17:
        titles=l[1:4]
        examples=[('4 + 7 + 2 = 13','3'),('4 + 8 + 2 = 14','4'),('5 + 6 + 2 = 13','3')]
        for j,(title,(formula,rem)) in enumerate(zip(titles,examples)):
            x=80+j*490
            a += [panel(x,180,450,505,title),txt(x+225,395,formula,34,'#27465e',600,'middle'),txt(x+225,530,l[4]+rem,31,'#b25936' if j==1 else '#28704c',700,'middle')]
        a += [txt(800,795,'4,7,2 ≠ 5,6,2    |    3 = 3',34,'#17485f',600,'middle')]
    elif n==18:
        footer={'ru':'3 задачи  |  3 id  |  1 true','kz':'3 тапсырма  |  3 id  |  1 true','en':'3 tasks  |  3 ids  |  1 true'}[lang]
        for j,title in enumerate(l[1:4]):
            x=80+j*490
            a += [panel(x,180,450,535,title)]
        a += [txt(140,360,'version: 0.2',35),txt(140,450,'tasks: [ … ]',35),txt(140,540,'encoding: UTF-8',30),
              txt(625,335,'t-01   false',34),txt(625,425,'t-02   true',34),txt(625,515,'t-03   false',34),
              txt(1110,340,l[4],30),txt(1110,430,l[5],30),txt(1110,520,l[6],29),
              txt(800,810,footer,37,'#17485f',700,'middle')]
    return '<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900" role="img" aria-label="'+escape(l[0])+'">\n'+'\n'.join(a)+'\n</svg>\n'

def main():
    for n in range(11,19):
        for lang in ('ru','kz','en'):
            stem={11:'bits-states',12:'binary-numbers',13:'text-unicode',14:'pixels-color',15:'sound-video-sampling',16:'compression',17:'integrity-errors',18:'representation-mastery'}[n]
            (OUT/f'map-{n:02d}-{stem}-{lang}.svg').write_text(draw(n,lang),encoding='utf-8')
    print('Generated 24 localized support diagrams')

if __name__=='__main__':main()
