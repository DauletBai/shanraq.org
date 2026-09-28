# Why Cloud & DevOps now: infrastructure, roles, and boundaries

_Lead (summary):_ **Separate data-centre growth from job-title growth and map the skills of cloud engineering, DevOps, platform engineering, and SRE.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

This is the course's first support. Seeing the whole path matters more now than memorising tool names.

## Lesson outcome

Separate data-centre growth from job-title growth and map the skills of cloud engineering, DevOps, platform engineering, and SRE. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

Imagine a modern hotel. The building, power, and lifts resemble a data centre. Booking a room through an app resembles cloud capacity. The team that checks guests in, spots a broken lift, and restores it resembles service operations. DevOps connects people who change a service with people responsible for running it.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**🏭 hardware → ☁️ resource → 🔁 change → 👀 observe → 🛟 recover**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Data centre** — A building containing servers, networking, power, and cooling.
- **Cloud** — A way to request computing resources through software or a control panel without buying every server.
- **DevOps** — A practice of shared work, automation, and fast feedback when releasing changes.
- **SRE** — A reliability approach that measures service goals and automates repeated manual work.

## Take it apart without rushing

DevOps is neither one tool nor one person responsible for everything. It reduces change risk through automation, short feedback loops, and shared responsibility. A data centre supplies physical capacity, cloud turns it into programmable resources, and operations keeps a service available.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
printf '%s\n' 'change -> test -> image -> deploy -> observe -> recover'
printf '%s\n' 'evidence: repeatable command, metric, backup, recovery record'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The two `printf` lines install nothing. They provide the first course map. Read the chain left to right, cover it, and rebuild it from memory.

## What you should observe

- The chain begins with a change and ends with verified recovery.
- Every stage leaves observable evidence.
- You can state the course boundary: software operations.

## Rebuild the whole from the support map

![Support map 1](/static/course/cloud-devops/map-01-why-now-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “DevOps is an administrator solely responsible for Docker, servers, and every outage.” Correction: the team shares responsibility; tools support the process.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh script that prints the six release-path stages on separate lines and stops on the first error. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Terminal and Git: a reproducible working point**.

[Course contents](/course/cloud-devops?lang=en)
