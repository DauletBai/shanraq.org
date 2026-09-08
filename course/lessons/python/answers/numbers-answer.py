"""The answer to lesson 3's exercise: one line of a receipt, counted out.

Whole loaves and change are what // and % are for; the money is printed with
two decimal places because that is how a price is read, and the string is put
together by an f-string rather than by adding a number to text.
"""

price = 260
count = 3
paid = 1000

total = price * count
print(f"за товар: {total} тенге")
print(f"сдача с {paid}: {paid - total} тенге")
print(f"целых буханок на {paid}: {paid // price}, останется {paid % price}")
print(f"цена одной в долларах: {price / 540:.2f}")
print(f"строка и число вместе: {'хлеб'} × {count} = {total}")
