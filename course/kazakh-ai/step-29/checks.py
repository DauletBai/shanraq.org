import json

with open("examples.json", encoding="utf-8") as file:
    rows = json.load(file)
weights = {}
club_forms = {"шахмат": "шахмат", "шахматтың": "шахмат",
              "сурет": "сурет", "суреттің": "сурет"}
for row in rows:
    if row["split"] == "train":
        words = row["text"].lower().replace("?", "").split()
        for word in words:
            for feature in [word, word[:4]]:
                if feature not in weights:
                    weights[feature] = {"уақыт": 0, "орын": 0}
                weights[feature][row["label"]] += 1


def classify(question):
    text = question.lower().replace("?", "").replace(",", "").replace(".", "")
    words = text.split()
    allowed = set(club_forms) | {"қашан", "қайда", "уақыты", "орны", "өтеді", "болады"}
    for word in words:
        if word not in allowed:
            return {"status": "refuse", "reason": "білмеймін: таныс емес сөз"}
    words = [club_forms.get(word, word) for word in words]
    clubs = []
    for word in words:
        if word in ["шахмат", "сурет"]:
            clubs.append(word)
    if len(clubs) != 1:
        return {"status": "refuse", "reason": "нақтылаңыз: бір үйірме керек"}
    scores = {"уақыт": 0, "орын": 0}
    for word in words:
        for feature in [word, word[:4]]:
            if feature in weights:
                scores["уақыт"] += weights[feature]["уақыт"]
                scores["орын"] += weights[feature]["орын"]
    time_cue = "қашан" in words or "уақыты" in words
    place_cue = "қайда" in words or "орны" in words
    if time_cue and not place_cue and scores["уақыт"] > scores["орын"]:
        return {"status": "ready", "club": clubs[0], "kind": "уақыт"}
    if place_cue and not time_cue and scores["орын"] > scores["уақыт"]:
        return {"status": "ready", "club": clubs[0], "kind": "орын"}
    return {"status": "refuse", "reason": "нақтылаңыз: бір сұрақ түрі керек"}


if __name__ == "__main__":
    for question in ["Шахмат қашан?", "Шахматтың уақыты қашан?", "Шахмаат қашан?", "Шахмат неге?",
                     "Шахмат қашан интернет?", "Шахмат қашан қайда?"]:
        print(question, "→", classify(question))
