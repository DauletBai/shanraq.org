"""6-қадам -- 17-сабақтан кейін, модульдер мен пакет туралы.

Бесінші қадамдағы бір файл үш модульге бөлінді: желі -- sholu/derekkoz.py-да,
дискі -- sholu/saqtau.py-да, санау мен есеп -- sholu/esep.py-да. Мұнда, main.py
ішінде, тек ретті шақыру қалды: не істелетіні көрініп тұрады, ал қалай
істелетіні -- өз орнында.

Іске қосу:  python3 main.py
"""

from pathlib import Path

from sholu import derekkoz, esep, saqtau

HERE = Path(__file__).parent
STORE = HERE / "data"
SERIES = STORE / "inflation-kz.csv"
REPORT = STORE / "report.csv"


def main():
    series = saqtau.load(SERIES)
    if series:
        print(f"дискіден оқылды: {SERIES.name}")
    else:
        series = derekkoz.fetch("KZ", derekkoz.INFLATION, 2021, 2025)
        saqtau.save(series, SERIES)
        print(f"банктен алынып, сақталды: {SERIES.name}")

    esep.write(series, REPORT)
    print(f"есеп жазылды: {REPORT.name}")
    print(REPORT.read_text(encoding="utf-8"), end="")


if __name__ == "__main__":
    main()
