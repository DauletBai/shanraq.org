# Release 1.1: a web version and an honest cloud model

_Lead (summary):_ **Release a local web view, verify the old rules, and state honestly what our cloud model does not yet do.**

## Where we are on the map

This is lesson 46 of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.

![Release 1.1: a web version and an honest cloud model](/static/course/informatics/map-46-web-cloud-release-en.svg)

## Situation and question

A classmate sees the browser page and says, 'So the project is already on the internet.' Test that claim. The program listens on 127.0.0.1, which refers to the same computer; a remote family cannot access it through the public internet. We can still test a real browser view and identify what the next step would require.

## New words without gaps

**Deployment** means running a version where its intended audience can reach it. A **cloud** consists of rented or managed computing resources; it is not a special physical kind of network. **Loopback** routes a computer to itself through 127.0.0.1. A **regression** breaks a rule that worked before. A **checkpoint** preserves code, data, instructions and verifiable outcomes. A **trust boundary** separates this local classroom process from a public site that handles accounts and private data.

## The lesson's support signal

Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.

## Work through it step by step

From `step-04`, run `python3 web_assistant.py --port 8765` and open `http://127.0.0.1:8765/?lang=en&today=2026-10-09`. See the three original tasks and `REMIND`, `DONE`, `NO_DATE`, as in release 1.0. Switch to `ru` and then `kz`: labels change, calculation does not. Enter 2026-02-29 and confirm a 400 response. Try POST and get 405. Compare the checksum of `tasks.json` before and after: the file does not change. Then run `python3 -m unittest discover -s tests`; this checks repeatable behaviour rather than judging a screen by eye.

## Predict and check

Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.

```python
from datetime import date
from assistant_core import load_document, reminder_status
tasks = load_document("tasks.json")["tasks"]
for task in tasks:
    print(task["id"], reminder_status(task, date(2026, 10, 9)))
```

## Expected output

```text
t-01 REMIND
t-02 DONE
t-03 NO_DATE
```

## Catch the error

Do not expose this `http.server` publicly. Python's own documentation says it has only basic security checks and is unsuitable for production. Local HTTP has no TLS, accounts, permission system, or protection for real user records. Thus 'web version' means a browser view on one computer, and 'cloud model' means understanding a future architecture, not claiming a completed deployment.

## Project change

Keep `step-04` as a checkpoint: an exact copy of the 1.0 core, the three fictional records, a new read-only web server, tests, and instructions in three languages. Do not change the `tasks.json` format or the five reminder outcomes. Record future work in the decision log: HTTPS, users, permissions, protected storage and a production-grade server belong to later blocks.

## Task and evidence

For mastery, answer at least 8 of 10 questions about addresses, DNS, delivery, TLS, HTTP, HTML and accessibility; demonstrate three languages and three statuses. Repeat at least 7 answers a week later. Give the instructions to someone who did not write the code: they should open the page, try an impossible date, and explain why the file stayed unchanged.

## Ten questions for self-check

1. Who starts a request and who replies when the server runs on your own computer?
2. What does 127.0.0.1 mean, and why can the page work without internet?
3. Why does one IP address not prove there is one human behind it?
4. What does DNS return, and why is that not enough to trust a site?
5. What do TCP and the simple UDP model do if part 2 is lost?
6. What does TLS check, and what does a padlock not promise?
7. When does our project return 200, 400, 404, and 405?
8. Why does GET not mark a task complete in our release?
9. Why do label/input pairing and task-title escaping matter?
10. What three changes are needed before families can use this through the internet?

## Transfer to a new setting

A school asks to show tasks to every family. What changes in addressing, access control, data protection and server environment? Name at least three requirements before you even consider moving this classroom program off localhost.

## Return after 1, 7, and 30 days

After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.

## Primary reference to check

[Official documentation](https://docs.python.org/3/library/http.server.html)

## All course lessons

[All lessons](/course/informatics?lang=en)
