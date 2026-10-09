# State and variables through a game

_Lead (summary):_ **Use a game with three lives to understand state, then describe how the assistant changes a task.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![State and variables through a game](/static/course/informatics/map-20-state-variables-en.svg)

## Situation and question

A game shows three lives. After a trap there are two; after a bonus there are four. The screen changed, but what must be remembered between events? The assistant has a similar transition: a task goes from unfinished to done.

## New words without gaps

**State** is data describing the system at a chosen moment. A **variable** names a value the algorithm may change. **Assignment** replaces an old value with a new one; `lives ← lives − 1` is not a mathematical equality but an instruction to compute the right side first. An **initial value** is the state before the first event. A **state transition** is a before/after pair caused by one action. An **invariant** is a condition that must keep holding, such as lives never being negative. Do not confuse a paper variable with a saved JSON field: a temporary value disappears on shutdown unless written to storage.

## The lesson's support signal

`initial state → event → new value → check constraint`

## Work through it step by step

Make a table `event | lives before | action | lives after`. Begin at 3. Trap: `3 − 1 = 2`. Bonus: `2 + 2 = 4`. Three more traps yield 3, then 2, then 1. If a trap arrives at 0, the algorithm must keep 0 or end the game; choose the rule in advance. For `t-01`, the table is shorter: `done=false` → “mark as done” → `done=true`. Repeating that action should not unexpectedly return `false` if the button says “Mark as done.”

## Check it by hand

Draw three score boxes: before, action, after. For `3 − 2 + 4`, record `3 → 1` and then `1 → 5`; do not jump straight to five. Beside them draw two cards for the same task `id`: `done=false` before “mark complete” and `done=true` afterwards. Its title stays unchanged. Ask what a second completion action should do. The contract must say whether `true` remains `true`; otherwise the state is ambiguous.

## Predict and check

Begin with two lives and apply trap, +2 bonus, trap. Name each state before checking: 1, 3, 2. After two identical “mark as done” actions, what is `done`? It remains `true` if the action sets a value rather than toggling it. This matters when a command is sent twice.

## Catch the error

“`lives ← lives − 1` says a number equals itself minus one, so the algorithm is impossible.” The arrow means update: read the old value, compute a new one, store it. Use equality in mathematics and an arrow for the paper instruction.

## Transfer to a new setting

Invent a counter of pages read, initially 0, increased by each “page read” event. What if the event is accidentally sent twice? State a rule or two taps will record false progress.

## Project change

In `ALGORITHM-en.md`, add a state table for one task with `id`, `title`, and `done`. State that marking it done does not change `id`. Keep “save new done value to file” as a separate step; a paper variable is not a disk write.

## Task and evidence

Trace three game events and two task actions by hand. Evidence: before/after table, rule for the lower bound of 0, and result of a repeated command. Another learner should recover the final state using only your table.

## Return after 1, 7, and 30 days

Tomorrow explain the assignment arrow. In seven days trace a new counter. In a month identify which state in the program is temporary and which is saved.

[Next lesson: Sequence and step tracing](/read/informatics-21-sequence-tracing?lang=en)
