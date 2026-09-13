"""Lesson 43 task: a line, its error, and a comparison with doing nothing."""

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
    """A line over the years: the plain one, or one over the log of the index.

    It returns the prediction in the units of the series, so that the error can
    be compared rather than only looked at.
    """
    X = series.index.to_numpy().reshape(-1, 1)
    y = np.log(series.to_numpy()) if log else series.to_numpy()
    model = LinearRegression().fit(X, y)
    guess = model.predict(X)
    return pd.Series(np.exp(guess) if log else guess, index=series.index), model


def runs(errors):
    """How many times the errors in a row change sign."""
    signs = ["+" if e > 0 else "-" for e in errors]
    return 1 + sum(1 for a, b in zip(signs, signs[1:]) if a != b), "".join(signs)


line, model = fit(CPI)
log_line, log_model = fit(CPI, log=True)
naive = CPI.shift(1)

print("== The line over the years")
table = pd.DataFrame({"index": CPI, "line": line.round(1)})
table["error"] = (table["index"] - table["line"]).round(1)
print(table.to_string())
print("  it gains a year:", round(float(model.coef_[0]), 2), "points")
print("  R²:", round(r2_score(CPI, line), 3))
worst = table["error"].abs().idxmax()
print("  the worst year:", worst, "— off by", abs(round(float(table.loc[worst, "error"]), 1)), "points")

count, signs = runs(table["error"])
print()
print("== The shape of the error")
print("  the signs:", signs)
print("  runs:", count, "against", round(len(signs) / 2 + 1), "for a random error")

print()
print("== Compared on the same years")
# The naive model cannot predict the first year: it has no year before it.
# Models can only be compared where every one of them has an answer.
same = naive.notna()
scores = pd.Series({
    "the line": mean_absolute_error(CPI[same], line[same]),
    "the line in logarithms": mean_absolute_error(CPI[same], log_line[same]),
    "same as last year": mean_absolute_error(CPI[same], naive[same]),
}).round(2)
print(scores.to_string())
print("  years compared:", int(same.sum()))
print("  the best:", scores.idxmin())
print("  the line beats doing nothing by:",
      round(float(scores["same as last year"] - scores["the line"]), 2), "points")
