# `pytest` тесттері: конвейер уәделерін тексеру

_Лид (summary):_ **Python курсының елу екінші сабағы. Жоба ережелерін автоматты тесттерге айналдырамыз: қалыпты жағдайды, шекараны және күтілетін қатені тексеріп, бүкіл жинақты бір пәрменмен іске қосамыз.**

## Бұл не үшін керек

Қолмен іске қосу бағдарламаның **бір рет** жұмыс істегенін ғана көрсетеді. Тест жоба уәдесін сақтап, әр өзгерістен кейін қайта тексереді. Мысалы, белгісіз дереккөз қабылданбауы, себеп туралы тұжырым редакторды күтуі, ал бос есеп бұрынғы файлды қайта жазбауы керек.

Тест барлық қатенің жоқтығын дәлелдемейді. Ол тек жазылған мысалды тексереді. Сондықтан жақсы жинақ код жолдарының санын емес, мінез-құлықтың маңызды шекараларын қамтиды.

## Бірден толығымен

Екі файл жасаңыз:

```python
# policy.py
def route(claim_type, evidence_ids, known_ids, formula_id=None):
    if claim_type not in {"observation", "calculation", "cause"}:
        return "reject"
    if not set(evidence_ids).issubset(known_ids):
        return "reject"
    if claim_type == "cause":
        return "review"
    if claim_type == "calculation" and not formula_id:
        return "review"
    return "accept"
```

```python
# test_policy.py
import pytest

from policy import route


@pytest.mark.parametrize(
    ("claim_type", "evidence", "formula_id", "expected"),
    [
        ("observation", ["cpi-2025"], None, "accept"),
        ("cause", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], None, "review"),
        ("calculation", ["cpi-2025"], "change-v1", "accept"),
        ("guess", ["cpi-2025"], None, "reject"),
        ("observation", ["missing"], None, "reject"),
    ],
)
def test_route(claim_type, evidence, formula_id, expected):
    assert route(claim_type, evidence, {"cpi-2025"}, formula_id) == expected
```

Әзірлеу құралын орнатып, тексеруді іске қосыңыз:

```console
python -m pip install pytest
python -m pytest -q
```

`6 passed` — алты жағдай өтті деген сөз. `python -m pytest` пәрмені дәл сол `python` интерпретаторын қолданады.

## Тесттің құрылысы

`test_policy.py` файлы мен `test_route` функциясын `pytest` өзі табады. Ішінде Arrange–Act–Assert реті бар:

1. **Arrange** — кіріс деректерін дайындаңыз;
2. **Act** — тексерілетін бір әрекетті орындаңыз;
3. **Assert** — нәтижені уәдемен салыстырыңыз.

Бір тестке бірнеше тәуелсіз себепті қоспаңыз: ол құласа, қай ереже бұзылғаны түсініксіз болады.

## Көшірменің орнына жағдайлар кестесі

`@pytest.mark.parametrize` бір функцияны кестенің әр жолымен іске қосады. Жаңа ереже қосқанда, алдымен құлауға тиіс мысалды жазыңыз, кейін жұмыс кодын өзгертіп, тесттерді қайта жүргізіңіз. Бұл — «қызыл → жасыл → жақсарту» айналымы.

Барлық ықтимал комбинациядан алып кесте жасамаңыз. Мінез-құлық сыныптарын таңдаңыз: қалыпты жағдай, бос мән, шекара, қате түр және бұрын табылған ақау.

## Күтілетін қате

Кейде қате — дұрыс нәтиже:

```python
import pytest

def percent(value):
    if not 0 <= value <= 100:
        raise ValueError("пайыз ауқымнан тыс")
    return value

def test_percent_rejects_101():
    with pytest.raises(ValueError, match="ауқымнан тыс"):
        percent(101)
```

Ерекшелік шықпаса немесе түрі басқа болса, тест құлайды. Кең `pytest.raises(Exception)` қолданбаңыз: ол бөтен ақауды дұрыс мінез-құлық деп қабылдауы мүмкін.

## Қоқыссыз файлдар

`tmp_path` кірістірілген фикстурасы әр тестке бөлек уақытша қалта береді:

```python
def test_write_report(tmp_path):
    report = tmp_path / "report.txt"
    report.write_text("дайын\n", encoding="utf-8")
    assert report.read_text(encoding="utf-8") == "дайын\n"
```

