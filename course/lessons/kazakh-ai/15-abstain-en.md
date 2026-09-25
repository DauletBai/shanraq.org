# Lesson 15. Ask for clarification instead of guessing

## Why this matters

A ticket clerk cannot issue a ticket if the buyer names two films and no showing. They ask which one the buyer wants. Our model must also stop when it cannot tell which class or which kind of information is requested. A **refusal** here means “not enough verified information for the next step,” not “your question is bad.”

## Join the previous steps

Before the code, learn its new marks. `text` stores the lowercase question, `clean` removes `?` and `,`, and `words` stores its words. `entities` will hold found club names; `intents` will hold question types. `len(list)` counts items. `> 1` means “more than one.” `[0]` takes the first item because Python numbers positions from zero; an empty list has no such item. We therefore count and check before taking an item.

We still **do not answer** with a time or place: there is no catalog of approved facts yet. `іздеу: шахмат уақыт` means “we may look for the chess time record,” not “the time was found.” Use the two names from lesson 13 and two question types from lesson 14. Create `request.py` in UTF-8.

```python
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
```

The [step file](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-15) matches this example.

`len(entities)` counts list items; `== 0` checks for none. `>` means “greater than,” so `> 1` detects more than one finding. `elif` branches are tried in order and only the first matching one runs. `entities[0]` takes the **first** item; Python counts list positions from zero. We reach this line only after confirming there is exactly one item in both lists. `intents[0]` works the same way. Taking `[0]` from an empty list would cause an error; the preceding checks prevent that. `==` compares while a single `=` stores a value, as in lesson 10.

`Шахмат қашан?` gives `іздеу: шахмат уақыт`. `Робот қашан?` gives `білмеймін: қай үйірме?` because that name is not yet in the dictionary. `Шахмат, сурет қайда?` asks for one club. `Шахмат қашан қайда?` asks for one question type. `Шахмат неге?` gives `білмеймін: уақыт па, орын ба?`. Repeating the same club name also counts as two findings; we will improve that behavior later instead of silently deleting the user's words.

## Memory map

Question → known names + known types → **exactly one** of each? → lookup key; otherwise → say what is missing or needs clarification. Lookup key → facts (from lesson 16), not an answer yet.

![Lesson 15 support map](/static/course/kazakh-ai/map-15-abstain-en.svg)

## Recall and task

Hide the code and recover the four reasons to stop. Run the examples above and match each message to its reason.

**Task:** predict the results for `СУРЕТ қайда?`, `Сурет қашан қайда?`, and an empty line. **Hint:** count both lists before checking conditions. **Answer:** `іздеу: сурет орын`; `нақтылаңыз: бір сұрақты таңдаңыз`; `білмеймін: қай үйірме?`. **Common mistake:** treating `іздеу:` as an answer about a location. The model has only identified two keys for a future search; it has not verified a fact.

## Gate 3: one unambiguous key

Test four cases: one club and one type; no club; two clubs; two types. Before each run, name the contents of both lists and the outcome. Pass when only the first case produces a key and the other three give distinct, clear reasons to stop. If not, return to the list length checks and retry with new names.
