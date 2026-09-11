#!/usr/bin/env python3
"""Check that the code in a lesson speaks the lesson's language.

The course is written in three languages, and for a long time the code blocks
were kept word-for-word identical in all three. That looked tidy and worked
against the reader: the lesson was in their language and the program's output
was not. This checks the rule that replaced it — comments, strings and sample
data belong to the language of the lesson they sit in.

Usage:
    python3 tools/coursecheck/langcheck.py article-go-forms.md ...
    python3 tools/coursecheck/langcheck.py lessons/*.md

The language comes from the file name: NAME.md and NAME-ru.md are Russian,
NAME-kz.md Kazakh, NAME-en.md English. Exit code 1 means something was found.
"""

import os
import re
import sys

KZ_LETTERS = set("әғқңөұүһіӘҒҚҢӨҰҮҺІ")

# Kazakh words spelled with no Kazakh-only letter. Without this list every one
# of them looks Russian to a letter test, and the check cries wolf.
KAZAKH_LOOKALIKES = {
    "дала", "туралы", "деген", "жады", "мин", "саны", "табылмады", "блогым",
    "жоба", "бос", "тым", "жазба", "басталады", "жазда", "сары", "жасыл",
    "болады", "адам", "жазады", "бет", "тек", "блог", "жол", "аты", "бар",
    "жоқ", "бір", "екі", "сан", "мен", "басты", "баяу", "автор", "веб",
    "дене", "алу", "жазу", "осы", "алды", "кез", "келген", "сөз", "мақала",
    "алдық", "не", "музыка", "ішінде",
}

# Of those, the ones that are not also Russian words. A Russian lesson has no
# business printing them, and a letter test alone would never notice.
KAZAKH_ONLY = KAZAKH_LOOKALIKES - {
    "мин", "веб", "музыка", "бар", "блог", "автор", "не", "текст",
}

# Words that mark a string as Russian rather than Kazakh.
RUSSIAN_MARKERS = re.compile(
    r"\b(?:заголовок|заголовки|прочитать|поставить|внешняя|внутренняя|ошибка|"
    r"статья|статей|статьи|слов|букв|добавлено|удалено|блокнот|команды|главная|"
    r"медленная|начало|абзац|строкой|доверенно|память|последние|взломано|"
    r"черновик|изменено|путь|метод|только|нет|номер|номера|нужен|такого|"
    r"ничего|найдено|неизвестная|команда|наберите|часть|текст|привет|мир|"
    r"степи|степь|шанырак|заметка|казахский|значение|тегов|за)\b",
    re.IGNORECASE,
)

# Every fence is matched, and the language decides what is read. Naming only
# some languages in the pattern looked cheaper and was wrong twice over: a
# ```python block was never opened, so the second course's code went unchecked,
# and the scanner then paired that block's closing fence with the next opening
# one -- reading the prose in between as if it were code.
FENCE = re.compile(r"```(\w*)\n(.*?)```", re.S)
CODE_FENCES = ("go", "html", "text", "python", "")
CYRILLIC_WORD = re.compile(r"[А-Яа-яЁёӘҒҚҢӨҰҮҺІәғқңөұүһі]{2,}")

# Deliberate exceptions, documented in docs/go-course.md: a word the lesson is
# about keeps its own language, because replacing it removes the subject.
ALLOWED = {
    ("go-joldar-men-runalar", "ru"): {"шаңырақ", "шаңыр"},
    ("article-go-runes", "ru"): {"шаңырақ", "шаңыр"},
    ("article-go-capstone", "en"): {"шаңырақ"},
}


def language_of(path):
    stem = os.path.basename(path)
    for suffix, lang in ((".kz.md", "kz"), ("-kz.md", "kz"),
                         (".en.md", "en"), ("-en.md", "en")):
        if stem.endswith(suffix):
            return stem[: -len(suffix)], lang
    return stem.removesuffix(".ru.md").removesuffix(".md"), "ru"


