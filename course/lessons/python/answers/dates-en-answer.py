"""The answer to lesson 14's exercise, English: three measurements on a time axis.

The dates arrive as text and become dates at once, because everything after
that -- the span, the interval, the next date -- is arithmetic that text cannot
do. The moment of the report is given a timezone at the point it is parsed: a
naive time is not a moment, it is a reading off somebody's wall.
"""

from datetime import datetime, timedelta, timezone

rows = [
    ("2026-01-15", 520.0),
    ("2026-02-01", 546.5),
    ("2026-03-01", 498.0),
]
report = "02.03.2026 08:00"

days = []
for text, price in rows:
    day = datetime.strptime(text, "%Y-%m-%d").date()
    days.append(day)
    print(f"{day.strftime('%d.%m.%Y')}: {price}")

span = days[-1] - days[0]
print("days between the first and the last:", span.days)
print(f"average interval: {span.days / (len(days) - 1):.2f} days")
print("the next measurement:", (days[-1] + timedelta(days=30)).strftime("%d.%m.%Y"))

almaty = timezone(timedelta(hours=5))
moment = datetime.strptime(report, "%d.%m.%Y %H:%M").replace(tzinfo=almaty)
print("the report in Almaty:", moment.isoformat())
print("the same in UTC:     ", moment.astimezone(timezone.utc).isoformat())
