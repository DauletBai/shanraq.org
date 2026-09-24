import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
weights = {}
for row in rows:
    if row["split"] == "train":
        words = row["text"].lower().replace("?", "").replace(",", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature not in weights:
                    weights[feature] = {"уақыт": 0, "орын": 0}
                weights[feature][row["label"]] += 1
answered = 0
correct = 0
known = 0
refused = 0
unknown_refused = 0
for row in rows:
    if row["split"] == "test":
        scores = {"уақыт": 0, "орын": 0}
        words = row["text"].lower().replace("?", "").replace(",", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature in weights:
                    scores["уақыт"] += weights[feature]["уақыт"]
                    scores["орын"] += weights[feature]["орын"]
        time_cue = "қашан" in words or "уақыты" in words
        place_cue = "қайда" in words or "орны" in words
        if time_cue and not place_cue and scores["уақыт"] > scores["орын"]:
            prediction = "уақыт"
        elif place_cue and not time_cue and scores["орын"] > scores["уақыт"]:
            prediction = "орын"
        else:
            prediction = "білмеймін"
        if row["label"] != "белгісіз":
            known += 1
        if prediction == "білмеймін":
            refused += 1
            if row["label"] == "белгісіз":
                unknown_refused += 1
        else:
            answered += 1
            if prediction == row["label"]:
                correct += 1
        print(row["text"], "→", prediction, "|", row["label"])
print("Жауап:", answered, "Дұрыс:", correct, "Белгілі:", known)
print("Бас тарту:", refused, "Белгісізге дұрыс бас тарту:", unknown_refused)
if answered == 0:
    print("Дәлдік: анықталмайды")
else:
    print("Дәлдік:", round(correct / answered, 2))
if known == 0:
    print("Толықтық: анықталмайды")
else:
    print("Толықтық:", round(correct / known, 2))
total = answered + refused
if total == 0:
    print("Бас тарту үлесі: анықталмайды")
else:
    print("Бас тарту үлесі:", round(refused / total, 2))
