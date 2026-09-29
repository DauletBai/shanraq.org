# Graphs and algorithms: nodes, connections, and the best route

_Lead (summary):_ **We model a network with vertices and edges, distinguish paths and cycles, execute Dijkstra's algorithm by hand, and verify that the route is truly shortest.**

## Where we are on the map

Logic supplied rules for valid transitions. A graph turns them into edges between objects, while an algorithm gives a reproducible order of actions.

![A weighted graph with six vertices shows a shortest A-to-F path of length 13](/static/course/mathematics/map-57-graphs-algorithms-en.svg)

## Begin with a familiar image

Road maps, social networks, task dependencies, and web links share one structure: objects and connections. Edge weights may mean time, distance, or cost, but their units must be consistent.

## Precise meaning

A graph has vertices and edges. Edges may be directed or undirected, weighted or unweighted. A path follows adjacent vertices. With nonnegative weights, Dijkstra repeatedly finalizes the nearest unsettled vertex and relaxes its neighbours. Negative edges require another algorithm.

## The lesson's support signal

`meaning → model → calculation → check → explain in your own words`

model nodes and edges → weight unit → initial distances → nearest vertex → update neighbours → repeat → reconstruct path

## Worked example

From A we get C=2 and B=4. Through C, B improves to 3. Through B, D=8; through D, E=10; through E, F=13. Predecessors give `A→C→B→D→E→F`, with `2+1+5+2+3=13`.

## Example with a fading prompt

Change the numbers in the example and repeat the solution without copying its steps. Estimate first, calculate exactly, and check against the original condition.

## Recall without a prompt

1. Define the main idea in one sentence.
2. Rebuild the support signal from memory.
3. Solve the example with different numbers and explain why every step is valid.

## Find and correct the mistake

Picking the shortest outgoing edge at every intersection and calling the result globally shortest is wrong. A cheap local step can lead to an expensive continuation; the algorithm compares full distances from the start.

## Transfer to a new setting

Apply the same relationship to something you can measure yourself. Name the quantities, units, and the boundary within which the model is valid. Check the answer by a second representation or an inverse operation.

## Exercise

**Required.** Use a Dijkstra table to find shortest distances from A to every vertex. Remove edge D–E and recompute the route to F. Give one example where unweighted breadth-first search is the right tool instead.

**With your own data.** Create one everyday example with the same structure and show it in words, a formula, and a check.

**Optional.** Seven days later, change the numbers and repeat without the support map.

## Open perspective

The next lesson joins data, probability, functions, matrices, and graphs in one testable modeling cycle.

[Course contents](/course/mathematics?lang=en)
