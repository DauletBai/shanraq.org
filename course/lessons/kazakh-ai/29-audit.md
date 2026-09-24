# Урок 29. Проверяем локальность и объясняем ответ

## Зачем это нужно

На продуктовой этикетке написано, откуда товар, когда проверен и до какого дня годен. На нашей карточке факта нужны те же сведения. Для каждого ответа запишем дорожку решений: что распознали, сколько карточек нашли и почему разрешили ответ.

## Перед кодом

Скопируйте вместе `engine.py`, `checks.py`, `examples.json` и `facts.json`. `from checks import classify` подключает функцию урока 27: при импорте её учебные примеры не печатаются благодаря `__main__`. `answer(question, on_date)` получает текст вопроса и календарную дату, возвращает словарь `text`/`trace`. `trace` — список коротких шагов, `str(...)` превращает число в текст, знак `+` сцепляет строки. `on_date.isoformat()` пишет дату как `год-месяц-день`. Сначала проверяется тип вопроса, затем одна карточка, дата проверки и срок действия. Учебный факт о шахматах вымышлен.

Файлы этого шага: `examples.json`, `checks.py`, `engine.py`, `facts.json`.

Из корня проекта перейдите в папку шага и запустите программу:

```text
cd course/kazakh-ai/step-29
python3 engine.py
```

В Windows замените `python3` на `py`.

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

[Файлы шага](https://github.com/DauletBai/shanraq.org/tree/main/course/kazakh-ai/step-29).

## Что делает программа

Для `Шахмат қашан?` на учебную дату 24 сентября 2026 года видны ключ, одна карточка, источник и две даты. `Сурет қайда?` не находит карточку; опечатка `Шахмаат` останавливается до поиска. Если день раньше `checked_on` или позже `valid_until`, ответ не выдаётся. Знак `<` означает «раньше/меньше», а `>` — «позже/больше». Программа читает локальные JSON-файлы и не содержит вызова сетевого API; это проверка видимого пути кода, не доказательство изоляции всей операционной системы. След можно проверить вручную, но неверная исходная карточка всё равно породит неверный ответ.

## Опорная карта

Вопрос → тип → ключ → 0/1/много карточек → две даты → ответ с источником + след или отказ.

## Вспомните и проверьте

Закройте код и восстановите порядок проверок. Эталон: шахматы дают ответ и пять шагов следа; у рисования нет факта, поэтому `білмеймін: дерек жоқ`; опечатка даёт `білмеймін: таныс емес сөз`. Подсказка: два ранних `return` завершают функцию до работы с датой. Поменяйте `valid_until` на `2026-09-01` и получите `білмеймін: бұл күнге жарамды дерек жоқ`. Частая ошибка: считать след доказательством истины исходного листа.
