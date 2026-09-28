# Secrets, configuration, and environments without leaks

_Lead (summary):_ **Separate public configuration from secrets, scope tokens narrowly, and prevent values from reaching logs.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **test once → build once → sign/scan → push digest → promote same digest**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Separate public configuration from secrets, scope tokens narrowly, and prevent values from reaching logs. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

An office address may be printed on a card; a safe key may not. Both support work but require different handling. Configuration describes an environment; a secret grants power. A secret does not belong in an image, log, or Git history.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**config may be shown | secret is restricted | both enter outside image**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Configuration** — Non-secret parameters controlling application behaviour.
- **Secret** — A value whose disclosure grants access or causes harm.
- **Environment variable** — One channel for passing a setting to a process at startup.
- **Rotation** — Replacing a secret and revoking the old value according to a plan.

## Take it apart without rushing

A secret is more than an environment variable: it has an owner, scope, lifetime, and rotation procedure. Log masking helps but does not remove the risk of passing a secret to a process or third-party action. Reduce privilege and lifetime first.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
test -f .gitignore
git grep -n -E '(BEGIN (RSA|OPENSSH) PRIVATE KEY|token=|password=)' -- ':!tools/course/generate_cloud_devops.py' || true
printf '%s\n' '.env must stay untracked'
git check-ignore .env || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The exercise uses only a dummy value, checks `.gitignore`, and searches for suspicious strings. The process receives the secret at startup. Do not print it even for demonstration; prove only its presence.

## What you should observe

- No real value is printed.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 17](/static/course/cloud-devops/map-17-secrets-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “Removing a secret from the latest commit is enough.” Correction: it remains in history and may have been copied; revoke and rotate it.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh script that requires `DEPLOY_TOKEN`, never prints it, checks that it is nonempty, and unsets it before exit. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Deployment, smoke testing, and a verified rollback**.

[Course contents](/course/cloud-devops?lang=en)
