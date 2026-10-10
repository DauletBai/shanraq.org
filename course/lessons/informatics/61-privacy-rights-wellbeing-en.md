# Privacy, authorship, and digital wellbeing

_Lead (summary):_ **A study report needs minutes and dates.**

## Where we are on the map

Lesson 61 of 72. Digital security block (55–62). We continue a local assistant with three fictional tasks. Each control has a defined boundary: the classroom version does not become a public service.

![Privacy, authorship, and digital wellbeing](/static/course/informatics/map-61-privacy-rights-wellbeing-en.svg)

## Situation and question

A study report needs minutes and dates. A person’s name, photograph and exact address do not answer “how long did the task take?” Less collected data means less harm if data escape.

## New words without gaps

**Data minimisation** means collecting only what a clear purpose needs. **Consent** must be informed and fit the context; children require additional safeguards and adult support. A **licence** defines permitted reuse of someone else’s material. **Digital wellbeing** includes breaks, sleep and the right to be offline. Our assistant data are fictional.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Inspect `study_sessions.csv`: session ID, task ID, date, minutes and source. That is enough for our classroom report. Propose a new “home address” column and show why it cannot improve the 45/0/25 calculation. Then write a checklist: do not paste other people’s private messages into an AI service, check an image licence, take breaks and tell a trusted adult about unwanted contact.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
import csv
with open("study_sessions.csv", encoding="utf-8") as file:
    print(next(csv.reader(file)))
```

## Expected output

```text
['session_id', 'task_id', 'observed_on', 'minutes', 'source']
```

## Catch the error

Public access to a file is not permission to republish it. An apparently anonymous ID may be linked to a person through other data; do not promise absolute anonymity. Deleting a name from the database does not automatically erase it from backups.

## Project change

Release 1.3 preserves the three tasks, four sessions and reminder rules. New operations live in `security_assistant.py`: validation, a separate backup and restore only into a new file. State which threat each change reduces and which it does not; all teaching data stay fictional.

## Task and evidence

Remove three unnecessary fields from a fictional form while preserving the ability to produce the report. Show where you would record the retention purpose and who could request deletion. Name one limit of your design.

## Transfer to a new setting

A club wants to publish children’s photographs and attendance times. What questions must it ask first?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.unicef.org/parenting/child-care/online-privacy)

## Next lesson

[Release 1.3: a protected and recoverable assistant](/read/informatics-62-security-release?lang=en)
