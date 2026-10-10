# Assistant 1.2: data in SQLite

This is a classroom checkpoint with three fictional tasks and four fictional study sessions. Do not enter real names, marks, messages or passwords. Python 3 from the previous block is enough.

From this folder run `python3 data_assistant.py init`. It reads and validates `tasks.json` and `study_sessions.csv`, then creates local `assistant.db`. A failed import leaves no completed DB file. Repeating `init` does not replace an existing database. Do not commit `assistant.db` to Git.

Run `python3 data_assistant.py report`: expect `t-01 45`, `t-02 0`, `t-03 25`. Zero beside `t-02` means no recorded sessions, not a measured zero-minute session. Run `python3 web_assistant.py --port 8765` and open `http://127.0.0.1:8765/?lang=en&today=2026-10-09`. Earlier statuses survive: `REMIND`, `DONE`, `NO_DATE`. For other languages use `lang=ru` or `lang=kz`. Stop the server with Ctrl+C.

Test with `python3 -m unittest discover -s tests -v`. The database and browser view are **local only**. Classroom `http.server` is not for the public internet and has no user accounts or access control.
