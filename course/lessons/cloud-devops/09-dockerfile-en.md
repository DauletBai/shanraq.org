# Dockerfile: small, tested, and unprivileged

_Lead (summary):_ **Build CloudLab as a multi-stage image, test during the build, and verify its unprivileged user.**

## Lesson outcome

Build CloudLab as a multi-stage image, test during the build, and verify its unprivileged user. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A multi-stage build leaves the compiler in the build stage and copies only binaries into the final image. This reduces size and attack surface. `USER 65532` constrains the process, while a static health check needs no shell in the image.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
docker build --build-arg VERSION=lesson-09 -t cloudlab:lesson-09 .
docker image inspect cloudlab:lesson-09 --format 'user={{.Config.User}} size={{.Size}}'
docker run --rm --entrypoint /cloudlab-healthcheck cloudlab:lesson-09 || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The image builds only after tests pass.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 9](/static/course/cloud-devops/map-09-dockerfile-en.svg)

## Exercise

**Required.** Write commands to build `cloudlab:practice` and fail via `docker image inspect` when the configured user is `0` or empty. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
