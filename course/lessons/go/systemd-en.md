# A service that survives a reboot: systemd and a tidy stop

_Лид (summary):_ **The forty-eighth lesson of the Go course. A program is stopped by a signal rather than a button, and the difference is measured: cut off, a reader gets an EOF; stopped tidily, an answer of 200 after 270 ms. Plus a unit file the real systemd checked, a start after every reboot, and the order in which things close.**

## Why this is needed

The blog runs from a terminal. Close the laptop and the blog is gone. Reboot the server and the blog is gone. Deploy a new version and whoever was reading an article at that moment got a broken connection instead of a page.

Two things cure all of it. The first is a **service**: the system starts the blog at boot and brings it back if it falls over. The second is a **tidy stop**: the program learns it has been asked to finish, writes out the answers already started, and only then goes.

We start with the second, because without it the first only makes things worse: the service will diligently restart the blog, cutting readers off at every deploy.

About the second half straight away: `systemd` exists only on Linux. On macOS or Windows the tidy stop works exactly the same — it is Go code, not a system setting — and the unit file is best tried on the server the blog is deployed to. The same job is done by `launchd` on macOS and by services on Windows; the file differs, the meaning does not.

## The whole thing at once

A new folder, `go mod init sabaq45`. The program brings a server up, sends a reader after a page, and halfway through the answer stops the server — in four different ways:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// A reader who has already clicked: the answer takes a third of a second.
const answerTime = 300 * time.Millisecond

func main() {
	fmt.Println("== cut off: the server was closed mid-sentence")
	report(stop(func(srv *http.Server) error { return srv.Close() }))

	fmt.Println()
	fmt.Println("== tidily: Shutdown waited for the answer")
	report(stop(func(srv *http.Server) error {
		return srv.Shutdown(context.Background())
	}))

	fmt.Println()
	fmt.Println("== less time than the answer needs")
	report(stop(func(srv *http.Server) error {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			// What is done in real life: there is no time left, so the connections go.
			srv.Close()
		}
		return err
	}))

	fmt.Println()
	fmt.Println("== the signal systemd sends")
	signalDemo()
}

// result is what the reader got and what the server was doing meanwhile.
type result struct {
	reader   string
	waited   time.Duration
	shutdown error
}

// stop brings a server up, sends a reader after a page, and halfway through
// the answer stops the server the way it was told to.
func stop(how func(*http.Server) error) result {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	srv := &http.Server{Handler: http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(answerTime)
			fmt.Fprint(w, "done")
		})}
	go srv.Serve(l)

	done := make(chan string, 1)
	go func() {
		resp, err := http.Get("http://" + l.Addr().String() + "/")
		if err != nil {
			done <- "cut off: " + shortErr(err)
			return
		}
		defer resp.Body.Close()
		done <- fmt.Sprintf("answer %d", resp.StatusCode)
	}()

	// Let the reader reach the handler: the stop has to happen in the middle of
	// an answer, or we measure the wrong thing.
	time.Sleep(50 * time.Millisecond)

	start := time.Now()
	shutdownErr := how(srv)
	waited := time.Since(start)

	return result{reader: <-done, waited: waited.Round(10 * time.Millisecond), shutdown: shutdownErr}
}

func report(r result) {
	fmt.Printf("the reader got:     %s\n", r.reader)
	fmt.Printf("stopping took:      %v\n", r.waited)
	if r.shutdown != nil {
		fmt.Printf("stopping returned:  %v\n", r.shutdown)
		return
	}
	fmt.Println("stopping returned:  nil")
}

// signalDemo shows how a program learns it has been asked to finish. systemd
// sends SIGTERM and waits; here we send it to ourselves.
func signalDemo() {
	ctx, stopSignals := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, os.Interrupt)
	defer stopSignals()

	start := time.Now()
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		fmt.Println("could not send the signal:", err)
		return
	}

	<-ctx.Done()
	fmt.Printf("the signal arrived:   %s\n", instant(time.Since(start)))
	fmt.Printf("the context says:     %v\n", ctx.Err())
	fmt.Println("next:                 Shutdown, then close the database")
}

