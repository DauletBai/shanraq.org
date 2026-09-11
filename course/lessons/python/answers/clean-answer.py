"""Решение задания урока 30: уборка с журналом.

Каждое действие оставляет след: сколько было, сколько стало, что именно
изменилось. Уборка без журнала — это данные, про которые никто уже не скажет,
что с ними сделали.
"""

import pandas as pd

rows = [
    ("Костанай", "еда", "4 200"),
    ("костанай ", "топливо", "18 500"),
    ("Рудный", "еда", "2 600"),
    ("Рудный", "еда", "2 600"),
    ("КОСТАНАЙ", "связь", "н/д"),
    ("Рудный", "аренда", "62 000"),
    ("Костанай", "аренда", "95 000"),
    ("Астана", "еда", "—"),
]
df = pd.DataFrame(rows, columns=["city", "kind", "amount"])

print("журнал уборки:")
print("  строк пришло:", len(df))
print("  городов пришло:", len(df["city"].unique()))

df["city"] = df["city"].str.strip().str.capitalize()
print("  городов стало:", len(df["city"].unique()), "—", sorted(df["city"].unique()))

df["amount"] = pd.to_numeric(df["amount"].str.replace(" ", "", regex=False), errors="coerce")
print("  не стали числами:", int(df["amount"].isna().sum()))

before = len(df)
df = df.drop_duplicates()
print(f"  дублей убрано: {before - len(df)}")

known = df.dropna(subset=["amount"])
print("  строк в отчёт:", len(known), "из", len(df))

print()
print("по городам:")
report = known.groupby("city").agg(
    чеков=("amount", "count"),
    всего=("amount", "sum"),
    средний=("amount", "mean"),
).round(0).sort_values("всего", ascending=False)
print(report.to_string())

print()
print("суммы сходятся:", int(report["всего"].sum()) == int(known["amount"].sum()))
