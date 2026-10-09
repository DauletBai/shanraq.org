#!/usr/bin/env python3
"""Build the trilingual Python block from reviewed teaching text and runnable examples."""

import json
from pathlib import Path

from informatics_python_content_a import CONTENT as FIRST
from informatics_python_content_b import CONTENT as SECOND

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "course/lessons/informatics"
LESSONS = {row["number"]: row for row in json.loads(
    (ROOT / "course/informatics/curriculum.json").read_text(encoding="utf-8"))["lessons"]}
CONTENT = {**FIRST, **SECOND}

# These snippets are executed by the release test in the step-03 project directory.
EXAMPLES = {
    27: ('done = False\nprint("Before:", done)\ndone = True\nprint("After:", done)', 'Before: False\nAfter: True'),
    28: ('raw = "2"\ndays = int(raw)\nprint(type(raw).__name__, type(days).__name__)\nprint(days + 1)', 'str int\n3'),
    29: ('days = 2\nlimit = 2\nprint(days + 1)\nprint(days <= limit)', '3\nTrue'),
    30: ('for days in (-1, 0, 2, 3):\n    if days < 0:\n        status = "OVERDUE"\n    elif days <= 2:\n        status = "REMIND"\n    else:\n        status = "NOT_YET"\n    print(days, status)', '-1 OVERDUE\n0 REMIND\n2 REMIND\n3 NOT_YET'),
    31: ('tasks = [{"done": False}, {"done": True}, {"done": True}]\ncount = 0\nfor task in tasks:\n    if task["done"]:\n        count += 1\nprint(count)', '2'),
    32: ('tasks = [{"id": "t-01", "done": False}, {"id": "t-02", "done": True}]\nfor task in tasks:\n    print(task["id"], task["done"])', 't-01 False\nt-02 True'),
    33: ('from datetime import date\ntoday = date.fromisoformat("2026-10-09")\ndue = date.fromisoformat("2026-10-11")\nprint("Әліппе оқу", (due - today).days)', 'Әліппе оқу 2'),
    34: ('from datetime import date\nfrom assistant_core import reminder_status\ntask = {"done": False, "due_date": "2026-10-11"}\nprint(reminder_status(task, date(2026, 10, 9)))', 'REMIND'),
    35: ('import json\nitem = {"title": "Әліппе оқу", "done": False}\nencoded = json.dumps(item, ensure_ascii=False)\nprint(encoded)\nprint(json.loads(encoded)["title"])', '{"title": "Әліппе оқу", "done": false}\nӘліппе оқу'),
    36: ('from assistant_core import parse_date\ntry:\n    parse_date("2026-02-29")\nexcept ValueError:\n    print("INVALID DATE")', 'INVALID DATE'),
    37: ('from datetime import date\nfrom assistant_core import reminder_status\ntask = {"done": False, "due_date": "2026-10-11"}\nactual = reminder_status(task, date(2026, 10, 9))\nassert actual == "REMIND"\nprint("Boundary test passed")', 'Boundary test passed'),
    38: ('from assistant import main\nmain(["reminders", "--today", "2026-10-09"])', 'Today: 2026-10-09 (local calendar date)\nt-01: REMIND\nt-02: DONE\nt-03: NO_DATE'),
}

HEAD = {
    "ru": ("Лид (summary)", "Где мы на карте", "Ситуация и вопрос", "Новые слова без пропусков", "Опорный сигнал", "Разбираем шаг за шагом", "Предскажите и проверьте", "Ожидаемый вывод", "Поймайте ошибку", "Изменение проекта", "Задание и доказательство", "Перенос в новую ситуацию", "Возврат через 1, 7 и 30 дней", "Следующий урок", "Все уроки курса"),
    "kz": ("Лид (summary)", "Картадағы орнымыз", "Жағдай және сұрақ", "Жаңа сөздерді анық түсінейік", "Сабақтың тірек сигналы", "Қадамдап талдаймыз", "Болжаңыз және тексеріңіз", "Күтілетін нәтиже", "Қатені табыңыз", "Жобаға енгізілетін өзгеріс", "Тапсырма және дәлел", "Жаңа жағдайға көшіру", "1, 7 және 30 күннен кейін қайталау", "Келесі сабақ", "Курстың барлық сабақтары"),
    "en": ("Lead (summary)", "Where we are on the map", "Situation and question", "New words without gaps", "The lesson's support signal", "Work through it step by step", "Predict and check", "Expected output", "Catch the error", "Project change", "Task and evidence", "Transfer to a new setting", "Return after 1, 7, and 30 days", "Next lesson", "All course lessons"),
}

