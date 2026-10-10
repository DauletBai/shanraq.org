# DNS: how a name leads to a server

_Lead (summary):_ **Find out how a site name leads to an address and why DNS does not prove the site's identity.**

## Where we are on the map

This is lesson 41 of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.

![DNS: how a name leads to a server](/static/course/informatics/map-41-dns-names-en.svg)

## Situation and question

You remember a friend's name but not their current number. An address book can help locate a number, yet it cannot prove who answers the phone. A site has a similar split: people remember names, while a network needs addresses. Examine what happens after entering `assistant.test` and before a connection is attempted.

## New words without gaps

A **domain name** is a readable name such as `assistant.test`. **DNS** is a system for finding records about names, often compared to an address book. An **A record** maps a name to an IPv4 address; an **AAAA record** maps it to IPv6. A **cache** temporarily remembers an answer, and **TTL** sets a normal lifetime for that cached answer. DNS is not identity verification: an answer may be stale or tampered with, and a TLS certificate must be checked separately.

## The lesson's support signal

Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.

## Work through it step by step

Look up `assistant.test → 203.0.113.7` in our classroom table. The address is reserved for documentation and is not a real assistant server; `.test` is reserved for testing. Now imagine the entry changes to 203.0.113.8. One computer still has the old answer cached while another has the new one. Why might their connection attempts differ? Compare a TTL measured in minutes, but do not promise every device updates at exactly the same instant. Our local server needs no such public record: 127.0.0.1 explicitly means the loopback interface on the same computer.

## Predict and check

Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.

```python
records = {"assistant.test": "203.0.113.7"}
print(records["assistant.test"])
```

## Expected output

```text
203.0.113.7
```

## Catch the error

It is wrong to say DNS delivers the page: it usually helps find an address, after which other steps occur. It is equally wrong to treat a correct IP as proof a site is safe. A name and address are distinct clues; remote server trust requires another check.

## Project change

On your future hosting diagram, draw a name-to-IP arrow labelled DNS lookup. Draw separate arrows for connecting and the later TLS check. Label localhost as the classroom release without public DNS. This keeps the learning model from pretending that our project is already live in a cloud.

## Task and evidence

Given `assistant.test → 203.0.113.7` and `library.test → 203.0.113.8`, find the answer for each name. After the first record changes, explain why one cache can still show the old address and why a DNS answer is not a certificate of trust.

## Transfer to a new setting

Someone sends a familiar-looking domain with one extra letter. Even if DNS returns a working IP, what would you check before entering a password? The classroom assistant asks for no password; use another real service for this thought experiment.

## Return after 1, 7, and 30 days

After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.

## Primary reference to check

[Official documentation](https://developer.mozilla.org/en-US/docs/Glossary/DNS)

## Next lesson

[TCP and UDP: different delivery promises](/read/informatics-42-tcp-udp-delivery?lang=en)
