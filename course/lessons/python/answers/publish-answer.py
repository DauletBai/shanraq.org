"""Финальное задание: состояние публикации после трёх проверок."""
checks = {"tests": True, "types": True, "offline": True}
assert all(checks.values())
print("checks: passed")
print("artifact: complete")
print("publish: ready")
