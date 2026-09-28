# Registry, immutable tags, SBOM, and supply-chain checks

_Lead (summary):_ **Tie an image to a commit SHA, produce an SBOM, and scan its contents before publishing.**

## Lesson outcome

Tie an image to a commit SHA, produce an SBOM, and scan its contents before publishing. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A registry stores and distributes images, but an upload does not create trust. A release needs an immutable version, traceable origin, and visible dependencies. An SBOM lists components; a scanner matches known vulnerabilities but cannot prove that no bug exists.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
sha=$(git rev-parse --short=12 HEAD)
docker build --build-arg VERSION="$sha" -t "cloudlab:$sha" .
docker image inspect "cloudlab:$sha" --format '{{.Id}}'
printf 'release=%s\n' "$sha"
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The tag equals the source commit SHA.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 13](/static/course/cloud-devops/map-13-registry-en.svg)

## Exercise

**Required.** Write a POSIX sh script that derives a tag from the 12-character Git SHA, builds the image, and prints its local image ID; do not use `latest`. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
