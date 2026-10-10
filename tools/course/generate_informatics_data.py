#!/usr/bin/env python3
"""Render the hand-localized Data, Tables and Databases block (47–54)."""

import json
from pathlib import Path

from informatics_data_content import CONTENT

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "course/lessons/informatics"
ROWS = {row["number"]: row for row in json.loads(
    (ROOT / "course/informatics/curriculum.json").read_text(encoding="utf-8"))["lessons"]}
PROJECT_URL = "https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/step-05"

HEAD = {
    "ru": ("Лид (summary)", "Где мы на карте", "Ситуация и вопрос", "Новые слова без пропусков", "Опорный сигнал", "Разбираем шаг за шагом", "Предскажите и проверьте", "Ожидаемый вывод", "Поймайте ошибку", "Изменение проекта", "Задание и доказательство", "Перенос в новую ситуацию", "Возврат через 1, 7 и 30 дней", "Следующий урок", "Все уроки курса"),
    "kz": ("Қысқаша мазмұн", "Картадағы орнымыз", "Жағдай және сұрақ", "Жаңа сөздерді анық түсінейік", "Сабақтың тірек сигналы", "Қадамдап талдаймыз", "Болжаңыз және тексеріңіз", "Күтілетін нәтиже", "Қатені табыңыз", "Жобаға енгізілетін өзгеріс", "Тапсырма және дәлел", "Жаңа жағдайға көшіру", "1, 7 және 30 күннен кейін қайталау", "Келесі сабақ", "Курстың барлық сабақтары"),
    "en": ("Lead (summary)", "Where we are on the map", "Situation and question", "New words without gaps", "The lesson's support signal", "Work through it step by step", "Predict and check", "Expected output", "Catch the error", "Project change", "Task and evidence", "Transfer to a new setting", "Return after 1, 7, and 30 days", "Next lesson", "All course lessons"),
}

EXAMPLES = {
    47: ('import csv\nwith open("study_sessions.csv", encoding="utf-8", newline="") as stream:\n    first = next(csv.DictReader(stream))\nprint(first["session_id"], first["task_id"], first["minutes"])', 's-01 t-01 20'),
    48: ('import csv\nwith open("study_sessions.csv", encoding="utf-8", newline="") as stream:\n    total = sum(int(row["minutes"]) for row in csv.DictReader(stream))\nprint(total)', '70'),
    49: ('from data_store import read_sessions\nrows = read_sessions("study_sessions.csv", {"t-01", "t-02", "t-03"})\nprint(len(rows), sum(row[3] for row in rows))', '4 70'),
    50: ('totals = [("t-01", 45), ("t-02", 0), ("t-03", 25)]\nfor task_id, minutes in totals:\n    print(task_id, "#" * (minutes // 5) or "(none)")', 't-01 #########\nt-02 (none)\nt-03 #####'),
    51: ('import sqlite3\ndb = sqlite3.connect(":memory:")\ndb.execute("PRAGMA foreign_keys = ON")\ndb.execute("CREATE TABLE tasks (task_id TEXT PRIMARY KEY)")\ndb.execute("CREATE TABLE sessions (task_id TEXT REFERENCES tasks(task_id))")\nprint(db.execute("PRAGMA foreign_keys").fetchone()[0])\ndb.close()', '1'),
    52: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom data_store import create_database, connect_readonly\nwith TemporaryDirectory() as folder:\n    db = Path(folder) / "test.db"\n    create_database(db, "tasks.json", "study_sessions.csv")\n    with connect_readonly(db) as connection:\n        rows = connection.execute("SELECT minutes FROM study_sessions WHERE task_id = ? ORDER BY observed_on", ("t-01",)).fetchall()\n    print(*(row[0] for row in rows))', '20 25'),
    53: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom data_store import create_database, minutes_report\nwith TemporaryDirectory() as folder:\n    db = Path(folder) / "test.db"\n    create_database(db, "tasks.json", "study_sessions.csv")\n    for task_id, minutes in minutes_report(db):\n        print(task_id, minutes)', 't-01 45\nt-02 0\nt-03 25'),
    54: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom datetime import date\nfrom data_store import create_database, load_document_from_db\nfrom assistant_core import reminder_status\nwith TemporaryDirectory() as folder:\n    db = Path(folder) / "test.db"\n    create_database(db, "tasks.json", "study_sessions.csv")\n    for task in load_document_from_db(db)["tasks"]:\n        print(task["id"], reminder_status(task, date(2026, 10, 9)))', 't-01 REMIND\nt-02 DONE\nt-03 NO_DATE'),
}

