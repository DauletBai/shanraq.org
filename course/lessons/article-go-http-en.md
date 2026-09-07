# HTTP for real: method, path, headers and status code

_Лид (summary):_ **The nineteenth lesson of the Go course and the start of the module on the web. A request is a struct you can read: method, path, headers, parameters. A response goes out one way — headers, status, body, in that order. What happens to a header set too late, why 404 and 405 are different answers, and what the superfluous WriteHeader message means.**

## Why this matters

In the first lesson you looked at HTTP from the outside with `curl`. In the third you wrote a server that answers with one line. Everything else lies between those two points: how to **read** a request and how to **assemble** a response deliberately rather than with one line for every case.

This is where the real blog begins. Until the server can tell what it was asked, it can have no pages, no languages and no errors.

## The whole thing at once

A new folder, `go mod init sabaq17`, `main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "GET only")
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no such page:", r.URL.Path)
		return
	}

	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "kz"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "method:", r.Method)
	fmt.Fprintln(w, "path:", r.URL.Path)
	fmt.Fprintln(w, "language:", lang)
	fmt.Fprintln(w, "browser:", r.Header.Get("User-Agent"))
}

func main() {
	http.HandleFunc("/", handle)
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Run `go run .` and open `localhost:8080` in a browser. Then look at the whole response — the `-i` flag shows the headers:

```
curl -i "http://localhost:8080/?lang=ru"
```

Before looking below, say out loud which status code will come back. Then compare:

```
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Sat, 05 Sep 2026 04:38:30 GMT
Content-Length: 51

method: GET
path: /
language: ru
browser: curl/8.4
```

Now ask for a page that does not exist, and with the wrong method:

```
curl -i http://localhost:8080/kokek
curl -i -X POST http://localhost:8080/
```

```
HTTP/1.1 404 Not Found
...
no such page: /kokek
```

```
HTTP/1.1 405 Method Not Allowed
Allow: GET
```

## Taking it apart

### A request is a struct, and you have seen it before

`r *http.Request` is that star from the pointers lesson. You are handed the address of the request rather than a copy of it: a request is large, and copying it on every call would be wasteful.

Inside are fields you reach through a dot, as in any struct:

```go
r.Method        // "GET"
r.URL.Path      // "/kokek"
r.Header        // all the headers
```

There is nothing new in the language here. Only the field names are new.

> **Picture it.** A filled-in form. At the top is what is being asked for (the method) and about what (the path), below are boxes with details (the headers), and at the end an attachment if there is one (the body). You read the form rather than guess.

### Method and path

The method says **what** to do: `GET` means show me, `POST` means take this. The path says **what with**.

You have to check them yourself. A handler registered at `/` receives every path nobody else claimed — `/kokek` included. That is why the program has an explicit `r.URL.Path != "/"` check; without it a page that does not exist would answer 200 with the front page.

Take method names from the constants — `http.MethodGet` rather than `"GET"`: the compiler will not catch a typo in a string, but it will catch one in a constant's name.

### Headers: reading and setting

```go
r.Header.Get("User-Agent")                        // read
w.Header().Set("Content-Type", "text/plain")      // set
```

Reading is from the request, writing is into the response, and there is little room to confuse them: on a request `Header` is a field, on a response `Header()` is a method.

If you do not set `Content-Type`, Go guesses it from the first bytes of the body. In our 404 we set nothing, and `Content-Type: text/plain; charset=utf-8` appeared in the response anyway. Convenient, but not something to rely on: it recognises HTML, and it will not recognise a format of your own.

### The status code, and the order

Here is the main rule of the lesson, and it is worth checking with your own hands.

```go
w.Header().Set("Content-Type", "application/json")   // 1. headers
w.WriteHeader(http.StatusNotFound)                   // 2. status
fmt.Fprintln(w, "{}")                                // 3. body
```

In that order and no other. The response leaves for the network as you write it, and what has gone cannot be called back.

Swap the first two lines and the header will not travel: the response will carry the type Go guessed rather than yours. Write the body before `WriteHeader` and the code stays 200 for good, while the server's log gets:

```
http: superfluous response.WriteHeader call from main.main.func3 (main.go:22)
```

> **Picture it.** An envelope in a postbox. While it is in your hand, write anything you like on it. Once it is through the slot it is too late to write, and there is no getting it back.

A warning of its own: if you decide to check this with a test using `httptest.NewRecorder`, the late header **will** show up there. The recorder keeps the headers in an ordinary map and hands it over at the end, which a real connection does not. This has to be checked against a live server.

### The body: `w` is an interface

`fmt.Fprintln(w, …)` works because `w` can accept bytes. It is an `io.Writer` — the interface from the lesson on contracts — and the same function writes to a file, to a buffer and to the network without knowing the difference.

Hence the trick we used all through the language module: write the logic so that it writes into an `io.Writer`, and it can be checked by a test with nothing running.

### Query parameters

```go
lang := r.URL.Query().Get("lang")
```

Everything after the `?` is a parameter. `Query()` parses them into a map, `Get` fetches one by name and returns an empty string when it is absent. No errors, no panics — which is why the default value is yours to set.

That is exactly how the language of an article is chosen on this site: `?lang=kz`.

### 404 and 405 are different answers

People confuse them, and the difference is simple. **404** means there is no such page at all. **405** means the page exists but not like that: you sent a `POST` where only `GET` is accepted.

And 405 comes with an obligation: an `Allow` header listing what is permitted. You saw it in the response above — it is not decoration but part of the protocol.

> **Picture it.** A locked door and the wrong door. "There is no such office" and "the office exists but opens on Wednesdays" are different answers, and a person acts differently on each.

## The lesson map

![The lesson map: a response goes out one way — headers, status, body](/static/course/go/map-http-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why must a handler registered at `/` check the path itself?
2. What happens if you set a header after `w.WriteHeader`?
3. When do you answer 404, and when 405?

## Exercise

**Required.** Extend the server so that it answers three paths: `/` with a greeting, `/about` with a line about yourself, and everything else with a 404 and a clear message. Allow only `GET` on `/about` and answer 405 with an `Allow` header to anything else. Check every case with `curl -i` and make sure the codes are the ones you intended.

**Optional.**

- Add `?lang=`: answer in the matching language for `kz`, `ru` and `en`, and with a 400 and an explanation for anything else.
- Set a header of your own, `X-Powered-By`, and find it in the `curl -i` output.
- Write the body before `WriteHeader` and look at the line in the server's log with your own eyes.

## Where this goes in your blog

This lesson is the foundation of the module. In the next one the paths separate into handlers and checking `r.URL.Path` by hand goes away. After that come templates, and HTML travels in the body instead of strings. But the order — headers, status, body — stays the same to the end of the course.

## Answers

1. Because `/` in the standard router means "everything nobody else claimed". Without the check, any non-existent page would return 200 with the content of the front page, and a search engine would conclude the site has thousands of identical pages.
2. Nothing happens — and that is the trouble. The headers have already gone out with the status code, and a late `Set` only changes a map in memory. The response will carry whatever Go guessed.
3. 404 when there is no such address at all. 405 when the address exists but the method is not allowed for it; a 405 must be accompanied by an `Allow` header listing the methods that are.

## Sources

- [net/http — documentation](https://pkg.go.dev/net/http)
- [http.Request — the request fields](https://pkg.go.dev/net/http#Request)
- [HTTP status codes — MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [RFC 9110: HTTP semantics](https://www.rfc-editor.org/rfc/rfc9110)
