# Lesson 4. Analyze a word without guessing

## Why this matters

One meaning can appear in several word forms. To recognize “house” in `үйде` and `үйлерде`, our program looks for a **stem**: the part from which forms are built. Other parts add meanings. A stem dictionary lists stems known to the program and what they mean.

## See the whole process

Ordinary spelling gives `үй`, `үйлер`, `үйлерде`. We insert hyphens only in our analysis: `үй`; `үй-лер`; `үй-лер-де`. `-лер` means plural; `-де` indicates location, interpreted with context. `үйлерде` can mean “in the houses”.

Another series is `мектеп` (“school”), `мектептер` (“schools”), `мектептерде` (“in the schools”). Why `-тер` instead of `-лер`? The sounds of an ending adapt to the stem. A **variant of an ending** is another sound form with the same job. We will construct the choice rule later; for now, do not swap endings by guesswork.

## Recall map

Written word → known stem → remaining parts → meaning of each part. Hyphens appear only in a teaching analysis.

## Recall without looking

Explain the difference between `үй` and `үйлерде`. Why should we not assume `мектеплер` is correct by copying `үйлер`?

## Task

Analyze `үйлер`, `үйлерде`, `мектептер`, and `мектептерде`. Name the stem, plural part if present, and location part if present. **Hint:** use only forms explained above. **Check:** `үй-лер`; `үй-лер-де`; `мектеп-тер`; `мектеп-тер-де`. The first and third have no location part. For an unknown stem, the honest result is “cannot analyze yet”.
