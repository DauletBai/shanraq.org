# TLS: encryption and checking the other party

_Lead (summary):_ **Separate encryption from server identity so a padlock is not mistaken for proof that a site is honest.**

## Where we are on the map

This is lesson 43 of 72 in the Internet, web, and cloud block (39–46). We continue one project: the three fictional records and five reminder outcomes from version 1.0 must survive. The support map shows the exact transition we can explain and test; an arrow does not claim that every intermediate device is drawn.

![TLS: encryption and checking the other party](/static/course/informatics/map-43-tls-trust-en.svg)

## Situation and question

Someone hands you a sealed envelope. A seal can deter reading or altering the letter, but it does not tell you whether you addressed it to the right person. TLS helps solve related problems in the browser. Name two risks separately: an observer on the route and a response from the wrong server.

## New words without gaps

**TLS** is a protocol for protecting a connection, often used with HTTPS. **Encryption** makes transmitted content unreadable to an ordinary observer on the network. A **certificate** binds a public key to a server name; a browser checks its trust chain and the name match. A **public key** can be shared, while its matching **private key** must remain with its owner. **HTTPS** is HTTP inside a protected connection. A padlock does not prove that a site's claims are true or control what the site itself does with received data.

## The lesson's support signal

Read the map from left to right. For every arrow, name its input, action, and verifiable result. Cover the next card's label and predict it in your own words. Return to the real example: distinguish a teaching analogy, localhost behaviour, and a possible remote site. If the picture promises more than the explanation, correct the explanation or the drawing.

## Work through it step by step

Suppose you enter `https://assistant.test/`. DNS supplies an address. The client then contacts the server and checks that its certificate is valid for `assistant.test` and chains to a trusted authority. The parties agree on secrets for protected communication; the HTTP request follows over that connection. If the certificate names a different site, continuing as if nothing happened is unsafe. `assistant.test` is fictional here: our project has no such certificate. The actual `step-04` classroom view uses `http://127.0.0.1:8765/`, without TLS, and listens only on the same computer.

## Predict and check

Do not run the code immediately. Write its exact expected result, including line order and characters. Explain what it proves and what it does not: this small example models one rule but does not create the internet, a TLS certificate, or a public server. Then run it from `step-04` and compare character by character. Change one input, predict the new result first, and only then test it.

```python
from urllib.parse import urlsplit
address = urlsplit("https://assistant.test/tasks")
print(address.scheme, address.hostname)
```

## Expected output

```text
https assistant.test
```

## Catch the error

Do not say DNS itself verifies the site owner, or that `https` guarantees that a shop is honest. TLS protects the channel and helps verify the server name under a trust model. If a user willingly sends a password to a deceptive domain that has its own valid certificate, a padlock cannot undo that choice.

## Project change

Use three different colours on your diagram: DNS finds an address; TCP or another transport moves data; TLS checks the remote party and protects content. Beside our local classroom server write 'loopback only, HTTP, no TLS'. Public hosting would require a different architecture that this block does not claim to provide.

## Task and evidence

Compare two cases: a certificate valid for `library.test` while you opened `assistant.test`; and a certificate matching its name, but the page asks for a suspicious money transfer. Which case fails the name check, and which requires evaluating the page and address themselves?

## Transfer to a new setting

The school Wi-Fi triggers a browser certificate warning. Why is 'just click through' a poor rule for every student? What information should a teacher gather before anyone enters credentials?

## Return after 1, 7, and 30 days

After 1 day, redraw the main transition from memory and give a counterexample to an oversimplified analogy. After 7 days, explain a fresh case to a classmate without reading the page and answer at least 7 of 10 questions. After 30 days, repeat the local web experiment and verify the three original tasks and five rules survive. Write down any term still unclear and return to its first explanation.

## Primary reference to check

[Official documentation](https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Transport_Layer_Security)

## Next lesson

[HTTP: a conversation between browser and server](/read/informatics-44-http-browser-server?lang=en)