// shortErr strips the address with its random port out of the error, so the
// output is the same on every run.
func shortErr(err error) string {
	var ue *net.OpError
	if errors.As(err, &ue) {
		return ue.Err.Error()
	}
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return shortErr(e.Unwrap())
	}
	return err.Error()
}

// instant hides the exact figure: what matters is that the signal arrives
// before the program can do anything, not how many microseconds it took.
func instant(d time.Duration) string {
	if d < time.Millisecond {
		return "at once, under a millisecond"
	}
	return d.Round(time.Millisecond).String()
}
```

The output:

```
== cut off: the server was closed mid-sentence
the reader got:     cut off: EOF
stopping took:      0s
stopping returned:  nil

== tidily: Shutdown waited for the answer
the reader got:     answer 200
stopping took:      270ms
stopping returned:  nil

== less time than the answer needs
the reader got:     cut off: EOF
stopping took:      50ms
stopping returned:  context deadline exceeded

== the signal systemd sends
the signal arrived:   at once, under a millisecond
the context says:     context canceled
next:                 Shutdown, then close the database
```

## The walk-through

### A cut and an answer are one call apart

The first two blocks differ by a single call. `srv.Close()` closes everything at once: stopping took `0s`, and instead of a page the reader got an `EOF` — the connection broke mid-sentence.

`srv.Shutdown(ctx)` does something else: it stops taking new requests, **waits for the ones already started**, and returns. Stopping took `270 ms` — exactly as long as our handler needed for the answer — and the reader got a `200`.

In real life the difference looks like this: on a deploy, half the readers see "this page is unavailable" — or notice nothing unusual at all.

> **Picture it.** A shop is closing. You can turn the lights off and lock the door with people inside, or you can stop letting new ones in and wait for those already at the till.

### The time has to be enough

The third block is the same tidy stop, but with 50 ms given to an answer that needs 300. `Shutdown` returned `context deadline exceeded`, and after that the program has nothing left to do but break the connections: the reader got an `EOF` again.

Hence the rule: **the stopping period has to be longer than the slowest answer**. In the blog that is an image upload, not a page with an article. Step 33 uses 15 seconds.

### A signal is how the system talks to a program

The fourth block answers the question of how a program learns that it is time. It is sent a **signal** — a short message from the system. `SIGTERM` means "finish up"; it is what systemd sends, and docker, and the stop button in a host's panel. `Ctrl+C` in a terminal is `SIGINT`, with the same meaning.

In Go that is one line:

```go
ctx, stopSignals := signal.NotifyContext(context.Background(),
	syscall.SIGTERM, os.Interrupt)
defer stopSignals()
```

`NotifyContext` turns the signal into a **cancelled context** — the very tool from the lesson on `context`. After that the program simply waits on `<-ctx.Done()` and needs no special knowledge of signals at all. Measured: the signal arrives in under a millisecond.

There is one signal that cannot be caught: `SIGKILL` (`kill -9`). It never reaches the program — the kernel kills the process. That is why `kill -9` is not a way of stopping something but a last resort.

### The order matters more than the code

Put together it looks like this, and the order here is not decoration:

1. the signal arrived — the context closed;
2. `srv.Shutdown(ctx)` — nothing new is taken, what started is waited for;
3. and only now `store.Close()` — the database closes last.

Closing the database first means breaking with your own hands the very answers you are waiting for. The mistake is common and looks harmless: `defer store.Close()` stands at the top of `main`, and that seems to be enough. It is — but only because `defer` runs **after** `Shutdown`, not before.

### The unit file

Now the first half of the job. systemd is what starts and watches programs on Linux. It needs one file; step 33 holds it in full, and the heart of it is these lines:

```ini
[Service]
User=blog
ExecStart=/opt/blog/blog
Environment=BLOG_ADDR=127.0.0.1:8080
Restart=always
RestartSec=2
KillSignal=SIGTERM
TimeoutStopSec=20

