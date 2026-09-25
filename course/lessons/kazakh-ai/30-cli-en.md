# Lesson 30. Release the command-line model and check our claims

## Why this matters

A first-aid kit is useful when each item has a name, expiry, and instructions. Our release is also a complete kit: program, rules, exercise cards, a run command, and checkable boundaries. It answers only two kinds of short question about a fictional club.

## Before the code

Copy **all five files** from step 30: `main.py`, `engine.py`, `checks.py`, `facts.json`, and `examples.json`. The command takes a date and question; quotes keep a spaced question as one argument. Put optional `--trace` **before the date** to print the decision trail. `sys.argv` stores arguments; `len` selects one of two accepted command forms. `sys.exit(2)` stops bad usage with an error code. `try` attempts to parse the date, while `except ValueError` catches a bad format and prints a clear message. The date is the **day being asked about**, not a claim that an exercise record is true now.

Files in this step: `examples.json`, `checks.py`, `engine.py`, `facts.json`, `main.py`.

From the project root, enter the step folder and run the program:

```text
cd course/kazakh-ai/step-30
python3 main.py 2026-09-24 "Шахмат қашан?"
```

On Windows, replace `python3` with `py`.

```python
import sys
from datetime import date

from engine import answer

if len(sys.argv) == 3:
    show_trace = False
    raw_date = sys.argv[1]
    question = sys.argv[2]
elif len(sys.argv) == 4 and sys.argv[1] == "--trace":
    show_trace = True
    raw_date = sys.argv[2]
    question = sys.argv[3]
else:
    print('Қолдану: python3 main.py [--trace] YYYY-MM-DD "Сұрақ"')
    sys.exit(2)
try:
    on_date = date.fromisoformat(raw_date)
except ValueError:
    print("Қате күн: YYYY-MM-DD түрінде жазыңыз")
    sys.exit(2)
result = answer(question, on_date)
print(result["text"])
if show_trace:
    for line in result["trace"]:
        print("Із:", line)
```

[Step files](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-30).

## How the program works

From the step folder, run `python3 main.py 2026-09-24 "Шахмат қашан?"`. It prints `шахмат үйірмесінің уақыты: бейсенбі, 15:00 | 2026-09-24 күнгі дерек | дереккөз: club-sheet-01`. `Сурет қайда?` yields `білмеймін: дерек жоқ`; `Шахмат қашан қайда?` asks for one type; date `2027-01-01` refuses an expired record. Add `--trace` to see five source-check steps. The small held-out set in lesson 25 tested question types only; it does not establish real service quality. We built a narrow local model, not universal AI. No paid API call in exercise code does not mean zero total cost; avoiding free-form generation does not eliminate source errors.

## The finished model's specification sheet

A model specification is like the plate on an electrical appliance: units and test conditions replace words such as “fast” and “small.” The step folder has two more files. `quality_cases.json` contains 16 checks written before the run; `benchmark.py` runs them and measures the finished system. Run:

```text
python3 benchmark.py --json
```

A **cold start** begins in a new Python process. A **warm request** runs inside an existing process. `p95` is the time within which 95 percent of requests finished. **RSS** is the peak physical memory of the entire process; `tracemalloc` counts only observed Python allocations. A **trainable parameter** is a number changed while a neural network learns. This system has none: it builds intent counts from 13 cards at startup.

**Throughput** is the number of requests completed per second. One **worker process** is one running copy of Python; in this experiment it handles requests in sequence and uses no more than one core at once. The program then warms 1, 2, 4, and 8 processes, gives each process 10,000 requests three times, and selects the middle total rate. One KiB is 1,024 bytes, and one MiB is 1,024 KiB. ARM is the processor family of the computer used for the control run.

A control run on 25 September 2026 used an ARM computer with macOS 26.5.2 and Python 3.14.5:

On a phone, swipe the following wide tables left and right.

