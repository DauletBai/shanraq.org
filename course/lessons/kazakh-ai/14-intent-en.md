# Lesson 14. Identify what the person asks for

## Why this matters

`Шахмат қашан?` (“When is chess?”) and `Шахмат қайда?` (“Where is chess?”) concern the same class but ask for different records. The **question type**, often called *intent*, is the kind of information wanted: time or place here. It is not the information itself. A form saying “item requested” is not the delivered item.

## The whole program first

Create `intent.py` in UTF-8. Use `қашан` (“when”) and `қайда` (“where”) as exact teaching clues. If both occur, preserve both types rather than silently choosing the first. `text` stores the lowercase question, `clean` removes `?` and `,`, and `words` stores the resulting words.

```python
text = input("Сұрақ: ").lower()
clean = text.replace("?", "").replace(",", "")
words = clean.split()
intents = []
if "қашан" in words:
    intents.append("уақыт")
if "қайда" in words:
    intents.append("орын")
print(intents)
```

The [ready-to-run file](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-14) has the same program.

`words` is the list produced by splitting, as in lesson 7. `intents = []` is an empty list of detected types. Two separate `if` statements test both words: the second still runs after the first matched. `append` from lesson 13 adds a label. `уақыт` means “time” and `орын` means “place”; these are **labels for looking up later records**, not a found time or location. `Шахмат қашан?` gives `['уақыт']`, `Шахмат қайда?` gives `['орын']`, and `Шахмат қашан қайда?` gives `['уақыт', 'орын']`. As in lesson 13, `replace` removes both `?` and commas, so `Шахмат қашан, қайда?` also gives both types. Other punctuation still remains.

`Шахмат неге?` gives `[]`: `неге` (“why?” or “to what?” depending on context) is outside our two supported types. An empty list does not make the question meaningless. Word order and endings are not analyzed yet; a match must be a whole word after `split()`.

## Memory map

Question → word list → `қашан` present? add “time” → `қайда` present? add “place” → keep every found type.

## Recall and task

Hide the code and explain why we use two `if` statements, not `if`/`elif`. Run `Сурет қашан?`, `Сурет қайда?`, `Сурет қашан қайда?`, and `Сурет неге?`.

**Task:** predict the outputs for `ҚАЙДА?` and an empty input. **Hint:** `lower()` changes case; the list may stay empty. **Answer:** `['орын']` and `[]`. **Common mistake:** seeing `['уақыт']` and thinking the program already knows the schedule. It knows only the kind of fact needed; the facts come in lessons 16–20.
