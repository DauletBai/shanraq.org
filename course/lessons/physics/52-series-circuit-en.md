# Why can one failed lamp darken a whole string?

_Лид (summary):_ **Explore one branch: equal current, added resistances and shared voltage.**

## Situation and prediction

An old-style series light string goes dark everywhere when one element breaks. Such a circuit has only one route. Trace it with your finger from source through both components and back. First predict whether less current leaves the first resistor than entered it. Then test the prediction against charge conservation.

## New words, step by step

A **series connection** has one route without a fork. Current I is therefore the same through every element; it is not consumed. For two ohmic resistors, Rtotal = R₁ + R₂. Their voltages add up to the source voltage: Usource = U₁ + U₂. This assumes ideal connections and negligible source internal resistance. Real lamps need not behave as fixed resistors.

## Working through supplied data

Let a 3 V source feed 100 Ω and 200 Ω in series. Rtotal = 300 Ω and I = 3/300 = 0.01 A = 10 mA everywhere. The first has U₁ = IR₁ = 1 V; the second has U₂ = 2 V. Check: 1 + 2 = 3 V. Removing the first resistor without replacing it by a wire creates an open circuit, not more current.

## A safe check

Draw the intact circuit and one with a break between components. If using a ready-made low-voltage kit, an adult disconnects the cells first and checks for a short path. Do not experiment with a mains light string; a drawing is enough.

## A reasoning trap

A mistake is to give each resistor the full 3 V and also add currents. In one branch, voltage drops add while current stays equal.

## Support map

![Why can one failed lamp darken a whole string?](/static/course/physics/map-52-series-circuit-en.svg)

## Transfer to a new case

Swap the two resistor positions. Will total current and the voltage across each physical resistor change? Explain without fresh arithmetic.

## Lab notebook entry

Add a page to My World Lab Notebook with four lines: question; prediction before reading the answer; observation or calculation from supplied numbers; revised conclusion and its conditions. If you did not do an experiment, do not write ‘measured’. After the diagram, explain where charge moves, where energy is transferred, and where the model simplifies a real device.

[Next lesson: Why can household lamps switch independently?](/read/physics-53-parallel-circuit?lang=en)
