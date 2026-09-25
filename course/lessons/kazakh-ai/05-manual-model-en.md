# Lesson 5. Build a model you can run by hand

## Why this matters

Before trusting a program to analyze words, we should be able to perform the steps ourselves. A **manual model** is a written procedure that lets any reader reproduce the result. It exposes missing rules.

## See the whole process

Our dictionary currently has two stems: `үй` (“house”) and `мектеп` (“school”). Our teaching table allows only these sequences: stem; stem + plural part; stem + plural part + location part. Here `+` means “the next part of the analysis”, not arithmetic addition. Use `-лер` after `үй` and `-тер` after `мектеп`; then add `-де`. This is a **small teaching table**, not a complete rule for Kazakh.

Analyze `үйлерде`: 1) recognize `үй`; 2) recognize `-лер`; 3) recognize `-де`; 4) read “in the houses”. `мектептерде` gives “in the schools”. For `кітаптарда`, say “unknown for now”: `кітап` is absent from our dictionary, even if a person knows the word.

## Recall map

Dictionary → find stem → check next part → check order → output an analysis or “unknown”.

![Lesson 5 support map](/static/course/kazakh-ai/map-05-manual-model-en.svg)

## Recall without looking

Describe every step for `үйлерде`. Name a boundary of our table: which stems and endings does it not yet know?

## Task

Use the table to analyze `үй`, `үйлер`, `мектептерде`, `кітаптарда`, and `мектеплерде`. Give the result or the first step at which checking stops.

**Check:** `үй` — house; `үйлер` — houses; `мектептерде` — in the schools; `кітаптарда` — unknown stem; `мектеплерде` — after the known stem, our table expects `-тер`, not `-лер`. We will learn the fuller sound rules later.

## Gate 1: ready to continue?

Without looking, draw `stem → plural → location`, analyze all five words, and explain both refusals aloud. Pass when at least four of five analyses are correct and you can distinguish an unknown stem from a forbidden ending. Otherwise, return to the table, repair only the failed step, and retry with `үйлерде`.
