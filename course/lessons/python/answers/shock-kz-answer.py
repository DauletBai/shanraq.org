"""42-сабақтың тапсырмасы: бұл байланыс кімдікі және артында неше ел тұр."""

import pandas as pd

# Дүниежүзілік банк, 2022 жыл: инфляция, доллар бағамының өзгерісі және ақша
# массасының өсімі, пайызбен. None — көрсеткіш жоқ дегені, нөл емес.
KORSHILER = pd.DataFrame(
    [
        ("Армения", 8.64, -13.52, 16.10),
        ("Өзбекстан", 11.45, 4.15, 30.17),
        ("Грузия", 11.90, -9.48, 11.04),
        ("Ресей", 13.74, -7.02, None),
        ("Әзірбайжан", 13.85, 0.00, 23.60),
        ("Қырғызстан", 13.92, -0.62, 30.59),
        ("Қазақстан", 15.03, 8.04, 13.94),
        ("Беларусь", 15.21, 3.44, None),
        ("Молдова", 28.74, 6.88, 5.27),
        ("Түркия", 72.31, 86.98, 60.34),
    ],
    columns=["ел", "инфляция", "бағам", "ақша"],
).set_index("ел")

COLUMNS = ["бағам", "ақша"]


def baylanys(table, column, target="инфляция"):
    """Байланыс және ол қанша елде саналғаны.

    Бірінсіз екіншісі басылмайды: он ел бойынша 0.97 мен үш ел бойынша 0.97 —
    әртүрлі тұжырым, ал pandas бос орны бар жұптарды үнсіз тастап кетеді.
    """
    pair = table[[target, column]].dropna()
    return pair[target].corr(pair[column]), len(pair)


def bir_elsiz(table, column):
    """Әр елді кезекпен алып тастап саналған байланыс."""
    out = {}
    for country in table.index:
        value, _ = baylanys(table.drop(country), column)
        out[country] = round(value, 3)
    return pd.Series(out)


print("== Байланыс және артында неше ел тұр")
base = {}
for column in COLUMNS:
    value, count = baylanys(KORSHILER, column)
    base[column] = value
    print(f"  {column:8} {value:+.3f}  {count} ел бойынша")

print()
print("== Әр елді кезекпен алып тастасақ")
table = pd.DataFrame({column: bir_elsiz(KORSHILER, column) for column in COLUMNS})
for column in COLUMNS:
    shift = (table[column] - base[column]).round(3)
    # Таңбасы минус нөл «-0.000» болып басылады, ал нөлден тек кітапхана ішіндегі
    # амалдар ретімен ерекшеленеді. Уәде етілген шығыс оған тәуелді болмауға тиіс.
    table[column + ", жылжу"] = shift.where(shift != 0, 0.0)
# Сұрыптау тұрақты: жылжу бірнеше елде бірдей, ал кәдімгі сұрыптауда олардың
# реті машинаға тәуелді болады, сабақ болса дәл шығысты уәде етеді.
order = table["бағам, жылжу"].abs().sort_values(ascending=False, kind="stable").index
print(table.loc[order].to_string())

print()
print("== Жауапты кім шешеді")
for column in COLUMNS:
    worst = (table[column] - base[column]).abs().idxmax()
    print(f"  {column:8} «{worst}» елісіз байланыс {base[column]:+.3f} дегеннен"
          f" {table.loc[worst, column]:+.3f} дегенге барады")

print()
print("== Тексеру")
signs = {column: (base[column] > 0) != (table[column] > 0) for column in COLUMNS}
for column in COLUMNS:
    flips = list(table.index[signs[column]])
    print(f"  {column:8} таңбаны өзгертетін елдер: {', '.join(flips) if flips else 'ешқайсысы'}")
