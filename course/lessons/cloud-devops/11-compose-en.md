# Docker Compose: services, networking, and declarative startup

_Lead (summary):_ **Start the application and reverse proxy from one declaration and reach a service by DNS name inside the network.**

## Lesson outcome

Start the application and reverse proxy from one declaration and reach a service by DNS name inside the network. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

Compose declares desired services, networks, and volumes. Inside a shared network, `app` resolves through built-in DNS; publishing a port is needed only at the host boundary. `depends_on` orders startup, while a health check proves readiness.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
docker compose config
docker compose up -d --build
docker compose ps
curl -fsS http://127.0.0.1:8080/healthz
docker compose logs --tail=10 app
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- Both services run and the app is healthy.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 11](/static/course/cloud-devops/map-11-compose-en.svg)

## Exercise

**Required.** Write commands that validate Compose, start the stack, wait at most 30 seconds for `/readyz`, and print the last 50 log lines on failure. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
