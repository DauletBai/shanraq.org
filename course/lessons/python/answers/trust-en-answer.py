"""Lesson 51 exercise: route claims by their evidence needs."""


def route(candidates, known_ids):
    decisions = {"accept": [], "review": [], "reject": []}
    for item in candidates:
        if item["claim_type"] not in {"observation", "calculation", "cause"}:
            decision = "reject"
        elif not set(item["evidence_ids"]).issubset(known_ids):
            decision = "reject"
        elif item["claim_type"] == "cause":
            decision = "review"
        elif item["claim_type"] == "calculation" and not item.get("formula_id"):
            decision = "review"
        else:
            decision = "accept"
        decisions[decision].append(item["id"])
    return decisions


items = [
    {"id": "obs-1", "claim_type": "observation", "evidence_ids": ["cpi-2025"]},
    {"id": "cause-1", "claim_type": "cause", "evidence_ids": ["cpi-2025"]},
    {"id": "calc-1", "claim_type": "calculation", "evidence_ids": ["cpi-2025"]},
    {"id": "obs-2", "claim_type": "observation", "evidence_ids": ["missing"]},
    {"id": "guess-1", "claim_type": "guess", "evidence_ids": ["cpi-2025"]},
]
result = route(items, {"cpi-2025"})
for decision in ("accept", "review", "reject"):
    print(f"{decision}:", len(result[decision]))
print("published:", result["accept"])
