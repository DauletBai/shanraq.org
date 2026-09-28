#!/usr/bin/env python3
"""Generate the trilingual Cloud & DevOps lessons and their small support maps."""
from pathlib import Path
from html import escape
import re

from cloud_devops_pedagogy import PACKS, PREFACE_TEACHING

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / "course/lessons/cloud-devops"
MAPS = ROOT / "web/static/course/cloud-devops"

LOCALE = {
    "ru": {
        "lead": "Лид (summary)", "outcome": "Результат урока", "why": "Почему это важно",
        "practice": "Практика", "read": "Разберите результат", "check": "Проверка",
        "task": "Задание", "required": "Обязательное.", "optional": "По желанию.",
        "map": "Опорная схема", "run": "Выполните команды из корня `course/cloud-devops-lab`.",
        "safe": "Команды рассчитаны на учебную среду. Перед командой, меняющей сервер или данные, прочитайте её целиком и проверьте текущий каталог.",
        "done": "Не переходите дальше, пока не можете объяснить, что проверяет каждая команда и какой отказ она обнаруживает.",
    },
    "kz": {
        "lead": "Лид (summary)", "outcome": "Сабақ нәтижесі", "why": "Бұл не үшін маңызды",
        "practice": "Тәжірибе", "read": "Нәтижені талдаңыз", "check": "Тексеру",
        "task": "Тапсырма", "required": "Міндетті.", "optional": "Қалауыңызша.",
        "map": "Тірек сызба", "run": "Пәрмендерді `course/cloud-devops-lab` түбірінен орындаңыз.",
        "safe": "Пәрмендер оқу ортасына арналған. Серверді не деректі өзгертетін пәрменді іске қоспас бұрын оны толық оқып, ағымдағы қалтаны тексеріңіз.",
        "done": "Әр пәрмен нені тексеретінін және қандай ақауды табатынын түсіндіре алмайынша келесі сабаққа өтпеңіз.",
    },
    "en": {
        "lead": "Lead (summary)", "outcome": "Lesson outcome", "why": "Why this matters",
        "practice": "Practice", "read": "Read the result", "check": "Verification",
        "task": "Exercise", "required": "Required.", "optional": "Optional.",
        "map": "Support map", "run": "Run the commands from the `course/cloud-devops-lab` root.",
        "safe": "The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.",
        "done": "Do not move on until you can explain what every command verifies and which failure it can reveal.",
    },
}