# Russian words that keep creeping into the Kazakh lessons because the author
# thinks in Russian while writing them. Each one has a Kazakh word that says the
# same thing, and a reader who meets "выгрузка" in a Kazakh sentence learns that
# the course was translated rather than written.
RUSSISMS_KZ = {
    "выгрузка": "деректер экспорты",
    "выгрузку": "деректер экспортын",
    "выгрузки": "деректер экспорты",
    "олқылық": "жетіспейтін мән",
    "аңғал уақыт": "белдеусіз уақыт",
    "саналы уақыт": "белдеуі бар уақыт",
}

# Words that are Kazakh and still wrong here. Some are book words a reader of
# this course has never met ("бөгде"); others break vowel harmony, which is the
# mistake a Russian-speaking writer makes without hearing it. Matched whole, so
# that "экспортының" is not reported for containing "экспортыны".
WRONG_KZ = {
    "бөгде": "бөтен, басқа",
    "мәндерға": "мәндерге",
    "мәнсыз": "мәнсіз",
    "экспортыны": "экспортын",
    "экспортыдағы": "экспортындағы",
    "экспортылар": "экспорттары",
}

# Second-person singular forms. The lessons address the reader as "сіз"
# throughout, and one "сен" in the middle of a paragraph reads as a different
# author.
INFORMAL_KZ = re.compile(
    r"\b(өзің|сенің|сені|саған|қатеңе|атауың|аласың|алмайсың|көресің|жазасың|"
    r"істейсің|қоясың|білесің|аларсың)\b", re.IGNORECASE)


def check_kazakh_prose(path):
    """Reports Russian loanwords, wrong words and informal address in Kazakh.

    Both are about the same thing: a Kazakh page that reads as a translation is
    a Kazakh page the reader will switch away from.
    """
    if not path.endswith("-kz.md"):
        return []
    with open(path, encoding="utf-8") as f:
        text = f.read()
    prose = "".join(FENCE.sub("", text))
    found = []
    for word, better in RUSSISMS_KZ.items():
        if word in prose.lower():
            found.append(f"русизм «{word}» — лучше «{better}»")
    for word, better in WRONG_KZ.items():
        if re.search(r"\b" + word + r"\b", prose, re.IGNORECASE):
            found.append(f"неудачное слово «{word}» — лучше «{better}»")
    for m in INFORMAL_KZ.finditer(prose):
        found.append(f"обращение на «сен»: «{m.group(0)}»")
        break
    return found


def check(path):
    name, lang = language_of(path)
    with open(path, encoding="utf-8") as f:
        code = "".join(body for lang, body in FENCE.findall(f.read())
                       if lang in CODE_FENCES)
    allowed = {w.lower() for w in ALLOWED.get((name, lang), set())}

    if lang == "ru":
        found = {
            w.lower() for w in CYRILLIC_WORD.findall(code)
            if set(w) & KZ_LETTERS or w.lower() in KAZAKH_ONLY
        }
        found -= allowed
        return sorted(found), "казахские слова в русском уроке"

    found = set()
    for word in CYRILLIC_WORD.findall(code):
        low = word.lower()
        if set(word) & KZ_LETTERS or low in KAZAKH_LOOKALIKES or low in allowed:
            continue
        if RUSSIAN_MARKERS.fullmatch(low) or RUSSIAN_MARKERS.search(low):
            found.add(low)
    return sorted(found), f"русские слова в уроке на «{lang}»"


def main(paths):
    if not paths:
        print(__doc__)
        return 2
    bad = 0
    for path in paths:
        words, why = check(path)
        if words:
            bad += 1
            print(f"{path}: {why}: {' '.join(words)}")
        for note in check_kazakh_prose(path):
            bad += 1
            print(f"{path}: {note}")
    if bad:
        print(f"\nзамечаний: {bad}")
        return 1
    print(f"проверено файлов: {len(paths)} — язык кода везде совпадает с языком урока")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
