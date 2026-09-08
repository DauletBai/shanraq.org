"""The answer to lesson 4's exercise: three weights as they come off a receipt.

They arrive as text -- with spaces around them and a comma where a point should
be -- and one of them is missing altogether, which is a gap and not a zero. No
loops and no branches here: neither has been taught yet, and the exercise asks
for nothing they would make easier.
"""

first = " 1,15 "
second = "0,25"
third = None

a = float(first.strip().replace(",", "."))
b = float(second.strip().replace(",", "."))

print(f"первый: {a} кг, тип {type(a).__name__}")
print(f"второй: {b} кг, тип {type(b).__name__}")
print("третий пуст:", third is None)
zero = 0.0
print("а ноль — это значение:", zero is None)
print(f"всего {a + b:.2f} кг, в среднем {(a + b) / 2:.3f} кг")
