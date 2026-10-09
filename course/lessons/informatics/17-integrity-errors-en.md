# How to detect damaged data

_Lead (summary):_ **Practice detecting changed data and see why a simple checksum cannot guarantee authenticity.**

## Where we are on the map

Two files open without an error, but one says the task is no longer done. How can you prove the copy changed when it looks almost the same?

![How to detect damaged data](/static/course/informatics/map-17-integrity-errors-en.svg)

## Everyday analogy and exact model

A cashier checks the number of items against a receipt. A mismatch shows a problem; a match does not prove that every item is unchanged. File transfers need a similar but more exact comparison.

For a classroom checksum, add 4, 7, 2 and take the remainder after division by 10: 13 mod 10 = 3. If 7 changes to 8, the checksum becomes 14 mod 10 = 4, so the change is noticed. Yet 5, 6, 2 also gives 13 mod 10 = 3. Different data with the same checksum form a collision. Cryptographic hashes such as SHA-256 make accidental matches vastly less likely, but a hash without a trusted source does not prove who created a file.

## Opening is not the same as matching

**Integrity** means data has not changed unexpectedly. A **checksum** is a compact result computed from file contents; changing bytes usually changes it. A **SHA-256 hash** is one such check, but matching hashes without a trusted original value do not prove who created the file. **JSON syntax** gives the rules for valid writing; a **data schema** states which fields and types are needed; a **semantic check** asks whether the record makes sense for the task.

Make two copies of `tasks.json`. In one, change `done: true` to `done: false` while keeping commas and quotes correct. Both may still pass a syntax check. Compare hashes: they differ. Next, change a title and restore it: the final files can match again even though their histories differ. A hash describes current bytes, not every action that happened.

For the assistant, perform three checks: open JSON, verify required `id`, `title`, and `done` fields with their types, then verify the three expected fictional records. Compare a backup with its source immediately after creation. If you record its own hash only after corruption, you cannot detect the earlier loss; you need a value saved in advance.

## Where the analogy ends

An integrity check detects change only against a retained reference. A checksum cannot repair an error by itself, replace a backup, or prove authorship if someone can replace both file and advertised hash.

## The lesson's support signal

`data → calculate sum → keep reference → calculate again → compare; a match is not proof of authorship`

## Worked example

For 4,7,2 obtain 3. Change 7 to 8 and obtain 4. Change 4 to 5 and 7 to 6 and obtain 3 again. These steps show both detection and the method’s limit.

## Predict before observing

Can the classroom checksum distinguish 4,7,2 from 5,6,2? Calculate both remainders first. Both are 3, so it cannot.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Invent another three-digit list totaling 13 without looking at the example. Check that data differ but the remainder is the same.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“The checksum matches, so the file is certainly authentic and unchanged.” Correct this: a match means only that this method found no difference. Weak checksums have collisions.

## Transfer to a new situation

You download a file and its hash from the same unknown website. Comparison can find accidental download damage. Does it prove the website is trustworthy? Explain.

## Project change

For three fictional assistant records, note a classroom checksum of numerical ids and label it as an illustration, not security. For genuine file transfers later, plan to use SHA-256 obtained through a trusted channel.

## Exercise

Show one change that the classroom checksum detects and another it misses. Say exactly what the method checks and what it does not promise.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow reconstruct the 4,7,2 example. In a week invent a collision. In a month explain why the reference must come from a trusted source.

[Next lesson: Mastery checkpoint: the assistant's data format](/read/informatics-18-representation-mastery?lang=en)
