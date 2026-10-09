# Decompose a large problem without losing the goal

_Lead (summary):_ **Break a reminder request into small testable actions and build the assistant's first paper algorithm.**

## Where we are on the map

This is the third block of one project. We use version 0.2 data but describe actions on paper first, where mistakes are easier to find before writing Python.

![Decompose a large problem without losing the goal](/static/course/informatics/map-19-problem-decomposition-en.svg)

## Situation and question

A parent says, “Make the assistant remind me about the school project.” Drawing a button immediately is tempting, but where is the due date stored, and what happens when it is missing? First state what result the person should see and which data that result requires.

## New words without gaps

A **problem** names the required result, not a list of buttons. **Decomposition** splits the problem into parts that can each be explained and tested. A **subtask** has input and output: “read the due date” returns a date or a message that none is set. A **dependency** means one step uses another step's result. A **success criterion** is observable: with a date, the assistant gives the correct reminder; without one, it honestly says that no date is set. An interface only displays the result; it is not the logic itself.

## The lesson's support signal

`need → input data → small steps → testable result`

## Work through it step by step

Start with fictional task `t-01`, “Кітап оқу.” A separate test card says it is due tomorrow; the older `tasks.json` has no due-date field yet. (1) Find the record by `id`, not title: titles may repeat. (2) Read its due date. (3) If none exists, return “no date set” and stop. (4) Compare the date with today. (5) If the task is unfinished and due soon, prepare a reminder. (6) Display the result. Draw six cards and dependency arrows. Explain why comparing dates cannot precede reading the date. This paper model is not yet code: it lacks a date format and a rule for “soon.” Later lessons will supply them.

## Check it by hand

Take four cards: `id`, `done`, due date, and output. Keep the task title on a separate card. Let a partner act as the assistant: hand over cards one by one, and ask which card is needed next. Remove the date card: the partner must stop rather than guess. Change the title while keeping the same `id`: the result should not change. This reveals which information actually affects the decision. Write the rule down before drawing a screen or writing code.

## Predict and check

A record has a title but no due date. Can date comparison run immediately? No: there is nothing to compare; end that branch with an understandable message. A record has a date but is already done: should it trigger a reminder? No, if the requirement is to remind only about unfinished tasks. Find the shared step in both cases: locating the record.

## Catch the error

“Let's make a beautiful screen first and decide on data later.” Repair the plan: the input and date check determine what the screen may honestly promise. An absent date must not look like a zero date.

## Transfer to a new setting

Decompose another request: “Show how many tasks are done.” Which steps are shared with reminding and which are new? Do not read a due date without a reason. Counting must examine every record, not only the first.

## Project change

In `step-02/ALGORITHM-en.md`, document the goal, input, output, and six reminder cards. Keep the three original fictional tasks unchanged. Add a separate assumptions table for today's date, date format, and the meaning of “soon.” Do not hide unchecked assumptions inside a diagram.

## Task and evidence

Give the six shuffled cards to a partner. They must restore dependencies, show the no-date branch, and explain the stop. Evidence: a diagram, two runs with different inputs, and one corrected assumption. If two incompatible orders both seem valid, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow reconstruct the six steps from memory. In seven days decompose a counting request unaided. In a month check that future code does not rely on an unstated date format.

[Next lesson: State and variables through a game](/read/informatics-20-state-variables?lang=en)
