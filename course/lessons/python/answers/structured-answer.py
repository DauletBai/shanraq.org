"""Задание урока 49: разделить ошибки JSON и ошибки схемы."""

import json
from json import JSONDecodeError

from jsonschema import Draft202012Validator, ValidationError

SCHEMA = {
    "type": "object",
    "properties": {
        "year": {"type": "integer"},
        "value": {"type": "number"},
        "unit": {"const": "percent"},
    },
    "required": ["year", "value", "unit"],
    "additionalProperties": False,
}
VALIDATOR = Draft202012Validator(SCHEMA)


def check_many(raw_items):
    accepted = []
    counts = {"invalid_json": 0, "invalid_schema": 0}
    for raw in raw_items:
        try:
            payload = json.loads(raw)
            VALIDATOR.validate(payload)
        except JSONDecodeError:
            counts["invalid_json"] += 1
        except ValidationError:
            counts["invalid_schema"] += 1
        else:
            accepted.append(payload)
    return accepted, counts


items = [
    '{"year":2025,"value":11.4,"unit":"percent"}',
    '{"year":2025,"value":11.4,}',
    '{"year":"2025","value":11.4,"unit":"percent"}',
    '{"year":2024,"value":8.7,"unit":"percent","note":"estimate"}',
]
accepted, counts = check_many(items)
print("принято:", len(accepted))
print("не JSON:", counts["invalid_json"])
print("не по схеме:", counts["invalid_schema"])
print("годы:", [item["year"] for item in accepted])
