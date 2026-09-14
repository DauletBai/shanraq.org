"""Cleaning: what arrives is not what can be counted.

The bank is a tidy source and still sends a year with no figure in it, because
the year is not over. A file read back off the disk can carry a row twice if a
run was interrupted halfway. Neither is a disaster; both have to be decided
about once, here, rather than guessed at in every report.

Every decision leaves a line in the log. A table nobody can say what was done to
is worse than a dirty one.

The year also becomes a date here. A year as a number sorts and subtracts by
luck: 2024 and 2025 happen to be numbers that behave. A date says what it is,
and the day the digest reads a monthly series the code above will not have to
change.
"""

import pandas as pd


def clean(table):
    """Returns the table fit to count, and a log of what was done to it."""
    log = [f"жол келді: {len(table)}"]

    # A row without a key cannot be counted and cannot be joined: it is not a
    # gap in the data, it is the absence of the row itself.
    keyed = table.dropna(subset=["country", "year"])
    if len(keyed) != len(table):
        log.append(f"кілтсіз жол алынып тасталды: {len(table) - len(keyed)}")

    before = len(keyed)
    unique = keyed.drop_duplicates(subset=["country", "year"], keep="last")
    if len(unique) != before:
        log.append(f"қайталанған жол алынып тасталды: {before - len(unique)}")

    clean = unique.copy()
    clean["year"] = clean["year"].astype("int64")
    # Anything that is not a number becomes a gap rather than a zero: a year
    # not yet counted is not a year of zero inflation.
    clean["value"] = pd.to_numeric(clean["value"], errors="coerce")
    log.append(f"саны жоқ жыл: {int(clean['value'].isna().sum())}")

    # The first day of the year: a date the rest of pandas understands.
    clean["date"] = pd.to_datetime(clean["year"], format="%Y")

    # Five words in a column of hundreds of rows: stored once as a category.
    clean["country"] = clean["country"].astype("category")

    log.append(f"есепке кететін жол: {len(clean)}")
    return clean.sort_values(["country", "date"]).reset_index(drop=True), log
