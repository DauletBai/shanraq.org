# Lesson 26. Compare our learning parser with Adam FST

## Why this matters

Two people inspect a bicycle: a learner names three familiar parts, while a workshop also records each part’s condition. Compare only the work they have in common. Our lesson 11 recognizes two preset forms; `adam_fst` is a separate Rust tool that analyzes a Kazakh word into a root and grammatical features. It is not a source for club schedules.

## Before the code

The program always shows the learning analysis `мектептерімізде → мектеп + тер + іміз + де`. The second experiment is optional. `sys.argv` is the list of command words: its first item names the program, and the next may be a path to the `adam_fst` executable. `Path(...).resolve()` makes it absolute, and `is_file()` checks that it exists. `subprocess.run` starts the external program with an argument list rather than a shell. `cwd` temporarily selects the Adam repository root because the tool reads its dictionary from the `data` directory. `capture_output=True` stores its output, `text=True` decodes it as text, and `check=False` lets us inspect failure ourselves. `timeout=10` limits waiting, and `returncode` reports success. The first experiment works without Rust installed.

[Adam source project](https://github.com/qazaq-ai/adam). Its morphology CLI lives in the `adam-kernel-fst` module.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-26
python3 compare.py
```

On Windows, replace `python3` with `py`.

```python
import subprocess
import sys
from pathlib import Path

word = "мектептерімізде"
ours = {"root": "мектеп", "parts": ["тер", "іміз", "де"]}
print("Оқу үлгісі:", word, "→", ours["root"], ours["parts"])
if len(sys.argv) == 1:
    print("adam_fst: қосымша салыстыру іске қосылмады")
else:
    binary = Path(sys.argv[1]).expanduser().resolve()
    if not binary.is_file():
        print("adam_fst: бағдарлама файлы табылмады")
    else:
        run = subprocess.run([str(binary), "analyse", word], cwd=binary.parents[2],
                             text=True, capture_output=True, timeout=10, check=False)
        if run.returncode != 0:
            print("adam_fst: іске қосу қатесі")
        else:
            print("adam_fst:", run.stdout.strip())
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-26).

## How the program works

With no second argument, the program reports that the optional comparison was skipped. To try it, clone the public Adam repository, run `cargo build --release -p adam-kernel-fst --bin adam_fst` there, and pass the path to the resulting `target/release/adam_fst` file (`adam_fst.exe` on Windows) while leaving it inside the repository. Our program runs `adam_fst analyse мектептерімізде` from the Adam root and displays the tool's answer beside the learning analysis. Installing Rust is a separate task; the Python course does not require it. Compare only the detected root and features. Agreement on one word does not establish which system is “better,” nor does it verify a timetable fact or its source.

## Support map

One word → learning parse → optional Adam FST CLI → compare root and features → state both limits.

![Lesson 26 support map](/static/course/kazakh-ai/map-26-qazaq-ir-en.svg)

## Recall and check

Hide the code and recall the two lines available without Rust: `Оқу үлгісі: мектептерімізде → мектеп ['тер', 'іміз', 'де']` and the skipped-comparison message. Hint: follow `len(sys.argv) == 1`. Then pass a nonexistent file path and find `бағдарлама файлы табылмады`. Common mistake: treating a matching root as proof of a schedule. Word analysis and fact checking are different stages.
