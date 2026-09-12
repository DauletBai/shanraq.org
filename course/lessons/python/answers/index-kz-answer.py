"""38-сабақтың тапсырмасы: индекс, көбейткіш және матрас астындағы мың теңге."""

import pandas as pd

rates = pd.Series(
    [6.68, 14.36, 7.44, 6.16, 5.33, 6.72, 8.04, 15.03, 14.53, 8.69],
    index=[2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024],
)


def price_index(rates, base=100.0):
    """Жылдық мөлшерлемелерден баға индексі: пайыздар көбейткішке айналады."""
    return ((1 + rates / 100).cumprod() * base).round(2)


def deflate(nominal, index, base=100.0):
    """Номиналды сомалар — базалық жылдың ақшасымен."""
    return (nominal / index[nominal.index] * base).round(2)


index = price_index(rates)
total = float((1 + rates / 100).prod())

print("== Индекс, 2014 = 100")
print(index.to_string())

print()
print("== Он жылдың қорытындысы")
print("  көбейткіш:", round(total, 3))
print("  өсім:", round((total - 1) * 100, 2), "%")
print("  мөлшерлемелер қосындысы (былай санауға болмайды):", round(rates.sum(), 2), "%")
print("  жылдық орташа:", round((total ** (1 / len(rates)) - 1) * 100, 2), "%")

print()
print("== 2014 жылғы мың теңге")
print("  сол кездегідей сатып алады:", round(1000 / total, 2), "теңге")
print("  бұрынғыны алу үшін керек:", round(1000 * total, 2), "теңге")

print()
print("== Үш есе өскен жалақы")
pay = pd.Series({2016: 150000, 2020: 250000, 2024: 450000}, name="номинал")
real = deflate(pay, index)
table = pd.DataFrame({"номинал": pay, "2014 теңгесімен": real})
table["2016-ға нақты, %"] = (real / real[2016] * 100 - 100).round(2)
print(table.to_string())
print("  номинал өсті:", round(pay[2024] / pay[2016], 2), "есе")
print("  нақты ақшамен:", round(real[2024] / real[2016], 2), "есе")
