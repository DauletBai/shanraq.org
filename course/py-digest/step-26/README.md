# Step 26 — a plausible cause waits for evidence

The state of the digest after lesson 51, *What a model cannot be trusted with*,
and the end of the AI module.

Step 25 proved that a model cannot replace a sourced number. This step adds a
second boundary in `sholu/model.py`:

- observations with known evidence may be accepted;
- calculations need a named, reproducible formula;
- causal prose is routed to `review` even when it cites the same data row;
- unknown sources and claim types are rejected.

The page states the refusal explicitly:

```
Себеп туралы тұжырым автоматты жарияланбады:
оған бөлек дәлел мен адам тексеруі керек.
```

That sentence is policy-owned, not model-written. `review` never means “publish
with a warning”; it means the claim stays outside the report until independent
evidence and a person support it.

```
pip install -r requirements.txt
python3 main.py --dry-run
python3 main.py --offline
```
