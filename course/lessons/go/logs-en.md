# Errors and logs in Go: seeing what broke

_Лид (summary):_ **The twenty-fifth lesson of the Go course. A panic in a handler does not bring the server down — it brings the answer down, and the reader gets nothing at all. recover in one wrapper, a reference number on the page and in the log, log/slog instead of printing to the console, and why an answer that has already left cannot be rescued.**

## Why this matters

While everything works, a log looks like a luxury. You need it on exactly the day a reader writes "your site does not open" and says nothing else.

Today there are two questions. What the reader sees when the program breaks. And what has to be left behind for you, so you can work out why.

## The whole thing at once

A new folder, `go mod init sabaq23`, `main.go`:

```go
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime/debug"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

// The number is read back out of the header we set a moment ago.
func reqID(w http.ResponseWriter) string { return w.Header().Get("X-Request-Id") }

func withID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", fmt.Sprintf("%08x", rand.Uint32()))
		next.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				w.Header().Set("Connection", "close")
				log.Error("panic", "id", reqID(w), "path", r.URL.Path,
					"cause", fmt.Sprint(v), "stack", string(debug.Stack()))
				serverError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// One sentence for every reader; the detail stays in the log.
func serverError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Something broke on our side. Reference number: %s\n", reqID(w))
}

var errNotFound = errors.New("article not found")

func find(slug string) (string, error) {
	titles := map[string]string{"dala": "About the steppe", "shanyraq": "What a shanyraq is"}
	if t, ok := titles[slug]; ok {
		return t, nil
	}
	return "", fmt.Errorf("find %q: %w", slug, errNotFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "My blog")
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		title, err := find(r.PathValue("slug"))
		switch {
		case errors.Is(err, errNotFound):
			log.Info("no such article", "id", reqID(w), "slug", r.PathValue("slug"))
			http.NotFound(w, r)
			return
		case err != nil:
			log.Error("could not read", "id", reqID(w), "err", err)
			serverError(w)
			return
		}
		fmt.Fprintln(w, title)
	})

	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		var titles map[string]string
		titles["dala"] = "About the steppe" // writing to a nil map
	})

	handler := withID(recoverPanic(mux))

	log.Info("server started", "addr", "http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
```

## Taking it apart

### A panic does not bring the server down, it brings the answer down

First, see what happens with no protection at all. Take `recoverPanic` out of the last line for a minute and call `/panic`:

```
curl -i http://localhost:8080/panic
```

```
curl: (52) Empty reply from server
```

Not a 500. Not a blank page. **Nothing**: the connection closed and there was no status code at all.

The server is alive, though — ask for `/` again and it answers. `net/http` catches the panic itself, at the connection level, and this is what lands in its own log:

```
2026/09/05 16:03:35 http: panic serving [::1]:55036: assignment to entry in nil map
goroutine 4 [running]:
net/http.(*conn).serve.func1()
```

So you have a stack trace and the reader has nothing. From where they stand the site is simply broken, and they cannot tell you anything useful: there is nothing to report.

> **Picture it.** The shop assistant faints without a word. The shop is still open, the next customer is served, and this one is left at the counter not even knowing who to ask.

### `recover` — one piece of middleware for every page

Put `recoverPanic` back and call the same address:

```
HTTP/1.1 500 Internal Server Error
Connection: close
X-Request-Id: 1a2b3c4d
Date: Sat, 05 Sep 2026 11:03:07 GMT
Content-Length: 56
Content-Type: text/plain; charset=utf-8

Something broke on our side. Reference number: 1a2b3c4d
```

Now there is a conversation: a status code, a page, and a number.

It works through `defer` and `recover`. `defer` runs a function when the handler finishes — normally or by panic, it makes no difference. `recover` inside such a function stops the fall and hands back whatever was panicked with.

```go
defer func() {
	if v := recover(); v != nil {
		...
	}
}()
```

`recover()` means something **only** inside a deferred function. Called on its own it returns `nil` and stops nothing.

We set `Connection: close` because a connection that carried half an answer does not deserve trust: let the browser open a new one.

And why middleware rather than a line in every handler: there will be thirty handlers, and forgetting it in one is enough. The wrapper goes on the whole `mux` once — the trick from the middleware lesson.

### A number for the reader, the detail for the log

The temptation to show the reader `err.Error()` is strong, and you must not. The text of an error is your insides: a table name, a file path, sometimes a fragment of a query. It helps a passer-by not at all, and it helps someone hunting for holes a great deal.

So the sentence is the same for everyone, and instead of detail there is a number:

```
Something broke on our side. Reference number: 1a2b3c4d
```

The same number is in the log:

```
level=ERROR msg=panic id=1a2b3c4d path=/panic cause="assignment to entry in nil map" stack="goroutine 4 [running..."
```

The time is gone from the excerpts and the request number is a sample: yours will differ, but it is the same in both lines.

The reader sends eight characters and you find the line. Without that number the search is "roughly yesterday evening somebody opened something".

We put the number in an `X-Request-Id` header and read it back from there. The grown-up way is to carry it in the request context, but context comes later; a header is something you already know, and it has a side benefit — the number is visible in the response.

