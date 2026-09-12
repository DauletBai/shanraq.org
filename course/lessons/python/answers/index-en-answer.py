"""Lesson 38 exercise: the index, the multiplier and a thousand under the mattress."""

import pandas as pd

rates = pd.Series(
    [6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
)


def price_index(rates, base=100.0):
    """A price index out of yearly rates: percentages become multipliers."""
    return ((1 + rates / 100).cumprod() * base).round(2)


def deflate(nominal, index, base=100.0):
    """Nominal sums, in the money of the base year."""
    return (nominal / index[nominal.index] * base).round(2)


index = price_index(rates)
total = float((1 + rates / 100).prod())

print("== The index, 2014 = 100")
print(index.to_string())

print()
print("== The ten years in total")
print("  multiplier:", round(total, 3))
print("  growth:", round((total - 1) * 100, 2), "%")
print("  the rates added up (never do this):", round(rates.sum(), 2), "%")
print("  yearly average:", round((total ** (1 / len(rates)) - 1) * 100, 2), "%")

print()
print("== A thousand of 2014")
print("  buys what it then bought for:", round(1000 / total, 2), "tenge")
print("  to buy the old basket you need:", round(1000 * total, 2), "tenge")

print()
print("== A salary that tripled")
pay = pd.Series({2016: 150000, 2020: 250000, 2024: 450000}, name="nominal")
real = deflate(pay, index)
table = pd.DataFrame({"nominal": pay, "in 2014 tenge": real})
table["real, % of 2016"] = (real / real[2016] * 100 - 100).round(2)
print(table.to_string())
print("  the nominal grew", round(pay[2024] / pay[2016], 2), "times")
print("  in real money:", round(real[2024] / real[2016], 2), "times")
