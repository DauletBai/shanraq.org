# `context`: stopping work in time

_Лид (summary):_ **The forty-third lesson of the Go course. The reader closed the tab, and the handler learns of it by itself: measured here, it drops the work after two hundred milliseconds. A deadline for a query, a key of your own type instead of a string, and the promised move of the signed-in person into the request's context.**

## Why this matters

The handlers of the past lessons assume the reader is waiting. The reader is not always.

They close the tab, lose the signal in a lift, press reload halfway through. The server has no idea: it keeps going to the database, computing and assembling an answer nobody will take. With ten such readers it goes unnoticed; with ten thousand it is half the server's work wasted.

There is the opposite trouble too: a query to the database that hangs. With no deadline it holds one connection from the pool, then a second, then all of them — and the blog stops answering everybody because of one slow query.

Both are solved by the same thing — **a context**.

## The whole thing at once

A new folder, `go mod init sabaq40`:

```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("== the reader closed the tab")
	cancelled()

	fmt.Println()
	fmt.Println("== a query given half a second")
	fmt.Println("fast:  ", count(db, 100*1000, 500*time.Millisecond))
	fmt.Println("slow:  ", count(db, 100*1000*1000, 500*time.Millisecond))
}

// cancelled starts a slow handler and cuts the request off halfway.
func cancelled() {
	done := make(chan string, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		select {
		case <-time.After(2 * time.Second):
			done <- "finished the work"
		case <-r.Context().Done():
			// The request's context closes by itself once the reader is gone.
			done <- fmt.Sprintf("dropped the work after %s: %v",
				time.Since(start).Round(10*time.Millisecond), r.Context().Err())
		}
	}))
	defer srv.Close()

	// A client that gives up after two hundred milliseconds.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if _, err := http.DefaultClient.Do(req); err != nil {
		fmt.Println("the reader: left —", context.Cause(ctx))
	}
	fmt.Println("the handler:", <-done)
}

// count counts in the database slowly on purpose: a recursive query of n steps.
func count(db *sql.DB, n int, limit time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()

	start := time.Now()
	var got int
	err := db.QueryRowContext(ctx, `with recursive nums(i) as (
		select 1 union all select i + 1 from nums where i < ?
	) select count(*) from nums`, n).Scan(&got)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Sprintf("cut off after %s: %v", time.Since(start).Round(10*time.Millisecond), err)
	case err != nil:
		return "error: " + err.Error()
	}
	return fmt.Sprintf("counted %d in %s", got, time.Since(start).Round(10*time.Millisecond))
}
```

The output:

```
== the reader closed the tab
the reader: left — context deadline exceeded
the handler: dropped the work after 200ms: context canceled

== a query given half a second
fast:   counted 100000 in 60ms
slow:   cut off after 500ms: context deadline exceeded
```

## Taking it apart

### A request's context closes by itself

Every request has a context of its own: `r.Context()`. Nobody makes it by hand — `http.Server` does, and it **closes by itself** once the connection to the reader is gone.

All a handler has to do is listen:

```go
select {
case <-time.After(2 * time.Second):
	// the work finished
case <-r.Context().Done():
	// the reader left
}
```

`Done()` is a channel that closes at the moment of cancellation. The output shows it: the reader gave up after two hundred milliseconds, and after exactly two hundred milliseconds the handler dropped two seconds of work.

`Err()` says why: `context canceled` — cancelled; `context deadline exceeded` — out of time.

> **Picture it.** A waiter who has been told the guest has left. They will not carry the hot dish to an empty table and wait for it to go cold — they turn back halfway.

### A deadline for the work: `WithTimeout`

The second half of the output is about something else: the reader did not leave, the work dragged on.

```go
ctx, cancel := context.WithTimeout(context.Background(), limit)
defer cancel()
```

`WithTimeout` makes a new context that cancels itself after the given time. The fast query made it in sixty milliseconds; the slow one was cut off at exactly five hundred.

`defer cancel()` is compulsory, and not for tidiness: `WithTimeout` starts a timer and `cancel` releases it. Forget it and the timer lives to its deadline; on a busy server such forgotten timers pile up. Go complains about it when the code is checked, and rightly.

Notice where the parent came from: `context.Background()` is the "empty" context, the root of a chain. In a handler the parent must be **`r.Context()`** rather than `Background`: then a cancellation by the reader cuts the query to the database off as well.

### Why `database/sql` has the `...Context` twins

Now it is clear what those twin methods are for: `QueryContext`, `QueryRowContext`, `ExecContext`.

An ordinary `QueryRow` knows nothing about cancellation: it will wait for the database as long as the database likes. The version with a context hands the driver a deadline, and when it passes the query is **cut off on the database's side** — the connection goes back to the pool instead of hanging.

The rule is simple: in a handler always `...Context`, and always with `r.Context()`. The plain versions stay for scripts and migrations, where there is nobody to cancel.

### A context is the first argument

`func (s *Store) Get(ctx context.Context, slug string)` — that is how any function that can wait looks.