### The limit of the trick: when the answer has already left

Now, honestly, about what `recover` cannot do.

A handler that wrote two hundred lines and only then fell over produces this:

```
HTTP/1.1 200 OK
...
line 199 — filling the buffer
Something broke on our side. Reference number: 1a2b3c4d
```

The status is `200`. Not `500`: the headers left when the `net/http` buffer filled up, and they cannot be changed now. A familiar line lands in the log:

```
superfluous response.WriteHeader call from main.serverError (main.go:41)
```

The reader sees a page that looks fine, with an apology at the end. A search engine sees `200` and indexes it.

Hence the rule already introduced in the templates lesson, which now pays for itself: **assemble the page in a buffer in full and send it in one piece**. Then the fall happens before the first byte and `recover` gets there in time.

### `log/slog`: lines you can search

`fmt.Println` in a log is text. `slog` is fields.

```go
log.Info("no such article", "id", reqID(w), "slug", r.PathValue("slug"))
```

```
level=INFO msg="no such article" id=1a2b3c4d slug=kokek
```

Name-value pairs follow the message. The difference shows up at a million lines: `slug=kokek` can be filtered, counted and grouped, while "could not find article kokek" can only be read by eye.

The handler is chosen when the logger is made:

```go
slog.New(slog.NewTextHandler(os.Stdout, nil))   // for a person
slog.New(slog.NewJSONHandler(os.Stdout, nil))   // for a machine
```

`Text` reads well over a developer's shoulder, `JSON` is collected by a log store. Swapping them is one line, because none of the rest of the code is tied to the choice.

The levels — `Debug`, `Info`, `Warn`, `Error` — exist so you do not have to read everything. `Debug` is not printed by default.

### An error is not always a `500`

Look at `/read/{slug}` again. An article that does not exist is not a breakage:

```go
case errors.Is(err, errNotFound):
	log.Info("no such article", ...)
	http.NotFound(w, r)
```

`Info` rather than `Error`, and `404` rather than `500`. The reader mistyped an address or the article was taken down — the program did its job correctly.

Log `Error` for every page that was not found and within a week the log holds ten thousand "errors", with the real one lost among them. A level is a promise: `Error` means a person needs to look.

The distinction comes from the errors lesson: `errors.Is` answers "which error is this", and the answer picks both the status code and the log level.

### Logs go to stdout, not to a file

`slog.NewTextHandler(os.Stdout, nil)` writes to the standard output, and that is not a simplification made for the lesson.

A program that writes to a file itself has to think about the path, the permissions, the size, the rotation and the deleting of old ones. Whatever launches it already does all of that: `systemd`, `docker`, any hosting. Write to `stdout` and you see the lines in your terminal at home, while on the server they go wherever it is configured to put them, with not a single change in the code.

## The lesson map

![Lesson map: the panic is caught, and the reader gets a code and a number](/static/course/go/map-errors2-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. A handler panics and there is no `recover` middleware. What does the reader see?
2. Why must the text of an error not be shown to the reader?
3. Why is an article that was not found logged at `Info` rather than `Error`?

## Exercise

**Required.** Bring the templates from the templates lesson back into the program and make the `500` page a template of its own, with the reference number on it. Assemble it in a buffer. Check that a panic really does produce a `500`, and that the number on the page matches the number in the log.

**Optional.**

- Replace `NewTextHandler` with `NewJSONHandler` and look at the same line again.
- Add middleware that writes one line per request: method, path, status and duration. You will have to remember the status yourself, as in the middleware lesson. Put it outside `recoverPanic`: otherwise a panic unwinds straight past it and the request line never appears.
- Write a handler that produces many lines and then panics, and confirm that there is no `500` in the response.

## Where this goes in your blog

The blog has stopped going silent when it breaks. Next comes configuration: the address, the port and the paths stop being written into the code.

Debts. The reference number lives in a header, and its proper home is the request context; that arrives with the concurrency lesson. The per-request line is still separate middleware, and one wrapper for everything would be better. And `debug.Stack()` in a log is a heavy thing: on a real site it is written for real panics only, not for every error.

## Answers

1. Nothing. The connection closes with no answer: `curl` says `Empty reply from server` and the browser shows its own page about the site being unreachable. There is no status code at all. The server keeps running and the stack goes into its own log — so you have the information and the reader does not.
2. Because the text of an error is the program's insides: table names, paths, fragments of queries. It helps the reader with nothing, and it tells anyone looking for weaknesses how the server is built. So the reader gets one general sentence and a reference number, and the detail goes to the log.
3. Because it is not a breakage: the program did the right thing, there simply is no such article. `Error` is a promise that a person has to step in. Mark every wrong address that way and the real errors drown among thousands of false ones.

## Sources

- [log/slog — documentation](https://pkg.go.dev/log/slog)
- [The Go blog: structured logging with slog](https://go.dev/blog/slog)
- [net/http: Server.ErrorLog](https://pkg.go.dev/net/http#Server)
- [runtime/debug: Stack](https://pkg.go.dev/runtime/debug#Stack)
- [Effective Go: defer, panic and recover](https://go.dev/doc/effective_go#recover)
