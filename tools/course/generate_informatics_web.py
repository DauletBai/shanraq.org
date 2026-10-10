#!/usr/bin/env python3
"""Render reviewed trilingual web-block lessons from hand-written content."""

import json
from pathlib import Path

from informatics_web_content import CONTENT

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "course/lessons/informatics"
ROWS = {row["number"]: row for row in json.loads(
    (ROOT / "course/informatics/curriculum.json").read_text(encoding="utf-8"))["lessons"]}

HEAD = {
    "ru": ("Лид (summary)", "Где мы на карте", "Ситуация и вопрос", "Новые слова без пропусков", "Опорный сигнал", "Разбираем шаг за шагом", "Предскажите и проверьте", "Ожидаемый вывод", "Поймайте ошибку", "Изменение проекта", "Задание и доказательство", "Перенос в новую ситуацию", "Возврат через 1, 7 и 30 дней", "Следующий урок", "Все уроки курса"),
    "kz": ("Лид (summary)", "Картадағы орнымыз", "Жағдай және сұрақ", "Жаңа сөздерді анық түсінейік", "Сабақтың тірек сигналы", "Қадамдап талдаймыз", "Болжаңыз және тексеріңіз", "Күтілетін нәтиже", "Қатені табыңыз", "Жобаға енгізілетін өзгеріс", "Тапсырма және дәлел", "Жаңа жағдайға көшіру", "1, 7 және 30 күннен кейін қайталау", "Келесі сабақ", "Курстың барлық сабақтары"),
    "en": ("Lead (summary)", "Where we are on the map", "Situation and question", "New words without gaps", "The lesson's support signal", "Work through it step by step", "Predict and check", "Expected output", "Catch the error", "Project change", "Task and evidence", "Transfer to a new setting", "Return after 1, 7, and 30 days", "Next lesson", "All course lessons"),
}

EXAMPLES = {
    39: ('hops = ["phone", "router", "provider", "server"]\nprint(" → ".join(hops))', 'phone → router → provider → server'),
    40: ('parts = {2: "AS", 1: "T", 3: "K"}\nprint("".join(parts[n] for n in sorted(parts)))', 'TASK'),
    41: ('records = {"assistant.test": "203.0.113.7"}\nprint(records["assistant.test"])', '203.0.113.7'),
    42: ('received = {3: "done", 1: "task", 2: ":"}\nprint("".join(received[n] for n in sorted(received)))', 'task:done'),
    43: ('from urllib.parse import urlsplit\naddress = urlsplit("https://assistant.test/tasks")\nprint(address.scheme, address.hostname)', 'https assistant.test'),
    44: ('from urllib.parse import urlsplit\naddress = urlsplit("http://127.0.0.1:8765/?today=2026-10-09")\nprint(address.path, address.query)', '/ today=2026-10-09'),
    45: ('from html import escape\nprint(escape("<b>Task</b>"))', '&lt;b&gt;Task&lt;/b&gt;'),
    46: ('from datetime import date\nfrom assistant_core import load_document, reminder_status\ntasks = load_document("tasks.json")["tasks"]\nfor task in tasks:\n    print(task["id"], reminder_status(task, date(2026, 10, 9)))', 't-01 REMIND\nt-02 DONE\nt-03 NO_DATE'),
}

