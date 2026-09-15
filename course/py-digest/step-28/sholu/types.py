"""Typed records at the digest's model boundary."""
from dataclasses import dataclass

@dataclass(frozen=True)
class Observation:
    source_id: str
    year: int
    value: float

def latest(rows: list[Observation]) -> Observation | None:
    """Return the newest observation, admitting an empty input honestly."""
    return max(rows, key=lambda row: row.year) if rows else None

def label(row: Observation | None) -> str:
    """Narrow the optional value before reading its fields."""
    if row is None:
        return "дерек жоқ"
    return f"{row.year}: {row.value:.1f} %"
