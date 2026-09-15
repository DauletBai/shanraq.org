from sholu.types import Observation, label, latest

def test_latest_and_empty_input():
    rows = [Observation("cpi", 2024, 8.7), Observation("cpi", 2025, 11.4)]
    assert label(latest(rows)) == "2025: 11.4 %"
    assert label(latest([])) == "дерек жоқ"
