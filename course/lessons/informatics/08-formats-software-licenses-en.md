# Formats, software, and the creator's permission

_Lead (summary):_ **We separate bytes from reading rules, choose an open format, and check a licence before adding someone else's material.**

## Where we are on the map

You renamed a photo from `.jpg` to `.png`, but the picture did not become a PNG. Why does changing the label leave the contents untouched?

![Formats, software, and the creator's permission](/static/course/informatics/map-08-formats-en.svg)

## A familiar image and the precise model

A format defines data layout and meaning. An encoding maps numbers to characters. A MIME type describes exchanged content. A licence states permitted use; no licence does not mean freedom. Access to a file does not grant copying rights, and an extension does not change bytes.

## Label, file rules, and the creator's permission

A **format** is a set of rules for data inside a file: where its header sits, which values are valid, and how to read them. An **extension** is the end of a filename; it hints at a format but does not convert it. A **program** can correctly read only formats it knows. **Conversion** actually rewrites the data under another format's rules. A **license** states which uses the creator allows; being able to copy a file does not give you permission to publish it.

Copy `tasks.json` as `tasks.txt`. Open both as plain text: the bytes remain the same. Rename the copy to `tasks.png`: an image viewer cannot turn the text into a picture. Real conversion needs software that reads the source and writes a new file using PNG rules. Compare a checksum or file size before and after renaming: the name changes, the contents do not. Next, break one quotation mark in the JSON and explain why the extension is still correct while the content no longer follows the format.

For the assistant, keep a table: `asset — format — program that opens it — source — permission to redistribute`. Mark your fictional records as your own work. Check the original source's license before adding another person's image to the project.

## The limit of the analogy

The familiar image opens the topic but does not replace the mechanism. State where the analogy stops being accurate.

## The lesson's support signal

`bytes + format + encoding → data | source + licence + attribution → permitted use`

## Worked example

We choose UTF-8 JSON with `id`, `title`, `date`, `state`, and `category`. Valid JSON does not guarantee valid dates or required fields; our schema checks them.

## Predict before observing

Before the experiment, hide the continuation and predict the result in writing. Give the chain of causes from the support signal as well as the answer. Perform the action, record the observation, and compare it with the prediction. A match without an explanation is luck rather than understanding; a mismatch is useful evidence about the link that needs rebuilding.

## Recall without a prompt

Hide the page, rebuild the signal, and explain every transition with a new example.

Explain the topic to someone who has not read the lesson. Do not use a new term until you have unpacked it in plain words. Then restore its precise name and show its place in the system. The listener should be able to predict the next step and give a different example. If they merely repeat a sentence, reduce the explanation to the support signal and rebuild it.

## Find and correct the mistake

“A search-result image is free to take.” Find the original source, creator, licence, date, and terms for the intended use.

## Transfer to a new setting

Apply the model to a device or file absent from the lesson. Separate observation from assumption.

## Project change

Fix the JSON schema. For external material, record source, creator, licence, and check date. External material is optional in the first release.

Keep four lines in the decision log: the intended result, the observation, the reason for the chosen action, and the verification. A screenshot can support evidence but cannot replace text and the working file. Do not use personal data: three fictional records provide the same testability and let you show the project safely to another person.

## Exercise

**Required.** Complete the example, correct the error, and make the project change with a reason.

**With your own data.** Test the decision on three fictional records.

**Optional.** Ask another person to rebuild the explanation from the map.

The readiness criterion is concrete: reproduce the signal without the page, correct the proposed error, transfer the model to an unfamiliar example, and show the project change. If one part fails, return to that link instead of rereading the whole lesson. After repairing it, repeat only an equivalent task with different data.

## Retrieval after 1, 7, and 30 days

Rebuild the signal tomorrow; solve an equivalent task after seven days; after thirty days, verify the decision in the project.

[Next lesson: versions, collaboration, and accessibility](/read/informatics-09-versions-collaboration-accessibility?lang=en)
