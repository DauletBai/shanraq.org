# Build and publish one verified image

_Lead (summary):_ **Publish to a registry only after tests and promote one digest across environments.**

## Lesson outcome

Publish to a registry only after tests and promote one digest across environments. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

Rebuilding for production may produce an artifact different from the one tested. Build once, verify, publish, and promote its digest. Registry credentials belong only in the publishing job and must not reach an untrusted fork pull request.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
sha=$(git rev-parse --short=12 HEAD)
image=ghcr.io/OWNER/cloudlab:$sha
printf 'would publish %s\n' "$image"
printf 'production must record a sha256 digest, not latest\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The name contains an immutable version.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 16](/static/course/cloud-devops/map-16-release-image-en.svg)

## Exercise

**Required.** Write a POSIX sh check for `REGISTRY`, `IMAGE`, and `GIT_SHA` that constructs the full image name and refuses to publish the tag `latest`. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
