#!/usr/bin/env python3
"""Render manually localized Informatics lessons 55–72 and precise support maps."""
import json
from html import escape
from pathlib import Path

from informatics_final_content import CONTENT

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/informatics'
MAPS = ROOT / 'web/static/course/informatics'
ROWS = {r['number']: r for r in json.loads((ROOT/'course/informatics/curriculum.json').read_text())['lessons']}

H = {
'ru': ('Лид (summary)','Где мы на карте','Ситуация и вопрос','Новые слова без пропусков','Опорный сигнал','Разбираем шаг за шагом','Предскажите и проверьте','Ожидаемый вывод','Поймайте ошибку','Изменение проекта','Задание и доказательство','Перенос в новую ситуацию','Возврат через 1, 7 и 30 дней','Первоисточник для проверки','Следующий урок','Все уроки курса'),
'kz': ('Қысқаша мазмұн','Картадағы орнымыз','Жағдай және сұрақ','Жаңа сөздерді анық түсінейік','Сабақтың тірек сигналы','Қадамдап талдаймыз','Болжаңыз және тексеріңіз','Күтілетін нәтиже','Қатені табыңыз','Жобаға енгізілетін өзгеріс','Тапсырма және дәлел','Жаңа жағдайға көшіру','1, 7 және 30 күннен кейін қайталау','Тексеруге арналған бастапқы құжат','Келесі сабақ','Курстың барлық сабақтары'),
'en': ('Lead (summary)','Where we are on the map','Situation and question','New words without gaps',"The lesson's support signal",'Work through it step by step','Predict and check','Expected output','Catch the error','Project change','Task and evidence','Transfer to a new setting','Return after 1, 7, and 30 days','Primary reference to check','Next lesson','All course lessons'),
}
INTRO = {
'ru': ('Урок {n} из 72. Блок «Цифровая безопасность» (55–62). Мы продолжаем локального помощника с тремя вымышленными задачами. Каждая мера имеет конкретную границу: учебная версия не становится готовым публичным сервисом.', 'Урок {n} из 72. Блок «ИИ и выпуск проекта» (63–72). Помощник сохраняет прежние правила напоминания и локальную базу; новая функция только предлагает категорию задачи и допускает отказ.'),
'kz': ('72 сабақтың {n}-сі. «Цифрлық қауіпсіздік» блогы (55–62). Үш ойдан алынған тапсырмасы бар жергілікті көмекшіні жалғастырамыз. Әр шараның шегі нақты: оқу нұсқасы дайын ашық сервиске айналмайды.', '72 сабақтың {n}-сі. «ЖИ және жоба шығарылымы» блогы (63–72). Көмекшінің мерзім ережелері мен жергілікті базасы сақталады; жаңа мүмкіндік тек тапсырма санатын ұсынып, бас тартуға мүмкіндік береді.'),
'en': ('Lesson {n} of 72. Digital security block (55–62). We continue a local assistant with three fictional tasks. Each control has a defined boundary: the classroom version does not become a public service.', 'Lesson {n} of 72. AI and project release block (63–72). The assistant keeps its earlier reminder rules and local database; the new feature only suggests a task category and can abstain.'),
}
SUPPORT = {
'ru': 'Прочитайте три клетки схемы слева направо, затем закройте третью. Назовите вход и действие в каждой клетке; предскажите, чем закончится путь и какой наблюдаемый факт подтвердит ответ. Сравните рисунок с файлом проекта или контрольным примером. Если число на схеме не получается из данных, исправьте схему, а не данные под красивую картинку.',
'kz': 'Сұлбаның үш ұяшығын солдан оңға оқып, үшіншісін жабыңыз. Әр ұяшықтың кірісі мен әрекетін айтыңыз; жолдың немен аяқталатынын және қандай байқалатын фактіні тексеретініңізді болжаңыз. Суретті жоба файлымен не бақылау мысалымен салыстырыңыз. Сұлбадағы сан деректен шықпаса, әдемі сурет үшін деректі емес, сұлбаны түзетіңіз.',
'en': 'Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.',
}
PREDICT = {
'ru': 'Запишите ожидаемый вывод до запуска кода. Выполните его в указанной папке `step-06` или `step-07`, сравните результат символ за символом и объясните каждую строку. Затем измените один безопасный вымышленный вход и повторите цикл «предсказать — запустить — изменить — объяснить — проверить». Наблюдение на маленьком наборе не превращайте в обещание для реальных людей.',
'kz': 'Кодты қоспай тұрып күтілетін нәтижені жазыңыз. Оны `step-06` не `step-07` бумасында орындап, таңба бойынша салыстырыңыз және әр жолды түсіндіріңіз. Кейін бір қауіпсіз ойдан алынған кірісті өзгертіп, «болжау — іске қосу — өзгерту — түсіндіру — тексеру» айналымын қайталаңыз. Шағын жиындағы бақылауды нақты адамдарға берілген уәдеге айналдырмаңыз.',
'en': 'Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.',
}
PROJECT = {
'ru': ('В проекте 1.3 сохраняются исходные три задачи, четыре занятия и правила напоминаний. Новые действия вынесены в `security_assistant.py`: проверка, создание отдельной копии и восстановление только в новый файл. Запишите, от какой угрозы помогает это изменение и от какой не помогает; учебные данные остаются вымышленными.', 'В проекте 2.1 сохраняются база, отчёт, напоминания и резервирование. Файл `labelled_tasks.csv` содержит только вымышленные примеры. `classifier.py` выдаёт совет `math`, `reading` или `review`, но не записывает решение в базу. Каждый вывод модели сравнивается с отложенными примерами и пересматривается человеком.'),
'kz': ('1.3 жобасында бұрынғы үш тапсырма, төрт сабақ және еске салу ережесі сақталады. Жаңа әрекеттер `security_assistant.py` ішіне бөлінді: тексеру, бөлек көшірме жасау және тек жаңа файлға қалпына келтіру. Өзгеріс қай қатерге көмектесетінін және қайсына көмектеспейтінін жазыңыз; оқу дерегі ойдан алынған күйде.', '`2.1` жобасында база, есеп, еске салу және көшірме сақтау қалады. `labelled_tasks.csv` тек ойдан алынған мысалдарды қамтиды. `classifier.py` `math`, `reading` не `review` кеңесін береді, бірақ шешімді базаға жазбайды. Әр ұсыныс қалдырылған мысалмен салыстырылып, адам қайта қарайды.'),
'en': ('Release 1.3 preserves the three tasks, four sessions and reminder rules. New operations live in `security_assistant.py`: validation, a separate backup and restore only into a new file. State which threat each change reduces and which it does not; all teaching data stay fictional.', 'Release 2.1 preserves the DB, report, reminders and backup. `labelled_tasks.csv` contains fictional examples only. `classifier.py` suggests `math`, `reading` or `review` but writes no decision to the database. Check each model output against held-out examples and let a person reconsider it.'),
}
REVIEW = {
'ru': 'Через 1 день нарисуйте схему по памяти и найдите в ней одну границу применения. Через 7 дней объясните её товарищу, не читая урок, и решите 7 из 10 вопросов блока; ошибочные ответы разберите на конкретных входах. Через 30 дней откройте свежую копию проекта, воспроизведите проверку и перенесите принцип в новый пример. Сохраните свой протокол с входом, ожидаемым и фактическим результатом.',
'kz': '1 күннен кейін сұлбаны жатқа салып, қолданудың бір шегін табыңыз. 7 күннен кейін оны мәтінсіз досыңызға түсіндіріп, блоктың он сұрағының жетеуін шешіңіз; қате жауапты нақты кіріспен талдаңыз. 30 күннен кейін жобаның жаңа көшірмесін ашып, тексеруді қайталап, қағиданы жаңа мысалға көшіріңіз. Кіріс, күтілетін және нақты нәтижесі бар хаттаманы сақтаңыз.',
'en': 'After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.',
}
SOURCES = {
55:'https://www.nist.gov/cyberframework',56:'https://pages.nist.gov/800-63-4/sp800-63b.html',57:'https://www.sqlite.org/uri.html',58:'https://www.cisa.gov/secure-our-world/recognize-and-report-phishing',59:'https://www.nist.gov/publications/guidelines-selection-configuration-and-use-transport-layer-security-tls-implementations',60:'https://docs.python.org/3/library/sqlite3.html#sqlite3.Connection.backup',61:'https://www.unicef.org/parenting/child-care/online-privacy',62:'https://www.sqlite.org/pragma.html#pragma_integrity_check',63:'https://www.nist.gov/itl/ai-risk-management-framework',64:'https://scikit-learn.org/stable/modules/cross_validation.html',65:'https://docs.python.org/3/library/collections.html#collections.Counter',66:'https://scikit-learn.org/stable/modules/model_evaluation.html#confusion-matrix',67:'https://www.nist.gov/itl/ai-risk-management-framework',68:'https://arxiv.org/abs/1706.03762',69:'https://www.nist.gov/itl/ai-risk-management-framework',70:'https://www.nist.gov/itl/ai-risk-management-framework',71:'https://docs.python.org/3/library/unittest.html',72:'https://www.nist.gov/itl/ai-risk-management-framework'}