| Measurement | Result |
|---|---:|
| Trainable neural parameters | 0 |
| Learning cards / approved facts | 13 / 2 |
| Available logical cores / worker processes / concurrent requests | 8 / 1 / 1 |
| Model code and working data | 6,907 bytes (6.7 KiB) |
| Median cold start | 34.682 ms |
| Median warm request / `p95` | 0.034 / 0.037 ms |
| Throughput of one sequential process | 42,520 requests/s |
| Peak `tracemalloc` allocations | 135,626 bytes |
| Peak RSS of the whole Python process | 25,608,192 bytes (24.4 MiB) |
| All scenarios / correct supported answers | 16 of 16 / 6 of 6 |
| Correct refusals | 10 of 10 |
| Confident answer where refusal was required | 0 of 10 |
| External API requests | 0 |
| Direct API fee | 0; device, electricity, and labor are not free |
| Energy per request | not measured; a hardware power meter is required |

### What additional processes delivered

**Speedup** compares total throughput with one worker. **Parallel efficiency** divides speedup by the worker count: 100% would mean that every new worker added its full share of speed. Summed RSS below adds the worker peaks; it excludes the coordinator process.

| Worker processes | Concurrent requests | Total requests/s | Speedup | Efficiency | Summed peak RSS |
|---:|---:|---:|---:|---:|---:|
| 1 | 1 | 42,953 | 1.00× | 100% | 25.4 MiB |
| 2 | 2 | 80,372 | 1.87× | 93.6% | 50.8 MiB |
| 4 | 4 | 146,241 | 3.41× | 85.1% | 101.7 MiB |
| 8 | 8 | 141,357 | 3.29× | 41.1% | 204.0 MiB |

Two and four workers raised total throughput. Eight were slightly slower than four and used roughly twice the worker memory. The experiment reveals a limit of this computer and program; it does not by itself identify the cause. Core types, shared files, memory, or operating-system scheduling may contribute. A different server needs a new measurement.

The last quality row gives “hallucination” an operational meaning for this experiment. An error occurs when the reference requires refusal but the program confidently produces another answer. `0 of 10` applies only to the published checks. It does not establish zero errors on arbitrary text and does not protect against a wrong fact in `facts.json`.

## Comparison with large LLMs

The makers of closed GPT, Claude, and Gemini models do not publish their weight counts or server memory. The claim that every parameter of every major LLM is known is therefore false. Open model scale is published; memory still depends on weight precision and deployment.

A stored parameter value is called a **weight**. FP8 and BF16 are two ways to store a weight, in about 1 or 2 bytes. Where the parameter count is known, multiplication therefore gives a reproducible lower estimate. These figures cover weights only; working memory, context cache, and serving data are additional. A **token** is a word or word part into which an LLM divides text; **context** is the set of tokens available to it in one request. The `≈` sign below marks a calculated figure, not one measured on our computer.

| System | Published scale | Memory or size | The same 16 checks |
|---|---|---|---|
| Our learning model | 0 neural parameters; 13 cards; 2 facts | measured: 6,907 file bytes; 24.4 MiB RSS for one Python process | measured: 16/16; unsupported confident answers 0/10 |
| Llama 3.1 405B | 405B parameters; 128K-token context | ≈405 GB FP8 or ≈810 GB BF16: 58.6–117 million times our file size | not run on our suite |
| DeepSeek-V3 | 671B total, 37B active per token; 128K context | ≈671 GB FP8 or ≈1.342 TB BF16: 97–194 million times larger; working memory is additional | not run on our suite |
| Qwen3-235B-A22B | 235B total, 22B active | ≈235 GB FP8 or ≈470 GB BF16: 34–68 million times larger; working memory is additional | not run on our suite |
| GPT-5.6 Sol / Gemini 3.7 Flash | parameters and server RAM undisclosed; context about 1M tokens | cannot be calculated from public data: context and price do not determine weight count | not run on our suite |

The “one thousand times” threshold is exceeded for the **size of these stated representations**. That factor cannot be transferred to intelligence, task range, speed, or quality: an LLM handles free text and many tasks, while our model knows two question types and two facts. We also do not substitute someone else's number for LLM speed. A valid ratio requires running a named version on the same hardware, or through a fixed API, with the same 16 requests.

