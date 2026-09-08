"""The answer to lesson 9's exercise: two functions over two series.

Both of them are given the gaps to deal with: average leaves them out of the
count, and above skips a year whose number is missing on either side -- a
comparison against None is a crash, not an answer.
"""


def average(series):
    """Среднее по ряду; годы без числа не в счёт."""
    total = 0.0
    count = 0
    for value in series.values():
        if value is None:
            continue
        total += value
        count += 1
    assert count > 0, "в ряду нет ни одного числа"
    return total / count


def above(series, other):
    """Годы, в которых первый ряд выше второго.

    Год без числа хоть в одном из рядов не в счёт.
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

print(f"среднее по Казахстану: {average(kz):.2f}")
print(f"среднее по миру: {average(world):.2f}")
print("выше мировой:", above(kz, world))
