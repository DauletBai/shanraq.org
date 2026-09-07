# Counting for yourself: prices, the tenge and one question

_Лид (summary):_ **The first lesson of the Python course. Nothing has to be installed: a program of thirty lines fetches the official numbers itself. Prices in Kazakhstan grew 3.48 times since 2010 — a thousand tenge of that year is worth 287 today. And after the same shock the neighbours' inflation is three times lower.**

## Why this is needed

Numbers are quoted at you. Inflation is this. Growth is that. The one to blame is him. You cannot check what was said, so all that is left is to believe it or not.

This course is about no longer having to choose between believing and not believing. The data on prices, money and the exchange rate lies in the open; to take it and count you need one tool — Python — and a few evenings.

Let us start from the end: with the answer you will have in fifteen minutes, knowing nothing about the language yet.

## The whole thing at once

Make a file called `tsena.py` and paste this in. Nothing has to be installed: everything used here comes with Python.

```python
"""The first program of the course: it asks the question the course was written for.

Nothing has to be installed -- everything used here comes with Python. The
program fetches the official numbers itself and counts with them.
"""

import json
import urllib.request

# The World Bank hands out country indicators with no key and no sign-up.
# FP.CPI.TOTL is the consumer price index, FP.CPI.TOTL.ZG inflation in percent.
API = "https://api.worldbank.org/v2/country/{}/indicator/{}?format=json&per_page=80&date={}:{}"


def rows(country, indicator, first, last):
    """Returns {year: value} -- whatever the World Bank answered."""
    url = API.format(country, indicator, first, last)
    with urllib.request.urlopen(url, timeout=30) as answer:
        data = json.load(answer)
    if len(data) < 2 or not data[1]:
        return {}
    return {int(item["date"]): item["value"] for item in data[1] if item["value"] is not None}


def main():
    print("== how many times prices grew in Kazakhstan")
    prices = rows("KZ", "FP.CPI.TOTL", 2010, 2025)
    base, last = min(prices), max(prices)
    times = prices[last] / prices[base]
    print(f"price index in {base}: {prices[base]:.1f}")
    print(f"price index in {last}: {prices[last]:.1f}")
    print(f"prices grew {times:.2f} times")
    print(f"1000 tenge of {base} is worth {1000 / times:.0f} tenge today")

    print()
    print("== one world, different prices: inflation by year, %")
    countries = {"KZ": "Kazakhstan", "WLD": "the world", "GE": "Georgia",
                 "AM": "Armenia", "PL": "Poland"}
    years = range(2021, 2026)
    print("country     " + "".join(f"{year:>7}" for year in years))
    for code, name in countries.items():
        inflation = rows(code, "FP.CPI.TOTL.ZG", 2021, 2025)
        line = "".join(f"{inflation[year]:7.1f}" if year in inflation else f"{'—':>7}"
                       for year in years)
        print(f"{name:<12}{line}")


main()
```

Run `python3 tsena.py`. The output:

```
== how many times prices grew in Kazakhstan
price index in 2010: 100.0
price index in 2025: 348.1
prices grew 3.48 times
1000 tenge of 2010 is worth 287 tenge today

== one world, different prices: inflation by year, %
country        2021   2022   2023   2024   2025
Kazakhstan      8.0   15.0   14.5    8.7   11.4
the world       3.5    8.1    5.8    3.0    3.0
Georgia         9.6   11.9    2.5    1.1    3.9
Armenia         7.2    8.6    2.0    0.3    3.3
Poland          5.1   14.4   11.5    3.8    3.8
```

## The walk-through

### What you have just counted

The first four lines are about the money in your pocket. The price index of 2010 is 100; today it is 348.1. That means prices grew **3.48 times**, and the thousand tenge you held in 2010 buys today what **287 tenge** buys.

Not "about three times" and not "everybody knows". The number was counted in front of you from an official series, and you can change the year from 2010 to 2015 and see what happens.

### The same shock, different numbers

The second table matters more than the first, and here is why.

When prices rise, the explanation usually sounds like this: the world got dearer, the war, logistics, grain, oil. Let us check. In 2022 things did indeed get dearer everywhere: the world 8.1%, Kazakhstan 15.0%.

After that the paths part. By 2025 the world is back at 3.0%, Armenia at 3.3%, Georgia at 3.9%, Poland at 3.8%. Kazakhstan is at **11.4%** and rising again.

The shock was shared. So it is not what makes the difference.

### What follows from this, and what does not

One thing follows: **"the outside world is to blame" is not a sufficient explanation**. It does not survive a comparison with neighbours who got the same oil, the same grain and the same logistics.

