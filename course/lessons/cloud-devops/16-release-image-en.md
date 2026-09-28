# Build and publish one verified image

_Lead (summary):_ **Publish to a registry only after tests and promote one digest across environments.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **event → workflow → job → runner → steps → evidence**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Publish to a registry only after tests and promote one digest across environments. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A medicine sample is tested, and that exact batch is distributed. If a new batch is made for each pharmacy, the earlier test no longer proves its contents. Likewise, one image is built once and promoted across environments by digest.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**test once → build once → sign/scan → push digest → promote same digest**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Artifact** — A built pipeline output such as a binary, image, or package.
- **Release** — An identified set of changes and artifacts ready for use.
- **Promotion** — Moving the same artifact to the next environment without rebuilding.
- **Provenance** — A record of where and how an artifact was produced.

## Take it apart without rushing

Rebuilding for production may produce an artifact different from the one tested. Build once, verify, publish, and promote its digest. Registry credentials belong only in the publishing job and must not reach an untrusted fork pull request.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
sha=$(git rev-parse --short=12 HEAD)
image=ghcr.io/OWNER/cloudlab:$sha
printf 'would publish %s\n' "$image"
printf 'production must record a sha256 digest, not latest\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The exercise builds an image, ties it to a commit SHA, and extracts its digest. Record that digest in the release log. The next stage must receive it, not newly built local content under the same tag.

## What you should observe

- The name contains an immutable version.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 16](/static/course/cloud-devops/map-16-release-image-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “Rebuilding the same commit in production is identical.” Correction: dependencies, base images, or builders may change; promote the verified digest.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh check for `REGISTRY`, `IMAGE`, and `GIT_SHA` that constructs the full image name and refuses to publish the tag `latest`. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Secrets, configuration, and environments without leaks**.

[Course contents](/course/cloud-devops?lang=en)
