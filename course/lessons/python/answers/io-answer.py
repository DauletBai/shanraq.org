"""Решение задания урока 26: выгрузка бухгалтера превращается в таблицу.

Файл пришёл из Excel: точка с запятой, запятая в дробях, метка «н/д» вместо
пропуска и BOM в начале. Всё это — не циклы, а параметры чтения.
"""

import sqlite3
from pathlib import Path

import pandas as pd

HERE = Path(__file__).resolve().parent
DIRTY = HERE / "vygruzka.csv"
CLEAN = HERE / "chisto.csv"
DB = HERE / "svodka.db"

# Так этот файл выглядит на диске: его пишет сама программа, чтобы задание было
# одинаковым у всех.
DIRTY.write_text(
    "год;категория;сумма;примечание\n"
    "2024;еда;1 200,50;\n"
    "2024;связь;4 000,00;тариф сменился\n"
    "2025;еда;1 380,75;\n"
    "2025;связь;н/д;счёт не пришёл\n",
    encoding="utf-8-sig",
)

df = pd.read_csv(
    DIRTY,
    sep=";",
    decimal=",",
    thousands=" ",
    encoding="utf-8-sig",
    na_values=["н/д"],
    usecols=["год", "категория", "сумма"],
)
print("прочитано:")
print(df)
print("типы:", dict(df.dtypes.astype(str)))
print("пропусков:", int(df["сумма"].isna().sum()))

print()
print("чисто:")
df.to_csv(CLEAN, index=False, encoding="utf-8")
print(CLEAN.read_text(encoding="utf-8"), end="")

print()
db = sqlite3.connect(DB)
df.to_sql("rashody", db, if_exists="replace", index=False)
answer = pd.read_sql(
    "SELECT год, SUM(сумма) AS сумма FROM rashody GROUP BY год ORDER BY год", db
)
db.close()
print("по годам:")
print(answer.to_string(index=False))
