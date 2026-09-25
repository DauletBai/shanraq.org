# Before you begin: AI without an LLM

Across 30 lessons you will build a small system for short questions about an imaginary school club. It analyzes familiar Kazakh forms, identifies the question type, finds an approved record, and shows an answer with its source. Without enough evidence it says `білмеймін` (“I don't know”).

The course assumes no prior study of artificial intelligence. It needs no paid API, permanent internet connection, or dedicated GPU. It is **not a first programming course**, however: from lesson 16 onward we combine files, JSON, functions, modules, and dates, then measure and run a multi-file program.

## What you will have at the end

You will run a finished local model in the console, inspect the path from question to source, and produce its technical specification. A control run on 25 September 2026 used an ARM computer with macOS 26.5.2 and Python 3.14.5:

| Measurement | Result |
|---|---:|
| Model code and data | 6,907 bytes (6.7 KiB) |
| Whole process memory including Python | 20.3 MiB |
| Fresh start / request in a ready process | 35.214 / 0.034 ms |
| Requests per second in the control run | 41,827 |
| Open quality scenarios | 16 of 16 |
| Confident answers in ten cases requiring refusal | 0 of 10 |
| External API requests | 0 |

One KiB is 1,024 bytes, one MiB is 1,024 KiB, and a millisecond is one thousandth of a second. A fresh start is like opening a workshop; a ready request uses tools already laid out. Time and memory will differ on another computer, so the program also reports its test environment. `0 of 10` is only the result of ten published refusal checks, not a promise of zero errors on arbitrary text. In [lesson 30](https://shanraq.org/read/kazakh-ai-30-cli?lang=en), you will repeat the experiment and use a table to identify where this narrow model is stronger than a large LLM, where it falls short, and when combining both approaches helps.

## What you should know first

The reliable beginner route is lessons 1–17 of the [Python course](https://shanraq.org/course/python?lang=en). You need variables, strings, lists, dictionaries and sets; conditions and loops; functions and `return`; files, JSON, imports, and errors. SQL is optional. The [SQL course](https://shanraq.org/course/sql?lang=en) helps with keys, constraints, and auditable records, but this course explains every required catalog idea again. Go and Rust are not prerequisites.

Check yourself before lesson 1. You are ready if you can open a terminal, run a Python file, explain `=` versus `==`, loop over a list, retrieve a dictionary value, write a function with `return`, and load JSON. If three or more tasks are unfamiliar, take the introductory Python block first. Lesson 2 helps set up the workspace but does not replace those foundations.

## Five terms that remove the AI black box

**Artificial intelligence (AI)** is an umbrella term for programs that perform tasks involving recognition, selection, or inference. It does not name one device or one technique.

An **algorithm** is an exact sequence of actions. **Data** are the records those actions use. A **model** is a stored way of turning input into a result: it may contain human-written rules or parameters found from examples. **Training** finds parameters from labeled examples. **Using a model** applies finished rules or parameters to a new input.

Lessons 1–20 build a **rule-based system**: a person specifies word parts, checks, and fact cards. Lessons 21–25 add a tiny **learned model** that counts features in labeled questions. Lessons 26–30 connect both parts to source checks. A **large language model (LLM)** is only one kind of learned model; AI and LLM are not synonyms.

The whole course in one support signal:

`question → word analysis → question type → key → approved card → sourced answer`
`ambiguity, conflict, or missing card → explained refusal`

An arrow means “pass the result to the next step.” On the maps, a red border marks an action or data, a green arrow marks order, and a gray note marks a limit. We decode these signals in words before asking you to recall them.

## How languages package one idea

Take “to our houses”: house, plurality, possession, and direction. **Grammar** supplies rules for combining such meanings. A **morpheme** is the smallest word part with a function. **Morphology** studies how words are built from morphemes.

- **Analytic structure** resembles separate cards on a table. English `to our houses` expresses direction and possession with the separate words `to` and `our`. Strong analytic traits occur in English, Mandarin Chinese, Vietnamese, Thai, and Khmer.
- **Agglutinative structure** resembles pieces placed in order on a skewer. For analysis, Kazakh `үйлерімізге` becomes `үй-лер-іміз-ге`: house — plural — ours — toward. Ordinary spelling has no hyphens.
- **Fusional structure** resembles one label carrying several messages. Russian `столом` uses `-ом` to express both singular number and instrumental case. Fusional traits are prominent in Russian, Polish, German, Spanish, and Lithuanian.

These are strategies, not sealed boxes. English has an ending in `houses`; Kazakh uses auxiliary verbs in analytic constructions. A Kazakh stem may also alternate: `кітап` becomes `кітабым`, changing `п` to `б`. The skewer illustrates order and visible functions, not permanent letters.

### What a percentage can actually mean

There is no standard scientific “percentage of agglutinativeness” for a whole language. A number requires a stated denominator and corpus. WALS examines selected case and tense, aspect, and mood markers. **Aspect** describes whether or how an action unfolds; **mood** marks its status, such as a fact, command, or condition. WALS assigns a category, not a percentage.

On a narrow screen, swipe wide tables left and right.

| Language | Branch | WALS fusion category | Meaning carried by the sampled case marker |
|---|---|---|---|
| Turkish | Turkic | exclusively concatenative | case only |
| Hungarian | Uralic | exclusively concatenative | case only |
| Japanese | Japonic | exclusively concatenative | case only |
| Korean | Koreanic | exclusively concatenative | case only |
| Basque | language isolate | exclusively concatenative | case only |
| Georgian | Kartvelian | exclusively concatenative | case only |
| Kannada | Dravidian | exclusively concatenative | case only |
| Imbabura Quechua | Quechuan | exclusively concatenative | case only |
| Finnish | Uralic | exclusively concatenative | case and number together |
| Evenki | Tungusic | exclusively concatenative | case plus whether the participant is identified |

Rows run from a marker with one job to markers that combine jobs; there is no scientific ranking within those groups. The last column reveals mixed properties in Finnish and Evenki. WALS even puts Russian and English in the same fusion category, although a Russian ending may carry several meanings and English widely uses separate words. One measure is insufficient.

Joseph Greenberg proposed a genuine numerical **agglutination index**: divide predictable morpheme junctions by all morpheme junctions in a selected text. His old 100-word samples produce the descending list below. These are percentages of **junctions in particular texts**, not portions of whole languages:

| Rank | Language and sample | Junction index |
|---:|---|---:|
| 1 | Swahili | 67% |
| 1 | spoken Turkish | 67% |
| 3 | written Turkish | 60% |
| 4 | Yakut | 51% |
| 5 | Greek | 40% |
| 6 | English | 30% |
| 7 | Inuit | 3% |

The low number for morphologically rich Inuit exposes the index's limit: it measures predictability at junctions, not word length or an amount of “agglutination.” Modern corpus research also finds that such indices depend on text genre and the rule used to divide words. These sources contain no comparable values obtained by one method for all ten languages above or for Kazakh. Assigning Kazakh 98%, 99%, or 100% and ranking it first would therefore invent evidence.

An academic grammar supports a different precise statement: Kazakh derivational and inflectional morphology uses suffixes. That strong fact does not erase sound alternations or auxiliary constructions. Persian or Russian origin also does not automatically lower a percentage; a borrowed stem can take ordinary Kazakh endings.

## Why this course uses Kazakh

No natural language is inherently more logical or mathematical than another. Kazakh is especially suitable **for this teaching project** because much grammatical information is expressed by ordered suffixes; a chain can be drawn as `stem → plural → possession → case`; many instructional boundaries are visible; sound variants follow a finite set of rules; and we have linguistic competence plus the `qazaq-ir` and `adam` projects for comparison.

Ordered suffixes map naturally to tables, conditions, and diagrams of allowed transitions between word parts. That makes Kazakh morphology clear for a programming lesson. It does not make the whole language a formula: context, meaning, alternations, and exceptions still need evidence.

## Why build AI without an LLM?

An LLM can generate free text, but a fluent answer can contain an unsupported claim: a **hallucination**. Our system creates no new factual claim. It answers only after finding one valid card and displays its source. Think of a librarian who shows the catalog card or admits no record was found.

Errors remain possible. A wrong card, incomplete rule, bad analysis, or software defect can still yield a false answer. The design reduces one specific risk—unsupported fact generation—and leaves a visible audit trail.

A small local program may use less computation and keep questions off an external service. Lesson 28 measures one classifier; lesson 30 measures the whole system. For open LLMs, the calculated weight size exceeds our file size by tens of millions, but that factor describes only the stated representations: the systems have different tasks and capabilities.

## How to study with support signals

The course adapts Viktor Shatalov's approach: see a whole block first, study the full explanation, compress it into a support map, recall it without looking, then act on a new example. Every five lessons ends with a mastery gate that can be retried. A mistake points back to a specific support signal rather than closing the path forward.

Use this loop: life analogy → complete result → terms and symbols → hidden-page recall → task on new data → reference answer → explain the mistake. Reading code alone is not completion.

The six gates are: parse a word by hand; perform that analysis in code; produce one unambiguous question key; answer only from a valid card; measure classifier quality honestly; run the final program and defend its limits.

Sources: [WALS on fusion](https://wals.info/chapter/20), [WALS on exponence](https://wals.info/chapter/21), [Greenberg's quantitative morphology paper](https://doi.org/10.1086/464575), [modern review of corpus indices and their limits](https://pmc.ncbi.nlm.nih.gov/articles/PMC9159679/), [Kazakh grammar description](https://slaviccenters.duke.edu/sites/slaviccenters.duke.edu/files/file-attachments/kazakh-grammar.pdf), [rule-based Kazakh morphological analysis](https://aclanthology.org/W14-2806/), [auxiliary verb constructions in modern spoken Kazakh](https://openresearch.surrey.ac.uk/esploro/outputs/doctoral/Auxiliary-verb-constructions-in-Modern-Spoken/99653765502346), [survey of factual errors in LLMs](https://aclanthology.org/2024.emnlp-main.1088/), [Viktor Shatalov's reference signals in the Russian State Library](https://search.rsl.ru/ru/record/01007628586).
