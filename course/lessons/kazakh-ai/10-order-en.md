# Lesson 10. Check the order of word parts

## Why this matters

On a letter, you write the address and then add the stamp; the familiar procedure has an order. Word parts have positions too. Our example follows stem → plural → possession → direction. **Possession** says “whose?”; **direction** says “toward what?” The Kazakh `үй-лер-іміз-ге` means “to our houses.” Hyphens help us see the parts; the written word is `үйлерімізге`.

## The order and a small check

Take one verified chain: `үй` (“house”) + `лер` (“several”) + `іміз` (“our”) + `ге` (“to”). These four parts make `үйлерімізге`. `үйімізлерге` does not follow this pattern: possession was placed before plurality. Endings can sound different with other words. We check **this** example, not all of Kazakh morphology.

Create `order.py` in UTF-8 and run it.

```python
expected = ["көптік", "тәуелдік", "барыс"]
chosen = ["көптік", "тәуелдік", "барыс"]
if chosen == expected:
    print("үй" + "лер" + "іміз" + "ге")
else:
    print("реті қате")
```

`expected` is the intended order of labels: plural, possession, direction. `chosen` is the proposed order. We saw lists and square brackets in lesson 7; now we create a list ourselves. `==` asks whether the two lists have equal values in the same order. This is a comparison, while one `=` stores a value. Lesson 8 introduced `if` and indentation. Finally we join four **previously checked** text parts; lesson 9 introduced `+` for strings. The output is `үйлерімізге`.

Change `chosen` to `["тәуелдік", "көптік", "барыс"]` and run again. You get `реті қате` (“wrong order”). The program does not yet construct a word from `chosen`: the list only allows or prevents printing one ready example. The next lesson block will inspect the actual parts of an entered word. An unfamiliar word does not become valid just because its parts appear to be ordered.

## Memory map

Meaning of parts → allowed order → compare with `==` → ready example or refusal. The three labels in the list are like labels on successive drawers.

## Task and self-check

Hide the code. Rebuild both lists and explain the difference between `=` and `==`. Then try `["көптік", "барыс", "тәуелдік"]`. **Hint:** positions two and three differ from `expected`. **Answer:** `реті қате`. Restore the correct list and see `үйлерімізге`: it means “to our houses,” not “in our houses.” **Common mistake:** reading `==` as assignment, or forgetting the spaces before `print`. If a form is unknown, the teaching model should say “I don't know” until it has a verified rule.
