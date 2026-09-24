import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
counts = {"уақыт": 0, "орын": 0, "белгісіз": 0}
for row in rows:
    label = row["label"]
    if label in counts:
        counts[label] += 1
    else:
        print("Қате белгі:", label)
print("Уақыт:", counts["уақыт"])
print("Орын:", counts["орын"])
print("Белгісіз:", counts["белгісіз"])
