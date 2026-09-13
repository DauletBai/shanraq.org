"""Formats: how a number reaches a reader.

Numbers are computed once and formatted once, and those are different places.
Everything above this file works with full precision; everything the reader sees
goes through here.

The separators are ours rather than Python's default: a space in the thousands
and a comma in the decimals, in the numbers and in the percentages alike -- a
table in which the two disagree looks like two tables.
"""

DASH = "—"


def number(value, digits=2):
    """1234567.891 -> "1 234 567,89". A gap becomes a dash."""
    if value is None or value != value:      # NaN is the only value unequal to itself
        return DASH
    return f"{value:,.{digits}f}".replace(",", "\u00a0").replace(".", ",")


def signed(value, digits=2):
    """A change: the sign is part of the number, and zero is not "+0"."""
    if value is None or value != value:
        return DASH
    if round(value, digits) == 0:
        return number(0, digits)
    return f"{value:+,.{digits}f}".replace(",", "\u00a0").replace(".", ",")


def percent(value, digits=1):
    """A share as 0.674 -> "67,4%", with the same comma as the numbers."""
    if value is None or value != value:
        return DASH
    return f"{value:.{digits}%}".replace(".", ",")
