"""The answer to lesson 8's exercise: a walk over the months with a gap in it.

The gap is skipped rather than counted, the numbering starts at one because
that is how months are named outside a program, and the search for the first
expensive month stops as soon as it finds one.
"""

prices = [520.0, 546.0, None, 498.0, 515.0]

total = 0.0
count = 0
for number, value in enumerate(prices, start=1):
    if value is None:
        print(f"{number}: дерек жоқ")
        continue
    total += value
    count += 1
    print(f"{number}: {value}")
print(f"саны бар ай: {count}, орташа: {total / count:.2f}")

for number, value in enumerate(prices, start=1):
    if value is not None and value > 530:
        print(f"530-дан қымбат алғашқы ай: {number}")
        break
