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

print(f"біріншісі: {a} кг, типі {type(a).__name__}")
print(f"екіншісі: {b} кг, типі {type(b).__name__}")
print("үшіншісі бос:", third is None)
zero = 0.0
print("ал нөл — бұл мән:", zero is None)
print(f"барлығы {a + b:.2f} кг, орташа {(a + b) / 2:.3f} кг")