But the opposite simple explanation — "they printed money" — does not add up on its own either, and we will check that too. Broad money grows by 13–18% a year in Armenia, 11–17% in Georgia and 12–21% in Kazakhstan. The growth is comparable while inflation differs three- to fourfold. So it is not one row of numbers: beside it stand the tenge's exchange rate, the tariffs, the taxes, and **where** the new money goes.

The course will not hand you a ready culprit. It will hand you the ability to put four series side by side and see which one does not add up. That is sturdier than anybody else's conclusion — ours included.

> **Picture it.** A doctor who names the diagnosis over the phone, and a doctor who gives you the scan and teaches you to read it. The first may be right. The second cannot be fooled.

### How we will learn

The Go course on this site is built on the method of **Viktor Fyodorovich Shatalov**, a Soviet teacher whose system let schoolchildren cover the syllabus several times faster. We consider it the best thing twentieth-century teaching produced, and the second course follows it too.

Four things are taken from it.

**The whole first, the details after.** Every lesson begins with a working program. You run it understanding nothing yet — as today — and only then take it apart.

**The supporting signal.** At the end of the walk-through comes the map of the lesson: one picture on one screen, where what matters is drawn as shapes and links. It can be photographed and kept on a phone; the point is that you can redraw it by hand.

**A picture instead of a definition.** Every unfamiliar term is explained with something from ordinary life. A price index is not "a basket deflator" but a ruler that had a hundred divisions in 2010 and has three hundred and forty-eight today.

**The right not to understand the first time.** No marks, no failed tests, no "you did not pass". One required exercise per lesson — small and always doable — and two more if you want them.

### What we will build

By the end of the course you have a **digest of your own**: a program that fetches fresh data itself, files it in a database, counts, draws charts, assembles a report page and runs on a schedule without you. Fifty-five lessons, from the first line to the schedule.

Along the way you will count your own personal inflation from your own receipts and compare it with the official figure; find out how much money there is per tenge of GDP and who created it; train a first model and see where it lies. In the last lessons a language model appears beside the numbers — and with it the checking of what it wrote.

One rule of the course stands apart: **you are invited to check us as well**. This site has already published pieces on money, banks and inflation. In the exercises you will take a claim from one of our articles and check it against open data. If it does not hold, write to us and the article gets corrected.

### What you need to start

A computer, the internet and Python. Check whether you have it: type `python3 --version` in a terminal. If it answers with a version number, you are ready for this lesson. If not, that is the next lesson, where the language and the workplace are set up from scratch.

Nothing else: no paid course, no keys, no sign-up. The data we count is open to everybody — that is the whole point.

## The map of the lesson

![The map of the lesson: prices, the neighbours and the question](/static/course/py/map-prices-en.svg)

## Say it in your own words

Without looking, answer aloud or on paper. The answers are at the end of the lesson.

1. What does "the price index of 2025 is 348.1" mean, when in 2010 it was 100?
2. Why is the table of neighbouring countries stronger than a single row for Kazakhstan?
3. Why can these two tables not yet name the one to blame?

## The exercise

**Required.** Run the program and change the starting year in it from 2010 to the year you were born — or any year that matters to you. Work out what a thousand tenge of that year is worth today.

**If you want more.**

- Add a second row of countries to the table: Kyrgyzstan (`KG`), Uzbekistan (`UZ`), Turkey (`TR`). See who is near us and who is not.
- Replace the indicator `FP.CPI.TOTL.ZG` with `PA.NUS.FCRF`, the average exchange rate to the dollar. What happened to the tenge over the same years?
- Find somebody's public claim about prices in Kazakhstan and check it with this program.

## Where this goes in the project

Today's program is the first version of our digest: it already fetches data from the world and counts something with it. Next we teach it to keep a history rather than ask for everything again, add the tenge's rate from the National Bank, put the data in a database, draw charts and assemble a report.

The debts are visible already. The program dies without the internet and says nothing human about it. It fetches the data afresh every time, though a yearly series changes once a year. And it has not one check in it: if the World Bank answers with emptiness, we will not notice. All of that is the next few lessons.

## The answers

1. That prices grew 3.48 times: the same basket that cost 100 units in 2010 costs 348.1 today. The other side of it is that a thousand tenge of that year buys what 287 buys now.
2. Because one row can be explained by anything. A comparison tests the explanation: if a shared external shock is to blame, the neighbours who got the same shock should show similar numbers. They do not.
3. Because a coincidence and a difference are not yet a cause. We have seen that an external shock cannot explain it; to name a cause you have to put money, the exchange rate, tariffs and taxes side by side — and that is the work of several lessons, not one table.

## Sources

- [World Bank: the consumer price index](https://data.worldbank.org/indicator/FP.CPI.TOTL)
- [World Bank: the indicators API](https://datahelpdesk.worldbank.org/knowledgebase/articles/889392)
- [Python: the urllib.request module](https://docs.python.org/3/library/urllib.request.html)
