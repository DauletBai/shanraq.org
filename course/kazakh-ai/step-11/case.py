known = {"үйлеріміз": "үй | лер | іміз",
         "мектептеріміз": "мектеп | тер | іміз"}
word = input("Сөз: ").lower()
if word.endswith("ге"):
    base = word.removesuffix("ге")
    case = "барыс"
elif word.endswith("де"):
    base = word.removesuffix("де")
    case = "жатыс"
else:
    base = ""
    case = ""
if base in known:
    print(known[base], "|", case)
else:
    print("білмеймін")