PREFACE = {
"ru": """# Прежде чем начать: почему Cloud & DevOps сейчас

_Лид (summary):_ **Подробное введение в бесплатный практический курс: почему инфраструктурные навыки растут вместе с облаками, ИИ и дата-центрами, что строит Казахстан и чему именно вы научитесь без обещаний мгновенной профессии.**

## Почему мы выбрали этот курс

Облачный сервис выглядит нематериальным, но всегда работает на реальных процессорах, дисках, сетях и электроснабжении. Чем больше компании используют ИИ, потоковую обработку, государственные цифровые услуги и онлайн-продукты, тем больше им нужны вычислительные мощности и люди, способные безопасно выпускать, наблюдать и восстанавливать программные системы.

Международное энергетическое агентство оценило, что мировое потребление электричества дата-центрами выросло на 17% в 2025 году. В его базовом прогнозе оно почти удвоится: с 485 ТВт·ч в 2025 году до 950 ТВт·ч в 2030-м, а потребление ИИ-ориентированных центров утроится. Это не прогноз количества вакансий, но сильный физический признак расширения инфраструктуры. [Источник: IEA, Key Questions on Energy and AI](https://www.iea.org/reports/key-questions-on-energy-and-ai/executive-summary).

Программный слой тоже стал зрелым. По опросу CNCF за 2025 год, 82% опрошенных пользователей контейнеров запускают Kubernetes в production, а 66% организаций, размещающих генеративные модели, используют Kubernetes хотя бы для части inference-нагрузки. Опрос относится к участникам cloud native экосистемы и не описывает все компании мира, но показывает, что контейнерная эксплуатация уже не эксперимент. [Источник: CNCF Annual Cloud Native Survey](https://www.cncf.io/reports/the-cncf-annual-cloud-native-survey/).

## Что меняется на рынке труда

В отчёте World Economic Forum работодатели поставили сети и кибербезопасность на второе место среди быстрее всего растущих групп навыков до 2030 года — после ИИ и больших данных. [Источник: Future of Jobs Report 2025, PDF](https://reports.weforum.org/docs/WEF_Future_of_Jobs_Report_2025.pdf). Однако нельзя превращать этот вывод в обещание, что любая должность с названием DevOps обязательно вырастет.

Статистика США хорошо показывает сдвиг. BLS ожидает на 2025–2035 годы рост занятости сетевых архитекторов примерно на 8%, специалистов по информационной безопасности — на 21%, разработчиков и тестировщиков — на 10%. При этом занятость классических сетевых и системных администраторов прогнозируется ниже на 4%: часть ручных задач автоматизируется, переходит к разработчикам DevOps и сервисным провайдерам. Это данные одной страны, а не прогноз для Казахстана, но они объясняют, почему курс строится вокруг кода, автоматизации, безопасности и надёжности, а не вокруг запоминания команд. [Источники: BLS — сетевые архитекторы](https://www.bls.gov/ooh/computer-and-information-technology/computer-network-architects.htm), [информационная безопасность](https://www.bls.gov/ooh/computer-and-information-technology/information-security-analysts.htm), [системное администрирование](https://www.bls.gov/ooh/computer-and-information-technology/network-and-computer-systems-administrators.htm).

ИИ умеет написать YAML или подсказать команду, но не несёт ответственность за простой, потерю данных, счёт облачного провайдера или открытый наружу секрет. Инженер всё ещё должен формулировать проверяемые требования, понимать границы доступа, читать метрики, выбирать откат и доказывать восстановление. В курсе ИИ можно использовать как помощника, но каждое предложение проверяется командой, тестом или наблюдаемым результатом.

## Почему это важно для Казахстана

Казахстан уже переводит стратегический интерес в физическую инфраструктуру. В 2025 году запущен национальный кластер Alem.Cloud на NVIDIA H200; Министерство искусственного интеллекта и цифрового развития сообщало примерно о 2 экзафлопс в FP8. [Источник: итоги министерства за 2025 год](https://www.gov.kz/memleket/entities/maidd/press/news/details/1136177?lang=ru).

В 2026 году Правительство назвало цифровую инфраструктуру одним из приоритетов и сообщило о соглашениях по проекту «Долина ЦОДов» в Экибастузе. Первый объект на 50 МВт строится, его ввод заявлен на июнь 2027 года; весь кампус рассматривается с потенциалом до 1 ГВт. Это планы и строящиеся объекты, поэтому мы не выдаём будущую мощность за уже работающую. [Источник: Правительство РК, ход проекта](https://primeminister.kz/ru/news/olzhas-bektenov-provel-soveshchanie-o-hode-realizacii-proekta-dolina-codov-v-pavlodarskoy-oblasti-31770), [KAZAKH INVEST, проект до 1 ГВт](https://invest.gov.kz/ru/amp/news/40804/).

Такая ставка может увеличить спрос на разные специальности: энергетиков, инженеров охлаждения и связи, сетевых инженеров, специалистов физической безопасности, операторов оборудования, cloud/platform-инженеров, SRE, DevOps и специалистов кибербезопасности. Этот курс охватывает программную эксплуатацию сервисов. Он не готовит электрика или проектировщика инженерных систем дата-центра и не заменяет практику на реальном производственном объекте.

## Почему в cloud native так часто встречается Go

У современной инфраструктуры нет одного языка. Ядро Linux в основном написано на C, автоматизация часто использует Python и shell, интерфейсы — JavaScript и TypeScript, а высокопроизводительные компоненты могут быть на Rust, C++ и Java. Но именно **управляющий слой cloud native экосистемы во многом построен на Go**. На нём развиваются Docker Engine, Kubernetes, Prometheus, Caddy, OpenTofu и многие Kubernetes operators. Репозитории проектов позволяют проверить это напрямую: [Moby/Docker Engine](https://github.com/moby/moby), [Kubernetes](https://github.com/kubernetes/kubernetes), [Prometheus](https://github.com/prometheus/prometheus), [Caddy](https://github.com/caddyserver/caddy), [OpenTofu](https://github.com/opentofu/opentofu).

Go подошёл этой области по практическим причинам: из проекта удобно получить один бинарный файл, кросс-компиляция встроена в инструменты, goroutines упрощают конкурентные сетевые программы, стандартная библиотека хорошо покрывает HTTP, а явная модель ошибок удобна там, где отказ нельзя скрывать. Это не доказывает превосходство языка для любой задачи. Это объясняет, почему чтение небольшого Go-сервиса помогает инженеру инфраструктуры понимать инструменты вокруг него.

Поэтому учебный CloudLab написан на Go и собирается в статический бинарный файл для минимального контейнера. **Знать Go до начала курса не требуется**: приложение уже готово, а здесь мы управляем его жизненным циклом. Если вы захотите разобраться в коде глубже, пройдите бесплатный курс [«Go: с нуля до своего блога»](https://shanraq.org/course/go?lang=ru).

## Какую роль вы сможете примерить

- **Cloud engineer** собирает облачные сети, вычислительные ресурсы, доступы и стоимость.
- **DevOps engineer** сокращает путь от изменения к безопасному выпуску через автоматизацию и общую ответственность команды.
- **Platform engineer** создаёт внутреннюю платформу и удобный стандартный путь для разработчиков.
- **SRE** управляет надёжностью через показатели, цели обслуживания, бюджет ошибок и инженерную реакцию на инциденты.
- **Инженер дата-центра** работает также с физической площадкой, питанием, охлаждением, стойками и каналами связи; это отдельная область.

Названия вакансий пересекаются. Поэтому итог курса — не ярлык должности, а портфолио с доказательствами: контейнерный образ, CI, защищённый сервер, HTTPS, инфраструктура как код, Kubernetes-манифесты, метрики, резервная копия и протокол восстановления.

## Почему здесь нет диплома ради диплома

Shanraq не выдаёт бесполезный сертификат только за то, что вкладка урока была открыта. Сам по себе красивый PDF не поднимет упавший сервис, не найдёт утечку секрета и не восстановит данные. Наши бесплатные курсы дают работу, которую можно показать и повторить: репозиторий с понятной историей, воспроизводимый выпуск, работающий адрес, автоматические проверки, журнал учебного инцидента и восстановленная из копии информация.

Лучшее доказательство — ваше желание разобраться и усердие, превращённые в работающий результат. К концу курса у вас будет не обещание профессии, а собственный инженерный инструмент и портфолио, с которыми можно решать реальные задачи, проходить техническое собеседование и дальше строить источник дохода. Доход не возникает автоматически после 24 уроков: его дают практика, ответственность и способность надёжно решать чужую проблему. Именно это курс и помогает доказать без формальной бумаги.

## Сквозной проект и правила курса

Мы эксплуатируем небольшой сервис CloudLab. Он хранит заметки на диске, отдаёт health/readiness endpoints, метрику запросов и структурированные журналы. Приложение намеренно простое: внимание направлено на путь выпуска и отказоустойчивость. Сначала всё работает локально и бесплатно. Для уроков с публичным DNS и TLS можно ненадолго арендовать маленькую виртуальную машину у любого провайдера; конкретный бренд не требуется. После упражнения сервер можно удалить.

24 урока проходят Linux, сеть, Docker, Compose, registry, CI/CD, секреты, облачные доступы и стоимость, SSH, firewall, HTTPS, OpenTofu, Ansible, Kubernetes, наблюдаемость, резервное копирование и учебный инцидент. Читайте команду до запуска, никогда не вставляйте реальные токены в репозиторий и не экспериментируйте на чужой или рабочей инфраструктуре.

Курс бесплатный и доступен на русском, казахском и английском. Для чтения регистрация не нужна. Для практики понадобится компьютер, терминал, Git и возможность установить Docker. Linux можно запустить в виртуальной машине или WSL2. Все эталонные файлы лежат в каталоге [`course/cloud-devops-lab`](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

[Перейти к оглавлению курса](/course/cloud-devops?lang=ru)
""",
"kz": """# Бастамас бұрын: неге Cloud & DevOps дәл қазір қажет

_Лид (summary):_ **Бұл тегін тәжірибелік курстың толық кіріспесі: бұлт, ЖИ және дата-орталықтармен бірге инфрақұрылым дағдылары неге өсіп келеді, Қазақстан не салып жатыр және курс жалған мансап уәдесінсіз нені үйретеді.**

## Неге осы курсты таңдадық

Бұлттық сервис көзге көрінбегенімен, ол нақты процессорда, дискіде, желіде және электр желісінде жұмыс істейді. Компаниялар ЖИ, ағындық өңдеу, мемлекеттік цифрлық қызметтер мен онлайн өнімдерді көбірек қолданған сайын есептеу қуаты және жүйені қауіпсіз шығарып, бақылап, қалпына келтіре алатын мамандар қажет болады.

Халықаралық энергетикалық агенттік дерегі бойынша дата-орталықтардың электр тұтынуы 2025 жылы 17% өсті. Негізгі болжамда ол 2025 жылғы 485 ТВт·сағ-тан 2030 жылы 950 ТВт·сағ-қа жуықтайды, ал ЖИ-ға бағытталған орталықтардың тұтынуы үш есе артады. Бұл жұмыс орындарының саны емес, инфрақұрылымның кеңеюін көрсететін физикалық белгі. [Дереккөз: IEA](https://www.iea.org/reports/key-questions-on-energy-and-ai/executive-summary).

CNCF-тің 2025 жылғы сауалнамасында контейнер пайдаланушыларының 82%-ы Kubernetes-ті production ортасында қолданатынын, ал генеративті модель орналастыратын ұйымдардың 66%-ы inference жүктемесінің кемінде бір бөлігін Kubernetes арқылы басқаратынын айтты. Сауалнама бүкіл әлем компанияларын түгел қамтымайды, бірақ cloud native тәсілінің тәжірибеден тұрақты инфрақұрылымға өткенін көрсетеді. [Дереккөз: CNCF](https://www.cncf.io/reports/the-cncf-annual-cloud-native-survey/).

## Еңбек нарығында не өзгеріп жатыр

World Economic Forum жұмыс берушілер сауалнамасында желілер мен киберқауіпсіздік 2030 жылға дейін ең жылдам өсетін дағдылар арасында екінші орында тұр. [Дереккөз: Future of Jobs 2025, PDF](https://reports.weforum.org/docs/WEF_Future_of_Jobs_Report_2025.pdf). Бірақ бұл DevOps атауы бар әр жұмыс міндетті түрде өседі деген кепілдік емес.

АҚШ BLS болжамы өзгерістің бағытын айқын көрсетеді: 2025–2035 жылдары желі сәулетшілері шамамен 8%, ақпараттық қауіпсіздік талдаушылары 21%, әзірлеушілер мен тестілеушілер 10% өседі деп күтіледі. Классикалық желі және жүйе әкімшілері 4% азаюы мүмкін, өйткені қол жұмысының бір бөлігі автоматтандырылып, DevOps әзірлеушілеріне және сервис провайдерлеріне өтеді. Бұл Қазақстанға арналған болжам емес, бірақ курста пәрмен жаттаудың орнына код, автоматтандыру, қауіпсіздік және сенімділік неге негізгі екенін түсіндіреді. [BLS: желі сәулетшілері](https://www.bls.gov/ooh/computer-and-information-technology/computer-network-architects.htm), [ақпараттық қауіпсіздік](https://www.bls.gov/ooh/computer-and-information-technology/information-security-analysts.htm), [жүйе әкімшілері](https://www.bls.gov/ooh/computer-and-information-technology/network-and-computer-systems-administrators.htm).

ЖИ YAML жаза алады немесе пәрмен ұсынады, бірақ тоқтап қалу, дерек жоғалуы, бұлт шоты немесе ашық қалған құпия үшін жауап бермейді. Инженер тексерілетін талап құрып, рұқсат шегін түсініп, метриканы оқып, кері қайтаруды таңдап, қалпына келтіруді дәлелдеуі керек. Курста ЖИ көмекші бола алады, бірақ әр ұсыныс пәрменмен, тестпен немесе байқалатын нәтижемен тексеріледі.

## Бұл Қазақстан үшін неге маңызды

Қазақстан стратегиялық ниетті нақты инфрақұрылымға айналдырып жатыр. 2025 жылы NVIDIA H200 негізіндегі Alem.Cloud ұлттық кластері іске қосылды; министрлік FP8 өлшемінде шамамен 2 экзафлопс қуат туралы хабарлады. [Дереккөз: министрліктің 2025 жылғы қорытындысы](https://www.gov.kz/memleket/entities/maidd/press/news/details/1136177?lang=kk).

2026 жылы Үкімет Екібастұздағы «Деректерді өңдеу орталықтарының алқабы» жобасын басым бағыт деп атады. Алғашқы 50 МВт нысан салынып жатыр, оны 2027 жылғы маусымда іске қосу жоспарланған; кампустың болашақ әлеуеті 1 ГВт-қа дейін қарастырылуда. Бұлар жоспар мен құрылыс барысындағы нысандар, сондықтан болашақ қуатты қазір жұмыс істеп тұрған қуат деп көрсетпейміз. [Үкімет: жоба барысы](https://primeminister.kz/ru/news/olzhas-bektenov-provel-soveshchanie-o-hode-realizacii-proekta-dolina-codov-v-pavlodarskoy-oblasti-31770), [KAZAKH INVEST: 1 ГВт-қа дейін](https://invest.gov.kz/ru/amp/news/40804/).

Бұл бағыт энергетиктерге, салқындату және байланыс инженерлеріне, желі мамандарына, жабдық операторларына, cloud/platform инженерлеріне, SRE, DevOps және киберқауіпсіздік мамандарына сұранысты арттыруы мүмкін. Бұл курс сервистің бағдарламалық эксплуатациясын қамтиды. Ол дата-орталық электр жүйесін жобалауды үйретпейді және нақты өндірістік тәжірибені алмастырмайды.

## Cloud native әлемінде Go неге жиі кездеседі

Қазіргі инфрақұрылымның жалғыз тілі жоқ. Linux ядросы негізінен C тілінде, автоматтандыруда Python мен shell, интерфейсте JavaScript және TypeScript, ал кейбір жоғары өнімді компонентте Rust, C++ және Java қолданылады. Бірақ **cloud native экожүйесінің басқару қабатының елеулі бөлігі Go тілінде жасалған**. Docker Engine, Kubernetes, Prometheus, Caddy, OpenTofu және көптеген Kubernetes operator осы тілде дамиды. Мұны жобалардың өз репозиторийінен тексеруге болады: [Moby/Docker Engine](https://github.com/moby/moby), [Kubernetes](https://github.com/kubernetes/kubernetes), [Prometheus](https://github.com/prometheus/prometheus), [Caddy](https://github.com/caddyserver/caddy), [OpenTofu](https://github.com/opentofu/opentofu).

Go бұл салаға тәжірибелік себептермен сай келді: жобадан бір binary файл алу ыңғайлы, cross-compilation құралға кіріктірілген, goroutine бәсекелі желілік бағдарламаны жеңілдетеді, standard library HTTP-ды жақсы қамтиды, ал айқын error моделі ақауды жасыруға болмайтын ортаға сәйкес. Бұл тіл барлық міндетте үздік деген сөз емес. Бұл шағын Go сервисін оқу инфрақұрылым инженеріне айналасындағы құралды түсінуге неге көмектесетінін көрсетеді.

Сондықтан CloudLab Go тілінде жазылып, минимал контейнер үшін статикалық binary болып жиналады. **Курсты бастау үшін Go білу міндетті емес**: қолданба дайын, біз оның өмір циклін басқарамыз. Кодты тереңірек түсінгіңіз келсе, [«Go: нөлден өз блогыңа дейін»](https://shanraq.org/course/go?lang=kz) тегін курсын өтіңіз.

## Қай рөлді байқап көресіз

- **Cloud engineer** бұлттық желілерді, есептеу ресурстарын, рұқсаттарды және шығынды басқарады.
- **DevOps engineer** автоматтандыру арқылы өзгерістен қауіпсіз шығарылымға дейінгі жолды қысқартады.
- **Platform engineer** әзірлеушілерге ішкі платформа мен стандартты ыңғайлы жол жасайды.
- **SRE** көрсеткіш, қызмет мақсаты, қате бюджеті және инцидентке инженерлік әрекет арқылы сенімділікті басқарады.
- **Дата-орталық инженері** физикалық алаң, қуат, салқындату, тірек және байланыс арнасымен де жұмыс істейді; бұл бөлек мамандық саласы.

Лауазым атаулары қиылысады. Сондықтан курс нәтижесі — атау емес, дәлелі бар портфолио: контейнер бейнесі, CI, қорғалған сервер, HTTPS, infrastructure as code, Kubernetes манифесттері, метрика, сақтық көшірме және қалпына келтіру хаттамасы.

## Неге мұнда диплом үшін диплом жоқ

Shanraq сабақ беті ашылғаны үшін ғана пайдасыз сертификат бермейді. Әдемі PDF құлаған сервисті көтермейді, құпияның ағуын таппайды және деректі қалпына келтірмейді. Біздің тегін курстарымыз көрсетуге және қайталауға болатын нәтиже береді: түсінікті тарихы бар репозиторий, қайталанатын шығарылым, жұмыс істейтін мекенжай, автоматты тексеру, оқу инцидентінің журналы және көшірмеден қалпына келген дерек.

Ең жақсы дәлел — түсінуге деген ынта мен еңбектің жұмыс істейтін нәтижеге айналуы. Курс соңында мамандық туралы бос уәде емес, нақты міндет шешуге, техникалық сұхбатта көрсетуге және табыс көзін әрі қарай құруға болатын инженерлік құрал мен портфолио қалады. 24 сабақ табысқа автоматты кепілдік бермейді: табысты тәжірибе, жауапкершілік және өзгенің мәселесін сенімді шешу қабілеті әкеледі. Курс осыны формалды қағазсыз дәлелдеуге көмектеседі.

## Ортақ жоба және оқу ережесі

Біз CloudLab атты шағын сервисті басқарамыз. Ол жазбаларды дискіге сақтайды, health/readiness endpoint, сұрау метрикасы және құрылымдалған журнал береді. Қолданба әдейі қарапайым: назар шығару жолы мен ақауға төзімділікке түседі. Алғашқы жұмыстар жергілікті және тегін. Ашық DNS пен TLS сабағына кез келген провайдерден шағын виртуалды серверді қысқа уақытқа алуға болады; нақты бренд міндетті емес. Тәжірибеден кейін серверді жойыңыз.

24 сабақ Linux, желі, Docker, Compose, registry, CI/CD, құпиялар, бұлттық рұқсат пен шығын, SSH, firewall, HTTPS, OpenTofu, Ansible, Kubernetes, бақылау, сақтық көшірме және оқу инцидентін қамтиды. Пәрменді іске қоспас бұрын оқыңыз, шын токенді репозиторийге салмаңыз және бөтен не жұмыс инфрақұрылымында тәжірибе жасамаңыз.

Курс қазақ, орыс және ағылшын тілінде тегін беріледі. Оқу үшін тіркелу керек емес. Тәжірибе үшін компьютер, терминал, Git және Docker орнату мүмкіндігі қажет. Linux-ті виртуалды машинада немесе WSL2-де іске қосуға болады. Эталон файлдар [`course/cloud-devops-lab`](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab) каталогында.

[Курс мазмұнына өту](/course/cloud-devops?lang=kz)
""",
"en": """# Before you begin: why Cloud & DevOps now

_Lead (summary):_ **A detailed introduction to this free practical course: why infrastructure skills are growing with cloud, AI, and data centres, what Kazakhstan is building, and what you will learn without a promise of an instant career.**

## Why we chose this course

A cloud service looks intangible, but it always runs on physical processors, disks, networks, cooling, and electricity. As organizations use more AI, streaming data, digital public services, and online products, they need more compute and people who can release, observe, secure, and recover software systems.

The International Energy Agency reports that global data centre electricity use rose 17% in 2025. Its central outlook almost doubles consumption from 485 TWh in 2025 to 950 TWh in 2030, while consumption by AI-focused facilities triples. This is not a forecast of job counts, but it is a physical signal that infrastructure is expanding. [Source: IEA, Key Questions on Energy and AI](https://www.iea.org/reports/key-questions-on-energy-and-ai/executive-summary).

The software layer is mature too. In CNCF's 2025 survey, 82% of container users reported running Kubernetes in production, and 66% of organizations hosting generative models used Kubernetes for at least some inference workloads. The respondents come from the cloud native ecosystem and do not represent every company, but the result shows that container operations are no longer merely experimental. [Source: CNCF Annual Cloud Native Survey](https://www.cncf.io/reports/the-cncf-annual-cloud-native-survey/).

## What is changing in work

Employers surveyed by the World Economic Forum ranked networks and cybersecurity as the second fastest-growing skill group through 2030, after AI and big data. [Source: Future of Jobs Report 2025, PDF](https://reports.weforum.org/docs/WEF_Future_of_Jobs_Report_2025.pdf). That finding must not be turned into a promise that every job titled DevOps will grow.

US labour projections illustrate the shift. For 2025–2035, BLS projects about 8% growth for network architects, 21% for information security analysts, and 10% for software developers, QA analysts, and testers. It projects a 4% decline for traditional network and systems administrators because some routine work is automated, moved to DevOps-focused developers, or outsourced to service providers. These are US figures, not a Kazakhstan forecast, but they explain why this course focuses on code, automation, security, and reliability instead of memorising commands. [Sources: BLS on network architects](https://www.bls.gov/ooh/computer-and-information-technology/computer-network-architects.htm), [information security](https://www.bls.gov/ooh/computer-and-information-technology/information-security-analysts.htm), and [systems administration](https://www.bls.gov/ooh/computer-and-information-technology/network-and-computer-systems-administrators.htm).

AI can draft YAML or suggest a command. It does not take responsibility for downtime, lost data, a cloud bill, or an exposed secret. An engineer still has to state testable requirements, understand access boundaries, read metrics, choose a rollback, and prove recovery. You may use AI as an assistant in this course, but verify every proposal with a command, test, or observable result.

## Why this matters in Kazakhstan

Kazakhstan is turning strategic interest into physical infrastructure. In 2025 it launched the national Alem.Cloud cluster on NVIDIA H200 accelerators; the Ministry of Artificial Intelligence and Digital Development reported performance of about two exaflops at FP8. [Source: the Ministry's 2025 results](https://www.gov.kz/memleket/entities/maidd/press/news/details/1136177?lang=en).

In 2026 the Government described digital infrastructure as a priority and announced agreements around the Data Center Valley in Ekibastuz. Its first 50 MW facility is under construction with commissioning announced for June 2027; the wider campus is being considered with capacity up to 1 GW. These are plans and construction milestones, so this course does not describe future capacity as already operating. [Source: Government project update](https://primeminister.kz/ru/news/olzhas-bektenov-provel-soveshchanie-o-hode-realizacii-proekta-dolina-codov-v-pavlodarskoy-oblasti-31770), [KAZAKH INVEST on the proposed 1 GW scale](https://invest.gov.kz/ru/amp/news/40804/).

This direction may increase demand across several distinct fields: power and cooling, telecommunications, network engineering, physical security, hardware operations, cloud and platform engineering, SRE, DevOps, and cybersecurity. This course covers software operations. It does not train a data centre electrician or facilities designer, and it cannot replace experience on a real production site.

## Why Go appears so often in cloud native infrastructure

Modern infrastructure has no single language. The Linux kernel is mostly C, automation often uses Python and shell, interfaces use JavaScript and TypeScript, and performance-sensitive components may use Rust, C++, or Java. Yet **a substantial part of the cloud native control plane is built in Go**. Docker Engine, Kubernetes, Prometheus, Caddy, OpenTofu, and many Kubernetes operators are developed in it. Their own repositories make that claim inspectable: [Moby/Docker Engine](https://github.com/moby/moby), [Kubernetes](https://github.com/kubernetes/kubernetes), [Prometheus](https://github.com/prometheus/prometheus), [Caddy](https://github.com/caddyserver/caddy), and [OpenTofu](https://github.com/opentofu/opentofu).

Go fits this field for practical reasons: a project can produce one convenient binary, cross-compilation is built into the toolchain, goroutines suit concurrent network services, the standard library has strong HTTP support, and explicit errors suit systems where failure must remain visible. This does not make Go the best language for every task. It explains why reading a small Go service helps an infrastructure engineer understand the tools around it.

CloudLab is therefore written in Go and built as a static binary for a minimal container. **You do not need to know Go before starting**: the application is supplied and this course operates its lifecycle. To understand its code more deeply, take the free [Go: From Zero to Your Own Blog](https://shanraq.org/course/go?lang=en) course.

## The roles you can try

- A **cloud engineer** manages cloud networks, compute, access, and cost.
- A **DevOps engineer** shortens the path from a change to a safe release through automation and shared responsibility.
- A **platform engineer** builds an internal platform and a paved path for developers.
- An **SRE** manages reliability through indicators, service objectives, error budgets, and engineering incident response.
- A **data centre engineer** also works with the physical site, power, cooling, racks, and links; this is a separate field.

Job titles overlap. The course outcome is therefore evidence rather than a label: a container image, CI, a secured server, HTTPS, infrastructure as code, Kubernetes manifests, metrics, a backup, and a recovery record.

## Why there is no diploma for its own sake

Shanraq does not issue a useless certificate merely because a lesson tab was opened. A polished PDF cannot restore a failed service, find a leaked secret, or recover data. Our free courses produce work that can be shown and repeated: a repository with a clear history, a reproducible release, a working endpoint, automated checks, an incident record, and data restored from a backup.

The strongest evidence is your desire to understand and your sustained effort turned into a working result. At the end you have an engineering tool and portfolio that can solve real problems, support a technical interview, and help you build a livelihood. Twenty-four lessons do not guarantee income automatically; practice, responsibility, and the ability to solve another person's problem reliably create value. The course helps you prove those abilities without relying on a decorative credential.

## The project and course rules

We operate a small service called CloudLab. It saves notes to disk and exposes health and readiness endpoints, a request metric, and structured logs. The application is deliberately simple so the deployment path and failure behaviour stay visible. The first part is local and free. For public DNS and TLS, you may rent a small virtual machine from any provider for a short time; no brand is required. Delete it after the exercise.

The 24 lessons cover Linux, networking, Docker, Compose, a registry, CI/CD, secrets, cloud access and cost, SSH, firewall rules, HTTPS, OpenTofu, Ansible, Kubernetes, observability, backups, and a recovery drill. Read commands before running them, never commit real tokens, and do not experiment on infrastructure you do not own or on a production system.

The course is free in Kazakh, Russian, and English. Reading needs no account. Practice needs a computer, a terminal, Git, and permission to install Docker. You can run Linux in a virtual machine or WSL2. The checked reference files live in [`course/cloud-devops-lab`](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

[Open the course contents](/course/cloud-devops?lang=en)
"""
}


