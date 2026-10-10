# Passwords, hashes, and a second factor

_Lead (summary):_ **Two people know the same short password.**

## Where we are on the map

Lesson 56 of 72. Digital security block (55–62). We continue a local assistant with three fictional tasks. Each control has a defined boundary: the classroom version does not become a public service.

![Passwords, hashes, and a second factor](/static/course/informatics/map-56-passwords-hashing-2fa-en.svg)

## Situation and question

Two people know the same short password. Even if a site stores no plain password, a stolen list might let an attacker guess it. Separate what a hash, a salt and a second factor each do.

## New words without gaps

A **password** is a secret you know. **Hashing** turns input into a fingerprint that is not normally reversed. A **salt** is a separate random value for each record; a deliberately costly key-derivation function slows mass guessing. A **second factor** proves access with another kind of evidence, such as a device. Two passwords are still one factor type. A **password manager** helps create long unique secrets.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

On paper compare “two identical passwords” with “two different salts”: the stored results differ. `step-06` collects no passwords and has no accounts; this lesson models the idea rather than inviting you to type a secret into the classroom assistant. Check the sign-in address before entering a real password and enable MFA on a real service where available. Store recovery codes separately from the phone.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
import hashlib, secrets
salt = secrets.token_bytes(16)
a = hashlib.scrypt(b"fictional long phrase", salt=salt, n=2**14, r=8, p=1)
b = hashlib.scrypt(b"fictional long phrase", salt=salt, n=2**14, r=8, p=1)
print(secrets.compare_digest(a, b))
```

## Expected output

```text
True
```

## Catch the error

Fast SHA-256 alone is not an appropriate password storage design. A salt is not a secret key and cannot make a weak password strong. A second factor does not replace checking the site you visit.

## Project change

Release 1.3 preserves the three tasks, four sessions and reminder rules. New operations live in `security_assistant.py`: validation, a separate backup and restore only into a new file. State which threat each change reduces and which it does not; all teaching data stay fictional.

## Task and evidence

Use a diagram to explain why identical passwords with different salts yield different records. Name an independent second factor and a way to recover when the device is lost.

## Transfer to a new setting

A classmate suggests one shared password for a project. How can each person have separate access so removing one person’s rights does not change everyone’s secret?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://pages.nist.gov/800-63-4/sp800-63b.html)

## Next lesson

[Who may do what: least privilege](/read/informatics-57-authorization-least-privilege?lang=en)
