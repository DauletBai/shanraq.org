"""Lesson 44 task: an honest check, a leak, and overfitting."""

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
CHECK = 4          # how many of the last years go to the check


def feature(years, degree=1):
    """Years near zero and to scale: powers of large numbers count badly."""
    return np.vander((years - 2010) / 10, degree + 1)


def split_by_time(years, values, check=CHECK):
    """The last check years go to the check. No randomness: the future must
    not end up in the training."""
    edge = np.sort(years)[-check]
    past = years < edge
    return years[past], values[past], years[~past], values[~past]


def errors(fit_years, fit_values, ask_years, ask_values, degree=None):
    """The error on training and on the check for one model."""
    if degree is None:
        model = LinearRegression().fit(feature(fit_years), np.log(fit_values))
        guess = lambda ask: np.exp(model.predict(feature(ask)))
    else:
        model = LinearRegression().fit(feature(fit_years, degree), fit_values)
        guess = lambda ask: model.predict(feature(ask, degree))
    return (mean_absolute_error(fit_values, guess(fit_years)),
            mean_absolute_error(ask_values, guess(ask_years)))


y_fit, v_fit, y_ask, v_ask = split_by_time(YEARS, VALUES)

print("== The cut by time")
print("  training:", [int(y) for y in y_fit])
print("  the check:", [int(y) for y in y_ask])
leak = sorted(set(y_ask) & set(y_fit))
print("  years in both parts:", leak if leak else "none")

on_train, on_test = errors(y_fit, v_fit, y_ask, v_ask)
print("  the line in logarithms: training", round(on_train, 2), "— check", round(on_test, 2))
print("  the check is worse by a factor of", round(on_test / on_train, 1))

print()
print("== The degrees of a polynomial")
rows = []
for degree in (1, 2, 3, 4, 5):
    train_error, test_error = errors(y_fit, v_fit, y_ask, v_ask, degree)
    rows.append({"degree": degree,
                 "on training": round(train_error, 2),
                 "on the check": round(test_error, 2)})
table = pd.DataFrame(rows).set_index("degree")
print(table.to_string())
by_train = table["on training"].idxmin()
by_test = table["on the check"].idxmin()
print("  chosen by training:", by_train, "— on the check it gives", table.loc[by_train, "on the check"])
print("  chosen by the check:", by_test, "— on the check it gives", table.loc[by_test, "on the check"])
print("  the price of the wrong choice:",
      round(float(table.loc[by_train, "on the check"] - table.loc[by_test, "on the check"]), 2), "points")

print()
print("== A random cut, three tries")
for seed in (0, 1, 42):
    fit_years, ask_years, fit_values, ask_values = train_test_split(
        YEARS, VALUES, test_size=CHECK, random_state=seed)
    _, test_error = errors(fit_years, fit_values, ask_years, ask_values)
    print(f"  seed {seed}: check years {[int(y) for y in sorted(ask_years)]}"
          f" — error {round(test_error, 2)}")
print("  by time:", round(on_test, 2), "— and that is the only honest figure")
