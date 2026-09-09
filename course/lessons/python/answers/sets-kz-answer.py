"""The answer to lesson 15's exercise, Kazakh: two exports, and what changed between them.

Every set is printed through sorted: a set has no order of its own, and printing
one directly gives an answer that can come out differently on another machine.
The duplicates disappear on the way in, which is the point -- the question is
what was on the list, not how many times it was written down.
"""

january = ["нан", "сүт", "май", "тұз", "нан", "жұмыртқа"]
february = ["нан", "сүт", "қант", "жұмыртқа", "жұмыртқа"]

a = set(january)
b = set(february)

print("қаңтар: позиция", len(january), "| әртүрлі", len(a))
print("ақпан:  позиция", len(february), "| әртүрлі", len(b))
print("екіншісінде жоқ:", sorted(a - b))
print("екіншісінде пайда болды:", sorted(b - a))
print("екеуінде де қалды:", sorted(a & b))
print("екі айда барлығы әртүрлі:", len(a | b))
