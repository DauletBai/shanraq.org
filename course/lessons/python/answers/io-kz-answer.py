"""26-сабақтың тапсырмасының шешімі: бухгалтердің экспорты кестеге айналады.

Файл Excel-ден келген: нүктелі үтір, бөлшекте үтір, жетіспейтін мәннің орнына
«дерек жоқ» белгісі және басында BOM. Мұның бәрі — цикл емес, оқу параметрлері.
"""

import sqlite3
from pathlib import Path

import pandas as pd

HERE = Path(__file__).resolve().parent
DIRTY = HERE / "eksport.csv"
CLEAN = HERE / "taza.csv"
DB = HERE / "sholu.db"

# Файл дискіде осылай көрінеді: тапсырма бәрінде бірдей болу үшін оны
# бағдарламаның өзі жазады.
DIRTY.write_text(
    "жыл;санат;сома;ескертпе\n"
    "2024;тамақ;1 200,50;\n"
    "2024;байланыс;4 000,00;тариф ауысты\n"
    "2025;тамақ;1 380,75;\n"
    "2025;байланыс;дерек жоқ;шот келмеді\n",
    encoding="utf-8-sig",
)

df = pd.read_csv(
    DIRTY,
    sep=";",
    decimal=",",
    thousands=" ",
    encoding="utf-8-sig",
    na_values=["дерек жоқ"],
    usecols=["жыл", "санат", "сома"],
)
print("оқылды:")
print(df)
print("түрлері:", dict(df.dtypes.astype(str)))
print("жетіспейтін мән:", int(df["сома"].isna().sum()))

print()
print("таза:")
df.to_csv(CLEAN, index=False, encoding="utf-8")
print(CLEAN.read_text(encoding="utf-8"), end="")

print()
db = sqlite3.connect(DB)
df.to_sql("shygys", db, if_exists="replace", index=False)
answer = pd.read_sql(
    "SELECT жыл, SUM(сома) AS сома FROM shygys GROUP BY жыл ORDER BY жыл", db
)
db.close()
print("жылдар бойынша:")
print(answer.to_string(index=False))
