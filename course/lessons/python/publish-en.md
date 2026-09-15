# Finale: verify, publish, and reproduce

_Lead (summary):_ **Lesson fifty-five and the end of the introductory Python course. Combine checks into one gate, build an immutable report directory, publish only after success, and expose verifiable source and build information to the reader.**

## What you have built

Across 55 lessons you moved from `print` to a program that fetches real data, stores it, computes measures, draws a chart, builds HTML, runs on a schedule, and limits trust in model output. The final step adds no library. It prevents half a result from becoming public.

## One gate before publication

`release.py` runs required checks and stops at the first failure:

```text
import subprocess
import sys

CHECKS = [
    [sys.executable, "-m", "pytest", "-q"],
    [sys.executable, "-m", "mypy", "--strict", "sholu", "main.py"],
    [sys.executable, "main.py", "--offline", "--dry-run"],
]

for command in CHECKS:
    print("+", " ".join(command), flush=True)
    subprocess.run(command, check=True)

print("READY")
```

`check=True` stops on a non-zero exit status. Do not swallow that exception to print a green message. Publication verification avoids the network: a cached snapshot makes it repeatable, while data refresh has separate checks.

## Publish an artifact, not a working directory

First build `index.html`, the image, an audit CSV, and `manifest.json` beneath a temporary directory. Record UTC build time, source identifiers, and every file's SHA-256. Only after verification atomically replace `public` with the new directory.

The reader can no longer receive new HTML with an old image. Preview before external publication:

```console
python release.py
python -m http.server 8000 --directory public
```

At `http://localhost:8000`, check mobile width, title, date, sources, table, chart, and the absence of causal prose still marked `review`.

## Final checklist

- the working tree contains only expected changes;
- tests and `mypy` pass in a clean environment;
- offline dry-run writes nothing;
- an empty or incomplete report cannot replace the previous one;
- HTML contains no secret or personal filesystem path;
- date, units, period, and sources are visible;
- user-controlled text is escaped;
- the previous artifact is backed up or quickly recoverable;
- the public URL returns `200` and files match the manifest;
- logs identify the data version and publication outcome.

Publication is a state transition, not “copy some files”: the verified artifact replaces the previous one as a whole, or nothing changes.

## Lesson map

![Lesson map: checks, artifact, publication, and verification](/static/course/py/map-publish-en.svg)

Recall cue: **checks → temporary artifact → atomic replacement → HTTP check → rollback**.

## Say it in your own words

1. Why not copy files individually into the public directory?
2. What belongs in `manifest.json`?
3. Why smoke-test the real URL after publication?
4. What is a rollback?

## Warm-up

**1. Predict.** Does the second command run after the first fails?

<!-- drill 1 -->
```python
import subprocess
import sys

try:
    subprocess.run([sys.executable, "-c", "raise SystemExit(1)"], check=True)
    print("publish")
except subprocess.CalledProcessError:
    print("stopped")
```

**2. Fill the safe digest mode.**

```text
python main.py --offline ...
```

**3. Fix it.** A script deletes the old `public` directory before verifying the new one.

## Final assignment

**Required.** Run one gate, build `public`, create a SHA-256 manifest, open a local HTTP server, and inspect the page with the checklist. Deliberately fail one check and confirm that the previous directory survives.

<!-- task out -->
```text
checks: passed
artifact: complete
publish: ready
```

**With your data.** Publish the digest to a static host and record its URL, check time, and rollback procedure.

**Optional.** Move the gate to CI with minimal permissions and publication restricted to a protected branch.

## Where to go next

The introductory course is complete. Choose a problem rather than more syntax: deeper data analysis, an API or web application, databases and SQL, machine-learning engineering, or a small contribution to an open Python project.

Keep the digest as a portfolio project. A README with the problem, reproducible setup, screenshot, sources, limitations, and working URL communicates skill better than a list of viewed topics.

## Answers

1. A reader could receive mixed versions; replace the checked directory as a whole.
2. Build time, sources, and file checksums.
3. A local build does not test DNS, server, route, or the bytes actually served.
4. Returning to the last known-good artifact.

<!-- drill 1 out -->
```text
stopped
```

2. The flag is `--dry-run`.

<!-- drill 2 -->
```python
command = "python main.py --offline --dry-run"
print(command)
```

<!-- drill 2 out -->
```text
python main.py --offline --dry-run
```

3. Build and verify a complete temporary directory; replace the old one only as the final atomic operation.

<!-- drill 3 -->
```python
states = ["build temporary", "verify", "replace public", "smoke test"]
print(" -> ".join(states))
```

<!-- drill 3 out -->
```text
build temporary -> verify -> replace public -> smoke test
```

## Sources

- [Python: subprocess](https://docs.python.org/3/library/subprocess.html) — exit statuses and `check=True`.
- [Python: hashlib](https://docs.python.org/3/library/hashlib.html) — SHA-256 manifests.
- [The Twelve-Factor App: build, release, run](https://12factor.net/build-release-run) — separating build and run.
