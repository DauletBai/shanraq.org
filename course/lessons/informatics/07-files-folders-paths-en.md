# Files, folders, and exact paths

_Lead (summary):_ **We build a project tree and distinguish a name from a path, a relative path from an absolute one, and a missing file from a wrong search location.**

## Where we are on the map

We build a project tree and distinguish a name from a path, a relative path from an absolute one, and a missing file from a wrong search location.

Start with an observation rather than a definition. Open an object on your device that relates to the lesson and record only what you can see: its name, location, action, and result. Write your explanation of the cause separately. This prevents a guess from becoming a fact. After the experiment, compare the explanation with the precise model and repair only the link that was wrong.

![Files, folders, and exact paths](/static/course/informatics/map-07-paths-en.svg)

## A familiar image and the precise model

Two files may be named `plan.txt`. A file is a named data sequence; a folder maps names; a path lists transitions. An absolute path begins at a root and a relative path at the working folder. Renaming `.txt` to `.jpg` does not turn text into a photo.

## The limit of the analogy

The familiar image opens the topic but does not replace the mechanism. State where the analogy stops being accurate.

## The lesson's support signal

`assistant/ → data/tasks.json | docs/passport.md | tests/`

## Worked example

From `assistant`, `data/tasks.json` is unambiguous. `../tasks.json` first moves to the parent and may identify another file. Inspect the working folder and resolved path before fixing.

## Predict before observing

Before the experiment, hide the continuation and predict the result in writing. Give the chain of causes from the support signal as well as the answer. Perform the action, record the observation, and compare it with the prediction. A match without an explanation is luck rather than understanding; a mismatch is useful evidence about the link that needs rebuilding.

## Recall without a prompt

Hide the page, rebuild the signal, and explain every transition with a new example.

Explain the topic to someone who has not read the lesson. Do not use a new term until you have unpacked it in plain words. Then restore its precise name and show its place in the system. The listener should be able to predict the next step and give a different example. If they merely repeat a sentence, reduce the explanation to the support signal and rebuild it.

## Find and correct the mistake

“File not found means no file exists.” The program may look elsewhere or the name may differ by one character. Separate observation from causes.

## Transfer to a new setting

Apply the model to a device or file absent from the lesson. Separate observation from assumption.

## Project change

Create `assistant/data`, `assistant/docs`, and `assistant/tests`. Keep the passport in `docs`, three fictional records in `data`, and record the working folder.

Keep four lines in the decision log: the intended result, the observation, the reason for the chosen action, and the verification. A screenshot can support evidence but cannot replace text and the working file. Do not use personal data: three fictional records provide the same testability and let you show the project safely to another person.

## Exercise

**Required.** Complete the example, correct the error, and make the project change with a reason.

**With your own data.** Test the decision on three fictional records.

**Optional.** Ask another person to rebuild the explanation from the map.

The readiness criterion is concrete: reproduce the signal without the page, correct the proposed error, transfer the model to an unfamiliar example, and show the project change. If one part fails, return to that link instead of rereading the whole lesson. After repairing it, repeat only an equivalent task with different data.

## Retrieval after 1, 7, and 30 days

Rebuild the signal tomorrow; solve an equivalent task after seven days; after thirty days, verify the decision in the project.

[Next lesson: formats, software, and the creator's permission](/course/informatics?lang=en)
