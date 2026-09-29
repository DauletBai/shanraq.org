# Diagnosis without grades and choosing the project problem

_Lead (summary):_ **Ten short situations will show where you need more support, while a completed passport will turn a broad assistant idea into a testable problem.**

## Where we are on the map

This is the second lesson of the first block. Everyone continues to computer systems afterwards, but learners receive different amounts of support. Diagnosis neither closes the path nor divides people into “technical” and “humanities” types.

![Four groups of diagnostic tasks leading to the project passport](/static/course/informatics/map-02-diagnostic-en.svg)

## A familiar situation: configuring a route

A navigation app needs more than the name of a city. It needs a starting point, a means of travel, and constraints. It offers different routes to a pedestrian and a bus. Diagnosis has a similar job: it identifies the starting point and the needed support, rather than a learner's worth.

The analogy has a limit. Navigation software can receive coordinates automatically, while understanding cannot be measured with one button. We need the learner's answer, explanation, and confidence.

## How to mark answers

Prepare a sheet of paper. Record a solution and one mark for each situation:

- `✓` — I can explain the reason to another person;
- `?` — I have an answer, but my explanation is not secure yet;
- `not yet` — I do not know which action to begin with.

A correct guess marked `?` is more useful than false confidence. Do not search for definitions before your first attempt. We are locating your initial support, not measuring search speed.

## Ten situations

### 1. What survives shutdown

You typed a note but have not deliberately saved it. The device then powered off completely. Can we know for certain that the note remains, disappears, or is there insufficient information? Name the condition on which the answer depends.

### 2. Where the file is

Two files are named `plan.txt`. One is in the `school` folder and the other in `sport`. How can you identify each file unambiguously? Write both addresses in any consistent form.

### 3. Which image contains more data points

Two uncompressed black-and-white images measure `10 × 10` and `100 × 100` dots. How many times more dots does the second contain? Explain the calculation.

### 4. Instructions for a very literal worker

Write steps for adding the task “Bring the book” to the assistant. The worker performs only written steps and does not guess missing ones. Include what happens if the task already exists.

### 5. Trace state

A game starts with `3` lives. The first event removes one life, the second adds two, and the third removes three. How many lives remain after each event? Write the entire sequence of states.

### 6. Check the boundary

A rule says, “A reminder is urgent when **no more than** two days remain.” What result does it give for 3, 2, 1, and 0 days? Explain “no more than”.

### 7. An error message

A program reports: `File tasks.json was not found`. Give two possible explanations and the first safe checking action. Do not immediately propose reinstalling the entire application.

### 8. An urgent link

A message says, “Your school account will be deleted in 10 minutes. Sign in at `shamraq-login.example`.” Name at least three signs to examine before following the link and a safe way to open the real service.

### 9. An honest claim from data

The assistant contained 2 tasks on Monday and 7 on Tuesday. A learner says, “My workload grew 3.5 times.” Are those counts enough to support a claim about real workload? Name at least two missing task properties.

### 10. An AI answer

AI claims, “A strong password must replace `a` with `@`.” It gives no source. What should happen before this advice becomes a project rule? Name at least three actions.

## Discussion and guide answers

1. There is insufficient information. Some applications save drafts automatically; others keep them only in working memory. We need to know whether the note reached persistent storage.
2. We need paths such as `school/plan.txt` and `sport/plan.txt`. A filename without its folder is ambiguous.
3. `10 × 10 = 100` and `100 × 100 = 10,000`; the second has 100 times, rather than 10 times, as many dots.
4. The answer should include input, an ordered sequence, a result, and a separate branch for a duplicate. There is no single required wording.
5. `3 → 2 → 4 → 1`. Intermediate states matter.
6. Three days is not urgent; 2, 1, and 0 are urgent. The boundary value 2 is included.
7. The file may be absent, or the program may be looking in the wrong folder. First check the exact path and whether the file exists, without deleting data.
8. Time pressure, a look-alike domain, and a request to sign in through a supplied link are warning signs. Open the real address from your own bookmark or type it yourself.
9. The counts are insufficient. We need at least duration, difficulty, deadline, or completion state. A count of records is not the same as workload.
10. Separate the claim, find independent reliable sources, check their date, compare the advice with the threat model, and test it on examples. Character substitution alone does not guarantee strength.

