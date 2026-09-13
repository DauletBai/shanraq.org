"""Задание урока 46: паспорт модели и граница, за которую она не отвечает."""

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
    """Прямая по логарифму индекса и границы, внутри которых она чему-то училась."""
    years = series.index.to_numpy()
    model = LinearRegression().fit(((years - 2010) / 10).reshape(-1, 1),
                                   np.log(series.to_numpy()))
    known = (int(years.min()), int(years.max()))

    def ask(year, outside=False):
        """Ответ модели. За краем данных — только если прямо попросили."""
        if not outside and not known[0] <= year <= known[1]:
            return None
        return float(np.exp(model.predict(np.array([[(year - 2010) / 10]])))[0])

    return ask, known, model


ask, known, model = trained_on(KAZAKHSTAN[KAZAKHSTAN.index <= CUT])
check = KAZAKHSTAN[KAZAKHSTAN.index > CUT]
errors = (check - np.array([ask(y, outside=True) for y in check.index])).abs()

print("== Паспорт модели")
print(f"  училась: {known[0]}-{known[1]}, точек {known[1] - known[0] + 1}")
print(f"  проверялась: {int(check.index.min())}-{int(check.index.max())},"
      f" ошибка от {errors.min():.0f} до {errors.max():.0f} пунктов")
print(f"  признаков на входе: {model.n_features_in_}")
print(f"  область применения: {known[0]}-{known[1]} — дальше начинается догадка")

print()
print("== Ответы внутри и снаружи")
for year in (2015, 2021, 2025, 2030):
    answer = ask(year)
    if answer is None:
        print(f"  {year}: за краем данных, модель видела {known[0]}-{known[1]}")
    else:
        print(f"  {year}: {answer:.1f}")

print()
print("== Что было бы без границы")
for year in (2025, 2030, 2100):
    print(f"  {year}: {ask(year, outside=True):.1f}")
print("  те же числа, только теперь видно, что это экстраполяция")

print()
print("== Честная строка для отчёта")
print(f"  индекс цен, модель обучена на {known[0]}-{known[1]},"
      f" ошибка на проверке {errors.min():.0f}-{errors.max():.0f} пунктов,"
      f" за {known[1]} годом не применима")
