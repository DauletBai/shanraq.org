"""Решение задания урока 35: отчёт, который сходится.

Суммы в местном виде, изменения со знаком, доли процентом — и итог, посчитанный
до округления, поэтому столбец складывается в точности в напечатанный итог.
"""

def money(value):
    """1234567.89 → «1 234 567,89»: пробел в тысячах, запятая в дробях."""
    return f"{value:,.2f}".replace(",", " ").replace(".", ",")


def percent(value):
    """Доля → «67,4%»: запятая та же, что и в суммах, иначе таблица разъедется."""
    return f"{value:.1%}".replace(".", ",")


rows = [
    ("аренда", 1_240_500.005, 62_500.0),
    ("еда", 348_210.335, -12_400.0),
    ("топливо", 190_880.665, 4_300.0),
    ("связь", 59_940.995, 0.0),
]

total = sum(amount for _, amount, _ in rows)

print(f"{'статья':<10}{'сумма, ₸':>16}{'доля':>8}{'изменение, ₸':>16}")
for name, amount, change in rows:
    share = amount / total
    print(f"{name:<10}{money(amount):>16}{percent(share):>8}{money(change) if change else '—':>16}")
print(f"{'итого':<10}{money(total):>16}{percent(1):>8}")

print()
print("проверка:")
print("  итог посчитан до округления:", round(total, 2) == round(sum(a for _, a, _ in rows), 2))
rounded = [round(amount, 2) for _, amount, _ in rows]
print("  сумма округлённых строк:", money(sum(rounded)))
print("  округлённый итог:       ", money(round(total, 2)))
print("  расхождение:", round(sum(rounded) - total, 2))
print("  доли складываются в:", percent(sum(amount / total for _, amount, _ in rows)))