SUPPORT = {
    "ru": "Вход → проверка правила → изменение состояния → наблюдаемый результат. Покажите пальцем, на каком шаге программа читает данные, где сравнивает их с договором и где только сообщает результат. Если шага не видно на схеме, найдите его в коде и добавьте в собственную трассу.",
    "kz": "Кіріс → ережені тексеру → күйді өзгерту → бақыланатын нәтиже. Бағдарлама қай жерде деректі оқитынын, келісіммен салыстыратынын және қай жерде тек хабарлайтынын көрсетіңіз. Қадам сызбада көрінбесе, кодтан тауып, өз трассаңызға қосыңыз.",
    "en": "Input → check the rule → change state → observable result. Point to where the program reads data, compares it with the contract, and merely reports the result. If a step is absent from the map, find it in the code and add it to your own trace.",
}
PREDICT = {
    "ru": "Закройте блок с выводом. По строкам проследите значения имён и напишите точный ожидаемый текст, включая регистр и порядок строк. Только затем выполните код из каталога `step-03` командой `python3` и сравните посимвольно. Поменяйте один вход; до запуска запишите новый прогноз. Если результат отличается, найдите первую строку расхождения, а не подгоняйте ответ задним числом.",
    "kz": "Нәтиже блогын жабыңыз. Атаулар мәнін жол сайын бақылап, әріп регистрі мен жол ретін қоса нақты мәтінді болжаңыз. Содан кейін ғана кодты `step-03` каталогында `python3` арқылы орындап, таңба бойынша салыстырыңыз. Бір кірісті өзгертіп, іске қоспай тұрып жаңа болжам жазыңыз. Айырма болса, жауабыңызды кейін ыңғайламай, алғашқы ауытқыған жолды табыңыз.",
    "en": "Cover the output block. Trace names line by line and write the exact predicted text, including case and line order. Only then run the code with `python3` from `step-03` and compare character by character. Change one input and predict again before running. If results differ, find the first divergent line instead of adjusting the prediction afterwards.",
}
REVIEW = {
    "ru": "Через 1 день восстановите правило и один граничный пример без страницы. Через 7 дней объясните ошибку в новом примере однокласснику. Через 30 дней повторите тест проекта, проверьте, что прежние записи и вывод сохранились, и снова перенесите правило на другую задачу.",
    "kz": "1 күннен кейін бетті ашпай ереже мен бір шекара мысалын еске түсіріңіз. 7 күннен кейін жаңа мысалдағы қатені сыныптасыңызға түсіндіріңіз. 30 күннен кейін жоба тестін қайталап, бұрынғы жазба мен нәтиженің сақталғанын тексеріп, ережені өзге міндетке қолданыңыз.",
    "en": "After 1 day, recall the rule and one boundary case without this page. After 7 days, explain a new error example to a classmate. After 30 days, rerun the project test, check earlier records and output still match, and transfer the rule to a different task again.",
}


def render(number, lang):
    row, item = LESSONS[number], CONTENT[number]
    labels = HEAD[lang]
    title = row[f"title_{lang}"]
    stem = f"{number:02d}-{row['id']}"
    suffix = "" if lang == "ru" else f"-{lang}"
    next_row = LESSONS.get(number + 1)
    next_link = (f"[{'Далее' if lang == 'ru' else 'Келесі' if lang == 'kz' else 'Continue'}: {next_row[f'title_{lang}']}](/{'read/' + 'informatics-' + f'{number+1:02d}-' + next_row['id']}?lang={lang})" if next_row and number < 38 else "")
    code, output = EXAMPLES[number]
    location = {
        "ru": f"Это урок {number} из 72 и часть блока Python (уроки 27–38). Опорная схема показывает проверяемые переходы. Продолжаем бумажный договор версии 0.3: пять исходов напоминания и исходные вымышленные задачи.",
        "kz": f"Бұл 72 сабақтың {number}-сі, Python блогының (27–38) бөлігі. Тірек сызба тексерілетін өтулерді көрсетеді. 0.3 қағаз келісіміндегі еске салудың бес нәтижесі мен ойдан алынған бастапқы тапсырмаларды жалғастырамыз.",
        "en": f"This is lesson {number} of 72 in the Python block (27–38). The support map shows verifiable transitions. We continue the 0.3 paper contract: five reminder outcomes and the original fictional tasks.",
    }[lang]
    sections = [f"# {title}", "", f"_{labels[0]}:_ **{item['lead'][lang]}**", "",
        f"## {labels[1]}", "", location, "", f"![{title}](/static/course/informatics/map-{stem}-{lang}.svg)", "",
        f"## {labels[2]}", "", item["hook"][lang], "",
        f"## {labels[3]}", "", item["terms"][lang], "",
        f"## {labels[4]}", "", SUPPORT[lang], "",
        f"## {labels[5]}", "", item["walk"][lang], "",
        f"## {labels[6]}", "", PREDICT[lang], "", "```python", code, "```", "",
        f"## {labels[7]}", "", "```text", output, "```", "",
        f"## {labels[8]}", "", item["trap"][lang], "",
        f"## {labels[9]}", "", item["project"][lang], "",
        f"## {labels[10]}", "", item["task"][lang], "",
        f"## {labels[11]}", "", item["transfer"][lang], "",
        f"## {labels[12]}", "", REVIEW[lang], "",
        f"## {labels[13] if next_link else labels[14]}", "",
        next_link or f"[{'Все уроки' if lang == 'ru' else 'Барлық сабақтар' if lang == 'kz' else 'All lessons'}](/course/informatics?lang={lang})", ""]
    return "\n".join(sections)


def main():
    assert set(CONTENT) == set(EXAMPLES) == set(range(27, 39))
    OUT.mkdir(parents=True, exist_ok=True)
    for number in range(27, 39):
        stem = f"{number:02d}-{LESSONS[number]['id']}"
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            (OUT / f"{stem}{suffix}.md").write_text(render(number, lang), encoding="utf-8")
    print("Wrote 36 manually localized Python-block lessons")


if __name__ == "__main__":
    main()
