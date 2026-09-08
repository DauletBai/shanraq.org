"""The answer to lesson 5's exercise: one receipt as a list of prices.

The built-in functions take the whole list at once, so nothing here needs a
loop. sorted returns a new list and leaves the receipt as it was, which is the
difference the lesson is about.
"""

prices = [260, 620, 1890, 90, 310]
names = ["хлеб", "молоко", "масло", "соль", "яйца"]

print("позиций:", len(prices), "| первая:", prices[0], "| последняя:", prices[-1])
print("сумма чека:", sum(prices))
print("самое дорогое:", max(prices), "| самое дешёвое:", min(prices))
print(f"средняя цена: {sum(prices) / len(prices):.2f}")
print("три первые:", prices[:3])
print("по возрастанию:", sorted(prices))
print("исходный не тронут:", prices)
print("названия:", names[:2], "…")
