"""The answer to lesson 6's exercise: a price list as a dictionary.

A missing item is asked for with get rather than by the square brackets, so the
program answers instead of falling over -- and the answer says outright that
there is no such item.
"""

prices = {"нан": 280, "сүт": 620, "май": 1890, "тұз": 90}

print("жазба саны:", len(prices))
print("сүттің бағасы:", prices["сүт"])
print("қанттың бағасы:", prices.get("қант"))
print("қанттың бағасы әдепкімен:", prices.get("қант", 0), "теңге")
print("қант бар:", "қант" in prices, "| тұз бар:", "тұз" in prices)

prices["қант"] = 450
print("қосқаннан кейін:", len(prices), "жазба, сома:", sum(prices.values()))
print("ең қымбат баға:", max(prices.values()))
print("атаулары:", list(prices.keys()))