INTRO = {
    "ru": "Это урок {n} из 72, блок «Интернет, веб и облако» (39–46). Мы продолжаем один проект: три вымышленные записи и пять результатов напоминания из версии 1.0 должны сохраниться. На опорной схеме показан ровно тот переход, который можно объяснить и проверить; стрелка не означает, что все промежуточные устройства нарисованы.",
    "kz": "Бұл 72 сабақтың {n}-сі, «Интернет, веб және бұлт» блогы (39–46). Бір жобаны жалғастырамыз: 1.0 нұсқасындағы ойдан алынған үш жазба мен еске салудың бес нәтижесі сақталуы тиіс. Тірек сызба тексеруге болатын өтуді ғана көрсетеді; жебе барлық аралық құрылғы салынғанын білдірмейді.",
    "en": "This is lesson {n} of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.",
}
SUPPORT = {
    "ru": "Прочитайте схему слева направо. Для каждой стрелки назовите вход, действие и проверяемый результат. Затем закройте подпись следующей карточки и предскажите её своими словами. Вернитесь к реальному примеру: различайте учебную аналогию, работу localhost и возможный удалённый сайт. Если картинка обещает больше, чем объяснение, исправьте объяснение или рисунок.",
    "kz": "Сызбаны солдан оңға қарай оқыңыз. Әр жебенің кірісін, әрекетін және тексерілетін нәтижесін атаңыз. Келесі карточканың жазуын жауып, оны өз сөзіңізбен болжаңыз. Нақты мысалға оралыңыз: оқу ұқсастығын, localhost жұмысын және болашақ қашық сайтты ажыратыңыз. Сурет түсіндірмеден артық уәде етсе, мәтінді не суретті түзетіңіз.",
    "en": "Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.",
}
PREDICT = {
    "ru": "Не запускайте код сразу. Выпишите точный ожидаемый результат, включая порядок строк и символы. Объясните, что код доказывает, а чего не доказывает: этот маленький пример моделирует одно правило и не создаёт настоящий интернет, TLS-сертификат или открытый сервер. Затем выполните фрагмент в каталоге `step-04` и сравните результат посимвольно. Измените один вход, сначала предскажите новый результат, потом проверьте.",
    "kz": "Кодты бірден іске қоспаңыз. Жол реті мен таңбаларын сақтап, нақты нәтиже жазыңыз. Бұл шағын мысал нені дәлелдейді, нені дәлелдемейді, түсіндіріңіз: ол бір ережені көрсетеді, бірақ шынайы интернет, TLS сертификаты немесе ашық сервер жасамайды. Содан соң үзіндіні `step-04` бумасында орындап, нәтижені таңба бойынша салыстырыңыз. Бір кірісті өзгертіп, алдымен болжап, кейін тексеріңіз.",
    "en": "Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.",
}
REVIEW = {
    "ru": "Через 1 день нарисуйте главный переход по памяти и приведите контрпример к слишком простой аналогии. Через 7 дней объясните новый случай однокласснику, не читая готовый текст; проверьте не менее 7 из 10 вопросов. Через 30 дней повторите опыт с локальной веб-версией и убедитесь, что три исходных задачи и пять правил сохранились. Запишите, какой термин всё ещё непонятен, и вернитесь к его первому объяснению.",
    "kz": "1 күннен кейін негізгі өтуді жатқа сызып, тым қарапайым ұқсастыққа қарсы мысал келтіріңіз. 7 күннен кейін мәтінге қарамай жаңа жағдайды досыңызға түсіндіріп, 10 сұрақтың кемінде 7-еуін тексеріңіз. 30 күннен кейін жергілікті веб-нұсқаны қайта ашып, үш бастапқы тапсырма мен бес ереженің сақталғанын көріңіз. Түсініксіз терминді жазып, оның алғашқы анықтамасына оралыңыз.",
    "en": "After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.",
}

RESOURCES = {
    39: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview",
    40: "https://www.rfc-editor.org/rfc/rfc791",
    41: "https://developer.mozilla.org/en-US/docs/Glossary/DNS",
    42: "https://www.rfc-editor.org/info/rfc9293",
    43: "https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Transport_Layer_Security",
    44: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Messages",
    45: "https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label",
    46: "https://docs.python.org/3/library/http.server.html",
}

