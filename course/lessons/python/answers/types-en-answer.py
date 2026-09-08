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

print(f"the first: {a} kg, type {type(a).__name__}")
print(f"the second: {b} kg, type {type(b).__name__}")
print("the third is empty:", third is None)
zero = 0.0
print("and zero is a value:", zero is None)
print(f"{a + b:.2f} kg in all, {(a + b) / 2:.3f} kg on average")
