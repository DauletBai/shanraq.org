# Kazakh AI: steps 6–25 / Қазақша ЖИ: 6–25-қадамдар / ИИ для казахского: шаги 6–25

## Русский

Это точные программы из уроков на трёх языках. Каждый шаг запускается отдельно: создайте папку `ai-course`, как во втором уроке, или запустите файл прямо из копии репозитория. Нужен Python 3.10 или новее, без дополнительных библиотек. Сначала прочитайте объяснение в уроке; файлы показывают учебные этапы, а не готовую модель.

## Қазақша

Бұл файлдар үш тілдегі сабақтарда берілген бағдарламалармен бірдей. Әр қадам жеке іске қосылады: екінші сабақтағыдай `ai-course` қалтасын жасаңыз немесе файлды жоба көшірмесінен тікелей іске қосыңыз. Python 3.10 не одан жаңа нұсқа жеткілікті, қосымша кітапхана керек емес. Әуелі сабақтағы түсіндірмені оқыңыз; файлдар дайын үлгіні емес, оқу кезеңдерін көрсетеді.

## English

These files exactly match the programs in all three lesson versions. Run each step separately, either from the `ai-course` folder described in lesson 2 or directly from a repository checkout. Python 3.10 or later is enough; no extra packages are needed. Read the lesson explanation first. These are learning stages, not a finished model.

| Step / Қадам / Шаг | File / Файл | Example input / Мысал / Ввод |
|---|---|---|
| 06 | `step-06/letters.py` | `үй` |
| 07 | `step-07/words.py` | `Үйлерде, мектептерде?` |
| 08 | `step-08/roots.py` | `үй`, then / содан соң / затем `кітап` |
| 09 | `step-09/plural.py` | `қала`, `үй`, `қалам`, `тіл`, `кітап`, `мектеп` |
| 10 | `step-10/order.py` | Change `chosen` / `chosen` мәнін өзгертіңіз / измените `chosen` |
| 11 | `step-11/case.py` | `үйлерімізге`, `мектептерімізде` |
| 12 | `step-12/meanings.py` | `ат`, `үй` |
| 13 | `step-13/entities.py` | `Шахмат, сурет қашан?` |
| 14 | `step-14/intent.py` | `Шахмат қашан, қайда?` |
| 15 | `step-15/request.py` | `Шахмат қашан?`, `Робот қашан?` |

From the repository root, run `python3 course/kazakh-ai/step-09/plural.py` on macOS/Linux or `py course/kazakh-ai/step-09/plural.py` on Windows. The files do not depend on one another; the complete model will be assembled in later lessons.

## Steps 16–20 / 16–20-қадамдар / Шаги 16–20

Each of these steps has its own `facts.json`. Enter that step's directory before running its Python file, because the file is read from the current directory. Copy both files when working outside the repository. The club and source IDs are fictional exercises.

Әр қадамның өз `facts.json` файлы бар. Python файлын іске қоспай тұрып, сол қадамның қалтасына кіріңіз: бағдарлама файлды ағымдағы қалтадан оқиды. Жоба сыртында жұмыс істесеңіз, екі файлды бірге көшіріңіз. Үйірме мен дереккөз белгілері оқу үшін ойдан алынған.

У каждого шага свой файл `facts.json`. Перед запуском Python-файла перейдите в папку шага: программа читает данные из текущей папки. Вне репозитория копируйте оба файла. Кружок и обозначения источников вымышлены для упражнения.

| Step / Қадам / Шаг | Program / Бағдарлама / Программа | Focus / Мақсат / Цель |
|---|---|---|
| 16 | `step-16/catalog.py` | JSON catalog / тізім / каталог |
| 17 | `step-17/search.py` | Two-key lookup / екі кілтпен іздеу / поиск по двум ключам |
| 18 | `step-18/source.py` | Source and check date / дереккөз бен күн / источник и дата |
| 19 | `step-19/answer.py` | Answer template / жауап үлгісі / шаблон ответа |
| 20 | `step-20/model.py` | Question, conflict, expiry / сұрақ, қайшылық, мерзім / вопрос, конфликт, срок |

## Steps 21–25 / 21–25-қадамдар / Шаги 21–25

Each step has a Python file and its own `examples.json`. Steps 21–24 expose only 13 training and tuning cards; step 25 opens six additional test cards. Enter the step directory before running the Python file. The questions are fictional and the small dataset is for learning, not a measured product benchmark. Keep the `test` labels sealed while designing the classifier and gate.

Әр қадамда Python файлы мен жеке `examples.json` бар. 21–24-қадамдарда тек оқыту мен баптауға арналған 13 карточка ашық, ал 25-қадамда тағы алты бақылау карточкасы ашылады. Бағдарламаны іске қоспай тұрып, сол қалтаға кіріңіз. Сұрақтар ойдан алынған; шағын жиын дайын өнімнің сапасын өлшемейді. Жіктегіш пен тоқтату ережесін жасағанда `test` белгілерін ашпаңыз.

У каждого шага есть Python-файл и собственный `examples.json`. На шагах 21–24 открыты только 13 карточек обучения и настройки; на шаге 25 добавляются шесть контрольных. Перед запуском перейдите в папку шага. Вопросы вымышлены, а маленькая выборка служит обучению, не оценке готового продукта. При настройке классификатора и правила остановки не смотрите метки `test`.

| Step / Қадам / Шаг | Program / Бағдарлама / Программа | Focus / Мақсат / Цель |
|---|---|---|
| 21 | `step-21/labels.py` | Human labels / адам белгілері / человеческая разметка |
| 22 | `step-22/splits.py` | Train, tune, test / оқыту, баптау, бақылау |
| 23 | `step-23/classifier.py` | Learned word counts / сөз санағы / счётчики слов |
| 24 | `step-24/gate.py` | Explicit stop rule / тоқтату ережесі / правило остановки |
| 25 | `step-25/evaluate.py` | Precision, recall, refusal / дәлдік, толықтық, бас тарту |
