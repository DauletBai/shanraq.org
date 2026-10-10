# Assistant 2.1: human-controlled model advice

This continues version 1.3. First follow the database and backup steps below. Also open `labelled_tasks.csv`: it has 12 fictional training rows and four separate held-out rows. Run `python3 classifier.py "Кітап оқу"`: expect `suggestion: reading`; the holdout has two correct reading and two correct math cases. Try `python3 classifier.py "Unknown task"`: expect `review`. This is a suggestion, not a task edit; check the decision yourself. A tiny 4/4 sample does not establish quality for real learners. Tests: `python3 -m unittest discover -s tests -v`. The DB and local server remain classroom tools and the backup is unencrypted.

## Earlier safety and recovery steps

This continues version 1.2. Three tasks and four sessions are fictional. Use Python 3; no external packages. Enter no real personal data, passwords or classmates' files. Download the repository and open a terminal in this directory.

1. `python3 data_assistant.py init` creates a new `assistant.db` only. Repeating it does not replace an existing DB.
2. `python3 data_assistant.py report` prints `t-01 45`, `t-02 0`, `t-03 25`. The zero for `t-02` means no recorded sessions.
3. `python3 security_assistant.py verify assistant.db` prints `tasks=3 sessions=4` after SQLite checks.
4. `python3 security_assistant.py backup assistant.db backups/first.db` creates a new copy and `first.db.sha256` checksum.
5. `python3 security_assistant.py restore backups/first.db restored.db` checks the copy and creates a new file, refusing to overwrite. Compare `python3 security_assistant.py verify restored.db` and the result of `minutes_report("restored.db")` in Python.
6. `python3 web_assistant.py --port 8765` opens a local page on `127.0.0.1` only. Stop with Ctrl+C.
7. `python3 -m unittest discover -s tests -v` checks older rules and recovery.

The backup is not encrypted. Keep it on a separate trusted device, not in a shared folder or Git. A checksum next to the copy detects change but cannot prove authenticity if an attacker can replace both files. The local `http.server` has no accounts and is not for publication on the internet.
