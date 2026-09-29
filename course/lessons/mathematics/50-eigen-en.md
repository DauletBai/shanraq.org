# Eigenvalues and eigenvectors: directions preserved by a transformation

_Lead (summary):_ **We find a matrix's special directions from Av=λv, interpret λ geometrically, and verify each pair by direct multiplication.**

## Where we are on the map

The matrix triples vector (1, 1) in the same direction and leaves vector (1, −1) unchanged

![A basis lets us choose directions for describing a space. A basis along which a matrix acts by simple scaling is especially useful.](/static/course/mathematics/map-50-eigen-en.svg)

## Begin with a familiar image

Stretching an image on a rubber sheet turns most arrows. Along a few lines, however, arrows stay on the same line and change only length or orientation.

## Precise meaning

A nonzero vector v is an eigenvector of A when `Av=λv`. The scalar λ is its eigenvalue: `|λ|` gives scale, a negative sign reverses direction, and λ=0 collapses it to zero. Find λ from `det(A−λI)=0`, then solve `(A−λI)v=0`. The zero vector is excluded.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

Av=λv → det(A−λI)=0 → find λ → find nonzero v → multiply to verify

## Worked example

For `A=[[2,1],[1,2]]`, `A(1,1)=(3,3)=3(1,1)`, so λ=3. Also `A(1,−1)=(1,−1)=1(1,−1)`, so λ=1.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Treating v=0 as an eigenvector is wrong. It supplies no direction and satisfies `A0=λ0` for every λ, so the definition explicitly excludes it.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** For `A=[[3,1],[1,3]]`, find the characteristic equation, eigenvalues, and one eigenvector for each. Verify every pair by direct multiplication.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

The checkpoint joins matrices, bases, and eigenvector directions through calculation, geometry, and meaning.

[Course contents](/course/mathematics?lang=en)
