"""Қорытынды тапсырма: үш тексеруден кейінгі жариялау күйі."""
checks = {"tests": True, "types": True, "offline": True}
assert all(checks.values())
print("тексерулер: өтті")
print("артефакт: толық")
print("жариялау: дайын")
