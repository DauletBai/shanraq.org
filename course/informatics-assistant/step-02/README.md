# My Digital Assistant — checkpoint 0.3

This release is a **paper algorithm**, not yet an executable reminder service.
The three fictional records in `../step-01/data/tasks.json` remain unchanged.
That file has no due-date field. The due-day values in `cases.json` are separate
test cards supplied to the algorithm, not hidden changes to the 0.2 format.

Read the specification in [Russian](ALGORITHM-ru.md),
[Kazakh](ALGORITHM-kz.md), or [English](ALGORITHM-en.md). Each describes the
same five reminder outcomes, loop, duplicate-ID rule, and acceptance cases.

To verify the fixture without installing anything:

```sh
python3 -m json.tool course/informatics-assistant/step-02/cases.json >/dev/null
python3 -m unittest tools.course.test_informatics_release
```

The next block translates this contract into Python. It must preserve every
accepted case, including the missing date, the two-day boundary, empty input,
and a repeated ID. No real pupil names, marks, contacts, or credentials belong
in this project.