TOPICS = [
{"stem":"01-why-now","title":{"ru":"Почему Cloud & DevOps сейчас: инфраструктура, роли и границы","kz":"Cloud & DevOps неге қазір қажет: инфрақұрылым, рөлдер және шекаралар","en":"Why Cloud & DevOps now: infrastructure, roles, and boundaries"},
 "summary":{"ru":"Отделяем рост дата-центров от роста конкретных профессий и строим карту навыков cloud engineer, DevOps, platform engineer и SRE.","kz":"Дата-орталық өсімін нақты мамандық өсімінен ажыратып, cloud engineer, DevOps, platform engineer және SRE дағдыларының картасын құрамыз.","en":"Separate data-centre growth from job-title growth and map the skills of cloud engineering, DevOps, platform engineering, and SRE."},
 "concept":{"ru":"DevOps — не один инструмент и не человек, который один отвечает за всё. Это способ уменьшать риск изменений с помощью автоматизации, короткой обратной связи и совместной ответственности. Дата-центр даёт физическую мощность, облако превращает её в программно заказываемый ресурс, а эксплуатационная практика удерживает сервис доступным.","kz":"DevOps — бір құрал да, бәріне жалғыз жауап беретін адам да емес. Ол автоматтандыру, қысқа кері байланыс және ортақ жауапкершілік арқылы өзгеріс тәуекелін азайтады. Дата-орталық физикалық қуат береді, бұлт оны бағдарламамен тапсырыс берілетін ресурсқа айналдырады, ал эксплуатация сервисті қолжетімді ұстайды.","en":"DevOps is neither one tool nor one person responsible for everything. It reduces change risk through automation, short feedback loops, and shared responsibility. A data centre supplies physical capacity, cloud turns it into programmable resources, and operations keeps a service available."},
 "cmd":"printf '%s\n' 'change -> test -> image -> deploy -> observe -> recover'\nprintf '%s\n' 'evidence: repeatable command, metric, backup, recovery record'",
 "checks":{"ru":["Цепочка начинается с изменения и заканчивается проверенным восстановлением.","Каждый этап оставляет наблюдаемое доказательство.","Вы можете назвать границу этого курса: программная эксплуатация."],"kz":["Тізбек өзгерістен басталып, тексерілген қалпына келтірумен аяқталады.","Әр кезең байқалатын дәлел қалдырады.","Курс шекарасын атай аласыз: бағдарламалық эксплуатация."],"en":["The chain begins with a change and ends with verified recovery.","Every stage leaves observable evidence.","You can state the course boundary: software operations."]},
 "task":{"ru":"Напишите POSIX sh-скрипт, который печатает шесть этапов пути выпуска по одному на строку и завершает работу при первой ошибке.","kz":"Шығарылым жолының алты кезеңін жеке жолға басатын және алғашқы қатеде тоқтайтын POSIX sh скриптін жазыңыз.","en":"Write a POSIX sh script that prints the six release-path stages on separate lines and stops on the first error."}},
{"stem":"02-service-map","title":{"ru":"CloudLab: сервис, среда и критерии готовности","kz":"CloudLab: сервис, орта және дайындық өлшемдері","en":"CloudLab: service, environment, and readiness criteria"},
 "summary":{"ru":"Запускаем сквозной сервис, исследуем его endpoints и превращаем расплывчатое «работает» в проверяемый контракт.","kz":"Ортақ сервисті іске қосып, endpoint-терін зерттейміз және көмескі «жұмыс істейді» сөзін тексерілетін келісімге айналдырамыз.","en":"Run the course service, inspect its endpoints, and turn a vague “works” into a testable contract."},
 "concept":{"ru":"CloudLab хранит заметки и специально открывает разные сигналы: `/healthz` говорит, что процесс жив, `/readyz` — что он готов принимать трафик, `/metrics` — сколько запросов увидел. Один HTTP 200 на главной странице ещё не доказывает сохранность данных, готовность зависимости или корректное завершение.","kz":"CloudLab жазбаларды сақтайды және әртүрлі сигнал ашады: `/healthz` үдерістің тірі екенін, `/readyz` трафик қабылдауға дайын екенін, `/metrics` қанша сұрау көргенін айтады. Басты беттегі бір HTTP 200 дерек сақталғанын не тәуелділік дайын екенін дәлелдемейді.","en":"CloudLab stores notes and exposes distinct signals: `/healthz` says the process is alive, `/readyz` says it can accept traffic, and `/metrics` counts requests. One HTTP 200 from the home page does not prove data durability, dependency readiness, or graceful shutdown."},
 "cmd":"go test ./...\nCLOUDLAB_ADDR=:8080 go run ./cmd/cloudlab &\npid=$!\ntrap 'kill \"$pid\" 2>/dev/null || true' EXIT\nsleep 1\ncurl -fsS http://127.0.0.1:8080/healthz\ncurl -fsS http://127.0.0.1:8080/metrics",
 "checks":{"ru":["Тесты проходят.","Health endpoint отвечает без HTML.","Счётчик запросов увеличивается после повторного curl."],"kz":["Тесттер өтеді.","Health endpoint HTML-сіз жауап береді.","Қайта curl жасағанда сұрау санауышы өседі."],"en":["The tests pass.","The health endpoint responds without HTML.","The request counter rises after another curl."]},
 "task":{"ru":"Напишите POSIX sh-проверку, которая с `curl -fsS` проверяет `/healthz` и `/readyz`, а затем печатает `cloudlab ready`.","kz":"`curl -fsS` арқылы `/healthz` және `/readyz` тексеріп, соңында `cloudlab ready` басатын POSIX sh тексеруін жазыңыз.","en":"Write a POSIX sh check that uses `curl -fsS` for `/healthz` and `/readyz`, then prints `cloudlab ready`."}},
{"stem":"03-terminal-git","title":{"ru":"Терминал и Git: воспроизводимая рабочая точка","kz":"Терминал және Git: қайталанатын жұмыс нүктесі","en":"Terminal and Git: a reproducible working point"},
 "summary":{"ru":"Учимся ориентироваться в каталоге, читать изменение до коммита и возвращаться к известной версии без удаления истории.","kz":"Каталогта бағдарлауды, commit алдында өзгерісті оқуды және тарихты өшірмей белгілі нұсқаға оралуды үйренеміз.","en":"Navigate the workspace, review a change before committing, and return to a known version without deleting history."},
 "concept":{"ru":"Инфраструктурная команда опасна, когда неизвестны каталог, версия и различие файлов. Git даёт не магическую отмену, а журнал намерений. Перед автоматизацией фиксируйте точку, смотрите diff и связывайте выпуск с конкретным commit SHA.","kz":"Каталог, нұсқа және файл айырмасы белгісіз болса, инфрақұрылым пәрмені қауіпті. Git — сиқырлы кері қайтару емес, ниет журналы. Автоматтандыру алдында нүктені бекітіп, diff қарап, шығарылымды нақты commit SHA-мен байланыстырыңыз.","en":"An infrastructure command is dangerous when its directory, version, and file differences are unknown. Git is a journal of intent rather than magical undo. Record a point, inspect the diff, and tie a release to an exact commit SHA."},
 "cmd":"pwd\ngit status --short\ngit diff --check\ngit rev-parse --short HEAD\ngit log -1 --format='%h %cs %s'",
 "checks":{"ru":["`pwd` указывает на репозиторий.","`git diff --check` не находит ошибок пробелов.","SHA однозначно связывает код и выпуск."],"kz":["`pwd` репозиторийді көрсетеді.","`git diff --check` бос орын қатесін таппайды.","SHA код пен шығарылымды бірмәнді байланыстырады."],"en":["`pwd` identifies the repository.","`git diff --check` reports no whitespace errors.","The SHA ties code to the release unambiguously."]},
 "task":{"ru":"Напишите POSIX sh-скрипт с `set -eu`, который отказывается работать вне Git-репозитория и печатает короткий SHA текущего commit.","kz":"`set -eu` қолданып, Git репозиторийінен тыс жұмыс істеуден бас тартатын және ағымдағы commit-тің қысқа SHA-сын басатын POSIX sh скриптін жазыңыз.","en":"Write a POSIX sh script with `set -eu` that refuses to run outside a Git repository and prints the current short commit SHA."}},
{"stem":"04-linux-access","title":{"ru":"Linux: пользователи, файлы и минимальные права","kz":"Linux: пайдаланушылар, файлдар және ең аз құқық","en":"Linux: users, files, and least privilege"},
 "summary":{"ru":"Разбираем владельца, группу, режим доступа и почему сервис не должен постоянно работать от root.","kz":"Ие, топ, қолжетімділік режимін және сервис неліктен үнемі root атынан жұмыс істемеуі керегін талдаймыз.","en":"Understand owners, groups, permission modes, and why a service should not run permanently as root."},
 "concept":{"ru":"Минимальные права ограничивают радиус последствий: скомпрометированный процесс получает только то, что нужно его задаче. Числа `750` и `640` — не ритуал; они описывают чтение, запись и выполнение для владельца, группы и остальных.","kz":"Ең аз құқық зардап аумағын шектейді: бұзылған үдеріс тек өз міндетіне қажетті рұқсатты алады. `750` және `640` сандары рәсім емес; олар ие, топ және басқалар үшін оқу, жазу, орындауды сипаттайды.","en":"Least privilege limits blast radius: a compromised process receives only what its task needs. Modes `750` and `640` are not rituals; they encode read, write, and execute access for owner, group, and others."},
 "cmd":"id\numask\nmkdir -p /tmp/cloudlab-permissions\nprintf '%s\n' secret > /tmp/cloudlab-permissions/config\nchmod 640 /tmp/cloudlab-permissions/config\nls -ld /tmp/cloudlab-permissions\nls -l /tmp/cloudlab-permissions/config",
 "checks":{"ru":["Файл не исполняемый.","У остальных пользователей нет доступа.","Вы можете расшифровать каждую цифру режима 640."],"kz":["Файл орындалмайды.","Басқа пайдаланушыларға рұқсат жоқ.","640 режимінің әр санын түсіндіре аласыз."],"en":["The file is not executable.","Other users have no access.","You can decode every digit in mode 640."]},
 "task":{"ru":"Напишите POSIX sh-скрипт, который создаёт каталог с режимом 750 и файл конфигурации с режимом 640, затем проверяет режимы через `stat`.","kz":"750 режимді каталог пен 640 режимді конфигурация файлын жасап, кейін режимдерді `stat` арқылы тексеретін POSIX sh скриптін жазыңыз.","en":"Write a POSIX sh script that creates a directory with mode 750 and a config file with mode 640, then verifies both modes with `stat`."}},
{"stem":"05-processes-logs","title":{"ru":"Процессы, сигналы, службы и журналы","kz":"Үдерістер, сигналдар, қызметтер және журналдар","en":"Processes, signals, services, and logs"},
 "summary":{"ru":"Наблюдаем PID, корректную остановку и структурированный журнал, а затем читаем unit-файл systemd как контракт запуска.","kz":"PID, дұрыс тоқтау және құрылымдалған журналды бақылаймыз, кейін systemd unit файлын іске қосу келісімі ретінде оқимыз.","en":"Observe a PID, graceful shutdown, and structured logs, then read a systemd unit as a startup contract."},
 "concept":{"ru":"Процесс имеет жизненный цикл. SIGTERM просит завершиться и даёт сохранить состояние; SIGKILL немедленно отбирает такую возможность. Менеджер служб добавляет пользователя, окружение, перезапуск и журнал, но не исправляет приложение, которое игнорирует завершение.","kz":"Үдерістің өмір циклі бар. SIGTERM аяқталуды сұрап, күйді сақтауға мүмкіндік береді; SIGKILL оны бірден тоқтатады. Қызмет менеджері пайдаланушы, орта, қайта іске қосу және журнал қосады, бірақ аяқталу сигналын елемейтін қолданбаны түзетпейді.","en":"A process has a lifecycle. SIGTERM requests shutdown and allows state to be saved; SIGKILL removes that chance. A service manager adds identity, environment, restart policy, and logs, but cannot fix an application that ignores shutdown."},
 "cmd":"CLOUDLAB_ADDR=:8080 go run ./cmd/cloudlab > /tmp/cloudlab.log 2>&1 &\npid=$!\nsleep 1\nps -p \"$pid\" -o pid=,ppid=,command=\nkill -TERM \"$pid\"\nwait \"$pid\"\ntail -n 5 /tmp/cloudlab.log",
 "checks":{"ru":["PID существует до SIGTERM.","`wait` завершается без принудительного убийства.","Журнал содержит поля времени, уровня и сообщения."],"kz":["SIGTERM-ге дейін PID бар.","`wait` күштеп өлтірусіз аяқталады.","Журналда уақыт, деңгей және хабар өрістері бар."],"en":["The PID exists before SIGTERM.","`wait` completes without a forced kill.","The log has time, level, and message fields."]},
 "task":{"ru":"Напишите POSIX sh-скрипт, который запускает фоновый процесс, сохраняет PID, устанавливает `trap` для SIGTERM и дожидается чистого завершения.","kz":"Фондық үдерісті іске қосып, PID сақтайтын, SIGTERM үшін `trap` орнататын және таза аяқталуды күтетін POSIX sh скриптін жазыңыз.","en":"Write a POSIX sh script that starts a background process, records its PID, installs a SIGTERM trap, and waits for clean shutdown."}},
{"stem":"06-networks-dns-http","title":{"ru":"Сеть без магии: IP, порт, DNS, HTTP и TLS","kz":"Сиқырсыз желі: IP, порт, DNS, HTTP және TLS","en":"Networking without magic: IP, port, DNS, HTTP, and TLS"},
 "summary":{"ru":"Прослеживаем запрос от имени домена до сокета процесса и учимся локализовать сетевой отказ по слоям.","kz":"Сұрауды домен атынан үдеріс сокетіне дейін қадағалап, желі ақауын қабат бойынша табуды үйренеміз.","en":"Trace a request from a domain name to a process socket and localize failures layer by layer."},
 "concept":{"ru":"DNS отвечает «какой адрес», TCP — «удаётся ли соединиться с портом», HTTP — «понял ли приложение запрос», TLS — «с кем установлено защищённое соединение». Проверяйте слои по порядку: иначе ошибка сертификата может маскироваться как «сервер не работает».","kz":"DNS «қай мекенжай», TCP «портпен байланыс бар ма», HTTP «қолданба сұрауды түсінді ме», TLS «қорғалған байланыс кіммен орнады» дегенге жауап береді. Қабаттарды ретімен тексеріңіз, әйтпесе сертификат қатесі «сервер істемейді» болып көрінеді.","en":"DNS answers which address, TCP whether a port accepts a connection, HTTP whether the application understood a request, and TLS whom the secure connection reached. Test layers in order so a certificate failure does not become a vague “server is down.”"},
 "cmd":"getent hosts example.com 2>/dev/null || nslookup example.com\ncurl -sS -o /dev/null -w 'status=%{http_code} ip=%{remote_ip} tls=%{ssl_verify_result}\\n' https://example.com\nprintf 'GET /healthz HTTP/1.1\\r\\nHost: localhost\\r\\nConnection: close\\r\\n\\r\\n' | nc 127.0.0.1 8080 || true",
 "checks":{"ru":["DNS возвращает адрес отдельно от HTTP.","curl показывает статус, удалённый IP и результат TLS.","Вы знаете, почему порт 443 не является самим HTTPS."],"kz":["DNS HTTP-ден бөлек мекенжай қайтарады.","curl мәртебе, қашық IP және TLS нәтижесін көрсетеді.","443 порттың өзі HTTPS емес екенін түсіндіре аласыз."],"en":["DNS returns an address independently of HTTP.","curl shows status, remote IP, and TLS verification.","You can explain why port 443 is not itself HTTPS."]},
 "task":{"ru":"Напишите POSIX sh-скрипт диагностики: принять URL аргументом, вывести HTTP-код и удалённый IP через `curl`, завершиться ошибкой при коде 400 и выше.","kz":"URL аргументін қабылдап, `curl` арқылы HTTP код пен қашық IP шығаратын, 400 не жоғары кодта қатемен аяқталатын POSIX sh диагностикасын жазыңыз.","en":"Write a POSIX sh diagnostic that accepts a URL, prints its HTTP status and remote IP with `curl`, and fails for status 400 or higher."}},
]

