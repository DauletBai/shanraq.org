# Probability: a measure of uncertainty from zero to one

_Lead (summary):_ **We build a sample space, distinguish an event from one outcome, and use complements, conditional probability, and a tree for sampling without replacement.**

## Where we are on the map

Combinatorics counted possibilities. Probability assigns weights and says how expected an event is before observing the result.

![A two-draw tree for a bag with three red and two blue balls gives probability 3/10 for two red balls](/static/course/mathematics/map-53-probability-en.svg)

## Begin with a familiar image

An opaque bag contains three red and two blue balls. We do not know the next colour, but we know the composition. Without replacement, that composition changes after the first draw.

## Precise meaning

An event probability lies from 0 to 1. For equally likely outcomes, `P(A)=m/n`. A complement has `P(not A)=1−P(A)`. Dependent steps use `P(A∩B)=P(A)P(B|A)`. Independence means learning A does not change the probability of B; it cannot be assumed for convenience.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

outcomes → event → branch weights → did the condition change? → multiply along paths → add matching paths → total probability 1

## Worked example

Without replacement, the first red has probability `3/5`. Then two of four balls are red, so `P(RR)=3/5·2/4=3/10`. Combinations check it: `C(5,2)=10` total pairs and `C(3,2)=3` red pairs.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Keeping the second-red probability at `3/5` after drawing a red ball is wrong. Only two red balls remain among four, so the conditional probability is `2/4`.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** A bag contains 4 white, 3 black, and 1 green ball. Draw two without replacement. Find probabilities that both are white, their colours differ, and at least one is green. Solve the last by a complement and check the complete set of cases.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

The next lesson moves from a random mechanism to an honest description of observed data.

[Course contents](/course/mathematics?lang=en)
