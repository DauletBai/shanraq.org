"""The answer to lesson 6's exercise: a price list as a dictionary.

A missing item is asked for with get rather than by the square brackets, so the
program answers instead of falling over -- and the answer says outright that
there is no such item.
"""

prices = {"хлеб": 280, "молоко": 620, "масло": 1890, "соль": 90}

print("записей:", len(prices))
print("цена молока:", prices["молоко"])
print("цена сахара:", prices.get("сахар"))
print("цена сахара с умолчанием:", prices.get("сахар", 0), "тенге")
print("сахар есть:", "сахар" in prices, "| соль есть:", "соль" in prices)

prices["сахар"] = 450
print("после добавления:", len(prices), "записей, сумма:", sum(prices.values()))
print("самая дорогая цена:", max(prices.values()))
print("названия:", list(prices.keys()))
