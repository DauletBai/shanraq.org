"""The picture: the report as a line per country.

A report is read by somebody who has thirty seconds. A table of fifteen numbers
needs all thirty; a chart needs two, and the table stays underneath for whoever
wants the rest.

No window is opened. The digest runs on a schedule, where there is no screen to
open one on, so the backend is chosen before pyplot is imported and the figure
goes straight to a file.
"""

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt

COLOURS = ["#b03a2e", "#2e6da4", "#7f8c8d"]


def draw(table, path, names=None):
    """Draws one line per country and saves the picture beside the report."""
    assert not table.empty, "кестеде бірде-бір жол жоқ"
    path.parent.mkdir(parents=True, exist_ok=True)

    fig, ax = plt.subplots(figsize=(8, 4))
    for number, (code, part) in enumerate(table.groupby("country", observed=True)):
        label = (names or {}).get(code, code)
        ax.plot(part["year"], part["value"], marker="o", linewidth=1.6,
                color=COLOURS[number % len(COLOURS)], label=label)

    ax.set_title("Инфляция, жылына %")
    ax.set_xlabel("жыл")
    ax.set_ylabel("%")
    ax.grid(True, linewidth=0.4, alpha=0.5)
    ax.legend()
    # The years are whole numbers, and a tick between two of them means nothing.
    ax.set_xticks(sorted(table["year"].unique()))

    fig.savefig(path, dpi=150, bbox_inches="tight")
    # A figure nobody closes stays in memory, and a program that draws one a day
    # finds that out on the twenty-first day.
    plt.close(fig)
    return path
