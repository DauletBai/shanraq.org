# TCP and UDP: different delivery promises

_Lead (summary):_ **Compare two delivery models and see why speed alone does not decide which one a task needs.**

## Where we are on the map

This is lesson 42 of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.

![TCP and UDP: different delivery promises](/static/course/informatics/map-42-tcp-udp-delivery-en.svg)

## Situation and question

A school club dictates three steps over a call: open, record, verify. If step two disappears, step three cannot repair the instruction. During a video call, however, skipping one lost frame may be better than waiting for it to be resent. What promises does each situation require?

## New words without gaps

**TCP** provides an ordered, reliable byte stream between endpoints of a connection; acknowledgements and retransmission help recover losses. **UDP** sends individual datagrams without built-in delivery or ordering guarantees. A **datagram** is a separate transport message. A **stream** is a sequence of bytes: TCP does not preserve an application's message boundaries by itself. Neither transport can fix a mistake in the task content itself.

## The lesson's support signal

Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.

## Work through it step by step

Imagine cards 1, 2, 3. In our teaching model of TCP, the receiver does not display a complete text while card 2 is missing; the sender learns of the loss and resends it. In a simple UDP model, 1 and 3 can arrive without 2, and the application decides whether to wait, ask again, or ignore the loss. Actual TCP numbers bytes, not our cards. Modern HTTP/3 can use QUIC over UDP; that does not leave its data without recovery, because QUIC implements needed guarantees above UDP.

## Predict and check

Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.

```python
received = {3: "done", 1: "task", 2: ":"}
print("".join(received[n] for n in sorted(received)))
```

## Expected output

```text
task:done
```

## Catch the error

The slogan 'UDP is always faster than TCP' is not a universal fact. Results depend on network conditions, the protocol above and the goal. Another mistake is treating one TCP send call as exactly one read on the other side. A byte stream may arrive in pieces or combined; the application must define its own data boundaries.

## Project change

Record a decision in the web assistant passport: a partial HTML page or task JSON must not be presented as a correct complete result. The classroom view uses ordinary HTTP over TCP on 127.0.0.1. Do not call it HTTPS or QUIC; those are not present in this release.

## Task and evidence

Cards 1 and 3 arrive but 2 is missing. Show the teaching-model TCP response and the simple UDP response. State an application rule for a whole JSON document: when may tasks be shown, and when should an incomplete response be reported?

## Transfer to a new setting

Choose a useful loss policy for a live video call and for transferring money. Justify it by the consequences of an error rather than a slogan about speed.

## Return after 1, 7, and 30 days

After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.

## Primary reference to check

[Official documentation](https://www.rfc-editor.org/info/rfc9293)

## Next lesson

[TLS: encryption and checking the other party](/read/informatics-43-tls-trust?lang=en)
