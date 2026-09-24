import json

with open("facts.json", encoding="utf-8") as file:
    facts = json.load(file)
print("Жазба саны:", len(facts))
for fact in facts:
    print(fact["club"], fact["kind"], fact["value"])
