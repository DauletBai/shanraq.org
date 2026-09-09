# Dates and time: 45 days, a timezone, and a format that reads two ways

_Лид (summary):_ **The fourteenth lesson of the Python course. There are 45 days between the 15th of January and the 1st of March, and `timedelta` counts them, not you. Plus two traps: `01.02.2026` parses under two formats without an error — as January and as February — and a naive time cannot be taken from an aware one, which Python says outright.**

## Why this is needed

In the previous lesson the year arrived as the string `"2025"` and we made a number of it, because it had to be compared afterwards. Dates are the same, only dearer.

A date in an export is text: `2026-01-15`, `15.01.2026`, `15 January`. While it is text, nothing can be done with it: you cannot subtract it, add a week to it, or ask how many days have passed.

And once it is parsed, it turns out that one and the same string reads two ways, and both of them look plausible. On top of that comes a question numbers never raised: half past nine — where?

## The whole thing at once

The file is `kunder.py`. Run it with `python kunder.py` from inside the environment.

The required part is the first two blocks: a string into a date, and the difference between dates. The third and the fourth show the format and the timezone the lesson was written for.

```python
"""Lesson 14: a date is not a string, it is a point on a time axis.

In the previous lesson the year arrived as the string "2025" and we made a
number of it. Dates are the same, only dearer: "01.02.2026" reads two ways, and
both of them look plausible.
"""

from datetime import date, datetime, timedelta, timezone

# Dates arrive from an export as text. The ISO format reads without explanation.
raw = ["2026-01-15", "2026-02-01", "2026-03-01"]

print("== from a string to a date")
days = []
for text in raw:
    day = date.fromisoformat(text)
    days.append(day)
    print(f"{text} → {day} | year {day.year}, month {day.month}, day {day.day}")

print()
print("== the difference between dates is a timedelta")
span = days[-1] - days[0]
print("between the first and the last:", span)
print("days:", span.days)
print("a week after the last one:", days[-1] + timedelta(days=7))

print()
print("== parse in one format, print in another")
moment = datetime.strptime("01.02.2026 09:30", "%d.%m.%Y %H:%M")
print("parsed:", moment)
print("printed:", moment.strftime("%d.%m.%Y, %H:%M"))
print("written into a file:", moment.isoformat())
try:
    datetime.strptime("2026-02-01", "%d.%m.%Y")
except ValueError as error:
    print("somebody else's format:", error)

print()
print("== the timezone: naive time and aware time")
naive = datetime(2026, 2, 1, 9, 30)
aware = datetime(2026, 2, 1, 9, 30, tzinfo=timezone.utc)
almaty = timezone(timedelta(hours=5))
print("naive:", naive, "| zone:", naive.tzinfo)
print("aware:", aware, "| zone:", aware.tzinfo)
print("the same instant in Almaty:", aware.astimezone(almaty))
try:
    print(aware - naive)
except TypeError as error:
    print("one cannot be taken from the other:", error)
```

It prints:

```
== from a string to a date
2026-01-15 → 2026-01-15 | year 2026, month 1, day 15
2026-02-01 → 2026-02-01 | year 2026, month 2, day 1
2026-03-01 → 2026-03-01 | year 2026, month 3, day 1

== the difference between dates is a timedelta
between the first and the last: 45 days, 0:00:00
days: 45
a week after the last one: 2026-03-08

== parse in one format, print in another
parsed: 2026-02-01 09:30:00
printed: 01.02.2026, 09:30
written into a file: 2026-02-01T09:30:00
somebody else's format: time data '2026-02-01' does not match format '%d.%m.%Y'

== the timezone: naive time and aware time
naive: 2026-02-01 09:30:00 | zone: None
aware: 2026-02-01 09:30:00+00:00 | zone: UTC
the same instant in Almaty: 2026-02-01 14:30:00+05:00
one cannot be taken from the other: can't subtract offset-naive and offset-aware datetimes
```

## Taking it apart

### A string becomes a date once — while it is parsed

```python
day = date.fromisoformat(text)
```

`fromisoformat` reads a date in the ISO 8601 format: year-month-day, four digits, two, two. It is the same move as `int(row["date"])` in the previous lesson: text becomes a value at once rather than "when it is needed".

It is needed at once. Everything a date is for is arithmetic: a difference, a shift, a comparison. Strings have none of it: `"2026-03-01" - "2026-01-15"` is an error.

Comparing the strings does work, mind: `"2026-03-01" > "2026-01-15"` is true. But it is true by accident: ISO is built so that the order of the strings matches the order of the dates. On `01.02.2026` the same comparison lies without a word.

