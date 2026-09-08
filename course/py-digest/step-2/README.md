# Step 2 — the same program, inside a project

The state of the digest after lesson 2, *A workplace*.

The code has not changed. What changed is around it: the program now lives in a
project with an environment of its own, a list of dependencies and a
`.gitignore`, so that a library installed for this project cannot break another
one.

```
python3 -m venv .venv
source .venv/bin/activate          # Windows: .venv\Scripts\Activate.ps1
python -m pip install -r requirements.txt
python main.py
```

`requirements.txt` is empty of packages on purpose: at this point the digest
still needs nothing beyond the standard library. The habit is what is being
set up, not the dependency.
