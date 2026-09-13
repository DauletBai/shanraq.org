"""Задание урока 43: прямая, её ошибка и сравнение с «ничего не делать»."""

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
    """Прямая по годам: обычная или по логарифму индекса.

    Возвращает предсказание в тех же единицах, что и ряд, — чтобы ошибку
    можно было сравнивать, а не только смотреть на неё.
    """
    X = series.index.to_numpy().reshape(-1, 1)
    y = np.log(series.to_numpy()) if log else series.to_numpy()
    model = LinearRegression().fit(X, y)
    guess = model.predict(X)
    return pd.Series(np.exp(guess) if log else guess, index=series.index), model


def runs(errors):
    """Сколько раз подряд идущие ошибки меняют знак."""
    signs = ["+" if e > 0 else "-" for e in errors]
    return 1 + sum(1 for a, b in zip(signs, signs[1:]) if a != b), "".join(signs)


line, model = fit(CPI)
log_line, log_model = fit(CPI, log=True)
naive = CPI.shift(1)

print("== Прямая по годам")
table = pd.DataFrame({"индекс": CPI, "прямая": line.round(1)})
table["ошибка"] = (table["индекс"] - table["прямая"]).round(1)
print(table.to_string())
print("  за год прибавляет:", round(float(model.coef_[0]), 2), "пункта")
print("  R²:", round(r2_score(CPI, line), 3))
worst = table["ошибка"].abs().idxmax()
print("  худший год:", worst, "— промах", abs(round(float(table.loc[worst, "ошибка"]), 1)), "пункта")

count, signs = runs(table["ошибка"])
print()
print("== Форма ошибки")
print("  знаки:", signs)
print("  групп:", count, "при", round(len(signs) / 2 + 1), "у случайной ошибки")

print()
print("== Сравнение на одних и тех же годах")
# Наивная модель не умеет предсказывать первый год: у него нет предыдущего.
# Сравнивать модели можно только там, где ответ есть у всех.
same = naive.notna()
scores = pd.Series({
    "прямая": mean_absolute_error(CPI[same], line[same]),
    "прямая в логарифмах": mean_absolute_error(CPI[same], log_line[same]),
    "как в прошлом году": mean_absolute_error(CPI[same], naive[same]),
}).round(2)
print(scores.to_string())
print("  лет в сравнении:", int(same.sum()))
print("  лучшая:", scores.idxmin())
print("  прямая обыгрывает «ничего не делать» на:",
      round(float(scores["как в прошлом году"] - scores["прямая"]), 2), "пункта")
