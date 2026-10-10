"""Behavioral checks for the classroom HTTP checkpoint."""

import json
from datetime import date
from http.server import ThreadingHTTPServer
from pathlib import Path
import sys
from tempfile import TemporaryDirectory
from threading import Thread
import unittest
from urllib.error import HTTPError
from urllib.request import Request, urlopen

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from assistant_core import load_document, reminder_status  # noqa: E402
from web_assistant import make_handler, render_page  # noqa: E402


ROOT = Path(__file__).resolve().parents[1]


class WebAssistantTests(unittest.TestCase):
    def test_web_view_keeps_existing_project_contract_in_all_languages(self):
        document = load_document(ROOT / "tasks.json")
        self.assertEqual([task["id"] for task in document["tasks"]], ["t-01", "t-02", "t-03"])
        for lang in ("ru", "kz", "en"):
            page = render_page(document, date(2026, 10, 9), lang)
            self.assertIn(f"lang='{lang}'", page)
            self.assertIn("1/3", page)
            for task in document["tasks"]:
                self.assertIn(reminder_status(task, date(2026, 10, 9)), page)
            self.assertIn("<label for='today'>", page)
            self.assertIn("<th scope='col'>", page)

    def test_escapes_titles_and_does_not_change_fictional_file(self):
        with TemporaryDirectory() as temporary:
            data_file = Path(temporary) / "tasks.json"
            document = load_document(ROOT / "tasks.json")
            document["tasks"][0]["title"] = "<script>alert(1)</script>"
            data_file.write_text(json.dumps(document, ensure_ascii=False), encoding="utf-8")
            before = data_file.read_bytes()
            server = ThreadingHTTPServer(("127.0.0.1", 0), make_handler(data_file))
            thread = Thread(target=server.serve_forever, daemon=True)
            thread.start()
            base = f"http://127.0.0.1:{server.server_port}"
            try:
                with urlopen(base + "/?lang=kz&today=2026-10-09") as response:
                    page = response.read().decode()
                    self.assertEqual(response.status, 200)
                    self.assertIn("&lt;script&gt;alert(1)&lt;/script&gt;", page)
                    self.assertIn("t-01", page)
                    self.assertIn("REMIND", page)
                    self.assertIn("default-src 'none'", response.headers["Content-Security-Policy"])
                with urlopen(base + "/health") as response:
                    self.assertEqual(response.read(), b"OK\n")
                for path in ("/?today=2026-02-29", "/?lang=de", "/tasks.json"):
                    with self.assertRaises(HTTPError) as caught:
                        urlopen(base + path)
                    caught.exception.close()
                with self.assertRaises(HTTPError) as caught:
                    urlopen(Request(base + "/", data=b"done=t-01", method="POST"))
                self.assertEqual(caught.exception.code, 405)
                caught.exception.close()
                self.assertEqual(data_file.read_bytes(), before)
            finally:
                server.shutdown()
                server.server_close()
                thread.join(timeout=3)


if __name__ == "__main__":
    unittest.main()
