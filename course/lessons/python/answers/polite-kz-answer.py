"""The answer to lesson 20's exercise, Kazakh: one polite fetch.

Politeness here is three concrete things and no goodwill at all: a timeout, so
the program stops waiting; retries only on the codes that can get better, with
the pause doubling, so a struggling server is not hit harder; and a cache, so
the same question is not asked twice. The count at the end is the whole point --
four readings, two requests.
"""

import json
import threading
import time
import urllib.error
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

HITS = {"rate": 0}
RETRY = (429, 500, 502, 503, 504)
TTL = 60
CACHE = {}


class Teacher(BaseHTTPRequestHandler):
    """Оқу сервері: бірінші рет бас тартады, содан кейін жауап береді."""

    def do_GET(self):
        HITS["rate"] += 1
        if HITS["rate"] == 1:
            self.answer(503, b"try later")
        else:
            self.answer(200, b'{"rate": 510.43}')

    def answer(self, code, body):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *args):
        return


server = ThreadingHTTPServer(("127.0.0.1", 0), Teacher)
threading.Thread(target=server.serve_forever, daemon=True).start()
URL = f"http://127.0.0.1:{server.server_address[1]}/rate"


def fetch(url, timeout=1.0):
    """Таймауты және өз қолтаңбасы бар бір сұраныс."""
    request = urllib.request.Request(url, headers={"User-Agent": "shanraq-course/1.0"})
    with urllib.request.urlopen(request, timeout=timeout) as answer:
        return json.load(answer)


def polite(url, tries=3):
    """Кэштен не желіден алынған жауап — қайталаумен және шегініспен."""
    now = time.monotonic()
    if url in CACHE and now - CACHE[url][0] < TTL:
        return CACHE[url][1], "кэштен"

    delay = 0.1
    for attempt in range(1, tries + 1):
        try:
            data = fetch(url)
        except urllib.error.HTTPError as error:
            if error.code not in RETRY:
                raise
            print(f"  {attempt}-талпыныс: {error.code}, {delay:.1f} с күтеміз")
            time.sleep(delay)
            delay *= 2
            continue
        except TimeoutError:
            print(f"  {attempt}-талпыныс: үнсіздік, {delay:.1f} с күтеміз")
            time.sleep(delay)
            delay *= 2
            continue
        CACHE[url] = (now, data)
        return data, "желіден"
    raise RuntimeError(f"дереккөз {tries} талпыныста жауап бермеді")


for number in range(1, 5):
    data, where = polite(URL)
    print(f"{number}: {data['rate']} — {where}")
print("серверге сұраныс:", HITS["rate"])

server.shutdown()
