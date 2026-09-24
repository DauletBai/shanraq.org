# Lesson 26. Compare our learning parser with qazaq-ir

## Why this matters

Two people inspect a bicycle: a learner names three familiar parts, while a workshop also records each part’s condition. Compare only the work they have in common. Our lesson 11 recognizes two preset forms; `qazaq-ir` is a separate Rust project that returns JSON with a root, suffix chain, and analysis status. It is not a source for club schedules.

## Before the code

The program always shows the learning analysis `мектептерімізде → мектеп + тер + іміз + де`. The second experiment is optional. `sys.argv` is the list of command words: its first item names the program, and the next may be a path to the `qazaq-ir` executable. `Path(...).is_file()` checks that the file exists. `subprocess.run` starts the external program with an argument list rather than a shell. `capture_output=True` stores its output, `text=True` decodes it as text, and `check=False` lets us inspect failure ourselves. `timeout=10` limits waiting, and `returncode` reports success. `json.loads` reads its JSON, and `tokens` is its list of analyzed words. The first experiment works without Rust installed.

[qazaq-ir source project](https://github.com/qazaq-ai/qazaq-ir).

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-26
python3 compare.py
```

On Windows, replace `python3` with `py`.

```python
import json
import subprocess
import sys
from pathlib import Path

word = "мектептерімізде"
ours = {"root": "мектеп", "parts": ["тер", "іміз", "де"]}
print("Оқу үлгісі:", word, "→", ours["root"], ours["parts"])
if len(sys.argv) == 1:
    print("qazaq-ir: қосымша салыстыру іске қосылмады")
else:
    binary = Path(sys.argv[1])
    if not binary.is_file():
        print("qazaq-ir: бағдарлама файлы табылмады")
    else:
        run = subprocess.run([str(binary), "analyze", "--format", "compact", word],
                             text=True, capture_output=True, timeout=10, check=False)
        if run.returncode != 0:
            print("qazaq-ir: іске қосу қатесі")
        else:
            result = json.loads(run.stdout)
            for token in result["tokens"]:
                print("qazaq-ir:", token["surface"], "→", token["root"],
                      token["analysis_status"])
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-26).

## How the program works

With no second argument, the program prints `qazaq-ir: қосымша салыстыру іске қосылмады`. For an optional comparison, build the project with `cargo build --release -p qazaq-ir-cli` and pass the resulting `target/release/qazaq-ir` path (`.exe` on Windows). Installing Rust is a separate task; the Python course does not require it. A local check with `qazaq-ir` 0.31.0 returned `root: мектеп` and `analysis_status: partial` for this form: the root was found, but a full parse was not confirmed. Another version may differ. Compare `root` and `analysis_status` in the project JSON. Agreement on one word does not establish which system is “better AI,” nor does it verify a timetable fact or its source.

## Support map

One word → learning parse → optional qazaq-ir CLI → compare root and status → state both limits.

## Recall and check

Hide the code and recall the two lines available without Rust: `Оқу үлгісі: мектептерімізде → мектеп ['тер', 'іміз', 'де']` and the skipped-comparison message. Hint: follow `len(sys.argv) == 1`. Then pass a nonexistent file path and find `бағдарлама файлы табылмады`. Common mistake: treating a matching root as proof of a schedule. Word analysis and fact checking are different stages.
