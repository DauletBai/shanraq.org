# Health checks, limits, and graceful shutdown

_Lead (summary):_ **Distinguish liveness from readiness, set resource boundaries, and test shutdown behaviour.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **services + network + volumes = one repeatable start**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Distinguish liveness from readiness, set resource boundaries, and test shutdown behaviour. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

In a hospital, a pulse says whether a person is alive but not whether a surgeon is ready to operate. An operating-room schedule limits concurrent work. Liveness, readiness, and resource limits likewise answer three different questions.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**liveness: restart? | readiness: send traffic? | limits: how much?**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Liveness** — A check asking whether a process is stuck enough to restart.
- **Readiness** — A check asking whether requests should be sent here now.
- **Resource limit** — An upper CPU or memory boundary for a process.
- **Graceful shutdown** — Shutdown that stops new traffic and finishes work already started.

## Take it apart without rushing

Liveness asks whether the process should restart; readiness asks whether it should receive new traffic. Confusing them turns temporary load into a restart loop. Limits protect neighbours, but a limit set too low creates its own failure.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
docker compose ps
docker inspect --format '{{json .State.Health}}' cloud-devops-lab-app-1 2>/dev/null || true
time docker compose stop -t 10 app
docker compose up -d app
curl -fsS http://127.0.0.1:8080/readyz
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

The exercise observes both endpoints, sends SIGTERM, and checks termination. Look beyond the command status: ask whether new work stopped and existing work had time to finish.

## What you should observe

- Shutdown completes within the grace period.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 12](/static/course/cloud-devops/map-12-health-resources-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “The more frequent and deep a liveness check, the more reliable the service.” Correction: an expensive or wrong check can itself trigger restarts and an outage.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh loop with a deadline that waits for a readiness URL, pauses between attempts, and returns nonzero after timeout. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Registry, immutable tags, SBOM, and supply-chain checks**.

[Course contents](/course/cloud-devops?lang=en)
