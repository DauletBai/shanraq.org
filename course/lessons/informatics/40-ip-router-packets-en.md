# IP addresses, routers, and packets

_Lead (summary):_ **Split a message into pieces, locate the destination address, and restore order without guessing.**

## Where we are on the map

This is lesson 40 of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.

![IP addresses, routers, and packets](/static/course/informatics/map-40-ip-router-packets-en.svg)

## Situation and question

You send a friend three numbered cards spelling a task title. A courier brings cards 2, 1, 3. Your friend must use the numbers rather than read them in arrival order. Real networks also need to distinguish where data is going from how an application should interpret what arrives.

## New words without gaps

An **IP address** identifies a network interface in a network; it is not a person's name or an apartment number. A **router** examines a destination and chooses the next part of a route. A **packet** carries a portion of data plus control information. The card numbers in our exercise belong to the teaching model: IP alone does not promise order, retransmission, or completeness. A **port**, introduced later, distinguishes services on one device: the address identifies a host, and a port identifies the destination process.

## The lesson's support signal

Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.

## Work through it step by step

Draw sender 192.0.2.10 and recipient 203.0.113.7 with a router between them. These are documentation addresses, not actual assistant servers. The router consults its forwarding table to send a packet toward 203.0.113.7. If one card goes missing, IP alone does not tell you it will be resent. If parts arrive 2, 1, 3, the application should not blindly display a scrambled title. In our exercise, the numbers are application labels; a later lesson introduces TCP's different way of ordering a byte stream.

## Predict and check

Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.

```python
parts = {2: "AS", 1: "T", 3: "K"}
print("".join(parts[n] for n in sorted(parts)))
```

## Expected output

```text
TASK
```

## Catch the error

Do not say an IP address uniquely identifies a human. Devices may share a public address, addresses may change, and servers may have several. Do not promise IP will deliver every piece: it helps address and forward packets but does not itself guarantee reliable delivery.

## Project change

Add three columns to the web project passport: `assistant.test`, teaching IP 203.0.113.7, and local address 127.0.0.1. The last refers to the same computer. The other two are models in this lesson; we are not running a public server under a fictional name.

## Task and evidence

Cards 3 and 1 arrive, but card 2 does not. Show what is known and explain why the message cannot yet be called complete. Name the destination address and next hop in an invented forwarding table.

## Transfer to a new setting

Two phones behind one home router open the same site. Why might the site see one public address even though two devices are involved? What other observations would be needed before inferring the number of people?

## Return after 1, 7, and 30 days

After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.

## Primary reference to check

[Official documentation](https://www.rfc-editor.org/rfc/rfc791)

## Next lesson

[DNS: how a name leads to a server](/read/informatics-41-dns-names?lang=en)
