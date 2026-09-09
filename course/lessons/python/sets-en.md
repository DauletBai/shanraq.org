# Sets: `&`, `-`, `|`, and what the second export does not have

_Лид (summary):_ **The fifteenth lesson of the Python course. Of the five names in the list only four are different: a set swallows the repeat without a word. In exchange, three questions about two exports take three signs: `&` for what is in both, `-` for what is gone, `|` for all of it. Plus `{}` is an empty dictionary, not an empty set.**

## Why this is needed

Lists and dictionaries answer "how many" and "what is under this key". There is a third question, and in work with data it is asked more often than either: **is it there**.

Two exports have arrived, one for January and one for February. What was in both? What is gone from the second? What is new? With lists that takes nested loops, and nested loops are easy to get wrong.

A set is a collection without repeats and without order. Exactly what "is it there" needs — and three questions about two exports take it three signs.

## The whole thing at once

The file is `jiyn.py`. Run it with `python jiyn.py` from inside the environment.

The required part is the first two blocks: turning a list into a set, and the three operations. The third and the fourth show the membership test and two traps.

```python
"""Lesson 15: a set is about "is it there", not about "how many and in what order".

Two exports of goods from two months. The questions are the same for both: what
is in each of them, what is gone from the second, what is new. With lists that
takes loops; with sets it takes one sign.
"""

january = ["bread", "milk", "butter", "salt", "bread"]
february = ["bread", "milk", "sugar", "eggs"]

print("== what is lost on the way into a set")
a = set(january)
b = set(february)
print("in January's list:", len(january), "| different names:", len(a))
print("a set is printed sorted:", sorted(a))

print()
print("== three questions, three signs")
print("in both (a & b):       ", sorted(a & b))
print("gone from the second (a - b):", sorted(a - b))
print("new (b - a):           ", sorted(b - a))
print("all of it (a | b):     ", sorted(a | b))
print("in one only (a ^ b):   ", sorted(a ^ b))

print()
print("== is it there, and does it fit inside")
print("'salt' in January:", "salt" in a, "| in February:", "salt" in b)
print("February inside January:", b <= a, "| they overlap:", not a.isdisjoint(b))
print("different names over two months:", len(a | b))

print()
print("== two traps")
print("the type of {} is", type({}).__name__, "| the type of set() is", type(set()).__name__)
try:
    {["bread", "milk"]}
except TypeError as error:
    print("a list cannot go inside:", error)
```

It prints:

```
== what is lost on the way into a set
in January's list: 5 | different names: 4
a set is printed sorted: ['bread', 'butter', 'milk', 'salt']

== three questions, three signs
in both (a & b):        ['bread', 'milk']
gone from the second (a - b): ['butter', 'salt']
new (b - a):            ['eggs', 'sugar']
all of it (a | b):      ['bread', 'butter', 'eggs', 'milk', 'salt', 'sugar']
in one only (a ^ b):    ['butter', 'eggs', 'salt', 'sugar']

== is it there, and does it fit inside
'salt' in January: True | in February: False
February inside January: False | they overlap: True
different names over two months: 6

== two traps
the type of {} is dict | the type of set() is set
a list cannot go inside: cannot use 'list' as a set element (unhashable type: 'list')
```

## Taking it apart

### The repeats vanish without a word

```
in January's list: 5 | different names: 4
```

`set(january)` throws the second "bread" away and says nothing about it. That is what a set is: **an element is either in it or not**, and how many times it turned up is not a question for it.

Hence the most honest common use: counting how many **different** values the data holds. `len(set(...))` against `len(...)` is one line, and it shows at once whether the export has duplicates in it.

Hence the mistake as well: a set must not be reached for where the repeats are the data. Three purchases of bread in a month are three purchases and not one; those need a list or a counter, not a set.

> **Picture it.** A guest list. What matters is who was invited, not how many times their name was written down. But if it is a shopping list, then two loaves are two loaves.

### There is no order, so it is printed through `sorted`

```python
print("a set is printed sorted:", sorted(a))
```

A set has no order of its own — not the order things were added in, and not any other. Print one directly and the order may come out differently on another machine or on another run: it depends on how the elements were laid out inside.

So the rule is simple: **a set is shown through `sorted`**. Then the output is the same for everybody and can be compared by eye. Inside the program the order is not needed — there a set is asked, not read.

### Three questions, three signs

```
in both (a & b):        ['bread', 'milk']
gone from the second (a - b): ['butter', 'salt']
new (b - a):            ['eggs', 'sugar']
```

`&` is the intersection: what is in both. `-` is the difference: what is in the first and not in the second. `|` is the union: all of it, once each. `^` is the symmetric difference: what is in exactly one of the two.

The order matters with `-`: `a - b` and `b - a` answer different questions — "what is gone" and "what is new". That is the commonest confusion at first, and it is cured by asking out loud which set is being subtracted from.

Every sign has a word for a twin: `a.intersection(b)`, `a.difference(b)`, `a.union(b)`, `a.symmetric_difference(b)`. The signs are shorter, the words are clearer in somebody else's code — take whichever reads better in place.

### "Is it there" is what sets were made for

```python
"salt" in a
```

