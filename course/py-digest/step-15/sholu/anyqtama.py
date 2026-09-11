"""The directory: what the bank does not send.

The bank answers with codes -- KAZ, UZB, RUS -- and a report full of codes is a
report somebody has to decode. The names are ours, so they live here, in a table
of their own with the code as the key.

Deliberately incomplete: a country the digest looks at may be missing from it,
and the report has to survive that rather than drop the country.
"""

import pandas as pd

NAMES = [
    ("KAZ", "Қазақстан", "Орталық Азия"),
    ("UZB", "Өзбекстан", "Орталық Азия"),
    ("KGZ", "Қырғызстан", "Орталық Азия"),
]


def table():
    """The directory as a table: code, name, region."""
    return pd.DataFrame(NAMES, columns=["country", "аты", "өңір"])
