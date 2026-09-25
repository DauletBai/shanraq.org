# Lesson 6. Read Kazakh Cyrillic without losing letters

## Why this matters

Change one letter on an envelope and it may go to the wrong address. For our model, `үй` and `уй` are also different spellings. Before analyzing a word, we need to receive exactly what the person typed.

## Setup and the whole program

Use the `ai-course` folder from [lesson 2](https://shanraq.org/read/kazakh-ai-02-first-program?lang=en). If you do not have it, follow that lesson's installation and launch steps first. Create a plain-text file called `letters.py` there. **UTF-8** is a way to save the letters in a file so Python reads them correctly. Select UTF-8 if your editor asks for an encoding.

```python
word = input("Сөз: ")
print(word)
```

The ready-to-run file is in the [course project folder](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-06); read the code explanation in this lesson first.

Run `python3 letters.py` on macOS/Linux or `py letters.py` on Windows. At `Сөз: `, type `үй` and press Enter. You will see `Сөз: үй` and then `үй` on the next line.

## Every step explained

`input(...)` pauses the program and waits for a line from the keyboard. The quoted text is a prompt for the person; it is not part of the answer. Enter ends the input. `word =` saves the answer under the name `word`; `print(word)` displays it. Lesson 2 explained variables, `=`, quotation marks, and `print`. The colon inside the quotes is just prompt text, not a Python command.

Kazakh Cyrillic has letters worth checking carefully: `ә ғ қ ң ө ұ ү һ і`, along with their uppercase forms. `ү` and `у` are distinct, like different digits in a phone number. Our program preserves what was entered; it does not yet judge whether the word is spelled correctly. Latin `y` does not stand for Cyrillic `у`.

## Memory map

Keyboard → `input` → text stored in `word` → `print` → the same text. UTF-8 preserves letters in the file; you still need to type the intended letters.

![Lesson 6 support map](/static/course/kazakh-ai/map-06-cyrillic-en.svg)

## Recall and check

Hide the code and rebuild both lines from the map. Explain where the word on the second output line came from. Run the program afresh with `үй`, `өнер`, and `қала`. Each time the second line should exactly match your input. If the letters look wrong, check your keyboard layout and the file encoding.

**Task:** change only the prompt to `Атау: ` and enter `мектеп`. **Hint:** change the string inside `input(...)`, not `word`. **Answer:** `word = input("Атау: ")`; the screen shows `Атау: мектеп`, then `мектеп`. **Common mistake:** `print("word")` displays the literal word `word`, because quotation marks make text, not a variable name.
