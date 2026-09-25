# Lesson 8. Build a dictionary of known stems

## Why this matters

A guest list tells a door attendant who is invited; a similar name is not enough. Our **stem dictionary** likewise lists only forms the model knows. For now we enter one bare stem without an ending. We will learn `үйлер` and `мектепте` later.

## The whole program

Create `roots.py` in the same folder and UTF-8 encoding as the previous files. Run `python3 roots.py` or `py roots.py`.

```python
roots = {"үй": "баспана", "мектеп": "оқу орны"}
word = input("Негіз: ").lower()
if word in roots:
    print(roots[word])
else:
    print("білмеймін")
```

The ready-to-run file is in the [course project folder](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-08); read the code explanation in this lesson first.

Input `үй` gives `баспана` (“home”), `мектеп` gives `оқу орны` (“place of study”), and `кітап` gives `білмеймін` (“I don't know”). Run the file again for each trial.

## Every new mark explained

Curly braces `{...}` create a **dictionary** of key → value pairs, like labeled drawers. The key `"үй"` points to `"баспана"`. A colon links a key to its value, and a comma separates pairs. In `roots[word]`, square brackets mean “take the value for the key held in `word`.” Lesson 7 also displayed square brackets for a list. Here, because `roots` is a dictionary, they perform a key lookup.

`if word in roots:` asks whether that key exists. `if` means “if” and `in` means “is contained in.” The colon after the condition starts a block of actions. Four spaces before `print` show that this line runs only when the condition holds. `else:` means “otherwise”; its indented line runs for an unknown key. We use `lower()` because our keys are lowercase. Matching is exact: a similar word does not become known automatically.

Empty input is unknown too. We have not yet removed spaces around the word or recognized endings. Refusing an unknown word is better than guessing a meaning that we cannot support.

## Memory map

One entered stem → lowercase → is the key in the dictionary? → known value or `білмеймін`.

![Lesson 8 support map](/static/course/kazakh-ai/map-08-roots-en.svg)

## Recall and task

Hide the code, draw two labeled drawers for `үй` and `мектеп`, and reconstruct the `if`/`else` block. Then try `ҮЙ`, `үй`, `үйлер`, and `кітап`. **Hint:** uppercase becomes lowercase; no ending is removed. **Answer:** the first two give `баспана`, the last two `білмеймін`. **Common mistake:** removing the spaces before `print(roots[word])` leaves Python unable to tell which line belongs to `if`. If the bracket line fails, check that you used ordinary `"` marks and closed every bracket.
