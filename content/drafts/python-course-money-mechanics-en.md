# Free Python course for beginners: 57 lessons and a real-data project

_In brief:_ **Shanraq’s open course teaches Python through one continuous project. Learners collect data from official sources, validate and analyse it, then turn their work into a reproducible report. Here is who the 57 lessons are for, how the learning path works, and why money is a subject to investigate rather than a conclusion supplied in advance.**

Python is often introduced through disconnected exercises: print a line, calculate a number, write a loop. Each topic may make sense on its own, yet a beginner can still struggle to see how the pieces become a real program. The Shanraq course follows a different route: each new skill becomes part of a working research project.

Learners move from setting up an environment and understanding basic syntax to requesting data over HTTP, reading JSON and XML, storing records in SQLite, analysing tables with pandas, creating charts, testing statistical claims, writing automated tests, and publishing a result. Every topic therefore answers two questions: “How do I write this in Python?” and “Why does the project need it?”

![Support map: the data path from source to conclusion](/static/course/py/map-pipeline-en.svg)

## An investigation with a verifiable result

The continuing subject is prices, money, and the operation of a modern monetary system. Learners work with open data from the National Bank of Kazakhstan and the World Bank. They practise distinguishing a source observation from a calculated measure and checking whether a conclusion can withstand comparison with the data.

The lessons use MMM as an abbreviation for *Modern Money Mechanics*. Here it refers to the mechanisms through which money is created and moves through a modern banking system. It does not refer to a pyramid scheme, nor does it imply an accusation decided in advance. The course asks learners to locate an official source, preserve data provenance, state a hypothesis, and disclose the limits of an analysis instead of accepting a ready-made narrative.

![Support map: money, measures, and causal claims](/static/course/py/map-money-en.svg)

This subject gives programming a concrete purpose. A variable becomes a measurement, a dictionary becomes a source record, a table becomes a set of observations, and a test becomes protection against a silent error in the conclusion. The initial hypothesis may turn out to be wrong. In a research project, that is not failure; it is a valid and useful result.

## Who the course is for

For a school student who already knows some computing, the course provides a bridge from classroom exercises to a project they can demonstrate and explain. University students can connect Python with data processing, databases, HTTP, statistics, and reproducibility. Adult learners can work through small units at their own pace instead of following a group timetable.

Previous programming experience helps but is not required. The opening lessons cover Python installation, virtual environments, values and types, conditions, loops, functions, and collections. Later lessons become substantially more demanding, so the intended route is sequential: run each example, answer the recall questions in your own words, and only then compare your attempt with the reference solution.

No account is required to read the lessons or complete the exercises. An account is useful for learners who want to discuss the material and ask questions in the comments. The course itself remains open and free of charge.

## How the learning path works

The lessons move from the whole to the details. Learners first see a working fragment and its expected result. They then examine its parts, predict what the program will do, complete the code, and repair an error. At the end of the lesson, the new component returns to the continuing project.

The visual support maps draw on the idea of support signals associated with Viktor Shatalov’s teaching method. A large topic is compressed into a small visual structure that learners can revisit during recall. A map does not replace explanation or practice; it helps reconstruct the sequence of actions and the relationships between concepts without rereading the entire lesson.

![Support map: checking sources, types, and assumptions](/static/course/py/map-trust-en.svg)

The same active learning loop appears at every important stage:

1. see the complete outcome;
2. predict what the code will do;
3. run it and compare expectation with evidence;
4. modify or complete the program;
5. explain the decision in your own words;
6. integrate the result and protect it with a test.

This keeps the learner from becoming a passive reader. An error becomes an observable clue rather than a reason to abandon the topic.

## What learners should be able to do at the end

The goal is not to memorise every Python function. Learners should be able to break a problem into steps: obtain data, validate the server response and schema, clean values, store them, calculate a result, separate evidence from assumption, test critical rules, and package the outcome so another person can reproduce it.

Dedicated lessons cover type hints, logging, scheduled execution, protection against concurrent runs, the limits of statistical evidence, and careful use of a language model. AI is treated as an assistant whose output needs schemas, source checks, and tests—not as a substitute for understanding.

![Support map: reproducible output and publication](/static/course/py/map-publish-en.svg)

## A practical way to study

One short lesson per session, followed by practice in a personal copy of the project, is a sustainable pace. Avoid copying the reference answer before making a first attempt. When the code works, change the input and test an edge case: an empty response, a wrong type, a missing date, or a duplicate key.

Start with the [Python course map](/course/python?lang=en). If an explanation makes an abrupt leap, an example fails, or a support map is ambiguous, leave a comment. That feedback can make the path clearer for the next learner.

## Sources and course materials

- [Python course on Shanraq](/course/python?lang=en)
- [Open Data Repository of the National Bank of Kazakhstan](https://data.nationalbank.kz/)
- [World Bank Open Data](https://data.worldbank.org/)
- [Official Python documentation](https://docs.python.org/3/)
