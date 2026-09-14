# A model on our data: verify every number

_Лид (summary):_ **Lesson fifty of the Python course. Give the model a small table with stable source identifiers, accept structured claims, and compare every number with its source row. The model proposes wording; our code inserts the verified value, unit, and year into the final text.**

## Why this matters

Lesson 49 checked JSON structure. A response containing `{"value": 17.9}` may satisfy the schema perfectly and still be invented. A schema knows the number's type, not our table.

Each fact therefore gets a stable `source_id`. The model returns that identifier with its value. Code locates the source row and compares its value, year, and unit. One mismatch rejects the claim.

An even safer boundary keeps digits out of model-written prose. The model supplies a verbal label; the program formats trusted numbers.

## The whole thing first

The file `grounded.py` extends the previous lesson's contract.

```python
"""Lesson 50: verify a model response against source data."""

import json
from decimal import Decimal

from jsonschema import Draft202012Validator

SOURCES = {
    "inflation-kz-2025": {
        "year": 2025,
        "value": Decimal("11.4"),
        "unit": "percent",
    }
}

CLAIM_SCHEMA = {
    "type": "object",
    "properties": {
        "source_id": {"type": "string"},
        "year": {"type": "integer"},
        "value": {"type": "number"},
        "unit": {"type": "string"},
        "label": {"type": "string", "minLength": 1},
    },
    "required": ["source_id", "year", "value", "unit", "label"],
    "additionalProperties": False,
}
VALIDATOR = Draft202012Validator(CLAIM_SCHEMA)


def verify_claim(raw):
    """Check shape, provenance, and every numeric field."""
    claim = json.loads(raw, parse_float=Decimal)
    VALIDATOR.validate(claim)
    source = SOURCES.get(claim["source_id"])
    if source is None:
        raise ValueError("unknown source")
    for field in ("year", "value", "unit"):
        if claim[field] != source[field]:
            raise ValueError(f"{field} does not match the source")
    if any(char.isdigit() for char in claim["label"]):
        raise ValueError("the model put a number in the label")
    return claim, source


raw = json.dumps({
    "source_id": "inflation-kz-2025",
    "year": 2025,
    "value": 11.4,
    "unit": "percent",
    "label": "Inflation in Kazakhstan",
})
claim, source = verify_claim(raw)
text = f'{claim["label"]}: {source["value"]} % ({source["year"]}).'

print("source:", claim["source_id"])
print("value verified:", claim["value"] == source["value"])
print("unit verified:", claim["unit"] == source["unit"])
print("result:", text)
```

Output:

```text
source: inflation-kz-2025
value verified: True
unit verified: True
result: Inflation in Kazakhstan: 11.4 % (2025).
```

## Give data an address

A `source_id` is neither a link to an entire website nor a row number that changes after sorting. It is a stable identifier for one record. It should lead to the dataset version, observation date, and unit.

Send the model only rows required for the task. This reduces noise, memory use, and opportunities to choose the wrong value. Preserve the dataset version or fingerprint beside the request.

## Why `Decimal`

JSON normally turns a fraction into `float`. Binary fractions are awkward for money and exact comparisons. `json.loads(..., parse_float=Decimal)` reads `11.4` as a decimal value that compares exactly with `Decimal("11.4")`.

That does not make the value true. It only removes an accidental representation error; agreement with the chosen source establishes provenance.

## Words and numbers have different jobs

If the model returns `label="Inflation rose by 15 percent"`, it has hidden a new number in prose. The deliberately strict `isdigit()` rule rejects it. The model names the measure; code prints its value.

Checking digits is not enough for editorial prose: “inflation fell” may be false without containing a number. This lesson therefore permits a neutral label. A later program can calculate “rose” from two verified values.

## Reject; do not silently repair

Do not silently replace the model's wrong value with the right one. Record the rejection and request another response or send it to a person. Otherwise the audit record and publication describe different outputs.

Useful rejection reasons are an unknown source, a mismatched value/year/unit, a digit inside the label, and invalid JSON structure.

## Lesson map

![Lesson map: source, claim, verification, and safe text](/static/course/py/map-grounded-en.svg)

Reconstruct the cue: **source → `source_id` → schema → compare fields → code inserts the number**.

## Say it in your own words

1. Why is a valid JSON Schema insufficient to verify a fact?
2. Why keep `source_id` when the response already contains a value?
3. Why does code, rather than the model, print the final number?
4. What happens to a response whose value does not match?

## Warm-up

**1. Predict.** Are these ordinary `float` values equal?

<!-- drill 1 -->
```python
print(0.1 + 0.2 == 0.3)
```

**2. Fill the gap.** Parse JSON fractions as `Decimal`.

```python
claim = json.loads(raw, ...=Decimal)
```

**3. Mend it.** This code trusts the model's number without finding its source.

```python
claim = json.loads(raw)
text = f'Inflation: {claim["value"]} %'
```

## Exercise

**Required.** Write `accept_many(raw_items, sources)`. Accept a claim only when its `source_id` exists and its `year`, `value`, and `unit` all match; its label must contain no digits. Return accepted identifiers and rejection counts by reason.

<!-- task out -->
```text
accepted: ['inflation-kz-2025']
unknown source: 1
value mismatch: 1
number in label: 1
```

**With your own data.** Give two digest rows stable identifiers and ask the model for neutral labels. Include one deliberately wrong response and demonstrate its rejection.

**Optional.** Add a `dataset_version` and file SHA-256 to each source. Reject claims that name another version.

## Where it enters the project

The model may now suggest a label, but it cannot silently replace the number. Lesson 51 examines a harder boundary: plausible causal explanations absent from the data.

## Answers

1. A schema checks shape and types, not agreement with our table.
2. It identifies the exact record and version used for comparison.
3. The publication then receives the trusted source value, not a model-typed copy.
4. Reject it, record the precise reason, and do not pass it onward.

<!-- drill 1 out -->
```text
False
```

2. The argument is `parse_float`.

<!-- drill 2 -->
```python
import json
from decimal import Decimal

raw = '{"value": 11.4}'
claim = json.loads(raw, parse_float=Decimal)
print(claim["value"] == Decimal("11.4"))
```

<!-- drill 2 out -->
```text
True
```

3. Look up the source, compare fields, and print the source value.

<!-- drill 3 -->
```python
claim = {"source_id": "cpi-2025", "value": 11.4}
sources = {"cpi-2025": {"value": 11.4}}
source = sources[claim["source_id"]]
assert claim["value"] == source["value"]
print(f'Inflation: {source["value"]} %')
```

<!-- drill 3 out -->
```text
Inflation: 11.4 %
```

## Sources

- [Python `json`](https://docs.python.org/3/library/json.html) — `parse_float` and JSON parsing.
- [Python `decimal`](https://docs.python.org/3/library/decimal.html) — exact decimal arithmetic.
- [OWASP: validating LLM output](https://genai.owasp.org/llmrisk/llm052025-improper-output-handling/) — validating output before passing it to another component.
