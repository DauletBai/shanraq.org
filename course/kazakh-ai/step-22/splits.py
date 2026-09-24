import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
counts = {"train": 0, "tune": 0, "test": 0}
seen = []
valid = True
for row in rows:
    if row["text"] in seen:
        print("Қайталанған сұрақ:", row["text"])
        valid = False
    seen.append(row["text"])
    split = row["split"]
    if split in counts:
        counts[split] += 1
    else:
        print("Қате бөлік:", split)
        valid = False
print("Оқыту:", counts["train"])
print("Баптау:", counts["tune"])
print("Бақылау:", counts["test"])
print("Жинақ жарамды:", valid)
