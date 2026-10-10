"""Read-only, loopback-only web view of the fictional assistant 1.0 data.

Classroom demonstration: Python's http.server is not a production web server.
Run: python3 web_assistant.py --port 8765
"""

import argparse
from datetime import date
from html import escape
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import parse_qs, urlsplit

from assistant_core import count_done, load_document, parse_date, reminder_status


TEXT = {
    "ru": ("Мой цифровой помощник", "Задачи", "Сегодня", "Показать", "Выполнено", "Состояние", "Учебные вымышленные данные"),
    "kz": ("Менің цифрлық көмекшім", "Тапсырмалар", "Бүгін", "Көрсету", "Орындалды", "Күйі", "Ойдан алынған оқу деректері"),
    "en": ("My Digital Assistant", "Tasks", "Today", "Show", "Completed", "Status", "Fictional learning data"),
}


def render_page(document, today, lang="ru"):
    """Build escaped, accessible HTML without changing the saved document."""
    if lang not in TEXT:
        raise ValueError("unsupported language")
    title, tasks_label, today_label, show, completed, status, notice = TEXT[lang]
    rows = []
    for task in document["tasks"]:
        rows.append("<tr><th scope='row'>" + escape(task["id"]) + "</th><td>" +
                    escape(task["title"]) + "</td><td>" +
                    escape(reminder_status(task, today)) + "</td></tr>")
    total = len(document["tasks"])
    rows_html = "".join(rows) or f"<tr><td colspan='3'>{escape(tasks_label)}: 0</td></tr>"
    return ("<!doctype html><html lang='" + lang + "'><head><meta charset='utf-8'>" +
            "<meta name='viewport' content='width=device-width,initial-scale=1'>" +
            "<title>" + escape(title) + "</title><style>" +
            "body{font:1.1rem/1.5 system-ui;max-width:56rem;margin:2rem auto;padding:0 1rem;color:#152338;background:#f5f8fc}" +
            "main{background:white;padding:1.5rem;border-radius:1rem;box-shadow:0 8px 24px #1a365822}" +
            "table{border-collapse:collapse;width:100%}th,td{border:1px solid #64748b;padding:.7rem;text-align:left}" +
            "th{background:#eaf2ff}label{display:block;margin:.8rem 0}button{padding:.5rem 1rem}" +
            "</style></head><body><main><h1>" + escape(title) + "</h1><p>" + escape(notice) +
            "</p><form method='get' action='/'><input type='hidden' name='lang' value='" + lang +
            "'><label for='today'>" + escape(today_label) +
            "</label><input id='today' name='today' type='date' required value='" +
            today.isoformat() + "'><button type='submit'>" + escape(show) +
            "</button></form><p>" + escape(completed) + ": " +
            str(count_done(document["tasks"])) + "/" + str(total) +
            "</p><h2>" + escape(tasks_label) + "</h2><table><thead><tr>" +
            "<th scope='col'>ID</th><th scope='col'>" + escape(tasks_label) +
            "</th><th scope='col'>" + escape(status) +
            "</th></tr></thead><tbody>" + rows_html + "</tbody></table></main></body></html>")


def make_handler(data_file):
    """Fix the data file at server creation; URL paths never select local files."""
    class Handler(BaseHTTPRequestHandler):
        def send_body(self, status, body, content_type):
            payload = body.encode("utf-8")
            self.send_response(status)
            self.send_header("Content-Type", content_type + "; charset=utf-8")
            self.send_header("Content-Length", str(len(payload)))
            self.send_header("Cache-Control", "no-store")
            self.send_header("X-Content-Type-Options", "nosniff")
            self.send_header("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'")
            self.end_headers()
            self.wfile.write(payload)

        def do_GET(self):
            url = urlsplit(self.path)
            if url.path == "/health":
                self.send_body(200, "OK\n", "text/plain")
                return
            if url.path != "/":
                self.send_body(404, "Not found\n", "text/plain")
                return
            query = parse_qs(url.query)
            try:
                lang = query.get("lang", ["ru"])[0]
                today = parse_date(query["today"][0]) if "today" in query else date.today()
                document = load_document(data_file)
                body = render_page(document, today, lang)
            except (ValueError, IndexError):
                self.send_body(400, "Invalid date, language, or classroom data\n", "text/plain")
                return
            self.send_body(200, body, "text/html")

        def do_POST(self):
            self.send_body(405, "Read-only classroom version\n", "text/plain")

    return Handler


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--port", type=int, default=8765)
    parser.add_argument("--file", type=Path, default=Path(__file__).with_name("tasks.json"))
    args = parser.parse_args()
    if not 0 <= args.port <= 65535:
        parser.error("port must be between 0 and 65535")
    server = ThreadingHTTPServer(("127.0.0.1", args.port), make_handler(args.file))
    print(f"Classroom view: http://127.0.0.1:{server.server_port}/")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        server.server_close()


if __name__ == "__main__":
    main()
