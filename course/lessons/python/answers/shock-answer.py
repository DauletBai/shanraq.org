"""Задание урока 42: чья это связь и сколько стран за ней стоит."""

import pandas as pd

# Всемирный банк, 2022 год: инфляция, изменение курса к доллару и рост
# денежной массы, в процентах. None — показателя нет, а не ноль.
SOSEDI = pd.DataFrame(
    [
        ("Армения", 8.64, -13.52, 16.10),
        ("Узбекистан", 11.45, 4.15, 30.17),
        ("Грузия", 11.90, -9.48, 11.04),
        ("Россия", 13.74, -7.02, None),
        ("Азербайджан", 13.85, 0.00, 23.60),
        ("Кыргызстан", 13.92, -0.62, 30.59),
        ("Казахстан", 15.03, 8.04, 13.94),
        ("Беларусь", 15.21, 3.44, None),
        ("Молдова", 28.74, 6.88, 5.27),
        ("Турция", 72.31, 86.98, 60.34),
    ],
    columns=["страна", "инфляция", "курс", "деньги"],
).set_index("страна")

COLUMNS = ["курс", "деньги"]


def svyaz(table, column, target="инфляция"):
    """Связь и число стран, на которых она посчитана.

    Одно без другого не печатают: 0.97 по десяти странам и 0.97 по трём —
    разные утверждения, а pandas молча выбрасывает пары с пропуском.
    """
    pair = table[[target, column]].dropna()
    return pair[target].corr(pair[column]), len(pair)


def bez_odnoy(table, column):
    """Связь, посчитанная без каждой страны по очереди."""
    out = {}
    for country in table.index:
        value, _ = svyaz(table.drop(country), column)
        out[country] = round(value, 3)
    return pd.Series(out)


print("== Связь и сколько стран за ней стоит")
base = {}
for column in COLUMNS:
    value, count = svyaz(SOSEDI, column)
    base[column] = value
    print(f"  {column:8} {value:+.3f}  по {count} странам")

print()
print("== Если убрать по одной стране")
table = pd.DataFrame({column: bez_odnoy(SOSEDI, column) for column in COLUMNS})
for column in COLUMNS:
    table[column + ", сдвиг"] = (table[column] - base[column]).round(3)
order = table["курс, сдвиг"].abs().sort_values(ascending=False).index
print(table.loc[order].to_string())

print()
print("== Кто решает ответ")
for column in COLUMNS:
    worst = (table[column] - base[column]).abs().idxmax()
    print(f"  {column:8} без страны «{worst}» связь идёт с {base[column]:+.3f}"
          f" на {table.loc[worst, column]:+.3f}")

print()
print("== Проверка")
signs = {column: (base[column] > 0) != (table[column] > 0) for column in COLUMNS}
for column in COLUMNS:
    flips = list(table.index[signs[column]])
    print(f"  {column:8} знак меняют страны: {', '.join(flips) if flips else 'никакие'}")
