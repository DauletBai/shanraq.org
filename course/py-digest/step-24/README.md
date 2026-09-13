# Step 24 — the report says what it does not say

The state of the digest after lesson 46, *Where a model ends*, and the last step
of the module.

Everything above this step counts. This one writes down what the counting does
not cover, at the foot of the page:

```
Бұл есеп нені айтпайды
 • болжам жоқ: тренд өз жылдарында 3.55 тармаққа дейін қателеседі
 • баға мен ақшаның байланысы саналмады: 2 ел, 8 керек
 • ақша массасы жоқ ел: RUS
 • жылдар: 2025 дейін; одан кейінгі жыл туралы есеп ештеңе айтпайды
```

Not one of those lines is typed by hand. `esep.shekteu()` takes each figure
from the table that produced it: the worst miss comes out of the trend, the two
countries out of the refused correlation, the blank money column out of the
money table, the last year out of the report itself. A table that changes
changes the line that limits it — which is the whole reason to build the limits
from data instead of writing them once in a footer and forgetting them.

```
pip install -r requirements.txt
python3 main.py --dry-run     # count, write nothing
python3 main.py               # rebuild the report
```

A page of numbers with nothing beside them reads as a page of facts. This is the
cheapest thing in the whole digest — one function, four lines of output — and it
is the difference between a report and a claim.
