# Step 25 — a model cannot replace a number

The state of the digest after lesson 50, *A model on our data*.

The previous step had no language-model boundary. This one adds it without
making Ollama a condition for rebuilding the report:

- `sholu/model.py` verifies `source_id`, year, value and unit against the report;
- decimal values are compared exactly and a model label may contain no digits;
- the final sentence receives its number from the report, not from the model;
- `main.py` uses a deterministic sample response so CI and `--offline` still
  run without a model. A later step may replace that sample with an Ollama call
  without moving the trust boundary.

The verified line is visible above the sources in `report.html`:

```
Қазақстандағы инфляция: 11.39 % (2025). [inflation-KAZ-2025]
```

Run it as before:

```
pip install -r requirements.txt
python3 main.py --dry-run
python3 main.py --offline
```

A failed comparison aborts the rebuild. The old report stays in place; an
unverified sentence never becomes a softer warning on a public page.
