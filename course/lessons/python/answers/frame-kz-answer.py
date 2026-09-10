"""25-сабақтың тапсырмасының шешімі: нөмірдің орнына қолтаңба.

Әр елдегі инфляцияның шыңы, төрт жылда бағаның неше есе өскені және ең тыныш
жыл — бұлардың бәрі позицияға емес, қолтаңбаға қойылған сұрақ.
"""

import pandas as pd

data = {
    "KZ": [8.0, 15.0, 14.5, 8.7],
    "UZ": [10.8, 11.4, 10.0, 9.6],
    "RU": [6.7, 13.7, 5.9, 8.4],
}
df = pd.DataFrame(data, index=[2021, 2022, 2023, 2024])
df.index.name = "жыл"

print("кесте:")
print(df)

print()
print("елдер бойынша шың:")
for country in df.columns:
    year = df[country].idxmax()
    print(f"  {country}: {year} ({df.loc[year, country]})")

print()
print("төрт жылда баға өсті:")
growth = (1 + df / 100).prod()
for country, times in growth.items():
    print(f"  {country}: {times:.2f} есе")

print()
average = df.mean(axis=1).round(2)
calm = average.idxmin()
print(f"ең тыныш жыл: {calm} (орташа {average[calm]:.2f})")
