# Registry, immutable tags, SBOM, and supply-chain checks

_Lead (summary):_ **Tie an image to a commit SHA, produce an SBOM, and scan its contents before publishing.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **liveness: restart? | readiness: send traffic? | limits: how much?**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Tie an image to a commit SHA, produce an SBOM, and scan its contents before publishing. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A warehouse accepts boxes, but storage alone does not prove their contents are good. A registry is an image warehouse. A digest resembles a seal for one exact box, an SBOM is the parts list, and a scanner compares that list with known problems.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**commit → image → digest → registry → SBOM/scan → deploy**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Registry** — A service that stores and distributes container images.
- **Immutable** — Unchangeable after publication; a new version gets a new identity.
- **SBOM** — A list of software components and versions inside an artifact.
- **Vulnerability** — A known weakness whose applicability must be checked in context.

## Take it apart without rushing

A registry stores and distributes images, but an upload does not create trust. A release needs an immutable version, traceable origin, and visible dependencies. An SBOM lists components; a scanner matches known vulnerabilities but cannot prove that no bug exists.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
sha=$(git rev-parse --short=12 HEAD)
docker build --build-arg VERSION="$sha" -t "cloudlab:$sha" .
docker image inspect "cloudlab:$sha" --format '{{.Id}}'
printf 'release=%s\n' "$sha"
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The commands name an image with a short SHA, obtain its digest, and produce an inventory. If a scanner is absent, the lesson does not pretend success: understand the expected artifact and install the tool separately.

## What you should observe

- The tag equals the source commit SHA.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 13](/static/course/cloud-devops/map-13-registry-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “No CVE was found, so the image is safe.” Correction: databases are incomplete, and configuration or your own code may still be vulnerable.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh script that derives a tag from the 12-character Git SHA, builds the image, and prints its local image ID; do not use `latest`. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **CI/CD as a feedback system**.

[Course contents](/course/cloud-devops?lang=en)
