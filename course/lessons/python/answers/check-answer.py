"""Задание урока 44: честная проверка, утечка и переобучение."""

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
CHECK = 4          # сколько последних лет уходит на проверку


def feature(years, degree=1):
    """Годы у нуля и в масштабе: степени больших чисел считаются плохо."""
    return np.vander((years - 2010) / 10, degree + 1)


def split_by_time(years, values, check=CHECK):
    """Последние check лет — проверочные. Никакой случайности: будущее
    не должно попасть в обучение."""
    edge = np.sort(years)[-check]
    past = years < edge
    return years[past], values[past], years[~past], values[~past]


def errors(fit_years, fit_values, ask_years, ask_values, degree=None):
    """Ошибка на обучении и на проверке для одной модели."""
    if degree is None:
        model = LinearRegression().fit(feature(fit_years), np.log(fit_values))
        guess = lambda ask: np.exp(model.predict(feature(ask)))
    else:
        model = LinearRegression().fit(feature(fit_years, degree), fit_values)
        guess = lambda ask: model.predict(feature(ask, degree))
    return (mean_absolute_error(fit_values, guess(fit_years)),
            mean_absolute_error(ask_values, guess(ask_years)))


y_fit, v_fit, y_ask, v_ask = split_by_time(YEARS, VALUES)

print("== Разрез по времени")
print("  обучение:", [int(y) for y in y_fit])
print("  проверка:", [int(y) for y in y_ask])
leak = sorted(set(y_ask) & set(y_fit))
print("  годы, попавшие в обе части:", leak if leak else "нет")

on_train, on_test = errors(y_fit, v_fit, y_ask, v_ask)
print("  прямая по логарифму: обучение", round(on_train, 2), "— проверка", round(on_test, 2))
print("  проверка хуже в", round(on_test / on_train, 1), "раза")

print()
print("== Степени многочлена")
rows = []
for degree in (1, 2, 3, 4, 5):
    train_error, test_error = errors(y_fit, v_fit, y_ask, v_ask, degree)
    rows.append({"степень": degree,
                 "на обучении": round(train_error, 2),
                 "на проверке": round(test_error, 2)})
table = pd.DataFrame(rows).set_index("степень")
print(table.to_string())
by_train = table["на обучении"].idxmin()
by_test = table["на проверке"].idxmin()
print("  выбор по обучению:", by_train, "— на проверке она даёт", table.loc[by_train, "на проверке"])
print("  выбор по проверке:", by_test, "— на проверке она даёт", table.loc[by_test, "на проверке"])
print("  цена неправильного выбора:",
      round(float(table.loc[by_train, "на проверке"] - table.loc[by_test, "на проверке"]), 2), "пункта")

print()
print("== Случайный разрез, три попытки")
for seed in (0, 1, 42):
    fit_years, ask_years, fit_values, ask_values = train_test_split(
        YEARS, VALUES, test_size=CHECK, random_state=seed)
    _, test_error = errors(fit_years, fit_values, ask_years, ask_values)
    print(f"  seed {seed}: проверочные годы {[int(y) for y in sorted(ask_years)]}"
          f" — ошибка {round(test_error, 2)}")
print("  по времени:", round(on_test, 2), "— и это единственная честная цифра")
