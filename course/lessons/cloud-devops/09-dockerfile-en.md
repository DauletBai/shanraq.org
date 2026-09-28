# Dockerfile: small, tested, and unprivileged

_Lead (summary):_ **Build CloudLab as a multi-stage image, test during the build, and verify its unprivileged user.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **Dockerfile → layers → image@digest → container → stop/remove**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Build CloudLab as a multi-stage image, test during the build, and verify its unprivileged user. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A travel case contains only what the trip needs, not the whole workshop. A multi-stage Dockerfile builds in a large workshop, then copies the finished binary into a small clean image.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**source → test → build → copy binary → non-root runtime**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Dockerfile** — A text recipe that builds an image step by step.
- **Build stage** — A temporary environment containing compiler and source.
- **Runtime stage** — The minimal environment in which the finished program runs.
- **Non-root** — Running as an identity without superuser powers.

## Take it apart without rushing

A multi-stage build leaves the compiler in the build stage and copies only binaries into the final image. This reduces size and attack surface. `USER 65532` constrains the process, while a static health check needs no shell in the image.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
docker build --build-arg VERSION=lesson-09 -t cloudlab:lesson-09 .
docker image inspect cloudlab:lesson-09 --format 'user={{.Config.User}} size={{.Size}}'
docker run --rm --entrypoint /cloudlab-healthcheck cloudlab:lesson-09 || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The build runs tests, compiles CloudLab, and carries only required files forward. `inspect` checks the user and a health request proves startup. Small size alone does not prove security.

## What you should observe

- The image builds only after tests pass.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 9](/static/course/cloud-devops/map-09-dockerfile-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “A small image is automatically secure.” Correction: a smaller surface helps, but updates, dependency checks, permissions, and safe code still matter.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write commands to build `cloudlab:practice` and fail via `docker image inspect` when the configured user is `0` or empty. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Volumes, bind mounts, and proof of persistence**.

[Course contents](/course/cloud-devops?lang=en)
