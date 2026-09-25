# Lesson 9. Six forms of the plural marker

## Why this matters

You can pack the same kind of item in boxes of different sizes. The meaning “more than one” stays, while the sound of the added part adjusts to its neighbors. In Kazakh, the **plural ending** has six forms: `-лар`, `-лер`, `-дар`, `-дер`, `-тар`, `-тер`. The hyphen shows a boundary for study; it is not written in the complete word.

## See the table first

On a narrow screen, swipe the table left to see its final column.

| Stem | Plural part | Whole word | Notice |
|---|---|---|---|
| `қала` — city | `-лар` | `қалалар` | last vowel is back; ends in a vowel |
| `үй` — house | `-лер` | `үйлер` | last vowel is front; ends in `й` |
| `қалам` — pen | `-дар` | `қаламдар` | last vowel is back; ends in `м` |
| `тіл` — language | `-дер` | `тілдер` | last vowel is front; ends in `л` |
| `кітап` — book | `-тар` | `кітаптар` | last vowel is back; ends in `п` |
| `мектеп` — school | `-тер` | `мектептер` | last vowel is front; ends in `п` |

A **vowel** is a sound made without blocking the air in the mouth. In these examples `а` and `ы` are back vowels; `е`, `і`, and `ү` are front vowels. The last vowel helps choose the ending's vowel, and the final sound of the stem affects the opening consonant. We do **not yet program** every sound rule: we store six verified pairs. A [Kazakh grammar](https://slaviccenters.duke.edu/sites/slaviccenters.duke.edu/files/file-attachments/kazakh-grammar.pdf) describes the full set of variants.

## A working table

Create `plural.py` in UTF-8 and run it with the familiar command.

```python
plural = {"қала": "лар", "үй": "лер", "қалам": "дар",
          "тіл": "дер", "кітап": "тар", "мектеп": "тер"}
word = input("Негіз: ").lower()
if word in plural:
    print(word + plural[word])
else:
    print("білмеймін")
```

The ready-to-run file is in the [course project folder](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-09); read the code explanation in this lesson first.

The dictionary continues on a second line after a comma. Its leading spaces just make the code readable; Python sees that the curly brace is still open. Here `+` **joins two strings**, like two strips of paper: `"үй" + "лер"` produces `"үйлер"`. It is not arithmetic. Input `үй` produces `үйлер`, and `тіл` produces `тілдер`. Input `дала` produces `білмеймін` because that stem is not in our table. Do not guess an ending for an unknown stem.

## Memory map and check

Known stem → stored plural form → join without a hyphen → word. Unknown stem → `білмеймін`. Hide the table and name one example starting with each of `л`, `д`, `т`; then check all six program entries.

![Lesson 9 support map](/static/course/kazakh-ai/map-09-plurals-en.svg)

**Task:** predict the outputs for `ҚАЛАМ`, `кітап`, and `жол` before running the program. **Hint:** `lower()` changes case but adds no keys. **Answer:** `қаламдар`, `кітаптар`, `білмеймін`. **Common mistake:** `print("word + plural[word]")` prints literal text, not the built word. We will study unmarked number and context dependent meanings later; this is still a table for six stems.
