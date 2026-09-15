# Step 30 — one checked artifact, then publication

The final digest after lesson 55. `release.py` copies the complete report into a
temporary directory, writes a SHA-256 manifest, and only then atomically replaces
`public`. If building fails, the previous directory remains untouched.

```console
python -m pytest -q
python -m mypy
python main.py --offline --dry-run
python release.py
python -m http.server 8000 --directory public
```

Inspect the real page before uploading the immutable `public` directory. Retain
the previous artifact for rollback.
