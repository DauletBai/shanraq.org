# Repetition without copying and a stopping condition

_Lead (summary):_ **Visit every task without copying commands and prove that repetition stops.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![Repetition without copying and a stopping condition](/static/course/informatics/map-23-loops-invariants-en.svg)

## Situation and question

The assistant has three tasks. Writing “check task” three times is easy. Tomorrow there may be thirty tasks, and copying the command thirty times would be absurd. We need one repeated step and a clear stopping point.

## New words without gaps

A **loop** repeats an action while elements remain or a condition holds. An **iteration** is one pass. A **counter** stores how many elements have been processed. A **stopping condition** says when repetition ends. A **loop invariant** is a statement true before and after every pass: “records before the current position have been checked; the others have not.” An **infinite loop** occurs when the exit condition never arrives, for example when position never advances. An **empty list** is an important test: the loop must stop without error and return zero completed tasks.

## The lesson's support signal

`position = 0 → while position < length → process → position + 1 → stop`

## Work through it step by step

For `t-01`, `t-02`, `t-03`, let `done` be false, true, false. Start at `position=0`, `completed=0`. First pass: false, count stays 0, position becomes 1. Second: true, count becomes 1, position 2. Third: false, count stays 1, position 3. Now `position < 3` is false and the answer is 1. Check the invariant after each pass: the count equals the number of true values in the part already visited. Advance before reading and the first task vanishes. Never advance and the loop never ends.

## Check it by hand

Line up three cards `false, true, false` and move a position token only to the right. Before each step, say how many cards have been processed and how many `true` values have been found; check both after the step. At position 3 nothing is left, so the token must stop. Repeat with an empty line: position already equals length, and the loop never starts. These two experiments show both the invariant and termination.

## Predict and check

What does an empty list return? Zero: the condition is false immediately. For `true, true, true`, the result is 3; for `false, false, false`, it is 0. Predict the passes for four records: exactly four if the index advances once per pass.

## Catch the error

“To include the last record, write `position <= length`.” At position equal to length there is no record: valid indexes run from 0 to `length − 1`. Fix the sign and test a one-item list.

## Transfer to a new setting

Instead of counting done tasks, find the first unfinished one. Can you stop before the end? Yes, once it is found; explain why early exit cannot count all done tasks.

## Project change

Add a `done=true` counting loop with a three-pass table to `ALGORITHM-en.md`. Include the invariant and stopping condition. Test empty list, one item, and the original three tasks. This algorithm reads data but does not change `tasks.json` yet.

## Task and evidence

Trace `false,true,false` and an empty list. Deliberately remove position advancement and show where the stopping proof fails. Evidence: table, invariant, and repaired loop.

## Return after 1, 7, and 30 days

Tomorrow explain why an empty list is safe. In seven days trace a search for the first unfinished item. In a month find a loop in your code and check its termination.

[Next lesson: A function as a promise: input, result, and conditions](/read/informatics-24-functions-contracts?lang=en)
