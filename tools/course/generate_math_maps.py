#!/usr/bin/env python3
"""Build the original, language-neutral support maps for the math pilot."""
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "web/static/course/mathematics"

HEAD = '''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 630" role="img" aria-labelledby="title desc">
<title id="title">{title}</title><desc id="desc">{desc}</desc>
<defs>
 <linearGradient id="bg" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#edf7ff"/><stop offset="1" stop-color="#fff7e5"/></linearGradient>
 <linearGradient id="red" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#ff746a"/><stop offset="1" stop-color="#b80f18"/></linearGradient>
 <linearGradient id="blue" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#6dd8f2"/><stop offset="1" stop-color="#2369b3"/></linearGradient>
 <linearGradient id="gold" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#ffe89a"/><stop offset="1" stop-color="#e39a16"/></linearGradient>
 <filter id="s" x="-30%" y="-30%" width="170%" height="180%"><feDropShadow dx="0" dy="10" stdDeviation="11" flood-color="#17324d" flood-opacity=".22"/></filter>
</defs><rect width="1200" height="630" rx="34" fill="url(#bg)"/>'''
TAIL = '</svg>\n'


def save(name, title, desc, body):
    OUT.mkdir(parents=True, exist_ok=True)
    (OUT / name).write_text(HEAD.format(title=title, desc=desc) + body + TAIL, encoding="utf-8")


def percent_grid(filled=37):
    """Draw a 10x10 grid and make every counted cell individually visible."""
    cells = []
    for index in range(100):
        classes = ["percent-cell"]
        if index < filled:
            classes.append("percent-cell--filled")
        if 30 <= index < 37:
            classes.append("percent-cell--final-seven")
        fill = "url(#gold)" if 30 <= index < 37 else "url(#red)" if index < filled else "#fff"
        cells.append(
            f'<rect class="{" ".join(classes)}" x="{index % 10 * 50}" y="{index // 10 * 50}" '
            f'width="50" height="50" fill="{fill}" stroke="#365675" stroke-width="2"/>'
        )
    return "".join(cells)


