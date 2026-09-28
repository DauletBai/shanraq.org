# Health checks, limits, and graceful shutdown

_Lead (summary):_ **Distinguish liveness from readiness, set resource boundaries, and test shutdown behaviour.**

## Lesson outcome

Distinguish liveness from readiness, set resource boundaries, and test shutdown behaviour. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

Liveness asks whether the process should restart; readiness asks whether it should receive new traffic. Confusing them turns temporary load into a restart loop. Limits protect neighbours, but a limit set too low creates its own failure.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
docker compose ps
docker inspect --format '{{json .State.Health}}' cloud-devops-lab-app-1 2>/dev/null || true
time docker compose stop -t 10 app
docker compose up -d app
curl -fsS http://127.0.0.1:8080/readyz
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- Shutdown completes within the grace period.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 12](/static/course/cloud-devops/map-12-health-resources-en.svg)

## Exercise

**Required.** Write a POSIX sh loop with a deadline that waits for a readiness URL, pauses between attempts, and returns nonzero after timeout. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
