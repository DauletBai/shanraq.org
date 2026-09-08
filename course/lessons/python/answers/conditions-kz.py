"""The answer to lesson 7's exercise, Kazakh: monthly prices against the average.

The order of the checks is set by the task -- the gap first, then "close to the
average", then above and below -- so that no comparison ever meets a None. The
average counts only the months that have a number.
"""

prices = {"қаңтар": 520.0, "ақпан": 580.0, "наурыз": None, "сәуір": 498.0}

total = 0.0
count = 0
for value in prices.values():
    if value is None:
        continue
    total += value
    count += 1
average = total / count
print(f"орташа: {average:.2f}")

for month, value in prices.items():
    if value is None:
        print(f"{month}: дерек жоқ")
    elif abs(value - average) <= average * 0.05:
        print(f"{month}: {value} — орташаға жуық")
    elif value > average:
        print(f"{month}: {value} — орташадан жоғары")
    else:
        print(f"{month}: {value} — орташадан төмен")
