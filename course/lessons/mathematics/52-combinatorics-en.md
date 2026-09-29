# Combinatorics: count possibilities without listing them all

_Lead (summary):_ **We identify stages of choice, distinguish an ordered selection from a set, and use sum, product, permutation, and combination rules with a small enumeration check.**

## Where we are on the map

Multiplication, fractions, and variables now help count all possible solutions before one outcome is chosen.

![Five objects produce 20 ordered pairs and 10 unordered pairs](/static/course/mathematics/map-52-combinatorics-en.svg)

## Begin with a familiar image

Take five books A, B, C, D, and E. For first and second place, AB and BA are different results. For choosing two books for one shelf, both orders describe the same pair.

## Precise meaning

The sum rule adds mutually exclusive cases. The product rule multiplies possibilities across successive stages. `P(n,k)=n!/(n−k)!` keeps order; `C(n,k)=n!/(k!(n−k)!)` ignores order. Before using a formula, ask whether repetition is allowed and whether order matters.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

describe the choice → split into stages → order? → repetition? → sum or product → formula → enumerate to check

## Worked example

Choosing a winner and runner-up from five books gives `P(5,2)=5·4=20`. Choosing two books without ranks counts every pair twice, so `C(5,2)=20/2=10`. The list AB, AC, AD, AE, BC, BD, BE, CD, CE, DE confirms ten pairs.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Using `C(5,2)=10` to choose a chair and secretary is wrong. The roles differ, so each pair supports two assignments; the correct count is `P(5,2)=20`.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** From 8 students, count ways to choose: (a) a class representative and deputy; (b) a two-person committee; (c) a three-person committee containing one specified student. Check every answer on a smaller four-student case.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

The next lesson turns counts of equally possible outcomes into a measure of uncertainty: probability.

[Course contents](/course/mathematics?lang=en)
