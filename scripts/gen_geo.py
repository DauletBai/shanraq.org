#!/usr/bin/env python3
"""Generate the geo_nodes migration (schema + seed) for the location cascade.

The tree is authored here in Python and emitted as SQL so hundreds of rows stay
maintainable and parent links are resolved by generated codes. Run:

    python3 scripts/gen_geo.py

Coverage: all Kazakhstan regions + major cities (primary market), and Russia's
internationally recognized federal subjects + capitals/major cities. Disputed
territories are intentionally excluded. The dataset is extensible — add nodes.
"""
import pathlib

OUT = pathlib.Path("pkg/modules/migrations/sql/20251107001100_create_geo.sql")

# node = (name_ru, name_kk, name_en, kind, [children])
def N(ru, kk="", en="", kind="", children=None):
    return (ru, kk, en, kind, children or [])

def city(ru, kk="", en=""):
    return N(ru, kk, en, "city", [])

def district(ru):
    return N(ru, "", "", "district", [])

def region(ru, kk, en, cities):
    return N(ru, kk, en, "region", [city(c) for c in cities])

def subject(ru, kind, cities):
    # A Russian federal subject with its cities (capital first).
    return N(ru, "", "", kind, [city(c) for c in cities])

def fedcity(ru, districts):
    return N(ru, "", "", "city", [district(d) for d in districts])

# ---------------- Kazakhstan ----------------
KZ = N("Казахстан", "Қазақстан", "Kazakhstan", "country", [
    fedcity("Астана", ["Есиль", "Алматы", "Сарыарка", "Байконур", "Нура"]),
    fedcity("Алматы", ["Алатау", "Алмалы", "Ауэзов", "Бостандык", "Жетысу", "Медеу", "Наурызбай", "Турксиб"]),
    fedcity("Шымкент", ["Абай", "Аль-Фараби", "Енбекши", "Каратау", "Туран"]),
    region("Абайская область", "Абай облысы", "Abai Region", ["Семей", "Аягоз", "Курчатов"]),
    region("Акмолинская область", "Ақмола облысы", "Akmola Region", ["Кокшетау", "Степногорск", "Атбасар", "Есиль"]),
    region("Актюбинская область", "Ақтөбе облысы", "Aktobe Region", ["Актобе", "Хромтау", "Кандыагаш", "Шалкар"]),
    region("Алматинская область", "Алматы облысы", "Almaty Region", ["Конаев", "Талгар", "Есик", "Каскелен", "Жаркент"]),
    region("Атырауская область", "Атырау облысы", "Atyrau Region", ["Атырау", "Кульсары"]),
    region("Западно-Казахстанская область", "Батыс Қазақстан облысы", "West Kazakhstan Region", ["Уральск", "Аксай"]),
    region("Жамбылская область", "Жамбыл облысы", "Jambyl Region", ["Тараз", "Каратау", "Шу", "Жанатас"]),
    region("Жетысуская область", "Жетісу облысы", "Jetisu Region", ["Талдыкорган", "Текели", "Уштобе", "Сарканд"]),
    region("Карагандинская область", "Қарағанды облысы", "Karaganda Region", ["Караганда", "Темиртау", "Сарань", "Шахтинск", "Балхаш"]),
    region("Костанайская область", "Қостанай облысы", "Kostanay Region", ["Костанай", "Рудный", "Лисаковск", "Аркалык"]),
    region("Кызылординская область", "Қызылорда облысы", "Kyzylorda Region", ["Кызылорда", "Аральск", "Казалинск"]),
    region("Мангистауская область", "Маңғыстау облысы", "Mangystau Region", ["Актау", "Жанаозен", "Форт-Шевченко"]),
    region("Северо-Казахстанская область", "Солтүстік Қазақстан облысы", "North Kazakhstan Region", ["Петропавловск", "Булаево", "Тайынша"]),
    region("Павлодарская область", "Павлодар облысы", "Pavlodar Region", ["Павлодар", "Экибастуз", "Аксу"]),
    region("Туркестанская область", "Түркістан облысы", "Turkestan Region", ["Туркестан", "Кентау", "Арыс", "Сарыагаш"]),
    region("Улытауская область", "Ұлытау облысы", "Ulytau Region", ["Жезказган", "Сатпаев", "Каражал"]),
    region("Восточно-Казахстанская область", "Шығыс Қазақстан облысы", "East Kazakhstan Region", ["Усть-Каменогорск", "Риддер", "Алтай", "Серебрянск"]),
])

