# Lesson 2. The first checkable program step

## Why this matters

A person can read a timetable. A program needs an exact instruction. First, we will display one approved record and distinguish source code from its output.

**Python** is the programming language for our learning model. **Code** is the text of instructions. A **terminal** is a window for computer commands. The command `python3` starts an installed Python; the name after it identifies the code file.

## Prepare a place for your file

Install Python from its [official downloads page](https://www.python.org/downloads/) if needed. On Windows, open PowerShell from Start and run `py --version`; it should show a version number. If `py` is missing, close and reopen PowerShell after installation. On macOS or Linux, open Terminal and run `python3 --version`; if it is missing, follow the [official platform guide](https://docs.python.org/3/using/index.html). Use Python 3.10 or newer.

Create a folder named `ai-course` on your desktop, then open it in the file manager. On Windows, type `powershell` into the folder's address bar and press Enter; PowerShell opens in that folder. On macOS or Linux, use the file manager's “Open in Terminal” action if available. Otherwise type `cd ` (including the space) in Terminal, drag the folder icon into the window, and press Enter. The **current folder** is where the command will look for `club.py`. Save the file there as UTF-8 plain text. In Windows Notepad, choose “All files” so its name does not become `club.py.txt`.

## See the whole process

Create a file named `club.py`:

```python
class_name = "Робот құрастыру"
time = "сейсенбі, 16:00"
print(class_name, time)
```

Run `python3 club.py` on macOS/Linux, or `py club.py` on Windows. Expected output: `Робот құрастыру сейсенбі, 16:00`. If the computer says it cannot find the file, check its name and the current folder.

## Every part explained

`class_name` and `time` are **variable names**: labels for stored values. Here `=` means “store the value on the right under the name on the left”; it is not an equation. Text between double quotation marks is a **string**, a sequence of characters. Parentheses in `print(...)` enclose what to display. The comma passes two values to `print`; a space appears between them in the output. The program does not yet analyze questions or verify sources. It only prints what we entered.

## Recall map

Code file → run command → Python executes lines from top to bottom → output. `name = value` stores; `print(...)` displays.

## Task and check

Without looking at the example, create a chess record for Thursday at 15:00. Predict the output before running it. Explain why this program cannot answer a question about art class.

**Hint:** change the values to the right of `=`. **Check:** `class_name = "Шахмат"`, `time = "бейсенбі, 15:00"`; output `Шахмат бейсенбі, 15:00`. No art class data was entered.
