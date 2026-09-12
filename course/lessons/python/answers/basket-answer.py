"""Задание урока 39: индекс своей корзины и сравнение с официальным."""

import pandas as pd

# Официальный индекс потребительских цен Казахстана (2010 = 100), Всемирный банк.
OFFICIAL = {2019: 189.30, 2024: 312.53}


def with_weights(basket):
    """Добавляет стоимости, веса и рост по каждой позиции."""
    out = basket.copy()
    out["тогда"] = out["сколько"] * out["было"]
    out["сейчас"] = out["сколько"] * out["стало"]
    out["вес"] = out["тогда"] / out["тогда"].sum()
    out["выросла в"] = out["стало"] / out["было"]
    out["вклад"] = out["вес"] * (out["выросла в"] - 1)
    return out


def basket_index(basket):
    """Индекс корзины: стоимость сегодня против стоимости тогда."""
    return basket["сейчас"].sum() / basket["тогда"].sum()


receipts = pd.DataFrame(
    [
        ("хлеб, буханка", 20, 90, 180),
        ("молоко, литр", 15, 260, 480),
        ("яйца, десяток", 8, 350, 1000),
        ("мясо, кг", 4, 1600, 3200),
        ("проезд, поездка", 40, 80, 100),
        ("интернет, месяц", 1, 5000, 6500),
        ("аренда, месяц", 1, 90000, 150000),
    ],
    columns=["позиция", "сколько", "было", "стало"],
)

basket = with_weights(receipts)
personal = basket_index(basket)
official = OFFICIAL[2024] / OFFICIAL[2019]

print("== Вклад в рост, по убыванию")
table = basket[["позиция", "вес", "выросла в", "вклад"]].sort_values("вклад", ascending=False)
print(table.round(3).to_string(index=False))

print()
print("== Итог")
print("  моя корзина:", round(personal, 3), "→", round((personal - 1) * 100, 1), "%")
print("  официальный индекс:", round(official, 3), "→", round((official - 1) * 100, 1), "%")
print("  разница:", round((personal - official) * 100, 1), "процентного пункта")

top = table.iloc[0]
print()
print("== Одна позиция объясняет")
print("  позиция:", top["позиция"])
print("  её вес:", round(top["вес"] * 100, 1), "% корзины")
print("  её вклад:", round(top["вклад"] / basket["вклад"].sum() * 100, 1), "% всего роста")

print()
print("== Проверка")
print("  сумма вкладов:", round(basket["вклад"].sum(), 4))
print("  рост корзины:", round(personal - 1, 4))
print("  сходится:", round(basket["вклад"].sum(), 6) == round(personal - 1, 6))
