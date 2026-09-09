"""Решение задания урока 24: один вопрос, заданный двумя способами.

По каждой категории — сумма, средний чек и доля в общей сумме. Сначала руками,
потом через pandas, и сверка: ответы обязаны совпасть до копейки.
"""

import pandas as pd

KINDS = ["еда", "топливо", "аренда", "связь", "прочее"]
BASE = {"еда": 2500, "топливо": 9000, "аренда": 60000, "связь": 4000, "прочее": 1500}
N = 300_000

rows = [(KINDS[i % 5], BASE[KINDS[i % 5]] + (i * 37) % 900 - 450) for i in range(N)]

# Руками: два словаря за один проход и деление в конце.
total, count = {}, {}
for kind, amount in rows:
    total[kind] = total.get(kind, 0) + amount
    count[kind] = count.get(kind, 0) + 1
whole = sum(total.values())
by_hand = {
    kind: (total[kind], total[kind] / count[kind], total[kind] / whole * 100)
    for kind in total
}

# pandas: та же таблица и те же три числа.
df = pd.DataFrame(rows, columns=["kind", "amount"])
table = df.groupby("kind")["amount"].agg(["sum", "mean"])
table["share"] = table["sum"] / table["sum"].sum() * 100
table = table.sort_values("sum", ascending=False)

print(f"{'категория':10} {'сумма':>12} {'средний':>10} {'доля':>7}")
for kind, row in table.iterrows():
    print(f"{kind:10} {row['sum']:12.0f} {row['mean']:10.2f} {row['share']:6.1f}%")
print(f"{'всего':10} {table['sum'].sum():12.0f}")

same = all(
    round(by_hand[kind][0], 6) == round(row["sum"], 6)
    and round(by_hand[kind][1], 6) == round(row["mean"], 6)
    and round(by_hand[kind][2], 6) == round(row["share"], 6)
    for kind, row in table.iterrows()
)
print("совпадают:", "да" if same else "нет")
