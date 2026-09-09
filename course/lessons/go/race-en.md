# Goroutines and channels: what the web needs them for

_Лид (summary):_ **The forty-second lesson of the Go course. Your server has been concurrent since the HTTP lesson and you did not choose it: every request is a goroutine of its own. Measured here: a counter without a lock first loses one reading out of five hundred, then brings the whole server down, and go test -race finds it beforehand.**

## Why this matters

Goroutines are usually introduced with an example where something has to be made faster. For the web that is the wrong way in: **your server is already concurrent**, and it became so without you.

`http.Server` runs every request in a goroutine of its own. While one reader waits for an answer, another is already being served. That is convenient and dangerous at once: everything the handlers touch together, they touch **at the same time**.

While the blog kept everything in the database there was no trouble: SQLite sorts out who writes. Add one variable in memory — a counter of readings, a cache, anything — and a race appears.

## The whole thing at once

A new folder, `go mod init sabaq39`:

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
)

// views counts readings. A handler runs in a goroutine of its own for every
// request, so this variable is touched by several at the same moment.
type views struct {
	mu sync.Mutex
	n  map[string]int
}

func newViews() *views { return &views{n: map[string]int{}} }

// Add without a lock — never write this; below is why.
func (v *views) AddUnsafe(slug string) { v.n[slug]++ }

// Add with a lock: one goroutine is let in at a time.
func (v *views) Add(slug string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.n[slug]++
}

func (v *views) Count(slug string) int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.n[slug]
}

func handler(v *views, safe bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if safe {
			v.Add("dala")
		} else {
			v.AddUnsafe("dala")
		}
		fmt.Fprint(w, "ok")
	})
}

// hammer sends many requests at once and waits for all of them.
func hammer(h http.Handler, n int) {
	srv := httptest.NewServer(h)
	defer srv.Close()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(srv.URL)
			if err == nil {
				resp.Body.Close()
			}
		}()
	}
	wg.Wait()
}

func main() {
	v := newViews()
	hammer(handler(v, true), 500)
	fmt.Println("with a lock:   ", v.Count("dala"), "of 500")

	u := newViews()
	hammer(handler(u, false), 500)
	fmt.Println("without a lock:", u.Count("dala"), "of 500")
}
```

Run it, and the answers differ:

```
with a lock:    500 of 500
without a lock: 499 of 500
```

```
with a lock:    500 of 500
fatal error: concurrent map writes

goroutine 5411 [running]:
```

The same program: in one run it lost one reading out of five hundred, in another it brought itself down entirely.

## Taking it apart

### A race: a program that sometimes lies and sometimes dies

`AddUnsafe` does `v.n[slug]++`. That looks like one action; it is three: read, add, write.

Two goroutines manage to read the same value, each adds its own one, each writes — and a reading disappears. That is how `499 of 500` comes about.

With a map it is worse. A concurrent write to a `map` is not merely a wrong number: Go checks for it deliberately and ends the program with `fatal error: concurrent map writes`. Not a panic that `recover` can catch, but the immediate death of the process. Your server stops answering everybody.

> **Picture it.** Two people correcting one line in a paper ledger. Each read "17", each crossed it out and wrote "18". There were two visitors, and one got written down.

The worst outcome is the third — the one in the first output: the program worked correctly. A race does not reproduce on demand. It waits for load, which means it fires on the working server rather than on your laptop.

### `-race`: it finds what cannot be seen

So a race is not looked for by eye. Go looks for it:

```
with a lock:    500 of 500
without a lock: 488 of 500
==================
WARNING: DATA RACE
Read at 0x00c000dd2780 by goroutine 3668:
  runtime.mapdelete_faststr()
      /usr/local/go/src/internal/runtime/maps/runtime_faststr_swiss.go:396 +0x8c
  main.(*views).AddUnsafe()
      main.go:20 +0x58
  main.main.handler.func2()
      main.go:40 +0x68
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/go/src/net/http/server.go:2294 +0x48
  net/http.serverHandler.ServeHTTP()
