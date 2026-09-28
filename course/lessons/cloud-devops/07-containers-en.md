# A container: a bounded process, not a tiny VM

_Lead (summary):_ **Run a first container, observe namespaces and an immutable image, and distinguish a process from a virtual machine.**

## Lesson outcome

Run a first container, observe namespaces and an immutable image, and distinguish a process from a virtual machine. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A container shares the host kernel while isolating process, network, and filesystem views. An image is a template; a container is one running instance. Removing a container must not equal losing important data.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
docker version
docker run --rm alpine:3.23 cat /etc/os-release
docker ps -a
docker image ls alpine:3.23
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- `--rm` removes the learning container after exit.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 7](/static/course/cloud-devops/map-07-containers-en.svg)

## Exercise

**Required.** Write commands that run `alpine:3.23`, print its hostname, and remove the container automatically after exit. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
