known = ["шахмат", "сурет"]
text = input("Сұрақ: ").lower()
clean = text.replace("?", "").replace(",", "")
words = clean.split()
entities = []
intents = []
for word in words:
    if word in known:
        entities.append(word)
    if word == "қашан":
        intents.append("уақыт")
    if word == "қайда":
        intents.append("орын")
if len(entities) == 0:
    print("білмеймін: қай үйірме?")
elif len(entities) > 1:
    print("нақтылаңыз: бір үйірмені атаңыз")
elif len(intents) == 0:
    print("білмеймін: уақыт па, орын ба?")
elif len(intents) > 1:
    print("нақтылаңыз: бір сұрақты таңдаңыз")
else:
    print("іздеу:", entities[0], intents[0])
