"""30-сабақтың тапсырмасының шешімі: журналы бар тазалау.

Әр әрекет із қалдырады: қанша болды, қанша болды, нақты не өзгерді. Журналсыз
тазалау — бұл онымен не істелгенін ешкім айта алмайтын дерек.
"""

import pandas as pd

rows = [
    ("Қостанай", "тамақ", "4 200"),
    ("қостанай ", "жанармай", "18 500"),
    ("Рудный", "тамақ", "2 600"),
    ("Рудный", "тамақ", "2 600"),
    ("ҚОСТАНАЙ", "байланыс", "дерек жоқ"),
    ("Рудный", "жалдау", "62 000"),
    ("Қостанай", "жалдау", "95 000"),
    ("Астана", "тамақ", "—"),
]
df = pd.DataFrame(rows, columns=["city", "kind", "amount"])

print("тазалау журналы:")
print("  жол келді:", len(df))
print("  қала келді:", len(df["city"].unique()))

df["city"] = df["city"].str.strip().str.capitalize()
print("  қала болды:", len(df["city"].unique()), "—", sorted(df["city"].unique()))

df["amount"] = pd.to_numeric(df["amount"].str.replace(" ", "", regex=False), errors="coerce")
print("  санға айналмағаны:", int(df["amount"].isna().sum()))

before = len(df)
df = df.drop_duplicates()
print(f"  алынған қайталау: {before - len(df)}")

known = df.dropna(subset=["amount"])
print("  есепке кететін жол:", len(known), "барлығынан", len(df))

print()
print("қалалар бойынша:")
report = known.groupby("city").agg(
    түбіртек=("amount", "count"),
    барлығы=("amount", "sum"),
    орташа=("amount", "mean"),
).round(0).sort_values("барлығы", ascending=False)
print(report.to_string())

print()
print("сомалар сәйкес:", int(report["барлығы"].sum()) == int(known["amount"].sum()))
