# What a model cannot be trusted with: fact, inference, cause

_Лид (summary):_ **Lesson fifty-one of the Python course. Separate an observation, a computed result, and a causal explanation. Code can verify a value and reproduce a calculation, but plausible causal prose remains a hypothesis until it has independent evidence and human review.**

## Why this matters

In lesson 50 the model could no longer replace a number: the program compared it with a `source_id`. A correct number can still acquire a false story:

> Inflation was 11.4% because the money supply grew.

The first part may be in a table. “Because” adds a causal relationship that one row does not establish. A language model can choose a convincing explanation; conviction is not evidence.

We need three kinds of claim:

- an **observation** appears directly in a source;
- a **calculation** is produced by reproducible code, such as a difference;
- a **cause** explains why something changed and needs separate evidence.

## The whole thing first

The file `trust.py` applies an admission policy. The model may classify a sentence, but code owns the final decision.

```python
"""Lesson 51: a causal story does not pass as an observation."""

EVIDENCE = {"inflation-kz-2025", "inflation-kz-2024"}


def decide(candidate):
    """Route a claim without guessing whether polished prose is true."""
    missing = set(candidate["evidence_ids"]) - EVIDENCE
    if missing:
        return "reject", "unknown source"
    if candidate["claim_type"] == "cause":
        return "review", "a cause needs independent evidence"
    if candidate["claim_type"] == "calculation":
        return "review", "a reproducible formula is required"
    if candidate["claim_type"] == "observation":
        return "accept", "the value can be checked against its source"
    return "reject", "unknown claim type"


candidates = [
    {
        "text": "Inflation in Kazakhstan has a published value.",
        "claim_type": "observation",
        "evidence_ids": ["inflation-kz-2025"],
    },
    {
        "text": "Inflation rose because the money supply grew.",
        "claim_type": "cause",
        "evidence_ids": ["inflation-kz-2025"],
    },
]

for candidate in candidates:
    decision, reason = decide(candidate)
    print(candidate["claim_type"], "->", decision, "—", reason)
```

Output:

```text
observation -> accept — the value can be checked against its source
cause -> review — a cause needs independent evidence
```

## Claim type is data

Constrain `claim_type` in the schema:

```python
"claim_type": {
    "type": "string",
    "enum": ["observation", "calculation", "cause"],
}
```

The model may still label a cause as an observation. Its label helps routing; it is not proof. For automatic publication, use a narrow observation template assembled from verified fields. Send free-form prose to review.

## Correlation does not answer why

Two series may move together because of a direct effect, reverse causation, a third factor, a shared trend, or chance. Even strong correlation does not choose among them.

An experiment, natural experiment, considered causal design, or reliable domain source may support a causal claim. One chart and one language model cannot.

Causal words such as “because”, “caused”, “led to”, and “thanks to” are useful warning signals, but a word list is not a safety boundary: the same meaning can be phrased differently. The safe policy is not to auto-publish free-form causal prose.

## Treat uncertainty as a state

A useful system needs at least three outcomes:

- `accept`: a narrow observation is rendered from verified data;
- `review`: the claim may be useful but needs a person or another source;
- `reject`: its structure, type, or evidence reference is invalid.

`review` is not disguised approval. The text remains unpublished until a person decides.

## Do not ask the model to certify itself

A second prompt—“check whether you invented anything”—may catch some errors, but it remains model output. Use it as an extra signal, never as independent verification.

Independent checks rely on something else: a source record, executable calculation, external document, or accountable human reviewer.

## Lesson map

![Lesson map: observation, calculation, cause, and decision](/static/course/py/map-trust-en.svg)

Reconstruct the cue: **observation → source; calculation → formula; cause → independent evidence and a person**.

## Say it in your own words

1. Why does a correct number not make its causal explanation correct?
2. How does a calculation differ from an observation?
3. Why can we not trust the model's own `claim_type` unconditionally?
4. What does `review` mean for publication?

## Warm-up

**1. Predict.** What type is “the value rose because of demand”?

<!-- drill 1 -->
```python
text = "the value rose because of demand"
print("cause" if "because" in text else "observation")
```

**2. Fill the gap.** Permit only three known types.

```python
schema = {"type": "string", ...: ["observation", "calculation", "cause"]}
```

**3. Mend it.** A causal claim is currently published automatically.

```python
if candidate["evidence_ids"]:
    publish(candidate["text"])
```

## Exercise

**Required.** Write `route(candidates, known_ids)`. An observation with known evidence receives `accept`; a cause receives `review`; a calculation without `formula_id` receives `review`; an unknown source or type receives `reject`. Count each outcome.

<!-- task out -->
```text
accept: 1
review: 2
reject: 2
published: ['obs-1']
```

**With your own data.** Write one observation, one calculation, and one causal hypothesis about the digest. State the evidence each requires.

**Optional.** Add a review queue containing the text, routing reason, evidence identifiers, and editor's decision. Never publish an item in `review`.

## Where it enters the project

The AI module is now closed: requests are reproducible, JSON is validated, numbers are checked, and causal prose cannot pass automatically. Lesson 52 begins engineering reliability and turns these rules into `pytest` tests.

## Answers

1. A number supports an observation, not the mechanism that produced it.
2. An observation lives in a source; a calculation must be reproducible from a named formula.
3. Classification is itself unverified model output.
4. The text stays outside publication pending a separate decision.

<!-- drill 1 out -->
```text
cause
```

2. The key is `"enum"`.

<!-- drill 2 -->
```python
schema = {"type": "string", "enum": ["observation", "calculation", "cause"]}
print("cause" in schema["enum"])
```

<!-- drill 2 out -->
```text
True
```

3. Route the cause to review.

<!-- drill 3 -->
```python
candidate = {"claim_type": "cause", "text": "Demand caused growth."}
decision = "review" if candidate["claim_type"] == "cause" else "accept"
print(decision)
```

<!-- drill 3 out -->
```text
review
```

## Sources

- [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework) — managing validity and AI risk.
- [OWASP: misinformation](https://genai.owasp.org/llmrisk/llm092025-misinformation/) — plausible false output and independent verification.
- [The Book of Why](https://www.basicbooks.com/titles/judea-pearl/the-book-of-why/9780465097616/) — the distinction between observation and causal inference.
