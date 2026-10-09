"""Acceptance cases from paper release 0.3 plus new file and CLI boundaries."""

from datetime import date, timedelta
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

from assistant_core import (add_task, complete_task, count_done, load_document,
                            parse_date, reminder_status, save_document,
                            to_version_one, validate_document)


HERE = Path(__file__).resolve().parents[1]
PAPER = json.loads((HERE.parent / "step-02/cases.json").read_text(encoding="utf-8"))
ORIGINAL = HERE.parent / "step-01/data/tasks.json"


class AssistantTests(unittest.TestCase):
    def test_paper_reminder_cases_are_preserved(self):
        today = date(2026, 10, 9)
        for case in PAPER["reminder_cases"]:
            days = case["days_to_due"]
            task = {"done": case["done"], "due_date": None if days is None
                    else (today + timedelta(days=days)).isoformat()}
            self.assertEqual(reminder_status(task, today), case["expected"], case)

    def test_paper_count_and_duplicate_cases_are_preserved(self):
        for case in PAPER["count_cases"]:
            tasks = [{"done": value} for value in case["done_values"]]
            self.assertEqual(count_done(tasks), case["expected"])
        for case in PAPER["duplicate_cases"]:
            document = {"version": "1.0", "tasks": [
                {"id": task_id, "title": f"Task {n}", "done": False}
                for n, task_id in enumerate(case["ids"])]}
            if case["has_duplicate"]:
                with self.assertRaisesRegex(ValueError, "duplicate id"):
                    validate_document(document)
            else:
                self.assertEqual(len(validate_document(document)["tasks"]), len(case["ids"]))

    def test_migration_preserves_three_records_and_does_not_invent_dates(self):
        old = load_document(ORIGINAL)
        self.assertEqual(old["version"], "0.2")
        new = to_version_one(old)
        self.assertEqual(new["version"], "1.0")
        self.assertEqual([task["id"] for task in new["tasks"]], ["t-01", "t-02", "t-03"])
        self.assertEqual([task["done"] for task in new["tasks"]], [False, True, False])
        self.assertTrue(all(task["due_date"] is None for task in new["tasks"]))

    def test_date_validation_catches_shape_and_impossible_date(self):
        self.assertEqual(parse_date("2026-10-09"), date(2026, 10, 9))
        for bad in ("20261009", "2026-02-29", "2026-13-01", 20261009, ""):
            with self.assertRaises(ValueError, msg=bad):
                parse_date(bad)

    def test_invalid_file_does_not_replace_valid_data(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "tasks.json"
            path.write_text("{broken", encoding="utf-8")
            with self.assertRaisesRegex(ValueError, "invalid JSON"):
                load_document(path)
            self.assertEqual(path.read_text(encoding="utf-8"), "{broken")
            with self.assertRaisesRegex(ValueError, "duplicate id"):
                save_document(path, {"version": "1.0", "tasks": [
                    {"id": "x", "title": "A", "done": False},
                    {"id": "x", "title": "B", "done": False}]})
            self.assertEqual(path.read_text(encoding="utf-8"), "{broken")

    def test_unicode_round_trip_and_repeat_completion(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "tasks.json"
            clean = add_task({"version": "1.0", "tasks": []}, "t-01", "Әліппе оқу", "2026-10-11")
            save_document(path, clean)
            self.assertIn("Әліппе оқу", path.read_text(encoding="utf-8"))
            loaded = load_document(path)
            self.assertEqual(loaded["tasks"][0]["title"], "Әліппе оқу")
            done_once = complete_task(loaded, "t-01")
            done_twice = complete_task(done_once, "t-01")
            self.assertEqual(done_once, done_twice)
            self.assertEqual(reminder_status(done_twice["tasks"][0], date(2026, 10, 9)), "DONE")

    def test_cli_uses_fixed_today_and_preserves_old_file_until_write(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "tasks.json"
            original = ORIGINAL.read_text(encoding="utf-8")
            path.write_text(original, encoding="utf-8")
            command = [sys.executable, str(HERE / "assistant.py"), "--file", str(path)]
            shown = subprocess.run(command + ["list"], capture_output=True, text=True, check=True)
            self.assertIn("Completed: 1/3", shown.stdout)
            self.assertEqual(path.read_text(encoding="utf-8"), original)
            added = subprocess.run(command + ["add", "t-04", "Жоба жазу", "--due", "2026-10-11"],
                                   capture_output=True, text=True, check=True)
            self.assertIn("Added: t-04", added.stdout)
            self.assertEqual(load_document(path)["version"], "1.0")
            output = subprocess.run(command + ["reminders", "--today", "2026-10-09"],
                                    capture_output=True, text=True, check=True).stdout
            self.assertIn("t-01: NO_DATE", output)
            self.assertIn("t-02: DONE", output)
            self.assertIn("t-04: REMIND", output)
            invalid = subprocess.run(command + ["add", "t-04", "Again"],
                                     capture_output=True, text=True)
            self.assertEqual(invalid.returncode, 2)
            self.assertIn("duplicate id", invalid.stderr)
            self.assertEqual(len(load_document(path)["tasks"]), 4)


if __name__ == "__main__":
    unittest.main()
