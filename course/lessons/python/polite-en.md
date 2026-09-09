# Politeness to somebody else's server: a timeout, a backoff and a cache

_Лид (summary):_ **The twentieth lesson of the Python course. Politeness to somebody else's server is three concrete things: a timeout, so you do not wait for ever; a retry with the pause doubling, so you do not hit what is already down; and a cache, so you do not ask twice. Measured on a server of our own: four readings, one request.**

## Why this is needed

The last three lessons went to the National Bank. It worked — and that is exactly why it is time to say how not to do it.

Somebody else's server owes you nothing. It may go quiet, answer "try later", go down for maintenance. A program written on the assumption that it always answers quickly will one day hang — and a program that answers a refusal with an immediate retry will finish off what was already struggling.

Politeness here is not manners but three concrete things: a **timeout**, a **retry with a backoff**, and a **cache**. We learn them on a server of our own, because practising on somebody else's is the very rudeness the lesson is about.

## The whole thing at once

The file is `sypaiy.py`. Run it with `python sypaiy.py` from inside the environment.

The first twenty lines are a teaching server that can stall, refuse, and count requests; knowing it by sight is enough, and the course will come back to how servers are built. The required part is the three blocks after it.

```python
"""Lesson 20: politeness to somebody else's server -- timeout, retry, cache.

The last three lessons went to a real bank over the network. Today we learn to
do that without waiting for ever and without getting in the way of whoever we
are asking -- and we learn it on a server of our own, because practising on
somebody else's is the very rudeness in question.
"""

import json
import threading
import time
import urllib.error
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

# How many times the server was asked for each address. That is the measure.
HITS = {"data": 0, "flaky": 0}


class Teacher(BaseHTTPRequestHandler):
    """A teaching server beside us: it can stall, fail, and count requests."""

    def do_GET(self):
        if self.path == "/slow":
            time.sleep(2)                      # the server is thinking
            self.answer(200, b'{"value": 1}')
        elif self.path == "/flaky":
            HITS["flaky"] += 1
            if HITS["flaky"] < 3:
                self.answer(503, b"try later")  # temporarily unable
            else:
                self.answer(200, b'{"value": 42}')
        else:
            HITS["data"] += 1
            self.answer(200, b'{"value": 510.43}')

    def answer(self, code, body):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *args):
        return                                  # keep its log out of our output


server = ThreadingHTTPServer(("127.0.0.1", 0), Teacher)
threading.Thread(target=server.serve_forever, daemon=True).start()
BASE = f"http://127.0.0.1:{server.server_address[1]}"


def get(url, timeout=1.0):
    """One request: a body, or an exception."""
    request = urllib.request.Request(url, headers={"User-Agent": "shanraq-course/1.0"})
    with urllib.request.urlopen(request, timeout=timeout) as answer:
        return json.load(answer)


print("== the timeout: how long to wait on somebody else's silence")
try:
    get(BASE + "/slow", timeout=0.3)
except TimeoutError:
    # Silence after the connection arrives exactly like this: TimeoutError.
    print("waited 0.3 s and stopped waiting")
except urllib.error.URLError as error:
    # And a connection that never happened is a URLError with the reason in it.
    print("could not connect:", type(error.reason).__name__)

print()
print("== retry with a pause: three attempts instead of one")
delay = 0.1
for attempt in range(1, 5):
    try:
        data = get(BASE + "/flaky")
    except urllib.error.HTTPError as error:
        # Only what can get better on its own is worth repeating.
        if error.code not in (429, 500, 502, 503, 504):
            print("not retrying:", error.code)
            break
        print(f"attempt {attempt}: {error.code}, waiting {delay:.1f} s")
        time.sleep(delay)
        delay *= 2
        continue
    print(f"attempt {attempt}: got {data['value']}")
    break
print("requests to the server:", HITS["flaky"])

print()
print("== the cache: the same answer is not asked for twice")
CACHE = {}
TTL = 60


def cached(url):
    """The answer from memory, until it goes stale."""
    now = time.monotonic()
    if url in CACHE and now - CACHE[url][0] < TTL:
        return CACHE[url][1], "from the cache"
    data = get(url)
    CACHE[url] = (now, data)
    return data, "from the network"


for _ in range(3):
    data, where = cached(BASE + "/data")
    print(f"{data['value']} — {where}")
print("requests to the server:", HITS["data"])

server.shutdown()
```

