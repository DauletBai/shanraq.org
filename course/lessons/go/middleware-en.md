# Handlers and middleware: one wrapper for every route

_Лид (summary):_ **The twenty-first lesson of the Go course. A wrapper around the router does the same thing for every request: a log line, a response time, a check. What http.Handler is, why a wrapper takes a handler and returns a handler, which order two of them run in, and how to learn the status code that ResponseWriter refuses to hand back.**

## Why this matters

A log line is wanted on every page. A sign-in check on half of them. A timing on all. Copying three lines into every handler means forgetting them in one of them some day, and spending half a day working out why "this page never shows up in the log".

Things like that are written once and put on everything at once.

## The whole thing at once

A new folder, `go mod init sabaq19`, `main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// logging wraps any handler and writes a line to the log.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s → %d in %s", r.Method, r.URL.Path, rec.code, time.Since(start).Round(time.Millisecond))
	})
}

// recorder remembers the status code: ResponseWriter does not keep it.
type recorder struct {
	http.ResponseWriter
	code int
}

func (rec *recorder) WriteHeader(code int) {
	rec.code = code
	rec.ResponseWriter.WriteHeader(code)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "the home page")
}

func slow(w http.ResponseWriter, r *http.Request) {
	time.Sleep(120 * time.Millisecond)
	fmt.Fprintln(w, "the slow page")
}

func gone(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "gone", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /slow", slow)
	mux.HandleFunc("GET /gone", gone)

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", logging(mux)))
}
```

Run it and visit all four addresses — `/`, `/slow`, `/gone` and the non-existent `/kokek`. Before looking below, say out loud whether the last one reaches the log. Then compare with the terminal:

```
GET / → 200 in 0s
GET /slow → 200 in 121ms
GET /gone → 404 in 0s
GET /kokek → 404 in 0s
```

That last line is the important one. `/kokek` never reached a handler of yours — the router answered 404 by itself — and the log has an entry all the same.

## Taking it apart

### `http.Handler` is an interface with one method

```go
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

Anything that can do `ServeHTTP` is a handler. And `ServeMux` is a handler too: it has that method, looks at the path and calls the right function.

The whole construction of this lesson follows from that. Since the router is a handler, it can be passed anywhere a handler is expected. Into another handler, for instance.

### `HandlerFunc` — the bridge from a function to the interface

Your `home` and `slow` are ordinary functions with no `ServeHTTP` method. `http.HandlerFunc` turns them into handlers: it is a type with a single method that simply calls the function itself.

```go
http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { ... })
```

Do not be put off by the spelling: it is a type conversion, like `int64(n)`. You declare a function as a function and dress it as the interface with that line.

### Middleware is a function that takes a handler and returns a handler

```go
func logging(next http.Handler) http.Handler
```

The whole point is in the signature. In goes what you had. Out comes the same thing, wrapped. Which is why wrappers can be put on one another, and neither the router nor the handlers know about them.

Inside a wrapper there are three parts: something before, the call to `next.ServeHTTP(w, r)`, something after.

> **Picture it.** The turnstile at a building's entrance. It does not know which office you are going to and does not decide for you. It notes the time you came in, lets you through, and notes the time you left — the same for everyone.

### The order: wrappers go on like nesting dolls

Put two on and see which order they run in:

```go
http.ListenAndServe(":8080", mark("outer")(mark("inner")(mux)))
```

```
outer in
inner in
the handler
inner out
outer out
```

The outer one meets the request first and sees the response out last. Hence a simple rule: whatever must happen before everything else goes on last, on the outside.

> **Picture it.** Nesting dolls. Take off the outer one and the next is underneath, down to the smallest; put them back in the reverse order. The request goes right through them and comes back the same way.

### The status code has to be remembered yourself

Here is where beginners trip. `ResponseWriter` has **no** way of being asked which code has gone out. It can only write.

So the wrapper slips in one of its own:

```go
type recorder struct {
	http.ResponseWriter
	code int
}

func (rec *recorder) WriteHeader(code int) {
	rec.code = code
	rec.ResponseWriter.WriteHeader(code)
}
```

The first line of the struct is embedding: `recorder` gets every method of `ResponseWriter` for nothing and remains a fully-fledged `ResponseWriter`. It has one method of its own — `WriteHeader`, which remembers the code and passes it on.

> **A caveat.** This `recorder` is a teaching one. A real one does more: it stops the status being written twice (Go silently ignores a second `WriteHeader`, which is worth knowing) and passes on the extra abilities of a `ResponseWriter` — `http.Flusher`, for one, without which streaming does not work. Ours is enough for a log; it is not enough for a library.

And `code: http.StatusOK` at creation is not decoration: a handler that simply writes a body never calls `WriteHeader` at all, and the response still goes out as a 200.

> **Picture it.** Carbon paper under a form. The original goes to the recipient, the copy stays with you — and nothing about what the recipient received has changed.

### A wrapper for every route, and a wrapper for one

```go
http.ListenAndServe(":8080", logging(mux))          // on everything
mux.Handle("GET /admin", onlyStaff(adminPage))      // on one route
```

The first is for the whole site: the log, the timing, recovering from a panic. The second is for what not everything needs: a sign-in check, a rate limit.

### What does not belong in a wrapper

A wrapper is tempting: it sees everything. So things that have no business being there move in quickly.

The rule is simple: **a wrapper does the same thing for everyone.** The moment an `if r.URL.Path == …` appears inside it, it is no longer a wrapper but a router written a second time and worse. That belongs either in a handler of its own or in a wrapper put on one route.

## The lesson map

![The lesson map: the request goes through the wrapper twice, in and out](/static/course/go/map-middleware-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why can `logging` wrap `mux` when `mux` is not a function?
2. Put two wrappers on. Which of them sees the request first?
3. Why does the wrapper need a `recorder` of its own when the status code has already gone out through `w`?

## Exercise

**Required.** Write a wrapper `withHeader` that adds an `X-Blog: go-oqu` header to every response, and put it on the whole router together with `logging`. Check with `curl -i` that the header is on all three addresses, the 404 included. Then write a second wrapper that prints a message before and after, and confirm the order is the one you predicted.

**Optional.**

- Have `recorder` count the bytes of the response as well and print that in the log.
- Write a `recover` wrapper that catches a handler's panic and answers 500 instead of dropping the connection.
- Put a wrapper on `/slow` alone, through `mux.Handle`, and check that it does not fire on the other addresses.

## Where this goes in your blog

You have just written the log, and it stays to the end of the course. Recovering from a panic and a sign-in check join it once there are users.

And in the next lesson HTML finally travels in the body: templates turn `[]Article` into a page, and `fmt.Fprintln` leaves the handlers.

## Answers

1. Because a wrapper takes not a function but the `http.Handler` interface — anything with a `ServeHTTP` method. `ServeMux` has one, so the router is a handler like any other.
2. The outer one: it gets the request first and hands the response over last. The order is the nesting dolls' — whatever is on the outside meets it first.
3. Because `ResponseWriter` can only write; the code it sent cannot be read back from it. The wrapper substitutes its own `recorder`, which notes the code on the way past and passes it to the real writer unchanged.

## Sources

- [http.Handler — documentation](https://pkg.go.dev/net/http#Handler)
- [http.HandlerFunc — documentation](https://pkg.go.dev/net/http#HandlerFunc)
- [Effective Go: embedding](https://go.dev/doc/effective_go#embedding)
