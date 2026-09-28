# CI/CD as a feedback system

_Lead (summary):_ **Design a pipeline from independent checks and separate continuous integration, delivery, and deployment.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **commit → image → digest → registry → SBOM/scan → deploy**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Design a pipeline from independent checks and separate continuous integration, delivery, and deployment. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

In a factory, a part is checked after each important stage, not only after the entire car is assembled. CI gives fast feedback after a change. Delivery prepares a verified release; deployment actually changes a running environment. These are different promises.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**commit → fast checks → build once → verify → deliver → deploy deliberately**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **CI** — Automated integration of changes with fast tests and checks.
- **Delivery** — A state in which a verified release is ready to deploy.
- **Deployment** — The actual installation of a version into an environment.
- **Pipeline** — A sequence of automated stages and transition conditions.

## Take it apart without rushing

CI quickly reports whether a change can be merged. Continuous delivery keeps a verified artifact releasable; continuous deployment releases every passing change automatically. Speed without trustworthy gates only delivers defects sooner.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

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

### What this experiment actually does

The local script mirrors the future pipeline: formatting, tests, build, and image verification. Run it after a deliberate small error and again after fixing it; feedback should be clear.

## What you should observe

- Checks run from fast to more expensive.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 14](/static/course/cloud-devops/map-14-cicd-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “CI/CD means every commit automatically deploys to production.” Correction: deployment frequency and approval are separate risk decisions.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh pipeline for formatting, tests, image build, and a smoke test. It must stop at the first failure and always remove its test container. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **GitHub Actions: verify every change**.

[Course contents](/course/cloud-devops?lang=en)
