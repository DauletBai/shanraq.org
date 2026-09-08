"""The answer to lesson 7's exercise, English: monthly prices against the average.

The order of the checks is set by the task -- the gap first, then "about the
average", then above and below -- so that no comparison ever meets a None. The
average counts only the months that have a number.
"""

prices = {"January": 520.0, "February": 580.0, "March": None, "April": 498.0}

total = 0.0
count = 0
for value in prices.values():
    if value is None:
        continue
    total += value
    count += 1
average = total / count
print(f"average: {average:.2f}")

for month, value in prices.items():
    if value is None:
        print(f"{month}: no data")
    elif abs(value - average) <= average * 0.05:
        print(f"{month}: {value} — about the average")
    elif value > average:
        print(f"{month}: {value} — above the average")
    else:
        print(f"{month}: {value} — below the average")