# Lessons 7–24 use the same carefully bounded lesson shape. Their facts and
# commands remain explicit here so generated pages stay reviewable.
MORE = [
("07-containers",("Контейнер: процесс с границами, а не маленькая VM","Контейнер: шекарасы бар үдеріс, шағын VM емес","A container: a bounded process, not a tiny VM"),("Запускаем первый контейнер, видим namespaces и неизменяемый образ, а также отделяем процесс от виртуальной машины.","Алғашқы контейнерді іске қосып, namespace пен өзгермейтін бейнені көреміз және үдерісті виртуалды машинадан ажыратамыз.","Run a first container, observe namespaces and an immutable image, and distinguish a process from a virtual machine."),("Контейнер разделяет ядро хоста и изолирует процессы, сеть и файловое представление. Образ — шаблон, контейнер — конкретный запущенный экземпляр. Удаление контейнера не должно быть равнозначно потере важных данных.","Контейнер хост ядросын бөлісіп, үдеріс, желі және файл көрінісін оқшаулайды. Бейне — үлгі, контейнер — оның іске қосылған данасы. Контейнерді жою маңызды деректі жоғалтумен тең болмауы керек.","A container shares the host kernel while isolating process, network, and filesystem views. An image is a template; a container is one running instance. Removing a container must not equal losing important data."),"docker version\ndocker run --rm alpine:3.23 cat /etc/os-release\ndocker ps -a\ndocker image ls alpine:3.23",("`--rm` оқу контейнерін аяқталған соң жояды.","`--rm` оқу контейнерін аяқталған соң жояды.","`--rm` removes the learning container after exit."),("Напишите команды, которые запускают контейнер `alpine:3.23`, печатают его hostname и автоматически удаляют контейнер после выхода.","`alpine:3.23` контейнерін іске қосып, hostname басатын және шыққан соң контейнерді автоматты жоятын пәрмендерді жазыңыз.","Write commands that run `alpine:3.23`, print its hostname, and remove the container automatically after exit.")),
("08-images-lifecycle",("Слои образа, теги и жизненный цикл контейнера","Бейне қабаттары, тегтер және контейнер өмір циклі","Image layers, tags, and the container lifecycle"),("Исследуем историю образа, фиксируем digest и различаем stop, start, remove и pull.","Бейне тарихын зерттеп, digest бекітеміз және stop, start, remove, pull айырмасын көреміз.","Inspect image history, pin a digest, and distinguish stop, start, remove, and pull."),("Тег — изменяемое имя, digest — идентификатор конкретного содержимого. Для воспроизводимого выпуска недостаточно написать `latest`: завтра это имя может указывать на другие байты. Жизненный цикл контейнера и образа тоже разный.","Тег — өзгеретін ат, digest — нақты мазмұн идентификаторы. Қайталанатын шығарылым үшін `latest` жеткіліксіз: ертең ол басқа байтты көрсетуі мүмкін. Контейнер мен бейненің өмір циклі де бөлек.","A tag is a movable name; a digest identifies exact content. `latest` is insufficient for a reproducible release because it may point to different bytes tomorrow. Containers and images also have separate lifecycles."),"docker pull alpine:3.23\ndocker image inspect alpine:3.23 --format '{{index .RepoDigests 0}}'\ndocker history alpine:3.23\ndocker create --name cloudlab-inspect alpine:3.23 true\ndocker start -a cloudlab-inspect\ndocker rm cloudlab-inspect",("Digest содержит sha256.","Digest ішінде sha256 бар.","The digest contains sha256."),("Напишите команды, которые получают образ, сохраняют его первый RepoDigest в переменную и отказываются продолжать, если digest пуст.","Бейнені алып, алғашқы RepoDigest-ті айнымалыға сақтайтын және digest бос болса тоқтайтын пәрмендерді жазыңыз.","Write commands that pull an image, save its first RepoDigest in a variable, and stop if the digest is empty.")),
("09-dockerfile",("Dockerfile: небольшая, проверяемая и непривилегированная сборка","Dockerfile: шағын, тексерілетін және артық құқығы жоқ жинақ","Dockerfile: small, tested, and unprivileged"),("Собираем CloudLab в multi-stage образ, запускаем тесты при сборке и проверяем непривилегированного пользователя.","CloudLab-ты multi-stage бейнеге жинап, жинақ кезінде тест өткіземіз және артық құқығы жоқ пайдаланушыны тексереміз.","Build CloudLab as a multi-stage image, test during the build, and verify its unprivileged user."),("Multi-stage сборка оставляет компилятор в строительном слое, а в итог переносит только бинарные файлы. Это уменьшает поверхность атаки и размер. `USER 65532` ограничивает процесс, а статический healthcheck не требует shell внутри образа.","Multi-stage жинақ компиляторды құрастыру қабатында қалдырып, соңғы бейнеге тек binary файлдарды өткізеді. Бұл шабуыл бетін және өлшемді азайтады. `USER 65532` үдерісті шектейді, ал статикалық healthcheck бейне ішінде shell талап етпейді.","A multi-stage build leaves the compiler in the build stage and copies only binaries into the final image. This reduces size and attack surface. `USER 65532` constrains the process, while a static health check needs no shell in the image."),"docker build --build-arg VERSION=lesson-09 -t cloudlab:lesson-09 .\ndocker image inspect cloudlab:lesson-09 --format 'user={{.Config.User}} size={{.Size}}'\ndocker run --rm --entrypoint /cloudlab-healthcheck cloudlab:lesson-09 || true",("Образ собирается только после тестов.","Бейне тек тест өткен соң жиналады.","The image builds only after tests pass."),("Напишите команды сборки `cloudlab:practice` и проверку через `docker image inspect`, которая завершается ошибкой, если пользователь образа равен `0` или пуст.","`cloudlab:practice` жинап, бейне пайдаланушысы `0` не бос болса қатемен аяқталатын `docker image inspect` тексеруін жазыңыз.","Write commands to build `cloudlab:practice` and fail via `docker image inspect` when the configured user is `0` or empty.")),
("10-persistent-data",("Тома, bind mounts и доказательство сохранности данных","Томдар, bind mount және деректің сақталуын дәлелдеу","Volumes, bind mounts, and proof of persistence"),("Выносим состояние из контейнера, пересоздаём его и доказываем, что заметка пережила замену процесса.","Күйді контейнерден шығарып, оны қайта жасап, жазбаның үдеріс ауысуынан кейін қалғанын дәлелдейміз.","Move state out of the container, recreate it, and prove a note survives process replacement."),("Контейнер должен быть заменяемым, а состояние — иметь явный жизненный цикл. Named volume удобен для данных Docker; bind mount удобен, когда оператору нужен конкретный путь хоста. Ни один из них сам по себе не является резервной копией.","Контейнер ауыстырылатын, ал күйдің өмір циклі айқын болуы керек. Named volume Docker дерегіне ыңғайлы; bind mount операторға нақты хост жолы керек кезде пайдалы. Екеуінің ешқайсысы өздігінен сақтық көшірме емес.","A container should be replaceable while state has an explicit lifecycle. A named volume suits Docker-managed data; a bind mount suits a known host path. Neither one is a backup by itself."),"docker volume create cloudlab-practice\ndocker run -d --name cloudlab-data -p 8080:8080 -v cloudlab-practice:/data cloudlab:lesson-09\nsleep 1\ncurl -fsS -X POST -d 'note=survives-recreate' http://127.0.0.1:8080/notes >/dev/null\ndocker rm -f cloudlab-data\ndocker run -d --name cloudlab-data -p 8080:8080 -v cloudlab-practice:/data cloudlab:lesson-09\nsleep 1\ncurl -fsS http://127.0.0.1:8080/ | grep survives-recreate\ndocker rm -f cloudlab-data",("Заметка остаётся после удаления контейнера.","Контейнер жойылған соң жазба қалады.","The note remains after the container is removed."),("Напишите команды, которые создают named volume, записывают в него файл через временный контейнер и читают тот же файл через второй контейнер.","Named volume жасап, уақытша контейнер арқылы файл жазып, екінші контейнер арқылы сол файлды оқитын пәрмендерді жазыңыз.","Write commands that create a named volume, write a file through one temporary container, and read it through a second container.")),
("11-compose",("Docker Compose: сервисы, сеть и декларативный запуск","Docker Compose: сервистер, желі және декларативті іске қосу","Docker Compose: services, networking, and declarative startup"),("Поднимаем приложение и reverse proxy одной декларацией и обращаемся к сервису по DNS-имени внутри сети.","Қолданба мен reverse proxy-ді бір декларациямен көтеріп, желі ішінде сервиске DNS атымен қатынаймыз.","Start the application and reverse proxy from one declaration and reach a service by DNS name inside the network."),("Compose описывает желаемые сервисы, сети и тома. Внутри общей сети имя `app` разрешается встроенным DNS; публикация порта нужна только на границе с хостом. `depends_on` упорядочивает запуск, а healthcheck доказывает готовность.","Compose қажетті сервис, желі және томды сипаттайды. Ортақ желіде `app` аты ішкі DNS арқылы шешіледі; порт жариялау тек хост шекарасында керек. `depends_on` іске қосуды реттейді, healthcheck дайындықты дәлелдейді.","Compose declares desired services, networks, and volumes. Inside a shared network, `app` resolves through built-in DNS; publishing a port is needed only at the host boundary. `depends_on` orders startup, while a health check proves readiness."),"docker compose config\ndocker compose up -d --build\ndocker compose ps\ncurl -fsS http://127.0.0.1:8080/healthz\ndocker compose logs --tail=10 app",("Оба сервиса запущены, приложение healthy.","Екі сервис те іске қосылған, қолданба healthy.","Both services run and the app is healthy."),("Напишите команды, которые валидируют Compose, запускают стек, ждут успешного `/readyz` не более 30 секунд и при ошибке печатают последние 50 строк журнала.","Compose тексеріп, стекті іске қосатын, `/readyz` жауабын 30 секундтан артық күтпейтін және қатеде соңғы 50 журнал жолын шығаратын пәрмендерді жазыңыз.","Write commands that validate Compose, start the stack, wait at most 30 seconds for `/readyz`, and print the last 50 log lines on failure.")),
("12-health-resources",("Healthchecks, лимиты и корректная остановка","Healthcheck, лимит және дұрыс тоқтау","Health checks, limits, and graceful shutdown"),("Различаем liveness и readiness, задаём ресурсные границы и проверяем поведение при остановке.","Liveness пен readiness айырмасын көріп, ресурс шекарасын қойып, тоқтау кезіндегі мінезді тексереміз.","Distinguish liveness from readiness, set resource boundaries, and test shutdown behaviour."),("Liveness отвечает, нужно ли перезапустить процесс; readiness — можно ли направлять ему новый трафик. Ошибка в этой границе превращает временно занятой сервис в цикл перезапусков. Лимиты защищают соседей, но слишком низкий лимит сам создаёт отказ.","Liveness үдерісті қайта бастау керек пе, readiness жаңа трафик жіберуге бола ма дегенге жауап береді. Бұл шекарадағы қате уақытша бос емес сервисті қайта іске қосу цикліне айналдырады. Лимит көршілерді қорғайды, бірақ тым төмен шек өзі ақау жасайды.","Liveness asks whether the process should restart; readiness asks whether it should receive new traffic. Confusing them turns temporary load into a restart loop. Limits protect neighbours, but a limit set too low creates its own failure."),"docker compose ps\ndocker inspect --format '{{json .State.Health}}' cloud-devops-lab-app-1 2>/dev/null || true\ntime docker compose stop -t 10 app\ndocker compose up -d app\ncurl -fsS http://127.0.0.1:8080/readyz",("Остановка укладывается в grace period.","Тоқтау grace period ішінде аяқталады.","Shutdown completes within the grace period."),("Напишите POSIX sh-цикл с дедлайном, который ждёт readiness URL, делает паузу между попытками и возвращает ненулевой код после таймаута.","Readiness URL-ды deadline-ға дейін күтетін, әрекет арасында үзіліс жасайтын және timeout соңында нөл емес код қайтаратын POSIX sh циклін жазыңыз.","Write a POSIX sh loop with a deadline that waits for a readiness URL, pauses between attempts, and returns nonzero after timeout.")),
("13-registry",("Registry, immutable tags, SBOM и проверка цепочки поставки","Registry, өзгермейтін тег, SBOM және жеткізу тізбегін тексеру","Registry, immutable tags, SBOM, and supply-chain checks"),("Связываем образ с commit SHA, получаем SBOM и сканируем содержимое до публикации.","Бейнені commit SHA-мен байланыстырып, SBOM алып, жариялау алдында мазмұнды сканерлейміз.","Tie an image to a commit SHA, produce an SBOM, and scan its contents before publishing."),("Registry хранит и раздаёт образы, но доверие не появляется от самого факта загрузки. Версия должна быть неизменяемой, происхождение — прослеживаемым, зависимости — видимыми. SBOM перечисляет компоненты; сканер сопоставляет их с известными уязвимостями, но не доказывает отсутствие всех ошибок.","Registry бейнені сақтап таратады, бірақ жүктеу фактісі сенім туғызбайды. Нұсқа өзгермейтін, шығу тегі қадағаланатын, тәуелділік көрінетін болуы керек. SBOM компонентті тізеді; сканер оны белгілі осалдықпен салыстырады, бірақ барлық қатенің жоқтығын дәлелдемейді.","A registry stores and distributes images, but an upload does not create trust. A release needs an immutable version, traceable origin, and visible dependencies. An SBOM lists components; a scanner matches known vulnerabilities but cannot prove that no bug exists."),"sha=$(git rev-parse --short=12 HEAD)\ndocker build --build-arg VERSION=\"$sha\" -t \"cloudlab:$sha\" .\ndocker image inspect \"cloudlab:$sha\" --format '{{.Id}}'\nprintf 'release=%s\n' \"$sha\"",("Тег равен исходному commit SHA.","Тег бастапқы commit SHA-ға тең.","The tag equals the source commit SHA."),("Напишите POSIX sh-скрипт, который строит тег из 12 символов Git SHA, собирает образ и печатает его локальный image ID; `latest` использовать нельзя.","12 таңбалы Git SHA-дан тег құрып, бейне жинап, оның жергілікті image ID-сын басатын POSIX sh скриптін жазыңыз; `latest` қолданбаңыз.","Write a POSIX sh script that derives a tag from the 12-character Git SHA, builds the image, and prints its local image ID; do not use `latest`.")),
("14-cicd",("CI/CD как система обратной связи","CI/CD кері байланыс жүйесі ретінде","CI/CD as a feedback system"),("Проектируем pipeline из независимых проверок и отделяем непрерывную интеграцию от доставки и развёртывания.","Pipeline-ды тәуелсіз тексерулерден құрып, үздіксіз интеграцияны жеткізу мен орналастырудан ажыратамыз.","Design a pipeline from independent checks and separate continuous integration, delivery, and deployment."),("CI быстро сообщает, можно ли объединять изменение. Continuous delivery держит проверенный артефакт готовым к выпуску; continuous deployment автоматически выпускает каждое прошедшее изменение. Скорость без надёжных ворот лишь быстрее доставляет дефект.","CI өзгерісті біріктіруге болатынын тез хабарлайды. Continuous delivery тексерілген артефактты шығаруға дайын ұстайды; continuous deployment әр өткен өзгерісті автоматты шығарады. Сенімді қақпасыз жылдамдық ақауды тек тез жеткізеді.","CI quickly reports whether a change can be merged. Continuous delivery keeps a verified artifact releasable; continuous deployment releases every passing change automatically. Speed without trustworthy gates only delivers defects sooner."),"go test ./...\ngofmt -l cmd internal\ndocker build -t cloudlab:ci .\ndocker run --rm -d --name cloudlab-ci -p 18080:8080 cloudlab:ci\ntrap 'docker rm -f cloudlab-ci >/dev/null 2>&1 || true' EXIT\nsleep 1\ncurl -fsS http://127.0.0.1:18080/readyz",("Проверки идут от быстрых к более дорогим.","Тексеру жылдамнан қымбатқа қарай жүреді.","Checks run from fast to more expensive."),("Напишите POSIX sh pipeline из форматирования, теста, сборки образа и smoke test. Он должен останавливаться при первом сбое и всегда удалять тестовый контейнер.","Пішім, тест, бейне жинау және smoke test кезеңдері бар POSIX sh pipeline жазыңыз. Ол алғашқы қатеде тоқтап, тест контейнерін әрқашан жоюы тиіс.","Write a POSIX sh pipeline for formatting, tests, image build, and a smoke test. It must stop at the first failure and always remove its test container.")),
("15-github-actions",("GitHub Actions: проверка каждого изменения","GitHub Actions: әр өзгерісті тексеру","GitHub Actions: verify every change"),("Переносим локальные проверки в workflow, фиксируем permissions и используем cache без хранения секретов.","Жергілікті тексеруді workflow-ға көшіріп, permissions бекітеміз және құпия сақтамай cache қолданамыз.","Move local checks into a workflow, pin permissions, and use caching without storing secrets."),("Workflow — исполняемый контракт репозитория. Триггер определяет момент, permissions — полномочия токена, job — чистую машину, steps — доказательства. Действия стороннего автора тоже являются кодом; фиксируйте доверенную версию и сокращайте права.","Workflow — репозиторийдің орындалатын келісімі. Trigger уақытты, permissions токен құқығын, job таза машинаны, steps дәлелді анықтайды. Үшінші тарап action-ы да код; сенімді нұсқаны бекітіп, құқықты азайтыңыз.","A workflow is the repository's executable contract. A trigger selects the moment, permissions constrain the token, a job supplies a clean machine, and steps produce evidence. Third-party actions are code too; pin trusted versions and minimize permissions."),"cp workflow.example.yml /tmp/cloudlab-workflow.yml\nsed -n '1,220p' /tmp/cloudlab-workflow.yml\ngrep -n 'permissions:[|]go test[|]build-push-action' /tmp/cloudlab-workflow.yml",("Workflow имеет явные permissions.","Workflow ішінде айқын permissions бар.","The workflow has explicit permissions."),("Напишите shell-команды шага CI, которые переходят в каталог проекта, запускают `go test ./...` и отклоняют неотформатированные Go-файлы.","Жоба каталогына өтіп, `go test ./...` іске қосатын және пішімделмеген Go файлдарын қабылдамайтын CI қадамының shell пәрмендерін жазыңыз.","Write the shell commands for a CI step that enters the project, runs `go test ./...`, and rejects unformatted Go files.")),
("16-release-image",("Сборка и публикация одного проверенного образа","Бір тексерілген бейнені жинау және жариялау","Build and publish one verified image"),("Настраиваем публикацию в registry только после тестов и переносим один digest между средами.","Registry-ге тек тесттен кейін жариялап, орта арасында бір digest өткіземіз.","Publish to a registry only after tests and promote one digest across environments."),("Повторная сборка для production может дать иной артефакт, чем тот, который тестировался. Надёжнее собрать один раз, проверить, опубликовать и продвигать его digest. Registry credentials нужны только job публикации и не должны быть доступны pull request из недоверенного fork.","Production үшін қайта жинау тесттелгеннен өзге артефакт беруі мүмкін. Бір рет жинап, тексеріп, жариялап, digest-ті орта арасында жылжыту сенімді. Registry құпиясы тек publish job-қа керек және сенімсіз fork pull request-іне берілмеуі тиіс.","Rebuilding for production may produce an artifact different from the one tested. Build once, verify, publish, and promote its digest. Registry credentials belong only in the publishing job and must not reach an untrusted fork pull request."),"sha=$(git rev-parse --short=12 HEAD)\nimage=ghcr.io/OWNER/cloudlab:$sha\nprintf 'would publish %s\n' \"$image\"\nprintf 'production must record a sha256 digest, not latest\n'",("Имя содержит неизменяемую версию.","Атауда өзгермейтін нұсқа бар.","The name contains an immutable version."),("Напишите POSIX sh-проверку переменных `REGISTRY`, `IMAGE` и `GIT_SHA`, которая строит полное имя образа и отказывается публиковать тег `latest`.","`REGISTRY`, `IMAGE`, `GIT_SHA` айнымалыларын тексеріп, толық бейне атын құратын және `latest` тегін жариялаудан бас тартатын POSIX sh тексеруін жазыңыз.","Write a POSIX sh check for `REGISTRY`, `IMAGE`, and `GIT_SHA` that constructs the full image name and refuses to publish the tag `latest`.")),
("17-secrets",("Секреты, конфигурация и среды без утечки","Құпия, конфигурация және орта: ағып кетусіз","Secrets, configuration, and environments without leaks"),("Разделяем открытую конфигурацию и секреты, ограничиваем область токена и предотвращаем вывод значения в журнал.","Ашық конфигурация мен құпияны бөліп, токен ауқымын шектеп, мәннің журналға шығуын болдырмаймыз.","Separate public configuration from secrets, scope tokens narrowly, and prevent values from reaching logs."),("Секрет — не просто переменная окружения, а значение с владельцем, областью, сроком и процедурой ротации. Маскирование журнала помогает, но не отменяет риск передачи секрета процессу или стороннему action. Сначала уменьшайте полномочия и время жизни.","Құпия — жай орта айнымалысы емес; оның иесі, ауқымы, мерзімі және ротация жолы бар. Журналды бүркеу көмектеседі, бірақ құпияны үдеріске не бөтен action-ға беру тәуекелін жоймайды. Алдымен құқық пен өмір мерзімін азайтыңыз.","A secret is more than an environment variable: it has an owner, scope, lifetime, and rotation procedure. Log masking helps but does not remove the risk of passing a secret to a process or third-party action. Reduce privilege and lifetime first."),"test -f .gitignore\ngit grep -n -E '(BEGIN (RSA|OPENSSH) PRIVATE KEY|token=|password=)' -- ':!tools/course/generate_cloud_devops.py' || true\nprintf '%s\n' '.env must stay untracked'\ngit check-ignore .env || true",("Реальные значения не выводятся.","Нақты мәндер шығарылмайды.","No real value is printed."),("Напишите POSIX sh-скрипт, который требует `DEPLOY_TOKEN`, никогда его не печатает, проверяет непустое значение и очищает переменную перед завершением.","`DEPLOY_TOKEN` талап ететін, оны ешқашан баспайтын, бос еместігін тексеретін және аяқталғанда айнымалыны тазалайтын POSIX sh скриптін жазыңыз.","Write a POSIX sh script that requires `DEPLOY_TOKEN`, never prints it, checks that it is nonempty, and unsets it before exit.")),
("18-deploy-rollback",("Развёртывание, smoke test и проверенный rollback","Орналастыру, smoke test және тексерілген rollback","Deployment, smoke testing, and a verified rollback"),("Выпускаем по версии, проверяем сервис снаружи и возвращаем предыдущий digest по заранее написанной процедуре.","Нұсқа бойынша шығарып, сервисті сырттан тексеріп, алдын ала жазылған рәсіммен алдыңғы digest-ке қайтарамыз.","Deploy by version, test from outside, and restore the previous digest with a prewritten procedure."),("Rollback должен существовать до инцидента. Он возвращает код, но не всегда возвращает совместимые данные: миграции требуют отдельного плана. После переключения проверяйте пользовательский путь и наблюдайте метрики, а не считайте успешный exit code доказательством.","Rollback инцидентке дейін дайын болуы керек. Ол кодты қайтарады, бірақ дерек әрқашан үйлесімді болмауы мүмкін: migration үшін бөлек жоспар керек. Ауыстырған соң пайдаланушы жолын және метриканы тексеріңіз; сәтті exit code жеткіліксіз.","A rollback must exist before an incident. It restores code but may not restore compatible data; migrations need a separate plan. After switching, test a user path and watch metrics rather than treating a successful exit code as proof."),"export CLOUDLAB_IMAGE=ghcr.io/OWNER/cloudlab:REPLACE_WITH_SHA\nexport CLOUDLAB_DOMAIN=cloudlab.example.com\ndocker compose -f compose.prod.yaml config >/tmp/cloudlab-prod.yml\ngrep -n 'image:[|]CLOUDLAB_ENV' /tmp/cloudlab-prod.yml\nprintf 'record previous digest before docker compose up -d\n'",("Конфигурация не использует latest.","Конфигурация latest қолданбайды.","The configuration does not use latest."),("Напишите POSIX sh-функцию deploy: сохранить текущий digest в файл, применить новый образ, выполнить smoke test и при его ошибке вернуть сохранённый образ.","Ағымдағы digest-ті файлға сақтап, жаңа бейнені қолданып, smoke test жасап, қате болса сақталған бейнеге қайтаратын deploy POSIX sh функциясын жазыңыз.","Write a POSIX sh deploy function that saves the current digest, applies a new image, runs a smoke test, and restores the saved image if that test fails.")),
("19-cloud",("Облако: IaaS, стоимость, IAM и общая ответственность","Бұлт: IaaS, шығын, IAM және ортақ жауапкершілік","Cloud: IaaS, cost, IAM, and shared responsibility"),("Проектируем минимальную виртуальную инфраструктуру, бюджет и доступ до создания оплачиваемого ресурса.","Ақылы ресурс жасамай тұрып минимал виртуалды инфрақұрылым, бюджет және қолжетімділікті жобалаймыз.","Design minimal virtual infrastructure, a budget, and access before creating a billable resource."),("В IaaS провайдер защищает физическую площадку и слой виртуализации, а клиент отвечает за ОС, доступы, firewall, данные и приложение. Бесплатный кредит не отменяет счёт после его окончания. Перед созданием ресурса запишите регион, размер, диск, трафик, резервные копии и условие удаления.","IaaS-та провайдер физикалық алаң мен виртуализация қабатын қорғайды, ал клиент ОС, рұқсат, firewall, дерек және қолданбаға жауап береді. Тегін кредит біткеннен кейінгі шотты жоймайды. Ресурс алдында аймақ, өлшем, диск, трафик, backup және жою шартын жазыңыз.","In IaaS, the provider protects the facility and virtualization layer, while the customer owns the OS, access, firewall, data, and application. Free credit does not cancel later billing. Before creating anything, record region, size, disk, traffic, backups, and the deletion condition."),"cat > /tmp/cloud-budget.txt <<'EOF'\nresource=one small Linux VM\nowner=student\nexpires=after lesson 22\npublic_ports=22,80,443\nbudget_alert=set before creation\nEOF\ncat /tmp/cloud-budget.txt\nprintf 'No cloud resource is created by this lesson.\n'",("План имеет владельца и срок удаления.","Жоспарда ие мен жою мерзімі бар.","The plan has an owner and deletion date."),("Напишите POSIX sh-проверку файла плана: обязательны строки `owner=`, `expires=` и `budget_alert=`; при отсутствии любой строки вернуть ошибку.","Жоспар файлын тексеретін POSIX sh жазыңыз: `owner=`, `expires=`, `budget_alert=` жолдары міндетті, бірі жоқ болса қате қайтарылсын.","Write a POSIX sh check for a plan file that requires `owner=`, `expires=`, and `budget_alert=` lines and fails if any is missing.")),
("20-server",("Ubuntu-сервер: SSH-ключи, обновления и firewall","Ubuntu сервері: SSH кілті, жаңарту және firewall","Ubuntu server: SSH keys, updates, and a firewall"),("Подключаемся непривилегированным пользователем, закрываем лишние порты и проверяем доступ второй сессией.","Артық құқығы жоқ пайдаланушымен қосылып, қажетсіз портты жауып, екінші сессиямен қолжетімділікті тексереміз.","Connect as an unprivileged user, close unnecessary ports, and verify access from a second session."),("Укрепление сервера делается в безопасном порядке: сначала новый пользователь и ключ, затем проверка второго входа, после этого запрет рискованного доступа. Если включить firewall или отключить root раньше проверки, можно закрыть себе сервер. Консоль провайдера остаётся аварийным каналом.","Серверді қауіпсіз ретпен нығайтыңыз: алдымен жаңа пайдаланушы мен кілт, кейін екінші кіруді тексеру, содан соң қауіпті қолжетімділікті жабу. Firewall-ды ерте қоссаңыз, өзіңізді серверден жауып қаласыз. Провайдер консолі апаттық арна болып қалады.","Harden a server in a safe order: create a user and key, verify a second login, then disable risky access. Enabling a firewall or disabling root too early can lock you out. The provider console remains an emergency path."),"ssh-keygen -t ed25519 -a 64 -f ~/.ssh/cloudlab_ed25519 -C cloudlab\nprintf 'Copy only cloudlab_ed25519.pub to the server.\n'\nprintf 'On Ubuntu: allow OpenSSH before enabling ufw.\n'\nprintf 'Keep the current session open while testing a second login.\n'",("Приватный ключ не копируется на сервер.","Жеке кілт серверге көшірілмейді.","The private key is never copied to the server."),("Напишите локальный POSIX sh-скрипт, который проверяет наличие private и public key, требует режим 600 у private key и печатает только fingerprint публичного ключа.","Private және public key барын, private key режимі 600 екенін тексеріп, тек public key fingerprint басатын жергілікті POSIX sh скриптін жазыңыз.","Write a local POSIX sh script that checks for private and public key files, requires mode 600 on the private key, and prints only the public-key fingerprint.")),
("21-https",("Домен, reverse proxy и автоматический HTTPS","Домен, reverse proxy және автоматты HTTPS","Domain, reverse proxy, and automatic HTTPS"),("Направляем DNS на сервер, отдаём CloudLab через Caddy и проверяем цепочку TLS снаружи.","DNS-ті серверге бағыттап, CloudLab-ты Caddy арқылы беріп, TLS тізбегін сырттан тексереміз.","Point DNS at the server, serve CloudLab through Caddy, and verify the TLS chain externally."),("Reverse proxy принимает публичное соединение, завершает TLS и передаёт запрос частному приложению. Для автоматического сертификата домен должен разрешаться в правильный IP, а порты 80 и 443 — быть доступны. Не публикуйте порт приложения наружу без необходимости.","Reverse proxy ашық байланысты қабылдап, TLS аяқтап, сұрауды жеке қолданбаға береді. Автоматты сертификат үшін домен дұрыс IP-ға шешіліп, 80 және 443 порттары ашық болуы керек. Қолданба портын қажетсіз сыртқа шығармаңыз.","A reverse proxy accepts the public connection, terminates TLS, and forwards the request to the private application. Automatic certificates require the domain to resolve to the right IP and ports 80 and 443 to be reachable. Do not expose the application port publicly without need."),"export CLOUDLAB_DOMAIN=cloudlab.example.com\ngetent hosts \"$CLOUDLAB_DOMAIN\" 2>/dev/null || true\nprintf 'Caddy route:\n'\nsed -n '1,80p' deploy/caddy/Caddyfile\nprintf 'After real DNS: curl -v https://%s/readyz\n' \"$CLOUDLAB_DOMAIN\"",("Приложение доступно через proxy, а не публичный порт 8080.","Қолданба 8080 ашық портымен емес, proxy арқылы қолжетімді.","The app is reached through the proxy, not a public port 8080."),("Напишите POSIX sh-проверку HTTPS URL: потребовать успешную проверку сертификата, HTTP 200 от `/readyz` и вывести дату окончания сертификата без `-k`.","HTTPS URL тексеруін жазыңыз: сертификат тексеруі сәтті, `/readyz` HTTP 200 болуы керек және `-k` қолданбай сертификаттың аяқталу күнін шығарыңыз.","Write a POSIX sh check for an HTTPS URL that requires valid certificate verification and HTTP 200 from `/readyz`, then prints the certificate expiry without `-k`.")),
("22-iac",("OpenTofu и Ansible: инфраструктура и настройка как код","OpenTofu және Ansible: инфрақұрылым мен баптау код ретінде","OpenTofu and Ansible: infrastructure and configuration as code"),("Валидируем описание ресурса, формируем inventory и применяем идемпотентную настройку хоста.","Ресурс сипаттамасын тексеріп, inventory құрып, хостқа идемпотентті баптау қолданамыз.","Validate a resource description, produce inventory, and apply idempotent host configuration."),("OpenTofu управляет жизненным циклом инфраструктурных ресурсов через state; Ansible приводит ОС к нужной конфигурации по SSH. Код делает изменение обозримым, но state и секреты требуют защиты. Идемпотентный playbook при повторном запуске не должен каждый раз сообщать об изменении.","OpenTofu state арқылы инфрақұрылым ресурсының өмір циклін басқарады; Ansible SSH арқылы ОС-ты қажетті күйге әкеледі. Код өзгерісті көрінетін етеді, бірақ state пен құпияны қорғау керек. Идемпотентті playbook қайта іске қосылған сайын өзгеріс жасамауы тиіс.","OpenTofu manages infrastructure resource lifecycles through state; Ansible brings the OS to the desired configuration over SSH. Code makes changes reviewable, but state and secrets need protection. An idempotent playbook should report no changes on a second run."),"tofu -chdir=infra/opentofu fmt -check 2>/dev/null || true\ntofu -chdir=infra/opentofu init -backend=false 2>/dev/null || true\ntofu -chdir=infra/opentofu validate 2>/dev/null || true\nansible-playbook --syntax-check -i infra/ansible/inventory.ini.example infra/ansible/site.yml 2>/dev/null || true",("Проверка не создаёт облачный ресурс.","Тексеру бұлт ресурсын жасамайды.","Validation creates no cloud resource."),("Напишите POSIX sh preflight, который требует `tofu` и `ansible-playbook`, выполняет fmt/validate и syntax-check, но никогда не вызывает `apply`.","`tofu` және `ansible-playbook` талап етіп, fmt/validate және syntax-check орындайтын, бірақ `apply` ешқашан шақырмайтын POSIX sh preflight жазыңыз.","Write a POSIX sh preflight that requires `tofu` and `ansible-playbook`, runs fmt/validate and syntax-check, and never calls `apply`.")),
("23-kubernetes",("Kubernetes: Deployment, Service, probes и rollout","Kubernetes: Deployment, Service, probe және rollout","Kubernetes: Deployment, Service, probes, and rollout"),("Читаем манифесты CloudLab, строим их через Kustomize и наблюдаем безопасное обновление и откат.","CloudLab манифесттерін оқып, Kustomize арқылы құрып, қауіпсіз жаңарту мен кері қайтаруды бақылаймыз.","Read the CloudLab manifests, build them with Kustomize, and observe a safe update and rollback."),("Deployment управляет желаемым числом Pod и обновлением, Service даёт стабильный адрес, probes управляют трафиком и перезапуском, ConfigMap отделяет открытую конфигурацию. PersistentVolumeClaim сохраняет данные, но один локальный JSON-файл не становится распределённой базой: учебный Deployment оставляет одну реплику.","Deployment қажетті Pod санын және жаңартуды, Service тұрақты мекенжайды, probe трафик пен қайта іске қосуды, ConfigMap ашық конфигурацияны басқарады. PersistentVolumeClaim деректі сақтайды, бірақ бір JSON файл таратылған база болмайды: оқу Deployment бір replica ұстайды.","A Deployment manages desired Pods and updates, a Service supplies a stable address, probes control traffic and restarts, and a ConfigMap separates public configuration. A PersistentVolumeClaim preserves data, but one JSON file does not become a distributed database, so this learning Deployment keeps one replica."),"kubectl kustomize deploy/k8s >/tmp/cloudlab-k8s.yaml\ngrep -n 'kind: Deployment[|]readinessProbe[|]runAsNonRoot[|]resources:' /tmp/cloudlab-k8s.yaml\nprintf 'With a local cluster: kubectl apply -k deploy/k8s\n'\nprintf 'Then: kubectl -n cloudlab rollout status deployment/cloudlab\n'",("Собранный YAML содержит probes и securityContext.","Жиналған YAML ішінде probe және securityContext бар.","The rendered YAML includes probes and a security context."),("Напишите POSIX sh-проверку, которая строит Kustomize YAML, требует наличие Deployment, Service, readinessProbe и `runAsNonRoot: true`, затем запускает client-side dry-run.","Kustomize YAML құрып, Deployment, Service, readinessProbe және `runAsNonRoot: true` болуын талап етіп, client-side dry-run жасайтын POSIX sh тексеруін жазыңыз.","Write a POSIX sh check that renders Kustomize YAML, requires a Deployment, Service, readinessProbe, and `runAsNonRoot: true`, then performs a client-side dry run.")),
("24-incident",("Наблюдаемость, резервная копия и учебный инцидент","Бақылау, сақтық көшірме және оқу инциденті","Observability, backups, and a recovery drill"),("Определяем показатели сервиса, создаём проверяемую копию, восстанавливаем данные и оформляем итоговый отчёт без поиска виноватого.","Сервис көрсеткішін анықтап, тексерілетін көшірме жасап, деректі қалпына келтіріп, кінә іздемейтін қорытынды есеп құрамыз.","Define service indicators, create a verifiable backup, restore data, and write a blameless final report."),("Логи отвечают, что произошло с отдельным событием; метрики показывают изменение чисел во времени; трассировка связывает путь запроса. Резервная копия считается полезной только после восстановления. Инцидент заканчивается не перезапуском, а восстановленной услугой, сохранёнными доказательствами и действиями, уменьшающими повторение.","Журнал жеке оқиғаға не болғанын, метрика санның уақыт бойынша өзгеруін, trace сұрау жолын көрсетеді. Сақтық көшірме тек қалпына келтірілген соң пайдалы деп саналады. Инцидент restart-пен емес, қызмет қалпына келіп, дәлел сақталып, қайталану қаупі азайғанда аяқталады.","Logs explain individual events, metrics show numbers changing over time, and traces connect a request path. A backup is useful only after a restore. An incident ends after service is restored, evidence is preserved, and follow-up work reduces recurrence, not merely after a restart."),"curl -fsS http://127.0.0.1:8080/metrics 2>/dev/null || true\nmkdir -p backups\n./scripts/backup.sh 2>/dev/null || true\nprintf 'A real drill restores into an isolated volume first.\n'\nprintf 'Record: detection, impact, timeline, recovery, evidence, follow-up owner and date.\n'",("Метрика связана с пользовательским результатом.","Метрика пайдаланушы нәтижесімен байланысты.","The metric is connected to a user outcome."),("Напишите POSIX sh-скрипт резервного копирования каталога: архив с UTC-временем, SHA-256 рядом, проверка архива и отказ удалить исходник. Затем перечислите команды тестового восстановления в новый каталог.","Каталогтың сақтық көшірмесін жасайтын POSIX sh жазыңыз: UTC уақыты бар архив, жанында SHA-256, архивті тексеру және бастапқыны жоймау. Кейін жаңа каталогқа сынақ қалпына келтіру пәрмендерін көрсетіңіз.","Write a POSIX sh directory backup script with a UTC-stamped archive, adjacent SHA-256, archive verification, and no source deletion. Then list commands for a test restore into a new directory.")),
]

