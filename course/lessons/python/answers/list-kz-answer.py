"""The answer to lesson 5's exercise: one receipt as a list of prices.

The built-in functions take the whole list at once, so nothing here needs a
loop. sorted returns a new list and leaves the receipt as it was, which is the
difference the lesson is about.
"""

prices = [260, 620, 1890, 90, 310]
names = ["нан", "сүт", "май", "тұз", "жұмыртқа"]

print("позиция:", len(prices), "| біріншісі:", prices[0], "| соңғысы:", prices[-1])
print("чектің сомасы:", sum(prices))
print("ең қымбаты:", max(prices), "| ең арзаны:", min(prices))
print(f"орташа баға: {sum(prices) / len(prices):.2f}")
print("алғашқы үшеуі:", prices[:3])
print("өсу ретімен:", sorted(prices))
print("бастапқысы қозғалмаған:", prices)
print("атаулары:", names[:2], "…")
