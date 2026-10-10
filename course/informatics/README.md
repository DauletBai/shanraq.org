# Informatics and AI course

The course is a dependency graph built around one product rather than a list of
unrelated software exercises. `curriculum.json` is the source of truth for its
72-lesson route and its eight atomic publication blocks.

The learner develops **My Digital Assistant** from a paper model to a tested
Python program, a web interface, a SQLite database, a protected application,
and a small measured AI feature. Every release must preserve the behaviour
proved by earlier tests. The product uses fictional learning data; learners do
not put their real names, marks, messages, contacts, or credentials into the
repository.

Every released lesson has hand-written Russian, Kazakh, and English versions.
The versions must keep the same facts, examples, expected outputs, project
release, and assessment. Russian files have no language suffix; Kazakh and
English files use `-kz` and `-en`.

The lesson cycle is:

1. a situation from the learner's digital life;
2. a question or contradiction;
3. a concrete model;
4. the precise term and the boundary of the analogy;
5. prediction before observation;
6. a worked example;
7. a fading example;
8. an independent task;
9. error diagnosis;
10. transfer to an unseen setting;
11. retrieval after 1, 7, and 30 days;
12. one tested change to the cross-course project.

Support maps are explanatory SVG assets. Decorative 3D treatment may clarify
layers and grouping, but cannot change the exact sequence, values, labels, or
relationships stated in the lesson.

Publication is block based. A block is committed and deployed only when all of
its lessons, three language versions, support maps, project checkpoint, and
release checks are complete. The blocks contain 10, 8, 8, 12, 8, 8, 8, and 10
lessons respectively. No individual lesson is published from an unfinished
block.

Run the structural release checks from the repository root:

```sh
python3 -m unittest tools.course.test_informatics_release
```

Release 0.2 adds lessons 11–18 about bits, binary numbers, Unicode and UTF-8,
pixels, sampled sound and video, compression, integrity, and a documented data
format. Its project checkpoint is in `course/informatics-assistant/step-01/`.

Release 0.3 adds lessons 19–26 about decomposition, state, tracing, conditions,
loops, functions, correctness, and efficiency. It keeps the three fictional
records in the 0.2 `tasks.json` unchanged. Due dates are separate test inputs:
the earlier file has no due-date field. The paper algorithms, trilingual
contracts, and independently checkable cases are in
`course/informatics-assistant/step-02/`. Python implementation follows in
the next block. The cover is
`/static/covers/school/informatics/algorithms/03-algorithmic-thinking.webp`.

Release 1.0 adds lessons 27–38. Learners turn the paper reminder contract
into a Python program, keep the three fictional 0.2 records when migrating to
1.0, validate ISO dates and UTF-8 JSON, diagnose errors, and run repeatable
tests. The checked command-line project and localized quickstarts are in
`course/informatics-assistant/step-03/`. The 4K 16:9 cover is
`/static/covers/school/informatics/python/04-python-assistant.webp`.

Release 1.1 adds lessons 39–46 about a message's network journey, IP,
packets, DNS, TCP and UDP, TLS, HTTP, semantic HTML, CSS, accessibility,
and an honest cloud model. Each lesson has manually localized Russian, Kazakh,
and English text and an exact support map. `step-04` is a real **local,
read-only** browser view of the same three fictional tasks and reminder rules.
It binds only to 127.0.0.1; this is not a public cloud deployment or a
production web server. The 4K cover is
`/static/covers/school/informatics/web/05-internet-web-cloud.webp`.

To prepare only this new block without rewriting the published 1–38 records:

```sh
python3 tools/course/prepare_informatics.py --start 39 \
  --sql /tmp/informatics-web-release.sql \
  --expected /tmp/informatics-web-expected.json
```

For a full rebuild, prepare all forty-six lessons as one reviewable transaction:

```sh
python3 tools/course/prepare_informatics.py \
  --sql /tmp/informatics-course.sql \
  --expected /tmp/informatics-course-expected.json
```

The preparation command does not contact production. Publication applies the
single transaction only after the source commit, application migration, static
assets, and course workflow have passed their checks.