The context goes **first** and is called `ctx`. That is not a matter of style but a convention of the whole ecosystem: the standard library is written that way, every driver is written that way, and when you do the same your code reads without explanation.

What is not done: a context is not kept as a field in a struct. A context lives exactly as long as one call; a struct lives longer, and a context stored inside it will one day be either already cancelled or, worse, eternal.

### A value in a context and a key of your own type

A context can carry values — and here it is easy to do badly:

```go
package main

import (
	"context"
	"fmt"
)

// A type of your own for the key is not a whim. A value in a context is found
// by a key compared together with its type: userKey(0) and another type with
// the same zero are different keys, while two "user" strings from different
// packages are the same key.
type userKey struct{}

type otherKey struct{}

func main() {
	ctx := context.Background()

	fmt.Println("== string keys")
	ctx = context.WithValue(ctx, "user", "aigul")     //lint:ignore SA1029 on purpose
	ctx = context.WithValue(ctx, "user", "a library") //lint:ignore SA1029 on purpose
	fmt.Println("what is left:", ctx.Value("user"))

	fmt.Println("== types of our own")
	ctx = context.WithValue(ctx, userKey{}, "aigul")
	ctx = context.WithValue(ctx, otherKey{}, "a library")
	fmt.Println("ours:  ", ctx.Value(userKey{}))
	fmt.Println("theirs:", ctx.Value(otherKey{}))
}
```

```
== string keys
what is left: a library
== types of our own
ours:   aigul
theirs: a library
```

The first pair of lines: two `WithValue` calls with the string key `"user"` — and our value is gone, somebody else's is left. That is what happens in life: your `"user"` and the `"user"` of a library are **the same key**, and whoever writes last wins.

The second pair: keys of our own types. A key is compared together with its type, and a type declared in your package cannot be repeated from outside. Both values are there.

So a key is always declared with a type of your own, and access to the value is hidden behind two small functions — `withUser` and `userFrom` — so that nothing outside needs to know the key, or that a context is involved at all.

### The signed-in person moves into the context

In the lessons on sessions and permissions I left a debt: every handler looked the signed-in person up again. This is where that belongs.

```go
func withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, ok := currentUser(r); ok {
			r = r.WithContext(context.WithValue(r.Context(), userKey{}, u))
		}
		next.ServeHTTP(w, r)
	})
}
```

The wrapper looks the person up **once** and puts what it found into the request's context. Any handler then takes it with `userFrom(r.Context())` — no trip to the database, no repeated code.

`r.WithContext` returns a **copy** of the request with a new context: a request is not changed in place. Hence the `r = r.WithContext(...)` — forget to assign it back and the wrapper does nothing, exactly like `append` in the slices lesson.

### What does not go into a context

A context is for what belongs to the **request as a whole**: cancellation, a deadline, the signed-in person, the request id for the log, the language of the page.

Not for the parameters of a function. If a function needs a `slug`, it is passed as an argument rather than hidden in a context: an argument is visible in the signature and a value from a context is not, and its absence shows up only on a running server.

And not for what the program needs rather than the request: the store, the settings, the logger. Those are passed explicitly — the way we did with the `app` struct.

## The lesson map

![Lesson map: a cancellation travelling down the chain](/static/course/go/map-context-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Who closes `r.Context()`, and when?
2. Why is a string not used as a key in a context?
3. Why is `defer cancel()` compulsory even where the work finished in time?

## Exercise

**Required.** Move the blog onto contexts: the store's methods take `ctx` first and call the `...Context` versions, the handlers pass `r.Context()`, and the signed-in person is found by a wrapper and put into the context under a key of its own type.

All of it is done in [step-28](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-28) — compare against it once you have written your own.

**Optional.**

- Give every query to the database a three-second deadline through `http.TimeoutHandler` or a wrapper of your own, and check that a slow query is cut off.
- Add the reason to the log: `context.Cause(ctx)` tells you more than `Err()`.
- Take `defer cancel()` out and see what `go vet` says.

## Where this goes in your blog

The blog stops working into the void: the reader left, so we stopped too. And the signed-in person is no longer looked up in every handler.

Still open. Nothing but the example has a deadline on its queries. The background work from the previous lesson cannot stop on a context. And there is no graceful shutdown, in which unfinished requests live out their seconds — that is the deploy lesson.

## Answers

1. `http.Server` closes it when the connection to the reader is gone — the tab was closed, the signal was lost, the request was cancelled. Nobody creates or closes it by hand.
2. Because a key is compared by value and by type, and the string `"user"` from your package and the same string from a library are the same key. The values overwrite each other in silence. A type of your own cannot be repeated from outside.
3. Because `WithTimeout` starts a timer and `cancel` releases it. Without that the timer lives to the end of its deadline, and on a busy server such timers pile up.

## Sources

- [Go: the context package](https://pkg.go.dev/context)
- [The Go blog: context](https://go.dev/blog/context)
- [Go: database/sql and contexts](https://pkg.go.dev/database/sql#DB.QueryContext)
