# Processes, signals, services, and logs

_Lead (summary):_ **Observe a PID, graceful shutdown, and structured logs, then read a systemd unit as a startup contract.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **👤 owner | 👥 group | 🌍 others × read/write/execute**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Observe a PID, graceful shutdown, and structured logs, then read a systemd unit as a startup contract. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

Sheet music is not yet music: music appears when a performer starts playing. A program file similarly becomes a process when launched. A PID is the performer's stage number, a signal is the conductor's request, and a log records what happened.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**program + start = process(PID); TERM → stop; log → what happened**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Process** — A running instance of a program with memory and state.
- **PID** — The numeric identifier of a specific process.
- **SIGTERM** — A polite request for a process to shut down and clean up.
- **Log** — A time-stamped event record with context for diagnosis.

## Take it apart without rushing

A process has a lifecycle. SIGTERM requests shutdown and allows state to be saved; SIGKILL removes that chance. A service manager adds identity, environment, restart policy, and logs, but cannot fix an application that ignores shutdown.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

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

### What this experiment actually does

The service starts in the background and its PID is saved. `ps` proves that exact process exists. `kill -TERM` asks it to stop, `wait` collects the result, and `tail` shows the last log events.

## What you should observe

- The PID exists before SIGTERM.
- `wait` completes without a forced kill.
- The log has time, level, and message fields.

## Rebuild the whole from the support map

![Support map 4](/static/course/cloud-devops/map-05-processes-logs-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “SIGKILL is a more reliable normal shutdown.” Correction: it gives the application no time to save data or finish work; it is an emergency measure.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh script that starts a background process, records its PID, installs a SIGTERM trap, and waits for clean shutdown. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Networking without magic: IP, port, DNS, HTTP, and TLS**.

[Course contents](/course/cloud-devops?lang=en)
