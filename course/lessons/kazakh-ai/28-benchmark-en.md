# Lesson 28. Measure time, memory, and cost

## Why this matters

A kettle takes time to heat from cold; pouring already-hot water is faster. A program likewise has a fresh-start cost and a cost for calling a ready function again. Measure both, file size, and observed memory allocations.

## Before the code

Copy `bench.py`, `checks.py`, and `examples.json` together. `subprocess.run` starts a fresh Python process five times; `sys.executable` selects the same Python version. `-c` runs the following code string, `capture_output=True` stores its output, `text=True` reads it as text, and `check=True` stops the measurement on failure. `range(5)` supplies five loop values; `range(200)` supplies 200. `time.perf_counter_ns()` reads a high-resolution clock in nanoseconds; one million nanoseconds make a millisecond. `statistics.median` chooses the middle measurement, reducing the effect of one slow outlier. After one warm-up, 200 calls measure the ready function. `tracemalloc.start()` tracks **Python** allocations, `get_traced_memory()` returns current and peak tracked bytes, and `stop()` ends tracking. `Path.stat().st_size` gives the byte size of two files.

Files in this step: `examples.json`, `checks.py`, `bench.py`.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-28
python3 bench.py
```

On Windows, replace `python3` with `py`.

```python
import statistics
import subprocess
import sys
import time
import tracemalloc
from pathlib import Path

from checks import classify

question = "Шахмат қашан?"
cold = []
for repeat in range(5):
    start = time.perf_counter_ns()
    subprocess.run([sys.executable, "-c", "from checks import classify; classify('Шахмат қашан?')"],
                   check=True, capture_output=True, text=True)
    cold.append(time.perf_counter_ns() - start)
classify(question)
times = []
tracemalloc.start()
for repeat in range(200):
    start = time.perf_counter_ns()
    classify(question)
    times.append(time.perf_counter_ns() - start)
current, peak = tracemalloc.get_traced_memory()
tracemalloc.stop()
files = Path("checks.py").stat().st_size + Path("examples.json").stat().st_size
print("Сұрау саны:", len(times))
print("Жаңа іске қосу, мс:", round(statistics.median(cold) / 1000000, 3))
print("Ортаңғы уақыт, мс:", round(statistics.median(times) / 1000000, 3))
print("Ең көп бақыланған бөлу, байт:", peak)
print("Код пен дерек, байт:", files)
print("Тікелей API төлемі: 0; құрылғы мен еңбек құны есептелмеді")
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-28).

## How the program works

Numbers depend on the machine, load, and Python version, so the reference result is the set of output labels and positive measurements, not fixed milliseconds. Warm-function time excludes loading `examples.json`, Python startup, and fact lookup. `tracemalloc` does not report total process memory; file bytes are not RAM use. Zero direct API fee means only that this experiment makes no paid external request. Hardware, electricity, labeling, and maintenance still cost money. A “thousand times” comparison with an LLM is not established here; comparable inputs, quality, and conditions would be required.

This lesson measures only the question classifier. Lesson 30 applies the same idea to the finished system, from the question through fact, date, and refusal checks.

## Support map

5 fresh processes → median startup; 200 calls → median function time + allocations; files → bytes; API fee → direct fee only.

![Lesson 28 support map](/static/course/kazakh-ai/map-28-benchmark-en.svg)

## Recall and check

Hide the code and explain why the two time measurements differ. Run the program twice; do not copy another machine’s numbers. Expected labels include `Сұрау саны: 200`, fresh-start time, middle call time, allocations, file bytes, and direct API fee. Hint: a fresh process loads the module and data again. Common mistake: calling `peak` total application RAM or declaring victory over an LLM from one number.
