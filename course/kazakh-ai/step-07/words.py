text = input("Сұрақ: ")
clean = text.lower().replace("?", "").replace(",", "")
words = clean.split()
print(words)
