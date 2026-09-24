# 28-сабақ. Уақытты, жадты және шығынды өлшейміз

## Бұл не үшін керек?

Шәйнекті алғаш қосқанда су біраз уақыт қызады, ал дайын ыстық суды құю тез болады. Бағдарламаның жаңа іске қосылуы мен дайын функцияның қайта шақырылуы да екі бөлек уақыт. Екеуін, файл көлемін және байқалған жад бөлінуін өлшейміз.

## Кодқа дейін

28-қадамдағы `bench.py`, `checks.py`, `examples.json` файлдарын бірге алыңыз. `subprocess.run` жаңа Python үрдісін бес рет ашады; `sys.executable` сол Python нұсқасын таңдайды. `-c` берілген код жолын орындайды, `capture_output=True` шығысты сақтайды, `text=True` оны мәтін ретінде оқиды, `check=True` қате болса өлшеуді тоқтатады. `range(5)` бес қайталауға бес сан береді; `range(200)` да 200 қайталауға сан береді. `time.perf_counter_ns()` уақытты наносекундпен өлшейді: бір миллисекундта миллион наносекунд бар. `statistics.median` реттелген өлшемдердің ортасын алады; бір кездейсоқ баяу нәтиже азырақ әсер етеді. Бір жылыту шақыруынан кейін 200 рет қайталау дайын функцияны өлшейді. `tracemalloc.start()` **Python** жад бөлулерін бақылайды, `get_traced_memory()` ағымдағы және ең үлкен байқалған көлемді қайтарады, `stop()` бақылауды тоқтатады. `Path.stat().st_size` екі файлдың байт санын береді.

Осы қадамның файлдары: `examples.json`, `checks.py`, `bench.py`.

Жоба түбірінен қадам қалтасына өтіп, бағдарламаны іске қосыңыз:

```text
cd course/kazakh-ai/step-28
python3 bench.py
```

Windows жүйесінде `python3` орнына `py` жазыңыз.

```python
import statistics
import subprocess
import sys
import time
import tracemalloc
from pathlib import Path

from checks import classify

question = "Шахмат қашан?"
cold = []
for repeat in range(5):
    start = time.perf_counter_ns()
    subprocess.run([sys.executable, "-c", "from checks import classify; classify('Шахмат қашан?')"],
                   check=True, capture_output=True, text=True)
    cold.append(time.perf_counter_ns() - start)
classify(question)
times = []
tracemalloc.start()
for repeat in range(200):
    start = time.perf_counter_ns()
    classify(question)
    times.append(time.perf_counter_ns() - start)
current, peak = tracemalloc.get_traced_memory()
tracemalloc.stop()
files = Path("checks.py").stat().st_size + Path("examples.json").stat().st_size
print("Сұрау саны:", len(times))
print("Жаңа іске қосу, мс:", round(statistics.median(cold) / 1000000, 3))
print("Ортаңғы уақыт, мс:", round(statistics.median(times) / 1000000, 3))
print("Ең көп бақыланған бөлу, байт:", peak)
print("Код пен дерек, байт:", files)
print("Тікелей API төлемі: 0; құрылғы мен еңбек құны есептелмеді")
```

[Қадам файлдары](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-28).

## Бағдарлама қалай жұмыс істейді?

Сандар компьютерге, жүктемеге және Python нұсқасына тәуелді. Сондықтан үлгі нәтиже — жол атаулары мен оң өлшемдер; нақты миллисекунд саны емес. Дайын функция уақыты `examples.json` оқуын, Python іске қосылуын және дерек іздеуді қамтымайды. `tracemalloc` үрдістің бүкіл жадын көрсетпейді; файл байттарының қосындысы бағдарлама жадына тең емес. Тікелей API төлемінің нөлі осы тәжірибеде ақылы сыртқы сұрау жоқ екенін ғана білдіреді. Құрылғы, электр қуаты, дерек белгілеу және сүйемелдеу де шығын. LLM-нен «мың есе артық» деген тұжырым дәлелденген жоқ: бірдей сұрауларды, сапаны және жағдайды салыстыру керек.

## Тірек сызба

5 жаңа үрдіс → іске қосу ортасы; 200 шақыру → функция уақытының ортасы + жад бөлінуі; файлдар → байттар; API төлемі → тек тікелей төлем.

## Еске түсіріп, тексеріңіз

Кодты жауып, екі уақыт өлшемі неге бөлек екенін түсіндіріңіз. Бағдарламаны екі рет іске қосыңыз; басқа компьютердің санын көшірмеңіз. Үлгі нәтиже: `Сұрау саны: 200` және жаңа іске қосу, ортаңғы уақыт, жад бөлінуі, файл байты мен тікелей API төлемі жолдары. Көмек: жаңа үрдіс модуль мен деректі қайта оқиды. Жиі қате: `peak` мәнін бағдарламаның бүкіл жады деу немесе бір санмен LLM-нен артықшылық жариялау.