# ---------------- Russia (recognized subjects; disputed excluded) ----------------
RU = N("Россия", "Ресей", "Russia", "country", [
    fedcity("Москва", ["ЦАО", "САО", "СВАО", "ВАО", "ЮВАО", "ЮАО", "ЮЗАО", "ЗАО", "СЗАО", "Зеленоградский АО"]),
    fedcity("Санкт-Петербург", ["Адмиралтейский", "Василеостровский", "Выборгский", "Калининский", "Кировский", "Центральный", "Невский", "Приморский", "Московский"]),
    # Republics
    subject("Республика Адыгея", "republic", ["Майкоп"]),
    subject("Республика Алтай", "republic", ["Горно-Алтайск"]),
    subject("Республика Башкортостан", "republic", ["Уфа", "Стерлитамак", "Салават", "Нефтекамск"]),
    subject("Республика Бурятия", "republic", ["Улан-Удэ"]),
    subject("Республика Дагестан", "republic", ["Махачкала", "Дербент", "Хасавюрт"]),
    subject("Республика Ингушетия", "republic", ["Магас", "Назрань"]),
    subject("Кабардино-Балкарская Республика", "republic", ["Нальчик"]),
    subject("Республика Калмыкия", "republic", ["Элиста"]),
    subject("Карачаево-Черкесская Республика", "republic", ["Черкесск"]),
    subject("Республика Карелия", "republic", ["Петрозаводск"]),
    subject("Республика Коми", "republic", ["Сыктывкар", "Ухта", "Воркута"]),
    subject("Республика Марий Эл", "republic", ["Йошкар-Ола"]),
    subject("Республика Мордовия", "republic", ["Саранск"]),
    subject("Республика Саха (Якутия)", "republic", ["Якутск", "Нерюнгри"]),
    subject("Республика Северная Осетия — Алания", "republic", ["Владикавказ"]),
    subject("Республика Татарстан", "republic", ["Казань", "Набережные Челны", "Альметьевск", "Нижнекамск"]),
    subject("Республика Тыва", "republic", ["Кызыл"]),
    subject("Удмуртская Республика", "republic", ["Ижевск", "Сарапул", "Воткинск"]),
    subject("Республика Хакасия", "republic", ["Абакан"]),
    subject("Чеченская Республика", "republic", ["Грозный", "Гудермес"]),
    subject("Чувашская Республика", "republic", ["Чебоксары", "Новочебоксарск"]),
    # Krais
    subject("Алтайский край", "krai", ["Барнаул", "Бийск", "Рубцовск"]),
    subject("Забайкальский край", "krai", ["Чита"]),
    subject("Камчатский край", "krai", ["Петропавловск-Камчатский"]),
    subject("Краснодарский край", "krai", ["Краснодар", "Сочи", "Новороссийск", "Армавир"]),
    subject("Красноярский край", "krai", ["Красноярск", "Норильск", "Ачинск"]),
    subject("Пермский край", "krai", ["Пермь", "Березники"]),
    subject("Приморский край", "krai", ["Владивосток", "Находка", "Уссурийск"]),
    subject("Ставропольский край", "krai", ["Ставрополь", "Пятигорск", "Невинномысск"]),
    subject("Хабаровский край", "krai", ["Хабаровск", "Комсомольск-на-Амуре"]),
    # Oblasts
    subject("Московская область", "oblast", ["Красногорск", "Балашиха", "Химки", "Подольск", "Мытищи", "Люберцы"]),
    subject("Ленинградская область", "oblast", ["Гатчина", "Выборг", "Всеволожск"]),
    subject("Свердловская область", "oblast", ["Екатеринбург", "Нижний Тагил", "Каменск-Уральский"]),
    subject("Новосибирская область", "oblast", ["Новосибирск", "Бердск"]),
    subject("Нижегородская область", "oblast", ["Нижний Новгород", "Дзержинск", "Арзамас"]),
    subject("Ростовская область", "oblast", ["Ростов-на-Дону", "Таганрог", "Шахты", "Волгодонск"]),
    subject("Челябинская область", "oblast", ["Челябинск", "Магнитогорск", "Златоуст"]),
    subject("Самарская область", "oblast", ["Самара", "Тольятти", "Сызрань"]),
    subject("Воронежская область", "oblast", ["Воронеж", "Борисоглебск"]),
    subject("Волгоградская область", "oblast", ["Волгоград", "Волжский", "Камышин"]),
    subject("Саратовская область", "oblast", ["Саратов", "Энгельс", "Балаково"]),
    subject("Тюменская область", "oblast", ["Тюмень", "Тобольск"]),
    subject("Омская область", "oblast", ["Омск"]),
    subject("Иркутская область", "oblast", ["Иркутск", "Братск", "Ангарск"]),
    subject("Кемеровская область", "oblast", ["Кемерово", "Новокузнецк", "Прокопьевск"]),
    subject("Оренбургская область", "oblast", ["Оренбург", "Орск"]),
    subject("Томская область", "oblast", ["Томск", "Северск"]),
    subject("Ульяновская область", "oblast", ["Ульяновск", "Димитровград"]),
    subject("Ярославская область", "oblast", ["Ярославль", "Рыбинск"]),
    subject("Владимирская область", "oblast", ["Владимир", "Ковров", "Муром"]),
    subject("Тульская область", "oblast", ["Тула", "Новомосковск"]),
    subject("Рязанская область", "oblast", ["Рязань"]),
    subject("Липецкая область", "oblast", ["Липецк", "Елец"]),
    subject("Пензенская область", "oblast", ["Пенза"]),
    subject("Кировская область", "oblast", ["Киров"]),
    subject("Архангельская область", "oblast", ["Архангельск", "Северодвинск"]),
    subject("Вологодская область", "oblast", ["Вологда", "Череповец"]),
    subject("Калининградская область", "oblast", ["Калининград"]),
    subject("Калужская область", "oblast", ["Калуга", "Обнинск"]),
    subject("Курская область", "oblast", ["Курск"]),
    subject("Белгородская область", "oblast", ["Белгород", "Старый Оскол"]),
    subject("Брянская область", "oblast", ["Брянск"]),
    subject("Смоленская область", "oblast", ["Смоленск"]),
    subject("Тверская область", "oblast", ["Тверь"]),
    subject("Астраханская область", "oblast", ["Астрахань"]),
    subject("Мурманская область", "oblast", ["Мурманск", "Североморск"]),
    subject("Новгородская область", "oblast", ["Великий Новгород"]),
    subject("Псковская область", "oblast", ["Псков"]),
    subject("Тамбовская область", "oblast", ["Тамбов", "Мичуринск"]),
    subject("Ивановская область", "oblast", ["Иваново"]),
    subject("Костромская область", "oblast", ["Кострома"]),
    subject("Орловская область", "oblast", ["Орёл"]),
    subject("Курганская область", "oblast", ["Курган"]),
    subject("Амурская область", "oblast", ["Благовещенск"]),
    subject("Сахалинская область", "oblast", ["Южно-Сахалинск"]),
    subject("Магаданская область", "oblast", ["Магадан"]),
    # Autonomous okrugs / oblast
    subject("Ханты-Мансийский автономный округ — Югра", "okrug", ["Ханты-Мансийск", "Сургут", "Нижневартовск"]),
    subject("Ямало-Ненецкий автономный округ", "okrug", ["Салехард", "Новый Уренгой", "Ноябрьск"]),
    subject("Чукотский автономный округ", "okrug", ["Анадырь"]),
    subject("Ненецкий автономный округ", "okrug", ["Нарьян-Мар"]),
    subject("Еврейская автономная область", "oblast", ["Биробиджан"]),
])

