# How binary notation stores numbers

_Lead (summary):_ **See why 1101₂ means thirteen and check binary numbers by place value instead of guessing.**

## Where we are on the map

The screen shows `1101`. Is that thirteen or one thousand one hundred and one? You need to know the number system first.

![How binary notation stores numbers](/static/course/informatics/map-12-binary-numbers-en.svg)

## Everyday analogy and exact model

In decimal notation, position changes a digit’s weight: the 1 in 13 means one ten. Binary notation also has positions, but their weights from right to left are 1, 2, 4, and 8. Imagine a set of weights where each next piece is twice as heavy.

The positions of 1101₂ have weights 8, 4, 2, and 1. Add the weights under ones: 8 + 4 + 0 + 1 = 13₁₀. Four places represent whole numbers from 0 through 15. A leading zero does not change the value: 01101₂ is still 13, although five positions are written. State the number base; otherwise 10 is ambiguous.

## A digit's place changes its value

A **number system** defines allowed digits and the values of places. Its **base** is the number of distinct digits: decimal has 10 and binary has 2. A **place** is a digit position; binary places from right to left have weights 1, 2, 4, 8, 16, and so on. Thus `1101₂` = 1 × 8 + 1 × 4 + 0 × 2 + 1 × 1 = 13₁₀. The subscript clarifies the base. The computer does not find such a subscript in every byte; the format tells a reader how to interpret the data.

Lay cards labelled 8, 4, 2, 1 on a table. To make 13, choose 8 + 4 + 1 and place `1 1 0 1` above them. To make 10, choose 8 + 2: `1010`. Hide the example and convert 9 yourself. Check by adding backwards: `1001₂` = 8 + 1. The mistake `1101₂ = 1 + 1 + 0 + 1 = 3` ignores place values.

In the project, do not turn task ID `t-01` into a number: it is a text label. Store the count of completed tasks in a separate numeric field. Similar-looking symbols can have different meanings under different type rules.

## Where the analogy ends

Weights help with arithmetic, but the computer does not store a number as four physical metal pieces. Bits form a sequence; an agreed interpretation tells us whether to read it as an integer, text, or part of an image.

## The lesson's support signal

`1101₂ → 8 + 4 + 0 + 1 → 13₁₀; 10₂ = 2₁₀`

## Worked example

For 9₁₀ take 8 and 1 to get 1001₂. For 6₁₀ take 4 and 2 to get 0110₂ in four places. Check each answer by adding the weights in positions containing 1.

## Predict before observing

What is 1010₂? Mark weights 8, 4, 2, 1 and circle selected ones before calculating. The answer is 8 + 2 = 10₁₀. In contrast, 10₂ means 2₁₀.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Write 7₁₀ in four places without looking. Then write 15₁₀. If you make a mistake, check one position at a time from the heaviest to the lightest.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“1010₂ = 1010₁₀ because the digits match.” Identify the mixed-up bases and recalculate using place weights. The same symbols can name different numbers in different bases.

## Transfer to a new situation

A panel shows 0011₂. Find its value and explain why four places may matter for field width although the first two zeros do not change the number.

## Project change

Add an example numerical id 13 to the project format note: 00001101₂ in eight places. Explain that a quoted “13” in JSON is a text identifier, not merely two binary places.

## Exercise

Build a three-place binary table for every number from 0 to 7. Check 5 and 6 by adding weights. Explain why 8 cannot fit into this table.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow reconstruct weights 1, 2, 4, 8. In a week convert 11₁₀ to binary. In a month explain the difference between 10₂ and 10₁₀ to a peer.

[Next lesson: How a computer distinguishes Ә, Я, and A](/read/informatics-13-text-unicode?lang=en)
