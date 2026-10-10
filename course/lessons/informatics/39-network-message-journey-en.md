# The journey of a message between two phones

_Lead (summary):_ **Follow one fictional task message from a phone to a classroom server and identify every stage of the journey.**

## Where we are on the map

This is lesson 39 of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.

![The journey of a message between two phones](/static/course/informatics/map-39-network-message-journey-en.svg)

## Situation and question

Imagine opening the assistant on a phone while a friend uses the same address on a laptop. The page appears quickly, so it feels as though the message jumped directly between screens. On paper, draw four stops: device, home router, provider network, server. Swap the phone and server in your drawing. Which side starts a request, and which side sends a response?

## New words without gaps

A **network** is a set of connected devices and links that carry data. A **node** is a participating device. A **client** starts a request; a **server** responds. A **provider** connects your home network to other networks. **Latency** is the time between sending and receiving. A letter moving through post offices is a useful image, but data is not a paper envelope and the return route need not be identical.

## The lesson's support signal

Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.

## Work through it step by step

First, the phone uses Wi-Fi to reach the home router. The router passes data to the provider and onward through other networks. The server receives the request and builds a response. The browser displays what comes back. On your drawing, put `tasks.json` on the server side only in the imagined remote version. In our real classroom release, the file sits on the same computer that runs localhost. Switching off Wi-Fi does not necessarily stop access to 127.0.0.1 on that computer: the request never reaches a provider.

## Predict and check

Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.

```python
hops = ["phone", "router", "provider", "server"]
print(" → ".join(hops))
```

## Expected output

```text
phone → router → provider → server
```

## Catch the error

A common mistake is treating the four drawn stops as the exact physical route for every site. There may be more networks between provider and server; a local request does not leave the computer at all. Another mistake is equating page loading with sending personal data. Our sample uses fictional records only.

## Project change

Keep two labelled paper routes: an imagined remote site and the real local view at `http://127.0.0.1:8765/`. Mark where `tasks.json` lives in each. Do not alter the program yet; this distinction sets the boundary for later security lessons.

## Task and evidence

Without the page, draw a request and response path for both localhost and a remote school site. Identify the data source, two causes of delay, and a broken link that would stop the remote response. Explain to a classmate why a local page can still work without internet access.

## Transfer to a new setting

If a classroom printer is reachable only inside the school network, which parts of your route remain? Do not call it a server without checking whether it accepts requests and produces replies itself.

## Return after 1, 7, and 30 days

After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.

## Primary reference to check

[Official documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview)

## Next lesson

[IP addresses, routers, and packets](/read/informatics-40-ip-router-packets?lang=en)