PREFLIGHT = {
    "ru": "Сначала откройте [папку учебного проекта](https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/step-04); скачайте репозиторий через Code → Download ZIP и откройте в нём `course/informatics-assistant/step-04`. Перед запуском выполните три знакомых действия из уроков 27–38: откройте терминал в `step-04`, выполните `python3 --version` и `python3 -c 'print(2 + 1)'`. Убедитесь, что видите номер Python и число `3`. Найдите в папке `web_assistant.py` и `tasks.json`; не открывайте JSON как программу. Если команда `python3` не найдена, вернитесь к инструкции установки Python из блока 1.0. Порт `8765` — номер службы на вашем компьютере; его не нужно покупать или настраивать у провайдера. После запуска оставьте терминал открытым и только тогда перейдите по адресу в браузере. Если адрес занят, остановите предыдущий учебный процесс Ctrl+C или запустите `--port 8766` и измените число в адресе.",
    "kz": "Алдымен [оқу жобасының бумасын](https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/step-04) ашыңыз; Code → Download ZIP арқылы репозиторийді жүктеп, ішіндегі `course/informatics-assistant/step-04` бумасына өтіңіз. Қоспас бұрын 27–38 сабақтан таныс үш әрекетті орындаңыз: терминалды `step-04` ішінде ашып, `python3 --version` және `python3 -c 'print(2 + 1)'` жазыңыз. Python нұсқасы мен `3` санын көріңіз. Бумада `web_assistant.py` және `tasks.json` табыңыз; JSON файлын бағдарлама ретінде іске қоспаңыз. `python3` табылмаса, 1.0 блогындағы орнату нұсқауына оралыңыз. `8765` порты — өз компьютеріңіздегі қызмет нөмірі; оны провайдерден сатып алу не баптау қажет емес. Сервер қосылған соң терминалды жаппай, браузердегі мекенжайға өтіңіз. Порт бос болмаса, алдыңғы үдерісті Ctrl+C арқылы тоқтатыңыз не `--port 8766` қолданып, мекенжайдағы санды өзгертіңіз.",
    "en": "First open the [classroom project folder](https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/step-04); download the repository with Code → Download ZIP and open `course/informatics-assistant/step-04` inside it. Before starting, repeat three familiar actions from lessons 27–38: open a terminal in `step-04`, run `python3 --version`, and run `python3 -c 'print(2 + 1)'`. Confirm that Python prints a version and then `3`. Find `web_assistant.py` and `tasks.json`; do not try to run JSON as a program. If `python3` is missing, return to the Python setup instructions in block 1.0. Port `8765` is a service number on your own computer; you do not buy it from a provider. Leave the terminal open after starting the server, then open the address in a browser. If the port is occupied, stop the earlier process with Ctrl+C or use `--port 8766` and update the browser address.",
}

MASTERY = {
    "ru": ["Кто начинает запрос и кто отвечает, если сервер запущен на вашем компьютере?", "Что означает 127.0.0.1 и почему страница может открываться без интернета?", "Почему один IP-адрес не доказывает, что за ним один человек?", "Что возвращает DNS и почему этого недостаточно для доверия сайту?", "Как TCP и простая UDP-модель ведут себя при потере части 2?", "Что проверяет TLS и чего не обещает значок замка?", "Какие ответы дают 200, 400, 404 и 405 в нашем проекте?", "Почему GET у нас не отмечает задачу выполненной?", "Зачем странице связка label/input и экранирование названия задачи?", "Какие три изменения нужны прежде, чем дать доступ семьям через интернет?"],
    "kz": ["Сервер өз компьютеріңізде болса, сұрауды кім бастап, жауапты кім береді?", "127.0.0.1 нені білдіреді және бет интернетсіз неге ашылуы мүмкін?", "Бір IP-ден бір ғана адамды неге анықтай алмаймыз?", "DNS не қайтарады және бұл сайтқа сену үшін неге жеткіліксіз?", "2-бөлік жоғалса, TCP мен қарапайым UDP моделі қалай әрекет етеді?", "TLS нені тексереді, құлып белгісі нені уәде етпейді?", "Жобада 200, 400, 404 және 405 жауаптары қашан беріледі?", "Неге бізде GET тапсырманы орындалды деп белгілемейді?", "label/input байланысы мен тапсырма атауын экранирлеу не үшін қажет?", "Отбасыларға интернет арқылы қолжетімділік бермес бұрын қандай үш өзгеріс керек?"],
    "en": ["Who starts a request and who replies when the server runs on your own computer?", "What does 127.0.0.1 mean, and why can the page work without internet?", "Why does one IP address not prove there is one human behind it?", "What does DNS return, and why is that not enough to trust a site?", "What do TCP and the simple UDP model do if part 2 is lost?", "What does TLS check, and what does a padlock not promise?", "When does our project return 200, 400, 404, and 405?", "Why does GET not mark a task complete in our release?", "Why do label/input pairing and task-title escaping matter?", "What three changes are needed before families can use this through the internet?"],
}


