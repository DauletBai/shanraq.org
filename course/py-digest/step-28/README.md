# Step 28 — types describe the model boundary

The digest after lesson 53. Runtime behaviour remains unchanged; `sholu/types.py`
adds a typed `Observation`, an honest optional result, and narrowing before field
access. `pyproject.toml` starts strict checking with this boundary module rather
than hiding a whole-project migration behind `Any`.

```console
python3 main.py --offline
python3 -m pytest -q
python3 -m mypy
```

A small strict boundary is more useful than a large configuration full of
ignored errors. Expand the checked set only after its annotations are real.
