# 29-сабақ. Жергілікті жұмысты тексеріп, жауапты түсіндіреміз

## Бұл не үшін керек?

Азық-түлік жапсырмасында оның қайдан келгені, қашан тексерілгені және қай күнге дейін жарамды екені жазылады. Нақты дерек карточкасына да осындай мәлімет керек. Әр жауапқа шешім жолын қосамыз: не танылды, қанша карточка табылды және неге жауап берілді.

## Кодқа дейін

`engine.py`, `checks.py`, `examples.json`, `facts.json` файлдарын бірге көшіріңіз. `from checks import classify` 27-сабақтағы функцияны қосады: `__main__` тексерісінің арқасында импорт кезінде оқу мысалдары басылмайды. `answer(question, on_date)` сұрақ мәтіні мен күнді қабылдап, `text`/`trace` өрістері бар сөздік қайтарады. `trace` — қысқа қадамдар тізімі, `str(...)` санды мәтінге айналдырады, `+` мәтіндерді біріктіреді. `on_date.isoformat()` күнді `жыл-ай-күн` түрінде жазады. Әуелі сұрақ түрі, содан соң бір карточка, тексерілген күн және жарамдылық мерзімі қаралады. Шахмат туралы дерек — ойдан алынған.

Осы қадамның файлдары: `examples.json`, `checks.py`, `engine.py`, `facts.json`.

Жоба түбірінен қадам қалтасына өтіп, бағдарламаны іске қосыңыз:

```text
cd course/kazakh-ai/step-29
python3 engine.py
```

Windows жүйесінде `python3` орнына `py` жазыңыз.

```python
import json
from datetime import date

from checks import classify


def answer(question, on_date):
    trace = []
    request = classify(question)
    if request["status"] != "ready":
        trace.append("Сұрақ түрі анықталмады")
        return {"text": request["reason"], "trace": trace}
    club = request["club"]
    kind = request["kind"]
    trace.append("Кілт: " + club + " / " + kind)
    with open("facts.json", encoding="utf-8") as file:
        facts = json.load(file)
    matches = []
    for fact in facts:
        if fact["club"] == club and fact["kind"] == kind:
            matches.append(fact)
    trace.append("Жазба саны: " + str(len(matches)))
    if len(matches) == 0:
        return {"text": "білмеймін: дерек жоқ", "trace": trace}
    if len(matches) > 1:
        return {"text": "тоқта: бірнеше дерек табылды", "trace": trace}
    fact = matches[0]
    checked = date.fromisoformat(fact["checked_on"])
    until = date.fromisoformat(fact["valid_until"])
    trace.append("Дереккөз: " + fact["source"])
    trace.append("Тексерілген: " + fact["checked_on"])
    trace.append("Жарамды: " + fact["valid_until"] + " дейін")
    if on_date < checked or on_date > until:
        return {"text": "білмеймін: бұл күнге жарамды дерек жоқ", "trace": trace}
    if kind == "уақыт":
        label = "уақыты"
    else:
        label = "орны"
    message = (fact["club"] + " үйірмесінің " + label + ": " + fact["value"]
               + " | " + on_date.isoformat() + " күнгі дерек"
               + " | дереккөз: " + fact["source"])
    return {"text": message, "trace": trace}


if __name__ == "__main__":
    day = date.fromisoformat("2026-09-24")
    for question in ["Шахмат қашан?", "Сурет қайда?", "Шахмаат қашан?"]:
        result = answer(question, day)
        print(question, "→", result["text"])
        print("Із:", result["trace"])
```

[Қадам файлдары](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-29).

## Бағдарлама қалай жұмыс істейді?

`Шахмат қашан?` сұрағы үшін 2026 жылғы 24 қыркүйекте кілт, жалғыз карточка, дереккөз және екі күн көрінеді. `Сурет қайда?` карточка таппайды; `Шахмаат` қатесі іздеуге жетпей тоқтайды. Күн `checked_on` мәнінен бұрын не `valid_until` мәнінен кейін болса, жауап берілмейді. `<` — «бұрын/кіші», `>` — «кейін/үлкен». Бағдарлама жергілікті JSON файлдарын оқиды, кодында желілік API шақыруы жоқ; бұл көрінетін код жолының тексерісі, бүкіл операциялық жүйе оқшаулығының дәлелі емес. Ізді қолмен тексеруге болады, бірақ бастапқы карточка қате болса, жауап та қате болуы мүмкін.

## Тірек сызба

Сұрақ → түр → кілт → 0/1/көп карточка → екі күн → дереккөзі мен ізі бар жауап не бас тарту.

![29-сабақтың тірек сызбасы](/static/course/kazakh-ai/map-29-audit-kz.svg)

## Еске түсіріп, тексеріңіз

Кодты жауып, тексерулер ретін айтыңыз. Үлгі нәтиже: шахматқа жауап пен іздің бес қадамы шығады; сурет үйірмесі туралы карточка жоқ, сондықтан `білмеймін: дерек жоқ`; қате атауға `білмеймін: таныс емес сөз` шығады. Көмек: екі ерте `return` күнді салыстыруға жеткізбейді. `valid_until` мәнін `2026-09-01` деп өзгертіп, `білмеймін: бұл күнге жарамды дерек жоқ` нәтижесін алыңыз. Жиі қате: ізді бастапқы парақтың шындығына дәлел деп санау.
