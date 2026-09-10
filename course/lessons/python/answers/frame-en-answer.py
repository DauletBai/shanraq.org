"""The exercise of lesson 25: labels instead of numbers.

The peak of inflation in each country, how many times prices grew over four
years, and the calmest year -- every one of them a question put to the labels
rather than to the positions.
"""

import pandas as pd

data = {
    "KZ": [8.0, 15.0, 14.5, 8.7],
    "UZ": [10.8, 11.4, 10.0, 9.6],
    "RU": [6.7, 13.7, 5.9, 8.4],
}
df = pd.DataFrame(data, index=[2021, 2022, 2023, 2024])
df.index.name = "year"

print("the table:")
print(df)

print()
print("the peak per country:")
for country in df.columns:
    year = df[country].idxmax()
    print(f"  {country}: {year} ({df.loc[year, country]})")

print()
print("prices over four years grew:")
growth = (1 + df / 100).prod()
for country, times in growth.items():
    print(f"  {country}: {times:.2f} times")

print()
average = df.mean(axis=1).round(2)
calm = average.idxmin()
print(f"the calmest year: {calm} (average {average[calm]:.2f})")
