# How sound and motion become data

_Lead (summary):_ **Understand how a microphone measures sound over time and how video shows a sequence of frames.**

## Where we are on the map

The previous block showed how to find files and separate a program from its data. Now we ask how data can carry meaning that another person and another program will read in the same way. State your first guess, test it with numbers, and record what you had to revise.

![How sound and motion become data](/static/course/informatics/map-15-sound-video-sampling-en.svg)

## Everyday analogy and exact model

If you record a river’s water level once an hour, the table shows separate measurements, not every change between them. Digitizing sound works similarly: the changing signal is measured many times each second. Video represents motion with a sequence of separate images.

An audio sample is a numerical measurement of the signal at one moment. Sample rate tells how many samples per second are taken for each channel. In a toy example at 4 samples/s, one second has sample times 0, 0.25, 0.5, 0.75 s: exactly four, not five. At 8000 samples/s, one channel, and 8 bits per sample, one second of raw data uses 8000 bytes. Video stores frames at a stated rate and normally carries sound separately.

## Where the analogy ends

A river table is not the river. The original wave between samples is not directly recorded; a reconstruction approximates it. Four samples/s is only a teaching diagram, not a rate for good speech. Headers and compression change actual file size.

## The lesson's support signal

`time → sample → number; rate × duration × bits/sample × channels → raw bits`

## Worked example

Use four measurements 2, 5, 3, 0 at 0, 0.25, 0.5, 0.75 s. Draw four points in the right order. A two-second clip at 4 samples/s with one channel has eight numbers.

## Predict before observing

How many raw bytes are in one second of mono audio at 8000 samples/s and 8 bits per sample? Write 8000 × 8 = 64000 bits, then divide by eight: 8000 bytes.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Hide the worked example and calculate two seconds under the same conditions. Explain why doubling the number of channels doubles raw data if everything else stays fixed.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“At 4 samples/s the points at 0, 0.25, 0.5, 0.75, and 1 s give five samples in the first second.” Explain the half-open interval [0,1): 1 s starts the next second.

## Transfer to a new situation

A toy video has 2 frames/s. Name six frame times over three seconds: 0, 0.5, 1, 1.5, 2, 2.5 s. Sound does not turn into frames simply because it accompanies video.

## Project change

Add to the assistant specification: if it later has a voice prompt, its format must say the sample rate, channel count, and sample representation. For now, store task text only; audio is unnecessary.

## Exercise

Draw two timelines: four sound samples per second and two video frames per second. Label the times and compare counts without confusing sound with images.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow recall all four sample times. In a week count mono samples over 2 s at 8000/s. In a month explain what the original signal between measurements does not directly preserve.

[Next lesson: Why ZIP and a messenger compress differently](/read/informatics-16-compression?lang=en)