> **Picture it.** A postmark. While the date is an impression on paper it is a picture; to count how long the letter travelled, somebody has to read it.

### The difference between dates is a span, not a number

```
between the first and the last: 45 days, 0:00:00
days: 45
```

Subtracting two dates gives not a number but a `timedelta` — a span. It prints itself its own way, and the number of days is asked of it: `.days`. Beside it live `.seconds` and `.total_seconds()`, for when the hours and minutes matter.

It works the other way too: `days[-1] + timedelta(days=7)` is the date a week later. What is added is a span rather than a number: dates cannot be added to numbers, and rightly so, because "seven" is seven of what.

### `strptime` reads, `strftime` prints

```python
moment = datetime.strptime("01.02.2026 09:30", "%d.%m.%Y %H:%M")
print(moment.strftime("%d.%m.%Y, %H:%M"))
```

The mnemonic is simple: **`p` for parse**, text into a date; **`f` for format**, a date into text. The format is a string of percent codes: `%d` the day, `%m` the month, `%Y` the four-digit year, `%H:%M` the hours and minutes. The full list is in [the documentation](https://docs.python.org/3/library/datetime.html#strftime-and-strptime-format-codes).

Both are needed, because a date has two addressees. Into a file and into an API it goes in ISO — `moment.isoformat()` — and anyone will read it. A person is shown the familiar `01.02.2026`.

There is a trap of politeness here as well: `%B` gives the name of the month — but `January` rather than «қаңтар», until the system's locale is set up. So a file keeps ISO, and the name of the month, where it is needed, is assembled by hand: the program has three languages and the locale has one.

### The format that reads two ways

```
somebody else's format: time data '2026-02-01' does not match format '%d.%m.%Y'
```

Here Python said it out loud: the format did not fit. But that is the **lucky** case.

The unlucky one is `01.02.2026`. In Kazakhstan and Russia it is the first of February; in the United States it is the second of January. Both dates exist, both parse without an error, and which one you get is decided not by the string but by the format you named. Name the wrong one and there will be no error and no warning: there will be a different date. Such a mistake is found a month later, when the wrong thing adds up.

Hence the rule: **ISO 8601 in files and APIs**, `2026-02-01`. It is unambiguous, it sorts as a string, and `fromisoformat` understands it. "Day.month.year" is for the screen.

### Naive time and aware time

```
naive: 2026-02-01 09:30:00 | zone: None
aware: 2026-02-01 09:30:00+00:00 | zone: UTC
```

A `datetime` has a `tzinfo` field, and it can be empty. Empty means **naive** time: half past nine, and where is unknown. The same `None` as in lesson seven, only what is missing now is the timezone.

Naive time does not answer the question "which instant is this". Half past nine in Almaty and half past nine in Warsaw are different moments, four hours apart. While the time is naive those four hours are not visible — they are simply not counted.

Aware time knows its zone, and it can be converted: `aware.astimezone(almaty)` gives the same instant written the Almaty way — 14:30. One instant, two writings.

Hence the attempt to subtract one from the other:

```
one cannot be taken from the other: can't subtract offset-naive and offset-aware datetimes
```

Python refuses — and that is the best thing it could do. A silent five-hour difference costs more than a loud error.

The rule is worth learning whole: **store in UTC, show in local time**. `datetime.now()` gives a naive local time, `datetime.now(timezone.utc)` an aware one; the second is what goes into a file and into a database.

A zone as an offset — `timezone(timedelta(hours=5))` — knows only "plus five". A zone by name — `ZoneInfo("Asia/Almaty")` from the `zoneinfo` module — also knows the history of the clock changes, and around the world that changes more often than one would think.

### What `timedelta` does not have

A month and a year. A span is measured in days, hours and seconds — in things that are always the same length. A "month" never is: it is 28 days one time and 31 another.

Adding a month means first agreeing what to do with the 31st of January. The standard library forces no such agreement on you; other libraries (`dateutil`) offer one along with their own rules.

## The map of the lesson

![The map of the lesson: the string, the span and the timezone](/static/course/py/map-dates-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why is a date turned into a `date` while it is parsed rather than left as a string?
2. How does `strptime` differ from `strftime`, and how is that remembered?
3. What is naive time, and why does Python refuse to subtract it from an aware one?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
from datetime import date

a = date(2026, 3, 1)
b = date(2026, 1, 15)
print((a - b).days, a > b)
```

**2. Fill in the gap.** In place of `...` put the format this string parses under.

```python
from datetime import datetime

moment = datetime.strptime("01.02.2026 09:30", ...)
print(moment.date())
```

**3. Fix it.** In an export from Kazakhstan, `01.02.2026` is the first of February. The program does not fall over, but it prints the wrong day. Find the cause.

```python
from datetime import datetime

day = datetime.strptime("01.02.2026", "%m.%d.%Y")
print(day.date())
```

## Exercise

**Required.** Given:

```python
rows = [
    ("2026-01-15", 520.0),
    ("2026-02-01", 546.5),
    ("2026-03-01", 498.0),
]
report = "02.03.2026 08:00"
```

Parse the dates of the measurements and print each one as `15.01.2026: 520.0`. Then print how many days there are between the first measurement and the last; the average interval between measurements to two decimal places; the date of the next measurement, thirty days after the last one. At the end parse the moment of the report under the format `day.month.year hours:minutes`, declare it Almaty time (UTC+5) and print it twice: in ISO in Almaty time, and in UTC.

The expected output:

<!-- task out -->
```
15.01.2026: 520.0
01.02.2026: 546.5
01.03.2026: 498.0
days between the first and the last: 45
average interval: 22.50 days
the next measurement: 31.03.2026
the report in Almaty: 2026-03-02T08:00:00+05:00
the same in UTC:      2026-03-02T03:00:00+00:00
```

Done when: the output matches line by line; the dates become a `date` or a `datetime` while being parsed rather than staying strings; the average interval is divided by the number of **intervals** rather than of measurements; the moment of the report is given its zone where it is parsed — otherwise it arrives in UTC with the same figures on it.

**On your own data.** Take three dates of your own — receipts, payments, anything with a figure on it. Work out how many days there are between the first and the last, and print each date in the form you are used to. Write one of them in a foreign format (`15/01/2026` or `15 Jan 2026`) and parse it with a `strptime` of your own.

**Optional.**

- Parse `01.02.2026` under both formats — `%d.%m.%Y` and `%m.%d.%Y` — and print the two dates side by side.
- Take `ZoneInfo("Asia/Almaty")` instead of a fixed offset and compare the result.
- Try adding `timedelta(months=1)` to a date and read what Python says.

## Where this goes in the project

The digest gets a time axis. The measurements stop being merely numbers: days are counted between them, and "when it was collected" is written in UTC — not because it looks better, but because the digest will one day move to a server, and that server's clock is not set to Almaty.

Debts. We still have nothing to add a month with, and a fixed offset knows nothing about clock changes. The first is settled by an agreement, the second by `zoneinfo`; we will come back to both when the digest starts running on a schedule.

## The answers

### To the questions

1. Because arithmetic comes next: a difference, a shift by a week, a comparison. A string can do none of it, and comparing strings matches comparing dates only in ISO and only by accident.
2. `p` for parse: `strptime` reads text and makes a `datetime`. `f` for format: `strftime` takes a `datetime` and makes text. The first is for the input, the second for the output.
3. Naive is a `datetime` whose `tzinfo` is empty: a time without a place. It cannot be taken from an aware one because how many hours lie between them is unknown; Python refuses out loud instead of working out a difference that is wrong by a whole zone.

### To the warm-up

1. `45 True`. The difference between dates is a `timedelta`, and `.days` is asked of it; dates can be compared directly, and the 1st of March is later than the 15th of January.

<!-- drill 1 out -->
```
45 True
```

2. `"%d.%m.%Y %H:%M"`. The format describes the **incoming** string rather than what you would like to get: day, dot, month, dot, year, space, hours, colon, minutes.

<!-- drill 2 -->
```python
from datetime import datetime

moment = datetime.strptime("01.02.2026 09:30", "%d.%m.%Y %H:%M")
print(moment.date())
```

<!-- drill 2 out -->
```
2026-02-01
```

3. The format was named the American way: `%m.%d.%Y` puts the month first. The string fitted it, there was no error, and what came out was the second of January instead of the first of February. That is what the most expensive mistake in dates looks like: not a crash, a different date.

<!-- drill 3 -->
```python
from datetime import datetime

day = datetime.strptime("01.02.2026", "%d.%m.%Y")
print(day.date())
```

<!-- drill 3 out -->
```
2026-02-01
```

## Sources

- [Python: the datetime module](https://docs.python.org/3/library/datetime.html)
- [Python: strftime and strptime format codes](https://docs.python.org/3/library/datetime.html#strftime-and-strptime-format-codes)
- [Python: zoneinfo — timezones by name](https://docs.python.org/3/library/zoneinfo.html)
- [RFC 3339: dates and times on the internet](https://www.rfc-editor.org/rfc/rfc3339)
