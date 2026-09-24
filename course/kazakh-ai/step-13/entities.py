known = ["шахмат", "сурет"]
text = input("Сұрақ: ").lower()
clean = text.replace("?", "").replace(",", "")
words = clean.split()
found = []
for word in words:
    if word in known:
        found.append(word)
print(found)
