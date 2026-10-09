# A photograph as a grid of numbers

_Lead (summary):_ **Break a photograph into pixels and check how numbers can describe the color at each position.**

## Where we are on the map

A task photo was enlarged tenfold and its edges became blocky. Did we gain detail, or merely enlarge the existing little squares?

![A photograph as a grid of numbers](/static/course/informatics/map-14-pixels-color-en.svg)

## Everyday analogy and exact model

Stand back from a mosaic and you see a picture. Come closer and you see its pieces. A raster image also has separate positions called pixels; the squares in our diagram simply make these positions easy to inspect.

Each pixel has a position and a color. In a simple RGB model, three numbers from 0 to 255 describe red, green, and blue channels. Our 2 × 2 grid has red (255,0,0) and green (0,255,0) on top, blue (0,0,255) and white (255,255,255) below. If each channel uses eight bits without compression, color data requires 2 × 2 × 3 = 12 bytes.

## A coloured cell and image size

A **pixel** is one cell of a digital image with an assigned colour. **Image resolution** counts pixels across and down, for example 4 × 3 = 12 pixels; it is not the screen's physical size. The **RGB** model uses three channels: red, green, and blue. With eight bits per channel, each value from 0 to 255 sets an intensity. `(255, 0, 0)` is red, `(0, 0, 0)` black, and `(255, 255, 255)` white. **Scaling** stretches or recalculates existing pixels; it cannot recover details the camera never recorded.

Draw a 4 × 3 grid. Colour the top row red, the middle green, and the bottom blue. There are 12 pixels, not three: each row colour appears in four cells. Change one cell to white and explain which three values change. Now draw an 8 × 6 grid by doubling each old pixel in both directions. You have 48 cells but no new information about the photographed scene.

Design a task icon for the assistant on a 4 × 3 grid and compare its legibility with a more detailed grid. If its meaning depends only on red versus green, add words too: accessibility applies to images.

## Where the analogy ends

A complete file is normally larger than 12 bytes because it needs format and size information; compression can change the size too. A screen pixel and a photograph pixel are not the same physical grain. We are modeling a discrete grid of data.

## The lesson's support signal

`row + column → pixel → (R,G,B); 2 × 2 × 3 = 12 bytes excluding headers`

## Worked example

Read the lower-left square in the 2 × 2 grid: blue (0,0,255). Change only that square to (255,255,255): the grid then has two white positions, and the other three positions stay as they were.

## Predict before observing

How many raw color bytes does a 3 × 2 grid with three eight-bit channels need? Find six pixels and multiply by three: 18. Do not call that the actual PNG file size.

Do not jump straight to the answer. Write your prediction and its reason first. Compare each intermediate step as well as the final result. If you were right by chance, repeat with different values.

## Reconstruct without a hint

Draw the 2 × 2 grid with all four colors, hide the diagram, and label its positions. Explain what changes if the top two pixels trade places.

Cover the worked example with paper. Reconstruct the chain from the short support signal, explain every transition aloud, and then reveal the example to check yourself. If stuck, look back at only the preceding step.

## Find and correct the mistake

“A 2 × 2 image file is always four bytes.” Point out three channels per pixel and the file’s metadata. The corrected raw RGB color count is 12 bytes.

## Transfer to a new situation

For a 4 × 4 assistant icon, calculate raw RGB data: 16 × 3 = 48 bytes. Why do a blurry photograph and a sharp icon of that same size have the same uncompressed amount in this model?

## Project change

Draw a 2 × 2 assistant status icon and write its four RGB triples row by row. Add a rule to FORMAT.md: width and height first, then colors from left to right in each row.

## Exercise

Change exactly one icon pixel. Record its position and old/new RGB values. Ask someone to reconstruct both images from the numbers and compare them.

Submit evidence that can be checked, rather than saying “I understood”: a table, calculation, file, or precise answer with a reason. Use fictional data. Ask another learner to repeat the action from your description; if they must guess, clarify the rule.

## Return after 1, 7, and 30 days

Tomorrow recall RGB order. In a week calculate raw bytes for a 5 × 2 grid. In a month explain why file size is not just the number of pixels.

[Next lesson: How sound and motion become data](/read/informatics-15-sound-video-sampling?lang=en)
