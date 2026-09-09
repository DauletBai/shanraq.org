"""The answer to lesson 15's exercise: two exports, and what changed between them.

Every set is printed through sorted: a set has no order of its own, and printing
one directly gives an answer that can come out differently on another machine.
The duplicates disappear on the way in, which is the point -- the question is
what was on the list, not how many times it was written down.
"""

january = ["хлеб", "молоко", "масло", "соль", "хлеб", "яйца"]
february = ["хлеб", "молоко", "сахар", "яйца", "яйца"]

a = set(january)
b = set(february)

print("январь:  позиций", len(january), "| разных", len(a))
print("февраль: позиций", len(february), "| разных", len(b))
print("пропало во втором:", sorted(a - b))
print("появилось во втором:", sorted(b - a))
print("осталось в обоих:", sorted(a & b))
print("всего разных за два месяца:", len(a | b))
