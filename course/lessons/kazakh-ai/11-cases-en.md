# Lesson 11. Tell “to our houses” from “in our houses”

## Why this matters

Imagine two tickets: one takes you **to** a museum; the other admits you **inside** it. Both mention the museum, but their jobs differ. A Kazakh word can likewise keep its stem and possession while changing its final part: `үйлерімізге` means “to our houses”; `үйлерімізде` means “in our houses.”

## See the whole pattern first

In lesson 10, we read `үй-лер-іміз-ге` as house + plural + our + direction. Here `-ге` marks direction; it is a **dative case** ending, answering “to what?”. `-де` marks location; it is a **locative case** ending, answering “where?”. A **case** marks a word's relationship to other words in a sentence: direction or location here. This part follows the possession part. For our second familiar stem, we get `мектеп-тер-іміз-ге` and `мектеп-тер-іміз-де`. The order and ending variants follow the [Kazakh grammar description](https://slaviccenters.duke.edu/sites/slaviccenters.duke.edu/files/file-attachments/kazakh-grammar.pdf).

## Runnable step

Before the code, meet its new marks. `endswith` checks the word ending; `removesuffix` removes an ending already checked; `elif` means “otherwise, if,” and the variable `case` stores the case name. The vertical bar `|` separates study parts in the output; it is not written in the word.

Create `case.py` in UTF-8 inside `ai-course`. Run `python3 case.py` on macOS/Linux or `py case.py` on Windows.

```python
known = {"үйлеріміз": "үй | лер | іміз",
         "мектептеріміз": "мектеп | тер | іміз"}
word = input("Сөз: ").lower()
if word.endswith("ге"):
    base = word.removesuffix("ге")
    case = "барыс"
elif word.endswith("де"):
    base = word.removesuffix("де")
    case = "жатыс"
else:
    base = ""
    case = ""
if base in known:
    print(known[base], "|", case)
else:
    print("білмеймін")
```

The [ready-to-run step file](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-11) has the same code. Read the explanation below before using it.

`endswith("ге")` asks whether the text ends with exactly these letters; the result is yes or no. `removesuffix("ге")` returns new text without that checked ending. We check before removing so that an unrelated ending does not become a known one. `elif` means “otherwise, if” and runs only if the first condition did not match. `else` handles the remaining cases. The empty string `""` means no recognized part was found. `base` and `case` are variable names; `case` is not a special Python command. The commas inside `print` place spaces between output parts. `|` is a visible boundary for study, not part of the Kazakh word.

We confirm only two “our …” bases and two final parts. The program does not validate all Kazakh forms. It does not accept `үйімізге`, `мектептерімізге?` with punctuation, or other endings. A refusal marks the table's limit; it does not say the word is wrong.

## Memory map

Entered word → check final `ге`/`де` → remove it → look up the rest → name direction/location or say `білмеймін`.

![Lesson 11 support map](/static/course/kazakh-ai/map-11-cases-en.svg)

## Recall and task

Hide the code. Using the map, explain why we cannot accept every word ending in `ге`. Then run `үйлерімізге`, `мектептерімізде`, `үйімізге`, and `үйлерімізді` separately.

**Task:** predict the outputs for `МЕКТЕПТЕРІМІЗГЕ` and `үйлерімізде` before running. **Hint:** `lower()` changes the case; the text left after removing the ending must occur in `known`. **Answer:** `мектеп | тер | іміз | барыс` and `үй | лер | іміз | жатыс`. `үйімізге` and `үйлерімізді` produce `білмеймін`. **Common mistake:** cutting off the last two letters without checking them is like tearing a ticket before reading its destination.
