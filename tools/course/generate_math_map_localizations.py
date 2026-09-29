#!/usr/bin/env python3
"""Create reviewed Kazakh and English SVG maps without external services."""
from pathlib import Path
import re
import xml.etree.ElementTree as ET

from math_localization_data import DATA

ROOT = Path(__file__).resolve().parents[2]
MAPS = ROOT / "web/static/course/mathematics"
ET.register_namespace("", "http://www.w3.org/2000/svg")

MAP_STEMS = {"map-full-ru.svg": "preface"}
for stem in DATA:
    if stem != "preface":
        number = stem[:2]
        candidates = list(MAPS.glob(f"map-{number}-*.svg"))
        candidates = [p for p in candidates if not p.stem.endswith(("-kz", "-en"))]
        if len(candidates) != 1:
            raise RuntimeError(f"map lookup for {stem}: {candidates}")
        MAP_STEMS[candidates[0].name] = stem

T = {
"Модель":("Модель","Model"), "Расчёт":("Есептеу","Calculation"), "Объяснение":("Түсіндіру","Explanation"),
"Исправление":("Қатені түзету","Error correction"), "ошибки":("",""), "Перенос":("Көшіру","Transfer"),
"ПЕРВАЯ ПОПЫТКА":("АЛҒАШҚЫ ӘРЕКЕТ","FIRST ATTEMPT"), "ЧЕРЕЗ 7 ДНЕЙ":("7 КҮННЕН КЕЙІН","AFTER 7 DAYS"),
"ещё 3":("тағы 3","3 more"), "группа × размер группы + остаток":("топ × топ көлемі + қалдық","groups × group size + remainder"),
"ТЫСЯЧИ":("МЫҢДЫҚ","THOUSANDS"), "СОТНИ":("ЖҮЗДІК","HUNDREDS"), "ДЕСЯТКИ":("ОНДЫҚ","TENS"), "ЕДИНИЦЫ":("БІРЛІК","ONES"),
"часть + часть = целое":("бөлік + бөлік = бүтін","part + part = whole"), "целое − часть = другая часть":("бүтін − бөлік = басқа бөлік","whole − part = other part"),
"группы × в группе = всего":("топ саны × топтағы саны = барлығы","groups × group size = total"), "сколько групп?":("неше топ?","how many groups?"), "сколько в группе?":("әр топта нешеу?","how many per group?"),
"оценка: 240 − примерно 150 ≈ 90":("бағалау: 240 − шамамен 150 ≈ 90","estimate: 240 − about 150 ≈ 90"), "+7 шагов вправо":("оңға 7 қадам","7 steps right"),
"тысячные":("мыңдық үлес","thousandths"), "0 единиц | 3 десятых | 7 сотых | 5 тысячных":("0 бірлік | 3 ондық үлес | 7 жүздік үлес | 5 мыңдық үлес","0 ones | 3 tenths | 7 hundredths | 5 thousandths"),
"Количество":("Мөлшер","Quantity"), "Разряды":("Разрядтар","Place value"), "Сложение":("Қосу","Addition"), "и вычитание":("және азайту","and subtraction"), "Умножение":("Көбейту","Multiplication"), "и деление":("және бөлу","and division"), "Порядок":("Амалдар реті","Order"), "и оценка":("және бағалау","and estimation"), "Отрицательные":("Теріс","Negative"), "числа":("сандар","numbers"), "Делимость":("Бөлінгіштік","Divisibility"), "Десятичные":("Ондық","Decimal"), "дроби":("бөлшектер","fractions"),
"плата":("төлем","boarding"), "за посадку":("отыру үшін","charge"), "120 тенге за каждый":("әр километрге","120 tenge per"), "километр расстояния d":("120 теңге, қашықтық d","kilometre, distance d"), "общая цена":("жалпы баға","total cost"), "в тенге":("теңгемен","in tenge"), "подстановка: d = 5 км":("мән қою: d = 5 км","substitute: d = 5 km"), "C = 700 + 120 × 5 = 1300 тенге":("C = 700 + 120 × 5 = 1300 теңге","C = 700 + 120 × 5 = 1300 tenge"),
"раскрываем скобки":("жақшаны ашамыз","expand parentheses"), "собираем подобные":("ұқсас мүшелерді жинаймыз","combine like terms"), "x = 4: исходное = 23":("x = 4: бастапқысы = 23","x = 4: original = 23"), "x = 4: новое = 23":("x = 4: жаңасы = 23","x = 4: new = 23"),
"−5 с обеих сторон":("екі жақтан да −5","−5 on both sides"), "÷3 с обеих сторон":("екі жағын да ÷3","÷3 on both sides"), "проверка":("тексеру","check"),
"|4 − (−3)| = 7 единиц":("|4 − (−3)| = 7 бірлік","|4 − (−3)| = 7 units"), "все числа до 4 включительно":("4-ке дейінгі барлық сан, 4 кіреді","all numbers up to and including 4"), "4 подходит: 3 × 4 + 2 = 14":("4 жарайды: 3 × 4 + 2 = 14","4 works: 3 × 4 + 2 = 14"),
"Переменные":("Айнымалылар","Variables"), "Выражения":("Өрнектер","Expressions"), "Уравнения":("Теңдеулер","Equations"), "Координаты":("Координаталар","Coordinates"), "Неравенства":("Теңсіздіктер","Inequalities"),
"старт b = 1 · шаг k = 2":("басы b = 1 · қадам k = 2","start b = 1 · step k = 2"), "каждый шаг x на 1 поднимает y на 2":("x әр 1 қадамда y-ті 2-ге арттырады","each x-step of 1 raises y by 2"),
"вычитаем":("азайтамыз","subtract"), "количество":("саны","quantity"), "стоимость в сотнях":("жүздікпен берілген құн","cost in hundreds"), "проверка в обоих уравнениях":("екі теңдеуде тексеру","check in both equations"),
"обратный вопрос":("кері сұрақ","inverse question"), "собираем":("біріктіреміз","combine"), "каждый член первой скобки × каждый член второй":("бірінші жақшаның әр мүшесі × екіншінің әр мүшесі","each term in the first factor × each term in the second"),
"√144 = 12, но x² = 144 → x = ±12":("√144 = 12, бірақ x² = 144 → x = ±12","√144 = 12, but x² = 144 → x = ±12"),
"корень 1":("1-түбір","root 1"), "корень 3":("3-түбір","root 3"), "вершина (2; −1)":("төбе (2; −1)","vertex (2, −1)"),
"равный шаг времени → одинаковый множитель":("бірдей уақыт қадамы → бірдей көбейткіш","equal time step → same factor"), "арифметическая":("арифметикалық","arithmetic"), "геометрическая":("геометриялық","geometric"), "номер n сообщает: выполнено n − 1 переходов":("n нөмірі: n − 1 ауысу жасалды","term n means n − 1 transitions"),
"Линейные":("Сызықтық","Linear"), "функции":("функциялар","functions"), "Системы":("Жүйелер","Systems"), "Степени":("Дәрежелер","Powers"), "и корни":("және түбірлер","and roots"), "Многочлены":("Көпмүшелер","Polynomials"), "Параболы":("Параболалар","Parabolas"), "Показательный":("Көрсеткіштік","Exponential"), "рост":("өсу","growth"), "Последовательности":("Тізбектер","Sequences"),
"Числа":("Сандар","Numbers"), "и действия":("және амалдар","and operations"), "Дроби":("Бөлшектер","Fractions"), "отношения":("қатынастар","ratios"), "проценты":("пайыздар","percentages"), "и предалгебра":("және предалгебра","and prealgebra"), "Алгебра":("Алгебра","Algebra"), "и функции":("және функциялар","and functions"), "Геометрия":("Геометрия","Geometry"), "и пространство":("және кеңістік","and space"), "Математический":("Математикалық","Mathematical"), "анализ":("анализ","analysis"), "Линейная":("Сызықтық","Linear"), "алгебра":("алгебра","algebra"), "Вероятность":("Ықтималдық","Probability"), "и статистика":("және статистика","and statistics"), "логика и графы":("логика және графтар","logic and graphs"), "7 узлов":("7 түйін","7 nodes"), "8 узлов":("8 түйін","8 nodes"), "5 узлов":("5 түйін","5 nodes"), "3 узла":("3 түйін","3 nodes"),
"3 из 8 = 3/8":("8-дің 3-еуі = 3/8","3 of 8 = 3/8"),
}

