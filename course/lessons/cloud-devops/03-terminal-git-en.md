# Terminal and Git: a reproducible working point

_Lead (summary):_ **Navigate the workspace, review a change before committing, and return to a known version without deleting history.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **🏭 hardware → ☁️ resource → 🔁 change → 👀 observe → 🛟 recover**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Navigate the workspace, review a change before committing, and return to a known version without deleting history. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A terminal is like speaking to a very literal assistant: it will not guess which room you intended to work in. Git resembles a workshop log: it records checkpoints and shows what changed. Before a risky command, ask “where am I?” and “what changed?”

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**📂 where? → 📝 what changed? → 📸 which version? → ✅ repeatable**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Terminal** — A window for text-based interaction with the operating system.
- **Shell** — A program that reads commands and starts other programs.
- **Repository** — A project directory together with its Git history.
- **Commit и SHA** — A saved history point and its unique short identifier.

## Take it apart without rushing

An infrastructure command is dangerous when its directory, version, and file differences are unknown. Git is a journal of intent rather than magical undo. Record a point, inspect the diff, and tie a release to an exact commit SHA.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
pwd
git status --short
git diff --check
git rev-parse --short HEAD
git log -1 --format='%h %cs %s'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

`pwd` prints the current directory. `status` shows unrecorded changes, `diff --check` finds some formatting errors, `rev-parse` identifies the exact version, and `log` adds the date and meaning of the latest commit.

## What you should observe

- `pwd` identifies the repository.
- `git diff --check` reports no whitespace errors.
- The SHA ties code to the release unambiguously.

## Rebuild the whole from the support map

![Support map 2](/static/course/cloud-devops/map-03-terminal-git-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “Git knows which change is good and can always undo everything safely.” Correction: Git stores file history; a person chooses meaning and a safe recovery method.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh script with `set -eu` that refuses to run outside a Git repository and prints the current short commit SHA. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Linux: users, files, and least privilege**.

[Course contents](/course/cloud-devops?lang=en)
