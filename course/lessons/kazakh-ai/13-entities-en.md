# Lesson 13. Find the subject of a question

## Why this matters

If you ask a librarian “When does the chess club open?”, they first need to know **which club** you mean. Our fictional study club has two lesson names: `Шахмат` (“chess”) and `Сурет` (“drawing”). Call the named thing the **subject of the question**. In software it is often called an *entity*: a specific item we may later look up in a catalog.

## The whole step first

For now we accept only bare names without endings. Lesson 11 showed how an ending can be checked, but we have not written those rules for club names. We collect every name found, so a question mentioning two lessons does not silently lose one.

`text` will store the lowercase question; `clean` will store it without `?` and `,`; `words` will be its list of words. Two new actions before the code: `for` visits list items in turn, and `append` adds an item to the end of another list. Indentation marks which lines repeat.

Create `entities.py` in UTF-8 and run it.

```python
known = ["шахмат", "сурет"]
text = input("Сұрақ: ").lower()
clean = text.replace("?", "").replace(",", "")
words = clean.split()
found = []
for word in words:
    if word in known:
        found.append(word)
print(found)
```

The [runnable step file](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-13) matches the printed code.

`found = []` creates an empty list, like a blank sheet for names we find. `for word in words:` means “take each word in the list, one after another.” A colon and four spaces mark the actions to repeat. The inner `if`, indented by **eight** spaces, checks each word. `found.append(word)` adds a match to the end of the list. After the loop, unindented `print(found)` runs once. Lesson 7 covered lowercase, two punctuation replacements, and `split()`; lesson 8 covered `in`. Repeated names remain repeated findings for now.

`Шахмат қашан?` gives `['шахмат']`. `Шахмат, сурет қашан?` gives `['шахмат', 'сурет']`. `Робот қашан?` gives `[]`: we have not put `Робот` in this tiny dictionary, even though a human may recognize it. An empty list means “not found in our catalog,” not “this club does not exist.” The program does not yet look up a lesson time.

## Memory map

Question → lowercase and two punctuation marks → words → visit each word → compare with two known names → list of findings.

![Lesson 13 support map](/static/course/kazakh-ai/map-13-entities-en.svg)

## Recall and task

Hide the code and explain why `print(found)` has no indentation. Then run `Сурет қайда?`, `Шахмат, сурет қашан?`, and `Робот қайда?`.

**Task:** predict the results for `СУРЕТ қашан?` and `Шахмат шахмат қашан?`. **Hint:** `lower()` does not remove repetitions. **Answer:** `['сурет']` and `['шахмат', 'шахмат']`. We will handle repetitions separately. **Common mistake:** placing `print(found)` inside the loop, which prints intermediate lists after every word instead of one final list.