for stem, title3, summary3, concept3, cmd, check3, task3 in MORE:
    title = dict(zip(("ru","kz","en"), title3))
    summary = dict(zip(("ru","kz","en"), summary3))
    concept = dict(zip(("ru","kz","en"), concept3))
    check_main = dict(zip(("ru","kz","en"), check3))
    checks = {
        "ru": [check_main["ru"], "Команда завершается предсказуемым кодом.", "Результат можно повторить на чистой среде."],
        "kz": [check_main["kz"], "Пәрмен болжамды кодпен аяқталады.", "Нәтижені таза ортада қайталауға болады."],
        "en": [check_main["en"], "The command exits with a predictable status.", "The result can be repeated in a clean environment."],
    }
    task = dict(zip(("ru","kz","en"), task3))
    TOPICS.append({"stem":stem,"title":title,"summary":summary,"concept":concept,"cmd":cmd,"checks":checks,"task":task})

# A beginner needs the workshop before the application and the network model
# before endpoint diagnostics.  Slugs keep their stable names, while the course
# order follows the learning dependency rather than the date they were drafted.
BEGINNER_ORDER = (
    "01-why-now", "03-terminal-git", "04-linux-access", "05-processes-logs",
    "06-networks-dns-http", "02-service-map", "07-containers",
    "08-images-lifecycle", "09-dockerfile", "10-persistent-data",
    "11-compose", "12-health-resources", "13-registry", "14-cicd",
    "15-github-actions", "16-release-image", "17-secrets",
    "18-deploy-rollback", "19-cloud", "20-server", "21-https", "22-iac",
    "23-kubernetes", "24-incident",
)
_topic_by_stem = {topic["stem"]: topic for topic in TOPICS}
assert set(_topic_by_stem) == set(BEGINNER_ORDER)
TOPICS = [_topic_by_stem[stem] for stem in BEGINNER_ORDER]


