# Step 27 — the pipeline's promises become tests

The digest after lesson 52, *`pytest` tests: verify the pipeline's promises*.
The runnable application remains the step-26 snapshot; this step adds a test
suite around its public behaviour instead of changing that behaviour.

The tests pin two boundaries:

- model claims with unknown evidence are rejected and causal prose waits for a
  person;
- generated HTML escapes text controlled by data and is written only beneath
  an isolated `tmp_path`.

Run the report and the tests from this directory:

```console
python3 main.py --offline
python3 -m pytest -q
```

A test directory is not production code. Keeping it beside the snapshot makes
the promise readable and lets every later step inherit a regression suite.