CYRILLIC = re.compile(r"[А-Яа-яЁё]")

def output_name(name, lang):
    if name == "map-full-ru.svg": return f"map-full-{lang}.svg"
    return name[:-4] + f"-{lang}.svg"

def generate(path, stem, lang):
    tree = ET.parse(path); root = tree.getroot()
    nodes = [e for e in root.iter() if e.tag.rsplit('}',1)[-1] in ("title","desc","text","tspan") and e.text and e.text.strip()]
    title_done = desc_done = False
    for e in nodes:
        tag=e.tag.rsplit('}',1)[-1]; source=e.text.strip()
        if tag == "title" and not title_done:
            e.text=DATA[stem][lang]["title"]; title_done=True; continue
        if tag == "desc" and not desc_done:
            e.text=DATA[stem][lang]["alt"]; desc_done=True; continue
        if source in T:
            e.text=T[source][0 if lang=="kz" else 1]
        elif CYRILLIC.search(source):
            raise ValueError(f"missing {lang} SVG translation in {path.name}: {source!r}")
    # Long translated labels need a little more breathing room in the same geometry.
    for e in root.iter():
        if e.tag.rsplit('}',1)[-1] == "text" and "font-size" in e.attrib:
            try: e.set("font-size", str(round(float(e.get("font-size"))*0.88, 2)))
            except ValueError: pass
    target=MAPS/output_name(path.name,lang)
    tree.write(target,encoding="unicode",xml_declaration=False)
    target.write_text(target.read_text(encoding="utf-8")+"\n",encoding="utf-8")
    return target

def main():
    for name,stem in sorted(MAP_STEMS.items()):
        for lang in ("kz","en"):
            print(generate(MAPS/name,stem,lang).relative_to(ROOT))

if __name__ == "__main__": main()
