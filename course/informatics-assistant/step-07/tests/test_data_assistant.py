"""Release 1.2 behaviour: migration, constraints, reports and web continuity."""

from datetime import date
from http.server import ThreadingHTTPServer
from pathlib import Path
import sqlite3
import sys
from tempfile import TemporaryDirectory
from threading import Thread
import unittest
from urllib.request import urlopen

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from assistant_core import load_document, reminder_status  # noqa: E402
from data_store import create_database, load_document_from_db, minutes_report, read_sessions  # noqa: E402
from web_assistant import make_handler  # noqa: E402


ROOT = Path(__file__).resolve().parents[1]


class DataAssistantTests(unittest.TestCase):
    def test_atomic_migration_preserves_tasks_and_earlier_reminders(self):
        with TemporaryDirectory() as directory:
            db = Path(directory) / "assistant.db"
            create_database(db, ROOT / "tasks.json", ROOT / "study_sessions.csv")
            self.assertEqual(load_document_from_db(db), load_document(ROOT / "tasks.json"))
            self.assertEqual(minutes_report(db), [("t-01", 45), ("t-02", 0), ("t-03", 25)])
            statuses = [reminder_status(task, date(2026, 10, 9))
                        for task in load_document_from_db(db)["tasks"]]
            self.assertEqual(statuses, ["REMIND", "DONE", "NO_DATE"])
            with self.assertRaisesRegex(ValueError, "already exists"):
                create_database(db, ROOT / "tasks.json", ROOT / "study_sessions.csv")
            with sqlite3.connect(db) as connection:
                connection.execute("PRAGMA foreign_keys = ON")
                with self.assertRaises(sqlite3.IntegrityError):
                    connection.execute("INSERT INTO study_sessions VALUES ('bad', 'unknown', '2026-10-09', 5, 'paper')")

    def test_bad_rows_do_not_create_database_and_web_is_read_only(self):
        with TemporaryDirectory() as directory:
            db = Path(directory) / "assistant.db"
            bad = Path(directory) / "bad.csv"
            bad.write_text("session_id,task_id,observed_on,minutes,source\n"
                           "s-01,t-01,2026-10-09,20,paper\n"
                           "s-01,t-03,2026-10-09,10,app\n", encoding="utf-8")
            with self.assertRaisesRegex(ValueError, "duplicate"):
                create_database(db, ROOT / "tasks.json", bad)
            self.assertFalse(db.exists())
            self.assertEqual(len(read_sessions(ROOT / "study_sessions.csv", {"t-01", "t-02", "t-03"})), 4)
            create_database(db, ROOT / "tasks.json", ROOT / "study_sessions.csv")
            before = db.read_bytes()
            server = ThreadingHTTPServer(("127.0.0.1", 0), make_handler(db))
            thread = Thread(target=server.serve_forever, daemon=True)
            thread.start()
            try:
                with urlopen(f"http://127.0.0.1:{server.server_port}/?lang=kz&today=2026-10-09") as response:
                    page = response.read().decode()
                    self.assertEqual(response.status, 200)
                    self.assertIn("REMIND", page)
                    self.assertIn("NO_DATE", page)
                self.assertEqual(db.read_bytes(), before)
            finally:
                server.shutdown()
                server.server_close()
                thread.join(timeout=3)


if __name__ == "__main__":
    unittest.main()
