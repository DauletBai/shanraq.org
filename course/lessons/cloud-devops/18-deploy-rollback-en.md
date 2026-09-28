# Deployment, smoke testing, and a verified rollback

_Lead (summary):_ **Deploy by version, test from outside, and restore the previous digest with a prewritten procedure.**

## Lesson outcome

Deploy by version, test from outside, and restore the previous digest with a prewritten procedure. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A rollback must exist before an incident. It restores code but may not restore compatible data; migrations need a separate plan. After switching, test a user path and watch metrics rather than treating a successful exit code as proof.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
export CLOUDLAB_IMAGE=ghcr.io/OWNER/cloudlab:REPLACE_WITH_SHA
export CLOUDLAB_DOMAIN=cloudlab.example.com
docker compose -f compose.prod.yaml config >/tmp/cloudlab-prod.yml
grep -n 'image:[|]CLOUDLAB_ENV' /tmp/cloudlab-prod.yml
printf 'record previous digest before docker compose up -d\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The configuration does not use latest.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 18](/static/course/cloud-devops/map-18-deploy-rollback-en.svg)

## Exercise

**Required.** Write a POSIX sh deploy function that saves the current digest, applies a new image, runs a smoke test, and restores the saved image if that test fails. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