An `in` test works on a list and on a set alike, but not in the same way. In a list Python walks the elements until it finds one: the longer the list, the longer it takes. In a set it works out at once where the element would have to be and looks only there — the size does not matter.

Over ten rows the difference is invisible. Over an export of a hundred thousand rows checked in a loop it is the difference between a second and half an hour — and that is the one reason a set is sometimes built purely for the checking.

Beside it stand the questions about whole collections: `b <= a` — is every element of `b` in `a`; `a.isdisjoint(b)` — do they not overlap at all.

### Two traps

```
the type of {} is dict | the type of set() is set
```

The braces are taken by the dictionary: `{}` is an **empty dictionary**, not an empty set. The empty set is only `set()`. The mistake does not show at once: `seen = {}` works until the first `seen.add(...)`, and there it says a dictionary has no `add`.

```
a list cannot go inside: cannot use 'list' as a set element (unhashable type: 'list')
```

Only what cannot be changed may go into a set: numbers, strings, tuples. A list can be changed — which means that, having put one inside, you could alter an element already laid out in its place, and the set would stop finding its own contents.

The reason is the one that keeps a list from being a dictionary key: in both cases the element is laid out by its hash. If an unchangeable collection is what you want, there is `frozenset`, and that one can go inside another set.

## The map of the lesson

![The map of the lesson: the repeats, the three signs and "is it there"](/static/course/py/map-sets-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What does a set do with repeats, and when is that harmful?
2. How does `a - b` differ from `b - a` on two exports?
3. Why is a set printed through `sorted`?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this line print?

<!-- drill 1 -->
```python
print(len([1, 2, 2, 3]), len({1, 2, 2, 3}))
```

**2. Fill in the gap.** In place of `...` put what leaves the names the second export does not have.

```python
a = {"bread", "milk", "salt"}
b = {"bread", "milk", "sugar"}
print(sorted(...))
```

**3. Fix it.** The program falls over on its second line. Read the error and make what is actually wanted.

```python
seen = {}
seen.add("bread")
print(sorted(seen))
```

## Exercise

**Required.** Given:

```python
january = ["bread", "milk", "butter", "salt", "bread", "eggs"]
february = ["bread", "milk", "sugar", "eggs", "eggs"]
```

For each month print how many items there were and how many of them were different. Then three answers: what is gone from the second month, what is new in it, what stayed in both. At the end, how many different names there were over the two months. Print every collection through `sorted`.

The expected output:

<!-- task out -->
```
January: items 6 | different 5
February: items 5 | different 4
gone from the second: ['butter', 'salt']
new in the second: ['sugar']
in both: ['bread', 'eggs', 'milk']
different names over two months: 6
```

Done when: the output matches line by line; there is not one nested loop — all three answers come from the signs; every set is printed through `sorted` rather than directly.

**On your own data.** Take two lists of your own — the goods on two receipts, the members of two chats, the files in two folders. Answer the same three questions, and say out loud which of them you would need more often in life.

**Optional.**

- Print a set without `sorted` several times in a row and look at the order.
- Build a `frozenset` out of the first month and put it inside another set.
- Count how many times each name occurs — and explain why a set is no good for that.

## Where this goes in the project

The digest starts noticing what went missing. It used to count what arrived; now it can say what the new export does **not** have compared with the last one — and that is the first sign that the source broke rather than that the world changed.

Still open. A set answers "is it there" and not "how many times" — that needs a counter, which arrives in the lesson on the `collections` module. And it keeps no order, so "the first five new ones" cannot be taken out of it without sorting.

## The answers

### To the questions

1. It throws them away without a word: an element is either in it or not. That is harmful where the repeat is the data: three purchases of bread turn into one, and the receipt no longer adds up.
2. `a - b` is what is gone from the second export, `b - a` is what appeared in it. Different questions, and the order of the subtraction decides which one you are answering.
3. Because a set has no order of its own: printed directly it may come out differently on another machine or on another run. `sorted` makes the output the same for everybody.

### To the warm-up

1. `4 3`. The list holds four elements and the set three: the second two is not in it, and nobody said so.

<!-- drill 1 out -->
```
4 3
```

2. `a - b`. The difference in that order leaves what is in the first collection and not in the second.

<!-- drill 2 -->
```python
a = {"bread", "milk", "salt"}
b = {"bread", "milk", "sugar"}
print(sorted(a - b))
```

<!-- drill 2 out -->
```
['salt']
```

3. `{}` is an empty dictionary rather than an empty set, so it had no `add`: `AttributeError: 'dict' object has no attribute 'add'`. An empty set is made with `set()`.

<!-- drill 3 -->
```python
seen = set()
seen.add("bread")
print(sorted(seen))
```

<!-- drill 3 out -->
```
['bread']
```

## Sources

- [Python: set types](https://docs.python.org/3/library/stdtypes.html#set-types-set-frozenset)
- [Python: the tutorial on sets](https://docs.python.org/3/tutorial/datastructures.html#sets)
- [Python: why elements have to be hashable](https://docs.python.org/3/glossary.html#term-hashable)
- [Python: the sorted function](https://docs.python.org/3/library/functions.html#sorted)
