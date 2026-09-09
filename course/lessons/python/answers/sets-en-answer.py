"""The answer to lesson 15's exercise, English: two exports, and what changed between them.

Every set is printed through sorted: a set has no order of its own, and printing
one directly gives an answer that can come out differently on another machine.
The duplicates disappear on the way in, which is the point -- the question is
what was on the list, not how many times it was written down.
"""

january = ["bread", "milk", "butter", "salt", "bread", "eggs"]
february = ["bread", "milk", "sugar", "eggs", "eggs"]

a = set(january)
b = set(february)

print("January: items", len(january), "| different", len(a))
print("February: items", len(february), "| different", len(b))
print("gone from the second:", sorted(a - b))
print("new in the second:", sorted(b - a))
print("in both:", sorted(a & b))
print("different names over two months:", len(a | b))
