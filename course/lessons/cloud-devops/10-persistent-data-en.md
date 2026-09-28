# Volumes, bind mounts, and proof of persistence

_Lead (summary):_ **Move state out of the container, recreate it, and prove a note survives process replacement.**

## Lesson outcome

Move state out of the container, recreate it, and prove a note survives process replacement. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A container should be replaceable while state has an explicit lifecycle. A named volume suits Docker-managed data; a bind mount suits a known host path. Neither one is a backup by itself.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
docker volume create cloudlab-practice
docker run -d --name cloudlab-data -p 8080:8080 -v cloudlab-practice:/data cloudlab:lesson-09
sleep 1
curl -fsS -X POST -d 'note=survives-recreate' http://127.0.0.1:8080/notes >/dev/null
docker rm -f cloudlab-data
docker run -d --name cloudlab-data -p 8080:8080 -v cloudlab-practice:/data cloudlab:lesson-09
sleep 1
curl -fsS http://127.0.0.1:8080/ | grep survives-recreate
docker rm -f cloudlab-data
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The note remains after the container is removed.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 10](/static/course/cloud-devops/map-10-persistent-data-en.svg)

## Exercise

**Required.** Write commands that create a named volume, write a file through one temporary container, and read it through a second container. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
