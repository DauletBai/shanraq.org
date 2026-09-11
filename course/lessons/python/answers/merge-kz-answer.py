"""29-сабақтың тапсырмасының шешімі: кілт бойынша біріктірілген дерек пен анықтама.

Деректің бірде-бір жолы жоғалмайтындай біріктіреміз, атауы табылмағандардың
санын шығарамыз, содан кейін ғана өңір бойынша топтаймыз.
"""

import pandas as pd

data = pd.DataFrame(
    [("KAZ", 2023, 14.5), ("KAZ", 2024, 8.7),
     ("UZB", 2023, 10.0), ("UZB", 2024, 9.6),
     ("RUS", 2023, 5.9), ("RUS", 2024, 8.4),
     ("KGZ", 2023, 10.8), ("KGZ", 2024, 6.3)],
    columns=["code", "year", "value"],
)
names = pd.DataFrame(
    [("KAZ", "Қазақстан", "Орталық Азия"),
     ("UZB", "Өзбекстан", "Орталық Азия"),
     ("KGZ", "Қырғызстан", "Орталық Азия"),
     ("TJK", "Тәжікстан", "Орталық Азия")],
    columns=["code", "name", "region"],
)

# Сол жақта — дерек, ол жоғалмауға тиіс: анықтама толық емес.
joined = data.merge(names, on="code", how="left", validate="many_to_one")
print("жол болатын:", len(data), "| болды:", len(joined))
print("атауы жоқ:", sorted(joined.loc[joined["name"].isna(), "code"].unique()))

print()
print("елдер бойынша орташа инфляция:")
by_country = joined.groupby(["code", "name"], dropna=False)["value"].mean().round(2)
print(by_country.reset_index().to_string(index=False))

print()
print("өңірлер бойынша:")
# Өңірі жоқ жолдарды топтауға алмаймыз: олардың өңірі бос емес, белгісіз, ал
# оларды «Орталық Азияға» жазып қою — дерек ойлап табу болар еді.
by_region = joined.dropna(subset=["region"]).groupby("region").agg(
    ел=("code", "nunique"),
    орташа=("value", "mean"),
).round(2)
print(by_region.to_string())
