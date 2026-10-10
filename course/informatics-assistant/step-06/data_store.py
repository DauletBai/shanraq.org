"""SQLite checkpoint for fictional classroom tasks and study sessions.

Only the explicit ``init`` command creates a database. Normal page requests
open an existing database read-only. The 1.0 reminder rules stay unchanged.
"""

import csv
from contextlib import closing
import os
from pathlib import Path
import sqlite3
import tempfile

from assistant_core import load_document, parse_date, validate_document


SCHEMA = """
CREATE TABLE tasks (
  task_id TEXT PRIMARY KEY,
  title TEXT NOT NULL CHECK (length(trim(title)) > 0),
  done INTEGER NOT NULL CHECK (done IN (0, 1)),
  due_date TEXT
);
CREATE TABLE study_sessions (
  session_id TEXT PRIMARY KEY,
  task_id TEXT NOT NULL REFERENCES tasks(task_id),
  observed_on TEXT NOT NULL,
  minutes INTEGER NOT NULL CHECK (minutes BETWEEN 0 AND 180),
  source TEXT NOT NULL CHECK (source IN ('paper', 'app'))
);
CREATE INDEX study_sessions_task_idx ON study_sessions(task_id);
"""
HEADERS = ["session_id", "task_id", "observed_on", "minutes", "source"]


def read_sessions(path, known_task_ids):
    """Validate the entire CSV before creating the persistent database."""
    with Path(path).open(newline="", encoding="utf-8-sig") as stream:
        reader = csv.DictReader(stream)
        if reader.fieldnames != HEADERS:
            raise ValueError("study CSV has incorrect columns")
        seen, rows = set(), []
        for number, row in enumerate(reader, 2):
            if None in row or any(value is None for value in row.values()):
                raise ValueError(f"CSV row {number} has extra or missing cells")
            session_id = row["session_id"].strip()
            task_id = row["task_id"].strip()
            date_string = row["observed_on"].strip()
            minutes_string = row["minutes"].strip()
            source = row["source"].strip()
            if not session_id or session_id in seen:
                raise ValueError(f"CSV row {number} has empty or duplicate session ID")
            seen.add(session_id)
            if task_id not in known_task_ids:
                raise ValueError(f"CSV row {number} references an unknown task")
            parse_date(date_string)
            if not minutes_string.isascii() or not minutes_string.isdecimal():
                raise ValueError(f"CSV row {number} has invalid minutes")
            minutes = int(minutes_string)
            if not 0 <= minutes <= 180:
                raise ValueError(f"CSV row {number} exceeds the classroom limit")
            if source not in {"paper", "app"}:
                raise ValueError(f"CSV row {number} has an unknown source")
            rows.append((session_id, task_id, date_string, minutes, source))
    return rows


def create_database(database_path, tasks_path, sessions_path):
    """Build a new database beside the target, then publish it atomically."""
    destination = Path(database_path)
    if destination.exists():
        raise ValueError("database already exists; keep your earlier checkpoint")
    document = load_document(tasks_path)
    if document["version"] != "1.0":
        raise ValueError("step-05 needs the 1.0 task format")
    observations = read_sessions(sessions_path, {task["id"] for task in document["tasks"]})
    handle, temporary_name = tempfile.mkstemp(prefix=".assistant-", suffix=".db", dir=destination.parent)
    os.close(handle)
    temporary = Path(temporary_name)
    try:
        with closing(sqlite3.connect(temporary)) as connection, connection:
            connection.execute("PRAGMA foreign_keys = ON")
            connection.executescript(SCHEMA)
            connection.executemany("INSERT INTO tasks VALUES (?, ?, ?, ?)", [
                (task["id"], task["title"], int(task["done"]), task["due_date"])
                for task in document["tasks"]])
            connection.executemany("INSERT INTO study_sessions VALUES (?, ?, ?, ?, ?)", observations)
            if connection.execute("PRAGMA foreign_key_check").fetchall():
                raise ValueError("database has broken references")
        try:
            os.link(temporary, destination)
        except FileExistsError as exc:
            raise ValueError("database already exists; keep your earlier checkpoint") from exc
    finally:
        temporary.unlink(missing_ok=True)


def connect_readonly(database_path):
    path = Path(database_path).resolve()
    if not path.is_file():
        raise ValueError("database missing; run the init command first")
    connection = sqlite3.connect(path.as_uri() + "?mode=ro", uri=True)
    connection.row_factory = sqlite3.Row
    connection.execute("PRAGMA foreign_keys = ON")
    return connection


def load_document_from_db(database_path):
    with closing(connect_readonly(database_path)) as connection:
        rows = connection.execute("SELECT task_id, title, done, due_date FROM tasks ORDER BY task_id").fetchall()
    return validate_document({"version": "1.0", "tasks": [
        {"id": row["task_id"], "title": row["title"], "done": bool(row["done"]),
         "due_date": row["due_date"]} for row in rows]})


def minutes_report(database_path):
    """LEFT JOIN keeps a task even when it has no study sessions."""
    with closing(connect_readonly(database_path)) as connection:
        return [tuple(row) for row in connection.execute("""
            SELECT t.task_id, COALESCE(SUM(s.minutes), 0) AS minutes
            FROM tasks AS t
            LEFT JOIN study_sessions AS s ON s.task_id = t.task_id
            GROUP BY t.task_id
            ORDER BY t.task_id
        """)]
