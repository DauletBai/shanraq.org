# Versions, collaboration, and accessibility

_Lead (summary):_ **We preserve meaningful history, combine changes without overwriting another person's work, and test accessibility.**

## Where we are on the map

Two people edited the same task at once. Which version should stay, and how do you see the difference? What if one person cannot distinguish the screen's colours?

![Versions, collaboration, and accessibility](/static/course/informatics/map-09-versions-en.svg)

## A familiar image and the precise model

Files named `final2-new` do not explain ancestry. A version is a defined product state. A small change has a reason and check. Comparison shows differences, revert restores known state, and merge combines independent edits. Accessibility requires more than colour, keyboard reachability, and clear errors.

## Make changes visible to everyone

A **version** is a saved state of the project at a particular time. A **diff** shows added, removed, and changed lines. A **conflict** occurs when two people change the same part differently and no automatic choice is justified. **Accessibility** means a person can complete the task with different ways of perceiving or controlling the interface; colour must not be the only signal. **Keyboard focus** is the place where the next key press will act without a mouse.

Copy `tasks.json` to `version-A.json` and `version-B.json`. In A, rename `t-01`; in B, change the `done` value of `t-01`. On paper, produce a final record containing both compatible edits. Now change the title differently in A and B: meaning must be discussed, not settled by “last save wins.” Record your decision and check that the final JSON opens correctly.

Show task state with words such as “done” or “not done” as well as colour. Move through the project with Tab: can you see which item has focus, open the file, and understand an error? Give the project to another learner and ask them to repeat the action without prompts. That tests collaboration and accessibility together.

## The limit of the analogy

The familiar image opens the topic but does not replace the mechanism. State where the analogy stops being accurate.

## The lesson's support signal

`small change → compare → review → test → version → recover`

## Worked example

The change “reject an empty title” contains one rule, an error example, and a check. It is not mixed with colour or folder renaming, so it is easy to verify and reverse.

## Predict before observing

Before the experiment, hide the continuation and predict the result in writing. Give the chain of causes from the support signal as well as the answer. Perform the action, record the observation, and compare it with the prediction. A match without an explanation is luck rather than understanding; a mismatch is useful evidence about the link that needs rebuilding.

## Recall without a prompt

Hide the page, rebuild the signal, and explain every transition with a new example.

Explain the topic to someone who has not read the lesson. Do not use a new term until you have unpacked it in plain words. Then restore its precise name and show its place in the system. The listener should be able to predict the next step and give a different example. If they merely repeat a sentence, reduce the explanation to the support signal and rebuild it.

## Find and correct the mistake

“We will add accessibility at the end.” Late structural repair costs more. Test contrast, headings, keyboard, focus, and messages now.

## Transfer to a new setting

Apply the model to a device or file absent from the lesson. Separate observation from assumption.

## Project change

Start a decision log and snapshot `v0.1-foundations`. Test the passport without colour and with keyboard only. Record one barrier and repair.

Keep four lines in the decision log: the intended result, the observation, the reason for the chosen action, and the verification. A screenshot can support evidence but cannot replace text and the working file. Do not use personal data: three fictional records provide the same testability and let you show the project safely to another person.

## Exercise

**Required.** Complete the example, correct the error, and make the project change with a reason.

**With your own data.** Test the decision on three fictional records.

**Optional.** Ask another person to rebuild the explanation from the map.

The readiness criterion is concrete: reproduce the signal without the page, correct the proposed error, transfer the model to an unfamiliar example, and show the project change. If one part fails, return to that link instead of rereading the whole lesson. After repairing it, repeat only an equivalent task with different data.

## Retrieval after 1, 7, and 30 days

Rebuild the signal tomorrow; solve an equivalent task after seven days; after thirty days, verify the decision in the project.

[Next lesson: the first-block mastery checkpoint](/read/informatics-10-systems-mastery?lang=en)
