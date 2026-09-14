"""Задание урока 50: принять только утверждения, совпавшие с источником."""

from decimal import Decimal

SOURCES = {"inflation-kz-2025": {"year": 2025, "value": Decimal("11.4"), "unit": "percent"}}


def accept_many(items, sources):
    accepted = []
    rejected = {"unknown": 0, "value": 0, "label": 0}
    for claim in items:
        source = sources.get(claim["source_id"])
        if source is None:
            rejected["unknown"] += 1
        elif any(char.isdigit() for char in claim["label"]):
            rejected["label"] += 1
        elif any(claim[field] != source[field] for field in ("year", "value", "unit")):
            rejected["value"] += 1
        else:
            accepted.append(claim["source_id"])
    return accepted, rejected


items = [
    {"source_id": "inflation-kz-2025", "year": 2025, "value": Decimal("11.4"), "unit": "percent", "label": "Инфляция в Казахстане"},
    {"source_id": "missing", "year": 2025, "value": Decimal("11.4"), "unit": "percent", "label": "Инфляция"},
    {"source_id": "inflation-kz-2025", "year": 2025, "value": Decimal("15.0"), "unit": "percent", "label": "Инфляция"},
    {"source_id": "inflation-kz-2025", "year": 2025, "value": Decimal("11.4"), "unit": "percent", "label": "Инфляция 2025 года"},
]
accepted, rejected = accept_many(items, SOURCES)
print("принято:", accepted)
print("неизвестный источник:", rejected["unknown"])
print("не совпало значение:", rejected["value"])
print("число в подписи:", rejected["label"])
