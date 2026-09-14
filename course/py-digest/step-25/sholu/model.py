"""The narrow boundary between a model response and the digest."""

import json
from decimal import Decimal

from jsonschema import Draft202012Validator

SCHEMA = {
    "type": "object",
    "properties": {
        "source_id": {"type": "string"},
        "year": {"type": "integer"},
        "value": {"type": "number"},
        "unit": {"const": "percent"},
        "label": {"type": "string", "minLength": 1},
    },
    "required": ["source_id", "year", "value", "unit", "label"],
    "additionalProperties": False,
}
VALIDATOR = Draft202012Validator(SCHEMA)


def source_from(report, country="KAZ"):
    """Make one stable source record from the report's latest observation."""
    rows = report.loc[report["country"] == country]
    if len(rows) != 1:
        raise ValueError(f"дереккөз бір жол болуы керек: {country}")
    row = rows.iloc[0]
    year = int(row["соңғы_жыл"])
    return {
        "source_id": f"inflation-{country}-{year}",
        "year": year,
        "value": Decimal(str(row["соңғы_мән"])),
        "unit": "percent",
    }


def sample_response(source):
    """A reproducible stand-in; a later lesson may replace only this function."""
    return json.dumps({
        **source,
        "value": float(source["value"]),
        "label": "Қазақстандағы инфляция",
    })


def verify(raw, source):
    """Return a sentence only after every source-owned field agrees."""
    claim = json.loads(raw, parse_float=Decimal)
    VALIDATOR.validate(claim)
    if claim["source_id"] != source["source_id"]:
        raise ValueError("модель белгісіз дереккөзді атады")
    for field in ("year", "value", "unit"):
        if claim[field] != source[field]:
            raise ValueError(f"модельдің {field} өрісі дереккөзге сәйкес емес")
    if any(char.isdigit() for char in claim["label"]):
        raise ValueError("модель атауға сан қосты")
    return (f'{claim["label"]}: {source["value"]} % ({source["year"]}). '
            f'[{source["source_id"]}]')
