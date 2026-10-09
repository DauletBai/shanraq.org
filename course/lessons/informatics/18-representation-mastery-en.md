# Mastery checkpoint: the assistant's data format

_Lead (summary):_ **Assemble version 0.2 of the Digital Assistant: a format for three fictional tasks, reading rules, and evidence that copies survive.**

## Where we are on the map

Another learner must now receive your data and understand it without talking to you. Can they open a Kazakh title, add a task, and notice a damaged copy?

![Mastery checkpoint: the assistant's data format](/static/course/informatics/map-18-representation-mastery-en.svg)

## Everyday analogy and exact model

Imagine handing the project to a classmate without an explanation. They must open the file, see three tasks, understand its fields, and recover the same data after copying. If they must call the author, the format is not ready.

A data format is an agreement about structure and meaning. Version 0.2 uses UTF-8 JSON with a version field and a tasks array. Each task has a unique text id, a title, and a Boolean done. Three fictional records are t-01 “Кітап оқу” false; t-02 “Сурет салу” true; t-03 “Есеп шығару” false. JSON object field order does not change meaning, while array element order may matter for display. Store the file losslessly and compare an exact copy after transfer.

## Defend the data contract

A **format specification** is a written agreement about valid files and fields. A **type** limits a value: `done` is Boolean `true` or `false`, not the string `"true"`; `title` is nonempty text. An **identifier** `id` distinguishes records even if their titles match. **Format version** `0.2` says which rules to use to read the file; it is not a task count. **Backward compatibility** means a new version reads older valid records, if that promise has been stated.

Give a partner only `tasks.json` and `FORMAT.md`. Without hints, they must add fictional task `t-04` with `Ә`, choose `done: false`, save in UTF-8, and explain why quoted `"false"` is wrong. They then change one byte in a copy: an integrity check should reveal the difference. If they must guess what `version` means or which `id` is valid, improve the contract and retry.

For an 8/10 pass, three forms of evidence are essential: a readable file, transfer to a new record, and a backup matching the original before deliberate damage. Test again on a different example after seven days. A screenshot alone proves neither encoding, types, nor preservation of bytes.

## Where the analogy ends

Textual JSON does not store done as a single physical bit. Equal line counts do not prove files equal. The classroom checksum from the prior lesson is not a security measure. For a practical exact-copy check use byte comparison or a cryptographic hash from a trusted source.

## The lesson's support signal

`three fictional tasks → UTF-8 JSON → reopen → check fields → compare copy`

## Worked example

Write {"version":"0.2","tasks":[{"id":"t-01","title":"Кітап оқу","done":false},{"id":"t-02","title":"Сурет салу","done":true},{"id":"t-03","title":"Есеп шығару","done":false}]}. Check three distinct ids, three titles, and exactly one true. Save as data/tasks.json in UTF-8.

## Predict before observing

If a copied file changes t-02 done to false, can it still be valid JSON? Yes. Does it preserve the original data? No. Syntax validation and exact-copy checking are separate tests.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Hide the sample and describe the required fields of one task yourself. Reopen the file and check that you have not turned Boolean false into the string "false".

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“The file opened, so its data are correct.” Correct this: opening is only step one. Check the schema, three ids, done values, and copy equality.

## Transfer to a new situation

Give the format to a classmate using another operating system. They should open Ә in a new fictional task title without corruption and explain the role of UTF-8; changing platforms must not change the meaning.

## Project change

Create data/tasks.json and FORMAT.md. In FORMAT.md describe UTF-8, version=0.2, required fields, valid done values, unique ids, and a check sequence. Keep a backup and compare it with the source byte by byte or with SHA-256.

## Exercise

Complete the checkpoint: 8/10 checks now, including independent transfer to an unfamiliar record, and 7/10 in seven days. Show the file, rules, and comparison result. If you miss a step, fix its rule and repeat on a comparable new example.

Check ten points: (1) UTF-8; (2) valid JSON; (3) version = 0.2; (4) the original three tasks; (5) all ids distinct; (6) title is nonempty text; (7) done is a Boolean, not a string; (8) exactly one true; (9) a new fictional record follows the same rules; (10) the backup matches the source byte for byte. Points 9 and 10 are required for a pass at 8/10.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow recall the three task fields without the page. In seven days ask someone else to add a fictional record using FORMAT.md and check at least 7/10 items. In a month reopen the copy and explain the whole path from character to byte.

[Next lesson: Decompose a large problem without losing the goal](/read/informatics-19-problem-decomposition?lang=en)
