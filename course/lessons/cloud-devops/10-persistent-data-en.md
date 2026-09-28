# Volumes, bind mounts, and proof of persistence

_Lead (summary):_ **Move state out of the container, recreate it, and prove a note survives process replacement.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **source → test → build → copy binary → non-root runtime**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Move state out of the container, recreate it, and prove a note survives process replacement. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A hotel room loses a guest's temporary items at checkout, while luggage in storage remains. A container filesystem resembles the room; a volume resembles separate storage. Even that storage does not replace a backup elsewhere.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**container changed; volume remained; restore proved data**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Volume** — Storage with a lifecycle separate from a container.
- **Bind mount** — A direct mount of a specific host directory.
- **State** — Data that must survive process restart or replacement.
- **Backup** — A separate copy from which restoration has been tested.

## Take it apart without rushing

A container should be replaceable while state has an explicit lifecycle. A named volume suits Docker-managed data; a bind mount suits a known host path. Neither one is a backup by itself.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
docker volume create cloudlab-practice
docker run -d --name cloudlab-data -p 8080:8080 -v cloudlab-practice:/data cloudlab:lesson-09
sleep 1
curl -fsS -X POST -d 'note=survives-recreate' http://127.0.0.1:8080/notes >/dev/null
docker rm -f cloudlab-data
docker run -d --name cloudlab-data -p 8080:8080 -v cloudlab-practice:/data cloudlab:lesson-09
sleep 1
curl -fsS http://127.0.0.1:8080/ | grep survives-recreate
docker rm -f cloudlab-data
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The exercise writes a note, recreates the container, and reads it again. This proves separation of data from process. Backup and restore prove another path; merely possessing an archive is insufficient.

## What you should observe

- The note remains after the container is removed.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 10](/static/course/cloud-devops/map-10-persistent-data-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “A volume is a backup.” Correction: it outlives a container, but deletion, disk failure, or data corruption can destroy it too.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write commands that create a named volume, write a file through one temporary container, and read it through a second container. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Docker Compose: services, networking, and declarative startup**.

[Course contents](/course/cloud-devops?lang=en)
