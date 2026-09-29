#!/usr/bin/env python3
"""Generate reviewed KZ/EN mathematics lessons from locally authored data."""
from pathlib import Path
import re

from math_localization_data import DATA

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/mathematics"
IMAGE = re.compile(r"!\[[^\]]*\]\((/static/course/mathematics/[^)]+\.svg)\)")

HEADINGS = {
    "kz": {
        "lead": "Лид (summary)", "where": "Картадағы орнымыз",
        "image": "Алдымен таныс бейне", "meaning": "Дәл мағынасы",
        "signal": "Сабақтың тірек сигналы", "worked": "Талданған мысал",
        "faded": "Көмегі азайған мысал", "recall": "Көмексіз жаңғыртыңыз",
        "mistake": "Қатені тауып түзетіңіз", "transfer": "Жаңа жағдайға көшіру",
        "task": "Тапсырма", "next": "Ашық перспектива",
        "contents": "Курс мазмұны",
    },
    "en": {
        "lead": "Lead (summary)", "where": "Where we are on the map",
        "image": "Begin with a familiar image", "meaning": "Precise meaning",
        "signal": "The lesson's support signal", "worked": "Worked example",
        "faded": "Example with a fading prompt", "recall": "Recall without a prompt",
        "mistake": "Find and correct the mistake", "transfer": "Transfer to a new setting",
        "task": "Exercise", "next": "Open perspective",
        "contents": "Course contents",
    },
}


def localized_map(source: str, lang: str) -> str:
    match = IMAGE.search(source)
    if not match:
        raise ValueError("lesson has no support map")
    path = match.group(1)
    if path.endswith("map-full-ru.svg"):
        return path.replace("map-full-ru.svg", f"map-full-{lang}.svg")
    return path[:-4] + f"-{lang}.svg"


def render(stem: str, lang: str, item: dict, source: str) -> str:
    h = HEADINGS[lang]
    signal = (
        "мағына → модель → есептеу → тексеру → өз сөзіңмен түсіндіру"
        if lang == "kz" else
        "meaning → model → calculation → check → explain in your own words"
    )
    if lang == "kz":
        faded = ("Мысалдағы сандарды өзгертіп, шешімді дайын жолға қарамай қайталаңыз. "
                 "Алдымен нәтижені шамалаңыз, содан кейін дәл есептеп, бастапқы шартпен тексеріңіз.")
        recall = ("1. Негізгі ұғымды бір сөйлеммен анықтаңыз.\n"
                  "2. Тірек сигналын жапқан күйде қайта жазыңыз.\n"
                  "3. Мысалды басқа сандармен шешіп, әр қадамның себебін айтыңыз.")
        transfer = ("Осы байланысты өзіңіз өлшей алатын жағдайға қолданыңыз. Шамаларды, "
                    "өлшем бірліктерін және модель жарамды болатын шекараны атаңыз. Жауапты "
                    "екінші тәсілмен немесе кері амалмен тексеріңіз.")
        task_prefix = "**Міндетті.** "
        own = ("\n\n**Өз дерегіңізбен.** Күнделікті өмірден осы құрылымға сай бір мысал құрып, "
               "оны сөзбен, формуламен және тексерумен көрсетіңіз.\n\n"
               "**Қалауыңызша.** Жеті күннен кейін сандарды ауыстырып, тірекке қарамай қайталаңыз.")
    else:
        faded = ("Change the numbers in the example and repeat the solution without copying its steps. "
                 "Estimate first, calculate exactly, and check against the original condition.")
        recall = ("1. Define the main idea in one sentence.\n"
                  "2. Rebuild the support signal from memory.\n"
                  "3. Solve the example with different numbers and explain why every step is valid.")
        transfer = ("Apply the same relationship to something you can measure yourself. Name the "
                    "quantities, units, and the boundary within which the model is valid. Check the "
                    "answer by a second representation or an inverse operation.")
        task_prefix = "**Required.** "
        own = ("\n\n**With your own data.** Create one everyday example with the same structure and "
               "show it in words, a formula, and a check.\n\n"
               "**Optional.** Seven days later, change the numbers and repeat without the support map.")
    map_path = localized_map(source, lang)
    return f"""# {item['title']}

_{h['lead']}:_ **{item['summary']}**

## {h['where']}

{item['where']}

![{item['alt']}]({map_path})

## {h['image']}

{item['image']}

## {h['meaning']}

{item['meaning']}

## {h['signal']}

`{signal}`

{item['signal']}

## {h['worked']}

{item['worked']}

## {h['faded']}

{faded}

## {h['recall']}

{recall}

## {h['mistake']}

{item['mistake']}

## {h['transfer']}

{transfer}

## {h['task']}

{task_prefix}{item['task']}{own}

## {h['next']}

{item['next']}

[{h['contents']}](/course/mathematics?lang={lang})
"""


def main():
    sources = sorted(
        path for path in LESSONS.glob("*.md")
        if not path.stem.endswith(("-kz", "-en"))
    )
    stems = {path.stem for path in sources}
    if stems != set(DATA):
        raise SystemExit(f"localization data mismatch: missing={stems-set(DATA)}, extra={set(DATA)-stems}")
    for path in sources:
        source = path.read_text(encoding="utf-8")
        for lang in ("kz", "en"):
            item = DATA[path.stem][lang]
            output = path.with_name(f"{path.stem}-{lang}.md")
            output.write_text(render(path.stem, lang, item, source), encoding="utf-8")
            print(output.relative_to(ROOT))


if __name__ == "__main__":
    main()
