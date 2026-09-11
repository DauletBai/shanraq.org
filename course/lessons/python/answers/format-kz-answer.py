"""35-сабақтың тапсырмасының шешімі: қорытындысы сәйкес келетін есеп.

Сомалар жергілікті түрде, өзгерістер таңбасымен, үлестер пайызбен — әрі
қорытынды дөңгелектеуге дейін есептелген.
"""

def money(value):
    """1234567.89 → «1 234 567,89»: мыңдықта бос орын, бөлшекте үтір."""
    return f"{value:,.2f}".replace(",", " ").replace(".", ",")


def percent(value):
    """Үлес → «67,4%»: үтір сомалардағыдай, әйтпесе кесте бұзылады."""
    return f"{value:.1%}".replace(".", ",")


rows = [
    ("жалдау", 1_240_500.005, 62_500.0),
    ("тамақ", 348_210.335, -12_400.0),
    ("жанармай", 190_880.665, 4_300.0),
    ("байланыс", 59_940.995, 0.0),
]

total = sum(amount for _, amount, _ in rows)

print(f"{'бап':<10}{'сома, ₸':>16}{'үлес':>8}{'өзгеріс, ₸':>16}")
for name, amount, change in rows:
    share = amount / total
    print(f"{name:<10}{money(amount):>16}{percent(share):>8}{money(change) if change else '—':>16}")
print(f"{'барлығы':<10}{money(total):>16}{percent(1):>8}")

print()
print("тексеру:")
print("  қорытынды дөңгелектеуге дейін есептелген:", round(total, 2) == round(sum(a for _, a, _ in rows), 2))
rounded = [round(amount, 2) for _, amount, _ in rows]
print("  дөңгелектелген жолдардың сомасы:", money(sum(rounded)))
print("  дөңгелектелген қорытынды:       ", money(round(total, 2)))
print("  айырма:", round(sum(rounded) - total, 2))
print("  үлестер қосындысы:", percent(sum(amount / total for _, amount, _ in rows)))
