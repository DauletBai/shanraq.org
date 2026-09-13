"""44-сабақтың тапсырмасы: адал тексеру, ағып кету және қайта оқып кету."""

import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_absolute_error
from sklearn.model_selection import train_test_split

CPI = pd.Series({
    2010: 100.00, 2011: 108.45, 2012: 114.09, 2013: 120.87, 2014: 129.15,
    2015: 137.77, 2016: 157.56, 2017: 169.28, 2018: 179.72, 2019: 189.30,
    2020: 202.02, 2021: 218.27, 2022: 251.07, 2023: 287.54, 2024: 312.53,
    2025: 348.12,
})
YEARS = CPI.index.to_numpy()
VALUES = CPI.to_numpy()
CHECK = 4          # тексеруге соңғы неше жыл кетеді


def feature(years, degree=1):
    """Жылдар нөлге жақын әрі масштабта: үлкен сандардың дәрежесі нашар саналады."""
    return np.vander((years - 2010) / 10, degree + 1)


def uaqyt_boyynsha(years, values, check=CHECK):
    """Соңғы check жыл — тексеруге. Ешқандай кездейсоқтық: болашақ
    үйретуге түспеуге тиіс."""
    edge = np.sort(years)[-check]
    past = years < edge
    return years[past], values[past], years[~past], values[~past]


def qateler(fit_years, fit_values, ask_years, ask_values, degree=None):
    """Бір модельдің үйретудегі және тексерудегі қатесі."""
    if degree is None:
        model = LinearRegression().fit(feature(fit_years), np.log(fit_values))
        guess = lambda ask: np.exp(model.predict(feature(ask)))
    else:
        model = LinearRegression().fit(feature(fit_years, degree), fit_values)
        guess = lambda ask: model.predict(feature(ask, degree))
    return (mean_absolute_error(fit_values, guess(fit_years)),
            mean_absolute_error(ask_values, guess(ask_years)))


y_fit, v_fit, y_ask, v_ask = uaqyt_boyynsha(YEARS, VALUES)

print("== Уақыт бойынша кесу")
print("  үйрету:", [int(y) for y in y_fit])
print("  тексеру:", [int(y) for y in y_ask])
leak = sorted(set(y_ask) & set(y_fit))
print("  екі бөлікке де түскен жылдар:", leak if leak else "жоқ")

on_train, on_test = qateler(y_fit, v_fit, y_ask, v_ask)
print("  логарифм бойынша түзу: үйрету", round(on_train, 2), "— тексеру", round(on_test, 2))
print("  тексеру", round(on_test / on_train, 1), "есе нашар")

print()
print("== Көпмүшенің дәрежелері")
rows = []
for degree in (1, 2, 3, 4, 5):
    train_error, test_error = qateler(y_fit, v_fit, y_ask, v_ask, degree)
    rows.append({"дәреже": degree,
                 "үйретуде": round(train_error, 2),
                 "тексеруде": round(test_error, 2)})
table = pd.DataFrame(rows).set_index("дәреже")
print(table.to_string())
by_train = table["үйретуде"].idxmin()
by_test = table["тексеруде"].idxmin()
print("  үйрету бойынша таңдау:", by_train, "— тексеруде ол береді", table.loc[by_train, "тексеруде"])
print("  тексеру бойынша таңдау:", by_test, "— тексеруде ол береді", table.loc[by_test, "тексеруде"])
print("  қате таңдаудың бағасы:",
      round(float(table.loc[by_train, "тексеруде"] - table.loc[by_test, "тексеруде"]), 2), "тармақ")

print()
print("== Кездейсоқ кесу, үш әрекет")
for seed in (0, 1, 42):
    fit_years, ask_years, fit_values, ask_values = train_test_split(
        YEARS, VALUES, test_size=CHECK, random_state=seed)
    _, test_error = qateler(fit_years, fit_values, ask_years, ask_values)
    print(f"  seed {seed}: тексеру жылдары {[int(y) for y in sorted(ask_years)]}"
          f" — қате {round(test_error, 2)}")
print("  уақыт бойынша:", round(on_test, 2), "— әрі бұл жалғыз адал цифр")
