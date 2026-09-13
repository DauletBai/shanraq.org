"""Lesson 42 task: whose link this is, and how many countries stand behind it."""

import pandas as pd

# The World Bank, 2022: inflation, the change in the rate against the dollar
# and the growth of broad money, in per cent. None means no figure, not zero.
NEIGHBOURS = pd.DataFrame(
    [
        ("Armenia", 8.64, -13.52, 16.10),
        ("Uzbekistan", 11.45, 4.15, 30.17),
        ("Georgia", 11.90, -9.48, 11.04),
        ("Russia", 13.74, -7.02, None),
        ("Azerbaijan", 13.85, 0.00, 23.60),
        ("Kyrgyzstan", 13.92, -0.62, 30.59),
        ("Kazakhstan", 15.03, 8.04, 13.94),
        ("Belarus", 15.21, 3.44, None),
        ("Moldova", 28.74, 6.88, 5.27),
        ("Turkiye", 72.31, 86.98, 60.34),
    ],
    columns=["country", "inflation", "rate", "money"],
).set_index("country")

COLUMNS = ["rate", "money"]


def link(table, column, target="inflation"):
    """The link, and the number of countries it was counted on.

    One is never printed without the other: 0.97 over ten countries and 0.97
    over three are different statements, and pandas drops the pairs with a gap
    without saying so.
    """
    pair = table[[target, column]].dropna()
    return pair[target].corr(pair[column]), len(pair)


def without_each(table, column):
    """The link counted without each country in turn."""
    out = {}
    for country in table.index:
        value, _ = link(table.drop(country), column)
        out[country] = round(value, 3)
    return pd.Series(out)


print("== The link and how many countries stand behind it")
base = {}
for column in COLUMNS:
    value, count = link(NEIGHBOURS, column)
    base[column] = value
    print(f"  {column:8} {value:+.3f}  over {count} countries")

print()
print("== Taking out one country at a time")
table = pd.DataFrame({column: without_each(NEIGHBOURS, column) for column in COLUMNS})
for column in COLUMNS:
    shift = (table[column] - base[column]).round(3)
    # A zero with a minus sign prints as "-0.000" and differs from a zero only
    # by the order of operations inside the library. A promised output must not
    # depend on that.
    table[column + ", shift"] = shift.where(shift != 0, 0.0)
# A stable sort: several countries share a shift, and with an ordinary sort
# their order depends on the machine, while the lesson promises an exact output.
order = table["rate, shift"].abs().sort_values(ascending=False, kind="stable").index
print(table.loc[order].to_string())

print()
print("== Who decides the answer")
for column in COLUMNS:
    worst = (table[column] - base[column]).abs().idxmax()
    print(f"  {column:8} without {worst} the link goes from {base[column]:+.3f}"
          f" to {table.loc[worst, column]:+.3f}")

print()
print("== The check")
signs = {column: (base[column] > 0) != (table[column] > 0) for column in COLUMNS}
for column in COLUMNS:
    flips = list(table.index[signs[column]])
    print(f"  {column:8} countries that flip the sign: {', '.join(flips) if flips else 'none'}")
