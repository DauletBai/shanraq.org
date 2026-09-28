# Image layers, tags, and the container lifecycle

_Lead (summary):_ **Inspect image history, pin a digest, and distinguish stop, start, remove, and pull.**

## Lesson outcome

Inspect image history, pin a digest, and distinguish stop, start, remove, and pull. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A tag is a movable name; a digest identifies exact content. `latest` is insufficient for a reproducible release because it may point to different bytes tomorrow. Containers and images also have separate lifecycles.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
docker pull alpine:3.23
docker image inspect alpine:3.23 --format '{{index .RepoDigests 0}}'
docker history alpine:3.23
docker create --name cloudlab-inspect alpine:3.23 true
docker start -a cloudlab-inspect
docker rm cloudlab-inspect
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The digest contains sha256.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 8](/static/course/cloud-devops/map-08-images-lifecycle-en.svg)

## Exercise

**Required.** Write commands that pull an image, save its first RepoDigest in a variable, and stop if the digest is empty. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
