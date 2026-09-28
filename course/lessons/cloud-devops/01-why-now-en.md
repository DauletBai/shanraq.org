# Why Cloud & DevOps now: infrastructure, roles, and boundaries

_Lead (summary):_ **Separate data-centre growth from job-title growth and map the skills of cloud engineering, DevOps, platform engineering, and SRE.**

## Lesson outcome

Separate data-centre growth from job-title growth and map the skills of cloud engineering, DevOps, platform engineering, and SRE. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

DevOps is neither one tool nor one person responsible for everything. It reduces change risk through automation, short feedback loops, and shared responsibility. A data centre supplies physical capacity, cloud turns it into programmable resources, and operations keeps a service available.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
printf '%s\n' 'change -> test -> image -> deploy -> observe -> recover'
printf '%s\n' 'evidence: repeatable command, metric, backup, recovery record'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The chain begins with a change and ends with verified recovery.
- Every stage leaves observable evidence.
- You can state the course boundary: software operations.

![Support map 1](/static/course/cloud-devops/map-01-why-now-en.svg)

## Exercise

**Required.** Write a POSIX sh script that prints the six release-path stages on separate lines and stops on the first error. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
