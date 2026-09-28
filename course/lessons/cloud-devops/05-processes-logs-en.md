# Processes, signals, services, and logs

_Lead (summary):_ **Observe a PID, graceful shutdown, and structured logs, then read a systemd unit as a startup contract.**

## Lesson outcome

Observe a PID, graceful shutdown, and structured logs, then read a systemd unit as a startup contract. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A process has a lifecycle. SIGTERM requests shutdown and allows state to be saved; SIGKILL removes that chance. A service manager adds identity, environment, restart policy, and logs, but cannot fix an application that ignores shutdown.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
CLOUDLAB_ADDR=:8080 go run ./cmd/cloudlab > /tmp/cloudlab.log 2>&1 &
pid=$!
sleep 1
ps -p "$pid" -o pid=,ppid=,command=
kill -TERM "$pid"
wait "$pid"
tail -n 5 /tmp/cloudlab.log
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The PID exists before SIGTERM.
- `wait` completes without a forced kill.
- The log has time, level, and message fields.

![Support map 5](/static/course/cloud-devops/map-05-processes-logs-en.svg)

## Exercise

**Required.** Write a POSIX sh script that starts a background process, records its PID, installs a SIGTERM trap, and waits for clean shutdown. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
