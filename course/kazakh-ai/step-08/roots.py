roots = {"үй": "баспана", "мектеп": "оқу орны"}
word = input("Негіз: ").lower()
if word in roots:
    print(roots[word])
else:
    print("білмеймін")
