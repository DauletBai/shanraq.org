"""Step 8 -- after lesson 26, on reading and writing formats.

The disk stops being written row by row. sholu/saqtau.py saves the table with
to_csv and reads it back with read_csv, so what comes off the disk is the table
itself; sholu/esep.py writes the report the same way.

The bank still answers with JSON over the network, and that part has not
changed: json.load is the right tool for one small answer.

Run:  python3 main.py
"""

from pathlib import Path

from sholu import derekkoz, esep, saqtau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation-kz.csv"
REPORT = STORE / "report.csv"


def main():
    table = saqtau.load(SERIES)
    if table is not None:
        print(f"дискіден оқылды: {SERIES.name}")
    else:
        series = derekkoz.fetch("KZ", derekkoz.INFLATION, 2021, 2025)
        table = esep.frame(series)
        saqtau.save(table, SERIES)
        print(f"банктен алынып, сақталды: {SERIES.name}")

    print(table)

    esep.write(table, REPORT)
    print(f"есеп жазылды: {REPORT.name}")
    print(REPORT.read_text(encoding="utf-8"), end="")


if __name__ == "__main__":
    main()
