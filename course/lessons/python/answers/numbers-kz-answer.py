"""The answer to lesson 3's exercise: one line of a receipt, counted out.

Whole loaves and change are what // and % are for; the money is printed with
two decimal places because that is how a price is read, and the string is put
together by an f-string rather than by adding a number to text.
"""

price = 260
count = 3
paid = 1000

total = price * count
print(f"тауар үшін: {total} теңге")
print(f"{paid}-нан қайтым: {paid - total} теңге")
print(f"{paid} теңгеге бүтін нан: {paid // price}, қалатыны {paid % price}")
print(f"біреуінің доллармен бағасы: {price / 540:.2f}")
print(f"жол мен сан бірге: {'нан'} × {count} = {total}")
