# How can you find a break without guessing?

_Лид (summary):_ **Practise safe fault-finding: diagram, prediction, measurement and a conclusion that respects the meter range.**

## Situation and prediction

A model lamp will not shine. You could swap parts at random, or narrow the cause by first asking whether a closed path exists. Draw the circuit and mark expected readings before checking. This is the course's original predict–test–explain cycle.

## New words, step by step

**Diagnosis** means finding a cause from observable signs. An **open circuit** breaks the path, so ideal current in its sole series branch is 0 A. A meter's **range** is the largest value in the selected setting. Measure current with an ammeter in a break; measure voltage with a voltmeter between two points. A battery showing about 3 V without a load does not prove the lamp is healthy. A reading at an open gap depends on exactly which two points the probes touch.

## Working through supplied data

Scenario: a 3 V source, switch, lamp and a hidden break between switch and lamp. Closing the switch still gives 0 A in the sole branch. In an ideal model, the whole source voltage may appear across the two ends of that break, whereas a healthy lamp carrying no current has nearly zero voltage drop. Real kits have contact and meter effects. One zero reading does not prove the lamp failed.

## A safe check

Start on paper: mark two probe positions and possible outcomes. On a ready-made safe kit, an adult turns off the source before moving wires, chooses mode and range, then checks the connections. Do not inspect household mains wiring.

## A reasoning trap

A dangerous mistake is to put a meter in ammeter mode across the battery: its low internal resistance can create a short circuit. Lack of light also does not identify a unique cause.

## Support map

![How can you find a break without guessing?](/static/course/physics/map-54-measuring-and-diagnosing-en.svg)

## Transfer to a new case

Make a three-row hypothesis table: open switch, broken wire, failed lamp. Name an observation that distinguishes at least two of these.

## Lab notebook entry

Add a page to My World Lab Notebook with four lines: question; prediction before reading the answer; observation or calculation from supplied numbers; revised conclusion and its conditions. If you did not do an experiment, do not write ‘measured’. After the diagram, explain where charge moves, where energy is transferred, and where the model simplifies a real device.

[Next lesson: How much energy does a tiny lamp use?](/read/physics-55-electric-power-energy?lang=en)
