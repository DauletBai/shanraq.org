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
import os
import re
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request

LINK = re.compile(r"\]\((https?://[^)\s]+)\)")
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


def main(argv):
    offline = "--offline" in argv
    paths = [a for a in argv[1:] if not a.startswith("--")]
    if not paths:
        print(__doc__.strip())
        return 2

    ours = repo_prefix()
    links = {}
    for p in paths:
        with open(p, encoding="utf-8") as f:
            for url in LINK.findall(f.read()):
                links.setdefault(url, set()).add(os.path.basename(p))

    bad = []
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
    print(f"проверено ссылок: {len(links)}, битых: {len(bad)}")
    return 1 if bad else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
