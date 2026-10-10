# Encryption, keys, and the limits of protection

_Lead (summary):_ **A student puts a note in a sealed envelope.**

## Where we are on the map

Lesson 59 of 72. Digital security block (55–62). We continue a local assistant with three fictional tasks. Each control has a defined boundary: the classroom version does not become a public service.

![Encryption, keys, and the limits of protection](/static/course/informatics/map-59-encryption-keys-en.svg)

## Situation and question

A student puts a note in a sealed envelope. You can see a letter exists but not read its text. Yet the envelope does not prove who sent it or replace keeping the key safe.

## New words without gaps

**Encryption** turns readable data into a form unreadable without a key. A **key** is a secret needed for transformation or verification. **TLS** protects transmission between browser and server when the name and certificate are checked correctly. **Disk encryption** protects a different part of the journey: the saved file. A **hash** used for integrity does not encrypt or hide data.

## The lesson's support signal

Read the three cells left to right, then cover the third. State the input and action in each cell; predict how the path ends and which observable fact checks it. Compare the drawing with a project file or worked example. If a number on the map cannot be derived from the data, correct the map rather than changing data for a pretty picture.

## Work through it step by step

Draw browser → network → server → disk → backup. Mark TLS on the network leg and disk protection at the file. Our `step-06` page uses HTTP only on `127.0.0.1`; it is not a protected public service example. Its `backup.db` is unencrypted, so keep it on a trusted device and never place it in a shared cloud folder.

## Predict and check

Write the expected output before running code. Execute it in `step-06` or `step-07`, compare every character and explain each line. Then change one safe fictional input and repeat “predict — run — change — explain — verify.” Do not turn an observation on a tiny dataset into a promise about real people.

```python
import hashlib
a = hashlib.sha256(b"20 minutes").hexdigest()
b = hashlib.sha256(b"200 minutes").hexdigest()
print(a == b)
```

## Expected output

```text
False
```

## Catch the error

HTTPS does not hide data from the destination server. A browser padlock alone does not prove the site owner is honest. A hash stored beside a file does not protect against someone who can replace both; authenticity needs another trusted mechanism.

## Project change

Release 1.3 preserves the three tasks, four sessions and reminder rules. New operations live in `security_assistant.py`: validation, a separate backup and restore only into a new file. State which threat each change reduces and which it does not; all teaching data stay fictional.

## Task and evidence

For each part of the route, name who could read the data. Explain what must change before a real multiuser service is published.

## Transfer to a new setting

A school transfers marks over HTTPS but leaves the file on a shared computer. Which boundary is protected, and which is not?

## Return after 1, 7, and 30 days

After 1 day redraw the map from memory and name one boundary of the idea. After 7 days explain it to a friend without reading the lesson and answer 7 of 10 block questions; analyse mistakes with concrete inputs. After 30 days open a fresh project copy, reproduce the check and transfer the idea to a new case. Keep a record of input, expected output and actual output.

## Primary reference to check

[Official reference](https://www.nist.gov/publications/guidelines-selection-configuration-and-use-transport-layer-security-tls-implementations)

## Next lesson

[Backups, updates, and event logs](/read/informatics-60-backup-updates-logs?lang=en)
