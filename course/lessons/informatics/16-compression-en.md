# Why ZIP and a messenger compress differently

_Lead (summary):_ **Compare lossless and lossy compression, then choose what is acceptable for our project data.**

## Where we are on the map

The previous block showed how to find files and separate a program from its data. Now we ask how data can carry meaning that another person and another program will read in the same way. State your first guess, test it with numbers, and record what you had to revise.

![Why ZIP and a messenger compress differently](/static/course/informatics/map-16-compression-en.svg)

## Everyday analogy and exact model

Writing “4А3Б” instead of “ААААБББ” describes repetition more briefly. Someone who knows the rule can restore the original exactly. Shrink a photograph, however, and fine details can disappear; merely unpacking it cannot bring them back.

Lossless compression permits exact reconstruction. ZIP is in this category. Our “4А3Б” is only a teaching example: real size depends on encoded bytes, headers, and algorithms, so ZIP may make a short file larger. Lossy compression discards some information according to a rule for a smaller representation; the recovered picture or sound need not match byte for byte. A messenger may recompress media but may transfer an attachment as a file instead. Check the actual sending mode.

## Where the analogy ends

“Compress” does not always mean “make smaller,” nor does it always mean “damage.” Choose according to what must be exact: task text and backups need lossless storage; a photographic preview may tolerate lost detail.

## The lesson's support signal

`lossless → exact reconstruction; lossy → some details cannot be recovered`

## Worked example

ААААБББ has seven visible characters and our symbolic 4А3Б has four. Expand the latter back to seven. Do not claim that this saves three bytes: Cyrillic letters use more than one byte in UTF-8 and a real format must define the rule.

## Predict before observing

Can you recover exact original RGB values after mapping every channel value 201 to 200? No: several originals might map to the same rounded result. Decide whether an exact copy is needed before choosing a format.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Invent a repeated string, shorten it using the classroom rule, and reconstruct it. Then give an example of a photo detail lost by reducing resolution.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“Sending a PNG in a messenger always damages it.” Correct the claim: treatment depends on sending mode. Compare source and received files, and transfer exact data as a file.

## Transfer to a new situation

How should you keep an assistant task table and a small decorative photo? The table needs exact copying; you may show a reduced photo while preserving its original separately.

## Project change

Write a FORMAT.md policy: tasks.json and its backup must be recoverable byte for byte; a preview cover may be optimized after checking readability. Do not replace the only original with a reduced copy.

## Exercise

Choose suitable treatment for three objects: task list, project archive, and thumbnail. Justify each choice, then test reversibility using your own short repeated string.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow name both compression types. In a week expand 4А3Б. In a month check whether another app preserved your transferred data.

[Next lesson: How to detect damaged data](/read/informatics-17-integrity-errors?lang=en)
