# A correct and sufficiently fast algorithm

_Lead (summary):_ **Prove an algorithm correct first, then compare its work on small and large lists.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![A correct and sufficiently fast algorithm](/static/course/informatics/map-25-correctness-efficiency-en.svg)

## Situation and question

Duplicate `id` values are easy to spot among three tasks. What about a thousand? We need a reliable check. Speed matters only after we define what counts as a duplicate and what answer is required.

## New words without gaps

**Correctness** means an algorithm meets its written contract for every valid input, not just one example. A **counterexample** is an input that breaks the promise. **Efficiency** compares work and memory as input grows. **Input size** `n` is the task count here. A **linear scan** visits each record roughly once. **Quadratic growth** occurs when each record is compared with every other: there are `n(n−1)/2` pairs. A **set** stores unique values and usually checks whether an `id` has appeared quickly, at the cost of extra memory. Complexity describes growth, not seconds guaranteed on one laptop.

## The lesson's support signal

`contract and counterexamples first → operations and memory second`

## Work through it step by step

Method A: compare each task `id` with every later one. Four tasks give 6 pairs; ten give 45. Method B: scan left to right with empty set `seen`; if `id` is already in it, report duplicate, otherwise add it. For `t-01,t-02,t-03,t-02`, B adds the first three and finds a duplicate on the fourth. Four records require at most four arriving values to be checked, but earlier IDs occupy memory. Compare `id`, not titles: distinct tasks may have the same title.

## Check it by hand

Write four IDs: `t-01, t-02, t-03, t-02`. For pairwise checking draw position pairs 1–2, 1–3, 1–4, 2–3, 2–4, 3–4; the duplicate appears in the fifth comparison. Now track `seen`: it holds three values after three steps, and `t-02` is already present at step four. When comparing work, name the operation being counted: pair comparisons or membership checks. The numbers 5 and 4 refer to different operations; actual device speed still needs measuring later.

## Predict and check

What about an empty list or one `id`? Neither has a duplicate. For `t-01,t-01`, the second is a duplicate. For `t-01` and `T-01`, the answer depends on a documented case rule; do not silently normalise case.

## Catch the error

“B is faster on every device and every input.” It usually needs fewer comparisons on large lists, but uses memory and depends on implementation. Choose by requirements and measure realistic sizes after writing code; a paper model is not a benchmark.

## Transfer to a new setting

Find the nearest due date in a list said to be sorted. Can you stop early? First establish that sorted order is guaranteed; otherwise early exit can give a wrong answer.

## Project change

In `ALGORITHM-en.md`, document both duplicate-ID checks, pair counts for `n=1,4,10`, and cases `t-01,t-02,t-01`, empty list, and equal titles with distinct IDs. Record the trade-off: fewer checks for extra memory.

## Task and evidence

Ask a classmate for a counterexample to comparing titles instead of IDs. They then justify B: after each step, `seen` contains exactly the IDs already visited. Evidence: a trace and invariant explanation.

## Return after 1, 7, and 30 days

Tomorrow reconstruct six pairs for four records. In seven days count comparisons for five. In a month compare the paper prediction with a measurement of your program.

[Next lesson: Mastery checkpoint: the assistant's algorithms on paper](/read/informatics-26-algorithms-mastery?lang=en)
