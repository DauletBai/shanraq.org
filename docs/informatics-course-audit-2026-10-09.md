# Informatics course review, 9 October 2026

## Scope and findings

The planned route has 72 lessons in eight dependency-linked blocks. The 18
previously released lessons were reviewed in all three languages. Lessons 1–4
already introduce the project and physical devices with extended explanations;
lessons 5–18 had a repeatable short template. Their chief weakness was the jump
from a named concept to an assignment without a small, observable experiment.
The recurring introductory paragraph also hid the specific question each lesson
was supposed to answer. Terminology was present but some relationships, such as
memory versus storage, a file versus its encoding, and sampling versus
compression, needed explicit boundaries.

The 5–18 manuscripts now open with topic-specific questions and add a concrete
bridge after the first model. Each bridge gives a learner an action, a visible
result, and an explanation connecting it to the digital assistant. These changes
apply separately to Russian, Kazakh, and English text. The existing project
examples, values, headings, SVGs, and revision cycle are preserved.

## New block 19–26

Each lesson now has a real-life request, terms defined in place, a worked trace,
a hands-on check, a prediction before revealing the answer, a fault to diagnose,
a transfer task, and a project change. The learner develops paper version 0.3:
reminders, completed-task counts, and duplicate-ID detection. The 0.2 data has
no due date, so the lessons explicitly provide due dates as **test cards**. No
lesson pretends the prior JSON already stores them. `cases.json` records the
expected outcomes, including boundaries, missing dates, empty lists, and
completed tasks; the test suite checks these outcomes against the contract.

Lessons 19–26 each have Russian, Kazakh, and English texts and maps. The maps
show four exact stages and the relevant boundary or outcome; the 4K cover is
kept with the informatics algorithm block. The checkpoint requires an
independent learner to reproduce the answers from the written contract. The
next Python block must implement the same contract before adding UI or storage.

## Teaching and editorial standard

The design uses short spaced reviews, worked examples followed by independent
practice, concrete diagrams alongside verbal explanation, and a question that
requires the learner to explain the result. These choices follow the US Institute
of Education Sciences practice guide *Organizing Instruction and Study to
Improve Student Learning*:
https://ies.ed.gov/ncee/wwc/PracticeGuide/1 . The guide supports these
techniques; it does not claim one fixed success rate for this course.

Before later blocks are released, review every new term at first use; verify
the diagram against its worked example; run the same test inputs in all three
languages; and have another reader execute the project instructions without the
author present. Code execution and actual notifications belong to the Python
block, not to this paper release.
