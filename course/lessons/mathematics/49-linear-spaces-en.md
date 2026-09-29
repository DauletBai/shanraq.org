# Vector spaces and bases: the directions that build a vector

_Lead (summary):_ **We understand linear combinations, span, independence, basis, and coordinates through two directions that build every plane vector uniquely.**

## Where we are on the map

Vector (6, 2) decomposed as four vectors (1, 1) and two vectors (1, −1)

![A matrix mixed the components of an input. We now ask why components describe a vector and which directions may replace the usual horizontal and vertical axes.](/static/course/mathematics/map-49-linear-spaces-en.svg)

## Begin with a familiar image

A city location can be reached along north–south and east–west streets, or along two other nonparallel roads. The selected independent directions act as a basis.

## Precise meaning

A linear combination of `b₁,…,bₖ` is `c₁b₁+…+cₖbₖ`. Their span contains every vector obtainable this way. The vectors are linearly independent when only zero coefficients produce the zero vector. A basis is independent and spans the whole space; its number of vectors is the dimension.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

vectors → linear combinations → span → independence → basis → coordinates

## Worked example

Take `b₁=(1,1)`, `b₂=(1,−1)`, and `v=(6,2)`. Solving `c₁+c₂=6`, `c₁−c₂=2` gives `c₁=4`, `c₂=2`. Thus `v=4b₁+2b₂`. The coordinates are unique because b₁ and b₂ are not parallel.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Calling `a=(1,2)` and `b=(2,4)` a basis of the plane merely because there are two vectors is wrong. They are parallel and span only one line.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** For `a=(1,2,0)`, `b=(0,1,1)`, and `c=(1,3,1)`, determine dependence. Find two independent vectors spanning the same plane and represent `r=(2,5,1)`.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

Some directions are not turned by a matrix; they are only stretched, compressed, or reversed. These are eigenvector directions.

[Course contents](/course/mathematics?lang=en)
