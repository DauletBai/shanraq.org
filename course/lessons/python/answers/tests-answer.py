"""Задание урока 52: тестируем правила допуска."""

import pytest


def route(claim_type, evidence_ids, known_ids, formula_id=None):
    if claim_type not in {"observation", "calculation", "cause"}:
        return "reject"
    if not set(evidence_ids).issubset(known_ids):
        return "reject"
    if claim_type == "cause" or (claim_type == "calculation" and not formula_id):
        return "review"
    return "accept"


@pytest.mark.parametrize(
    ("claim_type", "evidence", "formula_id", "expected"),
    [
        ("observation", ["cpi-2025"], None, "accept"),
        ("cause", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], "change-v1", "accept"),
        ("guess", ["cpi-2025"], None, "reject"),
        ("observation", ["missing"], None, "reject"),
    ],
)
def test_route(claim_type, evidence, formula_id, expected):
    assert route(claim_type, evidence, {"cpi-2025"}, formula_id) == expected


def percent(value):
    if not 0 <= value <= 100:
        raise ValueError("процент вне диапазона")
    return value


@pytest.mark.parametrize("value", [-1, 101])
def test_percent_rejects_outside_range(value):
    with pytest.raises(ValueError, match="вне диапазона"):
        percent(value)


if __name__ == "__main__":
    known = {"cpi-2025"}
    for claim_type, evidence, formula_id, expected in [
        ("observation", ["cpi-2025"], None, "accept"),
        ("cause", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], "change-v1", "accept"),
        ("guess", ["cpi-2025"], None, "reject"),
        ("observation", ["missing"], None, "reject"),
    ]:
        assert route(claim_type, evidence, known, formula_id) == expected
    for value in (-1, 101):
        with pytest.raises(ValueError):
            percent(value)
    print("8 passed")
