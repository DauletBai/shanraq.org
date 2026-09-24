import json

with open("facts.json", encoding="utf-8") as file:
    facts = json.load(file)
club = input("Үйірме: ").lower()
kind = input("Мәлімет түрі: ").lower()
matches = []
for fact in facts:
    if fact["club"] == club and fact["kind"] == kind:
        matches.append(fact)
if len(matches) == 1:
    fact = matches[0]
    if kind == "уақыт":
        print(fact["club"], "үйірмесінің уақыты:", fact["value"])
    elif kind == "орын":
        print(fact["club"], "үйірмесінің орны:", fact["value"])
    print("Дереккөз:", fact["source"], "| тексерілген күні:", fact["checked_on"])
elif len(matches) == 0:
    print("білмеймін: дерек жоқ")
else:
    print("тоқта: бірнеше дерек табылды")