It prints:

```
== the timeout: how long to wait on somebody else's silence
waited 0.3 s and stopped waiting

== retry with a pause: three attempts instead of one
attempt 1: 503, waiting 0.1 s
attempt 2: 503, waiting 0.2 s
attempt 3: got 42
requests to the server: 3

== the cache: the same answer is not asked for twice
510.43 — from the network
510.43 — from the cache
510.43 — from the cache
requests to the server: 1
```

## Taking it apart

### The timeout: how long to wait on somebody else's silence

```
waited 0.3 s and stopped waiting
```

Without a timeout `urlopen` waits for as long as the server stays quiet — minutes. From outside that is indistinguishable from a hung program, and it looks worst on a schedule: a task that should have taken a second is still hanging in the morning.

A timeout answers the question "how long are we prepared to wait", and that question **always has an answer**, even when nobody wrote one: then it is "as long as it takes".

Note what each failure arrives as:

```python
except TimeoutError:      # connected, and never got an answer
except urllib.error.URLError:   # never connected at all
```

Silence after a connection is a `TimeoutError`. A connection that never happened is a `URLError` with the reason inside it. Catch both: for the program they are different troubles, and both mean "no data".

> **Picture it.** A call to an enquiry desk. You put the phone down after the third ring not out of spite but because you have other things to do.

### Retry with a backoff: three attempts instead of one

```
attempt 1: 503, waiting 0.1 s
attempt 2: 503, waiting 0.2 s
attempt 3: got 42
requests to the server: 3
```

`503` means "I cannot right now, try later". That is temporary, and a retry belongs here. But a retry without a pause is a blow to somebody who is already failing: the server went down under load, and you added another hundred requests a second.

So the pause **doubles**: 0.1 → 0.2 → 0.4. That backoff — exponential backoff, in the usual phrase — gives the source time to get up and keeps your own run from stretching to infinity.

The second rule matters more than the first: **only what can get better on its own is worth repeating**.

```python
if error.code not in (429, 500, 502, 503, 504):
    ...   # not retrying
```

A `404` does not become a `200` because you asked again: the address is not there. A `403` does not become permission. But a `429` — "too often" — is asking you to wait outright, and often says how long, in a `Retry-After` header; when it is there, listen to it rather than to your own formula.

And retries need an end. Three or four attempts, and then the program says honestly that the source did not answer, instead of hanging until morning.

### The cache: the same answer is not asked for twice

```
510.43 — from the network
510.43 — from the cache
510.43 — from the cache
requests to the server: 1
```

Three readings, one request. That is the whole arithmetic of a cache, and it is also the answer to what a cache is for: not speed, but not asking somebody else's server about what you already know.

A cache must have a **life**. Today's exchange rate does not change all day, a year's inflation for a year, the weather for half an hour. That life is the decision about how fresh the data has to be; without it a cache turns into "the week before last", which is worse than slow.

Ours lives in memory and disappears with the program. For a digest that runs once a day that is not enough — there the cache goes on disk, as we already did with a file in lesson eleven; the key is the address, the contents are the answer and the time it arrived.

### Three more things that make you a welcome guest

**Sign your requests.** A `User-Agent` with a link to the project is a chance for somebody to write to you rather than ban your address.

**Do not ask in parallel what can be asked once.** A list of two hundred towns walked in a loop with no pause is two hundred requests a second; with a `time.sleep(0.2)` between them it is five, and nobody notices.

**Read the terms.** Many sources have a page about how often you may ask. The National Bank and the World Bank hand out data without a key — that is trust, not an invitation.

## The map of the lesson