Matching a guide answer does not yet prove understanding. Keep `?` if you cannot explain the reason or transfer the method to new data.

## How to use the result

Group your marks:

- tasks 1–3 — devices and data representation;
- 4–6 — algorithmic thinking;
- 7–8 — error diagnosis and security;
- 9–10 — data and AI.

A group containing `?` or `not yet` does not send you backwards. It tells you to perform every experiment and use gradually fading support in the relevant lessons. In a secure group, you may move through a familiar example more quickly, but you still complete transfer and the project change.

## The lesson's support signal

```text
situation → my answer → reason → confidence mark
                                  ↓
                    required amount of support
                                  ↓
                         a new check later
```

Diagnosis ends with a decision about the next action, rather than a score.

## Choosing the cross-course project theme

Choose one foundation:

1. study tasks and preparation;
2. sports training;
3. books, films, or games;
4. a club calendar;
5. an environmental or volunteer project.

The theme changes words and visual treatment, but not the course knowledge. Every project has records, dates, states, search, a database, protection, and a verifiable AI feature. A teacher can discuss one shared architecture while each learner creates a personally meaningful product.

## Worked passport example

```text
Name: Book Club Assistant
User: a member of a school club
Problem: a member forgets which book the club will discuss next
Three actions: add a meeting; see the nearest one; mark preparation
Does not store: real names, phone numbers, or private messages
Evidence of value: a new member finds the topic and date within 20 seconds
```

The user is a role, the problem is observable, the actions are testable, and prohibited data is named in advance.

## Example with a fading prompt

The idea is “an app for sport”. Make it precise in four lines:

- who uses it;
- what difficulty occurs;
- which three actions are needed;
- how value will be measured.

“The user will like it” is not yet an observable check. Replace it with completion time, the number of correctly found records, or a successful scenario.

## Recall without a prompt

Hide the page and explain:

1. why diagnosis is not a grade;
2. why an answer needs a reason and a confidence mark;
3. why the project does not need a user's real name;
4. how a testable problem differs from the name of a technology.

## Find and correct the mistake

Consider this statement: “Create a modern Python application with a database and artificial intelligence.”

It lists technologies but lacks a user, difficulty, and success condition. Rewrite it as one sentence:

> For [role] facing [observable difficulty], the product must produce [observable outcome], verified by [method].

## Transfer to a new setting

Imagine a school library without a computer. Can the same problem-framing method apply? Yes. Digital technology is not a required part of the purpose. We may improve a card, an ordering rule, or a process first and only then decide whether a computer is needed.

Name one problem better solved by changing a process than by creating another application.

## Project change: complete version 0.0

Complete the passport after diagnosis:

```text
Theme:
User role:
One observable problem:
Three actions in the first release:
Data in the first release:
Data we deliberately do not store:
How another person will verify value:
The group in which I need full support:
```

Use fictional records only. Create three examples with which we will later test algorithms.

## Exercise

**Required.** Complete all ten situations, retain your marks, and finish the passport. For one task marked `?`, state what observation would turn uncertainty into an explanation.

**With your own data.** Create three safe fictional records for your theme. Each needs a title, a deadline or date, a state, and a category.

**Optional.** Give another person only the problem statement and ask what they expect the product to do. Record the difference without persuading them to accept your version.

## Retrieval after 1, 7, and 30 days

- Tomorrow, repeat one task from your least secure group with different numbers or conditions.
- After seven days, complete tasks 3, 6, 8, and 10 without the old answers.
- After thirty days, compare the `0.0` passport with the working release and explain every change in purpose.

[Next lesson: phones and laptops as systems](/read/informatics-03-device-system?lang=en)

