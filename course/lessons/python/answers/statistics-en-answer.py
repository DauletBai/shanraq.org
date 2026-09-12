"""Lesson 41 task: recover the weights, find half the rise, change the base."""

import pandas as pd

# The Bureau of National Statistics, CPI, August 2026: the rise since the start
# of the year and each division's share of it. The names are shortened.
DIVISIONS = [
    ("Food and drink", 4.9, 1.90),
    ("Alcohol and tobacco", 10.4, 0.15),
    ("Clothing and footwear", 6.6, 0.61),
    ("Housing, water, energy", 8.6, 0.83),
    ("Household goods", 6.0, 0.33),
    ("Health", 11.5, 0.67),
    ("Transport", 3.6, 0.32),
    ("Communications", 6.5, 0.31),
    ("Recreation and culture", 10.4, 0.35),
    ("Education", 2.7, 0.07),
    ("Restaurants and hotels", 6.1, 0.11),
    ("Insurance and finance", 3.0, 0.02),
    ("Personal care and other", 9.6, 0.68),
]
HEADLINE = 6.4

# The World Bank, FP.CPI.TOTL: the price index, 2010 = 100.
SERIES = pd.Series({2019: 189.30, 2020: 202.02, 2021: 218.27, 2022: 251.07,
                    2023: 287.54, 2024: 312.53, 2025: 348.12})


def weights(divisions):
    """Recovers the weights from the contributions: contribution = weight x rise."""
    out = pd.DataFrame(divisions, columns=["division", "rise", "contribution"])
    out["weight"] = (out["contribution"] / out["rise"] * 100).round(1)
    return out.sort_values("contribution", ascending=False).reset_index(drop=True)


def half(table):
    """How many divisions from the top it takes to cover half of the rise."""
    total = table["contribution"].sum()
    running = table["contribution"].cumsum()
    return int((running < total / 2).sum() + 1), total


def rebase(series, year):
    """The same series on another base: the anchor year becomes a hundred."""
    return (series / series[year] * 100).round(2)


table = weights(DIVISIONS)
print("== The weights, recovered from the contributions")
print(table[["division", "weight", "rise", "contribution"]].to_string(index=False))
print("  the weights add up to:", round(table["weight"].sum(), 1), "%")

count, total = half(table)
print()
print("== Half of the rise")
print("  divisions needed:", count)
print("  they are:", ", ".join(table["division"].head(count)))
print("  their share of the rise:",
      round(table["contribution"].head(count).sum() / total * 100, 1), "%")

print()
print("== The headline and the sum of the contributions")
print("  the contributions add up to:", round(total, 2))
print("  the headline:", HEADLINE)
print("  the gap:", round(abs(HEADLINE - total), 2), "points — rounding")

print()
print("== Changing the base")
rebased = rebase(SERIES, 2019)
first, last = SERIES.index[0], SERIES.index[-1]
print("  base 2010:", round(SERIES[last] / SERIES[first], 3))
print("  base 2019:", round(rebased[last] / rebased[first], 3))
print("  the growth did not change:",
      round(SERIES[last] / SERIES[first], 3) == round(rebased[last] / rebased[first], 3))
