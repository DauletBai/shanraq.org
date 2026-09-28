# CloudLab: service, environment, and readiness criteria

_Lead (summary):_ **Run the course service, inspect its endpoints, and turn a vague “works” into a testable contract.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **name → DNS → IP → TCP:port → TLS → HTTP → response**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Run the course service, inspect its endpoints, and turn a vague “works” into a testable contract. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

A café may be open while its kitchen is not ready for orders: the cook has arrived but the oven is warming. The open sign resembles process health, the ready kitchen resembles readiness, and a served meal is the user outcome. One green light is insufficient.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**person → address → service → response; alive ≠ ready ≠ useful**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Service** — A program that does useful work and is reachable through a defined entry point.
- **Endpoint** — A specific address within a service, such as `/healthz`.
- **Health** — Evidence that the process exists and can respond.
- **Readiness** — Evidence that the service is ready for real user traffic.

## Take it apart without rushing

CloudLab stores notes and exposes distinct signals: `/healthz` says the process is alive, `/readyz` says it can accept traffic, and `/metrics` counts requests. One HTTP 200 from the home page does not prove data durability, dependency readiness, or graceful shutdown.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
go test ./...
CLOUDLAB_ADDR=:8080 go run ./cmd/cloudlab &
pid=$!
trap 'kill "$pid" 2>/dev/null || true' EXIT
sleep 1
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/metrics
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

Tests first verify known rules. The service then starts in the background and its number is stored in `pid`. `trap` promises to stop it even after an early exit. `curl` acts as a visitor and asks two different addresses.

## What you should observe

- The tests pass.
- The health endpoint responds without HTML.
- The request counter rises after another curl.

## Rebuild the whole from the support map

![Support map 6](/static/course/cloud-devops/map-02-service-map-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “If `/healthz` returns 200, data is certainly durable and the whole service is useful.” Correction: health proves only its narrow contract.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh check that uses `curl -fsS` for `/healthz` and `/readyz`, then prints `cloudlab ready`. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **A container: a bounded process, not a tiny VM**.

[Course contents](/course/cloud-devops?lang=en)
