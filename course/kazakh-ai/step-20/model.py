import json
from datetime import date

with open("facts.json", encoding="utf-8") as file:
    facts = json.load(file)
question = input("Сұрақ: ").lower().replace("?", "").replace(",", "")
words = question.split()
entities = []
intents = []
for word in words:
    if word in ["шахмат", "сурет"]:
        entities.append(word)
    if word == "қашан":
        intents.append("уақыт")
    if word == "қайда":
        intents.append("орын")
if len(entities) != 1 or len(intents) != 1:
    print("нақтылаңыз: бір үйірме және бір сұрақ түрі керек")
else:
    club = entities[0]
    kind = intents[0]
    matches = []
    for fact in facts:
        if fact["club"] == club and fact["kind"] == kind:
            matches.append(fact)
    if len(matches) == 0:
        print("білмеймін: дерек жоқ")
    elif len(matches) > 1:
        print("тоқта: бірнеше дерек табылды")
    else:
        fact = matches[0]
        today = date.fromisoformat("2026-09-24")
        until = date.fromisoformat(fact["valid_until"])
        if today > until:
            print("білмеймін: дерек ескірген")
        elif kind == "уақыт":
            print(fact["club"], "үйірмесінің уақыты:", fact["value"])
            print("Дереккөз:", fact["source"], "| тексерілген күні:", fact["checked_on"])
        else:
            print(fact["club"], "үйірмесінің орны:", fact["value"])
            print("Дереккөз:", fact["source"], "| тексерілген күні:", fact["checked_on"])
