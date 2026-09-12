"""Задание урока 41: восстановить веса, найти половину прироста и сменить базу."""

import pandas as pd

# Бюро национальной статистики, ИПЦ, август 2026: прирост с начала года и
# вклад раздела в этот прирост. Названия сокращены.
RAZDELY = [
    ("Продукты питания и напитки", 4.9, 1.90),
    ("Алкоголь и табак", 10.4, 0.15),
    ("Одежда и обувь", 6.6, 0.61),
    ("Жильё, вода, энергия", 8.6, 0.83),
    ("Дом и быт", 6.0, 0.33),
    ("Здравоохранение", 11.5, 0.67),
    ("Транспорт", 3.6, 0.32),
    ("Связь", 6.5, 0.31),
    ("Отдых и культура", 10.4, 0.35),
    ("Образование", 2.7, 0.07),
    ("Рестораны и гостиницы", 6.1, 0.11),
    ("Страхование и финансы", 3.0, 0.02),
    ("Личный уход и прочее", 9.6, 0.68),
]
ZAGOLOVOK = 6.4

# Всемирный банк, FP.CPI.TOTL: индекс цен, 2010 = 100.
RYAD = pd.Series({2019: 189.30, 2020: 202.02, 2021: 218.27, 2022: 251.07,
                  2023: 287.54, 2024: 312.53, 2025: 348.12})


def vesa(razdely):
    """Восстанавливает веса из вкладов: вклад = вес × прирост."""
    out = pd.DataFrame(razdely, columns=["раздел", "прирост", "вклад"])
    out["вес"] = (out["вклад"] / out["прирост"] * 100).round(1)
    return out.sort_values("вклад", ascending=False).reset_index(drop=True)


def polovina(table):
    """Сколько разделов сверху нужно, чтобы перекрыть половину прироста."""
    total = table["вклад"].sum()
    running = table["вклад"].cumsum()
    return int((running < total / 2).sum() + 1), total


def na_bazu(series, year):
    """Тот же ряд на другой базе: год-опора становится сотней."""
    return (series / series[year] * 100).round(2)


table = vesa(RAZDELY)
print("== Веса, восстановленные из вкладов")
print(table[["раздел", "вес", "прирост", "вклад"]].to_string(index=False))
print("  сумма весов:", round(table["вес"].sum(), 1), "%")

count, total = polovina(table)
print()
print("== Половина прироста")
print("  разделов хватило:", count)
print("  это:", ", ".join(table["раздел"].head(count)))
print("  их доля в приросте:", round(table["вклад"].head(count).sum() / total * 100, 1), "%")

print()
print("== Заголовок и сумма вкладов")
print("  сумма вкладов:", round(total, 2))
print("  заголовок:", ZAGOLOVOK)
print("  расхождение:", round(abs(ZAGOLOVOK - total), 2), "пункта — округление")

print()
print("== Смена базы")
rebased = na_bazu(RYAD, 2019)
first, last = RYAD.index[0], RYAD.index[-1]
print("  база 2010:", round(RYAD[last] / RYAD[first], 3))
print("  база 2019:", round(rebased[last] / rebased[first], 3))
print("  рост не изменился:", round(RYAD[last] / RYAD[first], 3) == round(rebased[last] / rebased[first], 3))
