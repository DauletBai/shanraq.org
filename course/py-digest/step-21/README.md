# Step 21 — the digest refuses what it cannot support

The state of the digest after lesson 42, *What comes from outside and what is
ours*.

The digest follows three countries. Asking whether prices and money moved
together across three countries always produces a number, and that number is
worthless: on the ten neighbours of lesson 42 a single country moved the same
kind of correlation by 1.3 and turned it from `+0.739` into `−0.538`.

So the fifth line of the report is a refusal:

```
елдер | керек | байланыс | қорытынды
    2 |     8 |        — | аз: 2 ел, 8 керек
```

- **`sholu/esep.py`, `baylanys()`** counts the countries that have both numbers,
  compares that with a floor it is given, and returns either the correlation or
  the reason it is not there. Two of the three countries qualify — the World
  Bank has no broad money for Russia — which is the second reason to say
  nothing.
- **The count travels with the answer.** `+0.97` over ten countries and `+0.97`
  over three are different statements, and pandas drops the pairs with a gap
  without saying so.
- **The refusal is printed, not skipped.** A line that disappears when the data
  is thin teaches the reader that the line was optional.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

The floor is the one number in `esep.py` chosen rather than measured: eight is
small enough for a regional digest to reach and large enough that no single
country decides the sign.
