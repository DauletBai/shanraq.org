"""The answer to lesson 5's exercise: one receipt as a list of prices.

The built-in functions take the whole list at once, so nothing here needs a
loop. sorted returns a new list and leaves the receipt as it was, which is the
difference the lesson is about.
"""

prices = [260, 620, 1890, 90, 310]
names = ["bread", "milk", "butter", "salt", "eggs"]

print("items:", len(prices), "| first:", prices[0], "| last:", prices[-1])
print("the receipt adds up to:", sum(prices))
print("dearest:", max(prices), "| cheapest:", min(prices))
print(f"average price: {sum(prices) / len(prices):.2f}")
print("the first three:", prices[:3])
print("in ascending order:", sorted(prices))
print("the original is untouched:", prices)
print("names:", names[:2], "…")