INTRO = {
    "ru": "Урок {n} из 72, блок «Данные, таблицы и базы» (47–54). Мы продолжаем помощника 1.1 с тремя вымышленными задачами. В конце блока его записи будут храниться в локальном SQLite, а проверенные правила напоминаний останутся прежними.",
    "kz": "72 сабақтың {n}-сі, «Деректер, кестелер және дерекқорлар» блогы (47–54). Үш ойдан алынған тапсырмасы бар 1.1 көмекшісін жалғастырамыз. Блок соңында жазбалар жергілікті SQLite ішінде сақталып, тексерілген еске салу ережелері өзгермейді.",
    "en": "Lesson {n} of 72, Data, tables and databases block (47–54). We continue assistant 1.1 with three fictional tasks. By the end, records live in local SQLite while tested reminder rules remain unchanged.",
}
SUPPORT = {
    "ru": "Прочитайте схему слева направо. Для каждой стрелки назовите вход, преобразование и проверку. Закройте третий шаг и предскажите его. Затем найдите соответствующую строку в CSV, SQL или отчёте: нарисованная стрелка не должна подменять реальные значения и связи.",
    "kz": "Сызбаны солдан оңға оқыңыз. Әр жебенің кірісін, өзгерісін және тексеруін айтыңыз. Үшінші қадамды жауып, нәтижесін болжаңыз. Содан соң CSV, SQL не есептегі тиісті жолды табыңыз: сурет нақты сан мен байланысты алмастырмауы керек.",
    "en": "Read the map left to right. Name each arrow's input, transformation and check. Cover the third step and predict it. Then locate the matching row in CSV, SQL or the report: an arrow must not replace real values or relationships.",
}
PREDICT = {
    "ru": "Сначала запишите ожидаемый вывод точно, включая порядок строк. Затем выполните код из папки `step-05` и сравните посимвольно. Измените один вход и повторите цикл «предсказать — запустить — объяснить». Пример доказывает только заявленное правило: он не превращает учебные вымышленные данные в измерение реальных людей.",
    "kz": "Алдымен жол ретіне дейін нақты күтілетін нәтижені жазыңыз. Кодты `step-05` бумасынан орындап, таңба бойынша салыстырыңыз. Бір кірісті өзгертіп, «болжау — іске қосу — түсіндіру» айналымын қайталаңыз. Мысал тек көрсетілген ережені дәлелдейді; ойдан алынған деректі нақты адамдар өлшеміне айналдырмайды.",
    "en": "Write the exact expected output, including row order, before running anything. Execute the snippet in `step-05` and compare characters. Change one input, then predict, run and explain again. The example proves only the stated rule; fictional teaching data do not become measurements of real people.",
}
REVIEW = {
    "ru": "Через 1 день восстановите схему по памяти и назовите случай, где простая аналогия не работает. Через 7 дней объясните новую таблицу товарищу без текста и ответьте минимум на 7 из 10 вопросов итогового урока. Через 30 дней повторите импорт во временную папку и проверьте три задачи, четыре занятия, сумму 70 и неизменные напоминания. Запишите оставшийся непонятным термин и вернитесь к его первому определению.",
    "kz": "1 күннен кейін сызбаны жатқа салып, қарапайым ұқсастық жұмыс істемейтін жағдайды атаңыз. 7 күннен кейін жаңа кестені мәтінсіз досыңызға түсіндіріп, соңғы сабақтағы он сұрақтың кемінде жетеуіне жауап беріңіз. 30 күннен кейін уақытша бумаға импорт жасап, үш тапсырма, төрт сабақ, 70 минут және өзгермеген еске салуды тексеріңіз. Түсініксіз терминді жазып, алғашқы анықтамасына оралыңыз.",
    "en": "After 1 day redraw the map from memory and name a case where the simple analogy fails. After 7 days explain a fresh table without the text and answer at least 7 of the final lesson's 10 questions. After 30 days repeat import into a temporary folder and check three tasks, four sessions, total 70 and unchanged reminders. Record any unclear term and return to its first definition.",
}
RESOURCES = {
    47: "https://docs.python.org/3/library/csv.html",
    48: "https://www.libreoffice.org/discover/calc/",
    49: "https://docs.python.org/3/library/csv.html",
    50: "https://service-manual.ons.gov.uk/data-visualisation/guidance/axes-and-gridlines",
    51: "https://www.sqlite.org/foreignkeys.html",
    52: "https://docs.python.org/3/library/sqlite3.html",
    53: "https://www.sqlite.org/lang_select.html",
    54: "https://www.sqlite.org/atomiccommit.html",
}
MASTERY = {
    "ru": ["Чем наблюдение отличается от вывода по памяти?", "Как отличить пустую длительность от нуля минут?", "Почему сумма четырёх записей равна 70, а не 60?", "Зачем хранить исходный CSV перед очисткой?", "Какой ключ запрещает повтор сессии, а какой — ссылку на несуществующую задачу?", "Почему SQLite требует включить проверку внешних ключей для каждого соединения?", "Чем COUNT(*) по пустому выбору отличается от SUM(minutes)?", "Почему LEFT JOIN сохраняет t-02, а INNER JOIN теряет?", "Какие три проверки делают диаграмму честной?", "Что сохранилось после миграции 1.2 и почему локальный сервер нельзя открыть всем?"],
    "kz": ["Бақылау естелік бойынша қорытындыдан несімен өзгеше?", "Бос ұзақтық пен нөл минутты қалай ажыратасыз?", "Төрт жазба қосындысы неге 60 емес, 70?", "Тазалау алдында бастапқы CSV не үшін сақталады?", "Қай кілт сабақтың қайталануына, қайсысы белгісіз тапсырмаға сілтемеге жол бермейді?", "SQLite байланыста сыртқы кілт тексеруін неге бөлек қосуды талап етеді?", "Бос іріктемедегі COUNT(*) пен SUM(minutes) айырмасы қандай?", "Неге LEFT JOIN t-02 сақтап, INNER JOIN жоғалтады?", "Қандай үш тексеру диаграмманы адал етеді?", "1.2 көшіруден кейін не сақталды және жергілікті серверді неге көпшілікке ашуға болмайды?"],
    "en": ["How does an observation differ from a conclusion recalled later?", "How do you distinguish unknown duration from zero minutes?", "Why do the four records total 70 rather than 60?", "Why keep the raw CSV before cleaning?", "Which key prevents duplicate sessions, and which prevents links to missing tasks?", "Why enable foreign-key checks for every SQLite connection?", "How do COUNT(*) and SUM(minutes) differ on an empty selection?", "Why does LEFT JOIN keep t-02 while INNER JOIN loses it?", "What three checks make a chart honest?", "What survived migration to 1.2, and why is the local server not public?"],
}


