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

Prepare the first ten lessons as one reviewable database transaction:

```sh
python3 tools/course/prepare_informatics.py \
  --sql /tmp/informatics-course.sql \
  --expected /tmp/informatics-course-expected.json
```

The preparation command does not contact production. Publication applies the
single transaction only after the source commit, application migration, static
assets, and course workflow have passed their checks.
