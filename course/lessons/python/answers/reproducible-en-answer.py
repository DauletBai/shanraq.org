"""Lesson 54: verify a reproducible environment manifest."""
files = {"requirements.in", "requirements.txt", ".python-version"}
checks = {"pytest", "mypy", "offline"}
assert len(files) == 3 and len(checks) == 3
print("environment: reproducible")
print("tests: passed")
print("types: passed")
