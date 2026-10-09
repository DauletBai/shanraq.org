# Digital Assistant 1.0

This is the working result of lessons 27–38. Use Python 3.12 or newer; no
extra packages or network are needed. Use fictional learning tasks only.

From this directory run:

```sh
python3 assistant.py list
python3 assistant.py reminders --today 2026-10-09
python3 -m unittest discover -s tests -v
```

Expect `Completed: 1/3`, then `t-01: REMIND`, `t-02: DONE`, and
`t-03: NO_DATE`. `--today` fixes the date so the experiment can be repeated
on any day. Without it, the program uses your computer's local calendar date.

To keep the example intact, copy `tasks.json` into another directory and put
its path after `--file` **before** the command:

```sh
python3 assistant.py --file my-tasks.json add t-04 "Жоба жазу" --due 2026-10-11
python3 assistant.py --file my-tasks.json done t-04
python3 assistant.py --file my-tasks.json list
```

The previous block's 0.2 file can be read unchanged. Only `add` or `done`
writes a 1.0 document, and neither invents a date for an old task. The new
format stores a date as `YYYY-MM-DD` or `null`. Invalid JSON, impossible dates,
and duplicate IDs are rejected before writing. `assistant_core.py` holds the
rules, `assistant.py` handles commands, and `tests` checks every promise of
paper version 0.3.
