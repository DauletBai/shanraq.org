analyses = {"ат": ["жылқы", "есім"],
            "үй": ["баспана"]}
word = input("Сөз: ").lower()
if word in analyses:
    print(analyses[word])
else:
    print("білмеймін")
