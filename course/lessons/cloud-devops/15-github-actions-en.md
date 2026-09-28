# GitHub Actions: verify every change

_Lead (summary):_ **Move local checks into a workflow, pin permissions, and use caching without storing secrets.**

## Lesson outcome

Move local checks into a workflow, pin permissions, and use caching without storing secrets. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A workflow is the repository's executable contract. A trigger selects the moment, permissions constrain the token, a job supplies a clean machine, and steps produce evidence. Third-party actions are code too; pin trusted versions and minimize permissions.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
cp workflow.example.yml /tmp/cloudlab-workflow.yml
sed -n '1,220p' /tmp/cloudlab-workflow.yml
grep -n 'permissions:[|]go test[|]build-push-action' /tmp/cloudlab-workflow.yml
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The workflow has explicit permissions.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 15](/static/course/cloud-devops/map-15-github-actions-en.svg)

## Exercise

**Required.** Write the shell commands for a CI step that enters the project, runs `go test ./...`, and rejects unformatted Go files. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
