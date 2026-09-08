"""The answer to lesson 10's exercise: five rows from somebody else's export.

Only the risky line sits inside the try -- everything done after it succeeds is
in the else -- and the error carries its own name, so a skipped row says what
was wrong with it instead of disappearing.
"""


class BadRow(ValueError):
    """Строка, которую не удалось разобрать."""


def to_number(text):
    """Текст в число: запятая считается точкой, остальное — BadRow."""
    clean = text.strip().replace(",", ".")
    try:
        return float(clean)
    except ValueError:
        raise BadRow(f"не число: {text!r}")


rows = ["520", "546,5", "", "нет цены", "498.0"]

total = 0.0
count = 0
for text in rows:
    try:
        value = to_number(text)
    except BadRow as error:
        print(f"пропущено — {error}")
    else:
        total += value
        count += 1
        print(f"взято: {value}")
print(f"разобрано {count} из {len(rows)}, среднее {total / count:.2f}")
