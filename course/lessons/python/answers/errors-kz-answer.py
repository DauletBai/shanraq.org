"""The answer to lesson 10's exercise: five rows from somebody else's export.

Only the risky line sits inside the try -- everything done after it succeeds is
in the else -- and the error carries its own name, so a skipped row says what
was wrong with it instead of disappearing.
"""


class BadRow(ValueError):
    """Талдауға келмеген жол."""


def to_number(text):
    """Мәтінді санға: үтір нүкте деп саналады, қалғаны — BadRow."""
    clean = text.strip().replace(",", ".")
    try:
        return float(clean)
    except ValueError:
        raise BadRow(f"сан емес: {text!r}")


rows = ["520", "546,5", "", "бағасы жоқ", "498.0"]

total = 0.0
count = 0
for text in rows:
    try:
        value = to_number(text)
    except BadRow as error:
        print(f"аттап өтілді — {error}")
    else:
        total += value
        count += 1
        print(f"алынды: {value}")
print(f"{len(rows)} жолдың {count}-і талданды, орташа {total / count:.2f}")