def render(number, lang):
    row, item, labels = ROWS[number], CONTENT[number][lang], HEAD[lang]
    title = row[f"title_{lang}"]
    stem = f"{number:02d}-{row['id']}"
    code, output = EXAMPLES[number]
    next_row = ROWS.get(number + 1)
    next_link = (f"[{next_row[f'title_{lang}']}](/read/informatics-{number+1:02d}-{next_row['id']}?lang={lang})"
                 if number < 54 else f"[{'Все уроки' if lang == 'ru' else 'Барлық сабақтар' if lang == 'kz' else 'All lessons'}](/course/informatics?lang={lang})")
    sections = [f"# {title}", "", f"_{labels[0]}:_ **{item['lead']}**", "",
                f"## {labels[1]}", "", INTRO[lang].format(n=number), "",
                f"![{title}](/static/course/informatics/map-{stem}-{lang}.svg)", "",
                f"## {labels[2]}", "", item["scene"], "",
                f"## {labels[3]}", "", item["terms"], "",
                f"## {labels[4]}", "", SUPPORT[lang], "",
                f"## {labels[5]}", "", item["walk"], "",
                f"## {labels[6]}", "", PREDICT[lang], "", "```python", code, "```", "",
                f"## {labels[7]}", "", "```text", output, "```", "",
                f"## {labels[8]}", "", item["trap"], "",
                f"## {labels[9]}", "", item["project"], "",
                f"## {labels[10]}", "", item["task"], "",
                f"## {labels[11]}", "", item["transfer"], "",
                f"## {labels[12]}", "", REVIEW[lang], "",
                f"## {'Первоисточник для проверки' if lang == 'ru' else 'Тексеруге арналған бастапқы құжат' if lang == 'kz' else 'Primary reference to check'}", "",
                f"[{'Официальная документация' if lang == 'ru' else 'Ресми құжаттама' if lang == 'kz' else 'Official documentation'}]({RESOURCES[number]})", "",
                f"## {labels[13] if number < 54 else labels[14]}", "", next_link, ""]
    if number == 47:
        start = {"ru": "Где взять файлы проекта", "kz": "Жоба файлдарын қайдан аламыз", "en": "Where to get the project files"}[lang]
        help_text = {
            "ru": f"Откройте [папку проекта 1.2]({PROJECT_URL}), скачайте репозиторий через Code → Download ZIP и найдите `course/informatics-assistant/step-05`. Откройте терминал в этой папке; команда `python3 --version` должна показать установленный Python. Начните с CSV и JSON, не запускайте файл базы как программу.",
            "kz": f"[1.2 жобасының бумасын]({PROJECT_URL}) ашып, Code → Download ZIP арқылы репозиторийді жүктеңіз де, `course/informatics-assistant/step-05` табыңыз. Терминалды осы бумада ашыңыз; `python3 --version` орнатылған Python нұсқасын көрсетуі тиіс. CSV мен JSON-нан бастаңыз, дерекқор файлын бағдарлама ретінде қоспаңыз.",
            "en": f"Open the [1.2 project folder]({PROJECT_URL}), download the repository with Code → Download ZIP and find `course/informatics-assistant/step-05`. Open a terminal there; `python3 --version` should show an installed Python. Begin with CSV and JSON; do not run the database file as a program.",
        }[lang]
        index = sections.index(f"## {labels[3]}")
        sections[index:index] = [f"## {start}", "", help_text, ""]
    if number == 54:
        heading = {"ru": "Десять вопросов для самопроверки", "kz": "Өзін тексеруге арналған он сұрақ", "en": "Ten questions for self-check"}[lang]
        index = sections.index(f"## {labels[11]}")
        sections[index:index] = [f"## {heading}", "", *[f"{i}. {q}" for i, q in enumerate(MASTERY[lang], 1)], ""]
    return "\n".join(sections)


def main():
    assert set(CONTENT) == set(EXAMPLES) == set(range(47, 55))
    for number in range(47, 55):
        stem = f"{number:02d}-{ROWS[number]['id']}"
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            (OUT / f"{stem}{suffix}.md").write_text(render(number, lang), encoding="utf-8")
    print("Wrote 24 hand-localized data-block lessons")


if __name__ == "__main__":
    main()
