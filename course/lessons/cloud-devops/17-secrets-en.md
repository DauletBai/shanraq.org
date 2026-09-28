# Secrets, configuration, and environments without leaks

_Lead (summary):_ **Separate public configuration from secrets, scope tokens narrowly, and prevent values from reaching logs.**

## Lesson outcome

Separate public configuration from secrets, scope tokens narrowly, and prevent values from reaching logs. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A secret is more than an environment variable: it has an owner, scope, lifetime, and rotation procedure. Log masking helps but does not remove the risk of passing a secret to a process or third-party action. Reduce privilege and lifetime first.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
test -f .gitignore
git grep -n -E '(BEGIN (RSA|OPENSSH) PRIVATE KEY|token=|password=)' -- ':!tools/course/generate_cloud_devops.py' || true
printf '%s\n' '.env must stay untracked'
git check-ignore .env || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- No real value is printed.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 17](/static/course/cloud-devops/map-17-secrets-en.svg)

## Exercise

**Required.** Write a POSIX sh script that requires `DEPLOY_TOKEN`, never prints it, checks that it is nonempty, and unsets it before exit. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
