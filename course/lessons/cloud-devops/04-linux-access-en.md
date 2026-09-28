# Linux: users, files, and least privilege

_Lead (summary):_ **Understand owners, groups, permission modes, and why a service should not run permanently as root.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **📂 where? → 📝 what changed? → 📸 which version? → ✅ repeatable**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Understand owners, groups, permission modes, and why a service should not run permanently as root. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

Workshop keys have different rights: one opens storage and another only the break room. Giving everyone a master key is convenient until one is lost. Linux similarly separates owner, group, and others; root is a master key that should not be used continuously.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**👤 owner | 👥 group | 🌍 others × read/write/execute**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **User** — The identity under which a process runs.
- **root** — The superuser with very few ordinary access restrictions.
- **Permissions** — Read, write, and execute rules for three groups of identities.
- **640** — Owner reads and writes; group reads; others have no access.

## Take it apart without rushing

Least privilege limits blast radius: a compromised process receives only what its task needs. Modes `750` and `640` are not rituals; they encode read, write, and execute access for owner, group, and others.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
id
umask
mkdir -p /tmp/cloudlab-permissions
printf '%s\n' secret > /tmp/cloudlab-permissions/config
chmod 640 /tmp/cloudlab-permissions/config
ls -ld /tmp/cloudlab-permissions
ls -l /tmp/cloudlab-permissions/config
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

`id` names your current identity. `umask` shows default restrictions. A learning file is created, `chmod 640` assigns rights, and `ls -l` lets you decode the result rather than trusting the command.

## What you should observe

- The file is not executable.
- Other users have no access.
- You can decode every digit in mode 640.

## Rebuild the whole from the support map

![Support map 3](/static/course/cloud-devops/map-04-linux-access-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “Running as root fixes permission problems.” Correction: it hides bad configuration and increases the impact of compromise or error.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh script that creates a directory with mode 750 and a config file with mode 640, then verifies both modes with `stat`. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Processes, signals, services, and logs**.

[Course contents](/course/cloud-devops?lang=en)