```

The paths in the excerpt are shortened to the file name.

`go run -race .` and `go test -race ./...` switch on the race detector. It watches who touches memory and when, and prints **both** sides of the conflict with a file and a line — here `main.go:20`, that very `v.n[slug]++`.

The detector finds a race even when this particular run came out right. It slows the program down several times over, so it is not kept on in production — it is run in the tests, and it is usually one line in the check before code is sent off.

The rule is simple: **`-race` in the tests is compulsory for any program with shared state**. A web program always has some.

### The lock: `sync.Mutex` and `defer Unlock`

The cure is a lock:

```go
func (v *views) Add(slug string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.n[slug]++
}
```

`Lock` lets one goroutine in; the rest wait on that line until `Unlock`. The `defer` is not decoration: without it any early `return` or panic leaves the lock closed for ever, and the server stops dead.

A mutex is put **next to what it guards** — a field in the same struct. Then it is visible what it protects, and nobody has to remember.

Worth keeping in mind: do as little as possible under a lock. Going to the database or the network while holding one is a sure way to turn a fast server into a queue.

### A channel: a queue of background work

A mutex is about a shared variable. A channel is about handing work from one goroutine to another, and on the web it is needed for something else: so that a reader does not wait for what they have no reason to wait for.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// A queue of background work. A buffer of a hundred: the handler drops a
	// job in and answers the reader instead of waiting for it to be done.
	jobs := make(chan string, 100)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// A loop over a channel ends by itself once the channel is closed.
		for slug := range jobs {
			time.Sleep(50 * time.Millisecond) // as if posting to telegram
			fmt.Println("  sent:", slug)
		}
	}()

	start := time.Now()
	for _, slug := range []string{"dala", "salem", "shanyraq"} {
		jobs <- slug
		fmt.Printf("answered the reader in %s\n", time.Since(start).Round(time.Millisecond))
	}

	// Closing the channel says there will be no more jobs. Without it the loop
	// in the goroutine would wait for ever and Wait would never return.
	close(jobs)
	wg.Wait()
	fmt.Println("all done in", time.Since(start).Round(10*time.Millisecond))
}
```

```
answered the reader in 0s
answered the reader in 0s
answered the reader in 0s
  sent: dala
  sent: salem
  sent: shanyraq
all done in 150ms
```

The reader was answered instantly three times, and a hundred and fifty milliseconds of work finished afterwards. That is how a blog posts to telegram, sends an email, clears a cache — everything whose result the reader does not need.

The buffer of a hundred is not a luxury but a decision: if the queue fills up, `jobs <- slug` **blocks** and the handler starts waiting. Sometimes that is right (better to wait than to lose), sometimes not, and then a `select` with a `default` drops the job with a line in the log.

### `WaitGroup` and a closed channel

The two lines at the end matter more than they look.

`close(jobs)` says there will be no more jobs. The loop `for slug := range jobs` ends by itself only after the channel is closed; without it the goroutine would wait for ever.

`wg.Wait()` holds `main` until the goroutine has finished. Without it the program ends and **takes every goroutine with it** — unfinished work included. That, by the way, answers a common question: a goroutine does not outlive its program.

### What is not done

**A goroutine for every little thing.** A `go doSomething()` in a handler looks free, but each such goroutine lives a life of its own: nobody waits for it, its panic brings the whole server down, and under a rush of requests there will be thousands of them. Background work goes to one or two workers through a channel, as above.

**Catching a panic in a goroutine with a wrapper.** The `recoverPanic` from the errors lesson protects only the goroutine it sits in. A panic in a background goroutine kills the whole process, so `recover` goes inside it.

**Thinking goroutines are about speed.** They are about waiting. While one request waits for the database, another computes — that is the whole benefit. They will not add numbers any faster.

## The lesson map

![Lesson map: shared memory under a lock, work through a channel](/static/course/go/map-race-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why is a race more dangerous than an ordinary bug?
2. What is the `defer` next to `Lock` for?
3. What happens to a background goroutine if `main` finishes?

## Exercise

**Required.** Add a counter of readings in memory to the blog: the article handler increases it, the about page shows it. First without a lock — and run `go test -race ./...` to see the report. Then close it with a lock and make sure the report is gone.

All of it is done in [step-27](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-27) — compare against it once you have written your own.

**Optional.**

- Swap the mutex for a `sync.RWMutex` and explain what `RLock` gives when reads are frequent.
- Send the log line to one worker through a channel and check that the answers got faster.
- Take `wg.Wait()` out of the worker example and see how many jobs get done.

## Where this goes in your blog

The blog gets its first shared state in memory — and with it the habit of running the tests with `-race`.

Still open. The counter lives until a restart: surviving one needs the database. There is one worker and no queue in the database — if the process dies, the jobs are lost. And `context`, which can stop work when a reader closes the tab, is the next lesson.

## Answers

1. Because it does not reproduce on demand: the same program sometimes works, sometimes loses data, sometimes dies with a `fatal error`. An ordinary bug breaks the same way every time; a race waits for load — that is, for the working server.
2. So the lock opens on every way out of the function, an early `return` and a panic included. A forgotten `Unlock` stops not one goroutine but every one that comes after it.
3. It disappears along with the program, its work unfinished. A goroutine does not outlive its program, which is why it is waited for — with a `WaitGroup`, for instance.

## Sources

- [Go: the race detector](https://go.dev/doc/articles/race_detector)
- [Go: the sync package](https://pkg.go.dev/sync)
- [Go: channels in the tour](https://go.dev/tour/concurrency/2)
