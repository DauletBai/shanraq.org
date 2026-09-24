# Before you begin: AI without an LLM

We will build a small program for short questions about an imaginary school club. It analyzes Kazakh words, finds an approved record, and answers from that record. Without a record it says “білмеймін” (“I don't know”). You need a computer and Python 3.10 or newer; Lesson 2 explains setup. The course is free and needs neither a paid API nor a dedicated GPU.

## How languages express the same idea

Consider “to our houses”: a house, plurality, possession, and direction. Languages place these meanings differently. **Grammar** is the rules for combining them; **morphology** concerns word structure. The labels below describe common patterns, not every word of a language.

- **Analytic structure** is like separate cards on a table: English `to our houses` puts direction and possession in the separate words `to` and `our`. Strong analytic traits occur in English, Mandarin Chinese, Vietnamese, Thai, and Khmer. A program must examine neighboring words to recover the whole meaning.
- **Agglutinative structure** is like food placed in order on a skewer: distinct parts attach to a meaningful stem. For analysis only, write Kazakh `үйлерімізге` as `үй-лер-іміз-ге`: house, plural, ours, toward. Ordinary spelling contains no hyphens.
- **Fusional structure** is like one card carrying two instructions: Russian `столом` uses `-ом` to express both singular number and instrumental case (“with what?”). Fusional traits are prominent in Russian, Polish, German, Spanish, and Lithuanian.

The skewer image does not mean that every stem letter stays fixed: `кітап` (“book”) becomes `кітабым` (“my book”), changing `п` to `б`. A program must include that rule too.

### Ten languages with agglutinative word forms

| Language | Main area | What matters here |
|---|---|---|
| Kazakh | Kazakhstan | Visible suffix chains; the language of our examples |
| Kyrgyz | Kyrgyzstan | Suffix chains with its own sound rules |
| Karakalpak | Karakalpakstan | A Turkic language close to Kazakh |
| Tatar | Tatarstan | Agglutinative forms alongside language contact |
| Bashkir | Bashkortostan | Agglutinative forms alongside language contact |
| Uzbek | Uzbekistan | Suffix chains; contact with Persian |
| Uyghur | China and Central Asia | Word parts attached in sequence |
| Turkmen | Turkmenistan | Suffixes attached in sequence |
| Azerbaijani | Azerbaijan and Iran | Stem, number, possession, case |
| Turkish | Türkiye | Clear suffix chains |

A word borrowed from Persian can take ordinary Turkic endings. Its origin does not give us a percentage of a language's “agglutinativeness”. We chose Kazakh because its forms make the steps visible, we can check examples ourselves, and our `qazaq-ir` and `adam` projects offer practical reference points. Other agglutinative languages would also work.

## Why build AI without an LLM?

A **large language model** (LLM) can freely generate text. This is useful for many tasks, but it can produce a plausible claim unsupported by a source: a **hallucination**. Our system gives a factual answer only when it finds an approved record. Think of a librarian showing a catalog card instead of inventing a missing book. Errors in analysis, records, or answer selection remain possible.

A small local program may use fewer resources and keep questions off external services. “A thousand times lighter, cheaper, and faster” remains an unproven hypothesis. In lesson 28 we measure our own program and explain what a fair comparison with an LLM would still require. The course has 30 lessons.

Sources: [WALS](https://wals.info/chapter/20), [Kazakh grammar description](https://slaviccenters.duke.edu/sites/slaviccenters.duke.edu/files/file-attachments/kazakh-grammar.pdf), [Turkic languages and language contact](https://www.iranicaonline.org/articles/iran-vii7-turkic-languages/), [survey of factual errors in LLMs](https://aclanthology.org/2024.emnlp-main.1088/).