Scale sources: [Meta on Llama 3.1 405B and 128K context](https://ai.meta.com/blog/meta-llama-3-1/), [DeepSeek-V3 model card](https://github.com/deepseek-ai/DeepSeek-V3), [NVIDIA requirements for DeepSeek-V3](https://github.com/NVIDIA/TensorRT-LLM/tree/main/examples/models/core/deepseek_v3), [Qwen3 description](https://qwenlm.github.io/blog/qwen3/), [official GPT-5.6 specifications](https://platform.openai.com/docs/models), [official Gemini 3.7 Flash specifications](https://ai.google.dev/gemini-api/docs/models/gemini-3.7-flash).

## Where our model is stronger and where it falls short

Comparison helps us choose a tool. A pocket calculator multiplies faster than a person but cannot summarize a book; a librarian understands an open question but uses more resources to search. Our model and an LLM differ in the same way.

| Property | Our deterministic model | Large LLM | Project decision |
|---|---|---|---|
| Narrow approved facts | returns only a catalog record and names its source | can compose fluent, plausible text | our model is stronger when an answer must match an approved card |
| Input outside the boundary | refuses and gives an explicit reason | may answer, but its factual support needs a separate check | refusal is more useful than a guess when errors are costly |
| Speed and resources for this task | measured locally: kilobytes of files, milliseconds, and tens of MiB including Python | large weights or a remote API | the narrow model fits a small fixed catalog better |
| Open language and broad topics | cannot do this | can explain, generalize, translate, and handle varied requests | an LLM is far stronger on open tasks |
| New words and paraphrases | needs a new rule or example | often understands without a program change | our model generalizes poorly to unfamiliar wording |
| Knowledge updates | a person edits a small auditable file | knowledge in weights cannot be fixed with one card; retrieval or a database can be added | targeted review is easier in our model, but the catalog needs manual care |
| Repeatability | identical code, data, and input give the same path and answer | output depends on version, settings, and serving system | our narrow model is easier to reproduce byte for byte |
| Request transmission | the learning code makes no network request | a cloud API receives the request; a locally deployed LLM can keep it on device | the advantage depends on deployment, not only model type |
| Decision audit | each step can be replayed line by line | the complete path through a huge network cannot be inspected | our model is more transparent inside its narrow boundary |
| Task scale | two intents, two facts, and a few known word forms | many languages, subjects, and formats | 16/16 cannot be extended beyond the learning suite |

There is a third option: a **hybrid system**. An LLM can interpret open wording, while a deterministic component checks allowed fields, date, and source before returning the fact. The combination is more complex than either component alone, but joins a broad input with a controlled fact.

Reconstruct the selection rule from three questions:

1. Is the fact set small and approved, with a high cost for invention? Start with a deterministic model.
2. Does the user need open text, with too many possible phrasings to list? Consider an LLM plus a separate fact check.
3. Are both properties required? Measure the hybrid on one sealed suite instead of adding together the promises of two technologies.

## Support map

Command + date + question → word guard → type → one card → valid date → sourced answer; otherwise an explained refusal.

![Lesson 30 support map](/static/course/kazakh-ai/map-30-cli-en.svg)

## Recall and check

Hide the code and reconstruct the support map. Try four commands: a known question, unknown club, expired date, and `--trace`. Hint: `2026-09-24` is between `checked_on` and `valid_until`, while `2027-01-01` is not. Expected for the expired date: `білмеймін: бұл күнге жарамды дерек жоқ`; for an unknown club: `білмеймін: таныс емес сөз`. Common mistake: presenting an exercise date as a real club’s current schedule.

## Gate 6: final defense

Show how an answer depends on a card and date; obtain the same answer for `Шахмат қашан?` and the known form `Шахматтың уақыты қашан?`; name two justified refusals; run the specification and explain the difference between RSS and Python allocations. Then explain why 16/16 is not perfection and why a size factor cannot be transferred to intelligence. Pass when you can guide a new learner through the chain without the lesson text. If you get stuck, return to the matching one of the six support signals and defend the project again.
