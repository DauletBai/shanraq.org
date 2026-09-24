import json
from datetime import date

from checks import classify


def answer(question, on_date):
    trace = []
    request = classify(question)
    if request["status"] != "ready":
        trace.append("Сұрақ түрі анықталмады")
        return {"text": request["reason"], "trace": trace}
    club = request["club"]
    kind = request["kind"]
    trace.append("Кілт: " + club + " / " + kind)
    with open("facts.json", encoding="utf-8") as file:
        facts = json.load(file)
    matches = []
    for fact in facts:
        if fact["club"] == club and fact["kind"] == kind:
            matches.append(fact)
    trace.append("Жазба саны: " + str(len(matches)))
    if len(matches) == 0:
        return {"text": "білмеймін: дерек жоқ", "trace": trace}
    if len(matches) > 1:
        return {"text": "тоқта: бірнеше дерек табылды", "trace": trace}
    fact = matches[0]
    checked = date.fromisoformat(fact["checked_on"])
    until = date.fromisoformat(fact["valid_until"])
    trace.append("Дереккөз: " + fact["source"])
    trace.append("Тексерілген: " + fact["checked_on"])
    trace.append("Жарамды: " + fact["valid_until"] + " дейін")
    if on_date < checked or on_date > until:
        return {"text": "білмеймін: бұл күнге жарамды дерек жоқ", "trace": trace}
    if kind == "уақыт":
        label = "уақыты"
    else:
        label = "орны"
    message = (fact["club"] + " үйірмесінің " + label + ": " + fact["value"]
               + " | " + on_date.isoformat() + " күнгі дерек"
               + " | дереккөз: " + fact["source"])
    return {"text": message, "trace": trace}


if __name__ == "__main__":
    day = date.fromisoformat("2026-09-24")
    for question in ["Шахмат қашан?", "Сурет қайда?", "Шахмаат қашан?"]:
        result = answer(question, day)
        print(question, "→", result["text"])
        print("Із:", result["trace"])
