# Sequence and step tracing

_Lead (summary):_ **Trace every algorithmic action and see how a small ordering mistake changes the final result.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![Sequence and step tracing](/static/course/informatics/map-21-sequence-tracing-en.svg)

## Situation and question

A courier receives three instructions: collect the parcel, go to the door, ring the bell. If they ring before collecting the parcel, the door opens but no delivery happens. A computer also follows written order unless we specify another rule.

## New words without gaps

An **algorithm** is a finite description of actions for a defined class of inputs. A **sequence** executes steps one after another. A **trace** records values after every step, not only the final answer. A **step** is one precisely named action. **Input** is data received by the algorithm; **output** is what it returns or displays. A **side effect** is an extra change such as writing a file. Two algorithms can display the same screen now but behave differently after restart because they saved in a different order.

## The lesson's support signal

`input → step 1 → state 1 → step 2 → state 2 → output`

## Work through it step by step

Let `count=0`, with two fictional records whose `done` values are `false` and `true`. Paper commands: (1) `count ← count + 1`; (2) display `count`; (3) increment again; (4) display. The trace displays 1, then 2. Move both increments before both displays and you get 2, then 2. Almost the same commands, different order. For the assistant, the honest sequence is: read file, find record, change `done`, save file, confirm after a second read. Saying “saved” before writing would be a false promise.

## Check it by hand

Run a table-top trace: one learner holds the `count` card, another follows `+1`, “display,” `+1` exactly from left to right. After every step, write the step number, current value, and visible output in a table. Swap “display” with the final `+1`: the final `count` remains 2, yet the person sees a different number. That is why checking only the final state cannot validate a program with visible output.

## Predict and check

Start at `count=2`; steps are `+1`, display, `−2`, display. Do not calculate only the final value. Record each step: 3, screen 3, 1, screen 1. Explain why the screen need not show every intermediate value.

## Catch the error

“The final `count` is 1, so the algorithm is right.” The earlier message might have shown 3 when the requirement asked for 1. Compare the full trace with each observable requirement.

## Transfer to a new setting

Order “open file,” “validate fields,” “display tasks,” and “report error” so invalid JSON cannot be displayed as trusted tasks. Draw an error branch that never reaches display.

## Project change

Add a trace table for `t-01` to `ALGORITHM-en.md`: state before, action, state after, visible message. Trace both a successful save and a read failure. Mark the step that actually changes the file.

## Task and evidence

Make two traces: normal path and missing-file path. Exchange two steps and explain what observation changes. Evidence: tables and the exact first wrong step.

## Return after 1, 7, and 30 days

Tomorrow reconstruct the counter trace. In seven days find an ordering error in a new everyday instruction. In a month compare real program order with your paper trace.

[Next lesson: Conditions and boundary cases](/read/informatics-22-conditions-boundaries?lang=en)
