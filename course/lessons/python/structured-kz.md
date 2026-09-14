# Құрылымды жауап: модель JSON-ы және оны тексеру

_Лид (summary):_ **Python курсының қырық тоғызыншы сабағы. Модельден еркін мәтін емес, берілген схемаға сай JSON сұраймыз, оны стандартты `json` модулімен талдап, `jsonschema` кітапханасымен тексереміз. Синтаксис, құрылым және факт деп бөлінген үш шекара әдемі жазылған қатені есепке өткізбейді.**

## Бұл не үшін керек

48-сабақта нақты сұранысты сақтадық. Бірақ қайталанатын сұраныстың жауабын бағдарлама пайдалана алмауы мүмкін: модель өрісті басқаша атады, санды жол түрінде жазды немесе JSON сыртына түсініктеме қосты.

«JSON қайтар» деген сөз — тілек. **JSON Schema** — код тексере алатын келісім. Ол міндетті өрістерді, олардың түрін, рұқсат етілген мәндерін және артық өріске тыйымды анықтайды.

Тексерудің үш сатысы бар:

1. `json.loads`: бұл шын мәнінде JSON ба?
2. `jsonschema`: оның құрылымы күткендей ме?
3. жоба ережесі: сандар бастапқы деректен алынған ба?

Бүгін алғашқы екі сатыны құрамыз. Үшіншісін келесі сабақта қосамыз.

## Бірден толығымен

Бекітілген нұсқаны орнатыңыз:

```bash
python -m pip install jsonschema==4.26.0
```

`structured.py` файлы сақталған жауаппен жұмыс істейді, сондықтан іске қосуға Ollama қажет емес.

```python
"""49-сабақ: құрылымды жауапты талдау және тексеру."""

import json

from jsonschema import Draft202012Validator

OUTPUT_SCHEMA = {
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
        "country": {"type": "string", "minLength": 1},
        "year": {"type": "integer", "minimum": 2000, "maximum": 2100},
        "value": {"type": "number"},
        "unit": {"type": "string", "enum": ["percent"]},
    },
    "required": ["country", "year", "value", "unit"],
    "additionalProperties": False,
}


def parse_and_validate(raw):
    """Келісімге сай JSON-ды ғана қайтару."""
    payload = json.loads(raw)
    Draft202012Validator(OUTPUT_SCHEMA).validate(payload)
    return payload


raw = '{"country":"Қазақстан","year":2025,"value":11.4,"unit":"percent"}'
result = parse_and_validate(raw)

print("ел:", result["country"])
print("жыл — бүтін сан:", isinstance(result["year"], int))
print("мән — сан:", isinstance(result["value"], (int, float)))
print("бірлік:", result["unit"])
print("артық өріс жоқ:", set(result) == {"country", "year", "value", "unit"})
```

Шығыс:

```text
ел: Қазақстан
жыл — бүтін сан: True
мән — сан: True
бірлік: percent
артық өріс жоқ: True
```

## Алдымен талдау, содан кейін сенім

`json.loads` JSON жолын Python мәндеріне айналдырады. Үтір қате тұрса немесе ашылатын жақшаның алдында мәтін болса, `JSONDecodeError` туады. Фигуралы жақша арасын тұрақты өрнекпен қиып алмаңыз: жауаптың бір сынығын толық құжат деп қабылдауыңыз мүмкін.

API-ден таза құрылымды нәтиже талап еткен дұрыс. Ollama сұранысында схеманы `format` өрісімен беруге болады:

```python
request["format"] = OUTPUT_SCHEMA
```

Бұл қате жауап санын азайтады, бірақ өз бағдарламаңыздағы тексеруді алмастырмайды. Модель мен желі сенім шекарасының сыртында қалады.

## Схема нені уәде етеді

`type: object` нысан талап етеді. `required` міндетті өрістерді атайды. `additionalProperties: False` жарияланбаған өрістерге тыйым салады: `county` деген қате атау дұрыс өрістің қасына жасырынбайды.

JSON саны `11.4` пен JSON жолы `"11.4"` — екі бөлек мән. Схема бірін екіншісіне айналдырмайды және модельдің ойын таппауы керек.

`enum` өлшем бірлігін `percent` сөзімен шектейді. Жыл аралығы анық қоқысты ұстайды, бірақ 2025 жылдың бастапқы кестеде барын дәлелдемейді.

Схеманың өзін де бір рет тексеріңіз:

```python
Draft202012Validator.check_schema(OUTPUT_SCHEMA)
```

## Қате — нәтиженің бір түрі

Бағдарлама шекарасында тек күтілетін қателерді ұстаңыз және бос сөздікпен жұмысты жалғастырмаңыз:

```python
from json import JSONDecodeError
from jsonschema import ValidationError

try:
    result = parse_and_validate(raw)
except JSONDecodeError:
    print("Жауап JSON емес")
except ValidationError:
    print("JSON схемаға сай емес")
else:
    save_for_review(result)
```

