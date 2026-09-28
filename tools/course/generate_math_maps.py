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
