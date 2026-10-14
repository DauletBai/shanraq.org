# Checkpoint: design lighting for a model room

_Лид (summary):_ **Complete the electrical page of your lab notebook: two independent lights, a calculation, a safe check and honest limits.**

## Situation and prediction

Your cardboard model room has a desk and an entrance. The desk light must switch off independently of the entrance light. Which circuit fits? Predict first; then compare series with parallel connections. Use closed paths, nodes, current, voltage and energy rather than memorising one formula.

## New words, step by step

A **design requirement** is a testable condition: each light has its own switch and the other remains on when one goes off. A **test criterion** is observable evidence that the requirement is met. In the calculation model, lights are represented by ohmic 300 Ω and 600 Ω resistors; real lamps and LEDs behave differently. State model limits and what you would measure before assembly: source voltage under load, branch currents and component current ratings. The model never connects to household mains.

## Working through supplied data

Let an ideal 3 V source feed two parallel branches. The 300 Ω branch takes 10 mA, the 600 Ω branch 5 mA, so total current is 15 mA. Turning off the second branch leaves the first at 3 V and 10 mA; total current becomes 10 mA. If both run for 60 s, Ptotal = 3 × 0.015 = 0.045 W and E = 2.7 J. These figures come from teaching data; they were not measured in a room.

## A safe check

On paper, draw a room plan, a two-branch circuit with separate switches, and a criterion–prediction–test table. If you wish, use a ready-made low-voltage kit only with an adult. LEDs need suitable resistors and correct polarity; our resistor model is not a ready wiring instruction for them. No mains lamps, outlets or direct cell shorts.

## A reasoning trap

A wrong conclusion is that one light went off because the other took its current. In the ideal parallel model, opening one branch does not interrupt the other. A real source may change voltage under load, so a physical device needs its own check.

## Support map

![Checkpoint: design lighting for a model room](/static/course/physics/map-56-electricity-investigation-mastery-en.svg)

## Transfer to a new case

Suppose the cell voltage sags to 2.8 V with both branches on. Which earlier numerical answers must be recalculated? Record why and one new measurement for the next block on magnetism.

## Lab notebook entry

Add a page to My World Lab Notebook with four lines: question; prediction before reading the answer; observation or calculation from supplied numbers; revised conclusion and its conditions. If you did not do an experiment, do not write ‘measured’. After the diagram, explain where charge moves, where energy is transferred, and where the model simplifies a real device.

[Return to the Physics course](/course/physics?lang=en)
