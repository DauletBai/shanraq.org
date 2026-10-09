"""Topic-specific bridges for the trilingual Informatics lessons.

Each bridge defines the words a beginner needs before asking them to use the
model, then makes the learner test that model on the continuing project.
"""

ENRICHMENT = {
    5: {
        "ru": {
            "hook": "Вы исправили название задачи, увидели его на экране и закрыли программу. После нового запуска вернулось старое название. Куда делась правка? Проследим путь записи, прежде чем искать виноватую кнопку.",
            "bridge": """## Разберём слова на одном опыте

**Инструкция** — отдельное действие, которое выполняет процессор; программа задаёт последовательность таких действий. **ОЗУ** — рабочее место для данных и работающей программы, содержимое которого обычно исчезает при выключении. **Накопитель** — место для файла, который должен пережить выключение. **Кэш** — небольшая область с недавно использованными данными: он ускоряет доступ, но сам по себе не подтверждает запись файла. **Сохранить** значит попросить программу передать изменения в постоянное хранилище; **проверить сохранение** значит закрыть и снова открыть файл.

Проведите безопасный опыт с копией `tasks.json`. Запишите в таблицу три момента: до открытия, после изменения одной вымышленной задачи и после повторного открытия. До нажатия «Сохранить» предскажите, где находится новая версия. Затем сохраните, закройте редактор и откройте файл вновь. Если новая запись видна только до закрытия, экран показывал временное состояние. Если осталась после запуска, данные дошли до накопителя. Проверяйте копию: так опыт не повредит исходный проект.

**Проверка понимания.** Ученик говорит: «Файл виден на экране, значит, он сохранён». Попросите показать повторное чтение. Теперь мысленно отключите питание до сохранения и после него. Результаты различаются потому, что экран, ОЗУ и накопитель выполняют разные роли. Это объяснение полезнее, чем рекламное число гигабайтов.""",
        },
        "kz": {
            "hook": "Тапсырманың атауын түзетіп, экраннан көрдіңіз де, бағдарламаны жаптыңыз. Қайта ашқанда ескі атау тұр. Өзгеріс қайда жоғалды? Батырманы кінәламай тұрып, жазбаның жолын бақылайық.",
            "bridge": """## Бір тәжірибеде жаңа сөздерді ашайық

**Нұсқау** — процессор орындайтын жеке әрекет; бағдарлама осындай әрекеттердің ретін береді. **Жедел жад** — жұмыс істеп тұрған бағдарлама мен деректердің уақытша орны; қуат өшкенде оның мазмұны әдетте жоғалады. **Жинақтауыш** — өшіргеннен кейін де қалуы тиіс файлдың орны. **Кэш** — жақында қолданылған деректерді сақтайтын шағын аймақ: ол қолжетімділікті жылдамдатады, бірақ файл жазылғанын дәлелдемейді. **Сақтау** — өзгерісті тұрақты сақтау құрылғысына жаздыру; **сақталғанын тексеру** — файлды жауып, қайта ашу.

`tasks.json` файлының көшірмесімен қауіпсіз тәжірибе жасаңыз. Кестеге үш сәтті жазыңыз: ашқанға дейін, ойдан шығарылған бір тапсырманы өзгерткеннен кейін және файлды қайта ашқаннан кейін. «Сақтау» батырмасына дейін жаңа нұсқа қайда тұрғанын болжаңыз. Содан кейін сақтап, редакторды жауып, файлды қайта ашыңыз. Жаңа жазба тек жабылғанға дейін көрінсе, экран уақытша күйді көрсеткен. Қайта ашқаннан кейін де тұрса, дерек жинақтауышқа жеткен. Түпнұсқаны бүлдірмеу үшін көшірмені қолданыңыз.

**Түсінгеніңізді тексеріңіз.** «Экранда көрініп тұр, демек сақталды» деген пікірге қайта ашып оқуды дәлел ретінде сұраңыз. Қуат сақтау алдында және одан кейін өшсе, нәтиже неге өзгереді? Өйткені экран, жедел жад және жинақтауыштың міндеттері бөлек.""",
        },
        "en": {
            "hook": "You changed a task title, saw it on screen, and closed the app. The old title returned when you reopened it. Where did the edit go? Follow the data before blaming a button.",
            "bridge": """## Unpack the words with one experiment

An **instruction** is one action the processor carries out; a program specifies an ordered set of instructions. **RAM** is the working area for a running program and its data; its ordinary contents disappear when power is removed. **Storage** keeps files that must survive shutdown. A **cache** is a small area holding recently used data: it can speed access but does not prove that a file was saved. To **save** is to ask the program to write changes to persistent storage; to **verify a save** is to close and reopen the file.

Try this safely on a copy of `tasks.json`. Make a table with three moments: before opening, after editing one fictional task, and after reopening. Before pressing Save, predict where the new version exists. Save, close the editor, and reopen the file. If the change appeared only before closing, the screen showed a temporary state. If it survived, the data reached storage. Use a copy so the experiment cannot damage the original project.

**Check your explanation.** A classmate says, “I can see the file, so it is saved.” Ask them to show a second read. Imagine cutting power before saving and after saving. The outcomes differ because screen, RAM, and storage have different jobs. This is a more useful explanation than any advertised number of gigabytes.""",
        },
    },
    6: {
        "ru": {
            "hook": "Иконка помощника осталась на экране, но после закрытия окно исчезло. Исчезла ли сама программа? Разберём запуск как последовательность проверяемых действий.",
            "bridge": """## Что именно запускается

**Файл программы** — сохранённые инструкции на накопителе; **приложение** — программа, предназначенная для задачи пользователя. **Процесс** — запущенный экземпляр программы с выделенными ресурсами; два открытых окна могут относиться к одному или нескольким процессам. **Операционная система (ОС)** — программа, которая распределяет процессорное время и память, открывает файлы и управляет доступом к устройствам. **Ресурс** здесь означает то, к чему процесс обращается: файл, память, камеру или сеть. **Разрешение** — право на конкретное действие, а не признак честности приложения.

Сыграйте запуск помощника на бумаге. Карточки: `приложение`, `ОС`, `процесс`, `tasks.json`, `экран`. Разложите: файл программы лежит на диске; ОС проверяет запуск и создаёт процесс; процесс просит прочитать `tasks.json`; ОС проверяет путь и право; процесс показывает результат. Уберите карточку `tasks.json`: процесс может запуститься, но не получить данные. Уберите право на чтение: файл существует, однако ОС откажет в доступе. Так вы различите «не найдено», «нет права» и «программа не запустилась».

В проекте запишите для каждой ошибки три строки: неудавшееся действие, наблюдаемое сообщение и первое безопасное исправление. Предоставить все разрешения сразу — плохая диагностика: это скрывает причину и расширяет доступ. **Принцип наименьших прав** означает давать лишь то право, без которого нужное действие не работает.""",
        },
        "kz": {
            "hook": "Көмекшінің белгішесі экранда қалды, бірақ жапқанда терезе жоғалды. Бағдарламаның өзі жоғалды ма? Іске қосуды тексерілетін қадамдарға бөлейік.",
            "bridge": """## Нақты не іске қосылады

**Бағдарлама файлы** — жинақтауыштағы сақталған нұсқаулар; **қолданба** — пайдаланушы міндетін орындайтын бағдарлама. **Процесс** — ресурстар бөлінген жұмыс істеп тұрған бағдарлама данасы; екі терезе бір немесе бірнеше процеске тиесілі болуы мүмкін. **Операциялық жүйе (ОЖ)** процессор уақытын және жадты бөліп, файлдарды ашып, құрылғыларға қолжетімділікті басқарады. **Ресурс** — процесс пайдаланатын файл, жад, камера немесе желі. **Рұқсат** — нақты әрекетке берілген құқық; ол қолданбаның қауіпсіз екенін дәлелдемейді.

Көмекшінің іске қосылуын қағаздағы карточкалармен көрсетіңіз: `қолданба`, `ОЖ`, `процесс`, `tasks.json`, `экран`. Бағдарлама файлы дискіде жатыр; ОЖ іске қосуды тексеріп, процесс құрады; процесс `tasks.json` файлын сұрайды; ОЖ жол мен құқықты тексереді; процесс нәтижені көрсетеді. `tasks.json` карточкасын алып тастаңыз: процесс іске қосылуы мүмкін, бірақ дерек табылмайды. Оқу құқығын алып тастаңыз: файл бар, бірақ ОЖ қолжетімділікті бермейді. Осылайша «табылмады», «рұқсат жоқ» және «бағдарлама ашылмады» дегенді ажыратасыз.

Жобада әр қатеге үш жол жазыңыз: қай әрекет өтпеді, қандай хабар көрінді және ең қауіпсіз алғашқы түзету қандай. Бірден барлық рұқсатты беру себепті жасырады. **Ең аз құқық қағидасы** қажетті әрекетке жететін құқықты ғана беруді білдіреді.""",
        },
        "en": {
            "hook": "The assistant icon stayed on screen, but its window vanished when you closed it. Did the program itself disappear? Let us break launch into observable steps.",
            "bridge": """## What actually starts

A **program file** is saved instructions on storage; an **application** is a program intended to serve a user's task. A **process** is one running instance with allocated resources; two windows may belong to one process or several. The **operating system (OS)** allocates processor time and memory, opens files, and controls access to devices. A **resource** here is something a process uses: a file, memory, camera, or network. A **permission** is the right to perform a particular action, not evidence that the app is trustworthy.

Act out the assistant's launch with cards labelled `application`, `OS`, `process`, `tasks.json`, and `screen`. The program file rests on storage; the OS checks launch and creates a process; the process asks for `tasks.json`; the OS checks the path and permission; the process displays the result. Remove the `tasks.json` card: the process may start but cannot find data. Remove read permission: the file exists, yet the OS denies access. You have now separated “not found,” “permission denied,” and “program did not start.”

In the project, give each failure three lines: the action that failed, the message observed, and the first safe repair. Giving every permission at once is poor diagnosis: it hides the cause and expands access. The **principle of least privilege** means granting only the right needed for the intended action.""",
        },
    },
    7: {
        "ru": {
            "hook": "Друг говорит: «Файл точно есть», а помощник отвечает: «Не найден». Это не магия и не каприз компьютера: у каждого файла есть точный адрес.",
            "bridge": """## Найдём файл, как квартиру по адресу

**Файл** — именованный набор данных. **Папка** группирует файлы и другие папки. **Путь** перечисляет папки от выбранной отправной точки до файла; например `project/data/tasks.json`. **Расширение** `.json` — часть имени, которая подсказывает формат, но не доказывает содержимое. **Текущая папка** — отправная точка для **относительного пути**; **абсолютный путь** начинается от корня файловой системы. Аналогия с адресом помогает, но компьютер не угадывает опечатки, как почтальон.

Нарисуйте дерево `project/` с папками `data/` и `notes/`. Поместите `tasks.json` в `data/`, а файл `plan.txt` — в `notes/`. Из `project/` путь к задачам равен `data/tasks.json`. Из `notes/` тот же файл находится по `../data/tasks.json`: `..` означает переход на одну папку вверх. Предскажите, что произойдёт, если открыть `data/task.json` без буквы `s`: это другой адрес, даже если человек понял замысел. Проверьте путь, не переименовывая исходник ради удобства.

В проекте запишите дерево папок и один точный путь к каждой из трёх вымышленных записей. Попросите другого учащегося найти файл только по вашему описанию. Если он спрашивает, где начинать, укажите текущую папку. Сообщение «файл не найден» исправляется проверкой адреса, а не установкой новой программы.""",
        },
        "kz": {
            "hook": "Досыңыз «файл анық бар» дейді, ал көмекші «табылмады» деп жауап береді. Компьютер қыңыр емес: әр файлдың нақты мекенжайы бар.",
            "bridge": """## Файлды мекенжаймен табайық

**Файл** — аты бар деректер жиыны. **Бума** файлдар мен басқа бумаларды топтайды. **Жол** таңдалған бастапқы орыннан файлға дейінгі бумаларды көрсетеді, мысалы `project/data/tasks.json`. `.json` **кеңейтімі** — форматты меңзейтін атаудың бөлігі, бірақ мазмұнның дәлелі емес. **Ағымдағы бума** — **салыстырмалы жолдың** басталу орны; **абсолют жол** файлдық жүйенің түбірінен басталады. Пошталық мекенжай ұқсатуы көмектеседі, бірақ компьютер адам сияқты қате жазуды жорамалдап түзетпейді.

`project/` ағашын салыңыз: ішінде `data/` мен `notes/` болсын. `tasks.json` файлын `data/` ішіне, `plan.txt` файлын `notes/` ішіне қойыңыз. `project/` бумасынан тапсырмаларға жол `data/tasks.json`. `notes/` бумасынан сол файлға жол `../data/tasks.json`: `..` бір деңгей жоғары шығуды білдіреді. `s` әрпі жоқ `data/task.json` ашылса не болатынын болжаңыз: бұл басқа мекенжай. Болжамды файлдың түпнұсқа атын өзгертпей тексеріңіз.

Жобада бума ағашын және үш ойдан шығарылған жазбаға жететін нақты жолдарды көрсетіңіз. Өзге оқушы файлды тек сипаттамаңыз бойынша тапсын. Ол «қайдан бастаймын?» десе, ағымдағы буманы көрсетіңіз. «Файл табылмады» қатесін жаңа бағдарлама орнатумен емес, жолды тексерумен түзетеміз.""",
        },
        "en": {
            "hook": "A friend insists that the file exists, yet the assistant says “not found.” The computer is not being stubborn: a file has an exact address.",
            "bridge": """## Find a file by its address

A **file** is a named collection of data. A **folder** groups files and other folders. A **path** lists the folders from a chosen starting point to the file, such as `project/data/tasks.json`. The `.json` **extension** is part of the name and hints at a format, but does not prove the contents. The **current folder** is the starting point for a **relative path**; an **absolute path** starts at the file system's root. A street-address analogy helps, but a computer does not guess misspellings like a helpful post worker.

Draw a `project/` tree with `data/` and `notes/`. Put `tasks.json` in `data/` and `plan.txt` in `notes/`. From `project/`, the task path is `data/tasks.json`. From `notes/`, the same file is at `../data/tasks.json`: `..` goes up one folder. Predict what happens if you open `data/task.json` without the `s`: that is a different address even if a person understands your intention. Check the path without renaming the source just to make the error disappear.

For the project, record the folder tree and one exact path to each of the three fictional records. Ask another learner to locate the file using only your description. If they ask where to start, state the current folder. Fix “file not found” by checking the address, not by installing another app.""",
        },
    },
    8: {
        "ru": {
            "hook": "Вы переименовали фотографию из `.jpg` в `.png`, но изображение не стало PNG. Почему смена этикетки не меняет содержимое?",
            "bridge": """## Этикетка, устройство файла и право автора

**Формат** — правило устройства данных внутри файла: где находится заголовок, какие значения допустимы и как их читать. **Расширение** — конец имени файла, полезная подсказка для ОС и человека, но не преобразователь. **Программа** умеет читать лишь форматы, для которых в ней есть соответствующие правила. **Конвертация** действительно переписывает данные в другой формат. **Лицензия** — разрешение автора на определённые способы использования; техническая возможность скопировать файл не означает права публиковать его.

Сделайте копию `tasks.json` под именем `tasks.txt`. Откройте оба файла как обычный текст: байты остались теми же. Переименуйте копию в `tasks.png`: программа просмотра изображений не получит из текста картинку. Для настоящего преобразования нужна программа, которая прочитает исходные данные и запишет новые по правилам PNG. Сравните хеши или размер файлов до и после переименования: имя меняется, содержимое нет. Затем намеренно испортите одну кавычку в JSON и объясните, почему расширение правильное, а содержимое уже не соответствует формату.

Для помощника заведите таблицу `объект — формат — чем открываем — откуда взят — разрешено ли распространять`. Для собственной вымышленной записи отметьте своё авторство. Для чужого изображения проверьте условия лицензии в первоисточнике прежде, чем добавлять его в проект.""",
        },
        "kz": {
            "hook": "Фотосуреттің `.jpg` жалғауын `.png` деп өзгерттіңіз, бірақ сурет PNG-ге айналмады. Неге жапсырманы ауыстыру ішкі мазмұнды өзгертпейді?",
            "bridge": """## Атау, файлдың ішкі ережесі және автор құқығы

**Формат** — файл ішіндегі деректің құрылыс ережесі: басы қайда, қандай мәндер рұқсат және олар қалай оқылады. **Кеңейтім** — файл атауының соңы; ОЖ мен адамға белгі береді, бірақ деректі түрлендірмейді. **Бағдарлама** өзіне белгілі форматтарды ғана дұрыс оқиды. **Түрлендіру** деректерді жаңа форматтың ережесімен қайта жазады. **Лицензия** — автордың белгілі пайдалану тәсілдеріне берген рұқсаты; файлды көшіре алу оны жариялауға құқық бермейді.

`tasks.json` файлының көшірмесін `tasks.txt` деп атаңыз. Екеуін де жай мәтін ретінде ашыңыз: ішіндегі байттар өзгермеді. Көшірмені `tasks.png` деп атағанда сурет қарау бағдарламасы мәтіннен кескін жасай алмайды. Шын түрлендіру үшін бастапқы деректі оқып, PNG ережесімен жаңа файл жазатын бағдарлама қажет. Атын өзгертпей тұрып және өзгерткеннен кейін файл өлшемін не хешін салыстырыңыз. Содан кейін JSON ішіндегі бір тырнақшаны бұзыңыз: кеңейтім дұрыс болғанымен, мазмұн форматқа сай емес.

Көмекші үшін `нысан — формат — ашатын бағдарлама — дереккөзі — таратуға рұқсат` кестесін жүргізіңіз. Өз ойдан шығарылған жазбаңызға өз авторлығыңызды белгілеңіз. Өзгенің суретін жобаға қоспас бұрын лицензия шартын бастапқы дереккөзден оқыңыз.""",
        },
        "en": {
            "hook": "You renamed a photo from `.jpg` to `.png`, but the picture did not become a PNG. Why does changing the label leave the contents untouched?",
            "bridge": """## Label, file rules, and the creator's permission

A **format** is a set of rules for data inside a file: where its header sits, which values are valid, and how to read them. An **extension** is the end of a filename; it hints at a format but does not convert it. A **program** can correctly read only formats it knows. **Conversion** actually rewrites the data under another format's rules. A **license** states which uses the creator allows; being able to copy a file does not give you permission to publish it.

Copy `tasks.json` as `tasks.txt`. Open both as plain text: the bytes remain the same. Rename the copy to `tasks.png`: an image viewer cannot turn the text into a picture. Real conversion needs software that reads the source and writes a new file using PNG rules. Compare a checksum or file size before and after renaming: the name changes, the contents do not. Next, break one quotation mark in the JSON and explain why the extension is still correct while the content no longer follows the format.

For the assistant, keep a table: `asset — format — program that opens it — source — permission to redistribute`. Mark your fictional records as your own work. Check the original source's license before adding another person's image to the project.""",
        },
    },
    9: {
        "ru": {
            "hook": "Два человека одновременно исправили одну задачу. Чью версию оставить и как понять, что изменилось? А если один из них не различает цвета на экране?",
            "bridge": """## Изменения должны быть видимы каждому

**Версия** — сохранённое состояние проекта в определённый момент. **Сравнение версий** показывает добавленные, удалённые и изменённые строки. **Конфликт** возникает, когда два человека меняют одну часть по-разному и выбор нельзя сделать автоматически. **Доступность** означает, что человек может выполнить задачу при разных способах восприятия и управления; цвет не должен быть единственным сигналом. **Фокус клавиатуры** — место, куда попадёт следующее нажатие без мыши.

Скопируйте `tasks.json` в `version-A.json` и `version-B.json`. В A поменяйте название `t-01`, в B — значение `done` у `t-01`. На бумаге составьте итоговую строку, которая сохраняет обе совместимые правки. Затем измените название `t-01` по-разному в A и B: теперь нужен разговор о смысле, а не слепое «последний сохранил — тот прав». Запишите решение в журнале и проверьте, что итоговый JSON открывается.

Покажите состояние задачи не только зелёным цветом, но и словом «готово» или «не готово». Попробуйте пройти проект клавишей Tab: видно ли, какой элемент выбран, можно ли открыть файл и понять сообщение об ошибке? Передайте файл товарищу и попросите повторить действие без ваших подсказок. Это одновременно проверка совместной работы и доступности.""",
        },
        "kz": {
            "hook": "Екі адам бір тапсырманы қатар түзетті. Қай нұсқаны қалдырамыз, айырманы қалай көреміз? Ал олардың бірі экрандағы түстерді ажыратпаса ше?",
            "bridge": """## Өзгеріс баршаға түсінікті болсын

**Нұсқа** — жобаның белгілі сәтте сақталған күйі. **Нұсқаларды салыстыру** қосылған, өшірілген және өзгерген жолдарды көрсетеді. **Қақтығыс** екі адам бір бөлікті әртүрлі өзгертіп, таңдауды автоматты түрде жасауға болмағанда туады. **Қолжетімділік** — адам ақпаратты әртүрлі қабылдап не құрылғыны әртүрлі басқарып, міндетті орындай алуы; түс жалғыз белгі болмауы керек. **Пернетақта фокусы** — тінтуірсіз басылған келесі перне әсер ететін орын.

`tasks.json` файлын `version-A.json` және `version-B.json` деп екі рет көшіріңіз. A нұсқасында `t-01` атауын, B нұсқасында `t-01` үшін `done` мәнін өзгертіңіз. Қағазда екі үйлесімді түзетуді сақтайтын қорытынды жол құрастырыңыз. Енді A мен B нұсқаларында атауды әртүрлі өзгертіңіз: «соңғы сақтаған дұрыс» деп емес, мағынасы туралы келісу керек. Шешімді журналға жазып, қорытынды JSON ашылатынын тексеріңіз.

Тапсырма күйін тек жасыл түспен емес, «дайын» немесе «дайын емес» деген сөзбен көрсетіңіз. Tab пернесімен жобаны аралаңыз: қай элемент таңдалғаны көріне ме, файлды ашуға және қате хабарын түсінуге бола ма? Басқа оқушы сіздің көмегіңізсіз әрекетті қайталап көрсін. Бұл бірлескен жұмыс пен қолжетімділікті қатар тексереді.""",
        },
        "en": {
            "hook": "Two people edited the same task at once. Which version should stay, and how do you see the difference? What if one person cannot distinguish the screen's colours?",
            "bridge": """## Make changes visible to everyone

A **version** is a saved state of the project at a particular time. A **diff** shows added, removed, and changed lines. A **conflict** occurs when two people change the same part differently and no automatic choice is justified. **Accessibility** means a person can complete the task with different ways of perceiving or controlling the interface; colour must not be the only signal. **Keyboard focus** is the place where the next key press will act without a mouse.

Copy `tasks.json` to `version-A.json` and `version-B.json`. In A, rename `t-01`; in B, change the `done` value of `t-01`. On paper, produce a final record containing both compatible edits. Now change the title differently in A and B: meaning must be discussed, not settled by “last save wins.” Record your decision and check that the final JSON opens correctly.

Show task state with words such as “done” or “not done” as well as colour. Move through the project with Tab: can you see which item has focus, open the file, and understand an error? Give the project to another learner and ask them to repeat the action without prompts. That tests collaboration and accessibility together.""",
        },
    },
    10: {
        "ru": {
            "hook": "Можно ли сказать «я освоил компьютер», если три файла открываются только на моём ноутбуке и только пока я подсказываю? Контрольная точка проверяет самостоятельную работу другого человека.",
            "bridge": """## Защита версии 0.1 как небольшое расследование

**Контрольная точка** — момент, когда мы проверяем результат по заранее известным признакам. **Артефакт** — оставленный работой файл, схема или таблица, которую можно открыть и проверить. **Воспроизводимость** означает, что другой человек получит тот же результат по вашей инструкции. **Перенос** — применение известного правила в новой ситуации, например на другом устройстве. **Критерий** — наблюдаемое условие зачёта, а не оценка «понравилось».

Дайте товарищу папку версии 0.1, но не касайтесь его устройства. Он должен: (1) найти `README`; (2) по нему открыть три вымышленные записи; (3) объяснить, почему изменение на экране ещё не доказательство сохранения; (4) назвать точный путь; (5) показать, какое право приложению действительно требуется; (6) прочесть состояние без опоры только на цвет. Если он спрашивает, где начать, это дефект инструкции, а не ученика. Исправьте README и повторите опыт с другим человеком.

Проверка знаний на 10 пунктов не заменяет работающий проект. Для зачёта покажите и ответы, и артефакты: дерево папок, таблицу ввода-вывода, журнал решений, реестр источников и воспроизводимую инструкцию. Через неделю решите аналогичную задачу с другим файлом; 7 из 10 — сигнал удержания знания, а не ярлык способности. Ошибка отправляет к конкретному шагу, который нужно пересобрать.""",
        },
        "kz": {
            "hook": "Үш файл тек өз ноутбугыңызда және сіз көрсетіп тұрғанда ғана ашылса, «компьютерді меңгердім» деуге бола ма? Бақылау нүктесі басқа адамның дербес жұмысын тексереді.",
            "bridge": """## 0.1 нұсқасын шағын зерттеу сияқты қорғау

**Бақылау нүктесі** — алдын ала белгілі белгілермен нәтижені тексеретін сәт. **Жұмыс нәтижесі** — ашып тексеруге болатын файл, сызба не кесте. **Қайталанымдылық** — өзге адамның сіздің нұсқауыңызбен дәл сондай нәтиже алуы. **Білімді көшіру** — таныс ережені жаңа жағдайда, мысалы басқа құрылғыда, қолдану. **Өлшем** — «ұнады» деген баға емес, бақыланатын қабылдау шарты.

0.1 нұсқасының бумасын досыңызға беріңіз, бірақ оның құрылғысына қол тигізбеңіз. Ол (1) `README` файлын табуы; (2) соған қарап үш ойдан шығарылған жазбаны ашуы; (3) экрандағы өзгеріс неге сақтау дәлелі емес екенін түсіндіруі; (4) нақты жолды атауы; (5) қолданбаға шынымен қандай құқық керек екенін көрсетуі; (6) күйді тек түске сүйенбей оқуы керек. «Қайдан бастаймын?» деген сұрақ — оқушының емес, нұсқаулықтың кемшілігі. README-ді түзетіп, басқа адаммен қайталаңыз.

Он тармақтық білім тексерісі жұмыс істейтін жобаны алмастырмайды. Жауаптармен бірге жұмыс нәтижесін көрсетіңіз: бума ағашы, енгізу-шығару кестесі, шешім журналы, дереккөз тізімі және қайта орындауға болатын нұсқаулық. Бір аптадан соң осыған ұқсас міндетті басқа файлмен орындаңыз. 7/10 — қабілетке таңба емес, білімнің сақталғанын көрсететін белгі. Қате нақты қай қадамды қалпына келтіру керегін көрсетеді.""",
        },
        "en": {
            "hook": "Can you say you understand the computer if three files open only on your own laptop and only while you give hints? This checkpoint tests another person's independent use.",
            "bridge": """## Defend version 0.1 as a small investigation

A **checkpoint** is a moment to test the result against known conditions. An **artifact** is a file, diagram, or table left by the work that someone else can inspect. **Reproducibility** means another person gets the same result from your instructions. **Transfer** is using a known rule in a new setting, such as another device. A **criterion** is an observable condition for success, not “I liked it.”

Give a classmate the version 0.1 folder but do not touch their device. They must (1) find `README`; (2) use it to open three fictional records; (3) explain why an on-screen change does not prove saving; (4) name an exact path; (5) show which permission the app really needs; and (6) read task state without relying only on colour. If they ask where to begin, the instruction needs repair. Revise README and repeat with someone else.

A ten-point quiz does not replace a working project. Show both answers and artifacts: the folder tree, input/output table, decision log, source register, and reproducible instructions. A week later solve an equivalent task with another file. A score of 7/10 signals retained knowledge, not a label for ability. Each error points to one step to rebuild.""",
        },
    },
    11: {
        "ru": {
            "hook": "Помощник должен различать «сделано» и «не сделано». Как передать эту разницу устройству, которое не понимает русских слов? Начнём с двух ясно различимых состояний.",
            "bridge": """## Смысл не живёт внутри нуля

**Состояние** — один из вариантов, в которых может находиться объект. **Код** — условная запись состояния по заранее принятому правилу. **Бит** — выбор между двумя различимыми состояниями, обычно записанный как 0 или 1. **Байт** — группа из восьми битов; у неё 2 × 2 × 2 × 2 × 2 × 2 × 2 × 2 = 256 сочетаний. Ноль является полноправным кодом, а не пустым местом. **Декодировать** значит применить правило и вернуть записи её смысл.

Напишите четыре карточки: `00`, `01`, `10`, `11`. Договоритесь: 00 — новая задача, 01 — в работе, 10 — выполнена; 11 пока не используется. Дайте карточки товарищу без легенды: он не сможет надёжно угадать смысл. Передайте легенду: теперь расшифровка воспроизводима. Добавьте пятое состояние «отложена»: четырёх кодов уже мало, понадобится третий бит. Не путайте число битов с числом символов слова `true` в JSON: текстовый формат хранит дополнительные знаки.

В паспорте помощника запишите соглашение `done: false/true` и рядом объясните, какие человеческие состояния оно отражает. Спросите, что делать с задачей «в работе»: новое поле или другое множество состояний? Один бит не может честно обозначить три варианта.""",
        },
        "kz": {
            "hook": "Көмекші «орындалды» мен «орындалмады» дегенді ажыратуы керек. Қазақша сөзді түсінбейтін құрылғыға осы айырманы қалай жеткіземіз? Екі айқын күйден бастайық.",
            "bridge": """## Мағына нөлдің ішінде өздігінен тұрмайды

**Күй** — нысанның болуы мүмкін жағдайының бірі. **Код** — алдын ала келісілген ережемен күйді жазу тәсілі. **Бит** — 0 не 1 деп жазылатын екі ажыратылатын күйдің бірін таңдау. **Байт** — сегіз биттің тобы; онда 2 × 2 × 2 × 2 × 2 × 2 × 2 × 2 = 256 комбинация бар. Нөл де жарамды код, бос орын емес. **Кодты ашу** — ережені қолданып, жазбаның мағынасын қалпына келтіру.

Төрт карточка жазыңыз: `00`, `01`, `10`, `11`. Келісім: 00 — жаңа тапсырма, 01 — орындалып жатыр, 10 — орындалды; 11 әзірге бос. Карточкаларды ережесіз досыңызға берсеңіз, ол мағынасын сенімді таба алмайды. Ережемен бірге берсеңіз, кодты аша алады. Бесінші күйді қосқанда төрт код жетпейді, үшінші бит керек. Бит санын JSON ішіндегі `true` сөзінің таңба санымен шатастырмаңыз: мәтіндік формат қосымша таңбаларды сақтайды.

Көмекшінің паспортында `done: false/true` келісімін және оның адамға түсінікті мағынасын жазыңыз. «Орындалып жатыр» деген үшінші күйге не керек: жаңа өріс пе, әлде басқа күй жүйесі ме? Бір бит үш бөлек мағынаны адал бере алмайды.""",
        },
        "en": {
            "hook": "The assistant must distinguish “done” from “not done.” How can we send that difference to a device that does not understand English words? Start with two reliably distinct states.",
            "bridge": """## Meaning is not inside the zero

A **state** is one possible condition of an object. A **code** records a state under an agreed rule. A **bit** chooses between two distinguishable states, usually written 0 or 1. A **byte** groups eight bits and has 2 × 2 × 2 × 2 × 2 × 2 × 2 × 2 = 256 combinations. Zero is a valid code, not an empty slot. To **decode** is to apply the rule and recover meaning.

Write four cards: `00`, `01`, `10`, `11`. Agree that 00 means new task, 01 in progress, and 10 done; leave 11 unused. Give the cards to a partner without the key: they cannot reliably guess the meanings. Add the key and decoding becomes repeatable. Add a fifth state: four codes are no longer enough, so a third bit is needed. Do not confuse bit count with the characters of `true` in JSON; a text format stores extra symbols.

In the assistant's passport, document `done: false/true` and the human states it represents. Ask how to show “in progress”: another field or a larger state set? One bit cannot faithfully represent three distinct choices.""",
        },
    },
    12: {
        "ru": {
            "hook": "На экране написано `1101`. Это тринадцать или тысяча сто один? Без указания системы счисления ответить нельзя.",
            "bridge": """## Место цифры меняет её цену

**Система счисления** задаёт допустимые цифры и стоимость разрядов. **Основание** — количество разных цифр: у десятичной системы 10, у двоичной 2. **Разряд** — место цифры; справа налево в двоичной записи он стоит 1, 2, 4, 8, 16 и так далее. Поэтому `1101₂` = 1 × 8 + 1 × 4 + 0 × 2 + 1 × 1 = 13₁₀. Нижний индекс уточняет основание; сам компьютер не видит индекс в каждом байте — формат сообщает, как читать данные.

Разложите на столе карточки с весами 8, 4, 2, 1. Чтобы получить 13, выберите 8 + 4 + 1 и положите над карточками `1 1 0 1`. Для 10 выберите 8 + 2: `1010`. Закройте пример и переведите 9 сами. Проверьте обратным сложением: `1001₂` = 8 + 1. Ошибка `1101₂ = 1 + 1 + 0 + 1 = 3` возникает, когда забывают веса разрядов.

В проекте не превращайте идентификатор задачи `t-01` в число: это текстовая метка. Запишите рядом число выполненных задач как отдельное числовое поле. У двух одинаково выглядящих строк разный смысл, если договорённость о типе различается.""",
        },
        "kz": {
            "hook": "Экранда `1101` тұр. Бұл он үш пе, әлде мың жүз бір ме? Санау жүйесін көрсетпей жауап бере алмаймыз.",
            "bridge": """## Цифрдың орны оның салмағын өзгертеді

**Санау жүйесі** рұқсат етілген цифрлар мен разряд салмағын белгілейді. **Негіз** — әртүрлі цифр саны: ондықта 10, екілікте 2. **Разряд** — цифр орны; екілік жазуда оңнан солға салмақтар 1, 2, 4, 8, 16 болып өседі. Сондықтан `1101₂` = 1 × 8 + 1 × 4 + 0 × 2 + 1 × 1 = 13₁₀. Төменгі индекс негізді нақтылайды; компьютер әр байттан индексті көрмейді, қалай оқу керегін формат айтады.

Үстелге 8, 4, 2, 1 салмақтары бар карточкалар қойыңыз. 13 алу үшін 8 + 4 + 1 таңдаңыз да, үстіне `1 1 0 1` жазыңыз. 10 үшін 8 + 2: `1010`. Мысалды жауып, 9 санын өзіңіз аударыңыз. Кері қосумен тексеріңіз: `1001₂` = 8 + 1. `1101₂ = 1 + 1 + 0 + 1 = 3` қатесі разряд салмағын ескермегенде шығады.

Жобада `t-01` тапсырма белгісін санға айналдырмаңыз: бұл мәтіндік идентификатор. Орындалған тапсырмалар санын бөлек сандық өріске жазыңыз. Түрі туралы келісім өзгеше болса, бірдей көрінетін таңбалардың мағынасы да өзгереді.""",
        },
        "en": {
            "hook": "The screen shows `1101`. Is that thirteen or one thousand one hundred and one? You need to know the number system first.",
            "bridge": """## A digit's place changes its value

A **number system** defines allowed digits and the values of places. Its **base** is the number of distinct digits: decimal has 10 and binary has 2. A **place** is a digit position; binary places from right to left have weights 1, 2, 4, 8, 16, and so on. Thus `1101₂` = 1 × 8 + 1 × 4 + 0 × 2 + 1 × 1 = 13₁₀. The subscript clarifies the base. The computer does not find such a subscript in every byte; the format tells a reader how to interpret the data.

Lay cards labelled 8, 4, 2, 1 on a table. To make 13, choose 8 + 4 + 1 and place `1 1 0 1` above them. To make 10, choose 8 + 2: `1010`. Hide the example and convert 9 yourself. Check by adding backwards: `1001₂` = 8 + 1. The mistake `1101₂ = 1 + 1 + 0 + 1 = 3` ignores place values.

In the project, do not turn task ID `t-01` into a number: it is a text label. Store the count of completed tasks in a separate numeric field. Similar-looking symbols can have different meanings under different type rules.""",
        },
    },
    13: {
        "ru": {
            "hook": "Одноклассник отправил файл с названием «Ән», но у друга появились непонятные значки. Разберём, на каком переходе от буквы к байтам потерялось правило.",
            "bridge": """## Буква, номер и байты — три разных вещи

**Символ** — единица текста, например `Ә`. **Кодовая точка Unicode** — назначенный символу номер; для `Ә` это U+04D8. **Кодировка UTF-8** задаёт, как этот номер записать байтами: U+04D8 становится `D3 98`. **Декодирование** превращает байты обратно в символ по той же кодировке. **Глиф** — видимая форма буквы в шрифте; если шрифт не содержит нужной формы, байты и кодировка могут быть правильными, а на экране появится пустой квадрат.

Проведите три проверки в порядке. (1) Скопируйте `Ә` и `A` в два поля: выглядят ли они одинаково? Нет. (2) Найдите их номера: U+04D8 и U+0041. (3) Сохраните `Ә` в UTF-8 и посмотрите байты: `D3 98`, тогда как латинская `A` занимает `41`. Если файл открылся с «кракозябрами», сначала проверьте кодировку чтения. Если вместо буквы квадрат при верных байтах, проверьте шрифт. Не называйте любую ошибку текста «плохим Unicode».

Добавьте вымышленную задачу с казахской буквой в `tasks.json`, закройте и вновь откройте файл. Если буква сохранилась, запись и чтение согласованы. В `FORMAT.md` запишите UTF-8 явно, чтобы другой человек не гадал.""",
        },
        "kz": {
            "hook": "Сыныптасыңыз «Ән» деген файл жіберді, ал досының экранында түсініксіз таңбалар шықты. Әріптен байтқа дейінгі қай қадамда ереже жоғалғанын анықтайық.",
            "bridge": """## Әріп, нөмір және байт — үш бөлек нәрсе

**Таңба** — мәтін бірлігі, мысалы `Ә`. **Unicode код нүктесі** — таңбаға берілген нөмір; `Ә` үшін ол U+04D8. **UTF-8 кодтауы** осы нөмірді байттармен жазу ережесі: U+04D8 мәні `D3 98` болады. **Кодты ашу** — сол ережемен байттан таңбаны қайта алу. **Глиф** — қаріптегі көрінетін әріп пішіні; қаріпте пішін болмаса, байт пен кодтау дұрыс болса да, экранда бос шаршы көрінуі мүмкін.

Үш тексерісті ретімен жасаңыз. (1) `Ә` мен `A` таңбаларын салыстырыңыз: бұлар бір әріп емес. (2) Нөмірлерін табыңыз: U+04D8 және U+0041. (3) `Ә` таңбасын UTF-8 түрінде сақтап, `D3 98` байттарын қараңыз; латынша `A` үшін `41`. Файл қате таңбалармен ашылса, әуелі оқу кодтауын тексеріңіз. Байт дұрыс болып, шаршы көрінсе, қаріпті тексеріңіз. Кез келген мәтін қатесін «Unicode нашар» деп атамаңыз.

`tasks.json` ішіне қазақ әрпі бар ойдан шығарылған тапсырма қосып, файлды жауып қайта ашыңыз. Әріп сақталса, жазу мен оқу ережелері үйлескен. Басқа адам болжауға мәжбүр болмас үшін `FORMAT.md` файлына UTF-8 талабын анық жазыңыз.""",
        },
        "en": {
            "hook": "A classmate sent a file named “Ән,” but their friend's screen showed strange marks. Find which step from character to bytes lost its rule.",
            "bridge": """## Character, number, and bytes are different things

A **character** is a unit of text such as `Ә`. A **Unicode code point** is its assigned number; `Ә` is U+04D8. **UTF-8 encoding** defines how to store that number as bytes: U+04D8 becomes `D3 98`. **Decoding** applies the same rule to recover a character. A **glyph** is the visible shape in a font; if the font lacks it, the bytes and encoding can be correct while the screen shows an empty box.

Check three things in order. (1) Compare `Ә` and `A`: they are different letters. (2) Look up their numbers: U+04D8 and U+0041. (3) Save `Ә` as UTF-8 and inspect its bytes, `D3 98`; Latin `A` uses `41`. If text opens as nonsense, check the decoding first. If the bytes are right but a box appears, check the font. Do not call every display problem “bad Unicode.”

Add a fictional task with a Kazakh letter to `tasks.json`, close the file, and reopen it. If the letter survives, writing and reading agree. State UTF-8 explicitly in `FORMAT.md` so another learner need not guess.""",
        },
    },
    14: {
        "ru": {
            "hook": "Фото задачи увеличили в десять раз, и его края стали ступенчатыми. Что именно увеличилось — количество деталей или размер уже существующих квадратиков?",
            "bridge": """## Цветной квадрат и размер изображения

**Пиксель** — ячейка цифрового изображения с заданным цветом. **Разрешение изображения** — число пикселей по ширине и высоте, например 4 × 3 = 12 пикселей; это не физические сантиметры экрана. В модели **RGB** цвет задают тремя каналами: красным, зелёным и синим. При восьми битах на канал число от 0 до 255 показывает интенсивность. `(255, 0, 0)` — красный, `(0, 0, 0)` — чёрный, `(255, 255, 255)` — белый. **Масштабирование** растягивает или пересчитывает имеющиеся пиксели; оно не восстанавливает не записанные камерой детали.

Нарисуйте сетку 4 × 3. Верхнюю строку заполните красным, среднюю зелёным, нижнюю синим. Получится 12 пикселей, а не три: цвет строки повторяется в четырёх ячейках. Закрасьте одну ячейку белым и объясните, какие три значения меняются. Затем нарисуйте 8 × 6, просто удвоив каждый старый пиксель в обе стороны. Ячеек стало 48, но новых сведений о сцене не появилось.

Для помощника сделайте иконку задачи сначала в сетке 4 × 3, затем сравните её читаемость с более детальной сеткой. Если смысл зависит только от красного и зелёного, добавьте подписи: доступность касается и изображений.""",
        },
        "kz": {
            "hook": "Тапсырма суретін он есе үлкейткенде шеттері сатыланып кетті. Не көбейді: жаңа бөлшек пе, әлде бұрынғы кішкене шаршылардың өлшемі ме?",
            "bridge": """## Түсті ұяшық және кескін өлшемі

**Пиксель** — түсі берілген цифрлық кескін ұяшығы. **Кескін ажыратымдылығы** — ені мен биіктігіндегі пиксель саны, мысалы 4 × 3 = 12 пиксель; бұл экранның сантиметрі емес. **RGB** үлгісінде түс үш арнамен беріледі: қызыл, жасыл, көк. Әр арнаға сегіз бит берілсе, 0–255 саны оның қарқындылығын көрсетеді. `(255, 0, 0)` — қызыл, `(0, 0, 0)` — қара, `(255, 255, 255)` — ақ. **Масштабтау** бар пиксельдерді созады не қайта есептейді; камера жазбаған бөлшекті қалпына келтірмейді.

4 × 3 тор салыңыз. Жоғарғы қатарды қызыл, ортаны жасыл, төменгіні көк түске бояңыз. Үш емес, 12 пиксель шықты: әр қатардың түсі төрт ұяшықта қайталанады. Бір ұяшықты аққа бояп, үш санның қалай өзгеретінін айтыңыз. Енді әр ескі пиксельді екі бағытта екі еселеп, 8 × 6 тор салыңыз. 48 ұяшық болды, бірақ көрініс туралы жаңа мәлімет пайда болған жоқ.

Көмекшіге 4 × 3 торда тапсырма белгісін жасап, кейін егжей-тегжейлі тормен оқылуын салыстырыңыз. Мағына тек қызыл мен жасылға тәуелді болса, сөздік белгі қосыңыз: қолжетімділік кескінге де қатысты.""",
        },
        "en": {
            "hook": "A task photo was enlarged tenfold and its edges became blocky. Did we gain detail, or merely enlarge the existing little squares?",
            "bridge": """## A coloured cell and image size

A **pixel** is one cell of a digital image with an assigned colour. **Image resolution** counts pixels across and down, for example 4 × 3 = 12 pixels; it is not the screen's physical size. The **RGB** model uses three channels: red, green, and blue. With eight bits per channel, each value from 0 to 255 sets an intensity. `(255, 0, 0)` is red, `(0, 0, 0)` black, and `(255, 255, 255)` white. **Scaling** stretches or recalculates existing pixels; it cannot recover details the camera never recorded.

Draw a 4 × 3 grid. Colour the top row red, the middle green, and the bottom blue. There are 12 pixels, not three: each row colour appears in four cells. Change one cell to white and explain which three values change. Now draw an 8 × 6 grid by doubling each old pixel in both directions. You have 48 cells but no new information about the photographed scene.

Design a task icon for the assistant on a 4 × 3 grid and compare its legibility with a more detailed grid. If its meaning depends only on red versus green, add words too: accessibility applies to images.""",
        },
    },
    15: {
        "ru": {
            "hook": "Голос звучит непрерывно, а файл хранит конечное число значений. Как превратить плавную волну в числа и что при этом можно потерять?",
            "bridge": """## Волна, измерения и кадры

**Звуковая волна** — изменение давления воздуха во времени. **Выборка** — измеренное значение в один момент; **частота дискретизации** — сколько таких измерений делают за секунду. Если шаг слишком редок, быстрые изменения теряются. **Амплитуда** — величина отклонения, связанная с громкостью, но не равная субъективному ощущению. **Кадр** — отдельное изображение в видеоряде; **частота кадров** — число кадров за секунду. Звук и видео поэтому хранят разные ряды измерений, которые проигрыватель согласует по времени.

Нарисуйте волну, проходящую через значения 0, 1, 0, −1, 0 за секунду. Отметьте измерения в 0, 0.25, 0.5, 0.75 и 1 секунды: получите пять чисел. Теперь измерьте только в 0, 0.5 и 1: выйдут одни нули, и движение волны исчезнет из записи. Это учебный пример недостаточно частой выборки, а не доказательство, что всякая реальная запись с таким шагом молчит. Для видео нарисуйте четыре кадра, где точка сдвигается; переставьте кадры и объясните изменение движения.

Помощнику голосовой ввод пока не нужен. В паспорте напишите, какую задачу он решил бы и какие новые данные пришлось бы сохранять. Не включайте микрофон лишь потому, что у телефона он есть.""",
        },
        "kz": {
            "hook": "Дауыс үздіксіз естіледі, ал файлда санаулы мән сақталады. Бірқалыпты толқынды санға қалай айналдырамыз және не жоғалуы мүмкін?",
            "bridge": """## Толқын, өлшем және кадр

**Дыбыс толқыны** — уақыт өте өзгеретін ауа қысымы. **Үлгі** — бір сәтте өлшенген мән; **дискреттеу жиілігі** — секундтағы осындай өлшемдер саны. Өлшемдер тым сирек болса, тез өзгерістер жоғалады. **Амплитуда** — ауытқу шамасы; ол дыбыс қаттылығына қатысты, бірақ адамның сезімімен бірдей емес. **Кадр** — бейнеқатардағы жеке кескін; **кадр жиілігі** — секундтағы кадр саны. Дыбыс пен видео әртүрлі өлшем қатарын сақтайды, ойнатқыш оларды уақыт бойынша үйлестіреді.

Бір секундта 0, 1, 0, −1, 0 мәндерінен өтетін толқын салыңыз. 0, 0.25, 0.5, 0.75 және 1 секундтағы өлшемдерді белгілеңіз: бес сан шығады. Енді тек 0, 0.5 және 1 мезеттерін өлшеңіз: тек нөл қалады, қозғалыс жазбада көрінбейді. Бұл — сирек өлшеудің оқу мысалы; әр шынайы дыбыс осылай жоғалады деген сөз емес. Видео үшін нүктесі жылжитын төрт кадр салыңыз; кадр ретін ауыстырып, қозғалыстың қалай өзгеретінін түсіндіріңіз.

Көмекшіге әзірге дауыспен енгізу қажет емес. Паспортта ол қандай міндетті шешерін және қандай жаңа дерек сақталатынын жазыңыз. Телефонда микрофон бар екен деп оны қоспаңыз.""",
        },
        "en": {
            "hook": "A voice sounds continuous, yet a file contains a finite list of values. How do we turn a smooth wave into numbers, and what might be lost?",
            "bridge": """## Wave, samples, and frames

A **sound wave** is changing air pressure over time. A **sample** is one measured value at one moment; the **sampling rate** counts such measurements per second. If samples are too sparse, fast changes can be missed. **Amplitude** is the size of a deviation; it relates to loudness but is not identical to perceived volume. A **frame** is one image in a video sequence; **frame rate** counts frames per second. Audio and video therefore store different sequences that the player synchronises in time.

Draw a wave passing through 0, 1, 0, −1, 0 over one second. Mark measurements at 0, 0.25, 0.5, 0.75, and 1 seconds: you get five numbers. Measure only at 0, 0.5, and 1: all readings are zero, and the motion vanishes from this recording. This is a teaching example of sampling too sparsely, not a claim that every real recording at that rate is silent. For video, draw four frames with a moving dot; reorder them and explain how the motion changes.

The assistant does not need voice input yet. In its passport, state which problem voice would solve and what new data would have to be stored. Do not enable a microphone just because the phone has one.""",
        },
    },
    16: {
        "ru": {
            "hook": "Почему ZIP-архив возвращает исходный текст точно, а отправленная через мессенджер фотография может стать менее чёткой? Оба файла стали меньше, но обещания у методов разные.",
            "bridge": """## Два смысла слова «сжать»

**Сжатие без потерь** позволяет восстановить исходные байты точно; так можно хранить документы и программы. **Сжатие с потерями** отбрасывает часть деталей, считая их менее заметными, поэтому точное восстановление невозможно. **Коэффициент сжатия** сравнивает объёмы до и после, но не измеряет качество. **Архив** — контейнер для одного или нескольких файлов; внутри него могут применяться методы сжатия. Смена расширения не выполняет ни один из этих методов.

Для учебной строки `ААААБББ` правило «число повторений + символ» даёт `4А3Б`. Чтобы разжать, прочитайте «четыре А, три Б» и восстановите исходную строку. На коротком тексте служебные правила могут занять больше места, поэтому не обещайте экономию всегда. Теперь представьте фотографию с оттенками 200 и 201, которые упрощённый метод округляет до 200: файл может стать меньше, но прежний 201 не вернуть. Это модель потерь, а не буквальное описание JPEG.

Скопируйте `tasks.json`, упакуйте копию в ZIP и распакуйте. Сравните её с исходником побайтно или через SHA-256. Если совпало, это доказательство точного восстановления для данного опыта. Фотографию проверяйте иначе: визуально и по размеру, не требуя побайтного совпадения от формата с потерями.""",
        },
        "kz": {
            "hook": "Неге ZIP мұрағаты бастапқы мәтінді дәл қайтарады, ал мессенджердегі сурет бұлыңғырлауы мүмкін? Екеуі де кішірейді, бірақ тәсілдердің уәдесі бөлек.",
            "bridge": """## «Сығудың» екі түрлі уәдесі

**Шығынсыз сығу** бастапқы байттарды дәл қалпына келтіреді; құжат пен бағдарламаларға керек. **Шығынды сығу** кей бөлшектерді елеусіз деп алып тастайды, сондықтан дәл кері қайтару мүмкін емес. **Сығу қатынасы** бұрынғы және кейінгі көлемді салыстырады, сапаны өлшемейді. **Мұрағат** — бір не бірнеше файлға арналған контейнер; ішінде сығу тәсілі қолданылуы мүмкін. Кеңейтімді ауыстыру бұл әрекеттерді орындамайды.

Оқу жолы `ААААБББ` үшін «қайталану саны + таңба» ережесі `4А3Б` береді. Кері ашқанда «төрт А, үш Б» деп бастапқы жолды қайтарыңыз. Қысқа мәтінде қосымша ережелердің өзі көп орын алуы мүмкін: үнем әрдайым бола бермейді. Енді суреттегі 200 және 201 реңктерін 200-ге дөңгелектейтін қарапайым тәсілді ойлаңыз: файл азаюы мүмкін, бірақ бұрынғы 201 мәнін қайтара алмаймыз. Бұл JPEG-тің дәл сипаттамасы емес, шығынды түсіндіретін үлгі.

`tasks.json` көшірмесін ZIP-ке салып, қайта шығарыңыз. Оны түпнұсқамен байт бойынша не SHA-256 арқылы салыстырыңыз. Сәйкестік — осы тәжірибеде дәл қалпына келгенінің дәлелі. Суретті көлемі мен көрінісі арқылы да бағалаңыз; шығынды форматтан байттық теңдік күтпеңіз.""",
        },
        "en": {
            "hook": "Why does a ZIP archive return the original text exactly, while a photo sent through a messenger may look softer? Both became smaller, but the methods make different promises.",
            "bridge": """## Two promises behind “compression”

**Lossless compression** lets you reconstruct the original bytes exactly; documents and programs need this. **Lossy compression** discards some detail considered less noticeable, so exact recovery is impossible. A **compression ratio** compares sizes before and after, but does not measure quality. An **archive** is a container for one or more files and may use compression internally. Renaming an extension performs neither method.

For the teaching string `ААААБББ`, the rule “count plus character” gives `4А3Б`. Decode it as “four А, three Б” and restore the original. On short text, the rules themselves can take extra space, so do not promise savings every time. Now imagine a photo containing shades 200 and 201 that a simple method rounds both to 200: size may fall, but the former 201 cannot be recovered. This models loss; it is not a literal description of JPEG.

Copy `tasks.json`, put the copy in a ZIP file, and extract it. Compare the extracted and original files byte for byte or by SHA-256. A match proves exact restoration for this experiment. Assess a photo by visible quality and size as well; do not expect byte equality from a lossy format.""",
        },
    },
    17: {
        "ru": {
            "hook": "Два файла открываются без ошибки, но в одном задача уже не выполнена. Как доказать, что копия изменилась, если внешний вид почти тот же?",
            "bridge": """## Открывается — ещё не значит совпадает

**Целостность** означает, что данные не изменились вопреки ожиданию. **Контрольная сумма** — короткий результат вычисления по содержимому файла; изменение байтов обычно меняет результат. **Хеш SHA-256** — один вид такой проверки, но совпадение хешей без надёжного исходного значения не доказывает, кто создал файл. **Синтаксис JSON** — правила допустимой записи; **схема данных** — договор о нужных полях и их типах; **смысловая проверка** выясняет, соответствует ли запись задаче.

Возьмите две копии `tasks.json`. В одной поменяйте `done: true` на `done: false`, сохранив правильные запятые и кавычки. Оба файла могут проходить синтаксическую проверку. Сравните хеши: они различаются. Затем поменяйте название задачи и верните обратно: конечные файлы снова совпадают, хотя история действий различалась. Хеш отвечает о текущих байтах, не рассказывает всю историю.

Для помощника настройте три шага: открыть JSON, проверить обязательные поля `id`, `title`, `done` и их типы, затем проверить ожидаемые три вымышленные записи. Резервную копию сравните с исходником сразу после создания. Если её собственный хеш записан уже после порчи, проверка не обнаружит прежнюю потерю: нужен заранее сохранённый эталон.""",
        },
        "kz": {
            "hook": "Екі файл да қатесіз ашылады, бірақ біреуінде тапсырма енді орындалмаған. Сыртқы көрініс ұқсас болса, көшірме өзгергенін қалай дәлелдейміз?",
            "bridge": """## Ашылуы — сәйкестік дәлелі емес

**Тұтастық** — деректің күтпеген түрде өзгермеуі. **Бақылау қосындысы** — файл мазмұнынан есептелетін қысқа нәтиже; байт өзгерсе, ол әдетте өзгереді. **SHA-256 хеші** — осындай тексерудің бір түрі, бірақ бастапқы сенімді мән болмаса, хештің сәйкестігі файл авторын дәлелдемейді. **JSON синтаксисі** — жазбаның дұрыс құрылу ережесі; **дерек сұлбасы** — қажет өрістер мен түрлері туралы келісім; **мағыналық тексеріс** — жазбаның міндетке сай келетінін анықтау.

`tasks.json` файлының екі көшірмесін алыңыз. Біреуінде үтір мен тырнақшаны сақтап, `done: true` мәнін `done: false` деп өзгертіңіз. Екі файл да синтаксистік тексерістен өтуі мүмкін. Хештерін салыстырыңыз: айырмашылық бар. Енді атауды өзгертіп, қайта орнына келтіріңіз: соңғы файлдар тең болуы мүмкін, бірақ әрекеттер тарихы басқа. Хеш қазіргі байтты көрсетеді, бүкіл тарихты емес.

Көмекшіде үш қадам жасаңыз: JSON ашу, `id`, `title`, `done` өрістерін және түрлерін тексеру, содан кейін күтілген үш ойдан шығарылған жазбаны тексеру. Қор көшірмесін жасаған бойда түпнұсқамен салыстырыңыз. Көшірме бүлінгеннен кейін ғана оның хешін жазсаңыз, бұрынғы жоғалтуды байқамайсыз: алдын ала сақталған үлгі қажет.""",
        },
        "en": {
            "hook": "Two files open without an error, but one says the task is no longer done. How can you prove the copy changed when it looks almost the same?",
            "bridge": """## Opening is not the same as matching

**Integrity** means data has not changed unexpectedly. A **checksum** is a compact result computed from file contents; changing bytes usually changes it. A **SHA-256 hash** is one such check, but matching hashes without a trusted original value do not prove who created the file. **JSON syntax** gives the rules for valid writing; a **data schema** states which fields and types are needed; a **semantic check** asks whether the record makes sense for the task.

Make two copies of `tasks.json`. In one, change `done: true` to `done: false` while keeping commas and quotes correct. Both may still pass a syntax check. Compare hashes: they differ. Next, change a title and restore it: the final files can match again even though their histories differ. A hash describes current bytes, not every action that happened.

For the assistant, perform three checks: open JSON, verify required `id`, `title`, and `done` fields with their types, then verify the three expected fictional records. Compare a backup with its source immediately after creation. If you record its own hash only after corruption, you cannot detect the earlier loss; you need a value saved in advance.""",
        },
    },
    18: {
        "ru": {
            "hook": "Теперь другой ученик должен получить ваши данные и понять их без разговора с вами. Сможет ли он открыть казахское название, добавить задачу и заметить повреждение копии?",
            "bridge": """## Защитим договор о данных

**Спецификация формата** — письменный договор, какие файлы и поля допустимы. **Тип** ограничивает значение: `done` — логическое `true` или `false`, а не строка `"true"`; `title` — непустой текст. **Идентификатор** `id` различает записи даже при одинаковых названиях. **Версия формата** `0.2` сообщает, по каким правилам читать файл; это не число выполненных задач. **Обратная совместимость** означает, что новая версия читает прежние допустимые записи, если это прямо обещано.

Передайте товарищу только `tasks.json` и `FORMAT.md`. Он должен без подсказок добавить вымышленную задачу `t-04` с буквой `Ә`, выбрать `done: false`, сохранить в UTF-8 и объяснить, почему `"false"` в кавычках не годится. Затем он меняет один байт в копии: проверка целостности должна обнаружить отличие. Если человек вынужден гадать о значении `version` или о допустимом `id`, дополните договор и повторите опыт.

При проверке 8/10 обязательны три вида доказательств: правильно читаемый файл, перенос на новой записи и совпадение резервной копии с исходником до намеренной порчи. Семидневная проверка проводится на другом примере. Показать только скриншот недостаточно: он не доказывает ни кодировку, ни типы, ни сохранность байтов.""",
        },
        "kz": {
            "hook": "Енді басқа оқушы сіздің деректеріңізді сізбен сөйлеспей түсінуі керек. Ол қазақша атауды ашып, тапсырма қосып, бүлінген көшірмені байқай ала ма?",
            "bridge": """## Дерек туралы келісімді қорғайық

**Формат сипаттамасы** — қандай файл мен өріс рұқсат екенін жазған келісім. **Түр** мәнді шектейді: `done` — `"true"` жолы емес, логикалық `true` не `false`; `title` — бос емес мәтін. **Идентификатор** `id` атаулары бірдей жазбаларды ажыратады. **Формат нұсқасы** `0.2` файлды қай ережемен оқу керегін көрсетеді; бұл орындалған тапсырмалар саны емес. **Кері үйлесімділік** — жаңа нұсқа бұрынғы дұрыс жазбаларды оқи алады деген айқын уәде.

Досыңызға тек `tasks.json` және `FORMAT.md` беріңіз. Ол көмексіз `Ә` әрпі бар ойдан шығарылған `t-04` тапсырмасын қосып, `done: false` таңдап, UTF-8 түрінде сақтап, неге тырнақшадағы `"false"` жарамсыз екенін түсіндірсін. Содан кейін көшірмедегі бір байтты өзгертсін: тұтастық тексерісі айырманы табуы тиіс. Адам `version` мағынасын немесе `id` ережесін болжауға мәжбүр болса, келісімді толықтырып, қайталаңыз.

8/10 үшін үш дәлел міндетті: дұрыс оқылатын файл, жаңа жазбаға білімді көшіру және әдейі бүлдірмей тұрып қор көшірмесінің түпнұсқамен теңдігі. Жеті күннен кейін басқа мысалмен тексеріңіз. Жалғыз скриншот кодтауды, түрлерді және байттардың сақталғанын дәлелдемейді.""",
        },
        "en": {
            "hook": "Another learner must now receive your data and understand it without talking to you. Can they open a Kazakh title, add a task, and notice a damaged copy?",
            "bridge": """## Defend the data contract

A **format specification** is a written agreement about valid files and fields. A **type** limits a value: `done` is Boolean `true` or `false`, not the string `"true"`; `title` is nonempty text. An **identifier** `id` distinguishes records even if their titles match. **Format version** `0.2` says which rules to use to read the file; it is not a task count. **Backward compatibility** means a new version reads older valid records, if that promise has been stated.

Give a partner only `tasks.json` and `FORMAT.md`. Without hints, they must add fictional task `t-04` with `Ә`, choose `done: false`, save in UTF-8, and explain why quoted `"false"` is wrong. They then change one byte in a copy: an integrity check should reveal the difference. If they must guess what `version` means or which `id` is valid, improve the contract and retry.

For an 8/10 pass, three forms of evidence are essential: a readable file, transfer to a new record, and a backup matching the original before deliberate damage. Test again on a different example after seven days. A screenshot alone proves neither encoding, types, nor preservation of bytes.""",
        },
    },
}
