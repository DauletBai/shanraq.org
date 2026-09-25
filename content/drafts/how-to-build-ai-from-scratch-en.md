# How to build your own AI from scratch: a free beginner course

_In brief:_ **Build a local AI model for Kazakh text in Shanraq’s free 30-lesson Python course. Learn to parse questions, find verified facts, check their dates, and show the source behind each answer. When the evidence is missing or outdated, the model explains why it cannot answer. Every example runs on an ordinary computer without a paid API.**

A search for “how to build your own AI” often leads to one of two extremes: instructions for connecting a ready-made chatbot or an advanced neural-network course. The central learning question gets lost between them: **what decisions make up an AI model, and how can each decision be tested?**

The new Shanraq course starts with that question. Learners do not rent somebody else’s model or hide the whole system behind one call to an **API**, an agreed way for programs to exchange data. They build a small local model for Kazakh text step by step: analyze words, identify what a question asks for, locate an approved fact, check its date and source, and then either answer or explain why an answer is not allowed.

> A useful learning model should reveal both its answer and the path that produced it.

[Open the free “AI without an LLM” course](/course/kazakh-ai?lang=en)

## Can you build AI without a neural network or an LLM?

Yes. Artificial intelligence is a broad name for systems that recognize input, select an action, or make an inference. A **neural network** stores learned patterns in many numerical parameters. A large language model, or LLM, is a large neural network for text and one way to build an AI system.

This course combines readable rules with a small classifier. A **classifier** is the part of a program that assigns a question to one of several known types. A question may ask for the time of a club, its location, or something outside the system’s scope. The model counts word features, but it earns permission to answer only after validating the source record.

The process resembles a careful librarian. The librarian understands which card a visitor needs, checks the date, and points to the source. If the card is missing, the librarian does not invent a shelf or opening time.

![Support map for the lesson on refusing when evidence is missing](/static/course/kazakh-ai/map-15-abstain-en.svg)

An LLM works differently. It can continue free-form text and handle a vast range of topics, which makes it good at explanation, translation, and conversation. A plausible sentence, however, is not a guarantee of a correct fact. A narrow model knows far less, but its answering boundary can be written down and tested.

## What learners build in 30 lessons

The final project accepts a date and a question in Kazakh. If someone asks when the chess club meets, the program must pass through several gates:

1. normalize the text;
2. recognize known forms of a word;
3. decide whether the question asks for a time or a place;
4. reject a question that combines two different requests;
5. find exactly one fact card;
6. check the record’s validity period;
7. return the answer with its date and source.

If any condition fails, the program does not cover the gap with a confident sentence. It reports that the question is ambiguous, the fact is absent, or the record has expired.

> An intelligent system begins with the ability to say, “I do not know, and I can explain why.”

By the end, learners have a runnable project with data files, quality checks, a typed launch command, and a repeatable performance measurement rather than a slide presentation.

## Why the project uses the Kazakh language

Kazakh makes it unusually easy to see how a language rule can become a sequence of program steps. It has strongly agglutinative morphology: recognizable parts are attached to a base in order, with each part contributing meaning.

A useful everyday image is a skewer. The base holds the central meaning, while suffixes are added in an ordered sequence. Sounds can alternate, ambiguity exists, and some meanings use separate words, so software must still test its assumptions. The visible sequence nevertheless makes Kazakh an effective language for a first morphological analyzer. **Morphological analysis** means dividing a word into meaningful parts.

![Support map for the main language structure types](/static/course/kazakh-ai/map-03-language-types-en.svg)

The course does not try to reduce a living language to a formula. It teaches a more practical lesson: some regularities can be encoded, while exceptions and boundaries must remain visible. This makes Kazakh a strong setting for learning algorithmic thinking.

## How the beginner AI course is organized

All 30 lessons contribute to one project, with difficulty increasing gradually.

