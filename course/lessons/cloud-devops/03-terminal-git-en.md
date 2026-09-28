# Terminal and Git: a reproducible working point

_Lead (summary):_ **Navigate the workspace, review a change before committing, and return to a known version without deleting history.**

## Lesson outcome

Navigate the workspace, review a change before committing, and return to a known version without deleting history. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

An infrastructure command is dangerous when its directory, version, and file differences are unknown. Git is a journal of intent rather than magical undo. Record a point, inspect the diff, and tie a release to an exact commit SHA.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
pwd
git status --short
git diff --check
git rev-parse --short HEAD
git log -1 --format='%h %cs %s'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- `pwd` identifies the repository.
- `git diff --check` reports no whitespace errors.
- The SHA ties code to the release unambiguously.

![Support map 3](/static/course/cloud-devops/map-03-terminal-git-en.svg)

## Exercise

**Required.** Write a POSIX sh script with `set -eu` that refuses to run outside a Git repository and prints the current short commit SHA. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
