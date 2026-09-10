"""Решение задания урока 25: подписи вместо номеров.

Пик инфляции по каждой стране, во сколько раз выросли цены за четыре года и
самый спокойный год — всё это вопросы к подписям, а не к позициям.
"""

import pandas as pd

data = {
    "KZ": [8.0, 15.0, 14.5, 8.7],
    "UZ": [10.8, 11.4, 10.0, 9.6],
    "RU": [6.7, 13.7, 5.9, 8.4],
}
df = pd.DataFrame(data, index=[2021, 2022, 2023, 2024])
df.index.name = "год"

print("таблица:")
print(df)

print()
print("пик по странам:")
for country in df.columns:
    year = df[country].idxmax()
    print(f"  {country}: {year} ({df.loc[year, country]})")

print()
print("цены за четыре года выросли:")
growth = (1 + df / 100).prod()
for country, times in growth.items():
    print(f"  {country}: в {times:.2f} раза")

print()
average = df.mean(axis=1).round(2)
calm = average.idxmin()
print(f"самый спокойный год: {calm} (среднее {average[calm]:.2f})")
