"""39-сабақтың тапсырмасы: өз себетіңіздің индексі және ресмимен салыстыру."""

import pandas as pd

# Қазақстанның тұтыну бағасы индексі (2010 = 100), Дүниежүзілік банк.
OFFICIAL = {2019: 189.30, 2024: 312.53}


def with_weights(basket):
    """Әр позицияға құн, салмақ және өсім қосады."""
    out = basket.copy()
    out["сонда"] = out["қанша"] * out["болған"]
    out["қазір"] = out["қанша"] * out["болды"]
    out["салмақ"] = out["сонда"] / out["сонда"].sum()
    out["өсті, есе"] = out["болды"] / out["болған"]
    out["үлесі"] = out["салмақ"] * (out["өсті, есе"] - 1)
    return out


def basket_index(basket):
    """Себет индексі: бүгінгі құн мен сол кездегі құн."""
    return basket["қазір"].sum() / basket["сонда"].sum()


receipts = pd.DataFrame(
    [
        ("нан, бөлке", 20, 90, 180),
        ("сүт, литр", 15, 260, 480),
        ("жұмыртқа, ондық", 8, 350, 1000),
        ("ет, кг", 4, 1600, 3200),
        ("жол жүру, сапар", 40, 80, 100),
        ("интернет, ай", 1, 5000, 6500),
        ("жалдау, ай", 1, 90000, 150000),
    ],
    columns=["атауы", "қанша", "болған", "болды"],
)

basket = with_weights(receipts)
personal = basket_index(basket)
official = OFFICIAL[2024] / OFFICIAL[2019]

print("== Өсімге үлес, кемуі бойынша")
table = basket[["атауы", "салмақ", "өсті, есе", "үлесі"]].sort_values("үлесі", ascending=False)
print(table.round(3).to_string(index=False))

print()
print("== Қорытынды")
print("  менің себетім:", round(personal, 3), "→", round((personal - 1) * 100, 1), "%")
print("  ресми индекс:", round(official, 3), "→", round((official - 1) * 100, 1), "%")
print("  айырма:", round((personal - official) * 100, 1), "пайыздық тармақ")

top = table.iloc[0]
print()
print("== Бір позиция түсіндіреді")
print("  позиция:", top["атауы"])
print("  салмағы:", round(top["салмақ"] * 100, 1), "% себеттің")
print("  үлесі:", round(top["үлесі"] / basket["үлесі"].sum() * 100, 1), "% бүкіл өсімнің")

print()
print("== Тексеру")
print("  үлестер қосындысы:", round(basket["үлесі"].sum(), 4))
print("  себеттің өсімі:", round(personal - 1, 4))
print("  сәйкес келді:", round(basket["үлесі"].sum(), 6) == round(personal - 1, 6))
