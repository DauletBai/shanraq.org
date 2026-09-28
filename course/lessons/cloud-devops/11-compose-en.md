# Docker Compose: services, networking, and declarative startup

_Lead (summary):_ **Start the application and reverse proxy from one declaration and reach a service by DNS name inside the network.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **container changed; volume remained; restore proved data**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Start the application and reverse proxy from one declaration and reach a service by DNS name inside the network. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

An orchestra does not start by calling every musician separately. A score lists participants, order, and relationships. A Compose file acts as that score: application, proxy, network, volume, and checks are described together, and service names become internal addresses.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**services + network + volumes = one repeatable start**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Docker Compose** — A tool that starts related containers from one YAML description.
- **Service** — A named container role and its settings in Compose.
- **Network** — A virtual network where services find one another by name.
- **YAML** — A human-readable data format in which indentation matters.

## Take it apart without rushing

Compose declares desired services, networks, and volumes. Inside a shared network, `app` resolves through built-in DNS; publishing a port is needed only at the host boundary. `depends_on` orders startup, while a health check proves readiness.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
docker compose config
docker compose up -d --build
docker compose ps
curl -fsS http://127.0.0.1:8080/healthz
docker compose logs --tail=10 app
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

`config` expands and checks the final description. `up -d` brings the environment to it, `ps` shows state, a request through the proxy tests the user path, and `down` removes the created processes and network.

## What you should observe

- Both services run and the app is healthy.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 11](/static/course/cloud-devops/map-11-compose-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “`depends_on` means a dependency is fully ready.” Correction: start order and readiness differ; readiness needs a separate check.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write commands that validate Compose, start the stack, wait at most 30 seconds for `/readyz`, and print the last 50 log lines on failure. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Health checks, limits, and graceful shutdown**.

[Course contents](/course/cloud-devops?lang=en)
