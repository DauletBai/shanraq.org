"""46-сабақтың тапсырмасы: модель паспорты және ол жауап бермейтін шекара."""

import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression

KAZAKHSTAN = pd.Series({
    2010: 100.00, 2011: 108.45, 2012: 114.09, 2013: 120.87, 2014: 129.15,
    2015: 137.77, 2016: 157.56, 2017: 169.28, 2018: 179.72, 2019: 189.30,
    2020: 202.02, 2021: 218.27, 2022: 251.07, 2023: 287.54, 2024: 312.53,
    2025: 348.12,
})
CUT = 2021


def trained_on(series):
    """Индекс логарифмі бойынша түзу және ол үйренген шекаралар."""
    years = series.index.to_numpy()
    model = LinearRegression().fit(((years - 2010) / 10).reshape(-1, 1),
                                   np.log(series.to_numpy()))
    known = (int(years.min()), int(years.max()))

    def ask(year, outside=False):
        """Модельдің жауабы. Дерек шетінен тыс — тек анық сұрағанда ғана."""
        if not outside and not known[0] <= year <= known[1]:
            return None
        return float(np.exp(model.predict(np.array([[(year - 2010) / 10]])))[0])

    return ask, known, model


ask, known, model = trained_on(KAZAKHSTAN[KAZAKHSTAN.index <= CUT])
check = KAZAKHSTAN[KAZAKHSTAN.index > CUT]
errors = (check - np.array([ask(y, outside=True) for y in check.index])).abs()

print("== Модель паспорты")
print(f"  үйренді: {known[0]}-{known[1]}, нүкте саны {known[1] - known[0] + 1}")
print(f"  тексерілді: {int(check.index.min())}-{int(check.index.max())},"
      f" қатесі {errors.min():.0f}-{errors.max():.0f} тармақ")
print(f"  кірістегі белгі саны: {model.n_features_in_}")
print(f"  қолдану аясы: {known[0]}-{known[1]} — одан әрі болжам басталады")

print()
print("== Іштегі және сырттағы жауаптар")
for year in (2015, 2021, 2025, 2030):
    answer = ask(year)
    if answer is None:
        print(f"  {year}: дерек шетінен тыс, модель {known[0]}-{known[1]} көрген")
    else:
        print(f"  {year}: {answer:.1f}")

print()
print("== Шекара болмаса не болар еді")
for year in (2025, 2030, 2100):
    print(f"  {year}: {ask(year, outside=True):.1f}")
print("  сол сандар, тек енді бұл экстраполяция екені көрініп тұр")

print()
print("== Есепке арналған адал жол")
print(f"  баға индексі, модель {known[0]}-{known[1]} жылдарда үйретілген,"
      f" тексерудегі қатесі {errors.min():.0f}-{errors.max():.0f} тармақ,"
      f" {known[1]} жылдан кейін қолданылмайды")