Тест ағымдағы қалтаға тәуелді емес және нақты есепті бұзбайды. Желіні, ағымдағы уақытты және сыртқы API-ды модульдік тестте берілген дерекпен не алмастырушы функциямен оқшаулаған дұрыс. Нақты шекараны бөлек интеграциялық тест тексере алады, бірақ ол баяу әрі тұрақсыздау.

## Құрылымды емес, мінез-құлықты тексеріңіз

Пайдалы тест байқалатын нәтижені тексереді: шешім, файл, есеп жолы немесе ерекшелік. Ішкі код өзгерсе де, уәде сақталса, тест өзгермеуі керек.

Табылған әр ақауға регрессиялық тест жазыңыз: алдымен қатені қайталап, тесттің ескі кодта құлағанын көріңіз, себепті түзетіңіз де, тестті сақтап қойыңыз.

## Сабақ картасы

![Сабақ картасы: уәде, мысалдар және нәтиже](/static/course/py/map-tests-kz.svg)

Тірек белгі: **уәде → қалыпты жағдай + шекара + қате → бір іске қосу → түсінікті нәтиже**.

## Өз сөзіңізбен айтыңыз

1. Неліктен бір сәтті қолмен іске қосу тестті алмастырмайды?
2. `@pytest.mark.parametrize` не береді?
3. `tmp_path` не үшін керек?
4. `pytest.raises(Exception)` неге тым кең?

## Жаттығу

**1. Болжаңыз.** Неше тест жағдайы орындалады?

<!-- drill 1 -->
```python
import pytest

values = [0, 50, 100]

@pytest.mark.parametrize("value", values)
def test_range(value):
    assert 0 <= value <= 100

print(len(values))
```

**2. Бос орынды толтырыңыз.**

```python
import pytest

with pytest.___(ValueError):
    percent(101)
```

**3. Түзетіңіз.** Тест нақты есепке жазып тұр:

```python
def test_report():
    path = Path("data/report.txt")
    path.write_text("test", encoding="utf-8")
```

## Тапсырма

**Міндетті.** `policy.py` және `test_policy.py` жасаңыз. Алғашқы мысалдағы алты жағдайды тексеріңіз. Бөлек тестте пайыз тексергішінің `-1` және `101` мәндерін `ValueError` арқылы қабылдамайтынын растаңыз. `python -m pytest -q` іске қосыңыз.

<!-- task out -->
```text
8 passed
```

**Өз дерегіңізбен.** Конвейеріңіздің үш уәдесін таңдаңыз. Әрқайсысына қалыпты жағдай, шекара және қате кіріс жазыңыз.

**Қосымша.** HTML жазуды `tmp_path` арқылы тексеріп, пайдаланушы мәтінінің экрандалатынын растаңыз.

## Бұл жобаның қай жеріне кіреді

Жобаның 27-қадамы модель ережелері мен бетке арналған регрессиялық тесттерді алады. Енді өзгерісті бір пәрменмен тексеруге болады. Келесі сабақ тип аннотациялары мен `mypy` арқылы кейбір сәйкессіздікті тестке дейін табады.

## Жауаптар

1. Қолмен іске қосу кіріс дерегін, күтілетін нәтижені және қайталанатын тексеруді сақтамайды.
2. Бір ережені бірнеше анық көрсетілген жағдайда іске қосады.
3. Тест файлдарын оқшаулап, жоба дерегін қозғамау үшін.
4. Кез келген бөтен ерекшелік күтілген қате деп қабылдануы мүмкін.

<!-- drill 1 out -->
```text
3
```

2. Әдіс — `raises`.

<!-- drill 2 -->
```python
import pytest

def percent(value):
    if not 0 <= value <= 100:
        raise ValueError("пайыз ауқымнан тыс")

with pytest.raises(ValueError):
    percent(101)
print("ValueError")
```

<!-- drill 2 out -->
```text
ValueError
```

3. Файлды `tmp_path` ішінде жасаңыз.

<!-- drill 3 -->
```python
def test_report(tmp_path):
    path = tmp_path / "report.txt"
    path.write_text("test", encoding="utf-8")
    assert path.read_text(encoding="utf-8") == "test"
```

<!-- drill 3 out -->
```text
```

## Дереккөздер

- [pytest: Get Started](https://docs.pytest.org/en/stable/getting-started.html) — тесттерді табу және іске қосу.
- [pytest: parametrization](https://docs.pytest.org/en/stable/how-to/parametrize.html) — тест жағдайларының кестелері.
- [pytest: temporary directories](https://docs.pytest.org/en/stable/how-to/tmp_path.html) — оқшауланған уақытша файлдар.
