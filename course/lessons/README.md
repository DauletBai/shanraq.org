# Course lessons — the sources

Every lesson of the courses lives here as Markdown, one file per language:

    article-go-http.md      Russian
    article-go-http-kz.md   Kazakh
    article-go-http-en.md   English

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

## The rules the checks enforce

- **Code speaks the lesson's language.** Comments, strings and sample data in a
  Kazakh lesson are Kazakh; in an English lesson, English. `langcheck.py`
  fails otherwise. Repository code is the opposite — English comments — and the
  two rules do not meet, because a lesson's program is read by the student.
- **Every program parses.** `gocheck.py` and `pycheck.py` hand each fenced block
  that starts a whole program to `gofmt -e` or to the Python compiler.
- **Links exist.** `linkcheck.py` compares links into this repository against the
  real remote and fetches the rest.
- **Nothing is silently missing.** `pysyllabus.py` holds every element of Python
  to the lesson that promised it — the answer to the Go course, where arrays and
  constants were not postponed but forgotten.

Run them all at once:

    python3 tools/coursecheck/langcheck.py course/lessons/*.md
    python3 tools/coursecheck/gocheck.py   course/lessons/article-go-*.md
    python3 tools/coursecheck/pycheck.py   course/lessons/article-py-*.md
    python3 tools/coursecheck/linkcheck.py --offline course/lessons/*.md
    python3 tools/coursecheck/pysyllabus.py plan

CI runs the same set on any change under `course/lessons/` or `tools/coursecheck/`.

## Publishing

`tools/course/publish.py` compares these files against the live site and updates
what differs. It never invents an article: a lesson must already exist, and the
mapping from file to slug is `tools/course/lesson-slugs.json`.

    python3 tools/course/publish.py --check      # what differs from the site
    python3 tools/course/publish.py              # send the differences

Numbers measured in a lesson are numbers measured, not numbers imagined: when a
figure changes, the program is run again and the output block is replaced with
what it printed. Where a measurement depends on the language — the length of a
Kazakh message, the BM25 rank of a Kazakh text — the three versions honestly
differ, and the studio's translation check says so. Where it does not — a
timestamp, an ephemeral port, a random token — the three versions carry the same
run.
