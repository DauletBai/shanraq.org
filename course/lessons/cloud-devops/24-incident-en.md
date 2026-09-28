# Observability, backups, and a recovery drill

_Lead (summary):_ **Define service indicators, create a verifiable backup, restore data, and write a blameless final report.**

## Lesson outcome

Define service indicators, create a verifiable backup, restore data, and write a blameless final report. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

Logs explain individual events, metrics show numbers changing over time, and traces connect a request path. A backup is useful only after a restore. An incident ends after service is restored, evidence is preserved, and follow-up work reduces recurrence, not merely after a restart.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
curl -fsS http://127.0.0.1:8080/metrics 2>/dev/null || true
mkdir -p backups
./scripts/backup.sh 2>/dev/null || true
printf 'A real drill restores into an isolated volume first.\n'
printf 'Record: detection, impact, timeline, recovery, evidence, follow-up owner and date.\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The metric is connected to a user outcome.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 24](/static/course/cloud-devops/map-24-incident-en.svg)

## Exercise

**Required.** Write a POSIX sh directory backup script with a UTC-stamped archive, adjacent SHA-256, archive verification, and no source deletion. Then list commands for a test restore into a new directory. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
