# A container: a bounded process, not a tiny VM

_Lead (summary):_ **Run a first container, observe namespaces and an immutable image, and distinguish a process from a virtual machine.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **person → address → service → response; alive ≠ ready ≠ useful**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Run a first container, observe namespaces and an immutable image, and distinguish a process from a virtual machine. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A restaurant kitchen is shared, but each cook has a station, ingredients, and space limit. A container resembles that station: it shares the Linux kernel but sees a bounded environment. A virtual machine is closer to a separate kitchen.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**container = process + boundaries; VM = separate guest OS**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Container** — An isolated process with defined files, networking, and limits.
- **Image** — An immutable file template used to start a container.
- **VM** — A virtual computer with its own guest operating system.
- **Namespace** — A Linux mechanism limiting what a process can see.

## Take it apart without rushing

A container shares the host kernel while isolating process, network, and filesystem views. An image is a template; a container is one running instance. Removing a container must not equal losing important data.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
docker version
docker run --rm alpine:3.23 cat /etc/os-release
docker ps -a
docker image ls alpine:3.23
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The commands start a known image, show its process from inside and outside, then remove the learning container. Compare the PID and files: this is a process viewed through different boundaries.

## What you should observe

- `--rm` removes the learning container after exit.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 7](/static/course/cloud-devops/map-07-containers-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “A container is a lightweight VM.” Correction: a container normally shares the host kernel; a VM runs a separate guest OS.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write commands that run `alpine:3.23`, print its hostname, and remove the container automatically after exit. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Image layers, tags, and the container lifecycle**.

[Course contents](/course/cloud-devops?lang=en)
