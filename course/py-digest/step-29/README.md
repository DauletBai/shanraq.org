# Step 29 — the environment is part of the project

The digest after lesson 54. `.python-version` names the interpreter,
`requirements.in` records direct intent, and `requirements.txt` is the exact
snapshot used by the course. Rebuild the lock with pip-tools in a clean environment.

```console
python -m pip install pip-tools==7.6.1
python -m piptools compile --generate-hashes requirements.in
python -m piptools sync requirements.txt
python -m pytest -q
python -m mypy
python main.py --offline --dry-run
```

Do not use an old workstation's `pip freeze` as project intent. When regenerated,
the lock's hashes and transitive dependencies belong in Git with its input.
