# Differential equations: a law of change as a process model

_Lead (summary):_ **We translate a verbal law into an equation with a derivative, apply an initial condition, and check the solution by substitution, units, and behavior.**

## Where we are on the map

The starting information is now a relationship between a quantity and its rate of change.

![Growth y′=0.05y starts at 100 and produces the curve y=100e^(0.05t)](/static/course/mathematics/map-46-differential-equations-en.svg)

## Begin with a familiar image

When an account repeatedly earns a percentage of its current balance, the addition grows with the amount already present. “Rate is proportional to amount” becomes `y′=ky`.

## Precise meaning

A differential equation relates an unknown function to its derivatives. The family `y=Ce^{kt}` solves `y′=ky`; the initial condition `y(0)=y₀` selects C=y₀. The sign of k determines growth or decay, and its unit is inverse time.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

quantities → verbal law → equation → solution family → initial condition → check → model limits

## Worked example

For `y′=0.05y`, `y(0)=100`, the solution is `y=100e^{0.05t}`. Check: `y′=5e^{0.05t}=0.05y`, and the value at t=0 is 100.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Using a positive k for decay and never checking behavior is wrong. Substitution, sign, units, and the long-term limit must agree with the process.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** A substance follows `M′=−0.3M`, `M(0)=50`. Find M(t), the time to reach 10 units, and check by differentiation and monotonicity.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

The checkpoint combines limits, derivatives, optimization, integrals, and differential models in one system.

[Course contents](/course/mathematics?lang=en)
