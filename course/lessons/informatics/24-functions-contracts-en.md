# A function as a promise: input, result, and conditions

_Lead (summary):_ **Extract the due-date check as one function with clear input, output, and failure behaviour.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![A function as a promise: input, result, and conditions](/static/course/informatics/map-24-functions-contracts-en.svg)

## Situation and question

The same “should we remind?” question appears in the task list and on a task card. Copy the rule twice and a later fix may reach only one copy. We need one function and a written contract for its behaviour.

## New words without gaps

A **function** is a named part of an algorithm that accepts input and returns a result. A **parameter** names an input inside its definition; an **argument** is the concrete value passed in. A **return value** is the answer, not necessarily a screen message. A **precondition** says which inputs are valid. A **postcondition** promises a result for a valid input. A **contract** collects these rules and the behaviour on failure. Functions may have side effects, but for a simple date check we explicitly forbid changing the file.

## The lesson's support signal

`input + precondition → one function → result + guarantee`

## Work through it step by step

Describe `should_remind(done, days_to_due)`. Input: Boolean `done` and an integer day count or missing date. Output: one of `REMIND`, `NOT_YET`, `OVERDUE`, `NO_DATE`, `DONE`. Order: if `done=true`, return `DONE`; else if date missing, `NO_DATE`; else if days below 0, `OVERDUE`; else if no more than 2, `REMIND`; otherwise `NOT_YET`. `(false,2)` gives `REMIND`, `(false,3)` gives `NOT_YET`, `(true,2)` gives `DONE`. Calling it does not change the file, so the screen can ask repeatedly.

## Check it by hand

Let a partner act as a function: hand over `done` and `days_to_due` cards, and ask for exactly one output card. Start with `done=true` and no date: return `DONE`, because completion is checked before the date. Then try `done=false` and no date: return `NO_DATE`. Draw a return arrow for each case. If the partner returns two answers or none, clarify the function contract. The function must not modify the source task merely to produce an answer.

## Predict and check

What does `(false,-1)` return? `OVERDUE`; `(false,None)`? `NO_DATE`; `(true,None)`? Under this order, `DONE`. Someone could choose `NO_DATE` for a done task instead, but that is a different contract that must be written and applied consistently.

## Catch the error

“If the date is absent, silently treat it as zero days.” Zero means today; absence means unknown. Mixing them creates false reminders. Repair the input type or add a separate branch.

## Transfer to a new setting

Specify `count_done(tasks)` with task-list input and integer output. Its precondition requires Boolean `done` for every record. Its postcondition says the result is between 0 and list length and the list remains unchanged.

## Project change

Add the `should_remind` contract to `ALGORITHM-en.md`: inputs, five outputs, missing-date rule, and “does not change file” guarantee. Link it to lesson 22's boundary cases.

## Task and evidence

Give another learner only the contract. They must predict five results without your diagram. Then swap two checks and find an input where the answer changes. Evidence: the contract and a counterexample.

## Return after 1, 7, and 30 days

Tomorrow distinguish parameter from argument. In seven days write a contract for another task. In a month compare Python behaviour with the promised cases.

[Next lesson: A correct and sufficiently fast algorithm](/read/informatics-25-correctness-efficiency?lang=en)
