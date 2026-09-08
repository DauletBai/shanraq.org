"""The answer to lesson 6's exercise: a price list as a dictionary.

A missing item is asked for with get rather than by the square brackets, so the
program answers instead of falling over -- and the answer says outright that
there is no such item.
"""

prices = {"bread": 280, "milk": 620, "butter": 1890, "salt": 90}

print("entries:", len(prices))
print("the price of milk:", prices["milk"])
print("the price of sugar:", prices.get("sugar"))
print("the price of sugar with a default:", prices.get("sugar", 0), "tenge")
print("sugar is there:", "sugar" in prices, "| salt is there:", "salt" in prices)

prices["sugar"] = 450
print("after adding:", len(prices), "entries, sum:", sum(prices.values()))
print("the dearest price:", max(prices.values()))
print("names:", list(prices.keys()))
