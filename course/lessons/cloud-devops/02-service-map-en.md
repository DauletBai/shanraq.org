# CloudLab: service, environment, and readiness criteria

_Lead (summary):_ **Run the course service, inspect its endpoints, and turn a vague “works” into a testable contract.**

## Lesson outcome

Run the course service, inspect its endpoints, and turn a vague “works” into a testable contract. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

CloudLab stores notes and exposes distinct signals: `/healthz` says the process is alive, `/readyz` says it can accept traffic, and `/metrics` counts requests. One HTTP 200 from the home page does not prove data durability, dependency readiness, or graceful shutdown.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

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

## Read the result

- The tests pass.
- The health endpoint responds without HTML.
- The request counter rises after another curl.

![Support map 2](/static/course/cloud-devops/map-02-service-map-en.svg)

## Exercise

**Required.** Write a POSIX sh check that uses `curl -fsS` for `/healthz` and `/readyz`, then prints `cloudlab ready`. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
