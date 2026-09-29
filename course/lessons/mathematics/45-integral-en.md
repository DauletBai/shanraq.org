# The integral: accumulation from small changes

_Lead (summary):_ **We add narrow contributions, understand a definite integral as a limit of sums, and link area, accumulation, and antiderivatives through the fundamental theorem.**

## Where we are on the map

The derivative broke change into an instantaneous rate; the integral assembles small changes into a total.

![The area under y=2x from 0 to 3 equals 9 by both a limiting sum and an exact integral](/static/course/mathematics/map-45-integral-en.svg)

## Begin with a familiar image

A water meter accumulates flow over time. Even when the rate changes each minute, small volumes add to total consumption.

## Precise meaning

The definite integral `∫[a,b]f(x)dx` is the limit of sums `f(xᵢ)Δx` as the partition gets finer. It measures signed accumulation. If `F′=f`, the fundamental theorem gives `∫[a,b]f(x)dx=F(b)−F(a)`.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

interval → partition → height × width → sum → limit → antiderivative → boundary difference

## Worked example

For `f(x)=2x`, an antiderivative is `F=x²`, so `∫[0,3]2x dx=9`. Geometry agrees: a triangle of base 3 and height 6 has area 9.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Writing only F(b) and forgetting the lower endpoint is wrong. A definite integral is a change in accumulation, so it requires `F(b)−F(a)`.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** Evaluate `∫[−1,2](2x+1)dx`, show signed parts on a sketch, and separately find the ordinary geometric area between graph and axis.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

When the function is unknown but a law for its change is known, we obtain a differential equation.

[Course contents](/course/mathematics?lang=en)
