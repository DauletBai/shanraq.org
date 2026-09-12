"""37-сабақтың тапсырмасы: қатарға бір өлшем, әрі оны көзбен емес, ереже таңдайды."""

import pandas as pd

kazakhstan = pd.Series(
    [6.85, 6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
)
year_2022 = pd.Series(
    {
        "Қазақстан": 15.03, "Ресей": 13.74, "Өзбекстан": 11.45,
        "Қырғызстан": 13.92, "Түркия": 72.31, "Грузия": 11.90,
        "Армения": 8.64, "Әзербайжан": 13.85, "Беларусь": 15.21,
        "Молдова": 28.74,
    },
)


def summary(row):
    """Қатарды сипаттайтын алты сан, бірде-бір артық сан жоқ."""
    quarters = row.quantile([0.25, 0.75])
    return {
        "бақылау": int(row.count()),
        "орташа": round(row.mean(), 2),
        "медиана": round(row.median(), 2),
        "шашырау": round(row.std(), 2),
        "ауқым": (row.min(), row.max()),
        "ширектер": round(quarters[0.75] - quarters[0.25], 2),
    }


def headline(row):
    """Талғам емес, ереже: орташа медианадан оның оннан бірінен артық
    ұзаса, қатардың құйрығы бар — тақырыпқа медиана барады."""
    mean, median = row.mean(), row.median()
    if abs(mean - median) / median > 0.1:
        return "медиана", round(median, 2)
    return "орташа", round(mean, 2)


# Үшінші қатар — сол жылдар, бірақ шыңдарсыз: ереже "орташа" деп те айта білуге
# тиіс, әйтпесе бұл ереже емес, киім ауыстырған әдет.
calm = kazakhstan.loc[2017:2021]

for title, row in (
    ("Қазақстан, 2014-2024", kazakhstan),
    ("Қазақстан, 2017-2021", calm),
    ("Он ел, 2022", year_2022),
):
    figures = summary(row)
    name, value = headline(row)
    print("==", title)
    print("  бақылау саны:", figures["бақылау"])
    print("  орташа:", figures["орташа"], "| медиана:", figures["медиана"])
    print(f"  шашырау: {figures['шашырау']} | {figures['ауқым'][0]}-тен {figures['ауқым'][1]}-ке")
    print("  ширекаралық ауқым:", figures["ширектер"])
    print("  тақырыпқа:", name, value, "%")
    print("  орташадан төмен:", int((row < row.mean()).sum()), "/", len(row))
    print()

print("== Неге бәріне бір өлшем емес")
both = pd.DataFrame({
    "11 жыл": kazakhstan.describe(),
    "5 тыныш": calm.describe(),
    "10 ел": year_2022.describe(),
})
print(both.round(2).to_string())