Қате туралы журналға пайдаланушының толық мәтінін қоспай жазуға болады. Сақтауға болмайтын дерек жоқ болса, қате жауаптың өзін жабық жерде талдау үшін қалдырыңыз.

## Схема нені білмейді

Схема `value` өрісінің сан екенін растайды. Бірақ `11.4` бастапқы кестеде болды ма, ел шатасты ма — оны білмейді. **Дұрыс құрылым фактінің дұрыстығын білдірмейді.**

Схемадан өткен нәтижені SQL, HTML немесе қабық пәрменінде бірден пайдаланбаңыз. Экрандау, параметрлі сұраныс және пәндік шектеу әлі де қажет.

## Сабақ картасы

![Сабақ картасы: жол, JSON, схема және тексерілген өрістер](/static/course/py/map-structured-kz.svg)

Опорлық сигналды қалпына келтіріңіз: **шикі жауап → `json.loads` → JSON Schema → фактіні тексеруге арналған дерек**.

## Өз сөзіңізбен айтыңыз

1. «JSON қайтар» деген нұсқау неге тексеруді алмастырмайды?
2. `JSONDecodeError` мен `ValidationError` айырмасы қандай?
3. `required` пен `additionalProperties: False` неге қатар керек?
4. Схемадан өткен санды неге бірден жариялауға болмайды?

## Жаттығу

**1. Болжаңыз.** Бұл мән `number` түрінің тексеруінен өте ме?

<!-- drill 1 -->
```python
import json

value = json.loads('{"value": "11.4"}')["value"]
print(isinstance(value, (int, float)))
```

**2. Бос орынды толтырыңыз.** Келісімде жоқ өріске тыйым салыңыз.

```python
schema = {"type": "object", ...: False}
```

**3. Жөндеңіз.** Құрылым қатесі қазір жарамды жауапқа айналып кетеді.

```python
try:
    result = parse_and_validate(raw)
except Exception:
    result = {}
save_for_review(result)
```

## Тапсырма

**Міндетті.** `check_many(raw_items)` функциясын жазыңыз. Әр жол үшін үш нәтиженің бірін санаңыз: `valid`, `invalid_json`, `invalid_schema`. Тек `JSONDecodeError` мен `ValidationError` ұстаңыз; қате жауап қабылданғандар тізіміне кірмесін.

<!-- task out -->
```text
қабылданды: 1
JSON емес: 1
схемаға сай емес: 2
жылдар: [2025]
```

**Өз дерегіңізбен.** Жиынтықтағы бір кесте нәтижесінің схемасын жазыңыз. Бірлікті `enum` арқылы атаңыз, артық өріске тыйым салыңыз және әр шекараны тексеретін төрт жауап құрыңыз.

**Қалауыңызша.** Осы схеманы Ollama сұранысының `format` өрісіне беріңіз. Схемамен және онсыз он рет орындап, жергілікті тексеруден өткен жауап үлесін салыстырыңыз.

## Бұл жобаға қайда кіреді

Модель жауабы енді құрылымы тексерілген нәтижелер тізіміне ғана түседі. 50-сабақта оған жиынтық сандарын беріп, қайтарған әр санды дереккөзбен салыстырамыз.

## Жауаптар

1. Модель түсініктеме қосуы, өрісті қалдыруы не санды жол етуі мүмкін; нұсқау код тексеруі емес.
2. Біріншісі — JSON синтаксисі қате, екіншісі — JSON дұрыс, бірақ құрылымы қате.
3. Бірі жоқ өрісті, екіншісі қате аталған не жарияланбаған өрісті табады.
4. Схема түр мен аралықты біледі, бірақ санның дереккөзі мен растығын білмейді.

<!-- drill 1 out -->
```text
False
```

2. Кілт — `"additionalProperties"`.

<!-- drill 2 -->
```python
schema = {"type": "object", "additionalProperties": False}
print(schema["additionalProperties"])
```

<!-- drill 2 out -->
```text
False
```

3. Күтілетін қатені ұстап, бас тартуды белгілеңіз және келесі кезеңді шақырмаңыз.

<!-- drill 3 -->
```python
from jsonschema import ValidationError, validate

raw = {}

try:
    validate(raw, {"type": "object", "required": ["year"]})
except ValidationError:
    print("қабылданбады")
else:
    print("қабылданды")
```

<!-- drill 3 out -->
```text
қабылданбады
```

## Дереккөздер

- [JSON Schema: алғашқы қадам](https://json-schema.org/learn/getting-started-step-by-step) — нысан, міндетті және қосымша қасиеттер.
- [`jsonschema` құжаттамасы](https://python-jsonschema.readthedocs.io/en/stable/validate/) — валидаторлар мен схеманы тексеру.
- [Ollama structured outputs](https://docs.ollama.com/capabilities/structured-outputs) — JSON Schema-ны `format` өрісіне беру.
- [Python `json`](https://docs.python.org/3/library/json.html) — JSON талдау және `JSONDecodeError`.
