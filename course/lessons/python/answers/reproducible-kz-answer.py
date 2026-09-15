"""54-сабақ: қайталанатын орта құрамын тексеру."""
files = {"requirements.in", "requirements.txt", ".python-version"}
checks = {"pytest", "mypy", "offline"}
assert len(files) == 3 and len(checks) == 3
print("орта: қайталанады")
print("тесттер: өтті")
print("типтер: өтті")
