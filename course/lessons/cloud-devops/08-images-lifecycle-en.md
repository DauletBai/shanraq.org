# Image layers, tags, and the container lifecycle

_Lead (summary):_ **Inspect image history, pin a digest, and distinguish stop, start, remove, and pull.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **container = process + boundaries; VM = separate guest OS**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Inspect image history, pin a digest, and distinguish stop, start, remove, and pull. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

An image resembles a sealed construction kit; a container is a model assembled from it. A tag is a label that can be moved. A digest fingerprints the content: change one part and the fingerprint changes.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**Dockerfile → layers → image@digest → container → stop/remove**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Layer** — A saved filesystem change in an image.
- **Tag** — A convenient movable version name such as `1.0`.
- **Digest** — A cryptographic identifier for exact image content.
- **Lifecycle** — Creation, start, stop, and removal as distinct states.

## Take it apart without rushing

A tag is a movable name; a digest identifies exact content. `latest` is insufficient for a reproducible release because it may point to different bytes tomorrow. Containers and images also have separate lifecycles.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
docker pull alpine:3.23
docker image inspect alpine:3.23 --format '{{index .RepoDigests 0}}'
docker history alpine:3.23
docker create --name cloudlab-inspect alpine:3.23 true
docker start -a cloudlab-inspect
docker rm cloudlab-inspect
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

`image history` exposes layers, `inspect` shows the digest, and `run`, `stop`, `start`, and `rm` demonstrate the difference between an image and an instance. After `stop` the container exists; after `rm` it does not.

## What you should observe

- The digest contains sha256.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 8](/static/course/cloud-devops/map-08-images-lifecycle-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “The `latest` tag always means the newest immutable image.” Correction: it is a movable label; a digest gives exact identity.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write commands that pull an image, save its first RepoDigest in a variable, and stop if the digest is empty. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Dockerfile: small, tested, and unprivileged**.

[Course contents](/course/cloud-devops?lang=en)
