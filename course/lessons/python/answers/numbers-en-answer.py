"""The answer to lesson 3's exercise: one line of a receipt, counted out.

Whole loaves and change are what // and % are for; the money is printed with
two decimal places because that is how a price is read, and the string is put
together by an f-string rather than by adding a number to text.
"""

price = 260
count = 3
paid = 1000

total = price * count
print(f"for the item: {total} tenge")
print(f"change from {paid}: {paid - total} tenge")
print(f"whole loaves for {paid}: {paid // price}, left over {paid % price}")
print(f"the price of one in dollars: {price / 540:.2f}")
print(f"text and number together: {'bread'} × {count} = {total}")