![The map of the lesson: the timeout, the retry and the cache](/static/course/py/map-polite-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What does a program with no timeout do when the server goes quiet?
2. Why does the pause between retries double rather than stay the same?
3. Which response codes are pointless to retry, and why?

## Warm-up

Three short steps before the exercise: predict, fill in, fix. The answers are at the end of the lesson, but answer them yourself first.

**1. Predict.** What does this program print?

<!-- drill 1 -->
```python
delay = 0.5
for attempt in range(1, 5):
    print(attempt, round(delay, 2))
    delay *= 2
```

**2. Fill in the gap.** In place of `...` put the condition: only these codes are worth retrying.

```python
RETRY = (429, 500, 502, 503, 504)

for code in (404, 500, 429, 403):
    print(code, "retrying" if ... else "not retrying")
```

**3. Fix it.** The source is not answering today, and this program will ask it for ever without a single pause. Give the retries an end and a backoff.

```python
attempts = 0


def fetch():
    """A source that is not answering today."""
    global attempts
    attempts += 1
    raise TimeoutError("timed out")


while True:
    try:
        print(fetch())
        break
    except TimeoutError:
        continue
```

## Exercise

**Required.** Write one polite function `polite(url)` and try it against the teaching server, which refuses the first time it is asked and answers the second (take the server from the lesson — it is in your own file).

The function must: set a `timeout`; retry **only** the codes `429, 500, 502, 503, 504`, up to three attempts, with the pause doubling from 0.1 s; print every failed attempt; put a successful answer into a cache with a life; and raise an error of its own once the attempts run out. Read the rate four times in a row and print how many times the server was troubled.

The expected output:

<!-- task out -->
```
  attempt 1: 503, waiting 0.1 s
1: 510.43 — from the network
2: 510.43 — from the cache
3: 510.43 — from the cache
4: 510.43 — from the cache
requests to the server: 2
```

Done when: the output matches line by line; the request has a `timeout`; `404` and `403` are not retried at all; the retries end rather than continue; four readings cost exactly two requests to the server.

**On your own data.** Take your request from lesson eighteen — the rate for a fixed day — and wrap it in this same function. Run the program three times in a row and look at the counter: that is the difference between a polite program and a rude one.

**Optional.**

- Make the cache a file: the key is the address, the contents the answer and the time it arrived.
- Add reading of the `Retry-After` header when the server answers `429`.
- Walk a list of ten addresses in a loop with `time.sleep(0.2)` and without, and compare how many requests a second each makes.

## Where this goes in the project

The digest stops being a problem for its sources. Every request it makes is now bounded in time, refusals are survived with a backoff, and a second run on the same day reads yesterday's answer off the disk instead of going back to the bank.

Still open. Our cache is in memory, and the digest needs it on disk — and there the question of what to do with old entries appears. We come back to it in the next lesson: where to put data so that it survives being switched off.

## The answers

### To the questions

1. It waits for as long as the server stays quiet — and from outside that is indistinguishable from a hang. A timeout is always there even when nobody wrote one: then it equals "as long as it takes".
2. Because a source usually refuses under load, and evenly spaced retries add to that load exactly when it is failing. Doubling gives it time to get up.
3. `404` and `403`: the address is not there and access is closed, and neither changes on a repeat. What is retried is the temporary: `429`, `500`, `502`, `503`, `504`.

### To the warm-up

1. `1 0.5`, `2 1.0`, `3 2.0`, `4 4.0`. The pause doubles, and by the fourth attempt the program waits four seconds — exactly the behaviour somebody else's server wants of it.

<!-- drill 1 out -->
```
1 0.5
2 1.0
3 2.0
4 4.0
```

2. `code in RETRY`. The list of temporary codes is written once and separately, so that the condition reads as words rather than as a handful of numbers.

<!-- drill 2 -->
```python
RETRY = (429, 500, 502, 503, 504)

for code in (404, 500, 429, 403):
    print(code, "retrying" if code in RETRY else "not retrying")
```

<!-- drill 2 out -->
```
404 not retrying
500 retrying
429 retrying
403 not retrying
```

3. The retries had neither an end nor a pause: a `while True` with a `continue` is requests without a stop, and they help neither you nor the source. The attempts must be a countable number, with a growing pause between them, and an honest answer after them that there is no data:

<!-- drill 3 -->
```python
import time

attempts = 0


def fetch():
    """A source that is not answering today."""
    global attempts
    attempts += 1
    raise TimeoutError("timed out")


delay = 0.1
for attempt in range(1, 4):
    try:
        print(fetch())
        break
    except TimeoutError:
        time.sleep(delay)
        delay *= 2
else:
    print("gave up; attempts:", attempts)
```

<!-- drill 3 out -->
```
gave up; attempts: 3
```

## Sources

- [Python: urlopen and timeout](https://docs.python.org/3/library/urllib.request.html#urllib.request.urlopen)
- [Python: http.server — the teaching server from the lesson](https://docs.python.org/3/library/http.server.html)
- [MDN: the Retry-After header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After)
- [MDN: response code 429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429)
