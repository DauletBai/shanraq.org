"""Lesson 45 task: the confusion matrix and a threshold chosen by cost."""

import pandas as pd
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import accuracy_score, confusion_matrix

# A country: the year its series starts and the rates themselves, in per cent.
INFLYACIYA = {
    "Azerbaijan": (2000, [1.81, 1.55, 2.77, 2.23, 6.71, 9.68, 8.33, 16.7, 20.85, 1.46,
                           5.73, 7.86, 1.07, 2.42, 1.37, 4.03, 12.44, 12.94, 2.27, 2.61,
                           2.76, 6.65, 13.85, 8.79, 2.21]),
    "Armenia": (2000, [-0.79, 3.15, 1.06, 4.72, 6.96, 0.64, 2.89, 4.41, 8.95, 3.41,
                       8.18, 7.65, 2.56, 5.79, 2.98, 3.73, -1.4, 0.97, 2.52, 1.44,
                       1.21, 7.18, 8.64, 1.98, 0.27]),
    "Belarus": (2000, [168.62, 61.13, 42.54, 28.4, 18.11, 10.34, 7, 8.43, 14.84, 12.95,
                        7.74, 53.23, 59.22, 18.31, 18.12, 13.53, 11.84, 6.03, 4.87, 5.6,
                        5.55, 9.46, 15.21, 5, 5.79]),
    "Georgia": (2000, [4.06, 4.65, 5.59, 0.84, 5.66, 8.25, 9.16, 9.24, 10, 1.73,
                      7.11, 8.54, -0.94, -0.51, 3.07, 4, 2.13, 6.04, 2.62, 4.85,
                      5.2, 9.57, 11.9, 2.49, 1.11]),
    "Kazakhstan": (2000, [13.18, 8.35, 5.84, 6.44, 6.88, 7.58, 8.72, 10.85, 17.14, 7.32,
                         7.4, 8.45, 5.2, 5.94, 6.85, 6.68, 14.36, 7.44, 6.16, 5.33,
                         6.72, 8.04, 15.03, 14.53, 8.69]),
    "Kyrgyzstan": (2000, [18.7, 6.92, 2.13, 2.97, 4.11, 4.34, 5.55, 10.23, 24.52, 6.84,
                          7.97, 16.64, 2.77, 6.61, 7.53, 6.5, 0.39, 3.18, 1.54, 1.13,
                          6.33, 11.91, 13.92, 10.75, 5]),
    "Moldova": (2000, [31.3, 9.76, 5.3, 11.75, 12.52, 11.96, 12.78, 12.37, 12.78, -0.06,
                       7.48, 7.69, 4.55, 4.6, 5.09, 9.68, 6.36, 6.57, 3.05, 4.84,
                       3.77, 5.11, 28.74, 13.42, 4.68]),
    "Russia": (2000, [20.8, 21.48, 15.79, 13.66, 10.89, 12.69, 9.67, 9.01, 14.11, 11.65,
                      6.85, 8.44, 5.07, 6.75, 7.82, 15.53, 7.04, 3.68, 2.88, 4.47,
                      3.38, 6.69, 13.74, 5.87, 8.43]),
    "Turkiye": (2000, [54.92, 54.4, 44.96, 21.6, 8.6, 8.18, 9.6, 8.76, 10.44, 6.25,
                      8.57, 6.47, 8.89, 7.49, 8.85, 7.67, 7.78, 11.14, 16.33, 15.18,
                      12.28, 19.6, 72.31, 53.86, 58.51]),
    "Uzbekistan": (2011, [13.78, 13.21, 11.84, 9.28, 8.75, 8.13, 13.88, 17.52, 14.53,
                          12.87, 10.85, 11.45, 9.96, 9.63]),
}
HIGH = 10.0        # what "high inflation" means here
CUT = 2015         # the last year of training
THRESHOLDS = [0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7]


def pairs(data):
    """Pairs "this year -> the next": the feature, the answer and the year."""
    rows = []
    for country, (first, values) in data.items():
        for i in range(len(values) - 1):
            rows.append({"country": country, "year": first + i,
                         "inflation": values[i], "next": values[i + 1]})
    table = pd.DataFrame(rows)
    table["high"] = (table["next"] >= HIGH).astype(int)
    return table


def matrix(truth, said):
    """The four numbers of the matrix, in an order easy to say in words."""
    tn, fp, fn, tp = confusion_matrix(truth, said, labels=[0, 1]).ravel()
    return {"correct no": tn, "false yes": fp, "missed": fn, "correct yes": tp}


table = pairs(INFLYACIYA)
train = table[table["year"] <= CUT]
test = table[table["year"] > CUT]

print("== The cut by time")
print("  training:", len(train), "pairs, of them high", int(train["high"].sum()))
print("  the check:", len(test), "pairs, of them high", int(test["high"].sum()))
leak = sorted(set(train["year"]) & set(test["year"]))
print("  years in both parts:", leak if leak else "none")

model = LogisticRegression().fit(train[["inflation"]], train["high"])
probability = model.predict_proba(test[["inflation"]])[:, 1]

print()
print("== The matrix at a threshold of 0.5")
cells = matrix(test["high"], (probability >= 0.5).astype(int))
for name, value in cells.items():
    print(f"  {name}: {value}")

print()
print("== Every threshold")
rows = []
for threshold in THRESHOLDS:
    said = (probability >= threshold).astype(int)
    cells = matrix(test["high"], said)
    rows.append({"threshold": threshold,
                 "false yes": cells["false yes"],
                 "missed": cells["missed"],
                 "accuracy": round(accuracy_score(test["high"], said), 3)})
grid = pd.DataFrame(rows).set_index("threshold")
print(grid.to_string())

print()
print("== Two different answers to one question")
# The cost comes from outside: how many false alarms one miss is worth.
for name, miss_costs in (("a miss costs five alarms", 5),
                         ("an alarm costs five misses", 0.2)):
    cost = grid["false yes"] + grid["missed"] * miss_costs
    best = cost.idxmin()
    print(f"  {name}: best threshold {best}"
          f" — false yes {int(grid.loc[best, 'false yes'])},"
          f" missed {int(grid.loc[best, 'missed'])}")
print("  one model, different answers: a threshold is not chosen by accuracy")
