# My Digital Assistant 1.0: executable checkpoint

This is the runnable result of informatics lessons 27–38. It uses Python's
standard library and three fictional tasks. No account, network, package
installation, or personal data is needed.

Choose a language: [Русский](README-ru.md) · [Қазақша](README-kz.md) ·
[English](README-en.md).

From this directory:

```sh
python3 assistant.py list
python3 assistant.py reminders --today 2026-10-09
python3 -m unittest discover -s tests -v
```

`assistant_core.py` contains testable rules; `assistant.py` is the command-line
interface. `tasks.json` is a version 1.0 sample. The original version 0.2 file
in `../step-01/data/tasks.json` stays unchanged. `add` or `done` explicitly
migrates a supplied 0.2 file to 1.0 while preserving all three tasks; it does
not invent due dates. `--today` makes reminder examples reproducible; without
it, the program uses the computer's local calendar date.
