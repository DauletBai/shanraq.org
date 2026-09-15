"""Lesson 53: find the latest observation with types."""
from dataclasses import dataclass

@dataclass(frozen=True)
class Observation:
    source_id: str
    year: int
    value: float

def latest(rows: list[Observation]) -> Observation | None:
    return max(rows, key=lambda row: row.year) if rows else None

def label(row: Observation | None) -> str:
    return "no data" if row is None else f"{row.year}: {row.value:.1f}%"

print(label(latest([Observation("cpi-kz", 2025, 11.4)])))
print(label(latest([])))
