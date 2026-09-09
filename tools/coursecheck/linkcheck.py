#!/usr/bin/env python3
"""Check every link a lesson hands the reader.

A dead link in a lesson is worse than a missing one: the reader trusts it,
follows it, and lands on a 404 that looks like the course lying. This came up
for real -- two lessons pointed at github.com/dake-edu/shanraq, an address that
never existed, while the project lives at DauletBai/shanraq.org.

So two checks. Links into our own repository are compared against the address
`git remote get-url origin` actually reports, which catches an invented one
without touching the network. Everything else is fetched, and anything that
does not answer 2xx or 3xx is named.

    python3 tools/coursecheck/linkcheck.py <файлы>...
    python3 tools/coursecheck/linkcheck.py --offline <файлы>...   # только свои ссылки
"""

import concurrent.futures
import json
import os
import re
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request

LINK = re.compile(r"\]\((https?://[^)\s]+)\)")
# Links inside the site: a lesson that sends the reader to /read/<slug> or shows
# an image from /static must not invent either. Neither is fetched -- both are
# checked against what the repository actually holds, so the check works with no
# network at all.
INTERNAL = re.compile(r"\]\((/[^)\s]+)\)")
SLUGS = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
                     "course", "lesson-slugs.json")
TIMEOUT = 20
AGENT = "shanraq-coursecheck/1.0 (+https://shanraq.org)"


def repo_prefix():
    """The tree URL of this repository, or None when git cannot say."""
    try:
        url = subprocess.check_output(
            ["git", "remote", "get-url", "origin"], text=True).strip()
    except (subprocess.CalledProcessError, FileNotFoundError):
        return None
    url = re.sub(r"^git@github\.com:", "https://github.com/", url)
    url = re.sub(r"\.git$", "", url)
    return url + "/tree/main/"


def fetch(url):
    """The status line for one link, or a short reason it never got one."""
    # A link may carry Kazakh or Russian letters; they have to be percent-encoded
    # before the request, or the failure is an encoding error rather than an
    # answer from the server.
    safe = urllib.parse.quote(url, safe=":/?#[]@!$&'()*+,;=%~")
    req = urllib.request.Request(safe, headers={"User-Agent": AGENT})
    try:
        with urllib.request.urlopen(req, timeout=TIMEOUT) as r:
            return r.status
    except urllib.error.HTTPError as e:
        return e.code
    except Exception as e:  # noqa: BLE001 -- a name that will not resolve, a timeout
        return type(e).__name__


def known_slugs():
    """Every lesson slug the publisher knows, or None when the map is missing."""
    try:
        with open(SLUGS, encoding="utf-8") as f:
            return {slug for slug, _ in json.load(f).values()}
    except (OSError, ValueError):
        return None


def internal_problems(path, slugs, root):
    """Links inside the site that lead nowhere.

    A lesson once sent a reader with no Python to "the next lesson" and the
    address was the course map, whose first button opens the lesson they were
    reading. The link answered 200 and was still wrong; a mistyped slug would
    have answered 404 and nobody would have noticed either.
    """
    out = []
    with open(path, encoding="utf-8") as f:
        text = f.read()
    for href in INTERNAL.findall(text):
        target = href.split("?")[0].split("#")[0]
        if target.startswith("/read/"):
            slug = target[len("/read/"):].strip("/")
            if slugs is not None and slug not in slugs:
                out.append(f"{href} — такого урока нет в lesson-slugs.json")
        elif target.startswith("/static/"):
            if not os.path.isfile(os.path.join(root, "web", target.lstrip("/"))):
                out.append(f"{href} — файла нет в web/static")
    return out


def main(argv):
    offline = "--offline" in argv
    paths = [a for a in argv[1:] if not a.startswith("--")]
    if not paths:
        print(__doc__.strip())
        return 2

    ours = repo_prefix()
    slugs = known_slugs()
    root = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    inside = 0
    links = {}
    for p in paths:
        with open(p, encoding="utf-8") as f:
            for url in LINK.findall(f.read()):
                links.setdefault(url, set()).add(os.path.basename(p))

    bad = []
    # Inside the site first: no network, and the commonest mistake -- a lesson
    # pointing at a page that does not exist -- is caught before anything is
    # fetched.
    for p in paths:
        for note in internal_problems(p, slugs, root):
            bad.append((note, "ссылка внутрь сайта", {os.path.basename(p)}))
        with open(p, encoding="utf-8") as f:
            inside += len(INTERNAL.findall(f.read()))

    # Ours first: a wrong repository address is a mistake in the text, not a
    # server having a bad day, and it is found without asking anyone.
    for url in sorted(links):
        if "github.com" in url and "/course/go-blog/" in url:
            if ours and not url.startswith(ours):
                bad.append((url, f"не наш репозиторий, ожидалось {ours}…", links[url]))

    if not offline:
        rest = [u for u in sorted(links) if not any(u == b[0] for b in bad)]
        with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
            for url, status in zip(rest, pool.map(fetch, rest)):
                if not (isinstance(status, int) and 200 <= status < 400):
                    bad.append((url, f"ответ {status}", links[url]))

    for url, why, where in sorted(bad, key=lambda b: b[0]):
        print(f"{url}\n    {why}\n    в файлах: {', '.join(sorted(where))}")
    print(f"проверено ссылок: {len(links)} наружу и {inside} внутрь, битых: {len(bad)}")
    return 1 if bad else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
