# The bit: two distinguishable states

_Lead (summary):_ **Turn two reliably distinct states into a bit, then count how many messages several bits can represent.**

## Where we are on the map

The previous block showed how to find files and separate a program from its data. Now we ask how data can carry meaning that another person and another program will read in the same way. State your first guess, test it with numbers, and record what you had to revise.

![The bit: two distinguishable states](/static/course/informatics/map-11-bits-states-en.svg)

## Everyday analogy and exact model

Imagine a lamp with two clearly distinguishable states: off and on. Agree to label them 0 and 1. Looking at two lamps in order gives 00, 01, 10, and 11: four messages. Three lamps give eight combinations. Order matters: 01 and 10 are different.

A bit is one choice between two distinguishable states. Zero is a state too; it does not mean no data. Each additional bit doubles the number of possible sequences: one bit gives two, two give four, three give eight. Eight bits make one byte: 256 different sequences, from 00000000 to 11111111. A byte alone does not tell us whether it means a letter, number, or color. A reading rule supplies the meaning.

## Where the analogy ends

The lamp is only an analogy. Computers implement states in physical devices; an electrical signal need not be exactly zero or one. The device separates ranges using an agreed threshold. One bit does not have to occupy one separate physical object.

## The lesson's support signal

`2 states → 1 bit; 2 × 2 × 2 = 8 combinations; 8 bits = 1 byte = 256 combinations`

## Worked example

One bit can represent “done / not done” if we agree beforehand that 1 means done. Four task states (“not started,” “in progress,” “done,” “postponed”) need at least two bits. Make a table with 00, 01, 10, and 11, assigning exactly one state to each combination.

## Predict before observing

A card has two bits. Can it represent five distinct states without another rule? Write every combination before answering. Check your list: there are only four, so a third bit is needed.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Hide the table and list all eight three-bit combinations yourself. Say why 001 differs from 010. If you miss one, begin with two prefixes, 0 and 1, then add every two-bit ending to each.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“Zero is empty, so 000 does not count.” Correct this: 000 is one of eight valid codes. An empty message has no recorded bits; that is a different idea.

## Transfer to a new situation

Design a signal for a door with three states: shut, open, stuck. Explain why one bit is too small and two suffice. Do not claim that 11 names a fourth real state unless the group agrees on one.

## Project change

Add a done field to the Digital Assistant specification with two values, false and true. State that false means “not done” and true means “done.” These written JSON values do not occupy just one bit on disk: the file format adds bytes of its own.

## Exercise

Draw a four-code table for task states. Check that two different meanings never share one code. Then explain to another learner why eight bits give 256 combinations even though each position has only two choices.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow reconstruct the three-bit table from memory. In seven days find the minimum bit count for five states. In thirty days ask someone else to interpret done without an oral hint.

[Next lesson: How binary notation stores numbers](/read/informatics-12-binary-numbers?lang=en)
