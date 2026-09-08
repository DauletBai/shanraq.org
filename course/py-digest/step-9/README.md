# Step 9 — the digest gets names

The state of the digest after lesson 9, *Functions*.

The walk over a series, written inside the program until now, has become two
functions with names: `average()` leaves the years without a figure out of the
count, and `above()` compares two series. An `assert` stops an empty series
where the mistake is.

```
python3 main.py
```

Needs the network. Nothing is stored yet: every run asks the World Bank again.

The averages differ slightly from the ones printed in the lesson: the lesson
works on the series rounded to one decimal, this program on the figures as the
bank returns them.
