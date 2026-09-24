plural = {"қала": "лар", "үй": "лер", "қалам": "дар",
          "тіл": "дер", "кітап": "тар", "мектеп": "тер"}
word = input("Негіз: ").lower()
if word in plural:
    print(word + plural[word])
else:
    print("білмеймін")
