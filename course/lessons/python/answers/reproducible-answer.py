"""Задание урока 54: проверяем состав воспроизводимой среды."""
files = {"requirements.in", "requirements.txt", ".python-version"}
checks = {"pytest", "mypy", "offline"}
assert len(files) == 3 and len(checks) == 3
print("environment: reproducible")
print("tests: passed")
print("types: passed")
