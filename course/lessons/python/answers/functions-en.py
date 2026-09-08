"""The answer to lesson 9's exercise: two functions over two series.

Both of them are given the gaps to deal with: average leaves them out of the
count, and above skips a year whose number is missing on either side -- a
comparison against None is a crash, not an answer.
"""


def average(series):
    """The average over a series; years without a number do not count."""
    total = 0.0
    count = 0
    for value in series.values():
        if value is None:
            continue
        total += value
        count += 1
    assert count > 0, "the series has no numbers at all"
    return total / count


def above(series, other):
    """The years in which the first series is above the second.

    A year without a number in either series does not count.
    """
    years = []
    for year, value in series.items():
        if value is None or other.get(year) is None:
            continue
        if value > other[year]:
            years.append(year)
    return years


kz = {2023: 14.5, 2024: 8.7, 2025: 11.4, 2026: None}
world = {2023: 5.8, 2024: 3.0, 2026: 2.9}

print(f"average for Kazakhstan: {average(kz):.2f}")
print(f"average for the world: {average(world):.2f}")
print("above the world:", above(kz, world))
