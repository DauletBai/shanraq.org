"""The small, testable core of My Digital Assistant 1.0.

Uses only Python's standard library and fictional classroom data.
"""

from datetime import date
import json
import os
from pathlib import Path
import re
import tempfile


DATE_PATTERN = re.compile(r"\d{4}-\d{2}-\d{2}\Z")
VERSIONS = {"0.2", "1.0"}


def parse_date(value):
    """Accept exactly YYYY-MM-DD and return a date, or explain the error."""
    if not isinstance(value, str) or not DATE_PATTERN.fullmatch(value):
        raise ValueError("date must have YYYY-MM-DD format")
    try:
        return date.fromisoformat(value)
    except ValueError as exc:
        raise ValueError("date does not exist") from exc


def validate_document(document):
    """Return a fresh validated document; never trust or mutate input data."""
    if not isinstance(document, dict) or document.get("version") not in VERSIONS:
        raise ValueError("unsupported document version")
    if set(document) != {"version", "tasks"}:
        raise ValueError("document has unknown or missing fields")
    records = document.get("tasks")
    if not isinstance(records, list):
        raise ValueError("tasks must be a list")
    seen = set()
    clean = []
    for position, task in enumerate(records, 1):
        if not isinstance(task, dict):
            raise ValueError(f"task {position} must be an object")
        allowed = {"id", "title", "done"}
        if document["version"] == "1.0":
            allowed.add("due_date")
        if set(task) - allowed:
            raise ValueError(f"task {position} has unknown fields")
        task_id = task.get("id")
        title = task.get("title")
        done = task.get("done")
        if not isinstance(task_id, str) or not task_id.strip():
            raise ValueError(f"task {position} needs a nonempty id")
        if task_id in seen:
            raise ValueError(f"duplicate id: {task_id}")
        seen.add(task_id)
        if not isinstance(title, str) or not title.strip():
            raise ValueError(f"task {position} needs a nonempty title")
        if not isinstance(done, bool):
            raise ValueError(f"task {position} needs a boolean done value")
        if document["version"] == "0.2" and "due_date" in task:
            raise ValueError("version 0.2 cannot store due_date")
        due = task.get("due_date")
        if due is not None:
            parse_date(due)
        checked = {"id": task_id, "title": title, "done": done}
        if document["version"] == "1.0":
            checked["due_date"] = due
        clean.append(checked)
    return {"version": document["version"], "tasks": clean}


def count_done(tasks):
    """Return the number of completed tasks; an empty list counts as zero."""
    count = 0
    for task in tasks:
        if task["done"]:
            count += 1
    return count


def reminder_status(task, today):
    """Implement the five outcomes of the 0.3 paper contract."""
    if task["done"]:
        return "DONE"
    if task.get("due_date") is None:
        return "NO_DATE"
    days = (parse_date(task["due_date"]) - today).days
    if days < 0:
        return "OVERDUE"
    if days <= 2:
        return "REMIND"
    return "NOT_YET"


def load_document(path):
    """Read UTF-8 JSON and reject malformed or partial data before display."""
    try:
        with Path(path).open(encoding="utf-8") as stream:
            document = json.load(stream)
    except OSError as exc:
        raise ValueError(f"cannot read file: {exc.strerror}") from exc
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid JSON at line {exc.lineno}") from exc
    return validate_document(document)


def save_document(path, document):
    """Write a validated 1.0 document via a temporary file in one directory."""
    clean = validate_document(document)
    if clean["version"] != "1.0":
        raise ValueError("save requires version 1.0")
    destination = Path(path)
    temporary = None
    try:
        with tempfile.NamedTemporaryFile("w", encoding="utf-8", dir=destination.parent,
                                         prefix=".tasks-", suffix=".tmp", delete=False) as stream:
            temporary = Path(stream.name)
            json.dump(clean, stream, ensure_ascii=False, indent=2)
            stream.write("\n")
            stream.flush()
            os.fsync(stream.fileno())
        os.replace(temporary, destination)
    except OSError as exc:
        raise ValueError(f"cannot save file: {exc.strerror}") from exc
    finally:
        if temporary is not None:
            temporary.unlink(missing_ok=True)


def to_version_one(document):
    """Explicitly migrate 0.2 records without inventing dates."""
    clean = validate_document(document)
    return {"version": "1.0", "tasks": [
        {**task, "due_date": task.get("due_date")} for task in clean["tasks"]]}


def add_task(document, task_id, title, due_date=None):
    """Return a new 1.0 document with one extra task."""
    clean = to_version_one(document)
    added = {"id": task_id, "title": title, "done": False, "due_date": due_date}
    clean["tasks"].append(added)
    return validate_document(clean)


def complete_task(document, task_id):
    """Mark one task done; repeating the command has the same result."""
    clean = to_version_one(document)
    for task in clean["tasks"]:
        if task["id"] == task_id:
            task["done"] = True
            return clean
    raise ValueError(f"unknown id: {task_id}")
