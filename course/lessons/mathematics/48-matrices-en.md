# Matrices: a table of numbers that transforms data

_Lead (summary):_ **We begin linear algebra with a concrete role for a matrix: storing coefficients, checking dimensions, and turning an input vector into a new set of numbers.**

## Where we are on the map

A two-by-two matrix transforms vector (4, 2) into vector (10, 10)

![Systems gave us rows of coefficients and geometric vectors gave us columns of coordinates. A matrix gathers that structure into one object.](/static/course/mathematics/map-48-matrices-en.svg)

## Begin with a familiar image

A control panel with two inputs and two displays mixes signals: each display combines the inputs with its own coefficients. Its settings table is a matrix.

## Precise meaning

A matrix is a rectangular table of numbers. An `m×n` matrix has m rows and n columns. In `A·x`, every row of A forms a dot product with column x. The operation exists only when the number of columns of A equals the number of components of x.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

dimensions → row · column → one component → every row → new vector → meaning check

## Worked example

Let `A=[[2,1],[1,3]]` and `x=(4,2)`. The first row gives `2·4+1·2=10`; the second gives `1·4+3·2=10`. Thus `Ax=(10,10)`. Dimensions confirm `(2×2)·(2×1)=(2×1)`.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Multiplying a 2×3 matrix by a two-component vector is invalid. The inner dimensions 3 and 2 do not match.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** For `A=[[3,−1],[2,4]]`, `u=(5,2)`, and `v=(−1,3)`, find Au, Av, and A(u+v). Verify `A(u+v)=Au+Av` and state every dimension.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

The next lesson uses a basis to explain which vectors can be built from selected directions.

[Course contents](/course/mathematics?lang=en)
