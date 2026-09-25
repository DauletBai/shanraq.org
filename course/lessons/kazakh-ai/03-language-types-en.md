# Lesson 3. Where does a language put grammatical meaning?

## Why this matters

A program sees letters, but we need meaning. “To our houses” contains a house, plurality, possession, and direction. These additions are **grammatical meanings**. Before writing rules, we must ask where a language stores them: in separate words, in a chain of parts, or together in one part.

## New terms before the examples

A **morpheme** is the smallest word part with a function. In the learning analysis `үй-лер-іміз-ге`, hyphens expose morpheme boundaries. A **root** carries the central dictionary meaning; here it is `үй`. A **suffix** follows a root or stem. Ordinary spelling does not contain these explanatory hyphens.

**Analytic**, **agglutinative**, and **fusional** name strategies for expressing grammatical meaning. They are not language families, levels of development, or rigid boxes. One language may combine them.

## Three strategies for one idea

| Strategy | Everyday image | Where the meaning sits | Example | Languages with a prominent trait |
|---|---|---|---|---|
| Analytic | separate cards on a table | mainly in separate function words and their order | English `to our houses`: `to` gives direction, `our` possession | English, Mandarin Chinese, Vietnamese, Thai, Khmer |
| Agglutinative | pieces placed in order on a skewer | a sequence of distinct morphemes | Kazakh `үй-лер-іміз-ге`: house — plural — ours — toward | Kazakh, Turkish, Kyrgyz, Finnish, Japanese |
| Fusional | one label with several marks | one morpheme expresses several meanings | Russian `книгами`: `-ами` expresses plural and instrumental case | Russian, Polish, German, Spanish, Lithuanian |

The skewer answers only two questions: are the parts distinct, and in what order do they attach? It does **not** promise that every root letter remains fixed. In `кітап → кітабым`, final `п` changes to `б`. A suffix also selects a sound variant: `үй-лер` but `кітап-тар`. Many such changes are predictable, but a program must still encode them.

## Why percentages do not make a simple ranking

The preface contains two different tables. The first shows a WALS category and meaning combinations for ten languages. The second gives Greenberg's old corpus index: the percentage of predictable morpheme junctions in one particular 100-word text. Neither measure means “the entire language is this percent agglutinative.”

Deeper comparison needs separate measures: morphemes per word, meanings per morpheme, stem alternation, and relations expressed by separate words. One invented percentage hides those reasons. A loan does not import its source language's grammar either: `мектеп` has Arabic-Persian history but accepts Kazakh `мектеп-тер-іміз-де`.

## Why Kazakh is convenient for a program

Kazakh derivational and inflectional morphology is built mainly with suffixes. The order `stem → plural → possession → case` can be represented as allowed transitions. In programming, this resembles a route with labeled stops: each stop restricts the next valid step.

That is a convenience for **morphological analysis**, not proof that Kazakh is more mathematical than other languages. Auxiliary constructions, ambiguity, sound alternation, and context remain. We use Kazakh because the suffix chain is visible, the examples matter to our readers, and they can be compared with `qazaq-ir`.

## Recall map

Separate words → analytic. Ordered distinct morphemes → agglutinative. Several meanings in one morpheme → fusional. One language may use several strategies.

![Lesson 3 support map](/static/course/kazakh-ai/map-03-language-types-en.svg)

## Recall and check

Hide the table and draw three images: cards, skewer, and a label with two marks. Explain **where grammatical meaning sits**, rather than merely naming the strategy.

Classify and justify: 1) English `in our school`; 2) Kazakh `мектеп-тер-іміз-де`; 3) Russian `книгами`; 4) Turkish `ev-ler-imiz-e` (“to our houses”). Then explain why English `-s` or a Kazakh auxiliary verb does not erase the language's predominant strategy.

**Check:** 1 is analytic because `in` and `our` are separate words; 2 and 4 are agglutinative chains with distinct jobs; 3 is fusional because `-ами` combines number and case. **Common mistake:** deciding from a word's origin or declaring a whole language “pure” from one form.

Sources: [WALS on fusion](https://wals.info/chapter/20), [WALS on exponence](https://wals.info/chapter/21), [Kazakh morphology description](https://slaviccenters.duke.edu/sites/slaviccenters.duke.edu/files/file-attachments/kazakh-grammar.pdf), [rule-based Kazakh morphology](https://aclanthology.org/W14-2806/).
