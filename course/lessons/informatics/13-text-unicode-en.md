# How a computer distinguishes Ә, Я, and A

_Lead (summary):_ **Find out how one file can preserve Kazakh Ә, Russian Я, and Latin A, and why a character is not necessarily one byte.**

## Where we are on the map

A classmate sent a file named “Ән,” but their friend's screen showed strange marks. Find which step from character to bytes lost its rule.

![How a computer distinguishes Ә, Я, and A](/static/course/informatics/map-13-text-unicode-en.svg)

## Everyday analogy and exact model

You send a friend the task name “Әлем.” If they see strange marks, the bytes may have been read using a different encoding from the one used to save them. Think of an address: one apartment number can take different amounts of space when written down.

Unicode assigns code points to characters: A is U+0041, Я is U+042F, and Ә is U+04D8. UTF-8 turns code points into bytes: A → 41, Я → D0 AF, Ә → D3 98 (hexadecimal notation). Thus A takes one byte; Я and Ә take two bytes each in UTF-8. Counts of visible marks, code points, and bytes can differ, especially for combined characters.

## Character, number, and bytes are different things

A **character** is a unit of text such as `Ә`. A **Unicode code point** is its assigned number; `Ә` is U+04D8. **UTF-8 encoding** defines how to store that number as bytes: U+04D8 becomes `D3 98`. **Decoding** applies the same rule to recover a character. A **glyph** is the visible shape in a font; if the font lacks it, the bytes and encoding can be correct while the screen shows an empty box.

Check three things in order. (1) Compare `Ә` and `A`: they are different letters. (2) Look up their numbers: U+04D8 and U+0041. (3) Save `Ә` as UTF-8 and inspect its bytes, `D3 98`; Latin `A` uses `41`. If text opens as nonsense, check the decoding first. If the bytes are right but a box appears, check the font. Do not call every display problem “bad Unicode.”

Add a fictional task with a Kazakh letter to `tasks.json`, close the file, and reopen it. If the letter survives, writing and reading agree. State UTF-8 explicitly in `FORMAT.md` so another learner need not guess.

## Where the analogy ends

“One character is one byte” is false. A Unicode code point is not already a byte sequence: an encoding must be chosen for storage. Even one visible character can involve several code points, such as a letter followed by a separate accent mark.

## The lesson's support signal

`visible mark → Unicode code point → UTF-8 encoding → bytes`

## Worked example

Record three rows: A / U+0041 / 41; Я / U+042F / D0 AF; Ә / U+04D8 / D3 98. The U+ value is a code point; the last value gives UTF-8 bytes. Read backwards: D3 98 in UTF-8 represents Ә.

## Predict before observing

How many UTF-8 bytes does “AӘ” require? Count two characters first, then add 1 + 2 = 3 bytes. Do not include the file name or header.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Without the table, explain why Я need not use the same number of bytes as A. Then check both concrete values against the table.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“U+04D8 is the bytes D4 D8.” Correct the confusion between code point and encoding: UTF-8 stores Ә as D3 98. Name the encoding explicitly.

## Transfer to a new situation

A friend saves Kazakh text as UTF-8, but another program reads it using a different encoding. Predict what might appear. Check the encoding on both sides, reopen the file, and compare it with the original text.

## Project change

Create a test file with task title “Әлем” and the line “Я и A.” State in FORMAT.md that project files use UTF-8. Close and reopen the file; verify all three letters. Use no real personal data.

## Exercise

Make a table for A, Я, Ә: code point, UTF-8 bytes, and byte count. Explain why counting visible characters cannot always give the exact text size.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow reconstruct the path from mark to bytes. In a week explain the difference between U+04D8 and D3 98. In a month check that a file survives a move between two programs.

[Next lesson: A photograph as a grid of numbers](/read/informatics-14-pixels-color?lang=en)
