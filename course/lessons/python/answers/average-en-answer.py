"""Lesson 37 exercise: one measure per series, chosen by a rule and not by eye."""

import pandas as pd

kazakhstan = pd.Series(
    [6.85, 6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
)
year_2022 = pd.Series(
    {
        "Kazakhstan": 15.03, "Russia": 13.74, "Uzbekistan": 11.45,
        "Kyrgyzstan": 13.92, "Turkiye": 72.31, "Georgia": 11.90,
        "Armenia": 8.64, "Azerbaijan": 13.85, "Belarus": 15.21,
        "Moldova": 28.74,
    },
)


def summary(row):
    """The six numbers a series is described with, and not one more."""
    quarters = row.quantile([0.25, 0.75])
    return {
        "count": int(row.count()),
        "mean": round(row.mean(), 2),
        "median": round(row.median(), 2),
        "spread": round(row.std(), 2),
        "range": (row.min(), row.max()),
        "quarters": round(quarters[0.75] - quarters[0.25], 2),
    }


def headline(row):
    """A rule, not a taste: if the mean has moved away from the median by more
    than a tenth of the median, the series has a tail and the headline takes
    the median."""
    mean, median = row.mean(), row.median()
    if abs(mean - median) / median > 0.1:
        return "median", round(median, 2)
    return "mean", round(mean, 2)


# A third series: the same years without the spikes. The rule has to be able to
# say "mean" too, or it is not a rule but a habit in disguise.
calm = kazakhstan.loc[2017:2021]

for title, row in (
    ("Kazakhstan, 2014-2024", kazakhstan),
    ("Kazakhstan, 2017-2021", calm),
    ("Ten countries, 2022", year_2022),
):
    figures = summary(row)
    name, value = headline(row)
    print("==", title)
    print("  observations:", figures["count"])
    print("  mean:", figures["mean"], "| median:", figures["median"])
    print(f"  spread: {figures['spread']} | from {figures['range'][0]} to {figures['range'][1]}")
    print("  interquartile range:", figures["quarters"])
    print("  for the headline:", name, value, "%")
    print("  below the mean:", int((row < row.mean()).sum()), "of", len(row))
    print()

print("== Why not one measure for everything")
both = pd.DataFrame({
    "11 years": kazakhstan.describe(),
    "5 calm": calm.describe(),
    "10 countries": year_2022.describe(),
})
print(both.round(2).to_string())
