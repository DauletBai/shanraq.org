# Ubuntu server: SSH keys, updates, and a firewall

_Lead (summary):_ **Connect as an unprivileged user, close unnecessary ports, and verify access from a second session.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **need → region/size/access → price alert → create → observe → delete**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Connect as an unprivileged user, close unnecessary ports, and verify access from a second session. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

When moving into a flat, you test the new key while the old door remains open, then change the lock fully. On a server, create a user and verify a second login before restricting root and the firewall, or you may lock yourself out.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**provider console → user+public key → second login → firewall → updates → app**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **SSH** — A secure protocol for remote login and command execution.
- **Key pair** — The private key stays with you; the public key may be placed on the server.
- **Firewall** — Rules controlling which inbound and outbound connections are allowed.
- **Patch** — An update that fixes defects or vulnerabilities.

## Take it apart without rushing

Harden a server in a safe order: create a user and key, verify a second login, then disable risky access. Enabling a firewall or disabling root too early can lock you out. The provider console remains an emergency path.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
ssh-keygen -t ed25519 -a 64 -f ~/.ssh/cloudlab_ed25519 -C cloudlab
printf 'Copy only cloudlab_ed25519.pub to the server.\n'
printf 'On Ubuntu: allow OpenSSH before enabling ufw.\n'
printf 'Keep the current session open while testing a second login.\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

`ssh-keygen` creates a pair. Only the `.pub` file goes to the server; the private file stays with mode 600. The first session deliberately remains open until a second login is proven.

## What you should observe

- The private key is never copied to the server.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 20](/static/course/cloud-devops/map-20-server-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “The private key must be copied to the server so it recognises me.” Correction: the server stores the public key; never transfer the private key.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a local POSIX sh script that checks for private and public key files, requires mode 600 on the private key, and prints only the public-key fingerprint. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Domain, reverse proxy, and automatic HTTPS**.

[Course contents](/course/cloud-devops?lang=en)
