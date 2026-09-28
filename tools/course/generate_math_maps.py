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


def main():
    save("map-01-diagnostic.svg", "Learning route", "A bridge of connected mathematical ideas.", '''
<path d="M100 475 C290 475 280 360 455 360 S650 245 795 245 S930 135 1090 135" fill="none" stroke="#d7b06a" stroke-width="34" stroke-linecap="round" filter="url(#s)"/>
<g fill="#fff" stroke="#365675" stroke-width="6">
 <circle cx="120" cy="475" r="54"/><circle cx="350" cy="390" r="54"/><circle cx="570" cy="310" r="54"/><circle cx="790" cy="245" r="54"/><circle cx="1015" cy="145" r="54"/>
</g><g font-family="system-ui,sans-serif" font-size="38" font-weight="700" text-anchor="middle" fill="#17324d">
 <text x="120" y="488">12</text><text x="350" y="403">3/4</text><text x="570" y="323">2:5</text><text x="790" y="258">25%</text><text x="1015" y="158">x</text>
</g><path d="M80 550h1040" stroke="#365675" stroke-width="3" stroke-dasharray="12 14" opacity=".45"/>
''')
    save("map-02-fraction.svg", "Three meanings of a fraction", "Equal parts, a point on a line, and division lead to the same fraction.", '''
<g filter="url(#s)"><circle cx="250" cy="300" r="132" fill="#fff" stroke="#365675" stroke-width="6"/>
<path d="M250 300V168A132 132 0 0 1 343.3 206.7Z" fill="url(#red)"/><path d="M250 300l93.3-93.3a132 132 0 0 1 38.7 93.3Z" fill="url(#red)"/>
<path d="M250 168v264M118 300h264M156.7 206.7l186.6 186.6M156.7 393.3l186.6-186.6" stroke="#365675" stroke-width="4"/>
</g><path d="M480 340h560" stroke="#365675" stroke-width="8" stroke-linecap="round"/>
<g stroke="#365675" stroke-width="5"><path d="M480 315v50M620 315v50M760 315v50M900 315v50M1040 315v50"/></g>
<path d="M480 340h280" stroke="url(#red)" stroke-width="18"/><circle cx="760" cy="340" r="18" fill="#b80f18"/>
<g font-family="system-ui,sans-serif" font-size="42" font-weight="700" text-anchor="middle" fill="#17324d"><text x="250" y="510">2/8</text><text x="760" y="420">2 ÷ 8</text></g>
''')
    save("map-03-equivalence.svg", "Equivalent fractions", "Bars cut into different numbers of equal pieces show the same length.", '''
<g filter="url(#s)" stroke="#365675" stroke-width="5"><rect x="130" y="155" width="900" height="100" rx="15" fill="#fff"/><rect x="130" y="155" width="450" height="100" rx="15" fill="url(#red)"/><path d="M580 155v100"/>
<rect x="130" y="365" width="900" height="100" rx="15" fill="#fff"/><rect x="130" y="365" width="450" height="100" rx="15" fill="url(#blue)"/><path d="M355 365v100M580 365v100M805 365v100"/></g>
<path d="M570 300l-24-24m24 24l24-24M570 300l-24 24m24-24l24 24" stroke="#e39a16" stroke-width="10" stroke-linecap="round"/>
<g font-family="system-ui,sans-serif" font-size="46" font-weight="700" text-anchor="end" fill="#17324d"><text x="110" y="220">1/2</text><text x="110" y="430">2/4</text></g>
''')
    save("map-04-compare.svg", "Compare fractions", "Both fractions are placed on one number line.", '''
<path d="M120 320h960" stroke="#365675" stroke-width="10" stroke-linecap="round"/>
<g stroke="#365675" stroke-width="5"><path d="M120 275v90M360 285v70M600 275v90M840 285v70M1080 275v90"/></g>
<circle cx="570" cy="320" r="28" fill="url(#blue)" filter="url(#s)"/><circle cx="760" cy="320" r="28" fill="url(#red)" filter="url(#s)"/>
<path d="M570 235v-70M760 235v-70" stroke="#365675" stroke-width="4"/>
<g font-family="system-ui,sans-serif" font-size="48" font-weight="700" text-anchor="middle" fill="#17324d"><text x="120" y="410">0</text><text x="1080" y="410">1</text><text x="570" y="145">1/2</text><text x="760" y="145">2/3</text><text x="665" y="520">1/2 &lt; 2/3</text></g>
''')
    save("map-05-operations.svg", "Fraction operations", "Two bars are recut into a shared unit before pieces are combined.", '''
<g filter="url(#s)" stroke="#365675" stroke-width="5"><rect x="85" y="150" width="360" height="90" rx="14" fill="#fff"/><rect x="85" y="150" width="180" height="90" rx="14" fill="url(#red)"/><path d="M265 150v90"/>
<rect x="85" y="345" width="360" height="90" rx="14" fill="#fff"/><rect x="85" y="345" width="120" height="90" rx="14" fill="url(#blue)"/><path d="M205 345v90M325 345v90"/>
<rect x="690" y="245" width="420" height="90" rx="14" fill="#fff"/><rect x="690" y="245" width="350" height="90" rx="14" fill="url(#gold)"/><path d="M760 245v90M830 245v90M900 245v90M970 245v90M1040 245v90"/></g>
<g font-family="system-ui,sans-serif" font-size="52" font-weight="700" text-anchor="middle" fill="#17324d"><text x="500" y="315">+</text><text x="610" y="315">→</text><text x="900" y="420">3/6 + 2/6 = 5/6</text></g>
''')
    save("map-06-ratio.svg", "Ratio", "Equal recipe groups preserve the same relationship.", '''
<g filter="url(#s)"><g fill="url(#red)"><circle cx="220" cy="220" r="55"/><circle cx="220" cy="410" r="55"/></g><g fill="url(#blue)"><circle cx="390" cy="170" r="55"/><circle cx="390" cy="315" r="55"/><circle cx="390" cy="460" r="55"/></g>
<path d="M510 315h100" stroke="#365675" stroke-width="12" stroke-linecap="round"/><path d="M590 280l40 35-40 35" fill="none" stroke="#365675" stroke-width="12" stroke-linecap="round" stroke-linejoin="round"/>
<g fill="url(#red)"><circle cx="745" cy="170" r="42"/><circle cx="745" cy="270" r="42"/><circle cx="745" cy="370" r="42"/><circle cx="745" cy="470" r="42"/></g><g fill="url(#blue)"><circle cx="920" cy="115" r="42"/><circle cx="920" cy="195" r="42"/><circle cx="920" cy="275" r="42"/><circle cx="920" cy="355" r="42"/><circle cx="920" cy="435" r="42"/><circle cx="920" cy="515" r="42"/></g></g>
<g font-family="system-ui,sans-serif" font-size="52" font-weight="700" text-anchor="middle" fill="#17324d"><text x="305" y="585">2 : 3</text><text x="835" y="585">4 : 6</text></g>
''')
    save("map-07-percent.svg", "Percent", "A hundred-cell grid connects fractions, decimals, and percent.", '''
<g transform="translate(115 65)" filter="url(#s)"><rect width="500" height="500" rx="12" fill="#fff" stroke="#365675" stroke-width="6"/>
<path d="M0 0h185v500H0z" fill="url(#red)" opacity=".9"/><g stroke="#365675" stroke-width="2" opacity=".55">'''+''.join(f'<path d="M{x} 0v500"/>' for x in range(50,500,50))+''.join(f'<path d="M0 {y}h500"/>' for y in range(50,500,50))+'''</g></g>
<path d="M700 315h120" stroke="#365675" stroke-width="12"/><path d="M800 280l40 35-40 35" fill="none" stroke="#365675" stroke-width="12" stroke-linecap="round"/>
<g font-family="system-ui,sans-serif" font-size="58" font-weight="700" text-anchor="middle" fill="#17324d"><text x="950" y="260">37/100</text><text x="950" y="340">0.37</text><text x="950" y="420">37%</text></g>
''')
    save("map-08-proportion.svg", "Proportion", "Two equivalent ratios balance because both terms scale together.", '''
<path d="M600 130v350M300 480h600" stroke="#365675" stroke-width="18" stroke-linecap="round" filter="url(#s)"/><path d="M255 250h690" stroke="#365675" stroke-width="14"/><path d="M350 250l-90 210h180zM850 250l-90 210h180z" fill="#fff" stroke="#365675" stroke-width="7"/>
<g font-family="system-ui,sans-serif" font-size="64" font-weight="700" text-anchor="middle" fill="#17324d"><text x="350" y="390">2/3</text><text x="850" y="390">6/9</text><text x="600" y="570">×3</text></g><circle cx="600" cy="250" r="34" fill="url(#gold)" stroke="#365675" stroke-width="5"/>
''')
    save("map-09-mastery.svg", "Mastery gate", "Understanding is checked by representation, calculation, explanation, error repair, and transfer.", '''
<path d="M110 475 C300 475 270 360 440 360 S620 250 760 250 S920 145 1080 145" fill="none" stroke="#d7b06a" stroke-width="42" stroke-linecap="round" filter="url(#s)"/>
<g fill="#fff" stroke="#365675" stroke-width="6"><rect x="85" y="425" width="160" height="100" rx="24"/><rect x="315" y="325" width="160" height="100" rx="24"/><rect x="545" y="225" width="160" height="100" rx="24"/><rect x="775" y="125" width="160" height="100" rx="24"/></g>
<g font-family="system-ui,sans-serif" font-size="42" font-weight="700" text-anchor="middle" fill="#17324d"><text x="165" y="489">◫</text><text x="395" y="389">=</text><text x="625" y="289">?</text><text x="855" y="189">↗</text><text x="1050" y="125">8/10</text><text x="1050" y="180">7/10</text><text x="1050" y="235">+7d</text></g>
''')
    save("map-full-ru.svg", "Полная карта курса", "Восемь зависимых этапов ведут от чисел к математическому моделированию.", '''
<path d="M125 165h950M1075 165v300H125M125 465h950" fill="none" stroke="#d7b06a" stroke-width="24" stroke-linecap="round" stroke-linejoin="round" filter="url(#s)"/>
<g filter="url(#s)" stroke="#365675" stroke-width="4"><rect x="55" y="90" width="240" height="145" rx="25" fill="#fff"/><rect x="335" y="90" width="240" height="145" rx="25" fill="#fff"/><rect x="615" y="90" width="240" height="145" rx="25" fill="#fff"/><rect x="895" y="90" width="240" height="145" rx="25" fill="#fff"/><rect x="895" y="390" width="240" height="145" rx="25" fill="#fff"/><rect x="615" y="390" width="240" height="145" rx="25" fill="#fff"/><rect x="335" y="390" width="240" height="145" rx="25" fill="#fff"/><rect x="55" y="390" width="240" height="145" rx="25" fill="#fff"/></g>
<g font-family="system-ui,sans-serif" text-anchor="middle" fill="#17324d"><g font-size="25" font-weight="750"><text x="175" y="145"><tspan x="175">Числа</tspan><tspan x="175" dy="34">и действия</tspan></text><text x="455" y="145"><tspan x="455">Дроби, отношения</tspan><tspan x="455" dy="34">и проценты</tspan></text><text x="735" y="145"><tspan x="735">Предалгебра:</tspan><tspan x="735" dy="34">переменные</tspan></text><text x="1015" y="145"><tspan x="1015">Алгебра</tspan><tspan x="1015" dy="34">и функции</tspan></text><text x="1015" y="445"><tspan x="1015">Геометрия</tspan><tspan x="1015" dy="34">и пространство</tspan></text><text x="735" y="445"><tspan x="735">Математический</tspan><tspan x="735" dy="34">анализ</tspan></text><text x="455" y="445"><tspan x="455">Линейная</tspan><tspan x="455" dy="34">алгебра</tspan></text><text x="175" y="435"><tspan x="175">Вероятность, данные,</tspan><tspan x="175" dy="34">логика и графы</tspan></text></g><g font-size="20" fill="#5b6b79"><text x="175" y="215">7 узлов</text><text x="455" y="215">8 узлов</text><text x="735" y="215">5 узлов</text><text x="1015" y="215">7 узлов</text><text x="1015" y="515">8 узлов</text><text x="735" y="515">5 узлов</text><text x="455" y="515">3 узла</text><text x="175" y="515">7 узлов</text></g></g>
<circle cx="45" cy="165" r="24" fill="url(#red)"/><path d="M1135 465h30l-18-17m18 17l-18 17" stroke="#b80f18" stroke-width="8" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
''')


if __name__ == "__main__":
    main()
