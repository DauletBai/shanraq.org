# Web Assistant 1.1: local classroom release

This continues version 1.0. `tasks.json` has the same three **fictional** tasks. The server displays them in a browser but never changes them. It uses only Python's standard library. Do not enter real student names, passwords, grades, or messages.

## Run

In this directory, run `python3 web_assistant.py --port 8765`, then open `http://127.0.0.1:8765/?lang=en&today=2026-10-09`. Change to `lang=ru` or `lang=kz` for the other interfaces. The form changes the viewed date. Press Ctrl+C in the terminal to stop the server.

Expected statuses: `t-01 REMIND`, `t-02 DONE`, `t-03 NO_DATE`. `/?today=2026-02-29` returns 400, `/missing` returns 404, and POST returns 405. Check with `python3 -m unittest discover -s tests -v`. Compare `shasum -a 256 tasks.json` before and after.

The server listens only on 127.0.0.1, uses HTTP without TLS, and is not intended for the public internet. Python documents `http.server` as unsuitable for production. A future cloud deployment would need a production server, HTTPS, accounts, access control, protected storage, and a separate deployment process. None of these is claimed here.

The checkpoint preserves the previous `assistant_core.py` rules, original records, web view, tests, and these instructions. Read-only behavior is deliberate: security and database lessons come later.
