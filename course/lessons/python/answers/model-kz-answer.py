"""43-сабақтың тапсырмасы: түзу, оның қатесі және «ештеңе істемеу» әдісімен салыстыру."""

import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_absolute_error, r2_score

CPI = pd.Series({
    2010: 100.00, 2011: 108.45, 2012: 114.09, 2013: 120.87, 2014: 129.15,
    2015: 137.77, 2016: 157.56, 2017: 169.28, 2018: 179.72, 2019: 189.30,
    2020: 202.02, 2021: 218.27, 2022: 251.07, 2023: 287.54, 2024: 312.53,
    2025: 348.12,
})


def fit(series, log=False):
    """Жылдар бойынша түзу: кәдімгісі немесе индекстің логарифмі бойынша.

    Болжамды қатармен бірдей бірлікте қайтарады — қатені тек қарап қана қоймай,
    салыстыруға болатындай.
    """
    X = series.index.to_numpy().reshape(-1, 1)
    y = np.log(series.to_numpy()) if log else series.to_numpy()
    model = LinearRegression().fit(X, y)
    guess = model.predict(X)
    return pd.Series(np.exp(guess) if log else guess, index=series.index), model


def runs(errors):
    """Қатар келген қателер таңбасын неше рет ауыстырады."""
    signs = ["+" if e > 0 else "-" for e in errors]
    return 1 + sum(1 for a, b in zip(signs, signs[1:]) if a != b), "".join(signs)


line, model = fit(CPI)
log_line, log_model = fit(CPI, log=True)
naive = CPI.shift(1)

print("== Жылдар бойынша түзу")
table = pd.DataFrame({"индекс": CPI, "түзу": line.round(1)})
table["қате"] = (table["индекс"] - table["түзу"]).round(1)
print(table.to_string())
print("  жылына қосады:", round(float(model.coef_[0]), 2), "тармақ")
print("  R²:", round(r2_score(CPI, line), 3))
worst = table["қате"].abs().idxmax()
print("  ең нашар жыл:", worst, "— қателік", abs(round(float(table.loc[worst, "қате"]), 1)), "тармақ")

count, signs = runs(table["қате"])
print()
print("== Қатенің пішіні")
print("  таңбалар:", signs)
print("  топтар:", count, "— кездейсоқ қатеде", round(len(signs) / 2 + 1))

print()
print("== Бір жылдарда салыстыру")
# Аңғал модель бірінші жылды болжай алмайды: оның алдыңғы жылы жоқ.
# Модельдерді бәрінде жауап бар жерде ғана салыстыруға болады.
same = naive.notna()
scores = pd.Series({
    "түзу": mean_absolute_error(CPI[same], line[same]),
    "логарифмдегі түзу": mean_absolute_error(CPI[same], log_line[same]),
    "өткен жылдағыдай": mean_absolute_error(CPI[same], naive[same]),
}).round(2)
print(scores.to_string())
print("  салыстырудағы жыл:", int(same.sum()))
print("  ең жақсысы:", scores.idxmin())
print("  түзу «ештеңе істемеуді» ұтады:",
      round(float(scores["өткен жылдағыдай"] - scores["түзу"]), 2), "тармаққа")
