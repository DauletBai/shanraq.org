"""Lesson 46 task: a model card, and the border it will not answer beyond."""

import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression

KAZAKHSTAN = pd.Series({
    2010: 100.00, 2011: 108.45, 2012: 114.09, 2013: 120.87, 2014: 129.15,
    2015: 137.77, 2016: 157.56, 2017: 169.28, 2018: 179.72, 2019: 189.30,
    2020: 202.02, 2021: 218.27, 2022: 251.07, 2023: 287.54, 2024: 312.53,
    2025: 348.12,
})
CUT = 2021


def trained_on(series):
    """A line through the log of the index, and the years it learned anything in."""
    years = series.index.to_numpy()
    model = LinearRegression().fit(((years - 2010) / 10).reshape(-1, 1),
                                   np.log(series.to_numpy()))
    known = (int(years.min()), int(years.max()))

    def ask(year, outside=False):
        """The model's answer. Past the edge of the data only if asked outright."""
        if not outside and not known[0] <= year <= known[1]:
            return None
        return float(np.exp(model.predict(np.array([[(year - 2010) / 10]])))[0])

    return ask, known, model


ask, known, model = trained_on(KAZAKHSTAN[KAZAKHSTAN.index <= CUT])
check = KAZAKHSTAN[KAZAKHSTAN.index > CUT]
errors = (check - np.array([ask(y, outside=True) for y in check.index])).abs()

print("== The model card")
print(f"  trained on: {known[0]}-{known[1]}, {known[1] - known[0] + 1} points")
print(f"  checked on: {int(check.index.min())}-{int(check.index.max())},"
      f" error from {errors.min():.0f} to {errors.max():.0f} points")
print(f"  features on the input: {model.n_features_in_}")
print(f"  in force for: {known[0]}-{known[1]} — past that a guess begins")

print()
print("== Answers inside and outside")
for year in (2015, 2021, 2025, 2030):
    answer = ask(year)
    if answer is None:
        print(f"  {year}: past the edge of the data, the model saw {known[0]}-{known[1]}")
    else:
        print(f"  {year}: {answer:.1f}")

print()
print("== What it would say without the border")
for year in (2025, 2030, 2100):
    print(f"  {year}: {ask(year, outside=True):.1f}")
print("  the same numbers, only now it shows that this is extrapolation")

print()
print("== An honest line for a report")
print(f"  price index, model trained on {known[0]}-{known[1]},"
      f" error on the check {errors.min():.0f}-{errors.max():.0f} points,"
      f" not in force past {known[1]}")
