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
for row in rows:
    if row["split"] == "tune":
        scores = {"уақыт": 0, "орын": 0}
        words = row["text"].lower().replace("?", "").replace(",", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature in weights:
                    scores["уақыт"] += weights[feature]["уақыт"]
                    scores["орын"] += weights[feature]["орын"]
        if scores["уақыт"] > scores["орын"]:
            prediction = "уақыт"
        elif scores["орын"] > scores["уақыт"]:
            prediction = "орын"
        else:
            prediction = "білмеймін"
        print(row["text"], "→", prediction)