def suffix(lang):
    return "" if lang == "ru" else f"-{lang}"


def render_lesson(topic, number, lang):
    l = LOCALE[lang]
    p = PACKS[topic["stem"]]
    points = "\n".join(f"- {x}" for x in topic["checks"][lang])
    glossary = "\n".join(
        f"- **{item['name'][lang]}** — {item['text'][lang]}" for item in p["terms"]
    )
    alt = f"{l['map']} {number}"
    map_url = f"/static/course/cloud-devops/map-{topic['stem']}-{lang}.svg"
    course_url = f"/course/cloud-devops?lang={lang}"
    repo = "https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab"
    command = re.sub(r"(printf '[^'\n]*)\n'", r"\1\\n'", topic["cmd"])
    previous = TOPICS[number - 2] if number > 1 else None
    following = TOPICS[number] if number < len(TOPICS) else None
    previous_signal = PACKS[previous["stem"]]["signal"][lang] if previous else ""
    previous_title = previous["title"][lang] if previous else ""
    next_title = following["title"][lang] if following else ""

    ui = {
        "ru": {
            "for_you": "Этот урок начинается с нуля. Незнакомое слово здесь не считается вашим пробелом: мы сначала создадим образ, затем дадим точное значение и только потом применим его.",
            "where": "Где мы находимся",
            "first": "Это первая опора курса. Сейчас важнее увидеть весь путь, чем запомнить названия инструментов.",
            "recall": "Полётное повторение. В прошлом уроке была опора",
            "recall_do": "Не подглядывая, назовите её четыре звена и только затем сравните с записью.",
            "image": "Сначала знакомый образ",
            "limit": "Аналогия помогает начать рассуждение, но не заменяет устройство технологии. Ниже мы уточним каждое слово.",
            "signal": "Опорный сигнал урока",
            "signal_do": "Прочитайте цепочку слева направо. Она отвечает не на вопрос «какую команду запомнить», а на вопрос «почему следующий шаг следует из предыдущего».",
            "words": "Новые слова простыми словами",
            "slow": "Разбираем без спешки",
            "slow2": "Сейчас не нужно запоминать формулировку дословно. Найдите в ней причинную связь и привяжите её к опорному сигналу выше.",
            "predict": "Сначала предскажите результат",
            "predict_text": "До запуска ответьте на бумаге: что должно измениться, что должно остаться прежним и какой вывод команды подтвердит ваш прогноз? Даже неверный прогноз полезен, если после опыта вы можете объяснить расхождение.",
            "steps": "Опыт маленькими шагами",
            "run": "Запускайте блок по одной смысловой группе. После каждой остановитесь и сопоставьте результат со своим прогнозом.",
            "explain": "Что именно делает этот опыт",
            "seen": "Что должно получиться",
            "map": "Соберите целое по опорной схеме",
            "map_do": "Проведите пальцем или карандашом по четырём блокам и расскажите весь урок одним связным объяснением.",
            "aloud": "Воспроизведите без подсказки",
            "aloud_items": ["Закройте схему и нарисуйте четыре блока по памяти.", "Объясните каждый переход словами «потому что».", "Дайте определение одному новому термину, не повторяя текст дословно.", "Назовите один сигнал, который доказывает результат."],
            "mistake": "Найдите и исправьте ошибку",
            "task_intro": "Теперь соберите тот же смысл самостоятельно. Можно смотреть на опорную схему; цель — не экзамен на память, а воспроизводимый результат.",
            "criteria": "Критерии готовности",
            "criteria_items": ["Вы понимаете каждую запускаемую строку.", "Повторный запуск даёт ожидаемый результат или безопасно объяснённое отличие.", "В выводе нет настоящих ключей, токенов и паролей.", "Вы можете показать факт, который подтверждает успех."],
            "retry": "Если проверка не прошла, зафиксируйте ожидаемое и фактическое, вернитесь к одному разорванному звену опоры и повторите. Ошибка не отнимает у вас право продолжать обучение.",
            "next": "Открытая перспектива",
            "next_text": "Следующий урок добавит новую опору",
            "finish": "Это финальный урок. Теперь восстановите всю дорогу курса: изменение, проверка, выпуск, наблюдение и доказанное восстановление.",
            "submit": "Отправьте только команды или скрипт POSIX sh без реальных ключей и токенов.",
            "optional_text": "Запишите один возможный отказ и сигнал, по которому вы его заметите.",
            "files": "Эталонные файлы проекта", "contents": "Оглавление курса",
        },
        "kz": {
            "for_you": "Бұл сабақ нөлден басталады. Бейтаныс сөз сіздің кемшілігіңіз емес: алдымен бейне құрамыз, кейін дәл мағына беріп, содан соң қолданамыз.",
            "where": "Қай жерде тұрмыз", "first": "Бұл курстың алғашқы тірегі. Қазір құрал атауын жаттаудан бұрын бүкіл жолды көру маңызды.",
            "recall": "Қысқа қайталау. Өткен сабақтың тірегі", "recall_do": "Қарамай төрт буынды атаңыз, содан кейін ғана жазбамен салыстырыңыз.",
            "image": "Алдымен таныс бейне", "limit": "Аналогия ойды бастауға көмектеседі, бірақ технология құрылысын алмастырмайды. Төменде әр сөзді нақтылаймыз.",
            "signal": "Сабақтың тірек сигналы", "signal_do": "Тізбекті солдан оңға оқыңыз. Ол «қай пәрменді жаттаймын?» емес, «келесі қадам неге алдыңғысынан шығады?» дегенге жауап береді.",
            "words": "Жаңа сөздер қарапайым тілмен", "slow": "Асықпай талдаймыз", "slow2": "Анықтаманы сөзбе-сөз жаттамаңыз. Себептік байланысты тауып, жоғарыдағы тірекпен қосыңыз.",
            "predict": "Алдымен нәтижені болжаңыз", "predict_text": "Іске қоспай тұрып қағазға жазыңыз: не өзгеруі, не сол күйде қалуы және қай output болжамды дәлелдеуі тиіс? Қате болжам да айырманы түсіндіре алсаңыз пайдалы.",
            "steps": "Шағын қадаммен тәжірибе", "run": "Блокты мағыналық топпен орындаңыз. Әр топтан кейін тоқтап, нәтижені болжаммен салыстырыңыз.",
            "explain": "Бұл тәжірибе нақты не істейді", "seen": "Не шығуы тиіс", "map": "Тірек сызбамен тұтасты жинаңыз", "map_do": "Төрт блок бойынша жүріп, бүкіл сабақты бір байланысқан әңгімемен түсіндіріңіз.",
            "aloud": "Көмексіз жаңғыртыңыз", "aloud_items": ["Сызбаны жауып, төрт блокты жатқа сызыңыз.", "Әр өтуді «себебі» сөзімен түсіндіріңіз.", "Бір жаңа терминді мәтінді қайталамай анықтаңыз.", "Нәтижені дәлелдейтін бір сигналды атаңыз."],
            "mistake": "Қатені тауып түзетіңіз", "task_intro": "Енді сол мағынаны өзіңіз жинаңыз. Тірек сызбаға қарауға болады; мақсат — жады емтиханы емес, қайталанатын нәтиже.",
            "criteria": "Дайындық өлшемдері", "criteria_items": ["Әр іске қосылатын жолды түсінесіз.", "Қайталау күтілетін нәтиже не қауіпсіз түсіндірілген айырма береді.", "Output ішінде нақты key, token және password жоқ.", "Success-ті дәлелдейтін факт көрсете аласыз."],
            "retry": "Тексеру өтпесе, күтілген және нақты нәтижені жазып, тіректің бір үзілген буынына оралып қайталаңыз. Қате оқу құқығын алып қоймайды.",
            "next": "Ашық перспектива", "next_text": "Келесі сабақ жаңа тірек қосады", "finish": "Бұл соңғы сабақ. Енді бүкіл жолды қалпына келтіріңіз: өзгеріс, тексеру, шығарылым, бақылау және дәлелденген restore.",
            "submit": "Тек POSIX sh пәрменін не скриптін, нақты кілт пен токенсіз жіберіңіз.", "optional_text": "Бір ықтимал ақауды және оны байқайтын сигналды жазыңыз.",
            "files": "Жобаның эталон файлдары", "contents": "Курс мазмұны",
        },
        "en": {
            "for_you": "This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.",
            "where": "Where we are", "first": "This is the course's first support. Seeing the whole path matters more now than memorising tool names.",
            "recall": "Flight review. The previous lesson used the support", "recall_do": "Without looking, name its four links, then compare with the written signal.",
            "image": "Begin with a familiar image", "limit": "An analogy starts reasoning but does not replace the technology. We make every word precise below.",
            "signal": "The lesson's support signal", "signal_do": "Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”",
            "words": "New words in plain language", "slow": "Take it apart without rushing", "slow2": "Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.",
            "predict": "Predict the result first", "predict_text": "Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.",
            "steps": "Experiment in small steps", "run": "Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.",
            "explain": "What this experiment actually does", "seen": "What you should observe", "map": "Rebuild the whole from the support map", "map_do": "Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.",
            "aloud": "Recall without a prompt", "aloud_items": ["Hide the map and draw its four blocks from memory.", "Explain every transition using the word “because”.", "Define one new term without repeating the text verbatim.", "Name one signal that proves the result."],
            "mistake": "Find and correct the mistake", "task_intro": "Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.",
            "criteria": "Completion criteria", "criteria_items": ["You understand every line you run.", "A repeat gives the expected result or a safely explained difference.", "Output contains no real keys, tokens, or passwords.", "You can point to evidence that proves success."],
            "retry": "If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.",
            "next": "Open perspective", "next_text": "The next lesson adds a new support", "finish": "This is the final lesson. Reconstruct the full course path: change, verification, release, observation, and proven recovery.",
            "submit": "Submit only POSIX sh commands or a script, with no real keys or tokens.", "optional_text": "Write down one possible failure and the signal that would reveal it.",
            "files": "Project reference files", "contents": "Course contents",
        },
    }[lang]

    recall = ui["first"] if previous is None else f"{ui['recall']} **{previous_signal}**. {ui['recall_do']}"
    aloud = "\n".join(f"{i}. {item}" for i, item in enumerate(ui["aloud_items"], 1))
    criteria = "\n".join(f"- {item}" for item in ui["criteria_items"])
    perspective = ui["finish"] if following is None else f"{ui['next_text']}: **{next_title}**."
    return f"""# {topic['title'][lang]}

_{l['lead']}:_ **{topic['summary'][lang]}**

> {ui['for_you']}

## {ui['where']}

{recall}

## {l['outcome']}

{topic['summary'][lang]} {l['done']}

## {ui['image']}

{p['analogy'][lang]}

{ui['limit']}

## {ui['signal']}

**{p['signal'][lang]}**

{ui['signal_do']}

## {ui['words']}

{glossary}

## {ui['slow']}

{topic['concept'][lang]}

{ui['slow2']}

## {ui['predict']}

{ui['predict_text']}

## {ui['steps']}

{l['run']} {l['safe']} {ui['run']}

```shell
{command}
```

[{ui['files']}]({repo}).

### {ui['explain']}

{p['command'][lang]}

## {ui['seen']}

{points}

## {ui['map']}

![{alt}]({map_url})

{ui['map_do']}

## {ui['aloud']}

{aloud}

## {ui['mistake']}

{p['mistake'][lang]}

## {l['task']}

{ui['task_intro']}

**{l['required']}** {topic['task'][lang]} {ui['submit']}

### {ui['criteria']}

{criteria}

{ui['retry']}

**{l['optional']}** {ui['optional_text']}

## {ui['next']}

{perspective}

[{ui['contents']}]({course_url})
"""


