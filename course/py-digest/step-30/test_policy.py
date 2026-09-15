"""Regression tests for the model trust boundary introduced in steps 25–26."""

from sholu.model import route_claim


def test_known_observation_is_accepted():
    assert route_claim("observation", ["cpi-2025"], {"cpi-2025"}) == "accept"


def test_cause_waits_for_human_review():
    assert route_claim("cause", ["cpi-2025"], {"cpi-2025"}) == "review"


def test_unknown_evidence_is_rejected():
    assert route_claim("observation", ["missing"], {"cpi-2025"}) == "reject"
