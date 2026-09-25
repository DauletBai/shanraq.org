# Lesson 7. Split a short question into words

## Why this matters

A librarian makes separate cards before looking up individual books. Our model likewise needs separate words before it can find stems. A **token** here means one separated word; **tokenization** is the act of obtaining those words from a line.

## The whole program

Create `words.py` in UTF-8 inside `ai-course`. Run it as you ran `letters.py` in lesson 6.

```python
text = input("Сұрақ: ")
clean = text.lower().replace("?", "").replace(",", "")
words = clean.split()
print(words)
```

The ready-to-run file is in the [course project folder](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-07); read the code explanation in this lesson first.

Enter `Үйлерде, мектептерде?`. The result is `['үйлерде', 'мектептерде']`. Square brackets mark a **list**: several values in order. A comma separates items, and quotation marks show that each item is text. Python prints these marks; you do not type them.

## Explain every step

`text.lower()` makes a new lowercase string. The dot says “apply the action to the text on the left.” `replace("?", "")` replaces each question mark with an empty string; `""` is text with zero characters. The next `replace` removes commas the same way. Each action produces new text, and we store the result as `clean`. `split()` separates it at spaces and other whitespace. Empty parentheses mean this action needs no extra instruction. The [Python reference](https://docs.python.org/3/library/stdtypes.html#string-methods) describes these methods.

This is a **small teaching rule**. It leaves full stops, exclamation marks, and quotation marks untouched; it also simply removes a comma inside a word. Do not treat it as a complete tokenizer for every text. An empty input gives an empty list, `[]`: no words were found.

## Memory map

Question → lowercase → remove `?` and `,` → split at whitespace → list of words. At each arrow, ask what could still remain.

![Lesson 7 support map](/static/course/kazakh-ai/map-07-words-en.svg)

## Recall and task

Hide the code and explain each dot and pair of parentheses in `clean = ...`. Then enter `Мектептерде үйлер бар ма?`. **Hint:** remove only `?` first, then split at spaces. **Answer:** `['мектептерде', 'үйлер', 'бар', 'ма']`. `бар` is only a separate token so far; the program does not know its meaning. **Common mistake:** without `lower()`, the capital `М` remains and will not match a later dictionary of lowercase stems.