rows = []
counter = [0]

def esc(s):
    return s.replace("'", "''")

def walk(node, parent_code, country, level, sort):
    ru, kk, en, kind, children = node
    counter[0] += 1
    code = "g%d" % counter[0]
    rows.append((code, parent_code, country, level, kind, ru, kk, en, sort))
    for i, ch in enumerate(children):
        walk(ch, code, country, level + 1, i)

walk(KZ, None, "KZ", 0, 0)
walk(RU, None, "RU", 0, 1)

lines = []
for code, parent, country, level, kind, ru, kk, en, sort in rows:
    p = "NULL" if parent is None else "'%s'" % parent
    lines.append(
        "('%s',%s,'%s',%d,'%s','%s','%s','%s',%d)"
        % (code, p, country, level, kind, esc(ru), esc(kk), esc(en), sort)
    )

sql = []
sql.append("-- +goose Up")
sql.append("-- Generated by scripts/gen_geo.py — do not edit by hand.")
sql.append("""CREATE TABLE geo_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT UNIQUE NOT NULL,
    parent_code TEXT,
    parent_id UUID,
    country CHAR(2) NOT NULL,
    level INT NOT NULL,
    kind TEXT NOT NULL,
    name_ru TEXT NOT NULL,
    name_kk TEXT NOT NULL DEFAULT '',
    name_en TEXT NOT NULL DEFAULT '',
    sort INT NOT NULL DEFAULT 0
);""")
sql.append("-- +goose StatementBegin")
sql.append("INSERT INTO geo_nodes (code,parent_code,country,level,kind,name_ru,name_kk,name_en,sort) VALUES")
sql.append(",\n".join(lines) + ";")
sql.append("-- +goose StatementEnd")
sql.append("UPDATE geo_nodes c SET parent_id = p.id FROM geo_nodes p WHERE c.parent_code = p.code;")
sql.append("ALTER TABLE geo_nodes ADD CONSTRAINT geo_nodes_parent_fk FOREIGN KEY (parent_id) REFERENCES geo_nodes(id) ON DELETE CASCADE;")
sql.append("CREATE INDEX idx_geo_parent ON geo_nodes(parent_id);")
sql.append("CREATE INDEX idx_geo_roots ON geo_nodes(country) WHERE parent_id IS NULL;")
sql.append("")
sql.append("-- +goose Down")
sql.append("DROP TABLE IF EXISTS geo_nodes;")
sql.append("")

OUT.write_text("\n".join(sql))
print("wrote", OUT, "with", len(rows), "nodes")
