text = input("Сұрақ: ").lower()
clean = text.replace("?", "").replace(",", "")
words = clean.split()
intents = []
if "қашан" in words:
    intents.append("уақыт")
if "қайда" in words:
    intents.append("орын")
print(intents)