def svg_label(value, x):
    """Centre a support label on one line, or split a long phrase over two."""
    value = escape(value)
    if len(value) <= 15 or " " not in value:
        return f'<text x="{x}" y="232" fill="#17233c">{value}</text>'
    words = value.split()
    cut = min(
        range(1, len(words)),
        key=lambda i: abs(len(" ".join(words[:i])) - len(" ".join(words[i:]))),
    )
    top, bottom = " ".join(words[:cut]), " ".join(words[cut:])
    return (f'<text x="{x}" fill="#17233c">'
            f'<tspan x="{x}" y="218">{top}</tspan>'
            f'<tspan x="{x}" y="248">{bottom}</tspan></text>')


def render_map(topic, number, lang):
    title = escape(topic["title"][lang])
    nodes = PACKS[topic["stem"]]["nodes"][lang]
    desc = escape(", ".join(nodes))
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="420" viewBox="0 0 1200 420" role="img" aria-labelledby="title desc">
<title id="title">{title}</title><desc id="desc">{desc}</desc>
<rect width="1200" height="420" rx="28" fill="#f4f7fb"/><text x="60" y="72" font-family="system-ui,sans-serif" font-size="30" font-weight="700" fill="#17233c">{number:02d}. {title}</text>
<g font-family="system-ui,sans-serif" font-size="23" text-anchor="middle">
<rect x="35" y="160" width="240" height="120" rx="22" fill="#dbeafe" stroke="#2563eb" stroke-width="3"/>{svg_label(nodes[0], 155)}
<path d="M285 220h35" stroke="#64748b" stroke-width="5"/><path d="m310 207 20 13-20 13" fill="#64748b"/>
<rect x="330" y="160" width="240" height="120" rx="22" fill="#dcfce7" stroke="#16a34a" stroke-width="3"/>{svg_label(nodes[1], 450)}
<path d="M580 220h35" stroke="#64748b" stroke-width="5"/><path d="m605 207 20 13-20 13" fill="#64748b"/>
<rect x="625" y="160" width="240" height="120" rx="22" fill="#fef3c7" stroke="#d97706" stroke-width="3"/>{svg_label(nodes[2], 745)}
<path d="M875 220h35" stroke="#64748b" stroke-width="5"/><path d="m900 207 20 13-20 13" fill="#64748b"/>
<rect x="920" y="160" width="245" height="120" rx="22" fill="#fee2e2" stroke="#dc2626" stroke-width="3"/>{svg_label(nodes[3], 1042)}
</g></svg>'''


def main():
    LESSONS.mkdir(parents=True, exist_ok=True)
    MAPS.mkdir(parents=True, exist_ok=True)
    for lang, text in PREFACE.items():
        marker = {"ru": "## Сквозной проект", "kz": "## Ортақ жоба", "en": "## The project"}[lang]
        expanded = text.replace(marker, PREFACE_TEACHING[lang].strip() + "\n\n" + marker)
        (LESSONS / f"preface{suffix(lang)}.md").write_text(expanded.strip()+"\n", encoding="utf-8")
    for number, topic in enumerate(TOPICS, 1):
        assert number <= 24
        for lang in ("ru", "kz", "en"):
            (LESSONS / f"{topic['stem']}{suffix(lang)}.md").write_text(render_lesson(topic, number, lang), encoding="utf-8")
            name = f"map-{topic['stem']}-{lang}.svg"
            (MAPS / name).write_text(render_map(topic, number, lang), encoding="utf-8")
    print(f"generated {1 + len(TOPICS)} pages x 3 languages and {len(TOPICS) * 3} maps")


if __name__ == "__main__":
    main()
