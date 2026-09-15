"""Задание урока 53: типизированный поиск последнего наблюдения."""

from dataclasses import dataclass


@dataclass(frozen=True)
class Observation:
    source_id: str
    year: int
    value: float


def latest(rows: list[Observation]) -> Observation | None:
    return max(rows, key=lambda row: row.year) if rows else None


def label(row: Observation | None) -> str:
    if row is None:
        return "нет данных"
    return f"{row.year}: {row.value:.1f} %"


print(label(latest([Observation("cpi-kz", 2025, 11.4)])))
print(label(latest([])))