EXAMPLES = {
55: ('from pathlib import Path\nfrom data_store import connect_readonly\nprint("local", Path("tasks.json").is_file())\nprint("read-only", "mode=ro" in Path("data_store.py").read_text())','local True\nread-only True'),
56: ('import hashlib, secrets\nsalt = secrets.token_bytes(16)\na = hashlib.scrypt(b"fictional long phrase", salt=salt, n=2**14, r=8, p=1)\nb = hashlib.scrypt(b"fictional long phrase", salt=salt, n=2**14, r=8, p=1)\nprint(secrets.compare_digest(a, b))','True'),
57: ('from pathlib import Path\ntext = Path("data_store.py").read_text()\nprint("mode=ro" in text, "destination.exists()" in text)','True True'),
58: ('from urllib.parse import urlsplit\nexpected = "school.example"\nfor url in ("https://school.example/login", "https://school.example.attacker.test/login"):\n    print(urlsplit(url).hostname == expected)','True\nFalse'),
59: ('import hashlib\na = hashlib.sha256(b"20 minutes").hexdigest()\nb = hashlib.sha256(b"200 minutes").hexdigest()\nprint(a == b)','False'),
60: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom data_store import create_database\nfrom security_assistant import backup, restore, verify\nwith TemporaryDirectory() as folder:\n    db, saved, restored = (Path(folder) / name for name in ("live.db", "backup.db", "restored.db"))\n    create_database(db, "tasks.json", "study_sessions.csv")\n    backup(db, saved)\n    print(verify(db), restore(saved, restored))','(3, 4) (3, 4)'),
61: ('import csv\nwith open("study_sessions.csv", encoding="utf-8") as file:\n    print(next(csv.reader(file)))','[\'session_id\', \'task_id\', \'observed_on\', \'minutes\', \'source\']'),
62: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom data_store import create_database, minutes_report\nfrom security_assistant import backup, restore, verify\nwith TemporaryDirectory() as folder:\n    a, b, c = (Path(folder) / name for name in ("a.db", "b.db", "c.db"))\n    create_database(a, "tasks.json", "study_sessions.csv")\n    backup(a, b); restore(b, c)\n    print(verify(c), minutes_report(c))','(3, 4) [(\'t-01\', 45), (\'t-02\', 0), (\'t-03\', 25)]'),
63: ('from datetime import date\nfrom assistant_core import reminder_status\nfrom classifier import read_examples, train, suggest\nrows = read_examples("labelled_tasks.csv")\nprint(reminder_status({"done": False, "due_date": "2026-10-10"}, date(2026,10,9)))\nprint(suggest("Кітап оқу", train(rows)))','REMIND\nreading'),
64: ('from classifier import read_examples\nrows = read_examples("labelled_tasks.csv")\nprint(sum(r["split"] == "train" for r in rows), sum(r["split"] == "test" for r in rows))','12 4'),
65: ('from classifier import read_examples, train, suggest\nmodel = train(read_examples("labelled_tasks.csv"))\nprint(suggest("Кітап оқу", model), suggest("Геометрия есебі", model), suggest("unknown", model))','reading math review'),
66: ('from classifier import read_examples, train, evaluate\nrows = read_examples("labelled_tasks.csv")\nmatrix = evaluate(rows, train(rows))\nprint(matrix[("reading", "reading")], matrix[("math", "math")], sum(matrix.values()))','2 2 4'),
67: ('from classifier import read_examples, train, suggest\nrows = read_examples("labelled_tasks.csv")\nmodel = train(rows)\nprint([suggest(r["title"], model) for r in rows if r["split"] == "test"])','[\'reading\', \'reading\', \'math\', \'math\']'),
68: ('from collections import Counter\nnext_words = Counter(("class", "holiday"))\nprint(next_words["class"], next_words["holiday"])','1 1'),
69: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom data_store import create_database, minutes_report\nwith TemporaryDirectory() as folder:\n    db = Path(folder) / "example.db"\n    create_database(db, "tasks.json", "study_sessions.csv")\n    print(dict(minutes_report(db))["t-02"])','0'),
70: ('from pathlib import Path\nfrom tempfile import TemporaryDirectory\nfrom data_store import create_database\nfrom classifier import read_examples, train, suggest\nwith TemporaryDirectory() as folder:\n    db = Path(folder) / "example.db"\n    create_database(db, "tasks.json", "study_sessions.csv")\n    before = db.read_bytes()\n    print(suggest("Кітап оқу", train(read_examples("labelled_tasks.csv"))))\n    print(before == db.read_bytes())','reading\nTrue'),
71: ('from classifier import read_examples, train, evaluate\nrows = read_examples("labelled_tasks.csv")\nprint(len(rows), sum(evaluate(rows, train(rows)).values()))','16 4'),
72: ('from classifier import read_examples, train, suggest\nmodel = train(read_examples("labelled_tasks.csv"))\nprint(*(suggest(title, model) for title in ("Кітап оқу", "Геометрия есебі", "unknown")))','reading math review'),
}

def page(number, lang):
    row, item, h = ROWS[number], CONTENT[number][lang], H[lang]
    block = 0 if number <= 62 else 1
    stem = f"{number:02d}-{row['id']}"
    code, output = EXAMPLES[number]
    nextrow = ROWS.get(number+1)
    nextlink = f"[{nextrow[f'title_{lang}']}](/read/informatics-{number+1:02d}-{nextrow['id']}?lang={lang})" if nextrow else f"[{'Все уроки' if lang=='ru' else 'Барлық сабақтар' if lang=='kz' else 'All lessons'}](/course/informatics?lang={lang})"
    folder = 'step-06' if block == 0 else 'step-07'
    repo = f'https://github.com/DauletBai/shanraq.org/tree/main/course/informatics-assistant/{folder}'
    lead = item['scene'].split('.')[0] + '.'
    common = {
      'ru': f'Откройте [папку контрольной версии]({repo}), скачайте репозиторий через Code → Download ZIP и найдите `course/informatics-assistant/{folder}`. Откройте терминал в этой папке. Команда `python3 --version` покажет Python; затем запустите указанную проверку. Все имена задач вымышлены, никаких реальных данных вводить не нужно.',
      'kz': f'[Бақылау нұсқасының бумасын]({repo}) ашып, Code → Download ZIP арқылы репозиторийді жүктеп, `course/informatics-assistant/{folder}` табыңыз. Терминалды осы бумада ашыңыз. `python3 --version` Python нұсқасын көрсетеді; содан кейін сабақтағы тексеруді орындаңыз. Барлық тапсырма атауы ойдан алынған, нақты дерек енгізбеңіз.',
      'en': f'Open the [checkpoint folder]({repo}), download the repository with Code → Download ZIP and find `course/informatics-assistant/{folder}`. Open a terminal there. `python3 --version` shows Python; then run the lesson check. All tasks are fictional; enter no real personal data.',
    }[lang]
    parts = [f'# {row[f"title_{lang}"]}', '', f'_{h[0]}:_ **{lead}**', '', f'## {h[1]}', '', INTRO[lang][block].format(n=number), '', f'![{row[f"title_{lang}"]}](/static/course/informatics/map-{stem}-{lang}.svg)', '', f'## {h[2]}', '', item['scene'], '', f'## {h[3]}', '', item['terms'], '', f'## {h[4]}', '', SUPPORT[lang], '', f'## {h[5]}', '', item['walk'], '', f'## {h[6]}', '', PREDICT[lang], '', '```python',code,'```','',f'## {h[7]}','','```text',output,'```','',f'## {h[8]}','',item['trap'],'',f'## {h[9]}','',PROJECT[lang][block],'',f'## {h[10]}','',item['task'],'',f'## {h[11]}','',item['transfer'],'',f'## {h[12]}','',REVIEW[lang],'',f'## {h[13]}','',f"[{'Официальный документ' if lang=='ru' else 'Ресми құжат' if lang=='kz' else 'Official reference'}]({SOURCES[number]})",'',f'## {h[14] if number<72 else h[15]}','',nextlink,'']
    if number in (55,63):
        idx=parts.index(f'## {h[3]}')
        title={'ru':'Где взять файлы проекта','kz':'Жоба файлдарын қайдан аламыз','en':'Where to get the project files'}[lang]
        parts[idx:idx]=[f'## {title}','',common,'']
    if number in (62,72):
        questions={
        62:{'ru':['Что именно защищает mode=ro?','Что является активом?','Чем отличается угроза от уязвимости?','Зачем соли разные?','Почему два пароля не два фактора?','Что делать с фишинговой ссылкой?','Чем хеш отличается от шифра?','Что означает 3/4?','Почему восстановление в новый файл?','Что остаётся вне защиты 1.3?'],
            'kz':['mode=ro нені қорғайды?','Құндылық деген не?','Қатер мен осалдық айырмасы қандай?','Тұз неге бөлек?','Екі құпиясөз неге екі фактор емес?','Фишинг сілтемесімен не істейміз?','Хеш пен шифр айырмасы қандай?','3/4 нені білдіреді?','Неге жаңа файлға қайтарамыз?','1.3 нені қорғамайды?'],
            'en':['What does mode=ro protect?','What is the asset?','How do threat and vulnerability differ?','Why use distinct salts?','Why are two passwords not two factors?','What should you do with a phishing link?','How does a hash differ from encryption?','What does 3/4 mean?','Why restore to a new file?','What remains outside release 1.3?']},
        72:{'ru':['Чем правило отличается от модели?','Что такое признак?','Что означает метка?','Почему test не участвует в обучении?','Какая матрица получилась на 4 строках?','Почему 4/4 не обещание?','Что означает review?','Кто принимает окончательное решение?','Как проверять факт ИИ?','Какая граница у версии 2.1?'],
            'kz':['Ереже мен модель айырмасы қандай?','Белгі деген не?','Таңба нені білдіреді?','Неге test үйретуге кірмейді?','Төрт жол матрицасы қандай?','Неге 4/4 кепіл емес?','review нені білдіреді?','Соңғы шешімді кім қабылдайды?','ЖИ фактісін қалай тексереміз?','2.1 шегі қандай?'],
            'en':['How does a rule differ from a model?','What is a feature?','What is a label?','Why is test excluded from training?','What matrix came from four rows?','Why is 4/4 no guarantee?','What does review mean?','Who makes the final decision?','How do you verify an AI fact?','What is the boundary of 2.1?']}}
        idx=parts.index(f'## {h[11]}')
        heading={'ru':'Десять вопросов для самопроверки','kz':'Өзін тексеруге арналған он сұрақ','en':'Ten questions for self-check'}[lang]
        parts[idx:idx]=[f'## {heading}','',*[f'{i}. {q}' for i,q in enumerate(questions[number][lang],1)],'']
    return '\n'.join(parts)

def map_svg(number,lang):
    row,item=ROWS[number],CONTENT[number][lang]
    title=row[f'title_{lang}']
    words=title.split(); lines=['']
    for word in words:
        nextline=(lines[-1]+' '+word).strip()
        if len(nextline)>43 and lines[-1]:lines.append(word)
        else:lines[-1]=nextline
    if len(lines)>2:raise ValueError((number,lang,'title too long',lines))
    heading=''.join(f'<text x="64" y="{147+i*55}" class="title">{escape(line)}</text>' for i,line in enumerate(lines))
    cards=[]
    for i,label in enumerate(item['map']):
        if len(label)>29:raise ValueError((number,lang,label))
        x=54+i*514
        color=('#66e6df','#ffd27d','#ff9584')[i]
        cards.append(f'<g class="final-stage" transform="translate({x} 307)"><rect x="9" y="13" width="482" height="354" rx="22" fill="#020d17" opacity=".67"/><rect width="482" height="354" rx="22" fill="url(#panel)" stroke="{color}" stroke-width="3"/><path d="M22 22 H460" stroke="#b7dcea" opacity=".4"/><rect x="25" y="24" width="59" height="57" rx="14" fill="{color}"/><text x="54" y="65" text-anchor="middle" class="number">{i+1}</text><path d="M27 107 H454" stroke="#6e91a9" stroke-width="2"/><rect x="27" y="136" width="428" height="89" rx="12" fill="#061b2b" stroke="#4b758b"/><text x="241" y="191" text-anchor="middle" class="card">{escape(label)}</text><text x="241" y="291" text-anchor="middle" class="small">{escape(("Кіріс → әрекет → тексеру" if lang=="kz" else "Input → action → check" if lang=="en" else "Вход → действие → проверка"))}</text></g>')
    foot={'ru':'Предскажи → запусти → измени → объясни → проверь','kz':'Болжа → іске қос → өзгерт → түсіндір → тексер','en':'Predict → run → change → explain → verify'}[lang]
    group='SECURITY' if number<=62 else 'AI · RELEASE'
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1600 900" role="img" aria-label="{escape(title)}"><defs><linearGradient id="bg" x2="1" y2="1"><stop stop-color="#071827"/><stop offset="1" stop-color="#164c60"/></linearGradient><linearGradient id="panel" x2="1" y2="1"><stop stop-color="#286076"/><stop offset="1" stop-color="#102b3f"/></linearGradient><pattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"><path d="M40 0H0V40" fill="none" stroke="#c7eaf2" stroke-opacity=".12"/></pattern></defs><style>.eyebrow{{font:700 22px Arial,sans-serif;fill:#86eee7;letter-spacing:3px}}.title{{font:700 43px Arial,sans-serif;fill:#f5f9ff}}.number{{font:700 30px Arial,sans-serif;fill:#102333}}.card{{font:700 27px Arial,sans-serif;fill:#fff}}.small{{font:600 24px Arial,sans-serif;fill:#e0f4f8}}.foot{{font:600 27px Arial,sans-serif;fill:#e0f0f7}}</style><rect width="1600" height="900" fill="url(#bg)"/><rect width="1600" height="900" fill="url(#grid)"/><text x="64" y="75" class="eyebrow">{group} · {number:02d} / 72</text>{heading}<path d="M64 253 H1536" stroke="#6b94ab" stroke-width="2"/>{''.join(cards)}<path d="M544 482 h18 m-9 -9 9 9 -9 9 M1058 482 h18 m-9 -9 9 9 -9 9" fill="none" stroke="#e2f6f9" stroke-width="5"/><path d="M64 755 H1536" stroke="#6b94ab" stroke-width="2"/><text x="800" y="812" text-anchor="middle" class="foot">{escape(foot)}</text></svg>'''

def main():
    assert set(CONTENT)==set(EXAMPLES)==set(range(55,73))
    for number in range(55,73):
        stem=f"{number:02d}-{ROWS[number]['id']}"
        for lang in ('ru','kz','en'):
            suffix='' if lang=='ru' else f'-{lang}'
            (LESSONS/f'{stem}{suffix}.md').write_text(page(number,lang),encoding='utf-8')
            (MAPS/f'map-{stem}-{lang}.svg').write_text(map_svg(number,lang),encoding='utf-8')
    print('Wrote 54 final lessons and 54 localized maps')

if __name__=='__main__':main()
