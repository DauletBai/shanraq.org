"""Lesson 39 task: the index of your own basket against the official one."""

import pandas as pd

# Kazakhstan consumer price index (2010 = 100), World Bank.
OFFICIAL = {2019: 189.30, 2024: 312.53}


def with_weights(basket):
    """Adds costs, weights and growth for every item."""
    out = basket.copy()
    out["then"] = out["count"] * out["was"]
    out["today"] = out["count"] * out["now"]
    out["weight"] = out["then"] / out["then"].sum()
    out["times"] = out["now"] / out["was"]
    out["contribution"] = out["weight"] * (out["times"] - 1)
    return out


def basket_index(basket):
    """The basket index: what it costs today against what it cost then."""
    return basket["today"].sum() / basket["then"].sum()


receipts = pd.DataFrame(
    [
        ("bread, loaf", 20, 90, 180),
        ("milk, litre", 15, 260, 480),
        ("eggs, ten", 8, 350, 1000),
        ("meat, kg", 4, 1600, 3200),
        ("transit, ride", 40, 80, 100),
        ("internet, month", 1, 5000, 6500),
        ("rent, month", 1, 90000, 150000),
    ],
    columns=["item", "count", "was", "now"],
)

basket = with_weights(receipts)
personal = basket_index(basket)
official = OFFICIAL[2024] / OFFICIAL[2019]

print("== Contribution to growth, largest first")
table = basket[["item", "weight", "times", "contribution"]].sort_values("contribution", ascending=False)
print(table.round(3).to_string(index=False))

print()
print("== Result")
print("  my basket:", round(personal, 3), "→", round((personal - 1) * 100, 1), "%")
print("  official index:", round(official, 3), "→", round((official - 1) * 100, 1), "%")
print("  difference:", round((personal - official) * 100, 1), "percentage points")

top = table.iloc[0]
print()
print("== One item explains")
print("  item:", top["item"])
print("  its weight:", round(top["weight"] * 100, 1), "% of the basket")
print("  its share:", round(top["contribution"] / basket["contribution"].sum() * 100, 1), "% of all growth")

print()
print("== Check")
print("  contributions add up to:", round(basket["contribution"].sum(), 4))
print("  basket growth:", round(personal - 1, 4))
print("  matches:", round(basket["contribution"].sum(), 6) == round(personal - 1, 6))
