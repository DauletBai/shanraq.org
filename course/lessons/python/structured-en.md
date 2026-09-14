# A structured response: model JSON and validation

_Лид (summary):_ **Lesson forty-nine of the Python course. Ask the model for JSON governed by a schema instead of free-form prose, parse it with the standard `json` module, and validate it with `jsonschema`. Three separate gates—syntax, structure, and facts—keep a plausible-looking error out of the report.**

## Why this matters

Lesson 48 preserved the exact request. A repeatable request can still produce an answer the program cannot use: the model renames a field, writes a number as a string, or wraps JSON in an explanation.

“Return JSON” is a request. **JSON Schema** is a contract code can check. It defines required fields, their types, permitted values, and whether undeclared fields are allowed.

There are three gates:

1. `json.loads`: is this JSON at all?
2. `jsonschema`: does it have the expected shape?
3. project rules: did the numbers actually come from our source data?

This lesson builds the first two. The next lesson adds the third.

## The whole thing first

Install the pinned version:

```bash
python -m pip install jsonschema==4.26.0
```

The file `structured.py` uses a stored response, so it runs without Ollama.

```python
"""Lesson 49: parse and validate a structured response."""

import json

from jsonschema import Draft202012Validator

OUTPUT_SCHEMA = {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
        "country": {"type": "string", "minLength": 1},
        "year": {"type": "integer", "minimum": 2000, "maximum": 2100},
        "value": {"type": "number"},
        "unit": {"type": "string", "enum": ["percent"]},
    },
    "required": ["country", "year", "value", "unit"],
    "additionalProperties": False,
}


def parse_and_validate(raw):
    """Return only JSON that satisfies our contract."""
    payload = json.loads(raw)
    Draft202012Validator(OUTPUT_SCHEMA).validate(payload)
    return payload


raw = '{"country":"Kazakhstan","year":2025,"value":11.4,"unit":"percent"}'
result = parse_and_validate(raw)

print("country:", result["country"])
print("year is an integer:", isinstance(result["year"], int))
print("value is a number:", isinstance(result["value"], (int, float)))
print("unit:", result["unit"])
print("no extra fields:", set(result) == {"country", "year", "value", "unit"})
```

Output:

```text
country: Kazakhstan
year is an integer: True
value is a number: True
unit: percent
no extra fields: True
```

## Parse first, trust later

`json.loads` turns a JSON string into Python values. A misplaced comma or prose before the opening brace raises `JSONDecodeError`. Do not extract text between braces with a regular expression: that can turn a fragment of a response into an apparently complete document.

It is better to request a clean structured result from the API. An Ollama request can carry the schema in its `format` field:

```python
request["format"] = OUTPUT_SCHEMA
```

This substantially reduces malformed responses; it does not replace validation in your program. The model and the network remain outside the trust boundary.

## What the schema promises

`type: object` requires an object. `required` lists fields without which the response is incomplete. `additionalProperties: False` rejects undeclared fields, so a misspelt `county` cannot hide beside the intended field.

The JSON number `11.4` and JSON string `"11.4"` are different values. A schema does not coerce one into the other and should not guess what the model meant.

`enum` restricts the unit to `percent`, preventing the program from mixing `%`, fractions, and percentage points. A year range catches obvious nonsense; it does not prove that 2025 exists in the source table.

Check the schema itself once when the application starts:

```python
Draft202012Validator.check_schema(OUTPUT_SCHEMA)
```

## Failure is an outcome

At the application boundary, catch only expected failures and do not continue the pipeline with an empty dictionary:

```python
from json import JSONDecodeError
from jsonschema import ValidationError

try:
    result = parse_and_validate(raw)
except JSONDecodeError:
    print("The response is not JSON")
except ValidationError:
    print("The JSON does not match the schema")
else:
    save_for_review(result)
```

The log can record the type of failure without copying all user input. Preserve the rejected response in restricted storage for diagnosis only when retaining its data is permitted.

## What the schema cannot know

The schema confirms that `value` is a number. It cannot know whether `11.4` appeared in the table, whether the model confused two countries, or whether a claimed cause is invented. **A valid shape is not a valid fact.**

Do not feed model output into SQL, HTML, or a shell command merely because it passed the schema. Escaping, parameterised queries, and domain constraints still apply.

## Lesson map

![Lesson map: a string, JSON, a schema, and checked fields](/static/course/py/map-structured-en.svg)

Reconstruct the cue: **raw response → `json.loads` → JSON Schema → data awaiting factual checks**.

## Say it in your own words

1. Why does “return JSON” not replace validation?
2. How does `JSONDecodeError` differ from `ValidationError`?
3. Why use both `required` and `additionalProperties: False`?
4. Why can a number that passed the schema not yet be published?

## Warm-up

**1. Predict.** Will this value satisfy a `number` check?

<!-- drill 1 -->
```python
import json

value = json.loads('{"value": "11.4"}')["value"]
print(isinstance(value, (int, float)))
```

**2. Fill the gap.** Reject fields outside the contract.

```python
schema = {"type": "object", ...: False}
```

**3. Mend it.** A structural failure currently becomes an apparently usable answer.

```python
try:
    result = parse_and_validate(raw)
except Exception:
    result = {}
save_for_review(result)
```

## Exercise

**Required.** Write `check_many(raw_items)`. Count one of three outcomes for each string: `valid`, `invalid_json`, or `invalid_schema`. Catch only `JSONDecodeError` and `ValidationError`; rejected responses must never enter the accepted list.

<!-- task out -->
```text
accepted: 1
not JSON: 1
wrong schema: 2
years: [2025]
```

**With your own data.** Write the result schema for one table in your digest. Name the unit with `enum`, reject extra fields, and construct four responses that exercise each boundary.

**Optional.** Pass the same schema as the Ollama request's `format`. Run ten requests with and without it, using a metric chosen in advance: the share that passes local validation.

## Where it enters the project

A model response can now enter only the list of structurally valid results. Lesson 50 gives the model numbers from our digest and compares every returned number with its source. Only then will text reach the report.

## Answers

1. A model may add prose, omit a field, or encode a number as a string; an instruction is not a code check.
2. The first means invalid JSON syntax; the second means valid JSON with the wrong structure.
3. One catches omissions, while the other catches misspelt or undeclared fields.
4. The schema knows types and ranges, not the provenance or truth of a number.

<!-- drill 1 out -->
```text
False
```

2. The key is `"additionalProperties"`.

<!-- drill 2 -->
```python
schema = {"type": "object", "additionalProperties": False}
print(schema["additionalProperties"])
```

<!-- drill 2 out -->
```text
False
```

3. Catch the expected failure, record rejection, and do not call the next stage.

<!-- drill 3 -->
```python
from jsonschema import ValidationError, validate

raw = {}

try:
    validate(raw, {"type": "object", "required": ["year"]})
except ValidationError:
    print("rejected")
else:
    print("accepted")
```

<!-- drill 3 out -->
```text
rejected
```

## Sources

- [JSON Schema: getting started](https://json-schema.org/learn/getting-started-step-by-step) — objects, required properties, and additional properties.
- [`jsonschema` documentation](https://python-jsonschema.readthedocs.io/en/stable/validate/) — validators and schema checking.
- [Ollama structured outputs](https://docs.ollama.com/capabilities/structured-outputs) — passing JSON Schema in `format`.
- [Python `json`](https://docs.python.org/3/library/json.html) — parsing JSON and `JSONDecodeError`.
