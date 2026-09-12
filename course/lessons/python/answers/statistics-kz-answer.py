"""41-сабақтың тапсырмасы: салмақты қалпына келтіру, өсімнің жартысы, база ауысуы."""

import pandas as pd

# Ұлттық статистика бюросы, ТБИ, 2026 тамыз: жыл басынан бергі өсім және
# бөлімнің осы өсімге қосқан үлесі. Атаулар қысқартылған.
BOLIMDER = [
    ("Азық-түлік және сусын", 4.9, 1.90),
    ("Алкоголь және темекі", 10.4, 0.15),
    ("Киім және аяқ киім", 6.6, 0.61),
    ("Тұрғын үй, су, энергия", 8.6, 0.83),
    ("Үй мен тұрмыс", 6.0, 0.33),
    ("Денсаулық сақтау", 11.5, 0.67),
    ("Көлік", 3.6, 0.32),
    ("Байланыс", 6.5, 0.31),
    ("Демалыс және мәдениет", 10.4, 0.35),
    ("Білім беру", 2.7, 0.07),
    ("Мейрамхана мен қонақүй", 6.1, 0.11),
    ("Сақтандыру және қаржы", 3.0, 0.02),
    ("Жеке күтім және өзгесі", 9.6, 0.68),
]
TAQYRYP = 6.4

# Дүниежүзілік банк, FP.CPI.TOTL: баға индексі, 2010 = 100.
QATAR = pd.Series({2019: 189.30, 2020: 202.02, 2021: 218.27, 2022: 251.07,
                   2023: 287.54, 2024: 312.53, 2025: 348.12})


def salmaqtar(bolimder):
    """Салмақты үлестен қалпына келтіреді: үлес = салмақ × өсім."""
    out = pd.DataFrame(bolimder, columns=["бөлім", "өсім", "үлес"])
    out["салмақ"] = (out["үлес"] / out["өсім"] * 100).round(1)
    return out.sort_values("үлес", ascending=False).reset_index(drop=True)


def jarty(table):
    """Өсімнің жартысын жабу үшін жоғарыдан неше бөлім керек."""
    total = table["үлес"].sum()
    running = table["үлес"].cumsum()
    return int((running < total / 2).sum() + 1), total


def bazaga(series, year):
    """Сол қатар басқа базада: тірек жыл жүзге айналады."""
    return (series / series[year] * 100).round(2)


table = salmaqtar(BOLIMDER)
print("== Үлестен қалпына келтірілген салмақтар")
print(table[["бөлім", "салмақ", "өсім", "үлес"]].to_string(index=False))
print("  салмақтар қосындысы:", round(table["салмақ"].sum(), 1), "%")

count, total = jarty(table)
print()
print("== Өсімнің жартысы")
print("  бөлім жетті:", count)
print("  олар:", ", ".join(table["бөлім"].head(count)))
print("  өсімдегі үлесі:", round(table["үлес"].head(count).sum() / total * 100, 1), "%")

print()
print("== Тақырып пен үлестер қосындысы")
print("  үлестер қосындысы:", round(total, 2))
print("  тақырып:", TAQYRYP)
print("  айырма:", round(abs(TAQYRYP - total), 2), "тармақ — дөңгелектеу")

print()
print("== База ауысуы")
rebased = bazaga(QATAR, 2019)
first, last = QATAR.index[0], QATAR.index[-1]
print("  2010 базасы:", round(QATAR[last] / QATAR[first], 3))
print("  2019 базасы:", round(rebased[last] / rebased[first], 3))
print("  өсім өзгерген жоқ:", round(QATAR[last] / QATAR[first], 3) == round(rebased[last] / rebased[first], 3))
