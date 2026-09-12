"""Задание урока 38: индекс, множитель и тысяча под матрасом."""

import pandas as pd

rates = pd.Series(
    [6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
)


def price_index(rates, base=100.0):
    """Индекс цен из годовых ставок: проценты становятся множителями."""
    return ((1 + rates / 100).cumprod() * base).round(2)


def deflate(nominal, index, base=100.0):
    """Номинальные суммы — в деньгах базового года."""
    return (nominal / index[nominal.index] * base).round(2)


index = price_index(rates)
total = float((1 + rates / 100).prod())

print("== Индекс, 2014 = 100")
print(index.to_string())

print()
print("== Итог десяти лет")
print("  множитель:", round(total, 3))
print("  рост:", round((total - 1) * 100, 2), "%")
print("  сумма ставок (так считать нельзя):", round(rates.sum(), 2), "%")
print("  среднегодовая:", round((total ** (1 / len(rates)) - 1) * 100, 2), "%")

print()
print("== Тысяча 2014 года")
print("  покупает столько же, сколько тогда:", round(1000 / total, 2), "тенге")
print("  чтобы купить прежнее, нужно:", round(1000 * total, 2), "тенге")

print()
print("== Зарплата, которая росла втрое")
pay = pd.Series({2016: 150000, 2020: 250000, 2024: 450000}, name="номинал")
real = deflate(pay, index)
table = pd.DataFrame({"номинал": pay, "в тенге 2014": real})
table["к 2016 в реальных, %"] = (real / real[2016] * 100 - 100).round(2)
print(table.to_string())
print("  номинал вырос в", round(pay[2024] / pay[2016], 2), "раза")
print("  в реальных деньгах — в", round(real[2024] / real[2016], 2), "раза")
