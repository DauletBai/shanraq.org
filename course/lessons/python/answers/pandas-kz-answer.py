"""24-сабақтың тапсырмасының шешімі: бір сұрақ, екі тәсілмен қойылған.

Әр санат бойынша — сома, орташа чек және жалпы сомадағы үлес. Алдымен қолмен,
содан кейін pandas арқылы, соңында салыстыру: жауаптар тиынға дейін сәйкес
келуге тиіс.
"""

import pandas as pd

KINDS = ["тамақ", "жанармай", "пәтерақы", "байланыс", "басқа"]
BASE = {"тамақ": 2500, "жанармай": 9000, "пәтерақы": 60000, "байланыс": 4000, "басқа": 1500}
N = 300_000

rows = [(KINDS[i % 5], BASE[KINDS[i % 5]] + (i * 37) % 900 - 450) for i in range(N)]

# Қолмен: бір өтуде екі сөздік, соңында бөлу.
total, count = {}, {}
for kind, amount in rows:
    total[kind] = total.get(kind, 0) + amount
    count[kind] = count.get(kind, 0) + 1
whole = sum(total.values())
by_hand = {
    kind: (total[kind], total[kind] / count[kind], total[kind] / whole * 100)
    for kind in total
}

# pandas: сол кесте және сол үш сан.
df = pd.DataFrame(rows, columns=["kind", "amount"])
table = df.groupby("kind")["amount"].agg(["sum", "mean"])
table["share"] = table["sum"] / table["sum"].sum() * 100
table = table.sort_values("sum", ascending=False)

print(f"{'санат':10} {'сома':>12} {'орташа':>10} {'үлес':>7}")
for kind, row in table.iterrows():
    print(f"{kind:10} {row['sum']:12.0f} {row['mean']:10.2f} {row['share']:6.1f}%")
print(f"{'барлығы':10} {table['sum'].sum():12.0f}")

same = all(
    round(by_hand[kind][0], 6) == round(row["sum"], 6)
    and round(by_hand[kind][1], 6) == round(row["mean"], 6)
    and round(by_hand[kind][2], 6) == round(row["share"], 6)
    for kind, row in table.iterrows()
)
print("сәйкес келеді:", "иә" if same else "жоқ")
