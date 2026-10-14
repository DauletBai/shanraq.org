# Why can household lamps switch independently?

_Лид (summary):_ **Find the fork: parallel branches share voltage, and their currents add at a junction.**

## Situation and prediction

At home you turn one lamp off while another stays on. Separate switches on separate branches explain this better than last lesson's series path. Draw two routes between the same pair of points. Electricity does not choose only the easier route: if both branches are closed, both carry current.

## New words, step by step

A **parallel connection** provides several routes between the same two nodes. A **node** is a joining point where wires meet and the path can fork. The voltage between that common pair of nodes is the same for every branch: U₁ = U₂ = Usource in the ideal model. Conservation of charge gives Itotal = I₁ + I₂. Do not add branch voltages. More branches ask more current of the source; a real cell may sag under heavy load.

## Working through supplied data

In supplied data, 3 V appears across parallel 300 Ω and 600 Ω resistors. I₁ = 3/300 = 0.01 A = 10 mA; I₂ = 3/600 = 0.005 A = 5 mA. Itotal = 15 mA. Equivalent resistance Req = U/Itotal = 200 Ω, less than either branch resistance. These results assume ohmic resistors and an ideal source.

## A safe check

Trace both branches with different colours and open the switch in only one: the other remains closed. For a real demonstration, use only a checked low-voltage kit with supervision. Do not add branches before checking the allowed current for cells and wires.

## A reasoning trap

A mistake is to say current always takes only the path of least resistance. With finite resistances both paths carry current. Another is to divide the 3 V among branches as if they were in series.

## Support map

![Why can household lamps switch independently?](/static/course/physics/map-53-parallel-circuit-en.svg)

## Transfer to a new case

What becomes of Itotal if the 600 Ω branch opens while the other stays at 3 V? Name your assumption about the source.

## Lab notebook entry

Add a page to My World Lab Notebook with four lines: question; prediction before reading the answer; observation or calculation from supplied numbers; revised conclusion and its conditions. If you did not do an experiment, do not write ‘measured’. After the diagram, explain where charge moves, where energy is transferred, and where the model simplifies a real device.

[Next lesson: How can you find a break without guessing?](/read/physics-54-measuring-and-diagnosing?lang=en)