def render(number, lang):
    row, item = ROWS[number], CONTENT[number][lang]
    labels = HEAD[lang]
    title = row[f"title_{lang}"]
    stem = f"{number:02d}-{row['id']}"
    next_row = ROWS.get(number + 1)
    code, output = EXAMPLES[number]
    next_link = (f"[{next_row[f'title_{lang}']}](/{'read/informatics-' + f'{number+1:02d}-' + next_row['id']}?lang={lang})"
                 if next_row and number < 46 else f"[{'Все уроки' if lang == 'ru' else 'Барлық сабақтар' if lang == 'kz' else 'All lessons'}](/course/informatics?lang={lang})")
    sections = [f"# {title}", "", f"_{labels[0]}:_ **{item['lead']}**", "",
                f"## {labels[1]}", "", INTRO[lang].format(n=number), "",
                f"![{title}](/static/course/informatics/map-{stem}-{lang}.svg)", "",
                f"## {labels[2]}", "", item['scene'], "",
                f"## {labels[3]}", "", item['terms'], "",
                f"## {labels[4]}", "", SUPPORT[lang], "",
                f"## {labels[5]}", "", item['walk'], "",
                f"## {labels[6]}", "", PREDICT[lang], "", "```python", code, "```", "",
                f"## {labels[7]}", "", "```text", output, "```", "",
                f"## {labels[8]}", "", item['trap'], "",
                f"## {labels[9]}", "", item['project'], "",
                f"## {labels[10]}", "", item['task'], "",
                f"## {labels[11]}", "", item['transfer'], "",
                f"## {labels[12]}", "", REVIEW[lang], "",
                f"## {'Первоисточник для проверки' if lang == 'ru' else 'Тексеруге арналған бастапқы құжат' if lang == 'kz' else 'Primary reference to check'}", "",
                f"[{'Официальная документация' if lang == 'ru' else 'Ресми құжаттама' if lang == 'kz' else 'Official documentation'}]({RESOURCES[number]})", "",
                f"## {labels[13] if number < 46 else labels[14]}", "", next_link, ""]
    if number == 44:
        heading = {"ru": "Перед первым запуском", "kz": "Алғашқы іске қосу алдында", "en": "Before the first run"}[lang]
        index = sections.index(f"## {labels[3]}")
        sections[index:index] = [f"## {heading}", "", PREFLIGHT[lang], ""]
    if number == 46:
        heading = {"ru": "Десять вопросов для самопроверки", "kz": "Өзін тексеруге арналған он сұрақ", "en": "Ten questions for self-check"}[lang]
        index = sections.index(f"## {labels[11]}")
        questions = [f"{n}. {question}" for n, question in enumerate(MASTERY[lang], 1)]
        sections[index:index] = [f"## {heading}", "", *questions, ""]
    return "\n".join(sections)


def main():
    assert set(CONTENT) == set(EXAMPLES) == set(range(39, 47))
    for number in range(39, 47):
        stem = f"{number:02d}-{ROWS[number]['id']}"
        for lang in ("ru", "kz", "en"):
            suffix = "" if lang == "ru" else f"-{lang}"
            (OUT / f"{stem}{suffix}.md").write_text(render(number, lang), encoding="utf-8")
    print("Wrote 24 hand-localized web-block lessons")


if __name__ == "__main__":
    main()
