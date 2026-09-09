"""The exercise of lesson 24: one question, asked two ways.

Per category: the sum, the average receipt and the share of the total. By hand
first, then through pandas, and a check at the end -- the two must agree down to
the last tiyn.
"""

import pandas as pd

KINDS = ["food", "fuel", "rent", "phone", "other"]
BASE = {"food": 2500, "fuel": 9000, "rent": 60000, "phone": 4000, "other": 1500}
N = 300_000

rows = [(KINDS[i % 5], BASE[KINDS[i % 5]] + (i * 37) % 900 - 450) for i in range(N)]

# By hand: two dictionaries in one pass, and the division at the end.
total, count = {}, {}
for kind, amount in rows:
    total[kind] = total.get(kind, 0) + amount
    count[kind] = count.get(kind, 0) + 1
whole = sum(total.values())
by_hand = {
    kind: (total[kind], total[kind] / count[kind], total[kind] / whole * 100)
    for kind in total
}

# pandas: the same table and the same three numbers.
df = pd.DataFrame(rows, columns=["kind", "amount"])
table = df.groupby("kind")["amount"].agg(["sum", "mean"])
table["share"] = table["sum"] / table["sum"].sum() * 100
table = table.sort_values("sum", ascending=False)

print(f"{'category':9} {'sum':>12} {'average':>10} {'share':>7}")
for kind, row in table.iterrows():
    print(f"{kind:9} {row['sum']:12.0f} {row['mean']:10.2f} {row['share']:6.1f}%")
print(f"{'total':9} {table['sum'].sum():12.0f}")

same = all(
    round(by_hand[kind][0], 6) == round(row["sum"], 6)
    and round(by_hand[kind][1], 6) == round(row["mean"], 6)
    and round(by_hand[kind][2], 6) == round(row["share"], 6)
    for kind, row in table.iterrows()
)
print("they agree:", "yes" if same else "no")
