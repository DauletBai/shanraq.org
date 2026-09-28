# GitHub Actions: verify every change

_Lead (summary):_ **Move local checks into a workflow, pin permissions, and use caching without storing secrets.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **commit → fast checks → build once → verify → deliver → deploy deliberately**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Move local checks into a workflow, pin permissions, and use caching without storing secrets. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

An automatic car wash starts when a sensor sees a car and runs a program at a specific bay. GitHub Actions is similar: an event starts a workflow, a job receives a clean runner, and steps execute commands in order and leave evidence.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**event → workflow → job → runner → steps → evidence**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Workflow** — A YAML file containing automation triggers, permissions, and jobs.
- **Job** — A set of steps executed on one runner.
- **Runner** — A temporary machine or agent that executes a job.
- **Permissions** — The minimum rights granted to a workflow's temporary token.

## Take it apart without rushing

A workflow is the repository's executable contract. A trigger selects the moment, permissions constrain the token, a job supplies a clean machine, and steps produce evidence. Third-party actions are code too; pin trusted versions and minimize permissions.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
cp workflow.example.yml /tmp/cloudlab-workflow.yml
sed -n '1,220p' /tmp/cloudlab-workflow.yml
grep -n 'permissions:[|]go test[|]build-push-action' /tmp/cloudlab-workflow.yml
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

Read the workflow as a story: when it begins, on which machine, with what rights, and which evidence it creates. Compare its commands with the local preflight so automation is not a secret second build path.

## What you should observe

- The workflow has explicit permissions.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 15](/static/course/cloud-devops/map-15-github-actions-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “A popular Marketplace action is automatically safe.” Correction: third-party actions execute code; review the source, pin versions, and reduce permissions.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write the shell commands for a CI step that enters the project, runs `go test ./...`, and rejects unformatted Go files. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Build and publish one verified image**.

[Course contents](/course/cloud-devops?lang=en)