| Stage | Work completed |
|---|---|
| Lessons 1–5 | project goal, language types, word parts, and a first manual analysis |
| Lessons 6–10 | text normalization and finding words, bases, and suffixes |
| Lessons 11–15 | case forms, ambiguity, entities, question intent, and refusal |
| Lessons 16–20 | approved facts, provenance, answer templates, and expiry |
| Lessons 21–25 | learning examples, data splits, a classifier, and quality measurement |
| Lessons 26–30 | model boundaries, comparison, audit, benchmarking, and the final command |

A new term appears only after its purpose has been explained. Learners first meet an everyday image and a finished result, then the term, code, and an independent task. Support maps help reconstruct the solution path without rereading the entire lesson.

Every fifth lesson is a mastery checkpoint. Learners reproduce the essential path from memory, check their result, and return to any unclear step. This structure draws on Viktor Shatalov’s idea of reference signals while keeping hands-on work central.

## How fast is the finished model?

The final lesson includes `benchmark.py`, allowing learners to repeat every measurement on their own computer. A **worker process** here means one running copy of Python that handles requests in sequence. A **logical core** is one execution path available to the operating system. One worker in this program uses no more than one logical core at a time, while several workers can run on different cores.

The recorded control run used an eight-logical-core ARM computer with macOS 26.5.2 and Python 3.14.5. ARM names the processor family; the operating system and Python version identify the rest of the test environment.

| Worker processes | Total throughput | Speedup | Summed process memory |
|---:|---:|---:|---:|
| 1 | 42,953 requests/s | 1.00× | 25.4 MiB |
| 2 | 80,372 requests/s | 1.87× | 50.8 MiB |
| 4 | 146,241 requests/s | 3.41× | 101.7 MiB |
| 8 | 141,357 requests/s | 3.29× | 204.0 MiB |

Four workers were fastest in this experiment. Eight did not increase throughput and used almost twice the memory of four.

The learning model’s code and data occupy 6.7 KiB. One KiB is 1,024 bytes, and one MiB is 1,024 KiB. A complete running process including Python used 24.4 MiB. The system produced all 16 expected outcomes on the published quality cases, including ten correct refusals.

That result is not a promise of perfect behavior on arbitrary text. The suite is small and public, and an incorrect source record remains incorrect. The figures teach learners to separate a measurement from a marketing claim. The complete method and large-model comparison are available in the [final course lesson](https://shanraq.org/read/kazakh-ai-30-cli?lang=en).

## Where this model beats an LLM, and where it falls short

A narrow model works well when the task is bounded in advance: a timetable, directory, industrial instruction, validated form, or local service that must not send data to an external API. Its rules can be read, a fact can be corrected in one record, and its refusals can be tested.

A large language model is far stronger when users write freely, topics cannot be listed ahead of time, and the system must produce coherent new text. It can explain and generalize in ways this learning classifier cannot.

Sometimes the right design is a mixed system. An LLM interprets a free-form request, while a rule-based component validates allowed fields, dates, and sources. A program is called **deterministic** when the same input and state follow the same path and produce the same result.

> Mature engineering begins with a precise task and an acceptable error, not with the largest available model.

## What you need before starting

The course is designed for people new to artificial intelligence, but it assumes basic Python: variables, conditions, loops, functions, lists, dictionaries, and reading JSON. If these topics are unfamiliar, begin with the early part of the [free Python course](/course/python?lang=en).

A broader learning path can continue through:

- the [Python course](/course/python?lang=en), which turns data into a verifiable project;
- the [SQL course](/course/sql?lang=en), which teaches reliable storage and reproducible answers;
- the [Rust course](/course/rust?lang=en), which develops an understanding of strict types, memory, and dependable command-line programs;
- the [Go course](/course/go?lang=en), which leads to a web application that could later host the model.

No graphics card or paid API is required. A regular computer, Python, and a willingness to run the examples yourself are enough. The course is open without registration. Start with [“Before you begin”](https://shanraq.org/read/kazakh-ai-before-start?lang=en), complete the entry check, and build the first part of the model yourself.
