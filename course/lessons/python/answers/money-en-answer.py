"""Lesson 40 task: money year by year, checked against a second source."""

import pandas as pd

# The National Bank, end of year, millions of tenge: cash, M3 and the base.
NBK = pd.DataFrame(
    [
        (2019, 2_300_505, 21_322_070, 6_893_176),
        (2020, 2_828_046, 24_917_785, 9_777_551),
        (2021, 2_997_723, 30_099_291, 10_957_714),
        (2022, 3_360_653, 34_295_955, 11_874_422),
        (2023, 3_639_662, 38_301_572, 11_539_973),
        (2024, 4_374_546, 45_660_003, 14_606_537),
        (2025, 4_749_199, 52_751_740, 15_678_721),
    ],
    columns=["year", "M0", "M3", "base"],
).set_index("year")

# The World Bank: GDP and broad money in tenge, the price index at 2010 = 100.
GDP = pd.Series({2019: 69_532_626.5, 2020: 70_649_033.2, 2021: 83_951_587.9,
                 2022: 103_765_518.2, 2023: 119_442_289.7, 2024: 136_693_318.3,
                 2025: 159_561_346.6})
BANK_M3 = pd.Series({2019: 21_322_070.31, 2020: 24_917_784.65, 2021: 30_099_290.60,
                     2022: 34_295_954.60, 2023: 38_301_571.83, 2024: 45_660_003.09,
                     2025: 52_751_740.08})
PRICES = pd.Series({2019: 189.30, 2020: 202.02, 2021: 218.27, 2022: 251.07,
                    2023: 287.54, 2024: 312.53, 2025: 348.12})


def count(nbk, gdp):
    """Year by year: the share of cash, money per tenge of GDP, the multiplier."""
    out = pd.DataFrame(index=nbk.index)
    out["M3"] = nbk["M3"]
    out["cash"] = (nbk["M0"] / nbk["M3"] * 100).round(1)
    out["per tenge of GDP"] = (nbk["M3"] / gdp).round(3)
    out["multiplier"] = (nbk["M3"] / nbk["base"]).round(3)
    return out


table = count(NBK, GDP)
print("== Money year by year")
print(table.to_string())

first, last = NBK.index[0], NBK.index[-1]
print()
print(f"== What grew over {first}-{last}, times")
growth = pd.Series({
    "prices": PRICES[last] / PRICES[first],
    "GDP": GDP[last] / GDP[first],
    "money M3": NBK.loc[last, "M3"] / NBK.loc[first, "M3"],
}).round(3)
print(growth.to_string())
print("  the fastest:", growth.idxmax())

print()
print("== Checked against a second source")
diff = (BANK_M3.round() - NBK["M3"]).abs()
print("  the largest gap, millions of tenge:", int(diff.max()))
print("  matches in every year:", bool((diff == 0).all()))