[Install]
WantedBy=multi-user.target
```

`Restart=always` — it fell, we bring it back; `RestartSec=2` — no more than once every two seconds, or a broken program restarts a thousand times a minute. `KillSignal=SIGTERM` is what we have just taught the program to accept. `TimeoutStopSec=20` is how long systemd waits before killing the process outright; this number has to be **larger** than our own stopping period, or the tidy stop never finishes and turns into the same cut.

`WantedBy=multi-user.target` is what "start it when the machine boots" actually is.

The file from the step was checked by the real systemd on a real server:

```
$ systemd-analyze verify blog.service
blog.service: Command /opt/blog/blog is not executable: No such file or directory
```

The only remark is that the program file itself is not on that machine; the unit itself drew no objection.

Three commands worth knowing:

```
systemctl enable --now blog    # start it and start it at boot
systemctl restart blog         # restart after a deploy
journalctl -u blog -f          # watch the log live
```

You have been writing that log since the lesson on errors and `slog`. systemd picks up everything a program prints to its output and files it into `journalctl`; no log files of your own are needed.

### A few more lines that cost a minute

The unit file holds lines that cost nothing and close a great deal:

- `User=blog` — the blog runs as its own user rather than as `root`. A blog that is broken into is then not a machine that is broken into;
- `NoNewPrivileges=true` — the process cannot gain more rights than it started with;
- `ProtectSystem=strict` with `ReadWritePaths=/var/lib/blog` — the whole system is read-only to it except the one folder with the database and the pictures;
- `PrivateTmp=true` — a `/tmp` of its own rather than the shared one.

## The map of the lesson

![The map of the lesson: a signal, Shutdown and a service](/static/course/go/map-systemd-en.svg)

## Say it in your own words

Without looking, answer aloud or on paper. The answers are at the end of the lesson.

1. How does `srv.Shutdown` differ from `srv.Close` from a reader's point of view?
2. Why is the database closed after `Shutdown` rather than before?
3. Why does `TimeoutStopSec` in the unit file have to be larger than the stopping period in the program?

## The exercise

**Required.** Teach your blog to stop tidily: the signal is caught with `signal.NotifyContext`, `Shutdown` is given a sensible period, and the database closes last. Check it: start the blog, send it `kill -TERM`, and look at the exit code — it should be `0`, with two lines about the stop in the log.

All of it is done in [step-33](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-33) — compare after you have done it yourself.

**If you want more.**

- Write your own unit file and check it: `systemd-analyze verify blog.service`.
- Set `TimeoutStopSec=1` and see what happens to a slow answer on a restart.
- Add a log line to the stop with the number of connections that had to be waited for.

## Where this goes in the blog

Step 33 differs from step 32 in one thing: it knows how to stop. The difference is measured right in the terminal:

```
step 32, SIGTERM:  exit code 143 — the process was killed where it stood,
                   the last line of the log is somebody's request

step 33, SIGTERM:  exit code 0
                   "the stop signal arrived, waiting 15s"
                   "every answer is finished, stopping"
```

Code `143` is `128 + 15`, that is, "the process was killed by signal number 15". Code `0` is "the program finished by itself and properly". In systemd's own view the first shows as `failed` and the second as `inactive (dead)`, and that line is where people look afterwards for the reason something died.

Still open. The blog still sets no timeouts on reading and writing a request, and without them one slow connection can hold a goroutine for a long time. That is the pre-launch checklist — the next lesson.

## The answers

1. `Close` breaks the connections immediately: the reader gets a cut instead of a page. `Shutdown` stops taking new requests, waits for the ones already started, and only then returns — the reader gets their answer.
2. Because the answers we are waiting for may go to the database. Closing it first breaks exactly the requests we are waiting on.
3. Because systemd waits for `TimeoutStopSec` and then kills the process. If it is smaller, the program does not finish its answers, and the tidy stop turns into the same cut.

## Sources

- [Go: http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- [Go: signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext)
- [systemd: the systemd.service manual](https://man7.org/linux/man-pages/man5/systemd.service.5.html)
