# Networking without magic: IP, port, DNS, HTTP, and TLS

_Lead (summary):_ **Trace a request from a domain name to a process socket and localize failures layer by layer.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **program + start = process(PID); TERM → stop; log → what happened**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Trace a request from a domain name to a process socket and localize failures layer by layer. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

Imagine a large business centre. A domain is the organisation's familiar name. DNS is the directory that finds its building address. An IP is the building's network address. A port is the required office number inside. TCP establishes a conversation with that office. TLS checks the other party's identity and creates a private corridor. HTTP is the language used to request a document.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**name → DNS → IP → TCP:port → TLS → HTTP → response**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **DNS** — The system that translates a domain name into an IP address.
- **IP address** — A numeric network address for a device or interface.
- **Port** — A number from 0 to 65535 that routes a connection to the right program at an address.
- **TCP / HTTP / TLS** — TCP carries an ordered stream; HTTP gives a request meaning; TLS encrypts it and verifies the peer with a certificate.

## Take it apart without rushing

DNS answers which address, TCP whether a port accepts a connection, HTTP whether the application understood a request, and TLS whom the secure connection reached. Test layers in order so a certificate failure does not become a vague “server is down.”

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
getent hosts example.com 2>/dev/null || nslookup example.com
curl -sS -o /dev/null -w 'status=%{http_code} ip=%{remote_ip} tls=%{ssl_verify_result}\n' https://example.com
printf 'GET /healthz HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n' | nc 127.0.0.1 8080 || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

First, `getent` asks DNS. Then `curl` prints the HTTP status, remote IP, and TLS verification result separately. The last command goes directly to a local IP and port. CloudLab starts in the next lesson, so connection refusal is useful evidence here: the address exists, but no process is accepting requests on that port.

## What you should observe

- DNS returns an address independently of HTTP.
- curl shows status, remote IP, and TLS verification.
- You can explain why port 443 is not itself HTTPS.

## Rebuild the whole from the support map

![Support map 5](/static/course/cloud-devops/map-06-networks-dns-http-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “DNS opens the website.” Correction: DNS only supplies an IP; a connection, optional TLS, an HTTP request, and the application must still work.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh diagnostic that accepts a URL, prints its HTTP status and remote IP with `curl`, and fails for status 400 or higher. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **CloudLab: service, environment, and readiness criteria**.

[Course contents](/course/cloud-devops?lang=en)
