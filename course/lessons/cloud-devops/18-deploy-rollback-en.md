# Deployment, smoke testing, and a verified rollback

_Lead (summary):_ **Deploy by version, test from outside, and restore the previous digest with a prewritten procedure.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **config may be shown | secret is restricted | both enter outside image**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Deploy by version, test from outside, and restore the previous digest with a prewritten procedure. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A plumber does not leave merely because a tap was installed. They run water, check for leaks, and know how to restore the old part. Deploy installs a version, a smoke test checks the main path, observation finds consequences, and rollback restores a known version.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**deploy X → smoke test → observe → keep X OR rollback Y → verify**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Deploy** — Installing a specific version into a specific environment.
- **Smoke test** — A short check of a critical user path after a change.
- **Rollback** — A return to a previously known version or state.
- **Migration** — A controlled change to data structure or meaning.

## Take it apart without rushing

A rollback must exist before an incident. It restores code but may not restore compatible data; migrations need a separate plan. After switching, test a user path and watch metrics rather than treating a successful exit code as proof.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
export CLOUDLAB_IMAGE=ghcr.io/OWNER/cloudlab:REPLACE_WITH_SHA
export CLOUDLAB_DOMAIN=cloudlab.example.com
docker compose -f compose.prod.yaml config >/tmp/cloudlab-prod.yml
grep -n 'image:[|]CLOUDLAB_ENV' /tmp/cloudlab-prod.yml
printf 'record previous digest before docker compose up -d\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The scenario records old and new digests, switches version, runs a smoke test, and observes a signal. Rollback is ready only after rehearsal. Returning code may not be compatible with data already changed.

## What you should observe

- The configuration does not use latest.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 18](/static/course/cloud-devops/map-18-deploy-rollback-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “A successful deploy exit code means a successful release.” Correction: test a user path and observe the system after the change.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh deploy function that saves the current digest, applies a new image, runs a smoke test, and restores the saved image if that test fails. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Cloud: IaaS, cost, IAM, and shared responsibility**.

[Course contents](/course/cloud-devops?lang=en)
