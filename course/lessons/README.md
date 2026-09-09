# Course lessons — the sources

Every lesson lives here as Markdown: a folder per course, and one file per
language inside it.

    go/http.md         Russian
    go/http-kz.md      Kazakh
    go/http-en.md      English
    python/prices.md   and the same three for the Python course

A folder per course rather than one flat pile: two courses make 150 files, and
SQL, Rust and machine learning are already named on the site as what comes
next. The language suffix is what the checks read, so it stays whatever the
folder is called.

The site is the published copy, not the source. Until this folder existed the
only copies were the production database and a scratch directory, so nothing
could be reviewed in a diff, nothing was checked by CI, and a lesson edited on
the site quietly diverged from the text it was written from. Both mistakes had
already happened: the lesson on packages was rewritten in the database and the
local copy stayed a version behind, and a Kazakh word was fixed in ten lessons
with no record of which ten.

## The shape of a file

    # Title
    <blank>
    _Лид (summary):_ **One paragraph that becomes the article's summary.**
    <blank>
    ## First section
    ...

The first line is the headline, the third is the summary the card and the search
result show, and everything from the fifth line on is the body. `tools/course/publish.py`
relies on exactly this shape.

The headline here may be shorter than the one on the site: the site's titles
carry a subtitle after a colon, and they are edited there, not here.

## The shape of an exercise

A lesson ends with a ladder, and the steps are what a reader climbs alone:

    ## Разминка / Жаттығу / Warm-up
    1. Предскажите   a program to answer before running it
    2. Заполните     one gap to fill in
    3. Почините      one broken program to mend
    ## Задание / Тапсырма / Exercise
    **Обязательное.**   the task on fixed data, with the output it must print
    **На своих данных.** the same task on the reader's own numbers
    **По желанию.**     what is worth trying and is not required

The exercise heading is the one the site reads to put the check box on the page
(`lessonExercise` in `pkg/modules/articles/course_check.go`), so it stays
exactly one of those three, and "На своих данных" comes after the required part
rather than instead of it.

Everything the ladder promises is run. A drill's program and its answer live
apart -- the reader is meant to answer first -- so they are tied by markers,
which are HTML comments and never reach the page:

    <!-- drill 1 -->        the program that is meant to run
    <!-- drill 1 out -->    the output it is meant to print
    <!-- task out -->       the output the required task promises

The required task's own solution lives in `course/lessons/python/answers/` as
`<lesson file>-answer.py` -- one per language, because the printed result is in
the lesson's language, and that is exactly where a translation drifts. The
`-answer` suffix is not decoration: a file called `csv.py` beside a program is
what `import csv` finds.

## The rules the checks enforce

- **Code speaks the lesson's language.** Comments, strings and sample data in a
  Kazakh lesson are Kazakh; in an English lesson, English. `langcheck.py`
  fails otherwise. Repository code is the opposite — English comments — and the
  two rules do not meet, because a lesson's program is read by the student.
- **Every program parses.** `gocheck.py` and `pycheck.py` hand each fenced block
  that starts a whole program to `gofmt -e` or to the Python compiler.
- **Every Python program prints what the page shows.** `pyrun.py` runs it and
  compares the output with the block printed under it, because compiling says
  nothing about what a program does. `--steps` does the same for the course
  project in `course/py-digest`.
- **Every answer is the one the machine printed.** The same `pyrun.py` runs the
  warm-up drills through their markers and the exercise's reference solution
  against the output the lesson promises. A promise nobody checks is the part
  that rots: the reader types towards it and gets something else.
- **Links exist.** `linkcheck.py` compares links into this repository against the
  real remote and fetches the rest.
- **Nothing is silently missing.** `pysyllabus.py` holds every element of Python
  to the lesson that promised it — the answer to the Go course, where arrays and
  constants were not postponed but forgotten.

- **The checks themselves are checked.** `selftest.py` hands each checker a
  small file with a known fault and fails when it is not reported. Two of these
  tools were silently broken before somebody read their source: langcheck's
  fence pattern never opened a ```python block, and pysyllabus took an equals
  sign as proof that a lesson explains functions.

Run them all at once:

    python3 tools/coursecheck/selftest.py
    python3 tools/coursecheck/langcheck.py course/lessons/*/*.md
    python3 tools/coursecheck/gocheck.py   course/lessons/go/*.md
    python3 tools/coursecheck/pycheck.py   course/lessons/python/*.md
    python3 tools/coursecheck/pyrun.py     course/lessons/python/*.md
    python3 tools/coursecheck/pyrun.py     --steps
    python3 tools/coursecheck/linkcheck.py --offline course/lessons/*/*.md
    python3 tools/coursecheck/pysyllabus.py plan
    python3 tools/coursecheck/pysyllabus.py published

CI runs the same set on any change under `course/lessons/` or `tools/coursecheck/`.

A new course adds a folder here and its own line to the checks that are
language-specific; nothing else moves. The code a course builds keeps its own
place — `course/go-blog` for the blog, `course/py-digest` for the digest —
because lessons link into it by URL and those links are already published.

## Publishing

`tools/course/publish.py` compares these files against the live site and updates
what differs. It never invents an article: a lesson must already exist, and the
mapping from file to slug is `tools/course/lesson-slugs.json`.

    python3 tools/course/publish.py --check       # what differs from the site
    python3 tools/course/publish.py               # send the differences
    python3 tools/course/publish.py python/prices.md   # one lesson

Numbers measured in a lesson are numbers measured, not numbers imagined: when a
figure changes, the program is run again and the output block is replaced with
what it printed. Where a measurement depends on the language — the length of a
Kazakh message, the BM25 rank of a Kazakh text — the three versions honestly
differ, and the studio's translation check says so. Where it does not — a
timestamp, an ephemeral port, a random token — the three versions carry the same
run.
