"""Задание урока 37: одна мера на ряд, и она выбрана правилом, а не на глаз."""

import pandas as pd

kazakhstan = pd.Series(
    [6.85, 6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
)
year_2022 = pd.Series(
    {
        "Казахстан": 15.03, "Россия": 13.74, "Узбекистан": 11.45,
        "Кыргызстан": 13.92, "Турция": 72.31, "Грузия": 11.90,
        "Армения": 8.64, "Азербайджан": 13.85, "Беларусь": 15.21,
        "Молдова": 28.74,
    },
)


def summary(row):
    """Шесть чисел, которыми ряд описывают, и ни одним больше."""
    quarters = row.quantile([0.25, 0.75])
    return {
        "наблюдений": int(row.count()),
        "среднее": round(row.mean(), 2),
        "медиана": round(row.median(), 2),
        "разброс": round(row.std(), 2),
        "размах": (row.min(), row.max()),
        "четверти": round(quarters[0.75] - quarters[0.25], 2),
    }


def headline(row):
    """Правило, а не вкус: если среднее ушло от медианы больше чем на десятую
    часть медианы, ряд тянет хвост — в заголовок идёт медиана."""
    mean, median = row.mean(), row.median()
    if abs(mean - median) / median > 0.1:
        return "медиана", round(median, 2)
    return "среднее", round(mean, 2)


# Третий ряд — те же годы, но без всплесков: правило должно уметь сказать и
# "среднее", иначе это не правило, а переодетая привычка.
calm = kazakhstan.loc[2017:2021]

for title, row in (
    ("Казахстан, 2014–2024", kazakhstan),
    ("Казахстан, 2017–2021", calm),
    ("Десять стран, 2022", year_2022),
):
    figures = summary(row)
    name, value = headline(row)
    print("==", title)
    print("  наблюдений:", figures["наблюдений"])
    print("  среднее:", figures["среднее"], "| медиана:", figures["медиана"])
    print("  разброс:", figures["разброс"], "| от", figures["размах"][0], "до", figures["размах"][1])
    print("  межквартильный размах:", figures["четверти"])
    print("  в заголовок:", name, value, "%")
    print("  ниже среднего:", int((row < row.mean()).sum()), "из", len(row))
    print()

print("== Почему не одна мера на всё")
both = pd.DataFrame({
    "11 лет": kazakhstan.describe(),
    "5 спокойных": calm.describe(),
    "10 стран": year_2022.describe(),
})
print(both.round(2).to_string())