def main():
    save("map-01-diagnostic.svg", "Маршрут входной диагностики", "Ответы из задач диагностики образуют путь от арифметики к алгебре.", '''
<path d="M100 475 C290 475 280 360 455 360 S650 245 795 245 S930 135 1090 135" fill="none" stroke="#d7b06a" stroke-width="34" stroke-linecap="round" filter="url(#s)"/>
<g fill="#fff" stroke="#365675" stroke-width="6">
 <circle cx="120" cy="475" r="54"/><circle cx="350" cy="390" r="54"/><circle cx="570" cy="310" r="54"/><circle cx="790" cy="245" r="54"/><circle cx="1015" cy="145" r="54"/>
</g><g data-diagnostic-route="12|3/8|2:3|25%|x=5" font-family="system-ui,sans-serif" font-size="38" font-weight="700" text-anchor="middle" fill="#17324d">
 <text x="120" y="488">12</text><text x="350" y="403">3/8</text><text x="570" y="323">2:3</text><text x="790" y="258">25%</text><text x="1015" y="158">x=5</text>
</g><path d="M80 550h1040" stroke="#365675" stroke-width="3" stroke-dasharray="12 14" opacity=".45"/>
''')
    save("map-02-fraction.svg", "Три смысла дроби 3/8", "Три из восьми равных частей, деление трех на восемь и точка три восьмых на прямой обозначают одно число.", '''
<g data-fraction="3/8" filter="url(#s)">
 <circle cx="250" cy="300" r="132" fill="#fff"/>
 <path class="fraction-sector--filled" d="M250 300L250 168A132 132 0 0 1 343.34 393.34Z" fill="url(#red)"/>
 <circle cx="250" cy="300" r="132" fill="none" stroke="#365675" stroke-width="6"/>
 <g class="fraction-spokes" stroke="#365675" stroke-width="3"><path d="M250 300L250 168"/><path d="M250 300L343.34 206.66"/><path d="M250 300L382 300"/><path d="M250 300L343.34 393.34"/><path d="M250 300L250 432"/><path d="M250 300L156.66 393.34"/><path d="M250 300L118 300"/><path d="M250 300L156.66 206.66"/></g>
</g>
<g font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d" font-weight="750"><text x="250" y="505" font-size="48">3 из 8 = 3/8</text><text x="800" y="175" font-size="50">3 ÷ 8 = 3/8</text></g>
<g data-number-line-fraction="3/8">
 <path d="M520 350H1080" stroke="#365675" stroke-width="8" stroke-linecap="round"/>
 <path d="M520 350H730" stroke="url(#red)" stroke-width="18"/>
 <g class="eighth-ticks" stroke="#365675" stroke-width="4"><path d="M520 322V378"/><path d="M590 328V372"/><path d="M660 328V372"/><path d="M730 318V382"/><path d="M800 328V372"/><path d="M870 328V372"/><path d="M940 328V372"/><path d="M1010 328V372"/><path d="M1080 322V378"/></g>
 <circle cx="730" cy="350" r="18" fill="#b80f18"/>
 <g font-family="system-ui,sans-serif" font-size="28" font-weight="700" text-anchor="middle" fill="#17324d"><text x="520" y="425">0</text><text x="730" y="425">3/8</text><text x="1080" y="425">1</text></g>
</g>
''')
    save("map-03-equivalence.svg", "Равные дроби 1/2 и 2/4", "Одинаковая половина длины разделена сначала на две, затем на четыре равные части.", '''
<g data-equivalence="1/2=2/4" filter="url(#s)" stroke="#365675" stroke-width="5"><rect x="130" y="155" width="900" height="100" rx="15" fill="#fff"/><path d="M130 155h450v100H130Z" fill="url(#red)"/><rect x="130" y="155" width="900" height="100" rx="15" fill="none"/><path d="M580 155v100"/>
<rect x="130" y="365" width="900" height="100" rx="15" fill="#fff"/><path d="M130 365h450v100H130Z" fill="url(#red)"/><rect x="130" y="365" width="900" height="100" rx="15" fill="none"/><path d="M355 365v100M580 365v100M805 365v100"/></g>
<path d="M535 282h70M535 318h70" stroke="#e39a16" stroke-width="10" stroke-linecap="round"/>
<g font-family="system-ui,sans-serif" font-size="46" font-weight="700" text-anchor="end" fill="#17324d"><text x="110" y="220">1/2</text><text x="110" y="430">2/4</text></g>
''')
    save("map-04-compare.svg", "Сравнение 1/2 и 2/3", "Обе дроби стоят на прямой, разделенной на шесть равных интервалов: три шестых левее четырех шестых.", '''
<g data-comparison="1/2&lt;2/3">
<path d="M120 320h960" stroke="#365675" stroke-width="10" stroke-linecap="round"/>
<g class="sixth-ticks" stroke="#365675" stroke-width="5"><path d="M120 275v90"/><path d="M280 285v70"/><path d="M440 285v70"/><path d="M600 275v90"/><path d="M760 275v90"/><path d="M920 285v70"/><path d="M1080 275v90"/></g>
<circle cx="600" cy="320" r="28" fill="url(#blue)" filter="url(#s)"/><circle cx="760" cy="320" r="28" fill="url(#red)" filter="url(#s)"/>
<path d="M600 235v-70M760 235v-70" stroke="#365675" stroke-width="4"/>
<g font-family="system-ui,sans-serif" font-weight="700" text-anchor="middle" fill="#17324d"><text x="120" y="410" font-size="40">0</text><text x="1080" y="410" font-size="40">1</text><text x="600" y="145" font-size="42">1/2 = 3/6</text><text x="760" y="85" font-size="42">2/3 = 4/6</text><text x="680" y="520" font-size="48">1/2 &lt; 2/3</text></g>
</g>
''')
    save("map-05-operations.svg", "Сложение 1/2 и 1/3", "Половина и треть сначала представлены шестыми, затем три шестых и две шестых объединены в пять шестых.", '''
<g data-equation="1/2+1/3=5/6" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4">
  <rect x="70" y="105" width="360" height="75" rx="12" fill="#fff"/><path d="M70 105h180v75H70Z" fill="url(#red)"/><rect x="70" y="105" width="360" height="75" rx="12" fill="none"/><path d="M250 105v75"/>
  <rect x="70" y="255" width="360" height="75" rx="12" fill="#fff"/><path d="M70 255h120v75H70Z" fill="url(#blue)"/><rect x="70" y="255" width="360" height="75" rx="12" fill="none"/><path d="M190 255v75M310 255v75"/>
  <rect x="700" y="105" width="420" height="75" rx="12" fill="#fff"/><path d="M700 105h210v75H700Z" fill="url(#red)"/><rect x="700" y="105" width="420" height="75" rx="12" fill="none"/><path d="M770 105v75M840 105v75M910 105v75M980 105v75M1050 105v75"/>
  <rect x="700" y="255" width="420" height="75" rx="12" fill="#fff"/><path d="M700 255h140v75H700Z" fill="url(#blue)"/><rect x="700" y="255" width="420" height="75" rx="12" fill="none"/><path d="M770 255v75M840 255v75M910 255v75M980 255v75M1050 255v75"/>
  <rect x="700" y="440" width="420" height="75" rx="12" fill="#fff"/><path d="M700 440h350v75H700Z" fill="url(#gold)"/><rect x="700" y="440" width="420" height="75" rx="12" fill="none"/><path d="M770 440v75M840 440v75M910 440v75M980 440v75M1050 440v75"/>
 </g>
 <g font-size="34" font-weight="750"><text x="250" y="225">1/2</text><text x="250" y="375">1/3</text><text x="910" y="225">3/6</text><text x="910" y="375">2/6</text><text x="910" y="575">3/6 + 2/6 = 5/6</text></g>
 <g font-size="58" font-weight="800"><text x="650" y="238">+</text><text x="910" y="425">↓</text></g>
 <g fill="none" stroke="#365675" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"><path d="M480 142H650l-22-18m22 18l-22 18"/><path d="M480 292H650l-22-18m22 18l-22 18"/></g>
</g>
''')
    save("map-06-ratio.svg", "Равные отношения 2:3 и 4:6", "Две красные меры и три синие меры увеличены вдвое: четыре красные и шесть синих.", '''
<g data-ratios="2:3=4:6" filter="url(#s)"><g fill="url(#red)"><circle cx="220" cy="220" r="55"/><circle cx="220" cy="410" r="55"/></g><g fill="url(#blue)"><circle cx="390" cy="170" r="55"/><circle cx="390" cy="315" r="55"/><circle cx="390" cy="460" r="55"/></g>
<path d="M510 315h100" stroke="#365675" stroke-width="12" stroke-linecap="round"/><path d="M590 280l40 35-40 35" fill="none" stroke="#365675" stroke-width="12" stroke-linecap="round" stroke-linejoin="round"/>
<g fill="url(#red)"><circle cx="745" cy="170" r="42"/><circle cx="745" cy="270" r="42"/><circle cx="745" cy="370" r="42"/><circle cx="745" cy="470" r="42"/></g><g fill="url(#blue)"><circle cx="920" cy="115" r="42"/><circle cx="920" cy="195" r="42"/><circle cx="920" cy="275" r="42"/><circle cx="920" cy="355" r="42"/><circle cx="920" cy="435" r="42"/><circle cx="920" cy="515" r="42"/></g></g>
<g font-family="system-ui,sans-serif" font-weight="700" text-anchor="middle" fill="#17324d"><text x="570" y="255" font-size="38">×2</text><text x="305" y="585" font-size="52">2 : 3</text><text x="835" y="585" font-size="52">4 : 6</text></g>
''')
    save("map-07-percent.svg", "37 клеток из 100", "Ровно 37 из 100 отдельных клеток закрашены: тридцать плюс семь.", '''
<g data-fraction="37/100" transform="translate(115 65)" filter="url(#s)">''' + percent_grid(37) + '''<rect width="500" height="500" rx="12" fill="none" stroke="#365675" stroke-width="6"/></g>
<path d="M700 315h120" stroke="#365675" stroke-width="12"/><path d="M800 280l40 35-40 35" fill="none" stroke="#365675" stroke-width="12" stroke-linecap="round"/>
<g font-family="system-ui,sans-serif" font-weight="700" text-anchor="middle" fill="#17324d"><text x="950" y="205" font-size="42">30 + 7 = 37</text><text x="950" y="290" font-size="58">37/100</text><text x="950" y="375" font-size="58">0.37</text><text x="950" y="460" font-size="58">37%</text></g>
''')
    save("map-08-proportion.svg", "Пропорция 2/3 = 6/9", "Обе части отношения увеличены в три раза, поэтому две дроби сохраняют одно значение.", '''
<g data-equivalence="2/3=6/9">
<path d="M600 130v350M300 480h600" stroke="#365675" stroke-width="18" stroke-linecap="round" filter="url(#s)"/><path d="M255 250h690" stroke="#365675" stroke-width="14"/><path d="M350 250l-90 210h180zM850 250l-90 210h180z" fill="#fff" stroke="#365675" stroke-width="7"/>
<g font-family="system-ui,sans-serif" font-weight="700" text-anchor="middle" fill="#17324d"><text x="350" y="390" font-size="64">2/3</text><text x="850" y="390" font-size="64">6/9</text><text x="600" y="575" font-size="38">2 × 3 = 6; 3 × 3 = 9</text></g><circle cx="600" cy="250" r="34" fill="url(#gold)" stroke="#365675" stroke-width="5"/>
</g>
''')
    save("map-09-mastery.svg", "Пять сторон мастерства", "Модель, вычисление, объяснение, исправление ошибки и перенос проверяются отдельно; цель сейчас восемь из десяти, через семь дней семь из десяти.", '''
<g data-mastery-skills="5" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <path d="M135 205H1055" stroke="#d7b06a" stroke-width="24" stroke-linecap="round" filter="url(#s)"/>
 <g fill="#fff" stroke="#365675" stroke-width="5" filter="url(#s)">
  <rect class="mastery-skill" x="45" y="95" width="180" height="220" rx="24"/><rect class="mastery-skill" x="275" y="95" width="180" height="220" rx="24"/><rect class="mastery-skill" x="505" y="95" width="180" height="220" rx="24"/><rect class="mastery-skill" x="735" y="95" width="180" height="220" rx="24"/><rect class="mastery-skill" x="965" y="95" width="180" height="220" rx="24"/>
 </g>
 <g fill="url(#red)" stroke="#fff" stroke-width="3"><circle cx="75" cy="125" r="21"/><circle cx="305" cy="125" r="21"/><circle cx="535" cy="125" r="21"/><circle cx="765" cy="125" r="21"/><circle cx="995" cy="125" r="21"/></g>
 <g fill="none" stroke="#17324d" stroke-width="4"><rect x="117" y="161" width="36" height="36" rx="3"/><path d="M135 161v36M117 179h36"/></g>
 <g font-size="38" font-weight="800"><text x="365" y="198">=</text><text x="595" y="198">?</text><text x="825" y="198">↻</text><text x="1055" y="198">↗</text></g>
 <g font-size="22" font-weight="750"><text x="135" y="260">Модель</text><text x="365" y="260">Расчёт</text><text x="595" y="260">Объяснение</text><text x="825" y="248"><tspan x="825">Исправление</tspan><tspan x="825" dy="28">ошибки</tspan></text><text x="1055" y="260">Перенос</text></g>
 <g font-size="18" font-weight="800" fill="#fff"><text x="75" y="132">1</text><text x="305" y="132">2</text><text x="535" y="132">3</text><text x="765" y="132">4</text><text x="995" y="132">5</text></g>
 <g filter="url(#s)"><rect x="245" y="405" width="310" height="130" rx="26" fill="#fff" stroke="#365675" stroke-width="5"/><rect x="645" y="405" width="310" height="130" rx="26" fill="#fff" stroke="#365675" stroke-width="5"/></g>
 <g font-weight="750"><text x="400" y="450" font-size="22">ПЕРВАЯ ПОПЫТКА</text><text x="400" y="505" font-size="48">8/10</text><text x="800" y="450" font-size="22">ЧЕРЕЗ 7 ДНЕЙ</text><text x="800" y="505" font-size="48">7/10</text></g>
 <path d="M575 470H625l-18-16m18 16l-18 16" fill="none" stroke="#b80f18" stroke-width="7" stroke-linecap="round" stroke-linejoin="round"/>
</g>
''')
    save("map-10-quantity-counting.svg", "Количество через равные группы", "Четыре полных лотка по шесть предметов и еще три предмета дают ровно двадцать семь.", '''
<g data-counting="4x6+3=27" font-family="system-ui,sans-serif" fill="#17324d" text-anchor="middle">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff">
  <rect x="55" y="105" width="215" height="300" rx="24"/><rect x="290" y="105" width="215" height="300" rx="24"/><rect x="525" y="105" width="215" height="300" rx="24"/><rect x="760" y="105" width="215" height="300" rx="24"/>
 </g>
 <g fill="url(#blue)">
  <circle cx="110" cy="175" r="25"/><circle cx="162" cy="175" r="25"/><circle cx="214" cy="175" r="25"/><circle cx="110" cy="250" r="25"/><circle cx="162" cy="250" r="25"/><circle cx="214" cy="250" r="25"/>
  <circle cx="345" cy="175" r="25"/><circle cx="397" cy="175" r="25"/><circle cx="449" cy="175" r="25"/><circle cx="345" cy="250" r="25"/><circle cx="397" cy="250" r="25"/><circle cx="449" cy="250" r="25"/>
  <circle cx="580" cy="175" r="25"/><circle cx="632" cy="175" r="25"/><circle cx="684" cy="175" r="25"/><circle cx="580" cy="250" r="25"/><circle cx="632" cy="250" r="25"/><circle cx="684" cy="250" r="25"/>
  <circle cx="815" cy="175" r="25"/><circle cx="867" cy="175" r="25"/><circle cx="919" cy="175" r="25"/><circle cx="815" cy="250" r="25"/><circle cx="867" cy="250" r="25"/><circle cx="919" cy="250" r="25"/>
 </g>
 <g fill="url(#red)" filter="url(#s)"><circle cx="1035" cy="175" r="27"/><circle cx="1100" cy="250" r="27"/><circle cx="1035" cy="325" r="27"/></g>
 <g font-size="25" font-weight="750"><text x="162" y="365">6</text><text x="397" y="365">6</text><text x="632" y="365">6</text><text x="867" y="365">6</text><text x="1068" y="390">ещё 3</text></g>
 <text x="600" y="505" font-size="50" font-weight="850">4 × 6 + 3 = 27</text>
 <text x="600" y="560" font-size="25" font-weight="650" fill="#5b6b79">группа × размер группы + остаток</text>
</g>
''')
    save("map-11-place-value.svg", "Разрядная запись числа 4 072", "Четыре тысячи, ноль сотен, семь десятков и две единицы образуют число четыре тысячи семьдесят два.", '''
<g data-place-value="4072=4000+70+2" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4">
  <rect x="55" y="120" width="260" height="300" rx="28" fill="#fff"/><rect x="330" y="120" width="260" height="300" rx="28" fill="#fff"/><rect x="605" y="120" width="260" height="300" rx="28" fill="#fff"/><rect x="880" y="120" width="260" height="300" rx="28" fill="#fff"/>
  <path d="M70 405h230v-40H70z" fill="url(#blue)"/><path d="M85 365h200v-40H85z" fill="url(#blue)"/><path d="M100 325h170v-40H100z" fill="url(#blue)"/><path d="M115 285h140v-40H115z" fill="url(#blue)"/>
  <rect x="640" y="365" width="190" height="26" rx="7" fill="url(#gold)"/><rect x="640" y="326" width="190" height="26" rx="7" fill="url(#gold)"/><rect x="640" y="287" width="190" height="26" rx="7" fill="url(#gold)"/><rect x="640" y="248" width="190" height="26" rx="7" fill="url(#gold)"/><rect x="640" y="209" width="190" height="26" rx="7" fill="url(#gold)"/><rect x="640" y="170" width="190" height="26" rx="7" fill="url(#gold)"/><rect x="640" y="131" width="190" height="26" rx="7" fill="url(#gold)"/>
  <rect x="955" y="335" width="65" height="65" rx="10" fill="url(#red)"/><rect x="1030" y="335" width="65" height="65" rx="10" fill="url(#red)"/>
 </g>
 <g font-size="23" font-weight="750"><text x="185" y="90">ТЫСЯЧИ</text><text x="460" y="90">СОТНИ</text><text x="735" y="90">ДЕСЯТКИ</text><text x="1010" y="90">ЕДИНИЦЫ</text></g>
 <g font-size="76" font-weight="850"><text x="185" y="505">4</text><text x="460" y="505">0</text><text x="735" y="505">7</text><text x="1010" y="505">2</text></g>
 <text x="600" y="580" font-size="34" font-weight="750">4 072 = 4 000 + 0 + 70 + 2</text>
</g>
''')
    save("map-12-addition-subtraction.svg", "Сложение и вычитание как прямой и обратный путь", "К числу двести шестьдесят восемь прибавляют сто пятьдесят семь и получают четыреста двадцать пять; вычитание возвращает начало.", '''
<g data-family="268+157=425" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <path d="M180 315C300 120 760 120 1015 280" fill="none" stroke="url(#red)" stroke-width="24" stroke-linecap="round" filter="url(#s)"/>
 <path d="M1015 350C760 515 300 515 180 345" fill="none" stroke="url(#blue)" stroke-width="24" stroke-linecap="round" filter="url(#s)"/>
 <path d="M992 247l50 49-67 13M204 379l-51-49 67-13" fill="none" stroke="#fff" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"/>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><circle cx="165" cy="330" r="92"/><circle cx="1035" cy="330" r="92"/></g>
 <g font-size="58" font-weight="850"><text x="165" y="348">268</text><text x="1035" y="348">425</text></g>
 <g font-size="33" font-weight="800"><text x="600" y="150">+157</text><text x="600" y="525">−157</text></g>
 <g font-size="27" font-weight="700"><text x="600" y="265">часть + часть = целое</text><text x="600" y="370">целое − часть = другая часть</text></g>
</g>
''')
    save("map-13-multiplication-division.svg", "Четыре ряда по шесть", "Прямоугольный массив содержит четыре ряда по шесть точек: всего двадцать четыре, а два деления возвращают неизвестный множитель.", '''
<g data-array="4x6=24" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g transform="translate(90 95)" filter="url(#s)">
  <rect width="510" height="390" rx="28" fill="#fff" stroke="#365675" stroke-width="5"/>
  <g fill="url(#blue)">
   <circle cx="65" cy="65" r="24"/><circle cx="140" cy="65" r="24"/><circle cx="215" cy="65" r="24"/><circle cx="290" cy="65" r="24"/><circle cx="365" cy="65" r="24"/><circle cx="440" cy="65" r="24"/>
   <circle cx="65" cy="150" r="24"/><circle cx="140" cy="150" r="24"/><circle cx="215" cy="150" r="24"/><circle cx="290" cy="150" r="24"/><circle cx="365" cy="150" r="24"/><circle cx="440" cy="150" r="24"/>
   <circle cx="65" cy="235" r="24"/><circle cx="140" cy="235" r="24"/><circle cx="215" cy="235" r="24"/><circle cx="290" cy="235" r="24"/><circle cx="365" cy="235" r="24"/><circle cx="440" cy="235" r="24"/>
   <circle cx="65" cy="320" r="24"/><circle cx="140" cy="320" r="24"/><circle cx="215" cy="320" r="24"/><circle cx="290" cy="320" r="24"/><circle cx="365" cy="320" r="24"/><circle cx="440" cy="320" r="24"/>
  </g>
 </g>
 <g filter="url(#s)" fill="#fff" stroke="#365675" stroke-width="4"><rect x="690" y="85" width="430" height="120" rx="24"/><rect x="690" y="255" width="430" height="120" rx="24"/><rect x="690" y="425" width="430" height="120" rx="24"/></g>
 <g font-weight="820"><text x="905" y="157" font-size="44">4 × 6 = 24</text><text x="905" y="327" font-size="44">24 ÷ 6 = 4</text><text x="905" y="497" font-size="44">24 ÷ 4 = 6</text></g>
 <g font-size="20" font-weight="650" fill="#5b6b79"><text x="905" y="188">группы × в группе = всего</text><text x="905" y="358">сколько групп?</text><text x="905" y="528">сколько в группе?</text></g>
</g>
''')
    save("map-14-order-estimation.svg", "Дерево выражения и оценка", "В выражении двести сорок минус шесть умножить на сумму восемнадцати и семи сначала получают двадцать пять, затем сто пятьдесят и итог девяносто.", '''
<g data-expression="240-6*(18+7)=90" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g stroke="#365675" stroke-width="7" fill="none"><path d="M600 135L380 270M600 135L820 270M820 270L710 410M820 270L930 410"/></g>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="430" y="65" width="340" height="105" rx="24"/><rect x="265" y="220" width="230" height="105" rx="24"/><rect x="705" y="220" width="230" height="105" rx="24"/><rect x="610" y="380" width="200" height="105" rx="24"/><rect x="840" y="380" width="180" height="105" rx="24"/></g>
 <g font-size="34" font-weight="820"><text x="600" y="130">240 − 150 = 90</text><text x="380" y="285">240</text><text x="820" y="285">6 × 25 = 150</text><text x="710" y="445">6</text></g><text x="930" y="445" font-size="29" font-weight="820">18 + 7 = 25</text>
 <g filter="url(#s)"><rect x="105" y="500" width="990" height="85" rx="22" fill="url(#gold)" stroke="#365675" stroke-width="4"/></g>
 <text x="600" y="555" font-size="31" font-weight="800">оценка: 240 − примерно 150 ≈ 90</text>
</g>
''')
    save("map-15-negative-numbers.svg", "От минус трех к четырем", "На числовой прямой семь шагов вправо от минус трех проходят через ноль и заканчиваются в точке четыре.", '''
<g data-integers="-3+7=4" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <path d="M90 330H1110" stroke="#365675" stroke-width="10" stroke-linecap="round"/>
 <g stroke="#365675" stroke-width="5"><path d="M120 290v80"/><path d="M200 300v60"/><path d="M280 300v60"/><path d="M360 285v90"/><path d="M440 300v60"/><path d="M520 300v60"/><path d="M600 280v100"/><path d="M680 300v60"/><path d="M760 300v60"/><path d="M840 300v60"/><path d="M920 285v90"/><path d="M1000 300v60"/><path d="M1080 290v80"/></g>
 <g font-size="24" font-weight="700"><text x="120" y="415">−6</text><text x="200" y="415">−5</text><text x="280" y="415">−4</text><text x="360" y="415">−3</text><text x="440" y="415">−2</text><text x="520" y="415">−1</text><text x="600" y="415">0</text><text x="680" y="415">1</text><text x="760" y="415">2</text><text x="840" y="415">3</text><text x="920" y="415">4</text><text x="1000" y="415">5</text><text x="1080" y="415">6</text></g>
 <path d="M360 245C485 90 795 90 920 245" fill="none" stroke="url(#red)" stroke-width="18" stroke-linecap="round" filter="url(#s)"/>
 <path d="M888 214l42 30-48 19" fill="none" stroke="#fff" stroke-width="7" stroke-linecap="round" stroke-linejoin="round"/>
 <circle cx="360" cy="330" r="21" fill="url(#blue)"/><circle cx="920" cy="330" r="21" fill="url(#red)"/>
 <text x="640" y="120" font-size="38" font-weight="850">+7 шагов вправо</text><text x="640" y="535" font-size="52" font-weight="850">−3 + 7 = 4</text>
</g>
''')
    save("map-16-divisibility-primes.svg", "Простые множители числа 84", "Дерево разбирает восемьдесят четыре на два, два, три и семь; произведение простых листьев возвращает исходное число.", '''
<g data-factorization="84=2^2*3*7" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g stroke="#365675" stroke-width="6" fill="none"><path d="M600 130L390 250M600 130L810 250M390 250L280 390M390 250L500 390M810 250L700 390M810 250L920 390"/></g>
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><circle cx="600" cy="105" r="64"/><circle cx="390" cy="250" r="58"/><circle cx="810" cy="250" r="58"/></g>
 <g filter="url(#s)" stroke="#fff" stroke-width="4"><circle cx="280" cy="410" r="56" fill="url(#red)"/><circle cx="500" cy="410" r="56" fill="url(#red)"/><circle cx="700" cy="410" r="56" fill="url(#blue)"/><circle cx="920" cy="410" r="56" fill="url(#gold)"/></g>
 <g font-size="42" font-weight="850"><text x="600" y="120">84</text><text x="390" y="265">4</text><text x="810" y="265">21</text></g>
 <g font-size="38" font-weight="850" fill="#fff"><text x="280" y="425">2</text><text x="500" y="425">2</text><text x="700" y="425">3</text><text x="920" y="425">7</text></g>
 <text x="600" y="555" font-size="43" font-weight="850">84 = 2 × 2 × 3 × 7 = 2² × 3 × 7</text>
</g>
''')
    save("map-17-decimal-fractions.svg", "Три восьмых как десятичная дробь", "Умножение числителя и знаменателя на сто двадцать пять превращает три восьмых в триста семьдесят пять тысячных и запись ноль целых триста семьдесят пять тысячных.", '''
<g data-decimal="3/8=375/1000=0.375" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="30" y="90" width="270" height="180" rx="28"/><rect x="455" y="90" width="290" height="180" rx="28"/><rect x="900" y="90" width="270" height="180" rx="28"/></g>
 <g font-size="52" font-weight="850"><text x="165" y="190">3/8</text><text x="1035" y="190">0,375</text></g><text x="600" y="190" font-size="48" font-weight="850">375/1000</text>
 <g fill="none" stroke="#b80f18" stroke-width="9" stroke-linecap="round" stroke-linejoin="round"><path d="M320 180H435l-18-16m18 16l-18 16"/><path d="M765 180H880l-18-16m18 16l-18 16"/></g>
 <g filter="url(#s)" fill="#fff" stroke="#365675" stroke-width="3"><rect x="335" y="105" width="86" height="40" rx="18"/><rect x="767" y="105" width="112" height="40" rx="18"/></g><g font-size="18" font-weight="800"><text x="378" y="132">×125</text><text x="823" y="132">тысячные</text></g>
 <defs><clipPath id="decimal-bar-clip"><rect width="960" height="92" rx="18"/></clipPath></defs><g transform="translate(120 365)" filter="url(#s)" stroke="#365675" stroke-width="4"><g clip-path="url(#decimal-bar-clip)"><rect width="960" height="92" fill="#fff"/><rect width="360" height="92" fill="url(#red)"/></g><rect width="960" height="92" rx="18" fill="none"/><path d="M120 0v92M240 0v92M360 0v92M480 0v92M600 0v92M720 0v92M840 0v92"/></g>
 <g font-size="24" font-weight="750"><text x="120" y="500">0</text><text x="480" y="500">3/8</text><text x="1080" y="500">1</text></g>
 <text x="600" y="575" font-size="31" font-weight="800">0 единиц | 3 десятых | 7 сотых | 5 тысячных</text>
</g>
''')
    save("map-18-foundations-mastery.svg", "Восемь опор числового фундамента", "Количество, разряды, сложение, умножение, порядок, отрицательные числа, делимость и десятичные дроби ведут к проверке восемь из десяти.", '''
<g data-foundations-skills="8" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff">
  <rect class="foundation-skill" x="40" y="65" width="255" height="115" rx="22"/><rect class="foundation-skill" x="325" y="65" width="255" height="115" rx="22"/><rect class="foundation-skill" x="610" y="65" width="255" height="115" rx="22"/><rect class="foundation-skill" x="895" y="65" width="255" height="115" rx="22"/>
  <rect class="foundation-skill" x="40" y="225" width="255" height="115" rx="22"/><rect class="foundation-skill" x="325" y="225" width="255" height="115" rx="22"/><rect class="foundation-skill" x="610" y="225" width="255" height="115" rx="22"/><rect class="foundation-skill" x="895" y="225" width="255" height="115" rx="22"/>
 </g>
 <g font-size="23" font-weight="780"><text x="167" y="132">Количество</text><text x="452" y="132">Разряды</text><text x="737" y="120"><tspan x="737">Сложение</tspan><tspan x="737" dy="29">и вычитание</tspan></text><text x="1022" y="120"><tspan x="1022">Умножение</tspan><tspan x="1022" dy="29">и деление</tspan></text><text x="167" y="280"><tspan x="167">Порядок</tspan><tspan x="167" dy="29">и оценка</tspan></text><text x="452" y="280"><tspan x="452">Отрицательные</tspan><tspan x="452" dy="29">числа</tspan></text><text x="737" y="292">Делимость</text><text x="1022" y="280"><tspan x="1022">Десятичные</tspan><tspan x="1022" dy="29">дроби</tspan></text></g>
 <path d="M170 365C260 430 390 445 600 445S940 430 1030 365" fill="none" stroke="#d7b06a" stroke-width="22" stroke-linecap="round" filter="url(#s)"/>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="url(#gold)"><rect x="230" y="415" width="340" height="145" rx="30"/><rect x="630" y="415" width="340" height="145" rx="30"/></g>
 <text x="400" y="470" font-size="24" font-weight="780">ПЕРВАЯ ПОПЫТКА</text><text x="400" y="530" font-size="48" font-weight="880">8/10</text><text x="800" y="470" font-size="24" font-weight="780">ЧЕРЕЗ 7 ДНЕЙ</text><text x="800" y="530" font-size="48" font-weight="880">7/10</text>
</g>
''')
    save("map-19-variables.svg", "Переменная хранит значение величины", "Цена поездки состоит из семисот тенге за посадку и ста двадцати тенге за каждый из d километров; при d равном пяти цена равна тысяче тремстам тенге.", '''
<g data-variable="C=700+120d;d=5;C=1300" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="45" y="100" width="270" height="245" rx="28"/><rect x="365" y="100" width="430" height="245" rx="28"/><rect x="845" y="100" width="310" height="245" rx="28"/></g>
 <g font-weight="850"><text x="180" y="205" font-size="58">700</text><text x="580" y="205" font-size="58">120 × d</text><text x="1000" y="205" font-size="58">C</text></g>
 <g font-size="22" font-weight="720" fill="#5b6b79"><text x="180" y="263"><tspan x="180">плата</tspan><tspan x="180" dy="30">за посадку</tspan></text><text x="580" y="263"><tspan x="580">120 тенге за каждый</tspan><tspan x="580" dy="30">километр расстояния d</tspan></text><text x="1000" y="263"><tspan x="1000">общая цена</tspan><tspan x="1000" dy="30">в тенге</tspan></text></g>
 <g font-size="48" font-weight="850"><text x="340" y="225">+</text><text x="820" y="225">=</text></g>
 <g filter="url(#s)"><rect x="145" y="420" width="910" height="125" rx="28" fill="url(#gold)" stroke="#365675" stroke-width="5"/></g>
 <text x="600" y="475" font-size="29" font-weight="760">подстановка: d = 5 км</text><text x="600" y="522" font-size="39" font-weight="860">C = 700 + 120 × 5 = 1300 тенге</text>
</g>
''')
    save("map-20-expressions.svg", "Равные пути преобразования выражения", "Выражение три умножить на два икс плюс пять минус четыре икс преобразуется в шесть икс плюс пятнадцать минус четыре икс, а затем в два икс плюс пятнадцать; при икс равном четырем оба пути дают двадцать три.", '''
<g data-expression="3(2x+5)-4x=2x+15;x=4;value=23" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="115" y="55" width="970" height="105" rx="26"/><rect x="115" y="230" width="970" height="105" rx="26"/><rect x="325" y="405" width="550" height="105" rx="26"/></g>
 <g font-size="40" font-weight="850"><text x="600" y="120">3(2x + 5) − 4x</text><text x="600" y="295">6x + 15 − 4x</text><text x="600" y="470">2x + 15</text></g>
 <g fill="none" stroke="#b80f18" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"><path d="M600 175v35m-15-17 15 17 15-17"/><path d="M600 350v35m-15-17 15 17 15-17"/></g>
 <g font-size="20" font-weight="760" fill="#5b6b79"><text x="945" y="202">раскрываем скобки</text><text x="950" y="377">собираем подобные</text></g>
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="url(#gold)"><rect x="85" y="535" width="460" height="65" rx="20"/><rect x="655" y="535" width="460" height="65" rx="20"/></g>
 <g font-size="23" font-weight="800"><text x="315" y="576">x = 4: исходное = 23</text><text x="885" y="576">x = 4: новое = 23</text></g>
</g>
''')
    save("map-21-equations-basic.svg", "Уравнение сохраняет равновесие", "Из уравнения три икс плюс пять равно двадцати шести одинаковое вычитание пяти и деление на три приводят к икс равному семи; подстановка подтверждает равенство.", '''
<g data-equation="3x+5=26;x=7" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <path d="M600 105v255M390 205h420M455 205l-90 145M745 205l90 145" fill="none" stroke="#365675" stroke-width="11" stroke-linecap="round"/>
 <path d="M265 350h200M735 350h200" stroke="#365675" stroke-width="9" stroke-linecap="round"/>
 <circle cx="600" cy="105" r="34" fill="url(#red)" stroke="#fff" stroke-width="5" filter="url(#s)"/>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="235" y="250" width="260" height="90" rx="22"/><rect x="705" y="250" width="260" height="90" rx="22"/></g>
 <g font-size="40" font-weight="850"><text x="365" y="307">3x + 5</text><text x="835" y="307">26</text></g>
 <g fill="none" stroke="#b80f18" stroke-width="7" stroke-linecap="round" stroke-linejoin="round"><path d="M250 410H950"/><path d="M930 395l20 15-20 15"/></g>
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect x="70" y="450" width="315" height="105" rx="23"/><rect x="442" y="450" width="315" height="105" rx="23"/><rect x="815" y="450" width="315" height="105" rx="23"/></g>
 <g font-size="27" font-weight="820"><text x="227" y="493">−5 с обеих сторон</text><text x="227" y="530">3x = 21</text><text x="600" y="493">÷3 с обеих сторон</text><text x="600" y="530">x = 7</text><text x="972" y="493">проверка</text><text x="972" y="530">3 × 7 + 5 = 26</text></g>
</g>
''')
    save("map-22-coordinates.svg", "Адреса точек на координатной плоскости", "Точки A минус три два и B четыре два лежат на одной горизонтали, а расстояние между ними равно семи единицам.", '''
<g data-coordinates="A(-3,2);B(4,2);distance=7" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g stroke="#b8ccda" stroke-width="2" opacity=".75"><path d="M150 80v450M250 80v450M350 80v450M450 80v450M550 80v450M650 80v450M750 80v450M850 80v450M950 80v450M1050 80v450"/><path d="M150 80h900M150 180h900M150 280h900M150 380h900M150 480h900"/></g>
 <g stroke="#365675" stroke-width="7" stroke-linecap="round"><path d="M115 380h995"/><path d="M550 555V50"/></g><g fill="#365675"><path d="M1110 380l-24-14v28z"/><path d="M550 50l-14 24h28z"/></g>
 <g font-size="20" font-weight="700"><text x="150" y="415">−4</text><text x="250" y="415">−3</text><text x="350" y="415">−2</text><text x="450" y="415">−1</text><text x="650" y="415">1</text><text x="750" y="415">2</text><text x="850" y="415">3</text><text x="950" y="415">4</text><text x="1050" y="415">5</text><text x="522" y="487">−1</text><text x="522" y="287">1</text><text x="522" y="187">2</text><text x="522" y="87">3</text><text x="1095" y="410">x</text><text x="580" y="68">y</text><text x="525" y="410">0</text></g>
 <path d="M250 180H950" stroke="url(#gold)" stroke-width="18" stroke-linecap="round" filter="url(#s)"/>
 <g filter="url(#s)" stroke="#fff" stroke-width="5"><circle cx="250" cy="180" r="22" fill="url(#red)"/><circle cx="950" cy="180" r="22" fill="url(#blue)"/></g>
 <g font-size="26" font-weight="830"><text x="250" y="140">A(−3; 2)</text><text x="950" y="140">B(4; 2)</text></g>
 <path d="M250 545v28M950 545v28M250 560h700" stroke="#b80f18" stroke-width="5"/><text x="600" y="603" font-size="25" font-weight="820">|4 − (−3)| = 7 единиц</text>
</g>
''')
    save("map-23-inequalities.svg", "Неравенство задает луч решений", "Три икс плюс два не больше четырнадцати преобразуется в икс не больше четырех; закрашенная граница четыре и луч влево показывают все решения.", '''
<g data-inequality="3x+2&lt;=14;x&lt;=4" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect x="45" y="70" width="330" height="110" rx="25"/><rect x="435" y="70" width="330" height="110" rx="25"/><rect x="825" y="70" width="330" height="110" rx="25"/></g>
 <g font-size="35" font-weight="850"><text x="210" y="138">3x + 2 ≤ 14</text><text x="600" y="138">3x ≤ 12</text><text x="990" y="138">x ≤ 4</text></g>
 <g fill="none" stroke="#b80f18" stroke-width="7" stroke-linecap="round" stroke-linejoin="round"><path d="M390 125h25m-15-14 15 14-15 14"/><path d="M780 125h25m-15-14 15 14-15 14"/></g>
 <path d="M120 395H1080" stroke="#365675" stroke-width="9" stroke-linecap="round"/>
 <g stroke="#365675" stroke-width="4"><path d="M180 365v60"/><path d="M300 365v60"/><path d="M420 365v60"/><path d="M540 365v60"/><path d="M660 365v60"/><path d="M780 350v90"/><path d="M900 365v60"/><path d="M1020 365v60"/></g>
 <g font-size="23" font-weight="720"><text x="180" y="465">−1</text><text x="300" y="465">0</text><text x="420" y="465">1</text><text x="540" y="465">2</text><text x="660" y="465">3</text><text x="780" y="465">4</text><text x="900" y="465">5</text><text x="1020" y="465">6</text></g>
 <path d="M145 395H780" stroke="#c9151e" stroke-width="22" stroke-linecap="round"/><path d="M145 395l38-25v50z" fill="#b80f18"/><circle cx="780" cy="395" r="27" fill="url(#red)" stroke="#fff" stroke-width="6"/>
 <text x="462" y="330" font-size="27" font-weight="800">все числа до 4 включительно</text>
 <g filter="url(#s)"><rect x="275" y="515" width="650" height="75" rx="22" fill="url(#gold)" stroke="#365675" stroke-width="4"/></g><text x="600" y="563" font-size="28" font-weight="820">4 подходит: 3 × 4 + 2 = 14</text>
</g>
''')
    save("map-24-prealgebra-mastery.svg", "Пять опор предалгебры", "Переменные, выражения, уравнения, координаты и неравенства ведут к проверке восемь из десяти сейчас и семь из десяти через семь дней.", '''
<g data-prealgebra-skills="5" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect class="prealgebra-skill" x="60" y="65" width="320" height="105" rx="23"/><rect class="prealgebra-skill" x="440" y="65" width="320" height="105" rx="23"/><rect class="prealgebra-skill" x="820" y="65" width="320" height="105" rx="23"/><rect class="prealgebra-skill" x="250" y="220" width="320" height="105" rx="23"/><rect class="prealgebra-skill" x="630" y="220" width="320" height="105" rx="23"/></g>
 <g font-size="25" font-weight="800"><text x="220" y="130">Переменные</text><text x="600" y="130">Выражения</text><text x="980" y="130">Уравнения</text><text x="410" y="285">Координаты</text><text x="790" y="285">Неравенства</text></g>
 <path d="M220 185C280 365 470 390 600 405M600 185V405M980 185C920 365 730 390 600 405M410 340C455 380 510 395 600 405M790 340C745 380 690 395 600 405" fill="none" stroke="#d7b06a" stroke-width="12" stroke-linecap="round"/>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="url(#gold)"><rect x="160" y="420" width="385" height="145" rx="28"/><rect x="655" y="420" width="385" height="145" rx="28"/></g>
 <text x="352" y="470" font-size="23" font-weight="780">ПЕРВАЯ ПОПЫТКА</text><text x="352" y="530" font-size="49" font-weight="880">8/10</text><text x="847" y="470" font-size="23" font-weight="780">ЧЕРЕЗ 7 ДНЕЙ</text><text x="847" y="530" font-size="49" font-weight="880">7/10</text>
</g>
''')
    save("map-25-linear-functions.svg", "Формула, таблица и график линейной функции", "Для функции игрек равно два икс плюс один таблица точек минус один минус один, ноль один и два пять лежит на одной прямой.", '''
<g data-linear="y=2x+1;(-1,-1);(0,1);(2,5)" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect x="45" y="65" width="320" height="500" rx="26"/><rect x="405" y="65" width="750" height="500" rx="26"/></g>
 <text x="205" y="125" font-size="34" font-weight="850">y = 2x + 1</text><text x="205" y="170" font-size="20" font-weight="720">старт b = 1 · шаг k = 2</text>
 <g font-size="24" font-weight="760"><text x="135" y="235">x</text><text x="270" y="235">y</text><text x="135" y="295">−1</text><text x="270" y="295">−1</text><text x="135" y="355">0</text><text x="270" y="355">1</text><text x="135" y="415">2</text><text x="270" y="415">5</text></g><g stroke="#b8ccda" stroke-width="2"><path d="M75 250h260M75 310h260M75 370h260M75 430h260M200 205v240"/></g>
 <g transform="translate(470 105)"><g stroke="#c5d7e4" stroke-width="2"><path d="M0 0v400M80 0v400M160 0v400M240 0v400M320 0v400M400 0v400M480 0v400M0 0h560M0 64h560M0 128h560M0 192h560M0 256h560M0 320h560M0 384h560"/></g><g stroke="#365675" stroke-width="6"><path d="M0 320h590"/><path d="M160 420V0"/></g><path d="M80 384L320 0" stroke="#c9151e" stroke-width="10"/><g fill="url(#blue)" stroke="#fff" stroke-width="4"><circle cx="80" cy="384" r="15"/><circle cx="160" cy="256" r="15"/><circle cx="320" cy="0" r="15"/></g></g>
 <text x="780" y="548" font-size="22" font-weight="760">каждый шаг x на 1 поднимает y на 2</text>
</g>
''')
    save("map-26-systems.svg", "Два условия имеют одну общую пару", "Система икс плюс игрек равно десяти и два икс плюс игрек равно шестнадцати после вычитания дает икс шесть и игрек четыре.", '''
<g data-system="x+y=10;2x+y=16;(6,4)" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="55" y="75" width="430" height="120" rx="26"/><rect x="715" y="75" width="430" height="120" rx="26"/><rect x="385" y="430" width="430" height="120" rx="26"/></g>
 <g font-size="38" font-weight="850"><text x="270" y="145">x + y = 10</text><text x="930" y="145">2x + y = 16</text><text x="600" y="500">(x; y) = (6; 4)</text></g>
 <path d="M270 220L515 405M930 220L685 405" fill="none" stroke="#d7b06a" stroke-width="16" stroke-linecap="round"/>
 <g filter="url(#s)"><circle cx="600" cy="315" r="105" fill="url(#blue)" stroke="#fff" stroke-width="6"/></g><text x="600" y="302" font-size="25" font-weight="760" fill="#fff">вычитаем</text><text x="600" y="345" font-size="34" font-weight="860" fill="#fff">x = 6</text>
 <g font-size="21" font-weight="720" fill="#5b6b79"><text x="270" y="180">количество</text><text x="930" y="180">стоимость в сотнях</text><text x="600" y="535">проверка в обоих уравнениях</text></g>
</g>
''')
    save("map-27-powers-roots.svg", "Степени собирают одинаковые множители", "Три множителя два и четыре множителя два образуют семь множителей, поэтому два в третьей умножить на два в четвертой равно двум в седьмой.", '''
<g data-powers="2^3*2^4=2^7=128;sqrt(144)=12" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect x="45" y="75" width="480" height="230" rx="26"/><rect x="675" y="75" width="480" height="230" rx="26"/><rect x="165" y="390" width="870" height="145" rx="28" fill="url(#gold)"/></g>
 <text x="285" y="135" font-size="31" font-weight="850">2³ × 2⁴</text><text x="915" y="135" font-size="31" font-weight="850">2⁷ = 128</text>
 <g font-size="25" font-weight="800"><text x="285" y="210">2·2·2 | 2·2·2·2</text><text x="915" y="210">2·2·2·2·2·2·2</text></g><g fill="none" stroke="#b80f18" stroke-width="9" stroke-linecap="round" stroke-linejoin="round"><path d="M545 190h110m-20-18 20 18-20 18"/></g>
 <text x="600" y="447" font-size="28" font-weight="780">обратный вопрос</text><text x="600" y="500" font-size="40" font-weight="860">√144 = 12, но x² = 144 → x = ±12</text>
</g>
''')
    save("map-28-polynomials.svg", "Умножение двух двучленов как четыре области", "Произведение двух икс плюс три и икс минус четыре дает области два икс квадрат, минус восемь икс, три икс и минус двенадцать, которые собираются в многочлен.", '''
<g data-polynomial="(2x+3)(x-4)=2x^2-5x-12" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g transform="translate(85 90)" filter="url(#s)" stroke="#365675" stroke-width="5"><rect width="650" height="360" rx="24" fill="#fff"/><path d="M430 0v360M0 235h650"/></g>
 <g font-size="34" font-weight="850"><text x="300" y="230">2x²</text><text x="625" y="230">−8x</text><text x="300" y="395">3x</text><text x="625" y="395">−12</text></g><g font-size="22" font-weight="760" fill="#5b6b79"><text x="300" y="70">2x</text><text x="625" y="70">+3</text><text x="55" y="220">x</text><text x="50" y="390">−4</text></g>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="url(#gold)"><rect x="800" y="115" width="350" height="330" rx="28"/></g>
 <text x="975" y="185" font-size="25" font-weight="780">собираем</text><text x="975" y="250" font-size="31" font-weight="850">−8x + 3x</text><text x="975" y="300" font-size="31" font-weight="850">= −5x</text><text x="975" y="380" font-size="31" font-weight="860">2x² − 5x − 12</text>
 <text x="600" y="555" font-size="29" font-weight="820">каждый член первой скобки × каждый член второй</text>
</g>
''')
    save("map-29-quadratics.svg", "Три формы одной параболы", "Функция икс квадрат минус четыре икс плюс три имеет корни один и три, вершину два минус один и одну параболу.", '''
<g data-quadratic="x^2-4x+3=(x-1)(x-3)=(x-2)^2-1" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect x="45" y="55" width="340" height="105" rx="24"/><rect x="430" y="55" width="340" height="105" rx="24"/><rect x="815" y="55" width="340" height="105" rx="24"/></g><g font-size="27" font-weight="840"><text x="215" y="120">x² − 4x + 3</text><text x="600" y="120">(x − 1)(x − 3)</text><text x="985" y="120">(x − 2)² − 1</text></g>
 <g transform="translate(180 190)"><g stroke="#c4d6e3" stroke-width="2"><path d="M0 0v330M140 0v330M280 0v330M420 0v330M560 0v330M700 0v330M840 0v330M0 55h840M0 165h840M0 275h840"/></g><path d="M0 165h860M280 340V0" stroke="#365675" stroke-width="6"/><path d="M350 28Q560 522 770 28" fill="none" stroke="#c9151e" stroke-width="10"/><g fill="url(#blue)" stroke="#fff" stroke-width="4"><circle cx="420" cy="165" r="15"/><circle cx="700" cy="165" r="15"/><circle cx="560" cy="275" r="17"/></g></g>
 <g fill="#fff" stroke="#365675" stroke-width="2"><rect x="530" y="365" width="140" height="38" rx="16"/><rect x="810" y="365" width="140" height="38" rx="16"/><rect x="640" y="478" width="200" height="38" rx="16"/></g><g font-size="20" font-weight="780"><text x="600" y="391">корень 1</text><text x="880" y="391">корень 3</text><text x="740" y="504">вершина (2; −1)</text></g>
</g>
''')
    save("map-30-exponential-log.svg", "Удвоение и обратный вопрос", "Пятьсот умножить на два в степени тэ дает через ноль один два три часа пятьсот тысячу две тысячи четыре тысячи, а логарифм по основанию два от восьми равен трем.", '''
<g data-exponential="N=500*2^t;log_2(8)=3" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect x="50" y="85" width="240" height="150" rx="24"/><rect x="335" y="85" width="240" height="150" rx="24"/><rect x="620" y="85" width="240" height="150" rx="24"/><rect x="905" y="85" width="240" height="150" rx="24"/></g>
 <g font-size="36" font-weight="850"><text x="170" y="165">500</text><text x="455" y="165">1000</text><text x="740" y="165">2000</text><text x="1025" y="165">4000</text></g><g font-size="20" font-weight="720"><text x="170" y="210">t = 0</text><text x="455" y="210">t = 1</text><text x="740" y="210">t = 2</text><text x="1025" y="210">t = 3</text></g><g font-size="29" font-weight="850" fill="#b80f18"><text x="312" y="170">×2</text><text x="597" y="170">×2</text><text x="882" y="170">×2</text></g>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="url(#gold)"><rect x="150" y="335" width="900" height="175" rx="30"/></g><text x="600" y="395" font-size="31" font-weight="820">500 × 2ᵗ = 4000 → 2ᵗ = 8</text><text x="600" y="463" font-size="45" font-weight="880">t = log₂8 = 3</text>
 <text x="600" y="570" font-size="25" font-weight="760">равный шаг времени → одинаковый множитель</text>
</g>
''')
    save("map-31-sequences.svg", "Два правила последовательностей", "Арифметическая последовательность пять восемь одиннадцать четырнадцать прибавляет три, а геометрическая два шесть восемнадцать пятьдесят четыре умножает на три.", '''
<g data-sequences="5,8,11,14;d=3|2,6,18,54;q=3" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="#fff"><rect x="55" y="70" width="1090" height="210" rx="28"/><rect x="55" y="350" width="1090" height="210" rx="28"/></g>
 <g font-size="23" font-weight="780" fill="#5b6b79"><text x="180" y="120">арифметическая</text><text x="180" y="400">геометрическая</text></g>
 <g font-size="45" font-weight="860"><text x="210" y="210">5</text><text x="470" y="210">8</text><text x="730" y="210">11</text><text x="990" y="210">14</text><text x="210" y="490">2</text><text x="470" y="490">6</text><text x="730" y="490">18</text><text x="990" y="490">54</text></g>
 <g font-size="25" font-weight="850" fill="#b80f18"><text x="340" y="205">+3 →</text><text x="600" y="205">+3 →</text><text x="860" y="205">+3 →</text><text x="340" y="485">×3 →</text><text x="600" y="485">×3 →</text><text x="860" y="485">×3 →</text></g>
 <text x="600" y="315" font-size="25" font-weight="780">номер n сообщает: выполнено n − 1 переходов</text>
</g>
''')
    save("map-32-algebra-mastery.svg", "Семь опор алгебры", "Линейные функции, системы, степени, многочлены, параболы, показательный рост и последовательности ведут к проверке восемь из десяти и повтору через семь дней.", '''
<g data-algebra-skills="7" font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff"><rect class="algebra-skill" x="35" y="55" width="250" height="95" rx="21"/><rect class="algebra-skill" x="330" y="55" width="250" height="95" rx="21"/><rect class="algebra-skill" x="625" y="55" width="250" height="95" rx="21"/><rect class="algebra-skill" x="920" y="55" width="250" height="95" rx="21"/><rect class="algebra-skill" x="180" y="205" width="250" height="95" rx="21"/><rect class="algebra-skill" x="475" y="205" width="250" height="95" rx="21"/><rect class="algebra-skill" x="770" y="205" width="250" height="95" rx="21"/></g>
 <g font-size="21" font-weight="790"><text x="160" y="101"><tspan x="160">Линейные</tspan><tspan x="160" dy="26">функции</tspan></text><text x="455" y="113">Системы</text><text x="750" y="101"><tspan x="750">Степени</tspan><tspan x="750" dy="26">и корни</tspan></text><text x="1045" y="113">Многочлены</text><text x="305" y="263">Параболы</text><text x="600" y="250"><tspan x="600">Показательный</tspan><tspan x="600" dy="27">рост</tspan></text><text x="895" y="263" font-size="17">Последовательности</text></g>
 <path d="M160 170C250 350 420 360 600 390M455 170C500 300 535 340 600 390M750 170C700 300 665 340 600 390M1045 170C950 350 780 360 600 390M305 320C380 370 470 385 600 390M600 320V390M895 320C820 370 730 385 600 390" fill="none" stroke="#d7b06a" stroke-width="10" stroke-linecap="round"/>
 <g filter="url(#s)" stroke="#365675" stroke-width="5" fill="url(#gold)"><rect x="155" y="410" width="390" height="145" rx="28"/><rect x="655" y="410" width="390" height="145" rx="28"/></g><text x="350" y="462" font-size="23" font-weight="780">ПЕРВАЯ ПОПЫТКА</text><text x="350" y="525" font-size="50" font-weight="880">8/10</text><text x="850" y="462" font-size="23" font-weight="780">ЧЕРЕЗ 7 ДНЕЙ</text><text x="850" y="525" font-size="50" font-weight="880">7/10</text>
</g>
''')
    save("map-full-ru.svg", "Полная карта курса", "Восемь пронумерованных этапов идут слева направо, затем сверху вниз и справа налево.", '''
<defs><marker id="course-arrowhead" markerWidth="18" markerHeight="18" refX="16" refY="9" orient="auto" markerUnits="userSpaceOnUse"><path d="M0 0L18 9L0 18Z" fill="#b80f18"/></marker></defs>
<g fill="none" stroke="#b80f18" stroke-width="7" stroke-linecap="round" marker-end="url(#course-arrowhead)">
 <path class="course-arrow" d="M305 168H325"/><path class="course-arrow" d="M585 168H605"/><path class="course-arrow" d="M865 168H885"/>
 <path class="course-arrow" d="M1015 263V365"/>
 <path class="course-arrow" d="M895 468H875"/><path class="course-arrow" d="M615 468H595"/><path class="course-arrow" d="M335 468H315"/>
</g>
<g filter="url(#s)" stroke="#365675" stroke-width="4" fill="#fff">
 <rect class="course-step" x="55" y="80" width="240" height="175" rx="25"/><rect class="course-step" x="335" y="80" width="240" height="175" rx="25"/><rect class="course-step" x="615" y="80" width="240" height="175" rx="25"/><rect class="course-step" x="895" y="80" width="240" height="175" rx="25"/>
 <rect class="course-step" x="895" y="380" width="240" height="175" rx="25"/><rect class="course-step" x="615" y="380" width="240" height="175" rx="25"/><rect class="course-step" x="335" y="380" width="240" height="175" rx="25"/><rect class="course-step" x="55" y="380" width="240" height="175" rx="25"/>
</g>
<g fill="url(#red)" stroke="#fff" stroke-width="3"><circle cx="82" cy="107" r="20"/><circle cx="362" cy="107" r="20"/><circle cx="642" cy="107" r="20"/><circle cx="922" cy="107" r="20"/><circle cx="922" cy="407" r="20"/><circle cx="642" cy="407" r="20"/><circle cx="362" cy="407" r="20"/><circle cx="82" cy="407" r="20"/></g>
<g font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d">
 <g font-size="22" font-weight="750">
  <text x="175" y="145"><tspan x="175">Числа</tspan><tspan x="175" dy="30">и действия</tspan></text>
  <text x="455" y="123"><tspan x="455">Дроби</tspan><tspan x="455" dy="29">отношения</tspan><tspan x="455" dy="29">проценты</tspan></text>
  <text x="735" y="145"><tspan x="735">Переменные</tspan><tspan x="735" dy="30">и предалгебра</tspan></text>
  <text x="1015" y="145"><tspan x="1015">Алгебра</tspan><tspan x="1015" dy="30">и функции</tspan></text>
  <text x="1015" y="445"><tspan x="1015">Геометрия</tspan><tspan x="1015" dy="30">и пространство</tspan></text>
  <text x="735" y="445"><tspan x="735">Математический</tspan><tspan x="735" dy="30">анализ</tspan></text>
  <text x="455" y="445"><tspan x="455">Линейная</tspan><tspan x="455" dy="30">алгебра</tspan></text>
  <text x="175" y="438"><tspan x="175">Вероятность</tspan><tspan x="175" dy="29">и статистика</tspan><tspan x="175" dy="29">логика и графы</tspan></text>
 </g>
 <g font-size="19" fill="#5b6b79"><text x="175" y="229">7 узлов</text><text x="455" y="229">8 узлов</text><text x="735" y="229">5 узлов</text><text x="1015" y="229">7 узлов</text><text x="1015" y="529">8 узлов</text><text x="735" y="529">5 узлов</text><text x="455" y="529">3 узла</text><text x="175" y="529">7 узлов</text></g>
 <g font-size="18" font-weight="800" fill="#fff"><text x="82" y="113">1</text><text x="362" y="113">2</text><text x="642" y="113">3</text><text x="922" y="113">4</text><text x="922" y="413">5</text><text x="642" y="413">6</text><text x="362" y="413">7</text><text x="82" y="413">8</text></g>
</g>
''')


if __name__ == "__main__":
    main()
