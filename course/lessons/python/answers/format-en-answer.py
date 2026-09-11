"""The exercise of lesson 35: a report that adds up.

Amounts in the local shape, changes with a sign, shares as per cent -- and a
total computed before any rounding.
"""

def money(value):
    """1234567.89 -> "1 234 567,89": a space in the thousands, a comma in the decimals."""
    return f"{value:,.2f}".replace(",", " ").replace(".", ",")


def percent(value):
    """A share -> "67,4%": the same comma as the amounts, or the table falls apart."""
    return f"{value:.1%}".replace(".", ",")


rows = [
    ("rent", 1_240_500.005, 62_500.0),
    ("food", 348_210.335, -12_400.0),
    ("fuel", 190_880.665, 4_300.0),
    ("phone", 59_940.995, 0.0),
]

total = sum(amount for _, amount, _ in rows)

print(f"{'item':<10}{'amount, tenge':>16}{'share':>8}{'change, tenge':>16}")
for name, amount, change in rows:
    share = amount / total
    print(f"{name:<10}{money(amount):>16}{percent(share):>8}{money(change) if change else '—':>16}")
print(f"{'total':<10}{money(total):>16}{percent(1):>8}")

print()
print("the check:")
print("  the total is computed before rounding:", round(total, 2) == round(sum(a for _, a, _ in rows), 2))
rounded = [round(amount, 2) for _, amount, _ in rows]
print("  the sum of the rounded rows:", money(sum(rounded)))
print("  the rounded total:          ", money(round(total, 2)))
print("  the gap:", round(sum(rounded) - total, 2))
print("  the shares add up to:", percent(sum(amount / total for _, amount, _ in rows)))
