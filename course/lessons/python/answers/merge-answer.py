"""Решение задания урока 29: данные и справочник, соединённые по ключу.

Соединяем так, чтобы ни одна строка данных не пропала, считаем, скольким не
нашлось названия, и только потом группируем по региону.
"""

import pandas as pd

data = pd.DataFrame(
    [("KAZ", 2023, 14.5), ("KAZ", 2024, 8.7),
     ("UZB", 2023, 10.0), ("UZB", 2024, 9.6),
     ("RUS", 2023, 5.9), ("RUS", 2024, 8.4),
     ("KGZ", 2023, 10.8), ("KGZ", 2024, 6.3)],
    columns=["code", "year", "value"],
)
names = pd.DataFrame(
    [("KAZ", "Казахстан", "Центральная Азия"),
     ("UZB", "Узбекистан", "Центральная Азия"),
     ("KGZ", "Кыргызстан", "Центральная Азия"),
     ("TJK", "Таджикистан", "Центральная Азия")],
    columns=["code", "name", "region"],
)

# Слева — данные, и они не должны пропасть: справочник неполный.
joined = data.merge(names, on="code", how="left", validate="many_to_one")
print("строк было:", len(data), "| стало:", len(joined))
print("без названия:", sorted(joined.loc[joined["name"].isna(), "code"].unique()))

print()
print("средняя инфляция по странам:")
by_country = joined.groupby(["code", "name"], dropna=False)["value"].mean().round(2)
print(by_country.reset_index().to_string(index=False))

print()
print("по регионам:")
# Строки без региона в группировку не берём: регион у них не пустой, а
# неизвестный, и приписывать их к «Центральной Азии» было бы выдумкой.
by_region = joined.dropna(subset=["region"]).groupby("region").agg(
    стран=("code", "nunique"),
    средняя=("value", "mean"),
).round(2)
print(by_region.to_string())
