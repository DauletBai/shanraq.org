# CI/CD as a feedback system

_Lead (summary):_ **Design a pipeline from independent checks and separate continuous integration, delivery, and deployment.**

## Lesson outcome

Design a pipeline from independent checks and separate continuous integration, delivery, and deployment. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

CI quickly reports whether a change can be merged. Continuous delivery keeps a verified artifact releasable; continuous deployment releases every passing change automatically. Speed without trustworthy gates only delivers defects sooner.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
go test ./...
gofmt -l cmd internal
docker build -t cloudlab:ci .
docker run --rm -d --name cloudlab-ci -p 18080:8080 cloudlab:ci
trap 'docker rm -f cloudlab-ci >/dev/null 2>&1 || true' EXIT
sleep 1
curl -fsS http://127.0.0.1:18080/readyz
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- Checks run from fast to more expensive.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 14](/static/course/cloud-devops/map-14-cicd-en.svg)

## Exercise

**Required.** Write a POSIX sh pipeline for formatting, tests, image build, and a smoke test. It must stop at the first failure and always remove its test container. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
