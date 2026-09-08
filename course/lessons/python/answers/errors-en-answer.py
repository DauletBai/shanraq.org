"""The answer to lesson 10's exercise: five rows from somebody else's export.

Only the risky line sits inside the try -- everything done after it succeeds is
in the else -- and the error carries its own name, so a skipped row says what
was wrong with it instead of disappearing.
"""


class BadRow(ValueError):
    """A row that could not be parsed."""


def to_number(text):
    """Text into a number: a comma counts as a point, anything else is a BadRow."""
    clean = text.strip().replace(",", ".")
    try:
        return float(clean)
    except ValueError:
        raise BadRow(f"not a number: {text!r}")


rows = ["520", "546,5", "", "no price", "498.0"]

total = 0.0
count = 0
for text in rows:
    try:
        value = to_number(text)
    except BadRow as error:
        print(f"skipped — {error}")
    else:
        total += value
        count += 1
        print(f"taken: {value}")
print(f"parsed {count} of {len(rows)}, average {total / count:.2f}")
