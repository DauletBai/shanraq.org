# Lesson 12. Keep more than one possible meaning

## Why this matters

A door marked “hall” might lead to a sports hall or a concert hall; you need more information. Kazakh `ат` likewise means “horse” and “name.” One isolated word cannot justify choosing either. A [dictionary of literary Kazakh](https://sozdikqor.kz/sozdik/?id=38) records both meanings.

## See the whole solution first

For one spelling we keep a **list of candidates**: possible analyses, not one guessed answer. For `ат`, store `жылқы` (“horse”) and `есім` (“name”). For `үй`, store only `баспана` (“home”) so far. These are not numerical probabilities: the first list item is not automatically more likely.

Create `meanings.py` in UTF-8 and run it with the familiar command.

```python
analyses = {"ат": ["жылқы", "есім"],
            "үй": ["баспана"]}
word = input("Сөз: ").lower()
if word in analyses:
    print(analyses[word])
else:
    print("білмеймін")
```

The [step file](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-12) matches the lesson code.

We already know dictionaries `{...}` from lesson 8 and lists `[...]` from lesson 7. Now one dictionary value is itself a whole list. Imagine a word card with two meanings on its reverse side. `analyses[word]` gets the list for a key, and `print` displays it with square brackets and quotation marks. Input `ат` gives `['жылқы', 'есім']`; `үй` gives `['баспана']`. Even one value stays in a list, so later we can ask the same question of every word: “how many candidates?” Input `жол` gives `білмеймін` because there is no such key.

This is **not** an analysis of every possible ending or a context-based choice of meaning. We entered two words by hand so ambiguity is preserved. If a later question needs to choose between “horse” and “name,” the model must find an explicit reason or ask for clarification. It must not treat the first entry as true merely because it is first.

## Memory map

Spelling → dictionary → list of possible meanings → keep all; no entry → `білмеймін`. The arrow does not jump straight to one answer.

## Recall and task

Hide the code and rebuild the entry for `ат` using the map. Explain why `['баспана']` stays a list even with one value. Then try `АТ`, `үй`, and `жол`.

**Task:** add `"кітап"` with one value, `"оқу құралы"` (“reading material”), and test `кітап`. **Hint:** put a comma after the `үй` pair and enclose the new value in square brackets. **Answer:** add `"кітап": ["оқу құралы"]` inside the dictionary; the output is `['оқу құралы']`. **Common mistake:** writing the two `ат` meanings under two identical dictionary keys: the second overwrites the first. One key must point to a list.
